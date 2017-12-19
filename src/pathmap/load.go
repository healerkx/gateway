
package pathmap

import (
	"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func LoadApiBindingInfo() bool {
	db, err := sql.Open("mysql", "root:root@/gateway?charset=utf8")
	if err != nil {
		println(err)
		return false
	}
	
	stmt, err := db.Prepare(`SELECT * from gw_api_bind where status=1;`)
	if err != nil {
		println(err.Error())
		return false
	}
	rows, err := stmt.Query()
	if err != nil {
		println(err)
	}

	cols, _ := rows.Columns();
    //这里表示一行所有列的值，用[]byte表示
    vals := make([][]byte, len(cols));
    //这里表示一行填充数据
    scans := make([]interface{}, len(cols));
    //这里scans引用vals，把数据填充到[]byte里
    for k, _ := range vals {
        scans[k] = &vals[k];
    }
 
    i := 0;
	result := []map[string]string{};
	
	for rows.Next() {
        //填充数据
        rows.Scan(scans...);
        //每行数据
        row := make(map[string]string);
        //把vals中的数据复制到row中
        for k, v := range vals {
            key := cols[k];
            //这里把[]byte数据转成string
            row[key] = string(v);
        }
        //放入结果集
        result = append(result, row)
        i++;
	}
	fmt.Printf("%v\n", result)
	return true
}