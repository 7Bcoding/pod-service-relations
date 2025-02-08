package model

import (
	"fmt"
	"time"

	"gorm.io/datatypes"

	"pod-service-relations/utils"
)

type Service struct {
	ID              uint32            `json:"id"`
	Name            string            `json:"name"`
	Arch            string            `json:"arch"`
	OS              string            `json:"os"`
	Tag             datatypes.JSON    `json:"tag"`
	ModuleName      string            `json:"module_name"`
	SyncType        string            `json:"sync_type"`
	BinRepo         string            `json:"bin_repo"`
	BinType         string            `json:"bin_type"`
	BinFtpAddr      string            `json:"bin_ftp_addr"`
	ConfRepo        string            `json:"conf_repo"`
	ConfType        string            `json:"conf_type"`
	ConfFtpAddr     string            `json:"conf_ftp_addr"`
	CreatedAt       time.Time         `json:"created_at"`
	UpdatedAt       time.Time         `json:"updated_at"`
	ServiceVersions []*ServiceVersion `json:"service_versions"`
	ProductServices []*ProductService `json:"product_services"`
}

func (s Service) GetBnses() []string {
	bnsSet := utils.NewSet()
	for _, ps := range s.ProductServices {
		for _, b := range ps.Product.Bnses {
			bnsSet.Add(b.Name)
		}
	}
	res := make([]string, bnsSet.Size())
	for i, v := range bnsSet.Elements() {
		res[i] = fmt.Sprint(v)
	}
	return res
}
