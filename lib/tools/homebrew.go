package tools

func IsHomebrewInstalled() bool {
	return ExistCommand("brew")
}

// TODO: Not tested if it installs homebrew correctly
func InstallHomebrew() {
	ExecCommand("/bin/bash", "-c", "\"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"")
}
