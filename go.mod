module github.com/LucasRoesler/openfaas-loki

go 1.14

require (
	github.com/fsnotify/fsnotify v1.4.9 // indirect
	github.com/go-chi/chi v4.1.1+incompatible
	github.com/golang/protobuf v1.4.2 // indirect
	github.com/google/go-cmp v0.4.1 // indirect
	github.com/grafana/loki v1.5.0
	github.com/mitchellh/mapstructure v1.3.1 // indirect
	github.com/openfaas/faas-provider v0.15.1
	github.com/pelletier/go-toml v1.8.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/prometheus v2.5.0+incompatible // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/cobra v1.0.0
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0-20200520182314-0ba52f642ac2 // indirect
	golang.org/x/sys v0.0.0-20200523222454-059865788121 // indirect
	google.golang.org/genproto v0.0.0-20200527145253-8367513e4ece // indirect
	google.golang.org/grpc v1.29.1 // indirect
	gopkg.in/ini.v1 v1.56.0 // indirect
)

// Override reference that causes an error from Go proxy - see https://github.com/golang/go/issues/33558
// copied from loki
replace k8s.io/client-go => k8s.io/client-go v0.0.0-20190620085101-78d2af792bab
