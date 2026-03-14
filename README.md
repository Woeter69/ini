# ⚡ ini

**The Blazing Fast Universal Project Initializer**

`ini` is a powerful, interactive CLI tool designed to scaffold new projects across multiple languages and domains in seconds. It combines a beautiful TUI experience with opinionated, high-quality templates to get you from zero to "Hello World" (or a full-blown API) instantly.

---

## 🚀 Features

- **Interactive TUI**: Built with [Charm.sh](https://charm.sh) libraries (`huh`, `lipgloss`) for a premium terminal experience.
- **Domain-Specific Scaffolding**: Intelligent categorization across domains (Web, DevOps, DB, Security, etc.) tailored per language.
- **Multi-Language Support**:
  - **Go**: 12+ domains with auto-dependency management.
  - **Python**: Integrated with `uv` for lightning-fast setup.
  - **Rust**: Native `cargo` integration with domain-specific templates.
  - **JS/TS (Bun)**: Native framework support (React, Vue, Svelte, Solid, Next.js) using instant embedded templates.
- **Zero Overhead**: Fully portable single binary powered by `go:embed`.
- **Git Ready**: Automatically initializes Git and professional `.gitignore` files.

---

## 🛠 Supported Languages & Domains

`ini` uses a structured set of domains to provide the right boilerplate for the right job.

### 🐹 Go
Natively supports: `web`, `devops`, `network`, `os`, `db`, `security`, `monitor`, `stream`, `comm`, `web3`, `lang`, `script`.

### 🐍 Python
Natively supports: `web` (FastAPI), `scraper`, `data` (Analytics), `cli`, `basic`.

### 🦀 Rust
Natively supports: `web`, `script`, `game`, `network`, `os`, `db`, `security`, `graphics`, `web3`, `lang`.

### ⚡ JavaScript / TypeScript (Bun)
Features a hierarchical picker for Frontend/Backend:
- **Frameworks**: React (JSX/TSX), Vue (JS/TS), Svelte (JS/TS), Solid (JS/TS), Next.js, Express.js, Vanilla.

---

## 📦 Installation

Initialize your project workspace:

```bash
go build -o ini .
sudo mv ini /usr/local/bin/ # Optional: add to PATH
```

---

## 🚀 Usage

### Interactive Mode
Simply run the command for your preferred language and follow the prompts:

```bash
ini go my-awesome-app
ini python data-tool
ini rust engine
ini bun website
```

### Non-Interactive (Flags)
You can scaffold projects instantly by providing the language, name, and domain flags. **A language argument is always required.**

```bash
# Syntax: ini [language] [project-name] --type [domain]
ini bun my-site --type web --framework react --variant ts
ini go my-service --type network
ini rust engine --type graphics
```

---

## 🛠 Development

### Prerequisites
- [Go](https://go.dev/dl/) 1.21+
- [Bun](https://bun.sh) (for JS/TS scaffolding)
- [uv](https://github.com/astral-sh/uv) (for Python scaffolding)

### Running Locally
```bash
go run main.go [lang] [name]
```

---

## 📄 License
MIT License. Built with ❤️ by **Woeter**.
