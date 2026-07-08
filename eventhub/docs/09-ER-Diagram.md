# Entity Relationship Diagram (ERD)

**Project:** EventHub

**Version:** 0.1

---

# Purpose

This document defines the logical structure of the EventHub database.

It identifies:

- Tables
- Columns
- Data Types
- Primary Keys
- Foreign Keys
- Relationships
- Constraints
- Future Indexes

This document serves as the blueprint for PostgreSQL schema creation.

---

# Entity Overview

The initial version of EventHub contains the following entities.

1. organizations
2. roles
3. users
4. venues
5. events
6. seat_layouts
7. seats
8. bookings
9. booking_seats
10. payments
11. tickets

---

# organizations

Represents a company or community using EventHub.

| Column | Type | Constraints |
|----------|----------|-------------|
| id | UUID | PK |
| name | VARCHAR(255) | NOT NULL |
| slug | VARCHAR(100) | UNIQUE |
| description | TEXT | NULL |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

Relationships

Organization

1

↓

N

Users

Organization

1

↓

N

Venues

---

# roles

Stores application roles.

| Column | Type | Constraints |
|----------|----------|-------------|
| id | UUID | PK |
| name | VARCHAR(100) | UNIQUE |

Examples

Platform Admin

Organization Admin

Organizer

Customer

Relationships

Role

1

↓

N

Users

---

# users

Stores all platform users.

| Column | Type | Constraints |
|----------|----------|-------------|
| id | UUID | PK |
| organization_id | UUID | FK |
| role_id | UUID | FK |
| first_name | VARCHAR(100) | NOT NULL |
| last_name | VARCHAR(100) | NOT NULL |
| email | VARCHAR(255) | UNIQUE |
| password_hash | TEXT | NOT NULL |
| phone | VARCHAR(20) | NULL |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

Relationships

Organization

1

↓

N

Users

Role

1

↓

N

Users

User

1

↓

N

Bookings

---

# venues

Stores event venues.

| Column | Type | Constraints |
|----------|----------|-------------|
| id | UUID | PK |
| organization_id | UUID | FK |
| name | VARCHAR(255) | NOT NULL |
| address | TEXT | NOT NULL |
| city | VARCHAR(100) | NOT NULL |
| state | VARCHAR(100) | NOT NULL |
| country | VARCHAR(100) | NOT NULL |
| capacity | INTEGER | CHECK(capacity > 0) |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

Relationships

Organization

1

↓

N

Venues

Venue

1

↓

N

Events

---

# events

Stores published events.

| Column | Type | Constraints |
|----------|----------|-------------|
| id | UUID | PK |
| venue_id | UUID | FK |
| organizer_id | UUID | FK(users.id) |
| title | VARCHAR(255) | NOT NULL |
| description | TEXT | NULL |
| start_time | TIMESTAMP | NOT NULL |
| end_time | TIMESTAMP | NOT NULL |
| status | VARCHAR(50) | NOT NULL |
| created_at | TIMESTAMP | NOT NULL |
| updated_at | TIMESTAMP | NOT NULL |

Relationships

Venue

1

↓

N

Events

User

1

↓

N

Events

Event

1

↓

1

Seat Layout

Event

1

↓

N

Bookings

---

# seat_layouts

Defines seating configuration.

| Column | Type | Constraints |
|----------|----------|-------------|
| id | UUID | PK |
| event_id | UUID | UNIQUE FK |
| name | VARCHAR(100) | NOT NULL |

Relationships

Event

1

↓

1

Seat Layout

Seat Layout

1

↓

N

Seats

---

# seats

Stores every seat.

| Column | Type | Constraints |
|----------|----------|-------------|
| id | UUID | PK |
| layout_id | UUID | FK |
| section | VARCHAR(50) | NOT NULL |
| row | VARCHAR(10) | NOT NULL |
| number | VARCHAR(10) | NOT NULL |
| category | VARCHAR(50) | NOT NULL |
| price | NUMERIC(10,2) | NOT NULL |

Relationship

Seat Layout

1

↓

N

Seats

---

# bookings

Stores customer bookings.

| Column | Type | Constraints |
|----------|----------|-------------|
| id | UUID | PK |
| event_id | UUID | FK |
| user_id | UUID | FK |
| status | VARCHAR(30) | NOT NULL |
| total_amount | NUMERIC(10,2) | NOT NULL |
| booked_at | TIMESTAMP | NOT NULL |

Relationships

User

1

↓

N

Bookings

Booking

1

↓

N

Booking Seats

Booking

1

↓

1

Payment

Booking

1

↓

1

Ticket

---

# booking_seats

Junction table.

Allows one booking to contain multiple seats.

| Column | Type | Constraints |
|----------|----------|-------------|
| booking_id | UUID | PK FK |
| seat_id | UUID | PK FK |

Relationship

Booking

N

↓

N

Seats

Implemented through booking_seats.

---

# payments

Stores payment information.

| Column | Type | Constraints |
|----------|----------|-------------|
| id | UUID | PK |
| booking_id | UUID | UNIQUE FK |
| amount | NUMERIC(10,2) | NOT NULL |
| payment_status | VARCHAR(50) | NOT NULL |
| payment_method | VARCHAR(50) | NOT NULL |
| transaction_reference | VARCHAR(255) | NULL |

Relationship

Booking

1

↓

1

Payment

---

# tickets

Stores generated tickets.

| Column | Type | Constraints |
|----------|----------|-------------|
| id | UUID | PK |
| booking_id | UUID | UNIQUE FK |
| qr_code | TEXT | NULL |
| issued_at | TIMESTAMP | NOT NULL |

Relationship

Booking

1

↓

1

Ticket

---

# Cardinality Summary

Organization

1 → N Users

Organization

1 → N Venues

Role

1 → N Users

Venue

1 → N Events

User

1 → N Events

Event

1 → 1 Seat Layout

Seat Layout

1 → N Seats

User

1 → N Bookings

Event

1 → N Bookings

Booking

N → N Seats

Booking

1 → 1 Payment

Booking

1 → 1 Ticket

---

# Initial Indexes

Unique

- users.email
- organizations.slug

Foreign Keys

- organization_id
- venue_id
- event_id
- booking_id
- role_id
- user_id

Future

Composite Indexes

- events(start_time, status)
- bookings(user_id, booked_at)
- seats(layout_id, category)

---

# Future Database Evolution

Version 1.1

- Coupons
- Notifications
- Audit Logs

Version 1.2

- Reviews
- Refunds
- Waitlists

Version 2

- Read Replicas
- Partitioning
- Materialized Views