Что выведет программа? Объяснить вывод программы. Объяснить как работают defer’ы и их порядок вызовов.

```go
package main

import (
	"fmt"
)


func test() (x int) {
	defer func() {
		x++
	}()
	x = 1
	return
}


func anotherTest() int {
	var x int
	defer func() {
		x++
	}()
	x = 1
	return x
}


func main() {
	fmt.Println(test())
	fmt.Println(anotherTest())
}
```

Ответ:
```
2
1

defer хранятся в стеке типа LIFO и выполняются в соответсвующем порядке.
defer выполняется прямо перед return.

В случае anotherTest() int, return можно разбить на след этапы:
Возвращаемым значение обявим result(в данном случае оно анонимно)
result = x
x++
return (возвращаем result, defer инкрементировал x а не result)

В случае test() (x int), return можно разбить на след этапы:
Возвращаемое значение x (оно не анонимно как в пред случае)
x++
return (возвращаем x, defer инкрементировал именно возвращаемое значение)
