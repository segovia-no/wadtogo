BINARY_NAME=wadtogo

build:
	GOARCH=amd64 GOOS=darwin go build -o ${BINARY_NAME}-darwin-x64 .
	GOARCH=arm64 GOOS=darwin go build -o ${BINARY_NAME}-darwin-arm .
	GOARCH=amd64 GOOS=linux go build -o ${BINARY_NAME}-linux-x64 .
	GOARCH=386 GOOS=linux go build -o ${BINARY_NAME}-linux-x86 .
	GOARCH=amd64 GOOS=windows go build -o ${BINARY_NAME}-windows-x64.exe .
	GOARCH=386 GOOS=windows go build -o ${BINARY_NAME}-windows-x86.exe .

run: build
	./${BINARY_NAME}

clean:
	go clean
	rm ${BINARY_NAME}-darwin
	rm ${BINARY_NAME}-linux
	rm ${BINARY_NAME}-windows
