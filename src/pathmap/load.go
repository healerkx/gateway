
package pathmap

import (
	//"fmt"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

func readRows(rows *sql.Rows) []map[string]string {
	cols, _ := rows.Columns();
    //这里表示一行所有列的值，用[]byte表示
    vals := make([][]byte, len(cols));
    //这里表示一行填充数据
    scans := make([]interface{}, len(cols));
    //这里scans引用vals，把数据填充到[]byte里
    for k, _ := range vals {
        scans[k] = &vals[k];
    }
 
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
            row[key] = string(v)
        }
        result = append(result, row)
	}
	return result
}

/**
 * updateTime == 0 for loading all
 */
func LoadApiBindingInfo(updateTime int32) ([]map[string]string, error) {
	db, err := sql.Open("mysql", "root:root@/gateway?charset=utf8")
	if err != nil {
		println(err)
		return nil, err
	}
	
	stmt, err := db.Prepare("SELECT * from gw_api_bind where status=1 and update_time>= ? ;")
	if err != nil {
		println(err.Error())
		return nil, err
	}
	rows, err := stmt.Query(updateTime)
	if err != nil {
		println(err)
	}
	return readRows(rows), nil
}


