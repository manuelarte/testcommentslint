package main

type MyStruct struct {
	id int
	name string
}

func NewMyStruct(id int, name string) MyStruct {
	return MyStruct{
		id: id,
		name: name,
	}
}
