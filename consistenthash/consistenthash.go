// 一致性hash
// 场景:缓存数据库扩容,缩容时,避免数据大量迁移,重新计算hash
// 明白几个基本概念:哈希环,节点,虚拟节点
package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type Map struct {
	// 计算hash的函数
	hash Hash
	// 虚拟节点的个数
	replicas int
	// 所有的虚拟节点的排序(哈希环是抽象的环,代码上面其实就是排了序的切片)
	keys []int
	//虚拟节点与真实节点的映射map
	hashMap map[int]string
}

func New(replicas int, fn Hash) *Map {
	// 初始化虚拟节点数量及hash函数
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}

	// 没有设置hash函数的话,默认选择crc32.ChecksumIEEE
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

func (m *Map) IsEmpty() bool {
	return len(m.keys) == 0
}

// 添加多个节点进hash
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		// 根据m.replices 添加虚拟节点进哈希环
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	sort.Ints(m.keys)
}

// 查找某条数据属于哪个节点
// key:缓存数据的key值
// 返回值:返回真实节点的名称
func (m *Map) Get(key string) string {
	if m.IsEmpty() {
		return ""
	}
	// 计算缓存数据key对应的hash
	hash := int(m.hash([]byte(key)))
	// 通过二分法查找m.keys中大于等于hash值的最小值(顺时针找最小的节点)
	index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	// 如果hash值大于m.keys中的所有值,则该key属于第一个虚拟节点(这就是哈希环为什么叫'环'了)
	if index == len(m.keys) {
		index = 0
	}
	return m.hashMap[m.keys[index]]
}
