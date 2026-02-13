# Test Plan for Lumen-PG

This document outlines the comprehensive test plan for implementing the Lumen-PG project using Test-Driven Development (TDD).

## Testing Strategy

### 1. Unit Tests
Unit tests verify individual components in isolation using mocks/stubs.

### 2. Integration Tests
Integration tests verify component interactions with real PostgreSQL database using testcontainers.

### 3. E2E Tests
End-to-end tests verify complete user flows through HTTP API with HTMX responses.

---

## Story 1: Setup & Configuration

### Unit Tests

#### UC-S1-01: Connection String Validation ✅
- **Given**: Invalid connection string format
- **When**: Application starts with `-db` flag
- **Then**: Returns validation error with clear message

#### UC-S1-02: Connection String Parsing ✅
- **Given**: Valid connection string with all parameters
- **When**: Connection repository parses the string
- **Then**: Extracts host, port, database, user, password correctly

#### UC-S1-03: Superadmin Connection Test Success ✅
- **Given**: Valid superadmin credentials
- **When**: Testing connectivity to PostgreSQL instance
- **Then**: Returns success without error

#### UC-S1-04: Superadmin Connection Test Failure ✅
- **Given**: Invalid superadmin credentials
- **When**: Testing connectivity to PostgreSQL instance
- **Then**: Returns connection error

#### UC-S1-05: Metadata Initialization - Roles and Permissions ✅
- **Given**: Valid superadmin connection
- **When**: Initialize metadata cache
- **Then**: Fetches all PostgreSQL roles and their permissions on databases, schemas, tables

#### UC-S1-06: In-Memory Metadata Storage - Per Role ✅
- **Given**: Fetched metadata with roles and permissions
- **When**: Storing in memory
- **Then**: Caches accessible databases, schemas, tables for each role

#### UC-S1-07: RBAC Initialization with User Accessibility ✅
- **Given**: PostgreSQL roles with varying permissions
- **When**: Initialize RBAC
- **Then**: Maps each role to accessible resources (databases, schemas, tables) and stores in memory

### Integration Tests

#### IT-S1-01: Connect to Real PostgreSQL ✅
- **Given**: Real PostgreSQL instance (testcontainers)
- **When**: Application starts with valid connection string
- **Then**: Successfully connects and initializes

#### IT-S1-02: Load Real Database Metadata with User Accessible Resources ✅
- **Given**: PostgreSQL with multiple databases, schemas, and roles with different permissions
- **When**: Loading global metadata during initialization
- **Then**: Returns all databases, schemas, tables, columns WITH associated role permissions

#### IT-S1-03: Load Real Relations and Role Access ✅
- **Given**: Tables with foreign key relationships and role-based access
- **When**: Loading metadata
- **Then**: Returns all foreign key relationships and which roles can access related tables

#### IT-S1-04: Cache Accessible Resources Per Role ✅
- **Given**: PostgreSQL with multiple roles (admin, editor, viewer)
- **When**: Loading metadata from roles
- **Then**: Caches accessible databases, schemas, and tables for each role separately in memory

---

## Story 2: Authentication & Identity

### Unit Tests

#### UC-S2-01: Login Form Validation - Empty Username ✅
- **Given**: Empty username
- **When**: Attempting login
- **Then**: Returns validation error

#### UC-S2-02: Login Form Validation - Empty Password ✅
- **Given**: Empty password
- **When**: Attempting login
- **Then**: Returns validation error

#### UC-S2-03: Login Connection Probe ✅
- **Given**: Valid PostgreSQL credentials
- **When**: User attempts login
- **Then**: Probes connection to first accessible database, schema, and table for that user

#### UC-S2-04: Login Connection Probe Failure ✅
- **Given**: User credentials but no accessible tables
- **When**: Probing connection to first accessible resource
- **Then**: Returns error (no accessible resources found)

#### UC-S2-05: Login Success After Probe ✅
- **Given**: Valid PostgreSQL credentials and accessible database/schema/table
- **When**: Connection probe succeeds
- **Then**: Creates session with cookies and returns first accessible table data

#### UC-S2-06: Session Cookie Creation - Username ✅
- **Given**: Successful login after probe
- **When**: Creating cookies
- **Then**: Sets long-lived cookie with username

#### UC-S2-07: Session Cookie Creation - Password ✅
- **Given**: Successful login after probe
- **When**: Creating cookies
- **Then**: Sets short-lived encrypted cookie with password

#### UC-S2-08: Session Validation - Valid Session ✅
- **Given**: Valid session cookie
- **When**: Validating session
- **Then**: Returns success

#### UC-S2-09: Session Validation - Expired Session ✅
- **Given**: Expired session cookie
- **When**: Validating session
- **Then**: Returns session expired error

#### UC-S2-10: Session Re-authentication ✅
- **Given**: Encrypted password cookie
- **When**: Each request
- **Then**: Re-authenticates with PostgreSQL

#### UC-S2-11: Data Explorer Population After Login ✅
- **Given**: Successful login with connection probe completed
- **When**: Session created and user redirected to Main View
- **Then**: Data Explorer sidebar populated with user's accessible databases, schemas, and tables

#### UC-S2-12: Logout Cookie Clearing ✅
- **Given**: Authenticated session
- **When**: User logs out
- **Then**: Clears both cookies

#### UC-S2-13: Header Username Display ✅
- **Given**: Authenticated session
- **When**: Rendering header
- **Then**: Displays username on the right

#### UC-S2-14: Navigation Menu Rendering ✅
- **Given**: Authenticated session
- **When**: Rendering header
- **Then**: Shows Main View, Manual Query Editor, ERD Viewer

#### UC-S2-15: Metadata Refresh Button ✅
- **Given**: Authenticated session with superadmin rights
- **When**: Clicking refresh button
- **Then**: Reloads global metadata from DBMS

### Integration Tests

#### IT-S2-01: Real PostgreSQL Connection Probe ✅
- **Given**: Real PostgreSQL user with accessible resources
- **When**: Logging in with valid credentials
- **Then**: Successfully probes connection to first accessible database, schema, and table

#### IT-S2-02: Real PostgreSQL Connection Probe Failure ✅
- **Given**: Real PostgreSQL user with no accessible resources
- **When**: Logging in with valid credentials
- **Then**: Connection probe fails, login rejected

#### IT-S2-03: Real Role-Based Resource Access ✅
- **Given**: Multiple PostgreSQL users with different role permissions
- **When**: Each user logs in
- **Then**: Returns only databases, schemas, tables accessible to that user

#### IT-S2-04: Session Persistence After Probe ✅
- **Given**: Logged in user after successful probe
- **When**: Making multiple requests
- **Then**: Session persists and accessible resources remain cached

#### IT-S2-05: Concurrent User Sessions with Isolated Resources ✅
- **Given**: Multiple users with different credentials and permissions
- **When**: Logging in simultaneously
- **Then**: Each has isolated session with separate Data Explorer populated for their accessible resources

### E2E Tests

#### E2E-S2-01: Login Flow with Connection Probe
- **Given**: User visits application with valid credentials
- **When**: Submitting login form
- **Then**: System probes first accessible resource, creates session, populates Data Explorer sidebar, redirects to Main View showing first accessible table (success indicator)

#### E2E-S2-02: Login Flow - No Accessible Resources
- **Given**: User visits application with valid credentials but no accessible tables
- **When**: Submitting login form
- **Then**: Shows error message "No accessible resources found"

#### E2E-S2-03: Login Flow - Invalid Credentials
- **Given**: User visits application
- **When**: Submitting login form with invalid credentials
- **Then**: Shows error message on login page

#### E2E-S2-04: Logout Flow
- **Given**: Authenticated user
- **When**: Clicking logout button
- **Then**: Redirects to login page and clears session

#### E2E-S2-05: Protected Route Access Without Auth
- **Given**: Unauthenticated user
- **When**: Accessing protected route directly
- **Then**: Redirects to login page

#### E2E-S2-06: Data Explorer Populated After Login
- **Given**: User successfully logged in and redirected to Main View
- **When**: Main page loads
- **Then**: Data Explorer sidebar is populated with user's accessible databases, schemas, and tables; Main View displays first accessible table data

---

## Story 3: ERD Viewer

### Unit Tests

#### UC-S3-01: ERD Data Generation
- **Given**: Cached metadata with tables and relations
- **When**: Generating ERD data
- **Then**: Returns JSON with tables and relationships

#### UC-S3-02: Table Box Representation
- **Given**: Table metadata
- **When**: Generating ERD
- **Then**: Includes table name, columns, data types

#### UC-S3-03: Relationship Lines
- **Given**: Foreign key relationships
- **When**: Generating ERD
- **Then**: Includes lines between related tables

#### UC-S3-04: Empty Schema ERD
- **Given**: Schema with no tables
- **When**: Generating ERD
- **Then**: Returns empty ERD

### Integration Tests

#### IT-S3-01: ERD from Real Schema
- **Given**: Real PostgreSQL schema with relationships
- **When**: Loading ERD
- **Then**: Returns accurate diagram data

#### IT-S3-02: Complex Relationships
- **Given**: Tables with multiple foreign keys
- **When**: Loading ERD
- **Then**: Shows all relationships correctly

### E2E Tests

#### E2E-S3-01: ERD Viewer Page Access
- **Given**: Authenticated user
- **When**: Clicking ERD Viewer in nav-menu
- **Then**: Displays ERD page with diagram

#### E2E-S3-02: ERD Zoom Controls
- **Given**: ERD viewer page
- **When**: Using zoom controls
- **Then**: Diagram zooms in/out

#### E2E-S3-03: ERD Pan
- **Given**: ERD viewer page
- **When**: Dragging diagram
- **Then**: Diagram pans around

#### E2E-S3-04: Table Click in ERD
- **Given**: ERD viewer page
- **When**: Clicking a table
- **Then**: Shows table details in side panel

---

## Story 4: Manual Query Editor

### Unit Tests

#### UC-S4-01: Single Query Execution ✅
- **Given**: Valid SELECT query
- **When**: Executing query
- **Then**: Returns result rows

#### UC-S4-02: Multiple Query Execution ✅
- **Given**: Multiple queries separated by semicolons
- **When**: Executing queries
- **Then**: Executes all queries in sequence

#### UC-S4-03: Query Result Offset Pagination ✅
- **Given**: SELECT query returning many rows
- **When**: Executing query
- **Then**: Returns first page (1000 rows max) with offset pagination

#### UC-S4-03a: Query Result Actual Size Display ✅
- **Given**: SELECT query returning 5000 rows
- **When**: Executing query
- **Then**: Shows "Data size: 5000 rows" indicator while only loading first 1000 rows

#### UC-S4-03b: Query Result Limit Hard Cap ✅
- **Given**: SELECT query returning more than 1000 rows
- **When**: Executing query
- **Then**: Hard limits to 1000 rows maximum (no more pages beyond this)

#### UC-S4-03c: Offset Pagination Next Page ✅
- **Given**: Query results with 1000 rows on first page
- **When**: Requesting next page with offset
- **Then**: Returns empty or message that limit reached

#### UC-S4-04: DDL Query Execution ✅
- **Given**: CREATE TABLE query
- **When**: Executing query
- **Then**: Returns success message

#### UC-S4-05: DML Query Execution ✅
- **Given**: INSERT query
- **When**: Executing query
- **Then**: Returns affected row count

#### UC-S4-06: Invalid Query Error ✅
- **Given**: Syntax error in query
- **When**: Executing query
- **Then**: Returns error message

#### UC-S4-07: Query Splitting ✅
- **Given**: SQL with semicolons in strings
- **When**: Splitting queries
- **Then**: Correctly identifies query boundaries

#### UC-S4-08: Parameterized Query Execution ✅
- **Given**: Query with parameters
- **When**: Executing query
- **Then**: Uses parameterized execution (no SQL injection)

### Integration Tests

#### IT-S4-01: Real SELECT Query ✅
- **Given**: Real table with data
- **When**: Executing SELECT query
- **Then**: Returns actual rows

#### IT-S4-02: Real DDL Query ✅
- **Given**: Real PostgreSQL connection
- **When**: Executing CREATE TABLE
- **Then**: Creates table in database

#### IT-S4-03: Real DML Query ✅
- **Given**: Real table
- **When**: Executing INSERT/UPDATE/DELETE
- **Then**: Modifies data and returns count

#### IT-S4-04: Query with Permission Denied ✅
- **Given**: User without table access
- **When**: Executing query on restricted table
- **Then**: Returns permission denied error

### E2E Tests

#### E2E-S4-01: Query Editor Page Access
- **Given**: Authenticated user
- **When**: Clicking Manual Query Editor in nav-menu
- **Then**: Displays query editor page

#### E2E-S4-02: Execute Single Query
- **Given**: Query editor with SELECT query
- **When**: Clicking Execute button
- **Then**: Displays results in results panel

#### E2E-S4-03: Execute Multiple Queries
- **Given**: Query editor with multiple queries
- **When**: Clicking Execute button
- **Then**: Displays all results in sequence

#### E2E-S4-04: Query Error Display
- **Given**: Query editor with invalid query
- **When**: Clicking Execute button
- **Then**: Displays error message in red

#### E2E-S4-05: Offset Pagination Results
- **Given**: Query with large result set (5000 rows)
- **When**: Viewing results
- **Then**: Shows first page (up to 1000 rows) with previous/next buttons and "Data size: 5000 rows" indicator

#### E2E-S4-05a: Offset Pagination Navigation
- **Given**: Query results displayed on first page with 1000 rows loaded
- **When**: Clicking next page button
- **Then**: Shows message that hard limit of 1000 rows reached (no additional pages beyond this)

#### E2E-S4-05b: Query Result Actual Size vs Display Limit
- **Given**: Query returning more than 1000 rows (e.g., 10000 rows total)
- **When**: Viewing results
- **Then**: Shows "Data size: 10000 rows" but displays maximum 1000 rows with indication that only first 1000 are accessible

#### E2E-S4-06: SQL Syntax Highlighting
- **Given**: Query editor with SQL
- **When**: Viewing editor
- **Then**: Keywords are syntax highlighted

---

## Story 5: Main View & Data Interaction

### Unit Tests

#### UC-S5-01: Table Data Loading ✅
- **Given**: Valid table name
- **When**: Loading table data
- **Then**: Returns first 50 rows

#### UC-S5-02: Cursor Pagination Next Page ✅
- **Given**: Table with 100 rows, user has loaded first 50
- **When**: Loading next page with cursor (infinite scroll)
- **Then**: Returns next 50 rows with new cursor

#### UC-S5-03: WHERE Clause Validation ✅
- **Given**: Valid WHERE clause fragment
- **When**: Validating clause
- **Then**: Returns success

#### UC-S5-04: WHERE Clause Injection Prevention ✅
- **Given**: WHERE clause with SQL injection attempt
- **When**: Validating clause
- **Then**: Returns validation error or uses parameterization

#### UC-S5-05: Column Sorting ASC ✅
- **Given**: Table data
- **When**: Sorting by column ascending
- **Then**: Returns sorted data

#### UC-S5-06: Column Sorting DESC ✅
- **Given**: Table data
- **When**: Sorting by column descending
- **Then**: Returns sorted data

#### UC-S5-07: Cursor Pagination Actual Size Display ✅
- **Given**: Table with 5000 total rows
- **When**: Loading table data
- **Then**: Shows "Data size: 5000 rows" indicator while only loading up to 1000 rows

#### UC-S5-08: Cursor Pagination Hard Limit ✅
- **Given**: Table with 1000+ rows loaded via infinite scroll
- **When**: Reaching hard limit of 1000 rows
- **Then**: Stops loading more data and shows indication that limit reached

#### UC-S5-09: Transaction Start ✅
- **Given**: User session without active transaction
- **When**: Starting transaction
- **Then**: Transaction becomes active with 1-minute timer

#### UC-S5-10: Transaction Already Active Error ✅
- **Given**: User session with active transaction
- **When**: Attempting to start another transaction
- **Then**: Returns error

#### UC-S5-11: Cell Edit Buffering ✅
- **Given**: Active transaction
- **When**: Editing cell value
- **Then**: Operation added to buffer, not committed

#### UC-S5-12: Transaction Commit ✅
- **Given**: Active transaction with buffered operations
- **When**: Committing transaction
- **Then**: Executes all operations atomically

#### UC-S5-13: Transaction Rollback ✅
- **Given**: Active transaction with buffered operations
- **When**: Rolling back transaction
- **Then**: Discards all buffered operations

#### UC-S5-14: Transaction Timer Expiration ✅
- **Given**: Active transaction 1 minute old
- **When**: Timer expires
- **Then**: Transaction automatically rolls back

#### UC-S5-15: Row Deletion Buffering ✅
- **Given**: Active transaction
- **When**: Deleting row
- **Then**: DELETE operation added to buffer

#### UC-S5-16: Row Insertion Buffering ✅
- **Given**: Active transaction
- **When**: Inserting new row
- **Then**: INSERT operation added to buffer

#### UC-S5-17: Foreign Key Navigation ✅
- **Given**: Cell with foreign key value
- **When**: Clicking cell (not in transaction mode)
- **Then**: Navigates to parent table

#### UC-S5-18: Primary Key Navigation ✅
- **Given**: Cell with primary key value
- **When**: Clicking cell (not in transaction mode)
- **Then**: Shows modal/panel with list of referencing tables and row counts for each

#### UC-S5-19: Read-Only Mode Enforcement ✅
- **Given**: Not in transaction mode
- **When**: Attempting to edit cell
- **Then**: Editing is disabled

### Integration Tests

#### IT-S5-01: Real Table Data Loading
- **Given**: Real table with data
- **When**: Loading main view
- **Then**: Returns actual data

#### IT-S5-02: Real Cursor Pagination
- **Given**: Real table with 100+ rows
- **When**: Paginating through data
- **Then**: Loads all pages correctly

#### IT-S5-03: Real WHERE Filter
- **Given**: Real table with data
- **When**: Applying WHERE clause
- **Then**: Returns filtered results

#### IT-S5-04: Real Transaction Commit ✅
- **Given**: Real table
- **When**: Starting transaction, editing cells, committing
- **Then**: Changes persisted to database

#### IT-S5-05: Real Transaction Rollback ✅
- **Given**: Real table
- **When**: Starting transaction, editing cells, rolling back
- **Then**: No changes persisted to database

#### IT-S5-06: Real Foreign Key Navigation ✅
- **Given**: Tables with foreign key relationships
- **When**: Clicking FK cell
- **Then**: Loads parent table data

#### IT-S5-07: Real Primary Key Navigation ✅
- **Given**: Tables with FK relationships
- **When**: Fetching referencing tables for a PK value
- **Then**: Returns list of referencing tables with accurate row counts for that FK relationship

### E2E Tests

#### E2E-S5-01: Main View Default Load
- **Given**: User logs in
- **When**: Login completes
- **Then**: Main view shows first accessible table

#### E2E-S5-02: Table Selection from Sidebar
- **Given**: Main view with sidebar
- **When**: Clicking table in sidebar
- **Then**: Loads selected table data

#### E2E-S5-03: WHERE Bar Filtering
- **Given**: Main view with table data
- **When**: Entering WHERE clause and submitting
- **Then**: Table data updates with filtered results

#### E2E-S5-04: Column Header Sorting
- **Given**: Main view with table data
- **When**: Clicking column header
- **Then**: Sorts data by that column

#### E2E-S5-05: Cursor Pagination Infinite Scroll with Actual Size
- **Given**: Main view with table data (5000 total rows, 50 rows per page, hard limit 1000 rows)
- **When**: Viewing table
- **Then**: Shows "Data size: 5000 rows" and loads first 50 rows automatically

#### E2E-S5-05a: Cursor Pagination Infinite Scroll Loading
- **Given**: Main view with first page loaded (50 rows shown)
- **When**: Scrolling to bottom of current page
- **Then**: Loads next 50 rows automatically via cursor pagination

#### E2E-S5-05b: Pagination Hard Limit Enforcement
- **Given**: Main view loaded with 1000 rows (20 pages × 50 rows)
- **When**: User scrolls to bottom after reaching hard limit
- **Then**: No more data loads and shows "Reached limit of 1000 rows. Total available: 5000 rows"

#### E2E-S5-06: Start Transaction Button
- **Given**: Main view not in transaction mode
- **When**: Clicking "Start Transaction"
- **Then**: Button changes to "Transaction Active" with timer

#### E2E-S5-07: Transaction Mode Cell Editing
- **Given**: Main view in transaction mode
- **When**: Clicking cell
- **Then**: Cell becomes editable

#### E2E-S5-08: Transaction Mode Edit Buffer Display
- **Given**: Transaction mode with edits
- **When**: Making changes
- **Then**: Buffered operations displayed somewhere

#### E2E-S5-09: Transaction Commit Button
- **Given**: Transaction mode with changes
- **When**: Clicking Commit
- **Then**: Changes saved, returns to read-only mode

#### E2E-S5-10: Transaction Rollback Button
- **Given**: Transaction mode with changes
- **When**: Clicking Rollback
- **Then**: Changes discarded, returns to read-only mode

#### E2E-S5-11: Transaction Timer Countdown
- **Given**: Active transaction
- **When**: Viewing page
- **Then**: Timer counts down from 60 seconds

#### E2E-S5-12: Transaction Row Delete Button
- **Given**: Transaction mode
- **When**: Clicking delete button on row
- **Then**: Row marked for deletion (buffered)

#### E2E-S5-13: Transaction New Row Button
- **Given**: Transaction mode
- **When**: Clicking "New Row" button
- **Then**: Empty editable row appears

#### E2E-S5-14: FK Cell Navigation (Read-Only)
- **Given**: Not in transaction mode
- **When**: Clicking FK cell
- **Then**: Navigates to parent table

#### E2E-S5-15: PK Cell Navigation (Read-Only)
- **Given**: Not in transaction mode, viewing table with PK that has foreign key references
- **When**: Clicking PK cell
- **Then**: Shows modal/panel with list of referencing tables and row counts (e.g., "orders: 5 rows", "invoices: 3 rows")

#### E2E-S5-15a: PK Cell Navigation - Table Click
- **Given**: Modal/panel open showing referencing tables
- **When**: Clicking on a referencing table in the modal
- **Then**: Main view navigates to that table with WHERE clause filtering to show only related rows

---

## Story 6: Isolation

### Unit Tests

#### UC-S6-01: Session Isolation
- **Given**: Two different user sessions
- **When**: Operating simultaneously
- **Then**: Sessions do not interfere with each other

#### UC-S6-02: Transaction Isolation ✅
- **Given**: User A and User B both have active transactions
- **When**: Both make changes
- **Then**: Changes are isolated per user

#### UC-S6-03: Cookie Isolation ✅
- **Given**: Two different browsers/sessions
- **When**: Each has different cookies
- **Then**: Each sees their own session data

### Integration Tests

#### IT-S6-01: Real Multi-User Connection
- **Given**: Two PostgreSQL users with different permissions
- **When**: Both log in simultaneously
- **Then**: Each connects with their own credentials

#### IT-S6-02: Real Permission Isolation
- **Given**: User A can see table X, User B cannot
- **When**: Both are logged in
- **Then**: User A sees table X, User B does not

#### IT-S6-03: Real Transaction Isolation ✅
- **Given**: User A starts transaction, edits data
- **When**: User B views same table
- **Then**: User B does not see User A's uncommitted changes

### E2E Tests

#### E2E-S6-01: Simultaneous Users Different Permissions
- **Given**: Two browser sessions logged in as different users
- **When**: Both navigate application
- **Then**: Each sees only their permitted data

#### E2E-S6-02: Simultaneous Transactions
- **Given**: Two users both start transactions on same table
- **When**: Both make edits
- **Then**: Changes are isolated until commit

#### E2E-S6-03: One User Cannot See Another's Session
- **Given**: User A logged in
- **When**: User B logs in with different credentials
- **Then**: User B does not see User A's session data

---

## Story 7: Security & Best Practices

### Unit Tests

#### UC-S7-01: SQL Injection Prevention - WHERE Clause ✅
- **Given**: WHERE clause with SQL injection attempt
- **When**: Executing query
- **Then**: Uses parameterized query or safely escapes

#### UC-S7-02: SQL Injection Prevention - Query Editor ✅
- **Given**: Query with injection attempt
- **When**: Executing query
- **Then**: Uses parameterized execution

#### UC-S7-03: Password Encryption in Cookie ✅
- **Given**: User password
- **When**: Storing in cookie
- **Then**: Password is encrypted

#### UC-S7-04: Password Decryption from Cookie ✅
- **Given**: Encrypted password cookie
- **When**: Reading cookie
- **Then**: Password is correctly decrypted

#### UC-S7-05: Cookie Tampering Detection
- **Given**: Tampered cookie
- **When**: Validating cookie
- **Then**: Rejects cookie as invalid

#### UC-S7-06: Session Timeout Short-Lived Cookie ✅
- **Given**: Session cookie created
- **When**: Time passes beyond timeout
- **Then**: Cookie expires automatically

#### UC-S7-07: Session Timeout Long-Lived Cookie ✅
- **Given**: Username cookie created
- **When**: Time passes
- **Then**: Cookie remains for long duration

### Integration Tests

#### IT-S7-01: Real SQL Injection Test ✅
- **Given**: Real PostgreSQL connection
- **When**: Attempting SQL injection in WHERE clause
- **Then**: Injection is prevented

#### IT-S7-02: Real Password Security ✅
- **Given**: Real password in cookie
- **When**: Inspecting cookie contents
- **Then**: Password is not visible in plaintext

#### IT-S7-03: Real Session Expiration ✅
- **Given**: Real session created
- **When**: Session timeout expires
- **Then**: Session is invalid and requires re-login

### E2E Tests

#### E2E-S7-01: SQL Injection via WHERE Bar ✅
- **Given**: Main view WHERE bar
- **When**: Entering SQL injection attempt
- **Then**: Injection is prevented, no unauthorized access

#### E2E-S7-02: SQL Injection via Query Editor ✅
- **Given**: Query editor
- **When**: Entering queries with injection attempts
- **Then**: Queries execute safely without escalation

#### E2E-S7-03: Cookie Tampering Prevention ✅
- **Given**: Authenticated session
- **When**: Manually tampering with cookie value
- **Then**: Session becomes invalid, requires re-login

#### E2E-S7-04: Session Timeout Enforcement ✅
- **Given**: Authenticated session
- **When**: Waiting beyond session timeout
- **Then**: Next request redirects to login

#### E2E-S7-05: HTTPS-Only Cookies (if HTTPS enabled)
- **Given**: Application running with HTTPS
- **When**: Creating cookies
- **Then**: Cookies have Secure flag set

#### E2E-S7-06: HTTPOnly Cookies ✅
- **Given**: Authentication cookies
- **When**: Creating cookies
- **Then**: Cookies have HTTPOnly flag set

---

## Test Execution Order (TDD)

### Phase 1: Core Domain Models
1. Unit tests for Transaction domain 
2. Unit tests for Session domain 
3. Unit tests for Cursor domain 
4. Unit tests for Query Splitter 
5. Unit tests for WHERE Validator 

### Phase 2: Repository Interfaces
1. Unit tests for Connection Repository interface
2. Unit tests for Metadata Repository interface
3. Unit tests for Query Repository interface

### Phase 3: Use Cases
1. Unit tests for Auth Usecase 
2. Unit tests for Metadata Usecase 
3. Unit tests for Data Explorer Usecase 
4. Unit tests for Query Usecase 
5. Unit tests for Transaction Usecase 

### Phase 4: Repository Implementations
1. Integration tests for PostgreSQL Connection Repository
2. Integration tests for PostgreSQL Metadata Repository
3. Integration tests for PostgreSQL Query Repository

### Phase 5: Web Layer
1. Unit tests for Session Manager middleware
2. Unit tests for Auth middleware
3. Integration tests for HTTP handlers

### Phase 6: E2E Tests
1. E2E tests for authentication flow
2. E2E tests for main view interactions
3. E2E tests for query editor
4. E2E tests for ERD viewer
5. E2E tests for transaction management
6. E2E tests for multi-user isolation
7. E2E tests for security features

---

## Test Coverage Goals

- **Unit Tests**: 90% coverage for domain logic and use cases
- **Integration Tests**: All database operations tested against real PostgreSQL
- **E2E Tests**: All user stories covered with happy path and error cases

---

## Test Data Setup

### Test Databases
- **testdb1**: Contains schema `public` with tables `users`, `posts`, `comments`
- **testdb2**: Contains schema `public` with tables `products`, `orders`

### Test Users
- **superuser**: Full access to all databases
- **user_readonly**: Read-only access to testdb1
- **user_limited**: Access to testdb1.public.users only
- **user_multi**: Access to both testdb1 and testdb2

### Test Tables

#### users (testdb1.public)
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100),
    created_at TIMESTAMP DEFAULT NOW()
);
```

#### posts (testdb1.public)
```sql
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    title VARCHAR(200),
    content TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

#### comments (testdb1.public)
```sql
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    post_id INTEGER REFERENCES posts(id),
    user_id INTEGER REFERENCES users(id),
    comment_text TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);
```

#### products (testdb2.public)
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100),
    price DECIMAL(10,2),
    stock INTEGER
);
```

#### orders (testdb2.public)
```sql
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER,
    ordered_at TIMESTAMP DEFAULT NOW()
);
```

---

## Continuous Testing

### Run All Tests
```bash
go test ./... -v
```

### Run Unit Tests Only
```bash
go test ./internal/domain/... ./internal/usecase/... -v
```

### Run Integration Tests Only
```bash
go test ./test/integration/... -v
```

### Run with Coverage
```bash
go test ./... -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Run Specific Story Tests
```bash
# Story 2 - Authentication
go test ./internal/usecase -run TestAuth -v
go test ./test/integration -run TestAuth -v

# Story 4 - Query Editor
go test ./internal/usecase -run TestQuery -v
go test ./test/integration -run TestQuery -v

# Story 5 - Main View & Transactions
go test ./internal/usecase -run TestDataExplorer -v
go test ./internal/usecase -run TestTransaction -v
go test ./test/integration -run TestTransaction -v
```

---

## Notes

- All tests must be written BEFORE implementing functionality (TDD)
- Each test should be independent and not rely on other tests
- Integration tests use testcontainers for real PostgreSQL instances
- E2E tests interact via HTTP endpoints and verify HTMX responses
- Mock repositories are used for unit tests
- Tests should clean up after themselves (no state pollution)
