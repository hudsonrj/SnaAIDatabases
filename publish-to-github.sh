#!/bin/bash

# Script para publicar repositório no GitHub

set -e

GITHUB_USER="hudsonrj"
REPO_NAME="SnaAIDatabases"

echo "🚀 Publicando repositório no GitHub..."
echo ""

# Verificar se já tem remote
if git remote get-url origin &> /dev/null; then
    echo "⚠️  Remote 'origin' já existe:"
    git remote -v
    read -p "Deseja substituir? (s/n): " REPLACE
    if [ "$REPLACE" = "s" ] || [ "$REPLACE" = "S" ]; then
        git remote remove origin
    else
        echo "Operação cancelada."
        exit 0
    fi
fi

# Tentar criar via API (requer token)
echo "📝 Para criar o repositório, você tem 2 opções:"
echo ""
echo "Opção 1: Via Web (Recomendado)"
echo "  1. Acesse: https://github.com/new"
echo "  2. Nome: $REPO_NAME"
echo "  3. Descrição: Ferramenta CLI para análise de bancos de dados com IA"
echo "  4. Público"
echo "  5. NÃO inicialize com README, .gitignore ou licença"
echo "  6. Clique em 'Create repository'"
echo "  7. Depois execute: git remote add origin https://github.com/$GITHUB_USER/$REPO_NAME.git"
echo "  8. E então: git push -u origin main"
echo ""
echo "Opção 2: Via API (requer token)"
read -p "Você tem um GitHub Personal Access Token? (s/n): " HAS_TOKEN

if [ "$HAS_TOKEN" = "s" ] || [ "$HAS_TOKEN" = "S" ]; then
    read -sp "Cole seu token: " GITHUB_TOKEN
    echo ""
    
    echo "📤 Criando repositório via API..."
    curl -X POST \
      -H "Authorization: token $GITHUB_TOKEN" \
      -H "Accept: application/vnd.github.v3+json" \
      https://api.github.com/user/repos \
      -d "{\"name\":\"$REPO_NAME\",\"description\":\"Ferramenta CLI para análise de bancos de dados com IA\",\"public\":true}" \
      && echo "✅ Repositório criado!" || echo "❌ Erro ao criar repositório"
    
    echo ""
    echo "📤 Adicionando remote e fazendo push..."
    git remote add origin https://github.com/$GITHUB_USER/$REPO_NAME.git
    git push -u origin main
    
    echo ""
    echo "🎉 Repositório publicado em: https://github.com/$GITHUB_USER/$REPO_NAME"
else
    echo ""
    echo "📋 Siga as instruções da Opção 1 acima."
    echo ""
    echo "Ou obtenha um token em: https://github.com/settings/tokens"
    echo "Permissões necessárias: repo (Full control of private repositories)"
fi
