# HostMK Security Java SDK

Typed Java 17+ client for the endpoints registered by HostMK Security Framework.

## Build and test

```bash
mvn clean test
```

## Bearer authentication

```java
HostMKClient client = HostMKClient.builder()
    .baseUrl("https://security.example.com")
    .token(System.getenv("HOSTMK_TOKEN"))
    .build();

ListResult<User> users = client.users().list();
```

## API key authentication

```java
HostMKClient client = HostMKClient.builder()
    .baseUrl("https://security.example.com")
    .apiKey(System.getenv("HOSTMK_API_KEY"))
    .build();

client.audit().list(100);
```

`login` and `refresh` store the returned access token in the shared client and clear the API key. JWTs are treated as opaque credentials; signing keys are never required by the SDK.

Available modules: authentication, users, tenants, sessions, RBAC roles, permissions, API keys, secrets, audit, fraud, settings, and monitoring.
