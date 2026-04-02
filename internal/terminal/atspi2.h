#ifndef PAIRADMIN_ATSPI2_H
#define PAIRADMIN_ATSPI2_H

#include <stdint.h>
#include <stdbool.h>

#ifdef __cplusplus
extern "C" {
#endif

// Callback for text changes
typedef void (*TextCallback)(const char* text);

// Initialize AT-SPI (returns 0 on success)
int pa_atspi_init(void);

// Cleanup
void pa_atspi_finalize(void);

// Check if AT-SPI service available
bool pa_atspi_available(void);

// List terminal sessions (returns null-terminated array, caller frees)
char** pa_atspi_list_sessions(int* count);

// Free session list
void pa_atspi_free_sessions(char** sessions);

// Capture text from terminal (caller frees)
char* pa_atspi_capture(const char* objpath);

// Free captured text
void pa_atspi_free_capture(char* text);

// Subscribe to text changes (stub)
bool pa_atspi_subscribe(const char* objpath, TextCallback cb);

// Unsubscribe
void pa_atspi_unsubscribe(const char* objpath);

// Get terminal dimensions
int pa_atspi_get_width(const char* objpath);
int pa_atspi_get_height(const char* objpath);

#ifdef __cplusplus
}
#endif

#endif // PAIRADMIN_ATSPI2_H
