import { writable, derived } from "svelte/store";

// --- CORE APP STATE ---
export const darkMode = writable(true);
export const folderPath = writable("");
export const rawMods = writable([]);
export const status = writable("");
export const targetLoader = writable("Fabric");
export const viewMode = writable("grid");

// --- UPDATE STATE ---
export const updateAvailable = writable(false);
export const latestVersion = writable("");
export const updateUrl = writable("");
export const isCheckingUpdate = writable(false);
export const isDownloading = writable(false);
export const downloadProgress = writable(0);
export const updateError = writable("");
export const updateStage = writable("");
export const updateMessage = writable("");
export const currentVersion = writable("");

// --- FILTER & SEARCH ---
export const searchQuery = writable("");
export const sortBy = writable("name");
export const filterBy = writable("all");

// --- COLLAPSIBLE SECTIONS ---
export const isAlertsOpen = writable(true);
export const isStatsOpen = writable(true);

// --- ALERT VISIBILITY STATE ---
export const showConflictAlert = writable(true);
export const showDependencyAlert = writable(true);
export const showSodiumAlert = writable(true);

// --- EXPANDED CARDS STATE ---
export const expandedMods = writable(new Set());

// --- DEPENDENCY CHECKER ---
export const issues = writable([]);
export const isChecking = writable(false);
export const checkStatus = writable("");

// --- IGNORED DEPENDENCIES ---
export const IGNORED_DEPS = new Set([
  "fabric",
  "fabricloader",
  "minecraft",
  "java",
]);

// --- DERIVED REACTIVE STATS ---
export const fabricCount = derived(
  rawMods,
  ($rawMods) => $rawMods.filter((m) => m.loader === "Fabric").length,
);

export const forgeCount = derived(
  rawMods,
  ($rawMods) => $rawMods.filter((m) => m.loader === "NeoForge").length,
);

export const conflictCount = derived(
  [targetLoader, fabricCount, forgeCount],
  ([$targetLoader, $fabricCount, $forgeCount]) =>
    $targetLoader === "Fabric" ? $forgeCount : $fabricCount,
);

export const hasSodium = derived(rawMods, ($rawMods) =>
  $rawMods.some((m) => m.name.toLowerCase().includes("sodium")),
);

// Calculate total active alerts
export const activeAlertCount = derived(
  [
    conflictCount,
    showConflictAlert,
    issues,
    showDependencyAlert,
    hasSodium,
    showSodiumAlert,
  ],
  ([
    $conflictCount,
    $showConflictAlert,
    $issues,
    $showDependencyAlert,
    $hasSodium,
    $showSodiumAlert,
  ]) =>
    ($conflictCount > 0 && $showConflictAlert ? 1 : 0) +
    ($issues.length > 0 && $showDependencyAlert ? 1 : 0) +
    ($hasSodium && $showSodiumAlert ? 1 : 0),
);

// --- FILTER LOGIC ---
export const filteredMods = derived(
  [rawMods, searchQuery, filterBy, sortBy, targetLoader],
  ([$rawMods, $searchQuery, $filterBy, $sortBy, $targetLoader]) => {
    return $rawMods
      .filter((mod) => {
        const searchMatch =
          mod.name.toLowerCase().includes($searchQuery.toLowerCase()) ||
          mod.fileName.toLowerCase().includes($searchQuery.toLowerCase());

        let statusMatch = true;
        if ($filterBy === "compatible") {
          statusMatch = mod.loader === $targetLoader;
        } else if ($filterBy === "incompatible") {
          statusMatch = mod.loader !== $targetLoader;
        }

        return searchMatch && statusMatch;
      })
      .sort((a, b) => {
        if ($sortBy === "name") {
          return a.name.localeCompare(b.name);
        } else if ($sortBy === "loader") {
          return a.loader.localeCompare(b.loader);
        } else if ($sortBy === "status") {
          const aStatus = a.loader === $targetLoader ? 0 : 1;
          const bStatus = b.loader === $targetLoader ? 0 : 1;
          return aStatus - bStatus;
        }
        return 0;
      });
  },
);

// --- UTILITY FUNCTIONS ---
export function resetApp() {
  folderPath.set("");
  rawMods.set([]);
  status.set("");
  issues.set([]);
  checkStatus.set("");
  expandedMods.set(new Set());

  // Reset alerts to default so they show up again next time
  showConflictAlert.set(true);
  showDependencyAlert.set(true);
  showSodiumAlert.set(true);
  isAlertsOpen.set(true);
}

export function toggleTheme() {
  darkMode.update((dm) => {
    const newDarkMode = !dm;
    if (typeof document !== "undefined") {
      document.documentElement.setAttribute(
        "data-theme",
        newDarkMode ? "dark" : "light",
      );
      if (newDarkMode) {
        document.body.classList.add("dark");
      } else {
        document.body.classList.remove("dark");
      }
    }
    return newDarkMode;
  });
}

export function toggleExpand(fileName) {
  expandedMods.update((set) => {
    const newSet = new Set(set);
    if (newSet.has(fileName)) {
      newSet.delete(fileName);
    } else {
      newSet.add(fileName);
    }
    return newSet;
  });
}
