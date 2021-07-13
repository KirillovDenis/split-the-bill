package main

import (
	"fmt"
)

type person struct {
	Name string
	Bill int
}

type tranche struct {
	From   string
	To     string
	Amount int
}

func main() {

	persons := []*person{
		{Name: "Denis", Bill: 5500},
		{Name: "Aigul", Bill: 0},
		{Name: "Natasha", Bill: 1600},
		{Name: "Inna", Bill: 0},
		{Name: "Andrey", Bill: 5000},
	}

	tranches := computeTranches(persons)
	fmt.Println(tranches)
}

func computeTranches(persons []*person) []tranche {
	sum := 0
	for _, person := range persons {
		sum += person.Bill
	}

	average := sum / len(persons)
	for i := range persons {
		persons[i].Bill -= average
	}

	var res []tranche
	for i, from := range persons {
		for j, to := range persons {
			if from.Bill >= 0 {
				continue
			}
			if to.Bill <= 0 {
				continue
			}

			t := tranche{From: from.Name, To: to.Name}
			tmp := from.Bill + to.Bill

			if tmp < 0 {
				t.Amount = to.Bill
				persons[i].Bill = tmp
				persons[j].Bill = 0
			} else {
				t.Amount = -from.Bill
				persons[i].Bill = 0
				persons[j].Bill = tmp
			}
			res = append(res, t)
		}
	}

	return res
}
