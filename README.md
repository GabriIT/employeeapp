# Dokku Subdirectory Deployment Guide (Frontend + Backend + PostgreSQL)

_Generated: 2025-08-10 04:34:17 UTC_

## 1) Overview & Architecture

```
We will build and deploy:
- Frontend: Next.js app — /client
- Backend: Go + Gin app — /server
- PostgreSQL: Dokku Postgres plugin (or existing Postgres)
- Domains:
  - emp.athenalabo.com  -> Frontend
  - api.emp.athenalabo.com -> Backend
- Ports inside containers:
  - Frontend: 3000
  - Backend:  8080
- Dokku 0.35.20 proxies to containers via nginx.
- Subdirectory deployment: each app has its own Dockerfile.
```

## 2) SSH key to VPS & Dokku

```
On local Ubuntu:
ssh-keygen -t ed25519 -C "you@example.com"
ssh-copy-id -i ~/.ssh/id_ed25519.pub ubuntu@209.145.61.113

Optional ~/.ssh/config:
Host vps
  HostName 209.145.61.113
  User ubuntu
  IdentityFile ~/.ssh/id_ed25519
  IdentitiesOnly yes

On VPS (add key to Dokku if needed):
sudo dokku ssh-keys:add yourname "<paste-public-key>"
Verify: sudo dokku ssh-keys:list
```

## 3) Create Dokku apps & set domains

```
On VPS:
dokku apps:create emp
dokku apps:create api-emp

dokku domains:set emp emp.athenalabo.com
dokku domains:set api-emp api.emp.athenalabo.com

Check:
dokku domains:report emp
dokku domains:report api-emp
```

## 4) Postgres (Dokku plugin) & seed data

```
Install plugin:
sudo dokku plugin:install https://github.com/dokku/dokku-postgres.git postgres

Create and link DB to backend:
dokku postgres:create employees-db
dokku postgres:link employees-db api-emp

Get DSN:
dokku postgres:info employees-db  # look for DSN

Seed:
psql "$(dokku postgres:info employees-db --dsn)" < /home/ubuntu/seed.sql
```

## 5) Local development

```
Frontend:
cd client && npm install && npm run dev   # http://localhost:3000

Backend (Go + Gin):
cd server && go mod tidy && go run main.go   # http://localhost:8080

Local Postgres tip:
postgres://postgres:postgres@127.0.0.1:5432/employees?sslmode=disable
```

## 6) Pin Go deps for Go 1.22.x

```
Newer libs forced Go 1.23+. Pin compatible versions in server/go.mod:

module server
go 1.22
require (
  github.com/gin-contrib/cors v1.7.0
  github.com/gin-gonic/gin v1.10.1
  github.com/lib/pq v1.10.9
)

Then run: go mod tidy
```

## 7) Dockerfiles (subdirectories)

```
client/Dockerfile (Next.js)
--------------------------------
FROM node:20 AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build
EXPOSE 3000
CMD ["npm","run","start"]

server/Dockerfile (Go + Gin)
--------------------------------
FROM golang:1.22
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]
```

## 8) Push via Subdirectory (UPDATED with deploy branches)

```
Add Dokku remotes (once):
git remote add dokku-emp dokku@209.145.61.113:emp
git remote add dokku-api dokku@209.145.61.113:api-emp

Option A — Persistent deploy branches (what we used)
----------------------------------------------------
# Frontend
git subtree split --prefix client vps -b deploy-client ||   (git branch -D deploy-client && git subtree split --prefix client vps -b deploy-client)
git push -f dokku-emp deploy-client:main

# Backend
git subtree split --prefix server vps -b deploy-server ||   (git branch -D deploy-server && git subtree split --prefix server vps -b deploy-server)
git push -f dokku-api deploy-server:main

Redeploy after changes:
git branch -D deploy-client || true
git subtree split --prefix client vps -b deploy-client
git push -f dokku-emp deploy-client:main

git branch -D deploy-server || true
git subtree split --prefix server vps -b deploy-server
git push -f dokku-api deploy-server:main

Option B — One-shot (no persistent branches)
--------------------------------------------
git push -f dokku-emp "$(git subtree split --prefix client vps)":main
git push -f dokku-api "$(git subtree split --prefix server vps)":main

Option C — `git subtree push` shorthand
---------------------------------------
git subtree push --prefix client dokku-emp main
git subtree push --prefix server dokku-api main
```

## 9) Map ports correctly (Dokku 0.35.20)

```
Clear & add proxy mappings:
dokku proxy:ports-clear emp
dokku proxy:ports-clear api-emp

# Frontend → 80:3000
dokku proxy:ports-add emp http:80:3000

# Backend  → 80:8080
dokku proxy:ports-add api-emp http:80:8080

Verify:
dokku proxy:report emp
dokku proxy:report api-emp
```

## 10) NGINX pitfalls & detection

```
If wrong site appears (Ubuntu default or another app):
- Ensure Dokku vhosts are included:
  sudo nginx -T | grep -n "include /home/dokku/*/nginx.conf"
- Inspect app vhost:
  sudo dokku nginx:show-config emp | sed -n '1,120p'
  sudo dokku nginx:show-config api-emp | sed -n '1,120p'
- Local host-header checks:
  curl -sI -H "Host: emp.athenalabo.com" http://127.0.0.1 | head -n 10
  curl -sI -H "Host: api.emp.athenalabo.com" http://127.0.0.1 | head -n 10
```

## 11) Force HTTP listen include (fallback trick)

```
If vhost shows listen 3000/8080 only, add a ports include:

Frontend:
sudo tee /home/dokku/emp/nginx.conf.d/ports-http.conf >/dev/null <<'EOF'
listen 80;
listen [::]:80;
EOF
sudo chown dokku:dokku /home/dokku/emp/nginx.conf.d/ports-http.conf

Backend:
sudo tee /home/dokku/api-emp/nginx.conf.d/ports-http.conf >/dev/null <<'EOF'
listen 80;
listen [::]:80;
EOF
sudo chown dokku:dokku /home/dokku/api-emp/nginx.conf.d/ports-http.conf

Reload:
sudo nginx -t && sudo systemctl reload nginx

Verify merged config:
sudo nginx -T | sed -n '/server_name emp\.athenalabo\.com/,+60p'
sudo nginx -T | sed -n '/server_name api\.emp\.athenalabo\.com/,+60p'
```

## 12) Frontend → Backend base URL & rebuild

```
Next.js inlines NEXT_PUBLIC_* at build time. After changing, rebuild emp.

Set + rebuild:
dokku config:set emp NEXT_PUBLIC_API_BASE="http://api.emp.athenalabo.com"
dokku ps:rebuild emp

Client snippet (page.tsx):
const API_BASE = process.env.NEXT_PUBLIC_API_BASE ?? "http://api.emp.athenalabo.com";
const url = new URL("/api/employee", API_BASE);
url.searchParams.set("name", name.trim());
const res = await fetch(url.toString());
```

## 13) Backend (Gin) — CORS, ILIKE, health

```
CORS:
r.Use(cors.New(cors.Config{
  AllowOrigins: []string{"http://emp.athenalabo.com", "https://emp.athenalabo.com"},
  AllowMethods: []string{"GET","OPTIONS"},
  AllowHeaders: []string{"Content-Type","Authorization"},
}))

Health:
r.GET("/healthz", func(c *gin.Context){ c.String(200,"ok") })
r.GET("/readyz",  func(c *gin.Context){ c.String(200,"ready") })

Query with ILIKE:
SELECT name, surname, birth_date, entry_date
FROM employees_data
WHERE name ILIKE $1;  -- e.g., "%bob%" 
```

## 14) Dokku healthchecks

```
Use app.json or CHECKS file (simple). Example CHECKS for backend:
WAIT=10
ATTEMPTS=30
/healthz:200
/readyz:200

Commit before pushing so Dokku waits for readiness.
```

## 15) Connectivity debugging in container

```
Enter container:
dokku enter api-emp web

Install tools if minimal image:
apt update && apt install -y iputils-ping netcat-openbsd postgresql-client

Ping / netcat:
ping 172.17.0.1
nc -zv 172.17.0.1 5432

psql with DSN:
psql "$(dokku postgres:info employees-db --dsn)" 
```

## 16) NGINX conflicts we fixed

```
- Removed Ubuntu default vhost on :80/:443 (sites-enabled/default).
- Deleted stray custom snippets that forced redirects/SSL:
  sudo rm -f /home/dokku/*/nginx.conf.d/redirect-http.conf
  sudo rm -f /home/dokku/*/nginx.conf.d/ssl-block.conf
- Always validate & reload:
  sudo nginx -t && sudo systemctl reload nginx
```

## 17) Let’s Encrypt vs Certbot (when hanging / rate-limited)

```
Try Dokku letsencrypt:
sudo dokku letsencrypt:enable emp
sudo dokku letsencrypt:enable api-emp

If it hangs, try staging then switch back:
dokku config:set --global DOKKU_LETSENCRYPT_SERVER=staging
dokku letsencrypt:enable emp && dokku letsencrypt:enable api-emp
dokku config:unset --global DOKKU_LETSENCRYPT_SERVER

If rate-limited, fallback to certbot:
sudo certbot certonly --webroot -w /var/lib/dokku/data/letsencrypt   -d emp.athenalabo.com   --non-interactive --agree-tos -m admin@athenalabo.com

Install into Dokku app:
dokku certs:add emp     < server.crt > < server.key >
dokku certs:add api-emp < server.crt > < server.key >
dokku ps:rebuild emp
dokku ps:rebuild api-emp
```

## 18) Common errors we solved

```
- 404 with '/undefined/api/...': NEXT_PUBLIC_API_BASE missing → set + rebuild.
- 400 from backend: empty 'name' → validate on client (trim, require).
- Gin wildcard static route clashing with '/employee': remove wildcard or adjust.
- GLIBC_2.34 missing: build & run on compatible base (golang:1.22).
- CORS blocked: allow exact frontend origins (HTTP+HTTPS during testing).
```

## 19) Local vs VPS differences

```
Local:
- http://localhost:3000 (frontend)
- http://localhost:8080 (backend)
- Calls are direct.

VPS:
- Use domains (emp / api.emp)
- Dokku/NGINX proxy to containers
- Must configure proxy maps + env + CORS
- TLS via letsencrypt / certbot (optional during testing)
```

## 20) Final checklist

```
✅ Apps created (emp, api-emp)
✅ Domains set
✅ Postgres created & linked; seeded
✅ Dockerfiles in /client and /server
✅ Subdirectory pushes working (deploy-client/deploy-server → main)
✅ Proxy ports: emp→80:3000, api-emp→80:8080
✅ Frontend NEXT_PUBLIC_API_BASE set; app rebuilt
✅ Backend CORS allows frontend
✅ Health endpoints OK
✅ NGINX shows correct vhosts
✅ HTTPS in place or planned

Generated: 2025-08-10 04:34:17 UTC
```

