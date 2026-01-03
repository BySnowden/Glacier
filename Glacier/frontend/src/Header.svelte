<script>
  export let darkMode;
  export let toggleTheme;
</script>

<header>
  <div class="brand">
    <div class="logo">
      <i class="fa-solid fa-mountain"></i>
    </div>
    <div class="text">
      <h1>Glacier</h1>
      <p>Mod Manager & Dependency Checker</p>
    </div>
  </div>

  <button class="theme-toggle" on:click={toggleTheme} title={darkMode ? "Switch to Light Mode" : "Switch to Dark Mode"}>
    {#if darkMode}
      <i class="fa-solid fa-moon"></i>
    {:else}
      <i class="fa-solid fa-sun"></i>
    {/if}
  </button>
</header>

<style>
  header {
    background: var(--card);
    border-bottom: 1px solid var(--card-border);
    padding: 16px 24px;
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .brand { display: flex; align-items: center; gap: 12px; }
  
  .logo { 
    width: 40px; 
    height: 40px; 
    background: var(--primary); 
    border-radius: 8px; 
    display: grid; 
    place-items: center; 
    box-shadow: 0 2px 8px rgba(37, 99, 235, 0.2);
  }

  .logo i {
    color: white;
    font-size: 1.2rem;
  }
  
  h1 { font-size: 1.25rem; margin: 0; font-weight: 700; letter-spacing: 0.5px; }
  p { font-size: 0.75rem; color: var(--muted-fg); margin: 0; }

  .theme-toggle {
    background: var(--background);
    border: 1px solid var(--card-border);
    color: var(--foreground);
    width: 40px; 
    height: 40px;
    border-radius: 8px;
    cursor: pointer;
    display: grid;
    place-items: center;
    font-size: 1.2rem;
    transition: all 0.2s ease;
  }

  .theme-toggle:hover {
    background: var(--card-border);
    transform: scale(1.05);
  }
</style>