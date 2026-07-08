# System Context

**Project:** EventHub

**Version:** 0.1

---

# Purpose

The System Context document defines the boundary of EventHub.

It answers:

- Who uses the system?
- Which external systems interact with EventHub?
- What responsibilities belong to EventHub?
- What responsibilities belong to external systems?

This serves as the highest-level architectural view before moving into the High Level Design (HLD).

---

# System Overview

EventHub is a multi-tenant SaaS platform that enables organizations to manage events and allows customers to discover and book tickets.

The platform exposes REST APIs and manages authentication, bookings, payments, and analytics while integrating with external services when necessary.

---

# Primary Actors

The following users interact directly with EventHub.

---

## Platform Administrator

Responsibilities

- Manage organizations
- View platform analytics
- Suspend accounts
- Monitor platform health

---

## Organization Admin

Responsibilities

- Create organization
- Invite organizers
- Manage organization settings

---

## Organizer

Responsibilities

- Create venues
- Create events
- Configure seat layouts
- Publish events
- Monitor bookings

---

## Customer

Responsibilities

- Register
- Login
- Browse events
- Search events
- Book tickets
- Cancel bookings
- View booking history

---

# External Systems

The following systems are outside the EventHub boundary.

---

## PostgreSQL

Purpose

Primary relational database.

Stores

- Users
- Organizations
- Venues
- Events
- Seats
- Bookings
- Payments

---

## Payment Gateway (Future)

Examples

- Razorpay
- Stripe

Responsibilities

- Process payments
- Verify transactions

---

## Email Service (Future)

Examples

- SendGrid
- Amazon SES

Responsibilities

- Booking confirmation
- Password reset
- Event notifications

---

## SMS Provider (Future)

Responsibilities

- OTP verification
- Booking notifications

---

## Redis (Future)

Responsibilities

- Caching
- Session storage
- Rate limiting

---

## Kafka (Future)

Responsibilities

- Event streaming
- Notifications
- Analytics
- Background processing

---

# System Boundary

Everything inside this boundary is owned by EventHub.

+-------------------------------------------------------------+

                    EventHub Platform

+-------------------------------------------------------------+

Authentication

Organization Management

Venue Management

Event Management

Seat Management

Booking Engine

Payment Module

Analytics

REST API

+-------------------------------------------------------------+

Everything outside this boundary is an external dependency.

---

# Context Diagram

                    +----------------------+
                    | Platform Admin       |
                    +----------+-----------+
                               |
                               |
                    +----------v-----------+
                    |                      |
                    |      EventHub        |
                    |                      |
                    +----------+-----------+
                               |
         +---------------------+----------------------+
         |                     |                      |
         |                     |                      |
+--------v--------+   +--------v--------+   +---------v--------+
| Organization    |   | Organizer       |   | Customer         |
+-----------------+   +-----------------+   +------------------+

                               |
                               |
                     +---------v---------+
                     | PostgreSQL        |
                     +-------------------+

Future Integrations

EventHub

↓

Redis

↓

Kafka

↓

Payment Gateway

↓

Email Service

↓

SMS Provider

---

# Responsibilities of EventHub

EventHub is responsible for:

- Authentication
- Authorization
- Business Rules
- Event Management
- Booking Engine
- Payment Orchestration
- Analytics
- Data Persistence

---

# Responsibilities NOT Owned by EventHub

The following are delegated to external services.

- Payment Processing
- Email Delivery
- SMS Delivery
- Cache Storage
- Event Streaming

---

# Technology Evolution

Version 0.1

Client

↓

REST API

↓

Go Backend

↓

PostgreSQL

---

Version 0.5

Client

↓

REST API

↓

Go Backend

↓

Redis

↓

PostgreSQL

---

Version 0.7

Client

↓

REST API

↓

Go Backend

↓

Kafka

↓

Redis

↓

PostgreSQL

---

Version 1.0

                     Internet
                          │
                          ▼
                  Load Balancer
                          │
                          ▼
                     EventHub API
                          │
        ┌─────────────────┼─────────────────┐
        ▼                 ▼                 ▼
 PostgreSQL            Redis             Kafka
        │                                   │
        ▼                                   ▼
 Payment Gateway                    Email Service

---

# Architectural Principles

1. EventHub owns the business logic.
2. External systems provide infrastructure capabilities.
3. Business logic must remain independent of external services.
4. External services can be replaced with minimal impact on the domain.
5. The platform starts as a modular monolith and evolves only when justified by business needs.

---

# Future Evolution

The current system is intentionally designed so that future versions can introduce:

- gRPC
- Microservices
- Event-Driven Architecture
- Distributed Caching
- Background Workers
- Cloud Deployment

without requiring a complete redesign.