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

var langConfigs = map[string]langConfig{
	"python": {
		FileName:   "main.py",
		CompileCmd: nil,
		RunCmd:     []string{"python3", "main.py"},
	},
	"java": {
		FileName:   "Main.java",
		CompileCmd: []string{"javac", "Main.java"},
		RunCmd:     []string{"java", "-cp", ".", "Main"},
	},
	"go": {
		FileName:   "main.go",
		CompileCmd: []string{"go", "build", "-o", "main_bin", "main.go"},
		RunCmd:     []string{"./main_bin"},
	},
	"cpp": {
		FileName:   "main.cpp",
		CompileCmd: []string{"g++", "-O2", "-o", "main_bin", "main.cpp"},
		RunCmd:     []string{"./main_bin"},
	},
	"c": {
		FileName:   "main.c",
		CompileCmd: []string{"gcc", "-O2", "-o", "main_bin", "main.c"},
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
// testInput：测试模式下直接传入的样例输入，非空时跳过查题目直接运行
func (s *CodeRunService) SubmitCodeRun(ctx context.Context, studentId string, problemId int64, language, code, runType, testInput string) (*codeModel.CodeRun, error) {
	// 校验语言
	if _, ok := langConfigs[language]; !ok {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "不支持的编程语言: "+language)
	}
	// 校验 runType
	if runType != "test" && runType != "submit" {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "run_type 必须为 test 或 submit")
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

	if runType == "test" && testInput != "" {
		// 测试模式：直接用前端传来的样例输入运行，不查数据库题目
		go s.executeCodeWithInput(record.Id, language, code, testInput)
	} else {
		// 提交模式：查询题目获取全部测试用例
		p, err := s.problemDAO.GetProblemById(problemId)
		if err != nil || p == nil {
			_ = s.codeRunDAO.UpdateCodeRun(record.Id, map[string]interface{}{
				"status":    "runtime_error",
				"error_msg": "题目不存在",
			})
			return record, nil
		}
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

// executeCodeWithInput 测试模式：直接用给定输入运行代码（不依赖数据库题目）
func (s *CodeRunService) executeCodeWithInput(runId int64, language, code, input string) {
	_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{"status": "running"})

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

	const defaultTimeLimitMs int64 = 5000
	status, output, errMsg, timeCost, memUsed := s.runSingleCase(tmpDir, cfg, input, defaultTimeLimitMs)
	_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{
		"status":      status,
		"output":      output,
		"error_msg":   errMsg,
		"time_cost":   timeCost,
		"memory_used": memUsed,
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
	var casesToRun []testCase
	if runType == "test" {
		// 测试模式：只跑样例用例（is_sample=1）
		for _, tc := range testCases {
			if tc.IsSample == 1 {
				casesToRun = append(casesToRun, tc)
			}
		}
		if len(casesToRun) == 0 {
			casesToRun = testCases[:1] // 至少跑第一个
		}
	} else {
		// 提交模式：跑所有测试用例
		casesToRun = testCases
	}

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
	var outputLines []string
	finalStatus := "accepted"

	// 默认时间限制 1000ms
	const defaultTimeLimitMs int64 = 1000

	for i, tc := range casesToRun {
		status, output, errMsg, timeCost, memUsed := s.runSingleCase(tmpDir, cfg, tc.Input, defaultTimeLimitMs)
		totalTimeCost += timeCost
		if memUsed > maxMemoryUsed {
			maxMemoryUsed = memUsed
		}

		if status != "accepted" {
			finalStatus = status
			_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{
				"status":      finalStatus,
				"output":      output,
				"error_msg":   fmt.Sprintf("第 %d 个测试点失败\n%s", i+1, errMsg),
				"time_cost":   totalTimeCost,
				"memory_used": maxMemoryUsed,
			})
			return
		}
		outputLines = append(outputLines, strings.TrimSpace(output))
	}

	// 全部通过
	_ = s.codeRunDAO.UpdateCodeRun(runId, map[string]interface{}{
		"status":      finalStatus,
		"output":      strings.Join(outputLines, "\n---\n"),
		"error_msg":   "",
		"time_cost":   totalTimeCost,
		"memory_used": maxMemoryUsed,
	})
}

// runSingleCase 运行单个测试用例，返回 (status, output, errMsg, timeCostMs, memUsedKB)
func (s *CodeRunService) runSingleCase(tmpDir string, cfg langConfig, input string, timeLimitMs int64) (string, string, string, int64, int64) {
	runCtx, cancel := context.WithTimeout(context.Background(), time.Duration(timeLimitMs+2000)*time.Millisecond)
	defer cancel()

	runArgs := make([]string, len(cfg.RunCmd))
	copy(runArgs, cfg.RunCmd)
	cmd := exec.CommandContext(runCtx, runArgs[0], runArgs[1:]...)
	cmd.Dir = tmpDir
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
