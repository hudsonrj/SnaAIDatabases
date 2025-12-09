![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white) ![SQLite](https://img.shields.io/badge/SQLite-003B57?style=for-the-badge&logo=sqlite&logoColor=white) ![License](https://img.shields.io/badge/license-MIT-green?style=for-the-badge) ![Version](https://img.shields.io/badge/version-1.1.0-blue?style=for-the-badge) ![GitHub stars](https://img.shields.io/github/stars/matheuzgomes/snip?style=for-the-badge&label=Stars)




<div align="center" style="margin-bottom: 15px; display: flex; align-items: center; justify-content: center; gap: 15px;">
  <img src="assets/snip_logo.png" alt="Snip Logo" width="120" height="130" style="border-radius: 16px; border: 2px solid #e0e0e0;">
  <h1 style="margin: 0;">Snip</h1>
</div>

A fast and efficient command-line note-taking tool built with Go. Snip helps you capture, organize, and search your notes with AI-powered features, project management, tasks, and checklists.

## 🎬 Demo

![Snip Demo](assets/snip_demo.gif)

## ✨ Features

### 📝 Notes Management

- **Create Notes**: Quickly create new notes with title and content
- **List Notes**: View all your notes with chronological sorting options
- **Search Notes**: Full-text search across all notes using SQLite FTS4
- **Edit Notes**: Update existing notes using your preferred editor
- **Get Notes**: Retrieve specific notes by ID with markdown rendering support
- **Delete Notes**: Remove notes you no longer need
- **Tags**: Organize notes with custom tags
- **Patch Notes**: Update note titles and manage tags
- **Export Notes**: Export notes to JSON and Markdown formats
- **Import Notes**: Import notes (markdown) from files and directories
- **Markdown Preview**: Render markdown content beautifully in the terminal
- **Fast Performance**: SQLite database with optimized indexes (90-127ns operations)
- **Editor Integration**: Supports nano, vim, vi, or custom `$EDITOR`
- **Comprehensive Testing**: Full test coverage with performance benchmarks

### 🤖 AI-Powered Features

- **AI Create Notes**: Generate notes with AI-powered content based on topics
- **AI Code Generation**: Generate code in multiple languages with AI
- **AI Search Enhancement**: Improve search queries using AI
- **AI Q&A**: Ask questions to AI based on your notes context
- **AI Project Planning**: Generate detailed project plans with AI
- **AI Checklist Generation**: Create checklists with AI-generated items

### 📁 Project Management

- **Projects**: Create and manage projects with descriptions and status
- **Tasks**: Create tasks within projects with priorities and due dates
- **Task Status**: Track tasks (pending, in_progress, completed)
- **Task Priorities**: Set task priorities (low, medium, high)
- **Checklists**: Create checklists for projects or tasks
- **Checklist Items**: Manage checklist items with completion tracking
- **Progress Tracking**: Visual progress indicators for checklists

### Command Examples

#### 📝 Basic Notes

```bash
# Create a new note
snip create "Meeting Notes"

# Create a new note quickly
snip create "World" --message "Hello!"

# List all notes (newest first)
snip list

# List notes chronologically (oldest first)
snip list --asc

# List with verbose information
snip list --verbose

# Search for notes containing specific terms
snip find "meeting"

# Edit an existing note
snip update 1

# Get a specific note by ID
snip show 1

# Get a note with markdown rendering
snip show 1 --render

# Delete a specific note by ID
snip delete 1

# Patch/update a note's title
snip patch 1 --title "New Title"

# Patch/update a note's tags
snip patch 1 --tag "work important"

# List notes with tags
snip list --tag "work"

# Export notes to JSON format
snip export --format json

# Export notes to Markdown format
snip export --format markdown

# Export notes created since a specific date
snip export --since "2024-01-01"

# Import notes from a directory
snip import /path/to/notes/directory

# Show editor information and available options
snip editor
```

#### 🤖 AI Features

```bash
# Create a note with AI-generated content
snip ai-create "Python Decorators" --tag "programming"

# Generate code with AI
snip ai-code "function to reverse a string" --lang "python"

# Improve search query with AI
snip ai-search "meeting notes"

# Ask questions to AI based on your notes
snip ai-ask "What did I write about Python?"
```

#### 📁 Project Management

```bash
# Create a project
snip project create "Web Application" --description "New web app project"

# Create a project with AI-generated plan
snip project ai-create "Mobile App" --description "iOS and Android app"

# List all projects
snip project list

# Show project details with tasks
snip project show 1

# Update project
snip project update 1 "Updated Name" --status "active"

# Delete a project
snip project delete 1
```

#### ✅ Tasks

```bash
# Create a task
snip task create "Implement authentication" --project 1 --priority high --due 2025-12-15

# List all tasks
snip task list

# List tasks for a specific project
snip task list --project 1

# List tasks by status
snip task list --status pending

# Show task details
snip task show 1

# Update a task
snip task update 1 "New Title" --status in_progress --priority medium

# Toggle task completion
snip task toggle 1

# Delete a task
snip task delete 1
```

#### 📋 Checklists

```bash
# Create a checklist
snip checklist create "Deployment Checklist" --project 1

# Create a checklist with AI-generated items
snip checklist ai-create "Pre-launch Checklist" --items 10 --project 1

# List all checklists
snip checklist list

# List checklists for a project
snip checklist list --project 1

# Show checklist with progress
snip checklist show 1

# Add item to checklist
snip checklist item-add 1 "Test database connection"

# Toggle checklist item completion
snip checklist item-toggle 5

# Delete checklist item
snip checklist item-delete 5

# Delete a checklist
snip checklist delete 1
```

### 🗄️ Database Analysis with AI

O Snip oferece um sistema completo e poderoso de análise de bancos de dados com integração de Inteligência Artificial. Este módulo permite realizar análises profundas, diagnósticos, tuning, e monitoramento de múltiplos tipos de bancos de dados, com interpretação inteligente dos resultados pela IA.

#### 🎯 Visão Geral

O sistema de análise de bancos de dados do Snip é uma solução abrangente que combina:

- **Análises Automatizadas**: Execução de queries e procedimentos nativos de cada banco
- **Inteligência Artificial**: Interpretação inteligente dos resultados e geração de recomendações
- **Armazenamento Persistente**: Todos os resultados são salvos para consulta posterior
- **Múltiplos Formatos**: Saída em JSON, Markdown, Texto ou HTML
- **Conexões Flexíveis**: Suporte a conexões locais, remotas, JDBC e connection strings

#### 🗃️ Bancos de Dados Suportados

- **Oracle**: Análises completas incluindo AWR, ASH, execution plans, tablespaces, e muito mais
- **SQL Server**: DMVs, locks, sessões ativas, queries em execução, planos de execução
- **MySQL**: Replicação, locks, fragmentação, slow queries, uso de índices
- **PostgreSQL**: Replicação, locks, bloqueios, fragmentação, estatísticas pg_stat
- **MongoDB**: Replicação, sharding, latência, performance, monitoramento de forks

#### 📊 Tipos de Análises Disponíveis

##### Análises Gerais (Todos os Bancos)

1. **Diagnostic** (`diagnostic`)
   - Análise completa do estado do banco de dados
   - Verificação de saúde geral, conexões, configurações
   - Identificação de problemas comuns

2. **Tuning** (`tuning`)
   - Recomendações de otimização de performance
   - Análise de queries lentas
   - Sugestões de índices e otimizações

3. **Query** (`query`)
   - Análise de consultas específicas
   - Identificação de queries problemáticas
   - Estatísticas de execução

4. **Tablespace** (`tablespace`) - Oracle
   - Uso de espaço em tablespaces
   - Crescimento e tendências
   - Alertas de espaço insuficiente

5. **Disk** (`disk`)
   - Análise de uso de disco
   - Espaço disponível e tendências
   - Recomendações de limpeza

6. **Tables** (`tables`)
   - Análise de tabelas e seus tamanhos
   - Crescimento de tabelas
   - Tabelas com maior uso

7. **Indexes** (`indexes`)
   - Análise de uso de índices
   - Índices não utilizados
   - Recomendações de criação/remoção

8. **Logs** (`logs`)
   - Análise detalhada de arquivos de log (.log e .xml)
   - Identificação de erros e warnings
   - Padrões e tendências em logs

9. **Predictive** (`predictive`)
   - Análises preditivas usando IA
   - Previsão de crescimento
   - Identificação de tendências problemáticas

10. **Error Knowledge** (`error_knowledge`)
    - Base de conhecimento de erros
    - Soluções sugeridas pela IA
    - Histórico de problemas similares

##### Análises Específicas Oracle

11. **AWR** (`awr`)
    - Geração de relatórios AWR (Automatic Workload Repository)
    - Análise de snapshots por período
    - Interpretação inteligente com IA
    - **Uso com IA**: Você pode solicitar em linguagem natural, por exemplo: "gerar AWR do período de ontem às 10h até hoje às 14h" e a IA transforma em parâmetros de snapshots

12. **ASH** (`ash`)
    - Análise de Active Session History
    - Identificação de tempos de espera
    - Planos de execução por SQL ID, Serial e SID
    - **Uso com IA**: A IA ajuda a construir queries ASH corretas e interpreta os resultados de forma inteligente

13. **Execution Plan** (`execution_plan`)
    - Análise de planos de execução
    - Identificação de problemas de performance
    - Recomendações de otimização

14. **PDBs** (`pdbs`)
    - Análise de todos os Pluggable Databases (PDBs)
    - Lista de PDBs disponíveis com status
    - Uso de espaço por PDB
    - Sessões e métricas de performance por PDB
    - **Específico para Oracle 12c+ com Multitenant**

15. **PDB** (`pdb`)
    - Análise detalhada de um PDB específico
    - Tabelas e objetos no PDB
    - Estatísticas de uso
    - **Específico para Oracle 12c+ com Multitenant**

16. **RAC Health** (`rac_health`)
    - Análise completa da saúde do cluster RAC
    - Status de todos os nós
    - Status dos serviços
    - Recursos do cluster
    - Verificação se está respondendo
    - **Específico para Oracle RAC**

17. **RAC Errors** (`rac_errors`)
    - Análise de erros no cluster RAC
    - Erros de instância
    - Erros de clusterware
    - Erros de interconnect
    - Deadlocks entre instâncias
    - **Específico para Oracle RAC**

18. **RAC Listener** (`rac_listener`)
    - Status do listener (lsnrctl)
    - Serviços registrados
    - Conexões ativas
    - Erros do listener
    - Análise de log do listener (com --log-path)
    - **Específico para Oracle RAC**

19. **RAC Latency** (`rac_latency`)
    - Latência de interconnect
    - Tempo de resposta por instância
    - Estatísticas de cache fusion
    - Bloqueios entre instâncias
    - **Específico para Oracle RAC**

##### Análises Específicas SQL Server

14. **Locks** (`locks`)
    - Análise de locks e bloqueios usando `sys.dm_tran_locks`
    - Identificação de deadlocks
    - Recomendações de resolução

15. **Active Sessions** (`active_sessions`)
    - Monitoramento de sessões ativas via DMVs
    - Análise tipo SQL Profile
    - Queries em execução e seus estados

16. **Running Queries** (`running_queries`)
    - Análise de queries em execução
    - Planos de execução ativos
    - Identificação de queries problemáticas

17. **Instance** (`instance`)
    - Análise completa da instância SQL Server
    - Informações gerais (versão, servidor, serviço)
    - Estatísticas de memória e CPU
    - Conexões ativas e status
    - **Análise ao nível da instância, não de database específico**

18. **Databases** (`databases`)
    - Lista todos os databases na instância
    - Status, recovery model, compatibility level
    - Tamanho de cada database
    - Estatísticas de I/O por database

19. **Database** (`database`)
    - Análise detalhada de um database específico
    - Tabelas e objetos no database
    - Uso de espaço (dados e log)
    - Estatísticas de performance

##### Análises Específicas MongoDB

17. **Replication** (`replication`)
    - Status de replicação
    - Lag de replicação
    - Saúde dos replicas

18. **Sharding** (`sharding`)
    - Status de sharding
    - Distribuição de dados
    - Balanceamento de shards

19. **Latency** (`latency`)
    - Análise de latência de operações
    - Identificação de gargalos
    - Recomendações de otimização

20. **Performance** (`performance`)
    - Análise geral de performance
    - Métricas de operações
    - Identificação de problemas

##### Análises Específicas PostgreSQL

21. **Postgres Replication** (`postgres_replication`)
    - Status de replicação usando `pg_stat_replication`
    - Lag de replicação (write, flush, replay)
    - Saúde dos replicas

22. **Postgres Locks** (`postgres_locks`)
    - Análise de locks usando `pg_locks` e `pg_stat_activity`
    - Identificação de bloqueios
    - Queries bloqueadas

23. **Postgres Fragmentation** (`postgres_fragmentation`)
    - Análise de fragmentação de tabelas
    - Recomendações de VACUUM
    - Uso de espaço e bloat

##### Análises Específicas MySQL

24. **MySQL Replication** (`mysql_replication`)
    - Status de replicação via `SHOW SLAVE STATUS`
    - Lag de replicação
    - Estado do slave

25. **MySQL Locks** (`mysql_locks`)
    - Análise de locks InnoDB
    - Processos ativos
    - Deadlocks recentes

26. **MySQL Fragmentation** (`mysql_fragmentation`)
    - Análise de fragmentação de tabelas
    - Espaço livre e fragmentação
    - Recomendações de otimização

##### Funcionalidades Especiais

27. **Checklist** (`checklist`)
    - Checklists diários, semanais e profundos
    - Verificações específicas por tipo de banco
    - Itens de prioridade alta, média e baixa
    - Rastreamento de status (pending, completed, skipped, failed)

28. **Backup** (`backup`)
    - Verificação de status de backups
    - Último backup realizado
    - Duração de backups
    - **Estimativa de RTO (Recovery Time Objective) usando IA**
    - Recomendações de estratégia de backup

29. **Dynamic** (`dynamic`)
    - **Análise dinâmica gerando queries com IA**
    - Você descreve o que quer analisar em linguagem natural
    - A IA gera a query SQL apropriada
    - Executa a query e interpreta os resultados
    - **Ideal para análises ad-hoc e exploração de dados**

30. **Chat** (`chat`)
    - **Chat interativo com o banco de dados usando IA**
    - Conversa natural sobre o banco
    - Geração automática de queries baseadas na conversa
    - Interpretação inteligente de resultados
    - **Use o comando `snip db-chat` para iniciar uma sessão interativa**

#### 🤖 Integração com Inteligência Artificial

A IA está integrada em múltiplos níveis do sistema de análise:

##### 1. Interpretação de Parâmetros em Linguagem Natural

**Exemplo - AWR Oracle:**
```bash
# Você pode solicitar em linguagem natural
snip db-analysis create \
  --title "AWR Período Manhã" \
  --db-type oracle \
  --analysis-type awr \
  --host localhost \
  --port 1521 \
  --database ORCL \
  --username sys \
  --password senha \
  --ai-query "gerar relatório AWR de ontem às 10h até hoje às 14h"
```

A IA interpreta a solicitação e:
- Identifica os snapshots necessários baseado nas datas/horários
- Constrói a query AWR correta
- Executa a análise
- Interpreta os resultados de forma inteligente

##### 2. Interpretação Inteligente de Resultados

Após cada análise, a IA:
- Analisa os dados coletados
- Identifica problemas e padrões
- Gera recomendações acionáveis
- Explica os resultados em linguagem clara e compreensível
- Sugere ações corretivas quando necessário

**Exemplo de saída com IA:**
```
# Resultado da Análise
[Resultados técnicos da análise...]

# Insights da IA
🤖 Análise Inteligente:

Identifiquei que o banco de dados está apresentando:
- Alto uso de CPU (85%) durante picos de carga
- Queries lentas relacionadas a falta de índices
- Fragmentação significativa em 3 tabelas principais

Recomendações:
1. Criar índice composto na tabela 'orders' nas colunas (customer_id, order_date)
2. Executar VACUUM FULL nas tabelas fragmentadas durante janela de manutenção
3. Considerar particionamento da tabela 'transactions' por data

Ações Imediatas:
- Monitorar a query com SQL ID 'abc123xyz' que está consumindo 40% do tempo de CPU
- Verificar locks na tabela 'inventory' que podem estar causando bloqueios
```

##### 3. Geração de Queries Inteligentes

Para análises como ASH (Oracle), a IA:
- Ajuda a construir queries corretas baseadas em SQL ID, Serial, SID
- Sugere filtros apropriados
- Otimiza a query para melhor performance

##### 4. Estimativa de RTO com IA

Para análises de backup, a IA:
- Analisa histórico de backups
- Considera tamanho dos dados
- Estima tempo de recuperação (RTO)
- Sugere melhorias na estratégia de backup

##### 5. Base de Conhecimento de Erros

A IA mantém contexto sobre:
- Erros comuns e suas soluções
- Padrões de problemas
- Histórico de resoluções bem-sucedidas
- Sugestões baseadas em experiências similares

#### 📝 Exemplos de Uso Passo a Passo

##### Exemplo 1: Análise Diagnóstica PostgreSQL

```bash
# Passo 1: Criar a análise
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

# Passo 2: Executar a análise
snip db-analysis run 1

# Passo 3: Ver os resultados
snip db-analysis get 1 --verbose
```

**O que acontece:**
1. Sistema conecta ao PostgreSQL
2. Executa queries de diagnóstico (versão, conexões, configurações, estatísticas)
3. IA analisa os resultados
4. Gera relatório com insights e recomendações
5. Salva tudo no banco de dados local

##### Exemplo 2: Análise de Locks SQL Server

```bash
# Criar e executar análise de locks
snip db-analysis create \
  --title "Análise de Locks SQL Server" \
  --db-type sqlserver \
  --analysis-type locks \
  --host sqlserver-prod \
  --port 1433 \
  --database AdventureWorks \
  --username sa \
  --password senha123

# Executar
snip db-analysis run 2
```

**Resultado inclui:**
- Lista de todos os locks ativos
- Sessões bloqueadas e bloqueantes
- Recomendações da IA sobre como resolver bloqueios
- Sugestões de queries para matar sessões problemáticas (se necessário)

##### Exemplo 3: Análise de Logs Oracle

```bash
# Analisar arquivo de log Oracle
snip db-analysis create \
  --title "Análise Log Alert Oracle" \
  --db-type oracle \
  --analysis-type logs \
  --log-path /u01/app/oracle/diag/rdbms/orcl/orcl/trace/alert_orcl.log

# Executar
snip db-analysis run 3
```

**A IA analisa:**
- Erros e warnings no log
- Padrões de problemas
- Ocorrências repetidas
- Sugestões de correção

##### Exemplo 4: Checklist Diário MySQL

```bash
# Criar checklist diário
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

# Executar
snip db-analysis run 4
```

**O checklist inclui:**
- Verificação de espaço em disco
- Status de replicação
- Queries lentas
- Locks ativos
- Status de backups
- E muito mais, tudo com status e recomendações

##### Exemplo 5: Análise de Backup com RTO

```bash
# Analisar backups SQL Server
snip db-analysis create \
  --title "Status Backups SQL Server" \
  --db-type sqlserver \
  --analysis-type backup \
  --host sqlserver-prod \
  --port 1433 \
  --database master \
  --username backup_admin \
  --password senha

# Executar
snip db-analysis run 5
```

**Resultado inclui:**
- Lista de backups encontrados (Full, Differential, Log)
- Data/hora do último backup
- Duração dos backups
- **Estimativa de RTO pela IA** baseada em:
  - Tamanho dos dados
  - Velocidade de restauração histórica
  - Tipo de backup disponível
- Recomendações para melhorar a estratégia de backup

##### Exemplo 6: Análise AWR Oracle com IA

```bash
# Criar análise AWR (a IA interpreta os parâmetros)
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

# Executar
snip db-analysis run 6
```

**A IA:**
1. Identifica os snapshots das últimas 24 horas
2. Foca nos períodos de maior carga
3. Gera o relatório AWR
4. Interpreta os resultados destacando:
   - Top wait events
   - Top SQL por tempo de execução
   - Problemas de performance identificados
   - Recomendações específicas

##### Exemplo 7: Análise ASH Oracle

```bash
# Analisar ASH para um SQL ID específico
snip db-analysis create \
  --title "ASH SQL ID abc123" \
  --db-type oracle \
  --analysis-type ash \
  --host oracle-prod \
  --port 1521 \
  --database ORCL \
  --username sys \
  --password senha \
  --ai-query "analisar SQL ID abc123xyz, mostrar tempos de espera e plano de execução"

# Executar
snip db-analysis run 7
```

**A IA:**
- Constrói queries ASH apropriadas
- Analisa tempos de espera
- Identifica o plano de execução
- Sugere otimizações baseadas nos dados coletados

##### Exemplo 8: Análise Dinâmica com IA

```bash
# Análise dinâmica - a IA gera a query baseada na solicitação
snip db-analysis create \
  --title "mostre as tabelas com mais de 1 milhão de linhas" \
  --db-type postgresql \
  --analysis-type dynamic \
  --host localhost \
  --port 5432 \
  --database mydb \
  --username user \
  --password pass

# Executar
snip db-analysis run 8
```

**A IA:**
1. Interpreta sua solicitação em linguagem natural
2. Gera a query SQL apropriada
3. Executa a query
4. Interpreta os resultados e fornece insights

##### Exemplo 9: Análise de PDBs Oracle

```bash
# Analisar todos os PDBs
snip db-analysis create \
  --title "Análise PDBs Oracle" \
  --db-type oracle \
  --analysis-type pdbs \
  --host oracle-prod \
  --port 1521 \
  --database CDB$ROOT \
  --username sys \
  --password senha

# Executar
snip db-analysis run 9
```

**Resultado inclui:**
- Lista de todos os PDBs com status
- Uso de espaço por PDB
- Sessões ativas por PDB
- Métricas de performance por PDB

##### Exemplo 10: Análise de Instância SQL Server

```bash
# Analisar instância completa
snip db-analysis create \
  --title "Análise Instância SQL Server" \
  --db-type sqlserver \
  --analysis-type instance \
  --host sqlserver-prod \
  --port 1433 \
  --username sa \
  --password senha

# Executar
snip db-analysis run 10
```

**Resultado inclui:**
- Informações da instância (versão, servidor)
- Estatísticas de memória e CPU
- Conexões ativas
- Status geral da instância

##### Exemplo 11: Chat Interativo

```bash
# Iniciar chat interativo
snip db-chat \
  --db-type postgresql \
  --host localhost \
  --port 5432 \
  --database mydb \
  --username user \
  --password pass

# No chat:
Você: quantas tabelas temos?
🤖 Assistente: Encontrei 25 tabelas no banco de dados...

Você: qual tabela tem mais registros?
🤖 Assistente: A tabela 'orders' tem 1.234.567 registros...

Você: mostre as 5 queries mais lentas
🤖 Assistente: [IA gera query, executa e interpreta]
```

##### Exemplo 12: Geração de Gráficos

```bash
# Gerar gráfico ASCII de uma análise
snip db-chart --analysis-id 1

# Gerar gráfico HTML interativo
snip db-chart --analysis-id 1 --type html --output chart.html

# Gerar gráfico de barras
snip db-chart --analysis-id 1 --type bar
```

**A IA:**
- Analisa os resultados da análise
- Sugere o melhor tipo de gráfico
- Extrai dados numéricos automaticamente
- Gera visualização apropriada

##### Exemplo 13: Plano de Manutenção

```bash
# Gerar plano de manutenção baseado em análise
snip db-maintenance --analysis-id 1

# Salvar plano em arquivo
snip db-maintenance --analysis-id 1 --output maintenance-plan.md
```

**O plano inclui:**
- Tarefas priorizadas com descrições detalhadas
- Passo a passo para cada tarefa
- Tempo estimado de execução
- Dependências entre tarefas
- Recomendações específicas baseadas na análise

##### Exemplo 14: Transformar Análise em Projeto

```bash
# Transformar análise em projeto
snip db-project --analysis-id 1
```

**O projeto gerado inclui:**
- Nome e descrição apropriados
- Tarefas priorizadas
- Passo a passo detalhado
- Tempo estimado e datas sugeridas
- Pode ser criado no sistema com `snip project create`

##### Exemplo 15: Transformar Incidente em Projeto

```bash
# Transformar incidente em projeto de resolução
snip db-project \
  --analysis-id 1 \
  --incident "Banco de dados apresentando lentidão durante picos de carga entre 14h e 16h"
```

**O projeto inclui:**
- Tarefas de resolução imediata
- Tarefas de prevenção
- Passo a passo para diagnóstico e correção
- Prioridade alta (incidentes)

#### 🔧 Métodos de Conexão

O sistema suporta múltiplas formas de conexão:

##### 1. Conexão por Parâmetros Individuais

```bash
snip db-analysis create \
  --title "Análise" \
  --db-type postgresql \
  --analysis-type diagnostic \
  --host localhost \
  --port 5432 \
  --database mydb \
  --username user \
  --password pass
```

**💡 Autenticação Local sem Senha:**

Para conexões locais, você pode omitir a senha se o usuário do sistema operacional for o owner do banco:

```bash
# PostgreSQL local (usuário do OS = usuário do banco)
snip db-analysis create \
  --title "Análise Local" \
  --db-type postgresql \
  --analysis-type diagnostic \
  --host localhost \
  --database mydb \
  --username postgres
  # --password não é necessário para conexões locais

# MySQL local (usuário root ou do OS)
snip db-analysis create \
  --title "Análise Local" \
  --db-type mysql \
  --analysis-type diagnostic \
  --host localhost \
  --database mydb \
  --username root
  # --password não é necessário

# Oracle local (OS Authentication)
snip db-analysis create \
  --title "Análise Local" \
  --db-type oracle \
  --analysis-type diagnostic \
  --host localhost \
  --database ORCL \
  --username sys
  # --password não é necessário (usa OS auth)

# SQL Server local (Windows Authentication)
snip db-analysis create \
  --title "Análise Local" \
  --db-type sqlserver \
  --analysis-type diagnostic \
  --host localhost \
  --database master
  # --username e --password não são necessários (usa Windows Auth)
```

**Como funciona:**
- **PostgreSQL**: Usa autenticação peer/trust se o usuário do OS for o mesmo do banco
- **MySQL**: Usa socket Unix ou autenticação sem senha para root/usuário do OS
- **Oracle**: Usa OS Authentication se o usuário estiver no grupo dba/oinstall/oracle
- **SQL Server**: Usa Windows Authentication (Integrated Security) localmente

##### 2. Conexão JDBC

```bash
snip db-analysis create \
  --title "Análise JDBC" \
  --db-type mysql \
  --analysis-type tuning \
  --jdbc-url "jdbc:mysql://localhost:3306/mydb?user=root&password=senha"
```

##### 3. Connection String

```bash
snip db-analysis create \
  --title "Análise Connection String" \
  --db-type sqlserver \
  --analysis-type locks \
  --conn-string "Server=localhost;Database=AdventureWorks;User Id=sa;Password=senha;"
```

##### 4. Conexão Remota

```bash
snip db-analysis create \
  --title "Análise Remota" \
  --db-type postgresql \
  --analysis-type diagnostic \
  --host remote-server.example.com \
  --port 5432 \
  --database mydb \
  --username user \
  --password pass \
  --remote
```

#### 📤 Formatos de Saída

Você pode escolher o formato de saída:

- **Markdown** (`markdown`) - Padrão, ideal para documentação
- **JSON** (`json`) - Para integração com outras ferramentas
- **Text** (`text`) - Formato simples de texto
- **HTML** (`html`) - Para visualização em navegador

```bash
snip db-analysis create \
  --title "Análise JSON" \
  --db-type mysql \
  --analysis-type diagnostic \
  --output json \
  --host localhost \
  --port 3306 \
  --database mydb
```

#### 📋 Comandos Disponíveis

```bash
# Criar uma nova análise
snip db-analysis create [opções]

# Criar análise com gráfico
snip db-analysis create --title "Análise" --db-type postgresql --analysis-type diagnostic --with-chart ...

# Criar análise e gerar plano de manutenção
snip db-analysis create --title "Análise" --db-type postgresql --analysis-type diagnostic --generate-plan ...

# Criar análise e transformar em projeto
snip db-analysis create --title "Análise" --db-type postgresql --analysis-type diagnostic --generate-project ...

# Listar todas as análises
snip db-analysis list

# Listar com filtros
snip db-analysis list --db-type postgresql
snip db-analysis list --analysis-type diagnostic
snip db-analysis list --limit 10

# Obter detalhes de uma análise
snip db-analysis get 1
snip db-analysis get 1 --verbose

# Executar uma análise
snip db-analysis run 1

# Deletar uma análise
snip db-analysis delete 1

# Gerar gráfico de uma análise
snip db-chart --analysis-id 1
snip db-chart --analysis-id 1 --type bar
snip db-chart --analysis-id 1 --type html --output chart.html

# Gerar plano de manutenção
snip db-maintenance --analysis-id 1
snip db-maintenance --analysis-id 1 --output maintenance-plan.md

# Transformar análise em projeto
snip db-project --analysis-id 1

# Transformar incidente em projeto
snip db-project --analysis-id 1 --incident "Banco de dados lento durante picos"

# Chat interativo com banco de dados
snip db-chat --db-type postgresql --host localhost --port 5432 --database mydb --username user --password pass

# Análises Oracle RAC
snip db-analysis create --title "Saúde RAC" --db-type oracle --analysis-type rac_health ...
snip db-analysis create --title "Erros RAC" --db-type oracle --analysis-type rac_errors ...
snip db-analysis create --title "Listener RAC" --db-type oracle --analysis-type rac_listener ...
snip db-analysis create --title "Latência RAC" --db-type oracle --analysis-type rac_latency ...
```

#### 💬 Chat Interativo com Banco de Dados

O Snip oferece um chat interativo onde você pode conversar com o banco de dados usando linguagem natural. A IA:

- **Gera queries SQL** baseadas em suas perguntas
- **Executa as queries** automaticamente
- **Interpreta os resultados** de forma clara e útil
- **Explica erros** e sugere correções
- **Mantém contexto** da conversa

**Exemplo de uso:**

```bash
# Iniciar chat
snip db-chat --db-type postgresql --host localhost --port 5432 --database mydb --username user --password pass

# No chat:
Você: quantas tabelas existem no banco?
🤖 Assistente: [IA gera query, executa e interpreta resultado]

Você: mostre as 10 tabelas com mais linhas
🤖 Assistente: [IA gera query SELECT, executa e mostra resultados interpretados]

Você: qual é a tabela que mais cresceu nos últimos 30 dias?
🤖 Assistente: [IA gera query complexa, executa e fornece análise]
```

**Comandos do chat:**
- Digite suas perguntas normalmente
- Digite `exit`, `quit` ou `sair` para encerrar o chat

#### 💾 Armazenamento e Consulta

Todas as análises são armazenadas no banco de dados SQLite local (`~/.snip/notes.db`), permitindo:

- **Histórico Completo**: Todas as análises ficam salvas para consulta posterior
- **Comparação**: Compare análises de diferentes períodos
- **Auditoria**: Mantenha registro de todas as análises realizadas
- **Relatórios**: Exporte análises para outros formatos

#### 🎯 Potencial e Benefícios

##### Para DBAs

- **Automação**: Reduz trabalho manual repetitivo
- **Inteligência**: IA identifica problemas que poderiam passar despercebidos
- **Documentação**: Todas as análises ficam documentadas automaticamente
- **Eficiência**: Análises complexas em minutos, não horas
- **Visualização**: Gráficos tornam dados mais compreensíveis
- **Planejamento**: Planos de manutenção estruturados e acionáveis
- **Gestão de Projetos**: Transforme análises e incidentes em projetos gerenciáveis

##### Para Desenvolvedores

- **Acesso Fácil**: Interface CLI simples para análises complexas
- **Aprendizado**: IA explica os resultados de forma compreensível
- **Debugging**: Identifica rapidamente problemas de performance
- **Visualização**: Gráficos ajudam a entender padrões e tendências
- **Projetos Estruturados**: Transforme problemas em projetos com tarefas claras

##### Para Equipes

- **Padronização**: Processos consistentes de análise
- **Colaboração**: Resultados compartilháveis e documentados
- **Histórico**: Rastreabilidade completa de análises
- **Visualização Compartilhada**: Gráficos HTML podem ser compartilhados
- **Gestão de Incidentes**: Transforme incidentes em projetos rastreáveis

##### Casos de Uso

1. **Monitoramento Proativo**: Execute checklists diários para identificar problemas antes que afetem produção
2. **Troubleshooting**: Análise rápida de problemas de performance com visualizações
3. **Capacitação**: Use a IA para aprender sobre bancos de dados
4. **Auditoria**: Mantenha registro de todas as análises realizadas
5. **Otimização Contínua**: Identifique oportunidades de melhoria regularmente
6. **Visualização de Dados**: Gráficos tornam análises mais didáticas e compreensíveis
7. **Planejamento de Manutenção**: Gere planos estruturados baseados em análises
8. **Gestão de Incidentes**: Transforme incidentes em projetos com tarefas e passo a passo
9. **Relatórios Executivos**: Gráficos HTML podem ser incluídos em apresentações
10. **Workflow Completo**: Análise → Gráfico → Plano → Projeto → Execução

#### 🔐 Segurança

- **Senhas**: Nunca são exibidas nos logs ou resultados
- **Conexões**: Suporte a SSL/TLS quando disponível
- **Armazenamento**: Configurações de conexão são armazenadas de forma segura
- **Permissões**: Respeita as permissões do usuário do banco de dados

#### ⚙️ Requisitos

- **Groq API Key**: Necessário para funcionalidades de IA (veja seção de configuração)
- **Drivers de Banco**: Alguns bancos podem requerer drivers específicos
- **Permissões**: Usuário do banco precisa de permissões apropriadas para as análises

#### 📊 Visualização com Gráficos

O Snip pode gerar gráficos automaticamente das análises para tornar os resultados mais visuais e didáticos:

**Tipos de Gráficos:**
- **ASCII**: Gráficos de texto para terminal
- **HTML**: Gráficos interativos usando Chart.js
- **Bar**: Gráficos de barras
- **Line**: Gráficos de linha
- **Pie**: Gráficos de pizza
- **Area**: Gráficos de área
- **Table**: Tabelas formatadas

**A IA:**
- Sugere automaticamente o melhor tipo de gráfico
- Extrai dados numéricos dos resultados
- Gera visualizações apropriadas

**Exemplo:**
```bash
# Gerar gráfico de uma análise
snip db-chart --analysis-id 1

# Gerar gráfico HTML interativo
snip db-chart --analysis-id 1 --type html --output chart.html
```

#### 🔧 Planos de Manutenção com IA

A IA pode gerar planos de manutenção detalhados baseados em análises:

**O plano inclui:**
- Tarefas priorizadas (high, medium, low)
- Passo a passo detalhado para cada tarefa
- Tempo estimado de execução
- Dependências entre tarefas
- Descrições claras e acionáveis

**Exemplo:**
```bash
# Gerar plano de manutenção
snip db-maintenance --analysis-id 1

# Salvar plano em arquivo
snip db-maintenance --analysis-id 1 --output maintenance-plan.md
```

#### 📁 Transformação em Projetos

Transforme análises e incidentes em projetos estruturados com tarefas:

**Funcionalidades:**
- Cria projeto com nome e descrição apropriados
- Gera tarefas priorizadas com passo a passo
- Estima tempo e sugere datas de vencimento
- Pode ser integrado ao sistema de projetos do Snip

**Exemplo:**
```bash
# Transformar análise em projeto
snip db-project --analysis-id 1

# Transformar incidente em projeto
snip db-project --analysis-id 1 --incident "Banco de dados lento durante picos de carga"
```

**Fluxo Completo:**
1. Execute uma análise: `snip db-analysis run 1`
2. Gere gráfico: `snip db-chart --analysis-id 1`
3. Crie plano de manutenção: `snip db-maintenance --analysis-id 1`
4. Transforme em projeto: `snip db-project --analysis-id 1`
5. Crie o projeto no sistema: `snip project create "Nome do Projeto" ...`

#### 🚀 Próximos Passos

1. Configure sua `GROQ_API_KEY` (veja seção de configuração)
2. Teste com uma análise simples: `snip db-analysis create --title "Teste" --db-type postgresql --analysis-type diagnostic ...`
3. Explore diferentes tipos de análise
4. Use checklists para monitoramento regular
5. Gere gráficos para visualizar resultados: `snip db-chart --analysis-id 1`
6. Crie planos de manutenção: `snip db-maintenance --analysis-id 1`
7. Transforme análises em projetos: `snip db-project --analysis-id 1`
8. Integre com seus processos de DevOps

## 🚀 Installation

### Package Managers

#### Scoop (Windows)
```bash
# Add the bucket
scoop bucket add snip https://github.com/matheuzgomes/Snip

# Install snip
scoop install snip

# Update snip
scoop update snip
```

#### Homebrew (macOS/Linux)
```bash
# Add the tap
brew tap matheuzgomes/homebrew-Snip

# Install snip
brew install --cask snip-notes

# Update snip
brew upgrade --cask snip-notes
```

**⚠️ macOS Security Note:**

If macOS blocks the app with "cannot be opened because the developer cannot be verified":

```bash
# Option 1: Remove quarantine attribute
xattr -d com.apple.quarantine /opt/homebrew/bin/snip

# Option 2: Allow in System Settings
# Go to: System Settings > Privacy & Security > Allow "snip"
```

### Direct Download

Pre-compiled binaries are available in the [releases](https://github.com/matheuzgomes/Snip/releases) page for:
- **Linux**: AMD64 and ARM64
- **Windows**: AMD64

### From Source

#### Prerequisites

- **Go 1.21 or later** - [Download Go](https://go.dev/dl/)
- **SQLite3 development libraries** (for CGO builds)
  - Windows: Included with Go or install via [SQLite](https://www.sqlite.org/download.html)
  - Linux: `sudo apt-get install libsqlite3-dev` (Debian/Ubuntu) or `sudo yum install sqlite-devel` (RHEL/CentOS)
  - macOS: Usually pre-installed or via Homebrew: `brew install sqlite`

#### Compilation

```bash
# Clone the repository
git clone https://github.com/hudsonrj/SnipAI.git
cd SnipAI

# Download dependencies
go mod download

# Build for your platform
go build -o snip.exe main.go

# For Windows (explicit)
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=1
go build -o snip.exe main.go

# For Linux
go build -o snip main.go

# For macOS
go build -o snip main.go

# Install to system path (Linux/macOS)
sudo mv snip /usr/local/bin/
```

#### Windows Build Notes

If you encounter issues running `snip.exe` directly, you can use:

```powershell
# Option 1: Use go run
go run main.go --help

# Option 2: Create an alias in PowerShell profile
# Add to $PROFILE:
function snip { 
    Set-Location "C:\repositorio\SnipAI\SnipAI"
    go run main.go $args
}
```

## 🗄️ Data Storage

Snip stores your notes in a SQLite database located at `~/.snip/notes.db`. The database includes:

- **Main Table**: Stores notes with metadata (ID, title, content, timestamps)
- **Tags Table**: Stores custom tags for organizing notes
- **Notes-Tags Table**: Many-to-many relationship between notes and tags
- **FTS Table**: Full-text search index for fast searching
- **Automatic Triggers**: Keeps search index synchronized with your notes

## 🔧 Configuration

### 🤖 AI Configuration (Groq API)

To use AI-powered features, you need to configure the `GROQ_API_KEY` environment variable.

#### Get Your API Key

1. Visit [Groq Console](https://console.groq.com/keys)
2. Sign up or log in
3. Generate a new API key
4. Copy the key

#### Set Environment Variable

**Windows (PowerShell):**
```powershell
# Temporary (current session only)
$env:GROQ_API_KEY="your_api_key_here"

# Permanent (add to user profile)
[Environment]::SetEnvironmentVariable("GROQ_API_KEY", "your_api_key_here", "User")
```

**Windows (CMD):**
```cmd
# Temporary
set GROQ_API_KEY=your_api_key_here

# Permanent: Control Panel > System > Advanced Settings > Environment Variables
```

**Linux/macOS:**
```bash
# Temporary
export GROQ_API_KEY="your_api_key_here"

# Permanent (add to ~/.bashrc or ~/.zshrc)
echo 'export GROQ_API_KEY="your_api_key_here"' >> ~/.bashrc
source ~/.bashrc
```

**Verify Configuration:**
```bash
# Windows PowerShell
echo $env:GROQ_API_KEY

# Linux/macOS
echo $GROQ_API_KEY
```

For detailed instructions, see [README_API_KEY.md](README_API_KEY.md).

### Editor Selection

Snip automatically detects your preferred editor with cross-platform support:

**Windows:**
- Visual Studio Code, Notepad++, Sublime Text, Atom, Micro, Nano, Vim, Notepad

**macOS:**
- Visual Studio Code, Sublime Text, Atom, Nano, Vim, Vi, Open

**Linux:**
- Nano, Vim, Vi, Micro, Visual Studio Code

**Priority Order:**
1. `$EDITOR` environment variable
2. Platform-specific editor detection
3. Smart fallback to basic editors

**Check Available Editors:**
```bash
snip editor
```

### Database Location

The database is automatically created at `~/.snip/notes.db`. The database includes:

- **Notes Table**: Your notes with metadata
- **Tags Table**: Custom tags
- **Projects Table**: Project information
- **Tasks Table**: Task details
- **Checklists Table**: Checklist definitions
- **Checklist Items Table**: Individual checklist items
- **FTS Table**: Full-text search index

You can backup your data by copying the `~/.snip/notes.db` file.

## 🛠️ Development

### Prerequisites

- Go 1.21 or later
- SQLite3 development libraries (for CGO builds)
- mingw-w64 (for Windows cross-compilation)

### Building

```bash
git clone https://github.com/matheuzgomes/Snip.git
cd Snip
go mod download
go build -o snip main.go
```

### Running Tests

```bash
# Run all tests
make test

# Run performance benchmarks
make bench

# Run tests with verbose output
go test -v ./internal/test/...
```

## 🗺️ Roadmap

### ✅ Completed Features

- ~~**🗑️ Delete Notes**: Remove notes you no longer need~~ ✅ Done!
- ~~**🏷️ Tags**: Organize notes with custom tags~~ ✅ Done!
- ~~**✏️ Patch Notes**: Update note titles and manage tags~~ ✅ Done!
- ~~**📤 Export**: Export notes to various formats (Markdown, JSON, etc.)~~ ✅ Done!
- ~~**📥 Import**: Import notes from files and directories~~ ✅ Done!
- ~~**🧪 Testing**: Comprehensive test suite with benchmarks~~ ✅ Done!
- ~~**🖼️ Markdown Preview**: Visualize rendered Markdown so you can see your notes as they'd appear formatted~~ ✅ Done!
- ~~**🤖 AI Features**: AI-powered note creation, code generation, search enhancement, and Q&A~~ ✅ Done!
- ~~**📁 Project Management**: Create and manage projects with tasks and checklists~~ ✅ Done!
- ~~**✅ Checklists**: Create checklists with AI-generated items and track progress~~ ✅ Done!

### Performance Metrics

Snip v1.1.0 delivers exceptional performance:

- **⚡ Sub-microsecond Operations**: Core operations run in 90-127 nanoseconds
- **💾 Memory Efficient**: Only 56 bytes per operation with 3 allocations
- **🧪 100% Test Coverage**: Comprehensive test suite with performance benchmarks
- **📊 Benchmarking**: Built-in performance monitoring with `make bench`

### Release Automation

We're using [GoReleaser](https://goreleaser.com/) for:

- ✅ **Automated Builds**: Cross-platform binary generation (Linux AMD64/ARM64, Windows AMD64)
- ✅ **Release Management**: Automated GitHub releases
- ✅ **Package Distribution**: Scoop, Homebrew, and Winget package managers
- ✅ **Cross-compilation**: Windows binaries built with mingw-w64
- ✅ **CGO Support**: SQLite integration with proper CGO compilation
- ✅ **CI/CD Pipeline**: Automated testing and release pipeline

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI functionality
- Uses [SQLite](https://sqlite.org/) with FTS4 for fast text search
- Inspired by modern note-taking tools and CLI utilities

**Made with ❤️ for anyone who wants to take notes**
