# Nethound API 文档

## 核心模块

### 1. Config 模块

配置管理模块，负责解析和验证配置参数。

#### 结构体

```go
type Config struct {
    // 测试目标参数
    Targets    []string `json:"targets,omitempty"`
    TargetFile string   `json:"target_file,omitempty"`
    IPRange    string   `json:"ip_range,omitempty"`
    PortRange  string   `json:"port_range,omitempty"`
    ExcelFile  string   `json:"excel_file,omitempty"`
    ExcelSheet string   `json:"excel_sheet,omitempty"`

    // 测试类型参数
    TestType       string `json:"test_type"`
    PingCount      int    `json:"ping_count"`
    PingInterval   int    `json:"ping_interval"`
    TelnetTimeout  int    `json:"telnet_timeout"`
    ReadTimeout    int    `json:"read_timeout"`

    // 输出参数
    Output     string `json:"output"`
    OutputFile string `json:"output_file,omitempty"`
    Verbose    bool   `json:"verbose,omitempty"`
    Quiet      bool   `json:"quiet,omitempty"`
    NoColor    bool   `json:"no_color,omitempty"`

    // 性能参数
    Concurrency int `json:"concurrency"`
    RateLimit   int `json:"rate_limit,omitempty"`
    Retry       int `json:"retry"`
    RetryDelay  int `json:"retry_delay"`

    // 其他参数
    DryRun     bool `json:"dry_run,omitempty"`
    SaveConfig bool `json:"save_config,omitempty"`
}
```

#### 方法

- `NewConfig() *Config`: 创建默认配置
- `Validate() error`: 验证配置参数
- `Save() error`: 保存配置到文件

### 2. Types 模块

定义核心数据类型。

#### 测试类型

```go
type TestType string

const (
    TestTypePing   TestType = "ping"
    TestTypeTelnet TestType = "telnet"
    TestTypeAll    TestType = "all"
)
```

#### 测试状态

```go
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
```

#### 测试结果

```go
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
```

### 3. Target 模块

目标解析模块，支持多种目标格式。

#### 支持的格式

- 单个目标: `host:port` 或 `host`
- 多个目标: `host1:port1,host2:port2`
- IP范围: `192.168.1.1-254:80`
- CIDR: `192.168.1.0/24:80`
- 端口范围: `host:80-90`

#### 方法

- `ParseTargets(cfg *config.Config) ([]types.Target, error)`: 解析所有目标
- `parseTarget(target string) ([]types.Target, error)`: 解析单个目标
- `parseRangeTargets(ipRange, portRange string) ([]types.Target, error)`: 解析范围目标

### 4. Tester 模块

测试执行模块，负责协调各种测试。

#### 主要方法

- `New(cfg *config.Config) *Tester`: 创建测试器
- `Run() ([]types.TestResult, error)`: 执行测试
- `runTests(ctx context.Context) error`: 执行所有测试
- `testTarget(ctx context.Context, target types.Target)`: 测试单个目标

#### Ping 测试器

```go
type Pinger struct {
    config *config.Config
}

// 方法
func NewPinger(cfg *config.Config) *Pinger
func (p *Pinger) Ping(ctx context.Context, target types.Target) types.TestResult
```

#### Telnet 测试器

```go
type Telneter struct {
    config *config.Config
}

// 方法
func NewTelneter(cfg *config.Config) *Telneter
func (t *Telneter) Test(ctx context.Context, target types.Target) types.TestResult
```

### 5. Output 模块

输出格式化模块，支持多种输出格式。

#### 支持的格式

- `table`: 表格格式（控制台输出）
- `json`: JSON格式
- `csv`: CSV格式
- `html`: HTML格式

#### 方法

- `New(cfg *config.Config) *Outputter`: 创建输出器
- `Output(results []types.TestResult) error`: 输出结果
- `outputTable(summary *types.TestSummary) error`: 表格输出
- `outputJSON(summary *types.TestSummary) error`: JSON输出
- `outputCSV(summary *types.TestSummary) error`: CSV输出
- `outputHTML(summary *types.TestSummary) error`: HTML输出

## 使用示例

### 基本用法

```go
package main

import (
    "fmt"
    "nethound/internal/config"
    "nethound/internal/tester"
    "nethound/internal/output"
)

func main() {
    // 创建配置
    cfg := config.NewConfig()
    cfg.Targets = []string{"baidu.com:80"}
    cfg.TestType = "all"
    cfg.Output = "json"
    
    // 验证配置
    if err := cfg.Validate(); err != nil {
        panic(err)
    }
    
    // 创建测试器
    testRunner := tester.New(cfg)
    
    // 执行测试
    results, err := testRunner.Run()
    if err != nil {
        panic(err)
    }
    
    // 输出结果
    outputter := output.New(cfg)
    if err := outputter.Output(results); err != nil {
        panic(err)
    }
}
```

### 自定义测试器

```go
// 创建自定义 ping 测试器
pinger := tester.NewPinger(cfg)
result := pinger.Ping(context.Background(), types.Target{
    Host: "baidu.com",
    Port: 0,
})

fmt.Printf("Ping result: %s\n", result.Status)
```

### 自定义输出

```go
// 创建测试摘要
summary := &types.TestSummary{
    Timestamp: time.Now(),
    Total:     len(results),
    Results:   results,
}

// 输出为JSON
outputter := output.New(cfg)
if err := outputter.outputJSON(summary); err != nil {
    panic(err)
}
```

## 错误处理

### 错误类型

```go
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
```

### 错误处理示例

```go
for _, result := range results {
    if !result.IsSuccess() {
        switch result.Error.Type {
        case types.ErrorTimeout:
            fmt.Printf("Host %s timeout\n", result.Host)
        case types.ErrorConnectionRefused:
            fmt.Printf("Host %s refused connection\n", result.Host)
        case types.ErrorDNS:
            fmt.Printf("DNS resolution failed for %s\n", result.Host)
        default:
            fmt.Printf("Unknown error for %s: %s\n", result.Host, result.Error.Message)
        }
    }
}
```

## 扩展开发

### 添加新的测试类型

1. 在 `types.go` 中定义新的测试类型
2. 在 `tester` 包中创建新的测试器
3. 在 `tester.go` 中添加调用逻辑
4. 更新配置验证逻辑

### 添加新的输出格式

1. 在 `output.go` 中添加新的输出方法
2. 在 `Output` 方法中添加格式判断
3. 更新配置验证逻辑

### 性能优化建议

1. 合理设置并发数 (`Concurrency`)
2. 根据网络环境调整超时时间
3. 使用速率限制避免网络拥塞
4. 对于大量目标，考虑分批处理