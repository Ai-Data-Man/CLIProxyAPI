# CLIProxyAPI — Ai-Data-Man 维护分支

本仓库是 [router-for-me/CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI) 的 **维护分支**。

上游项目是一个为 CLI 工具（Claude Code、Codex CLI、Gemini CLI 等）提供 OpenAI/Gemini/Claude/Codex 兼容 API 接口的代理服务器，支持 OAuth 登录、多账户负载均衡、流式响应等。

## 我们的改动

- **CI/CD**：增加 `workflow_dispatch` 手动触发发版
- **面板默认地址**：默认 `panel-github-repository` 指向 [Ai-Data-Man/Cli-Proxy-API-Management-Center](https://github.com/Ai-Data-Man/Cli-Proxy-API-Management-Center)
- **使用量统计**：恢复内置使用统计后端端点和内存统计插件（上游自 v6.10.0 起移除）
- **认证改进**：Home 模式 WebSocket 会话复用、作用域缓存、禁用状态持久化、单凭据 `disable_cooling` 覆盖
- **日志**：Home 请求日志转发支持
- **配置**：支持从远程 Home 控制面加载配置

## 文档

完整功能文档和使用指南请参见上游：

- [router-for-me/CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI) — 原始仓库
- [CLIProxyAPI 用户手册](https://help.router-for.me/cn/) — 官方指南

## 相关 Fork

- [Ai-Data-Man/Cli-Proxy-API-Management-Center](https://github.com/Ai-Data-Man/Cli-Proxy-API-Management-Center) — Web 管理面板 fork（恢复使用统计功能）
- [Ai-Data-Man/CLIProxyAPI-Tray](https://github.com/Ai-Data-Man/CLIProxyAPI-Tray) — Windows 托盘管理工具 fork

## 许可证

MIT
