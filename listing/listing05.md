Что выведет программа? Объяснить вывод программы.

```go
package main

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	{
		// do something
	}
	return nil
}

func main() {
	var err error
	err = test()
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
```

Ответ:
```
error

Ситуация аналогична listing03. 
test возвращет ссылку на структуру которая равна nil, но возвращаемое значение приравнивается к объекту, реализующиму интерфейс error. 
Таким образом оператор = не приравнивает сам err к возвращаемому значение, а приравнивает поинтер указателя на значение объекта err. 
То есть равен nil будет только указатель на занчение а не сам err
