.PHONY: etl
etl:
	@go run main.go

.PHONY: charts
charts: main.R bin/analytics.db
	@Rscript main.R

.PHONY: clean
clean:
	@rm -f bin/*
