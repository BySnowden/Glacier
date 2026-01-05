<script>
    import { onMount } from "svelte";
    import Header from "./components/ui/Header.svelte";
    import UpdateManager from "./components/ui/UpdateManager.svelte";
    import ModControls from "./components/panels/ModControls.svelte";
    import AlertPanel from "./components/panels/AlertPanel.svelte";
    import StatsPanel from "./components/panels/StatsPanel.svelte";
    import ModToolbar from "./components/ui/ModToolbar.svelte";
    import ModCard from "./components/ui/ModCard.svelte";

    import {
        darkMode,
        folderPath,
        rawMods,
        filteredMods,
        viewMode,
        expandedMods,
        targetLoader,
        issues,
        isChecking,
        checkStatus,
        IGNORED_DEPS,
        toggleTheme,
        toggleExpand,
    } from "./stores/index.js";

    import { ValidateMods } from "../wailsjs/go/main/App";
    import * as runtime from "../wailsjs/runtime/runtime";
    import logoIcon from "./assets/images/icon-transparent.png";

    import { resetApp } from "./stores/index.js";

    function openModrinthSearch(modName) {
        const searchUrl = `https://modrinth.com/mods?q=${encodeURIComponent(modName)}`;
        runtime.BrowserOpenURL(searchUrl);
    }

    function openCurseForgeSearch(modName) {
        const searchUrl = `https://www.curseforge.com/minecraft/search?search=${encodeURIComponent(modName)}`;
        runtime.BrowserOpenURL(searchUrl);
    }

    function openGitHub() {
        runtime.BrowserOpenURL("https://github.com/BySnowden/Glacier");
    }

    function openKofi() {
        runtime.BrowserOpenURL("https://ko-fi.com/bysnowden");
    }

    // Dependency check function
    async function runDependencyCheck() {
        if (!$folderPath) {
            console.log("No folder selected");
            return;
        }

        isChecking.set(true);
        checkStatus.set("checking");

        try {
            const allIssues = await ValidateMods($folderPath);

            // Filter out ignored dependencies
            const filteredIssues = allIssues.filter((issue) => {
                const dep = issue.MissingDep.toLowerCase();
                return !IGNORED_DEPS.has(dep);
            });

            issues.set(filteredIssues);
            checkStatus.set(filteredIssues.length > 0 ? "issues" : "clean");
        } catch (error) {
            console.error("Dependency check failed:", error);
            checkStatus.set("error");
        }

        isChecking.set(false);
    }

    onMount(() => {
        // Set initial theme
        if (typeof document !== "undefined") {
            document.documentElement.setAttribute(
                "data-theme",
                $darkMode ? "dark" : "light",
            );
            if ($darkMode) document.body.classList.add("dark");
            else document.body.classList.remove("dark");
        }

        // Prevent context menu
        document.addEventListener("contextmenu", (event) =>
            event.preventDefault(),
        );
    });
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
                <h1>Glacier</h1>
            </div>
            <div class="header-controls">
                <UpdateManager />
                <button class="theme-toggle" on:click={toggleTheme}>
                    <i class="fa-solid {$darkMode ? 'fa-sun' : 'fa-moon'}"></i>
                    <span>{$darkMode ? "Light" : "Dark"}</span>
                </button>
            </div>
        </div>
    </header>

    <main>
        <ModControls />

        <AlertPanel />

        <StatsPanel />

        {#if $rawMods.length > 0}
            <ModToolbar />

            <div class="grid {$viewMode === 'list' ? 'list-view' : ''}">
                {#each $filteredMods as mod (mod.fileName)}
                    {@const isExpanded = $expandedMods.has(mod.fileName)}
                    <div
                        class="mod-card {mod.loader.toLowerCase()} {$targetLoader &&
                        mod.loader !== $targetLoader
                            ? 'conflict'
                            : ''}"
                        class:expanded={isExpanded}
                    >
                        {#if $viewMode === "grid"}
                            <div class="card-header-grid">
                                <h3 class="mod-title" title={mod.name}>
                                    {mod.name}
                                </h3>
                                <span class="version-pill">v{mod.version}</span>
                            </div>

                            <div class="card-badges">
                                <span
                                    class="badge {mod.loader.toLowerCase()}"
                                    title="Mod Loader"
                                >
                                    {mod.loader}
                                </span>
                            </div>

                            {#if $targetLoader && mod.loader !== $targetLoader}
                                <div class="error-message">
                                    <i
                                        class="fa-solid fa-triangle-exclamation"
                                        style="margin-right: 8px;"
                                    ></i>
                                    Incompatible with {$targetLoader}
                                </div>
                            {/if}

                            <p
                                class="mod-description {isExpanded
                                    ? ''
                                    : 'clamp-desc'}"
                            >
                                {mod.description || "No description available."}
                            </p>

                            <button
                                class="show-more-btn"
                                on:click={() => toggleExpand(mod.fileName)}
                            >
                                {isExpanded ? "Show Less ▲" : "Show More ▼"}
                            </button>

                            <div class="card-footer-grid">
                                <div class="ext-links">
                                    <button
                                        class="ext-btn-text"
                                        title="Search on Modrinth"
                                        on:click={() =>
                                            openModrinthSearch(mod.name)}
                                    >
                                        <i class="fa-solid fa-hammer"></i>
                                        Modrinth
                                    </button>
                                    <button
                                        class="ext-btn-text"
                                        title="Search on CurseForge"
                                        on:click={() =>
                                            openCurseForgeSearch(mod.name)}
                                    >
                                        <i class="fa-solid fa-fire"></i>
                                        CurseForge
                                    </button>
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
                                        title="Mod Loader"
                                    >
                                        {mod.loader}
                                    </span>
                                    <span class="version-pill" title="Version">
                                        v{mod.version}
                                    </span>
                                    {#if $targetLoader && mod.loader !== $targetLoader}
                                        <span
                                            class="status-badge conflict"
                                            title="Loader Conflict"
                                        >
                                            Conflict
                                        </span>
                                        <span
                                            class="status-badge incompatible"
                                            title="Incompatible"
                                        >
                                            Incompatible
                                        </span>
                                    {/if}
                                </div>
                            </div>

                            <div class="col-desc">
                                <p class="mod-description truncate-list">
                                    {mod.description ||
                                        "No description available."}
                                </p>
                            </div>

                            <div class="col-actions">
                                <button
                                    class="ext-btn-text small"
                                    title="Search on Modrinth"
                                    on:click={() =>
                                        openModrinthSearch(mod.name)}
                                >
                                    <i class="fa-solid fa-hammer"></i>
                                    Modrinth
                                </button>
                                <button
                                    class="ext-btn-text small"
                                    title="Search on CurseForge"
                                    on:click={() =>
                                        openCurseForgeSearch(mod.name)}
                                >
                                    <i class="fa-solid fa-fire"></i>
                                    CurseForge
                                </button>
                            </div>
                        {/if}
                    </div>
                {/each}
            </div>

            <div class="footer-count">
                Showing {$filteredMods.length} of {$rawMods.length} mods
            </div>
        {:else}
            <div class="empty-state">
                <div class="icon-container">
                    <i class="fa-solid fa-gears"></i>
                </div>
                <h2>No Mods Found</h2>
                <p>Select a folder and scan for mods to get started.</p>
            </div>
        {/if}
    </main>

    <footer>
        <div class="footer-content">
            <span class="copyright">© 2026 Glacier Mod Manager</span>
            <div class="social-links">
                <button
                    class="social-btn github"
                    title="GitHub Repository"
                    on:click={openGitHub}
                >
                    <i class="fab fa-github"></i>
                </button>
                <button
                    class="social-btn kofi"
                    title="Support on Ko-fi"
                    on:click={openKofi}
                >
                    <i class="fas fa-coffee"></i>
                </button>
            </div>
        </div>
    </footer>
</div>

<style>
    /* Global styles */
    :global(body) {
        margin: 0;
        font-family:
            "Inter",
            -apple-system,
            BlinkMacSystemFont,
            "Segoe UI",
            system-ui,
            sans-serif;
        background: var(--background);
        color: var(--foreground);
    }

    .path-display,
    .alert p,
    .mod-title,
    .filename-grid,
    .filename-list,
    .mod-description {
        word-break: break-word;
        overflow-wrap: break-word;
        hyphens: auto;
    }

    button,
    a {
        font-family: inherit;
        transition: all 0.2s ease;
    }

    /* Main layout */
    .app-container {
        min-height: 100vh;
        display: flex;
        flex-direction: column;
        background: var(--background);
    }

    main {
        flex: 1;
        padding: 20px;
        max-width: 1200px;
        margin: 0 auto;
        width: 100%;
        box-sizing: border-box;
    }

    /* Header styles */
    header {
        background: var(--card);
        border-bottom: 1px solid var(--card-border);
        padding: 16px 20px;
        margin-bottom: 20px;
    }

    .header-content {
        max-width: 1200px;
        margin: 0 auto;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .header-controls {
        display: flex;
        align-items: center;
        gap: 16px;
    }

    .logo-section {
        display: flex;
        align-items: center;
        gap: 12px;
        cursor: pointer;
        user-select: none;
        padding: 8px 12px;
        border-radius: 8px;
        transition: all 0.2s ease;
    }

    .logo-section:hover {
        background: var(--muted);
    }

    .logo-section:active {
        transform: scale(0.98);
    }

    .logo-img {
        width: 32px;
        height: 32px;
        object-fit: contain;
        background: var(--primary);
        border-radius: 8px;
        padding: 6px;
        box-shadow: 0 2px 8px rgba(37, 99, 235, 0.2);
    }

    h1 {
        margin: 0;
        font-size: 1.5rem;
        font-weight: 700;
    }

    .theme-toggle {
        background: var(--card);
        border: 1px solid var(--card-border);
        color: var(--foreground);
        padding: 8px 12px;
        border-radius: 8px;
        cursor: pointer;
        font-size: 0.85rem;
        display: flex;
        align-items: center;
        gap: 6px;
        transition: all 0.2s ease;
    }

    .theme-toggle:hover {
        background: var(--muted);
        border-color: var(--primary);
    }

    /* Mod grid styles */
    .version-pill {
        background: var(--muted);
        color: var(--muted-fg);
        padding: 4px 8px;
        border-radius: 12px;
        font-size: 0.7rem;
        font-weight: 600;
        border: 1px solid var(--card-border);
        white-space: nowrap;
    }

    .badge {
        padding: 4px 8px;
        border-radius: 4px;
        font-size: 0.7rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
        white-space: nowrap;
    }

    .badge.fabric {
        background: var(--fabric-bg);
        color: var(--fabric-text);
    }

    .badge.neoforge {
        background: var(--forge-bg);
        color: var(--forge-text);
    }

    .status-badge {
        padding: 2px 6px;
        border-radius: 3px;
        font-size: 0.65rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    .status-badge.conflict {
        background: var(--destructive);
        color: white;
    }

    .status-badge.incompatible {
        background: var(--destructive);
        color: white;
    }

    /* Grid layout */
    .grid {
        display: grid;
        grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
        gap: 20px;
        margin-bottom: 24px;
    }

    .grid.list-view {
        grid-template-columns: 1fr;
        gap: 12px;
    }

    .grid.list-view .mod-card {
        display: grid;
        grid-template-columns: 2fr 3fr auto;
        align-items: center;
        gap: 20px;
        padding: 16px 20px;
    }

    /* Mod card styles */
    .mod-card {
        background: var(--card);
        border: 1px solid var(--card-border);
        border-radius: var(--radius);
        padding: 16px;
        display: flex;
        flex-direction: column;
        gap: 12px;
        transition: all 0.2s ease;
        position: relative;
    }

    .mod-card:hover {
        border-color: var(--primary);
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    }

    .mod-card.conflict {
        border-color: var(--destructive);
        background: rgba(239, 68, 68, 0.05);
    }

    /* Card content styles */
    .mod-title {
        margin: 0;
        font-size: 1.1rem;
        font-weight: 600;
        color: var(--foreground);
        line-height: 1.3;
    }

    .mod-description {
        color: var(--muted-fg);
        font-size: 0.9rem;
        line-height: 1.5;
        margin: 0;
    }

    .mod-description.clamp-desc {
        display: -webkit-box;
        -webkit-line-clamp: 3;
        line-clamp: 3;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }

    .mod-description.truncate-list {
        display: -webkit-box;
        -webkit-line-clamp: 2;
        line-clamp: 2;
        -webkit-box-orient: vertical;
        overflow: hidden;
    }

    .show-more-btn {
        background: none;
        border: none;
        color: var(--primary);
        font-size: 0.8rem;
        font-weight: 600;
        cursor: pointer;
        padding: 4px 0;
        text-align: left;
        align-self: flex-start;
    }

    .show-more-btn:hover {
        text-decoration: underline;
    }

    /* Grid view specific styles */
    .card-header-grid {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        gap: 12px;
    }

    .card-badges {
        display: flex;
        gap: 8px;
    }

    .error-message {
        color: var(--destructive);
        font-size: 0.85rem;
        font-weight: 500;
        display: flex;
        align-items: center;
        padding: 8px 12px;
        background: rgba(239, 68, 68, 0.1);
        border-radius: 6px;
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
        font-size: 0.75rem;
        color: var(--muted-fg);
        font-family: "Courier New", monospace;
        display: flex;
        align-items: center;
        gap: 6px;
        padding: 8px 12px;
        background: var(--muted);
        border-radius: 6px;
        margin-top: 8px;
    }

    /* List view specific styles */
    .col-title {
        display: flex;
        flex-direction: column;
        gap: 8px;
        align-items: flex-start;
        min-width: 0;
    }

    .col-title .mod-title {
        font-size: 1rem;
        margin-bottom: 4px;
    }

    .col-title > div {
        display: flex;
        gap: 6px;
        flex-wrap: wrap;
        align-items: center;
    }

    .col-desc {
        min-width: 0;
    }

    .col-desc .mod-description {
        margin: 0;
    }

    .col-actions {
        display: flex;
        gap: 8px;
        justify-content: flex-end;
    }

    /* External link buttons */
    .ext-btn-text {
        color: var(--primary);
        background: var(--background);
        border: 1px solid var(--card-border);
        font-size: 0.8rem;
        font-weight: 600;
        padding: 6px 12px;
        border-radius: 6px;
        display: flex;
        align-items: center;
        gap: 6px;
        transition: all 0.2s ease;
        white-space: nowrap;
        cursor: pointer;
        font-family: inherit;
        text-decoration: none;
    }

    .ext-btn-text:hover {
        background: var(--muted);
        border-color: var(--primary);
    }

    .ext-btn-text.small {
        font-size: 0.75rem;
        padding: 4px 8px;
    }

    /* Footer */
    .footer-count {
        text-align: center;
        color: var(--muted-fg);
        font-size: 0.9rem;
        margin: 20px 0;
    }

    footer {
        background: var(--card);
        border-top: 1px solid var(--card-border);
        padding: 20px;
        margin-top: auto;
    }

    .footer-content {
        max-width: 1200px;
        margin: 0 auto;
        display: flex;
        justify-content: space-between;
        align-items: center;
    }

    .copyright {
        color: var(--muted-fg);
        font-size: 0.85rem;
    }

    .social-links {
        display: flex;
        gap: 12px;
    }

    .social-btn {
        color: var(--muted-fg);
        background: none;
        border: none;
        font-size: 1.2rem;
        padding: 8px;
        border-radius: 6px;
        transition: all 0.2s ease;
        width: 36px;
        height: 36px;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        text-decoration: none;
    }

    .social-btn.github {
        color: var(--foreground);
    }

    .social-btn.kofi {
        color: #ff5722;
    }

    .social-btn:hover {
        background: var(--muted);
        transform: translateY(-2px);
    }

    /* Empty state */
    .empty-state {
        text-align: center;
        padding: 80px 20px;
        color: var(--muted-fg);
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 24px;
        animation: fadeIn 0.5s ease-out;
    }

    .icon-container {
        width: 120px;
        height: 120px;
        border-radius: 50%;
        background: var(--muted);
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 4rem;
        color: var(--muted-fg);
        transition: all 0.3s ease;
        opacity: 0.6;
    }

    .empty-state:hover .icon-container {
        background: var(--card-border);
        transform: scale(1.05) rotate(-5deg);
        opacity: 0.8;
        color: var(--primary);
    }

    @keyframes fadeIn {
        from {
            opacity: 0;
            transform: translateY(20px);
        }
        to {
            opacity: 1;
            transform: translateY(0);
        }
    }

    .empty-state h2 {
        margin: 0;
        font-size: 1.5rem;
        color: var(--foreground);
    }

    .empty-state p {
        margin: 0;
        font-size: 1rem;
        max-width: 400px;
    }
</style>
