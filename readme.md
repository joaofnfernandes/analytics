# Sentiment analysis on docs.docker.com

The Docker documentation tracks usage information with:
* Google Analytics
* Polldaddy (thumbs up and down on an article)

This repo allows you to import CSV data from these two sources and puts
the data into a single structured table in sqlite.
It also has some funtionality in R to do some basic statistical analysis
on the data, like distribution or thumbs up per page view.

## Use this repo

Git clone and

```
make
```

This will create a sqlite database that you can use with `sqlite3 bin/analytics.db`,
to do exploratory analysis. It also creates `bin/charts.pdf` with some statistical
analysis diagrams.

The CSV data on this repo is from August 2017, so you should get more recent data
and put it into the `data` directory.

## Under the hood

There are 2 entry points for this project:

* `main.go` parses the CSV data, does some ETL, and puts the data on a relational db
* `main.R` creates statistical analysis diagrams with the data that's on the relational
db
