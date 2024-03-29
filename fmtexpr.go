package motor

import (
	"regexp"
	"sort"
)

var (
	// normal verbs: %a %+b %#-1c
	normfmt = regexp.MustCompile(`%([#\+\-\d\.0\s]+)?[A-aZ-z%]`)
	// named verbs: %(name).2f
	namedfmt = regexp.MustCompile(`%\((\w+)\)(([#\+\-\d\.0\s]+)?([A-aZ-z]))`)
	// avoid repeated computation of the tag map
	fmtexprCache = map[string]tagmap{}
)

// Tag is a structured log element.
//
// They are extracted from named format string operands
// supported by out printf implementation.
//
// %(label)s, %(label)#+v
//
// See: Brr.Printf()
type Tag struct {
	Pos   int
	Label string
}

// tagmap is a (format string, tags) pair
type tagmap struct {
	format string
	tags   []Tag
}

// This function turns a format string and a bunch of named operands
// into a tagmap— a traditional fmt string along with the identified
// tags.
func fmtexpr(format string) tagmap {
	if pair, ok := fmtexprCache[format]; ok {
		return pair
	}

	// calculating positions of both normal and named ops
	normalpos := normfmt.FindAllStringIndex(format, -1)
	namedpos := namedfmt.FindAllStringIndex(format, -1)
	// extracting named ops names
	named := namedfmt.FindAllStringSubmatch(format, -1)

	type verbpos struct {
		pos   int
		named bool
	}

	verbs := make([]verbpos, 0, len(normalpos)+len(namedpos))
	for _, v := range normalpos {
		verbs = append(verbs, verbpos{v[0], false})
	}
	for _, v := range namedpos {
		verbs = append(verbs, verbpos{v[1], true})
	}

	// ordering operand positions
	sort.Slice(verbs, func(i, j int) bool {
		return verbs[i].pos < verbs[j].pos
	})

	// arranging named ops and their abs position
	var tags []Tag
	for i, verb := range verbs {
		if !verb.named {
			continue
		}
		// named[i][1] is the capture group of op's name
		tags = append(tags, Tag{i, named[len(tags)][1]})
		if len(tags) == len(namedpos) {
			break
		}
	}

	// converting named ops to regular ops
	regularFormat := namedfmt.ReplaceAllString(format, "%$2")

	fmtexprCache[format] = tagmap{regularFormat, tags}
	return fmtexprCache[format]
}
