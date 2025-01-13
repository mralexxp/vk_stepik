package main

import "context"

func NewBiz() *Biz {
	return &Biz{}
}

type Biz struct {
	UnimplementedBizServer
}

func (b *Biz) Check(context.Context, *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}
func (b *Biz) Add(context.Context, *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}
func (b *Biz) Test(context.Context, *Nothing) (*Nothing, error) {
	return &Nothing{}, nil
}
