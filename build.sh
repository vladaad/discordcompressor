if [ "$1" = "" ]
then
    version=dev
else
    version=$1
fi
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w -X 'github.com/vladaad/discordcompressor/build.VERSION=$version'" -o "discordcompressor.exe" -v
GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "-s -w -X 'github.com/vladaad/discordcompressor/build.VERSION=$version' -X 'github.com/vladaad/discordcompressor/build.BUILD=portable'" -o "discordcompressor-portable.exe" -v
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -X 'github.com/vladaad/discordcompressor/build.VERSION=$version'" -o "discordcompressor" -v
GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "-s -w -X 'github.com/vladaad/discordcompressor/build.VERSION=$version' -X 'github.com/vladaad/discordcompressor/build.BUILD=portable'" -o "discordcompressor-portable" -v