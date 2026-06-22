---
name: conventional-commit
description: Guida alla scrittura di messaggi commit secondo lo standard Conventional Commits.
---

# Conventional Commit

Usa questo skill per creare messaggi di commit chiari e consistenti secondo lo standard Conventional Commits 1.0.0.

## Formato

```
<type>(<scope>): <description>

[optional body]

[optional footer(s)]
```

## Tipi comuni

- `feat`: nuova funzionalita
- `fix`: correzione bug
- `docs`: sola documentazione
- `style`: formattazione, nessun cambiamento logico
- `refactor`: refactor senza fix/feat
- `test`: aggiunta o modifica test
- `chore`: task di manutenzione
- `perf`: miglioramento performance
- `build`: modifiche sistema di build o dipendenze
- `ci`: modifiche pipeline CI

## Regole rapide

- Usa descrizioni brevi all'imperativo, in minuscolo.
- Evita il punto finale nella subject line.
- Mantieni la subject line idealmente entro 72 caratteri.
- Usa `!` dopo type/scope per breaking change (es. `feat(api)!: ...`).
- In alternativa, indica breaking change nel footer con:
  - `BREAKING CHANGE: <descrizione>`

## Esempi

- `feat(auth): add refresh token rotation`
- `fix(warehouse): handle null supplier id`
- `docs(readme): update setup instructions`
- `refactor(db): split connection factory`
- `feat(api)!: remove legacy endpoint`

Con footer:

```
feat(order): add order cancellation endpoint

Allow users to cancel pending orders from dashboard.

Closes #42
```
