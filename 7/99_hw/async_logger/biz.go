package main

import (
	"context"
	"log"
)

func NewBiz() *Biz {
	const OP = "NewBiz"
	log.Print(OP)

	return &Biz{}
}

type Biz struct {
	UnimplementedBizServer
}

func (b *Biz) Check(context.Context, *Nothing) (*Nothing, error) {
	const OP = "Biz.Check"
	log.Print(OP)

	return &Nothing{}, nil
}
func (b *Biz) Add(context.Context, *Nothing) (*Nothing, error) {
	const OP = "Biz.Add"
	log.Print(OP)

	return &Nothing{}, nil
}
func (b *Biz) Test(context.Context, *Nothing) (*Nothing, error) {
	const OP = "Biz.Test"
	log.Print(OP)

	return &Nothing{}, nil
}
