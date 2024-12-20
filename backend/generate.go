package backend

//go:generate sqlc generate
//go:generate mockgen -source internal/db/wrapped_querier.go -destination internal/db/_mock/querier.go -package mock_db
//go:generate mockgen -source internal/services/jwt.go -destination internal/services/_mock/jwt.go -package mock_services
//go:generate mockgen -source internal/services/password.go -destination internal/services/_mock/password.go -package mock_services
