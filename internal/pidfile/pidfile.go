// Package pidfile provides structure and helper functions to create and remove
// PID file. A PID file is usually a file used to store the process ID of a
// running process.
package pidfile

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// PIDFile is a file used to store the process ID of a running process.
type PIDFile struct {
	path string
}

func checkPIDFileAlreadyExists(path string) error {
	if pidByte, err := os.ReadFile(path); err == nil {
		pidString := strings.TrimSpace(string(pidByte))
		if pid, err := strconv.Atoi(pidString); err == nil {
			if processExists(pid) {
				return fmt.Errorf("pid file found, ensure webhook is not running or delete %s", path)
			}
		}
	}
	return nil
}

// New creates a PIDfile using the specified path.
func New(path string) (*PIDFile, error) {
	if err := checkPIDFileAlreadyExists(path); err != nil {
		return nil, err
	}
	// Note MkdirAll returns nil if a directory already exists
	if err := MkdirAll(filepath.Dir(path), os.FileMode(0o755)); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, []byte(fmt.Sprintf("%d", os.Getpid())), 0o600); err != nil {
		return nil, err
	}

	return &PIDFile{path: path}, nil
}

// Remove removes the PIDFile.
func (file PIDFile) Remove() error {
	return os.Remove(file.path)
}

func (p *PIDFile) Write() error {
	return os.WriteFile(p.path, []byte(fmt.Sprintf("%d\n", os.Getpid())), 0o644)
}
