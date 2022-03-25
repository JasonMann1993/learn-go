package main

import (
	"database/sql/driver"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"strings"
)

var db *sqlx.DB

func initDB() (err error) {
	dsn := "homestead:secret@tcp(192.168.56.56)/mic?charset=utf8mb4&parseTime=True&loc=Local"
	// 也可以使用MustConnect连接不成功就panic
	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}
	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(10)
	return
}

type user struct {
	ID   int
	Name string
	Age  int
}

// 查询单条数据示例
func queryRowDemo() {
	sqlStr := "select id, name, age from user where id=?"
	var u user
	err := db.Get(&u, sqlStr, 1)
	if err != nil {
		fmt.Printf("get failed, err:%v\n", err)
		return
	}
	//fmt.Printf("id:%d name:%s age:%d\n", u.ID, u.Name, u.Age)
	fmt.Printf("%#v", u)
}

// 查询多条数据示例
func queryMultiRowDemo(){
	sqlStr := "select id, name, age from user where id > ?"
	var users []user
	err := db.Select(&users, sqlStr, 0)
	if err != nil {
		fmt.Printf("query failed, err:%v\n", err)
		return
	}
	fmt.Printf("users:%#v\n", users)
}

// 插入数据
func insertRowDemo() {
	sqlStr := "insert into user(name, age) values (?,?)"
	ret, err := db.Exec(sqlStr, "沙河小王子", 19)
	if err != nil {
		fmt.Printf("insert failed, err:%v\n", err)
		return
	}
	theID, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		fmt.Printf("get lastinsert ID failed, err:%v\n", err)
		return
	}
	fmt.Printf("insert success, the id is %d.\n", theID)
}

// 更新数据
func updateRowDemo() {
	sqlStr := "update user set age=? where id = ?"
	ret, err := db.Exec(sqlStr, 39, 6)
	if err != nil {
		fmt.Printf("update failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("update success, affected rows:%d\n", n)
}

// 删除数据
func deleteRowDemo() {
	sqlStr := "delete from user where id = ?"
	ret, err := db.Exec(sqlStr, 6)
	if err != nil {
		fmt.Printf("delete failed, err:%v\n", err)
		return
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		fmt.Printf("get RowsAffected failed, err:%v\n", err)
		return
	}
	fmt.Printf("delete success, affected rows:%d\n", n)
}

//DB.NamedExec方法用来绑定SQL语句与结构体或map中的同名字段。
func insertUserDemo()(err error){
	sqlStr := "insert into user (name,age) values (:name,:age)"
	_, err = db.NamedExec(sqlStr,
		map[string]interface{}{
			"name": "json",
			"age": "29",
		})
	return
}
// 与DB.NamedExec同理，这里是支持查询。
func namedQuery() {
	sqlStr := "select * from user where name=:name"
	rows, err := db.NamedQuery(sqlStr, map[string]interface{}{"name":"json"})
	if err != nil {
		fmt.Printf("db.NamedQuery failed, err:%v\n", err)
		return
	}
	defer rows.Close()
	for rows.Next(){
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}

	u := user{Name:"json"}
	rows,err = db.NamedQuery(sqlStr,u)
	if err != nil {
		fmt.Println("err:",err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var u user
		err := rows.StructScan(&u)
		if err != nil {
			fmt.Printf("scan failed, err:%v\n", err)
			continue
		}
		fmt.Printf("user:%#v\n", u)
	}
}

func transactionDemo2()(err error) {
	tx, err := db.Beginx() // 开启事物
	if err != nil {
		fmt.Printf("begin trans failed, err:%v\n", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			fmt.Println("rollback",err)
			tx.Rollback()
		} else {
			err = tx.Commit()
			fmt.Println("commit")
		}
	}()
	sqlStr1 := "Update user set age=20 where id=?"
	rs, err := tx.Exec(sqlStr1, 1)
	if err!= nil{
		return err
	}
	n, err := rs.RowsAffected()
	if err != nil {
		return err
	}

	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}

	sqlStr2 := "Update user set age=50 where id=?"
	rs, err = tx.Exec(sqlStr2, 2)
	if err!=nil{
		return err
	}
	n, err = rs.RowsAffected()
	if err != nil {
		return err
	}
	if n != 1 {
		return errors.New("exec sqlStr1 failed")
	}
	return
}

// BatchInsertUsers 自行构造批量插入的语句
func BatchInsertUsers(users []user) error {
	// 存放 （？，？） 的 slice
	valueStrings := make([]string, 0, len(users))
	// 存放value的slice
	valueArgs := make([]interface{}, 0, len(users) * 2)
	// 遍历users准备的数据
	for _, u := range users {
		// 此处占位符要与插入值的个数对应
		valueStrings = append(valueStrings, "(?, ?)")
		valueArgs = append(valueArgs, u.Name)
		valueArgs = append(valueArgs, u.Age)
	}
	// 自行拼接要执行的据体语句
	stmt := fmt.Sprintf("insert into user(name,age) values %s",strings.Join(valueStrings,","))
	_, err := db.Exec(stmt,valueArgs...)
	return err
}

// Value 使用 sqlx.In 实现批量插入
// 前提时需要实现 driver.Valuer 接口：
func (u user) Value() (driver.Value, error) {
	return []interface{}{u.Name, u.Age}, nil
}
// BatchInsertUsers2 使用sqlx.In帮我们拼接语句和参数, 注意传入的参数是[]interface{}
func BatchInsertUsers2(users []interface{}) error {
	query, args, _ := sqlx.In(
		"insert into user (name,age) values (?),(?),(?),(?)", // (?) 个数要和 len(users) 一样
		users..., // 如果 arg 实现了 driver.Valuer, sqlx,In 会通过调用 Value() 来展开它
		)
	fmt.Println(query) // 查看生成的querystring
	fmt.Println(args) // 场刊生成的args
	_, err := db.Exec(query, args...)
	return err
}

// BatchInsertUsers3 使用NamedExec实现批量插入
func BatchInsertUsers3(users []*user)error{
	_, err := db.NamedExec("insert into user (name,age) values(:name, :age)", users)
	return err
}

// QueryByIDs 根据给定ID查询
func QueryByIDs(ids []int)(users []user, err error){
	// 动态填充id
	query, args, err := sqlx.In("select name,age from user where id in (?)", ids)
	if err != nil {
		return
	}
	//fmt.Printf("%#v\n",query)
	query = db.Rebind(query)
	//fmt.Printf("%#v\n",query)


	err = db.Select(&users, query, args...)
	return
}

// QueryAndOrderByIDs 按照指定id查询并维护顺序
func QueryAndOrderByIds(ids []int)(users []user, err error){
	// 动态填充id
	strIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		strIDs = append(strIDs, fmt.Sprintf("%d", id))
	}
	query, args, err := sqlx.In("select name,age from user where id in (?) order by find_in_set(id,?)", ids, strings.Join(strIDs,","))
	if err != nil {
		return
	}

	err = db.Select(&users, query, args...)
	return
}


func main() {
	initDB()
	defer db.Close()
	//queryRowDemo()
	//queryMultiRowDemo()
	//insertRowDemo()
	//insertUserDemo()
	//namedQuery()
	//transactionDemo2()
	//BatchInsertUser2()


	// 批量插入
	//u1 := user{Name: "jason1", Age:18}
	//u2 := user{Name: "jason2", Age:19}
	//u3 := user{Name: "jason3", Age:28}
	// 方法1
	//users := []user{u1, u2, u3}
	//err := BatchInsertUsers(users)
	//if err != nil {
	//	fmt.Printf("BatchInsertUsers failed, err:%v\n", err)
	//}

	//users2 := []interface{}{u1, u2, u3, u1}
	//err := BatchInsertUsers2(users2)
	//if err != nil {
	//	fmt.Printf("BatchInsertUsers2 failed, err:%v\n", err)
	//}

	// 方法3
	//users3 := []*user{&u1, &u2, &u3}
	//err = BatchInsertUsers3(users3)
	//if err != nil {
	//	fmt.Printf("BatchInsertUsers3 failed, err:%v\n", err)
	//}

	//us,err := QueryByIDs([]int{18,19})
	us,err := QueryAndOrderByIds([]int{19,18})
	fmt.Println(us,err)
}
