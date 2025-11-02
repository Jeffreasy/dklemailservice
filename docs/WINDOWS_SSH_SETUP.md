# Windows SSH Setup Quick Guide

## Step 1: Open PowerShell or Command Prompt

Open PowerShell (recommended) or Command Prompt as a regular user (no admin needed).

## Step 2: Create SSH Directory (Already Done!)

The `.ssh` directory has been created at: `C:\Users\jeffrey\.ssh`

## Step 3: Generate SSH Key

Run this command in PowerShell:

```powershell
ssh-keygen -t ed25519 -f C:\Users\jeffrey\.ssh\id_ed25519
```

Or in Command Prompt:

```cmd
ssh-keygen -t ed25519 -f "%USERPROFILE%\.ssh\id_ed25519"
```

**You will be prompted:**
1. "Enter passphrase (empty for no passphrase):" - Press Enter for no passphrase, or type a secure passphrase
2. "Enter same passphrase again:" - Press Enter again, or retype the same passphrase

**This creates two files:**
- `C:\Users\jeffrey\.ssh\id_ed25519` - Private key (keep secret!)
- `C:\Users\jeffrey\.ssh\id_ed25519.pub` - Public key (share with Render)

## Step 4: Copy Your Public Key

### Option A: Using PowerShell
```powershell
Get-Content C:\Users\jeffrey\.ssh\id_ed25519.pub | Set-Clipboard
```

### Option B: Using Command Prompt
```cmd
type "%USERPROFILE%\.ssh\id_ed25519.pub"
```
Then manually copy the output (starts with `ssh-ed25519` and ends with your username)

### Option C: Open in Notepad
```cmd
notepad C:\Users\jeffrey\.ssh\id_ed25519.pub
```

## Step 5: Add Public Key to Render

1. Go to https://dashboard.render.com/account
2. Scroll to **SSH Public Keys** section
3. Click **+ Add SSH Public Key**
4. **Name**: Enter "Jeffrey's PC" or similar
5. **Key**: Paste your public key (entire contents from the .pub file)
6. Click **Add SSH Public Key**

## Step 6: Deploy Your Updated Docker Image

After committing and pushing your changes to Git, Render will automatically deploy the updated Dockerfile with SSH support.

## Step 7: Connect via SSH

### Using Render CLI (Recommended)

1. Install Render CLI:
   ```cmd
   npm install -g render-cli
   ```

2. Login:
   ```cmd
   render login
   ```

3. Connect to your service:
   ```cmd
   render ssh
   ```
   Use arrow keys to select `dklautomatie-backend` and press Enter.

### Using Direct SSH Command

Find your service ID in the Render Dashboard (looks like `srv-xxxxx`), then connect:

```cmd
ssh srv-YOUR_SERVICE_ID@ssh.oregon.render.com
```

Replace:
- `srv-YOUR_SERVICE_ID` with your actual service ID
- `oregon` with your region (ohio, virginia, frankfurt, or singapore)

## Common Commands in SSH Session

Once connected to your Render service:

```bash
# Check if app is running
ps aux | grep main

# View environment variables (be careful with sensitive data)
printenv | grep -E "DB_|SMTP_"

# Check disk usage
df -h

# Check memory usage
free -m

# Test health endpoint
wget -qO- http://localhost:8080/api/health

# Exit SSH session
exit
```

## Troubleshooting

### "bash: command not found"
Your app is running with `sh` shell. Use basic commands or switch to `bash` if available.

### "Permission denied (publickey)"
1. Verify key is in ssh-agent:
   ```cmd
   ssh-add -l
   ```
   
2. If not listed, add it:
   ```cmd
   ssh-add C:\Users\jeffrey\.ssh\id_ed25519
   ```

3. Verify it's added to Render account settings

### Connection times out
- Ensure service is running in Render Dashboard
- Check you're using correct region endpoint
- Verify service ID is correct

## Security Notes

⚠️ **Important:**
- Never share your **private key** (`id_ed25519` without .pub)
- Only share the **public key** (`id_ed25519.pub`)
- Consider using a passphrase for extra security
- Keep backups of your keys in a secure location

## Next Steps

For more detailed information, see the main documentation:
- [`docs/RENDER_SSH_SETUP.md`](RENDER_SSH_SETUP.md) - Complete SSH setup guide
- [Render SSH Documentation](https://render.com/docs/ssh)