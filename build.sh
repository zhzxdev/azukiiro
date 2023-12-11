#!/bin/bash

BUILD_DIR=$(dirname "$0")/build/releases
mkdir -p $BUILD_DIR
cd $BUILD_DIR

sum="sha256sum"
VERSION=$(git describe --tags --always --dirty)

# AMD64 
OSES=(linux darwin windows freebsd)
for os in ${OSES[@]}; do
	suffix=""
	if [ "$os" == "windows" ]
	then
		suffix=".exe"
	fi
	env GOOS=$os GOARCH=amd64 go build -o azukiiro_${os}_amd64${suffix} github.com/zhzxdev/azukiiro
	tar -zcf azukiiro-${os}-amd64-$VERSION.tar.gz azukiiro_${os}_amd64${suffix}
  rm azukiiro_${os}_amd64${suffix}
	HASH=$($sum azukiiro-${os}-amd64-$VERSION.tar.gz)
  echo "$HASH azukiiro-${os}-amd64-$VERSION.tar.gz" > azukiiro-$VERSION.sha256
done

# 386
OSES=(linux windows)
for os in ${OSES[@]}; do
	suffix=""
	if [ "$os" == "windows" ]
	then
		suffix=".exe"
	fi
	env GOOS=$os GOARCH=386 go build -o azukiiro_${os}_386${suffix} github.com/zhzxdev/azukiiro
	tar -zcf azukiiro-${os}-386-$VERSION.tar.gz azukiiro_${os}_386${suffix}
  rm azukiiro_${os}_386${suffix}
	HASH=$($sum azukiiro-${os}-386-$VERSION.tar.gz)
  echo "$HASH azukiiro-${os}-386-$VERSION.tar.gz" > azukiiro-$VERSION.sha256
done

# ARM
ARMS=(5 6 7)
for v in ${ARMS[@]}; do
	env GOOS=linux GOARCH=arm GOARM=$v go build -o azukiiro_linux_arm$v  github.com/zhzxdev/azukiiro
  tar -zcf azukiiro-linux-arm$v-$VERSION.tar.gz azukiiro_linux_arm$v
  rm azukiiro_linux_arm$v
  HASH=$($sum azukiiro-linux-arm$v-$VERSION.tar.gz)
  echo "$HASH azukiiro-linux-arm$v-$VERSION.tar.gz" > azukiiro-$VERSION.sha256
done

# ARM64
OSES=(linux darwin windows)
for os in ${OSES[@]}; do
	suffix=""
	if [ "$os" == "windows" ]
	then
		suffix=".exe"
	fi
	env GOOS=$os GOARCH=arm64 go build -o azukiiro_${os}_arm64${suffix} github.com/zhzxdev/azukiiro
	tar -zcf azukiiro-${os}-arm64-$VERSION.tar.gz azukiiro_${os}_arm64${suffix}
  rm azukiiro_${os}_arm64${suffix}
  HASH=$($sum azukiiro-${os}-arm64-$VERSION.tar.gz)
  echo "$HASH azukiiro-${os}-arm64-$VERSION.tar.gz" > azukiiro-$VERSION.sha256
done
