# ⚡ ini v1.0.0

**The Blazing Fast Universal Project Initializer**

`ini` is a powerful, interactive CLI tool designed to scaffold new projects across **39+ programming languages** and dozens of domains in seconds. It combines a beautiful TUI experience with opinionated, high-quality templates to get you from zero to a production-ready boilerplate instantly.

---

## 🚀 Features

- **Universal Interactive Picker**: Run `ini` without arguments to browse and select from all 39+ supported languages.
- **Dynamic TUI**: Built with [Charm.sh](https://charm.sh) libraries (`huh`, `lipgloss`) for a smooth, interactive terminal experience.
- **Global Taxonomy**: Standardized project categories (Web, API, CLI, OS, Embedded, AI, Data, etc.) applied consistently across the tool.
- **Mobile & Legacy Support**: First-class scaffolding for **Flutter, Dart, Objective-C**, and even low-level **Assembly** (x86_64, BIOS).
- **Toolchain Integrated**:
  - **Go**: Native module and dependency management.
  - **Python**: Powererd by `uv` for lightning-fast environment setup.
  - **Rust**: Full `cargo` integration with domain-specific templates.
  - **JS/TS**: Next.js, React, Vue, Svelte, and Solid support via Bun.
- **Zero Overhead**: Fully portable single binary powered by `go:embed`.

---

## 🛠 Supported Languages (Highlights)

`ini` supports a massive range of languages, each with specialized domains:

- **Modern**: `Go`, `Rust`, `Python`, `Zig`, `Swift`, `Kotlin`, `Nim`, `Bun (JS/TS)`, `Julia`.
- **Systems/Legacy**: `C`, `C++`, `Assembly`, `COBOL`, `Fortran`, `Ada`, `Pascal`, `Objective-C`.
- **Functional/Specialized**: `Haskell`, `Ocaml`, `Clojure`, `Erlang`, `Elixir`, `Scala`, `R`.
- **Scripting**: `Ruby`, `Perl`, `PHP`, `Lua`, `Shell (Bash)`, `V`, `Crystal`.

---

## 🚀 Usage

### Global Interactive Mode
Simply run `ini` to open the universal language picker:

```bash
ini
```

### Direct Scaffolding
Quickly start a project by specifying the language and name:

```bash
ini go my-app
ini assemble bootloader --type os
ini python ai-tool --type ai
```

### Non-Interactive (Flags)
Scaffold projects instantly with specific categories:

```bash
# Syntax: ini [language] [project-name] --type [category]
ini rust game-engine --type game
ini bun web-app --type web --framework next --variant ts
ini asm drive-controller --type embedded
```

---

## 🛠 Development

### Prerequisites
- [Go](https://go.dev/dl/) 1.21+
- Language-specific toolchains (e.g., `uv`, `cargo`, `nasm`, `bun`) depending on what you scaffold.

### Running Locally
```bash
go run main.go
```

### Internal Test Suite
We maintain an internal integration suite to verify all 39+ handlers:
```bash
./tests/verify_all.sh
```

---

## 📄 License
MIT License. Built with ❤️ by **Woeter**.
