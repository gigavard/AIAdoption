# scenario_catalog

## scenarios

| scenario_id | uc_id | scenario_name | scenario_type | related_feature_ids | trigger | expected_outcome | exception_notes |
|---|---|---|---|---|---|---|---|
| SCN-001 | UC-001 | Caricamento file completato con successo | Main | FEAT-003; FEAT-001 | Il team leader seleziona un file dal proprio dispositivo e avvia il caricamento. | Il file è salvato nel repository del team e appare nella lista file di tutti i membri. | — |
| SCN-002 | UC-001 | Tentativo di caricamento da parte di un membro | Exception | FEAT-001 | Un membro del team tenta di caricare un file. | Il sistema rifiuta l'operazione con messaggio di errore "operazione non consentita". | L'operazione non deve avere alcun effetto sul repository. |
| SCN-003 | UC-002 | Modifica file completata con successo | Main | FEAT-004; FEAT-001 | Il team leader seleziona un file esistente e carica una versione aggiornata. | Il file nel repository è aggiornato; la versione precedente è sostituita. | — |
| SCN-004 | UC-002 | Tentativo di modifica da parte di un membro | Exception | FEAT-001 | Un membro del team tenta di modificare un file. | Il sistema rifiuta l'operazione con messaggio di errore "operazione non consentita". | Il file non deve subire modifiche. |
| SCN-005 | UC-003 | Eliminazione file completata con successo | Main | FEAT-005; FEAT-009 | Il team leader seleziona un file e ne richiede l'eliminazione. | Il file è rimosso dal repository; lo storico delle condivisioni associate è mantenuto. | — |
| SCN-006 | UC-003 | Eliminazione file con condivisioni pregresse | Alternate | FEAT-005; FEAT-009 | Il team leader elimina un file che era stato precedentemente condiviso. | Il file non è più accessibile, ma lo storico delle condivisioni è consultabile. | Il record storico deve referenziare il file con indicazione "file eliminato". |
| SCN-007 | UC-004 | Visualizzazione lista file del team | Main | FEAT-006; FEAT-002 | L'utente accede alla sezione file dell'applicazione. | La lista mostra tutti i file del team con nome, data, autore e stato di lettura. | — |
| SCN-008 | UC-004 | Accesso negato a file di team diverso | Exception | FEAT-002 | Un utente tenta di accedere ai file di un team a cui non appartiene. | Il sistema nega l'accesso con messaggio di errore "accesso non consentito". | — |
| SCN-009 | UC-005 | Condivisione file con membro dello stesso team | Main | FEAT-007; FEAT-008; FEAT-009 | Un utente seleziona un file e uno o più destinatari dello stesso team, quindi avvia la condivisione. | L'evento è registrato nello storico; i destinatari ricevono notifica in-app. | — |
| SCN-010 | UC-005 | Tentativo di condivisione con membro di altro team | Exception | FEAT-002 | Un utente tenta di condividere un file con un utente appartenente a un team diverso. | Il sistema rifiuta l'operazione con messaggio di errore "destinatario non valido". | I destinatari dello stesso team eventualmente inclusi nella stessa richiesta non devono essere penalizzati. |
| SCN-011 | UC-006 | Ricezione notifica mentre l'utente è connesso | Main | FEAT-008 | Un file viene condiviso con l'utente mentre è attivo nell'applicazione. | Il contatore notifiche si aggiorna in tempo reale; il file appare come "non letto" nella lista. | — |
| SCN-012 | UC-006 | Ricezione notifica mentre l'utente non è connesso | Alternate | FEAT-008 | Un file viene condiviso con l'utente mentre non è connesso all'applicazione. | Al successivo accesso, il contatore mostra la notifica non letta; il file è marcato come "non letto". | La notifica deve essere persistita lato server fino alla consultazione. |
| SCN-013 | UC-007 | Consultazione storico condivisioni di un file esistente | Main | FEAT-009 | L'utente seleziona un file e richiede la visualizzazione del suo storico condivisioni. | La lista cronologica degli eventi (mittente, destinatario, data/ora) è visualizzata. | — |
| SCN-014 | UC-007 | Consultazione storico di un file eliminato | Alternate | FEAT-009 | L'utente consulta lo storico di un file che è stato eliminato dal repository. | Lo storico è ancora accessibile con indicazione "file eliminato". | Il record storico è disaccoppiato dall'esistenza fisica del file. |

## scenario_steps

| scenario_id | step_id | step_order | actor_id | action | system_response |
|---|---|---|---|---|---|
| SCN-001 | STEP-001 | 1 | role-001 | Accede alla sezione di caricamento file. | Il sistema verifica il ruolo dell'utente (team leader). |
| SCN-001 | STEP-002 | 2 | role-001 | Seleziona il file dal proprio dispositivo e conferma il caricamento. | Il sistema salva il file nel repository del team con metadati (nome, data, autore). |
| SCN-001 | STEP-003 | 3 | System | — | Il file appare nella lista file di tutti i membri del team. |
| SCN-002 | STEP-001 | 1 | role-002 | Tenta di accedere alla funzione di caricamento file. | Il sistema verifica il ruolo: membro non autorizzato. |
| SCN-002 | STEP-002 | 2 | System | — | Il sistema restituisce messaggio di errore "operazione non consentita". |
| SCN-003 | STEP-001 | 1 | role-001 | Seleziona un file esistente nel repository e richiede la modifica. | Il sistema verifica il ruolo (team leader) e carica il file per la modifica. |
| SCN-003 | STEP-002 | 2 | role-001 | Carica il nuovo contenuto o aggiorna i metadati e conferma. | Il sistema sostituisce il file precedente con la nuova versione. |
| SCN-003 | STEP-003 | 3 | System | — | La lista file mostra la versione aggiornata. |
| SCN-004 | STEP-001 | 1 | role-002 | Tenta di accedere alla funzione di modifica su un file. | Il sistema verifica il ruolo: membro non autorizzato. |
| SCN-004 | STEP-002 | 2 | System | — | Il sistema restituisce messaggio di errore "operazione non consentita". |
| SCN-005 | STEP-001 | 1 | role-001 | Seleziona un file nel repository e richiede l'eliminazione. | Il sistema verifica il ruolo (team leader) e chiede conferma. |
| SCN-005 | STEP-002 | 2 | role-001 | Conferma l'eliminazione. | Il sistema rimuove il file dal repository; il record storico delle condivisioni è mantenuto. |
| SCN-005 | STEP-003 | 3 | System | — | Il file non appare più nella lista file del team. |
| SCN-006 | STEP-001 | 1 | role-001 | Elimina un file che aveva condivisioni pregresse (come SCN-005). | Il sistema rimuove il file; i record storici associati sono marcati con "file eliminato". |
| SCN-006 | STEP-002 | 2 | System | — | Lo storico condivisioni del file è ancora consultabile con indicazione "file eliminato". |
| SCN-007 | STEP-001 | 1 | role-001; role-002 | Accede alla sezione file dell'applicazione. | Il sistema recupera la lista dei file del team dell'utente. |
| SCN-007 | STEP-002 | 2 | System | — | La lista mostra: nome file, data di caricamento, autore, stato di lettura (letto/non letto). |
| SCN-008 | STEP-001 | 1 | role-001; role-002 | Tenta di accedere ai file di un team diverso dal proprio. | Il sistema verifica l'appartenenza al team: accesso non consentito. |
| SCN-008 | STEP-002 | 2 | System | — | Il sistema restituisce messaggio di errore "accesso non consentito". |
| SCN-009 | STEP-001 | 1 | role-001; role-002 | Seleziona un file dalla lista e richiede la condivisione. | Il sistema presenta la lista dei membri del proprio team selezionabili come destinatari. |
| SCN-009 | STEP-002 | 2 | role-001; role-002 | Seleziona uno o più destinatari e conferma la condivisione. | Il sistema registra l'evento di condivisione (mittente, destinatari, file, data/ora). |
| SCN-009 | STEP-003 | 3 | System | — | Il sistema invia notifica in-app a ciascun destinatario; il contatore notifiche si aggiorna. |
| SCN-010 | STEP-001 | 1 | role-001; role-002 | Tenta di condividere un file con un utente di un altro team. | Il sistema verifica l'appartenenza al team del destinatario. |
| SCN-010 | STEP-002 | 2 | System | — | Il sistema rifiuta l'operazione con messaggio "destinatario non valido". |
| SCN-011 | STEP-001 | 1 | System | Un file viene condiviso con l'utente mentre è attivo nell'applicazione. | Il sistema aggiorna il contatore notifiche in tempo reale. |
| SCN-011 | STEP-002 | 2 | System | — | Il file appare nella lista con stato "non letto". |
| SCN-012 | STEP-001 | 1 | System | Un file viene condiviso con l'utente mentre non è connesso. | Il sistema persiste la notifica lato server. |
| SCN-012 | STEP-002 | 2 | role-001; role-002 | L'utente accede successivamente all'applicazione. | Il sistema mostra il contatore notifiche aggiornato; il file è marcato come "non letto". |
| SCN-013 | STEP-001 | 1 | role-001; role-002 | Seleziona un file e richiede la visualizzazione dello storico condivisioni. | Il sistema recupera tutti gli eventi di condivisione associati al file. |
| SCN-013 | STEP-002 | 2 | System | — | La lista cronologica è mostrata: mittente, destinatario, data/ora per ciascun evento. |
| SCN-014 | STEP-001 | 1 | role-001; role-002 | Richiede lo storico di un file che è stato eliminato. | Il sistema recupera i record storici dal database (disaccoppiati dal file fisico). |
| SCN-014 | STEP-002 | 2 | System | — | Lo storico è mostrato con indicazione "file eliminato" in corrispondenza del file. |

```
SCENARIOS_COMPLETED
```
