package base_demo

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSplit(t *testing.T) {
	got := Split("a:b:c", ":")
	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected:%v, got:%v", want, got)
	}
}

func TestSplitWithComplexSep(t *testing.T) {
	got := Split("abcd", "bc")
	want := []string{"a", "d"}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("expected:%v, got:%v", want, got)
	}
}

func TestTimeConsuming(t *testing.T) {
	if testing.Short() {
		t.Skip("short 模式下会跳过该测试用例")
	}
}

func TestSplitAll(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		input string
		sep string
		want []string
	}{
		{"base case", "a:b:c", ":", []string{"a", "b", "c"}},
		{"wrong sep", "a:b:c", ",", []string{"a:b:c"}},
		{"more sep", "abcd", "bc", []string{"a", "d"}},
		{"leading sep", "沙河有沙又有河", "沙", []string{"", "河有", "又有河"}},
	}
	// 遍历测试用例
	for _, tt := range tests {
		tt := tt // 注意这里重新声明tt变量（避免多个goroutine中使用了相同的变量）
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()  // 将每个测试用例标记为能够彼此并行运行
			got := Split(tt.input, tt.sep)
			assert.Equal(t, got,tt.want)
			//if !reflect.DeepEqual(got, tt.want) {
			//	t.Errorf("expected:%#v, got:%#v", tt.want, got)
			//}
		})
	}
}

//func BenchmarkSplit(b *testing.B){
//	time.Sleep(1 * time.Second)
//	b.ResetTimer() // 重置计时器
//	for i:= 0; i< b.N; i++ {
//		Split("沙河有沙又有河", "沙")
//	}
//}
//
//func benchmarkFib(b *testing.B, n int) {
//	for i := 0; i < b.N; i++ {
//		Fib(n)
//	}
//}
//
//func BenchmarkFib1(b *testing.B)  { benchmarkFib(b, 1) }
//func BenchmarkFib2(b *testing.B)  { benchmarkFib(b, 2) }
//func BenchmarkFib3(b *testing.B)  { benchmarkFib(b, 3) }
//func BenchmarkFib10(b *testing.B) { benchmarkFib(b, 10) }
//func BenchmarkFib20(b *testing.B) { benchmarkFib(b, 20) }
//func BenchmarkFib40(b *testing.B) { benchmarkFib(b, 40) }
//
//func BenchmarkSplitParallel(b *testing.B) {
//	// b.SetParallelism(1) // 设置使用的CPU数
//	b.RunParallel(func(pb *testing.PB) {
//		for pb.Next() {
//			Split("沙河有沙又有河", "沙")
//		}
//	})
//}

//func TestMain(m *testing.M) {
//	fmt.Println("write setup code here...") // 测试之前的做一些设置
//	// 如果 TestMain 使用了 flags，这里应该加上flag.Parse()
//	retCode := m.Run()                         // 执行测试
//	fmt.Println("write teardown code here...") // 测试之后做一些拆卸工作
//	os.Exit(retCode)                           // 退出测试
//}
