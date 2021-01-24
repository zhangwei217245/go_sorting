package datasorter

import (
	"container/heap"
	"fmt"
	"testing"
)

func TestSortByID(t *testing.T) {
	raw := []string{
		"1097,lielqtjhrq,yxyrq77v60fko5h93,Australia",
		"1096,ovvygwuzgpauj,4mcratsgk4nu6fu,Australia",
		"1099,lmyopgtaiuv,zokqpm p1pefz0h4x6n,Asia",
		"1098,agqsysyxsvs, f6m96 qckiu21s5,North America",
	}
	data := []Item{}
	for i := 0; i < len(raw); i++ {
		data = append(data, Item{raw[i], i})
	}
	SortByID(data)
	fmt.Println(data)
	fmt.Println("===========")
}

func TestSortByName(t *testing.T) {
	raw := []string{
		"1097,lielqtjhrq,yxyrq77v60fko5h93,Australia",
		"1096,ovvygwuzgpauj,4mcratsgk4nu6fu,Australia",
		"1099,lmyopgtaiuv,zokqpm p1pefz0h4x6n,Asia",
		"1098,agqsysyxsvs, f6m96 qckiu21s5,North America",
	}
	data := []Item{}
	for i := 0; i < len(raw); i++ {
		data = append(data, Item{raw[i], i})
	}
	SortByName(data)
	fmt.Println(data)
	fmt.Println("===========")
}

func TestSortByContinent(t *testing.T) {
	raw := []string{
		"1097,lielqtjhrq,yxyrq77v60fko5h93,Australia",
		"1096,ovvygwuzgpauj,4mcratsgk4nu6fu,Australia",
		"1099,lmyopgtaiuv,zokqpm p1pefz0h4x6n,Asia",
		"1098,agqsysyxsvs, f6m96 qckiu21s5,North America",
	}
	data := []Item{}
	for i := 0; i < len(raw); i++ {
		data = append(data, Item{raw[i], i})
	}
	SortByContinent(data)
	fmt.Println(data)
	fmt.Println("===========")
}

func TestByIDHeap(t *testing.T) {
	raw := []string{
		"1097,lielqtjhrq,yxyrq77v60fko5h93,Australia",
		"1096,ovvygwuzgpauj,4mcratsgk4nu6fu,Australia",
		"1099,lmyopgtaiuv,zokqpm p1pefz0h4x6n,Asia",
		"1098,lmyopgtaiuv, f6m96 qckiu21s5,North America",
	}
	data := []Item{}
	for i := 0; i < len(raw); i++ {
		data = append(data, Item{raw[i], i})
	}
	d := byID(data)
	h := &d

	heap.Init(h)
	heap.Push(h, Item{"1081,lmyopgtaiuv,pnmcm83qy d3hwrral,Africa", 5})
	for h.Len() > 0 {
		min := heap.Pop(h)
		fmt.Println(min)
	}
	fmt.Println("===========")
}

func TestByNameHeap(t *testing.T) {
	raw := []string{
		"1097,lielqtjhrq,yxyrq77v60fko5h93,Australia",
		"1096,ovvygwuzgpauj,4mcratsgk4nu6fu,Australia",
		"1099,lmyopgtaiuv,zokqpm p1pefz0h4x6n,Asia",
		"1098,lmyopgtaiuv, f6m96 qckiu21s5,North America",
	}
	data := []Item{}
	for i := 0; i < len(raw); i++ {
		data = append(data, Item{raw[i], i})
	}
	d := byName(data)
	h := &d

	heap.Init(h)
	heap.Push(h, Item{"1081,lmyopgtaiuv,pnmcm83qy d3hwrral,Africa", 5})
	for h.Len() > 0 {
		min := heap.Pop(h)
		fmt.Println(min)
	}
	fmt.Println("===========")
}

func TestByContinentHeap(t *testing.T) {
	raw := []string{
		"1097,lielqtjhrq,yxyrq77v60fko5h93,Australia",
		"1096,ovvygwuzgpauj,4mcratsgk4nu6fu,Australia",
		"1099,lmyopgtaiuv,zokqpm p1pefz0h4x6n,Asia",
		"1098,lmyopgtaiuv, f6m96 qckiu21s5,North America",
	}
	data := []Item{}
	for i := 0; i < len(raw); i++ {
		data = append(data, Item{raw[i], i})
	}
	d := byContinent(data)
	h := &d

	heap.Init(h)
	heap.Push(h, Item{"1081,lmyopgtaiuv,pnmcm83qy d3hwrral,Africa", 5})
	for h.Len() > 0 {
		min := heap.Pop(h)
		fmt.Println(min)
	}
	fmt.Println("===========")
}
