package model

import "time"

type ServiceVersion struct {
	ID                    uint32                   `json:"id"`
	ServiceID             uint32                   `json:"service_id"`
	BnsID                 uint32                   `json:"bns_id"`
	UUID                  string                   `json:"uuid"`
	BinVersion            string                   `json:"bin_version"`
	ConfVersion           string                   `json:"conf_version"`
	Env                   string                   `json:"env"`
	PackageAddr           string                   `json:"package_addr"`
	SourceLink            string                   `json:"source_link"`
	CollectType           string                   `json:"collect_type"`
	CreatedAt             time.Time                `json:"created_at"`
	UpdatedAt             time.Time                `json:"updated_at"`
	ServiceVersionHistory []*ServiceVersionHistory `json:"service_version_history"`
}
