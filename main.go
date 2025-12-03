package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

//go:embed wisim-nouveau/build
//go:embed wisim-nouveau/build/client/_app
//go:embed wisim-nouveau/build/server/chunks/_*
var content embed.FS

var (
	platform = runtime.GOOS
	arch     = runtime.GOARCH
)

func extractEmbedFS(fsys embed.FS, root string, destDir string) error {
	return fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(destDir, relPath)

		if d.IsDir() {
			return os.MkdirAll(destPath, 0o700)
		}

		data, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}

		return os.WriteFile(destPath, data, 0o644)
	})
}

func main() {
	tempDir, err := os.MkdirTemp("", "wisim-temp-*")
	if err != nil {
		log.Fatalf("Failed to create temporary directory: %s", err.Error())
	}

	println(tempDir)

	defer os.RemoveAll(tempDir)

	err = extractEmbedFS(content, "wisim-nouveau/build", tempDir)
	if err != nil {
		log.Fatalf("failed to extract embed FS: %s", err.Error())
	}

	println(platform)
	println(arch)

	if platform == "linux" || platform == "darwin" {
		os.Chmod(tempDir, 0755)
	}

	var wisimServerCMD *exec.Cmd
	switch platform {
	case "linux", "darwin":
		binPath := fmt.Sprintf("%s/client/wisimserver-%s-%s", tempDir, platform, arch)
		os.Chmod(binPath, 0755)
		wisimServerCMD = exec.Command(binPath, "8000", "10", "1")
	case "windows":
		wisimServerCMD = exec.Command(fmt.Sprintf("%s\\client\\wisimserver-%s-%s.exe", tempDir, platform, arch), "8000", "10", "1")
	}

	var nodeServerCMD *exec.Cmd
	switch platform {
	case "linux", "darwin":
		binPath := fmt.Sprintf("%s/client/node-%s-%s", tempDir, platform, arch)
		os.Chmod(binPath, 0755)
		nodeServerCMD = exec.Command(binPath, fmt.Sprintf("%s/index.js", tempDir))
	case "windows":
		nodeServerCMD = exec.Command(fmt.Sprintf("%s\\client\\node-%s-%s.exe", tempDir, platform, arch), fmt.Sprintf("%s\\index.js", tempDir))
	}

	done := make(chan bool, 1)

	go runPipe(done, nodeServerCMD)
	go runPipe(done, wisimServerCMD)

	switch platform {
	case "linux", "darwin":
		exec.Command("xdg-open", "http://0.0.0.0:3000").Run()
	case "windows":
		exec.Command("start", "http://0.0.0.0:3000").Run()
	}

	<-done

	os.Exit(0)
}

func runPipe(done chan bool, cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Run()
	if err != nil {
		log.Fatal(err.Error())
	}

	done <- true
}
