"""
Business logic layer per il warehouse management system.

Contiene la logica di business (use cases) per operazioni su prodotti.
Fa da intermediario tra main.py (API) e db.py (persistenza).
"""

from typing import List, Optional
from datetime import datetime
from .models import Prodotto, ProdottoRequest
from . import db


# ============= Custom Exceptions =============

class ProdottoException(Exception):
    """Eccezione base per gli errori di prodotto."""
    def __init__(self, message: str, error_code: str, details: dict = None):
        self.message = message
        self.error_code = error_code
        self.details = details or {}
        self.timestamp = datetime.now().isoformat()
        super().__init__(self.message)


class ProdottoValidationError(ProdottoException):
    """Errore di validazione dei dati del prodotto."""
    def __init__(self, message: str, field: str = None):
        super().__init__(
            message, 
            "VALIDATION_ERROR",
            {"field": field}
        )


class ProdottoNotFoundError(ProdottoException):
    """Prodotto non trovato."""
    def __init__(self, prodotto_id: int):
        super().__init__(
            f"Prodotto con ID {prodotto_id} non trovato",
            "NOT_FOUND",
            {"prodotto_id": prodotto_id}
        )


class InsufficientQuantityError(ProdottoException):
    """Quantità insufficiente nel magazzino."""
    def __init__(self, prodotto_id: int, available: int, requested: int):
        super().__init__(
            f"Quantità insufficiente: disponibili {available}, richiesti {requested}",
            "INSUFFICIENT_QUANTITY",
            {
                "prodotto_id": prodotto_id,
                "available": available,
                "requested": requested
            }
        )


class DuplicateProductError(ProdottoException):
    """Prodotto duplicato (stesso nome già esiste)."""
    def __init__(self, nome: str):
        super().__init__(
            f"Prodotto con nome '{nome}' già esiste",
            "DUPLICATE_PRODUCT",
            {"nome": nome}
        )


# ============= Logica di Business =============

def crea_prodotto(prodotto_data: ProdottoRequest) -> Prodotto:
    """
    Logica di business per creare un nuovo prodotto.
    
    Validazioni:
    - Nome non vuoto
    - Prezzo > 0
    - Quantità >= 0
    - Nome univoco nel magazzino
    
    Args:
        prodotto_data: Dati del prodotto (nome, prezzo, quantità)
    
    Returns:
        Prodotto creato con ID assegnato
    
    Raises:
        ProdottoValidationError: Se i dati non sono validi
        DuplicateProductError: Se il nome esiste già
    """
    if not prodotto_data.nome or not prodotto_data.nome.strip():
        raise ProdottoValidationError("Il nome del prodotto non può essere vuoto", "nome")
    
    if prodotto_data.prezzo <= 0:
        raise ProdottoValidationError("Il prezzo deve essere > 0", "prezzo")
    
    if prodotto_data.quantità < 0:
        raise ProdottoValidationError("La quantità non può essere negativa", "quantità")
    
    # Verifica unicità del nome
    for p in db.get_all():
        if p.nome.lower() == prodotto_data.nome.lower():
            raise DuplicateProductError(prodotto_data.nome)
    
    prodotto_creato = db.create(prodotto_data)
    return prodotto_creato


def get_tutti_prodotti() -> List[Prodotto]:
    """
    Logica di business per ottenere tutti i prodotti.
    
    Returns:
        Lista di tutti i Prodotto
    """
    try:
        return db.get_all()
    except Exception as e:
        raise ProdottoException(
            f"Errore durante il recupero dei prodotti: {str(e)}",
            "INTERNAL_ERROR"
        )


def get_prodotto(prodotto_id: int) -> Prodotto:
    """
    Logica di business per ottenere un prodotto per ID.
    
    Args:
        prodotto_id: ID del prodotto da cercare
    
    Returns:
        Prodotto se trovato
    
    Raises:
        ProdottoValidationError: Se l'ID non è valido
        ProdottoNotFoundError: Se il prodotto non esiste
    """
    if prodotto_id <= 0:
        raise ProdottoValidationError("L'ID deve essere > 0", "prodotto_id")
    
    prodotto = db.get_by_id(prodotto_id)
    if prodotto is None:
        raise ProdottoNotFoundError(prodotto_id)
    
    return prodotto


def aggiorna_prodotto(prodotto_id: int, prodotto_data: ProdottoRequest) -> Prodotto:
    """
    Logica di business per aggiornare un prodotto esistente.
    
    Validazioni:
    - Prodotto deve esistere
    - Nome non vuoto
    - Prezzo > 0
    - Quantità >= 0
    
    Args:
        prodotto_id: ID del prodotto da aggiornare
        prodotto_data: Nuovi dati (nome, prezzo, quantità)
    
    Returns:
        Prodotto aggiornato
    
    Raises:
        ProdottoValidationError: Se i dati non sono validi
        ProdottoNotFoundError: Se il prodotto non esiste
    """
    if prodotto_id <= 0:
        raise ProdottoValidationError("L'ID deve essere > 0", "prodotto_id")
    
    # Verifica che il prodotto esista
    prodotto = db.get_by_id(prodotto_id)
    if prodotto is None:
        raise ProdottoNotFoundError(prodotto_id)
    
    if not prodotto_data.nome or not prodotto_data.nome.strip():
        raise ProdottoValidationError("Il nome non può essere vuoto", "nome")
    
    if prodotto_data.prezzo <= 0:
        raise ProdottoValidationError("Il prezzo deve essere > 0", "prezzo")
    
    if prodotto_data.quantità < 0:
        raise ProdottoValidationError("La quantità non può essere negativa", "quantità")
    
    prodotto_aggiornato = db.update(prodotto_id, prodotto_data)
    return prodotto_aggiornato


def preleva_quantita(prodotto_id: int, quantita: int) -> Prodotto:
    """
    Logica di business per prelevare quantità da un prodotto.
    
    Caso d'uso: vendita di un prodotto
    
    Args:
        prodotto_id: ID del prodotto
        quantita: Quantità da prelevare
    
    Returns:
        Prodotto aggiornato
    
    Raises:
        ProdottoNotFoundError: Se il prodotto non esiste
        InsufficientQuantityError: Se la quantità è insufficiente
        ProdottoValidationError: Se quantità non valida
    """
    if prodotto_id <= 0:
        raise ProdottoValidationError("L'ID deve essere > 0", "prodotto_id")
    
    if quantita <= 0:
        raise ProdottoValidationError("La quantità da prelevare deve essere > 0", "quantita")
    
    # Verifica che il prodotto esista
    prodotto = db.get_by_id(prodotto_id)
    if prodotto is None:
        raise ProdottoNotFoundError(prodotto_id)
    
    # Verifica la quantità disponibile
    if prodotto.quantità < quantita:
        raise InsufficientQuantityError(
            prodotto_id,
            available=prodotto.quantità,
            requested=quantita
        )
    
    # Aggiorna la quantità
    nuovi_dati = ProdottoRequest(
        nome=prodotto.nome,
        prezzo=prodotto.prezzo,
        quantità=prodotto.quantità - quantita
    )
    
    prodotto_aggiornato = db.update(prodotto_id, nuovi_dati)
    return prodotto_aggiornato


def elimina_prodotto(prodotto_id: int) -> bool:
    """
    Logica di business per eliminare un prodotto.
    
    Args:
        prodotto_id: ID del prodotto da eliminare
    
    Returns:
        True se eliminato
    
    Raises:
        ProdottoValidationError: Se l'ID non è valido
        ProdottoNotFoundError: Se il prodotto non esiste
    """
    if prodotto_id <= 0:
        raise ProdottoValidationError("L'ID deve essere > 0", "prodotto_id")
    
    # Verifica che il prodotto esista prima di eliminarlo
    if db.get_by_id(prodotto_id) is None:
        raise ProdottoNotFoundError(prodotto_id)
    
    return db.delete(prodotto_id)
