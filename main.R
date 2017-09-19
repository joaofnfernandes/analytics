library(RSQLite)

# Open device
pdf("bin/charts.pdf")

db <- dbConnect(SQLite(), "bin/analytics.db")
sqlStmt <- "select url, pageviews, avg_time, bounce_rate, rating from page where rating not null order by pageviews desc"
table <- dbGetQuery(db, sqlStmt)

# Draw 2 diagrams side by side
par(mfrow=c(nrows=2, ncols=2))

# Sentiment boxplot
boxplot(table$rating)
title(sprintf("Sentiment distribution (N=%s)", length(table$rating)))

# Ratings per Pageviews
plot(table$rating, table$pageviews, main="Sentiment per page views",log="y", xlab="Sentiment", ylab="Page views")

# Ratings per Average time on page
plot(table$rating, table$avg_time, main="Sentiment per average time on page", xlab="Sentiment", ylab="Avg time on page (s)")

# Ratings per bounce rate
plot(table$rating, table$bounce_rate, main="Sentiment per bounce rate", xlab="Sentiment", ylab="Bounce rate")

# ------------------- Docker EE------------
sqlStmt <- "select url, pageviews, avg_time, bounce_rate, rating from page where rating not null and url like '/datacenter%' order by pageviews desc"
table <- dbGetQuery(db, sqlStmt)

# Sentiment boxplot
boxplot(table$rating)
title(sprintf("Sentiment distribution (N=%s)", length(table$rating)))

# Ratings per Pageviews
plot(table$rating, table$pageviews, main="Sentiment per page views",log="y", xlab="Sentiment", ylab="Page views")

# Ratings per Average time on page
plot(table$rating, table$avg_time, main="Sentiment per average time on page", xlab="Sentiment", ylab="Avg time on page (s)")

# Ratings per bounce rate
plot(table$rating, table$bounce_rate, main="Sentiment per bounce rate", xlab="Sentiment", ylab="Bounce rate")

# Close device
dev.off()
dbDisconnect(db)
