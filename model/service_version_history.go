package model

import "time"

type ServiceVersionHistory struct {
	ID               uint32    `json:"id"`
	ServiceVersionID uint32    `json:"service_version_id"`
	OldBinVersion    string    `json:"old_bin_version"`
	OldConfVersion   string    `json:"old_conf_version"`
	OldEnv           string    `json:"old_env"`
	OldPackageAddr   string    `json:"old_package_addr"`
	OldSource        string    `json:"old_source"`
	NewBinVersion    string    `json:"new_bin_version"`
	NewConfVersion   string    `json:"new_conf_version"`
	NewEnv           string    `json:"new_env"`
	NewSource        string    `json:"new_source"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
