package main

import (
	"fmt"
	"math"
)

//	Посетитель - поведенческий паттерн, позволяющий добавлять поведение в структуру без ее изменения.
//	Минусы: требует добавления доп метода, структур, интерфейса
//	плюсы: Не затрагивает исходный код объекта, нет риска сломать исходный код новым поведением

type shape interface {
	getType() string
	accept(visitor)
}

type square struct {
	Side float32
}

func (s *square) accept(v visitor) {
	v.visitForSquare(s)
}

func (s *square) getType() string {
	return "Квадрат"
}

type circle struct {
	Radius float32
}

func (c *circle) accept(v visitor) {
	v.visitForCircle(c)
}

func (c *circle) getType() string {
	return "Круг"
}

type visitor interface {
	visitForSquare(*square)
	visitForCircle(*circle)
}

type areaCalculator struct {
	area int
}

func (a *areaCalculator) visitForSquare(s *square) {
	fmt.Println("S =", s.Side*s.Side)
}

func (a *areaCalculator) visitForCircle(s *circle) {
	fmt.Println("S =", s.Radius*s.Radius*math.Pi)
}

func main() {
	square := &square{2}
	circle := &circle{3}

	areaCalculator := &areaCalculator{}

	square.accept(areaCalculator)
	circle.accept(areaCalculator)
}
