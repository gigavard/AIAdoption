"""
Test suite per i modelli Pydantic e i validatori.

Testa:
- Validazione del nome del prodotto
- Validazione del prezzo
- Creazione di modelli Pydantic validi
"""

import pytest
from pydantic import ValidationError
from src.warehouse.models import Prodotto, ProdottoRequest, validate_product_name, validate_product_price


class TestValidateProductName:
    """Test per il validatore validate_product_name."""
    
    def test_nome_valido(self):
        """Un nome valido passa la validazione."""
        assert validate_product_name("Laptop") == "Laptop"
        assert validate_product_name("Mouse Pro 2024") == "Mouse Pro 2024"
        assert validate_product_name("Prodotto-1_Test.v2") == "Prodotto-1_Test.v2"
    
    def test_nome_vuoto(self):
        """Un nome vuoto solleva ValueError."""
        with pytest.raises(ValueError, match="non può essere vuoto"):
            validate_product_name("")
    
    def test_nome_solo_spazi(self):
        """Un nome con solo spazi solleva ValueError."""
        with pytest.raises(ValueError, match="non può essere vuoto"):
            validate_product_name("   ")
    
    def test_nome_caratteri_speciali_non_ammessi(self):
        """Caratteri speciali non ammessi sollevano ValueError."""
        with pytest.raises(ValueError, match="può contenere solo lettere"):
            validate_product_name("Prodotto@Special!")
    
    def test_nome_trim_spazi(self):
        """Gli spazi iniziali e finali vengono rimossi."""
        assert validate_product_name("  Laptop  ") == "Laptop"


class TestValidateProductPrice:
    """Test per il validatore validate_product_price."""
    
    def test_prezzo_valido(self):
        """Un prezzo valido viene arrotondato a 2 decimali."""
        assert validate_product_price(99.99) == 99.99
        assert validate_product_price(100) == 100.0
        assert validate_product_price(29.995) == 30.0  # Arrotondamento
    
    def test_prezzo_zero(self):
        """Un prezzo di 0 solleva ValueError."""
        with pytest.raises(ValueError, match="deve essere > 0"):
            validate_product_price(0)
    
    def test_prezzo_negativo(self):
        """Un prezzo negativo solleva ValueError."""
        with pytest.raises(ValueError, match="deve essere > 0"):
            validate_product_price(-10.50)
    
    def test_prezzo_arrotondamento(self):
        """I prezzi vengono arrotondati correttamente a 2 decimali."""
        assert validate_product_price(99.999) == 100.0
        assert validate_product_price(19.996) == 20.0


class TestProdottoModel:
    """Test per il modello Prodotto."""
    
    def test_creazione_prodotto_valido(self):
        """Un prodotto valido viene creato correttamente."""
        p = Prodotto(nome="Laptop", prezzo=999.99, quantità=5)
        assert p.nome == "Laptop"
        assert p.prezzo == 999.99
        assert p.quantità == 5
        assert p.id is None  # ID non impostato
    
    def test_prodotto_con_id(self):
        """Un prodotto con ID viene creato correttamente."""
        p = Prodotto(id=1, nome="Mouse", prezzo=29.99, quantità=10)
        assert p.id == 1
        assert p.nome == "Mouse"
    
    def test_prodotto_quantita_default(self):
        """La quantità ha valore di default 0."""
        p = Prodotto(nome="Tastiera", prezzo=79.99)
        assert p.quantità == 0
    
    def test_prodotto_nome_non_valido(self):
        """Un nome non valido causa ValidationError."""
        with pytest.raises(ValidationError):
            Prodotto(nome="", prezzo=100)
    
    def test_prodotto_prezzo_non_valido(self):
        """Un prezzo non valido causa ValidationError."""
        with pytest.raises(ValidationError):
            Prodotto(nome="Prodotto", prezzo=0)
    
    def test_prodotto_quantita_fuori_range(self):
        """Una quantità fuori dal range [0, 999999] causa ValidationError."""
        with pytest.raises(ValidationError):
            Prodotto(nome="Prodotto", prezzo=100, quantità=-1)
        
        with pytest.raises(ValidationError):
            Prodotto(nome="Prodotto", prezzo=100, quantità=1000000)


class TestProdottoRequestModel:
    """Test per il modello ProdottoRequest."""
    
    def test_creazione_request_valida(self):
        """Un ProdottoRequest valido viene creato correttamente."""
        req = ProdottoRequest(nome="Prodotto Test", prezzo=150.50, quantità=20)
        assert req.nome == "Prodotto Test"
        assert req.prezzo == 150.50
        assert req.quantità == 20
    
    def test_request_senza_quantita(self):
        """La quantità ha valore di default 0."""
        req = ProdottoRequest(nome="Prodotto", prezzo=99.99)
        assert req.quantità == 0
    
    def test_request_nome_vuoto(self):
        """Un nome vuoto causa ValidationError."""
        with pytest.raises(ValidationError):
            ProdottoRequest(nome="", prezzo=100)
    
    def test_request_prezzo_negativo(self):
        """Un prezzo negativo causa ValidationError."""
        with pytest.raises(ValidationError):
            ProdottoRequest(nome="Prodotto", prezzo=-50)
