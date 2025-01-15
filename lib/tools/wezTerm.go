package tools

func IsWezTermInstalled() bool {
	return ExistApplication("WezTerm.app")
}

func InstallWezTerm() {
	ExecCommand("brew", "install", "--cask", "wezterm")
}
