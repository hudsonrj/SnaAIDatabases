# Instruções para Criar o Repositório no GitHub

## Opção 1: Usando GitHub CLI (gh)

```bash
# 1. Instalar GitHub CLI (se não tiver)
# macOS:
brew install gh

# Linux:
# Siga instruções em: https://cli.github.com/

# 2. Fazer login
gh auth login

# 3. Criar repositório
gh repo create SnaAIDatabases --public --source=. --remote=origin --push
```

## Opção 2: Usando Git e GitHub Web

```bash
# 1. Criar repositório no GitHub
# - Acesse: https://github.com/new
# - Nome: SnaAIDatabases
# - Descrição: "Ferramenta CLI para análise de bancos de dados com IA"
# - Público
# - NÃO inicialize com README, .gitignore ou licença

# 2. No diretório do projeto, execute:
git init
git add .
git commit -m "Initial commit: SnipAI Databases"
git branch -M main
git remote add origin https://github.com/hudsonrj/SnaAIDatabases.git
git push -u origin main
```

## Opção 3: Usando GitHub Desktop

1. Abra GitHub Desktop
2. File > Add Local Repository
3. Selecione o diretório do projeto
4. Publish repository
5. Nome: SnaAIDatabases
6. Marque "Keep this code private" como DESMARCADO (público)
7. Publish

## Após Criar o Repositório

1. Atualize o README.md substituindo hudsonrj pelo seu usuário GitHub
2. Adicione uma descrição no repositório
3. Adicione tópicos: go, database, cli, ai, oracle, sqlserver, mysql, postgresql, mongodb
4. Configure GitHub Actions (se necessário)

## Estrutura de Arquivos Recomendada

```
SnaAIDatabases/
├── README.md                 # README principal
├── SNAAIDATABASES_README.md  # README específico do repositório
├── INSTALL.md               # Guia de instalação detalhado
├── LICENSE                   # Licença (MIT recomendado)
├── .gitignore               # Arquivos ignorados
├── .github/
│   └── workflows/
│       └── release.yml      # GitHub Actions para releases
├── main.go
├── go.mod
├── go.sum
└── ... (resto dos arquivos)
```
