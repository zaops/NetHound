package tester

import (
	"context"
	"fmt"
	"sync"
	"time"

	"nethound/internal/config"
	"nethound/internal/target"
	"nethound/internal/types"
)

// Tester 测试器
type Tester struct {
	config  *config.Config
	targets []types.Target
	results []types.TestResult
	mu      sync.RWMutex
}

// New 创建新的测试器
func New(cfg *config.Config) *Tester {
	return &Tester{
		config:  cfg,
		results: make([]types.TestResult, 0),
	}
}

// Run 执行测试
func (t *Tester) Run() ([]types.TestResult, error) {
	// 解析目标
	parser := target.NewParser()
	targets, err := parser.ParseTargets(t.config)
	if err != nil {
		return nil, fmt.Errorf("解析目标失败: %v", err)
	}

	t.targets = targets

	if !t.config.Quiet {
		fmt.Printf("开始测试 %d 个目标...\n", len(targets))
	}

	// 创建上下文
	ctx := context.Background()

	// 执行测试
	if err := t.runTests(ctx); err != nil {
		return nil, fmt.Errorf("执行测试失败: %v", err)
	}

	return t.results, nil
}

// runTests 执行所有测试
func (t *Tester) runTests(ctx context.Context) error {
	// 创建工作池
	semaphore := make(chan struct{}, t.config.Concurrency)
	var wg sync.WaitGroup

	// 速率限制
	var rateLimiter <-chan time.Time
	if t.config.RateLimit > 0 {
		rateLimiter = time.Tick(time.Second / time.Duration(t.config.RateLimit))
	}

	startTime := time.Now()

	for _, target := range t.targets {
		// 速率限制
		if rateLimiter != nil {
			<-rateLimiter
		}

		wg.Add(1)
		go func(tgt types.Target) {
			defer wg.Done()

			// 获取信号量
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			// 执行测试
			t.testTarget(ctx, tgt)
		}(target)
	}

	wg.Wait()

	if !t.config.Quiet {
		duration := time.Since(startTime)
		fmt.Printf("测试完成，耗时: %v\n", duration)
	}

	return nil
}

// testTarget 测试单个目标
func (t *Tester) testTarget(ctx context.Context, target types.Target) {
	// 根据测试类型执行不同的测试
	switch t.config.TestType {
	case "ping":
		t.pingTest(ctx, target)
	case "telnet":
		t.telnetTest(ctx, target)
	case "all":
		t.pingTest(ctx, target)
		t.telnetTest(ctx, target)
	}
}

// pingTest 执行ping测试
func (t *Tester) pingTest(ctx context.Context, target types.Target) {
	pinger := NewPinger(t.config)
	result := pinger.Ping(ctx, target)
	
	t.mu.Lock()
	t.results = append(t.results, result)
	t.mu.Unlock()

	if t.config.Verbose {
		fmt.Printf("PING %s: %s\n", target.Host, result.Status)
	}
}

// telnetTest 执行telnet测试
func (t *Tester) telnetTest(ctx context.Context, target types.Target) {
	telneter := NewTelneter(t.config)
	result := telneter.Test(ctx, target)
	
	t.mu.Lock()
	t.results = append(t.results, result)
	t.mu.Unlock()

	if t.config.Verbose {
		fmt.Printf("TELNET %s:%d: %s\n", target.Host, target.Port, result.Status)
	}
}