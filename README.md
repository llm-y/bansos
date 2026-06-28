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
8. [Kiro Startup](https://startups.aws.com/credits/kiro)
9. [Agent Router](https://agentrouter.org)
10. [Unimodel](https://www.unimodel.ai/dashboard/overview)


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

### Prasyarat

Pastikan sudah melakukan instalasi [Claude Code CLI](https://docs.anthropic.com/en/docs/claude-code):

```bash
npm install -g @anthropic-ai/claude-code
```

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

## Cara Penggunaan

### 1. Daftar dan Generate API Key

Sebelum mulai, kamu perlu mendaftar dan membuat API key:

1. Buka link registrasi: [https://freemodel.dev/invite/FRE-e1e3c216](https://freemodel.dev/invite/FRE-e1e3c216)
2. Siapkan Telegram untuk proses verifikasi akun
3. Upayakan punya **5-10 akun**, masing-masing akun generate **1 API key**
4. Kumpulkan semua API key tersebut untuk dimasukkan ke file CSV di langkah berikutnya

> Semakin banyak API key yang kamu punya, semakin lancar rotasi key saat digunakan.

### 2. Siapkan file `bansos.csv`

Sebelum menjalankan `bansos`, kamu perlu membuat file `bansos.csv` di direktori tempat kamu akan menjalankan perintah. File ini berisi daftar API key dalam format CSV dengan header `id,key`.

ID bisa berupa angka atau huruf (bebas, sesuai keinginan kamu).

Contoh isi `bansos.csv`:

```csv
id,key
anu,sk-ant-api03-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
awg,sk-ant-api03-yyyyyyyyyyyyyyyyyyyyyyyyyyyyyyyy
bukped,sk-ant-api03-zzzzzzzzzzzzzzzzzzzzzzzzzzzzzzzz
```

> Baris pertama harus header `id,key`. Setiap baris berikutnya berisi ID dan API key dipisahkan koma. Pastikan tidak ada spasi tambahan.

### 3. Jalankan perintah

Buka terminal di direktori yang berisi `bansos.csv`, lalu jalankan:

```bash
bansos
```

### 4. Apa yang terjadi

Ketika dijalankan, `bansos` akan:

1. Membaca file `~/.claude/settings.json` untuk mendapatkan API key yang sedang aktif
2. Membaca daftar key dari `bansos.csv` di current directory
3. Mencari key aktif dalam daftar, lalu mengambil key berikutnya (jika sudah di akhir daftar, kembali ke key pertama)
4. Menampilkan ID key saat ini dan ID key yang dipilih berikutnya
5. Menulis `settings.json` baru dengan key yang sudah dirotasi

Contoh output:

```
Using bansos.csv: /home/user/bansos.csv
Current key: ID awg (sk-ant-a...yyyy)
Switched to: ID bukped (sk-ant-a...zzzz)
Successfully wrote settings to: /home/user/.claude/settings.json
```

Dengan begitu, setiap kali kamu menjalankan `bansos`, API key akan berganti secara otomatis ke key berikutnya dalam daftar, dan kamu bisa lihat ID mana yang sedang aktif.
