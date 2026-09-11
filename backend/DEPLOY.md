# Deploying to vaishnavraju.com

The whole site — HTML/CSS/JS frontend and the Go API — is compiled into a
single binary (the frontend is embedded via `go:embed`). That means
deployment is just: **build one thing, host it somewhere, point your domain
at it.**

Two paths below:

- **Render (free)** — no credit card, zero server management, custom domain
  with free HTTPS. Trade-off: the free tier sleeps after 15 minutes of no
  traffic, so the first visitor after a quiet spell waits ~30-60s for it to
  wake up. Good default for a low-traffic portfolio.
- **VPS + Caddy (~$4-6/month)** — always-on, full control, more to learn.
  Use this later if the cold start ever bothers you.

## Option A: Render (free)

Render builds straight from your GitHub repo, so first push your latest
changes:

```bash
git push origin main
```

### 1. Create the service

1. Sign up at [render.com](https://render.com) with your GitHub account (no card needed for the free tier).
2. **New +** → **Web Service** → pick the `Portfolio-X1` repo.
3. Since the Go project lives in the `backend/` folder, not the repo root, set:
   - **Root Directory**: `backend`
   - **Runtime**: `Go`
   - **Build Command**: `go build -o app ./cmd/server`
   - **Start Command**: `./app`
   - **Instance Type**: `Free`

### 2. Environment variables

In the service's **Environment** tab, add:

| Key | Value |
|---|---|
| `ALLOWED_ORIGINS` | `https://vaishnavraju.com,https://www.vaishnavraju.com` |
| `CONTACT_TO` | your real email, if you want the contact form to send you mail |
| `SMTP_HOST` / `SMTP_PORT` / `SMTP_USERNAME` / `SMTP_PASSWORD` | your SMTP provider's creds, if sending real email |

Don't set `PORT` — Render injects its own and our server already reads
`PORT` from the environment, so it just works.

### 3. Deploy

Click **Create Web Service**. Render builds and gives you a URL like
`portfolio-x1.onrender.com` — check that it loads before moving on.

### 4. Point vaishnavraju.com at it

1. In the Render service, go to **Settings → Custom Domains → Add Custom Domain**, enter `vaishnavraju.com` (and again for `www.vaishnavraju.com`).
2. Render shows you the exact DNS records to add (typically an `A`/`ALIAS`/`ANAME` record for the apex domain and a `CNAME` for `www`, pointing at Render's infrastructure).
3. Add those records in your domain registrar's DNS settings.
4. Wait for DNS to propagate (minutes to an hour) — Render will show the domain as verified and auto-provision a free TLS certificate.

Visit `https://vaishnavraju.com` — done.

### Updating the site later

Just push to `main` — Render auto-redeploys on every push (or trigger a
manual deploy from the dashboard).

## Option B: VPS + Caddy (paid, always-on)

Use this if you outgrow the free tier's cold starts.

### 1. Get a VPS

Spin up the cheapest Ubuntu 22.04/24.04 droplet/instance from DigitalOcean,
Hetzner, Linode, or similar (~$4-6/month). You'll get a public IP and root
SSH access.

### 2. Point your domain at it

| Type | Name | Value            |
|------|------|------------------|
| A    | @    | `<your VPS IP>`  |
| A    | www  | `<your VPS IP>`  |

### 3. Cross-compile the binary on your machine

No need to install Go on the server — build a Linux binary right from
Windows:

```powershell
cd backend
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o portfolio ./cmd/server
```

### 4. Copy it to the server

```powershell
scp .\portfolio root@<your VPS IP>:/opt/portfolio/portfolio
scp .env root@<your VPS IP>:/opt/portfolio/.env
```

(Create the directory first: `ssh root@<ip> "mkdir -p /opt/portfolio"`.)
In `.env` on the server, set `PORT=8080` and `ALLOWED_ORIGINS=https://vaishnavraju.com`.

### 5. Run it as a service (systemd)

Create `/etc/systemd/system/portfolio.service` on the server:

```ini
[Unit]
Description=Portfolio site
After=network.target

[Service]
WorkingDirectory=/opt/portfolio
ExecStart=/opt/portfolio/portfolio
EnvironmentFile=/opt/portfolio/.env
Restart=on-failure
User=www-data

[Install]
WantedBy=multi-user.target
```

```bash
chmod +x /opt/portfolio/portfolio
systemctl daemon-reload
systemctl enable --now portfolio
systemctl status portfolio   # should show "active (running)"
```

### 6. Install Caddy and point it at the app

```bash
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https curl
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update
sudo apt install -y caddy
```

Replace `/etc/caddy/Caddyfile` with:

```
vaishnavraju.com, www.vaishnavraju.com {
    reverse_proxy localhost:8080
}
```

```bash
sudo systemctl reload caddy
```

Caddy automatically requests and renews a Let's Encrypt TLS certificate the
first time it starts.

### Updating the site later

1. Edit the JSON files in `internal/data/` and/or the frontend in `internal/web/static/`.
2. Rebuild with `GOOS=linux GOARCH=amd64` set.
3. `scp` the new binary over the old one.
4. `ssh root@<ip> "systemctl restart portfolio"`.
