package main

import (
	"encoding/json"
	"fmt"
)

func main() {

	a1 := &Category{
		Name:          "1",
		Uid:           1,
		CategoryDepth: 1,
		Category1:     1,
	}

	a11 := &Category{
		Name:          "1-1",
		Uid:           2,
		CategoryDepth: 2,
		Category1:     1,
		Category2:     2,
	}

	a111 := &Category{
		Name:          "1-1-1",
		Uid:           3,
		CategoryDepth: 3,
		Category1:     1,
		Category2:     2,
		Category3:     3,
	}

	a12 := &Category{
		Name:          "1-2",
		Uid:           4,
		CategoryDepth: 2,
		Category1:     1,
		Category2:     4,
	}

	a13 := &Category{
		Name:          "1-3",
		Uid:           5,
		CategoryDepth: 2,
		Category1:     1,
		Category2:     5,
	}

	a131 := &Category{
		Name:          "1-3-1",
		Uid:           6,
		CategoryDepth: 3,
		Category1:     1,
		Category2:     5,
		Category3:     6,
	}

	a2 := &Category{
		Name:          "2",
		Uid:           7,
		CategoryDepth: 1,
		Category1:     7,
	}

	a21 := &Category{
		Name:          "2-1",
		Uid:           8,
		CategoryDepth: 2,
		Category1:     7,
		Category2:     8,
	}

	a22 := &Category{
		Name:          "2-2",
		Uid:           9,
		CategoryDepth: 2,
		Category1:     7,
		Category2:     9,
	}

	a23 := &Category{
		Name:          "2-3",
		Uid:           10,
		CategoryDepth: 2,
		Category1:     7,
		Category2:     10,
	}

	a231 := &Category{
		Name:          "2-3-1",
		Uid:           11,
		CategoryDepth: 3,
		Category1:     7,
		Category2:     10,
		Category3:     11,
	}

	a232 := &Category{
		Name:          "2-3-2",
		Uid:           12,
		CategoryDepth: 3,
		Category1:     7,
		Category2:     10,
		Category3:     12,
	}

	list := []*Category{a1, a11, a111, a12, a13, a131, a2, a21, a22, a23, a231, a232}

	tree := tree(list)

	fmt.Println("after")

	profJson, _ := json.Marshal(tree)
	fmt.Println(string(profJson))
	//	printNodes(tree)
}

type Category struct {
	Name          string `json:"name"`
	Uid           int    `json:"uid"`
	CategoryDepth int    `json:"depth"`
	Category1     int    `json:"category1"`
	Category2     int    `json:"category2"`
	Category3     int    `json:"category3"`

	// 自分の子ノードが設定される。子ノードには再帰的に孫ノードが設定される
	Childs []*Category `json:"child"`
}

func (c Category) parentNode() int {
	switch c.CategoryDepth {
	case 1:
		// depthが1の場合、親ノードが存在しないため、0を返却する
		return 0
	case 2:
		return c.Category1
	case 3:
		return c.Category2
	default:
		// depthが3までを想定する
		panic("illegal depth was set")
	}
}

func (c *Category) appendChild(child *Category) {
	c.Childs = append(c.Childs, child)
}

const MaxDepth = 3

// tree 一次元のList構造をTree構造に変更する
func tree(cs []*Category) []*Category {
	parent := map[int]*Category{}
	dummyRoot := &Category{}
	// ダミーのルートノードを設定する
	parent[0] = dummyRoot
	for i := 1; i <= MaxDepth; i++ {
		current := map[int]*Category{}
		for _, c := range cs {
			if c.CategoryDepth != i {
				continue
			}
			pUid := c.parentNode()

			if _, ok := parent[pUid]; !ok {
				fmt.Printf("parent node is not define. parents: %#v \n", pUid)
				fmt.Println(parent)
				continue
			}
			pNode := parent[pUid]
			pNode.appendChild(c)
			parent[pUid] = pNode
			current[c.Uid] = c
		}
		parent = current
	}
	return dummyRoot.Childs
}

func debugMap(m map[int]*Category) {
	for k, v := range m {
		fmt.Printf("K:%v, v:%v \n", k, v)
	}
}

func debugList(l []Category) {
	for _, v := range l {
		fmt.Printf("c:%v,\n", v)
	}
}

func printNodes(cs []*Category) {
	for _, n := range cs {
		fmt.Println(*n)
		printNodes(n.Childs)
	}
}
