package utils

import (
	"math/rand"
	"sort"
)

type Choice[T any] struct {
	Item   T
	Weight uint
}

type Chooser[T any] struct {
	choices []Choice[T]
	total   uint
}

func NewChooser[T any](choices []Choice[T]) *Chooser[T] {
	c := &Chooser[T]{
		choices: make([]Choice[T], len(choices)),
	}
	for i, choice := range choices {
		c.choices[i] = choice
		c.total += choice.Weight
	}
	sort.Slice(c.choices, func(i, j int) bool {
		return c.choices[i].Weight < c.choices[j].Weight
	})
	return c
}

func (c *Chooser[T]) Pick() *T {
	r := rand.Intn(int(c.total)) + 1
	for _, choice := range c.choices {
		r -= int(choice.Weight)
		if r <= 0 {
			return &choice.Item
		}
	}
	return nil
}
