#!/bin/sh
read -p "Please enter the ip: " ip
if [[ ! $ip ]]; then
    ip="localhost"
fi
os=`uname | tr "A-Z" "a-z"`
arch=`uname -m`
if [[ $arch == "x86_64" ]]; then
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
rm -rf build/$os-$arch
echo [i] Building binary...
sed -i "s/\\/[^:]*/\\/\\/$ip/g" web/.env.production
# echo "VITE_MGT_URL=http://localhost:3081" > web/.env.production
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
