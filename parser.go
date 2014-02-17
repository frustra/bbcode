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
	if t.ID == CLOSING_TAG {
		if n.Parent != nil && n.ID == OPENING_TAG && n.Value.(bbOpeningTag).Name == t.Value.(bbClosingTag).Name {
			return n.Parent
		}
	}

	// Join consecutive TEXT tokens
	if len(n.Children) == 0 && t.ID == TEXT && n.ID == TEXT {
		n.Value = n.Value.(string) + t.Value.(string)
		return n
	}

	node := &BBCodeNode{t, n, make([]*BBCodeNode, 0)}
	n.Children = append(n.Children, node)
	if t.ID == OPENING_TAG {
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
