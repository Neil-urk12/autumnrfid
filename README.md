# AutumnRFID System

## Overview
AutumnRFID is a student information system that uses RFID card scanning to display student data in real-time. The system provides a web interface for viewing student information, grades, and billing details when an RFID card is scanned.

## Features
- Real-time RFID Card Scanning: Instantly displays student information when an RFID card is scanned
- Student Profile View: Shows comprehensive student information including:
  - Personal details (name, ID, contact information)
  - Academic information (program, year level, block)
  - System access history
- Grades Summary: Displays academic performance across semesters and years
- Billing Information: Shows tuition fees, payment status, and remaining balances
- Server-Sent Events (SSE): Enables real-time updates without page refreshes

## Technology Stack
- Backend: Go with Fiber web framework
- Frontend: HTML, CSS, JavaScript with HTMX for dynamic content
- Database: SQL database (specific implementation details not shown)
- Real-time Updates: Server-Sent Events (SSE)

## Project Structure
```
autumnrfid/
├── cmd/
│   └── web/
│       └── main.go         # Application entry point
├── internal/
│   ├── config/             # Configuration management
│   ├── handlers/           # HTTP request handlers
│   └── repositories/       # Database access layer
└── ui/
    ├── html/
    │   └── pages/
    │       └── home.html   # Main application interface
    └── static/
        ├── images/         # Static images including profile pictures
        └── styles.css      # Application styling
```

## How It Works
1. The system listens for RFID card scans via the /card-scan endpoint
2. When a card is scanned, the system queries the database for student information
3. Student data is sent to the frontend via Server-Sent Events (SSE)
4. The frontend updates in real-time to display the student's information

## API Endpoints
- GET /: Main application interface
- GET /grades: View detailed grades
- GET /test-grades: Test endpoint for grades display
- GET /error: Error page display
- GET /stream: SSE endpoint for real-time updates
- POST /card-scan: Endpoint for receiving RFID card scan data
- GET /ping: Health check endpoint

## Development
The application includes development-friendly features:

- Template reloading for quick UI changes
- Debug mode for better error messages
- Automatic static folder creation

## Testing
The system includes test buttons in the UI for simulating:

- Student data updates
- 404 Not Found errors

## License


## Contributors
**Michi Bichi**
**Toma Toma**
**Manuk Cano**
