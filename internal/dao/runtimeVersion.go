package dao

import (
	"sync"

	"github.com/itering/subscan/model"
	"github.com/itering/substrate-api-rpc/metadata"
)

var mutex sync.Mutex

func (d *Dao) CreateRuntimeVersion(name string, specVersion int) int64 {
	mutex.Lock()
	defer mutex.Unlock()
	version, _ := findOne[model.RuntimeVersion](&d.ReadOnlyDao, "id, spec_version", where("spec_version = ?", specVersion), nil)
	if version != nil {
		return 0
	}
	query := d.db.Create(&model.RuntimeVersion{
		Name:        name,
		SpecVersion: specVersion,
	})
	return query.RowsAffected
}

func (d *Dao) SetRuntimeData(specVersion int, modules string, rawData string) int64 {
	query := d.db.Model(model.RuntimeVersion{}).Where("spec_version=?", specVersion).Updates(model.RuntimeVersion{
		Modules: modules,
		RawData: rawData,
	})
	return query.RowsAffected
}

func (d *ReadOnlyDao) RuntimeVersionList() []model.RuntimeVersion {
	var list []model.RuntimeVersion
	d.db.Select("spec_version,modules").Model(model.RuntimeVersion{}).Find(&list)
	return list
}

func (d *ReadOnlyDao) RuntimeVersionRecent() *model.RuntimeVersion {
	var list []model.RuntimeVersion
	_ = d.db.Select("spec_version,raw_data").Model(model.RuntimeVersion{}).Limit(1).Order("spec_version DESC").Find(&list)
	if len(list) == 0 {
		return nil
	}
	version, _ := findOne[model.RuntimeVersion](d, "spec_version,raw_data", whereClauses{}, "spec_version DESC")
	return version
}

func (d *ReadOnlyDao) RuntimeVersionRaw(spec int) *metadata.RuntimeRaw {
	var list []metadata.RuntimeRaw
	d.db.Model(model.RuntimeVersion{}).
		Select("spec_version as spec ,raw_data as raw").
		Where("spec_version = ?", spec).
		Limit(1).
		Find(&list)
	if len(list) == 0 {
		return nil
	}
	return &list[0]
}
