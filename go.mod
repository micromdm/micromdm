module micromdm.io/v2

go 1.14

require (
	crawshaw.io/sqlite v0.3.2
	github.com/felixge/httpsnoop v1.0.1
	github.com/go-kit/kit v0.10.0
	github.com/gorilla/csrf v1.7.0
	github.com/gorilla/mux v1.7.4
	github.com/jackc/pgconn v1.6.2
	github.com/jackc/pgx/v4 v4.7.2
	github.com/oklog/run v1.1.0
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/oklog/ulid/v2 v2.0.2
	github.com/peterbourgon/ff/v3 v3.0.0
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	rsc.io/goversion v1.2.0
)

replace crawshaw.io/sqlite => github.com/groob/sqlite v0.3.3-0.20200721040052-b46ed0907467
