package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"kandaoni.com/anqicms/config"
	"os"
)

type Module struct {
	Model
	TableName string       `json:"table_name" gorm:"column:table_name;type:varchar(50) not null;default:''"`
	UrlToken  string       `json:"url_token" gorm:"column:url_token;type:varchar(50) not null;default:''"` // 定义
	Title     string       `json:"title" gorm:"column:title;type:varchar(250) not null;default:''"`
	Fields    moduleFields `json:"fields" gorm:"column:fields;type:longtext default null"`
	IsSystem  int          `json:"is_system" gorm:"column:is_system;type:tinyint(1) unsigned not null;default:0"`
	TitleName string       `json:"title_name" gorm:"column:title_name;type:varchar(50) not null;default:''"`
	Status    uint         `json:"status" gorm:"column:status;type:tinyint(1) unsigned not null;default:0"`
}

type moduleFields []config.CustomField

func (a moduleFields) Value() (driver.Value, error) {
	return json.Marshal(a)
}

func (a *moduleFields) Scan(data interface{}) error {
	return json.Unmarshal(data.([]byte), &a)
}

func (m *Module) Migrate(tx *gorm.DB, focus bool) {
	if !tx.Migrator().HasTable(m.TableName) {
		tx.Exec("CREATE TABLE `?` (`id` int(10) unsigned NOT NULL AUTO_INCREMENT, PRIMARY KEY (`id`)) DEFAULT CHARSET=utf8mb4;", gorm.Expr(m.TableName))
	}
	// 根据表单字段，生成数据
	for _, field := range m.Fields {
		field.CheckSetFilter()
		column := field.GetFieldColumn()
		if !m.HasColumn(tx, field.FieldName) {
			//创建语句
			tx.Exec("ALTER TABLE ? ADD COLUMN ?", gorm.Expr(m.TableName), gorm.Expr(column))
		} else if focus {
			//更新语句
			tx.Exec("ALTER TABLE ? MODIFY COLUMN ?", gorm.Expr(m.TableName), gorm.Expr(column))
		}

		if field.IsFilter {
			idxName := fmt.Sprintf("idx_%s", field.FieldName)
			if !m.HasIndex(tx, idxName) {
				tx.Exec("CREATE INDEX `?` ON `?` (`?`)", gorm.Expr(idxName), gorm.Expr(m.TableName), gorm.Expr(field.FieldName))
			}
		}
	}
	// 检查模板文件夹是否存在，不存在，则创建
	dir := fmt.Sprintf("%stemplate/%s/%s", config.ExecPath, config.JsonData.System.TemplateName, m.TableName)
	_, err := os.Stat(dir)
	if err != nil && os.IsNotExist(err) {
		// 创建文件夹
		os.Mkdir(dir, os.ModePerm)
		// 创建文件
		os.WriteFile(dir + "/detail.html", []byte(m.TableName), os.ModePerm)
		os.WriteFile(dir + "/index.html", []byte(m.TableName), os.ModePerm)
		os.WriteFile(dir + "/list.html", []byte(m.TableName), os.ModePerm)
	}
}

func (m *Module) HasColumn(tx *gorm.DB, field string) bool {
	var count int64
	tx.Raw(
		"SELECT count(*) FROM INFORMATION_SCHEMA.columns WHERE table_schema = ? AND table_name = ? AND column_name = ?",
		config.JsonData.Mysql.Database, m.TableName, field,
	).Row().Scan(&count)

	return count > 0
}

func (m *Module) HasIndex(tx *gorm.DB,name string) bool {
	var count int64
	tx.Raw(
		"SELECT count(*) FROM information_schema.statistics WHERE table_schema = ? AND table_name = ? AND index_name = ?",
		config.JsonData.Mysql.Database, m.TableName, name,
	).Row().Scan(&count)

	return count > 0
}
