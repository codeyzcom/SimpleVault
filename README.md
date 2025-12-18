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

## Architecture Overview

### High-level design

