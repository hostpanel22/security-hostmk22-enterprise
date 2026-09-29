# Rapport d’avancement vérifiable

État inspecté le 2026-09-29. Ce rapport précède toute implémentation du SDK Go. Les statuts décrivent le code présent, pas une validation manuelle navigateur/production.

## Résumé chiffré

| Domaine | Avancement mesuré | Dénominateur / preuve |
| --- | ---: | --- |
| Backend | 100 % des chemins API du dashboard enregistrés | 17/17 chemins référencés par les services/pages présents dans `internal/api/router.go`; ce chiffre mesure l’enregistrement, pas l’exhaustivité globale du backend. |
| Frontend | 67,7 % des opérations disponibles reliées à l’UI | 21/31 opérations mutantes/opérationnelles du backend recensées dans la matrice ci-dessous. Le compte exclut les GET de lecture et le rafraîchissement React Query. |
| Dashboard | 0 % Production Ready | 0/12 pages satisfont simultanément l’exposition de leurs fonctions backend, les états attendus et une validation UI/API authentifiée. Toutes ont des appels API; cela ne suffit pas à valider la page en production. |
| SDK Go | 0 % implémenté | 0 fichier source `.go`; `sdk/go/go.mod` et `sdk/go/README.md` sont vides. Les sous-répertoires `client`, `auth`, `audit`, `fraud`, `permissions`, `tenancy` sont vides. |
| Tests | 1,6 % couverture statements Go agrégée; 0 fichier de test frontend | Sortie `go test ./... -cover` + `go tool cover -func`; `FRONTEND_TEST_FILES=0`. 35 fichiers `_test.go` trouvés sous le dépôt par le décompte PowerShell (voir note ci-dessous). |
| Déploiement | N/D | Aucun seuil/cible de déploiement n’est défini. La recherche de Dockerfile, Compose, manifests PaaS et workflow Go n’a retourné aucun candidat; cela ne permet pas de mesurer la disponibilité d’un déploiement externe. |

### Calcul frontend 21/31

Opérations distinctes comptées : tenants 4 (create/update/activate/suspend); users 7 (create/update/activate/suspend/verify-email/enable-MFA/delete); sessions 1 (revoke); rôles 3 (create/update/delete); permissions 3 (create/delete/check); API keys 3 (create/rotate/revoke); secrets 4 (create/rotate/history/revoke); fraud 4 (create/update/delete/analyze); audit 1 (create); settings 1 (save). Total 31. Celles reliées à l’UI actuelle : users 7, sessions 1, RBAC et permissions 6, API keys 3, secrets 4. Total 21. Les méthodes ajoutées dans `services/audit.ts` et `services/fraud.ts` ne comptent pas comme actions UI car leurs pages ne les appellent pas.

## Pages du dashboard

`GET` désigne le chargement observé dans le code frontend; la présence du chemin est confirmée dans le routeur backend. Les listes backend ne proposent pas de pagination universelle; la pagination locale ne prétend pas être serveur.

| Page | Endpoints utilisés | Champs retournés; affichés; manquants | Actions backend; actions UI; absentes de l’UI | Squelettes / TODO | Statut |
| --- | --- | --- | --- | --- | --- |
| Overview | `GET /users`, `/sessions`, `/api-keys`, `/secrets`, `/tenants`, `/fraud/violations`, `/audit/events` | L’agrégateur reçoit les objets User, Session, Key, Secret, Tenant, FraudViolation et AuditEvent; affiche 14 compteurs dérivés. Champs détaillés non affichés car les données sont réduites en compteurs. Les listes Audit/Fraud sont plafonnées (100 au service Overview; violations 100 au handler), donc pas de total garanti au-delà. | Le backend n’a pas de route d’agrégat Overview. UI: rafraîchissement manuel; erreur montre URL/statut/body si `StandardApiError`. Aucune action métier ou navigation vers détail. | État vide/chargement/erreur présents. Dépend d’un `Promise.all` de sept GET; un rejet bloque tout l’Overview. | Partiel |
| Tenants | `GET/POST /api/v1/tenants`; `POST /api/v1/tenants/action` pour update/suspend/activate | Tenant retourne `id,name,slug,domain,status,owner_id,plan,timezone,country,quotas,branding,metadata,created_at,updated_at,suspended_at,expires_at`. Affichés: name, slug, status, plan, branding, quotas. Manquants: id en colonne, domain, owner, timezone, country, metadata, dates/suspension/expiration. | Backend: create/update/activate/suspend. UI: GET, recherche et découpage local uniquement. Toutes les écritures backend manquent. | Pas de formulaire, ni filtre de statut, ni badge sémantique (statut est affiché en texte stylé); « Next » basé sur longueur de page, pas sur total serveur. | Partiel |
| Users | `GET/POST /api/v1/users`; `POST /api/v1/users/action` | User retourne `id,tenant_id,email,name,status,email_verified_at,mfa_enabled,role_ids,created_at,updated_at,last_login`; tous affichés. Hash/secret MFA exclus par le backend. | Backend: create/update/suspend/activate/verify-email/delete/enable-MFA. UI: ces 7 actions sont reliées avec confirmation sur actions sensibles. | Recherche/pagination locales. Formulaire de création ne permet que `active/pending`, alors que le modèle comprend aussi `inactive/suspended/blocked`. Erreur de chargement des rôles seulement visible quand création ouverte. Pas de refresh token automatique. | Partiel |
| Sessions | `GET/POST /api/v1/sessions?action=revoke&id=…` | Session retourne `id,tenant_id,user_id,user_agent,ip_address,expires_at,revoked_at,created_at`; tous affichés. | Backend: list/revoke. UI: revoke confirmé, recherche, filtre, pagination locale. | Pas de pagination serveur exposée par handler. | Partiel |
| RBAC | `GET/POST/PUT/DELETE /api/v1/rbac/roles`; `GET/POST/DELETE /api/v1/permissions`; `POST /api/v1/permissions/check` | Role `id,tenant_id,name,description,permissions,created_at,updated_at`; Permission `id,tenant_id,name,resource,action,description,created_at,updated_at`; ces champs sont affichés. | Backend: CRUD rôle, create/delete permission, check permission. UI: chacune de ces opérations, avec recherche/pagination locales et confirmations sur suppression. | Pas de tests UI/API authentifiés; liste non paginée côté backend. | Partiel |
| Permissions | Mêmes chemins que RBAC | Même modèle Permission et mêmes champs affichés dans `RbacPage`. | Create/delete/check reliés dans le panneau RBAC. | L’entrée de navigation Permissions pointe vers `/rbac`; pas de route ni page indépendante. | Partiel |
| API Keys | `GET/POST/DELETE /api/v1/api-keys`; rotation `POST?action=rotate` | Liste retourne `id,tenant_id,name,prefix,scopes,rate_limit,expires_at,revoked_at,last_used_at,created_at,updated_at`; tous affichés. Création/rotation peut retourner `key` en clair, affichée dans un panneau à usage ponctuel. | Backend: create/rotate/revoke; UI: les trois, avec confirmations de rotation/révocation. Aucun endpoint de modification générale n’existe. | Recherche/filtre/pagination locales; actions de mutation non testées avec credentials réels dans le navigateur. | Partiel |
| Secrets | `GET/POST/DELETE /api/v1/secrets`; rotation `POST?action=rotate`; historique `GET?id=…&history=true` | Liste: `id,tenant_id,name,type,masked_value,description,revoked_at,created_at,updated_at`; champs affichés. Historique: `id,secret_id,masked_value,created_at,created_by`; champs affichés. Pas de plaintext retourné par liste. | Backend: create/rotate/history/revoke; UI: les quatre, avec confirmation de révocation. | Recherche/filtre/pagination locales. Aucun endpoint de modification générale n’existe. Opérations sensibles non validées e2e. | Partiel |
| Fraud | `GET/POST/PUT/DELETE /api/v1/fraud/rules`; `GET /api/v1/fraud/violations`; `POST /api/v1/fraud/analyze` | Rule retourne `id,tenant_id,name,enabled,field,operator,threshold,window_seconds,action,score,metadata,created_at,updated_at`; l’écran n’affiche que name/condition/enabled/action/score. Violation retourne `id,tenant_id,user_id,signal_id,rule_id,score,level,reason,status,metadata,created_at,updated_at`; l’écran n’affiche que score/level/created_at. | Backend: CRUD règles + analyse; UI actuelle: lecture seule. Les méthodes de mutation/analyze existent dans le service mais `FraudPage` ne les utilise pas. Pas de mutation violation au backend. | Cartes en lecture seule; pas de formulaire, recherche, filtre, pagination ni confirmation de règle. | Partiel |
| Audit | `GET/POST /api/v1/audit/events`; `GET` accepte `limit` (défaut 50, max 500) | Event retourne `id,tenant_id,user_id,action,resource,severity,status,ip_address,user_agent,message,metadata,created_at,updated_at`. Affichés: action/resource/severity/status/created_at; tenant sert au filtre. ID, user, IP, agent, message, metadata, updated manquent. | Backend: list/create. UI: GET avec limit 100 et filtres locaux tenant/severity. Le `createAuditEvent` existe dans le service mais n’est pas appelé dans `AuditPage`. | Pas de formulaire de création, recherche texte générale, pagination ni export par cette page. Erreur de chargement affichée sans statut/URL/body. | Partiel |
| Monitoring | `GET /health`, `GET /api/v1/version` | Health: status/service/version/env/components(name/status/message); Version: name/version/env. UI montre name/status/detail construits à partir de ces champs. | Les deux handlers sont en lecture seule; aucune action mutante backend. | Pas d’auto-refresh ou bouton refresh. Le libellé `System online` du layout est statique et non lié au health endpoint. | Partiel |
| Settings | `GET/POST/PUT /api/v1/settings` | Setting retourne `id,tenant_id,key,value,updated_at`. UI montre seulement key et value. | Backend: upsert POST/PUT; UI: lecture seule. L’action save manque. | Pas de formulaire, recherche, filtre ou pagination. Erreur seulement message. | Partiel |

## Placeholders, mocks, erreurs

- Fichiers source frontend vides (`web/admin/src`) : **0**, vérifié par `Get-ChildItem ... | Where-Object Length -eq 0`.
- Recherche VS Code des termes `mock|TODO|FIXME|coming soon|not implemented|sample data|demo data` dans `web/admin/src/**` : **0 résultat**. Les données métier des pages auditées passent par les services API; aucun tableau mock détecté.
- Placeholders visuels vérifiés dans `DashboardLayout.tsx`: titre/fil d’Ariane toujours Overview, état `System online` et notification « Notification center is available… » statiques. `Authenticated` est basé sur la présence d’un access token, pas une vérification backend.
- Réponses non-JSON: `parseResponse()` les transforme en texte; `normalizeList()` retourne `[]` pour un payload non reconnu. Ce chemin peut produire un état vide sans erreur. Les erreurs HTTP des pages autres que Overview n’affichent généralement pas statut/URL/body.
- `web/admin` ne contient aucun fichier `.test.*` ou `.spec.*` selon le décompte exécuté: tests frontend 0.
- SDK Go: le `docs/SDK-Go.md` du dépôt décrit des API non présentes dans le module `sdk/go`; le `sdk/go/README.md` et son `go.mod` sont vides au snapshot audité.

## Validation réelle

Commandes exécutées depuis `web/admin`:

```text
npm run build
✓ 114 modules transformed.
dist/index.html                   0.44 kB │ gzip:  0.30 kB
dist/assets/index-BJ44hHiU.css   21.16 kB │ gzip:  4.78 kB
dist/assets/index-DknSqrlY.js   331.96 kB │ gzip: 96.79 kB
✓ built in 44.55s
BUILD_EXIT=0

npx tsc --noEmit
Exit code: 0 (aucune sortie TypeScript)

npm run lint
ESLint: aucune erreur; LINT_EXIT=0
Warning: TypeScript 5.9.3 hors plage prise en charge par @typescript-eslint/typescript-estree (>=4.8.4 <5.8.0)
```

Tests à la racine:

```text
go test ./...
GO_TEST_EXIT=0
```

`go test ./... -cover` et `go tool cover -func=<coverprofile>`:

```text
internal/api/middleware coverage: 38.2% of statements
internal/engine coverage: 48.8% of statements
total: (statements) 1.6%
COVERAGE_RUN_EXIT=0
```

Le total inclut le périmètre Go ./... et les paquets sans tests; les suites `tests/audit`, `tests/fraud`, `tests/permissions`, `tests/rbac`, `tests/secrets`, `tests/tenancy` rapportent `[no statements]` car ce sont des packages de tests d’intégration/helpers, pas des packages de production instrumentés.

## Ordre de suite

Le rapport est maintenant établi. Les pages incomplètes à prioriser avant retour dashboard : Tenants, Fraud, Audit, Settings; ensuite validation authentifiée e2e de toutes les écritures. Le SDK Go peut être commencé après ce rapport, dans `sdk/go` existant uniquement, en remplaçant ses fichiers vides sans ajouter un deuxième module à côté.