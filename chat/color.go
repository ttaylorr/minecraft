package chat

type Color string

var (
	Colors []Color = []Color{
		ColorBlack, ColorDarkBlue, ColorDarkGreen, ColorDarkAqua,
		ColorDarkRed, ColorPurple, ColorGold, ColorGray, ColorDarkGray,
		ColorBlue, ColorGreen, ColorAqua, ColorRed, ColorLightPurple,
		ColorYellow, ColorWhite,
	}
)

const (
	ColorBlack       Color = "black"
	ColorDarkBlue          = "dark_blue"
	ColorDarkGreen         = "dark_green"
	ColorDarkAqua          = "dark_aqua"
	ColorDarkRed           = "dark_red"
	ColorPurple            = "dark_purple"
	ColorGold              = "gold"
	ColorGray              = "gray"
	ColorDarkGray          = "dark_Gray"
	ColorBlue              = "blue"
	ColorGreen             = "green"
	ColorAqua              = "aqua"
	ColorRed               = "red"
	ColorLightPurple       = "light_purple"
	ColorYellow            = "yellow"
	ColorWhite             = "white"
)
