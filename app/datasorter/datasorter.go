package datasorter

import (
	"sort"
	"strconv"
	"strings"
)

type Item struct {
	Data   string
	Source int
}

type ByID []Item

type ByName []Item

type ByContinent []Item

func byIDLess(si []string, sj []string) bool {
	idi, _ := strconv.ParseUint(si[0], 10, 32)
	idj, _ := strconv.ParseUint(sj[0], 10, 32)
	return uint32(idi) < uint32(idj)
}

func (s ByID) Len() int {
	return len(s)
}
func (s ByID) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByID) Less(i, j int) bool {
	si := strings.Split(s[i].Data, ",")
	sj := strings.Split(s[j].Data, ",")
	return byIDLess(si, sj)
}

func (s ByName) Len() int {
	return len(s)
}
func (s ByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByName) Less(i, j int) bool {
	si := strings.Split(s[i].Data, ",")
	sj := strings.Split(s[j].Data, ",")
	if si[1] == sj[1] { // if two records are with the same name, sort by ID
		return byIDLess(si, sj)
	}
	return si[1] < sj[1]
}

func (s ByContinent) Len() int {
	return len(s)
}
func (s ByContinent) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s ByContinent) Less(i, j int) bool {
	si := strings.Split(s[i].Data, ",")
	sj := strings.Split(s[j].Data, ",")
	// if two continent are the same, sort by ID
	if si[3] == sj[3] {
		return byIDLess(si, sj)
	}
	return si[3] < sj[3]
}

// SortByID sorts the slice of rows by the ID column
func SortByID(s []Item) {
	sort.Sort(ByID(s))
}

// SortByName sorts the slice of rows by the name column
func SortByName(s []Item) {
	sort.Sort(ByName(s))
}

// SortByContinent sorts the slice of rows by the continent column
func SortByContinent(s []Item) {
	sort.Sort(ByContinent(s))
}

func (s *ByID) Push(x interface{}) {
	*s = append(*s, x.(Item))
}

func (s *ByID) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

func (s *ByName) Push(x interface{}) {
	*s = append(*s, x.(Item))
}

func (s *ByName) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

func (s *ByContinent) Push(x interface{}) {
	*s = append(*s, x.(Item))
}

func (s *ByContinent) Pop() interface{} {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}
