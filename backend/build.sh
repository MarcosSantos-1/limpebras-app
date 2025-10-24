#!/bin/bash
# Script de build para Render.com

echo "🚀 Iniciando build do backend..."

# Instalar Node.js (necessário para Puppeteer)
echo "📦 Instalando Node.js..."
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Verificar instalação do Node.js
echo "✅ Node.js instalado:"
node --version
npm --version

# Instalar dependências do Node.js
echo "📦 Instalando dependências do Node.js..."
npm install

# Build do Go
echo "🔧 Fazendo build do Go..."
go mod tidy
go build -o main .

echo "🎉 Build concluído com sucesso!"
