package treex

// TreeNode 树节点接口,必须提供id、Pid、children
type TreeNode interface {
	GetId() int64
	GetPid() int64
	GetChildren() *[]TreeNode
}

// BuildTree 构建树形结构
//
//	type Menu struct {
//	    Id       int64      `json:"id"`
//		Pid      int64      `json:"pid"`
//		Title    string     `json:"title"`
//		Children []TreeNode `json:"children"`
//	}
//
// // 实现TreeNode接口
// func (m *Menu) GetId() int64 { return m.Id }
// func (m *Menu) GetPid() int64 { return m.Pid }
//
//	func (m *Menu) GetChildren() *[]TreeNode {
//	    return (*[]TreeNode)(&m.Children)
//	}
//
//	menus := []Menu{
//	    {Id: 1, Pid: 0, Title: "系统管理"},
//	    {Id: 2, Pid: 1, Title: "用户管理"},
//	    {Id: 3, Pid: 1, Title: "角色管理"},
//	    {Id: 4, Pid: 2, Title: "创建用户"},
//	}
//
// tree := BuildTree(menus)
func BuildTree[T TreeNode](items []T) []TreeNode {
	nodeMap := make(map[int64]TreeNode)

	for i := range items {
		node := items[i]
		*node.GetChildren() = []TreeNode{}
		nodeMap[node.GetId()] = node
	}

	var roots []TreeNode

	for i := range items {
		node := items[i]
		if node.GetPid() == 0 || nodeMap[node.GetPid()] == nil {
			roots = append(roots, node)
			continue
		}

		parent := nodeMap[node.GetPid()]
		*parent.GetChildren() = append(*parent.GetChildren(), node)
	}

	return roots
}

// GetByPid 根据Pid获取所有节点
func GetByPid[T TreeNode](items []T, pid int64) []TreeNode {
	var result []TreeNode

	for i := range items {
		if items[i].GetPid() == pid {
			result = append(result, items[i])
		}
	}

	return result
}
