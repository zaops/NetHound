package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config 配置结构
type Config struct {
	// 测试目标参数
	Targets    []string `json:"targets,omitempty" mapstructure:"targets"`
	TargetFile string   `json:"target_file,omitempty" mapstructure:"target-file"`
	IPRange    string   `json:"ip_range,omitempty" mapstructure:"ip-range"`
	PortRange  string   `json:"port_range,omitempty" mapstructure:"port-range"`
	ExcelFile  string   `json:"excel_file,omitempty" mapstructure:"excel-file"`
	ExcelSheet string   `json:"excel_sheet,omitempty" mapstructure:"excel-sheet"`

	// 测试类型参数
	TestType       string `json:"test_type" mapstructure:"test-type"`
	PingCount      int    `json:"ping_count" mapstructure:"ping-count"`
	PingInterval   int    `json:"ping_interval" mapstructure:"ping-interval"`
	TelnetTimeout  int    `json:"telnet_timeout" mapstructure:"telnet-timeout"`
	ReadTimeout    int    `json:"read_timeout" mapstructure:"read-timeout"`

	// 输出参数
	Output     string `json:"output" mapstructure:"output"`
	OutputFile string `json:"output_file,omitempty" mapstructure:"output-file"`
	Verbose    bool   `json:"verbose,omitempty" mapstructure:"verbose"`
	Quiet      bool   `json:"quiet,omitempty" mapstructure:"quiet"`
	NoColor    bool   `json:"no_color,omitempty" mapstructure:"no-color"`

	// 性能参数
	Concurrency int `json:"concurrency" mapstructure:"concurrency"`
	RateLimit   int `json:"rate_limit,omitempty" mapstructure:"rate-limit"`
	Retry       int `json:"retry" mapstructure:"retry"`
	RetryDelay  int `json:"retry_delay" mapstructure:"retry-delay"`

	// 其他参数
	DryRun     bool `json:"dry_run,omitempty" mapstructure:"dry-run"`
	SaveConfig bool `json:"save_config,omitempty" mapstructure:"save-config"`
}

// NewConfig 创建默认配置
func NewConfig() *Config {
	return &Config{
		TestType:      "all",
		PingCount:     1,
		PingInterval:  1,
		TelnetTimeout: 3000,
		ReadTimeout:   1000,
		Output:        "table",
		Concurrency:   20,
		Retry:         0, 
		RetryDelay:    500,
		ExcelSheet:    "Sheet1",
	}
}

// Validate 验证配置
func (c *Config) Validate() error {
	// 验证测试类型
	validTestTypes := []string{"ping", "telnet", "all"}
	if !contains(validTestTypes, c.TestType) {
		return fmt.Errorf("无效的测试类型: %s, 支持的类型: %s", c.TestType, strings.Join(validTestTypes, ", "))
	}

	// 验证输出格式
	validOutputs := []string{"table", "json", "csv", "html"}
	if !contains(validOutputs, c.Output) {
		return fmt.Errorf("无效的输出格式: %s, 支持的格式: %s", c.Output, strings.Join(validOutputs, ", "))
	}

	// 验证数值参数
	if c.PingCount <= 0 {
		return fmt.Errorf("ping次数必须大于0")
	}

	if c.PingInterval <= 0 {
		return fmt.Errorf("ping间隔必须大于0")
	}

	if c.TelnetTimeout <= 0 {
		return fmt.Errorf("telnet超时时间必须大于0")
	}

	if c.ReadTimeout <= 0 {
		return fmt.Errorf("读取超时时间必须大于0")
	}

	if c.Concurrency <= 0 {
		return fmt.Errorf("并发数必须大于0")
	}

	if c.Retry < 0 {
		return fmt.Errorf("重试次数不能小于0")
	}

	if c.RetryDelay < 0 {
		return fmt.Errorf("重试间隔不能小于0")
	}

	// 验证目标配置
	if len(c.Targets) == 0 && c.TargetFile == "" && c.IPRange == "" && c.ExcelFile == "" {
		return fmt.Errorf("必须指定至少一个测试目标")
	}

	// 验证文件存在性
	if c.TargetFile != "" {
		if _, err := os.Stat(c.TargetFile); os.IsNotExist(err) {
			return fmt.Errorf("目标文件不存在: %s", c.TargetFile)
		}
	}

	if c.ExcelFile != "" {
		if _, err := os.Stat(c.ExcelFile); os.IsNotExist(err) {
			return fmt.Errorf("Excel文件不存在: %s", c.ExcelFile)
		}
	}

	return nil
}

// Save 保存配置到文件
func (c *Config) Save() error {
	configFile := "nethound.json"
	if c.OutputFile != "" {
		dir := filepath.Dir(c.OutputFile)
		configFile = filepath.Join(dir, "nethound.json")
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}

	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	fmt.Printf("配置已保存到: %s\n", configFile)
	return nil
}

// contains 检查字符串是否在切片中
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}