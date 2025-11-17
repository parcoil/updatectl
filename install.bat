@echo off
REM updatectl User Installer

echo Installing updatectl...

REM Combine user and system PATH for this session
for /f "tokens=2*" %%A in ('reg query "HKCU\Environment" /v Path 2^>nul') do set "USER_PATH=%%B"
for /f "tokens=2*" %%A in ('reg query "HKLM\SYSTEM\CurrentControlSet\Control\Session Manager\Environment" /v Path 2^>nul') do set "SYS_PATH=%%B"
set "PATH=%USER_PATH%;%SYS_PATH%"

REM Check if Go is installed
go version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Go is not installed or not in PATH.
    echo Please install Go or restart your terminal after installing it.
    pause
    exit /b 1
)

REM Build the binary
echo Building updatectl...
go build -o updatectl.exe main.go
if %errorlevel% neq 0 (
    echo Build failed.
    pause
    exit /b 1
)

REM Install location (user home)
set "INSTALL_DIR=%USERPROFILE%\updatectl"
if not exist "%INSTALL_DIR%" mkdir "%INSTALL_DIR%"

REM Copy binary
copy /Y updatectl.exe "%INSTALL_DIR%\" >nul

REM Add to user PATH if not already there
setlocal enabledelayedexpansion
for /f "tokens=2*" %%A in ('reg query "HKCU\Environment" /v Path 2^>nul') do (
    set "CUR_PATH=%%B"
)

echo !CUR_PATH! | find /I "%INSTALL_DIR%" >nul
if %errorlevel% neq 0 (
    echo Adding "%INSTALL_DIR%" to your user PATH...
    setx PATH "!CUR_PATH!;%INSTALL_DIR%" >nul
) else (
    echo Path already includes "%INSTALL_DIR%"
)

endlocal

echo.
echo updatectl installed successfully to "%INSTALL_DIR%"
echo.
echo You can now run "updatectl" from any new command prompt.
echo Run "updatectl init" to set up the daemon.
pause
