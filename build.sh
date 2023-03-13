#!/bin/sh
read -p "Please enter the os(Linux / Windows) and arch(aMd64 / aRm64): " os arch
if [[ ! $os ]]; then
    os=("linux" "windows")
fi
if [[ $os == "w" || $os =~ "win" ]]; then
    os="windows"
else
    os="linux"
fi
if [[ ! $arch ]]; then
    arch=("amd64" "arm64")
fi
if [[ $arch == "r" || $arch =~ "arm" ]]; then
    arch="arm64"
else
    arch="amd64"
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
rm -rf build/*
echo [i] Building web...
yarn --cwd web install
yarn --cwd web build
echo [i] Generating docs...
dir=("api" "mgt")
for d in ${dir[@]}; do
    if [[ -d $d ]]; then
        echo -e "\t$d"
        swag i -d $d -g $d.go --instanceName $d --pd -q
    fi
done
echo [i] Building binaries...
for o in ${os[@]}; do
    for a in ${arch[@]}; do
        echo -e "\t$o-$a"
        CGO_ENABLED=0 GOOS=$o GOARCH=$a go build -ldflags "-s -w" -o build/$o-$a/ $name
        cp -r conf build/$o-$a
    done
done
echo
echo --- End ---
