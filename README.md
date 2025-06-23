# Toney

**Toney** is a fast, lightweight, terminal-based note-taking app for the modern developer. Built with [Bubbletea](https://github.com/charmbracelet/bubbletea), Toney brings a sleek TUI interface with markdown rendering, file navigation, and native Neovim editing – all in your terminal.

---

## ✨ Features

- ⚡ **Blazingly Fast** – Minimal memory usage and snappy performance.
- 📝 **Markdown Renderer** – Styled previews via [`glamour`](https://github.com/charmbracelet/glamour).
- 🧠 **Neovim Integration** – Edit your notes using your favorite editor (`nvim`).
- 📂 **File Management** – Easily navigate, open, and manage markdown files.
- 🧩 **Component Architecture** – Modular codebase using Bubbletea’s `Model` system.
- 🎨 **TUI Styling** – Clean, responsive interface using `lipgloss`.

---

## 📦 Installation

> Prerequisites:  
> - Go 1.22+  
> - [Neovim](https://neovim.io/) installed (`nvim` command)

```bash
git clone https://github.com/NucleoFusion/Toney
cd Toney
go build -o toney ./cmd/toney
./toney
```

---

## 🗺 Roadmap

### Short Term Goals

- [ ] Overlay support
- [ ] Viewer style improvements
- [ ] Error popups
- [X] Separate package for messages
- [ ] Keybind refactor
- [ ] Config file support (`~/.config/toney/config.yaml`)
- [ ] Custom markdown renderer
- [ ] Custom components:  
  - [ ] [ ] Task Lists  
  - [ ] `code` blocks  
  - [ ] Tables  
- [ ] File Import/Export
- [ ] Configurable external editor support

### Long Term Goals

- [ ] Cross-platform **mobile app**
- [ ] **Server sync** with configuration & cloud storage

---

## 🛠️ Project Structure

```
toney/
├── cmd/              # Entry point (main.go)
├── internal/
│   ├── models/       # All UI models (Home, Viewer, Popups, etc.)
│   ├── enums/        # Typed enums (pages, popup types)
│   ├── messages/     # Message types for tea.Msg (will be modularized)
│   └── utils/        # Shared utility functions
```

---

## 🤝 Contributing

We welcome contributions! Toney follows Go and Bubbletea conventions.

### 🧾 Guidelines

- Follow idiomatic Go formatting (`go fmt`, `go vet`, `golint`).
- Use `Init`, `Update`, and `View` separation for all models.
- Keep component responsibilities well-isolated.
- All exported functions/types should be documented with Go-style comments.
- Prefer `tea.Msg` messages over direct cross-component function calls.

### ✅ How to contribute

1. Fork the repo and create a feature branch:
   ```bash
   git checkout -b feature/my-feature
   ```

2. Write your code and make sure it builds:
   ```bash
   go build ./...
   ```

3. Format your code:
   ```bash
   go fmt ./...
   ```

4. Commit and push:
   ```bash
   git commit -m "feat: add my awesome feature"
   git push origin feature/my-feature
   ```

5. Submit a Pull Request 🎉

Please open an issue or discussion for large changes before starting them.

---

## 📄 License

MIT License. See [LICENSE](./LICENSE).

---

## 💡 Inspiration

Toney is inspired by:
- [Glow](https://github.com/charmbracelet/glow) – for markdown rendering  
- [Lazygit](https://github.com/jesseduffield/lazygit) – for terminal UI polish  
- [Charm ecosystem](https://github.com/charmbracelet) – for all things delightful in the terminal

---

> Made with 💀 by [Nucleo](https://github.com/NucleoFusion) & [SourcewareLab](https://discord.gg/X69MUr2DKm)

