package domain

// Connection represents PostgreSQL connection configuration
type Connection struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
	SSLMode  string
}

// ConnectionString returns a formatted PostgreSQL connection string
func (c *Connection) ConnectionString() string {
	return "" // To be implemented
}
