package command

type Button struct {
	command Command
}

func (b *Button) press() {
	b.command.execute()
}

type Command interface {
	execute()
}

type OnCommand struct {
	device Device
}

func (c *OnCommand) execute() {
	c.device.on()
}

type OffCommand struct {
	device Device
}

func (c *OffCommand) execute() {
	c.device.off()
}

type Device interface {
	on()
	off()
}

type Plate struct {
	isHeated bool
}

func (t *Plate) on() {
	t.isHeated = true
}

func (t *Plate) off() {
	t.isHeated = false
}

func main() {
	plate := &Plate{}

	onCommand := &OnCommand{
		device: plate,
	}

	offCommand := &OffCommand{
		device: plate,
	}

	onButton := &Button{
		command: onCommand,
	}
	onButton.press()

	offButton := &Button{
		command: offCommand,
	}
	offButton.press()
}
