# Analysis Artifacts — TestEngineous

## Artefatti generati

| File | Descrizione |
|------|-------------|
| [01_business_requirements.md](01_business_requirements.md) | Business requirements e lista dei ruoli estratti dal PRD. Fonte primaria per la tracciabilità. |
| [02_functional_requirements.md](02_functional_requirements.md) | Aree funzionali e catalogo dei requisiti funzionali atomici, tracciati ai BR. |
| [03_features.md](03_features.md) | Catalogo delle feature: capacità funzionali discrete, tracciate a BR, FR e ruoli. |
| [04_use_cases.md](04_use_cases.md) | Catalogo degli use case: interazioni goal-oriented tra attori e sistema, tracciate a Feature e FR. |
| [04b_scenarios.md](04b_scenarios.md) | Catalogo degli scenari (Main, Alternate, Exception) con passi dettagliati per ogni use case. |
| [04c_non_functional_requirements.md](04c_non_functional_requirements.md) | Requisiti non funzionali per area di qualità, con target misurabili e tracciabilità ai BR. |

## Catena di tracciabilità

```
BR (Business Requirements)
 └── FR (Functional Requirements)  ← tracciati a BR
      └── Feature                  ← tracciata a BR + FR + Ruoli
           └── Use Case            ← tracciato a Feature + FR
                └── Scenario       ← tracciato a Use Case + Feature
NFR (Non-Functional Requirements)  ← tracciati a BR (qualità trasversale)
```

**Come leggere gli artefatti:**
- Parti sempre da `01_business_requirements.md` per capire il *perché* del progetto.
- Usa `02_functional_requirements.md` per il *cosa* il sistema deve fare.
- Usa `03_features.md` per pianificare e stimare le unità di lavoro.
- Usa `04_use_cases.md` e `04b_scenarios.md` per progettare e testare i flussi.
- Usa `04c_non_functional_requirements.md` come vincoli di qualità trasversali per architettura e testing.
