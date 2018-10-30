module im

require (
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/denisenkom/go-mssqldb v0.0.0-20181014144952-4e0d7dc8888f // indirect
	github.com/garyburd/redigo v2.0.0+incompatible
	github.com/go-sql-driver/mysql v1.4.0
	github.com/go-xorm/xorm v0.7.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/protobuf v1.2.0
	github.com/gomodule/redigo v2.0.0+incompatible // indirect
	github.com/lib/pq v1.0.0 // indirect
	github.com/mattn/go-sqlite3 v1.9.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/viper v1.2.1
	github.com/stretchr/testify v1.2.2 // indirect
	github.com/ziutek/mymysql v1.5.4 // indirect
	golang.org/x/crypto v0.0.0-20181015023909-0c41d7ab0a0e // indirect
	gopkg.in/mgo.v2 v2.0.0-20180705113604-9856a29383ce
)

replace github.com/garyburd/redigo v2.0.0+incompatible => github.com/gomodule/redigo v1.6.0+incompatible
