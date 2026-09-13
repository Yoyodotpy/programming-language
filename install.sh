#!/bin/sh
set -e

REPO="Yoyodotpy/programming-language"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    aarch64|arm64) ARCH="arm64" ;;
   *) echo "Unsupported Architecture: $ARCH" && exit 1 ;;
esac

case "$OS" in
    darwin) OS="darwin" ;;
    linux) OS="linux" ;;
    *) echo "Unsupported operating system: $OS" && exit 1 ;;
esac

URL="https://github.com/$REPO/releases/latest/download/gfpl_${OS}_${ARCH}"

echo "downloading gfpl from $URL..."

mkdir -p ./gfpl
cd ./gfpl

curl -sL "$URL" -o gfpl
chmod +x ./gfpl

echo "downloading example scripts..."
curl -sLO https://raw.githubusercontent.com/Yoyodotpy/programming-language/refs/heads/main/examples/fibonacci.gfpl
curl -sLO https://raw.githubusercontent.com/Yoyodotpy/programming-language/refs/heads/main/examples/hello_world_complicated.gfpl

echo "successfully downloaded gfpl and example scripts."
echo "to test gfpl, you can run:"
echo "./gfpl fibonacci.lamb"
echo "please look at the git repo if you need help: github.com/$REPO"
