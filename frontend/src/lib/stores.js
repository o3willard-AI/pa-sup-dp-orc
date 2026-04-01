import { writable } from 'svelte/store';

// Current terminal sessions
export const terminals = writable([]);

// Current chat messages
export const messages = writable([]);

// Command history for the active terminal
export const commandHistory = writable([]);

// Active terminal ID
export const activeTerminalId = writable('');

// Settings dialog open state
export const settingsOpen = writable(false);