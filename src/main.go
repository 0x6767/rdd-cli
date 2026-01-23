package main

import (
	"archive/zip"
	"bytes"
	"compress/flate"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const HOST = "https://setup-aws.rbxcdn.com"

var BIN = map[string]string{
	"WindowsPlayer":   "/",
	"WindowsStudio64": "/",
	"MacPlayer":       "/mac/",
	"MacStudio":       "/mac/",
}

var EXTRACT_ROOTS_PLAYER = map[string]string{
	"RobloxApp.zip":                     "",
	"redist.zip":                        "",
	"shaders.zip":                       "shaders/",
	"ssl.zip":                           "ssl/",
	"WebView2.zip":                      "",
	"WebView2RuntimeInstaller.zip":      "WebView2RuntimeInstaller/",
	"content-avatar.zip":                "content/avatar/",
	"content-configs.zip":               "content/configs/",
	"content-fonts.zip":                 "content/fonts/",
	"content-sky.zip":                   "content/sky/",
	"content-sounds.zip":                "content/sounds/",
	"content-textures2.zip":             "content/textures/",
	"content-models.zip":                "content/models/",
	"content-platform-fonts.zip":        "PlatformContent/pc/fonts/",
	"content-platform-dictionaries.zip": "PlatformContent/pc/shared_compression_dictionaries/",
	"content-terrain.zip":               "PlatformContent/pc/terrain/",
	"content-textures3.zip":             "PlatformContent/pc/textures/",
	"extracontent-luapackages.zip":      "ExtraContent/LuaPackages/",
	"extracontent-translations.zip":     "ExtraContent/translations/",
	"extracontent-models.zip":           "ExtraContent/models/",
	"extracontent-textures.zip":         "ExtraContent/textures/",
	"extracontent-places.zip":           "ExtraContent/places/",
}

var EXTRACT_ROOTS_STUDIO = map[string]string{
	"RobloxStudio.zip":                  "",
	"RibbonConfig.zip":                  "RibbonConfig/",
	"redist.zip":                        "",
	"Libraries.zip":                     "",
	"LibrariesQt5.zip":                  "",
	"WebView2.zip":                      "",
	"WebView2RuntimeInstaller.zip":      "",
	"shaders.zip":                       "shaders/",
	"ssl.zip":                           "ssl/",
	"Qml.zip":                           "Qml/",
	"Plugins.zip":                       "Plugins/",
	"StudioFonts.zip":                   "StudioFonts/",
	"BuiltInPlugins.zip":                "BuiltInPlugins/",
	"ApplicationConfig.zip":             "ApplicationConfig/",
	"BuiltInStandalonePlugins.zip":      "BuiltInStandalonePlugins/",
	"content-qt_translations.zip":       "content/qt_translations/",
	"content-sky.zip":                   "content/sky/",
	"content-fonts.zip":                 "content/fonts/",
	"content-avatar.zip":                "content/avatar/",
	"content-models.zip":                "content/models/",
	"content-sounds.zip":                "content/sounds/",
	"content-configs.zip":               "content/configs/",
	"content-api-docs.zip":              "content/api_docs/",
	"content-textures2.zip":             "content/textures/",
	"content-studio_svg_textures.zip":   "content/studio_svg_textures/",
	"content-platform-fonts.zip":        "PlatformContent/pc/fonts/",
	"content-platform-dictionaries.zip": "PlatformContent/pc/shared_compression_dictionaries/",
	"content-terrain.zip":               "PlatformContent/pc/terrain/",
	"content-textures3.zip":             "PlatformContent/pc/textures/",
	"extracontent-translations.zip":     "ExtraContent/translations/",
	"extracontent-luapackages.zip":      "ExtraContent/LuaPackages/",
	"extracontent-textures.zip":         "ExtraContent/textures/",
	"extracontent-scripts.zip":          "ExtraContent/scripts/",
	"extracontent-models.zip":           "ExtraContent/models/",
	"studiocontent-models.zip":          "StudioContent/models/",
	"studiocontent-textures.zip":        "StudioContent/textures/",
}

type Versions struct {
	Future struct {
		Windows     string `json:"Windows"`
		WindowsDate string `json:"WindowsDate"`
		Mac         string `json:"Mac"`
		MacDate     string `json:"MacDate"`
	} `json:"future"`
	Current struct {
		Windows     string `json:"Windows"`
		WindowsDate string `json:"WindowsDate"`
		Mac         string `json:"Mac"`
		MacDate     string `json:"MacDate"`
	} `json:"current"`
	Past struct {
		Windows     string `json:"Windows"`
		WindowsDate string `json:"WindowsDate"`
		Mac         string `json:"Mac"`
		MacDate     string `json:"MacDate"`
	} `json:"past"`
}

func help() {
	fmt.Println(`
rdd-cli

USAGE:
  rdd-cli -v <version> -t <binary> [-c <channel>] [--compress] [--level <1-9>]

OPTIONS:
  -v, --version         Version hash (example: version-xxxx)
  -t, --type            WindowsPlayer | WindowsStudio64 | MacPlayer | MacStudio
  -c, --channel         LIVE | zflag | production (default: LIVE)
  --compress            Compress output zip (default: false)
  --level               Compression level 1-9 (default: 6)
  -h, --help            Show help
  --version-cli         Show cli version

EXAMPLE:
  rdd-cli -v version-123abc -t WindowsPlayer
  rdd-cli -v version-123abc -t MacStudio -c LIVE
  rdd-cli -v version-123abc -t WindowsPlayer --compress --level 9
`)
}

func showVersions() {
	client := &http.Client{}
	versions := Versions{}

	urls := []string{
		"https://weao.xyz/api/versions/future",
		"https://weao.xyz/api/versions/current",
		"https://weao.xyz/api/versions/past",
	}

	for i, url := range urls {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "WEAO-3PService")
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			switch i {
			case 0:
				json.Unmarshal(body, &versions.Future)
			case 1:
				json.Unmarshal(body, &versions.Current)
			case 2:
				json.Unmarshal(body, &versions.Past)
			}
		}
	}

	fmt.Printf(`
Future Version
  Windows  %s (%s)
  Mac      %s (%s)

Current Version:
  Windows  %s (%s)
  Mac      %s (%s)

Past Version:
  Windows  %s (%s)
  Mac      %s (%s)
`, versions.Future.Windows, versions.Future.WindowsDate, versions.Future.Mac, versions.Future.MacDate,
		versions.Current.Windows, versions.Current.WindowsDate, versions.Current.Mac, versions.Current.MacDate,
		versions.Past.Windows, versions.Past.WindowsDate, versions.Past.Mac, versions.Past.MacDate)
}

func formatBytes(n int64) string {
	if n >= 1024*1024*1024 {
		return fmt.Sprintf("%.2f GB", float64(n)/1024/1024/1024)
	}
	if n >= 1024*1024 {
		return fmt.Sprintf("%.2f MB", float64(n)/1024/1024)
	}
	if n >= 1024 {
		return fmt.Sprintf("%.2f KB", float64(n)/1024)
	}
	return fmt.Sprintf("%d B", n)
}

func formatSpeed(bytesPerSec float64) string {
	if bytesPerSec >= 1024*1024 {
		return fmt.Sprintf("%.2f MB/s", bytesPerSec/1024/1024)
	}
	if bytesPerSec >= 1024 {
		return fmt.Sprintf("%.2f KB/s", bytesPerSec/1024)
	}
	return fmt.Sprintf("%.2f B/s", bytesPerSec)
}

func stream(url, name string) ([]byte, error) {
	fmt.Printf("  ↓ %s\n", name)

	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch: %s (%d)", url, resp.StatusCode)
	}

	total := resp.ContentLength
	reader := resp.Body

	var recv int64
	var parts [][]byte
	startTime := time.Now()

	buf := make([]byte, 32*1024)
	for {
		n, err := reader.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		recv += int64(n)
		parts = append(parts, append([]byte(nil), buf[:n]...))

		elapsed := time.Since(startTime).Seconds()
		if elapsed > 0 && total > 0 {
			pct := float64(recv) / float64(total)
			fill := int(pct * 28)
			bar := strings.Repeat("█", fill) + strings.Repeat("░", 28-fill)
			speed := float64(recv) / elapsed
			remaining := float64(total-recv) / speed
			eta := "--:--:--"
			if remaining > 0 {
				etaTime := time.Now().Add(time.Duration(remaining) * time.Second)
				eta = etaTime.Format("15:04:05")
			}
			fmt.Printf("\r  [%s] %5.1f%%  %10s / %10s  %10s  ETA %s",
				bar, pct*100, formatBytes(recv), formatBytes(total),
				formatSpeed(speed), eta)
		}
	}

	fmt.Println()
	result := make([]byte, 0, recv)
	for _, part := range parts {
		result = append(result, part...)
	}
	return result, nil
}

func downloadMac(version, binType, channel string) error {
	file := "RobloxApp.zip"
	if binType == "MacPlayer" {
		file = "RobloxPlayer.zip"
	} else if binType == "MacStudio" {
		file = "RobloxStudioApp.zip"
	}

	base := HOST + BIN[binType] + version + "-"
	data, err := stream(base+file, file)
	if err != nil {
		return err
	}

	out := fmt.Sprintf("%s-%s-%s.zip", channel, binType, version)
	err = os.WriteFile(out, data, 0644)
	if err != nil {
		return err
	}
	fmt.Printf("  ✓ Saved to %s\n", out)
	return nil
}

func downloadWindows(version, binType, channel string, compress bool, level int) error {
	base := HOST + BIN[binType] + version + "-"
	manifestUrl := base + "rbxPkgManifest.txt"

	resp, err := http.Get(manifestUrl)
	if err != nil || (resp != nil && resp.StatusCode != http.StatusOK) {
		if resp != nil {
			resp.Body.Close()
		}
		fallbackBase := HOST + "/channel/common" + BIN[binType] + version + "-"
		manifestUrl = fallbackBase + "rbxPkgManifest.txt"
		resp, err = http.Get(manifestUrl)
		if err != nil {
			return err
		}
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return fmt.Errorf("failed to fetch manifest: %d", resp.StatusCode)
	}

	manifestText, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	lines := strings.Split(string(manifestText), "\n")
	var packages []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && strings.HasSuffix(line, ".zip") {
			packages = append(packages, line)
		}
	}

	extractRoots := EXTRACT_ROOTS_PLAYER
	for _, pkg := range packages {
		if pkg == "RobloxStudio.zip" {
			extractRoots = EXTRACT_ROOTS_STUDIO
			break
		}
	}

	if binType == "WindowsStudio64" {
		for _, pkg := range packages {
			if pkg == "RobloxApp.zip" {
				return fmt.Errorf("BinaryType %s doesn't match manifest (RobloxApp.zip found)", binType)
			}
		}
	}

	if binType == "WindowsPlayer" {
		for _, pkg := range packages {
			if pkg == "RobloxStudio.zip" {
				return fmt.Errorf("BinaryType %s doesn't match manifest (RobloxStudio.zip found)", binType)
			}
		}
	}

	buf := new(bytes.Buffer)
	zipWriter := zip.NewWriter(buf)

	appSettings, _ := zipWriter.Create("AppSettings.xml")
	appSettings.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>
<Settings>
	<ContentFolder>content</ContentFolder>
	<BaseUrl>http://www.roblox.com</BaseUrl>
</Settings>
`))

	for _, pkg := range packages {
		pkgData, err := stream(base+pkg, pkg)
		if err != nil {
			return err
		}

		pkgZipReader, err := zip.NewReader(bytes.NewReader(pkgData), int64(len(pkgData)))
		if err != nil {
			fmt.Printf("  ⚠ Failed to read zip: %s\n", pkg)
			continue
		}

		extractRoot := extractRoots[pkg]
		if extractRoot == "" && pkg != "RobloxStudio.zip" && pkg != "RobloxApp.zip" && pkg != "shaders.zip" {
			fmt.Printf("  ⚠ Package not in extract roots, adding to root: %s\n", pkg)
			f, _ := zipWriter.Create(pkg)
			f.Write(pkgData)
			continue
		}

		for _, file := range pkgZipReader.File {
			if file.FileInfo().IsDir() {
				continue
			}

			rc, err := file.Open()
			if err != nil {
				continue
			}

			data, err := io.ReadAll(rc)
			rc.Close()
			if err != nil {
				continue
			}

			path := strings.ReplaceAll(file.Name, "\\", "/")
			zipPath := extractRoot + path
			f, err := zipWriter.Create(zipPath)
			if err != nil {
				continue
			}
			f.Write(data)
		}
	}

	zipWriter.Close()

	if compress {
		compressedBuf := new(bytes.Buffer)
		compressedWriter := zip.NewWriter(compressedBuf)
		compressedWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
			return flate.NewWriter(out, level)
		})

		zipReader, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
		if err != nil {
			return err
		}

		for _, file := range zipReader.File {
			rc, err := file.Open()
			if err != nil {
				continue
			}
			data, _ := io.ReadAll(rc)
			rc.Close()

			f, _ := compressedWriter.Create(file.Name)
			f.Write(data)
		}

		compressedWriter.Close()
		buf = compressedBuf
	}

	out := fmt.Sprintf("%s-%s-%s.zip", channel, binType, version)
	err = os.WriteFile(out, buf.Bytes(), 0644)
	if err != nil {
		return err
	}
	fmt.Printf("  ✓ Saved to %s\n", out)
	return nil
}

func main() {
	var version string
	var binType string
	var channel string
	var compress bool
	var level int
	var helpFlag bool
	var versionCli bool

	flag.StringVar(&version, "v", "", "Version hash (example: version-xxxx)")
	flag.StringVar(&binType, "t", "", "WindowsPlayer | WindowsStudio64 | MacPlayer | MacStudio")
	flag.StringVar(&channel, "c", "LIVE", "LIVE | zflag | production")
	flag.BoolVar(&compress, "compress", false, "Compress output zip")
	flag.IntVar(&level, "level", 6, "Compression level 1-9")
	flag.BoolVar(&helpFlag, "h", false, "Show help")
	flag.BoolVar(&versionCli, "version-cli", false, "Show cli version")

	flag.Parse()

	if helpFlag || len(os.Args) == 1 {
		help()
		showVersions()
		return
	}

	if versionCli {
		fmt.Println("rdd-cli v1.0.0")
		return
	}

	if version == "" || binType == "" {
		fmt.Println("Error: Missing required arguments")
		help()
		os.Exit(1)
	}

	if _, exists := BIN[binType]; !exists {
		fmt.Println("Error: Invalid binary type. Must be one of: WindowsPlayer, WindowsStudio64, MacPlayer, MacStudio")
		help()
		os.Exit(1)
	}

	fmt.Printf(`
┌─ rdd-cli
├─ Channel: %s
├─ Binary:  %s
├─ Version: %s
└─ Downloading...
`, channel, binType, version)

	isMac := strings.HasPrefix(binType, "Mac")

	var err error
	if isMac {
		err = downloadMac(version, binType, channel)
	} else {
		err = downloadWindows(version, binType, channel, compress, level)
	}

	if err != nil {
		fmt.Printf("  ✗ Error: %s\n", err.Error())
		os.Exit(1)
	}
}