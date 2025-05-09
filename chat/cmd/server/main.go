package main

import (
	"github.com/larek-tech/diploma/chat/pkg"
	"github.com/yogenyslav/pkg/errs"
)

func main() {
	if err := pkg.Run(); err != nil {
		panic(errs.WrapErr(err, "application fatal error"))
	}
}
