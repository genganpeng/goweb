package framework

import (
	"errors"
	"strings"
)

type Tree struct {
	root *node
}

type node struct {
	segment  string              // uri中的字符串，代表这个节点表示的路由中某个段的字符串
	childs   []*node             // 代表这个节点下的子节点
	isLast   bool                // 代表这个节点是否可以成为最终的路由规则。该节点是否能成为一个独立的uri, 是否自身就是一个终极节点
	handlers []ControllerHandler // 代表这个节点中包含的控制器，用于最终加载调用
	parent *node
}

func newNode() *node {
	return &node{
		segment: "",
		childs:  []*node{},
		isLast:  false,
	}
}

func NewTree() *Tree {
	root := newNode()
	return &Tree{root: root}
}

func isWildSegment(segment string) bool {
	return strings.HasPrefix(segment, ":")
}

func (n *node) filterChildNodes(segment string) []*node {
	if len(n.childs) == 0 {
		return nil
	}
	if isWildSegment(segment) {
		return n.childs
	}

	nodes := make([]*node, 0, len(n.childs))
	for _, child := range n.childs {
		if isWildSegment(child.segment) {
			nodes = append(nodes, child)
		} else {
			if segment == child.segment {
				nodes = append(nodes, child)
			}
		}
	}

	return nodes
}

// 判断路由是否已经在节点的所有子节点树中存在了
func (n *node) matchNode(uri string) *node {
	segments := strings.SplitN(uri, "/", 2)
	segment := segments[0]
	if !isWildSegment(segment) {
		segment = strings.ToUpper(segment)
	}
	cnodes := n.filterChildNodes(segment)
	if cnodes == nil || len(cnodes) == 0 {
		return nil
	}

	// 如果只有一个segment，则是最后一个标记
	if len(segments) == 1 {
		for _, child := range cnodes {
			// 如果segment已经是最后一个节点，判断这些cnode是否有isLast标志
			if child.isLast {
				return child
			}
		}
		// 都不是最后一个节点
		return nil
	}

	for _, child := range cnodes {
		matchChild := child.matchNode(segments[1])
		if matchChild != nil {
			return matchChild
		}
	}

	return nil
}

func (t *Tree) AddRouter(uri string, handlers []ControllerHandler) error {
	n := t.root
	if n.matchNode(uri) != nil {
		return errors.New("route exist: " + uri)
	}
	segments := strings.Split(uri, "/")
	for index, segment := range segments {
		if !isWildSegment(segment) {
			segment = strings.ToUpper(segment)
		}

		isLast := index == len(segments)-1

		var objNode *node
		cnodes := n.filterChildNodes(segment)
		for _, c := range cnodes {
			if segment == c.segment {
				objNode = c
				break
			}
		}

		if objNode == nil {
			cnode := newNode()
			if isLast {
				cnode.isLast = true
				cnode.handlers = handlers
			}
			cnode.segment = segment
			cnode.parent = n
			n.childs = append(n.childs, cnode)
			objNode = cnode
		}
		n = objNode

	}
	return nil
}

func (n *node) parseParamsFromEndNode(uri string) map[string]string {
	segments := strings.Split(uri, "/")
	ret := make(map[string]string)
	cur := n
	cnt := len(segments) - 1
	for i := cnt; i >= 0; i-- {
		if cur.segment == "" {
			break
		}
		if isWildSegment(segments[i]) {
			ret[cur.segment[1:]] = segments[i]
		}

		cur = n.parent
	}

	return ret
}

func (t *Tree) FindHandler(uri string) []ControllerHandler {
	mNode := t.root.matchNode(uri)
	if mNode == nil {
		return nil
	}
	return mNode.handlers
}
