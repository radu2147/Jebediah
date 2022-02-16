package utils

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	_ "golang.org/x/sys/windows/svc"
)

const RegistryKey = "Windows11_Scanner"

func AddToStartup(file string) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE|registry.SET_VALUE)
	defer k.Close()
	if err != nil {
		return err
	}
	err = k.SetStringValue(RegistryKey, fmt.Sprintf("\"%s\" /V", file))
	if err != nil {
		return err
	}
	return nil
}

func RemoveFromStartup() error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.QUERY_VALUE|registry.SET_VALUE)
	defer k.Close()
	if err != nil {
		return err
	}
	err = k.DeleteValue(RegistryKey)
	if err != nil {
		return err
	}
	return nil
}
