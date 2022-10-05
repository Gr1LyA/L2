package main

import "fmt"

// Фасад - структурный паттерн, реализует простой доступ к сложной системе
// Плюсы: предоставляет возможность легко манипулировать всей системой, изолирует клиентов от компонентов сложной подсистемы
// Минусы: есть риск что фасад станет божественным объектом(который хранит в себе слишком много или делает слишком много)

func main() {
	newMusicalGroup("вокал", "гитара", "барабаны", "басгитара").playStandartSong()
}

type musicalGroup struct {
	vocal  vocalist
	guitar guitarist
	drum   drummer
	bass   bassist
}

func newMusicalGroup(vocalName, guitarName, drumName, bassName string) musicalGroup {
	return musicalGroup{
		vocal:  vocalist{vocalName},
		guitar: guitarist{guitarName},
		drum:   drummer{drumName},
		bass:   bassist{bassName},
	}
}

func (m musicalGroup) playStandartSong() {
	m.drum.startPlaying()

	for i := 0; i < 3; i++ {
		m.bass.changeRhytm("куплет")
		m.guitar.playCouplet()
		m.vocal.singCouplet()
		m.bass.changeRhytm("припев")
		m.guitar.playChorus()
		m.vocal.singChorus()
	}

	m.drum.stopPlaying()
}

type vocalist struct {
	name string
}

func (v vocalist) singCouplet() {
	fmt.Println(v.name+":", "поет красивый куплет")
}

func (v vocalist) singChorus() {
	fmt.Println(v.name+":", "поет мощный припев")
}

type guitarist struct {
	name string
}

func (v guitarist) playCouplet() {
	fmt.Println(v.name+":", "играет куплет")
}

func (v guitarist) playChorus() {
	fmt.Println(v.name+":", "играет припев")
}

type drummer struct {
	name string
}

func (d drummer) startPlaying() {
	fmt.Println(d.name+":", "Начинает играть")
}

func (d drummer) stopPlaying() {
	fmt.Println(d.name+":", "Заканчивает")
}

type bassist struct {
	name string
}

func (v bassist) changeRhytm(part string) {
	fmt.Println(v.name+":", "Перешел на ритм "+part+"a")
}
