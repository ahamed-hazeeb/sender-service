# 📤 Points Sender Service

Transfer initiation and email notification microservice for Virtual Points Transfer System.

## Features

- Transfer initiation with validation
- Email notifications with HTML templates
- Transfer status management
- Integration with Auth Service

## API Endpoints

- `POST /transfer` - Initiate points transfer
- `GET /transfers/:userId` - Get user transfer history
- `POST /transfer/:id/complete` - Complete transfer (Saga pattern)

## Tech Stack

- **Go** with Gin framework
- **PostgreSQL** with GORM
- **SMTP** for email notifications

## Quick Start

```bash
git clone https://github.com/ahamed-hazeeb/sender-service.git
cd points-sender-service
go run cmd/server/main.go
```
