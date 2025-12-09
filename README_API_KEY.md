# Configuração da Chave de API Groq

Para usar as funcionalidades de IA do Snip, você precisa configurar a variável de ambiente `GROQ_API_KEY`.

## Como obter a chave

1. Acesse: https://console.groq.com/keys
2. Crie uma conta ou faça login
3. Gere uma nova chave de API
4. Copie a chave

## Como configurar

### Windows (PowerShell)

```powershell
# Temporário (apenas nesta sessão)
$env:GROQ_API_KEY="sua_chave_aqui"

# Permanente (adicionar ao perfil)
[Environment]::SetEnvironmentVariable("GROQ_API_KEY", "sua_chave_aqui", "User")
```

### Windows (CMD)

```cmd
# Temporário
set GROQ_API_KEY=sua_chave_aqui

# Permanente (via GUI)
# Painel de Controle > Sistema > Configurações Avançadas > Variáveis de Ambiente
```

### Linux/macOS

```bash
# Temporário
export GROQ_API_KEY="sua_chave_aqui"

# Permanente (adicionar ao ~/.bashrc ou ~/.zshrc)
echo 'export GROQ_API_KEY="sua_chave_aqui"' >> ~/.bashrc
source ~/.bashrc
```

## Verificar se está configurado

```powershell
# Windows PowerShell
echo $env:GROQ_API_KEY

# Linux/macOS
echo $GROQ_API_KEY
```

## Uso

Após configurar a variável de ambiente, todas as funcionalidades de IA estarão disponíveis:

### Comandos de IA para Notas

```powershell
# Criar nota com conteúdo gerado por IA
snip.exe ai-create "Python Decorators" --tag "programming"

# Gerar código com IA
snip.exe ai-code "função para calcular fatorial" --lang "python"

# Melhorar busca com IA
snip.exe ai-search "meeting notes"

# Fazer perguntas à IA baseadas nas suas notas
snip.exe ai-ask "O que escrevi sobre Python?"
```

### Comandos de IA para Projetos

```powershell
# Criar projeto com plano detalhado gerado por IA
snip.exe project ai-create "Aplicativo Web" --description "Sistema de gestão completo"
```

### Comandos de IA para Checklists

```powershell
# Criar checklist com itens gerados por IA
snip.exe checklist ai-create "Preparação para Deploy" --items 10 --project 1

# Criar checklist para uma tarefa específica
snip.exe checklist ai-create "Testes de Integração" --items 8 --task 5
```

## Troubleshooting

### Erro: "GROQ_API_KEY environment variable is not set"

Se você receber este erro, significa que a variável de ambiente não está configurada:

1. Verifique se configurou corretamente:
   ```powershell
   echo $env:GROQ_API_KEY
   ```

2. Se estiver vazio, configure novamente:
   ```powershell
   [Environment]::SetEnvironmentVariable("GROQ_API_KEY", "sua_chave_aqui", "User")
   ```

3. **Reinicie o PowerShell** para aplicar as mudanças permanentes

### Erro: "AI client not available"

Este erro indica que o cliente de IA não pôde ser inicializado. Verifique:
- A chave de API está configurada corretamente
- Você tem conexão com a internet
- A chave de API é válida e não expirou

## Segurança

⚠️ **Importante:**
- Nunca compartilhe sua chave de API
- Não commite a chave de API no código
- Use variáveis de ambiente para armazenar a chave
- Revogue e gere uma nova chave se suspeitar que foi comprometida

