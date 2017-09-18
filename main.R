library(RSQLite)

db <- dbConnect(SQLite(), "bin/analytics.db")
sqlStmt <- "select url, pageviews, avg_time, bounce_rate, rating from page where rating not null order by pageviews desc"
table <- dbGetQuery(db, sqlStmt)

# Draw 2 diagrams side by side
par(mfrow=c(nrows=2, ncols=2))

# Ratings per Pageviews
plot(table$rating, table$pageviews, xlab="Sentiment", ylab="Page views")
abline(lm(table$pageviews ~ table$rating))

# Ratings per Average time on page
plot(table$rating, table$avg_time, xlab="Sentiment", ylab="Avg time on page (s)")
abline(lm(table$avg_time ~ table$rating))

# Ratings per bounce rate
plot(table$rating, table$bounce_rate, xlab="Sentiment", ylab="Bounce rate")
abline(lm(table$bounce_rate ~ table$rating))
