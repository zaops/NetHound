package target

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/tealeg/xlsx"

	"nethound/internal/config"
	"nethound/internal/types"
)

// Parser 目标解析器
type Parser struct{}

// NewParser 创建新的解析器
func NewParser() *Parser {
	return &Parser{}
}

// ParseTargets 解析所有目标
func (p *Parser) ParseTargets(cfg *config.Config) ([]types.Target, error) {
	var allTargets []types.Target

	// 解析命令行目标
	if len(cfg.Targets) > 0 {
		targets, err := p.parseStringTargets(cfg.Targets)
		if err != nil {
			return nil, fmt.Errorf("解析命令行目标失败: %v", err)
		}
		allTargets = append(allTargets, targets...)
	}

	// 解析IP范围和端口范围
	if cfg.IPRange != "" || cfg.PortRange != "" {
		targets, err := p.parseRangeTargets(cfg.IPRange, cfg.PortRange)
		if err != nil {
			return nil, fmt.Errorf("解析范围目标失败: %v", err)
		}
		allTargets = append(allTargets, targets...)
	}

	// 解析目标文件
	if cfg.TargetFile != "" {
		targets, err := p.parseTargetFile(cfg.TargetFile)
		if err != nil {
			return nil, fmt.Errorf("解析目标文件失败: %v", err)
		}
		allTargets = append(allTargets, targets...)
	}

	// 解析Excel文件
	if cfg.ExcelFile != "" {
		targets, err := p.parseExcelFile(cfg.ExcelFile, cfg.ExcelSheet)
		if err != nil {
			return nil, fmt.Errorf("解析Excel文件失败: %v", err)
		}
		allTargets = append(allTargets, targets...)
	}

	// 去重
	return p.deduplicateTargets(allTargets), nil
}

// parseStringTargets 解析字符串目标列表
func (p *Parser) parseStringTargets(targets []string) ([]types.Target, error) {
	var result []types.Target

	for _, target := range targets {
		parsed, err := p.parseTarget(target)
		if err != nil {
			return nil, fmt.Errorf("解析目标 '%s' 失败: %v", target, err)
		}
		result = append(result, parsed...)
	}

	return result, nil
}

// parseTarget 解析单个目标字符串
func (p *Parser) parseTarget(target string) ([]types.Target, error) {
	// 分离主机和端口
	parts := strings.Split(target, ":")
	
	var host string
	var portStr string
	
	if len(parts) == 1 {
		// 只有主机名，没有端口（用于ping测试）
		host = parts[0]
		portStr = "0" // 使用0表示不需要端口
	} else if len(parts) >= 2 {
		// 有主机名和端口
		host = parts[0]
		portStr = strings.Join(parts[1:], ":")
	} else {
		return nil, fmt.Errorf("目标格式错误")
	}

	// 解析主机
	hosts, err := p.parseHosts(host)
	if err != nil {
		return nil, fmt.Errorf("解析主机失败: %v", err)
	}

	// 解析端口
	ports, err := p.parsePorts(portStr)
	if err != nil {
		return nil, fmt.Errorf("解析端口失败: %v", err)
	}

	// 组合结果
	var result []types.Target
	for _, h := range hosts {
		for _, port := range ports {
			result = append(result, types.Target{
				Host: h,
				Port: port,
			})
		}
	}

	return result, nil
}

// parseHosts 解析主机列表
func (p *Parser) parseHosts(hostStr string) ([]string, error) {
	var hosts []string

	// 处理逗号分隔的多个主机
	for _, host := range strings.Split(hostStr, ",") {
		host = strings.TrimSpace(host)
		if host == "" {
			continue
		}

		// 检查是否为IP范围
		if strings.Contains(host, "-") || strings.Contains(host, "/") {
			expanded, err := p.expandIPRange(host)
			if err != nil {
				return nil, fmt.Errorf("展开IP范围 '%s' 失败: %v", host, err)
			}
			hosts = append(hosts, expanded...)
		} else {
			hosts = append(hosts, host)
		}
	}

	return hosts, nil
}

// expandIPRange 展开IP范围
func (p *Parser) expandIPRange(ipRange string) ([]string, error) {
	var ips []string

	// 处理CIDR格式
	if strings.Contains(ipRange, "/") {
		_, ipNet, err := net.ParseCIDR(ipRange)
		if err != nil {
			return nil, fmt.Errorf("解析CIDR失败: %v", err)
		}

		// 遍历CIDR范围内的所有IP
		for ip := ipNet.IP.Mask(ipNet.Mask); ipNet.Contains(ip); p.incrementIP(ip) {
			ips = append(ips, ip.String())
		}
		return ips, nil
	}

	// 处理范围格式 (192.168.1.1-254)
	if strings.Contains(ipRange, "-") {
		return p.expandIPRangeWithDash(ipRange)
	}

	return []string{ipRange}, nil
}

// expandIPRangeWithDash 展开带连字符的IP范围
func (p *Parser) expandIPRangeWithDash(ipRange string) ([]string, error) {
	parts := strings.Split(ipRange, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("IP范围格式错误")
	}

	startIP := strings.TrimSpace(parts[0])
	endPart := strings.TrimSpace(parts[1])

	// 解析起始IP
	startIPParts := strings.Split(startIP, ".")
	if len(startIPParts) != 4 {
		return nil, fmt.Errorf("起始IP格式错误")
	}

	var ips []string

	// 如果结束部分是完整IP
	if strings.Contains(endPart, ".") {
		endIPParts := strings.Split(endPart, ".")
		if len(endIPParts) != 4 {
			return nil, fmt.Errorf("结束IP格式错误")
		}

		// 简单处理：只支持最后一段的范围
		start, err := strconv.Atoi(startIPParts[3])
		if err != nil {
			return nil, fmt.Errorf("解析起始IP最后一段失败: %v", err)
		}

		end, err := strconv.Atoi(endIPParts[3])
		if err != nil {
			return nil, fmt.Errorf("解析结束IP最后一段失败: %v", err)
		}

		for i := start; i <= end; i++ {
			ip := fmt.Sprintf("%s.%s.%s.%d", startIPParts[0], startIPParts[1], startIPParts[2], i)
			ips = append(ips, ip)
		}
	} else {
		// 结束部分是数字
		start, err := strconv.Atoi(startIPParts[3])
		if err != nil {
			return nil, fmt.Errorf("解析起始IP最后一段失败: %v", err)
		}

		end, err := strconv.Atoi(endPart)
		if err != nil {
			return nil, fmt.Errorf("解析结束数字失败: %v", err)
		}

		for i := start; i <= end; i++ {
			ip := fmt.Sprintf("%s.%s.%s.%d", startIPParts[0], startIPParts[1], startIPParts[2], i)
			ips = append(ips, ip)
		}
	}

	return ips, nil
}

// parsePorts 解析端口列表
func (p *Parser) parsePorts(portStr string) ([]int, error) {
	var ports []int

	for _, portRange := range strings.Split(portStr, ",") {
		portRange = strings.TrimSpace(portRange)
		if portRange == "" {
			continue
		}

		if strings.Contains(portRange, "-") {
			// 端口范围
			parts := strings.Split(portRange, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("端口范围格式错误: %s", portRange)
			}

			start, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				return nil, fmt.Errorf("解析起始端口失败: %v", err)
			}

			end, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, fmt.Errorf("解析结束端口失败: %v", err)
			}

			for port := start; port <= end; port++ {
				ports = append(ports, port)
			}
		} else {
			// 单个端口
			port, err := strconv.Atoi(portRange)
			if err != nil {
				return nil, fmt.Errorf("解析端口失败: %v", err)
			}
			ports = append(ports, port)
		}
	}

	return ports, nil
}

// parseRangeTargets 解析范围目标
func (p *Parser) parseRangeTargets(ipRange, portRange string) ([]types.Target, error) {
	var targets []types.Target

	// 解析IP范围
	hosts, err := p.parseHosts(ipRange)
	if err != nil {
		return nil, fmt.Errorf("解析IP范围失败: %v", err)
	}

	// 解析端口范围
	ports, err := p.parsePorts(portRange)
	if err != nil {
		return nil, fmt.Errorf("解析端口范围失败: %v", err)
	}

	// 组合结果
	for _, host := range hosts {
		for _, port := range ports {
			targets = append(targets, types.Target{
				Host: host,
				Port: port,
			})
		}
	}

	return targets, nil
}

// parseTargetFile 解析目标文件
func (p *Parser) parseTargetFile(filename string) ([]types.Target, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 根据文件扩展名判断格式
	if strings.HasSuffix(strings.ToLower(filename), ".json") {
		return p.parseJSONFile(file)
	}

	// 默认按文本格式解析
	return p.parseTextFile(file)
}

// parseJSONFile 解析JSON文件
func (p *Parser) parseJSONFile(file *os.File) ([]types.Target, error) {
	var data struct {
		Targets []string `json:"targets"`
	}

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("解析JSON失败: %v", err)
	}

	return p.parseStringTargets(data.Targets)
}

// parseTextFile 解析文本文件
func (p *Parser) parseTextFile(file *os.File) ([]types.Target, error) {
	var targets []types.Target

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析行格式: IP端口范围 端口范围 [描述]
		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("第%d行格式错误: %s", lineNum, line)
		}

		ipRange := parts[0]
		portRange := parts[1]
		description := ""
		if len(parts) > 2 {
			description = strings.Join(parts[2:], " ")
		}

		// 解析目标
		lineTargets, err := p.parseRangeTargets(ipRange, portRange)
		if err != nil {
			return nil, fmt.Errorf("第%d行解析失败: %v", lineNum, err)
		}

		// 添加描述
		for i := range lineTargets {
			lineTargets[i].Description = description
		}

		targets = append(targets, lineTargets...)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取文件失败: %v", err)
	}

	return targets, nil
}

// parseExcelFile 解析Excel文件
func (p *Parser) parseExcelFile(filename, sheetName string) ([]types.Target, error) {
	xlFile, err := xlsx.OpenFile(filename)
	if err != nil {
		return nil, fmt.Errorf("打开Excel文件失败: %v", err)
	}

	var sheet *xlsx.Sheet
	for _, s := range xlFile.Sheets {
		if s.Name == sheetName {
			sheet = s
			break
		}
	}

	if sheet == nil {
		return nil, fmt.Errorf("找不到工作表: %s", sheetName)
	}

	var targets []types.Target

	for i, row := range sheet.Rows {
		// 跳过标题行
		if i == 0 {
			continue
		}

		if len(row.Cells) < 2 {
			continue
		}

		ipRange := strings.TrimSpace(row.Cells[0].String())
		portRange := strings.TrimSpace(row.Cells[1].String())

		if ipRange == "" || portRange == "" {
			continue
		}

		description := ""
		if len(row.Cells) > 2 {
			description = strings.TrimSpace(row.Cells[2].String())
		}

		var tags []string
		if len(row.Cells) > 3 {
			tagStr := strings.TrimSpace(row.Cells[3].String())
			if tagStr != "" {
				tags = strings.Split(tagStr, ",")
				for j := range tags {
					tags[j] = strings.TrimSpace(tags[j])
				}
			}
		}

		// 解析目标
		rowTargets, err := p.parseRangeTargets(ipRange, portRange)
		if err != nil {
			return nil, fmt.Errorf("第%d行解析失败: %v", i+1, err)
		}

		// 添加描述和标签
		for j := range rowTargets {
			rowTargets[j].Description = description
			rowTargets[j].Tags = tags
		}

		targets = append(targets, rowTargets...)
	}

	return targets, nil
}

// deduplicateTargets 去重目标
func (p *Parser) deduplicateTargets(targets []types.Target) []types.Target {
	seen := make(map[string]bool)
	var result []types.Target

	for _, target := range targets {
		key := fmt.Sprintf("%s:%d", target.Host, target.Port)
		if !seen[key] {
			seen[key] = true
			result = append(result, target)
		}
	}

	return result
}

// incrementIP 递增IP地址
func (p *Parser) incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}