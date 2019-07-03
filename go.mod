module github.com/hikhvar/journaldtail

go 1.12

replace github.com/Sirupsen/logrus v1.4.2 => github.com/sirupsen/logrus v1.4.2

replace github.com/hashicorp/consul => github.com/hashicorp/consul v1.5.1

replace github.com/grafana/loki => /home/ry/src/github.com/grafana/loki

require (
	github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/cortexproject/cortex v0.0.0-20190702160911-795dd596d4d8
	github.com/go-ini/ini v1.25.4 // indirect
	github.com/go-kit/kit v0.8.0
	github.com/grafana/loki v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.8.1
	github.com/prometheus/common v0.4.1
)
