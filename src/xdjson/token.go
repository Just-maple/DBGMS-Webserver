package xdjson



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
