// Package numerizer parses words to integers.
package numerizer

import (
	"fmt"
	"strings"
)

var singleNumbers = map[string]int64{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
}

var directNumbers = map[string]int64{
	"ten":       10,
	"eleven":    11,
	"twelve":    12,
	"thirteen":  13,
	"fourteen":  14,
	"fifteen":   15,
	"sixteen":   16,
	"seventeen": 17,
	"eighteen":  18,
	"nineteen":  19,
	"ninteen":   19, // misspelling
}

var tenPrefixNumbers = map[string]int64{
	"twenty":  20,
	"thirty":  30,
	"forty":   40,
	"fourty":  40, // misspelling
	"fifty":   50,
	"sixty":   60,
	"seventy": 70,
	"eighty":  80,
	"ninety":  90,
}

var largeNumbers = map[string]int64{
	"thousand": 1000,
	"million":  1000000,
	"billion":  1000000000,
	"trillion": 1000000000000,
}

type item struct {
	typ itemType
	key string
	val int64
}

type itemType int

const (
	itemError itemType = iota
	itemAnd
	itemZero
	itemHundred
	itemSingle
	itemDirect
	itemTenPrefix
	itemLarge
	itemEOL
)

func newItem(s string) item {
	switch s {
	case "zero":
		return item{typ: itemZero, key: s}
	case "hundred":
		return item{typ: itemHundred, key: s, val: 100}
	case "and":
		return item{typ: itemAnd, key: s}
	}
	return newItemFromMap(s)
}

func newItemFromMap(s string) item {
	if v, ok := singleNumbers[s]; ok {
		return item{typ: itemSingle, key: s, val: v}
	} else if v, ok := directNumbers[s]; ok {
		return item{typ: itemDirect, key: s, val: v}
	} else if v, ok := tenPrefixNumbers[s]; ok {
		return item{typ: itemTenPrefix, key: s, val: v}
	} else if v, ok := largeNumbers[s]; ok {
		return item{typ: itemLarge, key: s, val: v}
	}
	return item{key: s}
}

func (i item) String() string {
	if i.typ == itemEOL {
		return "EOL"
	}
	return i.key
}

type parser struct {
	items []item
	pos   int
	prev  int64
	sum   int64
	err   error
}

func newParser(s string) *parser {
	s = preprocess(s)
	parts := strings.Fields(s)
	items := make([]item, 0, len(parts))
	for _, part := range parts {
		items = append(items, newItem(part))
	}
	return &parser{items: items}
}

func (p *parser) next() item {
	next := p.peek()
	p.pos++
	return next
}

func (p *parser) peek() item {
	if p.pos >= len(p.items) {
		return item{typ: itemEOL}
	}
	return p.items[p.pos]
}

func preprocess(s string) string {
	s = strings.Replace(s, ",", "", -1)
	s = strings.Replace(s, "-", " ", -1)
	return s
}

type parseFn func(*parser) parseFn

// Parse parses the provided string to an integer.
func Parse(s string) (int64, error) {
	p := newParser(s)
	for state := parse; state != nil; {
		state = state(p)
	}
	if p.err != nil {
		return 0, p.err
	}
	return p.sum, nil
}

func parse(p *parser) parseFn {
	i := p.peek()
	switch i.typ {
	case itemZero:
		return parseZero
	case itemSingle:
		return parseSingle
	case itemDirect:
		return parseDirect
	case itemTenPrefix:
		return parseTenPrefix
	}
	return parseError(p, i, item{})
}

func parseZero(p *parser) parseFn {
	i := p.next()
	next := p.peek()
	if next.typ != itemEOL {
		return parseError(p, next, i)
	}
	return nil
}

func parseSingle(p *parser) parseFn {
	i := p.next()
	next := p.peek()
	switch next.typ {
	case itemHundred:
		p.sum += p.prev
		p.prev = i.val
		return parseHundred
	case itemLarge:
		p.prev = i.val
		return parseLarge
	case itemEOL:
		p.sum += p.prev + i.val
		return nil
	}
	return parseError(p, next, i)
}

func parseDirect(p *parser) parseFn {
	i := p.next()
	next := p.peek()
	switch next.typ {
	case itemHundred:
		p.prev = i.val
		return parseLarge
	case itemLarge:
		p.prev = i.val
		return parseLarge
	case itemEOL:
		p.sum += p.prev + i.val
		return nil
	}
	return parseError(p, next, i)
}

func parseTenPrefix(p *parser) parseFn {
	i := p.next()
	next := p.peek()
	switch next.typ {
	case itemSingle:
		p.prev += i.val + next.val
		p.next()
		ahead := p.peek()
		switch ahead.typ {
		case itemHundred:
			p.sum += p.prev * ahead.val
			p.prev = 0
			return parseHundred
		case itemLarge:
			return parseLarge
		default:
			return parseSingle
		}
	case itemLarge:
		p.prev = i.val
		return parseLarge
	case itemEOL:
		p.sum += p.prev + i.val
		return nil
	}
	return parseError(p, next, i)
}

func parseHundred(p *parser) parseFn {
	i := p.next()
	p.prev *= i.val
	next := p.peek()
	switch next.typ {
	case itemSingle:
		return parseSingle
	case itemDirect:
		return parseDirect
	case itemTenPrefix:
		return parseTenPrefix
	case itemLarge:
		return parseLarge
	case itemAnd:
		return parseAnd
	case itemEOL:
		p.sum += p.prev
		return nil
	}
	return parseError(p, next, i)
}

func parseLarge(p *parser) parseFn {
	i := p.next()
	p.sum += p.prev * i.val
	p.prev = 0
	next := p.peek()
	switch next.typ {
	case itemSingle:
		return parseSingle
	case itemDirect:
		return parseDirect
	case itemTenPrefix:
		return parseTenPrefix
	case itemAnd:
		return parseAnd
	case itemEOL:
		return nil
	}
	return parseError(p, next, i)
}

func parseAnd(p *parser) parseFn {
	i := p.next()
	next := p.peek()
	switch next.typ {
	case itemSingle:
		return parseSingle
	case itemDirect:
		return parseDirect
	case itemTenPrefix:
		return parseTenPrefix
	}
	return parseError(p, next, i)
}

func parseError(p *parser, i item, after item) parseFn {
	switch i.typ {
	case itemError:
		p.err = fmt.Errorf("bad number: %q", i.key)
	default:
		if after.typ == itemError {
			p.err = fmt.Errorf("unexpected start %q", i.key)
		} else {
			p.err = fmt.Errorf("unexpected %q after %q", i.key, after.key)
		}
	}
	return nil
}
