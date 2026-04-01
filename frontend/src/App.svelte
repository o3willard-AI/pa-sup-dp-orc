<script>
  import ChatArea from './components/ChatArea.svelte';
  import CommandSidebar from './components/CommandSidebar.svelte';
  import TerminalTabs from './components/TerminalTabs.svelte';
  import StatusBar from './components/StatusBar.svelte';
  import SettingsDialog from './components/SettingsDialog.svelte';
  import { settingsOpen } from './lib/stores.js';
</script>

<main>
  <div class="app-layout">
    <!-- Left: Terminal tabs -->
    <div class="left-panel">
      <TerminalTabs />
    </div>

    <!-- Middle: Chat area -->
    <div class="middle-panel">
      <ChatArea />
    </div>

    <!-- Right: Command sidebar -->
    <div class="right-panel">
      <CommandSidebar />
    </div>
  </div>

  <!-- Status bar at bottom -->
  <StatusBar />

  <!-- Settings dialog (modal) -->
  {#if $settingsOpen}
    <SettingsDialog />
  {/if}
</main>

<style>
  * {
    box-sizing: border-box;
    margin: 0;
    padding: 0;
  }
  body {
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    overflow: hidden;
  }
  .app-layout {
    display: flex;
    height: calc(100vh - 40px); /* minus status bar */
  }
  .left-panel {
    width: 250px;
    border-right: 1px solid #ccc;
  }
  .middle-panel {
    flex: 1;
    border-right: 1px solid #ccc;
  }
  .right-panel {
    width: 300px;
  }
</style>