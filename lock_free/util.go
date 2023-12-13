package lock_free

func bigSplit(node *Node) ([]int, []interface{}) {
	newKeys := make([]int, 0)
	newNodes := make([]interface{}, 0)
	newNodeCount := ((node.NumKeys + minOrder - 1) / minOrder) - 1
	assert(newNodeCount > 0)
	currIndex := minOrder
	for nodeNum := 1; nodeNum < newNodeCount; nodeNum++ {
		newNode, _ := makeNode()
		newNode.Parent = node.Parent
		newNode.NumKeys = minOrder - 1
		for j := 0; j < minOrder; j++ {
			newNode.Keys[j] = node.Keys[j+currIndex]
		}
		for j := 0; j < minOrder; j++ {
			newNode.Pointers[j] = node.Pointers[j+minOrder]
		}
		newKeys = append(newKeys, newNode.Keys[nodeNum*minOrder])
		newNodes = append(newNodes, newNode)
		currIndex += minOrder
	}
	node.NumKeys = minOrder - 1
	node.Keys = node.Keys[:node.NumKeys]
	node.Pointers = node.Pointers[:node.NumKeys+1]
	return newKeys, newNodes
}

func assert(condition bool) {
	if !condition {
		panic("assertion failed")
	}
}
