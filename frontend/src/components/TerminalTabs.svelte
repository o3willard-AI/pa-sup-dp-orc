<script>
  import { terminals, activeTerminalId } from '../lib/stores.js';

  // Placeholder data
  terminals.set([
    { id: 'term1', name: 'bash' },
    { id: 'term2', name: 'zsh' },
  ]);
  activeTerminalId.set('term1');

  function selectTerminal(id) {
    activeTerminalId.set(id);
  }
</script>

<div class="terminal-tabs">
  <h3>Terminals</h3>
  {#each $terminals as term (term.id)}
    <div
      class="tab {term.id === $activeTerminalId ? 'active' : ''}"
      on:click={() => selectTerminal(term.id)}
    >
      {term.name}
    </div>
  {/each}
</div>

<style>
  .terminal-tabs {
    padding: 1rem;
  }
  h3 {
    margin-bottom: 1rem;
  }
  .tab {
    padding: 0.5rem;
    margin-bottom: 0.5rem;
    background: #f0f0f0;
    border-radius: 0.25rem;
    cursor: pointer;
  }
  .tab.active {
    background: #4caf50;
    color: white;
  }
  .tab:hover {
    background: #e0e0e0;
  }
</style>