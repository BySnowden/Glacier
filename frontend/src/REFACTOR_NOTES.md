# Component Refactoring Notes

## Overview
The original `App.svelte` file was becoming too large and difficult to maintain (1,864 lines). This refactoring breaks it down into smaller, focused, and reusable components while maintaining the same functionality.

## New Component Structure

### 1. `stores.js` - Centralized State Management
- **Purpose**: Manages all application state using Svelte stores
- **Contains**:
  - Core app state (darkMode, folderPath, rawMods, etc.)
  - Update management state
  - Filter and search state
  - Alert visibility state
  - Derived reactive computations (filteredMods, conflictCount, etc.)
  - Utility functions (resetApp, toggleTheme, toggleExpand)

### 2. `UpdateManager.svelte` - Update System
- **Purpose**: Handles all update checking and downloading functionality
- **Features**:
  - Automatic update checking on app startup
  - Progress tracking with visual progress bar
  - Error handling and retry logic
  - Clean UI for update status and actions

### 3. `ModControls.svelte` - Mod Management Controls
- **Purpose**: Handles folder selection, scanning, and dependency checking
- **Features**:
  - Target loader selection (Fabric/NeoForge)
  - Folder path selection and display
  - Mod scanning functionality
  - Dependency validation with proper error handling

### 4. `AlertPanel.svelte` - Alert System
- **Purpose**: Displays warnings and conflicts to the user
- **Features**:
  - Collapsible alert panel
  - Loader conflict warnings
  - Missing dependency alerts
  - Sodium compatibility warnings
  - Dismissible alerts with state management

### 5. `StatsPanel.svelte` - Statistics Display
- **Purpose**: Shows mod statistics and quick actions
- **Features**:
  - Total mod count, compatible count, conflict count
  - Visual stat boxes with icons
  - Dependency check trigger button
  - Collapsible panel

### 6. `ModToolbar.svelte` - Search and Filter Controls
- **Purpose**: Provides search, filtering, and view controls
- **Features**:
  - Real-time search functionality
  - Filter by compatibility status
  - Sort by name, loader, or status
  - Grid/List view toggle

## Key Benefits

1. **Maintainability**: Each component has a single, clear responsibility
2. **Reusability**: Components can be easily reused or swapped
3. **Readability**: Much easier to find and modify specific functionality
4. **Testing**: Individual components can be tested in isolation
5. **State Management**: Centralized store makes state predictable and consistent

## How to Use

The main `App.svelte` file now simply:
1. Imports all the component modules
2. Imports the stores for state management
3. Coordinates the components together
4. Handles the main mod grid display

### Example Usage:
```svelte
<script>
    import UpdateManager from "./UpdateManager.svelte";
    import ModControls from "./ModControls.svelte";
    // ... other imports
    
    import { darkMode, rawMods } from "./stores.js";
</script>

<UpdateManager />
<ModControls />
<!-- Components work together seamlessly -->
```

## File Sizes After Refactoring
- `App.svelte`: ~770 lines (was 1,864)
- `stores.js`: ~143 lines
- `UpdateManager.svelte`: ~275 lines
- `ModControls.svelte`: ~288 lines
- `AlertPanel.svelte`: ~259 lines
- `StatsPanel.svelte`: ~208 lines
- `ModToolbar.svelte`: ~167 lines

## Migration Notes

- All functionality remains exactly the same from the user perspective
- No breaking changes to the API calls or data structures
- All styling has been preserved and distributed appropriately
- The component props and events are minimal and focused
- Store subscriptions handle reactive updates automatically

## Future Improvements

With this modular structure, future enhancements become much easier:
- Add new alert types by extending `AlertPanel.svelte`
- Modify search/filter logic in `ModToolbar.svelte`
- Enhance statistics in `StatsPanel.svelte`
- Add new mod management features to `ModControls.svelte`
- Each component can be developed and tested independently