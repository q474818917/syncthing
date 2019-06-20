package train

import "fmt"

/**
len查看当前chan的size,每次range就会减少1
 */
func main() {

	c := make(chan int, 1)
	//c <- 1
	//for {
		select {
			case <-c:
		default:
			fmt.Println("默认")
		}
	//}

	c1 := make(chan int, 2)
	defer close(c1)

	c1 <- 1
	fmt.Println(len(c1))

	for cc := range c1 {
		fmt.Println(cc, "-", len(c1))
	}
}
