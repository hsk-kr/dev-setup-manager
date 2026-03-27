-- Minimal neovim init for tmux-md markdown viewer
-- Uses already-compiled treesitter parsers + tokyonight from lazy install
-- Renders markdown with render-markdown.nvim

local script_dir = vim.fn.fnamemodify(debug.getinfo(1, "S").source:sub(2), ":h")

-- Add lazy-installed plugins to runtimepath
local lazy_path = vim.fn.stdpath("data") .. "/lazy"
vim.opt.rtp:prepend(lazy_path .. "/tokyonight.nvim")
vim.opt.rtp:prepend(lazy_path .. "/nvim-treesitter")
vim.opt.rtp:prepend(lazy_path .. "/mini.icons")

-- Add render-markdown.nvim from tmux-md plugins dir
vim.opt.rtp:prepend(script_dir .. "/plugins/render-markdown.nvim")

-- Basic options
vim.opt.number = false
vim.opt.relativenumber = false
vim.opt.wrap = true
vim.opt.linebreak = true
vim.opt.cursorline = false
vim.opt.signcolumn = "no"
vim.opt.fillchars = { eob = " " }
vim.opt.termguicolors = true
vim.opt.laststatus = 0
vim.opt.cmdheight = 0
vim.opt.scrolloff = 5
vim.opt.conceallevel = 2
vim.opt.sidescrolloff = 8
vim.opt.foldcolumn = "4"

-- Colorscheme
pcall(function()
  vim.cmd.colorscheme("tokyonight")
end)

-- Enable treesitter highlighting for markdown
pcall(function()
  require("nvim-treesitter.configs").setup({
    highlight = { enable = true },
  })
end)

-- Setup render-markdown for pretty rendering
pcall(function()
  require("render-markdown").setup({
    render_modes = { "n", "v", "i", "c" },
    heading = {
      icons = { "# ", "## ", "### ", "#### ", "##### ", "###### " },
    },
  })
end)

-- Keymaps: q and Esc both close
vim.keymap.set("n", "q", "<cmd>q!<cr>", { buffer = false, silent = true })
vim.keymap.set("n", "<Esc>", "<cmd>q!<cr>", { buffer = false, silent = true })

-- Read-only settings applied after buffer loads
vim.api.nvim_create_autocmd("BufEnter", {
  pattern = "*.md",
  callback = function()
    vim.bo.filetype = "markdown"
    vim.bo.modifiable = false
    vim.bo.readonly = true
  end,
})
