//go:build windows
// +build windows

package main

import (
	"context"
	_ "embed"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/google/go-github/v69/github"
	"golang.org/x/sys/windows/registry"
)

var FLAGS_NO_UPDATE string = "NO_UPDATES"

func ForceUpdate() {
	client := github.NewClient(nil)

	files, _, err := client.Repositories.GetLatestRelease(context.Background(), "DHCPCD9", "go-steamvr-lighthouse-manager")

	if err != nil {
		log.Println("Failed to check updates: " + err.Error())
		return
	}

	//Finding the installer

	for _, v := range files.Assets {
		if strings.HasSuffix(*v.Name, "installer.exe") {
			//Downloading the release
			response, err := http.Get(*v.BrowserDownloadURL)

			if err != nil {
				log.Println("Failed to download latest release")
				return
			}

			bodyBytes, _ := io.ReadAll(response.Body)

			os.WriteFile(path.Join(GetConfigFolder(), "installer.exe"), bodyBytes, 0644)

			cmd := exec.Command("cmd", "/k", strings.ReplaceAll(path.Join(GetConfigFolder(), "installer.exe"), "/", "\\"))
			go cmd.Output()

			time.Sleep(time.Second)
			os.Exit(0)

		}
	}
}

func IsUpdatingSupported() bool {

	//probing regedit
	_, err := registry.OpenKey(registry.LOCAL_MACHINE, "SOFTWARE\\Microsoft\\Windows\\CurrentVersion\\Uninstall\\Alumi Inc.Base Station Manager", registry.QUERY_VALUE)

	return err == nil && !strings.Contains(FLAGS_NO_UPDATE, VERSION_FLAGS)
}
