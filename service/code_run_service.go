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
	"strconv"
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
	codeRunDAO     dao.CodeRunDAO
	problemDAO     dao.ProblemDAO
	profileService *StudentProfileService
}

// NewCodeRunService 创建代码运行服务
func NewCodeRunService() *CodeRunService {
	return &CodeRunService{
		codeRunDAO:     dao.NewCodeRunDAO(),
		problemDAO:     dao.NewProblemDAO(),
		profileService: NewStudentProfileService(),
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
		// 提交模式：执行 test_cases 全部用例，完成后异步更新学生画像
		go s.executeCodeAndUpdateProfile(record.Id, p, language, code, runType, studentId)
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

// SubmitTeacherCodeRun 教师提交代码运行任务（复用代码执行逻辑，记录 teacher_id）
func (s *CodeRunService) SubmitTeacherCodeRun(ctx context.Context, teacherId string, problemId int64, language, code, runType, testInput string) (*codeModel.CodeRun, error) {
	if _, ok := langConfigs[language]; !ok {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "不支持的编程语言: "+language)
	}
	if runType != "test" && runType != "submit" {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "run_type 必须为 test 或 submit")
	}

	p, err := s.problemDAO.GetProblemById(problemId)
	if err != nil || p == nil {
		return nil, errs.NewCommonError(errs.ErrBadRequest, "题目不存在")
	}

	record := &codeModel.CodeRun{
		ProblemId: problemId,
		StudentId: "",
		TeacherId: teacherId,
		Language:  language,
		Code:      code,
		RunType:   runType,
		Status:    "pending",
	}
	if err := s.codeRunDAO.CreateCodeRun(record); err != nil {
		return nil, errs.NewCommonError(errs.ErrInternal, "创建运行记录失败: "+err.Error())
	}

	if runType == "test" {
		go s.executeCodeWithShowcase(record.Id, p, language, code)
	} else {
		go s.executeCode(record.Id, p, language, code, runType)
	}

	return record, nil
}

// ListTeacherCodeRunRecords 查询教师某题的运行记录列表（倒序）
func (s *CodeRunService) ListTeacherCodeRunRecords(teacherId string, problemId int64, limit int) ([]*codeModel.CodeRun, error) {
	records, err := s.codeRunDAO.ListCodeRunsByTeacher(teacherId, problemId, limit)
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

// executeCodeAndUpdateProfile 执行代码并在完成后异步更新学生画像
func (s *CodeRunService) executeCodeAndUpdateProfile(runId int64, p *problem.Problem, language, code, runType, studentId string) {
	s.executeCode(runId, p, language, code, runType)

	// 异步更新学生画像（仅学生提交时触发）
	if studentId != "" && s.profileService != nil {
		go s.profileService.UpdateProfileAfterSubmit(studentId)
	}
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

// CheckCodeSyntax 仅编译代码检查语法错误（不运行），返回编译诊断信息
func (s *CodeRunService) CheckCodeSyntax(language, code string) (bool, []CompileDiagnostic, error) {
	cfg, ok := langConfigs[language]
	if !ok {
		return false, nil, errs.NewCommonError(errs.ErrBadRequest, "不支持的编程语言: "+language)
	}

	// 创建临时目录
	tmpDir, err := os.MkdirTemp("", "elysia_check_*")
	if err != nil {
		return false, nil, errs.NewCommonError(errs.ErrInternal, "创建临时目录失败")
	}
	defer os.RemoveAll(tmpDir)

	// 写入源代码文件
	srcFile := filepath.Join(tmpDir, cfg.FileName)
	if err := os.WriteFile(srcFile, []byte(code), 0644); err != nil {
		return false, nil, errs.NewCommonError(errs.ErrInternal, "写入代码文件失败")
	}

	// Python 是解释型语言，使用 py_compile 检查语法
	if language == "python" {
		return s.checkPythonSyntax(tmpDir, cfg.FileName)
	}

	// 编译型语言：执行编译命令
	if cfg.CompileCmd == nil {
		// 无编译命令（理论上不会到这里，Python 已在上面处理）
		return false, nil, nil
	}

	compileCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	compileCmd := exec.CommandContext(compileCtx, cfg.CompileCmd[0], cfg.CompileCmd[1:]...)
	compileCmd.Dir = tmpDir
	compileCmd.Env = buildEnv()
	var compileErr bytes.Buffer
	compileCmd.Stderr = &compileErr

	if err := compileCmd.Run(); err != nil {
		// 编译失败，解析错误信息
		errStr := compileErr.String()
		diagnostics := parseCompileErrors(errStr, language, cfg.FileName)
		return true, diagnostics, nil
	}

	// 编译成功，无错误
	return false, nil, nil
}

// CompileDiagnostic 编译诊断信息
type CompileDiagnostic struct {
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	EndLine   int    `json:"end_line"`
	EndColumn int    `json:"end_column"`
	Severity  string `json:"severity"` // error / warning
	Message   string `json:"message"`
}

// checkPythonSyntax 使用 pyflakes（优先）或 py_compile 检查 Python 语法
// pyflakes 可以一次性报告多处错误，py_compile 只能报告第一处
func (s *CodeRunService) checkPythonSyntax(tmpDir, fileName string) (bool, []CompileDiagnostic, error) {
	checkCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 优先尝试 pyflakes（能报告多处错误）
	pyflakesCmd := exec.CommandContext(checkCtx, "pyflakes", fileName)
	pyflakesCmd.Dir = tmpDir
	pyflakesCmd.Env = buildEnv()
	var pyflakesStdout, pyflakesStderr bytes.Buffer
	pyflakesCmd.Stdout = &pyflakesStdout
	pyflakesCmd.Stderr = &pyflakesStderr

	if err := pyflakesCmd.Run(); err != nil {
		// pyflakes 发现错误（exit code != 0）
		errOutput := pyflakesStdout.String() + pyflakesStderr.String()
		diagnostics := parsePyflakesErrors(errOutput, fileName)
		if len(diagnostics) > 0 {
			return true, diagnostics, nil
		}
		// pyflakes 可能未安装（command not found），回退到 py_compile
	}

	// pyflakes 无错误或不可用，回退到 py_compile
	checkCtx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()
	cmd := exec.CommandContext(checkCtx2, "python3", "-m", "py_compile", fileName)
	cmd.Dir = tmpDir
	cmd.Env = buildEnv()
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		errStr := stderr.String()
		diagnostics := parsePythonErrors(errStr, fileName)
		if len(diagnostics) == 0 {
			// 无法解析具体行号，返回编译器原始输出
			diagnostics = append(diagnostics, CompileDiagnostic{
				Line:      1,
				Column:    1,
				EndLine:   1,
				EndColumn: 1,
				Severity:  "error",
				Message:   strings.TrimSpace(errStr),
			})
		}
		return true, diagnostics, nil
	}

	return false, nil, nil
}

// parsePyflakesErrors 解析 pyflakes 错误输出（支持多处错误）
// 格式示例：
//   main.py:3:1 'os' imported but unused
//   main.py:5:5 undefined name 'foo'
func parsePyflakesErrors(errStr, fileName string) []CompileDiagnostic {
	var diagnostics []CompileDiagnostic
	baseName := filepath.Base(fileName)
	lines := strings.Split(errStr, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || !strings.Contains(line, baseName) {
			continue
		}

		// 格式：filename:line:col message 或 filename:line: message
		idx := strings.Index(line, baseName+":")
		if idx < 0 {
			continue
		}
		rest := line[idx+len(baseName)+1:]

		// 尝试解析 line:col message
		parts := strings.SplitN(rest, ":", 3)
		if len(parts) < 2 {
			continue
		}

		lineNum, err := strconv.Atoi(strings.TrimSpace(parts[0]))
		if err != nil {
			continue
		}

		colNum := 1
		message := ""
		if len(parts) >= 3 {
			if n, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
				colNum = n
				message = strings.TrimSpace(parts[2])
			} else {
				message = strings.TrimSpace(parts[1] + ":" + parts[2])
			}
		} else {
			message = strings.TrimSpace(parts[1])
		}

		if message == "" {
			continue
		}

		// pyflakes 的 "invalid syntax" 等是 error，其他（如 unused import）是 warning
		severity := "warning"
		if strings.Contains(strings.ToLower(message), "syntax") ||
			strings.Contains(strings.ToLower(message), "undefined") ||
			strings.Contains(strings.ToLower(message), "invalid") {
			severity = "error"
		}

		diagnostics = append(diagnostics, CompileDiagnostic{
			Line:      lineNum,
			Column:    colNum,
			EndLine:   lineNum,
			EndColumn: colNum,
			Severity:  severity,
			Message:   message,
		})
	}
	return diagnostics
}

// parsePythonErrors 解析 py_compile 的 Python 编译错误输出
// 格式示例：
//   File "main.py", line 3
//     print("hello"
//                  ^
//   SyntaxError: '(' was never closed
func parsePythonErrors(errStr, fileName string) []CompileDiagnostic {
	var diagnostics []CompileDiagnostic
	lines := strings.Split(errStr, "\n")

	for i, line := range lines {
		line = strings.TrimSpace(line)
		// 匹配 File "xxx", line N
		if strings.Contains(line, "File") && strings.Contains(line, "line") {
			lineNum := 1
			// 提取行号
			parts := strings.Split(line, "line ")
			if len(parts) >= 2 {
				numStr := strings.TrimSpace(parts[len(parts)-1])
				if n, err := strconv.Atoi(numStr); err == nil {
					lineNum = n
				}
			}
			// 查找后续的错误信息（通常是 SyntaxError: xxx）
			errMsg := ""
			for j := i + 1; j < len(lines); j++ {
				trimmed := strings.TrimSpace(lines[j])
				if strings.Contains(trimmed, "Error:") || strings.Contains(trimmed, "Error") {
					errMsg = trimmed
					break
				}
			}
			if errMsg == "" {
				errMsg = strings.TrimSpace(errStr)
			}
			diagnostics = append(diagnostics, CompileDiagnostic{
				Line:      lineNum,
				Column:    1,
				EndLine:   lineNum,
				EndColumn: 1,
				Severity:  "error",
				Message:   errMsg,
			})
		}
	}
	return diagnostics
}

// parseCompileErrors 解析编译型语言的编译错误输出
// 支持格式：
//   gcc/g++: main.cpp:10:5: error: expected ';' after expression
//   javac:   Main.java:5: error: ';' expected
//   go:      ./main.go:8:2: undefined: fmt.Printl
func parseCompileErrors(errStr, language, fileName string) []CompileDiagnostic {
	var diagnostics []CompileDiagnostic
	lines := strings.Split(errStr, "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		var d *CompileDiagnostic

		switch language {
		case "c", "cpp":
			d = parseGccError(line, fileName)
		case "java":
			d = parseJavacError(line, fileName)
		case "go":
			d = parseGoError(line, fileName)
		}

		if d != nil {
			diagnostics = append(diagnostics, *d)
		}
	}

	// 如果没有解析出具体行号，返回通用错误
	if len(diagnostics) == 0 && strings.TrimSpace(errStr) != "" {
		diagnostics = append(diagnostics, CompileDiagnostic{
			Line:      1,
			Column:    1,
			EndLine:   1,
			EndColumn: 1,
			Severity:  "error",
			Message:   strings.TrimSpace(errStr),
		})
	}

	return diagnostics
}

// parseGccError 解析 gcc/g++ 错误格式：main.cpp:10:5: error: expected ';'
// 也支持 fatal error、note 等格式
func parseGccError(line, fileName string) *CompileDiagnostic {
	baseName := filepath.Base(fileName)
	if !strings.Contains(line, baseName) {
		return nil
	}

	// 跳过 note: 行（补充说明，不是独立错误）
	if strings.Contains(line, ": note:") {
		return nil
	}

	// 格式：filename:line:col: severity: message
	idx := strings.Index(line, baseName+":")
	if idx < 0 {
		return nil
	}
	rest := line[idx+len(baseName)+1:]

	parts := strings.SplitN(rest, ":", 4)
	if len(parts) < 3 {
		return nil
	}

	lineNum, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	colNum, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil {
		return nil
	}
	if err2 != nil {
		colNum = 1
	}

	// 解析 severity 和 message
	remaining := strings.TrimSpace(parts[2])
	if len(parts) >= 4 {
		remaining = strings.TrimSpace(parts[2]) + ":" + parts[3]
	}

	severity := "error"
	message := remaining

	// 匹配 "error:"、"fatal error:"、"warning:" 前缀
	if strings.HasPrefix(remaining, "fatal error:") {
		severity = "error"
		message = strings.TrimPrefix(remaining, "fatal error:")
		message = strings.TrimSpace(message)
	} else if strings.HasPrefix(remaining, "error:") {
		severity = "error"
		message = strings.TrimPrefix(remaining, "error:")
		message = strings.TrimSpace(message)
	} else if strings.HasPrefix(remaining, "warning:") {
		severity = "warning"
		message = strings.TrimPrefix(remaining, "warning:")
		message = strings.TrimSpace(message)
	}

	return &CompileDiagnostic{
		Line:      lineNum,
		Column:    colNum,
		EndLine:   lineNum,
		EndColumn: colNum,
		Severity:  severity,
		Message:   message,
	}
}

// parseJavacError 解析 javac 错误格式：Main.java:5: error: ';' expected
// javac 输出中每个错误后面会跟一行代码和一行 ^ 指示符，这些行不包含文件名所以会被自动跳过
func parseJavacError(line, fileName string) *CompileDiagnostic {
	baseName := filepath.Base(fileName)
	if !strings.Contains(line, baseName) {
		return nil
	}

	// 跳过 javac 的汇总行，如 "2 errors" 或 "1 error"
	if strings.HasSuffix(strings.TrimSpace(line), "errors") || strings.HasSuffix(strings.TrimSpace(line), "error") {
		if !strings.Contains(line, ":") {
			return nil
		}
	}

	idx := strings.Index(line, baseName+":")
	if idx < 0 {
		return nil
	}
	rest := line[idx+len(baseName)+1:]

	parts := strings.SplitN(rest, ":", 3)
	if len(parts) < 2 {
		return nil
	}

	lineNum, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		return nil
	}

	severity := "error"
	message := strings.TrimSpace(strings.Join(parts[1:], ":"))
	if strings.HasPrefix(message, "error:") {
		severity = "error"
		message = strings.TrimPrefix(message, "error:")
		message = strings.TrimSpace(message)
	} else if strings.HasPrefix(message, "warning:") {
		severity = "warning"
		message = strings.TrimPrefix(message, "warning:")
		message = strings.TrimSpace(message)
	}

	return &CompileDiagnostic{
		Line:      lineNum,
		Column:    1,
		EndLine:   lineNum,
		EndColumn: 1,
		Severity:  severity,
		Message:   message,
	}
}

// parseGoError 解析 Go 编译错误格式：./main.go:8:2: undefined: fmt.Printl
func parseGoError(line, fileName string) *CompileDiagnostic {
	baseName := filepath.Base(fileName)
	// Go 错误可能以 ./ 开头
	if !strings.Contains(line, baseName) {
		return nil
	}

	idx := strings.Index(line, baseName+":")
	if idx < 0 {
		return nil
	}
	rest := line[idx+len(baseName)+1:]

	parts := strings.SplitN(rest, ":", 3)
	if len(parts) < 2 {
		return nil
	}

	lineNum, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err1 != nil {
		return nil
	}

	colNum := 1
	message := ""
	if len(parts) >= 3 {
		if n, err := strconv.Atoi(strings.TrimSpace(parts[1])); err == nil {
			colNum = n
			message = strings.TrimSpace(parts[2])
		} else {
			message = strings.TrimSpace(strings.Join(parts[1:], ":"))
		}
	} else {
		message = strings.TrimSpace(parts[1])
	}

	return &CompileDiagnostic{
		Line:      lineNum,
		Column:    colNum,
		EndLine:   lineNum,
		EndColumn: colNum,
		Severity:  "error",
		Message:   message,
	}
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
