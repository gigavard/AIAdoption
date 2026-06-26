# functional_requirements

## functional_areas

| area_id | area_name | area_description | related_br_ids |
|---|---|---|---|
| FA-001 | Gestione Utenti e Ruoli | Gestione dell'associazione utente-team-ruolo e applicazione delle autorizzazioni differenziate per ruolo. | BR-002 |
| FA-002 | Gestione File | Operazioni di caricamento, modifica ed eliminazione file eseguibili dal team leader; accesso in lettura per tutti i membri del team. | BR-001; BR-002 |
| FA-003 | Condivisione e Notifiche | Condivisione di file tra membri dello stesso team con notifica affidabile al destinatario. | BR-001 |
| FA-004 | Storico e Tracciabilità | Registrazione e consultazione dello storico degli eventi di condivisione, con persistenza anche dopo eliminazione del file. | BR-001 |

## requirement_catalog

| fr_id | area_id | requirement_name | requirement_description | primary_actor_ids | input_data | output_data | business_rules |
|---|---|---|---|---|---|---|---|
| FR-001 | FA-001 | Associazione Utente a Team e Ruolo | Il sistema deve associare ogni utente a un team e a un ruolo (team leader o membro). | role-001; role-002 | Identità utente (da Keycloak token) | Profilo utente con team e ruolo associati | Il ruolo è propagato dal token Keycloak; non è modificabile dall'applicazione. |
| FR-002 | FA-001 | Applicazione Autorizzazioni per Ruolo | Il sistema deve applicare le autorizzazioni in base al ruolo: i team leader possono caricare e modificare file; i membri possono solo leggere e condividere file. | role-001; role-002 | Ruolo utente (dal token) | Insieme di operazioni consentite | Le autorizzazioni sono determinate esclusivamente dal ruolo Keycloak. Nessuna logica parallela è ammessa. |
| FR-003 | FA-001 | Blocco Operazioni Non Autorizzate per Membri | Il sistema deve impedire ai membri di eseguire operazioni riservate ai team leader (caricamento di nuovi file, modifica). | role-002 | Richiesta di operazione non autorizzata | Risposta di errore (accesso negato) | Tentativo di caricamento o modifica da parte di un membro deve essere rifiutato con messaggio di errore appropriato. |
| FR-004 | FA-001 | Blocco Condivisione Inter-Team | Il sistema deve impedire la condivisione di file tra utenti appartenenti a team differenti. | role-001; role-002 | Richiesta di condivisione con utente di altro team | Risposta di errore (accesso negato) | La condivisione è consentita solo tra utenti dello stesso team. |
| FR-005 | FA-002 | Caricamento File da parte del Team Leader | Il sistema deve consentire al team leader di caricare nuovi file nel repository del team. | role-001 | File da caricare, metadati (nome, data) | File salvato nel repository del team; conferma di avvenuto caricamento | Solo il team leader può caricare nuovi file. |
| FR-006 | FA-002 | Modifica File Esistente da parte del Team Leader | Il sistema deve consentire al team leader di modificare i file precedentemente caricati. | role-001 | Identificativo file, nuovo contenuto o metadati | File aggiornato nel repository; conferma di avvenuta modifica | Solo il team leader può modificare file. La versione precedente viene sostituita. |
| FR-007 | FA-002 | Eliminazione File da parte del Team Leader | Il sistema deve consentire al team leader di eliminare file dal repository del team. | role-001 | Identificativo file da eliminare | File rimosso dal repository; conferma di avvenuta eliminazione | Solo il team leader può eliminare file. Lo storico delle condivisioni sul file eliminato deve essere mantenuto. |
| FR-008 | FA-002 | Accesso in Lettura ai File per i Membri del Team | Il sistema deve rendere i file del team disponibili in lettura a tutti i membri del team. | role-001; role-002 | Identificativo team / file | Lista file del team o contenuto del file selezionato | Ogni utente può accedere solo ai file del proprio team. |
| FR-009 | FA-002 | Blocco Accesso ai File da Utenti Esterni al Team | Il sistema deve impedire l'accesso ai file a utenti non appartenenti al team proprietario. | role-001; role-002 | Richiesta di accesso da utente di altro team | Risposta di errore (accesso negato) | L'appartenenza al team è verificata tramite il token Keycloak. |
| FR-010 | FA-003 | Condivisione File Intra-Team | Il sistema deve consentire a qualsiasi membro del team di condividere un file con uno o più altri membri dello stesso team. | role-001; role-002 | Identificativo file, lista destinatari (dello stesso team) | Evento di condivisione registrato; notifica inviata ai destinatari | La condivisione è possibile solo verso utenti dello stesso team. |
| FR-011 | FA-003 | Notifica Ricezione File Condiviso | Il sistema deve inviare una notifica al destinatario al momento della ricezione di un file condiviso. | role-001; role-002 | Evento di condivisione | Notifica in-app con icona e contatore aggiornati | La notifica deve essere visibile tramite icona con contatore nella barra di navigazione. |
| FR-012 | FA-003 | Notifica Garantita in Assenza di Connessione | Il sistema deve garantire che la notifica raggiunga il destinatario anche se non è connesso al momento della condivisione. | role-001; role-002 | Evento di condivisione avvenuto in assenza del destinatario | Notifica in-app disponibile al successivo accesso del destinatario | La notifica deve persistere fino alla consultazione del file da parte del destinatario. |
| FR-013 | FA-003 | Blocco Condivisione verso Utenti Esterni al Team | Il sistema deve impedire la condivisione di file con utenti esterni al team del mittente. | role-001; role-002 | Richiesta di condivisione verso utente di altro team | Risposta di errore (operazione non consentita) | Coerente con FR-004; la verifica avviene al momento dell'operazione di condivisione. |
| FR-014 | FA-004 | Registrazione Evento di Condivisione | Il sistema deve registrare ogni evento di condivisione, includendo mittente, destinatario, file condiviso e data/ora dell'operazione. | role-001; role-002 | Evento di condivisione | Record storico persistito | Il record deve essere immutabile e non cancellabile dall'utente. |
| FR-015 | FA-004 | Storico Condivisioni per File | Il sistema deve esporre uno storico consultabile delle condivisioni per ciascun file. | role-001; role-002 | Identificativo file | Lista cronologica degli eventi di condivisione relativi al file | Lo storico è visibile a tutti i membri del team. |
| FR-016 | FA-004 | Storico Condivisioni per Utente | Il sistema deve consentire a ogni utente di consultare lo storico delle condivisioni ricevute e inviate. | role-001; role-002 | Identificativo utente | Lista degli eventi di condivisione inviati e ricevuti dall'utente | Ogni utente vede solo le proprie condivisioni (inviate e ricevute). |
| FR-017 | FA-004 | Persistenza Storico dopo Eliminazione File | Il sistema deve mantenere lo storico delle condivisioni anche in caso di eliminazione del file originale. | role-001 | Evento di eliminazione file | Record storico mantenuto; file non più accessibile | Il record storico è disaccoppiato dall'esistenza fisica del file. |

```
FUNCTIONAL_REQUIREMENTS_COMPLETED
```
