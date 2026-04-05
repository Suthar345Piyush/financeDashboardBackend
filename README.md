# Finance Backend - Golang 

An Backend application desgined for a Finance application, where the implementation of RBAC (Role based access control), where viewer and analyst has access to some part of the application, and admin have full control on the application.

The application contains pages like records and dashboard, where the system is designed in such a way that only those have full access to the app, can only operate the whole application, like performing CRUD operations on the application and can impose the CRUD operations on the viewer and analyst who have limited and less access to the application.

# Tech Stack 

-- **Golang** - Primary backend language
--**Chi** - Go's routing framework
--**golang-JWT** - For Auth Tokens and Sessions
--**sqlite** - lightweight relational database


# API Routes 

API routes are based on role the User(viewer, analyst, admin) have  

POST - /auth/register - All(Viewer, Analyst, Admin)
POST - /auth/login - All(Viewer, Analyst, Admin)
GET - /users - Admin 
PUT - /users/:id/role - Admin 
GET - /records  - Analyst/Admin 
POST - /records - Admin 
PUT - /records/:id - Admin 
Delete - /records/:id - Admin
GET - /dashboard/summary - All(Viewer, Analyst, Admin)
GET - /dashboard/trends - Analyst, Admin


# High Level Architecture of the App 

<img width="1295" height="987" alt="image" src="https://github.com/user-attachments/assets/584d3e65-fce1-47a1-9294-df15b123036b" />



# Folder Structure 

```

finance-dashboard/
├── config/
│   └── config.go
├── db/
│   ├── db.go
│   └── migrations/
│       └── 001_init.sql
├── internal/
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── dashboard_handler.go
│   │   ├── record_handler.go
│   │   └── user_handler.go
│   ├── middleware/
│   │   ├── auth.go
│   │   └── role.go
│   ├── models/
│   │   ├── record.go
│   │   └── user.go
│   ├── repository/
│   │   ├── dashboard_repo.go
│   │   ├── record_repo.go
│   │   └── user_repo.go
│   └── service/
│       ├── auth_service.go
│       ├── dashboard_service.go
│       ├── record_service.go
│       └── user_service.go
├── pkg/
│   ├── jwt/
│   │   └── jwt.go
│   ├── response/
│   │   └── response.go
│   └── validator/
│       └── validator.go
├── .env
├── .gitignore
├── README.md
├── go.mod
├── go.sum
├── main.go
├── finance.db
├── finance.db-shm
└── finance.db-wal

```

