# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-08-07

### Added
- 初始版本发布
- 支持 PING 和 TELNET 连通性测试
- 支持多种输出格式：TABLE、JSON、CSV、HTML
- 支持多平台构建：Windows、Linux (amd64/arm64)
- 支持多种目标格式：单个目标、IP范围、CIDR、端口范围
- 支持配置文件和Excel文件输入
- 支持高并发测试（默认20个并发）
- 支持彩色表格输出
- 完整的命令行参数支持
- 详细的文档和使用指南

### Features
- **高性能**: 优化的并发测试，支持大规模网络扫描
- **多测试类型**: PING、TELNET、ALL 三种测试模式
- **灵活配置**: 命令行参数、配置文件、Excel文件多种配置方式
- **丰富输出**: 表格、JSON、CSV、HTML 多种输出格式
- **跨平台**: Windows 和 Linux 系统支持
- **批量测试**: IP范围、端口范围、CIDR 等批量测试功能

### Technical Details
- Go 1.19+ 支持
- 使用 Cobra 框架构建CLI
- 使用 Viper 进行配置管理
- 支持中文和英文环境
- 完整的错误处理和重试机制
- 自动文件输出和时间戳命名

### Documentation
- 完整的使用指南
- API 开发文档
- 配置示例文件
- 多平台构建脚本

## [Unreleased]

### Planned
- 支持更多测试协议（HTTP、HTTPS、FTP等）
- 支持测试结果的历史记录和对比
- 支持Web界面管理
- 支持分布式测试
- 支持更多输出格式（XML、YAML等）
- 支持测试结果的图表展示