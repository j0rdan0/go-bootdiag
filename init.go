package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
)

func downloadFile(url string) (string, error) {
	dir, err := os.MkdirTemp("", "viu")
	if err != nil {
		log.Fatalf("failed creating temp dir, %s\n", err)
		return "", err
	}
	defer os.RemoveAll(dir)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("failed getting viu binary, %s\n", err)
		return "", err
	}
	filename := dir + "/viu"
	f, err := os.Create(filename)
	defer f.Close()
	buff, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed getting file %s\n", err)
		return "", err
	}
	_, err = f.Write(buff)
	if err != nil {
		log.Fatalf("failed saving binary %s\n", err)
		return "", err
	}
	return filename, nil
}

func moveToPath(file string) {
	path := strings.Split(os.Getenv(strings.ToUpper("path")), ":")[0]
	newPath := path + "/" + file
	os.Rename(file, newPath)
}

func fetchLinuxVui() {
	filename, _ := downloadFile(viu_binary_x86_64_linux)
	os.Chmod(filename, 0744)
	moveToPath(filename)

}

func fetchWindowsInstaller() {
	filename, _ := downloadFile(viu_installer_windows)
	moveToPath(filename)

}

func setup() {

	switch os := runtime.GOOS; os {
	case "linux":
		{
			fetchLinuxVui()
			break
		}
	case "darwin":
		{
			log.Println("run brew install vui for installation")
			break
		}
	case "windows":
		{
			fetchWindowsInstaller()
			break
		}
	default:
		log.Println("unsupported OS")

	}
}
