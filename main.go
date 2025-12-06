package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"
)

//go:embed wisim-nouveau/build
//go:embed wisim-nouveau/build/client/_app
//go:embed wisim-nouveau/build/server/chunks/_page.*
//go:embed wisim-nouveau/build/server/chunks/_*
var content embed.FS

// go:embed wisim-nouveau/build/client/_app/immutable/chunks/_*

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
		err := os.Chmod(binPath, 0755)
		if err != nil {
			log.Fatalln(err)
		}
		wisimServerCMD = exec.Command(binPath, "8000", "10", "1")
	case "windows":
		wisimServerCMD = exec.Command(fmt.Sprintf("%s\\client\\wisimserver-%s-%s.exe", tempDir, platform, arch), "8000", "10", "1")
	}

	var nodeServerCMD *exec.Cmd
	switch platform {
	case "linux", "darwin":
		binPath := fmt.Sprintf("%s/client/node-%s-%s", tempDir, platform, arch)
		err := os.Chmod(binPath, 0755)
		if err != nil {
			log.Fatalln(err)
		}
		nodeServerCMD = exec.Command(binPath, fmt.Sprintf("%s/index.js", tempDir))
	case "windows":
		nodeServerCMD = exec.Command(fmt.Sprintf("%s\\client\\node-%s-%s.exe", tempDir, platform, arch), fmt.Sprintf("%s\\index.js", tempDir))
	}

	done := make(chan bool, 1)
	signals := make(chan os.Signal, 1)

	go runPipe(done, nodeServerCMD)
	go runPipe(done, wisimServerCMD)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	go func() {
		sig := <-signals
		fmt.Printf("\nRecieved signal: %s\n", sig.String())
		done <- true
	}()

	defer func() {
		nodeServerCMD.Process.Kill()
		wisimServerCMD.Process.Kill()
	}()

	time.Sleep(time.Second * 5)

	switch platform {
	case "linux":
		_ = exec.Command("xdg-open", "http://localhost:3000").Run()
	case "darwin":
		_ = exec.Command("open", "http://localhost:3000").Run()
	case "windows":
		_ = exec.Command("start", "http://localhost:3000").Run()
	}

	<-done
}

func runPipe(done chan bool, cmd *exec.Cmd) {
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stdout

	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		done <- false
		return
	}

	done <- true
}
