#!/bin/bash

BUILD_DIR=$(dirname "$0")/build/releases
rm -rf $BUILD_DIR
mkdir -p $BUILD_DIR

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
	env GOOS=$os GOARCH=amd64 go build -o $BUILD_DIR/azukiiro_${os}_amd64${suffix}
	tar -zcf $BUILD_DIR/azukiiro-${os}-amd64-$VERSION.tar.gz $BUILD_DIR/azukiiro_${os}_amd64${suffix}
  rm $BUILD_DIR/azukiiro_${os}_amd64${suffix}
	HASH=$($sum $BUILD_DIR/azukiiro-${os}-amd64-$VERSION.tar.gz)
  echo "$HASH" >> $BUILD_DIR/azukiiro-$VERSION.sha256
done

# 386
OSES=(linux windows)
for os in ${OSES[@]}; do
	suffix=""
	if [ "$os" == "windows" ]
	then
		suffix=".exe"
	fi
	env GOOS=$os GOARCH=386 go build -o $BUILD_DIR/azukiiro_${os}_386${suffix}
	tar -zcf $BUILD_DIR/azukiiro-${os}-386-$VERSION.tar.gz $BUILD_DIR/azukiiro_${os}_386${suffix}
  rm $BUILD_DIR/azukiiro_${os}_386${suffix}
	HASH=$($sum $BUILD_DIR/azukiiro-${os}-386-$VERSION.tar.gz)
  echo "$HASH" >> $BUILD_DIR/azukiiro-$VERSION.sha256
done

# ARM
ARMS=(5 6 7)
for v in ${ARMS[@]}; do
	env GOOS=linux GOARCH=arm GOARM=$v go build -o $BUILD_DIR/azukiiro_linux_arm$v 
  tar -zcf $BUILD_DIR/azukiiro-linux-arm$v-$VERSION.tar.gz $BUILD_DIR/azukiiro_linux_arm$v
  rm $BUILD_DIR/azukiiro_linux_arm$v
  HASH=$($sum $BUILD_DIR/azukiiro-linux-arm$v-$VERSION.tar.gz)
  echo "$HASH" >> $BUILD_DIR/azukiiro-$VERSION.sha256
done

# ARM64
OSES=(linux darwin windows)
for os in ${OSES[@]}; do
	suffix=""
	if [ "$os" == "windows" ]
	then
		suffix=".exe"
	fi
	env GOOS=$os GOARCH=arm64 go build -o $BUILD_DIR/azukiiro_${os}_arm64${suffix}
	tar -zcf $BUILD_DIR/azukiiro-${os}-arm64-$VERSION.tar.gz $BUILD_DIR/azukiiro_${os}_arm64${suffix}
  rm $BUILD_DIR/azukiiro_${os}_arm64${suffix}
  HASH=$($sum $BUILD_DIR/azukiiro-${os}-arm64-$VERSION.tar.gz)
  echo "$HASH" >> $BUILD_DIR/azukiiro-$VERSION.sha256
done
