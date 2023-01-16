package monitorutil

import (
	"errors"
	"github.com/bennsimon/uptimerobot-tooling/pkg/model"
	"github.com/bennsimon/uptimerobot-tooling/pkg/service"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

type MockMonitorService struct {
	mock.Mock
	service.IService
}

func (m *MockMonitorService) HandleRequest(dataMapInterface []map[string]interface{}, action model.Args) []map[string]interface{} {
	return m.Called(dataMapInterface, action).Get(0).([]map[string]interface{})
}

func TestExecuteMonitorActionShouldReturnErrorWhenIngressAnnotationIsNil(t *testing.T) {
	err := executeMonitorAction("", nil, model.Update, &service.NopService{})
	if err == nil {
		t.Errorf("got %v ,  want %v", nil, err)
	}
}

func TestExecuteMonitorActionShouldReturnErrorWhenAnnotationsIsEmpty(t *testing.T) {
	err := executeMonitorAction("", map[string]string{}, model.Update, &service.NopService{})
	if err == nil {
		t.Errorf("got %v ,  want %v", nil, err)
	}
}

func TestExecuteMonitorActionShouldReturnNil(t *testing.T) {
	err := executeMonitorAction("", map[string]string{
		GetUptimeRobotDomain(): "true",
	}, model.Update, &service.NopService{})
	if err != nil {
		t.Errorf("got %v ,  want %v", err, nil)
	}
}

func TestExecuteMonitorActionShouldErrorOnResult(t *testing.T) {
	testStruct := new(MockMonitorService)
	testStruct.On("HandleRequest", mock.IsType([]map[string]interface{}{}), mock.IsType(model.Args(""))).Return([]map[string]interface{}{
		{model.ErrorResultField: errors.New("some error")},
	})
	err := executeMonitorAction("", map[string]string{
		GetUptimeRobotDomain(): "true",
	}, model.Update, testStruct)
	if err == nil {
		t.Errorf("got %v ,  want %v", nil, err)
	}
}

func Test_buildDataMapFromAnnotations(t *testing.T) {
	type args struct {
		ingressAnnotations map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{name: "should return error if annotations are nil", args: args{ingressAnnotations: nil}, want: nil, wantErr: true},
		{name: "should return error if annotations are empty", args: args{ingressAnnotations: map[string]string{}}, want: nil, wantErr: true},
		{name: "should build map successfully", args: args{ingressAnnotations: map[string]string{
			GetUptimeRobotDomain():                 "true",
			GetUptimeRobotMonitorPrefix() + "type": "HTTP",
		}}, want: map[string]interface{}{
			"type": "HTTP",
		}, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildDataMapFromAnnotations(tt.args.ingressAnnotations)
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildDataMapFromAnnotations() = %v, want %v", got, tt.want)
			}

			if err != nil && !tt.wantErr {
				t.Errorf("buildDataMapFromAnnotations() = %v, want %v", err, tt.wantErr)
			}
		})
	}
}
