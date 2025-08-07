package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"

	"nethound/internal/config"
	"nethound/internal/types"
)

// Outputter 输出器
type Outputter struct {
	config *config.Config
}

// New 创建新的输出器
func New(cfg *config.Config) *Outputter {
	return &Outputter{
		config: cfg,
	}
}

// Output 输出结果
func (o *Outputter) Output(results []types.TestResult) error {
	// 创建摘要
	summary := o.createSummary(results)

	switch o.config.Output {
	case "table":
		return o.outputTable(summary)
	case "json":
		return o.outputJSON(summary)
	case "csv":
		return o.outputCSV(summary)
	case "html":
		return o.outputHTML(summary)
	default:
		return fmt.Errorf("不支持的输出格式: %s", o.config.Output)
	}
}

// createSummary 创建测试摘要
func (o *Outputter) createSummary(results []types.TestResult) *types.TestSummary {
	summary := &types.TestSummary{
		Timestamp: time.Now(),
		Total:     len(results),
		Results:   results,
	}

	// 计算成功和失败数量
	for _, result := range results {
		if result.IsSuccess() {
			summary.Success++
		} else {
			summary.Failed++
		}
	}

	// 计算总耗时（从第一个到最后一个结果的时间差）
	if len(results) > 0 {
		var earliest, latest time.Time
		for i, result := range results {
			if i == 0 {
				earliest = result.Timestamp
				latest = result.Timestamp
			} else {
				if result.Timestamp.Before(earliest) {
					earliest = result.Timestamp
				}
				if result.Timestamp.After(latest) {
					latest = result.Timestamp
				}
			}
		}
		summary.Duration = latest.Sub(earliest)
	}

	return summary
}

// outputTable 输出表格格式
func (o *Outputter) outputTable(summary *types.TestSummary) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Host", "Port", "Type", "Status", "RTT(ms)", "Error"})
	table.SetBorder(true)
	table.SetRowLine(true)

	// 设置颜色（如果启用）
	if !o.config.NoColor {
		table.SetHeaderColor(
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor},
		)
	}

	for _, result := range summary.Results {
		portStr := "-"
		if result.Port > 0 {
			portStr = fmt.Sprintf("%d", result.Port)
		}

		rttStr := "-"
		if result.RTT > 0 {
			rttStr = fmt.Sprintf("%.3f", result.RTTMilliseconds())
		}

		errorStr := "-"
		if result.Error != nil {
			errorStr = result.Error.Message
			// 限制错误信息长度
			if len(errorStr) > 50 {
				errorStr = errorStr[:47] + "..."
			}
		}

		row := []string{
			result.Host,
			portStr,
			result.Type,
			string(result.Status),
			rttStr,
			errorStr,
		}

		// 设置行颜色
		if !o.config.NoColor {
			var colors []tablewriter.Colors
			if result.IsSuccess() {
				colors = []tablewriter.Colors{
					{}, {}, {}, 
					{tablewriter.FgGreenColor},
					{}, {},
				}
			} else {
				colors = []tablewriter.Colors{
					{}, {}, {},
					{tablewriter.FgRedColor},
					{}, {tablewriter.FgRedColor},
				}
			}
			table.Rich(row, colors)
		} else {
			table.Append(row)
		}
	}

	table.Render()

	// 输出摘要
	if !o.config.Quiet {
		fmt.Printf("\n测试摘要:\n")
		fmt.Printf("总计: %d, 成功: %d, 失败: %d, 耗时: %v\n", 
			summary.Total, summary.Success, summary.Failed, summary.Duration)
	}

	return nil
}

// outputJSON 输出JSON格式
func (o *Outputter) outputJSON(summary *types.TestSummary) error {
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化JSON失败: %v", err)
	}

	return o.writeOutput(data)
}

// outputCSV 输出CSV格式
func (o *Outputter) outputCSV(summary *types.TestSummary) error {
	var output strings.Builder
	writer := csv.NewWriter(&output)

	// 写入标题行
	headers := []string{"Host", "Port", "Type", "Status", "RTT(ms)", "Error", "Timestamp", "Tags"}
	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("写入CSV标题失败: %v", err)
	}

	// 写入数据行
	for _, result := range summary.Results {
		portStr := ""
		if result.Port > 0 {
			portStr = fmt.Sprintf("%d", result.Port)
		}

		rttStr := ""
		if result.RTT > 0 {
			rttStr = fmt.Sprintf("%.3f", result.RTTMilliseconds())
		}

		errorStr := ""
		if result.Error != nil {
			errorStr = result.Error.Message
		}

		tagsStr := strings.Join(result.Tags, ",")

		row := []string{
			result.Host,
			portStr,
			result.Type,
			string(result.Status),
			rttStr,
			errorStr,
			result.Timestamp.Format(time.RFC3339),
			tagsStr,
		}

		if err := writer.Write(row); err != nil {
			return fmt.Errorf("写入CSV数据失败: %v", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("CSV写入错误: %v", err)
	}

	return o.writeOutput([]byte(output.String()))
}

// outputHTML 输出HTML格式
func (o *Outputter) outputHTML(summary *types.TestSummary) error {
	tmpl := `<!DOCTYPE html>
<html>
<head>
    <title>Nethound Network Test Results</title>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; margin: 20px; }
        h1 { color: #333; }
        .summary { background-color: #f5f5f5; padding: 15px; border-radius: 5px; margin-bottom: 20px; }
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; font-weight: bold; }
        .success { color: green; font-weight: bold; }
        .timeout { color: orange; font-weight: bold; }
        .error { color: red; font-weight: bold; }
        .refused { color: red; font-weight: bold; }
        .dns-error { color: purple; font-weight: bold; }
        .unreachable { color: darkred; font-weight: bold; }
        tr:nth-child(even) { background-color: #f9f9f9; }
        tr:hover { background-color: #f5f5f5; }
    </style>
</head>
<body>
    <h1>Nethound Network Test Results</h1>
    
    <div class="summary">
        <h2>测试摘要</h2>
        <p><strong>测试时间:</strong> {{.Timestamp.Format "2006-01-02 15:04:05"}}</p>
        <p><strong>总计:</strong> {{.Total}} | <strong>成功:</strong> {{.Success}} | <strong>失败:</strong> {{.Failed}}</p>
        <p><strong>耗时:</strong> {{.Duration}}</p>
    </div>

    <table>
        <thead>
            <tr>
                <th>Host</th>
                <th>Port</th>
                <th>Type</th>
                <th>Status</th>
                <th>RTT(ms)</th>
                <th>Error</th>
                <th>Timestamp</th>
                <th>Tags</th>
            </tr>
        </thead>
        <tbody>
            {{range .Results}}
            <tr>
                <td>{{.Host}}</td>
                <td>{{if gt .Port 0}}{{.Port}}{{else}}-{{end}}</td>
                <td>{{.Type}}</td>
                <td class="{{.Status | statusClass}}">{{.Status}}</td>
                <td>{{if gt .RTT 0}}{{printf "%.3f" .RTTMilliseconds}}{{else}}-{{end}}</td>
                <td>{{if .Error}}{{.Error.Message}}{{else}}-{{end}}</td>
                <td>{{.Timestamp.Format "15:04:05"}}</td>
                <td>{{range $i, $tag := .Tags}}{{if $i}}, {{end}}{{$tag}}{{end}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</body>
</html>`

	// 创建模板函数
	funcMap := template.FuncMap{
		"statusClass": func(status types.TestStatus) string {
			switch status {
			case types.StatusReachable:
				return "success"
			case types.StatusTimeout:
				return "timeout"
			case types.StatusRefused:
				return "refused"
			case types.StatusDNSError:
				return "dns-error"
			case types.StatusNetworkUnreachable, types.StatusHostUnreachable:
				return "unreachable"
			default:
				return "error"
			}
		},
	}

	t, err := template.New("html").Funcs(funcMap).Parse(tmpl)
	if err != nil {
		return fmt.Errorf("解析HTML模板失败: %v", err)
	}

	var output strings.Builder
	if err := t.Execute(&output, summary); err != nil {
		return fmt.Errorf("执行HTML模板失败: %v", err)
	}

	return o.writeOutput([]byte(output.String()))
}

// writeOutput 写入输出
func (o *Outputter) writeOutput(data []byte) error {
	outputFile := o.config.OutputFile
	
	// 如果没有指定输出文件，根据输出格式自动生成文件名
	if outputFile == "" {
		timestamp := time.Now().Format("20060102_150405")
		switch o.config.Output {
		case "json":
			outputFile = fmt.Sprintf("nethound_result_%s.json", timestamp)
		case "csv":
			outputFile = fmt.Sprintf("nethound_result_%s.csv", timestamp)
		case "html":
			outputFile = fmt.Sprintf("nethound_result_%s.html", timestamp)
		default:
			// 对于 table 格式，仍然输出到控制台
			fmt.Print(string(data))
			return nil
		}
	}

	// 确保输出目录存在
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(outputFile, data, 0644); err != nil {
		return fmt.Errorf("写入输出文件失败: %v", err)
	}

	if !o.config.Quiet {
		fmt.Printf("结果已保存到: %s\n", outputFile)
	}

	return nil
}
