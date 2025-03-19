package app

import "fmt"

func Guide() {
	fmt.Println("- Change the Homebrew click key to Shift + Command + F.")
	fmt.Println("- These commands assume you're using zsh. After installing wezterm, use your default terminal to install other software, then verify the results in wezterm.")
	fmt.Println("- When launching Karabiner Elements for the first time, it may reset the configuration, requiring you to reconfigure your dotfiles.")
	fmt.Println("- Homebrew should be installed by manually running the provided shell command.")
	fmt.Println("- You'll need to install Go, nvm, and Ruby to set up Language Server Protocol (LSP).")
	fmt.Println("- Remember to run source commands when needed (e.g., `source ~/.zshrc` or `NVM_PATH`).")
}
