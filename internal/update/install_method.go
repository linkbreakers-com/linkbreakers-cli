package update

import (
	"os"
	"path/filepath"
	"strings"
)

type InstallMethod string

const (
	InstallMethodDirect   InstallMethod = "direct"
	InstallMethodHomebrew InstallMethod = "homebrew"
	InstallMethodNpm      InstallMethod = "npm"
)

func DetectInstallMethod() InstallMethod {
	execPath, err := os.Executable()
	if err != nil {
		return InstallMethodDirect
	}
	if resolved, err := filepath.EvalSymlinks(execPath); err == nil {
		execPath = resolved
	}
	return installMethodForPath(execPath)
}

func installMethodForPath(path string) InstallMethod {
	p := strings.ReplaceAll(path, `\`, "/")
	switch {
	case strings.Contains(p, "/node_modules/linkbreakers-cli/"):
		return InstallMethodNpm
	case strings.Contains(p, "/Caskroom/"), strings.Contains(p, "/Cellar/"):
		return InstallMethodHomebrew
	default:
		return InstallMethodDirect
	}
}

func (m InstallMethod) UpdateCommand() string {
	switch m {
	case InstallMethodHomebrew:
		return "brew upgrade linkbreakers"
	case InstallMethodNpm:
		return "npm install -g linkbreakers-cli@latest"
	default:
		return "linkbreakers self-update"
	}
}
