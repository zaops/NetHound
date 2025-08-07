package tester

import (
	"context"
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"

	"nethound/internal/config"
	"nethound/internal/types"
)

// Telneter telnet测试器
type Telneter struct {
	config *config.Config
}

// NewTelneter 创建新的telnet测试器
func NewTelneter(cfg *config.Config) *Telneter {
	return &Telneter{
		config: cfg,
	}
}

// Test 执行telnet测试
func (t *Telneter) Test(ctx context.Context, target types.Target) types.TestResult {
	result := types.TestResult{
		Host:      target.Host,
		Port:      target.Port,
		Type:      "TELNET",
		Timestamp: time.Now(),
		Tags:      target.Tags,
	}

	// 重试机制
	for attempt := 0; attempt <= t.config.Retry; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(t.config.RetryDelay) * time.Millisecond)
		}

		// 执行telnet测试
		rtt, err := t.doTelnet(ctx, target)
		if err == nil {
			result.Status = types.StatusReachable
			result.RTT = rtt
			return result
		}

		// 最后一次尝试失败，设置错误信息
		if attempt == t.config.Retry {
			result.Status, result.Error = t.parseError(err)
		}
	}

	return result
}

// doTelnet 执行实际的telnet测试
func (t *Telneter) doTelnet(ctx context.Context, target types.Target) (time.Duration, error) {
	// 创建带超时的上下文
	timeout := time.Duration(t.config.TelnetTimeout) * time.Millisecond
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// 记录开始时间
	startTime := time.Now()

	// 创建dialer
	dialer := &net.Dialer{
		Timeout: timeout,
	}

	// 尝试连接
	address := fmt.Sprintf("%s:%d", target.Host, target.Port)
	conn, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	// 计算RTT
	rtt := time.Since(startTime)

	// 可选：尝试读取一些数据以确保连接完全建立
	if t.config.ReadTimeout > 0 {
		conn.SetReadDeadline(time.Now().Add(time.Duration(t.config.ReadTimeout) * time.Millisecond))
		buffer := make([]byte, 1024)
		_, err := conn.Read(buffer)
		// 忽略读取错误，因为很多服务不会立即发送数据
		_ = err
	}

	return rtt, nil
}

// parseError 解析错误类型
func (t *Telneter) parseError(err error) (types.TestStatus, *types.TestError) {
	errMsg := err.Error()

	// DNS解析错误
	if strings.Contains(errMsg, "no such host") || 
	   strings.Contains(errMsg, "dns") ||
	   strings.Contains(errMsg, "name resolution") {
		return types.StatusDNSError, &types.TestError{
			Type:    types.ErrorDNS,
			Message: errMsg,
		}
	}

	// 连接被拒绝
	if strings.Contains(errMsg, "connection refused") ||
	   strings.Contains(errMsg, "connect: connection refused") {
		return types.StatusRefused, &types.TestError{
			Type:    types.ErrorConnectionRefused,
			Message: errMsg,
		}
	}

	// 超时错误
	if strings.Contains(errMsg, "timeout") ||
	   strings.Contains(errMsg, "deadline exceeded") ||
	   strings.Contains(errMsg, "i/o timeout") {
		return types.StatusTimeout, &types.TestError{
			Type:    types.ErrorTimeout,
			Message: errMsg,
		}
	}

	// 网络不可达
	if strings.Contains(errMsg, "network unreachable") ||
	   strings.Contains(errMsg, "no route to host") {
		return types.StatusNetworkUnreachable, &types.TestError{
			Type:    types.ErrorNetworkUnreachable,
			Message: errMsg,
		}
	}

	// 主机不可达
	if strings.Contains(errMsg, "host unreachable") {
		return types.StatusHostUnreachable, &types.TestError{
			Type:    types.ErrorHostUnreachable,
			Message: errMsg,
		}
	}

	// 权限错误
	if strings.Contains(errMsg, "permission denied") {
		return types.StatusError, &types.TestError{
			Type:    types.ErrorPermission,
			Message: errMsg,
		}
	}

	// 检查系统错误码
	if opErr, ok := err.(*net.OpError); ok {
		if sysErr, ok := opErr.Err.(*net.DNSError); ok {
			return types.StatusDNSError, &types.TestError{
				Type:    types.ErrorDNS,
				Message: sysErr.Error(),
			}
		}

		if sysErr, ok := opErr.Err.(syscall.Errno); ok {
			switch sysErr {
			case syscall.ECONNREFUSED:
				return types.StatusRefused, &types.TestError{
					Type:    types.ErrorConnectionRefused,
					Message: "Connection refused",
				}
			case syscall.ETIMEDOUT:
				return types.StatusTimeout, &types.TestError{
					Type:    types.ErrorTimeout,
					Message: "Connection timeout",
				}
			case syscall.ENETUNREACH:
				return types.StatusNetworkUnreachable, &types.TestError{
					Type:    types.ErrorNetworkUnreachable,
					Message: "Network unreachable",
				}
			case syscall.EHOSTUNREACH:
				return types.StatusHostUnreachable, &types.TestError{
					Type:    types.ErrorHostUnreachable,
					Message: "Host unreachable",
				}
			}
		}
	}

	// 默认错误
	return types.StatusError, &types.TestError{
		Type:    types.ErrorInternal,
		Message: errMsg,
	}
}