<script>
  import { settingsOpen } from '../lib/stores.js';

  function close() {
    settingsOpen.set(false);
  }
</script>

<div class="settings-modal">
  <div class="settings-content">
    <h2>Settings</h2>
    <p>Settings dialog will be implemented in a later task.</p>
    <button on:click={close}>Close</button>
  </div>
</div>

<style>
  .settings-modal {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }
  .settings-content {
    background: white;
    padding: 2rem;
    border-radius: 0.5rem;
    max-width: 500px;
    width: 100%;
  }
</style>