package xdjson

type Json struct {
	Map         interface{}
	parser      parser
	listParser  listParser
	tokenParser *tokenParser
}

type parser struct {
	tokens    tokens
	pointer   int
	key       string
	subTokens tokens
	value     interface{}
	node      jsonNode
}

type jsonNode struct {
	Keys []string
	Map  map[string]interface{}
}

type listParser struct {
	tokens    tokens
	array     []interface{}
	subTokens tokens
	pointer   int
}

type tokens []string

func Init(str string) (j Json) {
	j.tokenParser = &tokenParser{raw: str}
	j.tokenParser.parse()
	j.parse()
	return
}

func (j *Json) parse() {
	var tks = j.tokenParser.tokens[1:]
	switch {
	case j.tokenParser.tokens.isJson():
		j.Map = j.parser.parse(tks).node
	case j.tokenParser.tokens.isArray():
		j.Map = j.listParser.parse(tks).array
	default:
		panic("invalid json string input")
	}
}

func (ps parser) parse(token tokens) *parser {
	ps.tokens = token
	ps.node = jsonNode{Map: make(map[string]interface{})}
	for ps.pointer < len(ps.tokens) {
		var tk = ps.tokens[ps.pointer]
		ps.subTokens = ps.tokens[ps.pointer:]
		switch {
		case ps.subTokens.isNewJson():
			ps.key = ps.subTokens[1]
			var tnode = new(parser).parse(ps.subTokens[5:])
			ps.value = tnode.node
			ps.pointer += tnode.pointer + 5
		case tk == "{":
			ps.key = ps.subTokens[2]
			var tnode = new(parser).parse(ps.subTokens[5:])
			ps.value = tnode.node
			ps.pointer += tnode.pointer + 5
		case tk == "[":
			var tnode = new(listParser).parse(ps.subTokens[1:])
			ps.value = tnode.array
			ps.pointer += tnode.pointer
		case ps.subTokens.isKey():
			ps.key = ps.subTokens[1]
			ps.pointer += 4
		case ps.subTokens.isNullString():
			ps.value = ""
			ps.pointer += 2
		case tk == ",":
			ps.pointer += 1
		case ps.subTokens.isString():
			ps.value = ps.subTokens[1]
			ps.pointer += 3
		case tk == "}":
			ps.pointer += 1
			return &ps
		default:
			ps.value = tk
			ps.pointer += 1
		}
		if ps.key != "" && ps.value != nil {
			ps.node.Map[ps.key] = ps.value
			ps.node.Keys = append(ps.node.Keys, ps.key)
			ps.key = ""
			ps.value = nil
		}
	}
	panic(ps)
	return &ps
}

func (ps listParser) parse(token tokens) (node *listParser) {
	ps.tokens = token
	for ps.pointer < len(ps.tokens) {
		var tk = ps.tokens[ps.pointer]
		ps.subTokens = ps.tokens[ps.pointer:]
		switch {
		case tk == "[":
			var tnode = new(listParser).parse(ps.subTokens[1:])
			ps.pointer += tnode.pointer + 1
			ps.array = append(ps.array, tnode.array)
		case tk == "]":
			ps.pointer += 1
			return &ps
		case tk == "{":
			var tnode = new(parser).parse(ps.subTokens[1:])
			ps.pointer += tnode.pointer + 1
			ps.array = append(ps.array, tnode.node)
		case ps.subTokens.isNullString():
			ps.array = append(ps.array, "")
		case tk == ",":
			ps.pointer += 1
		case ps.subTokens.isString():
			ps.array = append(ps.array, ps.subTokens[1])
			ps.pointer += 3
		default:
			ps.array = append(ps.array, tk)
			ps.pointer += 1
		}
	}
	panic(ps)
	return &ps
}
