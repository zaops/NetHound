package types

import (
	"fmt"
	"time"
)

// TestType 测试类型
type TestType string

const (
	TestTypePing   TestType = "ping"
	TestTypeTelnet TestType = "telnet"
	TestTypeAll    TestType = "all"
)

// TestStatus 测试状态
type TestStatus string

const (
	StatusReachable         TestStatus = "REACHABLE"
	StatusTimeout           TestStatus = "TIMEOUT"
	StatusRefused           TestStatus = "REFUSED"
	StatusDNSError          TestStatus = "DNS_ERROR"
	StatusNetworkUnreachable TestStatus = "NETWORK_UNREACHABLE"
	StatusHostUnreachable   TestStatus = "HOST_UNREACHABLE"
	StatusError             TestStatus = "ERROR"
)

// ErrorType 错误类型
type ErrorType string

const (
	ErrorTimeout            ErrorType = "timeout"
	ErrorConnectionRefused  ErrorType = "connection_refused"
	ErrorNetworkUnreachable ErrorType = "network_unreachable"
	ErrorHostUnreachable    ErrorType = "host_unreachable"
	ErrorDNS                ErrorType = "dns_error"
	ErrorPermission         ErrorType = "permission_denied"
	ErrorInternal           ErrorType = "internal_error"
)

// Target 测试目标
type Target struct {
	Host        string   `json:"host"`
	Port        int      `json:"port"`
	Description string   `json:"description,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// TestResult 测试结果
type TestResult struct {
	Host      string        `json:"host"`
	Port      int           `json:"port"`
	Type      string        `json:"type"`
	Status    TestStatus    `json:"status"`
	RTT       time.Duration `json:"rtt_ms,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
	Error     *TestError    `json:"error,omitempty"`
	Tags      []string      `json:"tags,omitempty"`
}

// TestError 测试错误
type TestError struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Code    int       `json:"code,omitempty"`
}

// TestSummary 测试摘要
type TestSummary struct {
	Timestamp   time.Time     `json:"timestamp"`
	Total       int           `json:"total"`
	Success     int           `json:"success"`
	Failed      int           `json:"failed"`
	Duration    time.Duration `json:"duration_ms"`
	Results     []TestResult  `json:"results"`
}

// String 返回目标的字符串表示
func (t Target) String() string {
	if t.Port > 0 {
		return fmt.Sprintf("%s:%d", t.Host, t.Port)
	}
	return t.Host
}

// IsSuccess 判断测试是否成功
func (r TestResult) IsSuccess() bool {
	return r.Status == StatusReachable
}

// RTTMilliseconds 返回RTT的毫秒数
func (r TestResult) RTTMilliseconds() float64 {
	if r.RTT == 0 {
		return 0
	}
	return float64(r.RTT.Nanoseconds()) / 1e6
}