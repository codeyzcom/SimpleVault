# SimpleVault

SimpleVault is a minimal self-hosted vault for storing sensitive information such as notes, credentials, and small files.
It focuses on simplicity, strong encryption, and a clean web interface without external UI frameworks.

The project is designed as a single-user or small multi-user personal vault, not as a cloud password manager.

---

## Features

- Encrypted storage (AES-256-GCM)
- Password-derived master key (no key files)
- Record types:
    - Notes
    - Website credentials
    - Small files (certificates, keys, documents)
- Web interface (Fiber + HTML templates)
- Session-based authentication
- Single active session per user
- Automatic session timeout
- Backup and restore of encrypted vault
- No JavaScript frameworks, no CSS frameworks

---

## Configuration

SimpleVault is configured via command-line flags at application startup.  
No configuration files are required — all settings are passed explicitly, keeping the setup simple and transparent.

### Available options

| Flag | Type | Default     | Description |
|------|------|-------------|-------------|
| `-host` | `string` | `localhost` | HTTP server host |
| `-port` | `int` | `7879`      | HTTP server port |
| `-session_ttl` | `duration` | `15m`       | Session lifetime (e.g. `10s`, `5m`, `1h`) |
| `-data_store` | `string` | `data/`     |  path to store vaults


### Notes
- session_ttl uses Go time.Duration format.
- When the session expires, the user is automatically logged out.
- Configuration is parsed once at startup and remains immutable during runtime.