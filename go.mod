module github.com/jinmukeji/jiujiantang-services

go 1.14

replace (
	github.com/mozillazg/go-pinyin => github.com/mozillazg/go-pinyin v0.15.0
	// FIXME: 由于 etcd 与 gRPC 的兼容问题，得降级 grpc 版本
	// https://github.com/etcd-io/etcd/issues/11721
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)

require (
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535
	github.com/aws/aws-sdk-go v1.30.17
	github.com/blang/semver v3.5.1+incompatible
	github.com/chanxuehong/rand v0.0.0-20180830053958-4b3aff17f488 // indirect
	github.com/chanxuehong/util v0.0.0-20200304121633-ca8141845b13 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-gomail/gomail v0.0.0-20160411212932-81ebce5c23df
	github.com/go-sql-driver/mysql v1.5.0
	github.com/golang/protobuf v1.4.0
	github.com/google/uuid v1.1.1
	github.com/iris-contrib/middleware/cors v0.0.0-20191219204441-78279b78a367
	github.com/jinmukeji/ae-v1 v1.0.2
	github.com/jinmukeji/ae/v2 v2.10.7
	github.com/jinmukeji/go-pkg/v2 v2.2.7
	github.com/jinmukeji/plat-pkg/v2 v2.2.0
	github.com/jinmukeji/proto/v3 v3.0.8
	github.com/jinzhu/gorm v1.9.12
	github.com/jpillora/ipfilter v1.2.1
	github.com/kataras/iris/v12 v12.1.8
	github.com/micro/cli/v2 v2.1.2
	github.com/micro/go-micro/v2 v2.5.0
	github.com/mozillazg/go-pinyin v0.0.0-00010101000000-000000000000
	github.com/sirupsen/logrus v1.5.0
	golang.org/x/net v0.0.0-20200425230154-ff2c4b7c35a0
	gopkg.in/chanxuehong/wechat.v2 v2.0.0-20190402080805-fa408c6cc20d
	gopkg.in/yaml.v2 v2.2.8
)
