# 🧠 AgentHub CLI

**AgentHub** is the universal CLI to install, run, and manage AI agents, tools, prompts, chains, and packages.

## 🚀 Install

### Homebrew (macOS & Linux):
```bash
brew tap agenthubcli/tap
brew install agenthub
```

### Scoop (Windows):
```powershell
scoop bucket add agenthub https://github.com/agenthubcli/scoop-agenthub
scoop install agenthub
```

### Chocolatey (Windows):
```powershell
choco install agenthub
```

## 🔧 Usage

```bash
agenthub init             # Start a new agent project
agenthub install agent    # Install an agent from registry
agenthub run my-agent     # Run an agent
agenthub publish          # Publish your agent
```

## 📦 Supported Packages

- `agentpkg.yaml` manifest
- Agents, Tools, Prompts, Chains, Datasets
- Built-in registry + local sandbox runner

## 📄 License

MIT — see [LICENSE](./LICENSE)
