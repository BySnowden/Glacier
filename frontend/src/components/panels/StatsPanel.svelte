<script>
    import {
        rawMods,
        fabricCount,
        forgeCount,
        conflictCount,
        isStatsOpen,
        isChecking,
        checkStatus,
    } from "../../stores/index.js";

    // Removed onRunDependencyCheck prop as button is now in ModControls
</script>

{#if $rawMods.length > 0}
    <div class="section-container">
        <button
            class="section-header-btn"
            on:click={() => isStatsOpen.update((open) => !open)}
            title="Toggle Statistics Panel"
        >
            <div class="title-group">
                <i class="fa-solid fa-chart-simple text-blue"></i>
                <h3>Statistics</h3>
            </div>
            <i
                class="fa-solid {$isStatsOpen
                    ? 'fa-chevron-up'
                    : 'fa-chevron-down'}"
            ></i>
        </button>

        {#if $isStatsOpen}
            <div class="stats-bar">
                <div class="stat-group">
                    <div class="stat-box">
                        <i class="fa-solid fa-cube"></i>
                        <div class="stat-info">
                            <span class="val">{$rawMods.length}</span>
                            <span class="label">Total Mods</span>
                        </div>
                    </div>
                    <div class="stat-box success">
                        <i class="fa-solid fa-check"></i>
                        <div class="stat-info">
                            <span class="val">
                                {$rawMods.length - $conflictCount}
                            </span>
                            <span class="label">Compatible</span>
                        </div>
                    </div>
                    {#if $conflictCount > 0}
                        <div class="stat-box danger">
                            <i class="fa-solid fa-xmark"></i>
                            <div class="stat-info">
                                <span class="val">{$conflictCount}</span>
                                <span class="label">Conflicts</span>
                            </div>
                        </div>
                    {/if}
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

    .text-blue {
        color: var(--primary);
    }

    .stats-bar {
        padding: 0 24px 24px;
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 24px;
    }

    .stat-group {
        display: flex;
        gap: 32px;
        flex-wrap: wrap;
        justify-content: center;
        width: 100%;
    }

    .stat-box {
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 20px 28px;
        background: var(--background);
        border: 1px solid var(--card-border);
        border-radius: 12px;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.08);
        min-width: 160px;
        flex: 1;
        max-width: 200px;
        transition: all 0.2s ease;
    }

    .stat-box:hover {
        transform: translateY(-2px);
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        border-color: var(--primary);
    }

    .stat-box i {
        font-size: 1.4rem;
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
        gap: 2px;
    }

    .stat-info .val {
        font-size: 1.3rem;
        font-weight: 700;
        color: var(--foreground);
    }

    .stat-info .label {
        font-size: 0.75rem;
        color: var(--muted-fg);
        font-weight: 500;
        text-transform: uppercase;
        letter-spacing: 0.5px;
    }

    @media (max-width: 768px) {
        .stats-bar {
            padding: 0 16px 20px;
            gap: 16px;
        }

        .stat-group {
            gap: 16px;
            flex-direction: column;
        }

        .stat-box {
            max-width: none;
            padding: 16px 20px;
        }

        .stat-box i {
            font-size: 1.2rem;
        }

        .stat-info .val {
            font-size: 1.1rem;
        }
    }
</style>
