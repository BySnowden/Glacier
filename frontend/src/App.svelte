<script>
    import Header from "./Header.svelte";
    import {
        SelectModFolder,
        ScanMods,
        ValidateMods,
        CheckForUpdates,
        TriggerUpdate,
    } from "../wailsjs/go/main/App";
    import { onMount } from "svelte";

    // Import the custom logo
    import logoIcon from "./assets/images/icon-transparent.png";

    function resetApp() {
        folderPath = "";
        rawMods = [];
        status = "";
        issues = [];
        checkStatus = "";
        expandedMods = new Set();

        // Reset alerts to default so they show up again next time
        showConflictAlert = true;
        showDependencyAlert = true;
        showSodiumAlert = true;
        isAlertsOpen = true;
    }

    // --- UPDATE STATE ---
    let updateAvailable = false;
    let latestVersion = "";
    let updateUrl = "";
    let isCheckingUpdate = false;

    let darkMode = true;
    let folderPath = "";
    let rawMods = [];
    let status = "";
    let targetLoader = "Fabric";
    let viewMode = "grid";

    // --- FILTER & SEARCH ---
    let searchQuery = "";
    let sortBy = "name";
    let filterBy = "all";

    // --- COLLAPSIBLE SECTIONS ---
    let isAlertsOpen = true;
    let isStatsOpen = true;

    // --- ALERT VISIBILITY STATE ---
    let showConflictAlert = true;
    let showDependencyAlert = true;
    let showSodiumAlert = true;

    // --- EXPANDED CARDS STATE ---
    let expandedMods = new Set();

    // --- DEPENDENCY CHECKER ---
    let issues = [];
    let isChecking = false;
    let checkStatus = "";

    // --- IGNORED DEPENDENCIES ---
    // These are "environment" dependencies that typically aren't .jar files
    // so we filter them out to prevent false alarms.
    const IGNORED_DEPS = new Set([
        "fabric",
        "fabricloader",
        "minecraft",
        "java",
    ]);

    // --- REACTIVE STATS ---
    $: fabricCount = rawMods.filter((m) => m.loader === "Fabric").length;
    $: forgeCount = rawMods.filter((m) => m.loader === "NeoForge").length;
    $: conflictCount = targetLoader === "Fabric" ? forgeCount : fabricCount;
    $: hasSodium = rawMods.some((m) => m.name.toLowerCase().includes("sodium"));

    // Calculate total active alerts
    $: activeAlertCount =
        (conflictCount > 0 && showConflictAlert ? 1 : 0) +
        (issues.length > 0 && showDependencyAlert ? 1 : 0) +
        (hasSodium && showSodiumAlert ? 1 : 0);

    // --- FILTER LOGIC ---
    $: filteredMods = rawMods
        .filter((mod) => {
            const searchMatch =
                mod.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                mod.fileName.toLowerCase().includes(searchQuery.toLowerCase());

            let statusMatch = true;
            const isConflict =
                (targetLoader === "Fabric" && mod.loader === "NeoForge") ||
                (targetLoader === "NeoForge" && mod.loader === "Fabric");

            if (filterBy === "compatible") statusMatch = !isConflict;
            if (filterBy === "incompatible") statusMatch = isConflict;

            return searchMatch && statusMatch;
        })
        .sort((a, b) => {
            if (sortBy === "name") return a.name.localeCompare(b.name);
            if (sortBy === "loader") return a.loader.localeCompare(b.loader);
            if (sortBy === "status") {
                const isConflictA =
                    targetLoader === "Fabric" && a.loader === "NeoForge";
                const isConflictB =
                    targetLoader === "Fabric" && b.loader === "NeoForge";
                return isConflictA === isConflictB ? 0 : isConflictA ? -1 : 1;
            }
            return 0;
        });

    function toggleTheme() {
        darkMode = !darkMode;
        if (darkMode) document.body.classList.add("dark");
        else document.body.classList.remove("dark");
    }

    function toggleExpand(fileName) {
        const newSet = new Set(expandedMods);
        if (newSet.has(fileName)) {
            newSet.delete(fileName);
        } else {
            newSet.add(fileName);
        }
        expandedMods = newSet;
    }

    onMount(async () => {
        if (darkMode) document.body.classList.add("dark");
        document.addEventListener("contextmenu", (event) =>
            event.preventDefault(),
        );
        checkForAppUpdates();
    });

    async function checkForAppUpdates() {
        isCheckingUpdate = true;
        try {
            const res = await CheckForUpdates();
            if (res.updateAvailable) {
                updateAvailable = true;
                latestVersion = res.latestVersion;
                updateUrl = res.downloadUrl;
            }
        } catch (err) {
            console.error("Failed to check updates", err);
        }
        isCheckingUpdate = false;
    }

    function handleUpdateClick() {
        if (!updateAvailable) return;
        TriggerUpdate(updateUrl);
    }

    async function pickFolder() {
        let path = await SelectModFolder();
        if (path) {
            folderPath = path;
            handleScan();
        }
    }

    async function handleScan() {
        status = "Scanning...";
        issues = [];
        checkStatus = "";
        expandedMods = new Set();

        // RESET ALERTS
        showConflictAlert = true;
        showDependencyAlert = true;
        showSodiumAlert = true;
        isAlertsOpen = true;

        try {
            rawMods = await ScanMods(folderPath);
            status = "";
        } catch (err) {
            status = "Error: " + err;
        }
    }

    async function runDependencyCheck() {
        if (!folderPath) return;
        isChecking = true;
        checkStatus = "";
        issues = [];
        showDependencyAlert = true;

        try {
            const allIssues = await ValidateMods(folderPath);

            // Filter out the annoying environment dependencies
            issues = allIssues.filter((issue) => {
                const dep = issue.MissingDep.toLowerCase();
                return !IGNORED_DEPS.has(dep);
            });

            checkStatus = issues.length === 0 ? "clean" : "issues";
            if (issues.length > 0) isAlertsOpen = true;
        } catch (err) {
            console.error(err);
        }
        isChecking = false;
    }
</script>

<div class="app-container">
    <header>
        <div class="header-content">
            <div
                class="logo-section"
                on:click={resetApp}
                on:keydown={(e) => e.key === "Enter" && resetApp()}
                role="button"
                tabindex="0"
                title="Reset and clear workspace"
            >
                <img src={logoIcon} alt="Glacier Logo" class="logo-img" />
                <h1>Glacier Mod Manager</h1>
            </div>

            <div class="header-controls">
                <button
                    class="update-btn {updateAvailable ? 'available' : ''}"
                    on:click={handleUpdateClick}
                    disabled={!updateAvailable}
                    title={updateAvailable
                        ? `Update available: ${latestVersion}`
                        : "Glacier is up to date"}
                >
                    {#if isCheckingUpdate}
                        <i class="fa-solid fa-circle-notch fa-spin"></i>
                    {:else if updateAvailable}
                        <i class="fa-solid fa-cloud-arrow-down"></i>
                        <span>Update v{latestVersion}</span>
                    {:else}
                        <i class="fa-solid fa-check"></i>
                        <span>Up to Date</span>
                    {/if}
                </button>

                <button class="theme-toggle" on:click={toggleTheme}>
                    <i class="fa-solid {darkMode ? 'fa-sun' : 'fa-moon'}"></i>
                    <span>{darkMode ? "Light Mode" : "Dark Mode"}</span>
                </button>
            </div>
        </div>
    </header>

    <main>
        <div class="controls-card">
            <div class="control-col">
                <label
                    ><i class="fa-solid fa-crosshairs"></i> Target Loader</label
                >
                <div class="toggle-group">
                    <button
                        class="toggle-btn {targetLoader === 'Fabric'
                            ? 'active fabric'
                            : ''}"
                        on:click={() => (targetLoader = "Fabric")}
                    >
                        Fabric
                    </button>
                    <button
                        class="toggle-btn {targetLoader === 'NeoForge'
                            ? 'active forge'
                            : ''}"
                        on:click={() => (targetLoader = "NeoForge")}
                    >
                        NeoForge
                    </button>
                </div>
            </div>

            <div class="control-col wide">
                <label
                    ><i class="fa-regular fa-folder-open"></i> Mods Folder</label
                >
                <div class="input-group">
                    <div class="path-display">
                        {folderPath || "No folder selected..."}
                    </div>
                    <button class="btn btn-primary" on:click={pickFolder}
                        >Select</button
                    >
                </div>
            </div>
        </div>

        {#if (conflictCount > 0 && showConflictAlert) || (issues.length > 0 && checkStatus === "issues" && showDependencyAlert) || (hasSodium && showSodiumAlert)}
            <div class="section-container">
                <button
                    class="section-header-btn"
                    on:click={() => (isAlertsOpen = !isAlertsOpen)}
                >
                    <div class="title-group">
                        <i class="fa-solid fa-triangle-exclamation text-warning"
                        ></i>
                        <h3>Active Alerts ({activeAlertCount})</h3>
                    </div>
                    <i
                        class="fa-solid {isAlertsOpen
                            ? 'fa-chevron-up'
                            : 'fa-chevron-down'}"
                    ></i>
                </button>

                {#if isAlertsOpen}
                    <div class="alerts-body">
                        {#if conflictCount > 0 && showConflictAlert}
                            <div class="alert critical">
                                <i class="fa-solid fa-circle-xmark icon"></i>
                                <div class="content">
                                    <h4>Incompatibility Detected</h4>
                                    <p>
                                        {conflictCount}
                                        {targetLoader === "Fabric"
                                            ? "NeoForge"
                                            : "Fabric"} mods detected.
                                    </p>
                                </div>
                                <button
                                    class="close-alert"
                                    on:click={() => (showConflictAlert = false)}
                                    title="Dismiss"
                                >
                                    <i class="fa-solid fa-xmark"></i>
                                </button>
                            </div>
                        {/if}

                        {#if checkStatus === "issues" && showDependencyAlert}
                            <div class="alert danger-list">
                                <i class="fa-solid fa-puzzle-piece icon"></i>
                                <div class="content">
                                    <h4>
                                        Missing Dependencies ({issues.length})
                                    </h4>
                                    <div class="issue-grid">
                                        {#each issues as issue}
                                            <div class="issue-item">
                                                <span class="mod-name"
                                                    >{issue.ModName}</span
                                                >
                                                <span class="arrow"
                                                    ><i
                                                        class="fa-solid fa-arrow-right-long"
                                                    ></i></span
                                                >
                                                <span class="missing-dep"
                                                    >{issue.MissingDep}</span
                                                >
                                            </div>
                                        {/each}
                                    </div>
                                </div>
                                <button
                                    class="close-alert"
                                    on:click={() =>
                                        (showDependencyAlert = false)}
                                    title="Dismiss"
                                >
                                    <i class="fa-solid fa-xmark"></i>
                                </button>
                            </div>
                        {/if}

                        {#if hasSodium && showSodiumAlert}
                            <div class="alert warning">
                                <i class="fa-solid fa-triangle-exclamation icon"
                                ></i>
                                <div class="content">
                                    <h4>Sodium Detected</h4>
                                    <p>Check for OptiFine conflicts.</p>
                                </div>
                                <button
                                    class="close-alert"
                                    on:click={() => (showSodiumAlert = false)}
                                    title="Dismiss"
                                >
                                    <i class="fa-solid fa-xmark"></i>
                                </button>
                            </div>
                        {/if}
                    </div>
                {/if}
            </div>
        {/if}

        {#if rawMods.length > 0}
            <div class="section-container">
                <button
                    class="section-header-btn"
                    on:click={() => (isStatsOpen = !isStatsOpen)}
                >
                    <div class="title-group">
                        <i class="fa-solid fa-chart-simple text-blue"></i>
                        <h3>Quick Stats</h3>
                    </div>
                    <i
                        class="fa-solid {isStatsOpen
                            ? 'fa-chevron-up'
                            : 'fa-chevron-down'}"
                    ></i>
                </button>

                {#if isStatsOpen}
                    <div class="stats-bar">
                        <div class="stat-group">
                            <div class="stat-box">
                                <i class="fa-solid fa-cube"></i>
                                <div class="stat-info">
                                    <span class="val">{rawMods.length}</span>
                                    <span class="label">Total</span>
                                </div>
                            </div>
                            <div class="stat-box success">
                                <i class="fa-solid fa-check"></i>
                                <div class="stat-info">
                                    <span class="val"
                                        >{rawMods.length - conflictCount}</span
                                    >
                                    <span class="label">Safe</span>
                                </div>
                            </div>
                            {#if conflictCount > 0}
                                <div class="stat-box danger">
                                    <i class="fa-solid fa-xmark"></i>
                                    <div class="stat-info">
                                        <span class="val">{conflictCount}</span>
                                        <span class="label">Conflicts</span>
                                    </div>
                                </div>
                            {/if}
                        </div>

                        <button
                            class="btn btn-secondary check-btn"
                            on:click={runDependencyCheck}
                            disabled={isChecking}
                        >
                            {#if isChecking}
                                <i class="fa-solid fa-spinner fa-spin"></i> Checking...
                            {:else}
                                <i class="fa-solid fa-shield-halved"></i> Validate
                            {/if}
                        </button>
                    </div>
                {/if}
            </div>

            <div class="toolbar">
                <div class="search-box">
                    <i class="fa-solid fa-magnifying-glass"></i>
                    <input
                        type="text"
                        placeholder="Search mods..."
                        bind:value={searchQuery}
                    />
                </div>

                <div class="filter-group">
                    <div class="select-wrapper">
                        <i class="fa-solid fa-filter"></i>
                        <select bind:value={filterBy}>
                            <option value="all">All</option>
                            <option value="compatible">Compatible</option>
                            <option value="incompatible">Conflicts</option>
                        </select>
                    </div>

                    <div class="select-wrapper">
                        <i class="fa-solid fa-arrow-down-a-z"></i>
                        <select bind:value={sortBy}>
                            <option value="name">Name</option>
                            <option value="loader">Loader</option>
                            <option value="status">Status</option>
                        </select>
                    </div>

                    <button
                        class="view-toggle-btn"
                        on:click={() =>
                            (viewMode = viewMode === "grid" ? "list" : "grid")}
                    >
                        <i
                            class="fa-solid {viewMode === 'grid'
                                ? 'fa-list'
                                : 'fa-border-all'}"
                        ></i>
                    </button>
                </div>
            </div>

            <div class="grid {viewMode === 'list' ? 'list-view' : ''}">
                {#each filteredMods as mod (mod.fileName)}
                    {@const isConflict =
                        (targetLoader === "Fabric" &&
                            mod.loader === "NeoForge") ||
                        (targetLoader === "NeoForge" &&
                            mod.loader === "Fabric")}
                    {@const isExpanded = expandedMods.has(mod.fileName)}

                    <div
                        class="mod-card {isConflict
                            ? 'conflict'
                            : ''} {mod.loader.toLowerCase()} {isExpanded
                            ? 'expanded'
                            : ''}"
                    >
                        {#if viewMode === "grid"}
                            <div class="card-header-grid">
                                <h3 class="mod-title" title={mod.name}>
                                    {mod.name}
                                </h3>
                                <span class="version-pill">{mod.version}</span>
                            </div>

                            <div class="card-badges">
                                <span class="badge {mod.loader.toLowerCase()}"
                                    >{mod.loader}</span
                                >
                            </div>

                            {#if isConflict}
                                <div class="error-message">
                                    <i class="fa-solid fa-triangle-exclamation"
                                    ></i>
                                    Incompatible with {targetLoader}
                                </div>
                            {/if}

                            <p
                                class="mod-description {isExpanded
                                    ? ''
                                    : 'clamp-desc'}"
                                title={mod.description}
                            >
                                {mod.description || "No description provided."}
                            </p>

                            <button
                                class="show-more-btn"
                                on:click|stopPropagation={() =>
                                    toggleExpand(mod.fileName)}
                            >
                                {isExpanded ? "Show Less" : "Show More"}
                            </button>

                            <div class="card-footer-grid">
                                <div class="ext-links">
                                    <a
                                        href="https://modrinth.com/mods?q={mod.name}"
                                        target="_blank"
                                        class="ext-btn-text"
                                        title="Search Modrinth"
                                    >
                                        <i class="fa-solid fa-hammer"></i> Modrinth
                                    </a>
                                    <a
                                        href="https://www.curseforge.com/minecraft/search?search={mod.name}"
                                        target="_blank"
                                        class="ext-btn-text"
                                        title="Search CurseForge"
                                    >
                                        <i class="fa-solid fa-fire"></i> CurseForge
                                    </a>
                                </div>
                            </div>

                            {#if isExpanded}
                                <div class="filename-grid">
                                    <i class="fa-regular fa-file"></i>
                                    {mod.fileName}
                                </div>
                            {/if}
                        {:else}
                            <div class="col-title">
                                <h3 class="mod-title" title={mod.name}>
                                    {mod.name}
                                </h3>
                                <div>
                                    <span
                                        class="badge {mod.loader.toLowerCase()}"
                                        >{mod.loader}</span
                                    >
                                    <span class="version-pill"
                                        >{mod.version}</span
                                    >
                                    {#if isConflict}
                                        <span class="status-badge conflict"
                                            >Conflict</span
                                        >
                                        <span class="status-badge incompatible"
                                            >Incompatible</span
                                        >
                                    {/if}
                                </div>
                            </div>

                            <div class="col-desc">
                                <p class="mod-description truncate-list">
                                    {mod.description ||
                                        "No description provided."}
                                </p>
                            </div>

                            <div class="col-actions">
                                <a
                                    href="https://modrinth.com/mods?q={mod.name}"
                                    target="_blank"
                                    class="ext-btn-text small"
                                    title="Search Modrinth"
                                >
                                    <i class="fa-solid fa-hammer"></i> Modrinth
                                </a>
                                <a
                                    href="https://www.curseforge.com/minecraft/search?search={mod.name}"
                                    target="_blank"
                                    class="ext-btn-text small"
                                    title="Search CurseForge"
                                >
                                    <i class="fa-solid fa-fire"></i> CurseForge
                                </a>
                            </div>
                        {/if}
                    </div>
                {/each}
            </div>

            <div class="footer-count">
                Showing {filteredMods.length} of {rawMods.length} mods
            </div>
        {:else}
            <div class="empty-state">
                <div class="icon-container">
                    <i class="fa-solid fa-gears"></i>
                </div>

                <h2>It's looking a little empty here...</h2>
                <p>Select your mods folder above to get started!</p>
            </div>
        {/if}
    </main>

    <footer>
        <div class="footer-content">
            <span class="copyright">© 2026 BySnowden</span>
            <div class="social-links">
                <a
                    href="https://github.com/BySnowden"
                    target="_blank"
                    class="social-btn github"
                >
                    <i class="fa-brands fa-github"></i>
                </a>
                <a
                    href="https://ko-fi.com/bysnowden"
                    target="_blank"
                    class="social-btn kofi"
                >
                    <i class="fa-solid fa-mug-hot"></i> Support
                </a>
            </div>
        </div>
    </footer>
</div>

<style>
    /* ============================================
      GLOBAL RESETS & BASE STYLES
      ============================================ */
    :global(body) {
        user-select: none;
        cursor: default;
    }

    .path-display,
    .alert p,
    .mod-title,
    .filename-grid,
    .filename-list,
    .mod-description {
        user-select: text;
        cursor: text;
    }

    button,
    select,
    a {
        cursor: pointer;
    }

    /* ============================================
      LAYOUT - MAIN CONTAINER
      ============================================ */
    .app-container {
        min-height: 100vh;
        display: flex;
        flex-direction: column;
        background: var(--background);
    }

    main {
        max-width: 1200px;
        margin: 0 auto;
        padding: 20px;
        width: 100%;
        box-sizing: border-box;
        flex: 1;
    }

    /* ============================================
      HEADER
      ============================================ */
    header {
        background: var(--card);
        border-bottom: 1px solid var(--card-border);
        padding: 16px 0;
        margin-bottom: 24px;
    }

    .header-content {
        max-width: 1200px;
        margin: 0 auto;
        padding: 0 20px;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .header-controls {
        display: flex;
        gap: 12px;
        align-items: center;
    }

    .logo-section {
        display: flex;
        align-items: center;
        gap: 12px;
        cursor: pointer;
        user-select: none;
        transition: opacity 0.2s;
    }

    .logo-section:hover {
        opacity: 0.8;
    }

    .logo-section:active {
        transform: scale(0.98);
    }

    .logo-img {
        height: 32px;
        width: auto;
        object-fit: contain;
    }

    h1 {
        margin: 0;
        font-size: 1.5rem;
        font-weight: 700;
    }

    .theme-toggle {
        background: transparent;
        border: 1px solid var(--card-border);
        color: var(--muted-fg);
        padding: 8px 16px;
        border-radius: 6px;
        font-weight: 600;
        display: flex;
        align-items: center;
        gap: 8px;
        transition: all 0.2s;
    }

    .theme-toggle:hover {
        background: var(--card-border);
        color: var(--foreground);
    }

    .update-btn {
        background: transparent;
        border: 1px solid var(--card-border);
        color: var(--muted-fg);
        padding: 8px 16px;
        border-radius: 6px;
        font-weight: 600;
        display: flex;
        align-items: center;
        gap: 8px;
        transition: all 0.3s ease;
        opacity: 0.5; /* Faded out by default */
        cursor: default;
    }

    .update-btn.available {
        background: var(--primary); /* Or use #22c55e for a specific green */
        color: white; /* Force white text on green button */
        border-color: var(--primary);
        opacity: 1;
        cursor: pointer;
        box-shadow: 0 0 10px rgba(var(--primary-rgb), 0.4);
        animation: pulse-green 2s infinite;
    }

    .update-btn.available:hover {
        filter: brightness(1.1);
        transform: translateY(-1px);
    }

    @keyframes pulse-green {
        0% {
            box-shadow: 0 0 0 0 rgba(34, 197, 94, 0.4);
        }
        70% {
            box-shadow: 0 0 0 6px rgba(34, 197, 94, 0);
        }
        100% {
            box-shadow: 0 0 0 0 rgba(34, 197, 94, 0);
        }
    }

    /* ============================================
      CONTROLS CARD
      ============================================ */
    .controls-card {
        background: var(--card);
        border: 1px solid var(--card-border);
        border-radius: 8px;
        padding: 20px;
        display: flex;
        gap: 24px;
        margin-bottom: 20px;
    }

    .control-col {
        flex: 1;
        display: flex;
        flex-direction: column;
        gap: 8px;
    }

    .control-col.wide {
        flex: 2;
    }

    label {
        font-size: 0.85rem;
        font-weight: 600;
        color: var(--muted-fg);
        display: flex;
        align-items: center;
        gap: 6px;
    }

    .toggle-group {
        display: flex;
        background: var(--background);
        padding: 4px;
        border-radius: 6px;
        border: 1px solid var(--card-border);
    }

    .toggle-btn {
        flex: 1;
        padding: 8px;
        border: none;
        background: transparent;
        color: var(--muted-fg);
        font-weight: 600;
        border-radius: 4px;
        transition: all 0.2s;
    }

    .toggle-btn.active {
        background: var(--card);
        color: var(--foreground);
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
    }

    .toggle-btn.active.fabric {
        color: var(--primary);
    }

    .toggle-btn.active.forge {
        color: #f97316;
    }

    .input-group {
        display: flex;
        gap: 10px;
    }

    .path-display {
        flex: 1;
        background: var(--background);
        border: 1px solid var(--card-border);
        padding: 8px 12px;
        border-radius: 6px;
        font-family: monospace;
        font-size: 0.9em;
        color: var(--muted-fg);
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .btn-primary {
        background: var(--primary);
        color: white;
        border: none;
        padding: 0 16px;
        border-radius: 6px;
        font-weight: 600;
    }

    /* ============================================
      COLLAPSIBLE SECTIONS
      ============================================ */
    .section-container {
        background: var(--card);
        border: 1px solid var(--card-border);
        border-radius: 8px;
        margin-bottom: 20px;
        overflow: hidden;
    }

    .section-header-btn {
        width: 100%;
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 12px 20px;
        background: transparent;
        border: none;
        color: var(--foreground);
        font-size: 1rem;
        border-bottom: 1px solid transparent;
        transition: background 0.2s;
    }

    .section-header-btn:hover {
        background: rgba(255, 255, 255, 0.02);
    }

    .title-group {
        display: flex;
        align-items: center;
        gap: 10px;
    }

    .title-group h3 {
        margin: 0;
        font-size: 0.95rem;
        font-weight: 600;
    }

    /* ============================================
      ALERTS SECTION
      ============================================ */
    .alerts-body {
        padding: 16px;
        border-top: 1px solid var(--card-border);
    }

    .alert {
        padding: 12px;
        padding-right: 40px;
        border-radius: 6px;
        display: flex;
        gap: 12px;
        margin-bottom: 8px;
        position: relative;
    }

    .alert:last-child {
        margin-bottom: 0;
    }

    .alert.critical {
        background: var(--alert-critical-bg);
        border: 1px solid var(--alert-critical-border);
        color: var(--alert-critical-text);
    }

    .alert.warning {
        background: var(--alert-warning-bg);
        border: 1px solid var(--alert-warning-border);
        color: var(--alert-warning-text);
    }

    .alert.danger-list {
        background: var(--alert-critical-bg);
        border: 1px solid var(--alert-critical-border);
        color: var(--alert-critical-text);
    }

    .alert .icon {
        font-size: 1.1rem;
        margin-top: 2px;
    }

    .alert h4 {
        margin: 0 0 4px 0;
        font-size: 0.95rem;
        font-weight: 700;
        color: inherit;
    }

    .alert p {
        margin: 0;
        font-size: 0.85rem;
        opacity: 1;
        color: inherit;
    }

    .close-alert {
        position: absolute;
        top: 10px;
        right: 10px;
        background: transparent;
        border: none;
        color: inherit;
        opacity: 0.6;
        padding: 4px;
        font-size: 1rem;
        cursor: pointer;
        transition: opacity 0.2s;
    }

    .close-alert:hover {
        opacity: 1;
    }

    /* ============================================
      ISSUE GRID (ALERTS SUBSECTION)
      ============================================ */
    .issue-grid {
        display: grid;
        gap: 6px;
        margin-top: 8px;
    }

    .issue-item {
        background: rgba(0, 0, 0, 0.1);
        padding: 8px 12px;
        border-radius: 4px;
        font-size: 0.85rem;
        display: flex;
        align-items: center;
        gap: 8px;
        border: 1px solid var(--alert-critical-border);
    }

    .mod-name {
        font-weight: bold;
        color: var(--alert-critical-text);
    }

    .missing-dep {
        color: var(--alert-critical-text);
        font-weight: bold;
    }

    .issue-item .arrow {
        color: var(--alert-critical-text);
        opacity: 0.6;
    }

    /* ============================================
      STATS BAR
      ============================================ */
    .stats-bar {
        padding: 16px;
        border-top: 1px solid var(--card-border);
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .stat-group {
        display: flex;
        gap: 20px;
    }

    .stat-box {
        display: flex;
        align-items: center;
        gap: 10px;
    }

    .stat-box i {
        font-size: 1.2rem;
        color: var(--muted-fg);
    }

    .stat-box.success i {
        color: var(--success);
    }

    .stat-box.danger i {
        color: var(--destructive);
    }

    .stat-info {
        display: flex;
        flex-direction: column;
        line-height: 1.1;
    }

    .stat-info .val {
        font-weight: 700;
        font-size: 1.1rem;
    }

    .stat-info .label {
        font-size: 0.75rem;
        color: var(--muted-fg);
        text-transform: uppercase;
        font-weight: 600;
    }

    .btn-secondary {
        background: var(--background);
        border: 1px solid var(--card-border);
        color: var(--foreground);
        padding: 8px 16px;
        border-radius: 6px;
        font-weight: 600;
        transition: all 0.2s;
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .btn-secondary:hover:not(:disabled) {
        border-color: var(--primary);
        color: var(--primary);
    }

    /* ============================================
      TOOLBAR & FILTERS
      ============================================ */
    .toolbar {
        display: flex;
        gap: 16px;
        margin-bottom: 16px;
        align-items: center;
    }

    .search-box {
        flex: 1;
        position: relative;
    }

    .search-box i {
        position: absolute;
        left: 12px;
        top: 50%;
        transform: translateY(-50%);
        color: var(--muted-fg);
        font-size: 0.9rem;
    }

    .search-box input {
        width: 100%;
        background: var(--card);
        border: 1px solid var(--card-border);
        padding: 10px 10px 10px 36px;
        border-radius: 6px;
        color: var(--foreground);
        outline: none;
        box-sizing: border-box;
    }

    .search-box input:focus {
        border-color: var(--primary);
    }

    .filter-group {
        display: flex;
        gap: 10px;
    }

    .select-wrapper {
        position: relative;
    }

    .select-wrapper i {
        position: absolute;
        left: 10px;
        top: 50%;
        transform: translateY(-50%);
        color: var(--muted-fg);
        font-size: 0.8rem;
        pointer-events: none;
    }

    select {
        background: var(--card);
        border: 1px solid var(--card-border);
        color: var(--foreground);
        padding: 9px 12px 9px 30px;
        border-radius: 6px;
        outline: none;
        appearance: none;
        font-size: 0.9rem;
        cursor: pointer;
    }

    select:hover {
        border-color: var(--muted-fg);
    }

    .view-toggle-btn {
        width: 38px;
        background: var(--card);
        border: 1px solid var(--card-border);
        color: var(--muted-fg);
        border-radius: 6px;
        font-size: 1rem;
        display: grid;
        place-items: center;
    }

    .view-toggle-btn:hover {
        color: var(--foreground);
        border-color: var(--muted-fg);
    }

    /* ============================================
      BADGES & PILLS
      ============================================ */
    .version-pill {
        font-size: 0.75rem;
        color: var(--muted-fg);
        font-family: monospace;
        background: rgba(255, 255, 255, 0.05);
        padding: 4px 8px;
        border-radius: 4px;
        width: fit-content;
        line-height: 1.2;
    }

    .badge {
        font-size: 0.65rem;
        padding: 4px 8px;
        border-radius: 3px;
        font-weight: 700;
        text-transform: uppercase;
        letter-spacing: 0.5px;
        width: fit-content;
        line-height: 1.2;
    }

    .badge.fabric {
        background: rgba(59, 130, 246, 0.2);
        color: #60a5fa;
    }

    .badge.neoforge {
        background: rgba(249, 115, 22, 0.2);
        color: #fb923c;
    }

    .status-badge {
        font-size: 0.65rem;
        padding: 4px 8px;
        border-radius: 4px;
        font-weight: 700;
        text-transform: uppercase;
        white-space: nowrap;
        display: inline-block;
    }

    .status-badge.conflict {
        background: rgba(220, 38, 38, 0.2);
        color: #f87171;
    }

    .status-badge.incompatible {
        background: rgba(220, 38, 38, 0.2);
        color: #f87171;
    }

    /* ============================================
      GRID SYSTEM
      ============================================ */
    .grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
        gap: 16px;
    }

    .grid.list-view {
        grid-template-columns: 1fr;
        gap: 12px;
    }

    .grid.list-view .mod-card {
        display: grid;
        grid-template-columns: 1fr;
        align-items: stretch;
        padding: 20px 24px;
        gap: 12px;
        min-height: 64px;
    }

    /* ============================================
      MOD CARD - BASE STYLES
      ============================================ */
    .mod-card {
        background: var(--card);
        border: 1px solid var(--card-border);
        border-radius: 8px;
        padding: 16px;
        display: flex;
        flex-direction: column;
        transition:
            transform 0.2s,
            box-shadow 0.2s,
            border-color 0.2s;
        position: relative;
        overflow: hidden;
        gap: 12px;
    }

    .mod-card:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        border-color: var(--muted-fg);
    }

    .mod-card.conflict {
        border-color: var(--destructive);
        background: rgba(40, 0, 0, 0.2);
    }

    /* ============================================
      MOD CARD - TEXT CONTENT
      ============================================ */
    .mod-title {
        font-size: 1rem;
        margin: 0;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        font-weight: 700;
    }

    .mod-description {
        font-size: 0.85rem;
        color: var(--muted-fg);
        margin: 0;
        line-height: 1.4;
    }

    .mod-description.clamp-desc {
        display: -webkit-box;
        -webkit-line-clamp: 2;
        line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }

    .mod-description.truncate-list {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
    }

    .show-more-btn {
        background: transparent;
        border: none;
        color: var(--primary);
        font-size: 0.75rem;
        font-weight: 600;
        padding: 0;
        margin-top: -4px;
        text-align: left;
        width: fit-content;
    }

    .show-more-btn:hover {
        text-decoration: underline;
    }

    /* ============================================
      MOD CARD - GRID VIEW
      ============================================ */
    .card-header-grid {
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .card-badges {
        display: flex;
        gap: 8px;
    }

    .error-message {
        display: flex;
        align-items: center;
        gap: 6px;
        color: #f87171;
        font-size: 0.8rem;
        font-weight: 600;
    }

    .card-footer-grid {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-top: auto;
    }

    .ext-links {
        display: flex;
        gap: 8px;
    }

    .filename-grid {
        font-size: 0.7rem;
        color: var(--muted-fg);
        opacity: 0.6;
        display: flex;
        align-items: center;
        gap: 4px;
        margin-top: 8px;
        border-top: 1px solid var(--card-border);
        padding-top: 8px;
    }

    /* ============================================
      MOD CARD - LIST VIEW
      ============================================ */
    .col-title {
        display: flex;
        flex-direction: column;
        align-items: flex-start;
        gap: 8px;
        overflow: hidden;
        flex: 1;
    }

    .col-title .mod-title {
        font-size: 1.05rem;
        font-weight: 700;
    }

    .col-title > div {
        display: flex;
        gap: 8px;
        flex-wrap: wrap;
        align-items: center;
    }

    .col-desc {
        display: block;
        margin-top: 8px;
    }

    .col-desc .mod-description {
        margin: 0;
    }

    .col-actions {
        display: flex;
        gap: 10px;
        margin-top: 12px;
    }

    /* ============================================
      EXTERNAL BUTTONS
      ============================================ */
    .ext-btn-text {
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 6px 12px;
        border-radius: 4px;
        font-size: 0.8rem;
        color: var(--foreground);
        background: var(--background);
        border: 1px solid var(--card-border);
        transition: all 0.2s;
        text-decoration: none;
        font-weight: 600;
        white-space: nowrap;
    }

    .ext-btn-text:hover {
        background: var(--card-border);
        border-color: var(--primary);
    }

    .ext-btn-text.small {
        padding: 4px 8px;
        font-size: 0.75rem;
    }

    /* ============================================
      FOOTER
      ============================================ */
    .footer-count {
        text-align: center;
        font-size: 0.8rem;
        color: var(--muted-fg);
        margin-top: 20px;
    }

    footer {
        margin-top: auto;
        padding: 20px 0;
        border-top: 1px solid var(--card-border);
        background: var(--card);
    }

    .footer-content {
        max-width: 1200px;
        margin: 0 auto;
        padding: 0 20px;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .copyright {
        font-size: 0.85rem;
        color: var(--muted-fg);
    }

    .social-links {
        display: flex;
        gap: 10px;
    }

    .social-btn {
        padding: 6px 12px;
        border-radius: 4px;
        font-size: 0.85rem;
        text-decoration: none;
        font-weight: 600;
        display: flex;
        align-items: center;
        gap: 6px;
        transition: all 0.2s;
    }

    .social-btn.github {
        background: var(--background);
        color: var(--foreground);
        border: 1px solid var(--card-border);
    }

    .social-btn.kofi {
        background: #ff5e5b;
        color: white;
    }

    .social-btn:hover {
        transform: translateY(-1px);
        filter: brightness(1.1);
    }

    /* ============================================
      EMPTY STATE
      ============================================ */
    .empty-state {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 80px 20px;
        text-align: center;
        color: var(--muted-fg);
        animation: fadeIn 0.5s ease-out;
    }

    .icon-container {
        font-size: 5rem;
        margin-bottom: 24px;
        opacity: 0.2;
        transition: transform 0.3s ease;
    }

    .empty-state:hover .icon-container {
        transform: scale(1.1) rotate(-5deg);
        opacity: 0.4;
        color: var(--primary);
    }

    .empty-state h2 {
        font-size: 1.5rem;
        font-weight: 700;
        margin: 0 0 8px 0;
        color: var(--foreground);
        opacity: 0.8;
    }

    .empty-state p {
        font-size: 0.95rem;
        max-width: 400px;
        line-height: 1.5;
    }

    /* ============================================
      ANIMATIONS
      ============================================ */
    @keyframes fadeIn {
        from {
            opacity: 0;
            transform: translateY(10px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    @keyframes pulse {
        0% {
            opacity: 1;
        }
        50% {
            opacity: 0.5;
        }
        100% {
            opacity: 1;
        }
    }
</style>
