mkdir -p docs
cat > docs/PLAN.md << 'EOF'
# Plan: CRUD Magazzino con FastAPI

## TL;DR
Creare un'API REST con **FastAPI** per gestire prodotti di un magazzino con operazioni CRUD (Create, Read, Update, Delete). I dati risiedono in memoria (lista Python) e l'app serve su `http://localhost:8000`. Struttura minimale con `src/` e `docs/` per documentazione.

---

## Steps

### **Fase 1: Setup iniziale del progetto**
1. Creare cartelle: `src/`, `docs/`, `tests/` (struttura standard Python)
2. Creare e attivare **venv**: `python -m venv venv` e `source venv/bin/activate`
3. Creare `requirements.txt` con dipendenze: `fastapi`, `uvicorn[standard]`
4. Installare dipendenze: `pip install -r requirements.txt`

### **Fase 2: Implementare la logica CRUD**
5. Creare `src/warehouse/models.py`:
   - Classe/dataclass `Prodotto` con campi: `id`, `nome`, `prezzo`, `quantità`
   - Pydantic BaseModel per serializzazione API (`ProdottoSchema`)

6. Creare `src/warehouse/db.py` (o inline in main.py):
   - Storage in-memory: lista `prodotti = []`
   - Funzioni CRUD: `create()`, `get_all()`, `get_by_id()`, `update()`, `delete()`
   - Logica di gestione ID auto-incremento

7. Creare `src/warehouse/main.py`:
   - Inizializzare app FastAPI
   - Definire routes (Fase 3)

### **Fase 3: Definire le API routes** (*dipende da Step 5-6*)
8. **GET `/products`**: Restituisce lista di tutti i prodotti
9. **GET `/products/{id}`**: Restituisce singolo prodotto by ID, errore 404 se non esiste
10. **POST `/products`**: Crea nuovo prodotto, restituisce il creato con ID assegnato
11. **PUT `/products/{id}`**: Aggiorna prodotto esistente, errore 404 se non esiste
12. **DELETE `/products/{id}`**: Elimina prodotto, errore 404 se non esiste
13. Aggiungere response models e status codes corretti (201 per POST, 204 o 200 per DELETE)

### **Fase 4: Documentazione**
14. Creare `README.md` nel root con:
    - Descrizione breve
    - Instructions: setup venv + pip install + come avviare l'app (`uvicorn src.warehouse.main:app --reload`)
    - Endpoints disponibili (lista con metodo HTTP, path, breve descrizione)
15. Aggiungere **docstrings** (formato Google o NumPy) a tutte le funzioni e routes FastAPI
16. Creare `docs/PLAN.md` (questo file, copia da /memories/session/plan.md)

### **Fase 5: Test manuale e validazione**
17. Avviare l'app: `uvicorn src.warehouse.main:app --reload`
18. Testare ogni endpoint (via curl, Postman, o FastAPI docs auto-generati su `/docs`)
19. Verificare: creazione, lettura, aggiornamento, cancellazione, errori 404

---

## Relevant files
- `requirements.txt` — FastAPI, uvicorn
- `src/warehouse/models.py` — `Prodotto` class/dataclass + `ProdottoSchema` Pydantic model
- `src/warehouse/db.py` — Storage in-memory + funzioni CRUD (o inline in main.py)
- `src/warehouse/main.py` — FastAPI app + routes (GET, POST, PUT, DELETE)
- `src/warehouse/__init__.py` — File vuoto per namespace package
- `README.md` — Setup instructions + endpoints overview
- `docs/PLAN.md` — Copia di questo piano

---

## Verification
1. **Setup**: Venv creato, dipendenze installate, `python -m fastapi --version` funziona
2. **Importi**: `python -c "from src.warehouse.main import app"` non restituisce errori
3. **API running**: `uvicorn src.warehouse.main:app --reload` avvia senza errori
4. **Endpoints**:
   - POST su `/products` con body `{"nome": "Prodotto 1", "prezzo": 10.0, "quantità": 5}` → ritorna 201 con ID
   - GET `/products` → ritorna lista (almeno quello creato)
   - GET `/products/1` → ritorna il prodotto creato
   - PUT `/products/1` con dati modificati → ritorna 200 con prodotto aggiornato
   - DELETE `/products/1` → ritorna 204 o 200
   - GET `/products/999` → ritorna 404
5. **Docs**: FastAPI Swagger UI su `http://localhost:8000/docs` visibile e funzionante
6. **Code quality**: Docstrings presenti, no errori di sintassi

---

## Decisions
- **Framework**: FastAPI (moderno, performante, auto-genera API docs)
- **Persistenza**: In-memory (lista Python) per semplicità educativa
- **Database**: Nessuno (no SQLAlchemy/SQLite per ora)
- **Testing**: Nessuno unit test (aggiungibile dopo come upgrade)
- **Struttura**: Minimale (`src/`, `docs/`, `tests/` dir creata ma vuota per ora)
- **Python version**: 3.8+ (FastAPI requirement)
- **Naming**: 
  - Modulo: `src.warehouse.main` (warehouse come dominio)
  - Endpoints in singolare: `/products` e `/products/{id}` (RESTful standard)
  - Modelli dati: `Prodotto` (italiano per coerenza con contesto)

---

## Further Considerations
1. **Pydantic versioning**: FastAPI funziona con Pydantic v1 e v2. Raccomandazione: usare Pydantic v2 (incluso in `fastapi>=0.100.0`). Se noti incompatibilità, specifica versione in requirements.txt.
2. **Port di default**: FastAPI di default usa porta 8000. Se occupata, il comando è `uvicorn src.warehouse.main:app --port 8001 --reload`.
3. **Validazione**: FastAPI valida automaticamente i dati con Pydantic. Aggiungere custom validators se necessario (es. prezzo > 0, quantità >= 0).
4. **Upgrade futuro**: Per persistenza reale, aggiungere SQLAlchemy + SQLite (10-15 min di refactoring).
EOF