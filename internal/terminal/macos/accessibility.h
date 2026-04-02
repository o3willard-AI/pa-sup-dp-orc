//go:build darwin
// +build darwin

#ifndef PAIRADMIN_ACCESSIBILITY_H
#define PAIRADMIN_ACCESSIBILITY_H

#ifdef __cplusplus
extern "C" {
#endif

// Check if accessibility is enabled
int accessibility_is_enabled(void);

// Request accessibility permission (returns 1 if granted, 0 if denied)
int accessibility_request_permission(void);

// Get Terminal.app window titles (returns malloc'd array, use accessibility_free_string_array)
char** accessibility_get_terminal_window_titles(int* count);

// Get Terminal.app window PIDs (returns malloc'd array, use accessibility_free_int_array)
int* accessibility_get_terminal_window_pids(int* count);

// Extract text from a specific Terminal.app window by PID
// Returns malloc'd string, use accessibility_free_string to free
char* accessibility_extract_text_from_window(int pid);

// Extract text from frontmost Terminal.app window
// Returns malloc'd string, use accessibility_free_string to free
char* accessibility_extract_text_from_frontmost_terminal(void);

// Memory management
void accessibility_free_string(char* str);
void accessibility_free_string_array(char** array, int count);
void accessibility_free_int_array(int* array);

#ifdef __cplusplus
}
#endif

#endif // PAIRADMIN_ACCESSIBILITY_H
