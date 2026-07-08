# High Level Design (HLD)

**Project:** EventHub

**Version:** 0.1

**Architecture Style:** Modular Monolith

---

# Purpose

This document describes the overall architecture of EventHub.

It defines:

- Major system components
- Responsibilities
- Communication between modules
- Deployment architecture
- Technology stack

This document intentionally avoids implementation details.

---

# Architectural Goals

The system should be:

- Simple to understand
- Easy to maintain
- Easy to test
- Easy to extend
- Ready for future scaling

The initial architecture prioritizes simplicity while preserving a clear migration path toward microservices.

---

# Architecture Decision

EventHub will initially be developed as a **Modular Monolith**.

Reasons:

- Faster development
- Easier debugging
- Simpler deployment
- Lower operational complexity
- Better learning experience
- Clear module boundaries

As the system evolves, modules can be extracted into independent services without changing the business logic.

---

# High Level Architecture

                          Internet
                              │
                              ▼
                     REST API (HTTP)
                              │
                              ▼
                     Gorilla Mux Router
                              │
                              ▼
                  +----------------------+
                  |      EventHub        |
                  |  Modular Monolith    |
                  +----------------------+
                              │
      ┌───────────────────────┼────────────────────────┐
      ▼                       ▼                        ▼
 Authentication         Organization             Event Module
      │                       │                        │
      ├───────────────┬────────┴──────────────┬────────┤
      ▼               ▼                       ▼
 Venue Module    Booking Module        Payment Module
                              │
                              ▼
                     PostgreSQL Database

Future Integrations

Booking Module

↓

Redis

↓

Kafka

↓

Email Service

↓

Payment Gateway

---

# Module Overview

The system consists of the following business modules.

---

## Authentication Module

Responsibilities

- Registration
- Login
- JWT
- Authorization

---

## Organization Module

Responsibilities

- Organization Management
- Organization Members
- Roles

---

## User Module

Responsibilities

- User Profile
- Account Management
- Permissions

---

## Venue Module

Responsibilities

- Venue CRUD
- Capacity Management

---

## Event Module

Responsibilities

- Event CRUD
- Event Publishing
- Event Search

---

## Seat Module

Responsibilities

- Seat Layout
- Seat Categories
- Availability

---

## Booking Module

Responsibilities

- Seat Reservation
- Booking
- Cancellation

---

## Payment Module

Responsibilities

- Payment Processing
- Payment Status

Initially

Mock Payment

Later

Stripe / Razorpay

---

## Analytics Module

Responsibilities

- Revenue Reports
- Ticket Sales
- Occupancy
- Dashboard

---

# Internal Layer Architecture

Every module follows the same layered architecture.

                 Router
                    │
                    ▼
                Handler
                    │
                    ▼
                Service
                    │
                    ▼
              Repository
                    │
                    ▼
              PostgreSQL

Every module remains independent.

Example

internal/

    booking/

        handler.go

        service.go

        repository.go

        routes.go

        model.go

---

# Request Flow

Customer

↓

HTTP Request

↓

Router

↓

Handler

↓

Service

↓

Repository

↓

PostgreSQL

↓

Repository

↓

Service

↓

Handler

↓

JSON Response

---

# Technology Stack

Language

- Go

Router

- Gorilla Mux

Database

- PostgreSQL

Database Driver

- pgx

Configuration

- godotenv

Authentication

- JWT

Password Hashing

- bcrypt

Validation

- go-playground/validator

---

# Deployment Architecture (Version 1)

                Client
                   │
                   ▼
           EventHub API
                   │
                   ▼
            PostgreSQL

One Backend

One Database

One Deployment

---

# Future Evolution

Version 1.1

Client

↓

EventHub API

↓

Redis

↓

PostgreSQL

---

Version 1.2

Client

↓

EventHub API

↓

Kafka

↓

Email Service

↓

PostgreSQL

---

Version 1.3

Client

↓

API Gateway

↓

Booking Service

↓

Payment Service

↓

Notification Service

↓

PostgreSQL

---

# Design Principles

1. Modular Monolith
2. Clean Architecture
3. Separation of Concerns
4. Dependency Injection
5. Repository Pattern
6. SOLID Principles
7. Business Logic independent of Infrastructure

---

# Module Communication

In Version 1

All modules communicate through method calls.

Booking Service

↓

Payment Service

↓

Analytics Service

No network communication exists between modules.

Future versions may replace internal method calls with gRPC or event-driven communication where appropriate.

---

# Scalability Strategy

Current

Scale vertically by increasing resources.

Future

- Load Balancer
- Multiple API Instances
- Redis Cache
- Kafka
- Background Workers
- gRPC
- Microservices

These improvements will only be introduced when justified by product requirements.

---

# Risks

Potential risks include:

- High booking concurrency
- Database bottlenecks
- Large analytics queries
- Payment failures
- Double booking

These risks will be addressed incrementally throughout the project.

---

# Success Criteria

The HLD is considered successful if:

- Module boundaries are clearly defined.
- Responsibilities are separated.
- Architecture supports future growth.
- Business logic remains independent of infrastructure.
- Future technologies can be integrated with minimal redesign.