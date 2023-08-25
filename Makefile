BINARY_NAME=azukiiro
OUT_DIR=build

.PHONY: build clean

build:
	go build -o ${OUT_DIR}/${BINARY_NAME}

clean:
	go clean
	rm -rf ${OUT_DIR}
