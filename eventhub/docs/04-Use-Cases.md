# Use Cases

**Project:** EventHub

**Version:** 0.1

---

# Purpose

This document describes how different users interact with the EventHub platform.

Each use case captures:

- Actor
- Goal
- Preconditions
- Main Flow
- Alternate Flows
- Postconditions

These use cases guide API design, database design, HLD, LLD, and implementation.

---

# Actors

The platform has four primary actors.

1. Platform Administrator
2. Organization Admin
3. Organizer
4. Customer

---

# UC-001 Register User

## Actor

Customer / Organizer

---

## Goal

Create an account on EventHub.

---

## Preconditions

- User is not registered.

---

## Main Flow

1. User opens registration page.
2. User enters details.
3. System validates input.
4. Password is hashed.
5. User record is created.
6. Success response returned.

---

## Alternate Flows

- Email already exists.
- Invalid input.

---

## Postconditions

User account exists.

---

# UC-002 Login

## Actor

Customer / Organizer

---

## Goal

Authenticate and receive an access token.

---

## Main Flow

1. User submits email and password.
2. System verifies credentials.
3. JWT token generated.
4. Token returned.

---

## Alternate Flows

- Invalid password.
- User not found.

---

## Postconditions

Authenticated session established.

---

# UC-003 Create Organization

## Actor

Organization Admin

---

## Goal

Create a new organization.

---

## Main Flow

1. Submit organization details.
2. Validate input.
3. Store organization.
4. Return created organization.

---

## Alternate Flows

- Duplicate organization name.

---

## Postconditions

Organization exists.

---

# UC-004 Create Venue

## Actor

Organizer

---

## Goal

Create a venue for hosting events.

---

## Main Flow

1. Select organization.
2. Enter venue details.
3. Save venue.
4. Return success.

---

## Alternate Flows

- Invalid capacity.
- Missing required fields.

---

## Postconditions

Venue is available for future events.

---

# UC-005 Create Event

## Actor

Organizer

---

## Goal

Publish an event.

---

## Main Flow

1. Select venue.
2. Enter event information.
3. Configure date and time.
4. Configure ticket prices.
5. Save event.

---

## Alternate Flows

- Venue unavailable.
- Invalid event date.

---

## Postconditions

Event exists.

---

# UC-006 Configure Seat Layout

## Actor

Organizer

---

## Goal

Create seat arrangement.

---

## Main Flow

1. Select event.
2. Define sections.
3. Define rows.
4. Define seats.
5. Save layout.

---

## Alternate Flows

- Duplicate seat numbers.

---

## Postconditions

Seats become available.

---

# UC-007 Browse Events

## Actor

Customer

---

## Goal

View available events.

---

## Main Flow

1. Open events page.
2. Search or filter.
3. View event list.
4. Open event details.

---

## Alternate Flows

- No matching events.

---

## Postconditions

Customer selects an event.

---

# UC-008 Book Ticket

## Actor

Customer

---

## Goal

Reserve seats for an event.

---

## Preconditions

- User is logged in.
- Event is active.
- Seats are available.

---

## Main Flow

1. Select event.
2. View available seats.
3. Select seats.
4. Create booking.
5. Process payment.
6. Confirm booking.
7. Generate ticket.

---

## Alternate Flows

- Seat already booked.
- Payment failed.
- Booking timeout.

---

## Postconditions

Booking confirmed.

Ticket generated.

---

# UC-009 Cancel Booking

## Actor

Customer

---

## Goal

Cancel an existing booking.

---

## Main Flow

1. View booking.
2. Select cancel.
3. Verify cancellation policy.
4. Cancel booking.
5. Release seats.

---

## Alternate Flows

- Event already started.
- Cancellation period expired.

---

## Postconditions

Booking cancelled.

Seats become available.

---

# UC-010 View Dashboard

## Actor

Organizer

---

## Goal

Monitor event performance.

---

## Main Flow

1. Open dashboard.
2. View revenue.
3. View tickets sold.
4. View occupancy.
5. View booking trends.

---

## Alternate Flows

- No events.

---

## Postconditions

Organizer understands event performance.

---

# UC-011 Platform Administration

## Actor

Platform Admin

---

## Goal

Manage the EventHub platform.

---

## Main Flow

1. View organizations.
2. View users.
3. Suspend accounts.
4. View platform analytics.

---

## Alternate Flows

- Organization already suspended.

---

## Postconditions

Platform remains healthy.

---

# Future Use Cases

Future releases may introduce:

- Coupons
- Waitlists
- Refunds
- Notifications
- Reviews
- QR Code Verification
- Event Recommendations
- Loyalty Programs

---

# Engineering Notes

These use cases directly drive:

- High Level Design (HLD)
- Database Schema
- API Design
- Low Level Design (LLD)
- Jira Stories
- Test Cases

Every new feature added to EventHub should begin with a corresponding use case.