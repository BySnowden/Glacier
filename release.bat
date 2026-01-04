@echo off
setlocal enabledelayedexpansion

REM Glacier Release Management Script (Batch Wrapper)
REM This batch file provides a simple interface to the PowerShell release script

echo.
echo ===============================================
echo   Glacier Mod Manager - Release Tool
echo ===============================================
echo.

REM Check if PowerShell is available
powershell -Command "Get-Host" >nul 2>&1
if errorlevel 1 (
    echo ERROR: PowerShell is not available or not in PATH
    echo This script requires PowerShell to run the release process.
    pause
    exit /b 1
)

REM Check if we're in the right directory
if not exist "app.go" (
    echo ERROR: app.go not found.
    echo Please run this script from the Glacier project root directory.
    pause
    exit /b 1
)

if not exist "release.ps1" (
    echo ERROR: release.ps1 not found.
    echo The PowerShell release script is missing.
    pause
    exit /b 1
)

REM Get version from user
set /p VERSION="Enter version number (e.g., 1.2.3): "

if "!VERSION!"=="" (
    echo ERROR: Version number cannot be empty.
    pause
    exit /b 1
)

REM Get optional release message
set /p MESSAGE="Enter release message (press Enter for default): "

if "!MESSAGE!"=="" (
    set MESSAGE=Release version !VERSION!
)

REM Ask for confirmation
echo.
echo Release Configuration:
echo =====================
echo Version: !VERSION!
echo Message: !MESSAGE!
echo.
set /p CONFIRM="Proceed with release? (Y/N): "

if /i not "!CONFIRM!"=="Y" (
    if /i not "!CONFIRM!"=="YES" (
        echo Release cancelled.
        pause
        exit /b 0
    )
)

REM Check for dry run option
echo.
echo Options:
echo 1. Normal release
echo 2. Dry run (preview only)
echo.
set /p OPTION="Choose option (1 or 2): "

set DRYRUN_FLAG=
if "!OPTION!"=="2" (
    set DRYRUN_FLAG=-DryRun
    echo.
    echo === DRY RUN MODE ENABLED ===
    echo No changes will be made to files or git repository
    echo.
)

REM Execute the PowerShell script
echo.
echo Starting release process...
echo ============================
echo.

if "!DRYRUN_FLAG!"=="" (
    powershell -ExecutionPolicy Bypass -File "release.ps1" -Version "!VERSION!" -Message "!MESSAGE!"
) else (
    powershell -ExecutionPolicy Bypass -File "release.ps1" -Version "!VERSION!" -Message "!MESSAGE!" -DryRun
)

set SCRIPT_EXIT_CODE=!errorlevel!

echo.
echo ============================

if !SCRIPT_EXIT_CODE! equ 0 (
    if "!DRYRUN_FLAG!"=="" (
        echo.
        echo ^✓ Release completed successfully!
        echo.
        echo The release has been pushed to GitHub and Actions will:
        echo - Build executables for all platforms
        echo - Create a new release with download links
        echo - Make it available for auto-updates
        echo.
        echo Monitor progress at:
        echo https://github.com/BySnowden/Glacier/actions
        echo.
        echo Release page:
        echo https://github.com/BySnowden/Glacier/releases/tag/v!VERSION!
    ) else (
        echo.
        echo ^✓ Dry run completed successfully!
        echo.
        echo This was a preview - no actual changes were made.
        echo To perform the real release, choose option 1.
    )
) else (
    echo.
    echo ^✗ Release failed with error code !SCRIPT_EXIT_CODE!
    echo.
    echo Check the output above for error details.
    echo You may need to:
    echo - Commit any pending changes
    echo - Check your Git configuration
    echo - Verify network connectivity
    echo - Use -Force flag if needed
)

echo.
pause
