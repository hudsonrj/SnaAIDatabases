#!/bin/bash

# Script para configurar e criar repositório GitHub SnaAIDatabases

set -e

echo "🚀 Configurando repositório SnaAIDatabases para GitHub..."

# Cores para output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Verificar se git está instalado
if ! command -v git &> /dev/null; then
    echo "❌ Git não está instalado. Por favor, instale o Git primeiro."
    exit 1
fi

# Verificar se estamos no diretório correto
if [ ! -f "main.go" ]; then
    echo "❌ Arquivo main.go não encontrado. Execute este script no diretório raiz do projeto."
    exit 1
fi

# Perguntar nome de usuário do GitHub
read -p "Digite seu nome de usuário do GitHub: " GITHUB_USER

if [ -z "$GITHUB_USER" ]; then
    echo "❌ Nome de usuário não pode estar vazio."
    exit 1
fi

echo ""
echo "📝 Substituindo hudsonrj por $GITHUB_USER nos arquivos..."

# Substituir hudsonrj nos arquivos
find . -type f \( -name "*.md" -o -name "*.yml" \) -exec sed -i '' "s/hudsonrj/$GITHUB_USER/g" {} +

echo "✅ Substituições concluídas"
echo ""

# Verificar se já é um repositório git
if [ -d ".git" ]; then
    echo "⚠️  Diretório já é um repositório Git."
    read -p "Deseja continuar mesmo assim? (s/n): " CONTINUE
    if [ "$CONTINUE" != "s" ] && [ "$CONTINUE" != "S" ]; then
        echo "Operação cancelada."
        exit 0
    fi
else
    echo "📦 Inicializando repositório Git..."
    git init
    git branch -M main
fi

# Adicionar arquivos
echo "📁 Adicionando arquivos ao Git..."
git add .

# Fazer commit inicial
echo "💾 Criando commit inicial..."
git commit -m "Initial commit: SnipAI Databases - Ferramenta CLI para análise de bancos de dados com IA" || {
    echo "⚠️  Nenhuma mudança para commitar ou commit falhou."
}

echo ""
echo "${GREEN}✅ Repositório local configurado!${NC}"
echo ""
echo "📤 Próximos passos:"
echo ""
echo "${YELLOW}Opção 1: Usando GitHub CLI (recomendado)${NC}"
echo "  1. Instale GitHub CLI: brew install gh (macOS) ou siga https://cli.github.com/"
echo "  2. Faça login: gh auth login"
echo "  3. Crie o repositório: gh repo create SnaAIDatabases --public --source=. --remote=origin --push"
echo ""
echo "${YELLOW}Opção 2: Usando GitHub Web${NC}"
echo "  1. Acesse: https://github.com/new"
echo "  2. Nome: SnaAIDatabases"
echo "  3. Descrição: Ferramenta CLI para análise de bancos de dados com IA"
echo "  4. Público"
echo "  5. NÃO inicialize com README, .gitignore ou licença"
echo "  6. Execute: git remote add origin https://github.com/$GITHUB_USER/SnaAIDatabases.git"
echo "  7. Execute: git push -u origin main"
echo ""
echo "${YELLOW}Opção 3: Usando GitHub Desktop${NC}"
echo "  1. Abra GitHub Desktop"
echo "  2. File > Add Local Repository"
echo "  3. Selecione este diretório"
echo "  4. Publish repository"
echo "  5. Nome: SnaAIDatabases, Público"
echo ""

