<script>
    import {
        searchQuery,
        sortBy,
        filterBy,
        viewMode,
    } from "../../stores/index.js";
</script>

<div class="toolbar">
    <div class="search-box">
        <i class="fa-solid fa-magnifying-glass"></i>
        <input
            type="text"
            placeholder="Search mods..."
            bind:value={$searchQuery}
        />
    </div>

    <div class="filter-group">
        <div class="select-wrapper">
            <i class="fa-solid fa-filter"></i>
            <select bind:value={$filterBy}>
                <option value="all">All Mods</option>
                <option value="compatible">Compatible Only</option>
                <option value="incompatible">Incompatible Only</option>
            </select>
        </div>

        <div class="select-wrapper">
            <i class="fa-solid fa-arrow-down-a-z"></i>
            <select bind:value={$sortBy}>
                <option value="name">Sort by Name</option>
                <option value="loader">Sort by Loader</option>
                <option value="status">Sort by Status</option>
            </select>
        </div>

        <button
            class="view-toggle-btn"
            on:click={() =>
                viewMode.update((mode) => (mode === "grid" ? "list" : "grid"))}
            title="Toggle View Mode"
        >
            <i class="fa-solid {$viewMode === 'grid' ? 'fa-list' : 'fa-grip'}"
            ></i>
        </button>
    </div>
</div>

<style>
    .toolbar {
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 20px;
        margin-bottom: 24px;
        padding: 16px;
        background: var(--card);
        border: 1px solid var(--card-border);
        border-radius: var(--radius);
    }

    .search-box {
        position: relative;
        flex: 1;
        max-width: 350px;
        min-width: 250px;
    }

    .search-box i {
        position: absolute;
        left: 12px;
        top: 50%;
        transform: translateY(-50%);
        color: var(--muted-fg);
        font-size: 0.9rem;
        z-index: 1;
    }

    .search-box input {
        width: 100%;
        padding: 10px 12px 10px 36px;
        border: 1px solid var(--card-border);
        border-radius: 8px;
        background: var(--card);
        color: var(--foreground);
        font-size: 0.9rem;
    }

    .search-box input:focus {
        outline: none;
        border-color: var(--primary);
    }

    .filter-group {
        display: flex;
        gap: 16px;
        align-items: center;
        flex-shrink: 0;
    }

    .select-wrapper {
        position: relative;
        display: flex;
        align-items: center;
    }

    .select-wrapper i {
        position: absolute;
        left: 12px;
        color: var(--muted-fg);
        font-size: 0.85rem;
        z-index: 1;
        pointer-events: none;
    }

    select {
        padding: 8px 12px 8px 36px;
        border: 1px solid var(--card-border);
        border-radius: 6px;
        background: var(--card);
        color: var(--foreground);
        font-size: 0.85rem;
        cursor: pointer;
        appearance: none;
        min-width: 160px;
    }

    select:hover {
        border-color: var(--primary);
    }

    .view-toggle-btn {
        background: var(--card);
        border: 1px solid var(--card-border);
        color: var(--muted-fg);
        width: 40px;
        height: 40px;
        border-radius: 8px;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 1rem;
        transition: all 0.2s ease;
    }

    .view-toggle-btn:hover {
        border-color: var(--primary);
        color: var(--primary);
    }

    @media (max-width: 768px) {
        .toolbar {
            flex-direction: column;
            gap: 16px;
            align-items: stretch;
            padding: 12px;
        }

        .search-box {
            max-width: none;
            min-width: auto;
        }

        .filter-group {
            justify-content: space-between;
            flex-wrap: wrap;
            gap: 12px;
        }

        select {
            min-width: 120px;
            flex: 1;
        }

        .view-toggle-btn {
            width: 36px;
            height: 36px;
        }
    }
</style>
