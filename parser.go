// Copyright 2014 Frustra. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package bbcode

type BBCodeNode struct {
	Token
	Parent     *BBCodeNode
	Children   []*BBCodeNode
	ClosingTag *BBClosingTag
}

func (n *BBCodeNode) appendChild(t Token) *BBCodeNode {
	if t.ID == CLOSING_TAG {
		curr := n
		closing := t.Value.(BBClosingTag)
		for curr.Parent != nil {
			if curr.ID == OPENING_TAG && curr.Value.(BBOpeningTag).Name == closing.Name {
				curr.ClosingTag = &closing
				return curr.Parent
			}
			curr = curr.Parent
		}
	}

	// Join consecutive TEXT tokens
	if len(n.Children) == 0 && t.ID == TEXT && n.ID == TEXT {
		n.Value = n.Value.(string) + t.Value.(string)
		return n
	}

	node := &BBCodeNode{t, n, make([]*BBCodeNode, 0), nil}
	n.Children = append(n.Children, node)
	if t.ID == OPENING_TAG {
		return node
	} else {
		return n
	}
}

func Parse(tokens chan Token) *BBCodeNode {
	root := &BBCodeNode{Token{TEXT, ""}, nil, make([]*BBCodeNode, 0), nil}
	curr := root
	for tok := range tokens {
		curr = curr.appendChild(tok)
	}
	return root
}
