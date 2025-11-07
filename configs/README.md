# Configuration Files

This directory contains configuration files for the Akeneo Migrator.

## Files

- `settings.local.json.example` - Example configuration file (committed to git)
- `settings.local.json` - Your actual configuration (ignored by git)

## Setup

1. Copy the example file:
```bash
cp settings.local.json.example settings.local.json
```

2. Edit `settings.local.json` with your credentials:
```json
{
  "akeneoSource": {
    "api": {
      "url": "https://your-source-akeneo.com",
      "credentials": {
        "clientId": "your_source_client_id",
        "secret": "your_source_secret",
        "username": "your_source_username",
        "password": "your_source_password"
      }
    }
  },
  "akeneoDest": {
    "api": {
      "url": "https://your-dest-akeneo.com",
      "credentials": {
        "clientId": "your_dest_client_id",
        "secret": "your_dest_secret",
        "username": "your_dest_username",
        "password": "your_dest_password"
      }
    }
  }
}
```

## Security

⚠️ **Important**: Never commit `settings.local.json` to git as it contains sensitive credentials.

The file is already included in `.gitignore` to prevent accidental commits.

## Environment Variables

You can also use environment variables by setting:
- `ENVIRONMENT=local` (default)
- `CONFIG_PATH=akeneo-migrator` (default)

The application will look for the config file at:
```
configs/akeneo-migrator/settings.local.json
```
