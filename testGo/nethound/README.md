# Nethound 网络连通性测试工具

[![Go Version](https://img.shields.io/badge/Go-1.24.5-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux-lightgrey.svg)](https://github.com)

Nethound 是一个轻量级、高性能的网络连通性测试工具，支持多目标并发测试，能够同时对多个IP地址和端口进行 ping 和 telnet 连通性测试。

## ✨ 特性

- 🚀 **高性能**: 支持高并发测试，默认20个并发连接
- 🎯 **多测试类型**: 支持 PING、TELNET 和组合测试
- 📊 **多输出格式**: 支持 TABLE、JSON、CSV、HTML 格式输出
- 🔧 **灵活配置**: 支持命令行参数、配置文件、Excel文件等多种配置方式
- 🌐 **跨平台**: 支持 Windows 和 Linux 系统
- 📁 **批量测试**: 支持IP范围、端口范围、CIDR等批量测试
- 🎨 **美观输出**: 支持彩色表格输出和详细的测试报告

## 🚀 快速开始

### 下载安装

从 [Releases](https://github.com/your-repo/nethound/releases) 页面下载对应平台的可执行文件：

- Windows: `nethound-windows-amd64.exe`
- Linux: `nethound-linux-amd64`

### 基本使用

```bash
# 测试单个主机（ping + telnet）
./nethound -t "baidu.com:80"

# 只进行 ping 测试
./nethound -t "baidu.com" -T ping

# 只进行 telnet 测试
./nethound -t "8.8.8.8:53" -T telnet

# 测试多个目标
./nethound -t "baidu.com:80,google.com:443,8.8.8.8:53"
```

### 输出示例

```
开始测试 1 个目标...
测试完成，耗时: 1.1315888s
+-----------+------+--------+-----------+---------+-------+
|   HOST    | PORT |  TYPE  |  STATUS   | RTT(MS) | ERROR |
+-----------+------+--------+-----------+---------+-------+
| baidu.com | -    | PING   | REACHABLE |   1.000 | -     |
| baidu.com |   80 | TELNET | REACHABLE |  41.522 | -     |
+-----------+------+--------+-----------+---------+-------+

测试摘要:
总计: 2 成功: 2 失败: 0 耗时: 89.7863ms
```

## 📖 文档

- [使用指南](docs/usage.md) - 详细的使用说明和示例
- [API 文档](docs/api.md) - 开发者API参考
- [需求文档](docs/nethound网络测试工具需求文档.md) - 项目需求和设计

## 🛠️ 构建

### 环境要求

- Go 1.19 或更高版本
- Git

### 克隆项目

```bash
git clone https://github.com/your-repo/nethound.git
cd nethound
```

### 单平台构建

```bash
# 构建当前平台
go build -o nethound

# Windows
go build -o nethound.exe

# Linux
GOOS=linux go build -o nethound
```

### 多平台构建

使用提供的构建脚本：

```bash
# Windows
./build.bat

# 生成的文件在 build/ 目录下：
# - nethound-windows-amd64.exe
# - nethound-windows-arm64.exe  
# - nethound-linux-amd64
# - nethound-linux-arm64
```

## 📋 命令行参数

### 基本参数

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--targets` | `-t` | - | 目标列表，格式: ip:port |
| `--test-type` | `-T` | `all` | 测试类型: ping, telnet, all |
| `--output` | `-o` | `table` | 输出格式: table, json, csv, html |
| `--concurrency` | `-C` | `20` | 并发测试数量 |
| `--verbose` | `-v` | `false` | 详细输出模式 |

### 完整参数列表

查看所有可用参数：

```bash
./nethound --help
```

## 📁 配置文件

### 目标文件 (targets.txt)

```
# IP范围 端口范围 描述
192.168.1.1-10 80,443 Web服务器
10.0.0.1 22,3389 管理服务器
baidu.com 80,443 外部网站
```

### JSON 配置文件

```json
{
  "targets": ["8.8.8.8:53", "baidu.com:80"],
  "test_type": "all",
  "output": "json",
  "concurrency": 20,
  "telnet_timeout": 3000
}
```

## 🎯 使用场景

- **网络诊断**: 快速检测网络连通性问题
- **服务监控**: 批量检测服务端口状态
- **网络扫描**: 发现网络中的活跃主机和服务
- **性能测试**: 测试网络延迟和响应时间
- **自动化运维**: 集成到监控和部署脚本中

## 🔧 高级功能

### 批量测试

```bash
# 测试整个网段
./nethound -r "192.168.1.1-254" -p "22,80,443"

# 从Excel文件读取目标
./nethound -e targets.xlsx --excel-sheet "生产环境"

# 高并发测试
./nethound -f targets.txt -C 50 -w 2000
```

### 结果导出

```bash
# 导出为JSON
./nethound -t "baidu.com:80" -o json

# 导出为CSV
./nethound -t "baidu.com:80" -o csv

# 导出为HTML报告
./nethound -t "baidu.com:80" -o html
```

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [cobra](https://github.com/spf13/cobra) - 强大的CLI框架
- [viper](https://github.com/spf13/viper) - 配置管理
- [tablewriter](https://github.com/olekukonko/tablewriter) - 表格输出
- [xlsx](https://github.com/tealeg/xlsx) - Excel文件处理

## 📞 联系

- 项目主页: [https://github.com/zaops/NetHound](https://github.com/your-repo/NetHound)
- 问题反馈: [Issues](https://github.com/zaops/NetHound/issues)

---

⭐ 如果这个项目对你有帮助，请给它一个星标，谢谢！