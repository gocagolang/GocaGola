# GocaGola

## Description

GocaGola est un framework léger pour créer des API GO avec un système de plugins dynamiques. Il permet de :

- Charger dynamiquement des handlers HTTP depuis des plugins Go
- Recharger automatiquement les plugins modifiés
- Attribuer des middlewares spécifiques à des routes
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

### Ajout de middlewares spécifiques aux routes

Vous pouvez définir des middlewares et les associer à des routes spécifiques en utilisant la fonction `MiddlewareResolver`. Voici un exemple :

```go
// filepath: /main.go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/gocagolang/GocaGola/routing"
)

func main() {
    // Définir la logique pour associer des middlewares à des routes
    routing.MiddlewareResolver = func(routePath string) []gin.HandlerFunc {
        if routePath == "/api/admin" {
            return []gin.HandlerFunc{AuthMiddleware, LoggingMiddleware}
        } else if routePath == "/api/public" {
            return []gin.HandlerFunc{LoggingMiddleware}
        }
        return []gin.HandlerFunc{}
    }

    // Initialiser le routeur
    routing.Initialize("api", "middlewares")
}

// Exemple de middleware d'authentification
func AuthMiddleware(c *gin.Context) {
    if c.GetHeader("Authorization") == "" {
        c.AbortWithStatusJSON(401, gin.H{"error": "Unauthorized"})
        return
    }
    c.Next()
}

// Exemple de middleware de logging
func LoggingMiddleware(c *gin.Context) {
    c.Next()
}
```

### Démarrage du serveur

```go
// filepath: /main.go
package main

import "github.com/gocagolang/GocaGola/routing"

func main() {
    routing.Initialize("api", "middlewares")
}
```

Lancer le serveur :
```bash
go run main.go
```

## Fonctionnalités

- **Système de Plugins Dynamique** : Ajoutez de nouveaux points d'accès sans redémarrer le serveur
- **Rechargement à Chaud** : Recompilation et rechargement automatique des plugins modifiés
- **Middlewares Personnalisables** : Associez des middlewares spécifiques à des routes via `MiddlewareResolver`
- **Convention plutôt que Configuration** : Structure de dossiers simple pour les points d'accès API
- **Basé sur Gin** : Framework HTTP haute performance
- **Facile à Développer** : Messages d'erreur clairs et journalisation détaillée

## Structure du Projet

```
/YourProject
├── api/
│   ├── admin/
│   │   └── main.go
│   ├── public/
│   │   └── main.go
│   └── users/
│       ├── main.go
│       └── :id/
│           └── main.go
├── middlewares/
│   ├── auth.go
│   └── logging.go
└── main.go
```

## Contribuer

Les pull requests sont les bienvenues. Pour les changements majeurs, veuillez d'abord ouvrir une issue pour discuter des modifications souhaitées.

