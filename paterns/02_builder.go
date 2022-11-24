package builder

import "fmt"

type Wheel struct {
}

type MachineBody struct {
}

type Glass struct {
}

type Engine struct {
}

type Car struct {
	wheels      []Wheel
	glasses     []Glass
	machineBody MachineBody
	engine      Engine
}

type IBuilder interface {
	BuildWheels()
	BuildGlasses()
	BuildMachineBody()
	BuildEngine()
	GetCar() Car
}

func NewBuilder() IBuilder {
	return NissanBuilder{}
}

type NissanBuilder struct {
	wheels      []Wheel
	glasses     []Glass
	machineBody MachineBody
	engine      Engine
}

func (nb NissanBuilder) BuildWheels() {
	nb.wheels = []Wheel{Wheel{}}
}

func (nb NissanBuilder) BuildGlasses() {
	nb.glasses = []Glass{Glass{}}
}

func (nb NissanBuilder) BuildMachineBody() {
	nb.machineBody = MachineBody{}
}

func (nb NissanBuilder) BuildEngine() {
	nb.engine = Engine{}
}

func (nb NissanBuilder) GetCar() Car {
	return Car{
		wheels:      nb.wheels,
		glasses:     nb.glasses,
		machineBody: nb.machineBody,
		engine:      nb.engine,
	}
}

//директор

type Mechanic struct {
	builder IBuilder
}

func NewMechanic(builder IBuilder) Mechanic {
	return Mechanic{builder: builder}
}

func (m Mechanic) make() Car {
	m.builder.BuildWheels()
	m.builder.BuildGlasses()
	m.builder.BuildMachineBody()
	m.builder.BuildEngine()
	return m.builder.GetCar()
}

func main() {
	builder := NewBuilder()
	mechanic := NewMechanic(builder)
	car := mechanic.make()
	fmt.Println(car)
}
