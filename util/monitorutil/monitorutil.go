package monitorutil

import (
	"errors"
	"github.com/bennsimon/uptimerobot-tooling/pkg/model"
	"github.com/bennsimon/uptimerobot-tooling/pkg/service"
	"github.com/bennsimon/uptimerobot-tooling/pkg/service/monitor"
	"os"
	"strings"
)

const (
	DomainLabelPrefixEnv = "DOMAIN_LABEL_PREFIX"
	Url                  = "url"
	LabelPrefix          = "uptimerobot-monitor"
)

func DeleteMonitor(host string, ingressLabels map[string]string) error {
	return executeMonitorAction(host, ingressLabels, model.Delete, monitor.New())
}

func CreateMonitor(host string, ingressLabels map[string]string) error {
	return executeMonitorAction(host, ingressLabels, model.Update, monitor.New())
}

func executeMonitorAction(host string, ingressLabels map[string]string, action model.Args, service service.IService) error {
	labels, err := buildDataMapFromLabels(ingressLabels)
	if err != nil {
		return err
	}
	labels[Url] = host

	resultArrayMap := service.HandleRequest([]map[string]interface{}{labels}, action)
	if resultArrayMap != nil && len(resultArrayMap) > 0 && resultArrayMap[0] != nil && resultArrayMap[0][model.ErrorResultField] != nil {
		return resultArrayMap[0][model.ErrorResultField].(error)
	}
	return nil
}

func buildDataMapFromLabels(ingressLabels map[string]string) (map[string]interface{}, error) {
	if ingressLabels == nil || len(ingressLabels) == 0 {
		return nil, errors.New("no ingress label provided")
	}
	dataMap := make(map[string]interface{})
	uptimeRobotLabelPrefix := GetUptimeRobotMonitorLabelPrefix()

	for labelKey, labelValue := range ingressLabels {
		if strings.HasPrefix(labelKey, uptimeRobotLabelPrefix) {
			_labelKey := strings.TrimPrefix(labelKey, uptimeRobotLabelPrefix)
			dataMap[_labelKey] = labelValue
		}
	}
	return dataMap, nil
}

func GetUptimeRobotLabelDomain() string {
	prefixStr := "/" + LabelPrefix
	envPrefix := getUptimeRobotDomain()
	if len(envPrefix) == 0 {
		return "my.domain" + prefixStr
	}
	return envPrefix + prefixStr
}

func GetUptimeRobotMonitorLabelPrefix() string {
	return GetUptimeRobotLabelDomain() + "-"
}

func getUptimeRobotDomain() string {
	return os.Getenv(DomainLabelPrefixEnv)
}

type MonitorUtil struct {
}
