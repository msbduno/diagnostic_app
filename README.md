# Mini Diagnostic App

> Application de diagnostic hardware automatisé pour Mac 

---
## Démonstration


https://github.com/user-attachments/assets/388aa4a2-e71f-4d1f-9655-27957fd3673f


---


## Vue d'ensemble

Cette application permet de réaliser un diagnostic automatisé des composants hardware d'un Mac. 
Elle collecte les informations essentielles (CPU, RAM, Stockage, Batterie) et les transmet automatiquement à un backend pour stockage et analyse.

### Objectifs
- **Rapidité** : Test complet en moins de 1 minute
- **Automatisation** : Zéro intervention manuelle après le lancement
- **Traçabilité** : Historique complet de chaque machine testée
- **Simplicité** : Interface intuitive en un clic

---

## Fonctionnalités

### Application macOS (Swift)
- Collecte automatique des informations hardware
- Interface utilisateur simple et claire
- Barre de progression en temps réel
- Affichage des résultats du diagnostic
- Envoi automatique au backend
- Gestion des erreurs réseau

### Backend (Go)
- API REST pour réception des diagnostics
- Stockage SQLite des résultats
- Validation des données entrantes
- Logs détaillés pour debug


### Tests Réalisés
| Composant | Informations Collectées |
|-----------|------------------------|
| **CPU** | Modèle, Nombre de cœurs, Architecture |
| **RAM** | Total, Utilisé, Disponible |
| **Stockage** | Capacité totale, Espace utilisé, Espace libre |
| **Batterie** | État de santé, Cycles de charge, Niveau actuel |
| **Système** | Numéro de série, Version macOS |

---

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                     Mac à diagnostiquer                     │
│                                                             │
│  ┌───────────────────────────────────────────────────────┐  │
│  │          Application Swift (SwiftUI)                  │  │
│  │                                                       │  │
│  │  1. Lance tests hardware (IOKit, ProcessInfo)         │  │
│  │  2. Affiche progression                               │  │
│  │  3. Présente résultats                                │  │
│  └──────────────────────┬────────────────────────────────┘  │
│                         │ HTTP POST                         │
└─────────────────────────┼───────────────────────────────────┘
                          │
                          ▼
         ┌────────────────────────────────────┐
         │      Backend API (Go + Gin)        │
         │                                    │
         │  • Endpoint POST /api/diagnostics  │
         │  • Validation des données          │
         │  • Logs et monitoring              │
         └────────────┬───────────────────────┘
                      │
                      ▼
         ┌────────────────────────────────────┐
         │     Base de données (SQLite)       │
         │                                    │
         │  • Stockage des diagnostics        │
         │  • Historique des tests            │
         │  • Fichier: diagnostics.db         │
         └────────────────────────────────────┘
```

---

## Stack Technique

### Frontend (Application macOS)
- **Langage** : Swift 5.9+
- **Framework UI** : SwiftUI
- **APIs système** :
    - `IOKit` : Informations hardware bas niveau
    - `ProcessInfo` : CPU et mémoire
    - `FileManager` : Stockage
    - `IOPowerSources` : Batterie
- **Réseau** : URLSession (HTTP natif)

### Backend
- **Langage** : Go 1.21+
- **Framework Web** : Gin (API REST)
- **Base de données** : SQLite3
- **Dépendances** :
  ```go
  github.com/gin-gonic/gin
  github.com/mattn/go-sqlite3
  github.com/google/uuid
  ```

### Outils de Développement
- **Xcode** 15+ (pour app Swift)
- **Go** 1.21+
- **SQLite** 3 (inclus avec macOS)
- **Git** pour versioning

---

## Installation

### Prérequis
- macOS 13.0 (Ventura) ou supérieur
- Xcode 15+ avec Command Line Tools
- Go 1.21+ installé ([télécharger](https://go.dev/dl/))

### 1. Cloner le repository

```bash
git clone https://github.com/votre-username/mini-diagnostic-app.git
cd mini-diagnostic-app
```

---

## Utilisation

### Démarrage du Backend

```bash
cd backend
./diagnostic-api
```

Le serveur démarre sur `http://localhost:8080`


### Lancement de l'Application macOS

1. Ouvrir l'app depuis Xcode ou depuis le dossier Build
2. Cliquer sur **"Lancer le Test"**
3. Patienter pendant le diagnostic (progression affichée)
4. Consulter les résultats à l'écran
5. Les données sont automatiquement envoyées au backend

---

## API Documentation

### Base URL
```
http://localhost:8080/api
```

### Endpoints

#### POST /api/diagnostics

Enregistre un nouveau diagnostic hardware.

**Request Body :**
```json
{
  "serial_number": "C02XYZ123ABC",
  "cpu_model": "Apple M2",
  "cpu_cores": 8,
  "ram_total_gb": 16.0,
  "ram_used_gb": 8.5,
  "storage_total_gb": 512.0,
  "storage_used_gb": 256.0,
  "battery_health": "Good",
  "battery_cycle_count": 42,
  "battery_percentage": 85,
  "test_duration_seconds": 45,
  "status": "completed"
}
```

**Response (201 Created) :**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "serial_number": "C02XYZ123ABC",
  "timestamp": "2025-12-09T14:32:15Z",
  "status": "completed"
}
```

**Codes d'erreur :**
- `400 Bad Request` : Données invalides
- `500 Internal Server Error` : Erreur serveur/BDD

#### GET /api/diagnostics/:serial_number

Récupère l'historique des diagnostics d'une machine.

**Response (200 OK) :**
```json
[
  {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "serial_number": "C02XYZ123ABC",
    "timestamp": "2025-12-09T14:32:15Z",
    "cpu_model": "Apple M2",
    "status": "completed"
  }
]
```

---

## Modèle de Données

### Table : `diagnostic_tests`

| Colonne | Type | Description |
|---------|------|-------------|
| `id` | TEXT PRIMARY KEY | UUID unique du test |
| `serial_number` | TEXT NOT NULL | Numéro de série du Mac |
| `timestamp` | DATETIME | Date/heure du diagnostic |
| `cpu_model` | TEXT | Modèle du processeur |
| `cpu_cores` | INTEGER | Nombre de cœurs CPU |
| `ram_total_gb` | REAL | RAM totale (Go) |
| `ram_used_gb` | REAL | RAM utilisée (Go) |
| `storage_total_gb` | REAL | Stockage total (Go) |
| `storage_used_gb` | REAL | Stockage utilisé (Go) |
| `battery_health` | TEXT | État batterie (Good/Fair/Poor) |
| `battery_cycle_count` | INTEGER | Nombre de cycles de charge |
| `battery_percentage` | INTEGER | Niveau de batterie (%) |
| `test_duration_seconds` | INTEGER | Durée du test (secondes) |
| `status` | TEXT | Statut (completed/failed) |

### Exemple SQL

```sql
CREATE TABLE diagnostic_tests (
    id TEXT PRIMARY KEY,
    serial_number TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    cpu_model TEXT,
    cpu_cores INTEGER,
    ram_total_gb REAL,
    ram_used_gb REAL,
    storage_total_gb REAL,
    storage_used_gb REAL,
    battery_health TEXT,
    battery_cycle_count INTEGER,
    battery_percentage INTEGER,
    test_duration_seconds INTEGER,
    status TEXT,
    UNIQUE(serial_number, timestamp)
);
```

---

## Structure du Projet

```
mini-diagnostic-app/
├── README.md
├── backend/
│   ├── main.go              # Point d'entrée API
│   ├── models/
│   │   └── diagnostic.go    # Modèle de données
│   ├── handlers/
│   │   └── diagnostics.go   # Handlers API
│   ├── database/
│   │   └── sqlite.go        # Connexion BDD
│   ├── go.mod
│   ├── go.sum
│   └── diagnostics.db       # BDD SQLite (généré)
│
├── macos-app/
│   ├── DiagnosticApp.xcodeproj
│   ├── DiagnosticApp/
│   │   ├── ContentView.swift       # Interface principale
│   │   ├── DiagnosticService.swift # Tests hardware
│   │   ├── APIClient.swift         # Connexion backend
│   │   ├── Models.swift            # Modèles de données
│   │   └── DiagnosticAppApp.swift  # Point d'entrée
│   └── Assets.xcassets/
│
└── docs/
    ├── architecture.png
    └── screenshots/
```

---
