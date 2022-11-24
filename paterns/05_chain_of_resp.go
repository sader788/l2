package main

import "fmt"

type Department interface {
	execute(*Order)
	setNext(Department)
}

type Order struct {
	isCreated  bool
	isRecieved bool
	isPayment  bool
}

type DrugStorage struct {
	next Department
}

func NewDrugStorage() *DrugStorage {
	return &DrugStorage{}
}

func (d *DrugStorage) execute(order *Order) {
	fmt.Println("Order created")
	order.isCreated = true
	d.next.execute(order)
}

func (d *DrugStorage) setNext(department Department) {
	d.next = department
}

type DrugStore struct {
	next Department
}

func (d *DrugStore) execute(order *Order) {
	if !order.isCreated {
		return
	}
	fmt.Println("Order shipped from storage")
	order.isRecieved = true
	d.next.execute(order)
}

func (d *DrugStore) setNext(department Department) {
	d.next = department
}

type Cashier struct {
	next Department
}

func (c *Cashier) execute(order *Order) {
	if !order.isRecieved {
		return
	}
	fmt.Println("Order is paid")
	order.isPayment = true
}

func (c *Cashier) setNext(next Department) {
	c.next = next
}

func main() {

	cashier := &Cashier{}

	//Set next for doctor department
	store := &DrugStore{}
	store.setNext(cashier)

	//Set next for medical department
	storage := &DrugStorage{}
	storage.setNext(store)

	order := &Order{}
	//Patient visiting
	storage.execute(order)
}
