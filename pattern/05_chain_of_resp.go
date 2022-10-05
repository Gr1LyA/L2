package main

import "fmt"

// Цепочка вызовов - поведенческий паттерн, суть заключается в передаче запроса по цепочке потенциальных обработчиков, пока один из них, или каждый, не обработает запрос.
/*
	Плюсы: Каждый обработчик ответсвенен за свою область
			клиент не зависит от обработчика
			Выстраивание цепи обработки на свое усмотрение

	Минусы: Создание доп объектов, усложнение
*/

type car struct {
	Body        bool
	Engine      bool
	Wheels      bool
	Electronics bool
}

// Интерфейс обработчиков
type department interface {
	execute(*car)
	setNext(department)
}

// Обработчик устанавливающий кузов машины
type bodyWorkshop struct {
	next department
}

func (w *bodyWorkshop) execute(c *car) {
	if c.Body {
		fmt.Println("body already installed")
	} else {
		fmt.Println("body installed")
		c.Body = true
	}
	if w.next != nil {
		w.next.execute(c)
	}
}

func (w *bodyWorkshop) setNext(next department) {
	w.next = next
}

// Обработчик устанавливающий двигатель машины
type engineWorkshop struct {
	next department
}

func (w *engineWorkshop) execute(c *car) {
	if c.Engine {
		fmt.Println("engine already installed")
	} else {
		fmt.Println("engine installed")
		c.Engine = true
	}
	if w.next != nil {
		w.next.execute(c)
	}
}

func (w *engineWorkshop) setNext(next department) {
	w.next = next
}

// Обработчик устанавливающий колеса машины
type wheelsWorkshop struct {
	next department
}

func (w *wheelsWorkshop) setNext(next department) {
	w.next = next
}

func (w *wheelsWorkshop) execute(c *car) {
	if c.Wheels {
		fmt.Println("wheels already installed")
	} else {
		fmt.Println("wheels installed")
		c.Wheels = true
	}
	if w.next != nil {
		w.next.execute(c)
	}
}

// Обработчик устанавливающий электронику машины
type electronicsWorkshop struct {
	next department
}

func (w *electronicsWorkshop) execute(c *car) {
	if c.Electronics {
		fmt.Println("electronics already installed")
	} else {
		fmt.Println("electronics installed")
		c.Electronics = true
	}
	if w.next != nil {
		w.next.execute(c)
	}
}

func (w *electronicsWorkshop) setNext(next department) {
	w.next = next
}

func main() {
	workShopElectric := &electronicsWorkshop{}
	workShopEngine := &engineWorkshop{}
	workShopWheels := &wheelsWorkshop{}
	workShopBody := &bodyWorkshop{}

	workShopBody.setNext(workShopWheels)
	workShopWheels.setNext(workShopEngine)
	workShopEngine.setNext(workShopElectric)

	workShopBody.execute(&car{})
}
