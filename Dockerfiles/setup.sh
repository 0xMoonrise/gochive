#!/usr/bin/env bash
set -e

ARCH=$(uname -m)
DOWNLOAD_DIR="/tmp/pdfium"
case "$ARCH" in
  x86_64)
    TAR="pdfium-linux-x86.tgz"
    ;;
  aarch64)
    TAR="pdfium-linux-arm64.tgz"
    ;;
  *)
    echo "Architecture not supported: $ARCH"
    exit 1
    ;;
esac

mkdir -p $DOWNLOAD_DIR
wget -q "https://github.com/bblanchon/pdfium-binaries/releases/latest/download/$TAR" -P /tmp/
tar -xzf "/tmp/$TAR" -C $DOWNLOAD_DIR

tee /etc/profile.d/pdfium.sh > /dev/null <<'EOF'
export PKG_CONFIG_PATH=/usr/local/lib/pkgconfig:$PKG_CONFIG_PATH
EOF

cp -r $DOWNLOAD_DIR/lib/* /usr/local/lib/
cp -r $DOWNLOAD_DIR/include/* /usr/local/include/

mkdir -p /usr/local/lib/pkgconfig

. $DOWNLOAD_DIR/VERSION

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
