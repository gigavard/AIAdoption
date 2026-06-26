# Product Requirements Document (PRD)
**Progetto:** TestEngineous
**Fase:** CONCEPT_DRAFT
**Data:** 2026-06-26

---

# 1. Panoramica del Progetto

## Contesto Organizzativo
L'iniziativa si inserisce in un contesto aziendale strutturato in team di lavoro, dove la collaborazione quotidiana richiede la condivisione frequente di documenti e file tra colleghi. L'organizzazione opera in un ambiente in cui la gestione dei contenuti digitali è parte integrante dei processi operativi interni. I destinatari principali della soluzione sono i dipendenti aziendali che collaborano su attività condivise all'interno dei rispettivi team.

## Sfide Attuali
Attualmente, lo scambio di file avviene prevalentemente tramite posta elettronica, con conseguente mancanza di un archivio condiviso e tracciabile. Questo approccio genera inefficienze significative: i destinatari spesso non si accorgono dei messaggi ricevuti, con il rischio di ritardi o mancate prese in carico. L'assenza di uno storico centralizzato rende difficile ricostruire l'evoluzione dei documenti e individuare la versione più aggiornata.

## Obiettivo Strategico
L'organizzazione intende superare i limiti del modello di condivisione informale via email, introducendo una piattaforma dedicata che garantisca ordine, tracciabilità e controllo nella gestione dei file condivisi. Il driver di cambiamento è la necessità di migliorare l'efficienza collaborativa interna, riducendo il rischio di comunicazioni perse e di lavoro basato su versioni obsolete dei documenti.

## Risultati Attesi
La nuova piattaforma dovrà consentire ai dipendenti di ricevere notifiche tempestive ogni volta che un file viene loro condiviso, eliminando il rischio di comunicazioni non rilevate. Ogni utente potrà gestire i propri file in autonomia — modificarli o eliminarli — garantendo la piena ownership dei contenuti. La disponibilità di uno storico consultabile migliorerà la trasparenza dei flussi di lavoro e la rintracciabilità delle informazioni nel tempo.

---

# 2. Business Requirements

| ID | Nome | Descrizione |
|----|------|-------------|
| BR-001 | Tracciabilità delle Comunicazioni | L'azienda deve eliminare il rischio di comunicazioni di file non ricevute o ignorate. Ogni condivisione deve essere notificata in modo affidabile al destinatario e deve essere consultabile tramite uno storico centralizzato, garantendo piena visibilità sui flussi di scambio documentale tra i team. |
| BR-002 | Responsabilità Individuale sui Contenuti | Ogni dipendente deve esercitare piena responsabilità sui file che ha condiviso, con la facoltà di aggiornarli o rimuoverli autonomamente. Questo risponde alla necessità aziendale di garantire la presa in carico individuale dei contenuti e di evitare la circolazione di file obsoleti o non più pertinenti. |

---

# 3. Functional Requirements

## 3.1 Gestione Utenti e Ruoli
> *Traces to: BR-002*

- Il sistema deve associare ogni utente a un team e a un ruolo (team leader o membro).
- Il sistema deve applicare le autorizzazioni in base al ruolo: i team leader possono caricare e modificare file; i membri possono solo leggere e condividere file.
- Il sistema deve impedire ai membri di eseguire operazioni riservate ai team leader (caricamento di nuovi file, modifica).
- Il sistema deve impedire la condivisione di file tra utenti appartenenti a team differenti.

## 3.2 Gestione File
> *Traces to: BR-001, BR-002*

- Il sistema deve consentire al team leader di caricare nuovi file nel repository del team.
- Il sistema deve consentire al team leader di modificare i file precedentemente caricati.
- Il sistema deve consentire al team leader di eliminare file dal repository del team.
- Il sistema deve rendere i file del team disponibili in lettura a tutti i membri del team.
- Il sistema deve impedire l'accesso ai file a utenti non appartenenti al team proprietario.

## 3.3 Condivisione e Notifiche
> *Traces to: BR-001*

- Il sistema deve consentire a qualsiasi membro del team di condividere un file con uno o più altri membri dello stesso team.
- Il sistema deve inviare una notifica al destinatario al momento della ricezione di un file condiviso.
- Il sistema deve garantire che la notifica raggiunga il destinatario anche se non è connesso al momento della condivisione.
- Il sistema deve impedire la condivisione di file con utenti esterni al team del mittente.

## 3.4 Storico e Tracciabilità
> *Traces to: BR-001*

- Il sistema deve registrare ogni evento di condivisione, includendo mittente, destinatario, file condiviso e data/ora dell'operazione.
- Il sistema deve esporre uno storico consultabile delle condivisioni per ciascun file.
- Il sistema deve consentire a ogni utente di consultare lo storico delle condivisioni ricevute e inviate.
- Il sistema deve mantenere lo storico anche in caso di eliminazione del file originale.

---

# 4. Vincoli Tecnici

## 4.1 Infrastruttura e Hosting

- **Piattaforma di deployment obbligatoria:** Il sistema deve essere distribuito su Amazon Elastic Kubernetes Service (EKS) su infrastruttura AWS, in conformità con la piattaforma cloud adottata dal cliente.
- **Architettura containerizzata:** Tutte le componenti del sistema devono essere pacchettizzate come container ed eseguite in ambiente Kubernetes.

## 4.2 Integrazione e Interoperabilità

- **Autenticazione tramite Keycloak:** Il sistema deve integrarsi con l'istanza Keycloak esistente del cliente per la gestione dell'autenticazione e dell'autorizzazione. Non è ammessa l'implementazione di un sistema di autenticazione proprietario.
- **Gestione ruoli via Keycloak:** I ruoli applicativi (team leader, membro) devono essere propagati attraverso i token Keycloak; il sistema deve leggerli dal token senza duplicare la logica di assegnazione ruoli.

## 4.3 Dati e Compliance

- **Conformità GDPR:** Il sistema deve rispettare il Regolamento Generale sulla Protezione dei Dati (GDPR). I dati personali dei dipendenti (identità, storico attività) devono essere trattati nel rispetto dei principi di minimizzazione, finalità e sicurezza.
- **Residenza dei dati:** I dati devono essere archiviati esclusivamente in regioni AWS nell'area UE (es. eu-west, eu-central). Il trasferimento di dati personali al di fuori dello Spazio Economico Europeo non è ammesso salvo adeguate garanzie contrattuali.

## 4.4 Sicurezza e Accesso

- **Autenticazione obbligatoria:** Nessuna funzionalità del sistema è accessibile in modalità anonima; l'accesso richiede autenticazione tramite Keycloak.
- **Autorizzazione basata su ruolo:** Le operazioni disponibili per ciascun utente sono determinate dal ruolo ricevuto via Keycloak; il sistema non deve implementare logiche di autorizzazione parallele o ridondanti.

## 4.5 Sviluppo e Delivery

- **Stack tecnologico backend obbligatorio:** Il backend deve essere sviluppato in Java 21 con framework Spring Boot. Linguaggi o framework alternativi per il backend non sono ammessi.
- **Stack tecnologico frontend obbligatorio:** Il frontend deve essere sviluppato con il framework Angular. Framework alternativi per il frontend non sono ammessi.

## 4.6 Vendor e Licenze

- **Nessun vincolo di vendor o licenza:** Non esistono contratti con vendor specifici, licenze preesistenti o tecnologie esplicitamente escluse. Le scelte tecnologiche non ancora vincolate rimangono aperte alla progettazione.

---

# 5. Requisiti Non Funzionali

## 5.1 Performance
> *Traces to: BR-001, BR-002*

| ID | Requisito | Target / Condizione | Applicabile Quando |
|---|---|---|---|
| NFR-001 | Reattività operazioni | Le operazioni principali (caricamento, download, consultazione storico) devono essere percepite come fluide e senza latenze evidenti `[TO BE REFINED]` | Orario lavorativo, fino a 200 utenti concorrenti |

## 5.2 Disponibilità e Affidabilità
> *Traces to: BR-001*

| ID | Requisito | Target / Condizione | Applicabile Quando |
|---|---|---|---|
| NFR-002 | Disponibilità operativa | Il sistema non deve subire interruzioni non pianificate durante le ore lavorative `[TO BE REFINED: SLA %]` | Orario lavorativo aziendale |

## 5.3 Sicurezza
> *Traces to: BR-001, BR-002*

| ID | Requisito | Target / Condizione | Applicabile Quando |
|---|---|---|---|
| NFR-003 | Scadenza sessione | Le sessioni utente devono scadere automaticamente dopo 30 minuti di inattività | Per tutti gli utenti autenticati, in qualsiasi momento |
| NFR-004 | Protezione file | I file devono essere accessibili esclusivamente agli utenti autorizzati del team proprietario, anche a livello di storage | In qualsiasi momento, per qualsiasi operazione di accesso |

## 5.4 Scalabilità
> *Traces to: BR-001, BR-002*

| ID | Requisito | Target / Condizione | Applicabile Quando |
|---|---|---|---|
| NFR-005 | Capacità utenti | Il sistema deve supportare fino a 200 utenti concorrenti senza degradazione delle prestazioni | Durante l'orario lavorativo |

## 5.5 Manutenibilità e Operabilità
> *Traces to: BR-001, BR-002*

| ID | Requisito | Target / Condizione | Applicabile Quando |
|---|---|---|---|
| NFR-006 | Alerting malfunzionamenti | Il sistema deve inviare notifiche automatiche all'amministratore di sistema e al capo progetto in caso di errori critici o indisponibilità | Al verificarsi di anomalie, in qualsiasi momento |

## 5.6 Usabilità
> *Traces to: BR-001, BR-002*

| ID | Requisito | Target / Condizione | Applicabile Quando |
|---|---|---|---|
| NFR-007 | Semplicità d'uso | Utenti con qualsiasi livello di esperienza digitale devono poter completare le operazioni principali senza formazione preliminare `[TO BE REFINED]` | Per tutti gli utenti, durante l'uso quotidiano |

*Nota: requisiti di accessibilità WCAG non applicabili — applicazione ad uso interno privato.*

## 5.7 Compliance e Auditabilità
> *Traces to: BR-001, BR-002*

| ID | Requisito | Target / Condizione | Applicabile Quando |
|---|---|---|---|
| NFR-008 | Diritti degli interessati GDPR | Il sistema deve consentire la gestione dei diritti degli interessati (accesso, cancellazione, portabilità dei dati personali) `[TO BE REFINED con DPO]` | Su richiesta degli interessati |
| NFR-009 | Conservazione dati | I dati personali e i log delle condivisioni devono essere conservati per il periodo minimo necessario, nel rispetto della normativa GDPR `[TO BE REFINED]` | Per tutta la durata del trattamento |

---

# 6. Requisiti di Customer Experience

## 6.1 Tipologie di Utente e Contesto d'Uso

- Gli utenti sono dipendenti aziendali organizzati in team, con due profili distinti: team leader (gestione file) e membro del team (lettura e condivisione).
- L'applicazione è utilizzata esclusivamente da PC desktop, in ufficio o in smart working; il supporto per dispositivi mobili o tablet non è richiesto.
- Tutti gli utenti devono poter accedere all'applicazione tramite browser desktop senza installazione di software aggiuntivo.

## 6.2 Azioni Principali dell'Utente

- Gli utenti devono poter visualizzare immediatamente, al momento dell'accesso, la lista completa dei file di propria pertinenza.
- I team leader devono poter caricare, modificare ed eliminare file tramite operazioni accessibili direttamente dalla lista file, senza navigare attraverso schermate multiple.
- I membri del team devono poter condividere un file con altri membri del proprio team tramite un'azione chiaramente identificabile nella lista file.

## 6.3 Contenuto Visibile nella Lista File

- Per ogni file, l'interfaccia deve mostrare: nome del file, data di caricamento, nome di chi lo ha caricato e stato di lettura (letto / non letto).
- L'informazione relativa ai destinatari della condivisione non deve essere visibile nella lista file.
- Le informazioni essenziali di ogni file devono essere leggibili senza necessità di aprire viste di dettaglio per le operazioni quotidiane.

## 6.4 Notifiche

- Gli utenti devono essere informati della ricezione di nuovi file tramite un'icona con contatore nella barra di navigazione, analoga al meccanismo delle notifiche email.
- La notifica deve restare visibile fino a quando l'utente non consulta il file ricevuto (stato "non letto").
- Il contatore deve aggiornarsi automaticamente senza necessità di ricaricare la pagina.

## 6.5 Semplicità d'Uso

- Gli utenti devono poter completare le operazioni principali (visualizzare file, condividere, consultare storico) in non più di tre interazioni a partire dalla schermata principale.
- L'interfaccia non deve richiedere formazione preliminare: le funzionalità devono essere riconoscibili dall'etichetta e dal contesto visivo senza necessità di guide o tutorial.

*Note: accessibilità WCAG non applicabile (uso interno privato). Requisiti di brand/visual identity non specificati.*
