package chat

type (
	TextComponent struct {
		Text string `json:"text"`

		Component
	}

	TranslateComponent struct {
		Translate string   `json:"translate"`
		With      []string `json:"with"`

		Component
	}

	ScoreComponent struct {
		Name      string `json:"name"`
		Objective string `json:"objective"`

		Component
	}

	SelectorComponent struct {
		Selector string   `json:"selector"`
		Args     []string `json:"extra"`

		Component
	}

	Component struct {
		Bold          bool `json:"bold"`
		Italic        bool `json:"italic"`
		Underlined    bool `json:"underlined"`
		Strikethrough bool `json:"strikethrough"`
		Obfuscated    bool `json:"obfuscated"`

		Color Color `json:"color"`

		ClickEvent ClickEvent `json:"clickEvent"`
		HoverEvent HoverEvent `json:"hoverEvent"`

		Insertion string `json:"insertion"`
	}
)
