package layer

import (
	"regexp"
	"strings"
)

const defaultRestrict = `[^/]+`

type StdRegExpMaker struct {
	regExpPlaceholders *regexp.Regexp
}

func (maker *StdRegExpMaker) initRegExp() *regexp.Regexp {
	if maker.regExpPlaceholders == nil {
		maker.regExpPlaceholders = regexp.MustCompile(`(?mU){(?P<name>.*)(:(?P<restrict>.*))?}`)
	}

	return maker.regExpPlaceholders
}

func (maker *StdRegExpMaker) MakeRegExp(l Layer) *regexp.Regexp {
	regexpPath := l.Path
	if regexpPath == `` {
		return nil
	}

	re := maker.initRegExp()
	nameIndex := re.SubexpIndex(`name`)
	restrictIndex := re.SubexpIndex(`restrict`)

	matched := re.FindAllStringSubmatch(regexpPath, -1)

	for _, match := range matched {
		name := match[nameIndex]
		restrict, ok := l.Restrictions[name]

		if !ok {
			restrict = match[restrictIndex]
		}

		if restrict == `` {
			restrict = defaultRestrict
		}

		replace := "(?P<" + name + ">" + restrict + ")"
		regexpPath = strings.ReplaceAll(regexpPath, match[0], replace)
	}

	return regexp.MustCompile(`(?mU)^` + regexpPath + `$`)
}
