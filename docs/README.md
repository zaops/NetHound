# Nethound 网络测试工具

Nethound 是一个基于 Go 开发的轻量级网络连通性测试工具，专为多云环境和混合网络架构设计。该工具支持多目标网络测试，能够同时对多个IP地址和端口进行ping和telnet连通性测试。

## 功能特性

- **多种测试类型**: 支持 ICMP Ping 和 TCP 端口连通性测试
- **批量测试**: 支持同时测试多个目标，提高测试效率
- **灵活的目标配置**: 支持命令行参数、配置文件、Excel文件等多种方式
- **丰富的输出格式**: 支持表格、JSON、CSV、HTML等输出格式
- **并发测试**: 支持自定义并发数量，提高测试速度
- **错误分类**: 精确识别网络不通的不同原因
- **重试机制**: 支持失败重试，提高测试准确性

## 安装

### 从源码编译

```bash
git clone <repository-url>
cd nethound
go mod tidy
go build -o nethound main.go
```

## 使用方法

### 基本用法

```bash
# 测试单个目标
./nethound -t "192.168.1.1:80"

# 测试多个目标
./nethound -t "192.168.1.1:80,192.168.1.1:443,google.com:80"

# 测试IP范围和端口范围
./nethound -r "192.168.1.1-254" -p "80,443,8080-8090"

# 从配置文件读取目标
./nethound -f targets.json

# 从Excel文件读取目标
./nethound -e targets.xlsx
```

### 测试类型

```bash
# 只执行ping测试
./nethound -t "192.168.1.1:80" -T ping

# 只执行telnet测试
./nethound -t "192.168.1.1:80" -T telnet

# 执行所有测试（默认）
./nethound -t "192.168.1.1:80" -T all
```

### 输出格式

```bash
# 表格格式输出（默认）
./nethound -t "192.168.1.1:80" -o table

# JSON格式输出
./nethound -t "192.168.1.1:80" -o json

# CSV格式输出
./nethound -t "192.168.1.1:80" -o csv

# HTML格式输出
./nethound -t "192.168.1.1:80" -o html

# 输出到文件
./nethound -t "192.168.1.1:80" -o json -of results.json
```

### 性能调优

```bash
# 设置并发数
./nethound -t "192.168.1.1-254:80" -C 50

# 设置重试次数和间隔
./nethound -t "192.168.1.1:80" --retry 3 --retry-delay 2000

# 设置超时时间
./nethound -t "192.168.1.1:80" --telnet-timeout 3000
```

## 配置文件格式

### JSON格式 (nethound.json)

```json
{
  "targets": [
    "192.168.1.1-254:80,443,8080-8090",
    "192.168.2.1-192.168.2.254:22,3306",
    "10.0.0.1:80,443"
  ],
  "settings": {
    "test_type": "all",
    "ping_count": 4,
    "telnet_timeout": 3000,
    "concurrency": 20,
    "output": "table",
    "retry": 2
  }
}
```

### 文本格式 (targets.txt)

```
# 格式：IP或范围 端口或范围 [描述]
192.168.1.1-254 80,443,8080-8090 内网服务器
192.168.2.1-192.168.2.254 22,3306 数据库服务器
10.0.0.1 80,443 生产Web服务器
```

### Excel格式

| A列(目标地址) | B列(端口) | C列(描述) | D列(标签) |
|--------------|----------|----------|----------|
| 192.168.1.1-254 | 80,443,8080 | 内网服务器 | web,internal |
| 192.168.2.1-192.168.2.254 | 22,3306 | 数据库集群 | db,production |
| 10.0.0.1 | 80,443 | 生产Web服务器 | web,production |

## 命令行参数

### 测试目标参数
- `--targets, -t`: 目标列表，格式: ip:port 或 host:port
- `--target-file, -f`: 目标配置文件路径
- `--ip-range, -r`: IP地址范围
- `--port-range, -p`: 端口范围
- `--excel-file, -e`: Excel文件路径
- `--excel-sheet`: Excel工作表名称（默认: Sheet1）

### 测试类型参数
- `--test-type, -T`: 测试类型: ping, telnet, all（默认: all）
- `--ping-count, -c`: ping测试次数（默认: 3）
- `--ping-interval, -i`: ping间隔时间(秒)（默认: 1）
- `--telnet-timeout, -to`: telnet连接超时(毫秒)（默认: 5000）
- `--read-timeout, -rt`: 读取超时(毫秒)（默认: 2000）

### 输出参数
- `--output, -o`: 输出格式: table, json, csv, html（默认: table）
- `--output-file, -of`: 输出文件路径
- `--verbose, -v`: 详细输出模式
- `--quiet, -q`: 静默模式，只显示错误
- `--no-color`: 禁用彩色输出

### 性能参数
- `--concurrency, -C`: 并发测试数量（默认: 10）
- `--rate-limit, -R`: 速率限制(次/秒)
- `--retry`: 失败重试次数（默认: 1）
- `--retry-delay`: 重试间隔(毫秒)（默认: 1000）

### 其他参数
- `--config, -c`: 配置文件路径
- `--dry-run`: 预运行模式，不实际执行测试
- `--save-config`: 将当前参数保存为配置文件

## 环境变量

支持以下环境变量：

```bash
export NETHOUND_TARGETS="192.168.1.1:80,192.168.1.1:443"
export NETHOUND_OUTPUT_FORMAT="json"
export NETHOUND_MAX_CONCURRENCY="20"
export NETHOUND_TIMEOUT="3000"
export NETHOUND_RETRY_COUNT="2"
export NETHOUND_RETRY_DELAY="1500"
```

## 输出示例

### 表格格式

```
┌─────────────────┬──────────┬──────┬──────────┬─────────┬────────────┐
│ Host            │ Port     │ Type │ Status   │ RTT(ms) │ Error      │
├─────────────────┼──────────┼──────┼──────────┼─────────┼────────────┤
│ 192.168.1.100   │ 80       │ TCP  │ REACHABLE│ 1.234   │ -          │
│ 192.168.1.100   │ 443      │ TCP  │ TIMEOUT  │ -       │ timeout    │
│ 192.168.1.100   │ 8080     │ TCP  │ REFUSED  │ -       │ connection refused │
│ example.com     │ -        │ ICMP │ REACHABLE│ 12.567  │ -          │
└─────────────────┴──────────┴──────┴──────────┴─────────┴────────────┘

测试摘要:
总计: 4, 成功: 2, 失败: 2, 耗时: 1.25s
```

### JSON格式

```json
{
  "timestamp": "2024-01-15T14:30:00Z",
  "summary": {
    "total": 4,
    "success": 2,
    "failed": 2,
    "duration_ms": 1250
  },
  "results": [
    {
      "host": "192.168.1.100",
      "port": 80,
      "type": "tcp",
      "status": "REACHABLE",
      "rtt_ms": 1.234,
      "timestamp": "2024-01-15T14:30:00Z",
      "error": null,
      "tags": ["production", "web"]
    }
  ]
}
```

## 错误状态说明

- `REACHABLE`: 目标可达
- `TIMEOUT`: 连接超时（网络不可达或服务无响应）
- `REFUSED`: 连接被拒绝（端口未开放或服务拒绝）
- `DNS_ERROR`: DNS解析失败
- `NETWORK_UNREACHABLE`: 网络不可达
- `HOST_UNREACHABLE`: 主机不可达
- `ERROR`: 其他错误

## 使用场景

1. **网络连通性检查**: 快速检查服务器或服务的可达性
2. **批量端口扫描**: 检查多个主机的多个端口状态
3. **网络故障排查**: 精确定位网络问题的原因
4. **监控脚本**: 集成到监控系统中进行定期检查
5. **网络规划**: 验证网络配置和防火墙规则

## 性能指标

- 支持最多1000个并发测试
- 单次测试延迟 < 100ms（本地网络）
- 内存占用 < 100MB（1000个目标）
- 支持大规模网络节点批量测试

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。

## 更新日志

### v1.0.0
- 基本ping和telnet测试功能
- 命令行参数支持
- JSON配置文件支持
- 表格、JSON、CSV、HTML输出格式
- 并发测试支持
- 错误分类和重试机制