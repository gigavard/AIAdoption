"""
Modelli dati per il warehouse management system.

Contiene le definizioni di Prodotto e gli schemi Pydantic
per serializzazione/deserializzazione API.
"""

from pydantic import BaseModel, Field, field_validator
from typing import Optional


# ============= Validatori Comuni =============

def validate_product_name(v: str) -> str:
    """
    Validatore comune per il nome del prodotto.
    
    Regole:
    - Non può essere vuoto o solo spazi
    - Contiene solo alfanumerici, spazi, '-', '_', '.'
    
    Args:
        v: Nome del prodotto
    
    Returns:
        Nome normalizzato (trim spazi)
    
    Raises:
        ValueError: Se il nome non rispetta le regole
    """
    v = v.strip()
    if not v:
        raise ValueError("Il nome non può essere vuoto o contenere solo spazi")
    if not all(c.isalnum() or c.isspace() or c in '-_.' for c in v):
        raise ValueError("Il nome può contenere solo lettere, numeri, spazi e i caratteri '-_.'")
    return v


def validate_product_price(v: float) -> float:
    """
    Validatore comune per il prezzo del prodotto.
    
    Regole:
    - Deve essere > 0
    - Arrotonda a 2 decimali
    
    Args:
        v: Prezzo unitario
    
    Returns:
        Prezzo arrotondato a 2 decimali
    
    Raises:
        ValueError: Se il prezzo è <= 0
    """
    if v <= 0:
        raise ValueError("Il prezzo deve essere > 0")
    return round(v, 2)


# ============= Modelli Pydantic =============

class Prodotto(BaseModel):
    """
    Schema Pydantic per rappresentare un Prodotto nell'API.
    
    Usato sia per request (POST/PUT) che per response (GET).
    
    Attributes:
        id: Identificativo univoco del prodotto (auto-generato dal server)
        nome: Nome del prodotto (es. "Laptop", "Mouse")
        prezzo: Prezzo unitario in € (max 2 decimali)
        quantità: Quantità disponibile in magazzino (0-999999)
    """
    id: Optional[int] = Field(None, description="ID auto-generato (server-side)")
    nome: str = Field(..., min_length=1, max_length=100, description="Nome del prodotto")
    prezzo: float = Field(..., gt=0, description="Prezzo unitario (deve essere > 0)")
    quantità: int = Field(default=0, ge=0, le=999999, description="Quantità in magazzino (0-999999)")

    @field_validator('nome')
    @classmethod
    def validate_nome(cls, v: str) -> str:
        """Valida il nome usando il validatore comune."""
        return validate_product_name(v)
    
    @field_validator('prezzo')
    @classmethod
    def validate_prezzo(cls, v: float) -> float:
        """Valida il prezzo usando il validatore comune."""
        return validate_product_price(v)

    class Config:
        json_schema_extra = {
            "example": {
                "id": 1,
                "nome": "Laptop",
                "prezzo": 999.99,
                "quantità": 5
            }
        }


class ProdottoRequest(BaseModel):
    """
    Schema Pydantic per request POST/PUT (creazione/aggiornamento).
    
    NON include l'ID perché viene generato dal server.
    """
    nome: str = Field(..., min_length=1, max_length=100, description="Nome del prodotto")
    prezzo: float = Field(..., gt=0, description="Prezzo unitario (deve essere > 0)")
    quantità: int = Field(default=0, ge=0, le=999999, description="Quantità in magazzino (0-999999)")

    @field_validator('nome')
    @classmethod
    def validate_nome(cls, v: str) -> str:
        """Valida il nome usando il validatore comune."""
        return validate_product_name(v)
    
    @field_validator('prezzo')
    @classmethod
    def validate_prezzo(cls, v: float) -> float:
        """Valida il prezzo usando il validatore comune."""
        return validate_product_price(v)

    class Config:
        json_schema_extra = {
            "example": {
                "nome": "Mouse",
                "prezzo": 29.99,
                "quantità": 10
            }
        }
