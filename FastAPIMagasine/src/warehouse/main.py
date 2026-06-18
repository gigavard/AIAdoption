"""
Applicazione FastAPI per il warehouse management system.

Gestisce le API REST per operazioni CRUD su prodotti.
Avviare con: uvicorn src.warehouse.main:app --reload
"""

from fastapi import FastAPI, HTTPException, Request
from fastapi.responses import JSONResponse
from pydantic import ValidationError

from .models import Prodotto, ProdottoRequest
from . import service
from .service import (
    ProdottoException,
    ProdottoValidationError,
    ProdottoNotFoundError,
    InsufficientQuantityError,
    DuplicateProductError
)


# Inizializza l'app FastAPI
app = FastAPI(
    title="Warehouse API",
    description="API per gestione prodotti di magazzino",
    version="0.1.0"
)


# ============= Exception Handlers con Logica di Business =============

@app.exception_handler(ProdottoValidationError)
async def validation_error_handler(request: Request, exc: ProdottoValidationError):
    """
    Handler per errori di validazione.
    
    Logica di business:
    - Status 422 (Unprocessable Entity)
    - Include il campo che ha causato l'errore
    - Include timestamp per tracciamento
    """
    return JSONResponse(
        status_code=422,
        content={
            "error_code": exc.error_code,
            "message": exc.message,
            "field": exc.details.get("field"),
            "timestamp": exc.timestamp,
            "path": str(request.url.path)
        }
    )


@app.exception_handler(ProdottoNotFoundError)
async def not_found_error_handler(request: Request, exc: ProdottoNotFoundError):
    """
    Handler per prodotto non trovato.
    
    Logica di business:
    - Status 404 (Not Found)
    - Include l'ID cercato per aiutare il debug
    """
    return JSONResponse(
        status_code=404,
        content={
            "error_code": exc.error_code,
            "message": exc.message,
            "prodotto_id": exc.details.get("prodotto_id"),
            "timestamp": exc.timestamp,
            "path": str(request.url.path)
        }
    )


@app.exception_handler(InsufficientQuantityError)
async def insufficient_quantity_handler(request: Request, exc: InsufficientQuantityError):
    """
    Handler per quantità insufficiente.
    
    Logica di business:
    - Status 409 (Conflict) - risorsa esiste ma non può essere usata
    - Include disponibilità e richiesto per decidere se ordinare
    - Suggerisce l'azione (ordina più stock, ajusta quantità, ecc)
    """
    return JSONResponse(
        status_code=409,
        content={
            "error_code": exc.error_code,
            "message": exc.message,
            "prodotto_id": exc.details.get("prodotto_id"),
            "available": exc.details.get("available"),
            "requested": exc.details.get("requested"),
            "shortfall": exc.details.get("requested") - exc.details.get("available"),
            "suggestion": "Aumentare lo stock o ridurre la quantità richiesta",
            "timestamp": exc.timestamp,
            "path": str(request.url.path)
        }
    )


@app.exception_handler(DuplicateProductError)
async def duplicate_product_handler(request: Request, exc: DuplicateProductError):
    """
    Handler per prodotto duplicato.
    
    Logica di business:
    - Status 409 (Conflict) - nome già esiste
    - Include il nome duplicato
    - Suggerisce di usare PUT per aggiornare se è lo stesso prodotto
    """
    return JSONResponse(
        status_code=409,
        content={
            "error_code": exc.error_code,
            "message": exc.message,
            "nome": exc.details.get("nome"),
            "suggestion": "Usa PUT /products/{id} per aggiornare un prodotto esistente",
            "timestamp": exc.timestamp,
            "path": str(request.url.path)
        }
    )


@app.exception_handler(ProdottoException)
async def generic_error_handler(request: Request, exc: ProdottoException):
    """
    Handler generico per eccezioni base.
    
    Logica di business:
    - Status 500 per errori imprevisti
    - Log dettagliato per debugging
    """
    return JSONResponse(
        status_code=500,
        content={
            "error_code": exc.error_code,
            "message": exc.message,
            "timestamp": exc.timestamp,
            "path": str(request.url.path)
        }
    )


@app.exception_handler(ValidationError)
async def pydantic_validation_handler(request: Request, exc: ValidationError):
    """
    Handler per errori di validazione Pydantic.
    
    Logica di business:
    - Status 422 (Unprocessable Entity)
    - Dettagli di ogni campo che ha fallito la validazione
    """
    errors = exc.errors()
    formatted_errors = [
        {
            "field": str(error["loc"][-1]),
            "message": error["msg"],
            "type": error["type"]
        }
        for error in errors
    ]
    
    return JSONResponse(
        status_code=422,
        content={
            "error_code": "PYDANTIC_VALIDATION_ERROR",
            "message": "Errore di validazione dei dati di input",
            "errors": formatted_errors,
            "path": str(request.url.path)
        }
    )


@app.get("/", tags=["Health"])
def root():
    """
    Endpoint di health check.
    
    Returns:
        Messaggio di benvenuto
    """
    return {
        "message": "Warehouse API v0.1.0",
        "docs": "/docs"
    }


@app.post("/products", response_model=Prodotto, status_code=201, tags=["Products"])
def crea_prodotto(prodotto_data: ProdottoRequest):
    """
    Crea un nuovo prodotto nel magazzino.
    
    Args:
        prodotto_data: Dati del prodotto (nome, prezzo, quantità)
    
    Returns:
        Prodotto creato con ID assegnato (status code 201)
    
    Raises:
        422: Se i dati non sono validi (validazione Pydantic)
        409: Se un prodotto con lo stesso nome esiste già
    
    Example:
        POST /products
        {
            "nome": "Laptop",
            "prezzo": 999.99,
            "quantità": 5
        }
    """
    prodotto_creato = service.crea_prodotto(prodotto_data)
    return prodotto_creato


@app.put("/products/{prodotto_id}", response_model=Prodotto, tags=["Products"])
def aggiorna_prodotto(prodotto_id: int, prodotto_data: ProdottoRequest):
    """
    Aggiorna un prodotto esistente nel magazzino.
    
    Args:
        prodotto_id: ID del prodotto da aggiornare (da URL)
        prodotto_data: Nuovi dati del prodotto (nome, prezzo, quantità)
    
    Returns:
        Prodotto aggiornato
    
    Raises:
        422: Se l'ID o i dati non sono validi
        404: Se il prodotto non esiste
    
    Example:
        PUT /products/1
        {
            "nome": "Laptop Pro",
            "prezzo": 1499.99,
            "quantità": 3
        }
    """
    prodotto_aggiornato = service.aggiorna_prodotto(prodotto_id, prodotto_data)
    return prodotto_aggiornato


@app.get("/products", response_model=list[Prodotto], tags=["Products"])
def lista_prodotti():
    """
    Ottiene la lista di tutti i prodotti nel magazzino.
    
    Returns:
        Lista di tutti i Prodotto (può essere vuota)
    
    Example:
        GET /products
    """
    prodotti = service.get_tutti_prodotti()
    return prodotti


@app.delete("/products/{prodotto_id}", status_code=204, tags=["Products"])
def elimina_prodotto(prodotto_id: int):
    """
    Elimina un prodotto dal magazzino.
    
    Args:
        prodotto_id: ID del prodotto da eliminare (da URL)
    
    Returns:
        Nessun contenuto (status code 204)
    
    Raises:
        422: Se l'ID non è valido
        404: Se il prodotto non esiste
    
    Example:
        DELETE /products/1
    """
    service.elimina_prodotto(prodotto_id)
    return None


@app.post("/products/{prodotto_id}/preleva", response_model=Prodotto, tags=["Products"])
def preleva_da_magazzino(prodotto_id: int, quantita: int):
    """
    Preleva quantità da un prodotto (simulazione di una vendita/prelievo).
    
    Logica di business:
    - Verifica che il prodotto esista
    - Verifica che la quantità sia disponibile
    - Decrementa la quantità se disponibile
    
    Args:
        prodotto_id: ID del prodotto
        quantita: Quantità da prelevare (query parameter)
    
    Returns:
        Prodotto aggiornato con la nuova quantità
    
    Raises:
        422: Se i parametri non sono validi
        404: Se il prodotto non esiste
        409: Se la quantità è insufficiente (include disponibilità e richiesto)
    
    Example:
        POST /products/1/preleva?quantita=3
        
        Response (200):
        {
            "id": 1,
            "nome": "Laptop",
            "prezzo": 999.99,
            "quantità": 2
        }
        
        Se quantità insufficiente (409):
        {
            "error_code": "INSUFFICIENT_QUANTITY",
            "message": "Quantità insufficiente: disponibili 2, richiesti 3",
            "prodotto_id": 1,
            "available": 2,
            "requested": 3,
            "shortfall": 1,
            "suggestion": "Aumentare lo stock o ridurre la quantità richiesta",
            "timestamp": "2026-06-18T...",
            "path": "/products/1/preleva"
        }
    """
    prodotto_aggiornato = service.preleva_quantita(prodotto_id, quantita)
    return prodotto_aggiornato
