# Nethound 网络测试工具需求文档

## 1. 项目概述

### 1.1 项目背景
Nethound 是一个基于 Go 1.24.5 开发的轻量级网络连通性测试工具，专为多云环境和混合网络架构设计。该工具支持多目标网络测试，能够同时对多个IP地址和端口进行ping和telnet连通性测试，并提供详细的测试结果分析。

### 1.2 项目目标
- 提供跨平台网络连通性测试解决方案
- 支持大规模网络节点批量测试
- 精确区分网络不通的不同原因
- 提供灵活的配置方式和丰富的输出格式
- 确保测试过程的高效性和准确性

## 2. 功能需求

### 2.1 核心功能
#### 2.1.1 网络连通性测试
- **ICMP Ping测试**: 支持对目标IP地址进行ICMP协议ping测试
- **TCP端口测试**: 支持对目标IP+端口进行TCP连接测试
- **并发测试**: 支持同时对多个目标进行测试
- **超时控制**: 支持自定义连接和响应超时时间

#### 2.1.2 结果分析
- **状态精确识别**:
  - `REACHABLE`: 目标可达
  - `TIMEOUT`: 连接超时（网络不可达或服务无响应）
  - `REFUSED`: 连接被拒绝（端口未开放或服务拒绝）
  - `DNS_ERROR`: DNS解析失败
  - `NETWORK_UNREACHABLE`: 网络不可达
  - `HOST_UNREACHABLE`: 主机不可达

#### 2.1.3 配置管理
- **命令行参数**: 支持通过命令行参数直接配置
- **配置文件**: 支持JSON格式的配置文件
- **环境变量**: 支持通过环境变量配置
- **配置优先级**: 命令行参数 > 配置文件 > 环境变量 > 默认值

### 2.2 高级功能
#### 2.2.1 测试模式
- **单次测试模式**: 执行一次测试后退出
- **持续监控模式**: 定期执行测试，持续监控网络状态
- **批量测试模式**: 从配置文件读取大量目标进行测试

#### 2.2.2 输出格式
- **标准输出**: 人类可读的表格格式
- **JSON格式**: 结构化数据输出
- **CSV格式**: 逗号分隔值格式，便于数据分析
- **日志文件**: 支持将结果输出到指定日志文件

#### 2.2.3 性能优化
- **连接池**: 复用TCP连接减少开销
- **并发控制**: 支持限制并发测试数量
- **结果缓存**: 缓存DNS解析结果减少重复查询
- **资源监控**: 监控自身CPU、内存使用情况

## 3. 参数设计

### 3.1 命令行参数

#### 3.1.1 测试目标参数
```bash
--targets, -t          # 目标列表，格式: ip:port 或 host:port
--target-file, -f      # 目标配置文件路径
--ip-range, -r         # IP地址范围，支持多种格式:
                      #   192.168.1.1-254
                      #   192.168.1.1-192.168.1.254
                      #   192.168.1.0/24
                      #   192.168.1.1,192.168.1.100-110
                      #   192.168.1-2.1-254
--port-range, -p       # 端口范围，如: 80,443,8080-8090,9000-9010
--excel-file, -e       # Excel文件路径，支持.xlsx和.xls格式
```

#### 3.1.2 测试类型参数
```bash
--test-type, -T        # 测试类型: ping, telnet, all (默认: all)
--ping-count, -c       # ping测试次数 (默认: 3)
--ping-interval, -i    # ping间隔时间(秒) (默认: 1)
--telnet-timeout, -to  # telnet连接超时(毫秒) (默认: 5000)
--read-timeout, -rt    # 读取超时(毫秒) (默认: 2000)
```

#### 3.1.3 输出参数
```bash
--output, -o           # 输出格式: table, json, csv, html (默认: table)
--output-file, -of     # 输出文件路径（默认：当前路径）
--verbose, -v          # 详细输出模式
--quiet, -q            # 静默模式，只显示错误
--no-color             # 禁用彩色输出
```

#### 3.1.4 性能参数
```bash
--concurrency, -C      # 并发测试数量 (默认: 10)
--rate-limit, -R       # 速率限制(次/秒)
--retry, -r            # 失败重试次数 (默认: 1)
--retry-delay, -rd     # 重试间隔(毫秒) (默认: 1000)
```

#### 3.1.5 配置参数
```bash
--config, -c           # 配置文件路径(默认：当前目录)
--config-format, -cf   # 配置文件格式: json, yaml, toml, txt (默认: json)
--excel-file, -e       # Excel文件路径 (.xlsx/.xls)
--excel-sheet          # Excel工作表名称 (默认: Sheet1)
--save-config, -sc     # 将当前参数保存为配置文件
--dry-run              # 预运行模式，不实际执行测试
```

### 3.2 配置文件格式

#### 3.2.1 极简配置文件格式
支持最简单的配置文件格式，每行一个目标：

**文本格式 (.txt)**
```
# 单行格式：IP或范围 端口或范围 [描述]
192.168.1.1-254 80,443,8080-8090 内网服务器
192.168.2.1-192.168.2.254 22,3306 数据库服务器
10.0.0.1 80,443 生产Web服务器
```

#### 3.2.2 Excel文件格式支持
支持从Excel文件读取配置，支持.xlsx和.xls格式：

**Excel文件格式要求**
- **Sheet1**作为默认读取的工作表
- **A列**: 目标IP地址或IP范围
- **B列**: 目标端口或端口范围
- **C列**: (可选) 描述信息
- **D列**: (可选) 标签信息，多个标签用逗号分隔

**Excel示例**
| A列(目标地址) | B列(端口) | C列(描述) | D列(标签) |
|--------------|----------|----------|----------|
| 192.168.1.1-254 | 80,443,8080 | 内网服务器 | web,internal |
| 192.168.2.1-192.168.2.254 | 22,3306 | 数据库集群 | db,production |
| 10.0.0.1 | 80,443 | 生产Web服务器 | web,production |
| example.com | 80,443,8080 | 外部服务 | external,web |

#### 3.2.3 JSON格式示例
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
    "concurrency": 20
  }
}
```

### 3.3 IP地址范围格式规范

#### 3.3.1 支持的IP地址格式
```bash
# 单IP地址
192.168.1.1

# 连续IP范围
192.168.1.1-254
192.168.1.1-192.168.1.254

# CIDR格式
192.168.1.0/24
10.0.0.0/8

# 多段IP组合
192.168.1.1-254,192.168.2.1-254
192.168.1-2.1-254
192.168.1.1,192.168.1.100-110,10.0.0.50

# 主机名
example.com
google.com
```

#### 3.3.2 支持的端口格式
```bash
# 单端口
80

# 连续端口范围
80-90

# 多个端口组合
80,443,8080
80,443,8080-8090

# 混合格式
80,443,8080-8090,9000-9010
```

## 4. 环境变量配置

### 4.1 支持的系统环境变量
```bash
NETHOUND_TARGETS          # 默认测试目标，支持上述所有格式
NETHOUND_CONFIG_FILE      # 默认配置文件路径
NETHOUND_OUTPUT_FORMAT    # 默认输出格式: table, json, csv, html
NETHOUND_LOG_LEVEL        # 日志级别: debug, info, warn, error
NETHOUND_MAX_CONCURRENCY  # 最大并发数
NETHOUND_TIMEOUT         # 默认超时时间
NETHOUND_RETRY_COUNT     # 默认重试次数
NETHOUND_RETRY_DELAY     # 默认重试间隔
NETHOUND_EXCEL_SHEET     # Excel默认工作表名称，默认为"Sheet1"
```

### 4.2 配置优先级
1. 命令行参数（最高优先级）
2. 配置文件
3. 环境变量
4. 默认值（最低优先级）

## 5. 输出格式规范

### 5.1 支持输出格式（已更新）
- **table**: 标准表格格式（终端显示）
- **json**: 结构化JSON格式
- **csv**: 逗号分隔值格式
- **html**: 网页格式，包含CSS样式和颜色标识
#### 5.1.1 标准表格输出
```
┌─────────────────┬──────────┬──────┬──────────┬─────────┬────────────┐
│ Host            │ Port     │ Type │ Status   │ RTT(ms) │ Error      │
├─────────────────┼──────────┼──────┼──────────┼─────────┼────────────┤
│ 192.168.1.100   │ 80       │ TCP  │ OK       │ 1.234   │ -          │
│ 192.168.1.100   │ 443      │ TCP  │ TIMEOUT  │ -       │ timeout    │
│ 192.168.1.100   │ 8080     │ TCP  │ REFUSED  │ -       │ connection refused │
│ example.com     │ 22       │ ICMP │ OK       │ 12.567  │ -          │
└─────────────────┴──────────┴──────┴──────────┴─────────┴────────────┘
```

#### 5.1.2 HTML格式输出
```html
<!DOCTYPE html>
<html>
<head>
    <title>Nethound Network Test Results</title>
    <style>
        table { border-collapse: collapse; width: 100%; }
        th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
        th { background-color: #f2f2f2; }
        .success { color: green; }
        .timeout { color: orange; }
        .error { color: red; }
    </style>
</head>
<body>
    <h1>Network Test Results</h1>
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
            </tr>
        </thead>
        <tbody>
            <tr>
                <td>192.168.1.100</td>
                <td>80</td>
                <td>TCP</td>
                <td class="success">OK</td>
                <td>1.234</td>
                <td>-</td>
                <td>2024-01-15 14:30:00</td>
            </tr>
        </tbody>
    </table>
</body>
</html>
```

### 5.2 JSON输出格式
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
      "status": "ok",
      "rtt_ms": 1.234,
      "timestamp": "2024-01-15T14:30:00Z",
      "error": null,
      "tags": ["production", "web"]
    },
    {
      "host": "192.168.1.100",
      "port": 443,
      "type": "tcp",
      "status": "timeout",
      "rtt_ms": null,
      "timestamp": "2024-01-15T14:30:01Z",
      "error": "connection timeout after 3000ms",
      "tags": ["production", "web"]
    }
  ]
}
```

## 6. 错误处理机制

### 6.1 错误类型定义
```go
type TestError struct {
    Type    ErrorType `json:"type"`
    Message string    `json:"message"`
    Code    int       `json:"code,omitempty"`
}

type ErrorType string

const (
    ErrorTimeout          ErrorType = "timeout"
    ErrorConnectionRefused ErrorType = "connection_refused"
    ErrorNetworkUnreachable ErrorType = "network_unreachable"
    ErrorHostUnreachable   ErrorType = "host_unreachable"
    ErrorDNS              ErrorType = "dns_error"
    ErrorPermission        ErrorType = "permission_denied"
    ErrorInternal          ErrorType = "internal_error"
)
```

### 6.2 错误恢复策略
- **指数退避重试**: 失败时采用指数退避算法重试
- **断路器模式**: 连续失败后暂停对该目标的测试
- **优雅降级**: 在资源不足时降低并发度
- **错误聚合**: 相同类型的错误进行聚合展示

## 7. 性能要求

### 7.1 性能指标
- **并发能力**: 支持最多1000个并发测试
- **测试延迟**: 单次测试延迟 < 100ms（本地网络）
- **内存使用**: 最大内存占用 < 100MB（1000个目标）
- **CPU使用**: 单核CPU占用 < 50%（1000个并发）
- **准确性**: 测试结果准确率 > 99.9%

### 7.2 资源限制
- **连接池大小**: 最大1000个并发连接
- **DNS缓存**: TTL 5分钟，最大1000条记录
- **结果缓存**: 最大10000条测试结果
- **日志轮转**: 单个日志文件最大100MB，保留7天

## 8. 安全要求

### 8.1 网络安全
- **连接超时**: 防止半开连接攻击
- **速率限制**: 防止对目标造成DoS攻击
- **DNS安全**: 支持DNSSEC验证
- **TLS支持**: 支持TLS连接测试（可选）

### 8.2 数据安全
- **敏感信息脱敏**: 配置文件中的密码等敏感信息脱敏显示
- **日志安全**: 日志中不包含敏感信息
- **配置加密**: 支持配置文件加密存储
- **权限控制**: 配置文件权限控制（600）

## 9. 扩展性设计

### 9.1 插件机制
- **测试插件**: 支持自定义测试类型插件
- **输出插件**: 支持自定义输出格式插件
- **认证插件**: 支持自定义认证方式插件

### 9.2 API接口
- **REST API**: 提供HTTP API进行远程控制
- **gRPC API**: 提供高性能RPC接口
- **WebSocket**: 支持实时结果推送

## 10. 部署和运维

### 10.1 安装方式
- **二进制包**: 提供预编译的二进制文件
- **容器镜像**: 提供Docker镜像
- **包管理器**: 支持apt、yum、brew等包管理器
- **源码编译**: 支持从源码编译安装

### 10.2 运维支持
- **健康检查**: 提供健康检查接口
- **日志级别**: 支持动态调整日志级别
- **热重载**: 支持配置文件热重载

## 11. 测试用例

### 11.1 功能测试
- **基本连通性测试**: 测试已知可达和不可达的目标
- **并发测试**: 测试不同并发度下的性能
- **错误处理测试**: 测试各种错误情况的处理
- **配置测试**: 测试不同配置格式的兼容性

### 11.2 性能测试
- **负载测试**: 测试最大并发能力
- **压力测试**: 测试长时间运行的稳定性
- **内存测试**: 测试内存泄漏情况
- **网络测试**: 测试不同网络环境下的性能

## 12. 版本规划

### 12.1 v1.0.0 (MVP)
- 基本ping和telnet测试功能
- 命令行参数支持
- JSON配置文件支持
- 表格和JSON输出格式

### 12.2 v1.1.0
- YAML和TOML配置文件支持
- Prometheus输出格式
- 并发测试优化
- 错误处理增强

### 12.3 v1.2.0
- 持续监控模式
- Web界面
- REST API支持
- 插件机制

### 12.4 v2.0.0
- 分布式测试
- 高级分析功能
- 机器学习异常检测
- 企业级特性

## 13. 验收标准

### 13.1 功能验收
- [ ] 能够正确执行ping测试
- [ ] 能够正确执行telnet测试
- [ ] 能够区分timeout和connection refused
- [ ] 支持命令行参数配置
- [ ] 支持JSON配置文件
- [ ] 支持多种输出格式
- [ ] 支持并发测试

### 13.2 性能验收
- [ ] 1000个并发测试下CPU占用<50%
- [ ] 1000个并发测试下内存占用<100MB
- [ ] 单次测试延迟<100ms
- [ ] 连续运行24小时无内存泄漏

### 13.3 安全验收
- [ ] 通过安全扫描
- [ ] 无已知安全漏洞
- [ ] 敏感信息正确处理
- [ ] 权限控制正确
