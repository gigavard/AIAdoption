"""
Test suite per la logica di business (service layer).

Testa:
- Creazione di prodotti con validazione
- Recupero di prodotti
- Aggiornamento di prodotti
- Gestione errori (duplicati, non trovati, etc.)
"""

import pytest
from unittest.mock import Mock, patch, MagicMock
from src.warehouse.models import Prodotto, ProdottoRequest
from src.warehouse.service import (
    crea_prodotto,
    get_tutti_prodotti,
    get_prodotto,
    ProdottoException,
    ProdottoValidationError,
    ProdottoNotFoundError,
    DuplicateProductError,
    InsufficientQuantityError,
)


class TestCreaProduct:
    """Test per la funzione crea_prodotto."""
    
    @patch('src.warehouse.service.db')
    def test_crea_prodotto_valido(self, mock_db):
        """Un prodotto valido viene creato correttamente."""
        # Setup
        mock_db.get_all.return_value = []
        prodotto_creato = Prodotto(id=1, nome="Laptop", prezzo=999.99, quantità=5)
        mock_db.create.return_value = prodotto_creato
        
        # Test
        req = ProdottoRequest(nome="Laptop", prezzo=999.99, quantità=5)
        result = crea_prodotto(req)
        
        # Assert
        assert result.id == 1
        assert result.nome == "Laptop"
        assert result.prezzo == 999.99
        assert result.quantità == 5
        mock_db.create.assert_called_once()
    
    @patch('src.warehouse.service.db')
    def test_crea_prodotto_nome_vuoto(self, mock_db):
        """Un nome vuoto viene validato da Pydantic e solleva ValidationError."""
        from pydantic import ValidationError
        
        with pytest.raises(ValidationError):
            req = ProdottoRequest(nome="", prezzo=100, quantità=5)
    
    @patch('src.warehouse.service.db')
    def test_crea_prodotto_prezzo_negativo(self, mock_db):
        """Un prezzo negativo viene validato da Pydantic e solleva ValidationError."""
        from pydantic import ValidationError
        
        with pytest.raises(ValidationError):
            req = ProdottoRequest(nome="Prodotto", prezzo=-50, quantità=5)
    
    @patch('src.warehouse.service.db')
    def test_crea_prodotto_quantita_negativa(self, mock_db):
        """Una quantità negativa viene validata da Pydantic e solleva ValidationError."""
        from pydantic import ValidationError
        
        with pytest.raises(ValidationError):
            req = ProdottoRequest(nome="Prodotto", prezzo=100, quantità=-5)
    
    @patch('src.warehouse.service.db')
    def test_crea_prodotto_duplicato(self, mock_db):
        """Un nome duplicato solleva DuplicateProductError."""
        # Setup: prodotto con lo stesso nome esiste già
        prodotto_esistente = Prodotto(id=1, nome="Laptop", prezzo=999.99, quantità=5)
        mock_db.get_all.return_value = [prodotto_esistente]
        
        # Test: tentativo di creare un nuovo prodotto con lo stesso nome
        req = ProdottoRequest(nome="Laptop", prezzo=899.99, quantità=3)
        
        with pytest.raises(DuplicateProductError):
            crea_prodotto(req)
    
    @patch('src.warehouse.service.db')
    def test_crea_prodotto_nome_case_insensitive(self, mock_db):
        """La verifica di duplicati è case-insensitive."""
        # Setup
        prodotto_esistente = Prodotto(id=1, nome="Laptop", prezzo=999.99, quantità=5)
        mock_db.get_all.return_value = [prodotto_esistente]
        
        # Test: "LAPTOP" è considerato duplicato di "Laptop"
        req = ProdottoRequest(nome="LAPTOP", prezzo=899.99, quantità=3)
        
        with pytest.raises(DuplicateProductError):
            crea_prodotto(req)


class TestGetTuttiProdotti:
    """Test per la funzione get_tutti_prodotti."""
    
    @patch('src.warehouse.service.db')
    def test_get_tutti_prodotti_lista_non_vuota(self, mock_db):
        """Restituisce la lista di tutti i prodotti."""
        # Setup
        prodotti = [
            Prodotto(id=1, nome="Laptop", prezzo=999.99, quantità=5),
            Prodotto(id=2, nome="Mouse", prezzo=29.99, quantità=10),
        ]
        mock_db.get_all.return_value = prodotti
        
        # Test
        result = get_tutti_prodotti()
        
        # Assert
        assert len(result) == 2
        assert result[0].nome == "Laptop"
        assert result[1].nome == "Mouse"
    
    @patch('src.warehouse.service.db')
    def test_get_tutti_prodotti_lista_vuota(self, mock_db):
        """Restituisce una lista vuota se non ci sono prodotti."""
        mock_db.get_all.return_value = []
        
        result = get_tutti_prodotti()
        
        assert result == []
    
    @patch('src.warehouse.service.db')
    def test_get_tutti_prodotti_errore_db(self, mock_db):
        """Un errore del database solleva ProdottoException."""
        mock_db.get_all.side_effect = Exception("DB connection error")
        
        with pytest.raises(ProdottoException, match="Errore durante il recupero"):
            get_tutti_prodotti()


class TestGetProdotto:
    """Test per la funzione get_prodotto."""
    
    @patch('src.warehouse.service.db')
    def test_get_prodotto_trovato(self, mock_db):
        """Un prodotto trovato viene restituito correttamente."""
        # Setup
        prodotto = Prodotto(id=1, nome="Laptop", prezzo=999.99, quantità=5)
        mock_db.get_by_id.return_value = prodotto
        
        # Test
        result = get_prodotto(1)
        
        # Assert
        assert result.id == 1
        assert result.nome == "Laptop"
        mock_db.get_by_id.assert_called_once_with(1)
    
    @patch('src.warehouse.service.db')
    def test_get_prodotto_non_trovato(self, mock_db):
        """Un prodotto non trovato solleva ProdottoNotFoundError."""
        mock_db.get_by_id.return_value = None
        
        with pytest.raises(ProdottoNotFoundError):
            get_prodotto(999)
    
    @patch('src.warehouse.service.db')
    def test_get_prodotto_id_non_valido(self, mock_db):
        """Un ID non valido solleva ProdottoValidationError."""
        with pytest.raises(ProdottoValidationError):
            get_prodotto(0)
        
        with pytest.raises(ProdottoValidationError):
            get_prodotto(-1)


class TestExceptionHandling:
    """Test per la gestione delle eccezioni."""
    
    def test_prodotto_exception_base(self):
        """ProdottoException contiene i campi corretti."""
        exc = ProdottoException("Test message", "TEST_CODE", {"field": "test"})
        
        assert exc.message == "Test message"
        assert exc.error_code == "TEST_CODE"
        assert exc.details == {"field": "test"}
        assert exc.timestamp is not None
    
    def test_prodotto_validation_error(self):
        """ProdottoValidationError contiene il campo."""
        exc = ProdottoValidationError("Invalid field", "nome")
        
        assert exc.error_code == "VALIDATION_ERROR"
        assert exc.details["field"] == "nome"
    
    def test_prodotto_not_found_error(self):
        """ProdottoNotFoundError contiene l'ID."""
        exc = ProdottoNotFoundError(123)
        
        assert exc.error_code == "NOT_FOUND"
        assert exc.details["prodotto_id"] == 123
    
    def test_insufficient_quantity_error(self):
        """InsufficientQuantityError contiene i dettagli della quantità."""
        exc = InsufficientQuantityError(1, 5, 10)
        
        assert exc.error_code == "INSUFFICIENT_QUANTITY"
        assert exc.details["prodotto_id"] == 1
        assert exc.details["available"] == 5
        assert exc.details["requested"] == 10
    
    def test_duplicate_product_error(self):
        """DuplicateProductError contiene il nome del prodotto."""
        exc = DuplicateProductError("Laptop")
        
        assert exc.error_code == "DUPLICATE_PRODUCT"
        assert exc.details["nome"] == "Laptop"
