I want to create Go + Embedded HTMX DBMS web based

This app is connect to 1 instance DB only. Which means 1 URL DB only.

I want this only for postgres and stateless.

This is a minimalist, lightweight, web-based DBMS client application built using Go and HTMX, designed to connect to a single PostgreSQL database instance. The application will provide users with a simple interface to explore database schemas, view and edit table data, execute manual SQL queries, and visualize entity-relationship diagrams (ERDs). The application will implement role-based access control (RBAC) based on PostgreSQL roles and permissions, ensuring that users can only access the data they are authorized to view or modify.

Here is the stories:

## User Stories

### Story 1: Setup & Configuration
When run, it have "-db" with superadmin connection string (the database name usually "postgres") to connect to Postgres instance. The application:
- Validates the connection string and tests connectivity to the Postgres instance
- Fetches and caches (in-memory) metadata about available databases, schemas, tables, columns, relations, and PostgreSQL roles
- For each role, determines and caches which databases, schemas, and tables are accessible based on role permissions
- Initializes role-based access control (RBAC) mapping each role to their accessible resources
- Store database url and role-based metadata in memory for the rest of the application lifecycle

### Story 2: Authentication & Identity

When visiting the application, users authenticate with their PostgreSQL credentials. The application:
- Displays a minimalist login form page with username and password
- Validates credentials and probes connection to the user's first accessible database, schema, and table before creating session
- If connection probe fails (user has no accessible resources), login is rejected with error message
- If connection probe succeeds, creates session with cookies
- Stores identity in a long-lived cookie for UI pre-filling (username)
- Encrypts the session password in a short-lived cookie and re-authenticates on each request
- After successful probe and session creation, populates Data Explorer sidebar with user's accessible databases, schemas, and tables (based on their role permissions from Story 1)
- Redirects to Main View which displays data from the first accessible table (indicator of successful login)
- Allows users to switch between accessible databases they have `CONNECT` permissions for
- There will be a Header with:
    - Application title on the left
    - Username on the right
    - Beside username, there will be a "Logout" button that clears cookies and returns to login page
    - Nav-menu: Main View, Manual Query Editor, ERD Viewer
    - There will be refresh button to reload global metadata from the DBMS
    - Is under the Data Explorer sidebar
    
### Story 3: ERD Viewer
- There will be an "ERD Viewer" page accessible from the nav-menu in the header
- The ERD Viewer displays an Entity-Relationship Diagram of the current selected database schema
- The diagram is generated dynamically based on the cached metadata of tables and their relationships
- The diagram visually represents tables as boxes with columns listed inside
- Relationships (foreign keys) are shown as lines connecting related tables
- Users can zoom in/out and pan around the diagram
- Clicking on a table in the diagram highlights it and displays its details (columns, data types, constraints) in a side panel

### Story 4: Manual Query Editor
- There will be a "Manual Query Editor" page accessible from the nav-menu in the header
- There will be 2 Panels: Query Editor (top) and Results View (bottom)
- The Manual Query Editor provides a text area for users to write and execute custom SQL queries
- The editor supports syntax highlighting for SQL
- An "Execute" button runs the highlighted query / on cursor query and displays results below the editor. It can be more than 1 query, separated by semicolons
- Results are shown with offset pagination in a table format for SELECT queries
  - Display shows actual total result size (e.g., "Data size: 5000 rows")
  - Hard limit: Only first 1000 rows are loaded and displayed (navigable with previous/next buttons)
  - User can see how many total rows exist but cannot page beyond the 1000 row limit
- For DDL/DML queries, the editor displays affected row counts or success messages
- Error messages are displayed for invalid queries

### Story 5: Main View & Data Interaction
- There will be a "Main View" page accessible from the nav-menu in the header or after login (home page)
- Data Explorer now shows tables the user has access to (role-aware)
- Clicking a table in the Data Explorer loads its data into the Main View
- In Main View:
    - There will be A "Where Bar" for entering SQL WHERE clause fragments to filter results
    - Besides "Where Bar" there will be "Start Transaction" button
    - Under "Where Bar" there will be table data with cursor pagination limit 50 per page with infinite scrolling
      - Display shows actual total result size (e.g., "Data size: 5000 rows")
      - Hard limit: Only first 1000 rows are loaded via infinite scroll (20 pages × 50 rows)
      - User can see how many total rows exist but can only access up to 1000 rows
    - Table columns are sortable by clicking on the column headers
    - NOT IN TRANSACTION MODE:
        - Data is read-only
        - No inline editing
        - No commit/rollback buttons
        - Clicking a Foreign Key cell navigates to parent table
        - Clicking a Primary Key cell shows a modal/panel with list of referencing child tables and row counts for each (e.g., "orders: 5 rows", "invoices: 3 rows")
        - Clicking on a referencing table in the modal navigates to that table with WHERE clause filtering to show only related rows
    - IN TRANSACTION MODE:
        - Inline cell editing enabled
        - Edited cells are buffered but not committed
        - "Start Transaction" button changes to "Transaction Active" with countdown timer (1 minutes) and "Commit" and "Rollback" buttons appear
        - Delete row button appears on each row
        - New row button appears above the table

### Story 6: Isolation
- User A and User B log in simultaneously with different Postgres credentials
- Each user sees only the schemas, tables, and data they have permissions for
- Transactions initiated by User A do not affect User B's view or data, and vice versa
- Session cookies ensure that each user's identity and permissions are isolated

### Story 7: Security & Best Practices
- All database interactions use parameterized queries to prevent SQL injection
- Sensitive information (like passwords) is encrypted in cookies
- Sessions have appropriate timeouts to enhance security
- The application adheres to security best practices for web applications and database interactions

===

Give me Test Cases for Unit Testing and E2E Testing for all the stories above.