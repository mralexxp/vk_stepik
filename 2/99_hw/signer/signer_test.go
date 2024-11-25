package main

import (
	"testing"
	"time"
)

func TestSingleHash(t *testing.T) {
	type args struct {
		in   chan interface{}
		out  chan interface{}
		data interface{}
		want string
	}
	in := make(chan interface{})
	out := make(chan interface{})

	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "#IntTest-0",
			args: args{
				in:   in,
				out:  out,
				data: 0,
			},
			want: "4108050209~502633748"},
		{name: "#IntTest-1",
			args: args{
				in:   in,
				out:  out,
				data: 1,
			},
			want: "2212294583~709660146"},
		//{name: "#StringTest",
		//	args: args{
		//		in:   in,
		//		out:  out,
		//		data: "abc",
		//	},
		//	// TODO: ...
		//	want: "какой-то хеш из стрингов"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go SingleHash(tt.args.in, tt.args.out)
			in <- tt.args.data

			// По очереди найти и сравнить полученные данные
			md5Data := <-out
			if md5Data.(string) != tt.want {
				t.Fail()
			}
		})
	}
}

func TestMultiHash(t *testing.T) {
	type args struct {
		in   chan interface{}
		out  chan interface{}
		data interface{}
		want string
	}

	in := make(chan interface{})
	out := make(chan interface{})

	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "#IntTest-0",
			args: args{
				in:   in,
				out:  out,
				data: 0,
			},
			want: "29568666068035183841425683795340791879727309630931025356555"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go SingleHash(tt.args.in, tt.args.out)
			go MultiHash(tt.args.out, tt.args.in)

			in <- tt.args.data

			time.Sleep(100 * time.Millisecond)

			md5Data := <-out
			t.Log(md5Data)
			if md5Data.(string) != tt.want {
				t.Fail()
			}
		})
	}
}
