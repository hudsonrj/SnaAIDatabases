# 🚀 Quick Start - SnipAI Databases

Guia rápido para começar a usar o SnipAI Databases.

## ⚡ Instalação Rápida

### Linux/macOS

```bash
# Clone o repositório
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases

# Compile
go build -o snip main.go

# Instale
sudo mv snip /usr/local/bin/
```

### Windows

```powershell
# Clone o repositório
git clone https://github.com/hudsonrj/SnaAIDatabases.git
cd SnaAIDatabases

# Compile
$env:CGO_ENABLED=1
go build -o snip.exe main.go
```

## 🔑 Configurar API Key

```bash
# Linux/macOS
export GROQ_API_KEY="sua_chave_aqui"

# Windows PowerShell
$env:GROQ_API_KEY="sua_chave_aqui"
```

Obtenha sua chave em: https://console.groq.com/keys

## 📝 Primeiro Uso

### 1. Análise Simples

```bash
# Criar análise
snip db-analysis create \
  --title "Minha Primeira Análise" \
  --db-type postgresql \
  --analysis-type diagnostic \
  --host localhost \
  --database mydb \
  --username user

# Executar
snip db-analysis run 1

# Ver resultados
snip db-analysis get 1
```

### 2. Chat Interativo

```bash
snip db-chat \
  --db-type postgresql \
  --host localhost \
  --port 5432 \
  --database mydb \
  --username user
```

### 3. Gerar Gráfico

```bash
snip db-chart --analysis-id 1
```

## 📚 Próximos Passos

- Leia o [README.md](README.md) completo
- Consulte o [INSTALL.md](INSTALL.md) para instalação detalhada
- Explore os exemplos na documentação

## 🆘 Precisa de Ajuda?

- [Documentação Completa](README.md)
- [Guia de Instalação](INSTALL.md)
- [Issues no GitHub](https://github.com/hudsonrj/SnaAIDatabases/issues)

