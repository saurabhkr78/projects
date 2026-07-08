# Domain Model

**Project:** EventHub

**Version:** 0.1

---

# Purpose

The Domain Model identifies the core business entities of EventHub and describes how they relate to each other.

This document is independent of:

- PostgreSQL
- Go
- REST APIs
- Frameworks

It represents the business itself.

---

# What is a Domain?

A domain is the real-world problem our software solves.

For EventHub, the domain is:

> Managing organizations, events, venues, bookings, and ticket sales.

Everything else exists to support this business.

---

# Core Business Entities

The following entities represent the core objects within the EventHub domain.

---

## Organization

Represents a company or community using EventHub.

Examples

- CodeConf India
- MusicFest
- TEDx Delhi
- AWS User Group

Responsibilities

- Own venues
- Own events
- Manage organizers

Relationships

Organization

↓

Users

↓

Venues

↓

Events

---

## User

Represents anyone using the platform.

Types

- Platform Admin
- Organization Admin
- Organizer
- Customer

Responsibilities

- Authenticate
- Create events (Organizer)
- Book events (Customer)

Relationships

User

↓

Bookings

---

## Role

Defines permissions within the system.

Examples

Platform Admin

Organization Admin

Organizer

Customer

A user has one role.

A role can belong to many users.

---

## Venue

Represents a physical location.

Examples

- Hall A
- Conference Center
- Auditorium
- Stadium

Responsibilities

- Host events

Relationships

Venue

↓

Events

---

## Event

Represents an activity customers can attend.

Examples

- Go Conference
- AI Summit
- Rock Concert
- Startup Meetup

Responsibilities

- Belongs to one venue
- Owned by one organization
- Has seat layout
- Accept bookings

Relationships

Organization

↓

Event

↓

Seat Layout

↓

Bookings

---

## Seat Layout

Defines the seating arrangement.

Examples

VIP

Gold

Silver

Balcony

Responsibilities

- Generate seats
- Categorize seats

Relationships

Seat Layout

↓

Seats

---

## Seat

Represents a single seat.

Examples

A1

A2

B5

VIP-12

Responsibilities

- Be reserved
- Be booked

Relationships

Seat

↓

Booking

---

## Booking

Represents a customer's reservation.

Responsibilities

- Reserve seats
- Track booking status
- Connect customer and event

Relationships

User

↓

Booking

↓

Payment

↓

Ticket

---

## Payment

Represents payment for a booking.

Responsibilities

- Store payment status
- Verify payment
- Track transaction

Relationships

Booking

↓

Payment

---

## Ticket

Represents the final ticket issued after successful booking.

Responsibilities

- Entry validation
- Ticket download
- QR Code (Future)

Relationships

Booking

↓

Ticket

---

# Business Relationships

Organization

1

↓

N

Users

---

Organization

1

↓

N

Venues

---

Venue

1

↓

N

Events

---

Event

1

↓

1

Seat Layout

---

Seat Layout

1

↓

N

Seats

---

User

1

↓

N

Bookings

---

Booking

1

↓

1

Payment

---

Booking

1

↓

1

Ticket

---

# Domain Rules

The following business rules define the behavior of the system.

## Organization

- Every event belongs to exactly one organization.
- Every venue belongs to exactly one organization.

---

## Venue

- A venue can host many events.
- An event can only be hosted in one venue.

---

## Event

- Every event has exactly one organizer.
- Every event has one seat layout.
- Events cannot be booked after they end.

---

## Seat

- A seat belongs to exactly one seat layout.
- A seat cannot be booked twice.

---

## Booking

- A booking belongs to one customer.
- A booking belongs to one event.
- A booking can contain one or more seats.
- Cancelled bookings release reserved seats.

---

## Payment

- Every successful booking must have a successful payment.
- Failed payments should not confirm bookings.

---

# Future Domain Expansion

The domain is intentionally designed for future growth.

Future entities may include:

- Coupons
- Reviews
- Notifications
- Audit Logs
- Loyalty Points
- Refunds
- Waitlists
- Recommendations

---

# Domain Diagram (Conceptual)

                    EventHub

                       │

      ┌────────────────┴────────────────┐

      ▼                                 ▼

Organization                      Platform Admin

      │

      ├────────────── Users

      │

      ├────────────── Venues

      │                    │

      │                    ▼

      │                 Events

      │                    │

      │                    ▼

      │               Seat Layout

      │                    │

      │                    ▼

      │                  Seats

      │                    │

      ▼                    ▼

Customers ─────────────► Bookings

                              │

                   ┌──────────┴──────────┐

                   ▼                     ▼

               Payments             Tickets

---

# Guiding Principles

- The domain should remain independent of technology.
- Business rules should exist within the domain, not the framework.
- Infrastructure should support the domain, never control it.
- Future architectural decisions must preserve domain integrity.