<script>
  export let mod;
  export let targetLoader;
  export let viewMode = 'grid';

  let expanded = false;
  
  $: isConflict = targetLoader && mod.loader !== targetLoader;
  $: modrinthLink = `https://modrinth.com/mods?q=${mod.name}`;
  $: curseLink = `https://www.curseforge.com/minecraft/search?search=${mod.name}`;
</script>

<div class="card {mod.loader.toLowerCase()} {isConflict ? 'conflict' : ''} {expanded ? 'expanded' : ''} {viewMode === 'list' ? 'list-mode' : ''}">
  
  <div class="card-top">
    <div class="title-section">
      <h3 class={isConflict ? 'text-danger' : ''}>{mod.name}</h3>
      <div class="badges">
        <span class="badge version">v{mod.version}</span>
        <span class="badge loader {mod.loader.toLowerCase()}">{mod.loader}</span>
        {#if isConflict}
          <span class="badge conflict">Conflict</span>
          <span class="badge incompatible">Incompatible</span>
        {/if}
      </div>
    </div>
  </div>

  <p class="desc" on:click={() => expanded = !expanded} on:keydown role="button" tabindex="0">
    {mod.description || "No description found."}
  </p>

  {#if expanded}
    <div class="links">
      <a href={modrinthLink} target="_blank" class="link-btn modrinth">Modrinth</a>
      <a href={curseLink} target="_blank" class="link-btn curse">CurseForge</a>
    </div>
    <div class="filename">📄 {mod.fileName}</div>
  {/if}

  <button class="expand-btn" on:click={() => expanded = !expanded}>
    {#if viewMode === 'list'}
       <i class="fa-solid {expanded ? 'fa-chevron-up' : 'fa-chevron-down'}"></i>
    {:else}
       {expanded ? "Show Less ▲" : "Show More ▼"}
    {/if}
  </button>
</div>

<style>
  .card {
    background: var(--card);
    border: 1px solid var(--card-border);
    border-radius: var(--radius);
    padding: 16px;
    transition: all 0.2s;
    display: flex; flex-direction: column; gap: 12px;
  }
  .card:hover { border-color: var(--primary); transform: translateY(-2px); box-shadow: 0 4px 12px rgba(0,0,0,0.05); }
  
  .card.list-mode {
    flex-direction: column;
    align-items: stretch;
    padding: 20px 24px;
    gap: 16px;
    height: auto;
    transform: none;
    border: 1px solid var(--card-border);
  }
  .card.list-mode:hover { 
    transform: none;
    border-color: var(--primary);
    box-shadow: 0 4px 12px rgba(0,0,0,0.08);
  }

  .card.list-mode .card-top {
    flex-grow: 1;
    margin-bottom: 0;
    align-items: space-between;
    flex-direction: row;
  }
  
  .card.list-mode:not(.expanded) .desc { display: block; }
  
  .card.list-mode.expanded {
    flex-wrap: nowrap; 
  }
  .card.list-mode.expanded .desc {
    width: 100%;
    margin-top: 0px;
    order: unset;
  }
  .card.list-mode.expanded .links { order: unset; width: 100%; margin-top: 12px; }
  .card.list-mode.expanded .filename { order: unset; width: 100%; }

  .card.conflict { border-color: var(--destructive); background: rgba(239, 68, 68, 0.05); opacity: 0.8; }
  .text-danger { color: var(--destructive); }
  .conflict-banner { font-size: 0.75rem; color: var(--destructive); font-weight: bold; display: flex; align-items: center; gap: 10px; letter-spacing: 0.5px; }

  .card-top { display: flex; justify-content: space-between; align-items: flex-start; gap: 16px; width: 100%; }
  .title-section { flex: 1; }
  h3 { margin: 0; font-size: 1.05rem; font-weight: 600; word-break: break-word; }

  .badges { display: flex; gap: 8px; flex-wrap: wrap; flex-shrink: 0; margin-top: 8px; }
  .badge { padding: 6px 10px; border-radius: 4px; font-size: 0.7rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.3px; }
  .badge.version { background: var(--background); color: var(--muted-fg); border: 1px solid var(--card-border); }
  .badge.loader.fabric { background: var(--fabric-bg); color: var(--fabric-text); }
  .badge.loader.neoforge { background: var(--forge-bg); color: var(--forge-text); }
  .badge.conflict { background: var(--alert-critical-bg); color: var(--alert-critical-text); border: 1px solid var(--alert-critical-border); }
  .badge.incompatible { background: var(--alert-critical-bg); color: var(--alert-critical-text); border: 1px solid var(--alert-critical-border); }

  .desc { font-size: 0.85rem; color: var(--muted-fg); line-height: 1.5; cursor: pointer; margin: 0; }
  .card.list-mode .desc { margin: 4px 0 0 0; }
  .card:not(.expanded) .desc { display: -webkit-box; -webkit-line-clamp: 2; line-clamp: 2; -webkit-box-orient: vertical; overflow: hidden; }

  .links { display: flex; gap: 10px; margin-top: 12px; }
  .link-btn { text-decoration: none; font-size: 0.75rem; padding: 6px 12px; border-radius: 4px; font-weight: 600; border: 1px solid transparent; transition: all 0.2s; }
  .link-btn.modrinth { background: rgba(16, 185, 129, 0.1); color: var(--success); border-color: rgba(16, 185, 129, 0.3); }
  .link-btn.modrinth:hover { background: rgba(16, 185, 129, 0.15); }
  .link-btn.curse { background: rgba(249, 115, 22, 0.1); color: #f97316; border-color: rgba(249, 115, 22, 0.3); }
  .link-btn.curse:hover { background: rgba(249, 115, 22, 0.15); }

  .filename { font-size: 0.7rem; font-family: monospace; color: var(--muted-fg); margin-top: 4px; }
  
  .expand-btn { background: none; border: none; color: var(--primary); font-size: 0.75rem; font-weight: 600; cursor: pointer; padding: 0; min-width: 24px; text-align: center;}
</style>