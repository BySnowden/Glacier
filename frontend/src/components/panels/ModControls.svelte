<script>
    import {
        folderPath,
        rawMods,
        status,
        targetLoader,
        issues,
        isChecking,
        checkStatus,
        IGNORED_DEPS,
        resetApp,
    } from "../../stores/index.js";
    import {
        SelectModFolder,
        ScanMods,
        ValidateMods,
    } from "../../../wailsjs/go/main/App";

    async function pickFolder() {
        let path = await SelectModFolder();
        if (path) {
            resetApp();
            folderPath.set(path);
        }
    }

    async function handleScan() {
        if (!$folderPath) {
            console.log("No folder selected");
            return;
        }

        status.set("Scanning mods...");
        try {
            let mods = await ScanMods($folderPath);
            rawMods.set(mods);
            status.set(`Found ${mods.length} mods`);
        } catch (error) {
            console.error("Failed to scan mods:", error);
            status.set("Failed to scan mods");
        }
    }

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
</script>

<div class="controls-card">
    <div class="main-controls">
        <div class="control-col">
            <label for="target-loader">
                <i class="fa-solid fa-crosshairs"></i>
                Target Loader
            </label>
            <div class="toggle-group">
                <button
                    class="toggle-btn {$targetLoader === 'Fabric'
                        ? 'active fabric'
                        : ''}"
                    on:click={() => targetLoader.set("Fabric")}
                >
                    Fabric
                </button>
                <button
                    class="toggle-btn {$targetLoader === 'NeoForge'
                        ? 'active forge'
                        : ''}"
                    on:click={() => targetLoader.set("NeoForge")}
                >
                    NeoForge
                </button>
            </div>
        </div>

        <div class="control-col wide">
            <label for="folder-path">
                <i class="fa-regular fa-folder-open"></i>
                Mods Folder
            </label>
            <div class="input-group">
                <div class="path-display">
                    {$folderPath || "No folder selected"}
                </div>
                <button class="btn btn-primary" on:click={pickFolder}>
                    Browse
                </button>
            </div>
        </div>
    </div>

    {#if $folderPath || $rawMods.length > 0}
        <div class="action-controls">
            {#if $folderPath}
                <button class="btn btn-secondary" on:click={handleScan}>
                    <i class="fa-solid fa-magnifying-glass"></i>
                    Scan Mods
                </button>
            {/if}

            {#if $rawMods.length > 0}
                <button
                    class="btn btn-secondary"
                    on:click={runDependencyCheck}
                    disabled={$isChecking}
                >
                    {#if $isChecking}
                        <i class="fa-solid fa-spinner fa-spin"></i>
                        Checking...
                    {:else}
                        <i class="fa-solid fa-shield-halved"></i>
                        Check Dependencies
                    {/if}
                </button>
            {/if}
        </div>
    {/if}

    {#if $status}
        <div class="status-text">
            <i class="fa-solid fa-info-circle"></i>
            {$status}
        </div>
    {/if}
</div>

<style>
    .controls-card {
        background: var(--card);
        border: 1px solid var(--card-border);
        border-radius: var(--radius);
        padding: 24px;
        display: flex;
        flex-direction: column;
        gap: 20px;
        margin-bottom: 20px;
    }

    .main-controls {
        display: flex;
        gap: 24px;
        align-items: flex-end;
        flex-wrap: wrap;
    }

    .action-controls {
        display: flex;
        gap: 16px;
        align-items: center;
        padding-top: 16px;
        border-top: 1px solid var(--card-border);
        flex-wrap: wrap;
    }

    .control-col {
        display: flex;
        flex-direction: column;
        gap: 8px;
        min-width: 140px;
    }

    .control-col.wide {
        flex: 1;
        min-width: 300px;
    }

    label {
        font-size: 0.9rem;
        font-weight: 600;
        color: var(--foreground);
        display: flex;
        align-items: center;
        gap: 8px;
    }

    .toggle-group {
        display: flex;
        background: var(--background);
        border-radius: 8px;
        padding: 2px;
        border: 1px solid var(--card-border);
    }

    .toggle-btn {
        flex: 1;
        padding: 8px 16px;
        background: transparent;
        border: none;
        border-radius: 6px;
        font-size: 0.85rem;
        font-weight: 500;
        color: var(--muted-fg);
        cursor: pointer;
        transition: all 0.2s ease;
    }

    .toggle-btn.active {
        background: var(--card);
        color: var(--foreground);
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.2);
    }

    .toggle-btn.active.fabric {
        color: var(--fabric-text);
    }

    .toggle-btn.active.forge {
        color: var(--forge-text);
    }

    .input-group {
        display: flex;
        gap: 12px;
    }

    .path-display {
        flex: 1;
        padding: 10px 12px;
        background: var(--background);
        border: 1px solid var(--card-border);
        border-radius: 6px;
        font-size: 0.85rem;
        color: var(--muted-fg);
        min-height: 20px;
        display: flex;
        align-items: center;
        font-family: "Courier New", monospace;
        word-break: break-all;
    }

    .btn {
        padding: 10px 16px;
        border-radius: 6px;
        font-size: 0.85rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s ease;
        display: flex;
        align-items: center;
        gap: 8px;
        border: none;
        text-decoration: none;
        white-space: nowrap;
    }

    .btn-primary {
        background: var(--primary);
        color: white;
    }

    .btn-primary:hover {
        background: #2563eb;
        transform: translateY(-1px);
    }

    .btn-secondary {
        background: var(--card);
        color: var(--foreground);
        border: 1px solid var(--card-border);
    }

    .btn-secondary:hover:not(:disabled) {
        background: var(--muted);
        border-color: var(--primary);
    }

    .btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
        transform: none;
    }

    .status-text {
        display: flex;
        align-items: center;
        gap: 8px;
        padding: 12px 16px;
        background: var(--muted);
        border: 1px solid var(--card-border);
        border-radius: 6px;
        font-size: 0.9rem;
        color: var(--foreground);
        border-top: 1px solid var(--card-border);
        margin-top: 8px;
    }

    .status-text i {
        color: var(--primary);
    }
</style>
