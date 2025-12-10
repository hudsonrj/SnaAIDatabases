# 📚 Exemplos Detalhados - SnipAI Databases

Este documento contém exemplos detalhados de uso, saídas de comandos e relatórios gerados pelo SnipAI Databases.

## 📋 Índice

1. [Análise Diagnóstica Completa](#análise-diagnóstica-completa)
2. [Análise de Performance com Gráficos](#análise-de-performance-com-gráficos)
3. [Chat Interativo com Banco de Dados](#chat-interativo-com-banco-de-dados)
4. [Chat com Histórico de Análises](#chat-com-histórico-de-análises)
5. [Análise Oracle AWR](#análise-oracle-awr)
6. [Análise de Locks SQL Server](#análise-de-locks-sql-server)
7. [Checklist Diário](#checklist-diário)
8. [Análise de Backup com RTO](#análise-de-backup-com-rto)
9. [Plano de Manutenção Gerado por IA](#plano-de-manutenção-gerado-por-ia)
10. [Exportação para Markdown](#exportação-para-markdown)
11. [Integração com Jira](#integração-com-jira)
12. [Integração com Confluence](#integração-com-confluence)

---

## 1. Análise Diagnóstica Completa

### Comando

```bash
snip db-analysis create \
  --title "Diagnóstico PostgreSQL Produção" \
  --db-type postgresql \
  --analysis-type diagnostic \
  --host db-prod.example.com \
  --port 5432 \
  --database myapp \
  --username dbadmin \
  --password minha_senha_segura

# Saída: Análise criada com ID 1

snip db-analysis run 1
```

### Saída do Comando

```
✓ Análise criada com sucesso!
  ID: 1
  Título: Diagnóstico PostgreSQL Produção
  Tipo: postgresql / diagnostic
  Status: pending

Executando análise...
✓ Conectado ao banco de dados
✓ Executando queries de diagnóstico...
✓ Gerando insights com IA...

✓ Análise concluída com sucesso!
  Status: completed

═══════════════════════════════════════════════════════════════
📊 RESULTADO DA ANÁLISE
═══════════════════════════════════════════════════════════════

# Diagnóstico PostgreSQL - db-prod.example.com:5432/myapp

## Informações Gerais
- **Versão**: PostgreSQL 14.9
- **Servidor**: db-prod.example.com
- **Database**: myapp
- **Data da Análise**: 2024-01-15 14:30:00

## Status do Banco
✅ Banco de dados está ONLINE
✅ Conexões ativas: 45/100 (45%)
⚠️  Conexões próximas do limite recomendado

## Estatísticas de Conexão
- **Total de conexões**: 45
- **Conexões ativas**: 38
- **Conexões idle**: 7
- **Máximo permitido**: 100
- **Conexões por usuário**:
  - app_user: 32
  - admin_user: 8
  - backup_user: 5

## Uso de Recursos
- **CPU**: 65% (moderado)
- **Memória**: 8.2 GB / 16 GB (51%)
- **Disco**: 450 GB / 1 TB (45%)
- **I/O Wait**: 12% (atenção)

## Tabelas Principais
| Tabela | Linhas | Tamanho | Índices |
|--------|--------|---------|---------|
| orders | 2,345,678 | 1.2 GB | 5 |
| customers | 125,432 | 45 MB | 3 |
| products | 8,765 | 12 MB | 2 |
| transactions | 5,678,901 | 2.8 GB | 8 |

## Queries Lentas (Top 5)
| Query | Tempo Médio | Execuções | Total Tempo |
|-------|-------------|-----------|-------------|
| SELECT * FROM orders WHERE... | 2.3s | 1,234 | 2,838s |
| JOIN customers c ON... | 1.8s | 567 | 1,021s |
| UPDATE products SET... | 1.2s | 890 | 1,068s |

## Índices Não Utilizados
⚠️  Encontrados 3 índices não utilizados:
- idx_orders_old_field (último uso: 2023-06-15)
- idx_customers_temp (último uso: nunca)
- idx_products_legacy (último uso: 2023-08-20)

───────────────────────────────────────────────────────────────
🤖 INSIGHTS DA IA
───────────────────────────────────────────────────────────────

## Análise Inteligente

Identifiquei os seguintes pontos de atenção:

### ⚠️ Problemas Identificados

1. **Conexões Próximas do Limite**
   - 45% das conexões disponíveis estão em uso
   - Durante picos, pode atingir o limite
   - **Recomendação**: Aumentar `max_connections` para 150 ou implementar pool de conexões

2. **I/O Wait Elevado (12%)**
   - Indica possível gargalo de disco
   - Pode impactar performance geral
   - **Recomendação**: Verificar configuração de I/O, considerar SSD ou otimizar queries

3. **Queries Lentas Identificadas**
   - Query mais lenta: `SELECT * FROM orders WHERE...` (2.3s médio)
   - **Recomendação**: 
     - Adicionar índice na coluna de filtro
     - Considerar particionamento da tabela `orders`
     - Revisar uso de `SELECT *` (trazer apenas colunas necessárias)

4. **Índices Não Utilizados**
   - 3 índices não são mais utilizados
   - Ocupam espaço e impactam performance de INSERT/UPDATE
   - **Recomendação**: Remover índices não utilizados após validação

### ✅ Pontos Positivos

- Banco de dados está saudável e operacional
- Uso de memória dentro do esperado (51%)
- Espaço em disco adequado (45% utilizado)
- Estrutura de índices bem configurada nas tabelas principais

### 📋 Ações Recomendadas (Prioridade)

**Alta Prioridade:**
1. Otimizar query `SELECT * FROM orders WHERE...` (impacto: alto)
2. Investigar I/O wait elevado (impacto: médio-alto)
3. Remover índices não utilizados (impacto: médio)

**Média Prioridade:**
4. Aumentar limite de conexões ou implementar pool (impacto: médio)
5. Monitorar crescimento da tabela `transactions` (impacto: baixo-médio)

**Baixa Prioridade:**
6. Revisar queries com tempo médio > 1s (impacto: baixo)

### 🔄 Próximos Passos

1. Executar `EXPLAIN ANALYZE` na query mais lenta
2. Verificar configurações de I/O do sistema
3. Agendar janela de manutenção para remoção de índices
4. Configurar alertas para uso de conexões > 80%
```

---

## 2. Análise de Performance com Gráficos

### Comando

```bash
snip db-analysis create \
  --title "Performance Analysis" \
  --db-type postgresql \
  --analysis-type tuning \
  --host localhost \
  --port 5432 \
  --database mydb \
  --username user

snip db-analysis run 1

# Gerar gráfico
snip db-chart --analysis-id 1 --type bar
```

### Saída do Gráfico (ASCII)

```
═══════════════════════════════════════════════════════════════
📈 VISUALIZAÇÃO - Performance por Tabela
═══════════════════════════════════════════════════════════════

Tempo Médio de Query (ms)
───────────────────────────────────────────────────────────────
orders
████████████████████████████████████████████████████ 2300 ms
customers
████████████████████████████████████ 1800 ms
products
████████████████████ 1200 ms
transactions
████████████████████████████████████████████████████████████ 2800 ms
───────────────────────────────────────────────────────────────
    0    500   1000   1500   2000   2500   3000
```

### Saída do Gráfico (HTML)

O gráfico HTML é gerado em `chart.html` e pode ser aberto no navegador:

```html
<!DOCTYPE html>
<html>
<head>
    <title>Performance Analysis - Chart</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
</head>
<body>
    <canvas id="performanceChart"></canvas>
    <script>
        // Gráfico interativo com Chart.js
        // Permite zoom, hover com detalhes, etc.
    </script>
</body>
</html>
```

---

## 3. Chat Interativo com Banco de Dados

### Comando

```bash
snip db-chat \
  --db-type postgresql \
  --host localhost \
  --port 5432 \
  --database mydb \
  --username user \
  --password pass
```

### Exemplo de Conversa

```
🤖 Chat com Banco de Dados iniciado!
Digite suas perguntas ou solicitações. A IA executará queries automaticamente e responderá com os resultados.
Digite 'exit', 'quit' ou 'sair' para sair.

Você: quantas tabelas existem no banco?

🤖 Assistente: Encontrei 25 tabelas no banco de dados 'mydb'. Aqui estão as principais:

| schema | tabela | linhas | tamanho |
|--------|--------|--------|---------|
| public | orders | 2,345,678 | 1.2 GB |
| public | customers | 125,432 | 45 MB |
| public | products | 8,765 | 12 MB |
| public | transactions | 5,678,901 | 2.8 GB |

As tabelas estão distribuídas principalmente no schema 'public', com a tabela 'transactions' sendo a maior em termos de linhas e a tabela 'orders' sendo a maior em tamanho de dados.

Você: qual tabela tem mais registros?

🤖 Assistente: A tabela 'transactions' tem o maior número de registros: **5,678,901 linhas**.

Esta tabela também ocupa 2.8 GB de espaço em disco, sendo a segunda maior em termos de tamanho. A tabela 'orders' tem menos linhas (2,345,678) mas ocupa mais espaço (1.2 GB), o que sugere que cada registro em 'orders' é maior em média.

Você: mostre as 5 queries mais lentas

🤖 Assistente: [Query gerada automaticamente]
```sql
SELECT 
    query,
    calls,
    total_time,
    mean_time,
    max_time
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 5
```

[Resultado da query executada]

| query | calls | total_time | mean_time | max_time |
|-------|-------|------------|-----------|----------|
| SELECT * FROM orders WHERE customer_id = $1 | 1,234 | 2,838s | 2.3s | 4.5s |
| SELECT c.*, o.* FROM customers c JOIN orders o... | 567 | 1,021s | 1.8s | 3.2s |
| UPDATE products SET stock = $1 WHERE id = $2 | 890 | 1,068s | 1.2s | 2.1s |

**Análise das Queries Mais Lentas:**

1. **Query mais lenta**: `SELECT * FROM orders WHERE customer_id = $1`
   - Tempo médio: 2.3 segundos
   - Executada 1,234 vezes
   - **Recomendação**: Adicionar índice em `customer_id` se não existir, ou revisar o uso de `SELECT *`

2. **Segunda mais lenta**: JOIN entre customers e orders
   - Tempo médio: 1.8 segundos
   - **Recomendação**: Verificar se os índices de foreign key estão otimizados

3. **Terceira mais lenta**: UPDATE em products
   - Tempo médio: 1.2 segundos
   - **Recomendação**: Verificar locks e considerar batch updates

Você: exit

Até logo! 👋
```

---

## 4. Chat com Histórico de Análises

### Comando

```bash
snip db-history chat
```

### Exemplo de Conversa

```
🤖 Chat com Histórico de Análises iniciado!
Digite suas perguntas sobre as análises armazenadas.
A IA executará queries automaticamente e responderá com os resultados.
Digite 'exit', 'quit' ou 'sair' para sair.

Você: quantas análises foram feitas este mês?

🤖 Assistente: [Query gerada e executada automaticamente]

Encontrei **12 análises** realizadas neste mês (Janeiro 2024):

| id | título | tipo_banco | tipo_análise | data |
|----|--------|------------|--------------|------|
| 15 | Diagnóstico PostgreSQL Produção | postgresql | diagnostic | 2024-01-15 14:30:00 |
| 14 | Análise de Performance | postgresql | tuning | 2024-01-14 10:20:00 |
| 13 | Checklist Diário MySQL | mysql | checklist | 2024-01-13 08:00:00 |
| ... | ... | ... | ... | ... |

**Distribuição por tipo de banco:**
- PostgreSQL: 8 análises
- MySQL: 3 análises
- Oracle: 1 análise

**Distribuição por tipo de análise:**
- Diagnostic: 5 análises
- Tuning: 4 análises
- Checklist: 3 análises

Você: compare as análises de diagnóstico do PostgreSQL entre janeiro e fevereiro

🤖 Assistente: [Query gerada para comparar análises]

**Comparação de Análises Diagnósticas - PostgreSQL**

**Janeiro 2024:**
- Total de análises: 5
- Status: 4 completed, 1 error
- Principais problemas identificados:
  - I/O wait elevado (12%)
  - Queries lentas em 3 análises
  - Conexões próximas do limite em 2 análises

**Fevereiro 2024:**
- Total de análises: 3
- Status: 3 completed
- Principais problemas identificados:
  - I/O wait melhorou para 8%
  - Queries otimizadas (tempo médio reduzido em 30%)
  - Conexões estáveis (média de 35%)

**Evolução:**
✅ **Melhorias identificadas:**
- I/O wait reduziu de 12% para 8% (melhoria de 33%)
- Performance de queries melhorou significativamente
- Uso de conexões mais estável

⚠️ **Atenção:**
- Número de análises reduziu (pode indicar menos monitoramento)
- Ainda há espaço para otimização de I/O

**Recomendação**: Continuar monitoramento regular e manter as otimizações implementadas.

Você: quais problemas foram identificados nas análises do SQL Server?

🤖 Assistente: [Query gerada para buscar problemas em análises SQL Server]

Encontrei **7 problemas principais** identificados nas análises do SQL Server:

1. **Deadlocks frequentes** (identificado em 3 análises)
   - Ocorrências: 12 deadlocks na última semana
   - Tabelas envolvidas: orders, inventory
   - **Solução sugerida**: Revisar ordem de locks, implementar retry logic

2. **Bloqueios prolongados** (identificado em 5 análises)
   - Tempo médio de bloqueio: 45 segundos
   - **Solução sugerida**: Otimizar queries longas, revisar isolamento de transações

3. **Espaço em disco insuficiente** (identificado em 2 análises)
   - Uso atual: 92% do disco
   - **Solução sugerida**: Limpar logs antigos, expandir storage

4. **Backup incompleto** (identificado em 1 análise)
   - Último backup completo: há 5 dias
   - **Solução sugerida**: Verificar agendamento de backups, investigar falhas

5. **Queries sem índice** (identificado em 4 análises)
   - 8 queries identificadas usando table scans
   - **Solução sugerida**: Criar índices apropriados

6. **Memória insuficiente** (identificado em 2 análises)
   - Uso de memória: 94%
   - **Solução sugerida**: Aumentar memória alocada ao SQL Server

7. **Replicação com lag** (identificado em 1 análise)
   - Lag de replicação: 15 minutos
   - **Solução sugerida**: Investigar latência de rede, otimizar queries no servidor principal

**Priorização:**
- 🔴 **Crítica**: Espaço em disco, backup incompleto
- 🟡 **Alta**: Deadlocks, bloqueios prolongados
- 🟢 **Média**: Queries sem índice, memória, replicação

Você: exit

Até logo! 👋
```

---

## 5. Análise Oracle AWR

### Comando

```bash
snip db-analysis create \
  --title "AWR Período Crítico" \
  --db-type oracle \
  --analysis-type awr \
  --host oracle-prod \
  --port 1521 \
  --database ORCL \
  --username sys \
  --password senha \
  --ai-query "gerar AWR das últimas 24 horas, focando em picos de carga"

snip db-analysis run 1
```

### Saída

```
═══════════════════════════════════════════════════════════════
📊 RESULTADO DA ANÁLISE - Oracle AWR
═══════════════════════════════════════════════════════════════

# Automatic Workload Repository (AWR) Report
# Período: 2024-01-15 00:00:00 até 2024-01-16 00:00:00
# Snapshots: 12345 (início) até 12350 (fim)

## Resumo Executivo
- **DB Time**: 1,234.56 horas
- **DB CPU**: 890.12 horas
- **Elapsed Time**: 24.00 horas
- **DB Time / Elapsed**: 51.44 (indica alta carga)

## Top 10 Wait Events
| Evento | Total Waits | Time Waited (s) | Avg Wait (ms) | % DB Time |
|--------|-------------|-----------------|---------------|-----------|
| db file sequential read | 1,234,567 | 45,678.90 | 37.0 | 12.5% |
| CPU time | - | 3,204,432.00 | - | 72.1% |
| db file scattered read | 567,890 | 12,345.67 | 21.7 | 3.4% |
| log file sync | 123,456 | 8,901.23 | 72.1 | 2.4% |

## Top 10 SQL por Tempo de Execução
| SQL ID | Tempo Total (s) | Execuções | Tempo Médio (s) | % DB Time |
|--------|-----------------|-----------|------------------|-----------|
| abc123xyz | 12,345.67 | 1,234 | 10.01 | 2.8% |
| def456uvw | 8,901.23 | 567 | 15.69 | 2.0% |
| ghi789rst | 6,789.01 | 890 | 7.62 | 1.5% |

───────────────────────────────────────────────────────────────
🤖 INSIGHTS DA IA
───────────────────────────────────────────────────────────────

## Análise do Período

### ⚠️ Problemas Críticos Identificados

1. **Alto DB Time / Elapsed (51.44)**
   - Indica que o banco está sob alta carga
   - Múltiplas sessões competindo por recursos
   - **Recomendação**: Investigar queries que estão consumindo mais tempo

2. **Top Wait Event: db file sequential read (12.5% do DB Time)**
   - 1.2 milhões de waits com tempo médio de 37ms
   - Indica possível falta de índices ou I/O lento
   - **Recomendação**: 
     - Verificar queries relacionadas ao SQL ID abc123xyz
     - Considerar adicionar índices apropriados
     - Verificar performance de I/O do storage

3. **SQL ID abc123xyz - Maior Consumidor de Tempo**
   - 12,345 segundos de tempo total
   - 1,234 execuções com tempo médio de 10 segundos
   - **Recomendação**: 
     - Revisar plano de execução
     - Considerar otimização ou reescrita da query
     - Verificar se índices estão sendo utilizados

### ✅ Pontos Positivos

- CPU time representa 72.1% do DB Time (boa utilização de CPU)
- Log file sync com tempo médio aceitável (72ms)
- Nenhum wait event crítico relacionado a locks ou deadlocks

### 📋 Ações Recomendadas

**Imediatas:**
1. Analisar plano de execução do SQL ID abc123xyz
2. Verificar índices das tabelas envolvidas
3. Monitorar I/O do storage durante picos

**Curto Prazo:**
4. Otimizar ou reescrever queries do Top 10
5. Considerar particionamento de tabelas grandes
6. Revisar estratégia de cache

**Médio Prazo:**
7. Implementar monitoramento contínuo de wait events
8. Criar alertas para DB Time / Elapsed > 50
9. Documentar queries críticas e seus planos de execução
```

---

## 6. Análise de Locks SQL Server

### Comando

```bash
snip db-analysis create \
  --title "Análise de Locks SQL Server" \
  --db-type sqlserver \
  --analysis-type locks \
  --host sqlserver-prod \
  --port 1433 \
  --database AdventureWorks \
  --username sa \
  --password senha123

snip db-analysis run 2
```

### Saída

```
═══════════════════════════════════════════════════════════════
📊 RESULTADO DA ANÁLISE - SQL Server Locks
═══════════════════════════════════════════════════════════════

# Análise de Locks e Bloqueios
# Data: 2024-01-15 15:45:00

## Locks Ativos
Total de locks: 156

### Por Tipo de Lock
| Tipo | Quantidade | % do Total |
|------|------------|------------|
| KEY | 89 | 57.1% |
| PAGE | 45 | 28.8% |
| OBJECT | 22 | 14.1% |

### Por Modo de Lock
| Modo | Quantidade | Descrição |
|------|------------|-----------|
| S (Shared) | 78 | Leitura compartilhada |
| X (Exclusive) | 34 | Escrita exclusiva |
| IS (Intent Shared) | 28 | Intenção de leitura |
| IX (Intent Exclusive) | 16 | Intenção de escrita |

## Bloqueios Identificados
⚠️  Encontrados 3 bloqueios ativos:

### Bloqueio 1
- **Sessão Bloqueada**: SPID 67
- **Sessão Bloqueante**: SPID 89
- **Tempo de Espera**: 45 segundos
- **Recurso**: KEY (1:12345:1)
- **Query Bloqueada**:
  ```sql
  SELECT * FROM orders WHERE customer_id = 12345
  ```
- **Query Bloqueante**:
  ```sql
  UPDATE orders SET status = 'processed' WHERE order_id = 67890
  ```

### Bloqueio 2
- **Sessão Bloqueada**: SPID 123
- **Sessão Bloqueante**: SPID 89
- **Tempo de Espera**: 23 segundos
- **Recurso**: PAGE (1:5678)
- **Query Bloqueada**:
  ```sql
  SELECT COUNT(*) FROM transactions WHERE date = '2024-01-15'
  ```

### Bloqueio 3
- **Sessão Bloqueada**: SPID 156
- **Sessão Bloqueante**: SPID 234
- **Tempo de Espera**: 12 segundos
- **Recurso**: OBJECT (AdventureWorks.orders)

## Deadlocks Recentes
⚠️  Encontrados 2 deadlocks nas últimas 24 horas:

### Deadlock 1 (2024-01-15 14:30:00)
- **Sessões envolvidas**: SPID 67 e SPID 89
- **Tabelas**: orders, inventory
- **Causa**: Conflito de locks em ordem diferente
- **Resolvido**: Sim (Oracle escolheu vítima)

### Deadlock 2 (2024-01-15 13:15:00)
- **Sessões envolvidas**: SPID 123 e SPID 234
- **Tabelas**: customers, orders
- **Causa**: Transações longas com múltiplos updates

───────────────────────────────────────────────────────────────
🤖 INSIGHTS DA IA
───────────────────────────────────────────────────────────────

## Análise de Bloqueios

### ⚠️ Problemas Críticos

1. **Bloqueio Prolongado (45 segundos)**
   - SPID 67 aguardando SPID 89 há 45 segundos
   - Query bloqueada é uma SELECT simples
   - Query bloqueante é um UPDATE
   - **Ação Imediata**: 
     ```sql
     -- Verificar status da sessão bloqueante
     SELECT * FROM sys.dm_exec_requests WHERE session_id = 89
     
     -- Se necessário, matar a sessão (com cuidado!)
     KILL 89
     ```

2. **Deadlocks Frequentes**
   - 2 deadlocks em 24 horas é um sinal de alerta
   - Padrão: conflitos entre orders e inventory/customers
   - **Recomendação**: 
     - Padronizar ordem de locks em todas as transações
     - Implementar retry logic com backoff exponencial
     - Considerar isolamento de transação mais baixo (READ COMMITTED SNAPSHOT)

3. **Query Bloqueante Problemática**
   - SPID 89 está bloqueando múltiplas sessões
   - UPDATE em orders sem WHERE específico pode estar causando lock de tabela
   - **Recomendação**: 
     - Revisar a query do SPID 89
     - Adicionar WHERE mais específico
     - Considerar batch updates em vez de atualização em massa

### 📋 Recomendações de Longo Prazo

1. **Implementar READ COMMITTED SNAPSHOT**
   ```sql
   ALTER DATABASE AdventureWorks SET READ_COMMITTED_SNAPSHOT ON
   ```
   - Reduz bloqueios de leitura
   - Melhora concorrência

2. **Otimizar Queries Longas**
   - Revisar queries com tempo de execução > 5 segundos
   - Adicionar índices apropriados
   - Considerar particionamento

3. **Monitoramento Contínuo**
   - Configurar alertas para bloqueios > 30 segundos
   - Monitorar deadlocks diariamente
   - Documentar padrões de bloqueio

### 🔧 Queries Úteis

```sql
-- Ver todas as sessões bloqueadas
SELECT 
    r.session_id AS blocked_session,
    r.blocking_session_id AS blocking_session,
    r.wait_type,
    r.wait_time,
    t.text AS query_text
FROM sys.dm_exec_requests r
CROSS APPLY sys.dm_exec_sql_text(r.sql_handle) t
WHERE r.blocking_session_id > 0

-- Ver deadlocks recentes
SELECT * FROM sys.dm_xe_session_events
WHERE event_name = 'xml_deadlock_report'
ORDER BY timestamp DESC
```
```

---

## 7. Checklist Diário

### Comando

```bash
snip db-analysis create \
  --title "Checklist Diário MySQL" \
  --db-type mysql \
  --analysis-type checklist \
  --host localhost \
  --port 3306 \
  --database mydb \
  --username root \
  --password senha \
  --checklist-type daily

snip db-analysis run 3
```

### Saída

```
═══════════════════════════════════════════════════════════════
📊 RESULTADO DA ANÁLISE - Checklist Diário MySQL
═══════════════════════════════════════════════════════════════

# Checklist Diário - MySQL
# Data: 2024-01-15 08:00:00

## ✅ Verificações Realizadas

### 1. Espaço em Disco ✅
- **Status**: OK
- **Uso**: 450 GB / 1 TB (45%)
- **Espaço Livre**: 550 GB
- **Tendência**: Estável (crescimento de 2 GB/dia)
- **Ação**: Nenhuma necessária

### 2. Status de Replicação ✅
- **Status**: OK
- **Master**: ONLINE
- **Slave**: ONLINE
- **Lag**: 0 segundos
- **Última sincronização**: 2024-01-15 07:59:58
- **Ação**: Nenhuma necessária

### 3. Queries Lentas ⚠️
- **Status**: ATENÇÃO
- **Queries > 2s**: 3 encontradas
- **Tempo total**: 12.5 segundos
- **Ação**: Revisar queries identificadas

**Queries Lentas:**
| Query | Tempo | Execuções |
|-------|-------|-----------|
| SELECT * FROM orders WHERE... | 3.2s | 45 |
| UPDATE products SET... | 2.8s | 12 |
| DELETE FROM logs WHERE... | 2.5s | 8 |

### 4. Locks Ativos ✅
- **Status**: OK
- **Locks ativos**: 5
- **Deadlocks (24h)**: 0
- **Bloqueios prolongados**: 0
- **Ação**: Nenhuma necessária

### 5. Status de Backups ✅
- **Status**: OK
- **Último backup completo**: 2024-01-15 02:00:00 (há 6 horas)
- **Último backup incremental**: 2024-01-15 08:00:00 (há 0 horas)
- **Tamanho backup completo**: 45 GB
- **Duração**: 1h 23min
- **Ação**: Nenhuma necessária

### 6. Conexões Ativas ✅
- **Status**: OK
- **Conexões ativas**: 45 / 200 (22.5%)
- **Conexões máximas**: 200
- **Tendência**: Estável
- **Ação**: Nenhuma necessária

### 7. Uso de Memória ✅
- **Status**: OK
- **Uso**: 8.2 GB / 16 GB (51%)
- **Buffer Pool**: 6.5 GB
- **Query Cache**: 512 MB
- **Ação**: Nenhuma necessária

### 8. Fragmentação de Tabelas ⚠️
- **Status**: ATENÇÃO
- **Tabelas fragmentadas**: 3
- **Fragmentação média**: 18%
- **Ação**: Agendar OPTIMIZE TABLE

**Tabelas Fragmentadas:**
| Tabela | Fragmentação | Tamanho |
|--------|--------------|---------|
| orders | 25% | 1.2 GB |
| transactions | 20% | 2.8 GB |
| logs | 15% | 450 MB |

### 9. Índices Não Utilizados ✅
- **Status**: OK
- **Índices não utilizados**: 0
- **Ação**: Nenhuma necessária

### 10. Erros no Log ⚠️
- **Status**: ATENÇÃO
- **Erros (24h)**: 5
- **Warnings (24h)**: 12
- **Ação**: Revisar log de erros

**Erros Encontrados:**
- 3 erros de conexão timeout
- 2 erros de query timeout

───────────────────────────────────────────────────────────────
🤖 INSIGHTS DA IA
───────────────────────────────────────────────────────────────

## Resumo do Checklist

### ✅ Status Geral: BOM

A maioria dos itens está em ordem. O banco de dados está operacional e saudável.

### ⚠️ Itens que Requerem Atenção

1. **Queries Lentas (3 encontradas)**
   - Impacto: Médio
   - **Ação Recomendada**: 
     - Executar EXPLAIN nas queries identificadas
     - Adicionar índices se necessário
     - Considerar otimização de queries

2. **Fragmentação de Tabelas (3 tabelas)**
   - Impacto: Baixo-Médio
   - **Ação Recomendada**: 
     - Agendar OPTIMIZE TABLE durante janela de manutenção
     - Priorizar tabela 'orders' (25% de fragmentação)

3. **Erros no Log (5 erros)**
   - Impacto: Baixo
   - **Ação Recomendada**: 
     - Investigar erros de timeout
     - Verificar configurações de timeout
     - Revisar carga do servidor durante picos

### 📋 Plano de Ação

**Hoje:**
- [ ] Revisar queries lentas identificadas
- [ ] Investigar erros de timeout no log

**Esta Semana:**
- [ ] Agendar OPTIMIZE TABLE para tabelas fragmentadas
- [ ] Implementar índices recomendados para queries lentas

**Este Mês:**
- [ ] Revisar estratégia de timeout
- [ ] Monitorar tendências de fragmentação
```

---

## 8. Análise de Backup com RTO

### Comando

```bash
snip db-analysis create \
  --title "Status Backups SQL Server" \
  --db-type sqlserver \
  --analysis-type backup \
  --host sqlserver-prod \
  --port 1433 \
  --database master \
  --username backup_admin \
  --password senha

snip db-analysis run 4
```

### Saída

```
═══════════════════════════════════════════════════════════════
📊 RESULTADO DA ANÁLISE - Status de Backups
═══════════════════════════════════════════════════════════════

# Análise de Backups - SQL Server
# Data: 2024-01-15 10:00:00

## Backups Encontrados

### Backup Completo (Full)
- **Último backup**: 2024-01-15 02:00:00 (há 8 horas)
- **Tamanho**: 45.2 GB
- **Duração**: 1h 23min
- **Status**: ✅ Sucesso
- **Localização**: \\backup-server\sql-backups\full\DB_20240115_020000.bak
- **Tipo de compressão**: Sim (taxa: 60%)

### Backup Diferencial (Differential)
- **Último backup**: 2024-01-15 08:00:00 (há 2 horas)
- **Tamanho**: 2.1 GB
- **Duração**: 8min 45s
- **Status**: ✅ Sucesso
- **Localização**: \\backup-server\sql-backups\diff\DB_20240115_080000.bak

### Backup de Log (Transaction Log)
- **Último backup**: 2024-01-15 10:00:00 (há 0 horas)
- **Tamanho**: 125 MB
- **Duração**: 45s
- **Status**: ✅ Sucesso
- **Frequência**: A cada 15 minutos
- **Localização**: \\backup-server\sql-backups\log\DB_20240115_100000.trn

## Histórico de Backups (Últimos 7 dias)

| Data | Tipo | Tamanho | Duração | Status |
|------|------|---------|---------|--------|
| 2024-01-15 02:00 | Full | 45.2 GB | 1h 23min | ✅ |
| 2024-01-15 08:00 | Differential | 2.1 GB | 8min 45s | ✅ |
| 2024-01-14 02:00 | Full | 44.8 GB | 1h 20min | ✅ |
| 2024-01-13 02:00 | Full | 44.5 GB | 1h 18min | ✅ |
| ... | ... | ... | ... | ... |

## Estratégia de Backup Atual

- **Full Backup**: Diário às 02:00
- **Differential Backup**: A cada 6 horas
- **Transaction Log Backup**: A cada 15 minutos
- **Retenção**: 30 dias para Full, 7 dias para Differential, 2 dias para Log

───────────────────────────────────────────────────────────────
🤖 INSIGHTS DA IA
───────────────────────────────────────────────────────────────

## Análise de Backups e Estimativa de RTO

### ✅ Status Geral: EXCELENTE

A estratégia de backup está bem configurada e funcionando corretamente.

### 📊 Estimativa de RTO (Recovery Time Objective)

Com base no histórico de backups e tamanho dos dados, a IA estima:

**Cenário 1: Falha Completa (Point-in-Time Recovery)**
- **RTO Estimado**: 2h 15min - 3h 30min
- **Processo**:
  1. Restaurar Full Backup mais recente: ~1h 23min
  2. Restaurar Differential Backup: ~8min 45s
  3. Restaurar Transaction Logs (últimas 2 horas): ~15-20min
  4. Verificação e validação: ~30-45min
- **RPO (Recovery Point Objective)**: Máximo de 15 minutos de perda de dados

**Cenário 2: Restauração de Database Específica**
- **RTO Estimado**: 45min - 1h 15min
- **Processo mais rápido** se apenas um database específico precisar ser restaurado

**Cenário 3: Restauração de Arquivo/Filegroup**
- **RTO Estimado**: 20min - 45min
- **Processo mais rápido** para falhas parciais

### 📋 Fatores que Impactam o RTO

**Positivos:**
- ✅ Backups completos diários (boa frequência)
- ✅ Backups diferenciais a cada 6 horas (reduz tempo de restauração)
- ✅ Transaction logs a cada 15 minutos (RPO baixo)
- ✅ Compressão de backups (60% de redução acelera transferência)
- ✅ Backups em localização de rede dedicada (boa performance)

**Melhorias Possíveis:**
- ⚠️ Considerar backup completo mais frequente (a cada 12 horas) para RTO menor
- ⚠️ Implementar backup em múltiplos locais (redundância)
- ⚠️ Considerar backup incremental para databases muito grandes
- ⚠️ Testar restauração regularmente para validar RTO real

### 🔧 Recomendações

**Curto Prazo:**
1. Testar processo de restauração completo (DR Drill)
2. Documentar procedimento de restauração passo a passo
3. Validar integridade dos backups regularmente

**Médio Prazo:**
4. Considerar aumentar frequência de backups diferenciais (a cada 4 horas)
5. Implementar monitoramento de espaço em disco de backup
6. Configurar alertas para falhas de backup

**Longo Prazo:**
7. Avaliar estratégia de backup contínuo (se necessário RTO < 1h)
8. Considerar backup em cloud para disaster recovery
9. Implementar backup automatizado de configurações do servidor

### 📝 Checklist de Validação de Backup

- [ ] Backups completos executando conforme agendamento
- [ ] Backups diferenciais executando conforme agendamento
- [ ] Transaction logs sendo fechados corretamente
- [ ] Espaço em disco de backup suficiente (atualmente: 2 TB disponível)
- [ ] Último teste de restauração: 2024-01-10 (há 5 dias) ⚠️ **Recomendado testar novamente**
- [ ] Documentação de procedimento de restauração atualizada
```

---

## 9. Plano de Manutenção Gerado por IA

### Comando

```bash
snip db-maintenance --analysis-id 1 --output maintenance-plan.md
```

### Saída (Arquivo Markdown)

```markdown
# Plano de Manutenção - Diagnóstico PostgreSQL Produção

**Gerado em**: 2024-01-15 15:00:00  
**Baseado em**: Análise ID 1  
**Prioridade Geral**: Média-Alta

---

## 📋 Tarefas de Manutenção

### 🔴 Alta Prioridade

#### Tarefa 1: Otimizar Query Mais Lenta
- **Descrição**: Otimizar a query `SELECT * FROM orders WHERE customer_id = $1` que está levando em média 2.3 segundos
- **Tempo Estimado**: 2-3 horas
- **Passo a Passo**:
  1. Executar `EXPLAIN ANALYZE` na query para identificar o plano de execução atual
  2. Verificar se existe índice em `customer_id` na tabela `orders`
  3. Se não existir, criar índice: `CREATE INDEX idx_orders_customer_id ON orders(customer_id)`
  4. Se existir, verificar se está sendo utilizado (pode estar desabilitado ou fragmentado)
  5. Revisar uso de `SELECT *` - considerar trazer apenas colunas necessárias
  6. Testar performance após otimização
  7. Monitorar impacto em produção
- **Dependências**: Nenhuma
- **Risco**: Baixo (criação de índice pode ser feita online)
- **Janela Recomendada**: Qualquer horário (índice pode ser criado online)

#### Tarefa 2: Investigar I/O Wait Elevado
- **Descrição**: I/O wait de 12% indica possível gargalo de disco que precisa ser investigado
- **Tempo Estimado**: 3-4 horas
- **Passo a Passo**:
  1. Verificar métricas de I/O do sistema operacional (iostat, sar)
  2. Identificar quais tabelas/arquivos estão gerando mais I/O
  3. Verificar se há queries fazendo table scans desnecessários
  4. Revisar configurações de PostgreSQL relacionadas a I/O (shared_buffers, effective_cache_size)
  5. Verificar se o storage está respondendo adequadamente
  6. Considerar otimizações de queries que estão gerando I/O excessivo
  7. Se necessário, escalar storage ou considerar SSD
- **Dependências**: Nenhuma
- **Risco**: Médio (mudanças de configuração podem requerer restart)
- **Janela Recomendada**: Janela de manutenção agendada

### 🟡 Média Prioridade

#### Tarefa 3: Remover Índices Não Utilizados
- **Descrição**: Remover 3 índices não utilizados que estão ocupando espaço e impactando performance de INSERT/UPDATE
- **Tempo Estimado**: 1 hora
- **Passo a Passo**:
  1. Validar que os índices realmente não são utilizados (verificar pg_stat_user_indexes)
  2. Documentar os índices a serem removidos
  3. Agendar janela de manutenção
  4. Remover índices um por um: `DROP INDEX idx_orders_old_field`
  5. Monitorar impacto após remoção
  6. Verificar se espaço foi liberado
- **Dependências**: Nenhuma
- **Risco**: Baixo (após validação)
- **Janela Recomendada**: Janela de manutenção agendada

#### Tarefa 4: Aumentar Limite de Conexões
- **Descrição**: Aumentar `max_connections` de 100 para 150 ou implementar pool de conexões
- **Tempo Estimado**: 2 horas
- **Passo a Passo**:
  1. Avaliar uso atual de memória por conexão
  2. Calcular impacto de aumentar conexões na memória total
  3. Se memória permitir, aumentar `max_connections` no postgresql.conf
  4. Se memória não permitir, considerar implementar pool (PgBouncer)
  5. Reiniciar PostgreSQL (se mudança direta) ou configurar pool
  6. Monitorar uso de conexões após mudança
- **Dependências**: Nenhuma
- **Risco**: Médio (mudança requer restart)
- **Janela Recomendada**: Janela de manutenção agendada

### 🟢 Baixa Prioridade

#### Tarefa 5: Revisar Queries com Tempo Médio > 1s
- **Descrição**: Revisar e otimizar queries identificadas com tempo médio superior a 1 segundo
- **Tempo Estimado**: 4-6 horas
- **Passo a Passo**:
  1. Listar todas as queries com tempo médio > 1s
  2. Priorizar por frequência de execução e tempo total
  3. Para cada query:
     - Executar EXPLAIN ANALYZE
     - Identificar gargalos
     - Propor otimizações (índices, reescrita, etc.)
     - Implementar e testar
  4. Documentar otimizações realizadas
- **Dependências**: Nenhuma
- **Risco**: Baixo
- **Janela Recomendada**: Trabalho contínuo, pode ser feito incrementalmente

---

## 📅 Cronograma Sugerido

### Semana 1
- **Segunda-feira**: Tarefa 1 (Otimizar Query Mais Lenta)
- **Quarta-feira**: Tarefa 2 (Investigar I/O Wait)

### Semana 2
- **Segunda-feira**: Tarefa 3 (Remover Índices)
- **Quarta-feira**: Tarefa 4 (Aumentar Conexões)

### Semana 3-4
- **Contínuo**: Tarefa 5 (Revisar Queries)

---

## ⚠️ Riscos e Mitigações

### Riscos Identificados
1. **Mudanças de Configuração**: Algumas mudanças requerem restart do PostgreSQL
   - **Mitigação**: Agendar em janela de manutenção, ter plano de rollback

2. **Remoção de Índices**: Pode impactar queries não monitoradas
   - **Mitigação**: Validar extensivamente antes de remover, manter backup dos índices

3. **Aumento de Conexões**: Pode impactar memória
   - **Mitigação**: Calcular impacto antes, monitorar após mudança

---

## 📊 Métricas de Sucesso

Após implementação, espera-se:
- Redução de tempo médio da query mais lenta: de 2.3s para < 1s
- Redução de I/O wait: de 12% para < 8%
- Aumento de capacidade de conexões: de 100 para 150
- Liberação de espaço: ~500 MB (remoção de índices)
- Melhoria geral de performance: 20-30%

---

## 🔄 Revisão e Acompanhamento

- **Revisar plano**: Após 2 semanas da implementação
- **Executar nova análise**: Após 1 mês
- **Comparar resultados**: Com análise inicial
- **Ajustar plano**: Conforme necessário
```

---

## 10. Exportação para Markdown

### Comando

```bash
snip db-analysis get 1 --export "relatorio-diagnostico.md"
```

### Saída do Arquivo

O arquivo `relatorio-diagnostico.md` é criado em `~/.snip/exports/` e contém:

```markdown
# Diagnóstico PostgreSQL Produção

**Tipo de Banco:** postgresql  
**Tipo de Análise:** diagnostic  
**Data:** 2024-01-15 14:30:00  

---

## 📊 Resultado da Análise

[Todo o conteúdo da análise formatado em Markdown]

## 📈 Visualização

[Gráficos e visualizações se houver]

## 🤖 Insights da IA

[Insights e recomendações da IA]
```

---

## 11. Integração com Jira

### Comando

```bash
# Configurar Jira
snip jira config \
  --base-url "https://empresa.atlassian.net" \
  --username "seu-email@empresa.com" \
  --api-key "sua-api-key" \
  --project-key "PROJ"

# Criar Epic a partir de análise
snip db-analysis run 1 --create-jira-epic

# Criar Issues a partir de problemas identificados
snip db-analysis run 1 --create-jira-issues
```

### Saída

```
✓ Configuração do Jira salva com sucesso!

Executando análise...
✓ Análise concluída!

📝 Criando Epic no Jira...
✓ Epic criado: PROJ-123

📝 Criando Issues no Jira...
✓ 3 Issue(s) criada(s) no Epic PROJ-123
   - PROJ-124: Otimizar Query Mais Lenta
   - PROJ-125: Investigar I/O Wait Elevado
   - PROJ-126: Remover Índices Não Utilizados
```

---

## 12. Integração com Confluence

### Comando

```bash
# Configurar Confluence
snip confluence config \
  --base-url "https://empresa.atlassian.net/wiki" \
  --username "seu-email@empresa.com" \
  --api-key "sua-api-key" \
  --space-key "DB"

# Exportar análise para Confluence
snip db-analysis get 1 --export-confluence "Análise PostgreSQL - Jan 2024"
```

### Saída

```
✓ Configuração do Confluence salva com sucesso!

Exportando análise para Confluence...
✓ Página Confluence 'Análise PostgreSQL - Jan 2024' criada com sucesso! ID: 123456
Link: https://empresa.atlassian.net/wiki/spaces/DB/pages/123456
```

---

## 📝 Notas Finais

Estes exemplos demonstram a versatilidade e poder do SnipAI Databases. A integração com IA torna as análises mais inteligentes e acionáveis, enquanto a exportação e integração com outras ferramentas facilita o compartilhamento e rastreamento de ações.

Para mais informações, consulte o [README.md](README.md) principal.

