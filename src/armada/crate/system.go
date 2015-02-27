package main

import (
	"fmt"
	//"os/exec"
	"strings"
)

func ip(args ...string) (string, error) {
	return fmt.Sprintf("ip %s", strings.Join(args, " ")), nil
}

func ipLink(args ...string) (string, error) {
	ipargs := append([]string{"link"}, args...)
	return ip(ipargs...)
}

//Wrapper around the ip command
/*
func ip(args ...string) (string, error) {
	path, err := exec.LookPath("ip")
	if err != nil {
		return "", fmt.Errorf("command not found: ip")
	}
	output, err := exec.Command(path, args...).CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ip failed: ip %v", strings.Join(args, " "))
	}
	return string(output), nil
}

// Wrapper around the iptables command
func iptables(args ...string) error {
	path, err := exec.LookPath("iptables")
	if err != nil {
		return fmt.Errorf("command not found: iptables")
	}
	if err := exec.Command(path, args...).Run(); err != nil {
		return fmt.Errorf("iptables failed: iptables %v", strings.Join(args, " "))
	}
	return nil
}
*/
