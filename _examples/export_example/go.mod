module export_example

go 1.21

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/watchakorn-18k/scalar-go v0.0.0
)

replace github.com/watchakorn-18k/scalar-go => ../..
