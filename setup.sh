#!/usr/bin/env bash
set -e

ROOT=/opt/gochive

sudo mkdir -p "${ROOT}/lib"
sudo chown 1000:1000 -R "${ROOT}"

PDF_VER="6.1.200"
ARCH=$(uname -m)
TMP=/tmp/pdfium
mkdir -p "$TMP"

NPM_PACKAGES=(
  "highlightjs|@highlightjs/cdn-assets|11.11.1|.|highlightjs"
  "mermaid|mermaid|11.16.0|dist|mermaid"
  "mathjax|mathjax|3.2.2|es5|mathjax"
  "marked|marked|18.0.6|lib|marked"
  "dompurify|dompurify|3.4.12|dist|dompurify"
)

for cmd in gcc pkg-config zip unzip wget make; do
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

wget -q "https://github.com/mozilla/pdf.js/releases/download/v${PDF_VER}/pdfjs-${PDF_VER}-dist.zip" -P /tmp/
unzip -q "/tmp/pdfjs-${PDF_VER}-dist.zip" -d "${ROOT}/lib/pdfjs/"

for entry in "${NPM_PACKAGES[@]}"; do
  IFS='|' read -r LABEL PKG_NAME VER SUBDIR DEST <<< "$entry"

  echo "Installing ${LABEL}@${VER}..."
  DEST_DIR="${ROOT}/lib/${DEST}"
  mkdir -p "$DEST_DIR"

  TARBALL="/tmp/${LABEL}-${VER}.tgz"
  wget -q \
    "https://registry.npmjs.org/${PKG_NAME}/-/$(basename "$PKG_NAME")-${VER}.tgz" \
    -O "$TARBALL"

  rm -rf /tmp/package
  tar -xzf "$TARBALL" -C /tmp

  if [ "$SUBDIR" = "." ]; then
    cp -r /tmp/package/* "$DEST_DIR"/
  else
    cp -r "/tmp/package/${SUBDIR}/"* "$DEST_DIR"/
  fi

  rm -rf /tmp/package "$TARBALL"
done

sudo chown 1000:1000 -R ${ROOT}
