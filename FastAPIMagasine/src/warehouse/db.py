"""
Database layer per il warehouse management system.

Gestisce lo storage in-memory e le operazioni CRUD.
Non usa database persistente, i dati risiedono solo in memoria.

Usa un dict come storage per O(1) lookup per ID.
"""

from typing import Dict, List
from .models import Prodotto, ProdottoRequest


# Storage in-memory: dict con ID come chiave per O(1) lookup
_db: Dict[int, Prodotto] = {}

# Contatore auto-incrementale per gli ID
_id_counter: int = 0


def create(prodotto_data: ProdottoRequest) -> Prodotto:
    """
    Crea un nuovo prodotto e lo salva nel database.
    
    Args:
        prodotto_data: Dati del prodotto (nome, prezzo, quantità)
    
    Returns:
        Prodotto creato con ID assegnato
    
    Example:
        >>> p = create(ProdottoRequest(nome="Laptop", prezzo=999.99, quantità=5))
        >>> print(p.id)  # 1
    """
    global _id_counter
    _id_counter += 1
    
    prodotto = Prodotto(
        id=_id_counter,
        nome=prodotto_data.nome,
        prezzo=prodotto_data.prezzo,
        quantità=prodotto_data.quantità
    )
    _db[_id_counter] = prodotto
    return prodotto


def get_all() -> List[Prodotto]:
    """
    Restituisce tutti i prodotti nel database.
    
    Returns:
        Lista di tutti i Prodotto (può essere vuota)
    
    Example:
        >>> prodotti = get_all()
        >>> print(len(prodotti))  # numero di prodotti
    """
    return list(_db.values())


def get_by_id(prodotto_id: int) -> Prodotto | None:
    """
    Restituisce un prodotto per ID - O(1) complexity.
    
    Args:
        prodotto_id: ID del prodotto da cercare
    
    Returns:
        Prodotto se trovato, None altrimenti
    
    Example:
        >>> p = get_by_id(1)
        >>> if p:
        >>>     print(p.nome)
    """
    return _db.get(prodotto_id)


def update(prodotto_id: int, prodotto_data: ProdottoRequest) -> Prodotto | None:
    """
    Aggiorna un prodotto esistente.
    
    Args:
        prodotto_id: ID del prodotto da aggiornare
        prodotto_data: Nuovi dati (nome, prezzo, quantità)
    
    Returns:
        Prodotto aggiornato se trovato, None altrimenti
    
    Example:
        >>> p = update(1, ProdottoRequest(nome="Laptop Pro", prezzo=1499.99, quantità=3))
        >>> if p:
        >>>     print(f"Aggiornato: {p.nome}")
    """
    prodotto = get_by_id(prodotto_id)
    if prodotto is None:
        return None
    
    # Aggiorna i campi del prodotto
    prodotto.nome = prodotto_data.nome
    prodotto.prezzo = prodotto_data.prezzo
    prodotto.quantità = prodotto_data.quantità
    
    return prodotto


def delete(prodotto_id: int) -> bool:
    """
    Elimina un prodotto dal database - O(1) complexity.
    
    Args:
        prodotto_id: ID del prodotto da eliminare
    
    Returns:
        True se eliminato, False se non trovato
    
    Example:
        >>> success = delete(1)
        >>> if success:
        >>>     print("Prodotto eliminato")
    """
    if prodotto_id in _db:
        del _db[prodotto_id]
        return True
    return False


def clear_all() -> None:
    """
    Cancella tutti i prodotti e resetta il contatore ID.
    
    Utile per testing e reset completo.
    """
    global _db, _id_counter
    _db = {}
    _id_counter = 0
