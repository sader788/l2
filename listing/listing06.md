Что выведет программа? Объяснить вывод программы. Рассказать про внутреннее устройство слайсов и что происходит при передачи их в качестве аргументов функции.

```go
package main

import (
	"fmt"
)

func main() {
	var s = []string{"1", "2", "3"}
	modifySlice(s)
	fmt.Println(s)
}

func modifySlice(i []string) {
	i[0] = "3"
	i = append(i, "4")
	i[1] = "5"
	i = append(i, "6")
}
```

Ответ:
```
Выведет [3, 2, 3], потому что когда мы передаем слайс мы копируем объект.
Внутри слайса есть указатель на кусок памяти, при операции индексации мы работаем с ним.
Поэтому когда мы изменяем нулевой элемент он меняется.

Дальше append затерает наш слайс новым, и мы уже работаем с другим слайсом.
```