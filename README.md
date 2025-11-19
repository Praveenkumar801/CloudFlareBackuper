# CloudFlare Backuper

A clean and efficient automated backup system that archives folders, uploads them to CloudFlare R2 storage, and sends Discord or Telegram notifications with download links.

## Features

- 🗜️ **Archive Multiple Folders**: Combines multiple directories into a single compressed tar.gz archive
- ☁️ **CloudFlare R2 Upload**: Automatically uploads backups to CloudFlare R2 storage
- 📢 **Multiple Notification Methods**: Supports Discord webhooks and Telegram bot notifications
- ⏰ **Flexible Scheduling**: Configure backup intervals using cron syntax
- 🔄 **Automatic Execution**: Runs in the background as a service
- 🎯 **Manual Backup**: Option to run a single backup on demand

## Installation

### Prerequisites

- Go 1.19 or higher
- CloudFlare R2 storage account with credentials
- Discord webhook URL and/or Telegram bot (at least one notification method required)

### Build from Source

```bash
git clone https://github.com/IndrajeethY/CloudFlareBackuper.git
cd CloudFlareBackuper
go build -o cloudflare-backuper
```

## Configuration

1. Copy the example configuration file:

```bash
cp config.example.yml config.yml
```

2. Edit `config.yml` with your settings:

```yaml
# CloudFlare R2 Storage Configuration
cloudflare:
  uri: "https://your_domain.com"
  bucket: "your_bucket_name"
  access_key_id: "your_access_key_id_here"
  secret_key: "your_secret_access_key_here"
  account_id: "your_account_id_here"


# Discord Webhook Configuration (optional)
discord:
  webhook_url: "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_WEBHOOK_TOKEN"

# Telegram Bot Configuration (optional)
# At least one notification method (Discord or Telegram) must be configured
telegram:
  bot_token: "YOUR_BOT_TOKEN"
  chat_id: "YOUR_CHAT_ID"

# Backup Configuration
backup:
  # Cron schedule format: "minute hour day month weekday"
  schedule: "0 */6 * * *"  # Every 6 hours
  
  # Folders to backup (will be combined into one archive)
  folders:
    - "/path/to/folder1"
    - "/path/to/folder2"
  
  # Backup filename prefix
  name_prefix: "backup"
  
  # Number of backups to keep (0 = keep all backups)
  # When set to 5, only the last 5 backups will be kept
  # Older backups are automatically deleted with notification
  retention_limit: 5
```

### Cron Schedule Examples

- `"0 */6 * * *"` - Every 6 hours
- `"0 0 * * *"` - Daily at midnight
- `"0 2 * * 0"` - Weekly on Sunday at 2 AM
- `"0 3 * * 1"` - Weekly on Monday at 3 AM
- `"*/30 * * * *"` - Every 30 minutes

## Usage

### Run as a Service (Continuous Backup)

```bash
# Run with default config file (config.yml)
./cloudflare-backuper

# Run with custom config file
./cloudflare-backuper -config /path/to/config.yml
```

The application will:
1. Run an initial backup immediately on startup
2. Schedule future backups based on the cron schedule
3. Continue running until stopped with Ctrl+C

### Run a Single Backup

```bash
# Run one backup and exit
./cloudflare-backuper -once

# Run one backup with custom config
./cloudflare-backuper -config /path/to/config.yml -once
```

### Run as a System Service

#### Using systemd (Linux)

Create a service file at `/etc/systemd/system/cloudflare-backuper.service`:

```ini
[Unit]
Description=CloudFlare Backuper Service
After=network.target

[Service]
Type=simple
User=your-username
WorkingDirectory=/path/to/CloudFlareBackuper
ExecStart=/path/to/CloudFlareBackuper/cloudflare-backuper
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

Enable and start the service:

```bash
sudo systemctl daemon-reload
sudo systemctl enable cloudflare-backuper
sudo systemctl start cloudflare-backuper
sudo systemctl status cloudflare-backuper
```

## Notifications

The application supports multiple notification methods. You can use Discord, Telegram, or both simultaneously. At least one notification method must be configured.

### Discord Notifications

The application sends rich embed notifications to Discord with:

#### Success Notification
- ✅ Green embed with "Backup Successful" title
- File name
- File size (human-readable format)
- Download link
- Timestamp

#### Failure Notification
- ❌ Red embed with "Backup Failed" title
- Error details
- Timestamp

#### Deletion Notification
- 🗑️ Yellow embed with "Old Backup Deleted" title
- Deleted file name
- Previous download link
- Timestamp

### Telegram Notifications

To set up Telegram notifications:

1. **Create a Telegram Bot**:
   - Talk to [@BotFather](https://t.me/botfather) on Telegram
   - Send `/newbot` command
   - Follow the instructions to create your bot
   - Copy the bot token provided

2. **Get Your Chat ID**:
   - Start a chat with your bot
   - Send any message to your bot
   - Visit: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
   - Find your `chat_id` in the JSON response

3. **Configure in config.yml**:
   ```yaml
   telegram:
     bot_token: "1234567890:ABCdefGHIjklMNOpqrsTUVwxyz"
     chat_id: "123456789"
   ```

#### Telegram Message Format
The application sends formatted messages to Telegram with:

- ✅ **Backup Successful**: File name, size, and download link
- ❌ **Backup Failed**: Error details
- 🗑️ **Old Backup Deleted**: Deleted file information

## Security Notes

- Never commit your `config.yml` file with real credentials
- Keep your CloudFlare secret key secure
- Restrict Discord webhook URL access
- Keep your Telegram bot token private
- Ensure proper file permissions on the config file (chmod 600)
- Consider using environment variables for sensitive credentials

## Troubleshooting

### "cloudflare.secret_key is required"

Make sure you've added the secret key to your `config.yml`:
```yaml
cloudflare:
  secret_key: "your_secret_access_key_here"
```

### "Failed to read config file"

Ensure `config.yml` exists in the current directory or specify the path with `-config`.

### "Failed to stat folder"

Verify that all folders listed in the config exist and are readable.

### Logs

The application logs all operations to stdout. To save logs to a file:

```bash
./cloudflare-backuper 2>&1 | tee backup.log
```

## Development

### Project Structure

```
CloudFlareBackuper/
├── .github/
│   └── workflows/   # GitHub Actions CI/CD workflows
├── backup/          # Archive creation logic
├── config/          # Configuration parsing
├── notification/    # Discord and Telegram notification integration
├── scheduler/       # Cron scheduling and backup orchestration
├── storage/         # CloudFlare R2 client
├── main.go          # Application entry point
├── config.example.yml
└── README.md
```

### Building

```bash
go build -o cloudflare-backuper
```

For optimized builds with smaller binaries:

```bash
go build -ldflags="-s -w" -trimpath -o cloudflare-backuper
```

### Testing

```bash
go test ./...
```

Run tests with race detector:

```bash
go test -race ./...
```

### CI/CD

The project uses GitHub Actions for continuous integration and automated releases:

- **CI Workflow**: Runs on every push and pull request
  - Builds on Linux, macOS, and Windows
  - Runs tests with race detector
  - Runs linting checks
  
- **Release Workflow**: Runs on tag push (e.g., `v1.0.0`)
  - Builds binaries for multiple platforms:
    - Linux (amd64, arm64)
    - macOS (amd64, arm64)
    - Windows (amd64, arm64)
  - Generates SHA256 checksums for all binaries
  - Creates GitHub release with all assets

To create a new release:

```bash
git tag v1.0.0
git push origin v1.0.0
```

### Memory Optimization

The application is optimized for efficient memory usage:

- **Streaming Uploads**: Files are streamed directly to R2 storage without loading entirely into memory
- **Streaming Archive Creation**: Large files are processed using streaming I/O operations
- **Resource Management**: Files are closed immediately after use to prevent descriptor leaks

This allows the application to handle large backups efficiently without excessive memory consumption.

## License

See [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.