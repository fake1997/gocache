package go_cache

// PeerPicker 通过传入的key选择响应的节点
type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// PeerGetter 从对应group查询缓存
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
