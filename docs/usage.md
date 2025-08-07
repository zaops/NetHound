# Nethound 使用指南

## 简介

Nethound 是一个轻量级网络连通性测试工具，支持多目标网络测试，能够同时对多个IP地址和端口进行ping和telnet连通性测试。

## 快速开始

### 基本用法

```bash
# 测试单个主机（ping + telnet）
.\nethound.exe -t "baidu.com:80"

# 只进行 ping 测试
.\nethound.exe -t "baidu.com" -T ping

# 只进行 telnet 测试
.\nethound.exe -t "8.8.8.8:53" -T telnet

# 测试多个目标
.\nethound.exe -t "baidu.com:80,google.com:443,8.8.8.8:53"
```

### 输出格式

```bash
# 默认表格输出（控制台）
.\nethound.exe -t "baidu.com:80"

# JSON 格式（自动保存到文件）
.\nethound.exe -t "baidu.com:80" -o json

# CSV 格式（自动保存到文件）
.\nethound.exe -t "baidu.com:80" -o csv

# HTML 格式（自动保存到文件）
.\nethound.exe -t "baidu.com:80" -o html

# 指定输出文件
.\nethound.exe -t "baidu.com:80" -o json -F result.json
```

## 命令行参数

### 测试目标参数

| 参数 | 简写 | 说明 | 示例 |
|------|------|------|------|
| `--targets` | `-t` | 目标列表，格式: ip:port | `-t "8.8.8.8:53,baidu.com:80"` |
| `--target-file` | `-f` | 目标配置文件路径 | `-f targets.txt` |
| `--ip-range` | `-r` | IP地址范围 | `-r "192.168.1.1-254"` |
| `--port-range` | `-p` | 端口范围 | `-p "80,443,8080-8090"` |
| `--excel-file` | `-e` | Excel文件路径 | `-e targets.xlsx` |

### 测试类型参数

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--test-type` | `-T` | `all` | 测试类型: ping, telnet, all |
| `--ping-count` | `-c` | `1` | ping测试次数 |
| `--ping-interval` | `-i` | `1` | ping间隔时间(秒) |
| `--telnet-timeout` | `-w` | `3000` | telnet连接超时(毫秒) |
| `--read-timeout` | `-d` | `1000` | 读取超时(毫秒) |

### 输出参数

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--output` | `-o` | `table` | 输出格式: table, json, csv, html |
| `--output-file` | `-F` | - | 输出文件路径 |
| `--verbose` | `-v` | `false` | 详细输出模式 |
| `--quiet` | `-q` | `false` | 静默模式 |
| `--no-color` | - | `false` | 禁用彩色输出 |

### 性能参数

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--concurrency` | `-C` | `20` | 并发测试数量 |
| `--rate-limit` | `-R` | `0` | 速率限制(次/秒) |
| `--retry` | `-y` | `0` | 失败重试次数 |
| `--retry-delay` | `-D` | `500` | 重试间隔(毫秒) |

## 配置文件

### 目标文件格式

创建 `targets.txt` 文件：

```
# IP范围 端口范围 描述
192.168.1.1-10 80,443 Web服务器
10.0.0.1 22,3389 管理服务器
baidu.com 80,443 外部网站
```

### JSON 配置文件

创建 `config.json` 文件：

```json
{
  "targets": ["8.8.8.8:53", "baidu.com:80"],
  "test_type": "all",
  "output": "json",
  "concurrency": 20,
  "telnet_timeout": 3000
}
```

使用配置文件：

```bash
.\nethound.exe --config config.json
```

## 输出示例

### 表格输出

```
+-----------+------+--------+-----------+---------+-------+
|   HOST    | PORT |  TYPE  |  STATUS   | RTT(MS) | ERROR |
+-----------+------+--------+-----------+---------+-------+
| baidu.com | -    | PING   | REACHABLE |   1.000 | -     |
| baidu.com |   80 | TELNET | REACHABLE |  41.522 | -     |
+-----------+------+--------+-----------+---------+-------+

测试摘要:
总计: 2 成功: 2 失败: 0 耗时: 89.7863ms
```

### 状态说明

- **REACHABLE**: 连接成功
- **TIMEOUT**: 连接超时
- **REFUSED**: 连接被拒绝
- **DNS_ERROR**: DNS解析失败
- **ERROR**: 其他错误

## 高级用法

### 批量测试

```bash
# 测试整个网段
.\nethound.exe -r "192.168.1.1-254" -p "22,80,443"

# 从Excel文件读取目标
.\nethound.exe -e targets.xlsx --excel-sheet "生产环境"

# 高并发测试
.\nethound.exe -f targets.txt -C 50 -w 2000
```

### 结果导出

```bash
# 导出为JSON并指定文件名
.\nethound.exe -t "baidu.com:80" -o json -F "test_$(date +%Y%m%d).json"

# 静默模式，只输出结果文件
.\nethound.exe -t "baidu.com:80" -o csv -q
```

## 构建和部署

### 多平台构建

使用提供的 `build.bat` 脚本：

```bash
.\build.bat
```

这将生成以下文件：
- `build/nethound-windows-amd64.exe`
- `build/nethound-windows-arm64.exe`
- `build/nethound-linux-amd64`
- `build/nethound-linux-arm64`

### 手动构建

```bash
# Windows
go build -o nethound.exe

# Linux
GOOS=linux go build -o nethound

# 交叉编译
GOOS=linux GOARCH=amd64 go build -o nethound-linux-amd64
```

## 故障排除

### 常见问题

1. **Ping 测试失败**
   - 检查防火墙设置
   - 确认目标主机允许ICMP

2. **权限问题**
   - Windows: 以管理员身份运行
   - Linux: 使用 sudo 或设置适当权限

3. **DNS 解析失败**
   - 检查网络连接
   - 尝试使用IP地址而不是域名

4. **性能问题**
   - 减少并发数 (`-C`)
   - 增加超时时间 (`-w`)
   - 使用速率限制 (`-R`)