// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

import (
	"fmt"
	"regexp"
	"strconv"
)

type TagCompilerFunc func(*BBCodeNode, BBOpeningTag) (*HTMLTag, bool)

type Compiler struct {
	tagCompilers               map[string]TagCompilerFunc
	AutoCloseTags              bool
	IgnoreUnmatchedClosingTags bool
}

func NewCompiler(autoCloseTags, ignoreUnmatchedClosingTags bool) Compiler {
	compiler := Compiler{make(map[string]TagCompilerFunc), autoCloseTags, ignoreUnmatchedClosingTags}

	for tag, compilerFunc := range DefaultTagCompilers {
		compiler.SetTag(tag, compilerFunc)
	}
	return compiler
}

func (c Compiler) Compile(str string) string {
	tokens := Lex(str)
	tree := Parse(tokens)
	return c.CompileTree(tree).String()
}

func (c Compiler) SetTag(tag string, compiler TagCompilerFunc) {
	if compiler == nil {
		delete(c.tagCompilers, tag)
	} else {
		c.tagCompilers[tag] = compiler
	}
}

// CompileTree transforms BBCodeNode into an HTML tag.
func (c Compiler) CompileTree(node *BBCodeNode) *HTMLTag {
	var out = NewHTMLTag("")
	if node.ID == TEXT {
		out.Value = node.Value.(string)
		InsertNewlines(out)
		for _, child := range node.Children {
			out.AppendChild(c.CompileTree(child))
		}
	} else if node.ID == CLOSING_TAG {
		if !c.IgnoreUnmatchedClosingTags {
			out.Value = node.Value.(BBClosingTag).Raw
			InsertNewlines(out)
		}
		for _, child := range node.Children {
			out.AppendChild(c.CompileTree(child))
		}
	} else if node.ClosingTag == nil && !c.AutoCloseTags {
		out.Value = node.Value.(BBOpeningTag).Raw
		InsertNewlines(out)
		for _, child := range node.Children {
			out.AppendChild(c.CompileTree(child))
		}
	} else {
		in := node.Value.(BBOpeningTag)

		compileFunc, ok := c.tagCompilers[in.Name]
		if ok {
			var appendExpr bool
			out, appendExpr = compileFunc(node, in)
			if appendExpr {
				if len(node.Children) == 0 {
					out.AppendChild(NewHTMLTag(""))
				} else {
					for _, child := range node.Children {
						out.AppendChild(c.CompileTree(child))
					}
				}
			}
		} else {
			out.Value = in.Raw
			InsertNewlines(out)
			if len(node.Children) == 0 {
				out.AppendChild(NewHTMLTag(""))
			} else {
				for _, child := range node.Children {
					out.AppendChild(c.CompileTree(child))
				}
			}
			if node.ClosingTag != nil {
				tag := NewHTMLTag(node.ClosingTag.Raw)
				InsertNewlines(tag)
				out.AppendChild(tag)
			}
		}
	}
	return out
}

func CompileText(in *BBCodeNode) string {
	out := ""
	if in.ID == TEXT {
		out = in.Value.(string)
	}
	for _, child := range in.Children {
		out += CompileText(child)
	}
	return out
}

func CompileRaw(in *BBCodeNode) *HTMLTag {
	out := NewHTMLTag("")
	if in.ID == TEXT {
		out.Value = in.Value.(string)
	} else if in.ID == CLOSING_TAG {
		out.Value = in.Value.(BBClosingTag).Raw
	} else {
		out.Value = in.Value.(BBOpeningTag).Raw
	}
	InsertNewlines(out)
	for _, child := range in.Children {
		out.AppendChild(CompileRaw(child))
	}
	if in.ID == OPENING_TAG && in.ClosingTag != nil {
		tag := NewHTMLTag(in.ClosingTag.Raw)
		InsertNewlines(tag)
		out.AppendChild(tag)
	}
	return out
}

var DefaultTagCompilers map[string]TagCompilerFunc

func init() {
	DefaultTagCompilers = make(map[string]TagCompilerFunc)
	DefaultTagCompilers["url"] = func(node *BBCodeNode, in BBOpeningTag) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		out.Name = "a"
		if in.Value == "" {
			text := CompileText(node)
			if len(text) > 0 {
				out.Attrs["href"] = ValidURL(text)
			}
		} else {
			out.Attrs["href"] = ValidURL(in.Value)
		}
		return out, true
	}

	DefaultTagCompilers["img"] = func(node *BBCodeNode, in BBOpeningTag) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		out.Name = "img"
		if in.Value == "" {
			out.Attrs["src"] = ValidURL(CompileText(node))
		} else {
			out.Attrs["src"] = ValidURL(in.Value)
			text := CompileText(node)
			if len(text) > 0 {
				out.Attrs["alt"] = text
				out.Attrs["title"] = out.Attrs["alt"]
			}
		}
		return out, false
	}

	var youtubeRegex = regexp.MustCompile(`(?:https?:\/\/)?(?:www\.)?(?:youtube\.com|youtu\.be)\/(?:watch\?v=)?([a-zA-Z0-9]+)`)
	DefaultTagCompilers["media"] = func(node *BBCodeNode, in BBOpeningTag) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		mediaUrl := CompileText(node)
		out.Name = "div"
		out.Attrs["class"] = "embedded-video"

		obj := NewHTMLTag("Embedded video")
		out.AppendChild(obj)

		matches := youtubeRegex.FindStringSubmatch(mediaUrl)
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
		return out, false
	}

	DefaultTagCompilers["center"] = func(node *BBCodeNode, in BBOpeningTag) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		out.Name = "div"
		out.Attrs["style"] = "text-align: center;"
		return out, true
	}

	DefaultTagCompilers["color"] = func(node *BBCodeNode, in BBOpeningTag) (*HTMLTag, bool) {
		return NewHTMLTag(""), true
	}

	DefaultTagCompilers["size"] = func(node *BBCodeNode, in BBOpeningTag) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		out.Name = "span"
		if size, err := strconv.Atoi(in.Value); err == nil {
			out.Attrs["style"] = fmt.Sprintf("font-size: %dpx;", size*4)
		}
		return out, true
	}

	DefaultTagCompilers["spoiler"] = func(node *BBCodeNode, in BBOpeningTag) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		out.Name = "div"
		out.Attrs["class"] = "expandable collapsed"
		return out, true
	}

	DefaultTagCompilers["quote"] = func(node *BBCodeNode, in BBOpeningTag) (*HTMLTag, bool) {
		out := NewHTMLTag("")
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
		return out.AppendChild(cite), true
	}

	DefaultTagCompilers["strike"] = func(node *BBCodeNode, in BBOpeningTag) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		out.Name = "s"
		return out, true
	}

	DefaultTagCompilers["code"] = func(node *BBCodeNode, in BBOpeningTag) (*HTMLTag, bool) {
		out := NewHTMLTag("")
		out.Name = "code"
		for _, child := range node.Children {
			out.AppendChild(CompileRaw(child))
		}
		return out, false
	}

	for _, tag := range []string{"i", "b", "u"} {
		DefaultTagCompilers[tag] = func(node *BBCodeNode, in BBOpeningTag) (*HTMLTag, bool) {
			out := NewHTMLTag("")
			out.Name = in.Name
			return out, true
		}
	}
}
