#!/bin/sh

targets="
    windows/amd64
    darwin/amd64
    darwin/arm64
    linux/amd64
    linux/arm64
"

mkdir -p bin

for target in $targets; do
    goos=${target%/*}
    goarch=${target#*/}

    output_name=gfpl_${goos}_${goarch}

    if [ "$goos" = "windows" ]; then
        output_name="${output_name}.exe"
    fi

    echo "building $output_name..."

    env GOOS="$goos" GOARCH="$goarch" go build -o "bin/$output_name" ./src/

    if [ $? -ne 0 ]; then
        echo "An error occured."
        exit 1
    fi
done

echo "done"
