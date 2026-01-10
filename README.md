# envsecrets

A lightweight CLI tool for securely managing environment variables and secrets.

envsecrets encrypts your environment variables using AES-GCM encryption with Argon2 key derivation. Your secrets are protected with a passphrase that can be stored securely in your system's keyring.

## Why envsecrets?

**Version control your secrets safely.** Unlike traditional approaches where `.env` files must be kept out of version control, envsecrets encrypts your secrets so they can be safely committed to your repository. This means:

- **Track secret changes over time** - See when secrets were added, modified, or removed through git history
- **Collaborate securely** - Share encrypted secrets with your team through version control
- **Environment parity** - Keep development, staging, and production secrets in the same repo
- **Disaster recovery** - Your encrypted secrets are backed up wherever your code is
- **Simplify deployment** - No need for separate secret management infrastructure

The `.envsecrets/*.vault` files are encrypted and safe to commit. Only those with the passphrase can decrypt them.

## Installation

### Prerequisites

- Go 1.22 or higher

### Install from GitHub

```bash
go install github.com/suvaidkhan/envsecrets@latest
```

This will download, build, and install the `envsecrets` binary to your `$GOPATH/bin` directory.

Make sure `$GOPATH/bin` is in your PATH:
```bash
export PATH=$PATH:$(go env GOPATH)/bin
```

### Build from Source

```bash
# Clone the repository
git clone https://github.com/suvaidkhan/envsecrets.git
cd envsecrets

# Build the binary
go build -o envsecrets

# Optionally, move to a directory in your PATH
sudo mv envsecrets /usr/local/bin/
```

### Verify Installation

```bash
envsecrets --version
```

## Commands

### Quick Reference

| Command | Description |
|---------|-------------|
| `init` | Initialize a new encrypted vault |
| `add` | Add or update a secret in a vault |
| `get` | Retrieve a specific secret |
| `delete` | Remove a secret from a vault |
| `export` | Export all secrets to dotenv or JSON |
| `import` | Import secrets from dotenv or JSON file |
| `rotate` | Change vault passphrase |
| `clear` | Clear cached passphrase from keyring |
| `destroy` | Permanently delete a vault |

---

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

### get - Retrieve a secret

Get a decrypted secret from a vault and print it to stdout.

```bash
envsecrets get --env prod --key API_KEY
```

**Flags:**
- `--env, -e` - Environment name (required)
- `--key, -k` - Secret key to retrieve (required)

**What it does:**
- Opens the vault with passphrase
- Retrieves and decrypts the specified secret
- Prints the value to stdout

**Use case:** Scripts, CI/CD pipelines, or exporting a single secret.

---

### delete - Delete a secret

Remove a secret from a vault.

```bash
# Delete with flags
envsecrets delete --env prod --key API_KEY

# Interactive mode
envsecrets delete --env prod
```

**Flags:**
- `--env, -e` - Environment name (required)
- `--key, -k` - Secret key to delete (prompts if not provided)

**What it does:**
- Opens the vault with passphrase
- Removes the specified entry
- Updates the vault file

---

### export - Export all secrets

Export all decrypted secrets from a vault in dotenv or JSON format.

```bash
# Export as dotenv format
envsecrets export --env prod > .env

# Export as JSON
envsecrets export --env staging --format json > env.json
```

**Flags:**
- `--env, -e` - Environment name (required)
- `--format` - Output format: `dotenv` or `json` (default: dotenv)

**What it does:**
- Opens the vault with passphrase
- Decrypts all entries
- Outputs to stdout in the specified format

**Use case:** Generate `.env` files for local development or deployment.

---

### import - Import secrets from file

Import secrets from a dotenv or JSON file into a vault.

```bash
# Import from dotenv file
envsecrets import .env --env prod --format dotenv

# Import from JSON file
envsecrets import config.json --env staging --format json

# Import from stdin
cat .env | envsecrets import --env local --format dotenv

# Overwrite existing keys
envsecrets import .env --env prod --format dotenv --overwrite
```

**Flags:**
- `--env, -e` - Environment name (required)
- `--format` - Input format: `dotenv` or `json` (required)
- `--overwrite` - Overwrite existing keys (default: false)

**What it does:**
- Parses the input file (dotenv or JSON format)
- Opens the vault with passphrase
- Encrypts and adds entries to the vault
- Skips existing keys unless `--overwrite` is used

**Use case:** Migrate existing `.env` files to encrypted storage.

---

### clear - Clear cached passphrase

Remove the cached passphrase for an environment from the system keyring.

```bash
envsecrets clear --env prod
```

**Flags:**
- `--env, -e` - Environment name (required)

**What it does:**
- Removes the cached passphrase from system keyring
- Next command will prompt for passphrase again

**Use case:** Security best practice, logout, or switching users.

---

### destroy - Destroy a vault

Permanently delete a vault and all its secrets.

```bash
envsecrets destroy --env prod
```

**Flags:**
- `--env, -e` - Environment name (required)

**What it does:**
- Clears cached passphrase
- Prompts for passphrase to verify authorization
- Asks for confirmation
- Permanently deletes the vault file

**Warning:** This action cannot be undone. All secrets will be lost unless backed up.

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

### Basic Workflow

```bash
# Create a production vault
envsecrets init --env production

# Add some secrets
envsecrets add --env production --key DATABASE_URL --value "postgres://..."
envsecrets add --env production --key API_KEY --value "sk-..." --secret

# Retrieve a secret
envsecrets get --env production --key API_KEY

# Delete a secret
envsecrets delete --env production --key OLD_KEY
```

### Migration from .env Files

```bash
# Import existing .env file
envsecrets init --env development
envsecrets import .env --env development --format dotenv

# Commit the encrypted vault
git add .envsecrets/development.vault
git commit -m "Add encrypted development secrets"
```

### Exporting Secrets

```bash
# Generate .env file for local development
envsecrets export --env development > .env

# Export as JSON for configuration
envsecrets export --env production --format json > config.json
```

### Security Maintenance

```bash
# Rotate passphrase periodically
envsecrets rotate --env production

# Clear cached passphrase when done
envsecrets clear --env production

# Destroy a vault (DANGEROUS)
envsecrets destroy --env old-env
```

## CI/CD Usage

Set the `ENVSECRET_PASSPHRASE` environment variable in your CI/CD pipeline:

```bash
# Set passphrase in CI environment
export ENVSECRET_PASSPHRASE="your-passphrase"

# Export secrets to .env for your application
envsecrets export --env prod > .env

# Or retrieve individual secrets
API_KEY=$(envsecrets get --env prod --key API_KEY)
DATABASE_URL=$(envsecrets get --env prod --key DATABASE_URL)
```

**Best practices:**
- Store the passphrase as a secret in your CI/CD platform (GitHub Secrets, GitLab CI Variables, etc.)
- Use different passphrases for each environment
- The encrypted `.envsecrets/*.vault` files can be safely committed to your repository
