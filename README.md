# flare

<img src="https://media1.tenor.com/m/R56Js9FnFmYAAAAC/flair-office-space.gif" />

Flare is an experiment in using ChatGPT to rewrite an existing open source project in a different language. In this case [Apache Spark](https://spark.apache.org/), which is written largely in Scala, to Go.

Read the [initial conversation](ChatGPT_Getting_Started.md) for an idea of how this started.

A [follow up conversation](ChatGPT_Spark_Connect.md) was started around implementing the [Spark Connect](https://spark.apache.org/docs/latest/spark-connect-overview.html) interface in Go.

## DeltaLake

Using [DeltaLake](https://delta.io/) from Go is a driving goal for this package. [Rivian wrote the Delta Connector](https://github.com/rivian/delta-go) we're using here specifically to [share a common logstore with Spark](https://delta.io/blog/rivian-delta-go/). That feels like a hack. We're hoping to remove Spark from the equation altogether. If this feels too heavy we could consider others, like [the one supported by treeverse](https://github.com/treeverse/delta-go)

## Spark

* [Architecture and Application Lifecycle](https://www.systemsltd.com/blogs/apache-spark-architecture-and-application-lifecycle)

## Arrow

[github.com/apache/arrow](https://github.com/apache/arrow)

## CoPilot prompts and results

* [connect/sessions](/session.copilot.md)
