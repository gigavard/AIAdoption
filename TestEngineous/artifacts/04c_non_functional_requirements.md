# non_functional_requirements

## quality_areas

| quality_area_id | quality_area_name | area_description | related_br_ids |
|---|---|---|---|
| QA-001 | Performance | Requisiti di reattività e tempi di risposta del sistema per le operazioni principali sotto carico operativo. | BR-001; BR-002 |
| QA-002 | Disponibilità e Affidabilità | Requisiti di continuità operativa e tolleranza ai guasti durante l'orario lavorativo. | BR-001 |
| QA-003 | Sicurezza | Requisiti di protezione dell'accesso, delle sessioni e dei dati archiviati. | BR-001; BR-002 |
| QA-004 | Scalabilità | Requisiti di capacità del sistema rispetto al numero di utenti concorrenti attesi. | BR-001; BR-002 |
| QA-005 | Manutenibilità e Operabilità | Requisiti di osservabilità, alerting e gestione operativa del sistema. | BR-001; BR-002 |
| QA-006 | Usabilità | Requisiti di facilità d'uso per utenti con diverso livello di esperienza digitale. | BR-001; BR-002 |
| QA-007 | Compliance e Auditabilità | Requisiti di conformità normativa (GDPR) e gestione dei diritti degli interessati. | BR-001; BR-002 |

## requirement_catalog

| nfr_id | quality_area_id | requirement_name | target_or_condition | applies_when | rationale |
|---|---|---|---|---|---|
| NFR-001 | QA-001 | Reattività operazioni principali | Le operazioni principali (caricamento, download, consultazione storico, visualizzazione lista file) devono essere percepite come fluide e senza latenze evidenti. [TO BE REFINED: target in ms da definire con il cliente] | Durante l'orario lavorativo, con fino a 200 utenti concorrenti. | Una risposta lenta durante l'uso quotidiano riduce l'adozione della piattaforma e vanifica i benefici attesi sulla collaborazione. |
| NFR-002 | QA-002 | Disponibilità durante l'orario lavorativo | Il sistema non deve subire interruzioni non pianificate durante le ore lavorative. [TO BE REFINED: SLA % da concordare con il cliente] | Orario lavorativo aziendale. | L'indisponibilità del sistema durante le ore operative è inaccettabile per il cliente e blocca i flussi di condivisione documentale. |
| NFR-003 | QA-003 | Scadenza sessione per inattività | Le sessioni utente devono scadere automaticamente dopo 30 minuti di inattività. | Per tutti gli utenti autenticati, in qualsiasi momento. | Riduce il rischio di accesso non autorizzato in caso di postazioni lasciate incustodite. |
| NFR-004 | QA-003 | Protezione file a livello di storage | I file archiviati devono essere accessibili esclusivamente agli utenti autorizzati del team proprietario, anche a livello di storage fisico. | In qualsiasi momento, per qualsiasi operazione di accesso ai file. | Garantisce che la segregazione dei dati per team sia applicata anche al livello infrastrutturale, in coerenza con i requisiti GDPR. |
| NFR-005 | QA-004 | Capacità utenti concorrenti | Il sistema deve supportare fino a 200 utenti concorrenti senza degradazione delle prestazioni. | Durante l'orario lavorativo. | La base utenti attuale è di circa 200 dipendenti; la piattaforma deve garantire prestazioni stabili per tutti gli utenti simultaneamente attivi. |
| NFR-006 | QA-005 | Alerting automatico in caso di malfunzionamento | Il sistema deve inviare notifiche automatiche all'amministratore di sistema e al capo progetto in caso di errori critici o indisponibilità del servizio. | Al verificarsi di anomalie o errori critici, in qualsiasi momento. | Consente una risposta operativa tempestiva ai malfunzionamenti, riducendo l'impatto sull'operatività aziendale. |
| NFR-007 | QA-006 | Semplicità d'uso per utenti non tecnici | Gli utenti devono poter completare le operazioni principali (visualizzare file, condividere, consultare storico) in non più di tre interazioni a partire dalla schermata principale, senza necessità di formazione preliminare. [TO BE REFINED: metriche di usabilità da validare in test utente] | Per tutti gli utenti, durante l'uso quotidiano. | La base utenti include dipendenti con diverso livello di esperienza digitale; un'interfaccia complessa ridurrebbe l'adozione. |
| NFR-008 | QA-007 | Gestione diritti degli interessati GDPR | Il sistema deve consentire la gestione dei diritti degli interessati (accesso, cancellazione, portabilità dei dati personali trattati). [TO BE REFINED con il DPO aziendale] | Su richiesta degli interessati, per tutta la durata del trattamento. | Obbligo normativo derivante dall'applicabilità del GDPR al trattamento dei dati personali dei dipendenti. |
| NFR-009 | QA-007 | Conservazione dati nel rispetto del GDPR | I dati personali e i log delle condivisioni devono essere conservati per il periodo minimo necessario alle finalità del trattamento e comunque nel rispetto dei termini previsti dal GDPR. [TO BE REFINED: periodo di retention da concordare con il DPO] | Per tutta la durata operativa del sistema. | Il GDPR impone il principio di limitazione della conservazione; la retention policy deve essere definita e applicata automaticamente. |

```
NON_FUNCTIONAL_REQUIREMENTS_COMPLETED
```
