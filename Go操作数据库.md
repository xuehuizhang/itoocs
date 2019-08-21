### Go操作数据库 

#### MySql介绍

mysql常用引擎：

	MyIASM引擎
	
		1，不支持事务
	
		2，不支持行锁
	
		3，读性能比较好
	
	Innodb引擎
	
		1，支持事务
	
		2，支持行锁
	
		3，整体性能比较好

#### MySql开发

单行查询

```go
package main

import (
	"database/sql" //这里其实是一些接口
	"fmt"

	_ "github.com/go-sql-driver/mysql" //第三方包实现了上面的接口 执行其中的init函数就为了注册mysql驱动
)

//User 用户
type User struct {
	Id   int    `db:"Id"`
	Name string `db:"Name"`  //如果数据库为null，这里可以使用  Name sql.NullString定义
	Pwd  string `db:"Pwd"`
}

var (
	Db *sql.DB
)

//InitDb 初始化数据库
func InitDb() (err error) {
	dns := "root:root@tcp(39.105.55.98:3306)/test"
	Db, err = sql.Open("mysql", dns)
	if err != nil {
		return
	}

	err = Db.Ping()
	if err != nil {
		return
	}
	return nil
}

func main() {
	err := InitDb()
	if err != nil {
		fmt.Println("init db err:", err)
	}
	sqlStr := "select id,name ,pwd from UserInfo where id=?"
	row := Db.QueryRow(sqlStr, 1)
	var user User

	err = row.Scan(&user.Id, &user.Name, &user.Pwd)  //获取到row之后如果不scan，就会一致堵塞
	if err != nil {                                 //因此一定要进行scan操作
		fmt.Println("row scan err:", err)
	}
	fmt.Println(user)
}
```

多行查询

```go
func queryMutile() {
	sqlStr := "select id,name,pwd from UserInfo where id>?"
	row, err := Db.Query(sqlStr, 0)
	if err != nil {
		fmt.Println("db query error", err)
	}
	defer row.Close() //这里一定要记得释放数据库连接，不然会占用连接

	for row.Next() {
		var user User
		err := row.Scan(&user.Id, &user.Name, &user.Pwd)
		if err != nil {
			fmt.Println("row scan error:", err)
		}
		fmt.Println(user)
	}
}
```

这里不知道为啥，在查询mysql时，表名必须和数据库大小写一致

案例：连接池耗尽

```go
var (
	Db *sql.DB
)

//InitDb 初始化数据库
func InitDb() (err error) {
	dns := "root:root@tcp(39.105.55.98:3306)/test"
	Db, err = sql.Open("mysql", dns)
	if err != nil {
		return
	}

	err = Db.Ping()
	if err != nil {
		return
	}
	Db.SetMaxOpenConns(100)
	Db.SetMaxIdleConns(16)
	return nil
}

func querySigle() {
	for i := 0; i < 101; i++ {
		fmt.Printf("query %d times\n", i)
		sqlStr := "select id,name ,pwd from UserInfo where id=?"
		row := Db.QueryRow(sqlStr, 2)

		if row != nil {
			continue   //这里获取到row之后没有scan，就会一致占用连接
		}              //因此在查询多行的时候也要记得关闭rows

		var user User
		err := row.Scan(&user.Id, &user.Name, &user.Pwd)
		if err != nil {
			fmt.Println("row scan err:", err)
		}
		fmt.Println(user)
	}
}
```



#### 增删改

```go
func insertData() {
	sqlStr := "insert into UserInfo(name,pwd) values (?,?)"
	result, err := Db.Exec(sqlStr, "张三", "123")
	if err != nil {
		fmt.Println("insert error:", err)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("get last insert id error:", err)
		return
	}
	fmt.Println("last insert id is ", id)
}

func updateData() {
	sqlStr := "update UserInfo set name=? where id=?"
	result, err := Db.Exec(sqlStr, "李四", 3)
	if err != nil {
		fmt.Println("update error:", err)
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		fmt.Println("get rows affected error", err)
		return
	}
	fmt.Println("affected N:", n)
}

func deleteData() {
	sqlStr := "delete from UserInfo where id=?"
	result, err := Db.Exec(sqlStr, 3)
	if err != nil {
		fmt.Println("delete error:", err)
		return
	}
	n, err := result.RowsAffected()
	if err != nil {
		fmt.Println("get rows affected error", err)
		return
	}
	fmt.Println("affected N:", n)
}
```

#### 预处理

```go
func preparSigleQuery() {
	sqlStr := "select * from UserInfo where id=?"
	stmt, err := Db.Prepare(sqlStr)
	if err != nil {
		fmt.Println("db prepare error:", err)
		return
	}

	row := stmt.QueryRow(1)
	var user User
	err = row.Scan(&user.Id, &user.Name, &user.Pwd)
	if err != nil {
		fmt.Println("row scan error:", err)
		return
	}
	fmt.Println(user)
}

func preparMutileQuery() {
	sqlStr := "select * from UserInfo where id>?"
	stmt, err := Db.Prepare(sqlStr)
	if err != nil {
		fmt.Println("db prepare error:", err)
		return
	}

	row, err := stmt.Query(0)
	if err != nil {
		fmt.Println("stmt query error:", err)
		return
	}
	defer row.Close()
	for row.Next() {
		var user User
		err = row.Scan(&user.Id, &user.Name, &user.Pwd)
		if err != nil {
			fmt.Println("row scan error:", err)
			return
		}
		fmt.Println(user)
	}
}

func preparInsert() {
	sqlStr := "insert into UserInfo (name,pwd) values(?,?)"
	stmt, err := Db.Prepare(sqlStr)
	if err != nil {
		fmt.Println("db prepare error:", err)
		return
	}

	result, err := stmt.Exec("李四", "321")
	if err != nil {
		fmt.Println("stmt exec error", err)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println("get last insert id error:", err)
		return
	}
	fmt.Println("last insert id ", id)
}
```

#### 事务

```go
func testTrans() {
	trans, err := Db.Begin()
	if err != nil {
		fmt.Println("db begin trans error:", err)
		return
	}

	sqlStr01 := "update UserInfo set pwd=? where id=?"
	_, err = trans.Exec(sqlStr01, "222", 1)
	if err != nil {
		fmt.Println("trans exec error:", err)
		trans.Rollback()
		return
	}
	sqlStr02 := "update UserInfo set pwd=? where id=?"
	_, err = trans.Exec(sqlStr02, "333", 2)
	if err != nil {
		fmt.Println("trans exec error:", err)
		trans.Rollback()
		return
	}
	err = trans.Commit()
	if err != nil {
		fmt.Println("trans commit error:", err)
		trans.Rollback()
		return
	}
}

```

特点：原子性  一致性  隔离性  持久性

#### sqlx

```go
package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var (
	Db *sqlx.DB
)

func initDb() (err error) {
	var dns = "root:root@tcp(39.105.55.98:3306)/test"
	Db, err = sqlx.Open("mysql", dns)
	if err != nil {
		return
	}

	err = Db.Ping()
	if err != nil {
		return
	}

	Db.SetMaxIdleConns(16)
	Db.SetMaxOpenConns(100)
	return nil
}

type user struct {
	Id   int            `db:"Id"`
	Name sql.NullString `db:"Name"`
	Pwd  string         `db:"Pwd"`
}

func sigleQuery() {
	sqlStr := "select * from UserInfo where id=?"
	var user user
	err := Db.Get(&user, sqlStr, 1)
	if err != nil {
		fmt.Println("db get error:", err)
	}
	fmt.Println(user)
}

func mutileQuery() {
	sqlStr := "select * from UserInfo where id>?"
	var users []user
	err := Db.Select(&users, sqlStr, 0)
	if err != nil {
		fmt.Println("db select error:", err)
	}
	fmt.Println(users)
}

func insertData() {
	sqlStr := "insert into UserInfo (name,pwd)values(?,?)"
	var user = user{
		Name: sql.NullString{
			Valid:  true,
			String: "马武",
		},
		Pwd: "123",
	}
	r, err := Db.DB.Exec(sqlStr, user.Name, user.Pwd)
	if err != nil {
		fmt.Println("db exec error:", err)
	}
	id, err := r.LastInsertId()
	if err != nil {
		fmt.Println("get last insert id error:", err)
	}
	fmt.Println(id)
}

func main() {
	err := initDb()
	if err != nil {
		fmt.Println("init db error:", err)
	}

	//sigleQuery()
	//mutileQuery()
	insertData()
}

```

sqlx事务类似于database/sql中，sqlx其实就是对database/sql的封装

为了防止sql注入漏洞，不要去拼接sql