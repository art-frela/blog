# blog is a simple web site with blog functions

## TODO

* [x] add MySQL storage
* [x] separete api methods and web
* [x] add MongoDB storage
* [-] add using BeeGo framework
* [ ] add config app
* [ ] add clear logging
* [ ] add metrics (Prometheus), business and techics
* [ ] GUI refactoring
* [ ] add user register
* [ ] add content valid and moderate
* [ ] add authoriazation from social networks
* [ ] etc...

### How to use

![important]

For storage your posts you can use MySQL or MongoDB.

#### MySQL

you must have connection to MySQL server > 8.X  
for creating database `blog` you need privileges to drop and create database and tables  

#### MongoDB

you must have connection to MongoDB server > 3.X  

- install to $GOPATH `go get -u github.com/art-frela/blog`
- help `blog -h`
- prepare database in MySQL
    - use mysql-client `mysql < ./db/blog.sql`
    - OR use go `go run ./db/dbmigrate.sql`
- set Storage connection string for MySQL `export DATABASE_URL=root:master@tcp(localhost:3306)/blog?parseTime=true`, for MongoDB `export DATABASE_URL=mongodb://locaslhost:27017`
- navigate to directory with `asset` folder  
- start `blog`
