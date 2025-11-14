# Quick Start Guide

This guide will help you get CloudFlare Backuper up and running in just a few minutes.

## Prerequisites

- CloudFlare R2 storage account
- CloudFlare R2 API credentials (Account ID, Access Key ID, Secret Access Key)
- Discord webhook URL
- Go 1.19+ (if building from source) OR Docker (for containerized deployment)

## Getting Your CloudFlare Credentials

1. Log in to your CloudFlare dashboard
2. Go to R2 storage
3. Create a bucket if you don't have one
4. Go to "Manage R2 API Tokens"
5. Create a new API token with read/write permissions
6. Note down:
   - Account ID
   - Access Key ID (UID)
   - Secret Access Key

## Getting Discord Webhook URL

1. Open Discord and go to your server
2. Right-click on the channel where you want notifications
3. Select "Edit Channel" → "Integrations" → "Webhooks"
4. Click "New Webhook" or use an existing one
5. Copy the webhook URL

## Installation Methods

### Method 1: Direct Binary (Recommended for simple setups)

1. **Clone and build:**
   ```bash
   git clone https://github.com/IndrajeethY/CloudFlareBackuper.git
   cd CloudFlareBackuper
   make build
   ```

2. **Configure:**
   ```bash
   cp config.example.yml config.yml
   nano config.yml  # Edit with your settings (including secret_key)
   ```

3. **Run:**
   ```bash
   ./cloudflare-backuper
   ```

### Method 2: Docker (Recommended for production)

1. **Clone the repository:**
   ```bash
   git clone https://github.com/IndrajeethY/CloudFlareBackuper.git
   cd CloudFlareBackuper
   ```

2. **Create config file:**
   ```bash
   cp config.example.yml config.yml
   nano config.yml  # Edit with your settings (including secret_key)
   ```

3. **Edit docker-compose.yml:**
   - Update the volume paths to point to your folders
   
4. **Start the service:**
   ```bash
   docker-compose up -d
   ```

5. **Check logs:**
   ```bash
   docker-compose logs -f
   ```

### Method 3: System Service (Linux)

1. **Build and install:**
   ```bash
   git clone https://github.com/IndrajeethY/CloudFlareBackuper.git
   cd CloudFlareBackuper
   make build
   sudo cp cloudflare-backuper /usr/local/bin/
   sudo cp config.example.yml /etc/cloudflare-backuper/config.yml
   ```

2. **Edit config:**
   ```bash
   sudo nano /etc/cloudflare-backuper/config.yml
   ```

3. **Create systemd service:**
   ```bash
   sudo nano /etc/systemd/system/cloudflare-backuper.service
   ```
   
   Paste:
   ```ini
   [Unit]
   Description=CloudFlare Backuper Service
   After=network.target

   [Service]
   Type=simple
   User=root
   WorkingDirectory=/usr/local/bin
   ExecStart=/usr/local/bin/cloudflare-backuper -config /etc/cloudflare-backuper/config.yml
   Restart=always
   RestartSec=10

   [Install]
   WantedBy=multi-user.target
   ```

4. **Start service:**
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl enable cloudflare-backuper
   sudo systemctl start cloudflare-backuper
   sudo systemctl status cloudflare-backuper
   ```

## Configuration

Edit your `config.yml` file:

```yaml
cloudflare:
  uri: "https://your-public-url.com"  # Your R2 public URL
  bucket: "your-bucket-name"
  uid: "your_access_key_id"
  secret_key: "your_secret_access_key"
  account_id: "your_account_id"

discord:
  webhook_url: "https://discord.com/api/webhooks/YOUR_WEBHOOK_ID/YOUR_TOKEN"

backup:
  schedule: "0 */6 * * *"  # Every 6 hours
  folders:
    - "/path/to/important/data"
    - "/home/user/documents"
    - "/var/www/html"
  name_prefix: "backup"
  retention_days: 30
```

## Common Schedules

- **Every hour:** `"0 * * * *"`
- **Every 6 hours:** `"0 */6 * * *"`
- **Daily at 2 AM:** `"0 2 * * *"`
- **Weekly (Sunday at 3 AM):** `"0 3 * * 0"`
- **Monthly (1st day at midnight):** `"0 0 1 * *"`

## Testing Your Setup

Run a one-time backup to test:

```bash
# Direct binary
./cloudflare-backuper -once

# Docker
docker-compose run --rm cloudflare-backuper -once

# System service
sudo systemctl start cloudflare-backuper
sudo journalctl -u cloudflare-backuper -f
```

## Verification

After running, you should:

1. ✅ See the backup file in your CloudFlare R2 bucket
2. ✅ Receive a Discord notification with the download link
3. ✅ Be able to download and extract the backup

## Troubleshooting

### "cloudflare.secret_key is required"

Make sure you've added the secret key to your `config.yml`:
```yaml
cloudflare:
  secret_key: "your_secret_access_key"
```

### "Failed to stat folder"

Ensure all folders in your config exist and are readable:
```bash
ls -la /path/to/folder
```

### "Failed to upload to R2"

- Verify your CloudFlare credentials are correct
- Check that the bucket exists and you have write permissions
- Ensure your internet connection is working

### Discord notification not received

- Verify the webhook URL is correct
- Test the webhook manually:
  ```bash
  curl -X POST -H "Content-Type: application/json" \
    -d '{"content":"Test message"}' \
    YOUR_WEBHOOK_URL
  ```

### Check logs

- **Direct binary:** Output is in console
- **Docker:** `docker-compose logs -f`
- **System service:** `sudo journalctl -u cloudflare-backuper -f`

## Next Steps

- Set up monitoring for the backup process
- Configure backup retention policies
- Test restoration procedures
- Set up alerts for backup failures

## Need Help?

- Check the main [README.md](README.md) for detailed documentation
- Open an issue on GitHub
- Review the logs for error messages

## Security Reminder

⚠️ **Never commit your `config.yml` or `.env` files to version control!**

The repository includes these files in `.gitignore` to prevent accidental commits.
