@echo off
if "%1" == "" (
    set version=dev
) else (
    set version=%1
)
echo This might screw up your GOOS env variable if it's not default
pause
go env -w GOOS=windows
go build -trimpath -ldflags "-s -w -X 'github.com/vladaad/discordcompressor/build.VERSION=%version%'" -o discordcompressor.exe
go env -w GOOS=linux
go build -trimpath -ldflags "-s -w -X 'github.com/vladaad/discordcompressor/build.VERSION=%version%'" -o discordcompressor
go env -u GOOS
set z="C:\Program Files\7-Zip\7z.exe"
set options=-mmt1 -mx9
%z% a %options% discordcompressor-windows.zip discordcompressor.exe
#%z% a %options% discordcompressor-windows-portable.zip discordcompressor-portable.exe
%z% a %options% discordcompressor-linux.zip discordcompressor
#%z% a %options% discordcompressor-linux-portable.zip discordcompressor-portable