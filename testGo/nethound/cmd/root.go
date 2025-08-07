package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"nethound/internal/config"
	"nethound/internal/output"
	"nethound/internal/tester"
)

var (
	cfgFile string
	cfg     *config.Config
)

var rootCmd = &cobra.Command{
	Use:   "nethound",
	Short: "Nethound - 网络连通性测试工具",
	Long: `Nethound 是一个轻量级网络连通性测试工具，支持多目标网络测试，
能够同时对多个IP地址和端口进行ping和telnet连通性测试。`,
	RunE: runTest,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// 全局配置参数
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "配置文件路径")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "详细输出模式")
	rootCmd.PersistentFlags().BoolP("quiet", "q", false, "静默模式")

	// 测试目标参数
	rootCmd.Flags().StringSliceP("targets", "t", []string{}, "目标列表，格式: ip:port")
	rootCmd.Flags().StringP("target-file", "f", "", "目标配置文件路径")
	rootCmd.Flags().StringP("ip-range", "r", "", "IP地址范围")
	rootCmd.Flags().StringP("port-range", "p", "", "端口范围")
	rootCmd.Flags().StringP("excel-file", "e", "", "Excel文件路径")
	rootCmd.Flags().String("excel-sheet", "Sheet1", "Excel工作表名称")

	// 测试类型参数
	rootCmd.Flags().StringP("test-type", "T", "all", "测试类型: ping, telnet, all")
	rootCmd.Flags().IntP("ping-count", "c", 3, "ping测试次数")
	rootCmd.Flags().IntP("ping-interval", "i", 1, "ping间隔时间(秒)")
	rootCmd.Flags().IntP("telnet-timeout", "w", 5000, "telnet连接超时(毫秒)")
	rootCmd.Flags().IntP("read-timeout", "d", 2000, "读取超时(毫秒)")

	// 输出参数
	rootCmd.Flags().StringP("output", "o", "table", "输出格式: table, json, csv, html")
	rootCmd.Flags().StringP("output-file", "F", "", "输出文件路径")
	rootCmd.Flags().Bool("no-color", false, "禁用彩色输出")

	// 性能参数
	rootCmd.Flags().IntP("concurrency", "C", 10, "并发测试数量")
	rootCmd.Flags().IntP("rate-limit", "R", 0, "速率限制(次/秒)")
	rootCmd.Flags().IntP("retry", "y", 1, "失败重试次数")
	rootCmd.Flags().IntP("retry-delay", "D", 1000, "重试间隔(毫秒)")

	// 其他参数
	rootCmd.Flags().Bool("dry-run", false, "预运行模式")
	rootCmd.Flags().Bool("save-config", false, "保存当前配置")

	// 绑定环境变量
	viper.BindEnv("targets", "NETHOUND_TARGETS")
	viper.BindEnv("output", "NETHOUND_OUTPUT_FORMAT")
	viper.BindEnv("concurrency", "NETHOUND_MAX_CONCURRENCY")
	viper.BindEnv("telnet-timeout", "NETHOUND_TIMEOUT")
	viper.BindEnv("retry", "NETHOUND_RETRY_COUNT")
	viper.BindEnv("retry-delay", "NETHOUND_RETRY_DELAY")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.SetConfigName("nethound")
		viper.SetConfigType("json")
	}

	viper.AutomaticEnv()
	viper.SetEnvPrefix("NETHOUND")

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("使用配置文件:", viper.ConfigFileUsed())
	}
}

func runTest(cmd *cobra.Command, args []string) error {
	// 创建默认配置
	cfg = config.NewConfig()
	if err := viper.Unmarshal(cfg); err != nil {
		return fmt.Errorf("解析配置失败: %v", err)
	}

	// 从命令行参数覆盖配置
	if err := parseFlags(cmd, cfg); err != nil {
		return fmt.Errorf("解析参数失败: %v", err)
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		return fmt.Errorf("配置验证失败: %v", err)
	}

	// 保存配置（如果需要）
	if cfg.SaveConfig {
		if err := cfg.Save(); err != nil {
			return fmt.Errorf("保存配置失败: %v", err)
		}
		fmt.Println("配置已保存")
		return nil
	}

	// 预运行模式
	if cfg.DryRun {
		fmt.Println("预运行模式 - 配置验证通过")
		return nil
	}

	// 创建测试器
	testRunner := tester.New(cfg)

	// 执行测试
	results, err := testRunner.Run()
	if err != nil {
		return fmt.Errorf("测试执行失败: %v", err)
	}

	// 输出结果
	outputter := output.New(cfg)
	if err := outputter.Output(results); err != nil {
		return fmt.Errorf("输出结果失败: %v", err)
	}

	return nil
}

func parseFlags(cmd *cobra.Command, cfg *config.Config) error {
	// 解析目标参数
	if targets, _ := cmd.Flags().GetStringSlice("targets"); len(targets) > 0 {
		cfg.Targets = targets
	}

	if targetFile, _ := cmd.Flags().GetString("target-file"); targetFile != "" {
		cfg.TargetFile = targetFile
	}

	if ipRange, _ := cmd.Flags().GetString("ip-range"); ipRange != "" {
		cfg.IPRange = ipRange
	}

	if portRange, _ := cmd.Flags().GetString("port-range"); portRange != "" {
		cfg.PortRange = portRange
	}

	if excelFile, _ := cmd.Flags().GetString("excel-file"); excelFile != "" {
		cfg.ExcelFile = excelFile
	}

	if excelSheet, _ := cmd.Flags().GetString("excel-sheet"); excelSheet != "" {
		cfg.ExcelSheet = excelSheet
	}

	// 解析测试类型参数
	if testType, _ := cmd.Flags().GetString("test-type"); testType != "" {
		cfg.TestType = testType
	}

	if pingCount, _ := cmd.Flags().GetInt("ping-count"); cmd.Flags().Changed("ping-count") {
		cfg.PingCount = pingCount
	}

	if pingInterval, _ := cmd.Flags().GetInt("ping-interval"); cmd.Flags().Changed("ping-interval") {
		cfg.PingInterval = pingInterval
	}

	if telnetTimeout, _ := cmd.Flags().GetInt("telnet-timeout"); cmd.Flags().Changed("telnet-timeout") {
		cfg.TelnetTimeout = telnetTimeout
	}

	if readTimeout, _ := cmd.Flags().GetInt("read-timeout"); cmd.Flags().Changed("read-timeout") {
		cfg.ReadTimeout = readTimeout
	}

	// 解析输出参数
	if output, _ := cmd.Flags().GetString("output"); output != "" {
		cfg.Output = output
	}

	if outputFile, _ := cmd.Flags().GetString("output-file"); outputFile != "" {
		cfg.OutputFile = outputFile
	}

	if noColor, _ := cmd.Flags().GetBool("no-color"); noColor {
		cfg.NoColor = noColor
	}

	// 解析性能参数
	if concurrency, _ := cmd.Flags().GetInt("concurrency"); cmd.Flags().Changed("concurrency") {
		cfg.Concurrency = concurrency
	}

	if rateLimit, _ := cmd.Flags().GetInt("rate-limit"); cmd.Flags().Changed("rate-limit") {
		cfg.RateLimit = rateLimit
	}

	if retry, _ := cmd.Flags().GetInt("retry"); cmd.Flags().Changed("retry") {
		cfg.Retry = retry
	}

	if retryDelay, _ := cmd.Flags().GetInt("retry-delay"); cmd.Flags().Changed("retry-delay") {
		cfg.RetryDelay = retryDelay
	}

	// 解析其他参数
	if verbose, _ := cmd.Flags().GetBool("verbose"); verbose {
		cfg.Verbose = verbose
	}

	if quiet, _ := cmd.Flags().GetBool("quiet"); quiet {
		cfg.Quiet = quiet
	}

	if dryRun, _ := cmd.Flags().GetBool("dry-run"); dryRun {
		cfg.DryRun = dryRun
	}

	if saveConfig, _ := cmd.Flags().GetBool("save-config"); saveConfig {
		cfg.SaveConfig = saveConfig
	}

	return nil
}