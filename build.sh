#!/bin/bash

# Script de build do frontend
# Copia e organiza arquivos de web/ para public/

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WEB_DIR="${SCRIPT_DIR}/web"
PUBLIC_DIR="${SCRIPT_DIR}/public"

echo "Building frontend from ${WEB_DIR} to ${PUBLIC_DIR}..."

# Remove pasta public se existir (exceto se estiver em git)
if [ -d "${PUBLIC_DIR}" ]; then
    echo "Cleaning existing public directory..."
    rm -rf "${PUBLIC_DIR}"
fi

# Cria pasta public
mkdir -p "${PUBLIC_DIR}"

# Copia arquivos da raiz de web/
echo "Copying root files..."
cp "${WEB_DIR}/index.html" "${PUBLIC_DIR}/"
cp "${WEB_DIR}/styles.css" "${PUBLIC_DIR}/"
cp "${WEB_DIR}/app.webmanifest" "${PUBLIC_DIR}/"

# Copia pasta images/
echo "Copying images..."
cp -r "${WEB_DIR}/images" "${PUBLIC_DIR}/"

# Copia conteúdo de src/ para raiz de public/
# (app.js e subpastas)
echo "Copying source files..."
cp -r "${WEB_DIR}/src"/* "${PUBLIC_DIR}/"

# Move app.js para raiz (se necessário, já que index.html referencia /app.js)
if [ -f "${PUBLIC_DIR}/app.js" ]; then
    # app.js já está na raiz, perfeito
    echo "app.js is in the root ✓"
else
    echo "Warning: app.js not found in expected location"
fi

echo "Build completed successfully!"
echo "Output directory: ${PUBLIC_DIR}"
