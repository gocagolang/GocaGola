# GocaGola

## Description

GocaGola est un framework léger pour créer des API GO avec un système de plugins dynamiques. Il permet de :

- Charger dynamiquement des handlers HTTP depuis des plugins Go
- Recharger automatiquement les plugins modifiés
- Simplifier la création d'API REST

## Installation

```bash
git clone https://github.com/YourUsername/YourProject.git
cd YourProject
go mod init gocagola
go get github.com/gocagolang/GocaGola/routing
```

## Utilisation

### Point d'accès simple

```go
// filepath: /api/users/main.go
package main

import "github.com/gin-gonic/gin"

func GET(c *gin.Context) {
    c.JSON(200, gin.H{
        "message": "Liste des utilisateurs",
    })
}
```

### Points d'accès avec paramètres

```go
// filepath: /api/users/:id/main.go
package main

import "github.com/gin-gonic/gin"

func GET(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{
        "message": "Détails de l'utilisateur",
        "id": id,
    })
}

func DELETE(c *gin.Context) {
    id := c.Param("id")
    c.JSON(200, gin.H{
        "message": "Utilisateur supprimé",
        "id": id,
    })
}
```

Exemple de réponse pour `GET /api/users/123` :
```json
{
    "message": "Détails de l'utilisateur",
    "id": "123"
}
```

Exemple de réponse pour `DELETE /api/users/123` :
```json
{
    "message": "Utilisateur supprimé",
    "id": "123"
}
```

### Démarrage du serveur

```go
// filepath: /main.go
package main

import "github.com/gocagolang/GocaGola/routing"

func main() {
    routing.Initialize("api")
}
```

Lancer le serveur :
```bash
go run main.go
```

## Fonctionnalités

- **Système de Plugins Dynamique** : Ajoutez de nouveaux points d'accès sans redémarrer le serveur
- **Rechargement à Chaud** : Recompilation et rechargement automatique des plugins modifiés
- **Convention plutôt que Configuration** : Structure de dossiers simple pour les points d'accès API
- **Basé sur Gin** : Framework HTTP haute performance
- **Facile à Développer** : Messages d'erreur clairs et journalisation détaillée

## Structure du Projet

```
/YourProject
├── api/
│   ├── ping/
│   │   └── main.go
│   └── users/
│       ├── main.go
│       └── :id/
│           └── main.go
└── main.go
```

## Contribuer

Les pull requests sont les bienvenues. Pour les changements majeurs, veuillez d'abord ouvrir une issue pour discuter des modifications souhaitées.

