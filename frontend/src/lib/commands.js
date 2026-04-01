import { GetCommandsByTerminal } from '../../wailsjs/go/main/App.js';
import { commandHistory } from './stores.js';

export async function fetchCommands(terminalId) {
    if (!terminalId) return;
    const commands = await GetCommandsByTerminal(terminalId);
    commandHistory.set(commands);
}