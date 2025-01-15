package tools

func IsNeovimInstalled() bool {
	return ExistCommand("nvim")
}

func InstallNeovim() {
	ExecCommand("brew", "install", "neovim")
}
