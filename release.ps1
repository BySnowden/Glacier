# Glacier Release Management Script
# This script helps automate the release process for Glacier Mod Manager

param(
    [Parameter(Mandatory=$true)]
    [string]$Version,

    [Parameter(Mandatory=$false)]
    [string]$Message = "Release version $Version",

    [Parameter(Mandatory=$false)]
    [switch]$DryRun,

    [Parameter(Mandatory=$false)]
    [switch]$Force
)

# Color functions for better output
function Write-Success {
    param([string]$Text)
    Write-Host $Text -ForegroundColor Green
}

function Write-Info {
    param([string]$Text)
    Write-Host $Text -ForegroundColor Cyan
}

function Write-Warning {
    param([string]$Text)
    Write-Host $Text -ForegroundColor Yellow
}

function Write-Error {
    param([string]$Text)
    Write-Host $Text -ForegroundColor Red
}

function Write-Step {
    param([string]$Text)
    Write-Host "`n=== $Text ===" -ForegroundColor Magenta
}

# Validate version format
if (-not ($Version -match '^v?\d+\.\d+\.\d+$')) {
    Write-Error "Invalid version format. Please use format: 1.2.3 or v1.2.3"
    exit 1
}

# Normalize version (ensure it starts with 'v')
$NormalizedVersion = if ($Version.StartsWith('v')) { $Version } else { "v$Version" }
$NumericVersion = $NormalizedVersion.TrimStart('v')

Write-Info "Glacier Release Manager"
Write-Info "======================"
Write-Info "Version: $NormalizedVersion"
Write-Info "Message: $Message"
if ($DryRun) { Write-Warning "DRY RUN MODE - No changes will be made" }
Write-Info ""

# Check if we're in the correct directory
if (-not (Test-Path "app.go")) {
    Write-Error "Error: app.go not found. Please run this script from the Glacier project root."
    exit 1
}

if (-not (Test-Path "wails.json")) {
    Write-Error "Error: wails.json not found. Please run this script from the Glacier project root."
    exit 1
}

Write-Step "Step 1: Checking Git Status"

# Check for uncommitted changes
$gitStatus = git status --porcelain
if ($gitStatus -and -not $Force) {
    Write-Error "Error: There are uncommitted changes in the repository."
    Write-Info "Uncommitted files:"
    $gitStatus | ForEach-Object { Write-Info "  $_" }
    Write-Info "Please commit or stash changes, or use -Force to proceed anyway."
    exit 1
}

if ($gitStatus -and $Force) {
    Write-Warning "Warning: Proceeding with uncommitted changes due to -Force flag"
}

# Check if tag already exists
$existingTag = git tag -l $NormalizedVersion
if ($existingTag -and -not $Force) {
    Write-Error "Error: Tag $NormalizedVersion already exists."
    Write-Info "Use -Force to overwrite existing tag, or choose a different version."
    exit 1
}

Write-Step "Step 2: Updating Version in Source Files"

# Update app.go version
$appGoPath = "app.go"
$appGoContent = Get-Content $appGoPath -Raw

# Find and update the AppVersion constant
$versionPattern = '(?<=const AppVersion = ")[\d.]+(?=")'
if ($appGoContent -match $versionPattern) {
    $currentVersion = $matches[0]
    Write-Info "Current version in app.go: $currentVersion"
    Write-Info "Updating to: $NumericVersion"

    if (-not $DryRun) {
        $newContent = $appGoContent -replace $versionPattern, $NumericVersion
        Set-Content -Path $appGoPath -Value $newContent -NoNewline
        Write-Success "Updated app.go"
    } else {
        Write-Info "Would update app.go version from $currentVersion to $NumericVersion"
    }
} else {
    Write-Error "Error: Could not find AppVersion constant in app.go"
    exit 1
}

# Update wails.json version
$wailsJsonPath = "wails.json"
$wailsJson = Get-Content $wailsJsonPath | ConvertFrom-Json
$currentWailsVersion = $wailsJson.info.productVersion
Write-Info "Current version in wails.json: $currentWailsVersion"
Write-Info "Updating to: $NumericVersion"

if (-not $DryRun) {
    $wailsJson.info.productVersion = $NumericVersion
    $wailsJson | ConvertTo-Json -Depth 10 | Set-Content $wailsJsonPath
    Write-Success "Updated wails.json"
} else {
    Write-Info "Would update wails.json version from $currentWailsVersion to $NumericVersion"
}

Write-Step "Step 3: Building Application"

if (-not $DryRun) {
    Write-Info "Running: wails build"
    $buildResult = & wails build 2>&1
    if ($LASTEXITCODE -ne 0) {
        Write-Error "Build failed:"
        $buildResult | ForEach-Object { Write-Error "  $_" }
        exit 1
    }
    Write-Success "Build completed successfully"
} else {
    Write-Info "Would run: wails build"
}

Write-Step "Step 4: Git Operations"

if (-not $DryRun) {
    # Stage version files
    Write-Info "Staging version changes..."
    git add $appGoPath $wailsJsonPath

    # Commit version changes
    Write-Info "Committing version changes..."
    git commit -m "Bump version to $NormalizedVersion"

    # Create or update tag
    if ($existingTag) {
        Write-Info "Updating existing tag $NormalizedVersion..."
        git tag -f $NormalizedVersion -m $Message
    } else {
        Write-Info "Creating tag $NormalizedVersion..."
        git tag $NormalizedVersion -m $Message
    }

    Write-Success "Git operations completed"
} else {
    Write-Info "Would stage: $appGoPath, $wailsJsonPath"
    Write-Info "Would commit: Bump version to $NormalizedVersion"
    if ($existingTag) {
        Write-Info "Would update tag: $NormalizedVersion"
    } else {
        Write-Info "Would create tag: $NormalizedVersion"
    }
}

Write-Step "Step 5: Push to Repository"

if (-not $DryRun) {
    Write-Info "Pushing commits to origin..."
    git push origin

    Write-Info "Pushing tag to origin..."
    if ($existingTag) {
        git push origin $NormalizedVersion --force
    } else {
        git push origin $NormalizedVersion
    }

    Write-Success "Pushed to repository"
} else {
    Write-Info "Would push commits to origin"
    Write-Info "Would push tag $NormalizedVersion to origin"
}

Write-Step "Release Complete!"

Write-Success "Successfully released Glacier Mod Manager $NormalizedVersion"
Write-Info ""
Write-Info "Next steps:"
Write-Info "1. GitHub Actions will automatically build and create the release"
Write-Info "2. Monitor the Actions tab: https://github.com/BySnowden/Glacier/actions"
Write-Info "3. Release will be available at: https://github.com/BySnowden/Glacier/releases/tag/$NormalizedVersion"
Write-Info ""
Write-Info "The auto-update system will now detect this new version!"

# Show final status
Write-Step "Summary"
Write-Success "✓ Updated version in source files"
Write-Success "✓ Built application successfully"
Write-Success "✓ Created git commit and tag"
Write-Success "✓ Pushed to repository"
Write-Info "🚀 GitHub Actions will handle the rest!"
