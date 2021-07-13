package main

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTranches(t *testing.T) {
	for _, tc := range []struct {
		Name     string
		Data     []*person
		Expected []tranche
	}{
		{
			Name: "1",
			Data: []*person{
				{Name: "Denis", Bill: 5500},
				{Name: "Aigul", Bill: 0},
				{Name: "Natasha", Bill: 1600},
				{Name: "Inna", Bill: 0},
				{Name: "Andrey", Bill: 5000},
			},
			Expected: []tranche{
				{From: "Aigul", To: "Denis", Amount: 2420},
				{From: "Natasha", To: "Denis", Amount: 660},
				{From: "Natasha", To: "Andrey", Amount: 160},
				{From: "Inna", To: "Andrey", Amount: 2420},
			},
		}, {
			Name: "2",
			Data: []*person{
				{Name: "A", Bill: 700},
				{Name: "B", Bill: 300},
			},
			Expected: []tranche{
				{From: "B", To: "A", Amount: 200},
			},
		},
	} {
		t.Run(tc.Name, func(t *testing.T) {
			actual := computeTranches(tc.Data)
			require.Equal(t, tc.Expected, actual)
		})
	}
}
