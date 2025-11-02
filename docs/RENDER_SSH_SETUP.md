# Render SSH Setup Guide

This guide explains how to set up and use SSH access to your DKL Email Service deployment on Render.

## Prerequisites

Your Dockerfile has been configured to support SSH connections with the following features:
- OpenSSH server installed
- Root `.ssh` directory created with proper permissions (0700)
- Root account unlocked for SSH access
- Bash shell available

## Step 1: Generate an SSH Key Pair

If you don't already have an SSH key pair on your machine:

### On Windows (PowerShell or Git Bash):
```bash
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519
```

### On macOS/Linux:
```bash
ssh-keygen -t ed25519 -f ~/.ssh/id_ed25519
```

**Notes:**
- You'll be prompted to provide an optional passphrase (recommended for added security)
- This creates two files:
  - `~/.ssh/id_ed25519` (private key - keep this secret!)
  - `~/.ssh/id_ed25519.pub` (public key - this is what you'll share with Render)

### Alternative Key Types
Render also supports:
- `rsa`: `ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa`
- `ecdsa`: `ssh-keygen -t ecdsa -b 521 -f ~/.ssh/id_ecdsa`
- Hardware keys: `ed25519-sk`, `ecdsa-sk` (YubiKey, etc.)

## Step 2: Add Your Public Key to Render

1. Go to your [Render Account Settings](https://dashboard.render.com/account)

2. Find the **SSH Public Keys** section and click **+ Add SSH Public Key**

3. Copy your public key to your clipboard:

   **On Windows (PowerShell):**
   ```powershell
   Get-Content ~/.ssh/id_ed25519.pub | Set-Clipboard
   ```

   **On macOS:**
   ```bash
   pbcopy < ~/.ssh/id_ed25519.pub
   ```

   **On Linux:**
   ```bash
   cat ~/.ssh/id_ed25519.pub
   # Then manually copy the output
   ```

4. In the Render dialog:
   - **Name**: Enter a descriptive name (e.g., "Work Laptop" or "Personal Computer")
   - **Key**: Paste your public key (entire contents of `.pub` file)

5. Click **Add SSH Public Key**

## Step 3: Starting an SSH Session

### Method 1: Using Render CLI (Recommended)

1. Install the Render CLI if you haven't already:
   ```bash
   # macOS/Linux
   npm install -g render-cli
   
   # Windows
   npm install -g render-cli
   ```

2. Log in to Render:
   ```bash
   render login
   ```

3. Start an SSH session:
   ```bash
   render ssh
   ```
   
   This will show an interactive menu of your services. Use arrow keys to select `dklautomatie-backend` and press Enter.

4. Or connect directly using your service ID:
   ```bash
   render ssh srv-YOUR_SERVICE_ID
   ```

### Method 2: Direct SSH Command

1. Find your service ID in the Render Dashboard (format: `srv-xxxxx`)

2. Determine your Render region:
   - Oregon: `ssh.oregon.render.com`
   - Ohio: `ssh.ohio.render.com`
   - Virginia: `ssh.virginia.render.com`
   - Frankfurt: `ssh.frankfurt.render.com`
   - Singapore: `ssh.singapore.render.com`

3. Connect using SSH:
   ```bash
   ssh srv-YOUR_SERVICE_ID@ssh.YOUR_REGION.render.com
   ```

   Example for Oregon:
   ```bash
   ssh srv-abc123@ssh.oregon.render.com
   ```

### Connecting to a Specific Instance

By default, SSH connects to a random instance. To connect to a specific instance:

```bash
ssh srv-abc123-d4e5f@ssh.oregon.render.com
```

Where `d4e5f` is the 5-character instance slug visible in your service logs or metrics.

## Common SSH Commands Once Connected

### Check Application Status
```bash
# View running processes
ps aux | grep main

# Check disk usage
df -h

# View memory usage
free -m

# Check environment variables (be careful with sensitive data)
printenv | grep -E "DB_|SMTP_|APP_"
```

### View Logs
```bash
# If your app logs to stdout/stderr, you can check recent logs
# Note: Actual log files depend on your logger configuration
cd /app
ls -la

# Check if log files exist
find . -name "*.log"
```

### Database Connectivity Test
```bash
# You can test database connectivity if psql is installed
# Note: It's not installed by default in the Alpine image
```

### Application Health Check
```bash
# Test the health endpoint from within the container
wget -qO- http://localhost:8080/api/health
```

## Important Docker-Specific Notes

✅ **Your configuration includes:**
- OpenSSH server (`openssh-server`)
- Bash shell for better SSH experience
- Properly configured `/root/.ssh` directory with 0700 permissions
- Unlocked root account for SSH access

❌ **Limitations:**
- Cannot mount persistent disks to `$HOME` (not applicable for this project)
- Cannot have a locked root account (already configured)

## Troubleshooting

### "Permission denied (publickey)" Error

1. **Verify your SSH key is being used:**
   ```bash
   ssh -v srv-YOUR_SERVICE_ID@ssh.YOUR_REGION.render.com
   ```
   
   Look for lines like:
   ```
   debug1: Offering public key: /Users/YOUR_NAME/.ssh/id_ed25519
   ```

2. **Check which keys are in your SSH agent:**
   ```bash
   ssh-add -l
   ```
   
   This should output:
   ```
   256 SHA256:YOUR_KEY_FINGERPRINT YOUR_NAME@YOUR_HOST (ED25519)
   ```

3. **Verify the key is added to your Render account:**
   - Go to [Render Account Settings](https://dashboard.render.com/account)
   - Check if your public key is listed under "SSH Public Keys"
   - Compare the fingerprint with the output from `ssh-add -l`

4. **Add your key to the SSH agent if needed:**
   ```bash
   ssh-add ~/.ssh/id_ed25519
   ```

### Connection Timeout

- Check that the service is running in the Render Dashboard
- Verify you're using the correct region endpoint
- Ensure your service ID is correct

### SSH Session Automatically Closes

SSH sessions are automatically closed when:
- The service is redeployed or restarted
- Render performs scheduled maintenance (with 1-hour notice for existing connections)

For long-running commands, consider using Render's one-off jobs instead.

## Render's Public Key Fingerprints

For security verification, Render's SSH public key fingerprints are:

| Region    | Fingerprint                                                |
|-----------|------------------------------------------------------------|
| Oregon    | `SHA256:KkZPgnApmttFYSkdJsCi7B01sgZPMI6kY53MDbbanGM`       |
| Ohio      | `SHA256:kRDsLlrHqOyqso58sEKyO6ZFMPj7p24zfNxYJ42yXGI`       |
| Virginia  | `SHA256:NCpSwboPnqL/Nvyy2Qc8Kgzpc3P/f3w5wDphhc+UZO0`       |
| Frankfurt | `SHA256:dBRrCEA0tBkvaYLzzDw/mzaANw6nUJO961Zx806spZs`       |
| Singapore | `SHA256:CUlRyv4TZ0vmHwmhsJkII/pz2cO4IgvR+ykqnRsOQFs`       |

## Memory Usage

SSH sessions use approximately:
- 2 MB for SSH access
- 3 MB per active SSH session
- Additional memory for any processes you run

Example: If you SSH from two computers and run bash on each:
- Total usage: ~15 MB (8 MB for SSH + 7 MB for bash processes)

## Security Best Practices

1. **Use a passphrase** for your SSH private key
2. **Never share** your private key file
3. **Regularly rotate** your SSH keys (every 6-12 months)
4. **Use unique keys** for different machines
5. **Remove old keys** from Render when devices are decommissioned
6. **Be cautious** when viewing environment variables (they may contain secrets)

## Additional Resources

- [Render SSH Documentation](https://render.com/docs/ssh)
- [Render CLI Documentation](https://render.com/docs/cli)
- [SSH Key Management Best Practices](https://www.ssh.com/academy/ssh/keygen)

## Support

If you encounter issues not covered in this guide:
1. Check the [Render Community Forum](https://community.render.com/)
2. Contact Render Support through the Dashboard
3. Review your service logs for any SSH-related errors