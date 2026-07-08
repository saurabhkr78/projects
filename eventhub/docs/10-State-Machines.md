# State Machines

**Project:** EventHub

**Version:** 0.1

---

# Purpose

This document defines the lifecycle of important business entities in EventHub.

Every business entity changes state during its lifetime.

State machines ensure that:

- Invalid transitions are prevented.
- Business rules remain consistent.
- Services implement predictable workflows.
- Future features are easier to extend.

---

# What is a State Machine?

A state machine describes:

- Current State
- Allowed Next States
- Invalid States
- Transition Rules

Instead of asking:

"What columns does Booking have?"

We ask:

"What can happen to a booking?"

---

# Event Lifecycle

An event moves through the following states.

Draft

↓

Published

↓

Bookings Open

↓

Bookings Closed

↓

Completed

↓

Archived

---

## State Descriptions

### Draft

The organizer is creating the event.

Customers cannot view it.

---

### Published

The event is publicly visible.

Bookings are still disabled.

---

### Bookings Open

Customers can reserve seats.

---

### Bookings Closed

New bookings are not accepted.

---

### Completed

The event has finished.

---

### Archived

Historical record only.

Cannot be modified.

---

## Allowed Transitions

Draft

↓

Published

Published

↓

Bookings Open

Bookings Open

↓

Bookings Closed

Bookings Closed

↓

Completed

Completed

↓

Archived

---

## Invalid Transitions

Draft

✗ Completed

Draft

✗ Archived

Completed

✗ Draft

Archived

✗ Published

---

# Booking Lifecycle

Pending

↓

Payment Pending

↓

Confirmed

↓

Checked In

↓

Completed

---

Alternative

Payment Pending

↓

Failed

↓

Cancelled

---

Alternative

Confirmed

↓

Cancelled

---

## State Descriptions

### Pending

Booking record created.

Seats are temporarily reserved.

---

### Payment Pending

Customer is paying.

---

### Confirmed

Payment successful.

Seats permanently allocated.

Ticket generated.

---

### Failed

Payment failed.

Seats released.

---

### Cancelled

Booking cancelled.

Seats released according to policy.

---

### Checked In

Customer entered the venue.

---

### Completed

Event completed.

Booking becomes historical.

---

# Allowed Booking Transitions

Pending

↓

Payment Pending

Payment Pending

↓

Confirmed

Payment Pending

↓

Failed

Confirmed

↓

Checked In

Checked In

↓

Completed

Confirmed

↓

Cancelled

Failed

↓

Cancelled

---

# Invalid Booking Transitions

Completed

✗ Pending

Cancelled

✗ Confirmed

Failed

✗ Confirmed

---

# Payment Lifecycle

Initiated

↓

Authorized

↓

Captured

↓

Completed

---

Alternative

Initiated

↓

Failed

---

Alternative

Authorized

↓

Refunded

---

## State Descriptions

### Initiated

Payment request created.

---

### Authorized

Payment provider approved payment.

---

### Captured

Money captured successfully.

---

### Completed

Booking finalized.

---

### Failed

Payment unsuccessful.

---

### Refunded

Money returned.

---

# Ticket Lifecycle

Generated

↓

Downloaded

↓

Checked In

↓

Expired

---

# Organization Lifecycle

Pending Verification

↓

Active

↓

Suspended

↓

Archived

---

# User Lifecycle

Registered

↓

Email Verified

↓

Active

↓

Suspended

↓

Deleted

---

# Business Rules

## Event Rules

- Archived events cannot be modified.
- Completed events cannot accept bookings.
- Bookings cannot open before publishing.

---

## Booking Rules

- A seat cannot belong to two confirmed bookings.
- Cancelled bookings release seats.
- Failed payments release seats.

---

## Payment Rules

- Every confirmed booking must have one successful payment.
- Failed payments never generate tickets.

---

## Ticket Rules

- Tickets are generated only after successful payment.
- Checked-in tickets cannot be reused.

---

# Implementation Strategy

Business state transitions belong in the **Service Layer**.

Example

BookingService

CreateBooking()

↓

ProcessPayment()

↓

ConfirmBooking()

↓

GenerateTicket()

Repositories only persist state.

Handlers only receive requests.

Business rules remain inside services.

---

# Future Evolution

Future versions may introduce:

- Waitlists
- Refund Approval
- Event Rescheduling
- Ticket Transfer
- Seat Upgrade
- Dynamic Pricing