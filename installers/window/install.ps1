$ErrorActionPreference = "Stop"

$Repo = "smtdfc/nagare"
$InstallPath = "$env:ProgramFiles\Nagare"

function Install-Binary {
    param (
        [string]$Src,
        [string]$DestName,
        [string]$Label
    )
    $DestPath = Join-Path $InstallPath $DestName
    try {
        if (-not (Test-Path $InstallPath)) {
            New-Item -ItemType Directory -Force -Path $InstallPath | Out-Null
        }
        Copy-Item -Path $Src -Destination $DestPath -Force
        Write-Host "Successfully installed $Label: $DestName" -ForegroundColor Green
        return $true
    }
    catch {
        Write-Host "Error: Failed to install $Label to $DestPath. $_" -ForegroundColor Red
        return $false
    }
}


$Arch = (Get-WmiObject Win32_OperatingSystem).OSArchitecture
if ([Environment]::Is64BitOperatingSystem) {
    $ArchTag = "amd64"
    if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") {
        $ArchTag = "arm64"
    }
} else {
    Write-Host "Unsupported architecture (32-bit is not supported)." -ForegroundColor Red
    exit 1
}

$OS = "windows"
Write-Host "Detected system: OS=$OS, Arch=$ArchTag"

$LocalDist = "dist\${OS}-${ArchTag}"

if ((Test-Path $LocalDist) -and (Test-Path "$LocalDist\nagare.exe")) {
    Write-Host "Found local build in $LocalDist. Installing from local files..."
    
    Install-Binary -Src "$LocalDist\nagare.exe" -DestName "nagare.exe" -Label "CLI (local)"
    if (Test-Path "$LocalDist\nagare-gateway.exe") {
        Install-Binary -Src "$LocalDist\nagare-gateway.exe" -DestName "nagare-gateway.exe" -Label "Gateway (local)"
    } else {
        Write-Host "Notice: Local Gateway binary not found in $LocalDist, skipping." -ForegroundColor Yellow
    }

    $UserPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
    if ($UserPath -notlike "*$InstallPath*") {
        [Environment]::SetEnvironmentVariable("Path", "$UserPath;$InstallPath", "Machine")
        Write-Host "Added $InstallPath to system PATH." -ForegroundColor Cyan
    }

    Write-Host "=========================================" -ForegroundColor Cyan
    Write-Host " Nagare installation completed!" -ForegroundColor Cyan
    Write-Host "=========================================" -ForegroundColor Cyan
    exit 0
}

Write-Host "Fetching latest release version from GitHub..." -ForegroundColor Cyan
try {
    $ReleaseInfo = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest"
    $Version = $ReleaseInfo.tag_name
}
catch {
    Write-Host "Failed to fetch version from GitHub. Please check your network connection." -ForegroundColor Red
    exit 1
}

if (-not $Version) {
    Write-Host "Failed to parse version from GitHub response." -ForegroundColor Red
    exit 1
}

Write-Host "Selected version for installation: $Version"

$ArchiveName = "nagare-${OS}-${ArchTag}.tar.gz"
$BaseUrl = "https://github.com/$Repo/releases/download/$Version"
$ArchiveUrl = "$BaseUrl/$ArchiveName"

$TmpDir = New-TemporaryFile | Remove-Item -PassThru | New-Item -ItemType Directory

try {
    $ArchivePath = Join-Path $TmpDir.FullName $ArchiveName
    Write-Host "Downloading archive from: $ArchiveUrl" -ForegroundColor Cyan
    Invoke-WebRequest -Uri $ArchiveUrl -OutFile $ArchivePath

    Write-Host "Extracting archive..."
    tar -xzf $ArchivePath -C $TmpDir.FullName

    $ExtractedNagare = Join-Path $TmpDir.FullName "nagare.exe"
    $ExtractedGateway = Join-Path $TmpDir.FullName "nagare-gateway.exe"

    if (Test-Path $ExtractedNagare) {
        Install-Binary -Src $ExtractedNagare -DestName "nagare.exe" -Label "CLI"
    } else {
        Write-Host "Warning: 'nagare.exe' binary not found in the downloaded archive." -ForegroundColor Yellow
    }

    if (Test-Path $ExtractedGateway) {
        Install-Binary -Src $ExtractedGateway -DestName "nagare-gateway.exe" -Label "Gateway"
    } else {
        Write-Host "Notice: 'nagare-gateway' binary not found in the archive, skipping." -ForegroundColor Yellow
    }
}
catch {
    Write-Host "Error: Could not download or extract release archive. $_" -ForegroundColor Red
    exit 1
}
finally {
    Remove-Item -Path $TmpDir.FullName -Recurse -Force -ErrorAction SilentlyContinue
}

$MachinePath = [Environment]::GetEnvironmentVariable("Path", "Machine")
if ($MachinePath -notlike "*$InstallPath*") {
    [Environment]::SetEnvironmentVariable("Path", "$MachinePath;$InstallPath", "Machine")
    Write-Host "Added $InstallPath to system PATH." -ForegroundColor Cyan
}

Write-Host "=========================================" -ForegroundColor Cyan
Write-Host " Nagare installation completed!" -ForegroundColor Cyan
Write-Host "=========================================" -ForegroundColor Cyan