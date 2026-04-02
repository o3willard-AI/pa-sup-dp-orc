//go:build linux

/* AT-SPI2 stub implementation for Linux terminal accessibility
 * Full implementation requires proper D-Bus and AT-SPI2 integration.
 * This stub provides the interface structure for future implementation.
 */

#include "atspi2.h"
#include <stdlib.h>
#include <string.h>

static int initialized = 0;

int pa_atspi_init(void) {
    initialized = 1;
    return 0;
}

void pa_atspi_finalize(void) {
    initialized = 0;
}

bool pa_atspi_available(void) {
    return initialized ? true : false;
}

char** pa_atspi_list_sessions(int* count) {
    if (!count) return NULL;
    *count = 0;
    /* Stub: return empty array */
    return (char**)malloc(sizeof(char*));
}

void pa_atspi_free_sessions(char** sessions) {
    if (sessions) free(sessions);
}

char* pa_atspi_capture(const char* objpath) {
    (void)objpath; /* unused in stub */
    return strdup("");
}

void pa_atspi_free_capture(char* text) {
    if (text) free(text);
}

bool pa_atspi_subscribe(const char* objpath, TextCallback cb) {
    (void)objpath; (void)cb; /* unused in stub */
    return false;
}

void pa_atspi_unsubscribe(const char* objpath) {
    (void)objpath; /* unused in stub */
}

int pa_atspi_get_width(const char* objpath) {
    (void)objpath; /* unused in stub */
    return 80;
}

int pa_atspi_get_height(const char* objpath) {
    (void)objpath; /* unused in stub */
    return 24;
}
