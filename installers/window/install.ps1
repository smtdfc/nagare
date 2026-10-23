
$ErrorActionPreference = "Stop"

$REPO = "smtdfc/nagare"
$INSTALL_PATH = Join-Path $HOME ".local\bin"

$NAGARE_DATA_DIR = Join-Path $HOME ".nagare"
$NAGARE_ASSETS_DIR = Join-Path $NAGARE_DATA_DIR "assets"
$BUILTIN_PLUGINS_STORE_DIR = Join-Path $NAGARE_ASSETS_DIR "plugins\builtins"
$CORE_PLUGINS_STORE_DIR = Join-Path $NAGARE_ASSETS_DIR "plugins\core"
$WEB_STORE_DIR = Join-Path $NAGARE_ASSETS_DIR "web"

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
        Copy-Item -Path (Join-Path $Source "*") `
            -Destination $Destination -Recurse -Force
    }
}

function Install-Distribution {
    param(
        [Parameter(Mandatory)]
        [string]$Source,

        [Parameter(Mandatory)]
        [string]$Label
    )

    Install-Binary `
        (Join-Path $Source "nagare.exe") `
        (Join-Path $INSTALL_PATH "nagare.exe") `
        "CLI ($Label)"

    $gateway = Join-Path $Source "nagare-gateway.exe"
    if (Test-Path -LiteralPath $gateway -PathType Leaf) {
        Install-Binary $gateway `
            (Join-Path $INSTALL_PATH "nagare-gateway.exe") `
            "Gateway ($Label)"
    }
    else {
        Write-Host "Notice: nagare-gateway.exe not found, skipping."
    }

    $desktop = Join-Path $Source "nagare-desktop.exe"
    if (Test-Path -LiteralPath $desktop -PathType Leaf) {
        Install-Binary $desktop `
            (Join-Path $INSTALL_PATH "nagare-desktop.exe") `
            "Nagare Desktop ($Label)"
    }
    else {
        Write-Host "Notice: nagare-desktop.exe not found, skipping."
    }

    Write-Host "Copying plugins and web assets..."

    Copy-Assets (Join-Path $Source "plugins\builtins") `
        $BUILTIN_PLUGINS_STORE_DIR

    Copy-Assets (Join-Path $Source "plugins\core") `
        $CORE_PLUGINS_STORE_DIR

    # Local build: dist/web
    # Release archive: plugins/web (same layout as the Bash script)
    $webSource = Join-Path $Source "web"
    if (-not (Test-Path -LiteralPath $webSource -PathType Container)) {
        $webSource = Join-Path $Source "plugins\web"
    }

    Copy-Assets $webSource $WEB_STORE_DIR
}

# Detect OS and architecture
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

# 2. Fetch latest GitHub release
$releaseUrl = "https://api.github.com/repos/$REPO/releases/latest"

try {
    $release = Invoke-RestMethod -Uri $releaseUrl -Method Get
    $VERSION = $release.tag_name
}
catch {
    throw "Failed to fetch the latest release from GitHub: $_"
}

if ([string]::IsNullOrWhiteSpace($VERSION)) {
    throw "Failed to determine the latest release version."
}

Write-Host "Selected version for installation: $VERSION"

$ARCHIVE_NAME = "nagare-${OS}-${ARCH}.tar.gz"
$ARCHIVE_URL = "https://github.com/$REPO/releases/download/$VERSION/$ARCHIVE_NAME"

$TMP_DIR = Join-Path ([System.IO.Path]::GetTempPath()) `
    ("nagare-install-" + [guid]::NewGuid().ToString("N"))

New-Item -ItemType Directory -Path $TMP_DIR -Force | Out-Null

try {
    $archivePath = Join-Path $TMP_DIR $ARCHIVE_NAME

    Write-Host "Downloading archive from: $ARCHIVE_URL"
    Invoke-WebRequest -Uri $ARCHIVE_URL -OutFile $archivePath

    Write-Host "Extracting archive..."
    tar.exe -xzf $archivePath -C $TMP_DIR

    if ($LASTEXITCODE -ne 0) {
        throw "Failed to extract release archive."
    }

    # Support either files at archive root or inside a single directory.
    $SOURCE_DIR = $TMP_DIR

    if (-not (Test-Path -LiteralPath (Join-Path $SOURCE_DIR "nagare.exe"))) {
        $candidate = Get-ChildItem -LiteralPath $TMP_DIR -Directory |
            Where-Object {
                Test-Path -LiteralPath (Join-Path $_.FullName "nagare.exe")
            } |
            Select-Object -First 1

        if ($candidate) {
            $SOURCE_DIR = $candidate.FullName
        }
    }

    if (-not (Test-Path -LiteralPath (Join-Path $SOURCE_DIR "nagare.exe"))) {
        throw "nagare.exe was not found in the downloaded archive."
    }

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
    if (Test-Path -LiteralPath $TMP_DIR) {
        Remove-Item -LiteralPath $TMP_DIR -Recurse -Force
    }
}
