// Copyright 2014 Frustra Sofware. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import (
	"fmt"
	"html"
	"net/url"
	"strings"
)

type htmlTag struct {
	name     string
	value    string
	attrs    map[string]string
	children []*htmlTag
}

func newHtmlTag(value string) *htmlTag {
	return &htmlTag{
		value:    value,
		attrs:    make(map[string]string),
		children: make([]*htmlTag, 0),
	}
}

func (t *htmlTag) string() string {
	if t.value != "" {
		return sanitize(t.value)
	}
	attrStrings := make([]string, 0, len(t.attrs))
	for key, value := range t.attrs {
		attrStrings = append(attrStrings, fmt.Sprintf(`%s="%s"`, key, escapeQuotes(sanitize(value))))
	}
	attrString := strings.Join(attrStrings, " ")
	if len(t.children) > 0 {
		var childrenString string
		for i, child := range t.children {
			if i == 0 {
				childrenString = child.string()
			} else {
				childrenString = fmt.Sprint(childrenString, " ", child.string())
			}
		}
		return fmt.Sprintf(`<%s %s>%s</%s>`, t.name, attrString, childrenString, t.name)
	} else {
		return fmt.Sprintf(`<%s %s/>`, t.name, attrString)
	}
}

func (t *htmlTag) appendChild(child *htmlTag) {
	t.children = append(t.children, child)
}

// compile transforms a tag and subexpression into an HTML string.
// It is only used by the generated parser code.
func compile(in bbTag, expr *htmlTag) *htmlTag {
	var out = newHtmlTag("")

	switch in.key {
	case "url":
		out.name = "a"
		if in.value == "" {
			out.attrs["href"] = safeURL(expr.value)
		} else {
			out.attrs["href"] = safeURL(in.value)
		}
		out.appendChild(expr)
	case "img":
		out.name = "img"
		if in.value == "" {
			out.attrs["src"] = safeURL(expr.value)
		} else {
			out.attrs["src"] = safeURL(in.value)
			out.attrs["alt"] = expr.value
		}
	}
	return out
}

func escapeQuotes(raw string) string {
	return strings.Replace(strings.Replace(raw, `"`, `\"`, -1), `\`, `\\`, -1)
}

func safeURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return ""
	}
	return strings.Replace(u.String(), `\`, "%5C", -1)
}

func sanitize(raw string) string {
	return html.EscapeString(raw)
}
