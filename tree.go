package main

import (
	"container/list"
	"errors"
)

type Valuable interface {
	GetID() int
	GetParentID() (int, error)
	GetIndent() int
}

type Tree struct {
	Root *Node
}

func NewTree(que *list.List) (*Tree, error) {
	root := &Node{}
	tree := &Tree{Root: root}
	errQue := list.New()

	for que.Len() != 0 {
		queLen := que.Len()
		for e := que.Front(); e != nil; e = e.Next() {
			if err := tree.Root.Insert(e.Value.(Valuable)); err != nil {
				errQue.PushBack(e.Value)
			}
		}
		if errQue.Len() == queLen {
			return tree, errors.New("New Tree Failed")
		}
		que = errQue
		errQue = list.New()
	}
	return tree, nil
}

func (tree *Tree) Traverse() []*Node {
	return tree.Root.Traverse()
}

func (tree *Tree) SearchById(id int) *Node {
	return tree.Root.SearchById(id)
}

type Node struct {
	Parent *Node
	Child  *Node
	Next   *Node
	Value  Valuable
}

func NewNode(val Valuable, parent *Node) *Node {
	return &Node{Value: val, Parent: parent}
}

func (node *Node) Traverse() []*Node {
	ret := []*Node{node}
	if node.Child != nil {
		ret = append(ret, node.Child.Traverse()...)
	}
	if node.Next != nil {
		ret = append(ret, node.Next.Traverse()...)
	}
	return ret
}

func (node *Node) SearchById(id int) *Node {
	if node.Value.GetID() == id {
		return node
	}
	if node.Child != nil {
		if res := node.Child.SearchById(id); res != nil {
			return res
		}
	}
	if node.Next != nil {
		if res := node.Next.SearchById(id); res != nil {
			return res
		}
	}
	return nil
}

func (node *Node) Parents() (ret []*Node) {
	if node.Parent == nil {
		return
	}
	return append(node.Parent.Parents(), node.Parent)
}

func (node *Node) Insert(val Valuable) (err error) {
	if node.Value == nil {
		node.Value = val
		return nil
	}

	valParentId, _ := val.GetParentID()

	if node.Child != nil {
		if err = node.Child.Insert(val); err == nil {
			return nil
		}
	}
	if node.Child == nil {
		if valParentId == node.Value.GetID() {
			node.Child = NewNode(val, node)
			return nil
		}
	}

	if node.Next != nil {
		if err = node.Next.Insert(val); err == nil {
			return nil
		}
	}
	if node.Next == nil {
		nodeParentId, _ := node.Value.GetParentID()
		if nodeParentId == valParentId {
			node.Next = NewNode(val, node.Parent)
			return nil
		}
	}

	return errors.New("No insert position")
}
