// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import (
	"fmt"
	"html"
	"net/url"
	"regexp"
	"strconv"
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
	var value string
	if t.value != "" {
		value = sanitize(t.value)
	}
	var attrString string
	for key, value := range t.attrs {
		attrString = fmt.Sprintf(`%s %s="%s"`, attrString, key, escapeQuotes(sanitize(value)))
	}
	if len(t.children) > 0 {
		var childrenString string
		for _, child := range t.children {
			childrenString = fmt.Sprint(childrenString, child.string())
		}
		if t.name != "" {
			return fmt.Sprintf(`%s<%s%s>%s</%s>`, value, t.name, attrString, childrenString, t.name)
		} else {
			return fmt.Sprint(value, childrenString)
		}
	} else if t.name != "" {
		return fmt.Sprintf(`%s<%s%s>`, value, t.name, attrString)
	} else {
		return value
	}
}

func (t *htmlTag) appendChild(child *htmlTag) *htmlTag {
	if child == nil {
		t.children = append(t.children, newHtmlTag(""))
	} else {
		t.children = append(t.children, child)
	}
	return t
}

var youtubeRegex = regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)\/(?:watch\?v=)?([a-zA-Z0-9]+)`)

// compile transforms a tag and subexpression into an HTML string.
// It is only used by the generated parser code.
func compile(node *BBCodeNode) *htmlTag {
	var out = newHtmlTag("")
	if node.id == TEXT {
		out.value = node.value.(string)
		if strings.ContainsRune(out.value, '\n') {
			parts := strings.Split(out.value, "\n")
			for i, part := range parts {
				if i == 0 {
					out.value = parts[i]
				} else {
					out.appendChild(newline()).appendChild(newHtmlTag(part))
				}
			}
		}
		for _, child := range node.Children {
			out.appendChild(compile(child))
		}
	} else if node.id == CLOSING_TAG {
		out.value = node.value.(bbClosingTag).raw
		for _, child := range node.Children {
			out.appendChild(compile(child))
		}
	} else {
		in := node.value.(bbOpeningTag)
		var expr *htmlTag
		if len(node.Children) == 1 {
			expr = compile(node.Children[0])
		} else if len(node.Children) > 1 {
			expr = newHtmlTag("")
			for _, child := range node.Children {
				expr.appendChild(compile(child))
			}
		}

		switch in.name {
		case "url":
			out.name = "a"
			if in.value == "" {
				if expr != nil {
					out.attrs["href"] = safeURL(expr.value)
				} else {
					out.attrs["href"] = ""
				}
			} else {
				out.attrs["href"] = safeURL(in.value)
			}
			out.appendChild(expr)
		case "img":
			out.name = "img"
			if in.value == "" {
				if expr != nil {
					out.attrs["src"] = safeURL(expr.value)
				} else {
					out.attrs["src"] = ""
				}
			} else {
				out.attrs["src"] = safeURL(in.value)
				if expr != nil {
					out.attrs["alt"] = expr.value
				}
			}
		case "media":
			if expr == nil {
				out.value = "Embedded video"
			} else {
				out.name = "div"
				out.attrs["class"] = "embedded-video"

				obj := newHtmlTag("Embedded video")
				out.appendChild(obj)

				matches := youtubeRegex.FindStringSubmatch(expr.value)
				if matches != nil {
					obj = newHtmlTag("")
					obj.name = "object"
					obj.attrs["width"] = "620"
					obj.attrs["height"] = "349"

					params := map[string]string{
						"movie":             fmt.Sprintf("//www.youtube.com/v/%s?version=3", matches[1]),
						"wmode":             "transparent",
						"allowFullScreen":   "true",
						"allowscriptaccess": "always",
					}

					embed := newHtmlTag("")
					embed.name = "embed"
					embed.attrs["type"] = "application/x-shockwave-flash"
					embed.attrs["width"] = "620"
					embed.attrs["height"] = "349"
					for name, value := range params {
						param := newHtmlTag("")
						param.name = "param"
						param.attrs["name"] = name
						param.attrs["value"] = value
						obj.appendChild(param)

						if name == "movie" {
							name = "src"
						}
						embed.attrs[name] = value
					}
					obj.appendChild(embed)
					out.appendChild(obj)
				}
			}
		case "center":
			out.name = "div"
			out.attrs["style"] = "text-align: center;"
			out.appendChild(expr)
		case "color":
			return expr
		case "size":
			out.name = "span"
			if size, err := strconv.Atoi(in.value); err == nil {
				out.attrs["style"] = fmt.Sprintf("font-size: %dpx;", size*4)
			}
			out.appendChild(expr)
		case "spoiler":
			out.name = "div"
			out.attrs["class"] = "expandable collapsed"
			out.appendChild(expr)
		case "quote":
			out.name = "blockquote"
			who := ""
			if name, ok := in.args["name"]; ok && name != "" {
				who = name
			} else {
				who = in.value
			}
			cite := newHtmlTag("")
			cite.name = "cite"
			if who != "" {
				cite.appendChild(newHtmlTag(who + " said:"))
			} else {
				cite.appendChild(newHtmlTag("Quote"))
			}
			out.appendChild(cite)
			out.appendChild(expr)
		case "strike":
			out.name = "s"
			out.appendChild(expr)
		case "i", "b", "u", "code":
			out.name = in.name
			out.appendChild(expr)
		default:
			out.value = in.raw
			out.appendChild(expr)
		}
	}
	return out
}

func newline() *htmlTag {
	var out = newHtmlTag("")
	out.name = "br"
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
