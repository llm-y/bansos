$ErrorActionPreference = "Stop"

$Repo = "llm-y/bansos"
$BinaryName = "bansos.exe"

# Detect Architecture
$Arch = if ([System.Environment]::Is64BitOperatingSystem) {
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64" -or $env:PROCESSOR_ARCHITEW6432 -eq "ARM64") {
        "arm64"
    } else {
        "amd64"
    }
} else {
    Write-Error "Unsupported: 32-bit OS is not supported."
    exit 1
}

$AssetName = "bansos-windows-${Arch}.exe"

Write-Host "Detected Architecture: $Arch"
Write-Host "Downloading latest release of $AssetName..."

# Get latest release info
$ReleaseUrl = "https://api.github.com/repos/$Repo/releases/latest"
$Release = Invoke-RestMethod -Uri $ReleaseUrl -Headers @{ "User-Agent" = "bansos-installer" }

# Find the download URL for the asset
$Asset = $Release.assets | Where-Object { $_.name -eq $AssetName }

if (-not $Asset) {
    Write-Error "Error: Could not find release asset for $AssetName"
    exit 1
}

$DownloadUrl = $Asset.browser_download_url
Write-Host "Download URL: $DownloadUrl"

# Determine install directory
$InstallDir = Join-Path $env:LOCALAPPDATA "Microsoft\WindowsApps"
if (-not (Test-Path $InstallDir)) {
    $InstallDir = Join-Path $env:USERPROFILE ".local\bin"
    if (-not (Test-Path $InstallDir)) {
        New-Item -ItemType Directory -Path $InstallDir -Force | Out-Null
    }
    # Check if directory is in PATH
    $UserPath = [Environment]::GetEnvironmentVariable("PATH", "User")
    if ($UserPath -notlike "*$InstallDir*") {
        [Environment]::SetEnvironmentVariable("PATH", "$UserPath;$InstallDir", "User")
        Write-Host "Added $InstallDir to user PATH. Please restart your terminal."
    }
}

$DestPath = Join-Path $InstallDir $BinaryName

# Download binary
Write-Host "Installing to $DestPath..."
Invoke-WebRequest -Uri $DownloadUrl -OutFile $DestPath -UseBasicParsing

Write-Host "Successfully installed $BinaryName to $DestPath"
Write-Host "Run 'bansos' to get started."
