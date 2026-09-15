# Secrets Module

Le module **Secrets** fournit un système sécurisé de gestion des secrets pour les applications multi-tenant.

Fonctionnalités principales :

- Stockage sécurisé des secrets
- Chiffrement AES-256-GCM
- Gestion des versions
- Rotation automatique des secrets
- Isolation multi-tenant
- Vault centralisé
- Middleware de protection
- Audit prêt à intégrer
- Compatible SaaS / Cloud Native

---

# Architecture

```text
core/secrets
│
├── secret.go
├── encryption.go
├── store.go
├── vault.go
├── rotation.go
├── middleware.go
├── errors.go
│
tests/secrets
├── secret_test.go
├── encryption_test.go
├── store_test.go
├── vault_test.go
├── rotation_test.go
└── middleware_test.go
```

---

# Secret

Un secret représente une donnée sensible :

- API Key
- Token JWT
- Password
- Database Password
- SMTP Password
- OAuth Secret
- Webhook Secret

Exemple :

```go
secret := secrets.NewSecret(
    "tenant-1",
    "stripe_api_key",
    "sk_live_xxxxx",
    "api_key",
)
```

---

# Validation

Chaque secret est validé avant stockage.

Vérifications :

- TenantID obligatoire
- Name obligatoire
- Value obligatoire
- Version > 0

Exemple :

```go
err := secret.Validate()
```

---

# Encryption Service

Le module utilise :

```text
AES-256-GCM
```

Taille de clé :

```text
32 bytes
```

Création d'une clé :

```go
key, err := secrets.GenerateKey()
```

Création du service :

```go
enc, err := secrets.NewEncryptionService(
    key,
)
```

---

# Chiffrement

```go
cipherText, err := enc.Encrypt(
    "my-secret",
)
```

---

# Déchiffrement

```go
plainText, err := enc.Decrypt(
    cipherText,
)
```

---

# Chiffrement d'un Secret

```go
err := enc.EncryptSecret(
    secret,
)
```

---

# Déchiffrement d'un Secret

```go
err := enc.DecryptSecret(
    secret,
)
```

---

# Store

Le stockage par défaut est :

```go
MemoryStore
```

Création :

```go
store := secrets.NewMemoryStore()
```

---

# Sauvegarde

```go
err := store.Save(
    ctx,
    secret,
)
```

---

# Recherche par ID

```go
secret, err := store.GetByID(
    ctx,
    id,
)
```

---

# Recherche par Nom

```go
secret, err := store.GetByName(
    ctx,
    "stripe_api_key",
)
```

---

# Liste

```go
items, err := store.List(
    ctx,
)
```

---

# Suppression

```go
err := store.Delete(
    ctx,
    secretID,
)
```

---

# Nettoyage complet

```go
err := store.Clear(
    ctx,
)
```

---

# Vault

Le Vault représente la couche métier principale.

Création :

```go
vault, err := secrets.NewVault(
    store,
    encryption,
)
```

---

# Création d'un Secret

```go
err := vault.CreateSecret(
    ctx,
    secret,
)
```

Le Vault :

- valide
- chiffre
- vérifie les doublons
- stocke

---

# Lecture d'un Secret

```go
secret, err := vault.GetSecret(
    ctx,
    tenantID,
    secretName,
)
```

Le secret est automatiquement déchiffré.

---

# Mise à jour

```go
err := vault.UpdateSecret(
    ctx,
    secret,
)
```

---

# Rotation

```go
err := vault.RotateSecret(
    ctx,
    tenantID,
    name,
    newValue,
)
```

Effets :

- nouvelle valeur
- version +1
- mise à jour timestamp
- rechiffrement

---

# Suppression

```go
err := vault.DeleteSecret(
    ctx,
    tenantID,
    name,
)
```

---

# Liste des secrets

```go
items, err := vault.ListSecrets(
    ctx,
    tenantID,
)
```

---

# Comptage

```go
count, err := vault.SecretCount(
    ctx,
    tenantID,
)
```

---

# Rotation Manager

Le Rotation Manager permet de détecter les secrets expirés.

Exemple :

```go
manager := secrets.NewRotationManager(
    30 * 24 * time.Hour,
)
```

Politique :

```text
30 jours
60 jours
90 jours
180 jours
```

---

# Vérification

```go
needsRotation :=
    manager.NeedsRotation(
        secret,
    )
```

---

# Exécution

```go
err := manager.Rotate(
    secret,
    newValue,
)
```

---

# Multi-Tenant Security

Chaque secret appartient à :

```go
TenantID
```

Isolation stricte :

```text
Tenant A
  └── ne voit jamais
      Tenant B
```

---

# Sécurité

Le module applique :

- AES-256-GCM
- UUID v4
- Versioning
- Validation stricte
- Isolation tenant
- Deep Copy
- Mutex RWLock
- Context Support

---

# Gestion des erreurs

Principales erreurs :

```go
ErrInvalidSecret
ErrSecretNotFound
ErrSecretAlreadyExists

ErrInvalidEncryptionKey
ErrEncryptionFailed
ErrDecryptionFailed

ErrRotationFailed

ErrAccessDenied
ErrUnauthorizedAccess

ErrStoreUnavailable
ErrVaultUnavailable
```

---

# Exemple complet

```go
ctx := context.Background()

key, _ := secrets.GenerateKey()

enc, _ := secrets.NewEncryptionService(
    key,
)

store := secrets.NewMemoryStore()

vault, _ := secrets.NewVault(
    store,
    enc,
)

secret := secrets.NewSecret(
    "tenant-1",
    "database_password",
    "super-secret-password",
    "password",
)

_ = vault.CreateSecret(
    ctx,
    secret,
)

loaded, _ := vault.GetSecret(
    ctx,
    "tenant-1",
    "database_password",
)

fmt.Println(
    loaded.Value,
)
```

---

# Status

```text
Module: Secrets
Version: v1.0.0
Status: Production Ready
Coverage: Unit Tested
Multi-Tenant: Yes
Encryption: AES-256-GCM
Rotation: Supported
Vault: Supported
```