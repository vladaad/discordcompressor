@echo off
if "%1" == "" (
    set version=dev
) else (
    set version=%1
)
go build -trimpath -ldflags "-s -w -X 'github.com/vladaad/discordcompressor/build.VERSION=%version%'"  -o discordcompressor.exe
powershell "Get-FileHash discordcompressor.exe"
pause