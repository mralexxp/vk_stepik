package main

import "testing"

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
		{name: "#IntTest",
			args: args{
				in:   in,
				out:  out,
				data: 0,
			},
			want: "4108050209~502633748"},
		{name: "#StringTest",
			args: args{
				in:   in,
				out:  out,
				data: "abc",
			},
			// TODO: ...
			want: "какой-то хеш из стрингов"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go SingleHash(tt.args.in, tt.args.out)
			in <- tt.args.data

			// По очереди найти и сравнить полученные данные
			md5Data := <-out
			t.Logf("md5Data:%v", md5Data)
			if md5Data.(string) != tt.want {
				t.Fail()
			}
		})
	}
}
