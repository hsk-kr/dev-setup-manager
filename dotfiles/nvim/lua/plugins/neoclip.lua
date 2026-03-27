return {
	{
		"AckslD/nvim-neoclip.lua",
		dependencies = {
			-- you'll need at least one of these
			{ "nvim-telescope/telescope.nvim" },
			{ "kkharji/sqlite.lua", module = "sqlite" },
			-- {'ibhagwan/fzf-lua'},
		},
		config = function()
			require("neoclip").setup({
				history = 100,
				enable_persistent_history = true,
			})
		end,
	},
}
