# Portfolio backend

Go API serving portfolio content and handling the contact form. Stdlib only —
no external dependencies to fetch.

## Run locally

```bash
cp .env.example .env
# edit .env if you want real email sending; otherwise contact messages
# just get logged to the console.

go run ./cmd/server
```

Server listens on `:8080` by default (`PORT` in `.env`).

## Endpoints

| Method | Path              | Description                          |
|--------|-------------------|---------------------------------------|
| GET    | /api/health       | Liveness check                        |
| GET    | /api/profile      | About/profile info                    |
| GET    | /api/skills       | Skills grouped by category            |
| GET    | /api/projects     | Project list                          |
| GET    | /api/experience   | Work experience list                  |
| POST   | /api/contact      | `{ name, email, message }` -> sends/logs an email |

## Content

`/api/profile`, `/api/skills`, `/api/projects`, and `/api/experience` are
served straight from the JSON files in `internal/data/`. Edit those files
with your real info and rebuild — no code changes needed.

## Config (env vars)

- `PORT` — listen port (default `8080`)
- `ALLOWED_ORIGINS` — comma-separated CORS allowlist (default `http://localhost:5173`)
- `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD` — SMTP credentials for sending contact emails
- `CONTACT_FROM` — from address used for outgoing mail
- `CONTACT_TO` — where contact form submissions get sent
- `CONTACT_RATE_LIMIT_PER_MIN` — max contact submissions per IP per minute (default `5`)

If `SMTP_HOST`/`CONTACT_TO` are unset, the contact form still works but just
logs the message instead of emailing it — convenient for local dev.
