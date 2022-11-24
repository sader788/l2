package strategy

import "fmt"

type CharacterActivity interface {
	execute()
}

type Character struct {
	name     string
	activity CharacterActivity
}

func (ch *Character) Do() {
	ch.activity.execute()
}

type Sleeping struct {
}

func (s *Sleeping) execute() {
	fmt.Println("npc sleep")
}

type Working struct {
}

func (w *Working) execute() {
	fmt.Println("npc work")
}

type Talking struct {
}

func (t Talking) execute() {
	fmt.Println("npc talk")
}

func main() {
	character := &Character{}

	character.activity = &Sleeping{}
	character.Do()

	character.activity = &Working{}
	character.Do()

	character.activity = &Talking{}
	character.Do()

}
