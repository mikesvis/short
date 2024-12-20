#!/bin/bash

# go build -o ./cmd/staticlint/staticlint ./cmd/staticlint/*.go
# go build -o ./cmd/shortener/shortener ./cmd/shortener/*.go
# go test ./... -coverprofile cover.out && go tool cover -func cover.out

~/bin/shortenertestbeta  -test.v -test.run=^TestIteration1$ -binary-path=cmd/shortener/shortener;
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration2$ -binary-path=cmd/shortener/shortener -source-path=./;
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration3$ -binary-path=cmd/shortener/shortener -source-path=./;
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration4$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081;
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration5$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081;
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration7$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081;
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration6$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081;
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration8$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081;
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration9$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081 -file-storage-path=/tmp/fsgo.json;
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration10$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081 -database-dsn="host=0.0.0.0 port=5433 user=postgres password=postgres dbname=short sslmode=disable";
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration11$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081 -database-dsn="host=0.0.0.0 port=5433 user=postgres password=postgres dbname=short sslmode=disable";
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration12$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081 -database-dsn="host=0.0.0.0 port=5433 user=postgres password=postgres dbname=short sslmode=disable";
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration13$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081 -database-dsn="host=0.0.0.0 port=5433 user=postgres password=postgres dbname=short sslmode=disable";
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration14$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081 -database-dsn="host=0.0.0.0 port=5433 user=postgres password=postgres dbname=short sslmode=disable";
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration15$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081 -database-dsn="host=0.0.0.0 port=5433 user=postgres password=postgres dbname=short sslmode=disable";
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration16$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081 -database-dsn="host=0.0.0.0 port=5433 user=postgres password=postgres dbname=short sslmode=disable";
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration17$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081 -database-dsn="host=0.0.0.0 port=5433 user=postgres password=postgres dbname=short sslmode=disable";
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration18$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081 -database-dsn="host=0.0.0.0 port=5433 user=postgres password=postgres dbname=short sslmode=disable";
~/bin/shortenertestbeta  -test.v -test.run=^TestIteration19$ -binary-path=cmd/shortener/shortener -source-path=./ -server-port=8081 -database-dsn="host=0.0.0.0 port=5433 user=postgres password=postgres dbname=short sslmode=disable";
