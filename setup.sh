#!/usr/bin/env bash
set -e

ARCH=$(uname -m)
TMP=/tmp/pdfium
mkdir -p $TMP

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

cp "$TMP/lib/libpdfium.so" /usr/local/lib/

cp -r "$TMP/include/"* /usr/local/include/

mkdir -p /usr/local/lib/pkgconfig
. "$TMP/VERSION"

cat > /usr/local/lib/pkgconfig/pdfium.pc <<EOF
prefix=/usr/local
exec_prefix=\${prefix}
libdir=\${exec_prefix}/lib
includedir=\${prefix}/include

Name: pdfium
Description: PDFium
Version: $MAJOR.$MINOR.$BUILD.$PATCH
Libs: -L\${libdir} -lpdfium
Cflags: -I\${includedir}
EOF

PDF_VER="5.4.530"

wget "https://github.com/mozilla/pdf.js/releases/download/v${PDF_VER}/pdfjs-${PDF_VER}-dist.zip" -P /opt/
mkdir -p /opt/pdfjs
unzip "/opt/pdfjs-${PDF_VER}-dist.zip" -d /opt/pdfjs/
