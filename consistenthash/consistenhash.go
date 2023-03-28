package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

//解决相同的key可以访问相同的服务器
// Hash maps bytes to uint32
//定义了函数类型Hash,采取依赖注入的方式，允许用于替代换成自定义的Hash函数，也方便测试时替换，默认为crc32.ChecksumIEEE算法

type Hash func(data []byte) uint32

// Map是一致性hash算法的主要数据结构

type Map struct {
	hash     Hash           //hash函数
	replicas int            //虚拟节点倍数
	keys     []int          //Sorts 哈希环
	hashMap  map[int]string //虚拟节点与真实节点的映射表hashMap,键是虚拟节点的哈希值，值是真实节点的名称
}

//构造函数New()允许自定义虚拟节点倍数和Hash函数

func New(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//接下来，试下添加真实节点/机器的Add()方法
//Add函数允许传入0或者多个真实节点的名称
//对每个真实节点key，对应创建m.replicas个虚拟节点，虚拟节点的名称是：strconv.Itoa(i)+key,即通过添加编号的方式区分不同虚拟节点
//使用m.hash()计算虚拟节点的哈希值，使用append(m,keys,hash)添加到环上
//最后一步，环上的哈希值排序。

func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	//Golang提供了sort.Ints()函数对int数组切片序列从小到大升序排序
	sort.Ints(m.keys)

}

//最后一步，试下选择节点的Get()方法,查找存放的数据放到那个节点里

func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	// Search函数采用二分找到切片中的元素
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
