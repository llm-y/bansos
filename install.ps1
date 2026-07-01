$ErrorActionPreference = "Stop"

# Auto-elevate to Administrator if not already running as admin
$currentPrincipal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
if (-not $currentPrincipal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    Write-Host "Requesting Administrator privileges..."
    $scriptPath = $MyInvocation.MyCommand.Path
    if (-not $scriptPath) {
        $scriptPath = $PSCommandPath
    }
    Start-Process PowerShell -Verb RunAs -ArgumentList "-ExecutionPolicy Bypass -File `"$scriptPath`""
    exit
}

Write-Host "==========================================="
Write-Host "  Bansos Installer (Administrator Mode)"
Write-Host "==========================================="
Write-Host ""

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

Write-Host "[1/7] Detected Architecture: $Arch"

# Check if bansos.exe is currently running and kill it
Write-Host "[2/7] Checking for running bansos.exe processes..."
$bansosProcs = Get-Process -Name "bansos" -ErrorAction SilentlyContinue
if ($bansosProcs) {
    Write-Host "  bansos.exe is currently running. Stopping process..."
    $bansosProcs | Stop-Process -Force
    Write-Host "  Process stopped. Waiting for file handle release..."
    Start-Sleep -Seconds 2
} else {
    Write-Host "  No running bansos.exe process detected."
}

Write-Host "[3/7] Downloading latest release of $AssetName..."

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
Write-Host "  Download URL: $DownloadUrl"

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
        Write-Host "  Added $InstallDir to user PATH. Please restart your terminal."
    }
}

$DestPath = Join-Path $InstallDir $BinaryName

# Download binary to temp file first, then move to destination
Write-Host "[4/7] Installing to $DestPath..."
$TempPath = Join-Path $env:TEMP "bansos_update_$([System.Guid]::NewGuid().ToString('N')).exe"
Invoke-WebRequest -Uri $DownloadUrl -OutFile $TempPath -UseBasicParsing

# Move temp file to destination (retry if file is still locked)
$maxRetries = 3
$retryCount = 0
$installed = $false
while (-not $installed -and $retryCount -lt $maxRetries) {
    try {
        Copy-Item -Path $TempPath -Destination $DestPath -Force -ErrorAction Stop
        $installed = $true
    } catch {
        $retryCount++
        if ($retryCount -lt $maxRetries) {
            Write-Host "  File is still locked, retrying in 2 seconds... (attempt $retryCount/$maxRetries)"
            Start-Sleep -Seconds 2
        } else {
            Remove-Item -Path $TempPath -Force -ErrorAction SilentlyContinue
            Write-Error "Error: Could not write to $DestPath after $maxRetries attempts. Please close bansos.exe and try again."
            exit 1
        }
    }
}
Remove-Item -Path $TempPath -Force -ErrorAction SilentlyContinue

# Bypass Smart App Control & Unblock binary
Write-Host "[5/7] Bypassing Smart App Control & unblocking binary..."
try {
    # Remove Mark of the Web (MOTW) so Windows does not flag as untrusted
    Unblock-File -Path "$DestPath" -ErrorAction Stop
    Write-Host "  Unblock-File applied: Mark of the Web removed."
} catch {
    Write-Warning "  Unblock-File failed: $_"
}
try {
    # Fallback: manually remove Zone.Identifier alternate data stream
    $zoneFile = "${DestPath}:Zone.Identifier"
    if (Test-Path $zoneFile) {
        Remove-Item -Path $zoneFile -Force -ErrorAction Stop
        Write-Host "  Zone.Identifier ADS removed manually."
    } else {
        Write-Host "  Zone.Identifier ADS not present (already clean)."
    }
} catch {
    Write-Warning "  Could not remove Zone.Identifier ADS: $_"
}
try {
    # Add process exclusion so Smart App Control/Defender won't block execution
    Add-MpPreference -ExclusionProcess "$BinaryName" -ErrorAction Stop
    Write-Host "  Defender ExclusionProcess added for: $BinaryName"
} catch {
    Write-Warning "  Could not add ExclusionProcess: $_"
    Write-Warning "  You may need to manually add $BinaryName as a process exclusion in Windows Security."
}

# Add Windows Firewall rule to allow bansos outbound network access
Write-Host "[6/7] Configuring Windows Firewall..."
try {
    # Remove existing rule if present (to avoid duplicates)
    $existingRule = Get-NetFirewallRule -DisplayName "Allow Bansos Outbound" -ErrorAction SilentlyContinue
    if ($existingRule) {
        Remove-NetFirewallRule -DisplayName "Allow Bansos Outbound" -ErrorAction SilentlyContinue
    }
    # Add firewall rule allowing outbound traffic for bansos
    New-NetFirewallRule -DisplayName "Allow Bansos Outbound" `
        -Direction Outbound `
        -Action Allow `
        -Program "$DestPath" `
        -Protocol TCP `
        -Enabled True `
        -Profile Any `
        -Description "Allow bansos application to send data over the network" | Out-Null
    Write-Host "  Firewall rule added: bansos.exe is allowed outbound network access."
} catch {
    Write-Warning "  Could not configure firewall rule: $_"
    Write-Warning "  You may need to manually allow bansos.exe through Windows Firewall."
}

# Add Windows Defender exclusion for the binary path
Write-Host "[7/7] Adding Windows Defender exclusion..."
try {
    Add-MpPreference -ExclusionPath "$DestPath" -ErrorAction Stop
    Write-Host "  Windows Defender exclusion added for: $DestPath"
} catch {
    Write-Warning "  Could not add Defender exclusion: $_"
    Write-Warning "  You may need to manually add an exclusion for $DestPath in Windows Security."
}

Write-Host ""
Write-Host "==========================================="
Write-Host "  Installation complete!"
Write-Host "==========================================="
Write-Host "  Binary: $DestPath"
Write-Host "  Run 'bansos' to get started."
Write-Host ""
Write-Host "Jalankan script ini kapan saja untuk update ke versi terbaru."
Write-Host ""
Read-Host "Press Enter to exit..."
