package ledserver

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

const (
	VoltageMinThreshold float64 = 4.75
	VoltageMaxThreshold float64 = 5.30
)

var (
	ext5vPathOnce sync.Once
	ext5vPath     string
	ext5vPathErr  error
)

func readPi5Ext5V(sysfsPath string) (float64, error) {
	data, err := os.ReadFile(sysfsPath)
	if err != nil {
		return 0, err
	}
	mv, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
	if err != nil {
		return 0, err
	}
	return mv / 1000.0, nil
}

func findExt5VPath() (string, error) {
	matches, err := filepath.Glob("/sys/class/hwmon/hwmon*/in*_input")
	if err != nil {
		return "", err
	}
	for _, p := range matches {
		if v, err := readPi5Ext5V(p); err == nil && v > 4.0 && v < 6.0 {
			return p, nil
		}
	}
	return "", fmt.Errorf("no 5V hwmon input found")
}

func Voltage() (float64, error) {
	ext5vPathOnce.Do(func() {
		ext5vPath, ext5vPathErr = findExt5VPath()
	})
	if ext5vPathErr != nil {
		return 0, ext5vPathErr
	}
	return readPi5Ext5V(ext5vPath)
}
