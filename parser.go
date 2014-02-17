// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

type BBCodeNode struct {
	Token
	Parent   *BBCodeNode
	Children []*BBCodeNode
}

func (n *BBCodeNode) appendChild(t Token) *BBCodeNode {
	if t.id == CLOSING_TAG {
		if n.Parent != nil && n.id == OPENING_TAG && n.value.(bbOpeningTag).name == t.value.(bbClosingTag).name {
			return n.Parent
		}
	}

	// Join consecutive TEXT tokens
	if len(n.Children) == 0 && t.id == TEXT && n.id == TEXT {
		n.value = n.value.(string) + t.value.(string)
		return n
	}

	node := &BBCodeNode{t, n, make([]*BBCodeNode, 0)}
	n.Children = append(n.Children, node)
	if t.id == OPENING_TAG {
		return node
	} else {
		return n
	}
}

func Parse(tokens chan Token) *BBCodeNode {
	var root *BBCodeNode
	var curr *BBCodeNode
	root = &BBCodeNode{Token{TEXT, ""}, nil, make([]*BBCodeNode, 0)}
	curr = root
	for tok := range tokens {
		curr = curr.appendChild(tok)
	}
	return root
}
