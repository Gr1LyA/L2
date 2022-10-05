package main

import (
	"fmt"
)

// Фабричный метод — порождающий паттерн проектирования, решает проблему создания различных продуктов,
// без указания конкретных классов продуктов.

/*
	ипользуется когда заране неизвестны типы и зависимости с которыми прийдется работать
	плюсы: Не нужно выстраивать бизнес логику отталкиваясь от того с каким типом прийдется работать

	как пример можно приветси кросплатформенную программу которая будет строить gui исходя из стиля ОС
*/

// Интерфейс продукта фабрики
type carI interface {
	getInfo() string
}

// сам класс машин
type car struct {
	color     string
	hp        int
	clearance int
}

func (c *car) getInfo() string {
	return fmt.Sprintf("color: %s, hp: %d, clearance: %d", c.color, c.hp, c.clearance)
}

type policeCar struct {
	car
}

func newPoliceCar() carI {
	return &policeCar{
		car{
			color:     "dark blue",
			hp:        220,
			clearance: 8,
		},
	}
}

type sportCar struct {
	car
}

func newSportCar() carI {
	return &sportCar{
		car{
			color:     "bloody red",
			hp:        500,
			clearance: 4,
		},
	}
}

type defaultCar struct {
	car
}

func newDefaultCar() carI {
	return &defaultCar{
		car{
			color:     "white",
			hp:        150,
			clearance: 12,
		},
	}
}

// Производство машин
func factory(typeCar string) carI {
	switch typeCar {
	case "sport car":
		return newSportCar()
	case "police car":
		return newPoliceCar()
	default:
		return newDefaultCar()
	}
}

func main() {
	car1 := factory("sport car")
	car2 := factory("police car")
	car3 := factory("def")

	fmt.Println(car1.getInfo())
	fmt.Println(car2.getInfo())
	fmt.Println(car3.getInfo())
}
