# envlens

Static analyzer for `.env` files that detects missing keys across deployment targets.

## Installation

```bash
go install github.com/yourname/envlens@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envlens.git && cd envlens && go build ./...
```

## Usage

Point `envlens` at a base `.env` file and one or more target environment files to check for missing keys:

```bash
envlens --base .env.example --targets .env.production,.env.staging
```

Example output:

```
.env.production
  ✗ STRIPE_SECRET_KEY   (missing)
  ✗ SENTRY_DSN          (missing)

.env.staging
  ✓ All keys present

2 issue(s) found across 2 target(s).
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--base` | Base `.env` file to use as the source of truth | `.env.example` |
| `--targets` | Comma-separated list of target env files | — |
| `--strict` | Exit with non-zero status if any keys are missing | `false` |

### CI Integration

```bash
# Fail the pipeline if any keys are missing
envlens --base .env.example --targets .env.production --strict
```

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

## License

[MIT](LICENSE)