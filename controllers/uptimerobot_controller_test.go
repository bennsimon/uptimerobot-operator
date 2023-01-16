package controllers

import (
	"context"
	"errors"
	"github.com/bennsimon/uptimerobot-operator/util/monitorutil"
	"github.com/stretchr/testify/mock"
	network "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sync"
	"testing"
)

type testClient struct {
	client.Client
	mock.Mock
}

type testUtilProvider struct {
	wg sync.WaitGroup
	mock.Mock
}

func (t *testClient) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	if key.Name == "ValidIngress" {
		ingress := obj.(*network.Ingress)
		ingress.Annotations = map[string]string{
			monitorutil.GetUptimeRobotDomain(): "true",
		}
		ingress.Spec.Rules = []network.IngressRule{{Host: "test.localhost"}}
	}
	args := t.Called(ctx, key, obj, opts)
	return args.Error(0)
}

func (r *testUtilProvider) CreateMonitor(host string, annotations map[string]string) error {
	args := r.Called(host, annotations)
	return args.Error(0)
}

func (r *testUtilProvider) DeleteMonitor(host string, annotations map[string]string) error {
	args := r.Called(host, annotations)
	r.wg.Done()
	return args.Error(0)
}

func TestUptimerobotReconciler_Reconcile(t *testing.T) {
	type args struct {
		host        string
		annotations map[string]string
	}
	var testclient *testClient
	var testutilprovider *testUtilProvider
	r := &UptimerobotReconciler{
		Scheme: &runtime.Scheme{},
	}
	tests := []struct {
		name        string
		args        args
		tn          types.NamespacedName
		wantError   bool
		setupMocks  func()
		verifyMocks func()
	}{
		{name: "should return err if get resource has error other than 404", args: args{host: "", annotations: map[string]string{}}, setupMocks: func() {
			testclient = &testClient{}
			testclient.On("Get", mock.IsType(context.Background()), mock.IsType(types.NamespacedName{Namespace: "default", Name: "uptimerobotoperator"}), mock.IsType(&network.Ingress{}), mock.Anything).Return(errors.New(""))
			r.Client = testclient
		}, tn: types.NamespacedName{Namespace: "default", Name: "world"}, verifyMocks: func() {
			testclient.AssertExpectations(t)
		}, wantError: true},
		{name: "should return nil when create monitor action fails with valid ingress", args: args{host: "", annotations: map[string]string{}}, setupMocks: func() {
			testclient = &testClient{}
			testclient.On("Get", mock.IsType(context.Background()), mock.IsType(types.NamespacedName{Namespace: "default", Name: "ValidIngress"}), mock.IsType(&network.Ingress{}), mock.Anything).Return(nil)
			r.Client = testclient

			testutilprovider = &testUtilProvider{}
			testutilprovider.On("CreateMonitor", mock.Anything, mock.IsType(map[string]string{})).Return(errors.New("some error"))
			r.UtilProvider = testutilprovider
		}, tn: types.NamespacedName{Namespace: "default", Name: "ValidIngress"}, verifyMocks: func() {
			testclient.AssertExpectations(t)
			testutilprovider.AssertExpectations(t)
		}, wantError: false},
		{name: "should return nil when create monitor action is successful with valid ingress", args: args{host: "", annotations: map[string]string{}}, setupMocks: func() {
			testclient = &testClient{}
			testclient.On("Get", mock.IsType(context.Background()), mock.IsType(types.NamespacedName{Namespace: "default", Name: "ValidIngress"}), mock.IsType(&network.Ingress{}), mock.Anything).Return(nil)
			r.Client = testclient

			testutilprovider = &testUtilProvider{}
			testutilprovider.On("CreateMonitor", mock.Anything, mock.IsType(map[string]string{})).Return(nil)
			r.UtilProvider = testutilprovider
		}, tn: types.NamespacedName{Namespace: "default", Name: "ValidIngress"}, verifyMocks: func() {
			testclient.AssertExpectations(t)
			testutilprovider.AssertExpectations(t)
		}, wantError: false},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			tt.setupMocks()
			defer tt.verifyMocks()
			_, err := r.Reconcile(context.Background(), ctrl.Request{NamespacedName: tt.tn})

			if err == nil && tt.wantError {
				t.Errorf("want %v got %v", nil, err)
			}
			if err != nil && !tt.wantError {
				t.Errorf("want %v got %v", err, nil)
			}
		})
	}
}

func Test_hasEnabledUptimeRobotMonitor(t *testing.T) {
	r := &UptimerobotReconciler{}
	type args struct {
		annotationMap map[string]string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "should return false due to bool conversion error", args: args{annotationMap: map[string]string{
			monitorutil.GetUptimeRobotDomain(): "falser",
		}}, want: false},
		{name: "should return false when set false", args: args{annotationMap: map[string]string{
			monitorutil.GetUptimeRobotDomain(): "false",
		}}, want: false},
		{name: "should return false due to bool conversion error", args: args{annotationMap: map[string]string{
			monitorutil.GetUptimeRobotDomain(): "falser",
		}}, want: false},
		{name: "should return false config is missing", args: args{annotationMap: map[string]string{
			monitorutil.GetUptimeRobotDomain() + "-": "true",
		}}, want: false},
		{name: "should return true config is set to true", args: args{annotationMap: map[string]string{
			monitorutil.GetUptimeRobotDomain(): "true",
		}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.hasEnabledUptimeRobotMonitor(tt.args.annotationMap); got != tt.want {
				t.Errorf("hasEnabledUptimeRobotMonitor() = %v, want %v", got, tt.want)
			}
		})
	}
}
func Test_buildHostSchemeMap(t *testing.T) {
	type args struct {
		ingress *network.Ingress
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{name: "should return empty map", args: args{ingress: &network.Ingress{}}, want: map[string]string{}},
		{name: "should return non empty map", args: args{ingress: &network.Ingress{Spec: network.IngressSpec{
			Rules: []network.IngressRule{{Host: "test.localhost"}},
		}}}, want: map[string]string{
			"test.localhost": "http",
		}},
		{name: "should override http rule if tls host exists", args: args{ingress: &network.Ingress{Spec: network.IngressSpec{
			Rules: []network.IngressRule{{Host: "test.localhost"}},
			TLS:   []network.IngressTLS{{Hosts: []string{"test.localhost"}}}}}}, want: map[string]string{
			"test.localhost": "https",
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := buildHostSchemeMap(tt.args.ingress); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("buildHostSchemeMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterGenericEvent(t *testing.T) {
	r := &UptimerobotReconciler{}
	type args struct {
		genericEvent event.GenericEvent
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "should return false if annotations empty", args: args{genericEvent: event.GenericEvent{Object: &network.Ingress{
			ObjectMeta: ctrl.ObjectMeta{
				Name:      "Ingress",
				Namespace: "default",
			},
			TypeMeta: ctrl.TypeMeta{
				Kind: "Ingress",
			},
		}}}, want: false},
		{name: "should return false if ingress is not enabled", args: args{genericEvent: event.GenericEvent{Object: &network.Ingress{
			ObjectMeta: ctrl.ObjectMeta{
				Name:      "Ingress",
				Namespace: "default",
				Annotations: map[string]string{
					monitorutil.GetUptimeRobotDomain(): "false",
				},
			},
			TypeMeta: ctrl.TypeMeta{
				Kind: "Ingress",
			},
		}}}, want: false},
		{name: "should return true if ingress is enabled", args: args{genericEvent: event.GenericEvent{Object: &network.Ingress{
			ObjectMeta: ctrl.ObjectMeta{
				Name:      "Ingress",
				Namespace: "default",
				Annotations: map[string]string{
					monitorutil.GetUptimeRobotDomain(): "true",
				},
			},
			TypeMeta: ctrl.TypeMeta{
				Kind: "Ingress",
			},
		}}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.filterGenericEvent(tt.args.genericEvent); got != tt.want {
				t.Errorf("filterGenericEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterCreateEvent(t *testing.T) {
	r := &UptimerobotReconciler{}
	type args struct {
		event event.CreateEvent
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "should return false if ingress is not enabled", args: args{event: event.CreateEvent{Object: &network.Ingress{
			ObjectMeta: ctrl.ObjectMeta{
				Name:      "Ingress",
				Namespace: "default",
				Annotations: map[string]string{
					monitorutil.GetUptimeRobotDomain(): "false",
				},
			},
			TypeMeta: ctrl.TypeMeta{
				Kind: "Ingress",
			},
		}}}, want: false},
		{name: "should return true if ingress is enabled", args: args{event: event.CreateEvent{Object: &network.Ingress{
			ObjectMeta: ctrl.ObjectMeta{
				Name:      "Ingress",
				Namespace: "default",
				Annotations: map[string]string{
					monitorutil.GetUptimeRobotDomain(): "true",
				},
			},
			TypeMeta: ctrl.TypeMeta{
				Kind: "Ingress",
			},
		}}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.filterCreateEvent(tt.args.event); got != tt.want {
				t.Errorf("filterCreateEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filterUpdateEvent(t *testing.T) {
	r := &UptimerobotReconciler{}
	type args struct {
		updateEvent event.UpdateEvent
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "should return false if ingress is not enabled", args: args{updateEvent: event.UpdateEvent{ObjectNew: &network.Ingress{
			ObjectMeta: ctrl.ObjectMeta{
				Name:      "Ingress",
				Namespace: "default",
				Annotations: map[string]string{
					monitorutil.GetUptimeRobotDomain(): "false",
				},
			},
			TypeMeta: ctrl.TypeMeta{
				Kind: "Ingress",
			},
		}}}, want: false},
		{name: "should return true if ingress is enabled", args: args{updateEvent: event.UpdateEvent{ObjectNew: &network.Ingress{
			ObjectMeta: ctrl.ObjectMeta{
				Name:      "Ingress",
				Namespace: "default",
				Annotations: map[string]string{
					monitorutil.GetUptimeRobotDomain(): "true",
				},
			},
			TypeMeta: ctrl.TypeMeta{
				Kind: "Ingress",
			},
		}}}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := r.filterUpdateEvent(tt.args.updateEvent); got != tt.want {
				t.Errorf("filterUpdateEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUptimerobotReconciler_cleanUpAfterIngressDeletion(t *testing.T) {
	var testutilprovider *testUtilProvider
	r := &UptimerobotReconciler{
		Scheme: &runtime.Scheme{},
	}

	type args struct {
		deleteEvent event.DeleteEvent
	}
	tests := []struct {
		name        string
		args        args
		want        bool
		setupMocks  func()
		verifyMocks func()
	}{
		{name: "should return false and call delete monitor if ingress resource enabled and return error", args: args{deleteEvent: event.DeleteEvent{Object: &network.Ingress{
			ObjectMeta: ctrl.ObjectMeta{
				Name:      "Ingress",
				Namespace: "default",
				Annotations: map[string]string{
					monitorutil.GetUptimeRobotDomain(): "true",
				},
			},
			TypeMeta: ctrl.TypeMeta{
				Kind: "Ingress",
			},
		}}}, setupMocks: func() {
			testutilprovider = &testUtilProvider{}
			testutilprovider.wg.Add(1)
			testutilprovider.On("DeleteMonitor", mock.IsType(""), mock.IsType(map[string]string{})).Return(errors.New(""))
			r.UtilProvider = testutilprovider
		}, verifyMocks: func() {
			testutilprovider.AssertExpectations(t)
		}, want: false},
		{name: "should return false and call delete monitor if ingress resource enabled and return nil", args: args{deleteEvent: event.DeleteEvent{Object: &network.Ingress{
			ObjectMeta: ctrl.ObjectMeta{
				Name:      "Ingress",
				Namespace: "default",
				Annotations: map[string]string{
					monitorutil.GetUptimeRobotDomain(): "true",
				},
			},
			TypeMeta: ctrl.TypeMeta{
				Kind: "Ingress",
			},
		}}}, setupMocks: func() {
			testutilprovider = &testUtilProvider{}
			testutilprovider.wg.Add(1)
			testutilprovider.On("DeleteMonitor", mock.IsType(""), mock.IsType(map[string]string{})).Return(nil)
			r.UtilProvider = testutilprovider
		}, verifyMocks: func() {
			testutilprovider.AssertExpectations(t)
		}, want: false},
		{name: "should return false and not call delete monitor if ingress resource disabled", args: args{deleteEvent: event.DeleteEvent{Object: &network.Ingress{
			ObjectMeta: ctrl.ObjectMeta{
				Name:      "Ingress",
				Namespace: "default",
				Annotations: map[string]string{
					monitorutil.GetUptimeRobotDomain(): "false",
				},
			},
			TypeMeta: ctrl.TypeMeta{
				Kind: "Ingress",
			},
		}}}, setupMocks: func() {
			testutilprovider = &testUtilProvider{}
			testutilprovider.On("DeleteMonitor", mock.IsType(""), mock.IsType(map[string]string{})).Return(errors.New(""))
			r.UtilProvider = testutilprovider
		}, verifyMocks: func() {
			testutilprovider.AssertNotCalled(t, "DeleteMonitor", mock.IsType(""), mock.IsType(map[string]string{}))
		}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			defer tt.verifyMocks()
			got := r.filterDeleteEvent(tt.args.deleteEvent)
			testutilprovider.wg.Wait()
			if got != tt.want {
				t.Errorf("filterDeleteEvent() = %v, want %v", got, tt.want)
			}
		})
	}
}
