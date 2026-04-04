package dao

import (
	"github.com/yzf120/elysia-backend/model/class"
)

// MaterialDAO 学习资料数据访问对象
type MaterialDAO interface {
	CreateMaterial(material *class.SectionMaterial) error
	GetMaterialById(materialId string) (*class.SectionMaterial, error)
	ListMaterialsBySectionId(sectionId string) ([]*class.SectionMaterial, error)
	DeleteMaterial(materialId string) error
	DeleteMaterialsBySectionId(sectionId string) error
}

type materialDAOImpl struct{}

// NewMaterialDAO 创建学习资料DAO
func NewMaterialDAO() MaterialDAO {
	return &materialDAOImpl{}
}

// CreateMaterial 创建学习资料
func (d *materialDAOImpl) CreateMaterial(material *class.SectionMaterial) error {
	return DB.Create(material).Error
}

// GetMaterialById 根据资料ID查询
func (d *materialDAOImpl) GetMaterialById(materialId string) (*class.SectionMaterial, error) {
	var material class.SectionMaterial
	err := DB.Where("material_id = ?", materialId).First(&material).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

// ListMaterialsBySectionId 查询小节下所有学习资料（按 sort_order 升序）
func (d *materialDAOImpl) ListMaterialsBySectionId(sectionId string) ([]*class.SectionMaterial, error) {
	var materials []*class.SectionMaterial
	err := DB.Where("section_id = ? AND status = 1", sectionId).Order("sort_order ASC, id ASC").Find(&materials).Error
	return materials, err
}

// DeleteMaterial 删除单条学习资料
func (d *materialDAOImpl) DeleteMaterial(materialId string) error {
	return DB.Where("material_id = ?", materialId).Delete(&class.SectionMaterial{}).Error
}

// DeleteMaterialsBySectionId 删除小节下所有学习资料
func (d *materialDAOImpl) DeleteMaterialsBySectionId(sectionId string) error {
	return DB.Where("section_id = ?", sectionId).Delete(&class.SectionMaterial{}).Error
}
