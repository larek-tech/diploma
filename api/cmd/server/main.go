package main

import (
	"github.com/larek-tech/diploma/api/pkg"
	"github.com/yogenyslav/pkg/errs"
)

// TODO: change host

// @title Diploma API
// @version 1.0
// @description Diploma RAG API service documentation.
// @license.name MIT
// @license.url https://github.com/larek-tech/diploma/blob/api/LICENSE
// @host localhost:9000
// @BasePath /
// @securitydefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	if err := pkg.Run(); err != nil {
		panic(errs.WrapErr(err, "application fatal error"))
	}
}
