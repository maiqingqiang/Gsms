package gsms

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPhoneNumber(t *testing.T) {
	as := assert.New(t)

	phoneNumber := NewPhoneNumber(18888888888, "86")
	as.Equal(86, phoneNumber.IDDCode())
	as.EqualValues(18888888888, phoneNumber.Number())
	as.Equal("+8618888888888", phoneNumber.UniversalNumber())
	as.Equal("+8618888888888", phoneNumber.String())
	as.Equal("008618888888888", phoneNumber.ZeroPrefixedNumber())
	as.Equal("+86", phoneNumber.PrefixedIDDCode("+"))
	as.Equal("0086", phoneNumber.PrefixedIDDCode("00"))
	as.True(phoneNumber.InChineseMainland())
}

func TestPhoneNumber_IDDCode(t *testing.T) {
	type fields struct {
		number  int
		iddCode string
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "+86",
			fields: fields{
				number:  18888888888,
				iddCode: "+86",
			},
			want: 86,
		},
		{
			name: "0086",
			fields: fields{
				number:  18888888888,
				iddCode: "0086",
			},
			want: 86,
		},
		{
			name: "empty",
			fields: fields{
				number:  18888888888,
				iddCode: "",
			},
			want: 0,
		},
		{
			name: "other",
			fields: fields{
				number:  18888888888,
				iddCode: "other",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewPhoneNumber(tt.fields.number, tt.fields.iddCode)
			assert.Equalf(t, tt.want, p.IDDCode(), "IDDCode()")
		})
	}
}
