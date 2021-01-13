package lts

import (
	"database/sql"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
	"fmt"
)

const (
	SQLCreateDatabaseLts = "CREATE DATABASE IF NOT EXISTS lts"

)


func taskCreateDatabaseOnce(lts *crdv1alpha1.LtsJobTracker)error{
	db,err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/",lts.Spec.Config.Db.User,
		lts.Spec.Config.Db.Password, lts.Spec.Config.Db.Host, lts.Spec.Config.Db.Port))
	if err != nil{
		return err
	}
	defer db.Close()
	_,err = db.Exec(SQLCreateDatabaseLts)
	if err != nil{
		return err
	}
	return nil
}