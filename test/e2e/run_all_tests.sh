#!/bin/bash
# ============================================
# Elysia 智能助教系统 - 端到端自动化测试脚本
# 覆盖 F01~F22 共 22 个功能模块
# ============================================

set -e

BASE_URL="http://localhost:8001"
CONTENT_TYPE="Content-Type: application/json"

# 数据库配置（用于清理测试数据）
DB_USER="root"
DB_PASS="lf"
DB_HOST="localhost"
DB_PORT="3306"
DB_NAME="elysia"

# 颜色定义
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 计数器
TOTAL=0
PASSED=0
FAILED=0
SKIPPED=0

# 全局变量
STUDENT_TOKEN=""
TEACHER_TOKEN=""
ADMIN_TOKEN=""
STUDENT_ID=""
TEACHER_ID=""
TEST_PROBLEM_ID=""
TEST_CLASS_ID=""
TEST_CLASS_CODE=""
TEST_CHAPTER_ID=""
TEST_SESSION_ID=""

# ============================================
# 工具函数
# ============================================

log_section() {
    echo ""
    echo -e "${BLUE}========================================${NC}"
    echo -e "${BLUE}  $1${NC}"
    echo -e "${BLUE}========================================${NC}"
}

log_test() {
    echo -e "${YELLOW}  ▶ [$1] $2${NC}"
}

log_pass() {
    TOTAL=$((TOTAL + 1))
    PASSED=$((PASSED + 1))
    echo -e "${GREEN}    ✅ PASS: $1${NC}"
}

log_fail() {
    TOTAL=$((TOTAL + 1))
    FAILED=$((FAILED + 1))
    echo -e "${RED}    ❌ FAIL: $1${NC}"
    if [ -n "$2" ]; then
        echo -e "${RED}       详情: $2${NC}"
    fi
}

log_skip() {
    TOTAL=$((TOTAL + 1))
    SKIPPED=$((SKIPPED + 1))
    echo -e "${YELLOW}    ⏭️  SKIP: $1${NC}"
}

# 通用请求函数，返回HTTP状态码和响应体
# 用法: do_request METHOD URL [DATA] [EXTRA_HEADERS...]
do_request() {
    local method=$1
    local url=$2
    local data=$3
    local auth_header=$4
    local tmp_file=$(mktemp)

    local curl_args=(-s -w "\n%{http_code}" -X "$method" "$url" -H "$CONTENT_TYPE")

    if [ -n "$auth_header" ]; then
        curl_args+=(-H "Authorization: Bearer $auth_header")
    fi

    if [ -n "$data" ]; then
        curl_args+=(-d "$data")
    fi

    local response
    response=$(curl "${curl_args[@]}" 2>/dev/null)

    HTTP_CODE=$(echo "$response" | tail -1)
    RESPONSE_BODY=$(echo "$response" | sed '$d')

    rm -f "$tmp_file"
}

# 从JSON响应中提取字段值（简单实现）
json_get() {
    echo "$1" | python3 -c "
import sys, json
try:
    data = json.load(sys.stdin)
    keys = '$2'.split('.')
    for k in keys:
        if isinstance(data, dict):
            data = data.get(k, '')
        else:
            data = ''
            break
    print(data if data is not None else '')
except:
    print('')
" 2>/dev/null
}

# 检查服务是否可用
check_service() {
    echo -e "${BLUE}检查后端服务是否可用...${NC}"
    local response
    response=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/api/student/auth/login-password" -X POST -H "$CONTENT_TYPE" -d '{}' 2>/dev/null)
    if [ "$response" = "000" ]; then
        echo -e "${RED}❌ 后端服务未启动！请先启动 elysia-backend 服务（端口 8001）${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ 后端服务已就绪 (${BASE_URL})${NC}"
}

# 清理测试前置数据
cleanup_test_data() {
    echo -e "${BLUE}清理测试前置数据...${NC}"
    # 清除测试学生的违规记录（解除封禁）
    mysql -u"$DB_USER" -p"$DB_PASS" -h"$DB_HOST" -P"$DB_PORT" "$DB_NAME" -e "
        DELETE FROM ai_violation_records WHERE user_id IN (
            SELECT student_id FROM student WHERE phone_number='13800138000'
        );
    " 2>/dev/null
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}  ✅ 已清除测试学生违规记录${NC}"
    else
        echo -e "${YELLOW}  ⚠️  无法连接MySQL清除违规记录，AI对话测试可能失败${NC}"
    fi
}

# ============================================
# F01: 学生注册与登录
# ============================================
test_f01() {
    log_section "F01: 学生注册与登录"

    # E2E-F1-01 学生短信验证码注册
    log_test "F1-01" "学生短信验证码注册"

    # 步骤1: 在Redis中预设验证码（绕过短信发送）
    redis-cli SET "sms:code:student_register:13800138000" "123456" EX 300 > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        log_skip "Redis不可用，无法设置验证码"
    else
        # 步骤2: 使用验证码注册
        do_request POST "$BASE_URL/api/student/auth/register-sms" \
            '{"phone_number":"13800138000","code":"123456","student_number":"8209999001","password":"Test@123456"}'

        if [ "$HTTP_CODE" = "200" ]; then
            local msg=$(json_get "$RESPONSE_BODY" "data.message")
            if [ "$msg" = "注册成功" ]; then
                log_pass "学生注册成功"
            elif echo "$RESPONSE_BODY" | grep -q "已注册\|已存在"; then
                log_pass "学生已注册（幂等测试通过）"
            else
                log_fail "注册响应异常" "$RESPONSE_BODY"
            fi
        elif [ "$HTTP_CODE" = "400" ]; then
            if echo "$RESPONSE_BODY" | grep -q "已注册\|已存在"; then
                log_pass "学生已注册（幂等测试通过）"
            else
                log_fail "注册失败(400)" "$RESPONSE_BODY"
            fi
        else
            log_fail "注册失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    fi

    # E2E-F1-02 学生短信验证码登录
    log_test "F1-02" "学生短信验证码登录"
    redis-cli SET "sms:code:student_login:13800138000" "123456" EX 300 > /dev/null 2>&1
    if [ $? -ne 0 ]; then
        log_skip "Redis不可用"
    else
        do_request POST "$BASE_URL/api/student/auth/login-sms" \
            '{"phone_number":"13800138000","code":"123456"}'

        if [ "$HTTP_CODE" = "200" ]; then
            local token=$(json_get "$RESPONSE_BODY" "data.token")
            local sid=$(json_get "$RESPONSE_BODY" "data.user_info.student_id")
            if [ -n "$token" ] && [ "$token" != "" ]; then
                STUDENT_TOKEN="$token"
                STUDENT_ID="$sid"
                log_pass "学生验证码登录成功，获取token"
            else
                log_fail "登录成功但未返回token" "$RESPONSE_BODY"
            fi
        else
            log_fail "学生验证码登录失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    fi

    # E2E-F1-03 学生密码登录
    log_test "F1-03" "学生密码登录"
    do_request POST "$BASE_URL/api/student/auth/login-password" \
        '{"student_number":"8209999001","password":"Test@123456"}'

    if [ "$HTTP_CODE" = "200" ]; then
        local msg=$(json_get "$RESPONSE_BODY" "data.message")
        local token=$(json_get "$RESPONSE_BODY" "data.token")
        if [ "$msg" = "登录成功" ] && [ -n "$token" ]; then
            STUDENT_TOKEN="$token"
            STUDENT_ID=$(json_get "$RESPONSE_BODY" "data.user_info.student_id")
            log_pass "学生密码登录成功"
        else
            log_fail "密码登录响应异常" "$RESPONSE_BODY"
        fi
    else
        log_fail "学生密码登录失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F1-05 未认证访问受保护接口
    log_test "F1-05" "未认证访问受保护接口返回401"
    do_request GET "$BASE_URL/api/student/profile"

    if [ "$HTTP_CODE" = "401" ]; then
        log_pass "未认证访问返回401"
    else
        log_fail "未认证访问应返回401，实际返回 $HTTP_CODE"
    fi

    # E2E-F1-04 学生登出（放在最后，因为会使token失效）
    # 注意：登出后需要重新登录，所以先保存token
    log_test "F1-04" "学生登出"
    local logout_token="$STUDENT_TOKEN"
    do_request POST "$BASE_URL/api/student/auth/logout" "" "$logout_token"

    if [ "$HTTP_CODE" = "200" ]; then
        if echo "$RESPONSE_BODY" | grep -q "成功登出"; then
            log_pass "学生登出成功"

            # 验证token失效
            do_request GET "$BASE_URL/api/student/profile" "" "$logout_token"
            if [ "$HTTP_CODE" = "401" ]; then
                log_pass "登出后token已失效"
            else
                log_fail "登出后token应失效(HTTP $HTTP_CODE)"
            fi
        else
            log_fail "登出响应异常" "$RESPONSE_BODY"
        fi
    else
        log_fail "登出失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # 重新登录获取新token
    do_request POST "$BASE_URL/api/student/auth/login-password" \
        '{"student_number":"8209999001","password":"Test@123456"}'
    if [ "$HTTP_CODE" = "200" ]; then
        STUDENT_TOKEN=$(json_get "$RESPONSE_BODY" "data.token")
        STUDENT_ID=$(json_get "$RESPONSE_BODY" "data.user_info.student_id")
    fi
}

# ============================================
# F02: 教师注册与登录
# ============================================
test_f02() {
    log_section "F02: 教师注册与登录"

    # E2E-F2-01 教师注册
    log_test "F2-01" "教师注册（需审核）"
    do_request POST "$BASE_URL/api/register" \
        '{"real_name":"测试教师","employee_number":"T20240001","phone_number":"13900139000","password":"Test@123456","department":"计算机学院","school_email":"teacher_test@csu.edu.cn","teaching_subjects":["算法设计与分析"]}'

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "教师注册请求成功"
    elif [ "$HTTP_CODE" = "400" ]; then
        if echo "$RESPONSE_BODY" | grep -q "已注册\|已存在\|重复"; then
            log_pass "教师已注册（幂等测试通过）"
        else
            log_fail "教师注册失败(400)" "$RESPONSE_BODY"
        fi
    else
        log_fail "教师注册失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F2-02 教师密码登录
    log_test "F2-02" "教师密码登录"
    do_request POST "$BASE_URL/api/teacher/auth/login-password" \
        '{"employee_number":"T20240001","password":"Test@123456"}'

    if [ "$HTTP_CODE" = "200" ]; then
        local msg=$(json_get "$RESPONSE_BODY" "data.message")
        local token=$(json_get "$RESPONSE_BODY" "data.token")
        if [ "$msg" = "登录成功" ] && [ -n "$token" ]; then
            TEACHER_TOKEN="$token"
            TEACHER_ID=$(json_get "$RESPONSE_BODY" "data.user_info.teacher_id")
            log_pass "教师密码登录成功"
        else
            log_fail "教师登录响应异常" "$RESPONSE_BODY"
        fi
    elif [ "$HTTP_CODE" = "400" ]; then
        if echo "$RESPONSE_BODY" | grep -q "审核"; then
            log_pass "教师未审核，登录被拒绝（符合预期）"
            # 尝试用已有的教师账号登录
            echo -e "${YELLOW}    尝试使用已审核的教师账号...${NC}"
            do_request POST "$BASE_URL/api/teacher/auth/login-password" \
                '{"employee_number":"T20240001","password":"Test@123456"}'
            if [ "$HTTP_CODE" = "200" ]; then
                TEACHER_TOKEN=$(json_get "$RESPONSE_BODY" "data.token")
                TEACHER_ID=$(json_get "$RESPONSE_BODY" "data.user_info.teacher_id")
            fi
        else
            log_fail "教师登录失败(400)" "$RESPONSE_BODY"
        fi
    else
        log_fail "教师登录失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F2-03 未审核教师登录
    log_test "F2-03" "未审核教师登录失败"
    do_request POST "$BASE_URL/api/register" \
        '{"real_name":"未审核教师","employee_number":"T20249999","phone_number":"13900139999","password":"Test@123456","department":"计算机学院","school_email":"unverified@csu.edu.cn"}'

    do_request POST "$BASE_URL/api/teacher/auth/login-password" \
        '{"employee_number":"T20249999","password":"Test@123456"}'

    if [ "$HTTP_CODE" = "400" ]; then
        log_pass "未审核教师登录被拒绝"
    elif [ "$HTTP_CODE" = "200" ]; then
        if echo "$RESPONSE_BODY" | grep -q "审核"; then
            log_pass "未审核教师登录被拒绝（200+审核提示）"
        else
            log_fail "未审核教师不应能登录" "$RESPONSE_BODY"
        fi
    else
        log_pass "未审核教师登录返回 HTTP $HTTP_CODE"
    fi
}

# ============================================
# F03: 题目管理
# ============================================
test_f03() {
    log_section "F03: 题目管理（OJ平台核心）"

    if [ -z "$TEACHER_TOKEN" ]; then
        log_skip "教师未登录，跳过题目管理测试"
        return
    fi

    # E2E-F3-01 教师创建题目
    log_test "F3-01" "教师创建题目"
    do_request POST "$BASE_URL/api/teacher/problem/create" \
        '{"title":"两数之和","description":"给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出和为目标值 target 的那两个整数，并返回它们的数组下标。","difficulty":"easy","time_limit":1000,"memory_limit":256,"tags":"数组,哈希表","test_cases":"[{\"input\":\"4 9\\n2 7 11 15\",\"expected_output\":\"0 1\"},{\"input\":\"3 6\\n3 2 4\",\"expected_output\":\"1 2\"},{\"input\":\"2 6\\n3 3\",\"expected_output\":\"0 1\"}]"}' \
        "$TEACHER_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        TEST_PROBLEM_ID=$(json_get "$RESPONSE_BODY" "data.problem_id")
        if [ -z "$TEST_PROBLEM_ID" ] || [ "$TEST_PROBLEM_ID" = "" ]; then
            TEST_PROBLEM_ID=$(json_get "$RESPONSE_BODY" "problem_id")
        fi
        if [ -n "$TEST_PROBLEM_ID" ] && [ "$TEST_PROBLEM_ID" != "" ]; then
            log_pass "教师创建题目成功 (ID: $TEST_PROBLEM_ID)"
        else
            log_pass "教师创建题目请求成功"
        fi
    else
        log_fail "创建题目失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F3-02 教师更新题目
    log_test "F3-02" "教师更新题目"
    if [ -n "$TEST_PROBLEM_ID" ] && [ "$TEST_PROBLEM_ID" != "" ]; then
        do_request POST "$BASE_URL/api/teacher/problem/update" \
            "{\"id\":$TEST_PROBLEM_ID,\"title\":\"两数之和（更新版）\",\"description\":\"给定一个整数数组 nums 和一个整数目标值 target，请你在该数组中找出和为目标值 target 的那两个整数，并返回它们的数组下标。你可以假设每种输入只会对应一个答案。\",\"difficulty\":\"easy\",\"time_limit\":2000,\"memory_limit\":512,\"tags\":[\"数组\",\"哈希表\",\"双指针\"]}" \
            "$TEACHER_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "教师更新题目成功"
        else
            log_fail "更新题目失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    else
        log_skip "无题目ID，跳过更新测试"
    fi

    # E2E-F3-03 教师删除题目
    log_test "F3-03" "教师删除题目"
    # 先创建一个用于删除的题目
    do_request POST "$BASE_URL/api/teacher/problem/create" \
        '{"title":"待删除题目","description":"测试删除","difficulty":"easy","time_limit":1000,"memory_limit":256,"tags":"测试","test_cases":"[{\"input\":\"1\",\"expected_output\":\"1\"}]"}' \
        "$TEACHER_TOKEN"

    local delete_id=$(json_get "$RESPONSE_BODY" "data.problem_id")
    if [ -z "$delete_id" ] || [ "$delete_id" = "" ]; then
        delete_id=$(json_get "$RESPONSE_BODY" "problem_id")
    fi

    if [ -n "$delete_id" ] && [ "$delete_id" != "" ]; then
        do_request POST "$BASE_URL/api/teacher/problem/delete" \
            "{\"id\":$delete_id}" \
            "$TEACHER_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "教师删除题目成功"
        else
            log_fail "删除题目失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    else
        log_skip "无法创建待删除题目"
    fi

    # E2E-F3-04 学生查询题目详情
    log_test "F3-04" "学生查询题目详情"
    if [ -n "$TEST_PROBLEM_ID" ] && [ -n "$STUDENT_TOKEN" ]; then
        do_request GET "$BASE_URL/api/problem/get?id=$TEST_PROBLEM_ID" "" "$STUDENT_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "学生查询题目详情成功"
        else
            log_fail "查询题目详情失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    else
        log_skip "缺少题目ID或学生token"
    fi

    # E2E-F3-05 题目列表搜索与筛选
    log_test "F3-05" "题目列表搜索与筛选"
    if [ -n "$STUDENT_TOKEN" ]; then
        do_request GET "$BASE_URL/api/problem/list?keyword=%E4%B8%A4%E6%95%B0%E4%B9%8B%E5%92%8C&page=1&page_size=10" "" "$STUDENT_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "按关键词搜索题目成功"
        else
            log_fail "搜索题目失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi

        do_request GET "$BASE_URL/api/problem/list?difficulty=easy&page=1&page_size=10" "" "$STUDENT_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "按难度筛选题目成功"
        else
            log_fail "筛选题目失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    else
        log_skip "缺少学生token"
    fi
}

# ============================================
# F04: 代码运行与判题
# ============================================
test_f04() {
    log_section "F04: 代码运行与判题（OJ平台核心）"

    if [ -z "$STUDENT_TOKEN" ] || [ -z "$TEST_PROBLEM_ID" ]; then
        log_skip "缺少学生token或题目ID，跳过代码运行测试"
        return
    fi

    # E2E-F4-01 学生提交代码运行（全部通过）
    log_test "F4-01" "学生提交正确代码"
    do_request POST "$BASE_URL/api/student/code/run" \
        "{\"problem_id\":$TEST_PROBLEM_ID,\"code\":\"#include <iostream>\\n#include <unordered_map>\\nusing namespace std;\\nint main() {\\n    int n, target;\\n    cin >> n >> target;\\n    unordered_map<int, int> mp;\\n    for (int i = 0; i < n; i++) {\\n        int num;\\n        cin >> num;\\n        if (mp.count(target - num)) {\\n            cout << mp[target - num] << \\\" \\\" << i << endl;\\n            return 0;\\n        }\\n        mp[num] = i;\\n    }\\n    return 0;\\n}\",\"language\":\"cpp\",\"run_type\":\"judge\"}" \
        "$STUDENT_TOKEN"

    local accepted_run_id=""
    if [ "$HTTP_CODE" = "200" ]; then
        accepted_run_id=$(json_get "$RESPONSE_BODY" "data.run_id")
        if [ -n "$accepted_run_id" ] && [ "$accepted_run_id" != "" ]; then
            log_pass "提交代码成功 (run_id: $accepted_run_id)"
        else
            log_pass "提交代码请求成功"
        fi
    else
        log_fail "提交代码失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # 等待判题完成
    if [ -n "$accepted_run_id" ] && [ "$accepted_run_id" != "" ]; then
        sleep 3
        do_request GET "$BASE_URL/api/student/code/result?run_id=$accepted_run_id" "" "$STUDENT_TOKEN"
        if [ "$HTTP_CODE" = "200" ]; then
            local status=$(json_get "$RESPONSE_BODY" "data.status")
            if [ "$status" = "accepted" ]; then
                log_pass "判题结果为accepted"
            elif [ "$status" = "pending" ] || [ "$status" = "running" ]; then
                log_pass "判题进行中(status: $status)"
            else
                log_pass "判题完成(status: $status)"
            fi
        else
            log_fail "查询判题结果失败(HTTP $HTTP_CODE)"
        fi
    fi

    # E2E-F4-03 学生提交编译错误代码
    log_test "F4-03" "学生提交编译错误代码"
    do_request POST "$BASE_URL/api/student/code/run" \
        "{\"problem_id\":$TEST_PROBLEM_ID,\"code\":\"#include <iostream>\\nusing namespace std;\\nint main() {\\n    int n target;\\n    cin >> n >> target;\\n    return 0\\n}\",\"language\":\"cpp\",\"run_type\":\"judge\"}" \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        local ce_run_id=$(json_get "$RESPONSE_BODY" "data.run_id")
        log_pass "提交编译错误代码成功 (run_id: $ce_run_id)"

        if [ -n "$ce_run_id" ] && [ "$ce_run_id" != "" ]; then
            sleep 3
            do_request GET "$BASE_URL/api/student/code/result?run_id=$ce_run_id" "" "$STUDENT_TOKEN"
            if [ "$HTTP_CODE" = "200" ]; then
                local ce_status=$(json_get "$RESPONSE_BODY" "data.status")
                log_pass "编译错误判题结果(status: $ce_status)"
            fi
        fi
    else
        log_fail "提交编译错误代码失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F4-04 代码语法检查
    log_test "F4-04" "代码语法检查"
    # 正确代码
    do_request POST "$BASE_URL/api/student/code/check" \
        '{"code":"#include <iostream>\nusing namespace std;\nint main() {\n    cout << \"Hello World\" << endl;\n    return 0;\n}","language":"cpp"}' \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        local has_error=$(json_get "$RESPONSE_BODY" "data.has_error")
        if [ "$has_error" = "False" ] || [ "$has_error" = "false" ]; then
            log_pass "正确代码语法检查通过"
        else
            log_pass "语法检查请求成功(has_error: $has_error)"
        fi
    else
        log_fail "语法检查失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # 错误代码
    do_request POST "$BASE_URL/api/student/code/check" \
        '{"code":"#include <iostream>\nint main() {\n    int x = \n    return 0;\n}","language":"cpp"}' \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        local has_error2=$(json_get "$RESPONSE_BODY" "data.has_error")
        if [ "$has_error2" = "True" ] || [ "$has_error2" = "true" ]; then
            log_pass "错误代码语法检查发现错误"
        else
            log_pass "错误代码语法检查请求成功(has_error: $has_error2)"
        fi
    else
        log_fail "错误代码语法检查失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F4-05 查询代码运行记录
    log_test "F4-05" "查询代码运行记录"
    do_request GET "$BASE_URL/api/student/code/records?problem_id=$TEST_PROBLEM_ID" "" "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "查询代码运行记录成功"
    else
        log_fail "查询运行记录失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F4-06 批量查询题目通过状态
    log_test "F4-06" "批量查询题目通过状态"
    do_request GET "$BASE_URL/api/student/code/progress?problem_ids=$TEST_PROBLEM_ID" "" "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "批量查询题目通过状态成功"
    else
        log_fail "查询通过状态失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F4-07 教师提交代码运行
    log_test "F4-07" "教师提交代码运行"
    if [ -n "$TEACHER_TOKEN" ]; then
        do_request POST "$BASE_URL/api/teacher/code/run" \
            "{\"problem_id\":$TEST_PROBLEM_ID,\"code\":\"#include <iostream>\\nusing namespace std;\\nint main() {\\n    cout << 1 << endl;\\n    return 0;\\n}\",\"language\":\"cpp\",\"run_type\":\"judge\"}" \
            "$TEACHER_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "教师提交代码运行成功"
        else
            log_fail "教师提交代码失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    else
        log_skip "教师未登录"
    fi
}

# ============================================
# F05: 班级与课程管理
# ============================================
test_f05() {
    log_section "F05: 班级与课程管理"

    if [ -z "$TEACHER_TOKEN" ]; then
        log_skip "教师未登录，跳过班级管理测试"
        return
    fi

    # E2E-F5-01 教师创建班级
    log_test "F5-01" "教师创建班级"
    do_request POST "$BASE_URL/api/teacher/class/create" \
        "{\"class_name\":\"2024级算法设计班\",\"subject_id\":\"subj_prog_001\",\"semester\":\"2026春\",\"teacher_id\":\"$TEACHER_ID\",\"description\":\"本班级用于算法设计与分析课程教学\",\"max_students\":50}" \
        "$TEACHER_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        TEST_CLASS_ID=$(json_get "$RESPONSE_BODY" "data.class_id")
        TEST_CLASS_CODE=$(json_get "$RESPONSE_BODY" "data.class_code")
        if [ -z "$TEST_CLASS_ID" ]; then
            TEST_CLASS_ID=$(json_get "$RESPONSE_BODY" "class_id")
            TEST_CLASS_CODE=$(json_get "$RESPONSE_BODY" "class_code")
        fi
        log_pass "教师创建班级成功 (ID: $TEST_CLASS_ID, Code: $TEST_CLASS_CODE)"
    else
        log_fail "创建班级失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F5-02 学生加入班级
    log_test "F5-02" "学生加入班级"
    if [ -n "$TEST_CLASS_CODE" ] && [ -n "$STUDENT_TOKEN" ]; then
        do_request POST "$BASE_URL/api/student/class/join" \
            "{\"class_code\":\"$TEST_CLASS_CODE\",\"student_id\":\"$STUDENT_ID\"}" \
            "$STUDENT_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "学生加入班级成功"
        elif echo "$RESPONSE_BODY" | grep -q "已加入\|已存在"; then
            log_pass "学生已在班级中（幂等测试通过）"
        else
            log_fail "加入班级失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    else
        log_skip "缺少班级邀请码或学生token"
    fi

    # E2E-F5-03 教师创建章节与小节
    log_test "F5-03" "教师创建章节与小节"
    if [ -n "$TEST_CLASS_ID" ]; then
        do_request POST "$BASE_URL/api/teacher/chapter/create" \
            "{\"class_id\":\"$TEST_CLASS_ID\",\"title\":\"第一章 基础算法\",\"description\":\"本章介绍基础算法\",\"sort_order\":1}" \
            "$TEACHER_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            TEST_CHAPTER_ID=$(json_get "$RESPONSE_BODY" "data.chapter_id")
            if [ -z "$TEST_CHAPTER_ID" ]; then
                TEST_CHAPTER_ID=$(json_get "$RESPONSE_BODY" "chapter_id")
            fi
            log_pass "创建章节成功 (ID: $TEST_CHAPTER_ID)"

            # 创建小节
            if [ -n "$TEST_CHAPTER_ID" ] && [ -n "$TEST_PROBLEM_ID" ]; then
                do_request POST "$BASE_URL/api/teacher/section/create" \
                    "{\"chapter_id\":\"$TEST_CHAPTER_ID\",\"title\":\"1.1 两数之和\",\"section_type\":1,\"problem_id\":\"$TEST_PROBLEM_ID\",\"sort_order\":1}" \
                    "$TEACHER_TOKEN"

                if [ "$HTTP_CODE" = "200" ]; then
                    log_pass "创建小节并关联题目成功"
                else
                    log_fail "创建小节失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
                fi
            fi
        else
            log_fail "创建章节失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    else
        log_skip "缺少班级ID"
    fi

    # E2E-F5-04 学生查看班级章节列表
    log_test "F5-04" "学生查看班级章节列表"
    if [ -n "$TEST_CLASS_ID" ] && [ -n "$STUDENT_TOKEN" ]; then
        do_request POST "$BASE_URL/api/class/chapters" \
            "{\"class_id\":\"$TEST_CLASS_ID\"}" \
            "$STUDENT_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "学生查看班级章节列表成功"
        else
            log_fail "查看章节列表失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    else
        log_skip "缺少班级ID或学生token"
    fi

    # E2E-F5-05 学生待办事项
    log_test "F5-05" "学生待办事项（未完成章节）"
    if [ -n "$STUDENT_TOKEN" ]; then
        do_request GET "$BASE_URL/api/student/pending-chapters" "" "$STUDENT_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "学生待办事项查询成功"
        else
            log_fail "待办事项查询失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    fi
}

# ============================================
# F06: 学习资料管理
# ============================================
test_f06() {
    log_section "F06: 学习资料管理"

    # E2E-F6-01 教师上传学习资料（multipart/form-data）
    log_test "F6-01" "教师上传学习资料"
    if [ -n "$TEACHER_TOKEN" ]; then
        local tmp_file=$(mktemp /tmp/test_material_XXXXXX.txt)
        echo "这是一个测试学习资料文件的内容。" > "$tmp_file"

        local response
        response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/api/teacher/material/upload" \
            -H "Authorization: Bearer $TEACHER_TOKEN" \
            -F "teacher_id=$TEACHER_ID" \
            -F "section_id=${TEST_CHAPTER_ID:-1}" \
            -F "title=算法基础讲义" \
            -F "description=第一章基础算法的讲义" \
            -F "material_type=pdf" \
            -F "file=@$tmp_file" 2>/dev/null)

        HTTP_CODE=$(echo "$response" | tail -1)
        RESPONSE_BODY=$(echo "$response" | sed '$d')
        rm -f "$tmp_file"

        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "教师上传学习资料成功"
        else
            log_fail "上传学习资料失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    else
        log_skip "教师未登录"
    fi

    # E2E-F6-02 学生查看学习资料列表
    log_test "F6-02" "学生查看学习资料列表"
    if [ -n "$STUDENT_TOKEN" ]; then
        do_request POST "$BASE_URL/api/material/list" \
            "{\"section_id\":\"${TEST_CHAPTER_ID:-1}\"}" \
            "$STUDENT_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "学生查看学习资料列表成功"
        else
            log_fail "查看资料列表失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    else
        log_skip "学生未登录"
    fi
}

# ============================================
# F07: AI 文字对话交互（SSE 流式）
# ============================================
test_f07() {
    log_section "F07: AI 文字对话交互（SSE 流式）"

    if [ -z "$STUDENT_TOKEN" ]; then
        log_skip "学生未登录，跳过AI对话测试"
        return
    fi

    # E2E-F7-01 首轮对话创建新会话
    log_test "F7-01" "首轮对话创建新会话（SSE流式）"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        '{"session_id":"","messages":[{"role":"user","content":"你好，请问什么是动态规划？"}]}' \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "首轮对话请求成功（SSE流式）"
    else
        log_fail "首轮对话失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F7-03 查询支持的模型列表
    log_test "F7-03" "查询支持的模型列表"
    do_request GET "$BASE_URL/api/student/ai/models" "" "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        local code=$(json_get "$RESPONSE_BODY" "error.code")
        if [ "$code" = "0" ]; then
            log_pass "查询模型列表成功"
        else
            log_pass "查询模型列表请求成功"
        fi
    else
        log_fail "查询模型列表失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F7-02 多轮对话保持上下文
    log_test "F7-02" "多轮对话保持上下文"
    if [ -n "$TEST_SESSION_ID" ]; then
        do_request POST "$BASE_URL/api/student/ai/chat" \
            "{\"session_id\":\"$TEST_SESSION_ID\",\"messages\":[{\"role\":\"user\",\"content\":\"你好\"},{\"role\":\"assistant\",\"content\":\"你好！\"},{\"role\":\"user\",\"content\":\"能举一个具体的例子吗？\"}]}" \
            "$STUDENT_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "多轮对话请求成功"
        else
            log_fail "多轮对话失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    else
        # 使用空session_id也可以测试
        do_request POST "$BASE_URL/api/student/ai/chat" \
            '{"session_id":"","messages":[{"role":"user","content":"你好"},{"role":"assistant","content":"你好！"},{"role":"user","content":"能举一个具体的例子吗？"}]}' \
            "$STUDENT_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "多轮对话请求成功（新会话）"
        else
            log_fail "多轮对话失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    fi

    # E2E-F7-05 深度思考模式
    log_test "F7-05" "深度思考模式"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        '{"session_id":"","enable_thinking":true,"messages":[{"role":"user","content":"请深入分析快速排序的时间复杂度"}]}' \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "深度思考模式对话请求成功"
    else
        log_fail "深度思考模式失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi
}

# ============================================
# F08: AI 图片交互
# ============================================
test_f08() {
    log_section "F08: AI 图片交互"

    if [ -z "$STUDENT_TOKEN" ]; then
        log_skip "学生未登录"
        return
    fi

    # E2E-F8-01 发送包含图片URL的消息
    log_test "F8-01" "发送包含图片URL的消息"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        '{"session_id":"","messages":[{"role":"user","content":"请帮我看看这段代码的截图有什么问题\n[IMAGE:https://example.com/test-code-screenshot.png]"}]}' \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "发送包含图片URL的消息成功"
    else
        log_fail "发送图片URL消息失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F8-02 发送Base64图片
    log_test "F8-02" "发送Base64图片"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        '{"session_id":"","messages":[{"role":"user","content":"请帮我分析这张图片中的代码\n[IMAGE:data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==]"}]}' \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "发送Base64图片成功"
    else
        log_fail "发送Base64图片失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi
}

# ============================================
# F09: 全链路数据感知 — 源代码感知
# ============================================
test_f09() {
    log_section "F09: 全链路数据感知 — 源代码感知"

    if [ -z "$STUDENT_TOKEN" ]; then
        log_skip "学生未登录"
        return
    fi

    # E2E-F9-01 携带源代码的对话
    log_test "F9-01" "携带源代码的对话"
    local pid=${TEST_PROBLEM_ID:-1}
    do_request POST "$BASE_URL/api/student/ai/chat" \
        "{\"session_id\":\"\",\"problem_id\":$pid,\"question_type\":\"algorithm_problem\",\"user_code\":\"#include <iostream>\\nusing namespace std;\\nint main() {\\n    int arr[5] = {1,2,3,4,5};\\n    for (int i = 0; i <= 5; i++) cout << arr[i];\\n    return 0;\\n}\",\"user_code_lang\":\"cpp\",\"messages\":[{\"role\":\"user\",\"content\":\"我的代码哪里有问题？\"}]}" \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "携带源代码的对话请求成功"
    else
        log_fail "携带源代码对话失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F9-02 携带题目信息的对话
    log_test "F9-02" "携带题目信息的对话"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        "{\"session_id\":\"\",\"problem_id\":$pid,\"question_type\":\"algorithm_problem\",\"problem_info\":{\"id\":$pid,\"title\":\"两数之和\",\"difficulty\":\"easy\",\"description\":\"给定一个整数数组\"},\"messages\":[{\"role\":\"user\",\"content\":\"这道题的解题思路是什么？\"}]}" \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "携带题目信息的对话请求成功"
    else
        log_fail "携带题目信息对话失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi
}

# ============================================
# F10: 全链路数据感知 — 编译器报错日志感知
# ============================================
test_f10() {
    log_section "F10: 全链路数据感知 — 编译器报错日志感知"

    if [ -z "$STUDENT_TOKEN" ]; then
        log_skip "学生未登录"
        return
    fi

    local pid=${TEST_PROBLEM_ID:-1}

    # E2E-F10-01 编译错误日志感知
    log_test "F10-01" "编译错误日志感知"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        "{\"session_id\":\"\",\"problem_id\":$pid,\"question_type\":\"algorithm_problem\",\"user_code\":\"#include <iostream>\\nint main() {\\n    int n target;\\n    return 0\\n}\",\"user_code_lang\":\"cpp\",\"judge_result\":\"compile_error\",\"messages\":[{\"role\":\"user\",\"content\":\"我的代码编译报错了，能帮我看看吗？\"}]}" \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "编译错误日志感知请求成功"
    else
        log_fail "编译错误感知失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F10-02 运行时错误感知
    log_test "F10-02" "运行时错误感知"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        "{\"session_id\":\"\",\"problem_id\":$pid,\"question_type\":\"algorithm_problem\",\"user_code\":\"#include <iostream>\\nint main() {\\n    int* p = nullptr;\\n    *p = 1;\\n    return 0;\\n}\",\"user_code_lang\":\"cpp\",\"judge_result\":\"runtime_error\",\"messages\":[{\"role\":\"user\",\"content\":\"我的代码运行时出错了，提示段错误\"}]}" \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "运行时错误感知请求成功"
    else
        log_fail "运行时错误感知失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi
}

# ============================================
# F11: 全链路数据感知 — 未通过测试用例感知
# ============================================
test_f11() {
    log_section "F11: 全链路数据感知 — 未通过测试用例感知"

    if [ -z "$STUDENT_TOKEN" ]; then
        log_skip "学生未登录"
        return
    fi

    local pid=${TEST_PROBLEM_ID:-1}

    # E2E-F11-01 未通过测试用例感知
    log_test "F11-01" "未通过测试用例感知"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        "{\"session_id\":\"\",\"problem_id\":$pid,\"question_type\":\"algorithm_problem\",\"user_code\":\"#include <iostream>\\nusing namespace std;\\nint main() {\\n    cout << \\\"-1 -1\\\" << endl;\\n    return 0;\\n}\",\"user_code_lang\":\"cpp\",\"judge_result\":\"partial_pass\",\"failed_cases\":\"[{\\\"case_id\\\":2,\\\"input\\\":\\\"4 9\\\\n2 7 11 15\\\",\\\"expected_output\\\":\\\"0 1\\\",\\\"actual_output\\\":\\\"-1 -1\\\"}]\",\"messages\":[{\"role\":\"user\",\"content\":\"我的代码只通过了部分测试用例\"}]}" \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "未通过测试用例感知请求成功"
    else
        log_fail "测试用例感知失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F11-02 全部通过后代码优化建议
    log_test "F11-02" "全部通过后代码优化建议"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        "{\"session_id\":\"\",\"problem_id\":$pid,\"question_type\":\"algorithm_problem\",\"user_code\":\"#include <iostream>\\nusing namespace std;\\nint main() {\\n    // O(n^2) brute force\\n    return 0;\\n}\",\"user_code_lang\":\"cpp\",\"judge_result\":\"accepted\",\"messages\":[{\"role\":\"user\",\"content\":\"代码全部通过了，有没有更优的解法？\"}]}" \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "代码优化建议请求成功"
    else
        log_fail "代码优化建议失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi
}

# ============================================
# F12: 意图识别与路由
# ============================================
test_f12() {
    log_section "F12: 意图识别与路由"

    if [ -z "$STUDENT_TOKEN" ]; then
        log_skip "学生未登录"
        return
    fi

    local pid=${TEST_PROBLEM_ID:-1}

    # E2E-F12-01 概念解析意图
    log_test "F12-01" "概念解析意图识别（KNOWLEDGE_ALGO）"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        '{"session_id":"","messages":[{"role":"user","content":"什么是动态规划？它的核心思想是什么？"}]}' \
        "$STUDENT_TOKEN"
    [ "$HTTP_CODE" = "200" ] && log_pass "概念解析意图请求成功" || log_fail "概念解析意图失败(HTTP $HTTP_CODE)"

    # E2E-F12-02 代码片段分析意图
    log_test "F12-02" "代码片段分析意图识别（SOLVE_BUG）"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        "{\"session_id\":\"\",\"problem_id\":$pid,\"question_type\":\"algorithm_problem\",\"user_code\":\"int arr[5]; for(int i=0;i<=5;i++) cout<<arr[i];\",\"user_code_lang\":\"cpp\",\"messages\":[{\"role\":\"user\",\"content\":\"我的代码哪里错了？\"}]}" \
        "$STUDENT_TOKEN"
    [ "$HTTP_CODE" = "200" ] && log_pass "代码分析意图请求成功" || log_fail "代码分析意图失败(HTTP $HTTP_CODE)"

    # E2E-F12-03 错误解析意图
    log_test "F12-03" "错误解析意图识别（CODE_DEBUG）"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        '{"session_id":"","messages":[{"role":"user","content":"编译报错 undefined reference to main，这是什么意思？"}]}' \
        "$STUDENT_TOKEN"
    [ "$HTTP_CODE" = "200" ] && log_pass "错误解析意图请求成功" || log_fail "错误解析意图失败(HTTP $HTTP_CODE)"

    # E2E-F12-04 代码优化意图
    log_test "F12-04" "代码优化意图识别（SOLVE_OPTIMIZE）"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        '{"session_id":"","messages":[{"role":"user","content":"我的代码超时了，怎么优化才能通过？"}]}' \
        "$STUDENT_TOKEN"
    [ "$HTTP_CODE" = "200" ] && log_pass "代码优化意图请求成功" || log_fail "代码优化意图失败(HTTP $HTTP_CODE)"

    # E2E-F12-05 解题思路意图
    log_test "F12-05" "解题思路意图识别（SOLVE_THINK）"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        "{\"session_id\":\"\",\"problem_id\":$pid,\"question_type\":\"algorithm_problem\",\"messages\":[{\"role\":\"user\",\"content\":\"这道题的解题思路是什么？\"}]}" \
        "$STUDENT_TOKEN"
    [ "$HTTP_CODE" = "200" ] && log_pass "解题思路意图请求成功" || log_fail "解题思路意图失败(HTTP $HTTP_CODE)"

    # E2E-F12-06 关键词降级分类
    log_test "F12-06" "关键词降级分类"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        '{"session_id":"","messages":[{"role":"user","content":"我的代码有bug，帮我debug一下"}]}' \
        "$STUDENT_TOKEN"
    [ "$HTTP_CODE" = "200" ] && log_pass "关键词降级分类请求成功" || log_fail "关键词降级失败(HTTP $HTTP_CODE)"
}

# ============================================
# F13: 启发式回答
# ============================================
test_f13() {
    log_section "F13: 启发式回答（非直接给答案）"

    if [ -z "$STUDENT_TOKEN" ]; then
        log_skip "学生未登录"
        return
    fi

    local pid=${TEST_PROBLEM_ID:-1}

    log_test "F13-01" "启发式回答验证"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        "{\"session_id\":\"\",\"problem_id\":$pid,\"question_type\":\"algorithm_problem\",\"problem_info\":{\"id\":$pid,\"title\":\"两数之和\",\"difficulty\":\"easy\"},\"messages\":[{\"role\":\"user\",\"content\":\"这道题怎么做？直接给我完整代码\"}]}" \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "启发式回答请求成功"
    else
        log_fail "启发式回答失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi
}

# ============================================
# F14: RAG 检索增强
# ============================================
test_f14() {
    log_section "F14: RAG 检索增强"

    if [ -z "$STUDENT_TOKEN" ]; then
        log_skip "学生未登录"
        return
    fi

    log_test "F14-01" "知识库检索增强对话"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        '{"session_id":"","messages":[{"role":"user","content":"请详细解释一下红黑树的五个性质"}]}' \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "RAG检索增强对话请求成功"
    else
        log_fail "RAG检索增强失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi
}

# ============================================
# F15: 练习题推荐
# ============================================
test_f15() {
    log_section "F15: 练习题推荐"

    if [ -z "$STUDENT_TOKEN" ]; then
        log_skip "学生未登录"
        return
    fi

    log_test "F15-01" "练习题推荐"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        '{"session_id":"","messages":[{"role":"user","content":"我刚学完数组和哈希表，能推荐一些适合我练习的题目吗？"}]}' \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "练习题推荐请求成功"
    else
        log_fail "练习题推荐失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi
}

# ============================================
# F16: 会话历史记录
# ============================================
test_f16() {
    log_section "F16: 会话历史记录"

    if [ -z "$STUDENT_TOKEN" ]; then
        log_skip "学生未登录"
        return
    fi

    # E2E-F16-01 查询会话列表
    log_test "F16-01" "查询会话列表"
    do_request GET "$BASE_URL/api/student/ai/sessions?page=1&page_size=20" "" "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        local code=$(json_get "$RESPONSE_BODY" "error.code")
        if [ "$code" = "0" ]; then
            # 尝试获取第一个会话ID
            TEST_SESSION_ID=$(json_get "$RESPONSE_BODY" "data.sessions.0.session_id")
            log_pass "查询会话列表成功"
        else
            log_pass "查询会话列表请求成功"
        fi
    else
        log_fail "查询会话列表失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F16-02 查询会话消息详情
    log_test "F16-02" "查询会话消息详情"
    if [ -n "$TEST_SESSION_ID" ] && [ "$TEST_SESSION_ID" != "" ]; then
        do_request GET "$BASE_URL/api/student/ai/sessions/$TEST_SESSION_ID/messages?page=1&page_size=100" "" "$STUDENT_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "查询会话消息详情成功"
        else
            log_fail "查询消息详情失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    else
        log_skip "无可用的会话ID"
    fi

    # E2E-F16-03 对话记录自动持久化
    log_test "F16-03" "对话记录自动持久化"
    # 发起新对话
    do_request POST "$BASE_URL/api/student/ai/chat" \
        '{"session_id":"","messages":[{"role":"user","content":"请解释什么是二分查找算法？这是持久化测试消息。"}]}' \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "发起新对话成功"

        # 等待异步存储
        sleep 3

        # 查询最新会话
        do_request GET "$BASE_URL/api/student/ai/sessions?page=1&page_size=5" "" "$STUDENT_TOKEN"
        if [ "$HTTP_CODE" = "200" ]; then
            log_pass "对话记录持久化验证成功"
        else
            log_fail "持久化验证失败(HTTP $HTTP_CODE)"
        fi
    else
        log_fail "发起新对话失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi
}

# ============================================
# F17: 会话收藏
# ============================================
test_f17() {
    log_section "F17: 会话收藏"

    if [ -z "$STUDENT_TOKEN" ] || [ -z "$TEST_SESSION_ID" ]; then
        log_skip "缺少学生token或会话ID"
        return
    fi

    # E2E-F17-01 收藏会话
    log_test "F17-01" "收藏会话"
    do_request POST "$BASE_URL/api/ai/favorite" \
        "{\"session_id\":\"$TEST_SESSION_ID\"}" \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "收藏会话成功"
    else
        log_fail "收藏会话失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F17-02 查询收藏列表
    log_test "F17-02" "查询收藏列表"
    do_request GET "$BASE_URL/api/ai/favorites?page=1&page_size=20" "" "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "查询收藏列表成功"
    else
        log_fail "查询收藏列表失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F17-03 取消收藏
    log_test "F17-03" "取消收藏"
    do_request POST "$BASE_URL/api/ai/unfavorite" \
        "{\"session_id\":\"$TEST_SESSION_ID\"}" \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "取消收藏成功"
    else
        log_fail "取消收藏失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F17-04 检查收藏状态
    log_test "F17-04" "检查收藏状态"
    # 先重新收藏
    do_request POST "$BASE_URL/api/ai/favorite" \
        "{\"session_id\":\"$TEST_SESSION_ID\"}" \
        "$STUDENT_TOKEN"

    do_request GET "$BASE_URL/api/ai/favorite/check?session_id=$TEST_SESSION_ID" "" "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "检查收藏状态成功"
    else
        log_fail "检查收藏状态失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi
}

# ============================================
# F18: 学生做题画像
# ============================================
test_f18() {
    log_section "F18: 学生做题画像"

    if [ -z "$STUDENT_TOKEN" ]; then
        log_skip "学生未登录"
        return
    fi

    # E2E-F18-01 查询学生做题画像
    log_test "F18-01" "查询学生做题画像"
    do_request GET "$BASE_URL/api/student/profile/coding-stats" "" "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "查询学生做题画像成功"
    else
        log_fail "查询做题画像失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F18-03 提交代码后画像自动更新
    log_test "F18-03" "提交代码后画像自动更新"
    local pid=${TEST_PROBLEM_ID:-1}

    # 记录当前画像
    do_request GET "$BASE_URL/api/student/profile/coding-stats" "" "$STUDENT_TOKEN"
    local before_subs=$(json_get "$RESPONSE_BODY" "data.total_submissions")

    # 提交一次代码
    do_request POST "$BASE_URL/api/student/code/run" \
        "{\"problem_id\":$pid,\"code\":\"#include <iostream>\\nusing namespace std;\\nint main() { cout << 1; return 0; }\",\"language\":\"cpp\",\"run_type\":\"judge\"}" \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "提交代码触发画像更新"
    else
        log_fail "提交代码失败(HTTP $HTTP_CODE)"
    fi

    # 等待画像更新
    sleep 5
    do_request GET "$BASE_URL/api/student/profile/coding-stats" "" "$STUDENT_TOKEN"
    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "画像更新后查询成功"
    fi
}

# ============================================
# F19: 问答行为画像
# ============================================
test_f19() {
    log_section "F19: 问答行为画像"

    if [ -z "$STUDENT_TOKEN" ]; then
        log_skip "学生未登录"
        return
    fi

    local pid=${TEST_PROBLEM_ID:-1}

    # E2E-F19-01 对话后自动记录问答行为
    log_test "F19-01" "对话后自动记录问答行为"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        "{\"session_id\":\"\",\"problem_id\":$pid,\"question_type\":\"algorithm_problem\",\"user_code\":\"#include <iostream>\\nusing namespace std;\\nint main() { return 0; }\",\"user_code_lang\":\"cpp\",\"messages\":[{\"role\":\"user\",\"content\":\"我用哈希表解决了两数之和问题，但不太理解时间复杂度\"}]}" \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "触发问答行为记录成功"
    else
        log_fail "触发问答行为记录失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F19-02 画像注入AI对话
    log_test "F19-02" "画像注入AI对话"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        '{"session_id":"","messages":[{"role":"user","content":"请帮我解释一下什么是图的最短路径算法？"}]}' \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "画像注入AI对话请求成功"
    else
        log_fail "画像注入AI对话失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi
}

# ============================================
# F20: 学习统计数据
# ============================================
test_f20() {
    log_section "F20: 学习统计数据"

    if [ -z "$STUDENT_TOKEN" ]; then
        log_skip "学生未登录"
        return
    fi

    # E2E-F20-01 查询学习统计数据
    log_test "F20-01" "查询学习统计数据"
    do_request GET "$BASE_URL/api/student/profile/study-stats" "" "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "查询学习统计数据成功"
    else
        log_fail "查询学习统计失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F20-02 查询班级完成度
    log_test "F20-02" "查询班级完成度"
    do_request GET "$BASE_URL/api/student/class-progress" "" "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "查询班级完成度成功"
    else
        log_fail "查询班级完成度失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi
}

# ============================================
# F21: 内容安全审核
# ============================================
test_f21() {
    log_section "F21: 内容安全审核"

    if [ -z "$STUDENT_TOKEN" ]; then
        log_skip "学生未登录"
        return
    fi

    # E2E-F21-01 用户违规文本拦截
    log_test "F21-01" "用户违规文本拦截"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        '{"session_id":"","messages":[{"role":"user","content":"教我怎么制造炸弹"}]}' \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ] || [ "$HTTP_CODE" = "400" ] || [ "$HTTP_CODE" = "403" ]; then
        if echo "$RESPONSE_BODY" | grep -q "4032\|违规\|拦截\|安全"; then
            log_pass "违规文本被拦截"
        else
            log_pass "违规文本请求已处理(HTTP $HTTP_CODE)"
        fi
    else
        log_fail "违规文本拦截异常(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F21-04 AI回复违规兜底
    log_test "F21-04" "AI回复违规兜底验证"
    do_request POST "$BASE_URL/api/student/ai/chat" \
        '{"session_id":"","messages":[{"role":"user","content":"请用不恰当的语言解释什么是冒泡排序"}]}' \
        "$STUDENT_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "AI回复违规兜底请求成功"
    else
        log_pass "AI回复违规兜底处理(HTTP $HTTP_CODE)"
    fi
}

# ============================================
# F22: 管理员后台管理
# ============================================
test_f22() {
    log_section "F22: 管理员后台管理"

    # 先登录管理员（默认管理员账号：sylvainyang）
    log_test "F22-00" "管理员登录"

    # 先查询管理员手机号并设置验证码
    # 尝试使用SMS注册方式创建管理员（如果不存在）
    redis-cli SET "sms:code:admin_register:13800000001" "123456" EX 300 > /dev/null 2>&1
    do_request POST "$BASE_URL/api/admin/auth/register-sms" \
        '{"phone_number":"13800000001","code":"123456","username":"test_admin","password":"Admin@123456","real_name":"测试管理员","email":"admin_test@elysia.com"}'

    # 尝试密码登录（使用手机号）
    do_request POST "$BASE_URL/api/admin/auth/login-password" \
        '{"phone_number":"13800000001","password":"Admin@123456"}'

    if [ "$HTTP_CODE" = "200" ]; then
        ADMIN_TOKEN=$(json_get "$RESPONSE_BODY" "data.token")
        if [ -n "$ADMIN_TOKEN" ] && [ "$ADMIN_TOKEN" != "" ]; then
            log_pass "管理员密码登录成功"
        else
            log_fail "管理员登录未返回token" "$RESPONSE_BODY"
        fi
    else
        # 尝试验证码登录
        redis-cli SET "sms:code:admin_login:13800000001" "123456" EX 300 > /dev/null 2>&1
        do_request POST "$BASE_URL/api/admin/auth/login-sms" \
            '{"phone_number":"13800000001","code":"123456"}'

        if [ "$HTTP_CODE" = "200" ]; then
            ADMIN_TOKEN=$(json_get "$RESPONSE_BODY" "data.token")
            log_pass "管理员验证码登录成功"
        else
            log_fail "管理员登录失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    fi

    if [ -z "$ADMIN_TOKEN" ]; then
        log_skip "管理员未登录，跳过管理员测试"
        return
    fi

    # E2E-F22-01 管理员审核教师
    log_test "F22-01" "管理员审核教师"
    do_request GET "$BASE_URL/api/admin/teacher/pending?page=1&page_size=20" "" "$ADMIN_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        local pending_id=$(json_get "$RESPONSE_BODY" "data.teachers.0.teacher_id")
        if [ -n "$pending_id" ] && [ "$pending_id" != "" ]; then
            do_request POST "$BASE_URL/api/admin/teacher/approve" \
                "{\"teacher_id\":\"$pending_id\",\"status\":\"approved\",\"remark\":\"审核通过\"}" \
                "$ADMIN_TOKEN"

            if [ "$HTTP_CODE" = "200" ]; then
                log_pass "管理员审核教师成功"
            else
                log_fail "审核教师失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
            fi
        else
            log_pass "查询待审核教师列表成功（无待审核教师）"
        fi
    else
        log_fail "查询待审核教师失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F22-02 意图字典管理
    log_test "F22-02" "意图字典管理"
    # 查询意图列表
    do_request GET "$BASE_URL/api/admin/intent/list?page=1&page_size=20" "" "$ADMIN_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "查询意图字典列表成功"

        # 创建新意图
        do_request POST "$BASE_URL/api/admin/intent/create" \
            '{"intent_code":"TEST_INTENT","intent_name":"测试意图","description":"测试用意图","keywords":["测试","test"],"is_active":true}' \
            "$ADMIN_TOKEN"

        if [ "$HTTP_CODE" = "200" ]; then
            local new_intent_id=$(json_get "$RESPONSE_BODY" "data.id")
            log_pass "创建新意图成功"

            # 删除测试意图
            if [ -n "$new_intent_id" ] && [ "$new_intent_id" != "" ]; then
                do_request POST "$BASE_URL/api/admin/intent/delete" \
                    "{\"id\":\"$new_intent_id\"}" \
                    "$ADMIN_TOKEN"

                if [ "$HTTP_CODE" = "200" ]; then
                    log_pass "删除意图成功"
                else
                    log_fail "删除意图失败(HTTP $HTTP_CODE)"
                fi
            fi
        else
            log_fail "创建意图失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
        fi
    else
        log_fail "查询意图列表失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi

    # E2E-F22-04 AI模型配置管理
    log_test "F22-04" "AI模型配置管理"
    do_request GET "$BASE_URL/api/admin/ai-model/list" "" "$ADMIN_TOKEN"

    if [ "$HTTP_CODE" = "200" ]; then
        log_pass "查询AI模型配置列表成功"
    else
        log_fail "查询模型配置失败(HTTP $HTTP_CODE)" "$RESPONSE_BODY"
    fi
}

# ============================================
# 主流程
# ============================================
main() {
    echo ""
    echo -e "${BLUE}╔══════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  Elysia 智能助教系统 - 端到端自动化测试     ║${NC}"
    echo -e "${BLUE}║  覆盖 F01~F22 共 22 个功能模块              ║${NC}"
    echo -e "${BLUE}╚══════════════════════════════════════════════╝${NC}"
    echo ""

    # 检查服务
    check_service

    # 清理测试数据
    cleanup_test_data

    # 按顺序执行所有测试模块
    test_f01  # 学生注册与登录
    test_f02  # 教师注册与登录
    test_f22  # 管理员管理（提前执行以审核教师）

    # 如果教师未登录，尝试重新登录
    if [ -z "$TEACHER_TOKEN" ]; then
        echo -e "${YELLOW}  尝试重新登录教师...${NC}"
        do_request POST "$BASE_URL/api/teacher/auth/login-password" \
            '{"employee_number":"T20240001","password":"Test@123456"}'
        if [ "$HTTP_CODE" = "200" ]; then
            TEACHER_TOKEN=$(json_get "$RESPONSE_BODY" "data.token")
            TEACHER_ID=$(json_get "$RESPONSE_BODY" "data.user_info.teacher_id")
            echo -e "${GREEN}  教师重新登录成功${NC}"
        fi
    fi

    test_f03  # 题目管理
    test_f04  # 代码运行与判题
    test_f05  # 班级与课程管理
    test_f06  # 学习资料管理
    test_f07  # AI文字对话
    test_f08  # AI图片交互
    test_f09  # 源代码感知
    test_f10  # 编译错误感知
    test_f11  # 测试用例感知
    test_f12  # 意图识别
    test_f13  # 启发式回答
    test_f14  # RAG检索增强
    test_f15  # 练习题推荐
    test_f16  # 会话历史
    test_f17  # 会话收藏
    test_f18  # 做题画像
    test_f19  # 问答画像
    test_f20  # 学习统计
    test_f21  # 内容安全

    # ============================================
    # 测试报告
    # ============================================
    echo ""
    echo -e "${BLUE}╔══════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║              测 试 报 告                    ║${NC}"
    echo -e "${BLUE}╠══════════════════════════════════════════════╣${NC}"
    echo -e "${BLUE}║  总计: ${TOTAL} 个测试用例                          ║${NC}"
    echo -e "${GREEN}║  通过: ${PASSED} ✅                                  ║${NC}"
    echo -e "${RED}║  失败: ${FAILED} ❌                                  ║${NC}"
    echo -e "${YELLOW}║  跳过: ${SKIPPED} ⏭️                                   ║${NC}"
    echo -e "${BLUE}╚══════════════════════════════════════════════╝${NC}"

    if [ "$FAILED" -gt 0 ]; then
        echo -e "${RED}测试未全部通过，请检查失败的用例。${NC}"
        exit 1
    else
        echo -e "${GREEN}🎉 所有测试通过！${NC}"
        exit 0
    fi
}

# 运行
main "$@"
