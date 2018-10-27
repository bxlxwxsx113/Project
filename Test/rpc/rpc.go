package rpc

import "errors"

type DemoService struct{}

type Args struct {
	A int
	B int
}

func (DemoService) Div(arg Args, result *float64) error {
	if arg.B == 0 {
		return errors.New("divsion is zero")
	}
	*result = float64(arg.A) / float64(arg.B)
	return nil
}
