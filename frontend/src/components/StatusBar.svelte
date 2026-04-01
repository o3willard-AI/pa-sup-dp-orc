<script>
  import { settingsOpen } from '../lib/stores.js';

  function openSettings() {
    settingsOpen.set(true);
  }
</script>

<div class="status-bar">
  <span class="model">Model: GPT‑4</span>
  <span class="connection">Connected</span>
  <span class="context">Context: 120/4000 tokens</span>
  <button class="settings-btn" on:click={openSettings}>Settings</button>
</div>

<style>
  .status-bar {
    position: fixed;
    bottom: 0;
    left: 0;
    right: 0;
    height: 40px;
    background: #333;
    color: white;
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 0 1rem;
    font-size: 0.9rem;
  }
  .settings-btn {
    background: #555;
    color: white;
    border: none;
    padding: 0.25rem 0.75rem;
    border-radius: 0.25rem;
    cursor: pointer;
  }
  .settings-btn:hover {
    background: #666;
  }
</style>