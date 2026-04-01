<script>
  import { settingsOpen } from '../lib/stores.js';
  import LLMTab from './settings/LLMTab.svelte';

  let activeTab = 'llm';

  function close() {
    settingsOpen.set(false);
  }
</script>

<div class="settings-modal" on:click={close}>
  <div class="settings-content" on:click|stopPropagation>
    <div class="header">
      <h2>Settings</h2>
      <button class="close-btn" on:click={close}>×</button>
    </div>

    <div class="tabs">
      <button class:active={activeTab === 'llm'} on:click={() => activeTab = 'llm'}>LLM</button>
      <button class:active={activeTab === 'terminals'} on:click={() => activeTab = 'terminals'}>Terminals</button>
      <button class:active={activeTab === 'hotkeys'} on:click={() => activeTab = 'hotkeys'}>Hotkeys</button>
      <button class:active={activeTab === 'appearance'} on:click={() => activeTab = 'appearance'}>Appearance</button>
    </div>

    <div class="tab-content">
      {#if activeTab === 'llm'}
        <LLMTab />
      {:else if activeTab === 'terminals'}
        <div class="tab">
          <h3>Terminal Settings</h3>
          <p>Configure terminal detection and capture settings.</p>
        </div>
      {:else if activeTab === 'hotkeys'}
        <div class="tab">
          <h3>Hotkey Configuration</h3>
          <p>Configure global hotkeys for PairAdmin.</p>
        </div>
      {:else if activeTab === 'appearance'}
        <div class="tab">
          <h3>Appearance</h3>
          <p>Choose light/dark theme and UI preferences.</p>
        </div>
      {/if}
    </div>
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
    padding: 0;
    border-radius: 0.5rem;
    max-width: 700px;
    width: 100%;
    max-height: 80vh;
    display: flex;
    flex-direction: column;
  }
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 1rem 1.5rem;
    border-bottom: 1px solid #ddd;
  }
  .close-btn {
    background: none;
    border: none;
    font-size: 2rem;
    cursor: pointer;
    color: #666;
  }
  .tabs {
    display: flex;
    border-bottom: 1px solid #ddd;
  }
  .tabs button {
    flex: 1;
    padding: 1rem;
    background: none;
    border: none;
    cursor: pointer;
    border-bottom: 3px solid transparent;
  }
  .tabs button.active {
    border-bottom-color: #4caf50;
    font-weight: bold;
  }
  .tab-content {
    flex: 1;
    overflow-y: auto;
    padding: 1.5rem;
  }
</style>