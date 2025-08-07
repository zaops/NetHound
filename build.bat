@echo off
echo 正在构建 Nethound 网络测试工具 - 多平台版本...
echo.

REM 创建构建目录
if not exist "build" mkdir build

REM 构建 Windows amd64
echo [1/3] 构建 Windows amd64...
set GOOS=windows
set GOARCH=amd64
go build -o build/nethound-windows-amd64.exe -ldflags "-s -w" .
if %ERRORLEVEL% NEQ 0 (
    echo Windows amd64 构建失败！
    goto :error
)

REM 构建 Linux amd64
echo [2/3] 构建 Linux amd64...
set GOOS=linux
set GOARCH=amd64
go build -o build/nethound-linux-amd64 -ldflags "-s -w" .
if %ERRORLEVEL% NEQ 0 (
    echo Linux amd64 构建失败！
    goto :error
)

REM 构建 Linux arm64
echo [3/3] 构建 Linux arm64...
set GOOS=linux
set GOARCH=arm64
go build -o build/nethound-linux-arm64 -ldflags "-s -w" .
if %ERRORLEVEL% NEQ 0 (
    echo Linux arm64 构建失败！
    goto :error
)

echo.
echo ========================================
echo 所有平台构建成功！
echo ========================================
echo 构建文件位置: build/ 目录
echo.
echo Windows amd64: build/nethound-windows-amd64.exe
echo Linux amd64:   build/nethound-linux-amd64
echo Linux arm64:   build/nethound-linux-arm64
echo.
echo 使用方法:
echo   Windows: .\build\nethound-windows-amd64.exe --help
echo   Linux:   ./build/nethound-linux-amd64 --help
echo.
goto :end

:error
echo.
echo 构建失败！请检查错误信息。
echo.

:end
pause
