# Glacier Mod Manager - Frontend Structure

## 📁 Project Organization

This document explains the organized file structure of the Glacier frontend codebase.

### 🗂️ Directory Structure

```
src/
├── components/           # Reusable UI components
│   ├── panels/          # Main application panels
│   │   ├── AlertPanel.svelte      # Conflict & dependency warnings
│   │   ├── ModControls.svelte     # Control panel for scanning/config
│   │   └── StatsPanel.svelte      # Mod statistics display
│   └── ui/              # UI components & utilities
│       ├── Header.svelte          # Basic header component (unused)
│       ├── ModCard.svelte         # Individual mod card (unused)
│       ├── ModToolbar.svelte      # Search, filters, view controls
│       └── UpdateManager.svelte   # App update system
├── stores/              # State management
│   └── index.js                   # Centralized app state & reactive logic
├── assets/              # Static assets
│   └── images/
│       └── icon-transparent.png  # Glacier logo
├── App.svelte           # Main application component
├── main.js             # Application entry point
├── style.css           # Global styles & theme system
└── vite-env.d.ts       # TypeScript definitions
```

## 🧩 Component Breakdown

### **Main Application**
- **`App.svelte`** - Root component that orchestrates all panels and handles the mod grid display

### **Panels** (`/components/panels/`)
These are the main functional areas of the application:

#### **`ModControls.svelte`**
- **Purpose**: Primary control panel for mod management
- **Features**:
  - Target loader selection (Fabric/NeoForge)
  - Folder browser and path display
  - Scan mods functionality
  - Dependency checking
  - Status display
- **State**: Uses stores for folder path, mods data, and checking status

#### **`AlertPanel.svelte`**
- **Purpose**: Displays warnings and conflicts to users
- **Features**:
  - Collapsible alert system
  - Loader conflict warnings (Fabric vs NeoForge mods)
  - Missing dependency alerts with detailed breakdown
  - Sodium compatibility warnings
  - Dismissible alerts with persistent state
- **State**: Monitors conflict counts, dependency issues, and alert visibility

#### **`StatsPanel.svelte`**
- **Purpose**: Shows mod statistics and overview
- **Features**:
  - Total mod count display
  - Compatible vs incompatible mod breakdown
  - Conflict count with visual indicators
  - Collapsible panel design
- **State**: Derives statistics from mod data and target loader

### **UI Components** (`/components/ui/`)
These are smaller, focused UI components:

#### **`UpdateManager.svelte`**
- **Purpose**: Handles application updates
- **Features**:
  - Automatic update checking on startup
  - Download progress tracking with visual progress bar
  - Update button with various states (checking, available, error)
  - Background update installation
- **State**: Manages update availability, progress, and error states

#### **`ModToolbar.svelte`**
- **Purpose**: Provides search and filtering controls
- **Features**:
  - Real-time search functionality
  - Mod filtering (all, compatible, incompatible)
  - Sorting options (name, loader, status)
  - Grid/List view toggle
  - Responsive design for mobile
- **State**: Controls search query, filters, and view mode

#### **Legacy Components** (Currently Unused)
- **`Header.svelte`** - Basic header component (replaced by inline header in App.svelte)
- **`ModCard.svelte`** - Individual mod card component (replaced by inline cards in App.svelte)

## 🗃️ State Management (`/stores/`)

### **`index.js`**
Centralized Svelte store that manages all application state:

#### **Core State**
- `darkMode` - Theme preference
- `folderPath` - Selected mods directory
- `rawMods` - Array of scanned mod data
- `targetLoader` - Selected mod loader (Fabric/NeoForge)
- `viewMode` - Display mode (grid/list)

#### **Update System State**
- `updateAvailable` - Whether an update is available
- `isDownloading` - Update download progress
- `currentVersion` - Installed app version
- `latestVersion` - Available update version

#### **Search & Filter State**
- `searchQuery` - Current search terms
- `sortBy` - Sort criteria
- `filterBy` - Filter criteria

#### **Alert System State**
- `issues` - Dependency validation results
- `conflictCount` - Number of loader conflicts
- `showConflictAlert` - Alert visibility flags
- `isAlertsOpen` - Panel collapse state

#### **Derived State**
- `filteredMods` - Processed mod list based on search/filter
- `fabricCount` / `forgeCount` - Loader-specific counts
- `activeAlertCount` - Total active alerts

#### **Utility Functions**
- `resetApp()` - Clears all mod data and resets state
- `toggleTheme()` - Switches between light/dark themes
- `toggleExpand()` - Manages mod card expansion state

## 🎨 Styling System

### **`style.css`**
- Global CSS variables for theming
- Light/dark theme definitions using `[data-theme]` attributes
- Base component styles (buttons, forms, cards)
- Responsive design utilities

### **Theme System**
- Uses CSS custom properties for consistent theming
- Supports both light and dark modes
- Automatic theme switching with smooth transitions
- Theme preference stored in component state

## 🔄 Data Flow

1. **Initialization**: `App.svelte` mounts and initializes theme
2. **User Action**: User selects folder via `ModControls`
3. **Scanning**: `ModControls` triggers mod scanning
4. **State Update**: `stores/index.js` updates `rawMods` array
5. **Reactive Updates**: All components automatically re-render based on new state
6. **Filtering**: `ModToolbar` controls filter mod display
7. **Display**: `App.svelte` renders filtered mods in grid/list format

## 🚀 Key Benefits of This Structure

- **Modularity**: Each component has a single, clear responsibility
- **Maintainability**: Easy to find and modify specific functionality
- **Reusability**: Components can be easily moved or reused
- **Scalability**: New features can be added without affecting existing code
- **Testing**: Individual components can be tested in isolation
- **State Management**: Centralized store prevents prop drilling and ensures consistency

## 📝 Development Notes

- All components use the centralized store system for state management
- Import paths use relative imports (`../../stores/index.js`)
- Components are designed to be self-contained with their own styling
- The main `App.svelte` focuses on layout and mod display logic
- Theme switching is handled globally through the store system