package strategies

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestOrderStrategy_Apply(t *testing.T) {
	type args struct {
		gateways []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Order strategy",
			args: args{
				gateways: []string{"yunpian", "aliyun", "aliyunrest", "aliyunintl", "submail", "huyi", "juhe"},
			},
			want: []string{"aliyun", "aliyunintl", "aliyunrest", "huyi", "juhe", "submail", "yunpian"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := &OrderStrategy{}

			newGateways := o.Apply(tt.args.gateways)
			t.Logf("newGateways: %v", newGateways)

			assert.Equalf(t, tt.want, newGateways, "Apply(%v)", tt.args.gateways)
		})
	}
}
