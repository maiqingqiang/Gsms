package core

import (
	"fmt"
	"strconv"
	"strings"
)

type PhoneNumberInterface interface {
	Number() int64
	IDDCode() int
	UniversalNumber() string
	ZeroPrefixedNumber() string
	PrefixedIDDCode(prefix string) string
	String() string
	InChineseMainland() bool
}

type PhoneNumber struct {
	number  int64
	iddCode int
}

func NewPhoneNumber(numberWithoutIDDCode int64, iddCodeStr string) (*PhoneNumber, error) {
	iddCodeStr = strings.TrimLeft(iddCodeStr, "+0")

	iddCode, err := strconv.Atoi(iddCodeStr)

	if err != nil {
		return nil, ErrInvalidIDDCode
	}

	return &PhoneNumber{
		number:  numberWithoutIDDCode,
		iddCode: iddCode,
	}, nil
}

func NewNewPhoneNumberWithoutIDDCode(numberWithoutIDDCode int64) *PhoneNumber {
	return &PhoneNumber{number: numberWithoutIDDCode}
}

// Number e.g. 86.
func (p *PhoneNumber) Number() int64 {
	return p.number
}

// IDDCode e.g. 13800138000.
func (p *PhoneNumber) IDDCode() int {
	return p.iddCode
}

//UniversalNumber  e.g. +8613800138000.
func (p *PhoneNumber) UniversalNumber() string {
	return fmt.Sprintf("%s%d", p.PrefixedIDDCode("+"), p.number)
}

// ZeroPrefixedNumber e.g. 008613800138000.
func (p *PhoneNumber) ZeroPrefixedNumber() string {
	return fmt.Sprintf("%s%d", p.PrefixedIDDCode("00"), p.number)
}

// PrefixedIDDCode  e.g. + -> +86
func (p *PhoneNumber) PrefixedIDDCode(prefix string) string {
	if p.iddCode == 0 {
		return ""
	}
	return fmt.Sprintf("%s%d", prefix, p.iddCode)
}

func (p *PhoneNumber) String() string {
	return p.UniversalNumber()
}

// InChineseMainland Check if the phone number belongs to chinese mainland.
func (p *PhoneNumber) InChineseMainland() bool {
	return p.iddCode == 86
}
