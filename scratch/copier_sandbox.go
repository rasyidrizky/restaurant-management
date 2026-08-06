package main

import (
	"fmt"
	"github.com/jinzhu/copier"
)

type UpdateReq struct {
	Name *string
	Age  *int
}

type Model struct {
	Name string
	Age  int
}

func main() {
	m := Model{Name: "Old", Age: 30}
	
	newName := "New"
	req := UpdateReq{Name: &newName} // Age is nil
	
	copier.CopyWithOption(&m, &req, copier.Option{IgnoreEmpty: true})
	
	fmt.Printf("Model: %+v\n", m)
}
