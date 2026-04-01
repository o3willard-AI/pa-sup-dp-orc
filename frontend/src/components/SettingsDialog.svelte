<script>
  import { settingsOpen } from '../lib/stores.js';
  import { tick } from 'svelte';
  import LLMTab from './settings/LLMTab.svelte';
  import TerminalsTab from './settings/TerminalsTab.svelte';
  import HotkeysTab from './settings/HotkeysTab.svelte';
  import AppearanceTab from './settings/AppearanceTab.svelte';

  let activeTab = 'llm';
  const tabOrder = ['llm', 'terminals', 'hotkeys', 'appearance'];

  function close() {
    settingsOpen.set(false);
  }

  function handleKeydown(event) {
    if (event.key === 'Escape') {
      close();
    }
  }

  async function handleTabKeydown(event) {
    const currentIndex = tabOrder.indexOf(activeTab);
    let newIndex = currentIndex;
    
    if (event.key === 'ArrowLeft') {
      newIndex = currentIndex > 0 ? currentIndex - 1 : tabOrder.length - 1;
    } else if (event.key === 'ArrowRight') {
      newIndex = currentIndex < tabOrder.length - 1 ? currentIndex + 1 : 0;
    } else if (event.key === 'Home') {
      newIndex = 0;
    } else if (event.key === 'End') {
      newIndex = tabOrder.length - 1;
    } else {
      return;
    }
    
    event.preventDefault();
    activeTab = tabOrder[newIndex];
    await tick();
    document.getElementById(`tab-${activeTab}`)?.focus();
  }

  function handleOverlayKeydown(event) {
    if (event.key === 'Enter' || event.key === ' ') {
      event.preventDefault();
      close();
    }
  }

  function noop() {}
</script>

<svelte:window on:keydown={handleKeydown} />

<div class="settings-modal" on:click={close} role="button" tabindex="0" aria-label="Close settings" on:keydown={handleOverlayKeydown}>
   <div class="settings-content" on:click|stopPropagation role="dialog" aria-modal="true" aria-labelledby="settings-title" on:keydown={noop}>
    <div class="header">
       <h2 id="settings-title">Settings</h2>
       <button class="close-btn" on:click={close} aria-label="Close">×</button>
    </div>

     <div class="tabs" role="tablist">
        <button id="tab-llm" role="tab" aria-selected={activeTab === 'llm'} aria-controls="settings-tabpanel" class:active={activeTab === 'llm'} on:click={() => activeTab = 'llm'} on:keydown={handleTabKeydown}>LLM</button>
        <button id="tab-terminals" role="tab" aria-selected={activeTab === 'terminals'} aria-controls="settings-tabpanel" class:active={activeTab === 'terminals'} on:click={() => activeTab = 'terminals'} on:keydown={handleTabKeydown}>Terminals</button>
        <button id="tab-hotkeys" role="tab" aria-selected={activeTab === 'hotkeys'} aria-controls="settings-tabpanel" class:active={activeTab === 'hotkeys'} on:click={() => activeTab = 'hotkeys'} on:keydown={handleTabKeydown}>Hotkeys</button>
        <button id="tab-appearance" role="tab" aria-selected={activeTab === 'appearance'} aria-controls="settings-tabpanel" class:active={activeTab === 'appearance'} on:click={() => activeTab = 'appearance'} on:keydown={handleTabKeydown}>Appearance</button>
    </div>

     <div id="settings-tabpanel" class="tab-content" role="tabpanel" aria-labelledby="tab-{activeTab}">
      {#if activeTab === 'llm'}
        <LLMTab />
      {:else if activeTab === 'terminals'}
        <TerminalsTab />
      {:else if activeTab === 'hotkeys'}
        <HotkeysTab />
      {:else if activeTab === 'appearance'}
        <AppearanceTab />
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