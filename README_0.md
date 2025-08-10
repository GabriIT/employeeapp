GO Version 1.22.2


employeeapp/
├── go.mod
├── go.sum
├── main.go               # Go backend main application
├── static/
│   └── index.html        # HTML frontend
└── README.md             # Optional documentation


# Run with:
go mod init employeeapp
go get github.com/lib/pq
go mod tidy
go run main.go


# Using Gin Gin as HTTP web framework
branch gin-basic


# Using NGINX with Client and Server

employeeapp-nginx/
├── docker-compose.yml
├── nginx.conf.prod
├── client/
│   ├── Dockerfile      ✅
│   ├── package.json
│   └── app/page.tsx
└── server/
    ├── Dockerfile      ✅
    ├── main.go
    ├── go.mod
    └── go.sum

# server/go.mod
cd server
go mod init server
go get github.com/gin-gonic/gin
go get github.com/lib/pq

# Step 2: Create client/ (Next.js frontend)

npx create-next-app@latest client --ts --app --no-tailwind --eslint

# Step 5: Build & Run
docker compose up --build


# Access to Postgres on VPS
Role is 'empapp', password is 'postgres'

# become postgres superuser
sudo -u postgres psql <<'SQL'
-- Create a dedicated login for the API (change the password!)
CREATE ROLE empapp WITH LOGIN PASSWORD 'CHANGEME-strong-password';

-- Create DB and make empapp the owner (skip if 'employees' already exists)
CREATE DATABASE employees OWNER empapp;

-- Make sure future objects are accessible to the owner (safe defaults)
ALTER DATABASE employees OWNER TO empapp;
SQL

# If you created the DB above:
psql -h 127.0.0.1 -U empapp -d employees -f ~/seed.sql

# If the file contains "CREATE DATABASE employees; \c employees" at the top,
# run it as superuser instead:
# PGPASSWORD=<postgres_password> psql -h 127.0.0.1 -U postgres -f ~/seed.sql

sudo -u postgres psql -d employees <<'SQL'
REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA public TO empapp;

GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO empapp;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO empapp;
SQL

sudo dokku config:set api-emp \
  DATABASE_URL="postgres://empapp:postgres@172.17.0.1:5432/employees?sslmode=disable" \
  PGSSLMODE=disable
sudo dokku ps:restart api-emp
