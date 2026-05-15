# CLIProxyAPI — Ai-Data-Man Fork

This is a **maintained fork** of [router-for-me/CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI).

Upstream is a proxy server providing OpenAI/Gemini/Claude/Codex compatible API interfaces for CLI tools (Claude Code, Codex CLI, Gemini CLI, etc.), supporting OAuth login, multi-account load balancing, streaming, and more.

## What We Changed

- **CI/CD**: Added `workflow_dispatch` trigger for manual releases
- **Panel default repository**: Default `panel-github-repository` points to [Ai-Data-Man/Cli-Proxy-API-Management-Center](https://github.com/Ai-Data-Man/Cli-Proxy-API-Management-Center)
- **Usage statistics**: Restored built-in usage statistics backend endpoints and in-memory stats plugin (removed upstream since v6.10.0)
- **Auth improvements**: WebSocket session reuse for home auths, scoped caching, disabled flag persistence, per-auth `disable_cooling` override
- **Logging**: Home request-log forwarding support
- **Config**: Remote home control plane loading support

## Documentation

For full upstream features, usage guides, and API docs:

- [router-for-me/CLIProxyAPI](https://github.com/router-for-me/CLIProxyAPI) — original repository
- [CLIProxyAPI Help](https://help.router-for.me/) — official guides

## Related Forks

- [Ai-Data-Man/Cli-Proxy-API-Management-Center](https://github.com/Ai-Data-Man/Cli-Proxy-API-Management-Center) — Web UI fork with restored usage statistics
- [Ai-Data-Man/CLIProxyAPI-Tray](https://github.com/Ai-Data-Man/CLIProxyAPI-Tray) — Windows tray manager fork

## License

MIT
