package metrics

import (
	"fmt"
	"hash/fnv"
	"sort"
	"strconv"
	"strings"
)

var orderMap = map[int]string{
	0: "B",
	1: "K",
	2: "M",
	3: "G",
}

func getMemoryString(mem float64) string {
	total := mem
	order := 0
	for total > 1000 {
		total = total / 1000
		order++
	}
	return fmt.Sprintf("%.2f%s", total, orderMap[order])
}

func keysFromMap(m map[string]string) []string {
	a := make([]string, len(m))
	i := 0
	for k := range m {
		a[i] = k
		i++
	}
	return a
}

func mapHash(m map[string]string) string {
	keys := keysFromMap(m)
	sort.Strings(keys)
	s := make([]string, len(m)*2)
	i := 0
	for _, k := range keys {
		s[i] = k
		i++
		s[i] = m[k]
		i++
	}
	h1 := fnv.New64a()
	h1.Write([]byte(strings.Join(s, "")))
	h := h1.Sum64()
	str := strconv.FormatUint(h, 10)
	return str
}
