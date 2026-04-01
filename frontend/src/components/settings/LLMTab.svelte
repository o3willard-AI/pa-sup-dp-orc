<script>
  let provider = 'openai';
  let openaiKey = '';
  let anthropicKey = '';
  let ollamaUrl = 'http://localhost:11434';
  let model = 'gpt-4';

  function save() {
    alert('Settings saved (placeholder)');
  }
</script>

<div class="tab">
  <h3>LLM Configuration</h3>
  <div class="field">
    <label for="llm-provider-select">Provider</label>
    <select id="llm-provider-select" bind:value={provider}>
      <option value="openai">OpenAI</option>
      <option value="anthropic">Anthropic</option>
      <option value="ollama">Ollama (local)</option>
    </select>
  </div>

  {#if provider === 'openai'}
    <div class="field">
      <label for="openai-api-key">OpenAI API Key</label>
      <input id="openai-api-key" type="password" bind:value={openaiKey} placeholder="sk-..." />
    </div>
    <div class="field">
      <label for="model-input">Model</label>
      <input id="model-input" bind:value={model} placeholder="gpt-4" />
    </div>
  {:else if provider === 'anthropic'}
    <div class="field">
      <label for="anthropic-api-key">Anthropic API Key</label>
      <input id="anthropic-api-key" type="password" bind:value={anthropicKey} placeholder="sk-ant-..." />
    </div>
  {:else if provider === 'ollama'}
    <div class="field">
      <label for="ollama-base-url">Ollama Base URL</label>
      <input id="ollama-base-url" bind:value={ollamaUrl} />
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