package errlog

import (
	"bufio"
	"hash/fnv"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func hashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func Write(errText string) (string, string) {
	file, errT := os.OpenFile("err.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if errT != nil {
		return "", errT.Error()
	}
	defer file.Close()
	write := bufio.NewWriter(file)

	// 获取错误字符串的哈希值
	hash := hashString(errText)

	// 创建随机数器
	randSource := rand.NewSource(int64(hash) + time.Now().UnixNano())
	rng := rand.New(randSource)

	// 生成8位随机数id
	id := 10000000 + rng.Intn(90000000)

	// 将id和时间转化为string类型
	strId := strconv.Itoa(id)
	strTime := time.Now().Format("2006-01-02 15:04:05")

	// 拼接字符串
	strs := []string{"[", strTime, "]", "[", strId, "]: ", errText, "\n"}
	result := strings.Join(strs, "")

	// 写入字符串
	if _, err := write.WriteString(result); err != nil {
		return "", err.Error()
	}

	// 确保字符串已经写入而不是缓存在内存里
	if err := write.Flush(); err != nil {
		return "", err.Error()
	}

	return strId, ""
}
