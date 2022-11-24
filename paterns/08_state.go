package main

import "fmt"

type state interface {
	buyStaff()
}

type Seller struct {
	hasStaff state
	noStaff  state

	current state

	staffCount int
}

func (s Seller) buyStaff() {
	s.current.buyStaff()
}

func NewSeller(staffCount int) *Seller {
	s := &Seller{staffCount: staffCount}

	ns := noStafState{seller: s}
	hs := hasStaffState{seller: s}

	s.hasStaff = &hs
	s.noStaff = &ns

	if staffCount <= 0 {
		s.current = &ns
		return s
	}

	s.current = &hs

	return s
}

type hasStaffState struct {
	seller *Seller
}

func (s hasStaffState) buyStaff() {
	s.seller.staffCount--
	fmt.Println("staff sold, staff count: ", s.seller.staffCount)
	if s.seller.staffCount == 0 {
		s.seller.current = s.seller.noStaff
	}
}

type noStafState struct {
	seller *Seller
}

func (s noStafState) buyStaff() {
	fmt.Println("no staff")
}

func main() {
	seller := NewSeller(1)

	seller.buyStaff()
	seller.buyStaff()
}
