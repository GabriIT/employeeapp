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


