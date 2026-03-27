return {
  "leoluz/nvim-dap-go",
  ft = "go",
  config = function()
    require("dap-go").setup({
      -- Windows quirk: set detached = false :contentReference[oaicite:1]{index=1}
      delve = { detached = vim.fn.has("win32") == 0 },
    })
  end,
  dependencies = { "mfussenegger/nvim-dap" },
}
