package utils

import "testing"

func TestCopyFields(t *testing.T) {
	type args struct {
		src interface{}
		dst interface{}
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "copy field test",
			args: args{
				src: &struct {
					Name string
					//Age  uint32
				}{
					Name: "test name1",
					//Age:  1,
				},
				dst: &struct {
					Name string
					//Age  uint32
				}{
					Name: "test name2",
					//Age:  5,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			CopyFields(tt.args.src, tt.args.dst)
		})
	}
}
