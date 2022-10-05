package main

import (
	"errors"
	"fmt"
)

// поведенческий паттерн, используется для:
// смены состояний объекта после, например, запроса, то есть применим к логике, где у объекта есть конкретные состояния;

// пример: есть какой-то магазин, покупать можно только когда есть товары, добавлять товары можно только когда они закончились
type magazine struct {
	hasItem state
	noItem  state

	currentState state

	itemCount int
	itemPrice int
}

func (m *magazine) buy(money int) error {
	return m.currentState.buy(money)
}

func (m *magazine) loadProducts(count int) error {
	return m.currentState.loadProducts(count)
}

func newMagazine(count, price int) *magazine {

	m := &magazine{itemCount: count, itemPrice: price}

	m.hasItem = hasItemState{m}
	m.noItem = noItemState{m}

	if m.itemCount > 0 {
		m.setState(m.hasItem)
	} else {
		m.setState(m.noItem)
	}

	return m
}

func (m *magazine) setState(s state) {
	m.currentState = s
}

type state interface {
	buy(int) error
	loadProducts(int) error
}

type hasItemState struct {
	m *magazine
}

func (s hasItemState) buy(money int) error {
	if money < s.m.itemPrice {
		return errors.New("not enough money")
	}
	s.m.itemCount--
	if s.m.itemCount <= 0 {
		s.m.setState(s.m.noItem)
	}
	fmt.Println("succesfully buy item")
	return nil
}

func (s hasItemState) loadProducts(count int) error {
	return errors.New("there are still products in stock")
}

type noItemState struct {
	m *magazine
}

func (s noItemState) buy(money int) error {
	return errors.New("Item out of stock")
}

func (s noItemState) loadProducts(count int) error {
	s.m.itemCount += count
	if s.m.itemCount > 0 {
		s.m.setState(s.m.hasItem)
	}
	fmt.Println("succesfully add items")
	return nil
}

func main() {
	m := newMagazine(1, 500)

	err := m.loadProducts(5)
	if err != nil {
		fmt.Println(err)
	}
	err = m.buy(500)
	if err != nil {
		fmt.Println(err)
	}
	err = m.buy(500)
	if err != nil {
		fmt.Println(err)
	}
	err = m.loadProducts(1)
	if err != nil {
		fmt.Println(err)
	}
	err = m.buy(499)
	if err != nil {
		fmt.Println(err)
	}
	err = m.buy(500)
	if err != nil {
		fmt.Println(err)
	}
}
