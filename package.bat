wsl ./build.sh %1
set z="C:\Program Files\7-Zip\7z.exe"
set options=-mmt1 -mx9
%z% a %options% discordcompressor-windows.zip discordcompressor.exe
%z% a %options% discordcompressor-windows-portable.zip discordcompressor-portable.exe
%z% a %options% discordcompressor-linux.zip discordcompressor
%z% a %options% discordcompressor-linux-portable.zip discordcompressor-portable