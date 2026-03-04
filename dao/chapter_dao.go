package dao

import (
	"github.com/yzf120/elysia-backend/model/class"
)

// ChapterDAO 章节数据访问对象
type ChapterDAO interface {
	// 章节操作
	CreateChapter(chapter *class.ClassChapter) error
	GetChapterById(chapterId string) (*class.ClassChapter, error)
	UpdateChapter(chapterId string, updates map[string]interface{}) error
	DeleteChapter(chapterId string) error
	ListChaptersByClassId(classId string) ([]*class.ClassChapter, error)
	BatchUpdateChapterOrder(classId string, orders []ChapterOrder) error

	// 小节操作
	CreateSection(section *class.ClassSection) error
	GetSectionById(sectionId string) (*class.ClassSection, error)
	UpdateSection(sectionId string, updates map[string]interface{}) error
	DeleteSection(sectionId string) error
	ListSectionsByChapterId(chapterId string) ([]*class.ClassSection, error)
	DeleteSectionsByChapterId(chapterId string) error
	BatchUpdateSectionOrder(chapterId string, orders []SectionOrder) error
}

// ChapterOrder 章节排序项
type ChapterOrder struct {
	ChapterId string
	SortOrder int32
}

// SectionOrder 小节排序项
type SectionOrder struct {
	SectionId string
	SortOrder int32
}

type chapterDAOImpl struct{}

// NewChapterDAO 创建章节DAO
func NewChapterDAO() ChapterDAO {
	return &chapterDAOImpl{}
}

// ==================== 章节操作 ====================

// CreateChapter 创建章节
func (d *chapterDAOImpl) CreateChapter(chapter *class.ClassChapter) error {
	return DB.Create(chapter).Error
}

// GetChapterById 根据章节ID查询章节
func (d *chapterDAOImpl) GetChapterById(chapterId string) (*class.ClassChapter, error) {
	var chapter class.ClassChapter
	err := DB.Where("chapter_id = ?", chapterId).First(&chapter).Error
	if err != nil {
		return nil, err
	}
	return &chapter, nil
}

// UpdateChapter 更新章节信息
func (d *chapterDAOImpl) UpdateChapter(chapterId string, updates map[string]interface{}) error {
	return DB.Model(&class.ClassChapter{}).Where("chapter_id = ?", chapterId).Updates(updates).Error
}

// DeleteChapter 删除章节（物理删除）
func (d *chapterDAOImpl) DeleteChapter(chapterId string) error {
	return DB.Where("chapter_id = ?", chapterId).Delete(&class.ClassChapter{}).Error
}

// ListChaptersByClassId 查询班级下所有章节（按 sort_order 升序）
func (d *chapterDAOImpl) ListChaptersByClassId(classId string) ([]*class.ClassChapter, error) {
	var chapters []*class.ClassChapter
	err := DB.Where("class_id = ? AND status = 1", classId).Order("sort_order ASC, id ASC").Find(&chapters).Error
	return chapters, err
}

// BatchUpdateChapterOrder 批量更新章节排序
func (d *chapterDAOImpl) BatchUpdateChapterOrder(classId string, orders []ChapterOrder) error {
	tx := DB.Begin()
	for _, o := range orders {
		if err := tx.Model(&class.ClassChapter{}).
			Where("chapter_id = ? AND class_id = ?", o.ChapterId, classId).
			Update("sort_order", o.SortOrder).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// ==================== 小节操作 ====================

// CreateSection 创建小节
func (d *chapterDAOImpl) CreateSection(section *class.ClassSection) error {
	return DB.Create(section).Error
}

// GetSectionById 根据小节ID查询小节
func (d *chapterDAOImpl) GetSectionById(sectionId string) (*class.ClassSection, error) {
	var section class.ClassSection
	err := DB.Where("section_id = ?", sectionId).First(&section).Error
	if err != nil {
		return nil, err
	}
	return &section, nil
}

// UpdateSection 更新小节信息
func (d *chapterDAOImpl) UpdateSection(sectionId string, updates map[string]interface{}) error {
	return DB.Model(&class.ClassSection{}).Where("section_id = ?", sectionId).Updates(updates).Error
}

// DeleteSection 删除小节（物理删除）
func (d *chapterDAOImpl) DeleteSection(sectionId string) error {
	return DB.Where("section_id = ?", sectionId).Delete(&class.ClassSection{}).Error
}

// ListSectionsByChapterId 查询章节下所有小节（按 sort_order 升序）
func (d *chapterDAOImpl) ListSectionsByChapterId(chapterId string) ([]*class.ClassSection, error) {
	var sections []*class.ClassSection
	err := DB.Where("chapter_id = ? AND status = 1", chapterId).Order("sort_order ASC, id ASC").Find(&sections).Error
	return sections, err
}

// DeleteSectionsByChapterId 删除章节下所有小节
func (d *chapterDAOImpl) DeleteSectionsByChapterId(chapterId string) error {
	return DB.Where("chapter_id = ?", chapterId).Delete(&class.ClassSection{}).Error
}

// BatchUpdateSectionOrder 批量更新小节排序
func (d *chapterDAOImpl) BatchUpdateSectionOrder(chapterId string, orders []SectionOrder) error {
	tx := DB.Begin()
	for _, o := range orders {
		if err := tx.Model(&class.ClassSection{}).
			Where("section_id = ? AND chapter_id = ?", o.SectionId, chapterId).
			Update("sort_order", o.SortOrder).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}
