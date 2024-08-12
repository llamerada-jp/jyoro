package jyoro

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"slices"
	"strings"
)

const (
	UHUBCTL_ACTION_ON   = "on"
	UHUBCTL_ACTION_OFF  = "off"
	UHUBCTL_ACTION_STAT = ""
)

type USB interface {
	IsON(location string, port uint) (bool, error)
	Power(location string, port uint, on bool) error
}

type usb struct {
}

func NewUSB() USB {
	u := &usb{}
	return u
}

// TODO: check real USB status
func (u *usb) IsON(location string, port uint) (bool, error) {
	out, err := u.runHubCtrl(location, port, UHUBCTL_ACTION_STAT)
	if err != nil {
		return false, err
	}
	return u.decodeStatus(out)
}

func (u *usb) Power(location string, port uint, on bool) error {
	act := UHUBCTL_ACTION_OFF
	if on {
		act = UHUBCTL_ACTION_ON
	}
	out, err := u.runHubCtrl(location, port, act)
	if err != nil {
		return err
	}
	after, err := u.decodeStatus(out)
	if err != nil {
		return err
	}

	if on != after {
		return fmt.Errorf("failed to turn %s USB", act)
	}
	return nil
}

func (u *usb) decodeStatus(out string) (bool, error) {
	lines := regexp.MustCompile(`\r\n|\n`).Split(out, -1)
	lines = slices.DeleteFunc(lines, func(s string) bool {
		return strings.TrimSpace(s) == ""
	})
	if len(lines) == 0 {
		return false, fmt.Errorf("no output from uhubctl")
	}
	tail := strings.TrimSpace(lines[len(lines)-1])
	cols := regexp.MustCompile(`\s+`).Split(tail, -1)
	if len(cols) < 4 {
		return false, fmt.Errorf("invalid output from uhubctl")
	}
	if cols[3] == "power" {
		return true, nil
	} else if cols[3] == "off" {
		return false, nil
	} else {
		return false, fmt.Errorf("unknown status from uhubctl")
	}
}

func (u *usb) runHubCtrl(location string, port uint, act string) (string, error) {
	uhubctlCmd := []string{
		"uhubctl",
		"-l", location,
		"-p", fmt.Sprintf("%d", port),
		"-f",
	}
	if act != "" {
		uhubctlCmd = append(uhubctlCmd, "-a", act)
	}

	out, err := exec.Command("sudo", uhubctlCmd...).Output()
	if err != nil {
		log.Println(string(out))
		return "", err
	}
	return string(out), nil
}
