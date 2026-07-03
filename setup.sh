#!/usr/bin/env bash

set -e

sudo mkdir -p /opt/gochive/lib
sudo chown 1000:1000 -R /opt/gochive/

PDF_VER="5.4.530"
HLJS_VER="11.11.1"
MERMAID_VER="11.16.0"
MATHJAX_VER="3.2.2"

ARCH=$(uname -m)
TMP=/tmp/pdfium

mkdir -p $TMP
for cmd in gcc pkg-config zip wget make; do
    command -v "$cmd" >/dev/null 2>&1 || need_pkgs+=("$cmd")
done

if [ "${#need_pkgs[@]}" -eq 0 ]; then
    echo "Necessary software are already installed."
else
    echo "Missing: ${need_pkgs[*]}"
    if command -v dnf >/dev/null 2>&1; then
        sudo dnf install -y "${need_pkgs[@]}"
    elif command -v apt >/dev/null 2>&1; then
        sudo apt update
        sudo apt install -y "${need_pkgs[@]}"
    else
        echo "Neither dnf nor apt found. Install manually: ${need_pkgs[*]}" >&2
        exit 1
    fi
fi

case "$ARCH" in
  x86_64)
    TAR="pdfium-linux-x64.tgz"
    ;;
  aarch64)
    TAR="pdfium-linux-arm64.tgz"
    ;;
  *)
    echo "Architecture not supported: $ARCH"
    exit 1
    ;;
esac

wget -q \
  "https://github.com/bblanchon/pdfium-binaries/releases/latest/download/$TAR" \
  -O "/tmp/$TAR"
tar -xzf "/tmp/$TAR" -C "$TMP"

sudo cp "$TMP/lib/libpdfium.so" /usr/local/lib/
sudo cp -r "$TMP/include/"* /usr/local/include/
sudo mkdir -p /usr/local/lib/pkgconfig

echo "/usr/local/lib" | sudo tee /etc/ld.so.conf.d/local.conf
sudo ldconfig

. "$TMP/VERSION"

sudo tee /usr/local/lib/pkgconfig/pdfium.pc > /dev/null <<EOF
prefix=/usr/local
exec_prefix=\${prefix}
libdir=\${exec_prefix}/lib
includedir=\${prefix}/include
Name: pdfium
Description: PDFium
Version: ${MAJOR}.${MINOR}.${BUILD}.${PATCH}
Libs: -L\${libdir} -lpdfium
Cflags: -I\${includedir}
EOF

wget "https://github.com/mozilla/pdf.js/releases/download/v${PDF_VER}/pdfjs-${PDF_VER}-dist.zip" -P /tmp/
unzip "/tmp/pdfjs-${PDF_VER}-dist.zip" -d /opt/gochive/lib/pdfjs/

mkdir -p /opt/gochive/lib/highlightjs
wget -q \
  "https://registry.npmjs.org/@highlightjs/cdn-assets/-/cdn-assets-${HLJS_VER}.tgz" \
  -O "/tmp/highlightjs-${HLJS_VER}.tgz"
tar -xzf "/tmp/highlightjs-${HLJS_VER}.tgz" -C /tmp
cp -r /tmp/package/* /opt/gochive/lib/highlightjs/
rm -rf /tmp/package

mkdir -p /opt/gochive/lib/mermaid
wget -q \
  "https://registry.npmjs.org/mermaid/-/mermaid-${MERMAID_VER}.tgz" \
  -O "/tmp/mermaid-${MERMAID_VER}.tgz"
tar -xzf "/tmp/mermaid-${MERMAID_VER}.tgz" -C /tmp
cp -r /tmp/package/dist/* /opt/gochive/lib/mermaid/
rm -rf /tmp/package

mkdir -p /opt/gochive/lib/mathjax
wget -q \
  "https://registry.npmjs.org/mathjax/-/mathjax-${MATHJAX_VER}.tgz" \
  -O "/tmp/mathjax-${MATHJAX_VER}.tgz"
tar -xzf "/tmp/mathjax-${MATHJAX_VER}.tgz" -C /tmp
cp -r /tmp/package/es5/* /opt/gochive/lib/mathjax/
rm -rf /tmp/package

sudo chown 1000:1000 -R /opt/gochive/
