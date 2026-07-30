package config

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/Thanga-tamil/noway_service/internal/logger"
	"github.com/jmoiron/sqlx"
)

var MasterDB *sqlx.DB
var TenantDBs map[string]*sqlx.DB


type TenantSource struct {
	ApiDomain string `json:"api_domain"`
	DataSource string `json:"datasource"`
}

func InitSql(c Cfg) {

	 db, err := sqlx.Open(c.Postgre.Type, get_db_uri(c))

	 if err != nil {
		logger.Error("error while opening tcp connection with postgreSQL: ", err.Error())
		 panic(err)
	 }

	 lookupDB := "SET SEARCH_PATH TO datasource;"

	 if _, err := db.Exec(lookupDB); err != nil {
		 logger.Error("error when executing search path for tenant config: ", err.Error())
		 panic(err)
	 } else {
		 MasterDB = db
	 }

	 row, err := db.Query("select api_domain, datasource from connection_config;")

	 if err != nil {
		logger.Error("error while executing Query func:: ", err.Error())
	 }

	 var tenantDBs []TenantSource

	 for row.Next() {
		 var tenantDB TenantSource
		 if err := row.Scan(&tenantDB.ApiDomain, &tenantDB.DataSource); err != nil {
			logger.Error("error while Scanning executed query for tenant search path:: ", err.Error())
		 }
		 tenantDBs = append(tenantDBs, tenantDB)
	 }

		logger.Infof("Tenant datasource from master config: %#v\n", tenantDBs)
	 loadTenantDBs(c, tenantDBs)
}

func get_db_uri(c Cfg) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	c.Postgre.Host, c.Postgre.Port, c.Postgre.User, c.Postgre.Password, c.Postgre.Dbns)
}

func loadTenantDBs(c Cfg, tenantDBs []TenantSource){

	dbs := make(map[string]*sqlx.DB)

	for _, tenant := range tenantDBs {

		db, err := sqlx.Open(c.Postgre.Type, get_db_uri(c))

		if err != nil {
			logger.Error("error while opening tcp connection with postgreSQL: ", err.Error())
			panic(err)
		}

		lookupDB := "SET SEARCH_PATH TO " + tenant.DataSource + " ;"

		if _, err := db.Exec(lookupDB); err != nil {
			logger.Error("error when executing search path for tenant config: ", err.Error())
			panic(err)
		} else {
			dbs[tenant.ApiDomain] = db
		}
		logger.Info("loading tenant db for tenant: ", tenant)
	}

	TenantDBs = dbs

	fmt.Println("Tenant dbs loaded in memory:: ", TenantDBs)
}
