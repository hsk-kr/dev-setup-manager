local wezterm = require("wezterm")

return {
	window_background_opacity = 1.0,
	macos_window_background_blur = 20,
	-- window_background_image = wezterm.config_dir .. "/backgrounds/zenitsu.jpg",
	background = {
		{
			source = {
				File = wezterm.config_dir .. "/backgrounds/zenitsu.jpg",
			},
			hsb = {
				brightness = 0.4, -- Adjust brightness (0.0 - 1.0)
				hue = 0.3,
				saturation = 0.7,
			},
			opacity = 0.2, -- Adjust transparency (0.0 - 1.0)
			-- width = "100%",
			-- height = "100%",
			repeat_x = "NoRepeat",
			repeat_y = "NoRepeat",
			vertical_align = "Middle",
			horizontal_align = "Center",
		},
	},
}
