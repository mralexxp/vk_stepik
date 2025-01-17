package main

import (
	"context"
)

func NewBiz() *Biz {
	const OP = "NewBiz"

	return &Biz{}
}

type Biz struct {
	UnimplementedBizServer
}

func (b *Biz) Check(context.Context, *Nothing) (*Nothing, error) {
	const OP = "Biz.Check"

	return &Nothing{}, nil
}
func (b *Biz) Add(context.Context, *Nothing) (*Nothing, error) {
	const OP = "Biz.Add"

	return &Nothing{}, nil
}
func (b *Biz) Test(context.Context, *Nothing) (*Nothing, error) {
	const OP = "Biz.Test"

	return &Nothing{}, nil
}
