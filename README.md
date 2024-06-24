# flare

<img src="https://media1.tenor.com/m/R56Js9FnFmYAAAAC/flair-office-space.gif" />

Flare is an experiment in using ChatGPT to rewrite an existing open source project in a different language. In this case [Apache Spark](https://spark.apache.org/), which is written largely in Scala, to Go.

Read the [initial conversation](ChatGPT_Getting_Started.md) for an idea of how this started.

A [follow up conversation](ChatGPT_Spark_Connect.md) was started around implementing the [Spark Connect](https://spark.apache.org/docs/latest/spark-connect-overview.html) interface in Go.

## DeltaLake

Using [DeltaLake]() from Go is a driving goal for this package. [Rivian wrote the Delta Connector](https://github.com/rivian/delta-go) we're using here specifically to [share a common logstore with Spark](https://delta.io/blog/rivian-delta-go/). That feels like a hack to me. We're hoping to remove Spark from the equation altogether. If this feels too heavy we could consider others, like [the one supported by treeverse](https://github.com/treeverse/delta-go)
