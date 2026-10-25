<p align="center">
  <h1 align="center">🌊 Nagare</h1>
  <p align="center">
    <strong>A lightweight, local AI desktop assistant that lives on your machine.</strong>
  </p>
  <p align="center">
    <a href="https://github.com/smtdfc/nagare/releases"><img src="https://img.shields.io/github/v/release/smtdfc/nagare?color=blue&style=flat-square" alt="Release" /></a>
    <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-square&logo=go" alt="Go" /></a>
    <a href="https://www.rust-lang.org/"><img src="https://img.shields.io/badge/Rust-dea584?style=flat-square&logo=rust" alt="Rust" /></a>
    <a href="https://react.dev/"><img src="https://img.shields.io/badge/React-19-61DAFB?style=flat-square&logo=react" alt="React 19" /></a>
    <a href="./LICENSE"><img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License" /></a>
  </p>
</p>

---

## What is Nagare?

**Nagare** (流れ - _flow_) is a personal AI assistant built directly for your desktop.

It runs quietly in the background on your computer (Linux, macOS, Windows), helping you automate daily workflows, control system tasks, schedule reminders and cron jobs, and interact seamlessly via an embedded web interface or chat channels like Telegram—all while keeping your conversations, memories, and secrets 100% local on your machine.

---

## Why Nagare?

- 🖐️ **Hands on Your OS:** Nagare isn't just a text bot. Out of the box, it can manage files, check system volume, launch processes, inspect system status, and run scheduled tasks.
- ⚡ **Lightweight Plugin System:** Write extensions in separate processes communicating over fast local IPC (Unix Sockets & Windows Named Pipes). A plugin crash never takes down your assistant.
- 🦀 **Fast Local Memory:** Semantic search and long-term memory are powered by Rust plugin
- 🖥️ **Built-in Web Interface:** Includes a clean, responsive React 19 chat dashboard accessible right from your browser.

---

## Quick Start

### 1. Install Nagare

#### Linux & macOS

```bash
curl -fsSL https://raw.githubusercontent.com/smtdfc/nagare/main/installers/linux/install.sh | sudo bash
```

#### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/smtdfc/nagare/main/installers/window/install.ps1 | iex
```

---

### 2. Start the Gateway Service

Start the Nagare Gateway:

```bash
nagare gateway start
```

To stop it at any time:

```bash
nagare gateway stop
```

---

### 3. Open the Web Dashboard

Launch the web UI:

```bash
nagare web
```

Then visit [http://localhost:3005](http://localhost:3005) to configure your API keys, chat with your assistant, and manage plugins.

---

### 4. Create an Auth Token

Generate a token for local API or script access:

```bash
nagare auth token
```

---

## Creating Custom Tools & Plugins

Nagare is designed to be easily extensible. Plugins run as independent executables (`.nagare_plugin`) and connect to Nagare via local IPC:

```go
package main

import (
	"context"
	"github.com/smtdfc/nagare/plugin/client"
)

func main() {
	plugin := client.NewPlugin()

	// Register a new tool the AI assistant can call
	plugin.RegisterTool(
		"get_battery_status",
		"Check the current battery level of the laptop",
		func(ctx context.Context, args struct{}) (string, error) {
			// Query OS battery status...
			return "Battery is at 85%, discharging", nil
		},
	)

	plugin.ConnectAndServe()
}
```

Build your plugin with the CLI:

```bash
nagare plugin build
```

### Signing plugin archives

Generate an Ed25519 key pair once, then pass the private key to the build command:

```bash
openssl genpkey -algorithm ED25519 -out plugin-signing-key.pem
openssl pkey -in plugin-signing-key.pem -pubout -out plugin-signing-key.pub.pem

nagare plugin build --sign plugin-signing-key.pem
```

The signature is written next to the archive as `<package>.nagare_plugin.sig`. Keep both files together when installing from a local path.

When installing from the web UI or updating an existing plugin, the `.sig` upload field is optional: choose the signature file alongside the package and both are verified together. Without it, the gateway falls back to the `.sig` file next to the package, if any.

The gateway trusts any Ed25519 public key placed in the trusted plugin key folder:

```bash
mkdir -p ~/.nagare/trusted
cp plugin-signing-key.pub.pem ~/.nagare/trusted/publisher.pem
```

Alternatively, copy any number of public key files into that folder. Key files may contain a PKIX PEM public key or a raw base64/hex Ed25519 public key. A signed archive is accepted when it verifies against at least one trusted key. If the trusted folder is empty, the archive installs with an **unverified** warning; an archive without a `.sig` file installs with an **unsigned** warning.

The installed binary checksum is persisted and checked again before every start. A binary modified after installation is refused.

---

## License

Distributed under the **MIT License**. See [`LICENSE`](./LICENSE) for details.
