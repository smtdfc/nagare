<p align="center">
  <h1 align="center">🌊 Nagare</h1>
  <p align="center">
    <strong>A lightweight, local-first autonomous AI daemon and plugin runtime for desktop automation.</strong>
  </p>
  <p align="center">
    <a href="https://github.com/smtdfc/nagare/releases"><img src="https://img.shields.io/github/v/release/smtdfc/nagare?color=blue&style=flat-square" alt="Release" /></a>
    <a href="https://go.dev/"><img src="https://img.shields.io/badge/Go-1.26+-00ADD8?style=flat-square&logo=go" alt="Go" /></a>
    <a href="https://www.rust-lang.org/"><img src="https://img.shields.io/badge/Rust-SIMD_FFI-dea584?style=flat-square&logo=rust" alt="Rust" /></a>
    <a href="https://react.dev/"><img src="https://img.shields.io/badge/React-19-61DAFB?style=flat-square&logo=react" alt="React 19" /></a>
    <a href="./LICENSE"><img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License" /></a>
  </p>
</p>

---

## What is Nagare?

**Nagare** (流れ - _flow_) is an extensible, privacy-preserving AI assistant daemon engineered to bridge large language models with your local operating system.

Unlike heavy cloud-centric orchestrators or bloated multi-process frameworks, Nagare runs as a compact, self-contained background service on Linux, macOS, and Windows. It gives autonomous agents real hands on your machine: executing system tasks, scheduling background jobs, managing memory through hardware-accelerated vector storage, and extending capabilities through decoupled out-of-process plugins.

---

## Why Nagare?

- 🔒 **Local-First & Private:** Your session history, vector memory, credentials, and execution logs live on your device. Zero telemetry, zero cloud lock-in.
- ⚡ **Microkernel Plugin Runtime:** Plugins run as isolated, separate processes communicating over low-latency binary IPC (Unix Domain Sockets & Windows Named Pipes). A crashed plugin will never bring down the host.
- 🦀 **PureGo Rust FFI:** Vector search and memory quantization are powered by native Rust SIMD kernels, dynamically linked without the cross-compilation baggage of CGO.
- ⚙️ **System-Level Control:** Out of the box, Nagare can manage processes, inspect user directories, manipulate system volume, execute scheduled cron jobs, and interface with external platforms like Telegram and Google Calendar.
- 🌐 **Modern Embedded Web Dashboard:** Includes a clean, real-time React 19 UI served straight from the CLI with bidirectional token streaming.

---

## Core Capabilities

### Autonomous ReAct Execution Loop

Nagare features a resilient **ReAct (Reasoning + Acting)** execution engine. The agent dynamically discovers tool categories, binds appropriate handlers, evaluates multi-turn tool outputs, and automatically recovers from model quirks—all backed by strict iteration limits and graceful `context.Context` cancellation.

### Out-of-Process Plugin Ecosystem

Plugins are compiled standalone executables (`.nagare_plugin`) that integrate via a length-prefixed binary framing protocol. Plugins can inject dynamic tools, stream messages, and bridge external event streams (e.g., chat services, notification webhooks) into Nagare's unified core event bus.

### Hardware-Accelerated Vector Memory

Built for fast semantic memory retrieval without spinning up heavyweight vector databases. Nagare embeds a custom native library built on top of Rust's `turbovec`—offering 2–4 bit quantization and SIMD-accelerated similarity searches with near-zero RAM overhead.

---

## Quick Start

### 1. Installation

#### Linux & macOS

```bash
curl -fsSL https://raw.githubusercontent.com/smtdfc/nagare/main/installers/linux/install.sh | sudo bash
```

#### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/smtdfc/nagare/main/installers/window/install.ps1 | iex
```

---

### 2. Running the Gateway Daemon

Start the Nagare Gateway in the background:

```bash
nagare gateway start -d
```

To stop the background gateway:

```bash
nagare gateway stop
```

---

### 3. Launching the Web Dashboard

Nagare comes with an embedded web interface:

```bash
nagare web
```

Open [http://localhost:3005](http://localhost:3005) in your browser to configure your LLM providers (OpenAI-compatible endpoints), manage chat sessions, and monitor active plugins.

---

### 4. Generating Client Credentials

Generate an authentication token for external API or CLI access:

```bash
nagare auth token
```

---

## Developing Plugins

Nagare makes building extensions straightforward. Because plugins communicate via IPC, you can write them using the provided Go SDK or implement the framing protocol in any language.

```go
package main

import (
	"context"
	"github.com/smtdfc/nagare/plugin/client"
)

func main() {
	plugin := client.NewPlugin("weather-extension", "1.0.0")

	plugin.RegisterTool(
		"get_forecast",
		"Fetch real-time weather information for a specific city",
		func(ctx context.Context, args ForecastInput) (ForecastOutput, error) {
			// Implementation logic
			return ForecastOutput{Temp: 24, Condition: "Sunny"}, nil
		},
	)

	plugin.ConnectAndServe()
}
```

Compile the plugin using the CLI:

```bash
nagare plugin build
```

---

## License

Distributed under the **MIT License**. See [`LICENSE`](./LICENSE) for more information.

Maintained with craft by [@smtdfc](https://github.com/smtdfc).
