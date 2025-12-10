#!/bin/bash

# Script de teste de conexão PostgreSQL

echo "🔌 Testando conexão com PostgreSQL..."
echo "Host: 100.123.115.38"
echo "Porta: 5432"
echo "Database: postgres"
echo "Usuário: postgres"
echo ""

# Verificar se psql está disponível
if command -v psql &> /dev/null; then
    echo "✅ psql encontrado, testando conexão..."
    PGPASSWORD=postgres psql -h 100.123.115.38 -p 5432 -U postgres -d postgres -c "SELECT version();" 2>&1 | head -5
else
    echo "⚠️  psql não encontrado. Instalando dependências Go para teste..."
    
    # Tentar compilar e testar conexão diretamente
    cat > /tmp/test_pg_conn.go << 'GOEOF'
package main

import (
    "database/sql"
    "fmt"
    _ "github.com/lib/pq"
)

func main() {
    connStr := "host=100.123.115.38 port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        fmt.Printf("❌ Erro ao abrir conexão: %v\n", err)
        return
    }
    defer db.Close()
    
    err = db.Ping()
    if err != nil {
        fmt.Printf("❌ Erro ao conectar: %v\n", err)
        return
    }
    
    var version string
    err = db.QueryRow("SELECT version();").Scan(&version)
    if err != nil {
        fmt.Printf("❌ Erro ao executar query: %v\n", err)
        return
    }
    
    fmt.Printf("✅ Conexão bem-sucedida!\n")
    fmt.Printf("📊 Versão: %s\n", version)
}
GOEOF
    
    echo "Compilando teste..."
    go run /tmp/test_pg_conn.go 2>&1
fi
