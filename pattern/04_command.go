package main

import "fmt"

//	Команда - поведенческий паттерн, в котором запросы или операции являются отдельными объектами
/*
	Плюсы: нет связью между отправителем и исполнителем комманды, удобно манипулировать операциями:
		повтор, выстраивание в очередь, отложенный запуск запросов
	Минусы: Усложнение кода из-за добавления доп интерфейса и структур
*/

type tv struct {
	On bool
}

// Включение телевизора
func (t *tv) start() {
	t.On = true
}

// Выключение телевизра
func (t *tv) shutdown() {
	t.On = false
}

// Команда, выполняющаяся по нажатию кнопки
type command interface {
	execute()
}

type shutdownCmd struct {
	tv *tv
}

// Конструктор комманды, аргументом передается то над чем будет проводиться действие
func newShutdownCmd(t *tv) *shutdownCmd {
	return &shutdownCmd{t}
}

// Выключает телевизор
func (s *shutdownCmd) execute() {
	s.tv.shutdown()
}

type startCmd struct {
	tv *tv
}

// Конструктор комманды, аргументом передается то над чем будет проводиться действие
func newStartCmd(t *tv) *startCmd {
	return &startCmd{t}
}

// Включает телевизор
func (s *startCmd) execute() {
	s.tv.start()
}

type button struct {
	cmd command
}

// Коснтруктор кнопки, аргументом передается комманда ,которая будет выполняться по нажатию
func newButton(cmd command) *button {
	return &button{cmd: cmd}
}

func (b *button) press() {
	b.cmd.execute()
}

func main() {
	t := &tv{}

	// Инициализируем комманды включения и выключения
	startCmd := newStartCmd(t)
	shutdownCmd := newShutdownCmd(t)

	// Инициализируем кнопки и передаем в них нужные комманды
	startbutton := newButton(startCmd)
	shutdownbutton := newButton(shutdownCmd)

	fmt.Println(t.On)
	startbutton.press()
	fmt.Println(t.On)
	shutdownbutton.press()
	fmt.Println(t.On)

}
