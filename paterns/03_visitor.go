package visitor

import "fmt"

type IVisitor interface {
	countMonthProfitSalon()
	countMonthProfitShop()
}

type ProfitCalculator struct {
}

func (pc ProfitCalculator) countMonthProfitSalon() {
	fmt.Println("profit is million")
}

func (pc ProfitCalculator) countMonthProfitShop() {
	fmt.Println("profit is billion")
}

type Establishment interface {
	accept(IVisitor)
}

type Saloon struct {
}

type Shop struct {
}

func (s Saloon) accept(v IVisitor) {
	v.countMonthProfitSalon()
}

func (s Shop) accept(v IVisitor) {
	v.countMonthProfitShop()
}

func NewShop() Shop {
	return Shop{}
}

func NewSaloon() Saloon {
	return Saloon{}
}

func NewVisitor() IVisitor {
	return ProfitCalculator{}
}

func main() {
	calculator := NewVisitor()

	shop := NewShop()
	saloon := NewSaloon()

	shop.accept(calculator)
	saloon.accept(calculator)
}