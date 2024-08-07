package sync

// TreeNode 树节点
type TreeNode[K comparable] interface {
	GetID() K
	GetPID() K
	AddChild(TreeNode[K])
	GetChild() []TreeNode[K]
}

// InitTree 初始化树
func InitTree[K comparable](root TreeNode[K], children []TreeNode[K]) {
	childMap := make(map[K]TreeNode[K])
	for _, child := range children {
		childMap[child.GetID()] = child
	}
	initTree(root, childMap)
}

func initTree[K comparable](root TreeNode[K], children map[K]TreeNode[K]) {
	// 先找出 root 的子节点
	for _, child := range children {
		// 正常
		if child.GetPID() == root.GetID() {
			root.AddChild(child)
			delete(children, child.GetID())
		}
	}
	// 递归
	for _, child := range root.GetChild() {
		initTree(child, children)
	}
}
