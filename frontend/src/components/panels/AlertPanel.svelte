<script>
    import {
        conflictCount,
        issues,
        checkStatus,
        hasSodium,
        showConflictAlert,
        showDependencyAlert,
        showSodiumAlert,
        isAlertsOpen,
        activeAlertCount,
    } from "../../stores/index.js";
    import * as runtime from "../../../wailsjs/runtime/runtime";

    let showLinkDialog = false;
    let selectedDependency = "";

    function openLinkDialog(depName) {
        selectedDependency = depName;
        showLinkDialog = true;
    }

    function closeLinkDialog() {
        showLinkDialog = false;
        selectedDependency = "";
    }

    function openModrinth() {
        const searchUrl = `https://modrinth.com/mods?q=${encodeURIComponent(selectedDependency)}`;
        runtime.BrowserOpenURL(searchUrl);
        closeLinkDialog();
    }

    function openCurseForge() {
        const searchUrl = `https://www.curseforge.com/minecraft/mc-mods/search?search=${encodeURIComponent(selectedDependency)}`;
        runtime.BrowserOpenURL(searchUrl);
        closeLinkDialog();
    }
</script>

{#if ($conflictCount > 0 && $showConflictAlert) || ($issues.length > 0 && $checkStatus === "issues" && $showDependencyAlert) || ($hasSodium && $showSodiumAlert)}
    <div class="section-container">
        <button
            class="section-header-btn"
            on:click={() => isAlertsOpen.update((open) => !open)}
            title="Toggle Alerts Panel"
        >
            <div class="title-group">
                <i class="fa-solid fa-triangle-exclamation text-warning"></i>
                <h3>Alerts ({$activeAlertCount})</h3>
            </div>
            <i
                class="fa-solid {$isAlertsOpen
                    ? 'fa-chevron-up'
                    : 'fa-chevron-down'}"
            ></i>
        </button>

        {#if $isAlertsOpen}
            <div class="alerts-body">
                {#if $conflictCount > 0 && $showConflictAlert}
                    <div class="alert critical">
                        <i class="fa-solid fa-circle-xmark icon"></i>
                        <div class="content">
                            <h4>Loader Conflicts Detected</h4>
                            <p>
                                You have {$conflictCount} mods that don't match your
                                target loader. These mods may not work correctly or
                                cause crashes.
                            </p>
                        </div>
                        <button
                            class="close-alert"
                            on:click={() => showConflictAlert.set(false)}
                            title="Dismiss Alert"
                        >
                            <i class="fa-solid fa-xmark"></i>
                        </button>
                    </div>
                {/if}

                {#if $checkStatus === "issues" && $showDependencyAlert}
                    <div class="alert danger-list">
                        <i class="fa-solid fa-puzzle-piece icon"></i>
                        <div class="content">
                            <h4>Missing Dependencies ({$issues.length})</h4>
                            <div class="issue-grid">
                                {#each $issues as issue}
                                    <div class="issue-item">
                                        <span
                                            class="mod-name"
                                            title={issue.ModName}
                                        >
                                            {issue.ModName}
                                        </span>
                                        <span class="arrow" title="requires">
                                            <i
                                                class="fa-solid fa-arrow-right"
                                                style="font-size: 0.7rem;"
                                            ></i>
                                        </span>
                                        <button
                                            class="missing-dep clickable-dep"
                                            title="Click to find {issue.MissingDep} on mod platforms"
                                            on:click={() =>
                                                openLinkDialog(
                                                    issue.MissingDep,
                                                )}
                                        >
                                            {issue.MissingDep}
                                            <i
                                                class="fa-solid fa-arrow-up-right-from-square"
                                            ></i>
                                        </button>
                                    </div>
                                {/each}
                            </div>
                        </div>
                        <button
                            class="close-alert"
                            on:click={() => showDependencyAlert.set(false)}
                            title="Dismiss Alert"
                        >
                            <i class="fa-solid fa-xmark"></i>
                        </button>
                    </div>
                {/if}

                {#if $hasSodium && $showSodiumAlert}
                    <div class="alert warning">
                        <i class="fa-solid fa-triangle-exclamation icon"></i>
                        <div class="content">
                            <h4>Sodium Detected</h4>
                            <p>
                                Sodium may cause compatibility issues with some
                                mods. Consider alternatives if you encounter
                                problems.
                            </p>
                        </div>
                        <button
                            class="close-alert"
                            on:click={() => showSodiumAlert.set(false)}
                            title="Dismiss Alert"
                        >
                            <i class="fa-solid fa-xmark"></i>
                        </button>
                    </div>
                {/if}
            </div>
        {/if}

        <!-- Link Choice Dialog -->
        {#if showLinkDialog}
            <div class="modal-overlay" on:click={closeLinkDialog}>
                <div class="modal-dialog" on:click|stopPropagation>
                    <div class="modal-header">
                        <h3>Choose Platform</h3>
                        <button class="close-modal" on:click={closeLinkDialog}>
                            <i class="fa-solid fa-xmark"></i>
                        </button>
                    </div>
                    <div class="modal-body">
                        <p>
                            Where would you like to search for <strong
                                >{selectedDependency}</strong
                            >?
                        </p>
                        <div class="platform-buttons">
                            <button
                                class="platform-btn modrinth-btn"
                                on:click={openModrinth}
                            >
                                <div class="platform-info">
                                    <strong>Modrinth</strong>
                                    <small>Modern mod platform</small>
                                </div>
                            </button>
                            <button
                                class="platform-btn curseforge-btn"
                                on:click={openCurseForge}
                            >
                                <div class="platform-info">
                                    <strong>CurseForge</strong>
                                    <small>Traditional mod platform</small>
                                </div>
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        {/if}
    </div>
{/if}

<style>
    .section-container {
        background: var(--card);
        border: 1px solid var(--card-border);
        border-radius: var(--radius);
        overflow: hidden;
        margin-bottom: 24px;
        box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
    }

    .section-header-btn {
        width: 100%;
        background: none;
        border: none;
        padding: 18px 24px;
        display: flex;
        justify-content: space-between;
        align-items: center;
        cursor: pointer;
        transition: all 0.2s ease;
        color: var(--foreground);
    }

    .section-header-btn:hover {
        background: rgba(0, 0, 0, 0.02);
    }

    [data-theme="dark"] .section-header-btn:hover {
        background: rgba(255, 255, 255, 0.05);
    }

    .title-group {
        display: flex;
        align-items: center;
        gap: 12px;
    }

    .title-group h3 {
        margin: 0;
        font-size: 1.1rem;
        font-weight: 600;
    }

    .text-warning {
        color: #f59e0b;
    }

    .alerts-body {
        padding: 0 24px 24px;
        display: flex;
        flex-direction: column;
        gap: 18px;
    }

    .alert {
        background: var(--background);
        border: 1px solid var(--card-border);
        border-radius: 8px;
        padding: 18px;
        display: flex;
        gap: 16px;
        align-items: flex-start;
        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    }

    .alert:last-child {
        margin-bottom: 0;
    }

    .alert.critical {
        background: var(--alert-critical-bg);
        border-color: var(--alert-critical-border);
        color: var(--alert-critical-text);
    }

    .alert.warning {
        background: var(--alert-warning-bg);
        border-color: var(--alert-warning-border);
        color: var(--alert-warning-text);
    }

    .alert.danger-list {
        background: var(--alert-critical-bg);
        border-color: var(--alert-critical-border);
        color: var(--alert-critical-text);
    }

    .alert .icon {
        flex-shrink: 0;
        font-size: 1.2rem;
        margin-top: 2px;
    }

    .alert h4 {
        margin: 0 0 8px 0;
        font-size: 1rem;
        font-weight: 600;
    }

    .alert p {
        margin: 0;
        font-size: 0.9rem;
        line-height: 1.5;
        opacity: 0.9;
    }

    .close-alert {
        background: none;
        border: none;
        color: inherit;
        cursor: pointer;
        padding: 4px 6px;
        border-radius: 4px;
        opacity: 0.7;
        transition: all 0.2s ease;
        flex-shrink: 0;
    }

    .close-alert:hover {
        opacity: 1;
        background: rgba(0, 0, 0, 0.1);
    }

    .content {
        flex: 1;
    }

    .issue-grid {
        display: grid;
        gap: 10px;
        margin-top: 14px;
    }

    .issue-item {
        display: grid;
        grid-template-columns: 1fr auto 1fr;
        gap: 14px;
        align-items: center;
        padding: 10px 14px;
        background: rgba(0, 0, 0, 0.1);
        border-radius: 6px;
        font-size: 0.85rem;
        border: 1px solid rgba(0, 0, 0, 0.05);
    }

    [data-theme="dark"] .issue-item {
        background: rgba(255, 255, 255, 0.05);
        border-color: rgba(255, 255, 255, 0.1);
    }

    .mod-name {
        font-weight: 600;
        text-align: left;
    }

    .missing-dep {
        font-weight: 500;
        text-align: left;
    }

    .clickable-dep {
        background: none;
        border: none;
        color: inherit;
        font-family: inherit;
        font-size: inherit;
        font-weight: 500;
        text-align: left;
        cursor: pointer;
        padding: 0;
        display: flex;
        align-items: center;
        gap: 6px;
        transition: all 0.2s ease;
        text-decoration: underline;
        text-decoration-color: transparent;
    }

    .clickable-dep:hover {
        text-decoration-color: currentColor;
        opacity: 0.8;
    }

    .clickable-dep i {
        font-size: 0.75rem;
        opacity: 0.7;
        transition: opacity 0.2s ease;
    }

    .clickable-dep:hover i {
        opacity: 1;
    }

    .issue-item .arrow {
        text-align: center;
        opacity: 0.6;
    }

    /* Modal Styles */
    .modal-overlay {
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
        backdrop-filter: blur(2px);
    }

    .modal-dialog {
        background: var(--card);
        border: 1px solid var(--card-border);
        border-radius: 12px;
        box-shadow: 0 10px 25px rgba(0, 0, 0, 0.2);
        min-width: 400px;
        max-width: 90vw;
        max-height: 90vh;
        overflow: hidden;
    }

    .modal-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 20px 24px;
        border-bottom: 1px solid var(--card-border);
        background: var(--background);
    }

    .modal-header h3 {
        margin: 0;
        font-size: 1.1rem;
        font-weight: 600;
        color: var(--foreground);
    }

    .close-modal {
        background: none;
        border: none;
        color: var(--foreground);
        cursor: pointer;
        padding: 6px;
        border-radius: 4px;
        opacity: 0.7;
        transition: all 0.2s ease;
    }

    .close-modal:hover {
        opacity: 1;
        background: var(--card-border);
    }

    .modal-body {
        padding: 24px;
    }

    .modal-body p {
        margin: 0 0 20px 0;
        color: var(--foreground);
        font-size: 0.95rem;
        line-height: 1.5;
    }

    .platform-buttons {
        display: flex;
        gap: 12px;
        flex-direction: column;
    }

    .platform-btn {
        background: var(--background);
        border: 1px solid var(--card-border);
        border-radius: 8px;
        padding: 16px 20px;
        cursor: pointer;
        transition: all 0.2s ease;
        text-align: left;
    }

    .platform-btn:hover {
        background: var(--card);
        border-color: var(--primary);
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
    }

    .modrinth-btn:hover {
        border-color: #00af5c;
    }

    .curseforge-btn:hover {
        border-color: #f16436;
    }

    .platform-info {
        display: flex;
        flex-direction: column;
        gap: 4px;
    }

    .platform-info strong {
        color: var(--foreground);
        font-size: 1rem;
        font-weight: 600;
    }

    .platform-info small {
        color: var(--foreground);
        opacity: 0.7;
        font-size: 0.85rem;
    }

    @media (min-width: 500px) {
        .platform-buttons {
            flex-direction: row;
        }
    }
</style>
