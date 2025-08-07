package tester

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"nethound/internal/config"
	"nethound/internal/types"
)

// Pinger ping测试器
type Pinger struct {
	config *config.Config
}

// NewPinger 创建新的ping测试器
func NewPinger(cfg *config.Config) *Pinger {
	return &Pinger{
		config: cfg,
	}
}

// Ping 执行ping测试
func (p *Pinger) Ping(ctx context.Context, target types.Target) types.TestResult {
	result := types.TestResult{
		Host:      target.Host,
		Port:      0, // ping不使用端口
		Type:      "PING",
		Timestamp: time.Now(),
		Tags:      target.Tags,
	}

	// 重试机制
	for attempt := 0; attempt <= p.config.Retry; attempt++ {
		if attempt > 0 {
			time.Sleep(time.Duration(p.config.RetryDelay) * time.Millisecond)
		}

		// 执行ping
		rtt, err := p.doPing(ctx, target.Host)
		if err == nil {
			result.Status = types.StatusReachable
			result.RTT = rtt
			return result
		}

		// 最后一次尝试失败，设置错误信息
		if attempt == p.config.Retry {
			result.Status, result.Error = p.parseError(err)
		}
	}

	return result
}

// doPing 执行实际的ping操作
func (p *Pinger) doPing(ctx context.Context, host string) (time.Duration, error) {
	// 首先尝试DNS解析
	if _, err := net.LookupHost(host); err != nil {
		return 0, fmt.Errorf("DNS解析失败: %v", err)
	}

	// 构建ping命令
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		cmd = exec.CommandContext(ctx, "ping", "-n", strconv.Itoa(p.config.PingCount), host)
	default:
		cmd = exec.CommandContext(ctx, "ping", "-c", strconv.Itoa(p.config.PingCount), host)
	}

	// 执行命令
	output, err := cmd.Output()
	outputStr := string(output)
	
	// 调试：输出ping命令的结果（仅在verbose模式下）
	if p.config.Verbose {
		fmt.Printf("DEBUG: ping命令输出:\n%s\n", outputStr)
		if err != nil {
			fmt.Printf("DEBUG: ping命令错误: %v\n", err)
		}
	}
	
	if err != nil {
		// 检查是否是因为主机不可达等原因导致的失败
		if exitError, ok := err.(*exec.ExitError); ok {
			stderrStr := string(exitError.Stderr)
			fullOutput := outputStr + stderrStr
			if strings.Contains(fullOutput, "请求超时") || strings.Contains(fullOutput, "Request timed out") ||
			   strings.Contains(fullOutput, "目标主机无法访问") || strings.Contains(fullOutput, "Destination host unreachable") {
				return 0, fmt.Errorf("ping超时或主机不可达")
			}
		}
		return 0, fmt.Errorf("ping命令执行失败: %v", err)
	}

	// 检查ping是否成功（主要看连通性，RTT是次要的）
	if p.isPingSuccessful(outputStr) {
		// 尝试解析RTT，如果解析失败也不影响连通性判断
		if rtt, err := p.parseRTT(outputStr); err == nil {
			return rtt, nil
		}
		// RTT解析失败，但ping成功，返回一个默认值
		return time.Millisecond, nil
	}

	return 0, fmt.Errorf("ping失败: 主机不可达")
}

// isPingSuccessful 判断ping是否成功（主要看连通性）
func (p *Pinger) isPingSuccessful(output string) bool {
	// Windows系统的成功标识
	if runtime.GOOS == "windows" {
		// 检查是否有成功的回复（中文）
		if strings.Contains(output, "来自") && strings.Contains(output, "的回复") {
			return true
		}
		// 检查是否有成功的回复（英文）
		if strings.Contains(output, "Reply from") {
			return true
		}
		// 检查统计信息中的成功率（中文）
		if strings.Contains(output, "丢失 = 0") || strings.Contains(output, "(0% 丢失)") {
			return true
		}
		// 检查统计信息中的成功率（英文）
		if strings.Contains(output, "lost = 0") || strings.Contains(output, "(0% loss)") {
			return true
		}
		// 检查是否包含时间信息（中文和英文）
		if strings.Contains(output, "时间=") || strings.Contains(output, "time=") {
			return true
		}
		// 处理编码问题：检查乱码中的关键模式
		if strings.Contains(output, "ʱ��=") && strings.Contains(output, "ms") {
			return true // 这是"时间="的乱码形式
		}
		if strings.Contains(output, "��ʧ = 0") || strings.Contains(output, "(0% ��ʧ)") {
			return true // 这是"丢失"的乱码形式
		}
		// 检查是否包含IP地址和TTL（这些通常不会乱码）
		if strings.Contains(output, "TTL=") && strings.Contains(output, "ms") {
			return true
		}
	} else {
		// Linux/Unix系统的成功标识
		if strings.Contains(output, "bytes from") {
			return true
		}
		// 检查统计信息
		if strings.Contains(output, "0% packet loss") {
			return true
		}
	}
	
	return false
}

// parseRTT 解析ping输出中的RTT时间
func (p *Pinger) parseRTT(output string) (time.Duration, error) {
	lines := strings.Split(output, "\n")
	
	var totalRTT time.Duration
	var count int

	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		// Windows格式: time=1ms 或 time<1ms
		if runtime.GOOS == "windows" {
			if strings.Contains(line, "time=") {
				parts := strings.Split(line, "time=")
				if len(parts) > 1 {
					timeStr := strings.Fields(parts[1])[0]
					timeStr = strings.TrimSuffix(timeStr, "ms")
					if rtt, err := strconv.ParseFloat(timeStr, 64); err == nil {
						totalRTT += time.Duration(rtt * float64(time.Millisecond))
						count++
					}
				}
			} else if strings.Contains(line, "time<") {
				totalRTT += time.Millisecond
				count++
			}
		} else {
			// Linux/Unix格式: time=1.234 ms
			if strings.Contains(line, "time=") {
				parts := strings.Split(line, "time=")
				if len(parts) > 1 {
					timeStr := strings.Fields(parts[1])[0]
					if rtt, err := strconv.ParseFloat(timeStr, 64); err == nil {
						totalRTT += time.Duration(rtt * float64(time.Millisecond))
						count++
					}
				}
			}
		}
	}

	if count == 0 {
		return 0, fmt.Errorf("未找到有效的RTT数据")
	}

	// 返回平均RTT
	return totalRTT / time.Duration(count), nil
}

// parseError 解析错误类型
func (p *Pinger) parseError(err error) (types.TestStatus, *types.TestError) {
	errMsg := err.Error()
	
	if strings.Contains(errMsg, "DNS解析失败") {
		return types.StatusDNSError, &types.TestError{
			Type:    types.ErrorDNS,
			Message: errMsg,
		}
	}

	if strings.Contains(errMsg, "timeout") || strings.Contains(errMsg, "超时") {
		return types.StatusTimeout, &types.TestError{
			Type:    types.ErrorTimeout,
			Message: errMsg,
		}
	}

	if strings.Contains(errMsg, "network unreachable") || strings.Contains(errMsg, "网络不可达") {
		return types.StatusNetworkUnreachable, &types.TestError{
			Type:    types.ErrorNetworkUnreachable,
			Message: errMsg,
		}
	}

	if strings.Contains(errMsg, "host unreachable") || strings.Contains(errMsg, "主机不可达") {
		return types.StatusHostUnreachable, &types.TestError{
			Type:    types.ErrorHostUnreachable,
			Message: errMsg,
		}
	}

	return types.StatusError, &types.TestError{
		Type:    types.ErrorInternal,
		Message: errMsg,
	}
}