package main

import "fmt"

//	Строитель — это порождающий паттерн проектирования, используется для построения сложных объектов покомпонентно,
//  дает возможность создания немного отличающихся в значениях, но одинаковых в конструкции объектов.
// 	К плюсам можно отнести пошаговое создание, один и тот же код создает различные объекты
//	К минусам можно отнести усложнение кода программы, введение доп типов

// Структура, которую можно по разному сконфигугировать
type car struct {
	EngineVolume  float32
	NumberOfDoors int
	RadiusWheel   int
	Clearance     int
}

// добавил возврат carBuilder из метода для того чтоб можно было вызывать методы через .
type carBuilderI interface {
	SetEngineVolume(val float32) carBuilderI
	SetNumberOfDoors(val int) carBuilderI
	SetRadiusWheel(val int) carBuilderI
	SetClearance(val int) carBuilderI

	Build() car
}

// Структура, позволяющая сконфигурировать любую машину
type carBuilder struct {
	engineVolume  float32
	numberOfDoors int
	radiusWheel   int
	clearance     int
}

func newCarBuilder() carBuilderI {
	return &carBuilder{}
}

func (b *carBuilder) SetEngineVolume(val float32) carBuilderI {
	b.engineVolume = val
	return b
}

func (b *carBuilder) SetNumberOfDoors(val int) carBuilderI {
	b.numberOfDoors = val
	return b
}

func (b *carBuilder) SetRadiusWheel(val int) carBuilderI {
	b.radiusWheel = val
	return b
}

func (b *carBuilder) SetClearance(val int) carBuilderI {
	b.clearance = val
	return b
}

func (b *carBuilder) Build() car {
	return car{
		EngineVolume:  b.engineVolume,
		NumberOfDoors: b.numberOfDoors,
		RadiusWheel:   b.radiusWheel,
		Clearance:     b.clearance,
	}
}

// Конструктор спорткара
func newSportCarBuilder() carBuilderI {
	return newCarBuilder().SetClearance(4).SetEngineVolume(5).SetNumberOfDoors(2).SetRadiusWheel(20)
}

func main() {
	//билдер без значений по умолчанию
	{
		builder := newCarBuilder()
		car := builder.SetClearance(10).SetEngineVolume(2.4).SetNumberOfDoors(5).SetRadiusWheel(16).Build()

		fmt.Println(car)
	}

	//билдер конкретной машины
	{
		builder := newSportCarBuilder()
		car := builder.Build()

		fmt.Println(car)

		// можно менять поля
		builder.SetNumberOfDoors(5)

		car = builder.Build()

		fmt.Println(car)
	}
}
