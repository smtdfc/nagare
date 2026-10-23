
$ErrorActionPreference = "Stop"

$REPO = "smtdfc/nagare"
$INSTALL_PATH = Join-Path $HOME ".local\bin"

$NAGARE_DATA_DIR = Join-Path $HOME ".nagare"
$NAGARE_ASSETS_DIR = Join-Path $NAGARE_DATA_DIR "assets"
$BUILTIN_PLUGINS_STORE_DIR = Join-Path $NAGARE_ASSETS_DIR "plugins\builtins"
$CORE_PLUGINS_STORE_DIR = Join-Path $NAGARE_ASSETS_DIR "plugins\core"
$WEB_STORE_DIR = Join-Path $NAGARE_ASSETS_DIR "web"

# Create installation directories
@(
    $INSTALL_PATH,
    $BUILTIN_PLUGINS_STORE_DIR,
    $CORE_PLUGINS_STORE_DIR,
    $WEB_STORE_DIR
) | ForEach-Object {
    New-Item -ItemType Directory -Path $_ -Force | Out-Null
}

function Install-Binary {
    param(
        [Parameter(Mandatory)]
        [string]$Source,

        [Parameter(Mandatory)]
        [string]$Destination,

        [Parameter(Mandatory)]
        [string]$Name
    )

    try {
        if (-not (Test-Path -LiteralPath $Source -PathType Leaf)) {
            throw "Source file does not exist: $Source"
        }

        Copy-Item -LiteralPath $Source -Destination $Destination -Force

        Write-Host "Successfully installed ${Name}: $(Split-Path $Destination -Leaf)"
    }
    catch {
        Write-Error "Failed to install ${Name}: $_"
        throw
    }
}

function Copy-Assets {
    param(
        [Parameter(Mandatory)]
        [string]$Source,

        [Parameter(Mandatory)]
        [string]$Destination
    )

    if (Test-Path -LiteralPath $Source -PathType Container) {
        New-Item -ItemType Directory -Path $Destination -Force | Out-Null

        Copy-Item `
            -Path (Join-Path $Source "*") `
            -Destination $Destination `
            -Recurse `
            -Force

        Write-Host "Copied assets: $Source -> $Destination"
    }
    else {
        Write-Host "Notice: asset directory not found, skipping: $Source"
    }
}

function Install-Distribution {
    param(
        [Parameter(Mandatory)]
        [string]$Source,

        [Parameter(Mandatory)]
        [string]$Label
    )

    # CLI
    Install-Binary `
        -Source (Join-Path $Source "nagare.exe") `
        -Destination (Join-Path $INSTALL_PATH "nagare.exe") `
        -Name "CLI ($Label)"

    # Gateway
    $gateway = Join-Path $Source "nagare-gateway.exe"

    if (Test-Path -LiteralPath $gateway -PathType Leaf) {
        Install-Binary `
            -Source $gateway `
            -Destination (Join-Path $INSTALL_PATH "nagare-gateway.exe") `
            -Name "Gateway ($Label)"
    }
    else {
        Write-Host "Notice: nagare-gateway.exe not found, skipping."
    }

    # Desktop
    $desktop = Join-Path $Source "nagare-desktop.exe"

    if (Test-Path -LiteralPath $desktop -PathType Leaf) {
        Install-Binary `
            -Source $desktop `
            -Destination (Join-Path $INSTALL_PATH "nagare-desktop.exe") `
            -Name "Nagare Desktop ($Label)"
    }
    else {
        Write-Host "Notice: nagare-desktop.exe not found, skipping."
    }

    # Plugins and web assets
    Write-Host "Copying plugins and web assets..."

    Copy-Assets `
        -Source (Join-Path $Source "plugins\builtins") `
        -Destination $BUILTIN_PLUGINS_STORE_DIR

    Copy-Assets `
        -Source (Join-Path $Source "plugins\core") `
        -Destination $CORE_PLUGINS_STORE_DIR

    # Local build: dist/web
    # Release archive: web or plugins/web
    $webSource = Join-Path $Source "web"

    if (-not (Test-Path -LiteralPath $webSource -PathType Container)) {
        $webSource = Join-Path $Source "plugins\web"
    }

    Copy-Assets `
        -Source $webSource `
        -Destination $WEB_STORE_DIR
}

# Detect operating system and architecture
$OS = "windows"

switch ([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture.ToString()) {
    "X64" {
        $ARCH = "amd64"
    }

    "Arm64" {
        $ARCH = "arm64"
    }

    default {
        throw "Unsupported CPU architecture: $([System.Runtime.InteropServices.RuntimeInformation]::OSArchitecture)"
    }
}

Write-Host "Detected system: OS=$OS, Arch=$ARCH"

# 1. Prefer local build
$LOCAL_DIST = Join-Path (Get-Location) "dist"
$LOCAL_CLI = Join-Path $LOCAL_DIST "nagare.exe"

if (Test-Path -LiteralPath $LOCAL_CLI -PathType Leaf) {
    Write-Host "Found local build in dist. Installing from local files..."

    Install-Distribution -Source $LOCAL_DIST -Label "local"

    Write-Host "========================================="
    Write-Host " Nagare local installation completed!"
    Write-Host "========================================="

    exit 0
}

# 2. Fetch the latest GitHub Release
$releaseUrl = "https://api.github.com/repos/$REPO/releases/latest"

try {
    Write-Host "Fetching latest release from GitHub..."

    $release = Invoke-RestMethod `
        -Uri $releaseUrl `
        -Method Get

    $VERSION = $release.tag_name
}
catch {
    throw "Failed to fetch the latest release from GitHub: $_"
}

if ([string]::IsNullOrWhiteSpace($VERSION)) {
    throw "Failed to determine the latest release version."
}

Write-Host "Selected version for installation: $VERSION"

# Download ZIP release archive
$ARCHIVE_NAME = "nagare-${OS}-${ARCH}.zip"
$ARCHIVE_URL = "https://github.com/$REPO/releases/download/$VERSION/$ARCHIVE_NAME"

$TMP_DIR = Join-Path `
    ([System.IO.Path]::GetTempPath()) `
    ("nagare-install-" + [guid]::NewGuid().ToString("N"))

New-Item -ItemType Directory -Path $TMP_DIR -Force | Out-Null

try {
    $archivePath = Join-Path $TMP_DIR $ARCHIVE_NAME

    Write-Host "Downloading archive from: $ARCHIVE_URL"

    Invoke-WebRequest `
        -Uri $ARCHIVE_URL `
        -OutFile $archivePath

    if (-not (Test-Path -LiteralPath $archivePath -PathType Leaf)) {
        throw "Downloaded archive was not found: $archivePath"
    }

    if ((Get-Item -LiteralPath $archivePath).Length -eq 0) {
        throw "Downloaded archive is empty."
    }

    # Extract ZIP archive
    Write-Host "Extracting archive..."

    Expand-Archive `
        -LiteralPath $archivePath `
        -DestinationPath $TMP_DIR `
        -Force

    # Support files at the archive root or inside one directory
    $SOURCE_DIR = $TMP_DIR

    if (-not (Test-Path -LiteralPath (Join-Path $SOURCE_DIR "nagare.exe") -PathType Leaf)) {
        $candidate = Get-ChildItem -LiteralPath $TMP_DIR -Directory |
            Where-Object {
                Test-Path -LiteralPath (Join-Path $_.FullName "nagare.exe") -PathType Leaf
            } |
            Select-Object -First 1

        if ($candidate) {
            $SOURCE_DIR = $candidate.FullName
        }
    }

    if (-not (Test-Path -LiteralPath (Join-Path $SOURCE_DIR "nagare.exe") -PathType Leaf)) {
        throw "nagare.exe was not found in the downloaded archive."
    }

    # Install binaries and assets
    Install-Distribution -Source $SOURCE_DIR -Label "release"

    Write-Host "========================================="
    Write-Host " Nagare installation completed!"
    Write-Host "========================================="
}
catch {
    Write-Error "Installation failed: $_"
    exit 1
}
finally {
    # Clean up temporary files
    if (Test-Path -LiteralPath $TMP_DIR) {
        Remove-Item `
            -LiteralPath $TMP_DIR `
            -Recurse `
            -Force
    }
}
