// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// Package bbcode implements a parser and HTML generator for BBCode.
package bbcode

type BBOpeningTag struct {
	Name  string
	Value string
	Args  map[string]string
	Raw   string
}

type BBClosingTag struct {
	Name string
	Raw  string
}

func (t *BBOpeningTag) String() string {
	str := t.Name
	if len(t.Value) > 0 {
		str += "=" + t.Value
	}
	for k, v := range t.Args {
		str += " " + k
		if len(v) > 0 {
			str += "=" + v
		}
	}
	return str
}
