# ❄️ Glacier

Hey there! I'm Snowden, and this is **Glacier**. I built this because I wanted a lightweight, fast, and easy-to-understand way to manage Minecraft mods without the bloat and heft of traditional launchers.

## 🛠️ Tech Stack

I wanted this to be fast. My original concept of this program in pure Java wasn't going to cut it, so I used a mix of backend power and frontend flexibility:

| Category       | Technology        |
| :------------- | :---------------- |
| **Backend**    | Go (Golang)       |
| **Frontend**   | Svelte + HTML/CSS |
| **Bridge**     | Wails v2          |
| **Build Tool** | Vite              |

## ✨ Features

- **ARM64 Native:** One of the few mod tools built specifically to run natively on ARM64 Windows devices (like the Surface & Surface Pro which I developed this on).
- **Conflict Detection:** Instantly see if you've accidentally mixed Fabric and NeoForge mods.
- **Lightning Fast:** The Go backend scans JAR files in milliseconds.
- **Portable:** No installer needed. Just one small `.exe` and you're good to go!

## 💡 The Origin Story

The idea for Glacier sparked last year, when my girlfriend and I were making modpacks together. We would occasionally get a little ahead of ourselves, and after installing quite a few mods, there would be some incompatibility error that would take too long to understand and fix. I wanted to build a "bridge" that could scan a folder and tell you exactly what was wrong before you even hit 'Play'.

Originally, I called it _CompatBridge_, but as the UI got cleaner and cooler, **Glacier** just felt right.

## 🧠 What I Learned

Building this was a massive learning curve. I had to figure out:

- How to embed frontend assets directly into a Go binary.
- How to handle cross-platform builds for ARM64 vs x64.
- The importance of a clean UI-user experience matters as much as the code.

## 🗺️ Planned Features

- [ ] **Automatic Dependency Checking**: Finding the missing API jars for you.
- [ ] **Modrinth Integration**: Exploring ways to browse/download mods directly while respecting fellow creators' rights.
- [ ] **Theme Support**: Letting other users select their own theme.

## 🚀 How to Download

1. Head over to the **Releases** tab on the right.
2. Download `Glacier.exe`.
3. Open it, select your desired Minecraft `mods` folder, and let it do the work!

## 🛠️ Local Development

If you want to modify the code or run it from source:

### 1. Clone the Repository

```bash
git clone https://github.com/BySnowden/Glacier.git
cd Glacier
```

### 2. Start Dev Mode

Wails has a built-in hot-reload system. This command will launch the app and automatically update whenever you save changes to your Go or Svelte files:

```bash
wails dev
```

### 3. Build for Production

Create the final optimized `.exe` in `build/bin/`:

```bash
wails build
```

## 🛡️ Transparency & Security

As an independent developer, I don't have a $400/year "Publisher Certificate" to sign my apps. Because of this, Windows SmartScreen will likely show a warning when you first run Glacier.

**Is it safe?**
Absolutely. Glacier is fully open-source under the **GPLv3 license**. This means:

- Every single line of code is available right here for you to audit.
- There are no hidden backdoors or telemetry; the app only interacts with your Minecraft folder.
- You can even build the executable yourself from the source if you want to be 100% sure of what you're running.

To bypass the warning: Click **"More Info"** -> **"Run Anyway"**.

---

**Licensed under GPLv3.**
