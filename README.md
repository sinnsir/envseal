# envseal

> Encrypt and version `.env` files using [age](https://age-encryption.org/) encryption, with per-environment key management and git-friendly diffs.

---

## Installation

```bash
go install github.com/yourusername/envseal@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/envseal/releases).

---

## Usage

### Initialize a new keychain

```bash
envseal init
```

### Seal (encrypt) an env file

```bash
envseal seal .env --env production
```

This produces a `.env.sealed` file safe to commit to your repository.

### Unseal (decrypt) an env file

```bash
envseal unseal .env.sealed --env production --out .env
```

### Rotate keys for an environment

```bash
envseal rotate --env staging
```

### Diff two sealed files

```bash
envseal diff .env.sealed .env.sealed.bak
```

Diffs are displayed in plaintext (keys only, values redacted) so changes remain reviewable in pull requests without exposing secrets.

---

## How It Works

- Each environment (`production`, `staging`, `development`, etc.) has its own [age](https://age-encryption.org/) keypair stored in `~/.config/envseal/keys/`.
- Sealed files are base64-encoded and structured for clean, line-level git diffs.
- No external services or cloud dependencies — everything runs locally.

---

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you'd like to change.

---

## License

[MIT](LICENSE)