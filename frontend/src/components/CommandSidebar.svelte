<script>
  import { commandHistory, activeTerminalId } from '../lib/stores.js';
  import { CopyCommandToClipboard } from '../../wailsjs/go/main/App.js';

  async function copyCommand(command) {
    try {
      await CopyCommandToClipboard(command.id, $activeTerminalId);
      alert(`Command "${command.command}" copied to clipboard`);
    } catch (error) {
      console.error('Failed to copy command:', error);
      alert('Error copying command: ' + error.message);
    }
  }
</script>

<div class="command-sidebar">
  <h3>Command History</h3>
  {#if $commandHistory.length === 0}
    <p class="empty">No commands yet. Ask the AI for help!</p>
  {:else}
    <div class="command-list">
      {#each $commandHistory as cmd (cmd.id)}
        <div class="command-card">
          <div class="command-text">{cmd.command}</div>
          <div class="command-meta">
            <span class="used">{cmd.usedCount} uses</span>
            <span class="time">{new Date(cmd.createdAt).toLocaleTimeString()}</span>
          </div>
          <button on:click={() => copyCommand(cmd)}>Copy</button>
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  .command-sidebar {
    padding: 1rem;
    background: #f9f9f9;
    border-left: 1px solid #ddd;
    height: 100%;
    overflow-y: auto;
  }
  h3 {
    margin-top: 0;
    margin-bottom: 1rem;
  }
  .empty {
    color: #999;
    font-style: italic;
  }
  .command-list {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }
  .command-card {
    padding: 0.75rem;
    background: white;
    border: 1px solid #ddd;
    border-radius: 0.5rem;
  }
  .command-text {
    font-family: monospace;
    font-size: 0.9rem;
    margin-bottom: 0.5rem;
    word-break: break-all;
  }
  .command-meta {
    display: flex;
    justify-content: space-between;
    font-size: 0.8rem;
    color: #666;
    margin-bottom: 0.5rem;
  }
  button {
    width: 100%;
    padding: 0.25rem;
    background: #2196f3;
    color: white;
    border: none;
    border-radius: 0.25rem;
    cursor: pointer;
  }
  button:hover {
    background: #0b7dda;
  }
</style>