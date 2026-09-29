# Audit UX/UI et fonctionnel du dashboard

## Portée et niveau de preuve

Audit du frontend et des routes backend présentes dans le dépôt. Les observations UX et les appels proviennent de `web/admin/src`; les routes HTTP et opérations proviennent de `internal/api/router.go` et `internal/api/handlers`. La présence d'une route dans le code ne prouve pas qu'elle répond dans un environnement déployé.

Le dev server Vite a été démarré sur `127.0.0.1:5173`. Une capture `curl.exe -i http://127.0.0.1:5173/api/v1/users` a reçu `HTTP/1.1 200 OK`, `Content-Type: text/html`, `Content-Length: 598`; le corps est l'`index.html` Vite (`HostMK Admin Foundation`), pas une réponse API JSON. La requête n'a donc pas reproduit le 404 rapporté. Aucun backend n'est lancé sur `localhost:8080`; démarrer `cmd/server` exécute des migrations PostgreSQL et n'a pas été fait. Le statut et le corps du 404 réel ainsi que son handler restent inconnus.

## Matrice

Dans « Fonctionne », « route au code » signifie que le frontend référence un chemin enregistré dans le routeur backend. Le succès HTTP et les données de l'environnement restent à vérifier avec une capture réseau authentique.

| Page | Endpoint(s) | Fonctionne | Champs affichés | Champs manquants | Actions manquantes |
| --- | --- | --- | --- | --- | --- |
| Overview | `GET /api/v1/users`, `/api/v1/sessions`, `/api/v1/api-keys`, `/api/v1/secrets`, `/api/v1/tenants`, `/api/v1/fraud/violations`, `/api/v1/audit/events` | Sept routes présentes au code. La page dérive des compteurs côté navigateur; aucune route agrégée Overview n'existe. Réponse runtime non mesurée. | Compteurs dérivés : utilisateurs actifs, sessions non révoquées, clés non expirées/révoquées, secrets non révoqués, tenants par statut, violations par niveau, événements par période. | Pas de champs détaillés : agrégats uniquement. Les totaux reflètent uniquement les enregistrements que les endpoints de liste renvoient (pas de pagination UI/API constatée). | Aucun bouton sur les cartes. Aucun lien vers les listes; aucun rafraîchissement manuel. |
| Tenants | `GET /api/v1/tenants`; backend expose aussi `POST /api/v1/tenants` et `/api/v1/tenants/action` | Route au code. Recherche et pagination sont des opérations locales après chargement; elles ne sont pas envoyées au serveur. | Nom, slug, statut, plan, branding, quotas. | ID comme colonne; domaine, propriétaire, fuseau horaire, pays, métadonnées, créé/modifié, suspension, expiration. | Créer, modifier, activer, suspendre. « Previous/Next » agit localement et n'est pas une pagination serveur. |
| Users | `GET /api/v1/users`; backend expose `POST /api/v1/users` et `/api/v1/users/action` | Routes au code; le frontend ne lit que la liste. | Nom, email, statut, MFA, dernière connexion, rôles, date de mise à jour. | ID comme colonne, tenant, date de création, date de vérification email. Le hash de mot de passe est volontairement exclu du JSON backend. | Créer un utilisateur, modifier son état/profil via UserAction, attribuer/mettre à jour les rôles. |
| Sessions | `GET /api/v1/sessions`; révocation via `POST /api/v1/sessions?action=revoke&id=…` | Routes au code; le bouton Révoquer appelle réellement l'API et invalide la query. | ID utilisateur, adresse IP, agent utilisateur, expiration, révocation. | ID de session comme colonne, tenant, date de création. | Aucune mutation supplémentaire exposée par le handler de sessions. Le bouton est inopérant si l'enregistrement n'a pas d'ID. |
| RBAC | `GET /api/v1/rbac/roles`, `GET /api/v1/permissions`; backend expose CRUD rôles, création/suppression de permissions et `POST /api/v1/permissions/check` | Routes au code. RBAC et Permissions sont le même écran `/rbac`; l'entrée Permissions de la navigation est un alias, pas une page distincte. | Rôles : nom, description, permissions. Permissions : nom, ressource. | Pour rôle : ID, tenant, dates. Pour permission : ID, tenant, action, description, dates. | Créer/modifier/supprimer rôles; créer/supprimer permissions; vérifier l'autorisation d'un utilisateur. |
| Permissions | Même endpoint `/api/v1/permissions`, affiché dans `RbacPage` | Pas de route frontend distincte; l'alias navigation fonctionne vers `/rbac`. Endpoint présent au code. | Nom, ressource. | ID, action, description, tenant et dates. | Toutes les actions permission disponibles au backend : création, suppression et vérification. |
| API Keys | `GET /api/v1/api-keys`; backend expose création `POST`, rotation `POST?action=rotate` et révocation `DELETE` sur la même route | Route au code; écran lecture seule. | Nom, préfixe, scopes, limite de débit, dernière utilisation, expiration. | ID comme colonne, tenant, statut révoqué, création/modification. La clé en clair, retournée uniquement à la création, n'est pas proposée car elle ne doit être affichée qu'une fois. | Créer, afficher une fois le secret de création, faire tourner, révoquer. |
| Secrets | `GET /api/v1/secrets`; backend expose création `POST`, rotation `POST?action=rotate`, historique `GET?id=…&history=true`, révocation `DELETE` | Route au code; écran lecture seule. | Nom, type, description, création, indicateur révoqué. | ID comme colonne, tenant, valeur masquée, mise à jour. Ne jamais afficher la valeur secrète persistée en clair. | Créer, faire tourner, consulter l'historique masqué, révoquer. |
| Fraud | `GET /api/v1/fraud/rules`, `GET /api/v1/fraud/violations`; backend expose CRUD règles et `POST /api/v1/fraud/analyze` | Routes au code; lecture seule. Les deux listes sont présentées en cartes, sans recherche ni filtre. | Règles : nom, champ, opérateur, seuil, activée, action, score. Violations : score, niveau, date. | Règle : ID, tenant, fenêtre, métadonnées, dates. Violation : ID, tenant, utilisateur, signal, règle, motif, statut, métadonnées, mise à jour. | Créer/modifier/supprimer une règle; lancer une analyse. Aucun endpoint de mutation de violation n'est déclaré. |
| Audit | `GET /api/v1/audit/events`; backend expose aussi `POST` | Route au code. Filtres tenant/sévérité exécutés localement; le service ne transmet pas de paramètres. Le handler accepte `limit`, mais le frontend ne l'envoie pas. | Action, ressource, sévérité, statut, date. | ID comme colonne, tenant, utilisateur, IP, agent utilisateur, message, métadonnées, date de mise à jour. | Export non exposé dans cette page; création d'événement backend non exposée. Pas de pagination ni de recherche serveur. |
| Monitoring | `GET /health`, `GET /api/v1/version` | Les deux routes sont enregistrées. Réponse runtime non vérifiée. L'écran transforme les composants de santé et les informations de version en cartes. | Nom du composant/service, statut, message/détail, nom/version/environnement du service. | Compteurs/mesures opérationnels non fournis par ces deux handlers. | Aucun acquittement, rafraîchissement manuel ou action opérationnelle. « Live » et « System online » dans le layout sont des libellés statiques, pas le résultat du health check. |
| Settings | `GET /api/v1/settings`; backend expose `POST` et `PUT` | Route au code; écran lecture seule. | Clé, valeur. | ID, tenant, date de mise à jour. | Ajouter ou modifier un réglage via les méthodes backend existantes. |

## Constats transverses

- Composants affichés : pages déclarées dans `web/admin/src/router/index.tsx`, montées dans `DashboardLayout`; la navigation est dans `DashboardLayout.tsx`. Le composant `Sidebar.tsx` distinct n'est pas utilisé par ce layout.
- Aucun mock métier n'a été trouvé dans les pages inspectées. Les listes proviennent des services API, et les compteurs Overview sont calculés sur les réponses de ces services.
- Les boutons visibles ont une action : recherche/pagination locale dans Tenants et filtres locaux dans Audit; révocation réelle dans Sessions. Les autres pages sont principalement des vues en lecture seule, sans boutons CRUD inertes identifiés.
- États de chargement, d'erreur et vides existent sur les pages inspectées. Les erreurs affichent surtout le message; API ajoute un statut et le payload à l'objet erreur, mais ces détails ne sont pas rendus systématiquement.
- Le fil d'Ariane et le titre « Overview » du layout restent identiques quelle que soit la route. La notification « Notification center is available… » et les statuts « Live »/« System online » sont statiques et peuvent donner une impression de données live non vérifiée.
- La route frontend Permissions n'est pas cassée au sens strict : elle mène à `/rbac`; elle duplique cependant la destination et ne permet pas d'ouvrir une page Permissions autonome. Aucun wildcard/not-found explicite n'est déclaré pour les routes inconnues.
- Pas de pagination serveur constatée. Tenants découpe localement après avoir chargé la liste. Audit envoie une requête de liste sans `limit`, bien que le handler backend le supporte. Les autres listes n'ont ni pagination ni recherche.

## Investigation du 404 Overview

### Appel frontend établi par le code

`OverviewPage` appelle `fetchOverviewCounts` avec React Query. `fetchOverviewCounts` déclenche en parallèle, via `Promise.all`, les services suivants :

| Requête | Service frontend | Handler / route backend enregistrée |
| --- | --- | --- |
| `GET /api/v1/users` | `services/users.ts` → `fetchUsers` | `API.Users` → `/api/v1/users` |
| `GET /api/v1/sessions` | `services/sessions.ts` → `fetchSessions` | `API.Sessions` → `/api/v1/sessions` |
| `GET /api/v1/api-keys` | `services/apiKeys.ts` → `fetchApiKeys` | `API.APIKeys` → `/api/v1/api-keys` |
| `GET /api/v1/secrets` | `services/secrets.ts` → `fetchSecrets` | `API.Secrets` → `/api/v1/secrets` |
| `GET /api/v1/tenants` | `services/tenants.ts` → `fetchTenants` | `API.Tenants` → `/api/v1/tenants` |
| `GET /api/v1/fraud/violations` | `services/fraud.ts` → `fetchFraudViolations` | `API.FraudViolations` → `/api/v1/fraud/violations` |
| `GET /api/v1/audit/events` | `services/audit.ts` → `fetchAuditEvents` | `API.AuditEvents` → `/api/v1/audit/events` |

Le client API compose chaque URL comme `VITE_API_BASE_URL` (vide par défaut) + le chemin ci-dessus. Quand cette variable est vide, l'URL de la première requête est l'origine courante + `/api/v1/users`. La capture locale Vite ci-dessus prouve qu'en développement sans proxy, ce chemin peut recevoir l'HTML de l'application plutôt que les données backend. Comme le client renvoie le texte des réponses non-JSON et `normalizeList` transforme une chaîne en liste vide, ce cas peut être affiché comme des compteurs à zéro plutôt qu'une erreur.

Pour toute réponse non-2xx, le client lève une erreur. `fetchOverviewCounts` ne capture pas individuellement les rejets : avant le correctif, le panneau ne permettait donc pas d'identifier lequel des sept GET avait échoué. Le client enrichit maintenant l'erreur avec méthode, URL, statut, texte de statut et corps; `OverviewPage` les affiche dans l'état d'erreur. Cela rendra la prochaine reproduction attribuable sans ajouter d'endpoint ni modifier le backend.

### Ce qui n'est pas prouvé

- Il n'y a pas de `GET /api/v1/overview` dans l'appel actuel, ni de handler Overview enregistré dans le routeur backend inspecté.
- Les sept chemins ci-dessus sont enregistrés dans le code backend présent. Cela ne prouve pas que le backend déployé ou ciblé par `VITE_API_BASE_URL` est la même version. Le résultat Vite pour `/api/v1/users` est **200 HTML**, non un 404; il ne peut donc pas être désigné comme endpoint responsable du message rapporté.
- La capture disponible ne contient aucune réponse du backend ciblé par l'application quand l'erreur 404 survient. Le statut, corps exact et handler effectivement responsable du 404 restent **indéterminés**. Attribuer le 404 à l'un des endpoints sans cette capture serait une hypothèse.
- Pour conclure avec les preuves demandées, il faut partager une capture Network de l'onglet où l'erreur apparaît (URL complète, méthode, statut, corps, sans bearer token) et l'URL de base API utilisée. Si le serveur est local, il faut le démarrer dans l'environnement déjà configuré puis reproduire l'erreur; aucune migration ni modification backend n'est nécessaire pour observer le navigateur si ce serveur est déjà lancé.

## Références de code

- Frontend Overview : `web/admin/src/pages/overview/OverviewPage.tsx`, `web/admin/src/services/overview.ts`, `web/admin/src/services/api.ts`.
- Navigation et routes : `web/admin/src/layouts/DashboardLayout.tsx`, `web/admin/src/router/index.tsx`.
- Enregistrement backend : `internal/api/router.go`.
- Handlers de ressources : `internal/api/handlers/api.go`, `internal/api/handlers/sessions.go`.
- Champs JSON affichables : `internal/models/*.go`.