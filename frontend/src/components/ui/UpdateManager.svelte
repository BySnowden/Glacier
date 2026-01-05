<script>
    import { onMount } from "svelte";
    import {
        updateAvailable,
        latestVersion,
        updateUrl,
        isCheckingUpdate,
        isDownloading,
        downloadProgress,
        updateError,
        updateStage,
        updateMessage,
        currentVersion,
    } from "../../stores/index.js";
    import {
        CheckForUpdates,
        DownloadAndInstall,
        GetAppVersion,
    } from "../../../wailsjs/go/main/App";
    import * as runtime from "../../../wailsjs/runtime/runtime";

    async function checkForAppUpdates() {
        isCheckingUpdate.set(true);
        updateError.set("");

        try {
            const res = await CheckForUpdates();
            if (res.error) {
                updateError.set(res.error);
                console.error("Update check failed:", res.error);
                // Reset update state on error
                updateAvailable.set(false);
                latestVersion.set("");
                updateUrl.set("");
            } else if (res.updateAvailable) {
                updateAvailable.set(true);
                latestVersion.set(res.latestVersion);
                updateUrl.set(res.downloadUrl);
            } else {
                // No update available - reset state
                updateAvailable.set(false);
                latestVersion.set(res.latestVersion || $currentVersion);
                updateUrl.set("");
            }
        } catch (err) {
            updateError.set("Failed to check for updates");
            console.error("Failed to check updates", err);
            // Reset update state on error
            updateAvailable.set(false);
            latestVersion.set("");
            updateUrl.set("");
        }

        isCheckingUpdate.set(false);
    }

    async function handleUpdateClick() {
        if (!$updateAvailable && !$updateError) {
            await checkForAppUpdates();
            return;
        }

        if (!$updateAvailable) return;

        isDownloading.set(true);
        updateError.set("");
        downloadProgress.set(0);
        updateStage.set("preparing");
        updateMessage.set("Preparing update...");

        try {
            // Listen for progress events
            const unsubscribeProgress = runtime.EventsOn(
                "update-progress",
                (data) => {
                    downloadProgress.set(data.percentage);
                    updateStage.set(data.stage);
                    updateMessage.set(data.message);
                },
            );

            const unsubscribeError = runtime.EventsOn(
                "update-error",
                (error) => {
                    updateError.set(error);
                    isDownloading.set(false);
                    unsubscribeProgress();
                    unsubscribeError();
                },
            );

            // This calls the Go function to download and install
            await DownloadAndInstall($updateUrl);

            // If we reach here, the app should be quitting soon
            updateMessage.set("Update installed successfully. Restarting...");
        } catch (err) {
            console.error("Update failed", err);
            updateError.set("Update failed: " + err.message);
            isDownloading.set(false);
        }
    }

    // Check for updates automatically on startup
    onMount(async () => {
        // Get current version
        try {
            const version = await GetAppVersion();
            currentVersion.set(version);
        } catch (err) {
            console.error("Failed to get app version:", err);
            currentVersion.set("Unknown");
        }

        // Wait a bit before checking to let the app fully load
        setTimeout(() => {
            checkForAppUpdates();
        }, 2000);
    });
</script>

<div class="update-controls">
    {#if $isDownloading}
        <div class="update-progress-container">
            <span class="progress-text">{$updateMessage}</span>
            <div class="progress-bar-bg">
                <div
                    class="progress-bar-fill"
                    style="width: {$downloadProgress}%"
                ></div>
            </div>
            {#if $updateStage === "downloading"}
                <span class="progress-percentage">
                    {Math.round($downloadProgress)}%
                </span>
            {/if}
        </div>
    {:else}
        <button
            class="update-btn {$updateAvailable
                ? 'available'
                : ''} {$updateError ? 'error' : ''}"
            on:click={handleUpdateClick}
            disabled={$isCheckingUpdate || $isDownloading}
            title={$updateAvailable
                ? `Update to ${$latestVersion}`
                : $updateError
                  ? $updateError
                  : `Current version: ${$currentVersion}`}
        >
            {#if $isCheckingUpdate}
                <i class="fa-solid fa-circle-notch fa-spin"></i>
                <span>Checking...</span>
            {:else if $updateError}
                <i class="fa-solid fa-exclamation-triangle"></i>
                <span>Update Error</span>
            {:else if $updateAvailable}
                <i class="fa-solid fa-cloud-arrow-down"></i>
                <span>Update v{$latestVersion}</span>
            {:else}
                <i class="fa-solid fa-check"></i>
                <span>Up to date</span>
            {/if}
        </button>
    {/if}
</div>

<style>
    .update-controls {
        display: flex;
        align-items: center;
        gap: 16px;
    }

    .update-btn {
        background: var(--card);
        border: 1px solid var(--card-border);
        color: var(--foreground);
        padding: 10px 16px;
        border-radius: 8px;
        cursor: pointer;
        font-size: 0.875rem;
        font-weight: 600;
        display: flex;
        align-items: center;
        gap: 8px;
        transition: all 0.2s ease;
        min-width: 160px;
        justify-content: center;
        white-space: nowrap;
        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    }

    .update-btn:hover:not(:disabled) {
        background: var(--muted);
        border-color: var(--primary);
        transform: translateY(-1px);
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
    }

    .update-btn.available {
        background: var(--success);
        color: white;
        border-color: var(--success);
        animation: pulse 2s infinite;
        box-shadow: 0 2px 8px rgba(16, 185, 129, 0.2);
    }

    .update-btn.available:hover {
        background: #059669;
        border-color: #059669;
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(16, 185, 129, 0.3);
    }

    .update-btn.error {
        background: var(--destructive);
        color: white;
        border-color: var(--destructive);
        box-shadow: 0 2px 8px rgba(239, 68, 68, 0.2);
    }

    .update-btn.error:hover {
        background: #dc2626;
        border-color: #dc2626;
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(239, 68, 68, 0.3);
    }

    .update-btn:disabled {
        opacity: 0.6;
        cursor: not-allowed;
        background: var(--muted);
        color: var(--muted-fg);
        border-color: var(--card-border);
        transform: none;
        box-shadow: none;
        animation: none;
    }

    .update-btn:disabled:hover {
        transform: none;
        box-shadow: none;
    }

    .update-progress-container {
        display: flex;
        flex-direction: column;
        gap: 8px;
        min-width: 220px;
        padding: 8px 12px;
        background: var(--card);
        border: 1px solid var(--card-border);
        border-radius: 8px;
        box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
    }

    .progress-text {
        font-size: 0.875rem;
        color: var(--foreground);
        font-weight: 500;
        text-align: center;
    }

    .progress-percentage {
        font-size: 0.8rem;
        color: var(--muted-fg);
        text-align: center;
        font-weight: 600;
    }

    .progress-bar-bg {
        height: 8px;
        background: var(--muted);
        border-radius: 4px;
        overflow: hidden;
        border: 1px solid var(--card-border);
    }

    .progress-bar-fill {
        height: 100%;
        background: linear-gradient(90deg, var(--primary), #60a5fa);
        transition: width 0.3s ease;
        position: relative;
        border-radius: 3px;
    }

    .progress-bar-fill::after {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: linear-gradient(
            90deg,
            transparent,
            rgba(255, 255, 255, 0.3),
            transparent
        );
        animation: shimmer 1.5s infinite;
    }

    @keyframes pulse {
        0%,
        100% {
            opacity: 1;
            box-shadow: 0 2px 8px rgba(16, 185, 129, 0.2);
        }
        50% {
            opacity: 0.9;
            box-shadow: 0 4px 16px rgba(16, 185, 129, 0.4);
        }
    }

    @keyframes shimmer {
        0% {
            transform: translateX(-100%);
        }
        100% {
            transform: translateX(100%);
        }
    }
</style>
