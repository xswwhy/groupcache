# groupcache
原项目地址 https://github.com/golang/groupcache
本项目为groupcache少部分代码的学习与解读

# 项目目录结构
* consistenthash: 一致性hash的使用案例
* lru: 最少使用淘汰机制
* singleflight: 单航班

# 概述  
groupcache是用go语言写的类似memcached的一个分布式缓存缓存数据库  
本项目只解读了groupcache少部分功能,代码也很少
一致性hash,lru,singleflight网上都有详细文章,code中配有详细注释,就不做过多解释了
