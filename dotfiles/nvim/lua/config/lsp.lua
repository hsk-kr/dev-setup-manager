vim.keymap.set("n", "gd", vim.lsp.buf.definition, { noremap = true, silent = true })
vim.keymap.set("n", "gr", vim.lsp.buf.references, { noremap = true, silent = true })
vim.keymap.set("n", "gk", vim.lsp.buf.hover, { noremap = true, silent = true })
vim.keymap.set("n", "ga", vim.lsp.buf.code_action, { noremap = true, silent = true })

vim.opt.winborder = "rounded"

local lsp_languages = { "typescript", "lua", "go", "tailwind" }

for _, lsp_language in pairs(lsp_languages) do
	vim.lsp.enable(lsp_language)
end

vim.diagnostic.config({
	virtual_text = true, -- show inline messages
	signs = true, -- show signs in the gutter
	underline = true, -- underline problematic text
	update_in_insert = false, -- don't update diagnostics while typing
	severity_sort = true, -- sort diagnostics by severity
})
