# Design: Command History Integration Fix

## Overview
Fix critical integration issues identified in code quality review for Task 8. The frontend `CommandSidebar` uses a `commandHistory` store but lacks backend binding to fetch commands. Additionally, chat messages lack command IDs needed for the copy functionality.

## Changes

### Backend Changes
1. **Add `GetCommandsByTerminal` method to `app.go`**
   - Delegates to `sessionStore.GetCommandsByTerminal(terminalID string) ([]SuggestedCommand, error)`
   - Wails automatically exposes this to frontend
   - Returns `[]SuggestedCommand` (already exported type)

2. **Modify `SendMessage` return type**
   - Current: `(string, error)` (content only)
   - New: `(SendMessageResponse, error)` where `SendMessageResponse` is a struct:
     ```go
     type SendMessageResponse struct {
         Content   string `json:"content"`
         CommandID string `json:"commandID"`
     }
     ```
   - `CommandID` is empty string if response is not stored as a command
   - Requires changes in `chat_handlers.go` to return `(content string, commandID string, err error)`
   - Update `app.go` to create and return `SendMessageResponse` struct

### Frontend Changes
1. **Update `ChatArea.svelte`**
   - Handle new return type from `SendMessage` (object with `content` and `commandID`)
   - Store `commandID` in message object (add field to message store)
   - Update copy button to call `CopyCommandToClipboard(commandID, terminalID)` using the stored ID

2. **Update `CommandSidebar.svelte`**
   - Add function to fetch command history: `GetCommandsByTerminal(terminalID)`
   - Update `commandHistory` store with fetched commands
   - Trigger fetch when `activeTerminalId` changes (reactive statement)
   - Also fetch after each successful `SendMessage` to keep history updated

### Data Flow
1. User sends message → `SendMessage` returns `SendMessageResponse` with `Content` and `CommandID` fields
2. Frontend adds message to store with `commandID` property (from `CommandID` field)
3. Frontend calls `GetCommandsByTerminal` to refresh command history
4. CommandSidebar displays updated history
5. Copy button uses stored `commandID` to call `CopyCommandToClipboard`

## Implementation Plan
1. Backend modifications (Go)
2. Frontend modifications (Svelte/JavaScript)
3. Testing: ensure Go compilation, frontend builds, integration works
4. Commit with message "fix: integrate frontend with backend command history"

## Notes
- `SuggestedCommand` struct is exported and will be serialized by Wails
- `GetCommandsByTerminal` already exists in `internal/session/store.go`
- Breaking change to `SendMessage` API is acceptable as project is early stage
- Wails bindings regenerate automatically on build
- New `SendMessageResponse` struct must be exported (capitalized) for Wails to serialize