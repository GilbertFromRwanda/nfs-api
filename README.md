# NFS API

NFS API Server built with Go, Gin, GORM, and MySQL.

## Getting Started

### Prerequisites

- Go 1.21+
- MySQL

### Environment Variables

Create a `.env` file in the root directory:

```env
DB_HOST=127.0.0.1
DB_PORT=3306
DB_USER=root
DB_PASSWORD=yourpassword
DB_NAME=nfs_share_db_v2
JWT_SECRET=your-secret-key
PORT=5002
```

### Run the Server

```bash
go run main.go
```

The server starts on `http://localhost:5002` by default.

## API

### Base URLs

| Path       | Description                                       |
| ---------- | ------------------------------------------------- |
| `/`        | Health check - returns `{"message": "Connected"}` |
| `/api/`    | Swagger UI documentation                          |
| `/api/v1/` | All API endpoints                                 |

### Auth

| Method | Endpoint                | Description         |
| ------ | ----------------------- | ------------------- |
| POST   | `/api/v1/auth/register` | Register a new user |
| POST   | `/api/v1/auth/login`    | Login               |

### Users

| Method | Endpoint                           | Description                     |
| ------ | ---------------------------------- | ------------------------------- |
| GET    | `/api/v1/users`                    | List all users                  |
| GET    | `/api/v1/users/:id`                | Get user by ID                  |
| PUT    | `/api/v1/users/:id`                | Update a user                   |
| DELETE | `/api/v1/users/:id`                | Delete a user                   |
| PATCH  | `/api/v1/users/:id/reset-password` | Reset a user's password (admin) |

### Offices

| Method | Endpoint              | Description      |
| ------ | --------------------- | ---------------- |
| POST   | `/api/v1/offices`     | Create an office |
| GET    | `/api/v1/offices`     | List all offices |
| GET    | `/api/v1/offices/:id` | Get office by ID |
| PUT    | `/api/v1/offices/:id` | Update an office |
| DELETE | `/api/v1/offices/:id` | Delete an office |

### Notaries

| Method | Endpoint                            | Description            |
| ------ | ----------------------------------- | ---------------------- |
| POST   | `/api/v1/notaries`                  | Create a notary        |
| GET    | `/api/v1/notaries`                  | List all notaries      |
| GET    | `/api/v1/notaries/office/:officeId` | Get notaries by office |
| GET    | `/api/v1/notaries/:id`              | Get notary by ID       |
| PUT    | `/api/v1/notaries/:id`              | Update a notary        |
| DELETE | `/api/v1/notaries/:id`              | Delete a notary        |

### Agreements

| Method | Endpoint                 | Description         |
| ------ | ------------------------ | ------------------- |
| POST   | `/api/v1/agreements`     | Create an agreement |
| GET    | `/api/v1/agreements`     | List all agreements |
| GET    | `/api/v1/agreements/:id` | Get agreement by ID |
| PUT    | `/api/v1/agreements/:id` | Update an agreement |
| DELETE | `/api/v1/agreements/:id` | Delete an agreement |

### NLA Files

| Method | Endpoint                        | Description        |
| ------ | ------------------------------- | ------------------ |
| POST   | `/api/v1/nla-files`             | Create an NLA file |
| GET    | `/api/v1/nla-files`             | List all NLA files |
| GET    | `/api/v1/nla-files/:id`         | Get NLA file by ID |
| PUT    | `/api/v1/nla-files/:id`         | Update an NLA file |
| DELETE | `/api/v1/nla-files/:id`         | Delete an NLA file |
| PATCH  | `/api/v1/nla-files/:id/next`    | Move to next step  |
| PATCH  | `/api/v1/nla-files/:id/approve` | Approve NLA file   |
| PATCH  | `/api/v1/nla-files/:id/reject`  | Reject NLA file    |
| GET    | `/api/v1/nla-files/counts`      | Get status counts  |

### NLA Checklists

| Method | Endpoint                     | Description         |
| ------ | ---------------------------- | ------------------- |
| POST   | `/api/v1/nla-checklists`     | Create a checklist  |
| GET    | `/api/v1/nla-checklists`     | List all checklists |
| GET    | `/api/v1/nla-checklists/:id` | Get checklist by ID |
| PUT    | `/api/v1/nla-checklists/:id` | Update a checklist  |
| DELETE | `/api/v1/nla-checklists/:id` | Delete a checklist  |

### Checklist Notes

| Method | Endpoint                      | Description    |
| ------ | ----------------------------- | -------------- |
| POST   | `/api/v1/checklist-notes`     | Create a note  |
| GET    | `/api/v1/checklist-notes`     | List all notes |
| GET    | `/api/v1/checklist-notes/:id` | Get note by ID |
| PUT    | `/api/v1/checklist-notes/:id` | Update a note  |
| DELETE | `/api/v1/checklist-notes/:id` | Delete a note  |

### Invoices

| Method | Endpoint                                          | Description           |
| ------ | ------------------------------------------------- | --------------------- |
| POST   | `/api/v1/invoices`                                | Create an invoice     |
| GET    | `/api/v1/invoices`                                | List all invoices     |
| GET    | `/api/v1/invoices/:id`                            | Get invoice by ID     |
| PUT    | `/api/v1/invoices/:id`                            | Update an invoice     |
| DELETE | `/api/v1/invoices/:id`                            | Delete an invoice     |
| PATCH  | `/api/v1/invoices/:id/services/:serviceId/status` | Update service status |
| GET    | `/api/v1/invoices/user-summaries`                 | Get user summaries    |

### Office Services

| Method | Endpoint                                   | Description            |
| ------ | ------------------------------------------ | ---------------------- |
| POST   | `/api/v1/office-services`                  | Create a service       |
| GET    | `/api/v1/office-services`                  | List all services      |
| GET    | `/api/v1/office-services/office/:officeId` | Get services by office |
| GET    | `/api/v1/office-services/:id`              | Get service by ID      |
| PUT    | `/api/v1/office-services/:id`              | Update a service       |
| DELETE | `/api/v1/office-services/:id`              | Delete a service       |

### Office Forms

| Method | Endpoint                                | Description         |
| ------ | --------------------------------------- | ------------------- |
| POST   | `/api/v1/office-forms`                  | Create a form       |
| GET    | `/api/v1/office-forms`                  | List all forms      |
| GET    | `/api/v1/office-forms/office/:officeId` | Get forms by office |
| GET    | `/api/v1/office-forms/:id`              | Get form by ID      |
| PUT    | `/api/v1/office-forms/:id`              | Update a form       |
| DELETE | `/api/v1/office-forms/:id`              | Delete a form       |

### Office Payments

| Method | Endpoint                      | Description       |
| ------ | ----------------------------- | ----------------- |
| POST   | `/api/v1/office-payments`     | Create a payment  |
| GET    | `/api/v1/office-payments`     | List all payments |
| GET    | `/api/v1/office-payments/:id` | Get payment by ID |
| PUT    | `/api/v1/office-payments/:id` | Update a payment  |
| DELETE | `/api/v1/office-payments/:id` | Delete a payment  |

### Permissions

| Method | Endpoint                  | Description          |
| ------ | ------------------------- | -------------------- |
| POST   | `/api/v1/permissions`     | Create a permission  |
| GET    | `/api/v1/permissions`     | List all permissions |
| GET    | `/api/v1/permissions/:id` | Get permission by ID |
| PUT    | `/api/v1/permissions/:id` | Update a permission  |
| DELETE | `/api/v1/permissions/:id` | Delete a permission  |

### User Permissions

| Method | Endpoint                                            | Description                 |
| ------ | --------------------------------------------------- | --------------------------- |
| POST   | `/api/v1/user-permissions`                          | Assign a permission         |
| GET    | `/api/v1/user-permissions/user/:userId`             | Get user's permissions      |
| DELETE | `/api/v1/user-permissions/:id`                      | Remove a permission         |
| DELETE | `/api/v1/user-permissions/user/:userId/revoke-many` | Revoke multiple permissions |
| POST   | `/api/v1/user-permissions/user/:userId/assign-many` | Assign multiple permissions |

### App Features

| Method | Endpoint                   | Description       |
| ------ | -------------------------- | ----------------- |
| POST   | `/api/v1/app-features`     | Create a feature  |
| GET    | `/api/v1/app-features`     | List all features |
| GET    | `/api/v1/app-features/:id` | Get feature by ID |
| PUT    | `/api/v1/app-features/:id` | Update a feature  |
| DELETE | `/api/v1/app-features/:id` | Delete a feature  |

### Office Features

| Method | Endpoint                                   | Description                |
| ------ | ------------------------------------------ | -------------------------- |
| POST   | `/api/v1/office-features/assign`           | Assign feature to office   |
| DELETE | `/api/v1/office-features/revoke`           | Revoke feature from office |
| GET    | `/api/v1/office-features/office/:officeId` | Get office's features      |

## Regenerating Swagger Docs

After making changes to annotations, run:

```bash
swag init --parseDependency
```

### Swagger UI

Once the server is running, visit: http://localhost:5002/api/index.html
