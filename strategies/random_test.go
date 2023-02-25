package strategies

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomStrategy_Apply(t *testing.T) {
	type args struct {
		gateways []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Random Strategy",
			args: args{
				gateways: []string{"yunpian", "aliyun", "aliyunrest", "aliyunintl", "submail", "huyi", "juhe"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &RandomStrategy{}

			newGateways := s.Apply(tt.args.gateways)
			t.Logf("newGateways: %v", newGateways)

			assert.Equal(t, len(tt.args.gateways), len(newGateways), "Apply(%v)", tt.args.gateways)
		})
	}
}
