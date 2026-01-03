<script>
  import Header from './Header.svelte';
  import ModCard from './ModCard.svelte';
  import { SelectModFolder, ScanMods } from '../wailsjs/go/main/App';

  let darkMode = true;
  let folderPath = "";
  let mods = [];
  let status = "";
  let targetLoader = "Fabric";
  let viewMode = 'grid';

  $: fabricCount = mods.filter(m => m.loader === "Fabric").length;
  $: forgeCount = mods.filter(m => m.loader === "NeoForge").length;
  $: conflictCount = (targetLoader === "Fabric" ? forgeCount : fabricCount);

  $: hasSodium = mods.some(m => m.name.toLowerCase().includes("sodium"));

  function toggleTheme() {
    darkMode = !darkMode;
    if (darkMode) {
      document.body.classList.add('dark');
    } else {
      document.body.classList.remove('dark');
    }
  }

  import { onMount } from 'svelte';
  onMount(() => { if(darkMode) document.body.classList.add('dark'); });

  async function pickFolder() {
    let path = await SelectModFolder();
    if (path) {
      folderPath = path;
      handleScan();
    }
  }

  async function handleScan() {
    status = "Scanning...";
    try {
      mods = await ScanMods(folderPath);
      status = "";
    } catch (err) {
      status = "Error: " + err;
    }
  }
</script>

<div class="app-container">
  <Header {darkMode} {toggleTheme} />

  <main>
    {#if mods.length > 0}
      <div class="status-indicator">
        <span class="pulse">●</span> Glacier Engine Active
      </div>
    {/if}

    <div class="controls-card">
      <div class="loader-selector">
        <label>Target Mod Loader</label>
        <div class="toggle-group">
          <button class="toggle-btn {targetLoader === 'Fabric' ? 'active fabric' : ''}" on:click={() => targetLoader = 'Fabric'}>
            Fabric
          </button>
          <button class="toggle-btn {targetLoader === 'NeoForge' ? 'active forge' : ''}" on:click={() => targetLoader = 'NeoForge'}>
            NeoForge
          </button>
        </div>
      </div>

      <div class="folder-selector">
        <label>Mods Folder</label>
        <div class="input-group">
          <div class="path-display">{folderPath || "No folder selected..."}</div>
          <button class="btn btn-primary" on:click={pickFolder}>Select Folder</button>
        </div>
      </div>
    </div>

    {#if conflictCount > 0}
      <div class="alert critical">
        <div class="icon">🛑</div>
        <div class="content">
          <h3>Incompatibility Detected</h3>
          <p>{conflictCount} {targetLoader === 'Fabric' ? 'NeoForge' : 'Fabric'} mods detected while targeting {targetLoader}. These will cause crashes.</p>
        </div>
      </div>
    {/if}

    {#if hasSodium}
      <div class="alert warning">
        <div class="icon">⚠️</div>
        <div class="content">
          <h3>Sodium Detected</h3>
          <p>Make sure to not have OptiFine alongside or things may break. If you have OptiFine, you may want to remove one or the other.</p>
        </div>
      </div>
    {/if}

    {#if mods.length > 0}
      <div class="stats-bar">
        <div class="stat"> <b>{mods.length}</b> Mods Scanned</div>
        <div class="stat text-success">✔ <b>{mods.length - conflictCount}</b> Compatible</div>
        {#if conflictCount > 0}
          <div class="stat text-danger">✖ <b>{conflictCount}</b> Blocked</div>
        {/if}
      </div>

      <div class="section-header">
        <h2 class="section-title">Installed Mods</h2>
        
        <button class="view-toggle-btn" on:click={() => viewMode = viewMode === 'grid' ? 'list' : 'grid'} title="Toggle View">
          {#if viewMode === 'grid'}
            <i class="fa-solid fa-bars"></i>
          {:else}
            <i class="fa-solid fa-border-all"></i>
          {/if}
        </button>
      </div>

      <div class="grid {viewMode === 'list' ? 'list-view' : ''}">
        {#each mods as mod (mod.fileName)}
          <ModCard {mod} {targetLoader} {viewMode} />
        {/each}
      </div>
    {/if}
  </main>

  <footer>
    <div class="footer-content">
      <span class="copyright">© 2026 BySnowden</span>
      <div class="social-links">
        <a href="https://github.com/BySnowden" target="_blank" class="social-btn github" title="GitHub">
          <i class="fa-brands fa-github"></i>
        </a>
        <a href="https://ko-fi.com/bysnowden" target="_blank" class="social-btn kofi" title="Support on Ko-fi">
          <i class="fa-solid fa-mug-hot"></i> <span>Support</span>
        </a>
      </div>
    </div>
  </footer>
</div>

<style>
  .app-container { min-height: 100vh; display: flex; flex-direction: column; }
  
  main { max-width: 1200px; margin: 0 auto; padding: 24px; width: 100%; box-sizing: border-box; flex: 1; }

  footer {
    border-top: 1px solid var(--card-border);
    padding: 20px 0;
    margin-top: auto;
    background: var(--card);
  }

  .footer-content {
    max-width: 1200px;
    margin: 0 auto;
    padding: 0 24px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .copyright {
    font-size: 0.85rem;
    color: var(--muted-fg);
    font-weight: 500;
  }

  .social-links { display: flex; gap: 12px; }

  .social-btn {
    text-decoration: none;
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 16px;
    border-radius: 6px;
    font-size: 0.85rem;
    font-weight: 600;
    transition: all 0.2s;
  }

  .social-btn.github {
    background: var(--background);
    color: var(--foreground);
    border: 1px solid var(--card-border);
  }
  .social-btn.github:hover {
    background: var(--card-border);
  }

  .social-btn.kofi {
    background: #FF5E5B;
    color: white;
    border: 1px solid #FF5E5B;
  }
  .social-btn.kofi:hover {
    background: #FF413D;
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(255, 94, 91, 0.3);
  }

  .section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
  .section-title { margin: 0; font-size: 1.25rem; font-weight: 600; }
  .view-toggle-btn { background: var(--card); border: 1px solid var(--card-border); color: var(--muted-fg); width: 36px; height: 36px; border-radius: 6px; cursor: pointer; display: grid; place-items: center; transition: all 0.2s; font-size: 1rem; }
  .view-toggle-btn:hover { background: var(--card-border); color: var(--foreground); transform: scale(1.05); }
  .grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(280px, 1fr)); gap: 16px; }
  .grid.list-view { grid-template-columns: 1fr; gap: 8px; }
  .status-indicator { display: inline-flex; align-items: center; gap: 8px; background: var(--card); border: 1px solid var(--card-border); padding: 6px 12px; border-radius: 20px; font-size: 0.85rem; font-weight: 600; color: var(--primary); margin-bottom: 16px; box-shadow: 0 2px 4px rgba(0,0,0,0.05); }
  .pulse { color: var(--success); animation: pulse 2s infinite; }
  @keyframes pulse { 0% { opacity: 1; } 50% { opacity: 0.5; } 100% { opacity: 1; } }
  .controls-card { background: var(--card); border: 1px solid var(--card-border); border-radius: var(--radius); padding: 24px; display: grid; grid-template-columns: 1fr 1.5fr; gap: 24px; margin-bottom: 24px; }
  label { display: block; font-size: 0.85rem; font-weight: 600; margin-bottom: 8px; color: var(--muted-fg); }
  .toggle-group { display: flex; background: var(--background); padding: 4px; border-radius: 8px; border: 1px solid var(--card-border); }
  .toggle-btn { flex: 1; padding: 10px; border: none; background: transparent; color: var(--muted-fg); font-weight: 600; border-radius: 6px; cursor: pointer; transition: all 0.2s; }
  .toggle-btn.active { background: var(--card); box-shadow: 0 2px 4px rgba(0,0,0,0.1); color: var(--foreground); }
  .toggle-btn.active.fabric { border-bottom: 2px solid var(--primary); color: var(--primary); }
  .toggle-btn.active.forge { border-bottom: 2px solid #F97316; color: #F97316; }
  .input-group { display: flex; gap: 12px; }
  .path-display { flex: 1; background: var(--background); border: 1px solid var(--card-border); padding: 10px; border-radius: 6px; font-family: monospace; color: var(--muted-fg); display: flex; align-items: center; }
  .alert { padding: 16px; border-radius: var(--radius); display: flex; gap: 12px; align-items: center; margin-bottom: 16px; }
  .alert.critical { background: rgba(239, 68, 68, 0.1); border: 1px solid var(--destructive); color: var(--destructive); }
  .alert.warning { background: var(--card); border: 1px solid var(--card-border); color: var(--muted-fg); }
  .alert.warning .icon { color: var(--warning); }
  .alert h3 { margin: 0; font-size: 1rem; }
  .alert p { margin: 4px 0 0 0; font-size: 0.9rem; opacity: 0.9; }
  .stats-bar { display: flex; gap: 24px; margin-bottom: 24px; background: var(--card); padding: 12px 24px; border-radius: var(--radius); border: 1px solid var(--card-border); }
  .stat { font-size: 0.9rem; }
  .text-success { color: var(--success); }
  .text-danger { color: var(--destructive); }
  @media (max-width: 768px) { .controls-card { grid-template-columns: 1fr; } }
</style>