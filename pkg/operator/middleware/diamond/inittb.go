package diamond

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	crdv1alpha1 "l0calh0st.cn/clickpaas-operator/pkg/apis/middleware/v1alpha1"
)

const (
	SQLCreateDatabase = `CREATE DATABASE if not exists diamond;`
	SQLCreateTableConfigInfo = "CREATE TABLE if not exists `config_info` (`id` bigint(64) unsigned NOT NULL AUTO_INCREMENT,`data_id` varchar(255) NOT NULL DEFAULT '',  `group_id` varchar(128) NOT NULL DEFAULT '',  `content` longtext NOT NULL,   `md5` varchar(32) NOT NULL DEFAULT '', `src_ip` varchar(20) DEFAULT NULL, `src_user` varchar(20) DEFAULT NULL,  `gmt_create` datetime NOT NULL,`gmt_modified` datetime NOT NULL,  `memo` varchar(1000) DEFAULT NULL, PRIMARY KEY (`id`) ) ENGINE=InnoDB AUTO_INCREMENT=550 DEFAULT CHARSET=utf8;"
	SQLCreateTableGroupInfo = "CREATE TABLE if not exists `group_info` (  `id` bigint(64) unsigned NOT NULL AUTO_INCREMENT, `address` varchar(70) NOT NULL DEFAULT '', `data_id` varchar(255) NOT NULL DEFAULT '', `group_id` varchar(128) NOT NULL DEFAULT '', `src_ip` varchar(20) DEFAULT NULL, `src_user` varchar(20) DEFAULT NULL,`gmt_create` datetime NOT NULL, `gmt_modified` datetime NOT NULL,    PRIMARY KEY (`id`),  UNIQUE KEY `uk_group_address` (`address`,`data_id`,`group_id`)  ) ENGINE=InnoDB DEFAULT CHARSET=utf8;"
	SQLCreateTableConfigHistory = "CREATE TABLE if not exists `config_info_history` ( `id` bigint(64) unsigned NOT NULL AUTO_INCREMENT, `data_id` varchar(255) NOT NULL DEFAULT '', `group_id` varchar(128) NOT NULL DEFAULT '', `content` longtext NOT NULL, `md5` varchar(32) NOT NULL DEFAULT '',`src_ip` varchar(20) DEFAULT NULL,  `src_user` varchar(20) DEFAULT NULL,`gmt_create` datetime NOT NULL, `gmt_modified` datetime NOT NULL, `memo` varchar(1000) DEFAULT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB AUTO_INCREMENT=550 DEFAULT CHARSET=utf8;"
)

type taskCreateDatabase struct {
	task *initDbTask
}

func NewTaskCreateDatabase(diamond *crdv1alpha1.Diamond)*taskCreateDatabase{
	return &taskCreateDatabase{task: &initDbTask{
		dbHost:     diamond.Spec.Config.Host,
		dbUser:     diamond.Spec.Config.User,
		dbPort:     diamond.Spec.Config.Port,
		dbPassword: diamond.Spec.Config.Password,
	}}
}

func(t *taskCreateDatabase)CreateDatabase()error{
	return t.task.createDatabase()
}

func(t *taskCreateDatabase)CreateTable()error{
	return t.task.createTable()
}

type initDbTask struct {
	dbHost string
	dbUser string
	dbPort int
	dbPassword string
}



func (t *initDbTask)createDatabase()error{
	db,err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/",t.dbUser, t.dbPassword, t.dbHost, t.dbPort))
	defer db.Close()
	if err != nil{
		return err
	}
	_,err = db.Exec(SQLCreateDatabase)
	if err != nil{
		return err
	}
	return nil
}

func(t *initDbTask)createTable()error{
	db,err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/diamond?charset=utf8",t.dbUser, t.dbPassword, t.dbHost, t.dbPort))
	defer db.Close()
	if err != nil{
		return err
	}
	_,err = db.Exec(SQLCreateTableConfigInfo)
	if err != nil{
		return fmt.Errorf("create table config_info failed; %s", err.Error())
	}
	_,err = db.Exec(SQLCreateTableGroupInfo)
	if err != nil {
		return fmt.Errorf("create table group_info failed: %s", err.Error())
	}
	_,err = db.Exec(SQLCreateTableConfigHistory)
	if err != nil{
		return fmt.Errorf("create config_history failed: %s", err.Error())
	}
	return nil
}