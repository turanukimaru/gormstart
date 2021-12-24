package dummydb

import (
	"fmt"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// Product gorm サンプルコード
type Product struct {
	gorm.Model
	Code  string
	Price uint
}

// DummyTable テスト用ダミーテーブル
type DummyTable struct {
	gorm.Model
	Dum uint
	Tex string
}

// github.com/denisenkom/go-mssqldb
func DbAccess() error {
	dsn := "sqlserver://gorm:pass@127.0.0.1:1433?database=gorm"
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	// ただ権限の関係か SQL server では機能しない…
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	// migration しないと auto increment の都合上 id が 1 でなくなってしまう…
	db.First(&product, 1)                 // find product with integer primary key
	db.First(&product, "code = ?", "D42") // find product with code D42

	// Update - update product's price to 200
	db.Model(&product).Update("Price", 200)
	// Update - update multiple fields
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // non-zero fields
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})
	var dummies []DummyTable
	// call は上手くいかない…
	db.Raw("call dbo.dummyProcedure(?,?)", 0, 3).Scan(dummies)
	sql := `DECLARE	@return_value int
			EXEC	@return_value = [dbo].[dummyProcedure]
			@Param1 = ?,
			@Param2 = ?
			SELECT	'Return Value' = @return_value
			`
	db.Raw(sql, 0, 3).Scan(&dummies)
	fmt.Println("get dummies!")
	for i, e := range dummies {
		fmt.Println(i)
		fmt.Printf("%d / %s", e.Dum, e.Tex)
	}
	// Delete - delete product
	deleted := db.Delete(&product, 1)
	return deleted.Error
}
