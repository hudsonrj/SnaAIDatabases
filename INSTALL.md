# Guia de Instalação - SnipAI Databases

Este guia fornece instruções detalhadas de instalação para todos os sistemas operacionais.

## 📋 Índice

- [Linux](#linux)
- [macOS](#macos)
- [Windows](#windows)
- [Verificação](#verificação)
- [Configuração](#configuração)
- [Solução de Problemas](#solução-de-problemas)

## 🐧 Linux

### Debian/Ubuntu

```bash
# 1. Atualizar sistema
sudo apt-get update

# 2. Instalar dependências
sudo apt-get install -y \
  build-essential \
  libsqlite3-dev \
  git \
  curl

# 3. Instalar Go (se não estiver instalado)
# Baixe de https://go.dev/dl/ ou use:
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 4. Clonar repositório
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases

# 5. Instalar dependências Go
go mod download

# 6. Compilar
go build -o snip main.go

# 7. Instalar no sistema
sudo mv snip /usr/local/bin/

# 8. Verificar
snip --version
```

### RHEL/CentOS/Fedora

```bash
# 1. Instalar dependências
# RHEL/CentOS:
sudo yum install -y gcc sqlite-devel git curl

# Fedora:
sudo dnf install -y gcc sqlite-devel git curl

# 2. Instalar Go (se necessário)
# Baixe de https://go.dev/dl/
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 3. Clonar e compilar
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases
go mod download
go build -o snip main.go
sudo mv snip /usr/local/bin/
snip --version
```

### Arch Linux

```bash
# 1. Instalar dependências
sudo pacman -S sqlite git base-devel curl

# 2. Instalar Go (se necessário)
sudo pacman -S go

# 3. Clonar e compilar
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases
go mod download
go build -o snip main.go
sudo mv snip /usr/local/bin/
snip --version
```

## 🍎 macOS

### Homebrew (Recomendado)

```bash
# 1. Instalar Homebrew (se não tiver)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 2. Instalar dependências
brew install sqlite git go

# 3. Clonar repositório
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases

# 4. Instalar dependências Go
go mod download

# 5. Compilar
go build -o snip main.go

# 6. Instalar no sistema
sudo mv snip /usr/local/bin/

# 7. Verificar
snip --version
```

### Instalação Manual

```bash
# 1. Instalar Go
# Baixe de https://go.dev/dl/
# Instale o .pkg

# 2. Instalar dependências
brew install sqlite git

# 3. Clonar e compilar
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases
go mod download
go build -o snip main.go
sudo mv snip /usr/local/bin/
snip --version
```

**⚠️ Nota de Segurança macOS:**

Se o macOS bloquear a execução:

```bash
# Remover atributo de quarentena
xattr -d com.apple.quarantine /usr/local/bin/snip

# Ou permitir nas Configurações do Sistema
# Configurações do Sistema > Privacidade e Segurança > Permitir "snip"
```

## 🪟 Windows

### PowerShell

```powershell
# 1. Verificar Go
go version
# Se não estiver instalado, baixe de https://go.dev/dl/

# 2. Instalar Git (se necessário)
# Baixe de https://git-scm.com/download/win

# 3. Clonar repositório
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases

# 4. Instalar dependências Go
go mod download

# 5. Compilar
$env:CGO_ENABLED=1
go build -o snip.exe main.go

# 6. Adicionar ao PATH (opcional)
# Copie snip.exe para C:\Program Files\SnipAI\
# Adicione C:\Program Files\SnipAI\ ao PATH nas Variáveis de Ambiente
```

### CMD

```cmd
REM 1. Verificar Go
go version

REM 2. Clonar repositório
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases

REM 3. Instalar dependências
go mod download

REM 4. Compilar
set CGO_ENABLED=1
go build -o snip.exe main.go
```

### Chocolatey

```powershell
# Instalar via Chocolatey (quando disponível)
choco install snaai-databases

# Atualizar
choco upgrade snaai-databases
```

### Scoop

```powershell
# Adicionar bucket (quando disponível)
scoop bucket add snaai https://github.com/hudsonrj/scoop-bucket

# Instalar
scoop install snaai-databases

# Atualizar
scoop update snaai-databases
```

## ✅ Verificação

Após a instalação, verifique se tudo está funcionando:

```bash
# Verificar versão
snip --version

# Verificar ajuda
snip --help

# Listar comandos disponíveis
snip db-analysis --help
```

## ⚙️ Configuração

### Configurar Groq API Key

```bash
# Linux/macOS
export GROQ_API_KEY="sua_chave_aqui"
echo 'export GROQ_API_KEY="sua_chave_aqui"' >> ~/.bashrc

# Windows PowerShell
$env:GROQ_API_KEY="sua_chave_aqui"
[Environment]::SetEnvironmentVariable("GROQ_API_KEY", "sua_chave_aqui", "User")

# Windows CMD
set GROQ_API_KEY=sua_chave_aqui
```

## 🔧 Solução de Problemas

### Erro: "command not found: snip"

**Linux/macOS:**
```bash
# Verificar se está no PATH
which snip

# Se não estiver, adicionar manualmente
export PATH=$PATH:/usr/local/bin
```

**Windows:**
- Verifique se `snip.exe` está em uma pasta no PATH
- Adicione a pasta ao PATH nas Variáveis de Ambiente

### Erro: "CGO_ENABLED=0"

**Solução:**
```bash
# Linux/macOS
export CGO_ENABLED=1
go build -o snip main.go

# Windows PowerShell
$env:CGO_ENABLED=1
go build -o snip.exe main.go
```

### Erro: "sqlite3.h: No such file"

**Linux:**
```bash
# Debian/Ubuntu
sudo apt-get install libsqlite3-dev

# RHEL/CentOS
sudo yum install sqlite-devel

# Fedora
sudo dnf install sqlite-devel
```

**macOS:**
```bash
brew install sqlite
```

**Windows:**
- SQLite geralmente vem com Go
- Se necessário, baixe de https://www.sqlite.org/download.html

### Erro de Permissão (macOS)

```bash
# Remover quarentena
xattr -d com.apple.quarantine /usr/local/bin/snip

# Dar permissão de execução
chmod +x /usr/local/bin/snip
```

## 📞 Suporte

Se encontrar problemas:

1. Verifique os [Issues](https://github.com/hudsonrj/SnaAIDatabases/issues)
2. Abra um novo Issue com detalhes do problema
3. Consulte a [Documentação](README.md)

