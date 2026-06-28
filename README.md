# bansos
info bansos ai
1. [Threads](https://www.threads.com/search?q=bansos&serp_type=default)
2. [x ](https://x.com/search?q=bansos%20claude&src=typed_query&f=live)
3. [Bansos AI](https://appverse.id/bansos-ai)
4. [Bansos Dev](https://bansos.dev/list/)

## Promo Provider
1. [Command Code](https://commandcode.ai/pricing)
2. [Mimo](https://mimo.mi.com/)
3. [Kiro](https://kiro.dev/pricing/)
4. [bobtrial+](https://bob.ibm.com/pricing)
5. Claude Code : [FreeModel](https://freemodel.dev/invite/FRE-e1e3c216) , [aerolink](https://aerolink.lat/register?ref=QE-9CMM)
6. [Nvidia](https://build.nvidia.com/models)
7. [Khusus OSS](https://claude.com/contact-sales/claude-for-oss)
8. [Agent Router](https://agentrouter.org)
9. [Unimodel](https://www.unimodel.ai/dashboard/overview)


C:\Users\<you>\.claude\settings.json

Base URL:
1. https://cc.freemodel.dev
2. https://capi.aerolink.lat/
3. https://agentrouter.org/

```json
{
    "env": {
        "ANTHROPIC_API_KEY": "YOUR_API_KEY",
        "ANTHROPIC_BASE_URL": "https://cc.freemodel.dev",
        "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1"
    },
    "permissions": {
        "allow": ["Bash", "Write", "Edit"],
        "deny": []
    },
    "apiKeyHelper": "echo 'YOUR_API_KEY'"
}
```

```json
{
    "env": {
        "ANTHROPIC_AUTH_TOKEN": "YOUR_API_KEY",
        "ANTHROPIC_API_KEY": "YOUR_API_KEY",
        "ANTHROPIC_BASE_URL": "https://agentrouter.org/",
        "CLAUDE_CODE_DISABLE_NONESSENTIAL_TRAFFIC": "1"
    },
    "permissions": {
        "allow": ["Bash", "Write", "Edit"],
        "deny": []
    }
}
```


## Broker Termurah
1. [Vietnam](https://t.me/tai_khoan_ai_bot)


## Installation

### Linux / macOS (Terminal)

Jalankan perintah berikut di terminal:

```bash
curl -fsSL https://raw.githubusercontent.com/llm-y/bansos/main/install.sh | bash
```

Script ini akan:
1. Mendeteksi OS dan arsitektur (amd64/arm64) secara otomatis
2. Download binary terbaru dari GitHub Releases
3. Rename menjadi `bansos` dan install ke `/usr/local/bin`
4. Siap dijalankan dari terminal manapun

### Windows (PowerShell)

Jalankan perintah berikut di PowerShell (Run as Administrator jika diperlukan):

```powershell
irm https://raw.githubusercontent.com/llm-y/bansos/main/install.ps1 | iex
```

Script ini akan:
1. Mendeteksi arsitektur (amd64/arm64) secara otomatis
2. Download binary terbaru dari GitHub Releases
3. Rename menjadi `bansos.exe` dan install ke folder yang ada di PATH
4. Siap dijalankan dari PowerShell atau Command Prompt manapun

### Verifikasi Instalasi

Setelah install, buka terminal/PowerShell baru dan jalankan:

```bash
bansos
```
