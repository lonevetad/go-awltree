package avltree

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

type errorTest struct {
	errText string
}

func (et *errorTest) Error() string {
	if et == nil {
		return "WAIT, THIS ERROR IS FRICKING NULL!"
	}
	return et.errText
}
func (et *errorTest) String() string {
	if et == nil {
		return "errorTest is nil"
	}
	return et.errText
}
func (et *errorTest) IsNil() bool {
	return et == nil
}
func (et *errorTest) IsEmpty() bool {
	return et == nil || et.errText == ""
}

func ne(s string) *errorTest {
	e := new(errorTest)
	(*e).errText = s
	return e
}

type TestData struct {
	Id   int
	Text string
}

func Extract(t *TestData) int {
	if t == nil {
		return -1
	}
	return t.Id
}

func IntCompare(i1 int, i2 int) int64 {
	if i1 > i2 {
		return 1
	}
	if i1 < i2 {
		return -1
	}
	return 0
}

func (td *TestData) String() string {
	if td == nil {
		return "null"
	}
	var sb strings.Builder
	sb.WriteString("<id= ")
	sb.WriteString(strconv.Itoa(td.Id))
	sb.WriteString("; text= \"")
	sb.WriteString(td.Text)
	sb.WriteString("\">")
	return sb.String()
}

func NewMetadata(keyZeroValue int, valueZeroValue *TestData) AVLTreeConstructorParams[int, *TestData] {
	avlTreeConstructorParams := AVLTreeConstructorParams[int, *TestData]{}
	avlTreeConstructorParams.KeyCollisionBehavior = Replace
	avlTreeConstructorParams.KeyZeroValue = keyZeroValue
	avlTreeConstructorParams.ValueZeroValue = valueZeroValue
	avlTreeConstructorParams.KeyExtractor = Extract
	avlTreeConstructorParams.Comparator = IntCompare
	return avlTreeConstructorParams
}

func NewTestData() *TestData {
	td := new(TestData)
	td.Id = -42
	td.Text = "HELLO NULL STRING"
	return td
}
func NewTestDataFilled(k int, v string) *TestData {
	td := new(TestData)
	td.Id = k
	td.Text = v
	return td
}

func NewTreeNodeFilled(tree *AVLTree[int, *TestData], key int, v string) *AVLTNode[int, *TestData] {
	return tree.newNode(key, NewTestDataFilled(key, v))
}

func NewTree() (*AVLTree[int, *TestData], error) {
	//td := NewTestData()
	avlTreeConstructorParams := NewMetadata(-1000, nil)
	return NewAVLTree(avlTreeConstructorParams)
}

//
// TESTS
//

func TestNewTree(t *testing.T) {
	tree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}

	if tree == nil {
		t.Fatal("the new tree should not be nil\n")
		return
	}

	if tree.root == nil {
		t.Fatalf("the tree's root should NOT be nil\n")
	}

	err = testEqualityPrimitive(true, tree.Size(), 0, "size should be 0")
	if err != nil {
		t.Fatal(err)
	}
	if tree._NIL == nil {
		t.Fatalf("the tree's \"_NIL\" should NOT be nil\n")
	}
	if tree.root != tree._NIL {
		t.Fatalf("the tree should NOT have a root AND should be \"_NIL\"\n")
	}
	err = testIsLeaf(tree, tree._NIL)
	if err != nil {
		t.Fatal(err)
	}
	err = testNIL(tree, true, tree.root, "root is not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testNIL(tree, true, tree.root.father, "father is not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree._NIL.height, DEPTH_INITIAL, fmt.Sprintf("NIL's height should be: %d", DEPTH_INITIAL))
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree._NIL.sizeLeft, 0, "NIL's sizeLeft should be 0")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree._NIL.sizeRight, 0, "NIL's sizeRight should be 0")
	if err != nil {
		t.Fatal(err)
	}
	err = testNIL(tree, true, tree.minValue, "minValue is not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testNIL(tree, true, tree.firstInserted, "firstInserted is not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testNIL(tree, true, tree.root.nextInOrder, "nextInOrder is not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testNIL(tree, true, tree.root.prevInOrder, "prevInOrder is not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testNIL(tree, true, tree.root.nextInserted, "nextInserted is not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testNIL(tree, true, tree.root.prevInserted, "prevInserted is not _NIL")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRotateLeftLeft_3nodes(t *testing.T) {
	tree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	dummyTree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	values := []int{
		//1, 2, 3,
		3, 2, 1,
	}
	nodesTree := []*AVLTNode[int, *TestData]{nil, nil, nil}
	nodesDummyTree := []*AVLTNode[int, *TestData]{nil, nil, nil}
	for i := 0; i < len(values); i++ {
		nodesTree[i] = NewTreeNodeFilled(tree, values[i], fmt.Sprintf("v_%d", values[i]))
		nodesDummyTree[i] = NewTreeNodeFilled(dummyTree, values[i], fmt.Sprintf("v_%d", values[i]))
	}
	for _, nt := range [][]*AVLTNode[int, *TestData]{nodesTree, nodesDummyTree} {
		for i := 0; i < len(values); i++ {
			linkNodes(nt[i], nt[(i+1)%len(values)], false)
		}
	}

	// seting up trees - original
	tree.root = nodesTree[0]
	nodesTree[0].left = nodesTree[1]
	nodesTree[1].father = nodesTree[0]
	nodesTree[1].left = nodesTree[2]
	nodesTree[2].father = nodesTree[1]
	tree.firstInserted = nodesTree[0]
	nodesTree[0].height = 2
	nodesTree[0].sizeLeft = 2
	nodesTree[0].sizeRight = 0
	nodesTree[1].height = 1
	nodesTree[1].sizeLeft = 1
	nodesTree[1].sizeRight = 0
	nodesTree[2].height = 0
	nodesTree[2].sizeLeft = 0
	nodesTree[2].sizeRight = 0
	tree.minValue = nodesTree[2]
	tree.size = 3
	linkNodes(nodesTree[2], nodesTree[1], true)
	linkNodes(nodesTree[1], nodesTree[0], true)
	linkNodes(nodesTree[0], nodesTree[2], true)

	// seting up trees - dummy
	dummyTree.root = nodesDummyTree[1]
	nodesDummyTree[1].left = nodesDummyTree[2]
	nodesDummyTree[1].right = nodesDummyTree[0]
	nodesDummyTree[2].father = nodesDummyTree[1]
	nodesDummyTree[0].father = nodesDummyTree[1]
	dummyTree.firstInserted = nodesDummyTree[0]
	nodesDummyTree[0].height = 0
	nodesDummyTree[0].sizeLeft = 0
	nodesDummyTree[0].sizeRight = 0
	nodesDummyTree[1].height = 1
	nodesDummyTree[1].sizeLeft = 1
	nodesDummyTree[1].sizeRight = 1
	nodesDummyTree[2].height = 0
	nodesDummyTree[2].sizeLeft = 0
	nodesDummyTree[2].sizeRight = 0
	dummyTree.minValue = nodesDummyTree[2]
	dummyTree.size = 3
	linkNodes(nodesDummyTree[2], nodesDummyTree[1], true)
	linkNodes(nodesDummyTree[1], nodesDummyTree[0], true)
	linkNodes(nodesDummyTree[0], nodesDummyTree[2], true)

	// rotating
	tree.insertFixup(nodesTree[1])
	tree.cleanNil()

	// checking

	expectSize := int64(len(values))
	err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateLeftLeft, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", tree.size, expectSize, err)
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, dummyTree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateLeftLeft, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", dummyTree.size, expectSize, err)
		t.Fatal(err)
	}

	areEquals, errText := CheckTrees(tree, dummyTree)
	if errText != nil {
		t.Fatal(errText.Error())
		return
	}
	if !areEquals {
		t.Fatal(fmt.Errorf("trees are not equal"))
		return
	}
}

func TestRotateRightRight_3nodes(t *testing.T) {
	tree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	dummyTree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	values := []int{
		1, 2, 3,
	}
	nodesTree := []*AVLTNode[int, *TestData]{nil, nil, nil}
	nodesDummyTree := []*AVLTNode[int, *TestData]{nil, nil, nil}
	for i := 0; i < len(values); i++ {
		nodesTree[i] = NewTreeNodeFilled(tree, values[i], fmt.Sprintf("v_%d", values[i]))
		nodesDummyTree[i] = NewTreeNodeFilled(dummyTree, values[i], fmt.Sprintf("v_%d", values[i]))
	}
	for _, nt := range [][]*AVLTNode[int, *TestData]{nodesTree, nodesDummyTree} {
		for i := 0; i < len(values); i++ {
			linkNodes(nt[i], nt[(i+1)%len(values)], false)
		}
	}

	// seting up trees - original
	tree.root = nodesTree[0]
	nodesTree[0].right = nodesTree[1]
	nodesTree[1].father = nodesTree[0]
	nodesTree[1].right = nodesTree[2]
	nodesTree[2].father = nodesTree[1]
	tree.firstInserted = nodesTree[0]
	nodesTree[0].height = 2
	nodesTree[0].sizeLeft = 0
	nodesTree[0].sizeRight = 2
	nodesTree[1].height = 1
	nodesTree[1].sizeLeft = 0
	nodesTree[1].sizeRight = 1
	nodesTree[2].height = 0
	nodesTree[2].sizeLeft = 0
	nodesTree[2].sizeRight = 0
	tree.minValue = nodesTree[0]
	tree.size = 3
	linkNodes(nodesTree[0], nodesTree[1], true)
	linkNodes(nodesTree[1], nodesTree[2], true)
	linkNodes(nodesTree[2], nodesTree[0], true)

	// seting up trees - dummy
	dummyTree.root = nodesDummyTree[1]
	nodesDummyTree[1].left = nodesDummyTree[0]
	nodesDummyTree[1].right = nodesDummyTree[2]
	nodesDummyTree[2].father = nodesDummyTree[1]
	nodesDummyTree[0].father = nodesDummyTree[1]
	dummyTree.firstInserted = nodesDummyTree[0]
	nodesDummyTree[0].height = 0
	nodesDummyTree[0].sizeLeft = 0
	nodesDummyTree[0].sizeRight = 0
	nodesDummyTree[1].height = 1
	nodesDummyTree[1].sizeLeft = 1
	nodesDummyTree[1].sizeRight = 1
	nodesDummyTree[2].height = 0
	nodesDummyTree[2].sizeLeft = 0
	nodesDummyTree[2].sizeRight = 0
	dummyTree.minValue = nodesDummyTree[0]
	dummyTree.size = 3
	linkNodes(nodesDummyTree[0], nodesDummyTree[1], true)
	linkNodes(nodesDummyTree[1], nodesDummyTree[2], true)
	linkNodes(nodesDummyTree[2], nodesDummyTree[0], true)

	// rotating
	tree.insertFixup(nodesTree[1])
	tree.cleanNil()

	// checking

	expectSize := int64(len(values))
	err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateRightRight, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", tree.size, expectSize, err)
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, dummyTree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateRightRight, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", dummyTree.size, expectSize, err)
		t.Fatal(err)
	}

	areEquals, errText := CheckTrees(tree, dummyTree)
	if errText != nil {
		t.Fatal(errText)
		return
	}
	if !areEquals {
		t.Fatal(fmt.Errorf("trees are not equal"))
		return
	}
}

func TestRotateLeftRight_3nodes(t *testing.T) {
	tree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	dummyTree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	values := []int{
		3, 1, 2,
	}
	nodesTree := []*AVLTNode[int, *TestData]{nil, nil, nil}
	nodesDummyTree := []*AVLTNode[int, *TestData]{nil, nil, nil}
	for i := 0; i < len(values); i++ {
		nodesTree[i] = NewTreeNodeFilled(tree, values[i], fmt.Sprintf("v_%d", values[i]))
		nodesDummyTree[i] = NewTreeNodeFilled(dummyTree, values[i], fmt.Sprintf("v_%d", values[i]))
	}
	for _, nt := range [][]*AVLTNode[int, *TestData]{nodesTree, nodesDummyTree} {
		for i := 0; i < len(values); i++ {
			linkNodes(nt[i], nt[(i+1)%len(values)], false)
		}
	}

	//	  3
	//	/
	// 1
	//  \
	//   2

	// seting up trees - original
	tree.root = nodesTree[0]
	nodesTree[0].left = nodesTree[1]
	nodesTree[1].father = nodesTree[0]
	nodesTree[1].right = nodesTree[2]
	nodesTree[2].father = nodesTree[1]
	tree.firstInserted = nodesTree[0]
	nodesTree[0].height = 2
	nodesTree[0].sizeLeft = 2
	nodesTree[0].sizeRight = 0
	nodesTree[1].height = 1
	nodesTree[1].sizeLeft = 1
	nodesTree[1].sizeRight = 0
	nodesTree[2].height = 0
	nodesTree[2].sizeLeft = 0
	nodesTree[2].sizeRight = 0
	tree.minValue = nodesTree[2]
	tree.size = 3
	linkNodes(nodesTree[1], nodesTree[2], true)
	linkNodes(nodesTree[2], nodesTree[0], true)
	linkNodes(nodesTree[0], nodesTree[1], true)

	// seting up trees - dummy
	dummyTree.root = nodesDummyTree[2]
	nodesDummyTree[2].left = nodesDummyTree[1]
	nodesDummyTree[2].right = nodesDummyTree[0]
	nodesDummyTree[1].father = nodesDummyTree[2]
	nodesDummyTree[0].father = nodesDummyTree[2]
	dummyTree.firstInserted = nodesDummyTree[0]
	nodesDummyTree[0].height = 0
	nodesDummyTree[0].sizeLeft = 0
	nodesDummyTree[0].sizeRight = 0
	nodesDummyTree[1].height = 0
	nodesDummyTree[1].sizeLeft = 0
	nodesDummyTree[1].sizeRight = 0
	nodesDummyTree[2].height = 1
	nodesDummyTree[2].sizeLeft = 1
	nodesDummyTree[2].sizeRight = 1
	dummyTree.minValue = nodesDummyTree[1]
	dummyTree.size = 3
	linkNodes(nodesDummyTree[1], nodesDummyTree[2], true)
	linkNodes(nodesDummyTree[2], nodesDummyTree[0], true)
	linkNodes(nodesDummyTree[0], nodesDummyTree[1], true)

	// rotating
	tree.insertFixup(nodesTree[1])
	tree.cleanNil()

	// checking

	expectSize := int64(len(values))
	err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateLeftRight, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", tree.size, expectSize, err)
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, dummyTree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateLeftRight, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", dummyTree.size, expectSize, err)
		t.Fatal(err)
	}

	areEquals, errText := CheckTrees(tree, dummyTree)
	if errText != nil {
		t.Fatal(errText.Error())
		return
	}
	if !areEquals {
		t.Fatal(fmt.Errorf("trees are not equal"))
		return
	}
}

func TestRotateRightLeft_3nodes(t *testing.T) {
	tree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	dummyTree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	values := []int{
		1, 3, 2,
	}
	nodesTree := []*AVLTNode[int, *TestData]{nil, nil, nil}
	nodesDummyTree := []*AVLTNode[int, *TestData]{nil, nil, nil}
	for i := 0; i < len(values); i++ {
		nodesTree[i] = NewTreeNodeFilled(tree, values[i], fmt.Sprintf("v_%d", values[i]))
		nodesDummyTree[i] = NewTreeNodeFilled(dummyTree, values[i], fmt.Sprintf("v_%d", values[i]))
	}
	for _, nt := range [][]*AVLTNode[int, *TestData]{nodesTree, nodesDummyTree} {
		for i := 0; i < len(values); i++ {
			linkNodes(nt[i], nt[(i+1)%len(values)], false)
		}
	}

	// seting up trees - original
	tree.root = nodesTree[0]
	nodesTree[0].right = nodesTree[1]
	nodesTree[1].father = nodesTree[0]
	nodesTree[1].left = nodesTree[2]
	nodesTree[2].father = nodesTree[1]
	tree.firstInserted = nodesTree[0]
	nodesTree[0].height = 2
	nodesTree[0].sizeLeft = 0
	nodesTree[0].sizeRight = 2
	nodesTree[1].height = 1
	nodesTree[1].sizeLeft = 1
	nodesTree[1].sizeRight = 0
	nodesTree[2].height = 0
	nodesTree[2].sizeLeft = 0
	nodesTree[2].sizeRight = 0
	tree.minValue = nodesTree[0]
	tree.size = 3
	linkNodes(nodesTree[0], nodesTree[2], true)
	linkNodes(nodesTree[2], nodesTree[1], true)
	linkNodes(nodesTree[1], nodesTree[0], true)

	// seting up trees - dummy
	dummyTree.root = nodesDummyTree[2]
	nodesDummyTree[2].left = nodesDummyTree[0]
	nodesDummyTree[2].right = nodesDummyTree[1]
	nodesDummyTree[0].father = nodesDummyTree[2]
	nodesDummyTree[1].father = nodesDummyTree[2]
	dummyTree.firstInserted = nodesDummyTree[0]
	nodesDummyTree[0].height = 0
	nodesDummyTree[0].sizeLeft = 0
	nodesDummyTree[0].sizeRight = 0
	nodesDummyTree[1].height = 0
	nodesDummyTree[1].sizeLeft = 0
	nodesDummyTree[1].sizeRight = 0
	nodesDummyTree[2].height = 1
	nodesDummyTree[2].sizeLeft = 1
	nodesDummyTree[2].sizeRight = 1
	dummyTree.minValue = nodesDummyTree[0]
	dummyTree.size = 3
	linkNodes(nodesDummyTree[0], nodesDummyTree[2], true)
	linkNodes(nodesDummyTree[2], nodesDummyTree[1], true)
	linkNodes(nodesDummyTree[1], nodesDummyTree[0], true)

	// rotating
	tree.insertFixup(nodesTree[1])
	tree.cleanNil()

	// checking

	expectSize := int64(len(values))
	err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateRightLeft, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", tree.size, expectSize, err)
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, dummyTree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateRightLeft, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", dummyTree.size, expectSize, err)
		t.Fatal(err)
	}

	areEquals, errText := CheckTrees(tree, dummyTree)
	if errText != nil {
		t.Fatal(errText)
		return
	}
	if !areEquals {
		t.Fatal(fmt.Errorf("trees are not equal"))
		return
	}
}

func TestRotateLeftLeft_5nodes(t *testing.T) {
	tree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	dummyTree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	values := []int{
		3, 2, 1, 4, 0,
	}
	nodesTree := []*AVLTNode[int, *TestData]{nil, nil, nil, nil, nil}
	nodesDummyTree := []*AVLTNode[int, *TestData]{nil, nil, nil, nil, nil}
	for i := 0; i < len(values); i++ {
		nodesTree[i] = NewTreeNodeFilled(tree, values[i], fmt.Sprintf("v_%d", values[i]))
		nodesDummyTree[i] = NewTreeNodeFilled(dummyTree, values[i], fmt.Sprintf("v_%d", values[i]))
	}
	for _, nt := range [][]*AVLTNode[int, *TestData]{nodesTree, nodesDummyTree} {
		for i := 0; i < len(values); i++ {
			linkNodes(nt[i], nt[(i+1)%len(values)], false)
		}
	}

	// seting up trees - original
	tree.root = nodesTree[0]
	nodesTree[0].left = nodesTree[1]
	nodesTree[1].father = nodesTree[0]
	nodesTree[1].left = nodesTree[2]
	nodesTree[2].father = nodesTree[1]
	tree.root.right = nodesTree[3]
	nodesTree[3].father = nodesTree[0]
	nodesTree[2].left = nodesTree[4]
	nodesTree[4].father = nodesTree[2]
	tree.firstInserted = nodesTree[0]
	nodesTree[0].height = 3
	nodesTree[0].sizeLeft = 3
	nodesTree[0].sizeRight = 1
	nodesTree[1].height = 2
	nodesTree[1].sizeLeft = 2
	nodesTree[1].sizeRight = 0
	nodesTree[2].height = 1
	nodesTree[2].sizeLeft = 1
	nodesTree[2].sizeRight = 0
	nodesTree[3].height = 0
	nodesTree[3].sizeLeft = 0
	nodesTree[3].sizeRight = 0
	nodesTree[4].height = 0
	nodesTree[4].sizeLeft = 0
	nodesTree[4].sizeRight = 0
	tree.minValue = nodesTree[4]
	tree.size = 5
	linkNodes(nodesTree[4], nodesTree[2], true)
	linkNodes(nodesTree[2], nodesTree[1], true)
	linkNodes(nodesTree[1], nodesTree[0], true)
	linkNodes(nodesTree[0], nodesTree[3], true)
	linkNodes(nodesTree[3], nodesTree[4], true)

	// seting up trees - dummy
	dummyTree.root = nodesDummyTree[0]
	nodesDummyTree[0].left = nodesDummyTree[2]
	nodesDummyTree[0].right = nodesDummyTree[3]
	nodesDummyTree[1].father = nodesDummyTree[2]
	nodesDummyTree[2].father = nodesDummyTree[0]
	nodesDummyTree[2].left = nodesDummyTree[4]
	nodesDummyTree[2].right = nodesDummyTree[1]
	nodesDummyTree[3].father = nodesDummyTree[0]
	nodesDummyTree[4].father = nodesDummyTree[2]
	dummyTree.firstInserted = nodesDummyTree[0]
	nodesDummyTree[0].height = 2
	nodesDummyTree[0].sizeLeft = 3
	nodesDummyTree[0].sizeRight = 1
	nodesDummyTree[1].height = 0
	nodesDummyTree[1].sizeLeft = 0
	nodesDummyTree[1].sizeRight = 0
	nodesDummyTree[2].height = 1
	nodesDummyTree[2].sizeLeft = 1
	nodesDummyTree[2].sizeRight = 1
	nodesDummyTree[3].height = 0
	nodesDummyTree[3].sizeLeft = 0
	nodesDummyTree[3].sizeRight = 0
	nodesDummyTree[4].height = 0
	nodesDummyTree[4].sizeLeft = 0
	nodesDummyTree[4].sizeRight = 0
	dummyTree.minValue = nodesDummyTree[4]
	dummyTree.size = 5
	linkNodes(nodesDummyTree[4], nodesDummyTree[2], true)
	linkNodes(nodesDummyTree[2], nodesDummyTree[1], true)
	linkNodes(nodesDummyTree[1], nodesDummyTree[0], true)
	linkNodes(nodesDummyTree[0], nodesDummyTree[3], true)
	linkNodes(nodesDummyTree[3], nodesDummyTree[4], true)

	// rotating
	tree.insertFixup(nodesTree[2])
	tree.cleanNil()

	// checking

	expectSize := int64(len(values))
	err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateLeftLeft, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", tree.size, expectSize, err)
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, dummyTree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateLeftLeft, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", dummyTree.size, expectSize, err)
		t.Fatal(err)
	}

	areEquals, errText := CheckTrees(tree, dummyTree)
	if errText != nil {
		t.Fatalf("%s", errText.Error())
		return
	}
	if !areEquals {
		t.Fatal(fmt.Errorf("trees are not equal"))
		return
	}
}

func TestRotateRightRight_5nodes(t *testing.T) {
	tree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	dummyTree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	values := []int{
		1, 2, 0, 3, 4,
	}
	nodesTree := []*AVLTNode[int, *TestData]{nil, nil, nil, nil, nil}
	nodesDummyTree := []*AVLTNode[int, *TestData]{nil, nil, nil, nil, nil}
	for i := 0; i < len(values); i++ {
		nodesTree[i] = NewTreeNodeFilled(tree, values[i], fmt.Sprintf("v_%d", values[i]))
		nodesDummyTree[i] = NewTreeNodeFilled(dummyTree, values[i], fmt.Sprintf("v_%d", values[i]))
	}
	for _, nt := range [][]*AVLTNode[int, *TestData]{nodesTree, nodesDummyTree} {
		for i := 0; i < len(values); i++ {
			linkNodes(nt[i], nt[(i+1)%len(values)], false)
		}
	}

	// seting up trees - original
	tree.root = nodesTree[0]
	nodesTree[0].right = nodesTree[1]
	nodesTree[0].left = nodesTree[2]
	nodesTree[1].father = nodesTree[0]
	nodesTree[1].right = nodesTree[3]
	nodesTree[2].father = nodesTree[0]
	nodesTree[3].right = nodesTree[4]
	nodesTree[3].father = nodesTree[1]
	nodesTree[4].father = nodesTree[3]
	tree.firstInserted = nodesTree[0]
	nodesTree[0].height = 3
	nodesTree[0].sizeLeft = 1
	nodesTree[0].sizeRight = 3
	nodesTree[1].height = 2
	nodesTree[1].sizeLeft = 0
	nodesTree[1].sizeRight = 2
	nodesTree[2].height = 0
	nodesTree[2].sizeLeft = 0
	nodesTree[2].sizeRight = 0
	nodesTree[3].height = 1
	nodesTree[3].sizeLeft = 0
	nodesTree[3].sizeRight = 1
	nodesTree[4].height = 0
	nodesTree[4].sizeLeft = 0
	nodesTree[4].sizeRight = 0
	tree.minValue = nodesTree[2]
	tree.size = 5
	linkNodes(nodesTree[2], nodesTree[0], true)
	linkNodes(nodesTree[0], nodesTree[1], true)
	linkNodes(nodesTree[1], nodesTree[3], true)
	linkNodes(nodesTree[3], nodesTree[4], true)
	linkNodes(nodesTree[4], nodesTree[2], true)

	// seting up trees - dummy
	dummyTree.root = nodesDummyTree[0]
	nodesDummyTree[0].left = nodesDummyTree[2]
	nodesDummyTree[0].right = nodesDummyTree[3]
	nodesDummyTree[1].father = nodesDummyTree[3]
	nodesDummyTree[2].father = nodesDummyTree[0]
	nodesDummyTree[3].father = nodesDummyTree[0]
	nodesDummyTree[3].left = nodesDummyTree[1]
	nodesDummyTree[3].right = nodesDummyTree[4]
	nodesDummyTree[4].father = nodesDummyTree[3]
	dummyTree.firstInserted = nodesDummyTree[0]
	nodesDummyTree[0].height = 2
	nodesDummyTree[0].sizeLeft = 1
	nodesDummyTree[0].sizeRight = 3
	nodesDummyTree[1].height = 0
	nodesDummyTree[1].sizeLeft = 0
	nodesDummyTree[1].sizeRight = 0
	nodesDummyTree[2].height = 0
	nodesDummyTree[2].sizeLeft = 0
	nodesDummyTree[2].sizeRight = 0
	nodesDummyTree[3].height = 1
	nodesDummyTree[3].sizeLeft = 1
	nodesDummyTree[3].sizeRight = 1
	nodesDummyTree[4].height = 0
	nodesDummyTree[4].sizeLeft = 0
	nodesDummyTree[4].sizeRight = 0
	dummyTree.minValue = nodesDummyTree[2]
	dummyTree.size = 5
	linkNodes(nodesDummyTree[2], nodesDummyTree[0], true)
	linkNodes(nodesDummyTree[0], nodesDummyTree[1], true)
	linkNodes(nodesDummyTree[1], nodesDummyTree[3], true)
	linkNodes(nodesDummyTree[3], nodesDummyTree[4], true)
	linkNodes(nodesDummyTree[4], nodesDummyTree[2], true)

	// rotating
	tree.insertFixup(nodesTree[3])
	tree.cleanNil()

	// checking

	expectSize := int64(len(values))
	err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateRightRight, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", tree.size, expectSize, err)
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, dummyTree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateRightRight, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", dummyTree.size, expectSize, err)
		t.Fatal(err)
	}

	areEquals, errText := CheckTrees(tree, dummyTree)
	if errText != nil {
		t.Fatal(errText)
		return
	}
	if !areEquals {
		t.Fatal(fmt.Errorf("trees are not equal"))
		return
	}
}

func TestRotateLeftRight_5nodes(t *testing.T) {
	tree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	dummyTree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	values := []int{
		3, 4, 2, 0, 1,
	}
	nodesTree := []*AVLTNode[int, *TestData]{nil, nil, nil, nil, nil}
	nodesDummyTree := []*AVLTNode[int, *TestData]{nil, nil, nil, nil, nil}
	for i := 0; i < len(values); i++ {
		nodesTree[i] = NewTreeNodeFilled(tree, values[i], fmt.Sprintf("v_%d", values[i]))
		nodesDummyTree[i] = NewTreeNodeFilled(dummyTree, values[i], fmt.Sprintf("v_%d", values[i]))
	}
	for _, nt := range [][]*AVLTNode[int, *TestData]{nodesTree, nodesDummyTree} {
		for i := 0; i < len(values); i++ {
			linkNodes(nt[i], nt[(i+1)%len(values)], false)
		}
	}

	//			  3
	//		  /	   \
	//	   2			 4
	//	/
	// 0
	//  \
	//   1
	// ->
	//		 3
	//	  /	 \
	//   1		   4
	//  / \
	// 0   2

	// seting up trees - original
	tree.root = nodesTree[0]
	nodesTree[0].left = nodesTree[2]
	nodesTree[0].right = nodesTree[1]
	nodesTree[1].father = nodesTree[0]
	nodesTree[2].father = nodesTree[0]
	nodesTree[2].left = nodesTree[3]
	nodesTree[3].father = nodesTree[2]
	nodesTree[3].right = nodesTree[4]
	nodesTree[4].father = nodesTree[3]
	tree.firstInserted = nodesTree[0]
	nodesTree[0].height = 3
	nodesTree[0].sizeLeft = 3
	nodesTree[0].sizeRight = 1
	nodesTree[1].height = 0
	nodesTree[1].sizeLeft = 0
	nodesTree[1].sizeRight = 0
	nodesTree[2].height = 2
	nodesTree[2].sizeLeft = 2
	nodesTree[2].sizeRight = 0
	nodesTree[3].height = 1
	nodesTree[3].sizeLeft = 0
	nodesTree[3].sizeRight = 1
	nodesTree[4].height = 0
	nodesTree[4].sizeLeft = 0
	nodesTree[4].sizeRight = 0
	tree.minValue = nodesTree[3]
	tree.size = 5
	linkNodes(nodesTree[3], nodesTree[4], true)
	linkNodes(nodesTree[4], nodesTree[2], true)
	linkNodes(nodesTree[2], nodesTree[0], true)
	linkNodes(nodesTree[0], nodesTree[1], true)
	linkNodes(nodesTree[1], nodesTree[3], true)

	// seting up trees - dummy
	dummyTree.root = nodesDummyTree[0]
	nodesDummyTree[0].left = nodesDummyTree[4]
	nodesDummyTree[0].right = nodesDummyTree[1]
	nodesDummyTree[1].father = nodesDummyTree[0]
	nodesDummyTree[2].father = nodesDummyTree[4]
	nodesDummyTree[3].father = nodesDummyTree[4]
	nodesDummyTree[4].father = nodesDummyTree[0]
	nodesDummyTree[4].left = nodesDummyTree[3]
	nodesDummyTree[4].right = nodesDummyTree[2]
	dummyTree.firstInserted = nodesDummyTree[0]
	nodesDummyTree[0].height = 2
	nodesDummyTree[0].sizeLeft = 3
	nodesDummyTree[0].sizeRight = 1
	nodesDummyTree[1].height = 0
	nodesDummyTree[1].sizeLeft = 0
	nodesDummyTree[1].sizeRight = 0
	nodesDummyTree[2].height = 0
	nodesDummyTree[2].sizeLeft = 0
	nodesDummyTree[2].sizeRight = 0
	nodesDummyTree[3].height = 0
	nodesDummyTree[3].sizeLeft = 0
	nodesDummyTree[3].sizeRight = 0
	nodesDummyTree[4].height = 1
	nodesDummyTree[4].sizeLeft = 1
	nodesDummyTree[4].sizeRight = 1
	dummyTree.minValue = nodesDummyTree[3]
	dummyTree.size = 5
	linkNodes(nodesDummyTree[3], nodesDummyTree[4], true)
	linkNodes(nodesDummyTree[4], nodesDummyTree[2], true)
	linkNodes(nodesDummyTree[2], nodesDummyTree[0], true)
	linkNodes(nodesDummyTree[0], nodesDummyTree[1], true)
	linkNodes(nodesDummyTree[1], nodesDummyTree[3], true)

	// rotating
	tree.insertFixup(nodesTree[3])
	tree.cleanNil()

	// checking

	expectSize := int64(len(values))
	err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateLeftRight, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", tree.size, expectSize, err)
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, dummyTree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateLeftRight, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", dummyTree.size, expectSize, err)
		t.Fatal(err)
	}

	areEquals, errText := CheckTrees(tree, dummyTree)
	if errText != nil {
		t.Fatal(errText)
		return
	}
	if !areEquals {
		t.Fatal(fmt.Errorf("trees are not equal"))
		return
	}
}

func TestRotateRightLeft_5nodes(t *testing.T) {
	tree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	dummyTree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	values := []int{
		3, 0, 4, 2, 1,
	}
	nodesTree := []*AVLTNode[int, *TestData]{nil, nil, nil, nil, nil}
	nodesDummyTree := []*AVLTNode[int, *TestData]{nil, nil, nil, nil, nil}
	for i := 0; i < len(values); i++ {
		nodesTree[i] = NewTreeNodeFilled(tree, values[i], fmt.Sprintf("v_%d", values[i]))
		nodesDummyTree[i] = NewTreeNodeFilled(dummyTree, values[i], fmt.Sprintf("v_%d", values[i]))
	}
	for _, nt := range [][]*AVLTNode[int, *TestData]{nodesTree, nodesDummyTree} {
		for i := 0; i < len(values); i++ {
			linkNodes(nt[i], nt[(i+1)%len(values)], false)
		}
	}

	//		3
	//	/	   \
	// 0			 4
	//   \
	//	 2
	//	/
	//   1
	// ->
	//		 3
	//	  /	 \
	//   1		   4
	//  / \
	// 0   2

	// seting up trees - original
	tree.root = nodesTree[0]
	nodesTree[0].left = nodesTree[1]
	nodesTree[0].right = nodesTree[2]
	nodesTree[1].father = nodesTree[0]
	nodesTree[1].right = nodesTree[3]
	nodesTree[2].father = nodesTree[0]
	nodesTree[3].father = nodesTree[1]
	nodesTree[3].left = nodesTree[4]
	nodesTree[4].father = nodesTree[3]
	tree.firstInserted = nodesTree[0]
	nodesTree[0].height = 3
	nodesTree[0].sizeLeft = 3
	nodesTree[0].sizeRight = 1
	nodesTree[1].height = 2
	nodesTree[1].sizeLeft = 0
	nodesTree[1].sizeRight = 2
	nodesTree[2].height = 0
	nodesTree[2].sizeLeft = 0
	nodesTree[2].sizeRight = 0
	nodesTree[3].height = 1
	nodesTree[3].sizeLeft = 1
	nodesTree[3].sizeRight = 0
	nodesTree[4].height = 0
	nodesTree[4].sizeLeft = 0
	nodesTree[4].sizeRight = 0
	tree.minValue = nodesTree[1]
	tree.size = 5
	linkNodes(nodesTree[1], nodesTree[4], true)
	linkNodes(nodesTree[4], nodesTree[3], true)
	linkNodes(nodesTree[3], nodesTree[0], true)
	linkNodes(nodesTree[0], nodesTree[2], true)
	linkNodes(nodesTree[2], nodesTree[1], true)

	// seting up trees - dummy
	dummyTree.root = nodesDummyTree[0]
	nodesDummyTree[0].left = nodesDummyTree[4]
	nodesDummyTree[0].right = nodesDummyTree[2]
	nodesDummyTree[1].father = nodesDummyTree[4]
	nodesDummyTree[2].father = nodesDummyTree[0]
	nodesDummyTree[3].father = nodesDummyTree[4]
	nodesDummyTree[4].father = nodesDummyTree[0]
	nodesDummyTree[4].left = nodesDummyTree[1]
	nodesDummyTree[4].right = nodesDummyTree[3]
	dummyTree.firstInserted = nodesDummyTree[0]
	nodesDummyTree[0].height = 2
	nodesDummyTree[0].sizeLeft = 3
	nodesDummyTree[0].sizeRight = 1
	nodesDummyTree[1].height = 0
	nodesDummyTree[1].sizeLeft = 0
	nodesDummyTree[1].sizeRight = 0
	nodesDummyTree[2].height = 0
	nodesDummyTree[2].sizeLeft = 0
	nodesDummyTree[2].sizeRight = 0
	nodesDummyTree[3].height = 0
	nodesDummyTree[3].sizeLeft = 0
	nodesDummyTree[3].sizeRight = 0
	nodesDummyTree[4].height = 1
	nodesDummyTree[4].sizeLeft = 1
	nodesDummyTree[4].sizeRight = 1
	dummyTree.minValue = nodesDummyTree[1]
	dummyTree.size = 5
	linkNodes(nodesDummyTree[1], nodesDummyTree[4], true)
	linkNodes(nodesDummyTree[4], nodesDummyTree[3], true)
	linkNodes(nodesDummyTree[3], nodesDummyTree[0], true)
	linkNodes(nodesDummyTree[0], nodesDummyTree[2], true)
	linkNodes(nodesDummyTree[2], nodesDummyTree[1], true)

	// rotating
	tree.insertFixup(nodesTree[3])
	tree.cleanNil()

	// checking

	expectSize := int64(len(values))
	err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateRightLeft, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", tree.size, expectSize, err)
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, dummyTree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		err = fmt.Errorf("on test TestRotateRightLeft, checking tree size (%d) and expected size (%d) falied\n\t-- error: %s", dummyTree.size, expectSize, err)
		t.Fatal(err)
	}

	areEquals, errText := CheckTrees(tree, dummyTree)
	if errText != nil {
		t.Fatal(errText)
		return
	}
	if !areEquals {
		t.Fatal(fmt.Errorf("trees are not equal"))
		return
	}
}

func Test_AddOne(t *testing.T) {
	tree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}

	data := NewTestData()
	data.Id = 0
	data.Text = "First"
	oldData, err := tree.Put(data.Id, data)
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, oldData, tree.avlTreeConstructorParams.ValueZeroValue, EqualTestData, //
		fmt.Sprintf("putting a value on empty tree should return the \"value's zero-value\", but we have: %v", oldData))
	if err != nil {
		t.Fatal(err)
	}
	if tree.root == nil {
		t.Fatalf("the tree's root should NOT be nil\n")
	}
	err = testIsLeaf(tree, tree._NIL)
	if err != nil {
		t.Fatal(err)
	}
	err = testNIL(tree, false, tree.root, "root is _NIL; should not be NIL")
	if err != nil {
		t.Fatal(err)
	}

	//

	err = testEqualityPrimitive(true, tree.Size(), 1, "size should be 1")
	if err != nil {
		t.Fatal(err)
	}

	if tree._NIL == nil {
		t.Fatalf("the tree's \"_NIL\" should NOT be nil\n")
	}

	// internal nodes disposition

	err = testNIL(tree, true, tree.root.father, "father is not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.height, 0, "new node's height should be: 0")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeLeft, 0, "new node's sizeLeft should be 0")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeRight, 0, "new node's sizeRight should be 0")
	if err != nil {
		t.Fatal(err)
	}

	err = testEqualityObj(true, tree.minValue, tree.root, EqualData, "min value node should be equal to root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.firstInserted, tree.root, EqualData, "first inserted node should be equal to root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}

	err = testEqualityObj(true, tree.root.nextInOrder, tree.root, EqualData, "nextInOrder should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.prevInOrder, tree.root, EqualData, "prevInOrder should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.nextInserted, tree.root, EqualData, "nextInserted should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.prevInserted, tree.root, EqualData, "prevInserted should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_AddOne_2WithSameKey_Replace(t *testing.T) {

	tree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}

	data := NewTestData()
	data.Id = 0
	data.Text = "First"
	_, err = tree.Put(data.Id, data)
	if err != nil {
		t.Fatal(err)
	}

	d2 := NewTestData()
	d2.Id = data.Id
	d2.Text = "Second"
	oldData, err := tree.Put(d2.Id, d2)
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, oldData, data, EqualTestData, "should be the first data-value inserted")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.keyVal.value, d2, EqualTestData, "should be the second data-value inserted, since it should be replaced")
	if err != nil {
		t.Fatal(err)
	}

	//

	err = testEqualityPrimitive(true, tree.Size(), 1, "size should be 1")
	if err != nil {
		t.Fatal(err)
	}

	if tree._NIL == nil {
		t.Fatalf("the tree's \"_NIL\" should NOT be nil\n")
	}

	// internal nodes disposition

	err = testNIL(tree, true, tree.root.father, "father is not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.height, 0, "new node's height should be: 0")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeLeft, 0, "new node's sizeLeft should be 0")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeRight, 0, "new node's sizeRight should be 0")
	if err != nil {
		t.Fatal(err)
	}

	err = testEqualityObj(true, tree.minValue, tree.root, EqualData, "min value node should be equal to root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.firstInserted, tree.root, EqualData, "first inserted node should be equal to root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}

	err = testEqualityObj(true, tree.root.nextInOrder, tree.root, EqualData, "nextInOrder should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.prevInOrder, tree.root, EqualData, "prevInOrder should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.nextInserted, tree.root, EqualData, "nextInserted should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.prevInserted, tree.root, EqualData, "prevInserted should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
}
func Test_AddOne_2WithSameKey_Ignore(t *testing.T) {

	tree, err := NewTree()
	tree.avlTreeConstructorParams.KeyCollisionBehavior = IgnoreInsertion
	if err != nil {
		t.Fatal(err)
	}

	data := NewTestData()
	data.Id = 0
	data.Text = "First"
	_, err = tree.Put(data.Id, data)
	if err != nil {
		t.Fatal(err)
	}

	d2 := NewTestData()
	d2.Id = data.Id
	d2.Text = "Second"
	oldData, err := tree.Put(d2.Id, d2)
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, oldData, data, EqualTestData, "should be the first data-value inserted")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.keyVal.value, data, EqualTestData, "should be the second data-value inserted, as if the Put would be rejected")
	if err != nil {
		t.Fatal(err)
	}

	//

	err = testEqualityPrimitive(true, tree.Size(), 1, "size should be 1")
	if err != nil {
		t.Fatal(err)
	}

	if tree._NIL == nil {
		t.Fatalf("the tree's \"_NIL\" should NOT be nil\n")
	}

	// internal nodes disposition

	err = testNIL(tree, true, tree.root.father, "father is not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.height, 0, "new node's height should be: 0")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeLeft, 0, "new node's sizeLeft should be 0")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeRight, 0, "new node's sizeRight should be 0")
	if err != nil {
		t.Fatal(err)
	}

	err = testEqualityObj(true, tree.minValue, tree.root, EqualData, "min value node should be equal to root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.firstInserted, tree.root, EqualData, "first inserted node should be equal to root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}

	err = testEqualityObj(true, tree.root.nextInOrder, tree.root, EqualData, "nextInOrder should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.prevInOrder, tree.root, EqualData, "prevInOrder should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.nextInserted, tree.root, EqualData, "nextInserted should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.prevInserted, tree.root, EqualData, "prevInserted should loop to itself, i.e. root, since it's the only node here")
	if err != nil {
		t.Fatal(err)
	}
}

// adding 2: [2,1], [2,3]

func Test_AddOne_2_InOrder(t *testing.T) {

	tree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}

	keys := []int{2, 1}
	var datas = make([]*TestData, len(keys))
	for i, k := range keys {

		data := NewTestData()
		data.Id = k
		data.Text = fmt.Sprintf("v_%d", i)
		datas[i] = data

		_, err = tree.Put(data.Id, data)
		if err != nil {
			t.Fatal(err)
		}
	}

	err = testNIL(tree, false, tree.root, "root should not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.keyVal.value, datas[0], EqualTestData, //
		fmt.Sprintf("root (%v) should be: %v", tree.root.keyVal.value, datas[0]))
	if err != nil {
		t.Fatal(err)
	}

	err = testNIL(tree, false, tree.root.left, "root's left should not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.left.keyVal.value, datas[1], EqualTestData, //
		fmt.Sprintf("root's left (%v) should be: %v", tree.root.left.keyVal.value, datas[1]))
	if err != nil {
		t.Fatal(err)
	}

	//

	expectSize := int64(len(datas))
	err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		t.Fatal(err)
	}

	// internal nodes disposition

	err = testNIL(tree, true, tree.root.father, "father is not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testNIL(tree, true, tree.root.right, "root's right should be _NIL, but it's not")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.height, 1, "new node's height should be: 1")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeLeft, 1, "new node's sizeLeft should be 1")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeRight, 0, "new node's sizeRight should be 0")
	if err != nil {
		t.Fatal(err)
	}

	err = testEqualityObj(true, tree.minValue, tree.root.left, EqualData, "min value node should be equal to root's left")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.firstInserted, tree.root, EqualData, "first inserted node should be equal to root")
	if err != nil {
		t.Fatal(err)
	}

	err = testEqualityObj(true, tree.root.nextInOrder, tree.root.left, EqualData, "root's nextInOrder should be root's left")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.prevInOrder, tree.root.left, EqualData, "root's prevInOrder should be root's left")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.nextInserted, tree.root.left, EqualData, "root's nextInserted should be root's left")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.prevInserted, tree.root.left, EqualData, "root's prevInserted should be root's left")
	if err != nil {
		t.Fatal(err)
	}

	err = testEqualityObj(true, tree.root.left.father, tree.root, EqualData, "the second node's father should be root")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.left.nextInOrder, tree.root, EqualData, "the second node's nextInOrder should be root")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.left.prevInOrder, tree.root, EqualData, "the second node's prevInOrder should be root")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.left.nextInserted, tree.root, EqualData, "the second node's nextInserted should be root")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.left.prevInserted, tree.root, EqualData, "the second node's prevInserted should be root")
	if err != nil {
		t.Fatal(err)
	}
}

func Test_AddOne_2_ReverseOrder(t *testing.T) {

	tree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}

	keys := []int{2, 3}
	var datas = make([]*TestData, len(keys))
	for i, k := range keys {

		data := NewTestData()
		data.Id = k
		data.Text = fmt.Sprintf("v_%d", i)
		datas[i] = data

		_, err = tree.Put(data.Id, data)
		if err != nil {
			t.Fatal(err)
		}
	}

	err = testNIL(tree, false, tree.root, "root should not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.keyVal.value, datas[0], EqualTestData, //
		fmt.Sprintf("root (%v) should be: %v", tree.root.keyVal.value, datas[0]))
	if err != nil {
		t.Fatal(err)
	}

	err = testNIL(tree, false, tree.root.right, "root's right should not be _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.right.keyVal.value, datas[1], EqualTestData, //
		fmt.Sprintf("root's right (%v) should be: %v", tree.root.right.keyVal.value, datas[1]))
	if err != nil {
		t.Fatal(err)
	}

	//

	expectSize := int64(len(datas))
	err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
	if err != nil {
		t.Fatal(err)
	}

	// internal nodes disposition

	err = testNIL(tree, true, tree.root.father, "father is not _NIL")
	if err != nil {
		t.Fatal(err)
	}
	err = testNIL(tree, true, tree.root.left, "root's left should be _NIL, but it's not")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.height, 1, "new node's height should be: 1")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeLeft, 0, "new node's sizeleft should be 0")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityPrimitive(true, tree.root.sizeRight, 1, "new node's sizeRight should be 1")
	if err != nil {
		t.Fatal(err)
	}

	err = testEqualityObj(true, tree.minValue, tree.root, EqualData, "min value node should be equal to root")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.firstInserted, tree.root, EqualData, "first inserted node should be equal to root")
	if err != nil {
		t.Fatal(err)
	}

	err = testEqualityObj(true, tree.root.nextInOrder, tree.root.right, EqualData, "root's nextInOrder should be root's right")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.prevInOrder, tree.root.right, EqualData, "root's prevInOrder should be root's right")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.nextInserted, tree.root.right, EqualData, "root's nextInserted should be root's right")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.prevInserted, tree.root.right, EqualData, "root's prevInserted should be root's right")
	if err != nil {
		t.Fatal(err)
	}

	err = testEqualityObj(true, tree.root.right.father, tree.root, EqualData, "the second node's father should be root")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.right.nextInOrder, tree.root, EqualData, "the second node's nextInOrder should be root")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.right.prevInOrder, tree.root, EqualData, "the second node's prevInOrder should be root")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.right.nextInserted, tree.root, EqualData, "the second node's nextInserted should be root")
	if err != nil {
		t.Fatal(err)
	}
	err = testEqualityObj(true, tree.root.right.prevInserted, tree.root, EqualData, "the second node's prevInserted should be root")
	if err != nil {
		t.Fatal(err)
	}
}

//

type AddSetup struct {
	name string
	keys []int
	/* Predicted node order if the tree would be traversed in a "breadth first" way.
	As a clarifying example, the first value of this array, i.e. the number at index 0,
	refers to the tree's root. The stored values are the indexes of the "keys" above,
	which could be used to retrieve the aforementioned nodes.
	*/
	onBreadthVisit_IndexKeyData []int
	datas                       []*TestData
}

func Test_Add_3(t *testing.T) {

	/* adding 3: - 6 tests
	   -) no rotation in order: [2,1,3]
	   -) no rotation anti- order: [2,3,1]
	   -) "left left": [3,2,1]
	   - - albero prima della rotazione: ( ((1) -> 2 <- ()) -> 3 <- () )
	   - - albero dopo della rotazione: ((1) -> 2 <- (3))
	   -) "left right": [3,1,2] ... rotazione articolata
	   - - albero prima della rotazione: ( (() -> 1 <- (2)) -> 3 <- () )
	   - - albero dopo della rotazione: ((1) -> 2 <- (3))
	   -) "right right": [1,2,3]
	   - - albero prima della rotazione: ( () -> 1 <- (() -> 2 <- (3)) )
	   - - albero dopo della rotazione: ((1) -> 2 <- (3))
	   -) "right left": [1,3,2] ... rotazione articolata
	   - - albero prima della rotazione: ( () -> 1 <- ((2) -> 3 <- ()) )
	   - - albero dopo della rotazione: ((1) -> 2 <- (3))
	*/
	setups := []AddSetup{
		{"no rotation - in order", []int{2, 1, 3}, []int{
			0, 1, 2,
		}, nil},
		{"no rotation - reverse order", []int{2, 3, 1}, []int{
			0, 2, 1,
		}, nil},
		{"left left", []int{3, 2, 1}, []int{
			1, 2, 0, //2, 0, 1,
		}, nil},
		{"left right", []int{3, 1, 2}, []int{2, 1, 0}, nil},
		{"right right", []int{1, 2, 3}, []int{1, 0, 2}, nil},
		{"right left", []int{1, 3, 2}, []int{2, 0, 1}, nil},
	}

	dummyTree, err := NewTree()
	if err != nil {
		t.Fatal(err)
		return
	}
	dummyNodes := []*AVLTNode[int, *TestData]{ // keep the Breadth's first visit order
		dummyTree.newNode(2, NewTestDataFilled(2, "root")),
		dummyTree.newNode(1, NewTestDataFilled(1, "left")),
		dummyTree.newNode(3, NewTestDataFilled(3, "right")),
	}
	dummyTree.root = dummyNodes[0]
	dummyTree.minValue = dummyNodes[1]
	dummyTree.root.left = dummyNodes[1]
	dummyTree.root.left.father = dummyTree.root
	dummyTree.root.right = dummyNodes[2]
	dummyTree.root.right.father = dummyTree.root
	dummyTree.size = 3
	dummyTree.root.height = 1
	dummyTree.root.sizeLeft = 1
	dummyTree.root.sizeRight = 1
	// adjust in-order
	/*for i, n := range dummyNodes {
		n.nextInOrder = dummyNodes[(i+1)%len(dummyNodes)]
		i_prev := i - 1
		if i_prev < 0 {
			i_prev += len(dummyNodes)
		}
		n.prevInOrder = dummyNodes[i_prev]
	}*/
	dummyTree.root.nextInOrder = dummyTree.root.right
	dummyTree.root.prevInOrder = dummyTree.root.left
	dummyTree.root.left.nextInOrder = dummyTree.root
	dummyTree.root.left.prevInOrder = dummyTree.root.right
	dummyTree.root.right.nextInOrder = dummyTree.root.left
	dummyTree.root.right.prevInOrder = dummyTree.root

	for _, data := range setups {

		// adjust inserted-chronological order

		for i := 0; i < len(dummyNodes); i++ {
			n := dummyNodes[data.onBreadthVisit_IndexKeyData[i]]
			if i == 0 {
				dummyTree.firstInserted = n
			}
			nextNode := dummyNodes[data.onBreadthVisit_IndexKeyData[(i+1)%len(dummyNodes)]]
			n.nextInserted = nextNode
			nextNode.prevInserted = n
		}

		tree, err := NewTree()
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Fatal(err)
		}

		data.datas = make([]*TestData, len(data.keys))

		for i, id := range data.keys {
			dataTest := NewTestData()
			dataTest.Id = id
			dataTest.Text = fmt.Sprintf("v_%d", i)
			data.datas[i] = dataTest

			_, err = tree.Put(dataTest.Id, dataTest)
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
		}

		// early definitions
		//var nodeNextInserted, nodePrevInserted *AVLTNode[int, *TestData]

		// root checks

		err = testNIL(tree, false, tree.root, "root should not _NIL")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Fatal(err)
		}

		expectSize := int64(len(data.datas))
		err = testEqualityPrimitive(true, tree.Size(), expectSize, fmt.Sprintf("size should be %d", expectSize))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Fatal(err)
		}

		areEquals, errText := CheckTrees(tree, dummyTree)
		if errText != nil {
			t.Fatal(errText)
			return
		}
		if !areEquals {
			t.Fatal(fmt.Errorf("trees are not equal"))
			return
		}

		err = testNIL(tree, true, tree.root.father, "father is not _NIL")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Fatal(err)
		}

		err = testEqualityPrimitive(true, tree.root.height, 1, "new node's height should be: 1")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Fatal(err)
		}
		err = testEqualityPrimitive(true, tree.root.sizeLeft, 1, "new node's sizeleft should be 1")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Fatal(err)
		}
		err = testEqualityPrimitive(true, tree.root.sizeRight, 1, "new node's sizeRight should be 1")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Fatal(err)
		}

		err = testEqualityObj(true, tree.minValue, tree.root.left, EqualData, "min value node should be equal to root 's left")
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Fatal(err)
		}
		/*
			err = testEqualityObj(true, tree.firstInserted, firstNodeInserted, EqualData, fmt.Sprintf("first inserted node (%v) should be equal to : %v", tree.firstInserted.keyVal.value, firstNodeInserted.keyVal.value))
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
		*/

		err = testEqualityObj(true, tree.root.nextInOrder, tree.root.right, EqualData, fmt.Sprintf("root's nextInOrder (whose value is: %v) should be root's right, with value: %v", tree.root.nextInOrder.keyVal.value, tree.root.right.keyVal.value))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Fatal(err)
		}
		err = testEqualityObj(true, tree.root.prevInOrder, tree.root.left, EqualData, fmt.Sprintf("root's prevInOrder (whose value is: %v) should be root's left, with value: %v", tree.root.prevInOrder.keyVal.value, tree.root.left.keyVal.value))
		if err != nil {
			err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
			t.Fatal(err)
		}

		/*
			indexRoot := 0
			rootsKeyIndex := data.onBreadthVisit_IndexKeyData[indexRoot]
			// indexLeft := data.onBreadthVisit_IndexKeyData[1]
			// indexRight := data.onBreadthVisit_IndexKeyData[2]
			dataRootExpected := data.datas[rootsKeyIndex]
			err = testEqualityObj(true, tree.root.keyVal.key, dataRootExpected.Id, EqualKey, //
				fmt.Sprintf("root key (%d) should be: %d", tree.root.keyVal.key, dataRootExpected.Id))
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
			err = testEqualityObj(true, tree.root.keyVal.value, dataRootExpected, EqualTestData, //
				fmt.Sprintf("root value (%v) should be: %v", tree.root.keyVal.value, dataRootExpected))
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}

			index_onBreadthVisit_IndexKeyData_root := 0
			index_onBreadthVisit_IndexKeyData_left := index_onBreadthVisit_IndexKeyData_root + 1
			index_onBreadthVisit_IndexKeyData_right := index_onBreadthVisit_IndexKeyData_root + 2

			indexKey_WhenRootWasAdded := data.onBreadthVisit_IndexKeyData[index_onBreadthVisit_IndexKeyData_root]
			dataRoot := data.datas[data.onBreadthVisit_IndexKeyData[index_onBreadthVisit_IndexKeyData_root]]
			rootRecalcolated := tree.getNode(dataRoot.Id)
			if rootRecalcolated != tree.root {
				t.Fatal("the tests are wrong! root has been wrongly indexed")
			}
			firstNodeInserted := tree.getNode(data.keys[0])

			indexRootNextInserted := data.onBreadthVisit_IndexKeyData[(indexKey_WhenRootWasAdded+1)%len(data.keys)]
			dataRootNextInserted := data.datas[indexRootNextInserted]
			tempIndexPrev := (indexKey_WhenRootWasAdded - 1)
			if tempIndexPrev < 0 {
				tempIndexPrev += len(data.keys)
			}
			indexRootPrevInserted := data.onBreadthVisit_IndexKeyData[tempIndexPrev]
			dataRootPrevInserted := data.datas[indexRootPrevInserted]

			nodeNextInserted = tree.getNode(dataRootNextInserted.Id)
			nodePrevInserted = tree.getNode(dataRootPrevInserted.Id)
			// secondNodeInserted = tree.getNode(data.keys[1])
			// thirdNodeInserted = tree.getNode(data.keys[2]) // 3 elements -> "3-1" then preceeds "0"


				// firstNodeInserted := tree.getNode(dataRootExpected.Id)
				// secondNodeNextInserted = tree.getNode(data.datas[indexLeft].Id)
				// thirdNodePrevInserted = tree.getNode( data.datas[indexRight].Id ) // 3 elements -> "3-1" then preceeds "0"


			if err != nil {
				err = testEqualityObj(true, tree.root.nextInserted, nodeNextInserted, EqualData, fmt.Sprintf("root's nextInserted (whose value is: %v) should be the node with value: %v", tree.root.nextInserted.keyVal.value, nodeNextInserted.keyVal.value))
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
			err = testEqualityObj(true, tree.root.prevInserted, nodePrevInserted, EqualData,
				fmt.Sprintf("\n\troot's prevInserted (whose value is: %v)\n\tshould be the node with value: %v\n\t(fetched with key: %d; index indexRootPrevInserted: %d, tempIndexPrev: %d, indexKey_WhenRootWasAdded:%d)",
					tree.root.prevInserted.keyVal.value, nodePrevInserted.keyVal.value,
					dataRootPrevInserted.Id, indexRootPrevInserted, tempIndexPrev, indexKey_WhenRootWasAdded))
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error:\n\t %s", data.name, data.keys, err)
				t.Fatal(err)
			}

			// root's left checks
			left := tree.root.left
			indexKey_WhenLeftWasAdded := data.onBreadthVisit_IndexKeyData[index_onBreadthVisit_IndexKeyData_left]
			// dataLeft := data.datas[data.onBreadthVisit_IndexKeyData[index_onBreadthVisit_IndexKeyData_left]]
			dataLeftNextInserted := data.datas[data.onBreadthVisit_IndexKeyData[(indexKey_WhenLeftWasAdded+1)%len(data.keys)]]
			tempIndexPrev = (indexKey_WhenLeftWasAdded - 1)
			if tempIndexPrev < 0 {
				tempIndexPrev += len(data.keys)
			}
			dataLeftPrevInserted := data.datas[data.onBreadthVisit_IndexKeyData[tempIndexPrev]]

			nodeNextInserted = tree.getNode(dataLeftNextInserted.Id)
			nodePrevInserted = tree.getNode(dataLeftPrevInserted.Id)

			err = testNIL(tree, false, left, "root's left should not be _NIL")
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
			err = testEqualityObj(true, left.father, tree.root, EqualData, fmt.Sprintf( //
				"left's father (left value: %v) should be root (value: %v), but we have as father: %v", //
				left.keyVal.value, tree.root.keyVal.value, left.father.keyVal.value))
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}

			err = testNIL(tree, true, left.left, "left's left should be _NIL")
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
			err = testNIL(tree, true, left.right, "left's right should be _NIL")
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}

			err = testEqualityObj(true, left.nextInOrder, tree.root, EqualData, fmt.Sprintf("root's left's nextInOrder should be root, with value: %v", tree.root.keyVal.value))
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
			err = testEqualityObj(true, left.prevInOrder, tree.root.right, EqualData, fmt.Sprintf("root's left's prevInOrder should be root's right, with value: %v", tree.root.right.keyVal.value))
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
			if err != nil {
				err = testEqualityObj(true, left.nextInserted, nodeNextInserted, EqualData, fmt.Sprintf("left's nextInserted should be the node with value: %v", nodeNextInserted.keyVal.value))
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
			err = testEqualityObj(true, left.prevInserted, nodePrevInserted, EqualData, fmt.Sprintf("left's prevInserted should be the node with value: %v", nodePrevInserted.keyVal.value))
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}

			// root's right checks
			right := tree.root.right
			indexKey_WhenRightWasAdded := data.onBreadthVisit_IndexKeyData[index_onBreadthVisit_IndexKeyData_right]
			// dataLeft := data.datas[data.onBreadthVisit_IndexKeyData[index_onBreadthVisit_IndexKeyData_right]]
			dataRightNextInserted := data.datas[data.onBreadthVisit_IndexKeyData[(indexKey_WhenRightWasAdded+1)%len(data.keys)]]
			tempIndexPrev = (indexKey_WhenRightWasAdded - 1)
			if tempIndexPrev < 0 {
				tempIndexPrev += len(data.keys)
			}
			dataRightPrevInserted := data.datas[data.onBreadthVisit_IndexKeyData[tempIndexPrev]]

			nodeNextInserted = tree.getNode(dataRightNextInserted.Id)
			nodePrevInserted = tree.getNode(dataRightPrevInserted.Id)

			err = testNIL(tree, false, right, "root's right should not be _NIL")
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
			err = testEqualityObj(true, right.father, tree.root, EqualData, fmt.Sprintf( //
				"right's father (right value: %v) should be root (value: %v), but we have as father: %v", //
				right.keyVal.value, tree.root.keyVal.value, right.father.keyVal.value))
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}

			err = testNIL(tree, true, right.left, "right's left should be _NIL")
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
			err = testNIL(tree, true, right.right, "right's right should be _NIL")
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}

			err = testEqualityObj(true, right.nextInOrder, tree.root.left, EqualData, fmt.Sprintf("root's right's nextInOrder should be root's left, with value: %v", tree.root.left.keyVal.value))
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
			err = testEqualityObj(true, right.prevInOrder, tree.root, EqualData, fmt.Sprintf("root's right's prevInOrder should be root, with value: %v", tree.root.keyVal.value))
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
			if err != nil {
				err = testEqualityObj(true, right.nextInserted, nodeNextInserted, EqualData, fmt.Sprintf("right's nextInserted should be the node with value: %v", nodeNextInserted.keyVal.value))
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
			err = testEqualityObj(true, right.prevInserted, nodePrevInserted, EqualData, fmt.Sprintf("right's prevInserted should be the node with value: %v", nodePrevInserted.keyVal.value))
			if err != nil {
				err = fmt.Errorf("on test {{\"%s\" - %v}} -- error: %s", data.name, data.keys, err)
				t.Fatal(err)
			}
		*/
	}
}

type TreeAlterationTestFunction func(t *testing.T, treeOriginal *AVLTree[int, *TestData], treeDummy *AVLTree[int, *TestData], index int, data *TestData)

func Test_Add_Massivo(t *testing.T) {
	tree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	dummyTree, err := NewTree()
	if err != nil {
		t.Fatal(err)
	}
	values := []int{
		20, 10, 30, //
		3, 5, //
		50, 100, //
		2, 1, //
		42, 37, //
	}
	// at the end
	//   .   .   .   .   .   . 20.
	//   .   .   .   .   . / .   . \ .   .
	//   .   .   .   . / .   .   .   . \ .
	//   .   .   . / .   .   .   .   .   . \ .
	//   .   . 5 .   .   .   .   .   .   .   .50 .
	//   .  /.   .\  .   .   .   .   .   . / .   .\  .
	//   2   .   .  10   .   .   .   .  42   .   .  100  .
	//1  .  3.   .   .   .   .   .   .37 .50 .

	valuesInTotal := len(values)
	// nodesTree := make([]*AVLTNode[int, *TestData], valuesInTotal)
	nodesDummyTree := make([]*AVLTNode[int, *TestData], valuesInTotal)
	alterationFns := make([]TreeAlterationTestFunction, valuesInTotal)

	for i := 0; i < valuesInTotal; i++ {
		alterationFns[i] = nil
		// nodesTree[i] = NewTreeNodeFilled(tree, values[i], fmt.Sprintf("v_%d", values[i]))
		nodesDummyTree[i] = NewTreeNodeFilled(dummyTree, values[i], fmt.Sprintf("v_%d", values[i]))
	}
	/*
		for _, nt := range [][]*AVLTNode[int, *TestData]{nodesTree, nodesDummyTree} {
			for i := 0; i < valuesInTotal; i++ {
				linkNodes(nt[i], nt[(i+1)%len(values)], false)
			}
		}
	*/
	for i := 0; i < valuesInTotal; i++ {
		linkNodes(nodesDummyTree[i], nodesDummyTree[(i+1)%len(values)], false)
	}

	// test functions

	newBaseErrorText := func(index int, data *TestData, additionalText string) string {
		return fmt.Sprintf("Error upon ADD MASSIVO on index %d with value %v, upon Put:\n\t%s\n", index, data.String(), additionalText)
	}

	indexNode := 0
	alterationFns[indexNode] = func(t *testing.T, treeOriginal *AVLTree[int, *TestData], treeDummy *AVLTree[int, *TestData], index int, data *TestData) {
		oldV, err := treeOriginal.Put(data.Id, data)
		if err != nil {
			t.Fatal(newBaseErrorText(index, data, err.Error()))
		}
		if oldV != treeOriginal.avlTreeConstructorParams.ValueZeroValue {
			t.Fatalf(newBaseErrorText(index, data, fmt.Sprintf("returned value upon first Put is not the zero value:\n\t <%v> - <%v>\n", oldV,
				treeOriginal.avlTreeConstructorParams.ValueZeroValue)))
		}

		// set-up tree dummy
		nd := nodesDummyTree[index] // 0-th
		treeDummy.root = nd
		linkNills(treeDummy, nd, true)
		treeDummy.size = 1
		treeDummy.minValue = nd
		treeDummy.firstInserted = nd
		linkNodes(nd, nd, true)
		linkNodes(nd, nd, false)

		// check equality

		err = nil
		equal, checkError := CheckTrees(treeOriginal, treeDummy)
		if checkError != nil {
			errorText := checkError.Error()
			additionalText := fmt.Sprintf("after %d-th Put, error: %s", index, errorText)
			t.Fatal(newBaseErrorText(index, data, additionalText))
		}
		if !equal {
			t.Fatal(newBaseErrorText(index, data,
				fmt.Sprintf("trees should be equal:\n\toriginal tree: %s\n\tdummy tree: %s\n", treeOriginal, treeDummy)))
		}
	}
	indexNode++ // 1
	alterationFns[indexNode] = func(t *testing.T, treeOriginal *AVLTree[int, *TestData], treeDummy *AVLTree[int, *TestData], index int, data *TestData) {
		oldV, err := treeOriginal.Put(data.Id, data)
		if err != nil {
			t.Fatal(newBaseErrorText(index, data, err.Error()))
			return
		}
		if oldV != treeOriginal.avlTreeConstructorParams.ValueZeroValue {
			t.Fatalf(newBaseErrorText(index, data, fmt.Sprintf("returned value upon second Put is not the zero value:\n\t <%v> - <%v>\n", oldV,
				treeOriginal.avlTreeConstructorParams.ValueZeroValue)))
		}

		// set-up tree dummy
		//   . 20
		//   /
		// 10
		nd := nodesDummyTree[index] // 1-th == 10
		linkNills(treeDummy, nd, true)
		treeDummy.root.left = nd
		nd.father = treeDummy.root
		treeDummy.size++
		linkNodes(nd, treeDummy.minValue, true)
		linkNodes(treeDummy.minValue, nd, true) // close the loop
		treeDummy.minValue = nd
		linkNodes(treeDummy.firstInserted, nd, false)
		linkNodes(nd, treeDummy.firstInserted, false) // close the loop

		treeDummy.root.sizeLeft = 1
		treeDummy.root.height = 1

		// check equality

		equal, checkError := CheckTrees(treeOriginal, treeDummy)
		if checkError != nil {
			errorText := checkError.Error()
			additionalText := fmt.Sprintf("after %d-th Put, error: %s", index, errorText)
			t.Fatal(newBaseErrorText(index, data, additionalText))
		}
		if !equal {
			t.Fatal(newBaseErrorText(index, data,
				fmt.Sprintf("trees should be equal:\n\toriginal tree: %s\n\tdummy tree: %s\n", treeOriginal, treeDummy)))
		}
	}
	indexNode++ // 2
	alterationFns[indexNode] = func(t *testing.T, treeOriginal *AVLTree[int, *TestData], treeDummy *AVLTree[int, *TestData], index int, data *TestData) {
		oldV, err := treeOriginal.Put(data.Id, data)
		if err != nil {
			t.Fatal(newBaseErrorText(index, data, err.Error()))
		}
		if oldV != treeOriginal.avlTreeConstructorParams.ValueZeroValue {
			t.Fatalf(newBaseErrorText(index, data, fmt.Sprintf("returned value upon third Put is not the zero value:\n\t <%v> - <%v>\n", oldV,
				treeOriginal.avlTreeConstructorParams.ValueZeroValue)))
		}

		// set-up tree dummy
		//   . 20.
		//   /   .\
		// 10.   . 30
		nd := nodesDummyTree[index] // 2-th == 30
		linkNills(treeDummy, nd, true)
		treeDummy.root.right = nd
		nd.father = treeDummy.root
		treeDummy.size++
		linkNodes(treeDummy.root, nd, true)
		linkNodes(nd, treeDummy.root.left, true) // close the loop
		if treeDummy.root.right != treeDummy.minValue.prevInOrder {
			t.Fatalf(newBaseErrorText(index, data, fmt.Sprintf("returned value upon third Put is not the zero value:\n\t <%v> - <%v>\n", oldV,
				treeOriginal.avlTreeConstructorParams.ValueZeroValue)))
		}
		linkNodes(treeDummy.root.left, nd, false)
		linkNodes(nd, treeDummy.firstInserted, false) // close the loop

		treeDummy.root.sizeRight = 1

		// check equality
		equal, checkError := CheckTrees(treeOriginal, treeDummy)

		if checkError != nil {
			errorText := checkError.Error()
			additionalText := fmt.Sprintf("after %d-th Put, error: %s", index, errorText)
			t.Fatal(newBaseErrorText(index, data, additionalText))
		}
		if !equal {
			t.Fatal(newBaseErrorText(index, data,
				fmt.Sprintf("trees should be equal:\n\toriginal tree: %s\n\tdummy tree: %s\n", treeOriginal, treeDummy)))
		}
	}
	indexNode++ // 3
	alterationFns[indexNode] = func(t *testing.T, treeOriginal *AVLTree[int, *TestData], treeDummy *AVLTree[int, *TestData], index int, data *TestData) {
		oldV, err := treeOriginal.Put(data.Id, data)
		if err != nil {
			t.Fatal(newBaseErrorText(index, data, err.Error()))
		}
		if oldV != treeOriginal.avlTreeConstructorParams.ValueZeroValue {
			t.Fatalf(newBaseErrorText(index, data, fmt.Sprintf("returned value upon fourth Put is not the zero value:\n\t <%v> - <%v>\n", oldV,
				treeOriginal.avlTreeConstructorParams.ValueZeroValue)))
		}

		// set-up tree dummy
		//   . 20.
		//   /   .\
		// 10.   . 30
		//3
		nd := nodesDummyTree[index] // 3-th == 3
		linkNills(treeDummy, nd, true)
		treeDummy.root.left.left = nd
		nd.father = treeDummy.root.left
		treeDummy.size++
		linkNodes(treeDummy.minValue.prevInOrder, nd, true) // 10->30 & 3
		linkNodes(nd, treeDummy.minValue, true)             // close the loop
		treeDummy.minValue = nd
		linkNodes(treeDummy.root.right, nd, false)
		linkNodes(nd, treeDummy.firstInserted, false) // close the loop

		treeDummy.root.sizeLeft = 2
		treeDummy.root.height = 2
		treeDummy.root.left.sizeLeft = 1
		treeDummy.root.left.height = 1

		// check equality

		equal, checkError := CheckTrees(treeOriginal, treeDummy)
		if checkError != nil {
			errorText := checkError.Error()
			additionalText := fmt.Sprintf("after %d-th Put, error: %s", index, errorText)
			t.Fatal(newBaseErrorText(index, data, additionalText))
		}
		if !equal {
			t.Fatal(newBaseErrorText(index, data,
				fmt.Sprintf("trees should be equal:\n\toriginal tree: %s\n\tdummy tree: %s\n", treeOriginal, treeDummy)))
		}
	}
	indexNode++ // 4
	alterationFns[indexNode] = func(t *testing.T, treeOriginal *AVLTree[int, *TestData], treeDummy *AVLTree[int, *TestData], index int, data *TestData) {
		oldV, err := treeOriginal.Put(data.Id, data)

		if err != nil {
			t.Fatal(newBaseErrorText(index, data, err.Error()))
		}
		if oldV != treeOriginal.avlTreeConstructorParams.ValueZeroValue {
			t.Fatalf(newBaseErrorText(index, data, fmt.Sprintf("returned value upon fourth Put is not the zero value:\n\t <%v> - <%v>\n", oldV,
				treeOriginal.avlTreeConstructorParams.ValueZeroValue)))
		}

		// set-up tree dummy
		//   . 20.
		//   /   .\
		// 10.   . 30
		//3
		// 5
		// -> rotate! LEFT-RIGHT

		//   . 20.
		//   /   .\
		//  5.   . 30
		//3 10
		nd := nodesDummyTree[index] // 4-th == 5
		treeDummy.cleanNil()
		linkNills(treeDummy, nd, true)
		treeDummy.size++
		node_10 := treeDummy.root.left
		node_3 := node_10.left
		linkNills(treeDummy, node_3, false)
		linkNills(treeDummy, node_10, false)
		treeDummy.root.left = nd
		nd.father = treeDummy.root
		nd.left = node_3
		node_3.father = nd
		nd.right = node_10
		node_10.father = nd

		nd.father = treeDummy.root.left
		linkNodes(node_3, nd, true)
		linkNodes(nd, node_10, true) // close the loop
		linkNodes(node_3, nd, false)
		linkNodes(nd, treeDummy.firstInserted, false) // close the loop

		treeDummy.root.sizeLeft = 3
		treeDummy.root.left.sizeLeft = 1 // nd
		treeDummy.root.left.sizeRight = 1
		nd.height = 1
		node_3.sizeLeft = 0
		node_3.sizeRight = 0
		node_3.height = 0
		node_10.sizeLeft = 0
		node_10.sizeRight = 0
		node_10.height = 0

		// check equality
		equal, checkError := CheckTrees(treeOriginal, treeDummy)
		if checkError != nil {
			errorText := checkError.Error()
			additionalText := fmt.Sprintf("after %d-th Put, error: %s", index, errorText)
			t.Fatal(newBaseErrorText(index, data, additionalText))

		}
		if !equal {
			t.Fatal(newBaseErrorText(index, data,
				fmt.Sprintf("trees should be equal:\n\toriginal tree: %s\n\tdummy tree: %s\n", treeOriginal, treeDummy)))
		}
	}
	indexNode++ // 5
	alterationFns[indexNode] = func(t *testing.T, treeOriginal *AVLTree[int, *TestData], treeDummy *AVLTree[int, *TestData], index int, data *TestData) {
		oldV, err := treeOriginal.Put(data.Id, data)
		if err != nil {
			t.Fatal(newBaseErrorText(index, data, err.Error()))
		}
		if oldV != treeOriginal.avlTreeConstructorParams.ValueZeroValue {
			t.Fatalf(newBaseErrorText(index, data, fmt.Sprintf("returned value upon fifth Put is not the zero value:\n\t <%v> - <%v>\n", oldV,
				treeOriginal.avlTreeConstructorParams.ValueZeroValue)))
		}

		// set-up tree dummy
		//   . 20.
		//   /   .\
		//  5.   . 30
		//3  10  .   50
		nd := nodesDummyTree[index] // 5-th == 50

		treeDummy.cleanNil()
		linkNills(treeDummy, nd, true)
		treeDummy.size++
		node_30 := treeDummy.root.right
		node_30.right = nd
		nd.father = node_30
		linkNodes(node_30, nd, true)
		linkNodes(nd, treeDummy.minValue, true) // close the loop
		if treeDummy.root.left != treeDummy.firstInserted.prevInserted {
			t.Fatalf("somehow the last inserted (%v) does not coincide with the expected one (%v)\n", treeDummy.firstInserted.prevInserted, treeDummy.root.left)
		}
		lastInserted := treeDummy.firstInserted.prevInserted // node "5"
		linkNodes(lastInserted, nd, false)
		linkNodes(nd, treeDummy.firstInserted, false) // close the loop
		node_30.height = 1
		node_30.sizeRight = 1
		dummyTree.root.sizeRight++

		// check equality
		equal, checkError := CheckTrees(treeOriginal, treeDummy)
		if checkError != nil {
			errorText := checkError.Error()
			additionalText := fmt.Sprintf("after %d-th Put, error: %s", index, errorText)
			t.Fatal(newBaseErrorText(index, data, additionalText))
		}
		if !equal {
			t.Fatal(newBaseErrorText(index, data,
				fmt.Sprintf("trees should be equal:\n\toriginal tree: %s\n\tdummy tree: %s\n", treeOriginal, treeDummy)))
		}
	}
	indexNode++ // 6
	alterationFns[indexNode] = func(t *testing.T, treeOriginal *AVLTree[int, *TestData], treeDummy *AVLTree[int, *TestData], index int, data *TestData) {
		oldV, err := treeOriginal.Put(data.Id, data)
		if err != nil {
			t.Fatal(newBaseErrorText(index, data, err.Error()))
		}
		if oldV != treeOriginal.avlTreeConstructorParams.ValueZeroValue {
			t.Fatalf(newBaseErrorText(index, data, fmt.Sprintf("returned value upon sixth Put is not the zero value:\n\t <%v> - <%v>\n", oldV,
				treeOriginal.avlTreeConstructorParams.ValueZeroValue)))
		}

		// set-up tree dummy
		//   . 20.
		//   /   .\
		//  5.   . 50 // a left rotation happened here
		//3  10  .30 100
		nd := nodesDummyTree[index] // 6-th == 100

		treeDummy.cleanNil()
		linkNills(treeDummy, nd, true)
		treeDummy.size++
		node_30 := treeDummy.root.right
		node_50 := node_30.right
		// TODO FROM HERE
		treeDummy.root.right = node_50
		treeDummy.root.sizeRight = 3
		node_30.sizeRight = 0
		node_30.right = dummyTree._NIL
		node_30.father = node_50
		node_30.height = 0
		node_50.father = dummyTree.root
		node_50.left = node_30
		node_50.right = nd
		node_50.father = dummyTree.root
		node_50.height = 1
		node_50.sizeLeft = 1
		node_50.sizeRight = 1
		nd.father = node_50

		linkNodes(node_50, nd, true)
		linkNodes(nd, treeDummy.minValue, true) // close the loop
		//if treeDummy.root.left != treeDummy.firstInserted.prevInserted {
		//	t.Fatalf("somehow the last inserted (%v) does not coincide with the expected one (%v)\n", treeDummy.firstInserted.prevInserted, treeDummy.root.left)
		//}
		lastInserted := treeDummy.firstInserted.prevInserted // node "5"
		linkNodes(lastInserted, nd, false)
		linkNodes(nd, treeDummy.firstInserted, false) // close the loop

		// check equality
		equal, checkError := CheckTrees(treeOriginal, treeDummy)
		if checkError != nil {
			errorText := checkError.Error()
			additionalText := fmt.Sprintf("after %d-th Put, error: %s", index, errorText)
			t.Fatal(newBaseErrorText(index, data, additionalText))
		}
		if !equal {
			t.Fatal(newBaseErrorText(index, data,
				fmt.Sprintf("trees should be equal:\n\toriginal tree: %s\n\tdummy tree: %s\n", treeOriginal, treeDummy)))
		}
	}
	indexNode++ // 7
	alterationFns[indexNode] = func(t *testing.T, treeOriginal *AVLTree[int, *TestData], treeDummy *AVLTree[int, *TestData], index int, data *TestData) {
		oldV, err := treeOriginal.Put(data.Id, data)
		if err != nil {
			t.Fatal(newBaseErrorText(index, data, err.Error()))
		}
		if oldV != treeOriginal.avlTreeConstructorParams.ValueZeroValue {
			t.Fatalf(newBaseErrorText(index, data, fmt.Sprintf("returned value upon seventh Put is not the zero value:\n\t <%v> - <%v>\n", oldV,
				treeOriginal.avlTreeConstructorParams.ValueZeroValue)))
		}

		// set-up tree dummy
		//   .   . 20.
		//   .  /.   .  \
		//   .5  .   .   . 50
		// 3 .10 .   .  30  100
		//2
		nd := nodesDummyTree[index] // 7-th == 2

		treeDummy.cleanNil()
		linkNills(treeDummy, nd, true)
		treeDummy.size++
		node_3 := treeDummy.root.left.left
		treeDummy.root.height++
		treeDummy.root.sizeLeft++
		treeDummy.root.left.height++
		treeDummy.root.left.sizeLeft++
		node_3.height++
		node_3.sizeLeft++
		node_3.left = nd
		nd.father = node_3
		maxValue := treeDummy.minValue.prevInOrder
		linkNodes(nd, node_3, true)
		linkNodes(maxValue, nd, true) // close the loop
		treeDummy.minValue = nd
		//if treeDummy.root.left != treeDummy.firstInserted.prevInserted {
		//	t.Fatalf("somehow the last inserted (%v) does not coincide with the expected one (%v)\n", treeDummy.firstInserted.prevInserted, treeDummy.root.left)
		//}
		lastInserted := treeDummy.firstInserted.prevInserted // node "5"
		linkNodes(lastInserted, nd, false)
		linkNodes(nd, treeDummy.firstInserted, false) // close the loop

		// check equality
		equal, checkError := CheckTrees(treeOriginal, treeDummy)
		if checkError != nil {
			errorText := checkError.Error()
			additionalText := fmt.Sprintf("after %d-th Put, error: %s", index, errorText)
			t.Fatal(newBaseErrorText(index, data, additionalText))
		}
		if !equal {
			t.Fatal(newBaseErrorText(index, data,
				fmt.Sprintf("trees should be equal:\n\toriginal tree: %s\n\tdummy tree: %s\n", treeOriginal, treeDummy)))
		}
	}
	indexNode++ // 8
	alterationFns[indexNode] = func(t *testing.T, treeOriginal *AVLTree[int, *TestData], treeDummy *AVLTree[int, *TestData], index int, data *TestData) {
		oldV, err := treeOriginal.Put(data.Id, data)
		if err != nil {
			t.Fatal(newBaseErrorText(index, data, err.Error()))
		}
		if oldV != treeOriginal.avlTreeConstructorParams.ValueZeroValue {
			t.Fatalf(newBaseErrorText(index, data, fmt.Sprintf("returned value upon eight Put is not the zero value:\n\t <%v> - <%v>\n", oldV,
				treeOriginal.avlTreeConstructorParams.ValueZeroValue)))
		}

		// set-up tree dummy
		//   .   . 20.
		//   .  /.   .  \
		//   .5  .   .   . 50
		// 2 .10 .   .  30  100
		//1 3 // rotation happened on "2"
		nd := nodesDummyTree[index] // 8-th == 1

		treeDummy.cleanNil()
		linkNills(treeDummy, nd, true)
		treeDummy.size++
		//node_3_path := []int{0, 0, 0}			// []bool{true,true,true}
		node_3 := treeDummy.root.left.left
		//node_2, _ := gnp(treeDummy, node_3_path) // treeDummy.root.left.left.left
		node_2 := node_3.left
		node_2.left = nd
		nd.father = node_2
		node_2.right = node_3
		node_2.father = node_3.father
		node_3.father.left = node_2
		node_3.father = node_2
		node_3.left = treeDummy._NIL
		node_3.right = treeDummy._NIL
		treeDummy.root.sizeLeft++
		treeDummy.root.left.sizeLeft++
		node_3.height = 0
		node_3.sizeLeft = 0
		node_3.sizeRight = 0
		node_2.height = 1
		node_2.sizeLeft = 1
		node_2.sizeRight = 1
		maxValue := treeDummy.minValue.prevInOrder
		linkNodes(nd, node_2, true)
		linkNodes(maxValue, nd, true) // close the loop
		treeDummy.minValue = nd
		//if treeDummy.root.left != treeDummy.firstInserted.prevInserted {
		//	t.Fatalf("somehow the last inserted (%v) does not coincide with the expected one (%v)\n", treeDummy.firstInserted.prevInserted, treeDummy.root.left)
		//}
		lastInserted := treeDummy.firstInserted.prevInserted // node "5"
		linkNodes(lastInserted, nd, false)
		linkNodes(nd, treeDummy.firstInserted, false) // close the loop

		// check equality
		equal, checkError := CheckTrees(treeOriginal, treeDummy)
		if checkError != nil {
			errorText := checkError.Error()
			additionalText := fmt.Sprintf("after %d-th Put, error: %s", index, errorText)
			t.Fatal(newBaseErrorText(index, data, additionalText))
		}
		if !equal {
			t.Fatal(newBaseErrorText(index, data,
				fmt.Sprintf("trees should be equal:\n\toriginal tree: %s\n\tdummy tree: %s\n", treeOriginal, treeDummy)))
		}
	}
	indexNode++ // 9
	alterationFns[indexNode] = func(t *testing.T, treeOriginal *AVLTree[int, *TestData], treeDummy *AVLTree[int, *TestData], index int, data *TestData) {
		oldV, err := treeOriginal.Put(data.Id, data)
		if err != nil {
			t.Fatal(newBaseErrorText(index, data, err.Error()))
		}
		if oldV != treeOriginal.avlTreeConstructorParams.ValueZeroValue {
			t.Fatalf(newBaseErrorText(index, data, fmt.Sprintf("returned value upon nineth Put is not the zero value:\n\t <%v> - <%v>\n", oldV,
				treeOriginal.avlTreeConstructorParams.ValueZeroValue)))
		}

		// set-up tree dummy
		//   .   . 20.
		//   .  /.   .  \
		//   .5  .   .   . 50
		// 2 .10 .   .  30  100
		//1 3.   .   .   42
		nd := nodesDummyTree[index] // 9-th == 42

		treeDummy.cleanNil()
		linkNills(treeDummy, nd, true)
		treeDummy.size++
		node_30 := treeDummy.root // used as temp
		node_30.sizeRight++
		node_30 = node_30.right //50
		node_50 := node_30
		node_50.height++
		node_50.sizeLeft++
		node_30 = node_50.left // 30
		node_30.height++
		node_30.sizeRight++
		node_30.right = nd
		nd.father = node_30

		linkNodes(nd, node_50, true)
		linkNodes(node_30, nd, true) // close the loop
		//if treeDummy.root.left != treeDummy.firstInserted.prevInserted {
		//	t.Fatalf("somehow the last inserted (%v) does not coincide with the expected one (%v)\n", treeDummy.firstInserted.prevInserted, treeDummy.root.left)
		//}
		lastInserted := treeDummy.firstInserted.prevInserted // node "5"
		linkNodes(lastInserted, nd, false)
		linkNodes(nd, treeDummy.firstInserted, false) // close the loop

		// check equality
		equal, checkError := CheckTrees(treeOriginal, treeDummy)
		if checkError != nil {
			errorText := checkError.Error()
			additionalText := fmt.Sprintf("after %d-th Put, error: %s", index, errorText)
			t.Fatal(newBaseErrorText(index, data, additionalText))
		}
		if !equal {
			t.Fatal(newBaseErrorText(index, data,
				fmt.Sprintf("trees should be equal:\n\toriginal tree: %s\n\tdummy tree: %s\n", treeOriginal, treeDummy)))
		}
	}

	indexNode++ // 10 TODO
	alterationFns[indexNode] = func(t *testing.T, treeOriginal *AVLTree[int, *TestData], treeDummy *AVLTree[int, *TestData], index int, data *TestData) {
		oldV, err := treeOriginal.Put(data.Id, data)
		if err != nil {
			t.Fatal(newBaseErrorText(index, data, err.Error()))
		}
		if oldV != treeOriginal.avlTreeConstructorParams.ValueZeroValue {
			t.Fatalf(newBaseErrorText(index, data, fmt.Sprintf("returned value upon tenth Put is not the zero value:\n\t <%v> - <%v>\n", oldV,
				treeOriginal.avlTreeConstructorParams.ValueZeroValue)))
		}

		// set-up tree dummy
		//   .   . 20.
		//   .  /.   .  \
		//   .5  .   .   . 50
		// 2 .10 .   .  37  100
		//1 3.   .   .30 42 // right-left rotation happened here
		nd := nodesDummyTree[index] // 10-th == 37

		treeDummy.cleanNil()
		linkNills(treeDummy, nd, true)
		treeDummy.size++
		node_30 := treeDummy.root // used as temp
		node_30.sizeRight++
		node_50 := node_30.right //50
		node_50.sizeLeft++
		node_30 = node_50.left // 30
		node_30.father = nd
		node_50.left = nd
		nd.father = node_50
		nd.left = node_30
		node_30.height = 0
		node_30.sizeRight = 0
		node_42 := node_30.right // 42
		nd.right = node_42
		node_42.father = nd
		node_30.right = dummyTree._NIL
		nd.height = 1
		nd.sizeLeft = 1
		nd.sizeRight = 1

		linkNodes(nd, node_42, true)
		linkNodes(node_30, nd, true)                         // close the loop
		lastInserted := treeDummy.firstInserted.prevInserted // node "5"
		linkNodes(lastInserted, nd, false)
		linkNodes(nd, treeDummy.firstInserted, false) // close the loop

		// check equality
		equal, checkError := CheckTrees(treeOriginal, treeDummy)
		if checkError != nil {
			errorText := checkError.Error()
			additionalText := fmt.Sprintf("after %d-th Put, error: %s", index, errorText)
			t.Fatal(newBaseErrorText(index, data, additionalText))
		}
		if !equal {
			t.Fatal(newBaseErrorText(index, data,
				fmt.Sprintf("trees should be equal:\n\toriginal tree: %s\n\tdummy tree: %s\n", treeOriginal, treeDummy)))
		}
	}

	// TODO : do the remaining 10 test functions

	//

	var value *TestData
	for i := 0; i < valuesInTotal; i++ {
		value = NewTestDataFilled(values[i], fmt.Sprintf("v_%d", i))
		if alterationFns[i] != nil {
			(alterationFns[i])(t, tree, dummyTree, i, value)
		}
	}
}

//
//
//

func EqualKey[K int | int32 | int64](d1 K, d2 K) bool {
	return d1 == d2
}
func EqualData[V any](d1 *V, d2 *V) bool {
	return d1 == d2
}
func EqualTestData(d1 *TestData, d2 *TestData) bool {
	return EqualData[TestData](d1, d2)
}
func EqualTestDataDeep(d1 *TestData, d2 *TestData) bool {
	return d1 == d2 || ((d1 != nil) && (d2 != nil) && //
		(d1.Id == d2.Id) && (d1.Text == d2.Text))
}

func testEqualityObj[V any](shouldBeEqual bool, actual V, expected V, equalityPredicate func(V, V) bool, additionalErrorText string) error {
	shouldOrShouldNot := ""
	if !shouldBeEqual {
		shouldOrShouldNot = " NOT"
	}
	if shouldBeEqual != equalityPredicate(actual, expected) {
		err := fmt.Errorf("actual value %v should%s be equal to expected value %v; %s\n", actual, shouldOrShouldNot, expected, additionalErrorText)
		if err != nil {
			return err
		}
	}
	return nil
}

func testEqualityPrimitive[V int | int64 | int32](shouldBeEqual bool, actual V, expected V, additionalErrorText string) error {
	shouldOrShouldNot := ""
	if !shouldBeEqual {
		shouldOrShouldNot = " NOT"
	}
	if shouldBeEqual != (actual == expected) {
		err := fmt.Errorf("actual value %d should%s be equal to expected value %d; %s\n", actual, shouldOrShouldNot, expected, additionalErrorText)
		if err != nil {
			return err
		}
	}
	return nil
}
func testLinkEquality[K any, V any](shouldBeEqual bool, n *AVLTNode[K, V], nodeLinkedTo *AVLTNode[K, V], additionalErrorText string) error {
	shouldOrShouldNot := ""
	if !shouldBeEqual {
		shouldOrShouldNot = " NOT"
	}
	if shouldBeEqual != (n == nodeLinkedTo) {
		err := fmt.Errorf("Node should%s be equal ; %s ;\n\tkey node = %v\n\tkey link = %v\n", shouldOrShouldNot, additionalErrorText, n.keyVal.key, nodeLinkedTo.keyVal.key)
		if err != nil {
			return err
		}
	}
	return nil
}
func testNIL[K any, V any](tree *AVLTree[K, V], shouldBeNil bool, n *AVLTNode[K, V], additionalErrorText string) error {
	return testLinkEquality(shouldBeNil, n, tree._NIL, fmt.Sprintf("should be nil; %s", additionalErrorText))
}
func testIsLeaf[K any, V any](tree *AVLTree[K, V], n *AVLTNode[K, V]) error {
	err := testNIL(tree, true, n.left, "therefore, has a left")
	if err != nil {
		return err
	}
	err = testNIL(tree, true, n.right, "therefore, has a right")
	return err
}

func GetNodePath[K any, V any](t *AVLTree[K, V], path []bool) (*AVLTNode[K, V], error) {
	n := t.root
	l := len(path)
	for i := 0; i < l && n != t._NIL; i++ {
		if path[i] {
			n = n.left
		} else {
			n = n.right
		}
	}
	if n == t._NIL { // error
		return nil, fmt.Errorf("node not found with path: %v", path)
	}
	return n, nil
}

func gnp[K any, V any](t *AVLTree[K, V], path []int) (*AVLTNode[K, V], error) {
	p := make([]bool, len(path))
	for i, isLeft := range path {
		p[i] = isLeft == 0
	}
	return GetNodePath(t, p)
}

func DumpTreesForErrorsPrinter[K any, V any](t1 *AVLTree[K, V], t2 *AVLTree[K, V], additionalPreText string, printer func(s string)) {
	printer(additionalPreText)
	printer("\nt1:\n")
	t1.StringInto(true, printer)
	printer("\nt2:\n")
	t2.StringInto(true, printer)
}
func DumpTreesForErrorsBuilder[K any, V any](t1 *AVLTree[K, V], t2 *AVLTree[K, V], additionalPreText string, sb strings.Builder) {
	DumpTreesForErrorsPrinter(t1, t2, additionalPreText, func(s string) { sb.WriteString(s) })
}
func DumpTreesForErrors[K any, V any](t1 *AVLTree[K, V], t2 *AVLTree[K, V], additionalPreText string) string {
	var sb strings.Builder
	DumpTreesForErrorsBuilder(t1, t2, additionalPreText, sb)
	return sb.String()
}
func CheckTrees[K any, V any](t1 *AVLTree[K, V], t2 *AVLTree[K, V]) (bool, *errorTest) {
	if t1 == t2 {
		return true, nil
	}

	if t1.size != t2.size {
		errText := fmt.Sprintf("different sizes: %d and %d\n", t1.size, t2.size)
		return false, ne(errText)
	}

	if t1.IsEmpty() && t2.IsEmpty() {
		return true, nil
	}
	if t1.IsEmpty() != t2.IsEmpty() {
		if t1.IsEmpty() {
			errText := fmt.Sprintf("t1 is empty but t2 is not: t2 has %d nodes", t2.size)
			return false, ne(errText)
		}
		errText := fmt.Sprintf("t1 is not empty but t2 is: t1 it has %d nodes", t2.size)
		return false, ne(errText)
	}
	//fmt.Println("on CheckTrees, checking height")
	if t1.root.height != t2.root.height {
		errText := DumpTreesForErrors(t1, t2, //
			fmt.Sprintf("they have different heights: t1's %d, t2's %d\n", t1.root.height, t2.root.height))
		return false, ne(errText)
	}

	maxHeight := t1.root.height
	pathRun := make([]bool, maxHeight+1)
	for i := 0; i < (int(maxHeight) + 1); i++ {
		pathRun[i] = false
	}

	equal, err := checkTreesEquality(t1, t2, t1.root, t2.root, pathRun, 0)
	if (!equal) || (err != nil) {
		return false, err
	}

	// look for nodes held by "in order / inserted-chronological" pointers BUT without fathers, children
	// (or a father who has them as child, which is granted by the tree equality)

	forEaches := []ForEachMode{
		InOrder,
		ReverseInOrder,
		Queue,
		Stack,
	}

	for _, fe := range forEaches {
		errs1 := []error{}
		errs2 := []error{}
		nodes_count1 := 0
		nodes_count2 := 0

		accumulator := func(currentTree *AVLTree[K, V], isOne bool, forEachM ForEachMode) func(node *AVLTNode[K, V]) {
			io := isOne
			ct := currentTree
			fem := forEachM
			return func(node *AVLTNode[K, V]) {
				isDangling := (ct.size > 1) && (node.father == ct._NIL) && (node.left == ct._NIL) && (node.right == ct._NIL)

				if isDangling {
					var sb strings.Builder

					sb.WriteString("while iterating on tree ")
					if io {
						sb.WriteString("one")
					} else {
						sb.WriteString("two")
					}
					sb.WriteString(" with ForEachMode < ")
					sb.WriteString(fem.String())
					sb.WriteString(" >, current node is dangling; the node:\n\t--")

					if node == ct._NIL {
						sb.WriteString("WAIT ! current node is NIL ! WTF???")
					} else {
						node.toStringTabbed(true, func(s string) { sb.WriteString(s) })
					}

					errCurrent := ne(sb.String())
					if io {
						errs1 = append(errs1, errCurrent)
						nodes_count1++
					} else {
						errs2 = append(errs2, errCurrent)
						nodes_count2++
					}
				}
			}
		}
		t1.forEachNode(fe, accumulator(t1, true, fe))
		t2.forEachNode(fe, accumulator(t2, false, fe))

		var sb strings.Builder
		hasErrors := false
		if len(errs1) > 0 {
			hasErrors = true
			sb.WriteString(fmt.Sprintf("on tree one, on ForEachType %s, %d errors:", fe.String(), nodes_count1))
			for _, e := range errs1 {
				sb.WriteString("\n\t- : ")
				sb.WriteString(e.Error())
			}
			sb.WriteString("\n")
		}
		errs1 = nil
		if len(errs2) > 0 {
			hasErrors = true
			sb.WriteString(fmt.Sprintf("on tree two, on ForEachType %s, %d errors:", fe.String(), nodes_count2))
			for _, e := range errs2 {
				sb.WriteString("\n\t- : ")
				sb.WriteString(e.Error())
			}
			sb.WriteString("\n")
		}
		errs2 = nil
		if hasErrors {
			errText := sb.String()
			sb.Reset()
			return false, ne(errText)
		}

	}
	return true, nil
}

func composeErrorOnCheckTree[K any, V any](t1 *AVLTree[K, V], t2 *AVLTree[K, V], n1 *AVLTNode[K, V], n2 *AVLTNode[K, V], pathRun []bool, depthCurrent int, additionalText string) string {
	var sb strings.Builder
	var branchText string
	if pathRun[depthCurrent] {
		branchText = "left"
	} else {
		branchText = "right"
	}
	sb.WriteString(fmt.Sprintf("\twhile exploring %s branch at depth %d (complete path: %v), an error occour:\n\t", branchText, depthCurrent, pathRun))
	sb.WriteString(additionalText)
	sb.WriteString("\n\tdumping nodes")
	printer := func(s string) { sb.WriteString(s) }
	n1.toStringTabbed(true, printer)
	n2.toStringTabbed(true, printer)
	printer = nil // clear the memory
	sb.WriteString("\n\tdumping trees")

	return DumpTreesForErrors(t1, t2, sb.String())
}

/*
path true == left, false == right
*/
func checkTreesEquality[K any, V any](t1 *AVLTree[K, V], t2 *AVLTree[K, V], n1 *AVLTNode[K, V], n2 *AVLTNode[K, V], pathRun []bool, depthCurrent int) (bool, *errorTest) {

	if n1 == nil && n2 == nil || (n1 == t1._NIL && n2 == t2._NIL) {
		return true, nil
	}

	if n1 == nil {
		// ERROR: SHOULD NOT BE NIL
		var nullity string = "null"
		errText := composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, //
			fmt.Sprintf("node of first tree is %s (the second one didn't)", nullity))
		return false, ne(errText)
	}
	if n2 == nil {
		// ERROR: SHOULD NOT BE NIL
		var nullity string = "null"
		errText := composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, //
			fmt.Sprintf("node of second tree is %s (the first one didn't)", nullity))
		return false, ne(errText)
	}

	integers := []string{"height", "size left", "size right"}
	int_1 := []int64{n1.height, n1.sizeLeft, n1.sizeRight}
	int_2 := []int64{n2.height, n2.sizeLeft, n2.sizeRight}
	i := 0
	l := len(integers)
	for ; i < l; i++ {
		if int_1[i] != int_2[i] {
			errText := composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, //
				fmt.Sprintf("checking integers; %s comparison failed: node1 ones= %d, node2 ones= %d", integers[i], int_1[i], int_2[i]))
			return false, ne(errText)
		}
	}

	var comp1, comp2 int

	// comparing node's keys

	/*
		keep track of the "pointers to pointers" in order to generalize the handling of
		the 4 nodes: the current node itself, its father, its left and its right.
		All of them are "pointers to nodes" whose node's key should be equal amont "1" and "2"
	*/
	pointersNode1 := []**AVLTNode[K, V]{&n1, &(n1.father), &(n1.left), &(n1.right)}
	pointersNode2 := []**AVLTNode[K, V]{&n2, &(n2.father), &(n2.left), &(n2.right)}
	pointerName := []string{"", "'s father", "'s left", "'s right"}
	i = 0
	l = len(pointerName)
	for ; i < l; i++ {
		pointer1 := pointersNode1[i]
		pointer2 := pointersNode2[i]
		nameNode := pointerName[i]

		node1 := *pointer1 // this node could be any among "pointerName"
		node2 := *pointer2

		// checking : NIL-ity

		if (node1 == t1._NIL) != (node2 == t2._NIL) { // the "XOR" ("^") does not exists, "!=" is equivalent
			errText := composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, //
				fmt.Sprintf("while comparing nodes%s NIL-ity, they are different: the nil-comparison results in < %t > for 1 and in < %t > for 2\n\t the checked node 1: %v\n\t the checked node 2: %v\n", //
					nameNode, (node1 == t1._NIL), (node2 == t2._NIL), node1, node2))
			return false, ne(errText)
		}
		// if both nodes are NOT "NIL", then they will be checked in the for loop below

	}

	// NOTE: those checks are shifted outside the for loop above because:
	// -) the father should already be checked (ERROR: IT WON'T IF N1 & N2 ARE THE ROOTS! HOW TO DEAL WITH THIS SITUATION [== "those roots"] ?)
	// -) the children (left & right) will be checked by the recursion -> no need to check tem __twice__ ["thrice", actually, due to the "father" check]
	// - -) "checking children" and "at first checking the nodes and THEN shifting the perspective onto the childer" are pratically the same thing:
	// .....the latter is just NOT reduntant
	whoseKeys := []string{"'s self", "nextInOrder", "prevInOrder", "nextInserted", "prevInserted"}
	keyOwners1 := []*AVLTNode[K, V]{n1, n1.nextInOrder, n1.prevInOrder, n1.nextInserted, n1.prevInserted}
	keyOwners2 := []*AVLTNode[K, V]{n2, n2.nextInOrder, n2.prevInOrder, n2.nextInserted, n2.prevInserted}
	i = 0
	l = len(whoseKeys)
	for ; i < l; i++ {
		keyOwnername := whoseKeys[i]
		node1 := keyOwners1[i]
		node2 := keyOwners2[i]

		comp1 = int(t1.avlTreeConstructorParams.Comparator(node1.keyVal.key, node2.keyVal.key))
		comp2 = int(t2.avlTreeConstructorParams.Comparator(node1.keyVal.key, node2.keyVal.key))

		if comp1 != 0 {
			errText := composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, //
				fmt.Sprintf("while comparing nodes%s key with tree 1 comparator, the comparison should be 0, but is: %d", keyOwnername, comp1))
			return false, ne(errText)
		}
		if comp2 != 0 {
			errText := composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, //
				fmt.Sprintf("while comparing nodes%s key with tree 2 comparator, the comparison should be 0, but is: %d", keyOwnername, comp2))
			return false, ne(errText)
		}
	}

	//fmt.Printf("nodes << n1:%v ; n2:%v >>, recursion on children: left\n", n1.keyVal.key, n2.keyVal.key)
	pathRun[depthCurrent] = true
	equal, err := checkTreesEquality(t1, t2, n1.left, n2.left, pathRun, depthCurrent+1)
	if (!equal) || (err != nil) {
		errText := composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, err.Error())

		return false, ne(errText)
	}
	pathRun[depthCurrent] = false
	//fmt.Printf("nodes << n1:%v ; n2:%v >>, recursion on children: right\n", n1.keyVal.key, n2.keyVal.key)
	equal, err = checkTreesEquality(t1, t2, n1.right, n2.right, pathRun, depthCurrent+1)
	if (!equal) || (err != nil) {
		errText := composeErrorOnCheckTree(t1, t2, n1, n2, pathRun, depthCurrent, err.Error())
		return false, ne(errText)
	}
	/*
		var sb strings.Builder
		for i := 0; i < depthCurrent; i++ {
			sb.WriteRune('\t')
		}
		tabsText := sb.String()
		fmt.Printf("\n\n%sunder the sea\n%sat depth %d you'll see\n%si'll return true\n%ssince I've been through\n%sa code bug-freeeee\n\n",
			tabsText, tabsText, depthCurrent, tabsText, tabsText, tabsText)
	*/
	return true, nil
}

/**
 * set the first node as the "previous" of the second one and the second node as the "next"
 * node of the first one.
 */
func linkNodes[K any, V any](n1 *AVLTNode[K, V], n2 *AVLTNode[K, V], isInOrder bool) {
	if isInOrder {
		n1.nextInOrder = n2
		n2.prevInOrder = n1
	} else {
		n1.nextInserted = n2
		n2.prevInserted = n1
	}
}

func linkNills[K any, V any](t *AVLTree[K, V], n *AVLTNode[K, V], shouldClearOrderings bool) {
	n.father = t._NIL
	n.left = t._NIL
	n.right = t._NIL
	n.sizeLeft = 0
	n.sizeRight = 0
	n.height = 0
	if !shouldClearOrderings {
		return
	}
	n.nextInOrder = t._NIL
	n.prevInOrder = t._NIL
	n.nextInserted = t._NIL
	n.prevInserted = t._NIL
}
