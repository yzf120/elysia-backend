package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/yzf120/elysia-backend/dao"
	"github.com/yzf120/elysia-backend/errs"
	codeModel "github.com/yzf120/elysia-backend/model/code"
	"github.com/yzf120/elysia-backend/model/problem"
)

// testCase 测试用例结构（与 problem.test_cases JSON 对应）
type testCase struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	IsSample       int    `json:"is_sample"`
	Explanation    string `json:"explanation"`
}

// langConfig 语言执行配置
type langConfig struct {
	FileName   string   // 源文件名
	CompileCmd []string // 编译命令（nil 表示解释型语言）
	RunCmd     []string // 运行命令（%s 占位符替换为文件路径）
}

// javaHome Java 安装目录（Homebrew OpenJDK 17）
const javaHome = "/opt/homebrew/opt/openjdk@17"

// sandboxInclude 沙箱头文件目录（包含 bits/stdc++.h 等 macOS 缺失的头文件）
const sandboxInclude = "/Users/sylvainyang/project/elysia/elysia-backend/sandbox/include"

var langConfigs = map[string]langConfig{
	"python": {
		FileName:   "main.py",
		CompileCmd: nil,
		RunCmd:     []string{"python3", "main.py"},
	},
	"java": {
		FileName:   "Main.java",
		CompileCmd: []string{javaHome + "/bin/javac", "Main.java"},
		RunCmd:     []string{javaHome + "/bin/java", "-cp", ".", "Main"},
	},
	"go": {
		FileName:   "main.go",
		CompileCmd: []string{"go", "build", "-o", "main_bin", "main.go"},
		RunCmd:     []string{"./main_bin"},
	},
	"cpp": {
		FileName:   "main.cpp",
		CompileCmd: []string{"g++", "-O2", "-std=c++17", "-I", sandboxInclude, "-o", "main_bin", "main.cpp"},
		RunCmd:     []string{"./main_bin"},
	},
	"c": {
		FileName:   "main.c",
		CompileCmd: []string{"gcc", "-O2", "-I", sandboxInclude, "-o", "main_bin", "main.c"},
		RunCmd:     []string{"./main_bin"},
	},
}

// CodeRunService 代码运行服务
type CodeRunService struct {
	codeRunDAO dao.CodeRunDAO
	problemDAO dao.ProblemDAO
}

// NewCodeRunService 创建代码运行服务
func NewCodeRunService() *CodeRunService {
	return &CodeRunService{
		codeRunDAO: dao.NewCodeRunDAO(),
		problemDAO: dao.NewProblemDAO(),
	}
}

// SubmitCodeRun 提交代码运行任务（异步执行）
// testInput：测试模式下直接传入的样例输入（已废弃，改为从 showcase 字段读取）
func (s *CodeRunService) SubmitCodeRun(ctx context.Context, studentId string, problemId int64, language, code, runType, testInput string) (*codeModel.CodeRun, error) {
	// 校验语言
	if _, ok := langConfigs[language]; !ok {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "不支持的编程语言: "+language)
	}
	// 校验 runType
	if runType != "test" && runType != "submit" {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "run_type 必须为 test 或 submit")
	}

	// 查询题目（test 和 submit 都需要）
	p, err := s.problemDAO.GetProblemById(problemId)
	if err != nil || p == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "题目不存在")
	}

	// 创建运行记录（pending 状态）
	record := &codeModel.CodeRun{
		ProblemId: problemId,
		StudentId: studentId,
		Language:  language,
		Code:      code,
		RunType:   runType,
		Status:    "pending",
	}
	if err := s.codeRunDAO.CreateCodeRun(record); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建运行记录失败: "+err.Error())
	}

	if runType == "test" {
		// 测试模式：使用 showcase 字段的用例（不记录到运行记录，直接执行后更新结果）
		go s.executeCodeWithShowcase(record.Id, p, language, code)
	} else {
		// 提交模式：执行 test_cases 全部用例
		go s.executeCode(record.Id, p, language, code, runType)
	}

	return record, nil
}

// GetCodeRunResult 查询代码运行结果
func (s *CodeRunService) GetCodeRunResult(runId int64) (*codeModel.CodeRun, error) {
	record, err := s.codeRunDAO.GetCodeRunById(runId)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "运行记录不存在")
	}
	return record, nil
}

// ListCodeRunRecords 查询学生某题的运行记录列表（倒序，只查 submit 类型）
func (s *CodeRunService) ListCodeRunRecords(studentId string, problemId int64, limit int) ([]*codeModel.CodeRun, error) {
	records, err := s.codeRunDAO.ListCodeRunsByStudent(studentId, problemId, limit)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询运行记录失败: "+err.Error())
	}
	return records, nil
}

// BatchGetAcceptedProblems 批量查询学生已完全通过的题目ID集合
func (s *CodeRunService) BatchGetAcceptedProblems(studentId string, problemIds []int64) (map[int64]bool, error) {
	result, err := s.codeRunDAO.BatchGetAcceptedProblems(studentId, problemIds)
	if err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "查询完成状态失败: "+err.Error())
	}
	return result, nil
}

// showcaseCaseResult 单个 showcase 用例的执行结果
type showcaseCaseResult struct {
	Index          int    `json:"index"`           // 用例序号（从1开始）
	Input          string `json:"input"`           // 输入
	ExpectedOutput string `json:"expected_output"` // 预期输出
	ActualOutput   string `json:"actual_output"`   // 实际输出
	Passed         bool   `json:"passed"`          // 是否通过
	Status         string `json:"status"`          // accepted / wrong_answer / runtime_error / time_limit_exceeded 等
	ErrorMsg       string `json:"error_msg"`       // 错误信息（编译/运行错误时）
	TimeCost       int64  `json:"time_cost"`       // 执行耗时 ms
}

// executeCodeWithShowcase 测试模式：使用题目 showcase 字段的用例运行代码
// output 字段存储 JSON 格式的每个 case 详细结果，供前端可视化展示
func (s *CodeRunService) executeCodeWithShowcase(runId int64, p *problem.Problem, language, code string) {
	_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{"status": "running"})

	// 解析 showcase 用例
	var showcaseCases []testCase
	if err := json.Unmarshal([]byte(p.Showcase), &showcaseCases); err != nil || len(showcaseCases) == 0 {
		_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{
			"status":    "runtime_error",
			"error_msg": "题目 showcase 格式错误或为空",
		})
		return
	}

	cfg := langConfigs[language]
	tmpDir, err := os.MkdirTemp("", "elysia_code_*")
	if err != nil {
		_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{
			"status":    "runtime_error",
			"error_msg": "创建临时目录失败",
		})
		return
	}
	defer os.RemoveAll(tmpDir)

	srcFile := filepath.Join(tmpDir, cfg.FileName)
	if err := os.WriteFile(srcFile, []byte(code), 0644); err != nil {
		_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{
			"status":    "runtime_error",
			"error_msg": "写入代码文件失败",
		})
		return
	}

	if cfg.CompileCmd != nil {
		compileCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		compileCmd := exec.CommandContext(compileCtx, cfg.CompileCmd[0], cfg.CompileCmd[1:]...)
		compileCmd.Dir = tmpDir
		compileCmd.Env = buildEnv()
		var compileErr bytes.Buffer
		compileCmd.Stderr = &compileErr
		if err := compileCmd.Run(); err != nil {
			errStr := compileErr.String()
			_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{
				"status":    "compile_error",
				"error_msg": errStr,
			})
			return
		}
	}

	const defaultTimeLimitMs int64 = 5000
	var totalTimeCost int64
	var maxMemoryUsed int64
	var caseResults []showcaseCaseResult
	finalStatus := "accepted"

	for i, tc := range showcaseCases {
		status, output, errMsg, timeCost, memUsed := s.runSingleCase(tmpDir, cfg, tc.Input, defaultTimeLimitMs)
		totalTimeCost += timeCost
		if memUsed > maxMemoryUsed {
			maxMemoryUsed = memUsed
		}

		actualOutput := strings.TrimSpace(output)
		expectedOutput := strings.TrimSpace(tc.ExpectedOutput)
		passed := status == "accepted" && actualOutput == expectedOutput

		caseResult := showcaseCaseResult{
			Index:          i + 1,
			Input:          tc.Input,
			ExpectedOutput: expectedOutput,
			ActualOutput:   actualOutput,
			Passed:         passed,
			Status:         status,
			ErrorMsg:       errMsg,
			TimeCost:       timeCost,
		}
		// 若运行本身出错（非 wrong_answer），status 直接用后端返回的
		if status != "accepted" {
			caseResult.Status = status
			finalStatus = status
		} else if !passed {
			caseResult.Status = "wrong_answer"
			finalStatus = "wrong_answer"
		}
		caseResults = append(caseResults, caseResult)
	}

	// 将所有 case 结果序列化为 JSON 存入 output 字段
	outputJSON, _ := json.Marshal(caseResults)

	_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{
		"status":      finalStatus,
		"output":      string(outputJSON),
		"error_msg":   "",
		"time_cost":   totalTimeCost,
		"memory_used": maxMemoryUsed,
	})
}

// executeCode 在沙箱中执行代码（goroutine 中运行）
func (s *CodeRunService) executeCode(runId int64, p *problem.Problem, language, code, runType string) {
	// 更新状态为 running
	_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{"status": "running"})

	// 解析测试用例
	var testCases []testCase
	if err := json.Unmarshal([]byte(p.TestCases), &testCases); err != nil || len(testCases) == 0 {
		_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{
			"status":    "runtime_error",
			"error_msg": "题目测试用例格式错误",
		})
		return
	}

	// 根据 runType 决定使用哪些测试用例
	// submit 模式：跑所有测试用例
	casesToRun := testCases

	cfg := langConfigs[language]

	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "elysia_code_*")
	if err != nil {
		_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{
			"status":    "runtime_error",
			"error_msg": "创建临时目录失败",
		})
		return
	}
	defer os.RemoveAll(tmpDir)

	// 写入源代码文件
	srcFile := filepath.Join(tmpDir, cfg.FileName)
	if err := os.WriteFile(srcFile, []byte(code), 0644); err != nil {
		_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{
			"status":    "runtime_error",
			"error_msg": "写入代码文件失败",
		})
		return
	}

	// 编译（如果需要）
	if cfg.CompileCmd != nil {
		compileCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		compileCmd := exec.CommandContext(compileCtx, cfg.CompileCmd[0], cfg.CompileCmd[1:]...)
		compileCmd.Dir = tmpDir
		compileCmd.Env = buildEnv()
		var compileErr bytes.Buffer
		compileCmd.Stderr = &compileErr
		if err := compileCmd.Run(); err != nil {
			_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{
				"status":    "compile_error",
				"error_msg": compileErr.String(),
			})
			return
		}
	}

	// 逐个运行测试用例
	var totalTimeCost int64
	var maxMemoryUsed int64
	finalStatus := "accepted"
	var caseResults []showcaseCaseResult

	// 默认时间限制 1000ms
	const defaultTimeLimitMs int64 = 1000

	for i, tc := range casesToRun {
		status, output, errMsg, timeCost, memUsed := s.runSingleCase(tmpDir, cfg, tc.Input, defaultTimeLimitMs)
		totalTimeCost += timeCost
		if memUsed > maxMemoryUsed {
			maxMemoryUsed = memUsed
		}

		actualOutput := strings.TrimSpace(output)
		expectedOutput := strings.TrimSpace(tc.ExpectedOutput)

		caseStatus := status
		passed := false
		if status == "accepted" {
			if actualOutput == expectedOutput {
				passed = true
			} else {
				caseStatus = "wrong_answer"
			}
		}

		caseResults = append(caseResults, showcaseCaseResult{
			Index:          i + 1,
			Input:          tc.Input,
			ExpectedOutput: expectedOutput,
			ActualOutput:   actualOutput,
			Passed:         passed,
			Status:         caseStatus,
			ErrorMsg:       errMsg,
			TimeCost:       timeCost,
		})

		if !passed && finalStatus == "accepted" {
			finalStatus = caseStatus
		}
	}

	// 将所有 case 结果序列化为 JSON 存入 output 字段
	outputJSON, _ := json.Marshal(caseResults)

	_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{
		"status":      finalStatus,
		"output":      string(outputJSON),
		"error_msg":   "",
		"time_cost":   totalTimeCost,
		"memory_used": maxMemoryUsed,
	})
}

// buildEnv 构建子进程环境变量，确保 Java 等工具路径可被找到
func buildEnv() []string {
	path := os.Getenv("PATH")
	// 将 Homebrew Java 路径注入，避免 Go 服务进程 PATH 缺失
	javaBin := javaHome + "/bin"
	if !strings.Contains(path, javaBin) {
		path = javaBin + ":" + path
	}
	env := os.Environ()
	for i, e := range env {
		if strings.HasPrefix(e, "PATH=") {
			env[i] = "PATH=" + path
			return env
		}
	}
	return append(env, "PATH="+path)
}

// runSingleCase 运行单个测试用例，返回 (status, output, errMsg, timeCostMs, memUsedKB)
func (s *CodeRunService) runSingleCase(tmpDir string, cfg langConfig, input string, timeLimitMs int64) (string, string, string, int64, int64) {
	runCtx, cancel := context.WithTimeout(context.Background(), time.Duration(timeLimitMs+2000)*time.Millisecond)
	defer cancel()

	runArgs := make([]string, len(cfg.RunCmd))
	copy(runArgs, cfg.RunCmd)
	cmd := exec.CommandContext(runCtx, runArgs[0], runArgs[1:]...)
	cmd.Dir = tmpDir
	cmd.Env = buildEnv()
	cmd.Stdin = strings.NewReader(input)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	startTime := time.Now()
	err := cmd.Run()
	timeCost := time.Since(startTime).Milliseconds()

	// 获取内存使用（近似值，通过 /proc/self/status 或 runtime）
	var memUsed int64
	if cmd.ProcessState != nil {
		memUsed = getMemoryUsage(cmd)
	}

	if runCtx.Err() == context.DeadlineExceeded || timeCost >= timeLimitMs {
		return "time_limit_exceeded", "", fmt.Sprintf("执行时间超过限制 %dms", timeLimitMs), timeCost, memUsed
	}

	if err != nil {
		errOutput := stderr.String()
		if errOutput == "" {
			errOutput = err.Error()
		}
		return "runtime_error", stdout.String(), errOutput, timeCost, memUsed
	}

	output := stdout.String()
	return "accepted", output, "", timeCost, memUsed
}

// getMemoryUsage 获取进程内存使用量（KB）
func getMemoryUsage(cmd *exec.Cmd) int64 {
	if runtime.GOOS == "linux" && cmd.ProcessState != nil {
		// Linux 下通过 /proc/{pid}/status 获取
		// 由于进程已结束，使用 runtime 统计近似值
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		return int64(ms.Alloc / 1024)
	}
	// 其他系统返回近似值
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	return int64(ms.Alloc / 1024)
}
