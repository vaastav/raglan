module github.com/vaastav/iridescent/apps/userservice/workflow

go 1.22.1

require github.com/blueprint-uservices/blueprint/runtime v0.0.0-20260314172942-77bfbde575a7

require github.com/vaastav/iridescent/iridescent_rt v0.0.0

replace github.com/vaastav/iridescent/iridescent_rt => ../../../iridescent_rt


require (
	github.com/pkg/errors v0.9.1 // indirect
	go.mongodb.org/mongo-driver v1.15.0 // indirect
	go.opentelemetry.io/otel v1.26.0 // indirect
	go.opentelemetry.io/otel/metric v1.26.0 // indirect
	go.opentelemetry.io/otel/trace v1.26.0 // indirect
	golang.org/x/exp v0.0.0-20240416160154-fe59bbe5cc7f // indirect
)
