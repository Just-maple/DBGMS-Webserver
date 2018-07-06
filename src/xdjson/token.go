package xdjson

type tokenParser struct {
	tokens     tokens
	subTokens  tokens
	tkStack    string
	strFlag    bool
	escapeFlag bool
	raw        string
}

func (tp *tokenParser) parse() {
	for _, c := range tp.raw {
		char := string(c)
		switch {
		case char == "\\":
			if tp.escapeFlag {
				tp.escapeFlag = true
			} else {
				tp.pushToken(char)
			}
		case char == `"` ||
			char == `{` ||
			char == `:` ||
			char == `}` ||
			char == `,` ||
			char == `]` ||
			char == `[`:
			{
				tp.pushToken(char)
			}
		default:
			tp.tkStack += char
		}

	}
	return
}

func (tp *tokenParser) clearStack(c string) {
	if len(tp.tkStack) > 0 {
		tp.tokens = append(tp.tokens, tp.tkStack)
	}
	tp.tkStack = ""
	tp.tokens = append(tp.tokens, c)
}

func (tp *tokenParser) pushToken(c string) {
	if tp.escapeFlag {
		tp.tkStack += c
		tp.escapeFlag = false
	} else if tp.strFlag {
		if c == "" {
			tp.strFlag = false
			tp.clearStack(c)
		} else {
			tp.tkStack += c
		}
	} else {
		tp.clearStack(c)
		if c == "" {
			tp.strFlag = true
		}
	}
}

func (tk tokens) isJson() bool {
	return tk[0] == "{" && tk[1] == `"` && tk[3] == `"` && tk[4] == `:`
}

func (tk tokens) isNullString() bool {
	return tk[0] == `"` && tk[1] == `"`
}
func (tk tokens) isArray() bool {
	return tk[0] == `[`
}
func (tk tokens) isString() bool {
	return tk[0] == `"` && tk[2] == `"`
}
func (tk tokens) isKey() bool {
	return tk[0] == `"` && tk[2] == `"` && tk[3] == `:`
}
func (tk tokens) isNewJson() bool {
	return tk.isKey() && tk[4:].isJson()
}
