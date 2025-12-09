# Como usar o Snip

## Problema conhecido com snip.exe

Se o `snip.exe` não executar diretamente (erro "not a valid application"), use uma das alternativas abaixo:

## Solução 1: Usar go run (Recomendado)

```powershell
go run main.go --help
go run main.go project create "Meu Projeto"
go run main.go checklist ai-create "Preparação" --items 5
```

## Solução 2: Usar o script PowerShell

```powershell
.\snip.ps1 --help
.\snip.ps1 project list
.\snip.ps1 checklist show 1
```

## Solução 3: Criar alias permanente

Adicione ao seu perfil do PowerShell (`notepad $PROFILE`):

```powershell
function snip { 
    Set-Location "C:\repositorio\SnipAI\SnipAI"
    go run main.go $args
}
```

Depois use simplesmente:
```powershell
snip --help
snip project create "Meu Projeto"
```

## Compilar novamente

Se precisar recompilar:

```powershell
cd C:\repositorio\SnipAI\SnipAI
go build -o snip.exe main.go
```

## Todas as funcionalidades disponíveis

- ✅ Notas (create, list, show, find, update, delete)
- ✅ IA (ai-create, ai-code, ai-search, ai-ask)
- ✅ Projetos (create, list, show, update, delete, ai-create)
- ✅ Tarefas (create, list, show, update, delete, toggle)
- ✅ Checklists (create, ai-create, list, show, item-add, item-toggle)

