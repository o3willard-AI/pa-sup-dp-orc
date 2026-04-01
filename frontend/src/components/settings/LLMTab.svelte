<script>
  let provider = 'openai';
  let openaiKey = '';
  let anthropicKey = '';
  let ollamaUrl = 'http://localhost:11434';
  let model = 'gpt-4';

  function save() {
    console.log('Save LLM settings', { provider, openaiKey, anthropicKey, ollamaUrl, model });
    alert('Settings saved (placeholder)');
  }
</script>

<div class="tab">
  <h3>LLM Configuration</h3>
  <div class="field">
    <label>Provider</label>
    <select bind:value={provider}>
      <option value="openai">OpenAI</option>
      <option value="anthropic">Anthropic</option>
      <option value="ollama">Ollama (local)</option>
    </select>
  </div>

  {#if provider === 'openai'}
    <div class="field">
      <label>OpenAI API Key</label>
      <input type="password" bind:value={openaiKey} placeholder="sk-..." />
    </div>
    <div class="field">
      <label>Model</label>
      <input bind:value={model} placeholder="gpt-4" />
    </div>
  {:else if provider === 'anthropic'}
    <div class="field">
      <label>Anthropic API Key</label>
      <input type="password" bind:value={anthropicKey} placeholder="sk-ant-..." />
    </div>
  {:else if provider === 'ollama'}
    <div class="field">
      <label>Ollama Base URL</label>
      <input bind:value={ollamaUrl} />
    </div>
  {/if}

  <button on:click={save}>Save</button>
</div>

<style>
  .tab {
    padding: 1rem;
  }
  .field {
    margin-bottom: 1rem;
  }
  label {
    display: block;
    margin-bottom: 0.25rem;
    font-weight: bold;
  }
  input, select {
    width: 100%;
    padding: 0.5rem;
    border: 1px solid #ccc;
    border-radius: 0.25rem;
  }
  button {
    padding: 0.5rem 1rem;
    background: #4caf50;
    color: white;
    border: none;
    border-radius: 0.25rem;
    cursor: pointer;
  }
</style>