go build -trimpath -ldflags "-s -w"  -o discordcompressor.exe
powershell "Get-FileHash discordcompressor.exe"
pause