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
	DomainPrefixEnv  = "DOMAIN_PREFIX"
	Url              = "url"
	AnnotationPrefix = "uptimerobot-monitor"
)

func DeleteMonitor(host string, ingressAnnotations map[string]string) error {
	return executeMonitorAction(host, ingressAnnotations, model.Delete, monitor.New())
}

func CreateMonitor(host string, ingressAnnotations map[string]string) error {
	return executeMonitorAction(host, ingressAnnotations, model.Update, monitor.New())
}

func executeMonitorAction(host string, ingressAnnotations map[string]string, action model.Args, service service.IService) error {
	annotations, err := buildDataMapFromAnnotations(ingressAnnotations)
	if err != nil {
		return err
	}

	if _, exists := annotations[Url]; !exists {
		annotations[Url] = host
	}

	resultArrayMap := service.HandleRequest([]map[string]interface{}{annotations}, action)
	if resultArrayMap != nil && len(resultArrayMap) > 0 && resultArrayMap[0] != nil && resultArrayMap[0][model.ErrorResultField] != nil {
		return resultArrayMap[0][model.ErrorResultField].(error)
	}
	return nil
}

func buildDataMapFromAnnotations(ingressAnnotations map[string]string) (map[string]interface{}, error) {
	if ingressAnnotations == nil || len(ingressAnnotations) == 0 {
		return nil, errors.New("no ingress annotation provided")
	}
	dataMap := make(map[string]interface{})
	uptimeRobotPrefix := GetUptimeRobotMonitorPrefix()

	for key, value := range ingressAnnotations {
		if strings.HasPrefix(key, uptimeRobotPrefix) {
			_Key := strings.TrimPrefix(key, uptimeRobotPrefix)
			dataMap[_Key] = value
		}
	}
	return dataMap, nil
}

func GetUptimeRobotDomain() string {
	prefixStr := "/" + AnnotationPrefix
	envPrefix := getUptimeRobotDomain()
	if len(envPrefix) == 0 {
		return "my.domain" + prefixStr
	}
	return envPrefix + prefixStr
}

func GetUptimeRobotMonitorPrefix() string {
	return GetUptimeRobotDomain() + "-"
}

func getUptimeRobotDomain() string {
	return os.Getenv(DomainPrefixEnv)
}

type MonitorUtil struct {
}
