.PHONY: etl
etl:
	@go run main.go

.PHONY: clean
clean:
	@rm -f bin/analytics.db
