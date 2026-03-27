return {
  cmd = { "tailwindcss-language-server", "--stdio" },
  root_markers = {
    "package.json",
    "tailwind.config.js",
    "tailwind.config.cjs",
    "tailwind.config.mjs",
    "tailwind.config.ts",
    "postcss.config.js",
    "postcss.config.cjs",
    "postcss.config.mjs",
    "postcss.config.ts",
  },
  filetypes = { "html", "typescriptreact", "javascriptreact" },
}
