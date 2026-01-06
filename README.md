# envsecrets

A lightweight CLI tool for securely managing environment variables and secrets.

envsecrets encrypts your environment variables using AES-GCM encryption with Argon2 key derivation. Your secrets are protected with a passphrase that can be stored securely in your system's keyring.

## Installation

```bash
go build -o envsecrets
```

## Commands

### init - Initialize a new vault

Create a new encrypted vault for an environment.

```bash
# Initialize a vault
envsecrets init --env prod

# Using short flag
envsecrets init -e staging
```

**What it does:**
- Prompts for a passphrase
- Creates `.envsecrets/{env}.vault` file
- Generates a salt and fingerprint
- Sets secure file permissions (0o600)

---

### add - Add or update entries

Add a new secret or update an existing one in a vault.

```bash
# Add with flags
envsecrets add --env prod --key API_KEY --value secret123

# Interactive mode
envsecrets add --env prod

# Hide value input (for sensitive data)
envsecrets add --env prod --secret
```

**Flags:**
- `--env, -e` - Environment name (required)
- `--key, -k` - Entry key (prompts if not provided)
- `--value, -v` - Entry value (prompts if not provided)
- `--secret, -s` - Hide value input in terminal

**What it does:**
- Opens the vault with passphrase
- Encrypts the value with AES-GCM
- Stores encrypted entry with timestamps
- Updates existing entries automatically

---

### rotate - Rotate passphrase

Change the passphrase for a vault by re-encrypting all entries.

```bash
envsecrets rotate --env prod
```

**What it does:**
- Opens vault with current passphrase
- Prompts for new passphrase (with confirmation)
- Generates new salt
- Re-encrypts all entries with new passphrase
- Updates keyring cache
- Preserves all timestamps

**Use case:** Periodic security rotation or when passphrase is compromised.

---

## Passphrase Management

Passphrases are retrieved in this order:

1. **Environment variable**: `ENVSECRET_PASSPHRASE`
2. **System keyring**: Cached from previous use
3. **Interactive prompt**: Asks user for input

The passphrase is cached in your system's keyring after first use for convenience.

## Security Features

- **Encryption**: AES-256-GCM (authenticated encryption)
- **Key Derivation**: Argon2id (memory-hard, GPU-resistant)
- **Authentication**: GCM provides authenticity verification
- **File Permissions**: Only owner can read/write vault files (0o600)
- **Memory Safety**: Sensitive data cleared after use
- **Keyring Integration**: Secure passphrase caching

## Vault Structure

Vaults are stored in `.envsecrets/{env}.vault`:

```json
{
  "meta": {
    "env": "production",
    "salt": "base64-encoded-salt",
    "fingerprint": "bcrypt-hash-of-passphrase"
  },
  "entries": {
    "API_KEY": {
      "value": "base64-encrypted-value",
      "created_at": "2025-01-05T10:00:00Z",
      "updated_at": "2025-01-05T10:00:00Z"
    }
  }
}
```

## Examples

```bash
# Create a production vault
envsecrets init --env production

# Add some secrets
envsecrets add --env production --key DATABASE_URL --value "postgres://..."
envsecrets add --env production --key API_KEY --value "sk-..." --secret

# Rotate passphrase periodically
envsecrets rotate --env production
```

## CI/CD Usage

Set the `ENVSECRET_PASSPHRASE` environment variable in your CI/CD pipeline:

```bash
export ENVSECRET_PASSPHRASE="your-passphrase"
envsecrets add --env prod --key KEY --value value
```
