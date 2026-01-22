#!/usr/bin/env bun

import JSZip from "jszip";

const HOST = "https://setup-aws.rbxcdn.com";

const BIN = {
  WindowsPlayer: "/",
  WindowsStudio64: "/",
  MacPlayer: "/mac/",
  MacStudio: "/mac/",
};

const EXTRACT_ROOTS = {
  player: {
    "RobloxApp.zip": "",
    "redist.zip": "",
    "shaders.zip": "shaders/",
    "ssl.zip": "ssl/",
    "WebView2.zip": "",
    "WebView2RuntimeInstaller.zip": "WebView2RuntimeInstaller/",
    "content-avatar.zip": "content/avatar/",
    "content-configs.zip": "content/configs/",
    "content-fonts.zip": "content/fonts/",
    "content-sky.zip": "content/sky/",
    "content-sounds.zip": "content/sounds/",
    "content-textures2.zip": "content/textures/",
    "content-models.zip": "content/models/",
    "content-platform-fonts.zip": "PlatformContent/pc/fonts/",
    "content-platform-dictionaries.zip":
      "PlatformContent/pc/shared_compression_dictionaries/",
    "content-terrain.zip": "PlatformContent/pc/terrain/",
    "content-textures3.zip": "PlatformContent/pc/textures/",
    "extracontent-luapackages.zip": "ExtraContent/LuaPackages/",
    "extracontent-translations.zip": "ExtraContent/translations/",
    "extracontent-models.zip": "ExtraContent/models/",
    "extracontent-textures.zip": "ExtraContent/textures/",
    "extracontent-places.zip": "ExtraContent/places/",
  },
  studio: {
    "RobloxStudio.zip": "",
    "RibbonConfig.zip": "RibbonConfig/",
    "redist.zip": "",
    "Libraries.zip": "",
    "LibrariesQt5.zip": "",
    "WebView2.zip": "",
    "WebView2RuntimeInstaller.zip": "",
    "shaders.zip": "shaders/",
    "ssl.zip": "ssl/",
    "Qml.zip": "Qml/",
    "Plugins.zip": "Plugins/",
    "StudioFonts.zip": "StudioFonts/",
    "BuiltInPlugins.zip": "BuiltInPlugins/",
    "ApplicationConfig.zip": "ApplicationConfig/",
    "BuiltInStandalonePlugins.zip": "BuiltInStandalonePlugins/",
    "content-qt_translations.zip": "content/qt_translations/",
    "content-sky.zip": "content/sky/",
    "content-fonts.zip": "content/fonts/",
    "content-avatar.zip": "content/avatar/",
    "content-models.zip": "content/models/",
    "content-sounds.zip": "content/sounds/",
    "content-configs.zip": "content/configs/",
    "content-api-docs.zip": "content/api_docs/",
    "content-textures2.zip": "content/textures/",
    "content-studio_svg_textures.zip": "content/studio_svg_textures/",
    "content-platform-fonts.zip": "PlatformContent/pc/fonts/",
    "content-platform-dictionaries.zip":
      "PlatformContent/pc/shared_compression_dictionaries/",
    "content-terrain.zip": "PlatformContent/pc/terrain/",
    "content-textures3.zip": "PlatformContent/pc/textures/",
    "extracontent-translations.zip": "ExtraContent/translations/",
    "extracontent-luapackages.zip": "ExtraContent/LuaPackages/",
    "extracontent-textures.zip": "ExtraContent/textures/",
    "extracontent-scripts.zip": "ExtraContent/scripts/",
    "extracontent-models.zip": "ExtraContent/models/",
    "studiocontent-models.zip": "StudioContent/models/",
    "studiocontent-textures.zip": "StudioContent/textures/",
  },
};

const COLORS = {
  reset: "\x1b[0m",
  bright: "\x1b[1m",
  dim: "\x1b[2m",
  cyan: "\x1b[36m",
  green: "\x1b[32m",
  yellow: "\x1b[33m",
  red: "\x1b[31m",
  blue: "\x1b[34m",
};

function help() {
  console.log(`
${COLORS.bright}${COLORS.cyan}rdd-cli${COLORS.reset}

${COLORS.bright}USAGE:${COLORS.reset}
  rdd download -v <version> -t <binary> [--channel <n>]

${COLORS.bright}OPTIONS:${COLORS.reset}
  ${COLORS.blue}-d, --download${COLORS.reset}        Download roblox binary
  ${COLORS.blue}-v, --version${COLORS.reset}         Version hash (example: version-xxxx)
  ${COLORS.blue}-t, --type${COLORS.reset}            WindowsPlayer | WindowsStudio64 | MacPlayer | MacStudio
  ${COLORS.blue}-c, --channel${COLORS.reset}         LIVE | zflag | production (default: LIVE)
  ${COLORS.blue}--compress${COLORS.reset}            Compress output zip (default: false)
  ${COLORS.blue}--level${COLORS.reset}               Compression level 1-9 (default: 6)
  ${COLORS.blue}-h, --help${COLORS.reset}            Show help
  ${COLORS.blue}--version-cli${COLORS.reset}         Show cli version

${COLORS.bright}EXAMPLE:${COLORS.reset}
  rdd download -v version-123abc -t WindowsPlayer
  rdd download -v version-123abc -t MacStudio --channel LIVE
  rdd download -v version-123abc -t WindowsPlayer --compress --level 9
`);
}

const args = process.argv.slice(2);

async function showVersions() {
  try {
    const versions = {
      future: await fetch("https://weao.xyz/api/versions/future", {
        headers: { "User-Agent": "WEAO-3PService" },
      }).then((r) => r.json()),
      current: await fetch("https://weao.xyz/api/versions/current", {
        headers: { "User-Agent": "WEAO-3PService" },
      }).then((r) => r.json()),
      past: await fetch("https://weao.xyz/api/versions/past", {
        headers: { "User-Agent": "WEAO-3PService" },
      }).then((r) => r.json()),
    };

    console.log(
      `Source: ${COLORS.reset}${COLORS.blue}https://weao.xyz/${COLORS.reset}${COLORS.dim}`,
    );

    console.log(`
${COLORS.bright}${COLORS.cyan}Future Version${COLORS.reset} ${COLORS.dim}${COLORS.reset}
  ${COLORS.blue}Windows${COLORS.reset}  ${versions.future.Windows} ${COLORS.dim}(${versions.future.WindowsDate})${COLORS.reset}
  ${COLORS.blue}Mac${COLORS.reset}      ${versions.future.Mac} ${COLORS.dim}(${versions.future.MacDate})${COLORS.reset}

${COLORS.bright}${COLORS.cyan}Current Version:${COLORS.reset}
  ${COLORS.green}Windows${COLORS.reset}  ${versions.current.Windows} ${COLORS.dim}(${versions.current.WindowsDate})${COLORS.reset}
  ${COLORS.green}Mac${COLORS.reset}      ${versions.current.Mac} ${COLORS.dim}(${versions.current.MacDate})${COLORS.reset}

${COLORS.bright}${COLORS.cyan}Past Version:${COLORS.reset}
  ${COLORS.red}Windows${COLORS.reset}  ${versions.past.Windows} ${COLORS.dim}(${versions.past.WindowsDate})${COLORS.reset}
  ${COLORS.red}Mac${COLORS.reset}      ${versions.past.Mac} ${COLORS.dim}(${versions.past.MacDate})${COLORS.reset}
`);
  } catch (err) {
    console.error(`${COLORS.red}Failed to fetch versions${COLORS.reset}`);
  }
}

if (args.length === 0 || args.includes("--help") || args.includes("-h")) {
  help();
  await showVersions();
  process.exit(0);
}

if (args.includes("--version-cli")) {
  console.log(`${COLORS.cyan}rdd-cli${COLORS.reset} v1.0.0`);
  process.exit(0);
}

function getArg(flags: string[]) {
  for (const flag of flags) {
    const idx = args.indexOf(flag);
    if (idx !== -1 && idx + 1 < args.length) {
      return args[idx + 1];
    }
  }
  return null;
}

function hasFlag(flags: string[]) {
  return flags.some((f) => args.includes(f));
}

const isDownload = hasFlag(["download", "-d"]);

if (!isDownload) {
  help();
  process.exit(0);
}

const version = getArg(["-v", "--version"]);
const type = getArg(["-t", "--type"]) as keyof typeof BIN;
const channel = getArg(["-c", "--channel"]) || "LIVE";
const compress = hasFlag(["--compress"]);
const level = parseInt(getArg(["--level"]) || "6");

if (!version || !type || !(type in BIN)) {
  console.error(
    `${COLORS.red}Error: Missing or invalid arguments${COLORS.reset}`,
  );
  help();
  process.exit(1);
}

const base =
  channel === "LIVE" ? HOST : `${HOST}/channel/${channel.toLowerCase()}`;

const root = `${base}${BIN[type]}${version}-`;

function formatBytes(n: number) {
  if (n >= 1024 * 1024 * 1024)
    return (n / 1024 / 1024 / 1024).toFixed(2) + " GB";
  if (n >= 1024 * 1024) return (n / 1024 / 1024).toFixed(2) + " MB";
  if (n >= 1024) return (n / 1024).toFixed(2) + " KB";
  return n.toFixed(2) + " B";
}

function formatSpeed(bytesPerSec: number) {
  if (bytesPerSec >= 1024 * 1024)
    return (bytesPerSec / 1024 / 1024).toFixed(2) + " MB/s";
  if (bytesPerSec >= 1024) return (bytesPerSec / 1024).toFixed(2) + " KB/s";
  return bytesPerSec.toFixed(2) + " B/s";
}

function bar(cur: number, total: number, elapsed: number) {
  const w = 28;
  const pct = total ? cur / total : 0;
  const fill = Math.floor(w * pct);
  const b = "█".repeat(fill) + "░".repeat(w - fill);

  const speed = elapsed > 0 ? cur / elapsed : 0;
  const remaining = speed > 0 ? (total - cur) / speed : 0;
  const eta =
    remaining > 0
      ? new Date(Date.now() + remaining * 1000).toLocaleTimeString()
      : "--:--:--";

  process.stdout.write(
    `\r  ${COLORS.cyan}[${b}]${COLORS.reset} ${(pct * 100).toFixed(1).padStart(5)}%  ${formatBytes(cur).padStart(10)} / ${formatBytes(total).padStart(10)}  ${COLORS.green}${formatSpeed(speed).padStart(10)}${COLORS.reset}  ${COLORS.dim}ETA ${eta}${COLORS.reset}`,
  );
}

async function stream(url: string, name: string) {
  console.log(`  ${COLORS.dim}↓${COLORS.reset} ${name}`);

  const res = await fetch(url);
  if (!res.ok) {
    console.error(
      `  ${COLORS.red}✗${COLORS.reset} Failed to fetch: ${url} (${res.status})`,
    );
    process.exit(1);
  }

  const total = Number(res.headers.get("content-length") ?? 0);
  const reader = res.body!.getReader();

  let recv = 0;
  const parts: Uint8Array[] = [];
  const startTime = Date.now();

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;

    recv += value.length;
    parts.push(value);

    const elapsed = (Date.now() - startTime) / 1000;
    bar(recv, total, elapsed);
  }

  process.stdout.write("\n");
  return Buffer.concat(parts);
}

async function downloadMac() {
  const file =
    type === "MacPlayer"
      ? "RobloxPlayer.zip"
      : type === "MacStudio"
        ? "RobloxStudioApp.zip"
        : "RobloxApp.zip";

  const data = await stream(root + file, file);
  const out = `${channel}-${type}-${version}.zip`;

  await Bun.write(out, data);
  console.log(
    `  ${COLORS.green}✓${COLORS.reset} Saved to ${COLORS.bright}${out}${COLORS.reset}`,
  );
}

async function downloadWindows() {
  let manifestUrl = root + "rbxPkgManifest.txt";
  let manifestRes = await fetch(manifestUrl);

  if (!manifestRes.ok) {
    const fallbackRoot = `${HOST}/channel/common${BIN[type]}${version}-`;
    manifestUrl = fallbackRoot + "rbxPkgManifest.txt";
    manifestRes = await fetch(manifestUrl);
  }

  if (!manifestRes.ok) {
    console.error(
      `  ${COLORS.red}✗${COLORS.reset} Failed to fetch manifest: ${manifestRes.status}`,
    );
    process.exit(1);
  }

  const manifestText = await manifestRes.text();
  const packages = manifestText
    .split("\n")
    .map((l) => l.trim())
    .filter((l) => l && l.endsWith(".zip"));

  let extractRoots =
    packages.includes("RobloxApp.zip") || packages.includes("RobloxApp.zip")
      ? EXTRACT_ROOTS.player
      : EXTRACT_ROOTS.studio;

  if (packages.includes("RobloxApp.zip") && type === "WindowsStudio64") {
    console.error(
      `  ${COLORS.red}✗${COLORS.reset} BinaryType ${type} doesn't match manifest (RobloxApp.zip found)`,
    );
    process.exit(1);
  }

  if (packages.includes("RobloxStudio.zip") && type === "WindowsPlayer") {
    console.error(
      `  ${COLORS.red}✗${COLORS.reset} BinaryType ${type} doesn't match manifest (RobloxStudio.zip found)`,
    );
    process.exit(1);
  }

  const outputZip = new JSZip();

  outputZip.file(
    "AppSettings.xml",
    `<?xml version="1.0" encoding="UTF-8"?>
<Settings>
	<ContentFolder>content</ContentFolder>
	<BaseUrl>http://www.roblox.com</BaseUrl>
</Settings>
`,
  );

  for (const pkg of packages) {
    const pkgData = await stream(root + pkg, pkg);
    const extractRoot = extractRoots[pkg as keyof typeof extractRoots] || "";

    if (!extractRoot && !(pkg in extractRoots)) {
      console.log(
        `  ${COLORS.yellow}⚠${COLORS.reset} Package not in extract roots, adding to root: ${pkg}`,
      );
      outputZip.file(pkg, pkgData);
      continue;
    }

    const pkgZip = await JSZip.loadAsync(pkgData);

    for (const [path, file] of Object.entries(pkgZip.files)) {
      if (file.dir) continue;

      const fixedPath = path.replace(/\\/g, "/");
      const data = await file.async("arraybuffer");
      outputZip.file(extractRoot + fixedPath, data);
    }
  }

  console.log(
    `  ${COLORS.dim}⟳${COLORS.reset} Generating output zip${compress ? ` (compression level ${level}/9)` : ""}...`,
  );

  const outputData = await outputZip.generateAsync({
    type: "arraybuffer",
    compression: compress ? "DEFLATE" : "STORE",
    compressionOptions: { level },
  });

  const out = `${channel}-${type}-${version}.zip`;
  await Bun.write(out, Buffer.from(outputData));
  console.log(
    `  ${COLORS.green}✓${COLORS.reset} Saved to ${COLORS.bright}${out}${COLORS.reset}`,
  );
}

console.log(`
${COLORS.bright}${COLORS.cyan}┌─ rdd-cli${COLORS.reset}
${COLORS.cyan}├─${COLORS.reset} Channel: ${COLORS.bright}${channel}${COLORS.reset}
${COLORS.cyan}├─${COLORS.reset} Binary:  ${COLORS.bright}${type}${COLORS.reset}
${COLORS.cyan}├─${COLORS.reset} Version: ${COLORS.bright}${version}${COLORS.reset}
${COLORS.cyan}└─${COLORS.reset} Downloading...
`);

async function main() {
  const isMac = type.startsWith("Mac");

  if (isMac) {
    await downloadMac();
  } else {
    await downloadWindows();
  }
}

main().catch((err) => {
  console.error(`  ${COLORS.red}✗${COLORS.reset} Error: ${err.message}`);
  process.exit(1);
});
