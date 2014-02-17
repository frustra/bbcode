// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import (
	"fmt"
	"regexp"
	"strconv"
)

// Compiler represents the base interface for a BBCode compiler.
// Implement this and call CompileCustom to override default behaviour.
// You may embed DefaultCompiler if you just need to make some small tweaks.
type Compiler interface {
	Compile(node *BBCodeNode) *HTMLTag
	CompileRaw(node *BBCodeNode) *HTMLTag
}

type DefaultCompiler struct{}

var youtubeRegex = regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)\/(?:watch\?v=)?([a-zA-Z0-9]+)`)

// compile transforms a tag and subexpression into an HTML string.
// It is only used by the generated parser code.
func (c DefaultCompiler) Compile(node *BBCodeNode) *HTMLTag {
	var out = NewHTMLTag("")
	if node.ID == TEXT {
		out.Value = node.Value.(string)
		insertNewlines(out)
		for _, child := range node.Children {
			out.AppendChild(c.Compile(child))
		}
	} else if node.ID == CLOSING_TAG {
		out.Value = node.Value.(bbClosingTag).Raw
		for _, child := range node.Children {
			out.AppendChild(c.Compile(child))
		}
	} else if node.ClosingTag == nil {
		out.Value = node.Value.(bbOpeningTag).Raw
		for _, child := range node.Children {
			out.AppendChild(c.Compile(child))
		}
	} else {
		in := node.Value.(bbOpeningTag)
		var expr *HTMLTag
		if len(node.Children) == 1 {
			expr = c.Compile(node.Children[0])
		} else if len(node.Children) > 1 {
			expr = NewHTMLTag("")
			for _, child := range node.Children {
				expr.AppendChild(c.Compile(child))
			}
		}

		switch in.Name {
		case "url":
			out.Name = "a"
			if in.Value == "" {
				if expr != nil {
					out.Attrs["href"] = safeURL(expr.Value)
				} else {
					out.Attrs["href"] = ""
				}
			} else {
				out.Attrs["href"] = safeURL(in.Value)
			}
			out.AppendChild(expr)
		case "img":
			out.Name = "img"
			if in.Value == "" {
				if expr != nil {
					out.Attrs["src"] = safeURL(expr.Value)
				} else {
					out.Attrs["src"] = ""
				}
			} else {
				out.Attrs["src"] = safeURL(in.Value)
				if expr != nil {
					out.Attrs["alt"] = expr.Value
				}
			}
		case "media":
			if expr == nil {
				out.Value = "Embedded video"
			} else {
				out.Name = "div"
				out.Attrs["class"] = "embedded-video"

				obj := NewHTMLTag("Embedded video")
				out.AppendChild(obj)

				matches := youtubeRegex.FindStringSubmatch(expr.Value)
				if matches != nil {
					obj = NewHTMLTag("")
					obj.Name = "object"
					obj.Attrs["width"] = "620"
					obj.Attrs["height"] = "349"

					params := map[string]string{
						"movie":             fmt.Sprintf("//www.youtube.com/v/%s?version=3", matches[1]),
						"wmode":             "transparent",
						"allowFullScreen":   "true",
						"allowscriptaccess": "always",
					}

					embed := NewHTMLTag("")
					embed.Name = "embed"
					embed.Attrs["type"] = "application/x-shockwave-flash"
					embed.Attrs["width"] = "620"
					embed.Attrs["height"] = "349"
					for name, value := range params {
						param := NewHTMLTag("")
						param.Name = "param"
						param.Attrs["name"] = name
						param.Attrs["value"] = value
						obj.AppendChild(param)

						if name == "movie" {
							name = "src"
						}
						embed.Attrs[name] = value
					}
					obj.AppendChild(embed)
					out.AppendChild(obj)
				}
			}
		case "center":
			out.Name = "div"
			out.Attrs["style"] = "text-align: center;"
			out.AppendChild(expr)
		case "color":
			return expr
		case "size":
			out.Name = "span"
			if size, err := strconv.Atoi(in.Value); err == nil {
				out.Attrs["style"] = fmt.Sprintf("font-size: %dpx;", size*4)
			}
			out.AppendChild(expr)
		case "spoiler":
			out.Name = "div"
			out.Attrs["class"] = "expandable collapsed"
			out.AppendChild(expr)
		case "quote":
			out.Name = "blockquote"
			who := ""
			if name, ok := in.Args["name"]; ok && name != "" {
				who = name
			} else {
				who = in.Value
			}
			cite := NewHTMLTag("")
			cite.Name = "cite"
			if who != "" {
				cite.AppendChild(NewHTMLTag(who + " said:"))
			} else {
				cite.AppendChild(NewHTMLTag("Quote"))
			}
			out.AppendChild(cite)
			out.AppendChild(expr)
		case "strike":
			out.Name = "s"
			out.AppendChild(expr)
		case "code":
			out.Name = "code"
			for _, child := range node.Children {
				out.AppendChild(c.CompileRaw(child))
			}
		case "i", "b", "u":
			out.Name = in.Name
			out.AppendChild(expr)
		default:
			out.Value = in.Raw
			insertNewlines(out)
			out.AppendChild(expr).AppendChild(NewHTMLTag("[/" + in.Name + "]"))
		}
	}
	return out
}

func (c DefaultCompiler) CompileRaw(in *BBCodeNode) *HTMLTag {
	out := NewHTMLTag("")
	if in.ID == TEXT {
		out.Value = in.Value.(string)
	} else if in.ID == OPENING_TAG {
		out.Value = in.Value.(bbOpeningTag).Raw
	} else if in.ID == CLOSING_TAG {
		out.Value = in.Value.(bbClosingTag).Raw
	}
	insertNewlines(out)
	for _, child := range in.Children {
		out.AppendChild(c.CompileRaw(child))
	}
	if in.ID == OPENING_TAG {
		out.AppendChild(NewHTMLTag("[/" + in.Value.(bbOpeningTag).Name + "]"))
	}
	return out
}
