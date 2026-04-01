import { GetCommandsByTerminal } from '../../wailsjs/go/main/App.js';
import { commandHistory } from './stores.js';

export async function fetchCommands(terminalId) {
    if (!terminalId) return;
    try {
        const commands = await GetCommandsByTerminal(terminalId);
        commandHistory.set(commands);
    } catch (error) {
        console.error('Failed to fetch commands:', error);
    }
}