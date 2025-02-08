package model

import "time"

const (
	PodServiceTABLE = "pod_services"
)

type PodService struct {
	ID           uint      `gorm:"primaryKey;AUTO_INCREMENT" json:"id"`
	PodName      string    `gorm:"not null;index" json:"pod_name"`
	Cluster      string    `json:"cluster"`
	PodStatus    string    `gorm:"index" json:"pod_status"`
	Namespace    string    `json:"namespace"`
	ServiceName  string    `json:"service_name"`
	PodToService string    `gorm:"type:longtext;comment:pod调用的service" json:"pod_to_service"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type AbnormalPod struct {
	ID             uint      `gorm:"primaryKey;AUTO_INCREMENT" json:"id"`
	PodName        string    `gorm:"not null;index" json:"pod_name"`
	Cluster        string    `json:"cluster"`
	Namespace      string    `json:"namespace"`
	ServiceName    string    `json:"service_name"`
	AffectPods     string    `gorm:"type:longtext;comment:影响的pod" json:"affect_pods"`
	AffectServices string    `gorm:"type:longtext;comment:影响的service" json:"affect_services"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
