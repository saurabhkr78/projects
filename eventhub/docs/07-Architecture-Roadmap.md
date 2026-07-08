# Architecture Roadmap (North Star)

Project: EventHub

Purpose

This document describes how the architecture of EventHub is expected to evolve over time.

The purpose is not to implement everything immediately.

Instead, it provides a long-term engineering vision that guides architectural decisions throughout the project's lifecycle.

Every new technology introduced into the system must solve a real business or technical problem.

---

# Engineering Philosophy

EventHub follows one core principle:

> Never introduce complexity before it becomes necessary.

Every architectural change must answer:

- What problem are we solving?
- Why is the current architecture insufficient?
- Why is this the right solution?
- What trade-offs are we accepting?

---

# Version 0.1

Goal

Build a clean modular monolith.

Architecture

Client

↓

REST API

↓

Go Backend

↓

PostgreSQL

Focus

- Clean Architecture
- PostgreSQL
- SQL
- pgx
- Interfaces
- Repository Pattern

---

# Version 0.2

Problem

Users need authentication.

Introduce

- JWT
- Password Hashing
- Authorization Middleware

---

# Version 0.3

Problem

Organizations need to manage venues and events.

Introduce

- Organization Module
- Venue Module
- Event Module

---

# Version 0.4

Problem

Customers need seat booking.

Introduce

- Booking Engine
- Transactions
- Row Locking
- ACID
- Concurrency Control

---

# Version 0.5

Problem

Frequently requested event data increases database load.

Introduce

Redis

Why?

- Reduce database traffic
- Faster reads
- Lower response time

Architecture

Client

↓

API

↓

Redis

↓

PostgreSQL

---

# Version 0.6

Problem

Booking API waits for email notifications.

Introduce

Kafka

Why?

Move slow work into background processing.

Architecture

Booking

↓

Kafka

↓

Notification Worker

↓

Email

---

# Version 0.7

Problem

Payment becomes an independently deployable module.

Introduce

gRPC

Architecture

Booking Service

↓

gRPC

↓

Payment Service

---

# Version 0.8

Problem

Traffic increases.

Introduce

Load Balancer

Multiple API Instances

Redis

Connection Pooling

---

# Version 0.9

Problem

Operational maturity.

Introduce

Docker

GitHub Actions

CI/CD

Monitoring

Logging

Metrics

---

# Version 1.0

Production Architecture

                    Internet
                         │
                         ▼
                  Load Balancer
                         │
          ┌──────────────┴──────────────┐
          ▼                             ▼
    API Instance 1                API Instance 2
          │                             │
          └──────────────┬──────────────┘
                         ▼
                 EventHub Backend
                         │
      ┌──────────────────┼──────────────────┐
      ▼                  ▼                  ▼
 PostgreSQL           Redis              Kafka
      │                                      │
      ▼                                      ▼
Payment Gateway                   Notification Worker

---

# Beyond Version 1

Potential future evolution

- Microservices
- Event Sourcing
- CQRS
- Kubernetes
- Service Mesh
- Distributed Tracing
- Multi-region Deployment