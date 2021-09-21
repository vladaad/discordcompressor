@echo off
if "%1" == "" (
    set version=0.4
) else (
    set version=%1
)
go build -trimpath -ldflags "-s -w -X 'github.com/vladaad/discordcompressor/build.VERSION=%version%'"  -o discordcompressor.exe
powershell "Get-FileHash discordcompressor.exe"
pause