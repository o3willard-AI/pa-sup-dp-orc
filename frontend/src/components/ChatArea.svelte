<script>
  import { messages, activeTerminalId, commandHistory } from '../lib/stores.js';
  import { SendMessage, CopyCommandToClipboard, GetCommandsByTerminal } from '../../wailsjs/go/main/App.js';

  let inputText = '';
  let isLoading = false;

  function generateId() {
    return Date.now().toString(36) + Math.random().toString(36).substr(2);
  }

  async function handleSend() {
    if (!inputText.trim() || isLoading) return;
    const terminalId = $activeTerminalId;
    if (!terminalId) {
      alert('Please select a terminal first');
      return;
    }

    const userMessage = inputText;
    inputText = '';
    isLoading = true;

    // Add user message to UI
    messages.update(msgs => [...msgs, { id: generateId(), role: 'user', content: userMessage }]);

    try {
      const response = await SendMessage(terminalId, userMessage);
      // response is now an object with content and commandID fields
      messages.update(msgs => [...msgs, { 
        id: generateId(),
        role: 'assistant', 
        content: response.content,
        commandID: response.commandID 
      }]);
      
      // Refresh command history for this terminal
      await fetchCommands(terminalId);
    } catch (error) {
      console.error('Failed to send message:', error);
      alert('Error: ' + error.message);
    } finally {
      isLoading = false;
    }
  }
  
  async function copyCommand(commandID) {
    const terminalId = $activeTerminalId;
    if (!terminalId) {
      alert('Please select a terminal first');
      return;
    }
    try {
      await CopyCommandToClipboard(commandID, terminalId);
      alert('Command copied to clipboard');
    } catch (error) {
      console.error('Failed to copy command:', error);
      alert('Error copying command: ' + error.message);
    }
  }
  
  async function fetchCommands(terminalId) {
    if (!terminalId) return;
    try {
      const commands = await GetCommandsByTerminal(terminalId);
      commandHistory.set(commands);
    } catch (error) {
      console.error('Failed to fetch commands:', error);
    }
  }

  function handleKeyDown(e) {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  }
</script>

<div class="chat-area">
  <div class="messages">
    {#each $messages as msg (msg.id)}
      <div class="message {msg.role}">
        <div class="role">{msg.role === 'user' ? 'You' : 'AI'}</div>
        <div class="content">
          {#if msg.role === 'assistant' && msg.content.startsWith('$')}
            <pre><code>{msg.content}</code></pre>
            {#if msg.commandID}
              <button on:click={() => copyCommand(msg.commandID)}>
                Copy to Terminal
              </button>
            {/if}
          {:else}
            {msg.content}
          {/if}
        </div>
      </div>
    {/each}
  </div>

  <div class="input-area">
    <textarea
      bind:value={inputText}
      on:keydown={handleKeyDown}
      placeholder="Ask the AI for a command..."
      rows="3"
      disabled={isLoading}
    />
    <button on:click={handleSend} disabled={isLoading || !inputText.trim()}>
      {isLoading ? 'Sending…' : 'Send'}
    </button>
  </div>
</div>

<style>
  .chat-area {
    display: flex;
    flex-direction: column;
    height: 100%;
    padding: 1rem;
  }
  .messages {
    flex: 1;
    overflow-y: auto;
    margin-bottom: 1rem;
  }
  .message {
    margin-bottom: 1rem;
    padding: 0.75rem;
    border-radius: 0.5rem;
    background: #f5f5f5;
  }
  .message.user {
    background: #e3f2fd;
  }
  .message.assistant {
    background: #f1f8e9;
  }
  .role {
    font-weight: bold;
    font-size: 0.8rem;
    margin-bottom: 0.25rem;
    color: #666;
  }
  .content pre {
    margin: 0;
    padding: 0.5rem;
    background: #2d2d2d;
    color: #f8f8f2;
    border-radius: 0.25rem;
    overflow-x: auto;
  }
  .input-area {
    display: flex;
    gap: 0.5rem;
  }
  textarea {
    flex: 1;
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 0.25rem;
    font-family: monospace;
  }
  button {
    padding: 0.5rem 1rem;
    background: #4caf50;
    color: white;
    border: none;
    border-radius: 0.25rem;
    cursor: pointer;
  }
  button:disabled {
    background: #ccc;
    cursor: not-allowed;
  }
</style>