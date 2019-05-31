package model

import "fmt"

func main() {

	c := make(chan int, 1)
	close(c)
	//c <- 1
	//for {
		select {
			case <-c:
				fmt.Println("获取到")
		default:
			fmt.Println("默认")
		}
	//}
}
