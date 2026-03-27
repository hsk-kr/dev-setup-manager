return {
	cmd = { "typescript-language-server", "--stdio" },
	cmd_env = { NODE_OPTIONS = "--max-old-space-size=4096" },
	root_markers = { "package.json", ".git", "tsconfig.json" },
	filetypes = { "typescript", "typescriptreact", "javascript", "javascriptreact" },
}
