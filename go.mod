module yaddd

go 1.14

replace pkcs7 => github.com/andviro/pkcs7 v0.0.0-20190605221235-59c41cd306ad

require (
	github.com/robfig/cron/v3 v3.0.1
	github.com/sirupsen/logrus v1.6.0
	go.uber.org/automaxprocs v1.3.0
	golang.org/x/sync v0.0.0-20200317015054-43a5402ce75a
	gopkg.in/yaml.v2 v2.3.0
)
