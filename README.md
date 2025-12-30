# 1-to-1 Real-Time Chat Application

A secure 1-to-1 chat application using JWT authentication and WebSockets, built with Go and Next.js.

<img width="7509" height="6340" alt="Architecture" src="https://github.com/user-attachments/assets/b7fd8bf5-1269-41dc-8e2b-2a27fb632e46" />


---

## Overview

This project is a real-time 1-to-1 chat application focused on backend architecture, security, and real-time communication. It combines JWT-based authentication, RESTful APIs, and authenticated WebSocket connections to enable secure messaging between users. The authentication logic was reused from a standalone JWT authentication project and extended to support real-time chat functionality.

---

## Tech Stack

### Backend
- Go (Golang)
- Gin Framework (RESTful APIs)
- Gorilla WebSocket (Real-time communication)
- JWT Authentication (Access Token & Refresh Token)

### Frontend
- Next.js

### Communication
- RESTful APIs
- WebSocket over TCP

### Deployment
- Docker
- Render.com (Backend)
- Vercel (Frontend)

---

## Architecture & Technical Details

### Authentication
- JWT-based authentication using access and refresh tokens
- Authentication logic reused from a separate JWT project
- Middleware validates JWTs and extracts user identity for protected routes
- All protected HTTP endpoints require a valid access token

### RESTful API (Gin)
- Handles user authentication
- Provides protected routes
- Manages token validation and refresh
- Secured using JWT middleware

### WebSocket Communication (Gorilla WebSocket)
- WebSocket connections start as standard HTTP requests
- JWT is validated from request headers before upgrading the connection
- Connection is upgraded to a persistent, bidirectional WebSocket over TCP
- Only authenticated users can establish WebSocket connections
- Each WebSocket connection is mapped to an authenticated user for secure 1-to-1 message routing

### Frontend
- Built with Next.js
- Minimal UI by design
- Focused on connecting to backend APIs and WebSocket server
- Provides basic functionality for sending and receiving messages

---

## Deployment

- Backend is containerized using Docker and deployed on Render.com
- Frontend is deployed on Vercel
- The application is publicly accessible via the Vercel frontend link

> ⚠️ **Note**  
> The backend runs on a free Render.com instance, which automatically spins down during inactivity.  
> The first request may take **50 seconds or longer** while the server wakes up. Subsequent requests will be fast once the backend is running.

---

## How to Test

1. Open the Vercel frontend link -> https://chatapp1to1.vercel.app/
2. Wait for the backend to wake up if it has been inactive
3. Register or log in
4. Start a 1-to-1 chat session
