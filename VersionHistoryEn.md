## v2.0.5
> - Optimize code generator queue creation issues.

## v2.0.4
> - Optimize the code generator queue, events, and listening output format.

## v2.0.3
> - Optimize the code generator to create support table names and database connection names, and automatically generate request structures based on table names.

## v2.0.2
> - Command line style beautification.

## v2.0.1
> - Optimize queue auto loading and optimize command-line output.

## v2.0.0
> - Add facade and provider 
> - Optimize Command Line,log,context,middleware
> - Command line create file template adjustment
> - Complete documentation

## v1.8.7
> - Context optimization, database connection optimization.

## v1.8.6
> - Add request timeout middleware and improve test cases for request timeout and flow limiting middleware.

## v1.8.5
> - Optimize flow limiting middleware, optimize log support for dynamic setter and getter of log levels, optimize login information for user IDs to be stored and retrieved using context keys, and improve middleware usage documentation.

## v1.8.4
> - Validate request optimization.

## v1.8.3
> - Public error codes can be added with error code prefixes and error code test cases.

## v1.8.2
> - Add a time assistant package and corresponding test cases.

## v1.8.1
> - Adjustment of log error level.

## v1.8.0
> - Optimize the specification configuration file to reduce the occurrence of circular dependencies in the later stage.

## v1.7.9
> - Optimize controllers and services, optimize database connections and connection pools.

## v1.7.8
> - Remove old model generation commands and add custom model generation commands.

## v1.7.7
> - Add new test cases.

## v1.7.6
> - Adjust the command tool path to the cmd directory.

## v1.7.5
> - Add database documents and database connections, which can be switched to MySQL, pgSQL, SQLite, and SQLSRV databases.

## v1.7.4
> - Cancel global variables, initialize new containers through bootstrap, bind context through middleware, and obtain container instances wherever there is context. Databases, caches, logs, and configurations can all be obtained through container instances.

## v1.7.3
> - Adjust RabbitMQ to remove unmaintained packages and use new packages

## v1.7.2
> - Optimize flow limiting middleware, add user flow limiting and IP flow limiting maps for automatic cleaning

## v1.7.1
> - Add global exception capture middleware.

## v1.7.0
> - Update the package name of `utils` to pkg, add the bootstrap directory as the boot directory, optimize the code, and improve the documentation.

## v1.6.0
> - Optimize context link logging (SQL, HTTP, listener, Redis, Kafka, RabbitMQ...)

## v1.5.4
> - Optimize logging stack SQL information, HTTP requests redis、kafka、rabbitmq、 Waiting is optional.

## v1.5.3
> - New flow limiting middleware, default flow limiting, user flow limiting, IP flow limiting

## v1.5.2
> - New database configuration supports MySQL, SQLite, pgSQL, and SQLSRV

## v1.5.1
> - Add command line shortcut to create data migration and data filling

## v1.5.0
> - Gorm dynamic query optimization and readme document improvement
> - Release Package v1.5.0

## v1.4.1
> - Command line data migration, adjustment, and optimization

## v1.4.0
> - Model validator command line creation optimization
> - Add Gorm dynamic query
> - Release Package v1.4.0

## v1.3.0
> Improve Kafka and RabbitMQ message queue command line shortcut to create consumers and producers
> Improve the command line to create a message queue document
> Release Package v1.3.0

## v1.2.4
> - Add Kafka and RabbitMQ message queues and configurations
> - New Assistant Function - Tree Structure Generation

## v1.2.3
> - Optimize context processing and log processing, adjust readme document update records

## v1.2.2
> - Optimized the time consumption of logging SQL, Redis, and HTTP

## v1.2.1
> - Optimize context processing
> - Optimize log processing and log processing for Redis, HTTP, MySQL
> - Improved readme document after optimization

## v1.2.0
> - Optimize context processing
> - Optimize log processing
> - Add Messages Release Subscription
> - Improved readme document after optimization

## v1.1.0
> Improve log debugging and user documentation, and complete version v1.0.0.

## v1.0.3
> Improve public response usage documentation.

## v1.0.2
> Error code optimization.

## v1.0.1
> Add the public package function 'FilterFields', adjust the public package function `StructToMap`, and modify the JSON serialization to use the package `go-json`.

## v1.0.0
> Except for incomplete response, error handling, and log documentation, all other updates have been completed.