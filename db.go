package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	//_ "odbc/driver"
	//_ "github.com/alexbrainman/odbc"
	_ "github.com/go-sql-driver/mysql"
)

var dbConnection map[string]*sql.DB

// DisconnectDB close down db connection and clear memory
func DisconnectDB() {

	for key, _ := range dbConnection {
		dbConnection[key].Close()
	}
	//dbConnection["default"].Close()
}

// ConnectDB make a connection to database using provide DSN as envoirent variable ConnectDB()
func ConnectDB() bool {
	var err error
	dbConnection = make(map[string]*sql.DB)
	var value *sql.DB

	dbType := os.Getenv("DBTYPE")
	if dbType == "" {
		Panic("can't get DBTYPE value from envoirenment variable")
		return false
	}
	switch dbType {
	case "ODBC":
		DSN := os.Getenv("DBDSN")
		if DSN == "" {
			Panic("can't get DBDSN value from envoirenment variable")
			return false
		}
		value, err = sql.Open("odbc", "DSN="+DSN)
		dbConnection["default"] = value
		break
	case "MYSQL":
		DBNAME := os.Getenv("DBNAME")
		if DBNAME == "" {
			Panic("can't get DBNAME value from envoirenment variable")
			return false
		}
		DBUSER := os.Getenv("DBUSER")
		if DBUSER == "" {
			Panic("can't get DBUSER value from envoirenment variable")
			return false
		}
		DBPASS := os.Getenv("DBPASS")
		if DBPASS == "" {
			Panic("can't get DBPASS value from envoirenment variable")
			return false
		}
		strConnection := ""
		DBHOST := os.Getenv("DBHOST")
		if DBHOST == "" {
			strConnection = fmt.Sprintf("%s:%s@/%s", DBUSER, DBPASS, DBNAME)
		} else {
			strConnection = fmt.Sprintf("%s:%s@%s/%s", DBUSER, DBPASS, DBHOST, DBNAME)
		}
		value, err = sql.Open("mysql", strConnection)
		dbConnection["default"] = value

		break
	default:
		Panic("Invalid DBTYPE, currently ODBC / MYSQL supported")
		break
	}
	if err != nil {
		Panic("DB Connection Error " + err.Error())
		return false
	}

	params := make([]interface{}, 0)
	dbConnections, _ := GetAllRows("SELECT * FROM `db_connections` WHERE `status`='ONLINE'", params, "default")
	for k := 0; k < len(dbConnections); k++ {

		strConnection := ""
		if dbConnections[k]["DBHOST"] == "" {
			strConnection = fmt.Sprintf("%s:%s@/%s", dbConnections[k]["DBUSER"], dbConnections[k]["DBPASS"], dbConnections[k]["DBNAME"])
		} else {
			strConnection = fmt.Sprintf("%s:%s@%s/%s", dbConnections[k]["DBUSER"], dbConnections[k]["DBPASS"], dbConnections[k]["DBHOST"], dbConnections[k]["DBNAME"])
		}

		value, err2 := sql.Open("mysql", strConnection)
		if err2 != nil {
			Panic("DB Connection Error " + err.Error())
			return false
		}
		dbConnection[dbConnections[k]["name"]] = value

		params := make([]interface{}, 0)
		Row, ok := GetSingleRow("SELECT NOW() AS 'Time'", params, "default")
		if ok {
			Info("DB Time " + Row["Time"])
		}

		Log(dbConnections[k]["name"] + " DB Connected")

		MaxOpenConns, _ := strconv.Atoi(dbConnections[k]["maxOpenConns"])
		MaxIdleConns, _ := strconv.Atoi(dbConnections[k]["maxIdleConns"])

		dbConnection[dbConnections[k]["name"]].SetMaxOpenConns(MaxOpenConns)
		dbConnection[dbConnections[k]["name"]].SetMaxIdleConns(MaxIdleConns)
		dbConnection[dbConnections[k]["name"]].SetConnMaxLifetime(time.Hour)

	}

	dbConnection["default"].SetMaxOpenConns(100)
	dbConnection["default"].SetMaxIdleConns(90)
	dbConnection["default"].SetConnMaxLifetime(time.Hour)
	keepAliveDB()
	return true
}

func keepAliveDB() {
	timer := time.NewTimer(time.Second * 60)
	go func() {
		<-timer.C
		Log("Refresh DB Connection")
		params := make([]interface{}, 0)
		Row, ok := GetSingleRow("SELECT NOW() AS 'Time'", params, "default")
		if ok {
			Info("DB Time " + Row["Time"])
		}
		keepAliveDB()
	}()
}

// GetSingleRow (Query as string, QueyParameters as []interface{}) return (Row as map[string]string, staus as bool)
func GetSingleRow(Query string, params []interface{}, DBConnectionName string) (map[string]string, bool) {
	Debug("GetSingleRow : Executing Query :: " + Query)
	var rows *sql.Rows
	var err error
	if len(params) > 0 {
		paramPointers := make([]interface{}, len(params))
		for i := range params {
			paramPointers[i] = &params[i]
		}
		rows, err = dbConnection[DBConnectionName].Query(Query, paramPointers...)
	} else {
		rows, err = dbConnection[DBConnectionName].Query(Query)
	}
	if err != nil {
		Warning("GetSingleRow : Error:" + err.Error() + "\tExecuting Query: " + Query)
		return nil, false
	}
	defer rows.Close()
	columnNames, err := rows.Columns()
	if err != nil {
		Warning("GetSingleRow : Error:" + err.Error() + "\tFetching Column Names Query: " + Query)
		return nil, false
	}
	if len(columnNames) < 1 {
		Info("GetSingleRow : Error:" + err.Error() + "\tNo Column in Resultset Query: " + Query)
		return nil, false
	}
	columns := make([]interface{}, len(columnNames))
	for i := range columnNames {
		columns[i] = new(sql.RawBytes)
	}
	Debug("GetSingleRow :\tscaning rows for Query: " + Query)
	ThisRow := make(map[string]string)
	ReturnStatus := false
	for rows.Next() {
		err := rows.Scan(columns...)
		if err != nil {
			Warning("GetSingleRow : Error:" + err.Error() + "\tScannig rows for Query: " + Query)
			return nil, false
		}
		for i, colName := range columnNames {
			if rb, ok := columns[i].(*sql.RawBytes); ok {
				colValue := string(*rb)
				ThisRow[colName] = colValue
				*rb = nil // reset pointer to discard current value to avoid a bug
				Debug("Column : " + colName + "\tValue : " + colValue)
			} else {
				Warning("GetSingleRow : Column " + colName + " contains nil value for Query: " + Query)
				ThisRow[colName] = ""
			}
		}
		ReturnStatus = true
		return ThisRow, true
	}
	return ThisRow, ReturnStatus
}

// GetAllRows (Query as string, QueyParameters as []interface{}) return (Row as []map[string]string, staus as bool)
func GetAllRows(Query string, params []interface{}, DBConnectionName string) ([]map[string]string, bool) {
	Debug("GetAllRows : Executing Query :: " + Query)
	var rows *sql.Rows
	var err error
	if len(params) > 0 {
		paramPointers := make([]interface{}, len(params))
		for i := range params {
			paramPointers[i] = &params[i]
		}
		rows, err = dbConnection[DBConnectionName].Query(Query, paramPointers...)
	} else {
		rows, err = dbConnection[DBConnectionName].Query(Query)
	}
	if err != nil {
		Warning("GetAllRows : Error:" + err.Error() + "\tExecuting Query: " + Query)
		return nil, false
	}
	defer rows.Close()
	columnNames, err := rows.Columns()
	if err != nil {
		Warning("GetAllRows : Error:" + err.Error() + "\tFetching Column Names Query: " + Query)
		return nil, false
	}
	if len(columnNames) < 1 {
		Info("GetAllRows : Error:" + err.Error() + "\tNo Column in Resultset Query: " + Query)
		return nil, false
	}
	columns := make([]interface{}, len(columnNames))
	for i := range columnNames {
		columns[i] = new(sql.RawBytes)
	}
	Debug("GetAllRows :\tscaning rows for Query: " + Query)
	var RecordsSet []map[string]string

	RowCount := 0
	for rows.Next() {
		RowCount++
		ThisRow := make(map[string]string)
		err := rows.Scan(columns...)
		if err != nil {
			Warning("GetSingleRow : Error:" + err.Error() + "\tScannig rows for Query: " + Query)
			return nil, false
		}
		for i, colName := range columnNames {
			if rb, ok := columns[i].(*sql.RawBytes); ok {
				colValue := string(*rb)
				ThisRow[colName] = colValue
				*rb = nil // reset pointer to discard current value to avoid a bug
				Debug("Column : " + colName + "\tValue : " + colValue)
			} else {
				Warning("GetSingleRow : Column " + colName + " contains nil value for Query: " + Query)
				ThisRow[colName] = ""
			}
		}
		RecordsSet = append(RecordsSet, ThisRow)
	}
	if RowCount > 0 {
		return RecordsSet, true
	}
	return nil, false
}

// UpdateDB (Query string, params []interface{}) and return (LastInsertID for INSERT OR Updated rows as int, status as bool)
func UpdateDB(Query string, params []interface{}, DBConnectionName string) (int64, bool) {
	Debug("UpdateDB : Executing Query :: " + Query)
	var Res sql.Result
	var err error
	if len(params) > 0 {
		paramPointers := make([]interface{}, len(params))
		for i := range params {
			paramPointers[i] = &params[i]
		}
		Res, err = dbConnection[DBConnectionName].Exec(Query, paramPointers...)
	} else {
		Res, err = dbConnection[DBConnectionName].Exec(Query)
	}
	if err != nil {
		Debug("UpdateDB : Error: " + err.Error() + " Executing Query: " + Query)
		return 0, false
	}
	var Num int64
	Num = 0
	switch strings.ToUpper(Query[0:6]) {
	case "UPDATE":
		count, err := Res.RowsAffected()
		if err == nil {
			Num = count
			Debug(fmt.Sprintf("UpdateDB : %d Rows(s) Updated for Query: %s", Num, Query))
		} else {
			Warning(fmt.Sprintf("UpdateDB: can't get updated rows, Error: %s", err.Error()))
		}
		break
	case "INSERT":
		count, err := Res.LastInsertId()
		if err == nil {
			Num = count
			Debug(fmt.Sprintf("UpdateDB : Last Insert ID is %d for Query: %s", Num, Query))
		} else {
			Warning(fmt.Sprintf("UpdateDB: can't get LastInsertID, Error: %s", err.Error()))
		}
		break
	default:
		break
	}
	return Num, true
}

// IsSQLSafe (strValue string) check if the string is safe to put in SQL
func IsSQLSafe(strValue string) bool {
	if (strings.Index(strValue, "'") > -1) || (strings.Index(strValue, ";") > -1) || (strings.Index(strValue, "`") > -1) {
		return false
	}
	return true
}

// AdvancedSQL (RequestJSON map[string]interface{}, ResJSON map[string]interface{}, ConditionState string) pass thourgh RequestJSON
// and build advanced where cluase and limit clause and return (bool status, string SQL)
func AdvancedSQL(RequestJSON map[string]interface{}, ResJSON map[string]interface{}, ConditionState string) (string, bool) {
	ReturnSTR := ""

	/*
		RequestData, ok := RequestJSON["data"]
		if !ok {
			Debug("AdvancedSQL:: No data found in RequestJSIN")
			return ReturnSTR, true
		}
		if reflect.ValueOf(RequestData).Kind() != reflect.Map {
			ResJSON["status"] = 486
			ResJSON["message"] = "inavlid data should be a valid JSON Object"
			return ReturnSTR, false
		}
		DataMap := RequestData.(map[string]interface{})
	*/
	DataMap := RequestJSON

	ok, searchField := GetKeyFromJSON(DataMap, ResJSON, "search_field")
	if ok {

		if strings.Index(searchField, ".") > -1 {
			searchField = strings.Replace(searchField, ".", "`.`", -1)
		}

		ok, searchValue := GetKeyFromJSON(DataMap, ResJSON, "search_value")
		if !ok {
			return "", false
		}
		ok, condition := GetKeyFromJSON(DataMap, ResJSON, "condition")
		if !ok {
			return "", false
		}
		switch condition {
		case "start_with":
			ReturnSTR += fmt.Sprintf(" %s `%s` LIKE '%s'", ConditionState, searchField, "%"+searchValue)
			break
		case "end_with":
			ReturnSTR += fmt.Sprintf(" %s `%s` LIKE '%s'", ConditionState, searchField, searchValue+"%")
			break
		case "contains":
			ReturnSTR += fmt.Sprintf(" %s `%s` LIKE '%s'", ConditionState, searchField, "%"+searchValue+"%")
			break
		case "equal":
			ReturnSTR += fmt.Sprintf(" %s `%s` = '%s'", ConditionState, searchField, searchValue)
			break
		default:
			ResJSON["status"] = 406
			ResJSON["message"] = "Invalid Condition " + condition
			return "", false
		}
	}

	ok, OrderClause := GetKeyFromJSON(DataMap, ResJSON, "sorting")
	if ok {
		ReturnSTR += fmt.Sprintf(" %s ", OrderClause)
	}

	ok, rowCount := GetKeyFromJSON(DataMap, ResJSON, "row_count")
	if ok {
		if !isNumeric(rowCount) {
			ResJSON["status"] = 406
			ResJSON["message"] = "row_count must be numeric"
			return "", false
		}
		ok, startRow := GetKeyFromJSON(DataMap, ResJSON, "start_row")
		if ok {
			if !isNumeric(startRow) {
				ResJSON["status"] = 406
				ResJSON["message"] = "start_row must be numeric"
				return "", false
			}
			ReturnSTR += fmt.Sprintf(" LIMIT %s, %s", startRow, rowCount)
		} else {
			ReturnSTR += fmt.Sprintf(" LIMIT %s", rowCount)
		}
	}

	return ReturnSTR, true
}
