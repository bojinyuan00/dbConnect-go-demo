package drivers

import (
	"fmt"
	"github.com/godoes/gorm-dameng"
	"gorm.io/gorm"
)

// NewDmDialector 创建达梦数据库方言
func NewDmDialector(dsn string) gorm.Dialector {
	//options := map[string]string{
	//	"schema":         "SYSDBA",
	//	"appName":        "GORM 连接达梦数据库示例",
	//	"connectTimeout": "30000",
	//}

	// dm://user:password@host:port?schema=SYSDBA[&...]
	//dsn = dameng.BuildUrl("SYSDBA", "admin123456", "127.0.0.1", 5236, options)
	fmt.Println(dsn)
	return dameng.Open(dsn)
}

//// VARCHAR 类型大小为字符长度
////db, err := gorm.Open(dameng.New(dameng.Config{DSN: dsn, VarcharSizeIsCharLength: true}))
//// VARCHAR 类型大小为字节长度（默认）
//db, err := gorm.Open(dameng.Open(dsn), &gorm.Config{})
//if err != nil {
//	// panic error or log error info
//}
//
//// do somethings
//var versionInfo []map[string]interface{}
//db.Table("SYS.V$VERSION").Find(&versionInfo)
//if err := db.Error; err == nil {
//	versionBytes, _ := json.MarshalIndent(versionInfo, "", "  ")
//	fmt.Printf("达梦数据库版本信息：\n%s\n", versionBytes)
//}

/****************** 控制台输出内容 *****************

达梦数据库版本信息：
[
  {
    "BANNER": "DM Database Server 64 V8"
  },
  {
    "BANNER": "DB Version: 0x7000c"
  },
  {
    "BANNER": "03134284094-20230927-******-*****"
  }
]

*************************************************/
