package factory

type ICharacter interface {
	setClass(class string)
	setWeapon(weapon string)
	getClass() string
	getWeapon() string
}

type Character struct {
	class  string
	weapon string
}

func (c Character) setClass(class string) {
	c.class = class
}

func (c Character) setWeapon(weapon string) {
	c.weapon = weapon
}

func (c Character) getClass() string {
	return c.class
}

func (c Character) getWeapon() string {
	return c.weapon
}

type Wizard struct {
	Character
}

type Warrior struct {
	Character
}

func main() {
	wizard := Wizard{Character{}}
	warrior := Warrior{Character{}}

	wizard.setClass("Wizard")
	wizard.setWeapon("Magic staff")

	warrior.setClass("Warrior")
	warrior.setWeapon("Sword")

}
