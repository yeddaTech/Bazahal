# 🛒 Progetto Bazahal – Guida per il Team
link webapp: https://bazahal.vercel.app/
l'app dopo 20 min di inattivita di utenza va in sleep, se accedi e non trovi il sito 
aspetta 30 secondi e riprova e andra.

VERCEL hosta il sito ed è collegato a github, ogni modificata applicata su github si applichera al sito.
neon hosta il database che dovrai connettere al tuo dbeaver (contatta yedda).

> ⚠️ **Leggere con attenzione prima di lavorare sul progetto**

Ciao!  
Se stai leggendo questo file, sei un collaboratore del progetto **Bazahal** (un aggregatore di prodotti Halal in affiliazione).

Attualmente il sito è **ONLINE e PUBBLICO** su VERCEL.

Questo significa che:

> ogni volta che il codice viene aggiornato sul ramo principale (`main`), il sito vero si aggiorna automaticamente.

Per evitare di **rompere il sito per tutti gli utenti**, devi configurare il progetto **in locale sul tuo PC** per fare le prove.

Segui questa guida **passo per passo senza saltare nulla**.

---

# 🧰 Passo 0 – Programmi necessari

Prima di iniziare assicurati di avere installato:

1. **Go (Golang)**  
   → linguaggio usato per il backend

2. **Git**  
   → per scaricare e inviare il codice

3. **PostgreSQL**  
   → database del progetto

4. **DBeaver** oppure **pgAdmin**  
   → per visualizzare il database

5. **VS Code**  
   → per scrivere il codice

---

# 📥 Passo 1 – Scaricare il codice

Apri il **terminale** (oppure il terminale di VS Code) e scrivi:

```bash
git clone https://github.com/yeddaTech/Bazahal.git
cd Bazahal
```

Questo scaricherà il progetto sul tuo computer.

---

# 🗄️ Passo 2 – Configurare il database (fondamentale)

Il codice capisce automaticamente se sta girando:

- sul **server online**
- oppure sul **tuo PC**

Se sei sul tuo PC cercherà un **database locale** con credenziali precise.

Se il database **non esiste**, il sito **non partirà**.

---

## Creazione del database

1️⃣ Apri **DBeaver** oppure **pgAdmin**

2️⃣ Crea una nuova connessione PostgreSQL

```
Host: localhost
User: postgres
Password: aicha
```

⚠️ Se hai usato una password diversa durante l'installazione di PostgreSQL, devi:

- cambiarla temporaneamente
- oppure avvisare l'amministratore del progetto

---

3️⃣ Crea un nuovo database chiamato **esattamente**

```
halalshop
```

---

4️⃣ Apri un **Nuovo Script SQL**

Assicurati di essere **dentro il database `halalshop`**.

Copia e incolla questo codice e premi **Run (▶)**:

```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    image_url TEXT
);
```

Se vedi un messaggio di successo, il database è pronto.

---

# 🚀 Passo 3 – Avviare il sito in locale

Ora possiamo accendere il sito sul tuo PC.

1️⃣ Apri **VS Code**

2️⃣ Apri il **terminale integrato**

```
Terminale → Nuovo terminale
```

3️⃣ Assicurati di essere nella cartella:

```
Bazahal
```

4️⃣ Avvia il server con:

```bash
go run .
```

Se non funziona prova:

```bash
go run main.go
```

---

Se nel terminale compare:

```
Connesso al database con successo!
```

allora funziona tutto.

Apri il browser e vai su:

```
http://localhost:8080
```

Questo è il **tuo sito privato di sviluppo**.

Qualsiasi modifica qui **non rompe il sito online**.

---

# 🗺️ Passo 4 – Orientarsi nel codice

Se non sai dove mettere le mani, ecco la mappa:

### main.go
Il cuore dell'applicazione.  
Avvia il server.

### cartella `database/`
Contiene `db.go` che gestisce la connessione al database.

⚠️ Non modificarlo se non sai cosa stai facendo.

### cartella `templates/` o `static/`

Qui trovi:

- HTML
- CSS
- eventuale JavaScript

Puoi modificare:

- colori
- bottoni
- grafica del sito

Salva il file e **ricarica la pagina nel browser** per vedere i cambiamenti.

---

# ⛔ Passo 5 – Regola d’oro (NON DISTRUGGERE IL SITO)

Il sito online ascolta il ramo:

```
main
```

⚠️ **NON devi MAI lavorare direttamente su `main`.**

⚠️ **NON fare mai:**

```bash
git push origin main
```

Se carichi codice rotto:

> il sito online si rompe entro circa **2 minuti**.

---

# 🌿 Git Flow da seguire

## Fase A – Aggiornare il progetto

Prima di lavorare:

```bash
git pull origin main
```

---

## Fase B – Creare il tuo branch

Crea una copia privata del progetto:

```bash
git checkout -b aggiunta-nuova-grafica
```

Esempi di nomi:

```
fix-bottone-rosso
aggiunta-pagina-prodotti
miglioramento-navbar
```

---

## Fase C – Lavorare

Modifica i file su **VS Code**.

Controlla sempre che funzioni su:

```
http://localhost:8080
```

---

## Fase D – Salvare il lavoro

Quando tutto funziona:

```bash
git add .
git commit -m "Spiega cosa hai modificato"
```

Esempio:

```
git commit -m "Cambio colore titolo homepage"
```

---

## Fase E – Inviare il branch

```bash
git push origin aggiunta-nuova-grafica
```

---

## Fase F – Creare la Pull Request

1️⃣ Vai su **GitHub**

2️⃣ Vedrai il bottone:

```
Compare & Pull Request
```

3️⃣ Cliccalo

4️⃣ Scrivi cosa hai fatto

5️⃣ Premi:

```
Create Pull Request
```

---

# ✅ Dopo la Pull Request

L'amministratore controllerà il codice.

Se è tutto ok:

```
Merge
```

La modifica verrà pubblicata sul sito ufficiale.

---

# 🚀 Buon lavoro

Buon coding e benvenuto nel progetto **Bazahal**!