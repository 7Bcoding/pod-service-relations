package test_client

import (
	"context"
	"k8s.io/client-go/kubernetes"
	"pod-service-relations/client"
	"reflect"
	"testing"
)

func TestGetConfigMaps(t *testing.T) {
	type args struct {
		ctx             context.Context
		masterURL       string
		clusterFilePath string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.GetConfigMaps(tt.args.ctx, tt.args.masterURL, tt.args.clusterFilePath)
		})
	}
}

func TestGetConfigMapsDataMapping(t *testing.T) {
	type args struct {
		ctx             context.Context
		masterURL       string
		clusterFilePath string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.GetConfigMapsDataMapping(tt.args.ctx, tt.args.masterURL, tt.args.clusterFilePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetConfigMapsDataMapping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPodConditionStatus(t *testing.T) {
	type args struct {
		ctx             context.Context
		masterURL       string
		clusterFilePath string
		namespace       string
		podName         string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.GetPodConditionStatus(tt.args.ctx, tt.args.masterURL, tt.args.clusterFilePath, tt.args.namespace, tt.args.podName); got != tt.want {
				t.Errorf("GetPodConditionStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPodNamespaceMap(t *testing.T) {
	type args struct {
		ctx             context.Context
		masterURL       string
		clusterFilePath string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.GetPodNamespaceMap(tt.args.ctx, tt.args.masterURL, tt.args.clusterFilePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPodNamespaceMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPodServiceMap(t *testing.T) {
	type args struct {
		ctx             context.Context
		masterURL       string
		clusterFilePath string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.GetPodServiceMap(tt.args.ctx, tt.args.masterURL, tt.args.clusterFilePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPodServiceMap() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPodStatusPhase(t *testing.T) {
	type args struct {
		ctx             context.Context
		masterURL       string
		clusterFilePath string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.GetPodStatusPhase(tt.args.ctx, tt.args.masterURL, tt.args.clusterFilePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPodStatusPhase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPodVolumeConfigMapMapping(t *testing.T) {
	type args struct {
		ctx             context.Context
		masterURL       string
		clusterFilePath string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.GetPodVolumeConfigMapMapping(tt.args.ctx, tt.args.masterURL, tt.args.clusterFilePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPodVolumeConfigMapMapping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPods(t *testing.T) {
	type args struct {
		ctx             context.Context
		masterURL       string
		clusterFilePath string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client.GetPods(tt.args.ctx, tt.args.masterURL, tt.args.clusterFilePath)
		})
	}
}

func TestGetPodsContainersEnv(t *testing.T) {
	type args struct {
		ctx             context.Context
		masterURL       string
		clusterFilePath string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.GetPodsContainersEnv(tt.args.ctx, tt.args.masterURL, tt.args.clusterFilePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetPodsContainersEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceLabelSelectorMapping(t *testing.T) {
	type args struct {
		ctx             context.Context
		masterURL       string
		clusterFilePath string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.GetServiceLabelSelectorMapping(tt.args.ctx, tt.args.masterURL, tt.args.clusterFilePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetServiceLabelSelectorMapping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetServiceList(t *testing.T) {
	type args struct {
		ctx             context.Context
		masterURL       string
		clusterFilePath string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.GetServiceList(tt.args.ctx, tt.args.masterURL, tt.args.clusterFilePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetServiceList() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewKubeClientSet(t *testing.T) {
	type args struct {
		ctx             context.Context
		masterURL       string
		clusterFilePath string
	}
	tests := []struct {
		name string
		args args
		want *kubernetes.Clientset
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := client.NewKubeClientSet(tt.args.ctx, tt.args.masterURL, tt.args.clusterFilePath); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewKubeClientSet() = %v, want %v", got, tt.want)
			}
		})
	}
}
