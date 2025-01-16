package tools

func IsTmuxInstalled() bool {
	return ExistCommand("tmux")
}

func InstallTmux() {
	ExecCommand("brew", "install", "tmux")
}
