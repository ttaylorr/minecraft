package chat

type (
	ClickAction string
	HoverAction string

	ClickEvent struct {
		Action ClickAction
		Value  string
	}

	HoverEvent struct {
		Action HoverAction
		Value  interface{}
	}
)

const (
	OpenUrlClickAction        ClickAction = "open_url"
	OpenFileClickAction                   = "open_file"
	RunCommandClickAction                 = "run_command"
	SuggestCommandClickAction             = "suggest_command"

	ShowTextHoverAction HoverAction = "show_text"
	ShowAchievement                 = "show_achievement"
	ShowItem                        = "show_item"
)
