package network

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func (network *Network) postConfigure(ctx context.Context, config *Config, tincBin string) error {
	var interfaces = map[string]bool{}

	list, err := net.Interfaces()
	if err != nil {
		return err
	}

	for _, iface := range list {
		if iface.Name == config.Interface {
			return nil
		}
		log.Println("found interface:", iface.Name)
		interfaces[iface.Name] = true
	}

	// install tap
	tapInstaller, err := network.findTapInstall(tincBin)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(ctx, tapInstaller, "install", "OemWin2k.inf", "tap0901")
	cmd.Dir = filepath.Dir(tapInstaller)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return err
	}

	// find new interface
	var newInterface string
	for newInterface == "" {
		log.Println("looking for a new interface")
		select {
		case <-time.After(1 * time.Second):
		case <-ctx.Done():
		}
		list, err = net.Interfaces()
		if err != nil {
			return err
		}
		for _, iface := range list {
			if !interfaces[iface.Name] {
				newInterface = iface.Name
				log.Println("new interface:", iface.Name)
				break
			}
		}
	}

	if newInterface == "" {
		return fmt.Errorf("new interface not found")
	}

	// rename
	cmd = exec.CommandContext(ctx, "netsh", "interface", "set", "interface",
		"name", "=", newInterface, "newname", "=", config.Interface)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	return cmd.Run()
}

func (network *Network) findTapInstall(tincBin string) (string, error) {
	var res string
	err := filepath.Walk(filepath.Dir(tincBin), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if info.Name() == "tapinstall.exe" {
			res = path
			return os.ErrExist
		}
		return nil
	})
	if err == os.ErrExist {
		err = nil
	} else if err == nil {
		err = os.ErrNotExist
	}
	return res, err
}

func (network *Network) beforeConfigure(config *Config) error {
	return nil
}
