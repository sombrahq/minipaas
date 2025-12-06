---
title: Install MiniPaaS CLI
---

## Installation Options

You can install **MiniPaaS CLI** using one of the following methods:

---

### 🧑‍💻 Option 1: Install via Go

If you have Go 1.21+ installed, use:

```bash
go install github.com/sombrahq/minipaas/minipaas-cli/cmd/minipaas@main
```

> This places the `minipaas` binary in `$(go env GOPATH)/bin`.
> Make sure that directory is in your system `PATH`:

```bash
export PATH="$(go env GOPATH)/bin:$PATH"
```

---

### 📦 Option 2: Download Prebuilt Binaries

1. Go to the [GitHub Releases](https://github.com/sombrahq/minipaas/releases).

2. Download the binary for your platform:

	* `minipaas-linux-amd64`
	* `minipaas-darwin-arm64`
	* `minipaas-windows-amd64.exe`

3. Make it executable (Linux/macOS):

```bash
chmod +x minipaas-*
mv minipaas-* /usr/local/bin/minipaas
```

4. Confirm installation:

```bash
minipaas --help
```

---

### 🛠 Option 3: Build from Source

If you prefer to build manually:

```bash
git clone https://github.com/sombrahq/minipaas.git
cd minipaas/minipaas-cli/cmd/minipaas
go mod tidy
go build
```

Output is saved in the `build/` directory.

---

## Requirements

* Go **1.21+**
* Docker CLI and OpenSSL (for TLS certs)
* Unix-like shell or terminal (Linux, macOS, or Windows PowerShell)

---

Need help? Open an [issue](https://github.com/sombrahq/minipaas/issues) or [contact me](../contact.md).
