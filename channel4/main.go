package main

import (
	"fmt"
	"math/rand"
)

func main()  {
	chs := [3]chan int{
		make(chan int, 3),
		make(chan int, 3),
		make(chan int, 3),
	}

	index1 := rand.Intn(3)  // 随机生成0-2之间的数字
	fmt.Printf("随机索引/数值: %d\n", index1)
	chs[index1] <- rand.Int()   // 向通道发送随机数字

	index2 := rand.Intn(3)
	fmt.Printf("随机索引/数值: %d\n", index2)
	chs[index2] <- rand.Int()

	index3 := rand.Intn(3)
	fmt.Printf("随机索引/数值: %d\n", index3)
	chs[index3] <- rand.Int()

	// 哪一个通道中有值，哪个对应的分支就会被执行
	for i := 0; i < 4; i++ {
		select {
		case num, ok := <- chs[0]:
			if !ok {
				fmt.Println("1break")
				break
			}
			fmt.Println("第一个条件分支被选中: chs[0]=>", num)
		case num, ok := <- chs[1]:
			if !ok {
				fmt.Println("2break")
				break
			}
			fmt.Println("第二个条件分支被选中: chs[1]=>", num)
		case num, ok := <- chs[2]:
			if !ok {
				fmt.Println("3break")
				break
			}
			fmt.Println("第三个条件分支被选中: chs[2]=>", num)
		default:
			fmt.Println("没有分支被选中")
		}
	}
}