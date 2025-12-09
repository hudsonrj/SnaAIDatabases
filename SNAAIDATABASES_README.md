# SnipAI Databases

Ferramenta CLI poderosa para análise de bancos de dados com Inteligência Artificial, suportando Oracle, SQL Server, MySQL, PostgreSQL e MongoDB.

## 🚀 Características Principais

- **Análise Inteligente**: IA integrada para interpretação de resultados e geração de recomendações
- **Multi-Banco**: Suporte para Oracle, SQL Server, MySQL, PostgreSQL, MongoDB
- **Oracle RAC**: Análises específicas para clusters RAC (saúde, erros, listener, latência)
- **Visualização**: Geração automática de gráficos (ASCII, HTML, Bar, Line, Pie)
- **Planos de Manutenção**: IA gera planos detalhados baseados em análises
- **Transformação em Projetos**: Converte análises e incidentes em projetos estruturados
- **Chat Interativo**: Conversa com o banco de dados usando linguagem natural
- **Análise Dinâmica**: Gera queries SQL baseadas em solicitações em linguagem natural
- **Autenticação Local**: Conexões locais sem senha quando o usuário é owner do banco

## 📋 Requisitos

- **Go 1.21 ou superior** - [Download Go](https://go.dev/dl/)
- **SQLite3** (bibliotecas de desenvolvimento)
- **Groq API Key** (para funcionalidades de IA) - [Obter chave](https://console.groq.com/keys)

## 🔧 Instalação

### Linux

#### Debian/Ubuntu

```bash
# Instalar dependências
sudo apt-get update
sudo apt-get install -y build-essential libsqlite3-dev git

# Clonar repositório
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases

# Instalar dependências Go
go mod download

# Compilar
go build -o snip main.go

# Instalar no sistema
sudo mv snip /usr/local/bin/

# Verificar instalação
snip --version
```

#### RHEL/CentOS/Fedora

```bash
# Instalar dependências
sudo yum install -y gcc sqlite-devel git
# ou para Fedora:
# sudo dnf install -y gcc sqlite-devel git

# Clonar repositório
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases

# Instalar dependências Go
go mod download

# Compilar
go build -o snip main.go

# Instalar no sistema
sudo mv snip /usr/local/bin/

# Verificar instalação
snip --version
```

#### Arch Linux

```bash
# Instalar dependências
sudo pacman -S sqlite git base-devel

# Clonar repositório
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases

# Instalar dependências Go
go mod download

# Compilar
go build -o snip main.go

# Instalar no sistema
sudo mv snip /usr/local/bin/

# Verificar instalação
snip --version
```

### macOS

#### Homebrew (Recomendado)

```bash
# Instalar dependências
brew install sqlite git

# Clonar repositório
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases

# Instalar dependências Go
go mod download

# Compilar
go build -o snip main.go

# Instalar no sistema
sudo mv snip /usr/local/bin/

# Verificar instalação
snip --version
```

#### Instalação Manual

```bash
# Instalar dependências (se necessário)
brew install sqlite git

# Clonar repositório
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases

# Instalar dependências Go
go mod download

# Compilar
go build -o snip main.go

# Instalar no sistema
sudo mv snip /usr/local/bin/

# Verificar instalação
snip --version
```

**⚠️ Nota de Segurança macOS:**

Se o macOS bloquear a execução:

```bash
# Opção 1: Remover atributo de quarentena
xattr -d com.apple.quarantine /usr/local/bin/snip

# Opção 2: Permitir nas Configurações do Sistema
# Ir em: Configurações do Sistema > Privacidade e Segurança > Permitir "snip"
```

### Windows

#### PowerShell

```powershell
# Verificar se Go está instalado
go version

# Se não estiver instalado, baixe de: https://go.dev/dl/

# Clonar repositório
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases

# Instalar dependências Go
go mod download

# Compilar para Windows
$env:CGO_ENABLED=1
go build -o snip.exe main.go

# Adicionar ao PATH (opcional)
# Copie snip.exe para uma pasta no PATH, por exemplo:
# C:\Program Files\SnipAI\
# Depois adicione ao PATH nas Variáveis de Ambiente do Sistema
```

#### CMD

```cmd
REM Verificar se Go está instalado
go version

REM Clonar repositório
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases

REM Instalar dependências Go
go mod download

REM Compilar
set CGO_ENABLED=1
go build -o snip.exe main.go
```

#### Scoop (Gerenciador de Pacotes)

```powershell
# Adicionar bucket (quando disponível)
scoop bucket add snaai https://github.com/hudsonrj/scoop-bucket

# Instalar
scoop install snaai

# Atualizar
scoop update snaai
```

#### Chocolatey (Gerenciador de Pacotes)

```powershell
# Instalar (quando disponível)
choco install snaai

# Atualizar
choco upgrade snaai
```

### Instalação via Binários Pré-compilados

Binários pré-compilados estão disponíveis nas [Releases](https://github.com/hudsonrj/SnaAIDatabases/releases):

#### Linux

```bash
# Baixar binário
wget https://github.com/hudsonrj/SnaAIDatabases/releases/latest/download/snip-linux-amd64

# Tornar executável
chmod +x snip-linux-amd64

# Instalar
sudo mv snip-linux-amd64 /usr/local/bin/snip

# Verificar
snip --version
```

#### macOS

```bash
# Baixar binário
wget https://github.com/hudsonrj/SnaAIDatabases/releases/latest/download/snip-darwin-amd64

# Tornar executável
chmod +x snip-darwin-amd64

# Instalar
sudo mv snip-darwin-amd64 /usr/local/bin/snip

# Verificar
snip --version
```

#### Windows

1. Baixe `snip-windows-amd64.exe` das [Releases](https://github.com/hudsonrj/SnaAIDatabases/releases)
2. Renomeie para `snip.exe`
3. Adicione ao PATH ou coloque em uma pasta acessível

## ⚙️ Configuração

### Configurar Groq API Key

Para usar funcionalidades de IA, configure a variável de ambiente `GROQ_API_KEY`:

#### Linux/macOS

```bash
# Temporário (apenas sessão atual)
export GROQ_API_KEY="sua_chave_aqui"

# Permanente (adicionar ao ~/.bashrc ou ~/.zshrc)
echo 'export GROQ_API_KEY="sua_chave_aqui"' >> ~/.bashrc
source ~/.bashrc
```

#### Windows PowerShell

```powershell
# Temporário
$env:GROQ_API_KEY="sua_chave_aqui"

# Permanente
[Environment]::SetEnvironmentVariable("GROQ_API_KEY", "sua_chave_aqui", "User")
```

#### Windows CMD

```cmd
REM Temporário
set GROQ_API_KEY=sua_chave_aqui

REM Permanente: Painel de Controle > Sistema > Configurações Avançadas > Variáveis de Ambiente
```

**Obter Chave API:**
1. Visite [Groq Console](https://console.groq.com/keys)
2. Faça login ou crie uma conta
3. Gere uma nova chave API
4. Copie a chave

### Verificar Configuração

```bash
# Linux/macOS
echo $GROQ_API_KEY

# Windows PowerShell
echo $env:GROQ_API_KEY

# Windows CMD
echo %GROQ_API_KEY%
```

## 📖 Uso Rápido

### Análise de Banco de Dados

```bash
# Criar análise
snip db-analysis create \
  --title "Análise PostgreSQL" \
  --db-type postgresql \
  --analysis-type diagnostic \
  --host localhost \
  --port 5432 \
  --database mydb \
  --username user

# Executar análise
snip db-analysis run 1

# Ver resultados
snip db-analysis get 1
```

### Chat Interativo

```bash
snip db-chat \
  --db-type postgresql \
  --host localhost \
  --port 5432 \
  --database mydb \
  --username user
```

### Análise Oracle RAC

```bash
snip db-analysis create \
  --title "Saúde RAC" \
  --db-type oracle \
  --analysis-type rac_health \
  --host localhost \
  --port 1521 \
  --database RACDB \
  --username sys
```

## 📚 Documentação Completa

Para documentação completa, consulte o [README.md](README.md) principal do projeto.

## 🛠️ Desenvolvimento

### Pré-requisitos

- Go 1.21 ou superior
- SQLite3 (bibliotecas de desenvolvimento)
- Git

### Compilar

```bash
# Clonar repositório
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases

# Instalar dependências
go mod download

# Compilar
go build -o snip main.go

# Executar testes
go test ./...
```

## 🤝 Contribuindo

Contribuições são bem-vindas! Por favor:

1. Faça um Fork do projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📝 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 🙏 Agradecimentos

- [Cobra](https://github.com/spf13/cobra) - Framework CLI
- [SQLite](https://sqlite.org/) - Banco de dados local
- [Groq](https://groq.com/) - API de IA

## 📞 Suporte

- **Issues**: [GitHub Issues](https://github.com/hudsonrj/SnaAIDatabases/issues)
- **Documentação**: [README.md](README.md)

---

**Feito com ❤️ para DBAs, Desenvolvedores e Equipes de Banco de Dados**

