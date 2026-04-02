//go:build windows

#ifndef PAIRADMIN_UIA_H
#define PAIRADMIN_UIA_H

#include <stddef.h>

#ifdef __cplusplus
extern "C" {
#endif

/**
 * UI Automation error codes.
 * All functions returning UIA_ErrorCode use these codes.
 */
typedef enum {
    UIA_OK = 0,                       ///< Operation succeeded
    UIA_INIT_FAILED = 1,              ///< UI Automation initialization failed
    UIA_COM_FAILURE = 2,              ///< COM operation failed
    UIA_WINDOW_NOT_FOUND = 3,         ///< Terminal window not found
    UIA_TEXT_PATTERN_UNAVAILABLE = 4, ///< Text pattern not available on window
    UIA_MEMORY_ERROR = 5,             ///< Memory allocation failure
    UIA_INVALID_ARGUMENT = 6,         ///< Invalid argument passed to function
    UIA_UNKNOWN_ERROR = 7             ///< Unknown/unexpected error
} UIA_ErrorCode;

/**
 * Initialize UI Automation subsystem.
 * Must be called before any other functions.
 * 
 * @return UIA_OK on success, or an error code on failure.
 */
UIA_ErrorCode uia_initialize(void);

/**
 * Check if UI Automation is available.
 * 
 * @return 0 if UI Automation is not available, non-zero if available.
 */
int uia_is_available(void);

/**
 * Cleanup UI Automation resources.
 * Should be called when UI Automation is no longer needed.
 */
void uia_cleanup(void);

/**
 * Get titles of all terminal windows.
 * 
 * @param count Output parameter for number of titles returned. Must not be NULL.
 * @param error_code Output parameter for error code. Can be NULL if errors are ignored.
 * @param error_message Output parameter for error message. Can be NULL. If non-NULL and an error occurs,
 *                      will be set to an allocated string that must be freed with free_wide_string().
 * @return Array of allocated wide strings (one per window). Each string must be freed with free_wide_string(),
 *         and the array itself must be freed with free_wide_string_array(). Returns NULL on error.
 */
wchar_t** get_terminal_window_titles(int* count, UIA_ErrorCode* error_code, wchar_t** error_message);

/**
 * Get process IDs of all terminal windows.
 * 
 * @param count Output parameter for number of PIDs returned. Must not be NULL.
 * @param error_code Output parameter for error code. Can be NULL if errors are ignored.
 * @param error_message Output parameter for error message. Can be NULL. If non-NULL and an error occurs,
 *                      will be set to an allocated string that must be freed with free_wide_string().
 * @return Array of allocated integers (one PID per window). The array must be freed with free_int_array().
 *         Returns NULL on error.
 */
int* get_terminal_window_pids(int* count, UIA_ErrorCode* error_code, wchar_t** error_message);

/**
 * Extract text content from a specific terminal window by process ID.
 * 
 * @param pid Process ID of the terminal window.
 * @param error_code Output parameter for error code. Can be NULL if errors are ignored.
 * @param error_message Output parameter for error message. Can be NULL. If non-NULL and an error occurs,
 *                      will be set to an allocated string that must be freed with free_wide_string().
 * @return Allocated wide string containing the extracted text. Must be freed with free_wide_string().
 *         Returns NULL on error.
 */
wchar_t* extract_text_from_window(int pid, UIA_ErrorCode* error_code, wchar_t** error_message);

/**
 * Extract text content from the currently focused terminal window.
 * 
 * @param error_code Output parameter for error code. Can be NULL if errors are ignored.
 * @param error_message Output parameter for error message. Can be NULL. If non-NULL and an error occurs,
 *                      will be set to an allocated string that must be freed with free_wide_string().
 * @return Allocated wide string containing the extracted text. Must be freed with free_wide_string().
 *         Returns NULL on error.
 */
wchar_t* extract_text_from_focused_terminal(UIA_ErrorCode* error_code, wchar_t** error_message);

/**
 * Free an array of wide strings.
 * 
 * @param array Array of wide strings to free. Can be NULL.
 * @param count Number of strings in the array. The function frees each
 *              non-NULL string in the array, then frees the array pointer.
 */
void free_wide_string_array(wchar_t** array, int count);

/**
 * Free a single wide string allocated by UI Automation functions.
 * 
 * @param str Wide string to free. Can be NULL.
 */
void free_wide_string(wchar_t* str);

/**
 * Free an array of integers allocated by get_terminal_window_pids().
 * 
 * @param array Array of integers to free. Can be NULL.
 */
void free_int_array(int* array);

#ifdef __cplusplus
}
#endif

#endif // PAIRADMIN_UIA_H
