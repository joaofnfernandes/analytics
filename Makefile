.PHONY: all
all: bin/analytics.db bin/charts.pdf

bin/analytics.db:
	@go run main.go

bin/charts.pdf: main.R bin/analytics.db
	@Rscript main.R

.PHONY: clean
clean:
	@rm -f bin/*
