module github.com/evergreen-ci/evergreen

go 1.16

// We need to keep this old YAML version because upgrading from this specific revision to any newer one somehow breaks
// project validation.
replace gopkg.in/20210107192922/yaml.v3 => gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b

require (
	github.com/99designs/gqlgen v0.17.20
	github.com/PuerkitoBio/rehttp v1.1.0
	github.com/aws/aws-sdk-go v1.44.109
	github.com/cheynewallace/tabby v1.1.1
	github.com/docker/docker v20.10.12+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/dustin/go-humanize v1.0.0
	github.com/evergreen-ci/birch v0.0.0-20220401151432-c792c3d8e0eb
	github.com/evergreen-ci/certdepot v0.0.0-20211117185134-dbedb3d79a10
	github.com/evergreen-ci/cocoa v0.0.0-20221021153157-d4b9d1734588
	github.com/evergreen-ci/gimlet v0.0.0-20220419172609-b882e01673e7
	github.com/evergreen-ci/go-test2json v0.0.0-20180702150328-5b6cfd2e8cb0
	github.com/evergreen-ci/juniper v0.0.0-20220118233332-0813edc78908
	github.com/evergreen-ci/pail v0.0.0-20220908201135-8a2090a672b7
	github.com/evergreen-ci/poplar v0.0.0-20220908212406-a5e2aa799def
	github.com/evergreen-ci/shrub v0.0.0-20211025143051-a8d91b2e29fd
	github.com/evergreen-ci/timber v0.0.0-20220809215118-028567bbef72
	github.com/evergreen-ci/utility v0.0.0-20220725171106-4730479c6118
	github.com/google/go-github/v34 v34.0.0
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510
	github.com/gophercloud/gophercloud v0.1.0
	github.com/gorilla/csrf v1.7.1
	github.com/gorilla/sessions v1.2.1
	github.com/jpillora/backoff v1.0.0
	github.com/kardianos/osext v0.0.0-20190222173326-2bc1f35cddc0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.4.2
	github.com/mongodb/amboy v0.0.0-20221024224945-afcee6f9ec39
	github.com/mongodb/anser v0.0.0-20220318141853-005b8ead5b8f
	github.com/mongodb/ftdc v0.0.0-20220401165013-13e4af55e809
	github.com/mongodb/grip v0.0.0-20221003151935-9c3eaaa08384
	github.com/pkg/errors v0.9.1
	github.com/robbiet480/go.sns v0.0.0-20210223081447-c7c9eb6836cb
	github.com/robfig/cron v1.2.0
	github.com/sabhiram/go-gitignore v0.0.0-20210923224102-525f6e181f06
	github.com/smartystreets/goconvey v1.7.2
	github.com/stretchr/testify v1.8.0
	github.com/urfave/cli v1.22.10
	github.com/vektah/gqlparser/v2 v2.5.1
	github.com/vmware/govmomi v0.27.1
	go.mongodb.org/mongo-driver v1.10.2
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	golang.org/x/oauth2 v0.0.0-20211104180415-d3ed0bb246c8
	golang.org/x/tools v0.1.12
	gonum.org/v1/gonum v0.12.0
	google.golang.org/api v0.60.0
	google.golang.org/grpc v1.49.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/mholt/archiver v2.0.1-0.20180417220235-e4ef56d48eb0+incompatible
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
)

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/evergreen-ci/evg-lint v0.0.0-20211115144425-3b19c8e83a57
	github.com/evergreen-ci/tarjan v0.0.0-20170824211642-fcd3f3321826
	github.com/mongodb/jasper v0.0.0-20220214215554-82e5a72cff6b
	github.com/shirou/gopsutil/v3 v3.22.9
	github.com/trinodb/trino-go-client v0.300.0
	google.golang.org/genproto v0.0.0-20211129164237-f09f9a12af12 // indirect
	gopkg.in/20210107192922/yaml.v3 v3.0.0-00010101000000-000000000000
)
