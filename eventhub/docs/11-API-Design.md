# API Design

**Project:** EventHub

**Version:** 0.1

**API Style:** REST

---

# Purpose

This document defines the REST API exposed by EventHub.

It specifies:

- Resources
- Endpoints
- HTTP Methods
- Request Bodies
- Response Bodies
- Status Codes
- Authentication Requirements

The API contract should remain stable even if the backend implementation changes.

---

# API Principles

- RESTful design
- JSON request/response
- Stateless communication
- JWT Authentication
- Resource-oriented URLs
- Consistent HTTP status codes

---

# Base URL

/api/v1

---

# Authentication

Protected endpoints require:

Authorization: Bearer <JWT_TOKEN>

---

# Standard Response Format

## Success

{
    "success": true,
    "message": "Operation successful",
    "data": {}
}

---

## Error

{
    "success": false,
    "message": "Validation failed",
    "error": {}
}

---

# Authentication APIs

## Register

POST /auth/register

Authentication

No

Request

{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john@example.com",
    "password": "secret123"
}

Response

201 Created

{
    "success": true,
    "message": "User registered successfully"
}

---

## Login

POST /auth/login

Authentication

No

Request

{
    "email": "john@example.com",
    "password": "secret123"
}

Response

200 OK

{
    "token": "<JWT>"
}

---

# Organization APIs

## Create Organization

POST /organizations

Authentication

Organization Admin

---

## Get Organization

GET /organizations/{id}

---

## Update Organization

PUT /organizations/{id}

---

## Delete Organization

DELETE /organizations/{id}

---

# Venue APIs

## Create Venue

POST /venues

Request

{
    "name": "Main Hall",
    "capacity": 500
}

---

## Get Venues

GET /venues

---

## Get Venue

GET /venues/{id}

---

## Update Venue

PUT /venues/{id}

---

## Delete Venue

DELETE /venues/{id}

---

# Event APIs

## Create Event

POST /events

---

## Get Events

GET /events

Query Parameters

?page=1

&limit=20

&city=Delhi

&status=published

---

## Get Event

GET /events/{id}

---

## Update Event

PUT /events/{id}

---

## Delete Event

DELETE /events/{id}

---

## Publish Event

PATCH /events/{id}/publish

---

# Seat APIs

## Generate Seat Layout

POST /events/{id}/layout

---

## View Seats

GET /events/{id}/seats

---

## Reserve Seat

POST /events/{id}/reserve

---

# Booking APIs

## Create Booking

POST /bookings

Request

{
    "event_id": "...",
    "seat_ids": [
        "...",
        "..."
    ]
}

---

## Get Booking

GET /bookings/{id}

---

## Cancel Booking

DELETE /bookings/{id}

---

## Booking History

GET /bookings

---

# Payment APIs

## Create Payment

POST /payments

---

## Payment Status

GET /payments/{id}

---

# Ticket APIs

## Download Ticket

GET /tickets/{id}

---

## Verify Ticket

POST /tickets/verify

---

# Analytics APIs

## Dashboard

GET /analytics/dashboard

---

## Revenue

GET /analytics/revenue

---

## Bookings

GET /analytics/bookings

---

# HTTP Status Codes

200 OK

Successful request

201 Created

Resource created

204 No Content

Successfully deleted

400 Bad Request

Validation error

401 Unauthorized

Authentication failed

403 Forbidden

Permission denied

404 Not Found

Resource does not exist

409 Conflict

Duplicate booking

422 Unprocessable Entity

Business rule violation

500 Internal Server Error

Unexpected server error

---

# Pagination

Endpoints returning collections support:

?page=

&limit=

Example

GET /events?page=1&limit=20

---

# Filtering

Examples

GET /events?city=Delhi

GET /events?category=Music

GET /events?status=published

---

# Sorting

GET /events?sort=start_time

GET /events?sort=-created_at

---

# Future APIs

Version 1.1

- Notifications
- Coupons
- Reviews

Version 1.2

- Waitlists
- Refunds
- Recommendations

---

# API Versioning

The initial API version is:

/api/v1

Future versions

/api/v2

Breaking changes should never be introduced without creating a new API version.

---

# API Design Principles

1. Resource-oriented endpoints.
2. Proper HTTP methods.
3. Consistent response format.
4. Clear error messages.
5. Secure by default.
6. Backward compatibility whenever possible.