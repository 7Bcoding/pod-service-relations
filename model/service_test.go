package model

import (
	"reflect"
	"testing"
	"time"

	"gorm.io/datatypes"
)

func TestService_GetBnses(t *testing.T) {
	type fields struct {
		ID              uint32
		Name            string
		Arch            string
		OS              string
		Tag             datatypes.JSON
		SyncType        string
		BinRepo         string
		BinType         string
		BinFtpAddr      string
		ConfRepo        string
		ConfType        string
		ConfFtpAddr     string
		CreatedAt       time.Time
		UpdatedAt       time.Time
		ServiceVersions []*ServiceVersion
		ProductServices []*ProductService
	}
	tests := []struct {
		name   string
		fields fields
		want   []string
	}{
		// TODO: Add test cases.
		{
			name:   "empty test",
			fields: fields{ProductServices: []*ProductService{}},
			want:   []string{},
		},
		{
			name: "normal test",
			fields: fields{
				ProductServices: []*ProductService{
					{
						Product: &Product{
							Bnses: []*Bns{
								{
									Name: "bns1",
								},
							},
						},
					},
					{
						Product: &Product{
							Bnses: []*Bns{
								{
									Name: "bns1",
								},
							},
						},
					},
				},
			},
			want: []string{"bns1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Service{
				ID:              tt.fields.ID,
				Name:            tt.fields.Name,
				Arch:            tt.fields.Arch,
				OS:              tt.fields.OS,
				Tag:             tt.fields.Tag,
				SyncType:        tt.fields.SyncType,
				BinRepo:         tt.fields.BinRepo,
				BinType:         tt.fields.BinType,
				BinFtpAddr:      tt.fields.BinFtpAddr,
				ConfRepo:        tt.fields.ConfRepo,
				ConfType:        tt.fields.ConfType,
				ConfFtpAddr:     tt.fields.ConfFtpAddr,
				CreatedAt:       tt.fields.CreatedAt,
				UpdatedAt:       tt.fields.UpdatedAt,
				ServiceVersions: tt.fields.ServiceVersions,
				ProductServices: tt.fields.ProductServices,
			}
			if got := s.GetBnses(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBnses() = %v, want %v", got, tt.want)
			}
		})
	}
}
