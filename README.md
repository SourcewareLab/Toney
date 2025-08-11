# Toney

**Toney** is a fast, lightweight, terminal-based note-taking app for the modern developer. Built with [Bubbletea](https://github.com/charmbracelet/bubbletea), Toney brings a sleek TUI interface with markdown rendering, file navigation, and native Neovim editing – all in your terminal.

Refer to the docs [here](https://sourcewarelab.github.io/Toney/). To view all the details!





https://github.com/user-attachments/assets/71d4555f-8ea1-4a5c-b1df-c67acbd9248a





---

## ✨ Features

- ⚡ **Fast** – Minimal memory usage and snappy performance.
- 📝 **Markdown Renderer** – Styled previews via [`glamour`](https://github.com/charmbracelet/glamour).
- 🧠 **Neovim Integration** – Edit your notes using your favorite editor (`nvim`).
- 📂 **File Management** – Easily navigate, open, and manage markdown files.
- ✅ **Daily Tasks** – Keep a record of your daily tasks in the `Daily Tasks` section.
- 📓 **Diary Entries** - Organize and manage your thoughts with the `Diary` section.
- ⚙️ **Extensive Config** - Configurate the app to look and work exactly as `you` want.
- 🎨 **TUI Styling** – Clean, responsive interface using `lipgloss`.

---

## 🚀 Installation

Refer to the docs [here](https://sourcewarelab.github.io/Toney/guide/install). To view all the installation options!

Run this command to ensure Toney is setup perfectly.

```
  toney init
```

In order to use the GitHub integration, you need to create a GitHub Personal Access Token [here](https://github.com/settings/tokens/new?scopes=repo&description=Toney+GitHub+Integration).
In the `Toney` configuration file, set the following:
```
[github]
enabled = true
token = "your-token-here"
owner = "your-github-username"
repo = "your-github-repo"
```
In order to use toney to set up the GitHub integration, you need to run the following command:
```
  toney github setup
``` 
In order to check if the GitHub integration is working, you can run:
```
  toney github sync
```


### ✅ Verify Installation

Once installed, you can run:

```
toney init
```

---

## 🔑 Keybinds

Refer to the docs [here](https://sourcewarelab.github.io/Toney/config/keybinds/). To view all the Default Keybinds!

Keybinds are also visible in the app itself thanks to the `bubbles` component library.

---

## 🗺 Roadmap

### v2.0.0 Goals

- [ ] Overlay Error Popups
- [ ] File Import/Export
- [ ] Fuzzy Finder in notes
- [ ] Recurring and Unique Daily tasks
- [ ] Nvim TODO Integration
- [ ] Github Issues Integration
- [ ] Bug solve

### Long Term Goals

- [ ] Cross-platform **mobile app**
- [ ] **Server sync** with configuration & cloud storage

---

## 🛠️ Project Structure

```
├── cmd          // CLI 
└── internal
    ├── colors   // Styles Config
    ├── config   // Config 
    ├── enums    // All enums
    ├── fileTree // Filtree - manages tree creation for directory
    ├── keymap   // Keymap 
    ├── messages // Messages - has all bubbletea message declarations
    ├── models   // Models - Has all bubbletea models
    └── styles   // Styles - has all lipgloss styles declarations
```

---

## 🤝 Contributing

We welcome contributions! Toney follows Go and Bubbletea conventions. You can also view more details [here](https://sourcewarelab.github.io/Toney/guide/contribute).

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

