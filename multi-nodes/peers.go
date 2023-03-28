package multi_nodes

//peerPicker的pickPeer()方法用于根据传入的key选择相应节点PeerGetting
//PeerGetter的Get()方法用于从对应group查找缓存值。PeerGetter就对应于上述流程中的HTTP客户端

//选择节点

type PeerPicker interface {
	PickPeer(key string) (peer PeerGetter, ok bool)
}

// 根据group和key选择group核心操作数据
type PeerGetter interface {
	Get(group string, key string) ([]byte, error)
}
