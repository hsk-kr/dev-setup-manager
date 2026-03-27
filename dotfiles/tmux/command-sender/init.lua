local file = vim.fn.expand("~/.tmux_command_sender")
local target_window = vim.env.CS_TARGET_WINDOW or ""

-- Open the persistent command file
vim.cmd("edit " .. vim.fn.fnameescape(file))

-- Minimal UI
vim.opt.number = false
vim.opt.relativenumber = false
vim.opt.signcolumn = "no"
vim.opt.laststatus = 0
vim.opt.cmdheight = 1
vim.opt.shortmess:append("I")

local function send_to_pane(pane_index)
  local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
  local text = table.concat(lines, "\n")
  if text == "" then
    return
  end
  vim.cmd("silent write")
  local target = string.format("%s.%d", target_window, pane_index)
  vim.fn.system({ "tmux", "load-buffer", "-" }, text)
  vim.fn.system({ "tmux", "paste-buffer", "-t", target })
  vim.fn.system({ "tmux", "send-keys", "-t", target, "C-m" })
  vim.cmd("quit")
end

-- 0-5 sends to that pane
for i = 0, 5 do
  vim.keymap.set("n", tostring(i), function()
    send_to_pane(i)
  end, { silent = true })
end

-- q or Esc to close
vim.keymap.set("n", "q", "<cmd>silent write<cr><cmd>quit<cr>", { silent = true })
vim.keymap.set("n", "<Esc>", "<cmd>silent write<cr><cmd>quit<cr>", { silent = true })
