package domain_test

import (
	"testing"

	"github.com/kamil5b/lumen-pg/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// UC-S4-07: Query Splitting - Simple
func TestSplitQueries_Single(t *testing.T) {
	queries := domain.SplitQueries("SELECT * FROM users")
	require.Len(t, queries, 1)
	assert.Equal(t, "SELECT * FROM users", queries[0])
}

// UC-S4-02/07: Multiple Query Execution
func TestSplitQueries_Multiple(t *testing.T) {
	queries := domain.SplitQueries("SELECT 1; SELECT 2; SELECT 3")
	require.Len(t, queries, 3)
	assert.Equal(t, "SELECT 1", queries[0])
	assert.Equal(t, "SELECT 2", queries[1])
	assert.Equal(t, "SELECT 3", queries[2])
}

// UC-S4-07: Query Splitting - Semicolons in Strings
func TestSplitQueries_SemicolonInSingleQuotedString(t *testing.T) {
	queries := domain.SplitQueries("SELECT 'hello;world'; SELECT 2")
	require.Len(t, queries, 2)
	assert.Equal(t, "SELECT 'hello;world'", queries[0])
	assert.Equal(t, "SELECT 2", queries[1])
}

func TestSplitQueries_SemicolonInDoubleQuotedString(t *testing.T) {
	queries := domain.SplitQueries(`SELECT "col;name" FROM t; SELECT 2`)
	require.Len(t, queries, 2)
	assert.Equal(t, `SELECT "col;name" FROM t`, queries[0])
	assert.Equal(t, "SELECT 2", queries[1])
}

func TestSplitQueries_Empty(t *testing.T) {
	queries := domain.SplitQueries("")
	assert.Empty(t, queries)
}

func TestSplitQueries_WhitespaceOnly(t *testing.T) {
	queries := domain.SplitQueries("   ;  ;  ")
	assert.Empty(t, queries)
}

func TestSplitQueries_TrailingSemicolon(t *testing.T) {
	queries := domain.SplitQueries("SELECT 1;")
	require.Len(t, queries, 1)
	assert.Equal(t, "SELECT 1", queries[0])
}

// UC-S5-03: WHERE Clause Validation
func TestValidateWhereClause_Valid(t *testing.T) {
	err := domain.ValidateWhereClause("id = 1")
	assert.NoError(t, err)
}

func TestValidateWhereClause_Empty(t *testing.T) {
	err := domain.ValidateWhereClause("")
	assert.NoError(t, err)
}

func TestValidateWhereClause_ComplexValid(t *testing.T) {
	err := domain.ValidateWhereClause("name = 'John' AND age > 25")
	assert.NoError(t, err)
}

// UC-S5-04: WHERE Clause Injection Prevention
func TestValidateWhereClause_InjectionDrop(t *testing.T) {
	err := domain.ValidateWhereClause("1=1; DROP TABLE users")
	assert.ErrorIs(t, err, domain.ErrSQLInjectionDetected)
}

func TestValidateWhereClause_InjectionAlter(t *testing.T) {
	err := domain.ValidateWhereClause("1=1; ALTER TABLE users")
	assert.ErrorIs(t, err, domain.ErrSQLInjectionDetected)
}

func TestValidateWhereClause_InjectionComment(t *testing.T) {
	err := domain.ValidateWhereClause("1=1 -- comment")
	assert.ErrorIs(t, err, domain.ErrSQLInjectionDetected)
}

func TestValidateWhereClause_InjectionBlockComment(t *testing.T) {
	err := domain.ValidateWhereClause("1=1 /* comment */")
	assert.ErrorIs(t, err, domain.ErrSQLInjectionDetected)
}

func TestValidateWhereClause_InjectionDelete(t *testing.T) {
	err := domain.ValidateWhereClause("1=1; DELETE FROM users")
	assert.ErrorIs(t, err, domain.ErrSQLInjectionDetected)
}

func TestValidateWhereClause_InjectionInsert(t *testing.T) {
	err := domain.ValidateWhereClause("1=1; INSERT INTO users VALUES (1)")
	assert.ErrorIs(t, err, domain.ErrSQLInjectionDetected)
}

func TestValidateWhereClause_InjectionUpdate(t *testing.T) {
	err := domain.ValidateWhereClause("1=1; UPDATE users SET name='x'")
	assert.ErrorIs(t, err, domain.ErrSQLInjectionDetected)
}

// UC-S7-01: SQL Injection Prevention - WHERE Clause
func TestValidateWhereClause_InjectionGrant(t *testing.T) {
	err := domain.ValidateWhereClause("1=1; GRANT ALL ON users TO public")
	assert.ErrorIs(t, err, domain.ErrSQLInjectionDetected)
}

func TestValidateWhereClause_InjectionRevoke(t *testing.T) {
	err := domain.ValidateWhereClause("1=1; REVOKE ALL ON users FROM public")
	assert.ErrorIs(t, err, domain.ErrSQLInjectionDetected)
}

func TestValidateWhereClause_InjectionCreate(t *testing.T) {
	err := domain.ValidateWhereClause("1=1; CREATE TABLE evil (id INT)")
	assert.ErrorIs(t, err, domain.ErrSQLInjectionDetected)
}

func TestValidateWhereClause_InjectionTruncate(t *testing.T) {
	err := domain.ValidateWhereClause("1=1; TRUNCATE TABLE users")
	assert.ErrorIs(t, err, domain.ErrSQLInjectionDetected)
}
