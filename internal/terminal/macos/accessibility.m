//go:build darwin
// +build darwin

// Objective-C implementation for macOS Accessibility API
// This file provides C-compatible functions for CGO to call

#include "accessibility.h"
#include <ApplicationServices/ApplicationServices.h>
#include <CoreFoundation/CoreFoundation.h>
#include <stdlib.h>
#include <string.h>

// Check if accessibility is enabled
int accessibility_is_enabled(void) {
    return AXIsProcessTrusted() ? 1 : 0;
}

// Request accessibility permission
int accessibility_request_permission(void) {
    CFDictionaryRef options = CFDictionaryCreate(
        kCFAllocatorDefault,
        (const void**)&kAXTrustedCheckOptionPrompt,
        (const void**)&kCFBooleanTrue,
        1,
        &kCFTypeDictionaryKeyCallBacks,
        &kCFTypeDictionaryValueCallBacks
    );
    
    Boolean result = AXIsProcessTrustedWithOptions(options);
    CFRelease(options);
    return result ? 1 : 0;
}

// Helper: Convert CFString to malloc'd C string
static char* cfstring_to_cstring(CFStringRef cfstr) {
    if (!cfstr) return NULL;
    
    CFIndex length = CFStringGetLength(cfstr);
    CFIndex max_size = CFStringGetMaximumSizeForEncoding(length, kCFStringEncodingUTF8) + 1;
    char* buffer = (char*)malloc(max_size);
    
    if (CFStringGetCString(cfstr, buffer, max_size, kCFStringEncodingUTF8)) {
        return buffer;
    }
    
    free(buffer);
    return NULL;
}

// Get Terminal.app window titles
char** accessibility_get_terminal_window_titles(int* count) {
    if (!count) return NULL;
    *count = 0;
    
    // Get Terminal.app process ID
    NSArray* apps = [NSRunningApplication runningApplicationsWithBundleIdentifier:@"com.apple.Terminal"];
    if ([apps count] == 0) return NULL;
    
    pid_t pid = [[apps objectAtIndex:0] processIdentifier];
    AXUIElementRef app = AXUIElementCreateApplication(pid);
    if (!app) return NULL;
    
    // Get windows
    CFArrayRef windows = NULL;
    AXError err = AXUIElementCopyAttributeValue(app, kAXWindowsAttribute, (CFTypeRef*)&windows);
    if (err != kAXErrorSuccess || !windows) {
        CFRelease(app);
        return NULL;
    }
    
    CFIndex window_count = CFArrayGetCount(windows);
    if (window_count == 0) {
        CFRelease(windows);
        CFRelease(app);
        return NULL;
    }
    
    // Allocate array of strings
    char** titles = (char**)malloc(sizeof(char*) * window_count);
    int valid_count = 0;
    
    for (CFIndex i = 0; i < window_count; i++) {
        AXUIElementRef window = (AXUIElementRef)CFArrayGetValueAtIndex(windows, i);
        CFStringRef title_ref = NULL;
        err = AXUIElementCopyAttributeValue(window, kAXTitleAttribute, (CFTypeRef*)&title_ref);
        
        if (err == kAXErrorSuccess && title_ref) {
            titles[valid_count] = cfstring_to_cstring(title_ref);
            if (titles[valid_count]) {
                valid_count++;
            }
            CFRelease(title_ref);
        }
    }
    
    CFRelease(windows);
    CFRelease(app);
    
    *count = valid_count;
    return titles;
}

// Get Terminal.app window PIDs
int* accessibility_get_terminal_window_pids(int* count) {
    if (!count) return NULL;
    *count = 0;
    
    NSArray* apps = [NSRunningApplication runningApplicationsWithBundleIdentifier:@"com.apple.Terminal"];
    if ([apps count] == 0) return NULL;
    
    pid_t pid = [[apps objectAtIndex:0] processIdentifier];
    
    // For now, return single PID (can be extended for multiple instances)
    int* pids = (int*)malloc(sizeof(int));
    pids[0] = pid;
    *count = 1;
    return pids;
}

// Extract text from a Terminal.app window
char* accessibility_extract_text_from_window(int pid) {
    AXUIElementRef app = AXUIElementCreateApplication(pid);
    if (!app) return NULL;
    
    // Get focused window
    AXUIElementRef window = NULL;
    AXError err = AXUIElementCopyAttributeValue(app, kAXFocusedWindowAttribute, (CFTypeRef*)&window);
    
    if (err != kAXErrorSuccess || !window) {
        // Try first window
        CFArrayRef windows = NULL;
        err = AXUIElementCopyAttributeValue(app, kAXWindowsAttribute, (CFTypeRef*)&windows);
        if (err == kAXErrorSuccess && windows && CFArrayGetCount(windows) > 0) {
            window = (AXUIElementRef)CFArrayGetValueAtIndex(windows, 0);
            // Don't release window - it's owned by the array
        }
    }
    
    if (!window) {
        CFRelease(app);
        return NULL;
    }
    
    // Get focused text area
    AXUIElementRef textarea = NULL;
    err = AXUIElementCopyAttributeValue(window, kAXFocusedTextAreaAttribute, (CFTypeRef*)&textarea);
    
    if (err != kAXErrorSuccess || !textarea) {
        CFRelease(app);
        return NULL;
    }
    
    // Get text value
    CFStringRef text_ref = NULL;
    err = AXUIElementCopyAttributeValue(textarea, kAXValueAttribute, (CFTypeRef*)&text_ref);
    
    char* result = NULL;
    if (err == kAXErrorSuccess && text_ref) {
        result = cfstring_to_cstring(text_ref);
        CFRelease(text_ref);
    }
    
    CFRelease(textarea);
    CFRelease(app);
    
    return result;
}

// Extract text from frontmost Terminal.app window
char* accessibility_extract_text_from_frontmost_terminal(void) {
    NSArray* apps = [NSRunningApplication runningApplicationsWithBundleIdentifier:@"com.apple.Terminal"];
    if ([apps count] == 0) return NULL;
    
    pid_t pid = [[apps objectAtIndex:0] processIdentifier];
    return accessibility_extract_text_from_window(pid);
}

// Memory management
void accessibility_free_string(char* str) {
    if (str) free(str);
}

void accessibility_free_string_array(char** array, int count) {
    if (!array) return;
    for (int i = 0; i < count; i++) {
        if (array[i]) free(array[i]);
    }
    free(array);
}

void accessibility_free_int_array(int* array) {
    if (array) free(array);
}
