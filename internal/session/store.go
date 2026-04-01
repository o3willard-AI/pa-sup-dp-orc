package session

import (
	"database/sql"
	"fmt"
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
		return nil, fmt.Errorf("open database: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	// Enable foreign key constraints (SQLite defaults to OFF)
	if _, err := db.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		db.Close()
		return nil, fmt.Errorf("enable foreign keys: %w", err)
	}
	s := &Store{db: db}
	if err := s.createTables(); err != nil {
		db.Close()
		return nil, fmt.Errorf("create tables: %w", err)
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
	if err != nil {
		return err
	}
	// Indexes for performance
	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_commands_terminal_id ON commands(terminal_id)`)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_commands_session_id ON commands(session_id)`)
	if err != nil {
		return err
	}
	_, err = s.db.Exec(`CREATE INDEX IF NOT EXISTS idx_sessions_terminal_id ON sessions(terminal_id)`)
	if err != nil {
		return err
	}
	return nil
}

// AddSession creates a new session.
func (s *Store) AddSession(sessionID, terminalID string) error {
	_, err := s.db.Exec(
		"INSERT INTO sessions (id, terminal_id, created_at) VALUES (?, ?, ?)",
		sessionID, terminalID, time.Now(),
	)
	if err != nil {
		return fmt.Errorf("add session %q: %w", sessionID, err)
	}
	return nil
}

// GetSession retrieves a session by ID.
func (s *Store) GetSession(sessionID string) (*Session, error) {
	var sess Session
	err := s.db.QueryRow("SELECT id, terminal_id, created_at FROM sessions WHERE id = ?", sessionID).Scan(&sess.ID, &sess.TerminalID, &sess.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("get session %q: %w", sessionID, err)
	}
	return &sess, nil
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
	if err != nil {
		return fmt.Errorf("add command %q: %w", cmd.ID, err)
	}
	return nil
}

// GetCommandsByTerminal returns all commands for a terminal.
func (s *Store) GetCommandsByTerminal(terminalID string) ([]SuggestedCommand, error) {
	rows, err := s.db.Query(`
		SELECT id, session_id, terminal_id, command, description, context, created_at, used_count, last_used_at
		FROM commands WHERE terminal_id = ?
		ORDER BY created_at DESC
	`, terminalID)
	if err != nil {
		return nil, fmt.Errorf("query commands for terminal %q: %w", terminalID, err)
	}
	defer rows.Close()

	var commands []SuggestedCommand
	for rows.Next() {
		var cmd SuggestedCommand
		err := rows.Scan(&cmd.ID, &cmd.SessionID, &cmd.TerminalID, &cmd.Command, &cmd.Description, &cmd.Context,
			&cmd.CreatedAt, &cmd.UsedCount, &cmd.LastUsedAt)
		if err != nil {
			return nil, fmt.Errorf("scan command row: %w", err)
		}
		commands = append(commands, cmd)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w", err)
	}
	return commands, nil
}

// GetCommandByID retrieves a command by ID.
func (s *Store) GetCommandByID(commandID string) (*SuggestedCommand, error) {
	var cmd SuggestedCommand
	err := s.db.QueryRow("SELECT id, session_id, terminal_id, command, description, context, created_at, used_count, last_used_at FROM commands WHERE id = ?", commandID).Scan(&cmd.ID, &cmd.SessionID, &cmd.TerminalID, &cmd.Command, &cmd.Description, &cmd.Context, &cmd.CreatedAt, &cmd.UsedCount, &cmd.LastUsedAt)
	if err != nil {
		return nil, fmt.Errorf("get command %q: %w", commandID, err)
	}
	return &cmd, nil
}

// IncrementUsedCount increments the used count of a command.
func (s *Store) IncrementUsedCount(commandID string) error {
	_, err := s.db.Exec(`
		UPDATE commands 
		SET used_count = used_count + 1, last_used_at = ?
		WHERE id = ?
	`, time.Now(), commandID)
	if err != nil {
		return fmt.Errorf("increment used count for command %q: %w", commandID, err)
	}
	return nil
}

// DeleteSession removes a session and its commands (cascade).
func (s *Store) DeleteSession(sessionID string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	if err != nil {
		return fmt.Errorf("delete session %q: %w", sessionID, err)
	}
	return nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}
