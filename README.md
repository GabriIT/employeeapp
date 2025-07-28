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
