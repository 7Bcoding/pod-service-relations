package test_dao

import (
	"gorm.io/gorm"
	"pod-service-relations/dao"
	"pod-service-relations/model"
	"reflect"
	"testing"
)

func TestNewAbnormalPodDao(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want abnormalPodDao
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dao.NewAbnormalPodDao(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAbnormalPodDao() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_abnormalPodDao_ALL(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		filters map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.AbnormalPod
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := abnormalPodDao{
				db: tt.fields.db,
			}
			got, err := ap.ALL(tt.args.filters)
			if (err != nil) != tt.wantErr {
				t.Errorf("ALL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ALL() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_abnormalPodDao_Create(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		AbnormalPod *model.AbnormalPod
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := abnormalPodDao{
				db: tt.fields.db,
			}
			if err := ap.Create(tt.args.AbnormalPod); (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_abnormalPodDao_Delete(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		id   uint
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := abnormalPodDao{
				db: tt.fields.db,
			}
			if err := ap.Delete(tt.args.id, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_abnormalPodDao_Get(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		PodName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.AbnormalPod
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := abnormalPodDao{
				db: tt.fields.db,
			}
			got, err := ap.Get(tt.args.PodName)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_abnormalPodDao_Update(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		Pod    *model.AbnormalPod
		fields map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := abnormalPodDao{
				db: tt.fields.db,
			}
			if err := ap.Update(tt.args.Pod, tt.args.fields); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
