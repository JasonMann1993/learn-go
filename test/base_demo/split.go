package base_demo

import "strings"

func Split(s, sep string) (result []string) {
	//result = make([]string, 0, strings.Count(s, sep)+1)
	i := strings.Index(s, sep)
	l := len(sep)
	for i > -1 {
		result = append(result, s[:i])
		s = s[i+l:]
		i = strings.Index(s, sep)
	}
	result = append(result, s)
	return
}

// Fib 是一个计算第n个斐波那契数的函数
func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) + Fib(n-2)
}
