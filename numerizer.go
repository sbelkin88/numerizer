// Package numerizer parses words to integers.
package numerizer

import (
	"fmt"
	"regexp"
	"strings"
)

var singleNumbers = map[string]int{
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

var directNumbers = map[string]int{
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
}

var tenPrefixNumbers = map[string]int{
	"twenty":  20,
	"thirty":  30,
	"forty":   40,
	"fifty":   50,
	"sixty":   60,
	"seventy": 70,
	"eighty":  80,
	"ninety":  90,
}

var largeNumbers = map[string]int{
	"thousand": 1000,
	"million":  1000000,
	"billion":  1000000000,
}

type item struct {
	typ itemType
	key string
	val int
}

type itemType int

const (
	itemError itemType = iota
	itemDollars
	itemCents
	itemZero
	itemHundred
	itemSingle
	itemDirect
	itemTenPrefix
	itemLarge
	itemEOL
	itemDefault
)

var (
	dollarRegex = regexp.MustCompile(`dollars?.?`)
	centsRegex  = regexp.MustCompile(`cents?.?`)
)

func newItem(s string) item {
	switch {
	case s == "zero":
		return item{typ: itemZero, key: s}
	case s == "hundred":
		return item{typ: itemHundred, key: s, val: 100}
	case dollarRegex.Match([]byte(s)):
		return item{typ: itemDollars, key: s}
	case centsRegex.Match([]byte(s)):
		return item{typ: itemCents, key: s}
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
	return item{typ: itemDefault, key: s}
}

type parser struct {
	items       []item
	pos         int
	prev        int
	sum         int
	dollars     int
	cents       int
	dollarsDone bool
	err         error
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
	s = strings.ToLower(s)
	s = strings.Replace(s, ",", "", -1)
	s = strings.Replace(s, ".", "", -1)
	s = strings.Replace(s, "-", " ", -1)
	return s
}

type parseFn func(*parser) parseFn

// Parse parses the provided string to an integer.
func Parse(s string) (int, error) {
	p := newParser(s)
	for state := parse; state != nil; {
		state = state(p)
	}
	if p.err != nil {
		return 0, p.err
	}

	amountInt := p.dollars * 100
	amountInt += p.cents
	return amountInt, nil
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
	case itemDefault:
		return parseDefault
	}
	return parseError(p, i, item{})
}

func parseZero(p *parser) parseFn {
	i := p.next()
	next := p.peek()
	switch next.typ {
	case itemDollars:
		return parseDollars
	case itemCents:
		return parseCents
	case itemEOL:
		return nil
	}
	return parseError(p, next, i)
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
	case itemDollars:
		p.sum += p.prev + i.val
		return parseDollars
	case itemCents:
		p.sum += p.prev + i.val
		return parseCents
	case itemDefault:
		return parseDefault
	case itemEOL:
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
	case itemDollars:
		p.sum += p.prev + i.val
		return parseDollars
	case itemCents:
		p.sum += p.prev + i.val
		return parseCents
	case itemDefault:
		return parseDefault
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
		case itemDollars:
			p.sum += p.prev
			return parseDollars
		case itemCents:
			p.sum += p.prev
			return parseCents
		case itemDefault:
			p.sum += p.prev
			return parseDefault
		default:
			return parseSingle
		}
	case itemLarge:
		p.prev = i.val
		return parseLarge
	case itemDefault:
		return parseDefault
	case itemDollars:
		p.sum += p.prev + i.val
		return parseDollars
	case itemCents:
		p.sum += p.prev + i.val
		return parseCents
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
	case itemDollars:
		p.sum += p.prev
		return parseDollars
	case itemDefault:
		return parseDefault
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
	case itemDollars:
		return parseDollars
	case itemDefault:
		return parseDefault
	}
	return parseError(p, next, i)
}

func parseDollars(p *parser) parseFn {
	p.dollarsDone = true
	p.dollars = p.sum
	p.sum = 0
	p.prev = 0
	i := p.next()
	next := p.peek()
	switch next.typ {
	case itemSingle:
		return parseSingle
	case itemDirect:
		return parseDirect
	case itemTenPrefix:
		return parseTenPrefix
	case itemDefault:
		return parseDefault
	case itemEOL:
		return nil
	}
	return parseError(p, next, i)
}

func parseCents(p *parser) parseFn {
	p.cents = p.sum
	p.sum = 0
	p.prev = 0

	return nil
}

func parseDefault(p *parser) parseFn {
	i := p.next()
	next := p.peek()
	switch next.typ {
	case itemZero:
		return parseZero
	case itemSingle:
		return parseSingle
	case itemDirect:
		return parseDirect
	case itemTenPrefix:
		return parseTenPrefix
	case itemHundred:
		return parseHundred
	case itemDefault:
		return parseDefault
	case itemDollars:
		return parseDollars
	case itemCents:
		return parseCents
	case itemEOL:
		if p.dollarsDone {
			p.cents = p.sum
		}
		return nil
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
