package main

import "fmt"

//	Стратегия - это поведенческий паттерн проектирования, который определяет семейство схожих алгоритмов и помещает каждый из них в собственный класс,
//	после чего алгоритмы можно взаимозаменять прямо во время исполнения программы.

// Плюсы:
// 	- Возможность замены алгоритмов в рантайме
// 	- Отделение алгоритмов от остальной логики, сокрытие самих алгоритмов
// Минусы:
// 	- Усложнение кода, засчет введения дополнительных объектов

type number struct {
	num int
	act action
}

func (s *number) setAction(a action) {
	s.act = a
}

func (s *number) runAction() {
	s.act.actionTo(s)
}

type action interface {
	actionTo(*number)
}

type square struct{}

func (s square) actionTo(n *number) {
	n.num *= n.num
}

type inkrement struct{}

func (s inkrement) actionTo(n *number) {
	n.num++
}

type dekrement struct{}

func (s dekrement) actionTo(n *number) {
	n.num--
}

func main() {
	n := &number{num: 1}

	fmt.Println(n.num)

	n.setAction(inkrement{})
	n.runAction()
	fmt.Println(n.num)

	n.setAction(square{})
	n.runAction()
	fmt.Println(n.num)

	n.setAction(dekrement{})
	n.runAction()
	fmt.Println(n.num)
}
