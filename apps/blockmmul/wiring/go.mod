module github.com/vaastav/raglan/apps/blockmmul/wiring

go 1.22.1

require github.com/vaastav/raglan/apps/blockmmul/workflow v0.0.0

replace github.com/vaastav/raglan/apps/blockmmul/workflow => ../workflow

require (
	github.com/blueprint-uservices/blueprint/blueprint v0.0.0-20241113113418-f54e1bbd9997
	github.com/blueprint-uservices/blueprint/plugins v0.0.0-20241113113418-f54e1bbd9997
	github.com/vaastav/raglan/plugins v0.0.0-20260428122744-a06b6965a2ea
)

require (
	github.com/blueprint-uservices/blueprint/runtime v0.0.0-20240619221802-d064c5861c1e // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/otiai10/copy v1.14.0 // indirect
	github.com/vaastav/raglan/iridescent_rt v0.0.0-20260428105604-0ea5820e202a // indirect
	go.mongodb.org/mongo-driver v1.15.0 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdoutmetric v1.26.0 // indirect
	go.opentelemetry.io/otel/exporters/stdout/stdouttrace v1.26.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/sdk v1.26.0 // indirect
	go.opentelemetry.io/otel/sdk/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.19.0 // indirect
	golang.org/x/tools v0.20.0 // indirect
)
