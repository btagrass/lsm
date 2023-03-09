#!/bin/sh
read -p "Please enter the os(Linux / Windows), arch(aMd64 / aRm64) and ip: " os arch ip
if [[ ! $os ]]; then
    os=`uname | tr "A-Z" "a-z"`
fi
if [[ $os == "w" || $os =~ "win" ]]; then
    os="windows"
else
    os="linux"
fi
if [[ ! $arch ]]; then
    arch=`uname -m`
fi
if [[ $arch == "r" || $arch =~ "arm" ]]; then
    arch="arm64"
else
    arch="amd64"
fi
if [[ ! $ip ]]; then
    ip="localhost"
fi
name=${PWD##*/}

echo
echo --- Prepare ---
echo [i] Configuring proxy...
go env -w GOPROXY=https://goproxy.cn,direct
echo [i] Installing cmds...
go install github.com/swaggo/swag/cmd/swag@latest
echo [i] Upgrading packages...
go get -u
go mod tidy
echo
echo --- Build ---
echo [i] Target: $os/$arch
echo [i] Cleaning dirs...
rm -rf build/$os-$arch
echo [i] Building binary...
sed -i "s/\\/[^:]*/\\/\\/$ip/g" web/.env.production
yarn --cwd web install
yarn --cwd web build
for dir in {api,mgt}; do
    if [[ -d $dir ]]; then
        swag i -d $dir -g $dir.go --instanceName $dir --pd -q
    fi
done
CGO_ENABLED=0 GOOS=$os GOARCH=$arch go build -ldflags "-s -w" -o build/$os-$arch/ $name
echo [i] Copying files...
cp -r conf build/$os-$arch
echo
echo --- End ---
