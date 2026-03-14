package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hsk-kr/dev-setup-manager/lib/config"
	"github.com/hsk-kr/dev-setup-manager/lib/spinner"
)

func SetupDotfiles(dotCfg config.DotfilesConfig) error {
	homePath, err := os.UserHomeDir()
	if err != nil {
		WarningMessage(err.Error())
		return err
	}

	devSetupManagerHomePath := filepath.Join(homePath, "dev-setup-manager")
	configDirPath := filepath.Join(homePath, ".config")
	devSetupManagerDotfilesPath := filepath.Join(devSetupManagerHomePath, "dotfiles")

	if err := ExecCommand("mkdir", "-p", devSetupManagerHomePath); err != nil {
		return err
	}

	if exist, _ := existFile(devSetupManagerDotfilesPath); exist {
		if err := ExecCommand("rm", "-rf", devSetupManagerDotfilesPath); err != nil {
			return err
		}
	}

	sp := spinner.New("Cloning dotfiles...")
	sp.Start()
	err = ExecCommandQuiet("git", "clone", dotCfg.Repo, devSetupManagerDotfilesPath)
	sp.Stop()
	if err != nil {
		return err
	}

	if err := ExecCommand("mkdir", "-p", configDirPath); err != nil {
		return err
	}

	// Symlink config directories
	for _, item := range dotCfg.ConfigLinks {
		target := filepath.Join(configDirPath, item)
		// Remove existing directory (not symlink) so ln doesn't create a link inside it
		if info, err := os.Lstat(target); err == nil && info.IsDir() && info.Mode()&os.ModeSymlink == 0 {
			if err := os.RemoveAll(target); err != nil {
				return err
			}
		}
		if err := ExecCommand("ln", "-sfn", filepath.Join(devSetupManagerDotfilesPath, item), target); err != nil {
			return err
		}
	}

	// Symlink home directories
	for source, target := range dotCfg.HomeLinks {
		if err := ExecCommand("ln", "-sfn", filepath.Join(devSetupManagerDotfilesPath, source), filepath.Join(homePath, target)); err != nil {
			return err
		}
	}

	// Extra links (e.g., claude/skills -> ~/.claude/skills)
	for _, link := range dotCfg.ExtraLinks {
		targetPath := config.ExpandPath(link.Target)
		targetDir := filepath.Dir(targetPath)
		if err := ExecCommand("mkdir", "-p", targetDir); err != nil {
			return err
		}
		if err := ExecCommand("ln", "-sfn", filepath.Join(devSetupManagerDotfilesPath, link.Source), targetPath); err != nil {
			return err
		}
	}

	// Run post-setup scripts
	for _, script := range dotCfg.PostScripts {
		scriptPath := filepath.Join(devSetupManagerDotfilesPath, script)
		sp := spinner.New(fmt.Sprintf("Running %s...", script))
		sp.Start()
		err := ExecCommandQuiet("bash", scriptPath)
		sp.Stop()
		if err != nil {
			WarningMessage(fmt.Sprintf("Post script %s failed: %s", script, err.Error()))
		}
	}

	// Add zsh source
	if dotCfg.ZshSource != "" {
		zshSource := dotCfg.ZshSource
		// Expand ~ in the source path
		zshSource = config.ExpandPath(zshSource)
		return AddZshSource(fmt.Sprintf("source %s", zshSource))
	}

	return nil
}
