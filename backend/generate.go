package backend

//go:generate sqlc generate
//go:generate mockgen -source internal/db/wrapped_querier.go -destination internal/db/_mock/querier.go -package mock_db
//go:generate mockgen -source internal/services/services.go -destination internal/services/_mock/services.go -package mock_services
//go:generate swag init -g internal/router/router.go --parseDependency
