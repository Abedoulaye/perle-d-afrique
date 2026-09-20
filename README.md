# Diariata's Fabrics

A full-stack e-commerce platform for a local African clothing business.

![Screenshot](./screenshot.png)

## Live Demo

[perledafrique.com](https://perledafrique.com)

## Features

- User authentication with JWT + refresh tokens
- Product browsing with pagination
- Shopping cart with quantity management
- Secure checkout with Stripe
- Order history and email verification
- Admin panel for product management

## Tech Stack

**Backend:** Go, net/http, PostgreSQL, pgx, Stripe, Resend
**Frontend:** React, TypeScript, React Router, Vite
**Infrastructure:** Railway (backend + DB), Vercel (frontend)

## Architecture

- REST API with middleware (auth, rate limiting, CORS)
- Row-level locking for stock management during checkout
- Webhook-driven payment confirmation
- Environment-based configuration

## Local Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- PostgreSQL 14+
- Stripe account (for payments)
- Resend account (for email verification)

### Backend

1. Clone the repository:
2. Install Go dependencies and npm dependencies
3. Create a .env file in your backend directory the following:
   a variable for username, password, database host, database port, database name, jwt secret, stripe secret key, stripe webhook secret, resend api key, cors origin
4. create a .env file in your frontend directory with the following:
   vite api url and vite stripe publishable key
5. run the database migrations
6. start the server

## What I Learned

- Backend Architecture
- Database transactions
- **Storing money as integers** - I learned to store prices as `int64` cents
  instead of floats to avoid precision errors in arithmetic.
- JWT with refresh tokens
- Rate Limiting
- input validation
- Stripe integration
- Metadata linking
- Structured logging
