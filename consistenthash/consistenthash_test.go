package consistenthash

import (
	"fmt"
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	// 创建虚拟节点数为3  hash函数为string->uint32的简单装换  key必须为数字
	hash := New(3, func(key []byte) uint32 {
		i, err := strconv.Atoi(string(key))
		if err != nil {
			panic(err)
		}
		return uint32(i)
	})

	// 添加三个节点
	// 添加完之后:
	// hash.keys [2, 4 ,6, 12, 14, 16, 22, 24, 26]
	// hash.hashMap{2:2, 4:4, 6:6, 12:2, 14:4, 16:6, 22:2, 24:4, 26:6}
	// 对hash.keys  hash.hashMap有困惑的,再仔细看看Add函数中对多个虚拟节点的处理方式 strconv.Itoa(i) + key
	hash.Add("6", "4", "2")

	// 测试数据
	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("查找%s的真实节点错误", k)
		}
	}

	// 添加新的节点
	// 添加完之后:
	// hash.keys [2, 4 ,6, 8, 12, 14, 16, 18, 22, 24, 26, 28]
	// hash.hashMap{2:2, 4:4, 6:6, 8:8, 12:2, 14:4, 16:6, 18:8, 22:2, 24:4, 26:6, 28:8}
	// 对hash.keys  hash.hashMap有困惑的,再仔细看看Add函数中对多个虚拟节点的处理方式 strconv.Itoa(i) + key
	hash.Add("8")

	// 修改测试数据
	testCases["27"] = "8" // 添加了新的节点后,27不再属于2  而是属于8
	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("查找%s的真实节点错误", k)
		}
	}
}

func TestConsistency(t *testing.T) {
	hash1 := New(1, nil)
	hash2 := New(1, nil)

	hash1.Add("Bill", "Bob", "Bonny")
	hash2.Add("Bob", "Bonny", "Bill")

	if hash1.Get("Ben") != hash2.Get("Ben") {
		t.Errorf("相同的key应该在同一节点上")
	}

	// 这里源码有问题,hash2新添加了3个节点, 相同的key在hash1 hash2是有可能属于不同节点的
	hash2.Add("Becky", "Ben", "Bobby")
	if hash1.Get("Ben") != hash2.Get("Ben") ||
		hash1.Get("Bob") != hash2.Get("Bob") ||
		hash1.Get("Bonny") != hash2.Get("Bonny") {
		t.Errorf("Direct matches should always return the same entry")
	}

}

// 压力测试,有兴趣的自行了解哟
func BenchmarkGet8(b *testing.B)   { benchmarkGet(b, 8) }
func BenchmarkGet32(b *testing.B)  { benchmarkGet(b, 32) }
func BenchmarkGet128(b *testing.B) { benchmarkGet(b, 128) }
func BenchmarkGet512(b *testing.B) { benchmarkGet(b, 512) }

func benchmarkGet(b *testing.B, shards int) {

	hash := New(50, nil)

	var buckets []string
	for i := 0; i < shards; i++ {
		buckets = append(buckets, fmt.Sprintf("shard-%d", i))
	}

	hash.Add(buckets...)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		hash.Get(buckets[i&(shards-1)])
	}
}
