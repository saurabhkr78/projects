# Requirements Traceability Matrix (RTM)

**Project:** EventHub

**Version:** 0.1

---

## Purpose

The Requirements Traceability Matrix (RTM) ensures that every business requirement can be traced throughout the software development lifecycle.

Each requirement is linked to:

- Product Requirement
- Jira Epic
- Sprint
- Module
- API
- Database
- Implementation Status

---

| Requirement ID | Requirement | Priority | Epic | Sprint | Module | API | Database | Status |
|----------------|-------------|----------|------|---------|--------|-----|----------|--------|
| FR-001 | User Registration | High | Authentication | Sprint 2 | User | POST /auth/register | users | Planned |
| FR-002 | User Login | High | Authentication | Sprint 2 | Auth | POST /auth/login | users | Planned |
| FR-003 | Create Organization | High | Organization | Sprint 3 | Organization | POST /organizations | organizations | Planned |
| FR-004 | Create Venue | High | Venue Management | Sprint 4 | Venue | POST /venues | venues | Planned |
| FR-005 | Update Venue | Medium | Venue Management | Sprint 4 | Venue | PUT /venues/{id} | venues | Planned |
| FR-006 | Delete Venue | Medium | Venue Management | Sprint 4 | Venue | DELETE /venues/{id} | venues | Planned |
| FR-007 | Create Event | High | Event Management | Sprint 5 | Event | POST /events | events | Planned |
| FR-008 | Update Event | Medium | Event Management | Sprint 5 | Event | PUT /events/{id} | events | Planned |
| FR-009 | Publish Event | High | Event Management | Sprint 5 | Event | PATCH /events/{id}/publish | events | Planned |
| FR-010 | Search Events | High | Event Management | Sprint 5 | Event | GET /events | events | Planned |
| FR-011 | View Event Details | High | Event Management | Sprint 5 | Event | GET /events/{id} | events | Planned |
| FR-012 | Configure Seat Layout | High | Seat Management | Sprint 6 | Seat | POST /events/{id}/seats | seats | Planned |
| FR-013 | View Available Seats | High | Seat Management | Sprint 6 | Seat | GET /events/{id}/seats | seats | Planned |
| FR-014 | Reserve Seat | Critical | Booking Engine | Sprint 7 | Booking | POST /bookings | bookings | Planned |
| FR-015 | Cancel Booking | High | Booking Engine | Sprint 7 | Booking | DELETE /bookings/{id} | bookings | Planned |
| FR-016 | Booking History | Medium | Booking Engine | Sprint 7 | Booking | GET /bookings | bookings | Planned |
| FR-017 | Mock Payment | High | Payment | Sprint 8 | Payment | POST /payments | payments | Planned |
| FR-018 | Booking Analytics | Medium | Analytics | Sprint 9 | Analytics | GET /analytics | bookings | Planned |