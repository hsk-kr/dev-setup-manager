local M = {}

local command_file = vim.fn.expand("~/.command_sender_tmux")

local function get_buffer_text()
  local lines = vim.api.nvim_buf_get_lines(0, 0, -1, false)
  return table.concat(lines, "\n")
end

local function get_visual_text()
  local start_pos = vim.fn.getpos("v")
  local end_pos = vim.fn.getpos(".")
  local start_row, start_col = start_pos[2], start_pos[3]
  local end_row, end_col = end_pos[2], end_pos[3]
  local mode = vim.fn.mode()

  if start_row > end_row or (start_row == end_row and start_col > end_col) then
    start_row, end_row = end_row, start_row
    start_col, end_col = end_col, start_col
  end

  local lines = vim.api.nvim_buf_get_lines(0, start_row - 1, end_row, false)
  if #lines == 0 then
    return ""
  end

  local is_linewise = mode == "V"
  if not is_linewise then
    lines[1] = string.sub(lines[1], start_col)
    lines[#lines] = string.sub(lines[#lines], 1, end_col)
  end

  return table.concat(lines, "\n")
end

local function send_text_to_tmux(pane_index, text)
  local target = string.format(":.%d", pane_index)
  vim.fn.system({ "tmux", "load-buffer", "-" }, text)
  vim.fn.system({ "tmux", "paste-buffer", "-t", target })
  vim.fn.system({ "tmux", "send-keys", "-t", target, "C-m" })
end

local function send_from_mode(pane_index, use_visual)
  local text = use_visual and get_visual_text() or get_buffer_text()
  if text == "" then
    return
  end
  vim.cmd("silent write")
  send_text_to_tmux(pane_index, text)
end

local function set_buffer_keymaps(bufnr)
  local opts = { buffer = bufnr, silent = true }
  vim.keymap.set("n", "0", function()
    send_from_mode(0, false)
  end, opts)
  vim.keymap.set("x", "0", function()
    send_from_mode(0, true)
  end, opts)
  vim.keymap.set("n", "1", function()
    send_from_mode(1, false)
  end, opts)
  vim.keymap.set("x", "1", function()
    send_from_mode(1, true)
  end, opts)
  vim.keymap.set("n", "2", function()
    send_from_mode(2, false)
  end, opts)
  vim.keymap.set("x", "2", function()
    send_from_mode(2, true)
  end, opts)
  vim.keymap.set("n", "3", function()
    send_from_mode(3, false)
  end, opts)
  vim.keymap.set("x", "3", function()
    send_from_mode(3, true)
  end, opts)
  vim.keymap.set("n", "4", function()
    send_from_mode(4, false)
  end, opts)
  vim.keymap.set("x", "4", function()
    send_from_mode(4, true)
  end, opts)
end

local function open_command_buffer()
  vim.cmd("rightbelow vsplit " .. vim.fn.fnameescape(command_file))
  local bufnr = vim.api.nvim_get_current_buf()
  set_buffer_keymaps(bufnr)
end

-- Autocommand to ensure keymaps are set whenever the command buffer is entered
local augroup = vim.api.nvim_create_augroup("CommandSenderTmux", { clear = true })
vim.api.nvim_create_autocmd("BufEnter", {
  group = augroup,
  pattern = vim.fn.expand("~/.command_sender_tmux"),
  callback = function(args)
    set_buffer_keymaps(args.buf)
  end,
})

vim.api.nvim_create_user_command("SendCommandToTmuxPanel", function()
  open_command_buffer()
end, {})

vim.keymap.set("n", "<leader>st", "<cmd>SendCommandToTmuxPanel<cr>", { silent = true })

-- Global keybindings to send text from any buffer
for i = 0, 5 do
  vim.keymap.set("n", "<leader>cs" .. i, function()
    send_from_mode(i, false)
  end, { silent = true, desc = "Send buffer to tmux pane " .. i })
  vim.keymap.set("x", "<leader>cs" .. i, function()
    send_from_mode(i, true)
  end, { silent = true, desc = "Send selection to tmux pane " .. i })
end

return M
