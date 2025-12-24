module github.com/QFServer

go 1.25.5

replace github.com/QFServer/client => ./client

require (
	github.com/QFServer/client v0.0.0-00010101000000-000000000000
	github.com/QFServer/log v0.0.0-00010101000000-000000000000
)

require github.com/QFServer/server v0.0.0-00010101000000-000000000000 // indirect

replace github.com/QFServer/log => ./log

replace github.com/QFServer/server => ./server
