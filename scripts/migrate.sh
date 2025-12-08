#!/bin/bash

migrate -path migrations -database "postgres://tanmoreuser:tanmorepass@localhost:5454/tanmoredb?sslmode=disable" up
