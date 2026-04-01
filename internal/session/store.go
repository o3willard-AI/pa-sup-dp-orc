package session

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

// Store manages persistence of sessions and suggested commands.
type Store struct {
	db *sql.DB
}

// NewStore creates a new session store at the given path.
func NewStore(path string) (*Store, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	// Enable foreign key constraints (SQLite defaults to OFF)
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return nil, err
	}
	s := &Store{db: db}
	if err := s.createTables(); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) createTables() error {
	// Sessions table
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			terminal_id TEXT NOT NULL,
			created_at DATETIME NOT NULL
		)
	`)
	if err != nil {
		return err
	}
	// Commands table
	_, err = s.db.Exec(`
		CREATE TABLE IF NOT EXISTS commands (
			id TEXT PRIMARY KEY,
			session_id TEXT NOT NULL,
			terminal_id TEXT NOT NULL,
			command TEXT NOT NULL,
			description TEXT,
			context TEXT,
			created_at DATETIME NOT NULL,
			used_count INTEGER DEFAULT 0,
			last_used_at DATETIME,
			FOREIGN KEY (session_id) REFERENCES sessions(id) ON DELETE CASCADE
		)
	`)
	return err
}

// AddSession creates a new session.
func (s *Store) AddSession(sessionID, terminalID string) error {
	_, err := s.db.Exec(
		"INSERT INTO sessions (id, terminal_id, created_at) VALUES (?, ?, ?)",
		sessionID, terminalID, time.Now(),
	)
	return err
}

// AddCommand adds a suggested command to a session.
func (s *Store) AddCommand(cmd SuggestedCommand) error {
	_, err := s.db.Exec(
		`INSERT INTO commands 
		(id, session_id, terminal_id, command, description, context, created_at, used_count, last_used_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		cmd.ID, cmd.SessionID, cmd.TerminalID, cmd.Command, cmd.Description, cmd.Context,
		cmd.CreatedAt, cmd.UsedCount, cmd.LastUsedAt,
	)
	return err
}

// GetCommandsByTerminal returns all commands for a terminal.
func (s *Store) GetCommandsByTerminal(terminalID string) ([]SuggestedCommand, error) {
	rows, err := s.db.Query(`
		SELECT id, session_id, terminal_id, command, description, context, created_at, used_count, last_used_at
		FROM commands WHERE terminal_id = ?
		ORDER BY created_at DESC
	`, terminalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var commands []SuggestedCommand
	for rows.Next() {
		var cmd SuggestedCommand
		err := rows.Scan(&cmd.ID, &cmd.SessionID, &cmd.TerminalID, &cmd.Command, &cmd.Description, &cmd.Context,
			&cmd.CreatedAt, &cmd.UsedCount, &cmd.LastUsedAt)
		if err != nil {
			return nil, err
		}
		commands = append(commands, cmd)
	}
	return commands, nil
}

// IncrementUsedCount increments the used count of a command.
func (s *Store) IncrementUsedCount(commandID string) error {
	_, err := s.db.Exec(`
		UPDATE commands 
		SET used_count = used_count + 1, last_used_at = ?
		WHERE id = ?
	`, time.Now(), commandID)
	return err
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}
