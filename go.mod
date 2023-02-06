module github.com/nmramorov/gowatcher

go 1.19

require (
	github.com/stretchr/testify v1.8.1
	internal/metrics v0.0.0-00010101000000-000000000000
)

require (
	github.com/caarlos0/env/v6 v6.10.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-chi/chi/v5 v5.0.8 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace internal/metrics => ./internal
