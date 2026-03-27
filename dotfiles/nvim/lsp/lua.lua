return {
	cmd = { "lua-language-server" },
	root_markers = { "init.lua", ".git" },
	filetypes = { "lua" },
	settings = {
		Lua = {
			diagnostics = {
				-- Get the language server to recognize the `vim` global
				globals = {
					"vim",
					"require",
				},
			},
		},
	},
}
