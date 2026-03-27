-- tmux-md neovim integration
-- Provides :TmuxMd command and <leader>md visual mapping

local TMUX_MD_SCRIPT = vim.fn.expand("~/dev-setup-manager/dotfiles/tmux-md/scripts/tmux-md.sh")

-- :TmuxMd — open current file in tmux popup
vim.api.nvim_create_user_command("TmuxMd", function()
  local file = vim.fn.expand("%:p")
  if file == "" then
    vim.notify("No file open", vim.log.levels.WARN)
    return
  end
  vim.fn.system(string.format("'%s' '%s' &", TMUX_MD_SCRIPT, file))
end, { desc = "Open current markdown file in tmux popup viewer" })

-- <leader>md in visual mode — select a file path, open it in tmux popup
vim.keymap.set("v", "<leader>md", function()
  -- Get visual selection
  vim.cmd('noau normal! "vy')
  local selected = vim.fn.getreg("v")
  selected = vim.trim(selected)

  if selected == "" then
    vim.notify("No text selected", vim.log.levels.WARN)
    return
  end

  -- Resolve relative paths from current buffer's directory
  if not selected:match("^/") then
    local buf_dir = vim.fn.expand("%:p:h")
    selected = buf_dir .. "/" .. selected
  end

  if vim.fn.filereadable(selected) == 0 then
    vim.notify("File not found: " .. selected, vim.log.levels.ERROR)
    return
  end

  vim.fn.system(string.format("'%s' '%s' &", TMUX_MD_SCRIPT, selected))
end, { desc = "Open selected file path in tmux-md popup" })
