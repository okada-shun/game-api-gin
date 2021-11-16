module local.packages/database

go 1.17

require (
	gorm.io/driver/mysql v1.1.3
	gorm.io/gorm v1.22.2
	local.packages/config v0.0.0-00010101000000-000000000000
	local.packages/model v0.0.0-00010101000000-000000000000
)

require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.2 // indirect
)

replace local.packages/config => ../config

replace local.packages/model => ../model

replace local.packages/util => ../util
