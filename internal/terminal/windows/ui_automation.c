//go:build windows && !mock

#include <windows.h>
#include <uiautomation.h>
#include <stdlib.h>
#include <string.h>
#include <wchar.h>
#include <oleauto.h>
#include <stdint.h>
#include "ui_automation.h"

// Global COM instance
static IUIAutomation* g_pAutomation = NULL;
static int g_initialized = 0;
static int g_com_owned = 0;

UIA_ErrorCode uia_initialize(void) {
    if (g_initialized) {
        return UIA_OK;
    }
    
    HRESULT hr = CoInitializeEx(NULL, COINIT_APARTMENTTHREADED);
    if (FAILED(hr)) {
        return UIA_INIT_FAILED;
    }
    
    // Track whether we own the COM initialization (S_OK) or it was already initialized (S_FALSE)
    g_com_owned = SUCCEEDED(hr) ? 1 : 0;
    
    hr = CoCreateInstance(&CLSID_CUIAutomation, NULL, CLSCTX_INPROC_SERVER,
                          &IID_IUIAutomation, (void**)&g_pAutomation);
    if (FAILED(hr)) {
        if (g_com_owned) {
            CoUninitialize();
            g_com_owned = 0;
        }
        return UIA_COM_FAILURE;
    }
    
    g_initialized = 1;
    return UIA_OK;
}

int uia_is_available(void) {
    return g_initialized && g_pAutomation != NULL;
}

void uia_cleanup(void) {
    if (g_pAutomation) {
        g_pAutomation->lpVtbl->Release(g_pAutomation);
        g_pAutomation = NULL;
    }
    if (g_initialized) {
        if (g_com_owned) {
            CoUninitialize();
            g_com_owned = 0;
        }
        g_initialized = 0;
    }
}

// Helper: Convert BSTR to malloc'd wchar_t string (TODO: Will be used in Tasks 3-4)
static wchar_t* bstr_to_wchar(BSTR bstr) {
    if (!bstr) return NULL;
    
    size_t len = SysStringLen(bstr);
    // Guard against overflow in (len + 1) * sizeof(wchar_t)
    if (len == SIZE_MAX || len > (SIZE_MAX / sizeof(wchar_t)) - 1) {
        return NULL;
    }
    wchar_t* result = (wchar_t*)malloc((len + 1) * sizeof(wchar_t));
    if (!result) return NULL;
    
    wmemcpy(result, bstr, len);
    result[len] = L'\0';
    return result;
}

// Helper: Duplicate wide string (allocates with malloc)
static wchar_t* duplicate_wstring(const wchar_t* src) {
    if (!src) return NULL;
    size_t len = wcslen(src);
    // Guard against overflow in (len + 1) * sizeof(wchar_t)
    if (len == SIZE_MAX || len > (SIZE_MAX / sizeof(wchar_t)) - 1) {
        return NULL;
    }
    wchar_t* dst = (wchar_t*)malloc((len + 1) * sizeof(wchar_t));
    if (!dst) return NULL;
    wcscpy(dst, src);
    return dst;
}

// Helper: Convert wide string to UTF-8 (caller must free)
static char* wchar_to_utf8(const wchar_t* wstr) {
    if (!wstr) return NULL;
    
    int required_size = WideCharToMultiByte(CP_UTF8, 0, wstr, -1, NULL, 0, NULL, NULL);
    if (required_size <= 0) return NULL;
    
    char* utf8 = (char*)malloc(required_size);
    if (!utf8) return NULL;
    
    int converted = WideCharToMultiByte(CP_UTF8, 0, wstr, -1, utf8, required_size, NULL, NULL);
    if (converted <= 0) {
        free(utf8);
        return NULL;
    }
    
    return utf8;
}

// Helper: Create condition for terminal window classes
static IUIAutomationCondition* create_terminal_condition(IUIAutomation* pAutomation) {
    if (!pAutomation) return NULL;
    
    // Terminal window class names
    const wchar_t* class_names[] = {
        L"CASCADIA_HOSTING_WINDOW_CLASS",  // Windows Terminal
        L"ConsoleWindowClass",              // PowerShell/cmd.exe
        L"PuTTY",                           // PuTTY
        NULL
    };
    
    IUIAutomationCondition* pOrCondition = NULL;
    
    // Create OR condition from individual class name conditions
    for (int i = 0; class_names[i]; i++) {
        VARIANT varClassName;
        varClassName.vt = VT_BSTR;
        varClassName.bstrVal = SysAllocString(class_names[i]);
        if (!varClassName.bstrVal) {
            varClassName.vt = VT_EMPTY;
            VariantClear(&varClassName);
            continue;  // Skip this class if memory allocation fails
        }
        
        IUIAutomationCondition* pClassCondition = NULL;
        HRESULT hr = pAutomation->lpVtbl->CreatePropertyCondition(
            pAutomation,
            UIA_ClassNamePropertyId,
            varClassName,
            &pClassCondition
        );
        
        VariantClear(&varClassName);
        
        if (SUCCEEDED(hr) && pClassCondition) {
            if (!pOrCondition) {
                pOrCondition = pClassCondition;
            } else {
                IUIAutomationCondition* pNewOr = NULL;
                hr = pAutomation->lpVtbl->CreateOrCondition(
                    pAutomation,
                    pOrCondition,
                    pClassCondition,
                    &pNewOr
                );
                if (SUCCEEDED(hr) && pNewOr) {
                    pOrCondition->lpVtbl->Release(pOrCondition);
                    pClassCondition->lpVtbl->Release(pClassCondition);
                    pOrCondition = pNewOr;
                } else {
                    // Failed to create OR condition, release the class condition
                    pClassCondition->lpVtbl->Release(pClassCondition);
                    // Keep existing pOrCondition
                }
            }
        }
    }
    
    return pOrCondition;
}

// Helper: Find terminal window by PID
static IUIAutomationElement* find_window_by_pid(int pid) {
    if (!g_initialized || !g_pAutomation || pid <= 0) {
        return NULL;
    }
    
    // Get desktop element
    IUIAutomationElement* pDesktop = NULL;
    HRESULT hr = g_pAutomation->lpVtbl->GetRootElement(g_pAutomation, &pDesktop);
    if (FAILED(hr) || !pDesktop) {
        return NULL;
    }
    
    // Create terminal condition
    IUIAutomationCondition* pCondition = create_terminal_condition(g_pAutomation);
    if (!pCondition) {
        pDesktop->lpVtbl->Release(pDesktop);
        return NULL;
    }
    
    // Find all terminal windows
    IUIAutomationElementArray* pElements = NULL;
    hr = g_pAutomation->lpVtbl->FindAll(
        g_pAutomation,
        TreeScope_Children,
        pCondition,
        &pElements
    );
    
    pCondition->lpVtbl->Release(pCondition);
    pDesktop->lpVtbl->Release(pDesktop);
    
    if (FAILED(hr) || !pElements) {
        return NULL;
    }
    
    // Get element count
    int element_count = 0;
    hr = pElements->lpVtbl->get_Length(pElements, &element_count);
    if (FAILED(hr)) {
        pElements->lpVtbl->Release(pElements);
        return NULL;
    }
    
    // Find window with matching PID
    IUIAutomationElement* pFoundElement = NULL;
    for (int i = 0; i < element_count && !pFoundElement; i++) {
        IUIAutomationElement* pElement = NULL;
        hr = pElements->lpVtbl->GetElement(pElements, i, &pElement);
        if (SUCCEEDED(hr) && pElement) {
            int element_pid = 0;
            hr = pElement->lpVtbl->GetCurrentProcessId(pElement, &element_pid);
            if (SUCCEEDED(hr) && element_pid == pid) {
                pFoundElement = pElement;
                pFoundElement->lpVtbl->AddRef(pFoundElement); // Keep reference
            }
            pElement->lpVtbl->Release(pElement);
        }
    }
    
    pElements->lpVtbl->Release(pElements);
    return pFoundElement;
}

// Helper: Free single string
void free_wide_string(wchar_t* str) {
    if (str) free(str);
}

// Helper: Free string array
void free_wide_string_array(wchar_t** array, int count) {
    if (!array) return;
    if (count <= 0) {
        free(array);
        return;
    }
    for (int i = 0; i < count; i++) {
        free_wide_string(array[i]);
    }
    free(array);
}

// Helper: Free integer array
void free_int_array(int* array) {
    if (array) free(array);
}

// Stub implementations for now (will be replaced in later tasks)
wchar_t** get_terminal_window_titles(int* count, UIA_ErrorCode* error_code, wchar_t** error_message) {
    if (error_code) *error_code = UIA_OK;
    if (error_message) *error_message = NULL;
    *count = 0;
    
    if (!g_initialized || !g_pAutomation) {
        if (error_code) *error_code = UIA_INIT_FAILED;
        if (error_message) *error_message = duplicate_wstring(L"UI Automation not initialized");
        return NULL;
    }
    
    // Get desktop element
    IUIAutomationElement* pDesktop = NULL;
    HRESULT hr = g_pAutomation->lpVtbl->GetRootElement(g_pAutomation, &pDesktop);
    if (FAILED(hr) || !pDesktop) {
        if (error_code) *error_code = UIA_COM_FAILURE;
        if (error_message) *error_message = duplicate_wstring(L"Failed to get desktop element");
        return NULL;
    }
    
    // Create terminal condition
    IUIAutomationCondition* pCondition = create_terminal_condition(g_pAutomation);
    if (!pCondition) {
        pDesktop->lpVtbl->Release(pDesktop);
        if (error_code) *error_code = UIA_COM_FAILURE;
        if (error_message) *error_message = duplicate_wstring(L"Failed to create terminal condition");
        return NULL;
    }
    
    // Find all terminal windows
    IUIAutomationElementArray* pElements = NULL;
    hr = g_pAutomation->lpVtbl->FindAll(
        g_pAutomation,
        TreeScope_Children,
        pCondition,
        &pElements
    );
    
    pCondition->lpVtbl->Release(pCondition);
    pDesktop->lpVtbl->Release(pDesktop);
    
    if (FAILED(hr) || !pElements) {
        if (error_code) *error_code = UIA_COM_FAILURE;
        if (error_message) *error_message = duplicate_wstring(L"Failed to find terminal windows");
        return NULL;
    }
    
    // Get element count
    int element_count = 0;
    hr = pElements->lpVtbl->get_Length(pElements, &element_count);
    if (FAILED(hr)) {
        pElements->lpVtbl->Release(pElements);
        if (error_code) *error_code = UIA_COM_FAILURE;
        if (error_message) *error_message = duplicate_wstring(L"Failed to get element count");
        return NULL;
    }
    
    // If no terminal windows found, return empty result
    if (element_count == 0) {
        pElements->lpVtbl->Release(pElements);
        *count = 0;
        return NULL;
    }
    
    // Allocate result arrays
    wchar_t** titles = (wchar_t**)malloc(element_count * sizeof(wchar_t*));
    if (!titles) {
        pElements->lpVtbl->Release(pElements);
        if (error_code) *error_code = UIA_MEMORY_ERROR;
        if (error_message) *error_message = duplicate_wstring(L"Memory allocation failed");
        return NULL;
    }
    
    int valid_count = 0;
    for (int i = 0; i < element_count; i++) {
        IUIAutomationElement* pElement = NULL;
        hr = pElements->lpVtbl->GetElement(pElements, i, &pElement);
        if (SUCCEEDED(hr) && pElement) {
            BSTR bstrName = NULL;
            hr = pElement->lpVtbl->GetCurrentName(pElement, &bstrName);
            if (SUCCEEDED(hr) && bstrName) {
                titles[valid_count] = bstr_to_wchar(bstrName);
                if (titles[valid_count]) {
                    valid_count++;
                }
                SysFreeString(bstrName);
            }
            pElement->lpVtbl->Release(pElement);
        }
    }
    
    pElements->lpVtbl->Release(pElements);
    
    *count = valid_count;
    return titles;
}

int* get_terminal_window_pids(int* count, UIA_ErrorCode* error_code, wchar_t** error_message) {
    if (error_code) *error_code = UIA_OK;
    if (error_message) *error_message = NULL;
    *count = 0;
    
    if (!g_initialized || !g_pAutomation) {
        if (error_code) *error_code = UIA_INIT_FAILED;
        if (error_message) *error_message = duplicate_wstring(L"UI Automation not initialized");
        return NULL;
    }
    
    // Get desktop element
    IUIAutomationElement* pDesktop = NULL;
    HRESULT hr = g_pAutomation->lpVtbl->GetRootElement(g_pAutomation, &pDesktop);
    if (FAILED(hr) || !pDesktop) {
        if (error_code) *error_code = UIA_COM_FAILURE;
        if (error_message) *error_message = duplicate_wstring(L"Failed to get desktop element");
        return NULL;
    }
    
    // Create terminal condition
    IUIAutomationCondition* pCondition = create_terminal_condition(g_pAutomation);
    if (!pCondition) {
        pDesktop->lpVtbl->Release(pDesktop);
        if (error_code) *error_code = UIA_COM_FAILURE;
        if (error_message) *error_message = duplicate_wstring(L"Failed to create terminal condition");
        return NULL;
    }
    
    // Find all terminal windows
    IUIAutomationElementArray* pElements = NULL;
    hr = g_pAutomation->lpVtbl->FindAll(
        g_pAutomation,
        TreeScope_Children,
        pCondition,
        &pElements
    );
    
    pCondition->lpVtbl->Release(pCondition);
    pDesktop->lpVtbl->Release(pDesktop);
    
    if (FAILED(hr) || !pElements) {
        if (error_code) *error_code = UIA_COM_FAILURE;
        if (error_message) *error_message = duplicate_wstring(L"Failed to find terminal windows");
        return NULL;
    }
    
    // Get element count
    int element_count = 0;
    hr = pElements->lpVtbl->get_Length(pElements, &element_count);
    if (FAILED(hr)) {
        pElements->lpVtbl->Release(pElements);
        if (error_code) *error_code = UIA_COM_FAILURE;
        if (error_message) *error_message = duplicate_wstring(L"Failed to get element count");
        return NULL;
    }
    
    // If no terminal windows found, return empty result
    if (element_count == 0) {
        pElements->lpVtbl->Release(pElements);
        *count = 0;
        return NULL;
    }
    
    // Allocate result array
    int* pids = (int*)malloc(element_count * sizeof(int));
    if (!pids) {
        pElements->lpVtbl->Release(pElements);
        if (error_code) *error_code = UIA_MEMORY_ERROR;
        if (error_message) *error_message = duplicate_wstring(L"Memory allocation failed");
        return NULL;
    }
    
    int valid_count = 0;
    for (int i = 0; i < element_count; i++) {
        IUIAutomationElement* pElement = NULL;
        hr = pElements->lpVtbl->GetElement(pElements, i, &pElement);
        if (SUCCEEDED(hr) && pElement) {
            int pid = 0;
            hr = pElement->lpVtbl->GetCurrentProcessId(pElement, &pid);
            if (SUCCEEDED(hr) && pid > 0) {
                pids[valid_count] = pid;
                valid_count++;
            }
            pElement->lpVtbl->Release(pElement);
        }
    }
    
    pElements->lpVtbl->Release(pElements);
    
    *count = valid_count;
    return pids;
}

wchar_t* extract_text_from_window(int pid, UIA_ErrorCode* error_code, wchar_t** error_message) {
    if (error_code) *error_code = UIA_OK;
    if (error_message) *error_message = NULL;
    
    if (!g_initialized || !g_pAutomation) {
        if (error_code) *error_code = UIA_INIT_FAILED;
        if (error_message) *error_message = duplicate_wstring(L"UI Automation not initialized");
        return NULL;
    }
    
    // Find window by PID
    IUIAutomationElement* pElement = find_window_by_pid(pid);
    if (!pElement) {
        if (error_code) *error_code = UIA_WINDOW_NOT_FOUND;
        if (error_message) *error_message = duplicate_wstring(L"Terminal window not found");
        return NULL;
    }
    
    // Get text pattern
    ITextPattern* pTextPattern = NULL;
    HRESULT hr = pElement->lpVtbl->GetCurrentPattern(
        pElement,
        UIA_TextPatternId,
        (IUnknown**)&pTextPattern
    );
    
    if (FAILED(hr) || !pTextPattern) {
        pElement->lpVtbl->Release(pElement);
        if (error_code) *error_code = UIA_TEXT_PATTERN_UNAVAILABLE;
        if (error_message) *error_message = duplicate_wstring(L"Text pattern not available for window");
        return NULL;
    }
    
    // Get text range
    ITextRange* pTextRange = NULL;
    hr = pTextPattern->lpVtbl->GetSelection(pTextPattern, &pTextRange);
    if (FAILED(hr)) {
        // Fallback to document range
        hr = pTextPattern->lpVtbl->GetDocumentRange(pTextPattern, &pTextRange);
    }
    
    wchar_t* result = NULL;
    if (SUCCEEDED(hr) && pTextRange) {
        BSTR bstrText = NULL;
        hr = pTextRange->lpVtbl->GetText(pTextRange, -1, &bstrText);
        if (SUCCEEDED(hr) && bstrText) {
            result = bstr_to_wchar(bstrText);
            SysFreeString(bstrText);
        }
        pTextRange->lpVtbl->Release(pTextRange);
    }
    
    pTextPattern->lpVtbl->Release(pTextPattern);
    pElement->lpVtbl->Release(pElement);
    
    if (!result) {
        if (error_code) *error_code = UIA_COM_FAILURE;
        if (error_message) *error_message = duplicate_wstring(L"Failed to extract text from window");
    }
    
    return result;
}

wchar_t* extract_text_from_focused_terminal(UIA_ErrorCode* error_code, wchar_t** error_message) {
    if (error_code) *error_code = UIA_OK;
    if (error_message) *error_message = NULL;
    
    if (!g_initialized || !g_pAutomation) {
        if (error_code) *error_code = UIA_INIT_FAILED;
        if (error_message) *error_message = duplicate_wstring(L"UI Automation not initialized");
        return NULL;
    }
    
    // Get focused element
    IUIAutomationElement* pFocused = NULL;
    HRESULT hr = g_pAutomation->lpVtbl->GetFocusedElement(g_pAutomation, &pFocused);
    if (FAILED(hr) || !pFocused) {
        if (error_code) *error_code = UIA_COM_FAILURE;
        if (error_message) *error_message = duplicate_wstring(L"Failed to get focused element");
        return NULL;
    }
    
    // Check if focused element is a terminal window
    int pid = 0;
    hr = pFocused->lpVtbl->GetCurrentProcessId(pFocused, &pid);
    if (FAILED(hr) || pid <= 0) {
        pFocused->lpVtbl->Release(pFocused);
        if (error_code) *error_code = UIA_WINDOW_NOT_FOUND;
        if (error_message) *error_message = duplicate_wstring(L"Focused window is not a terminal");
        return NULL;
    }
    
    // Use the PID-based extraction
    wchar_t* result = extract_text_from_window(pid, error_code, error_message);
    pFocused->lpVtbl->Release(pFocused);
    
    return result;
}
