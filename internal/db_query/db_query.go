package dbquery

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"hleb_flip/internal/config"
	dbrequests "hleb_flip/internal/db_requests"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
)

type DB struct {
	dbCfg    config.DBConfig
	certPool *x509.CertPool
	connCfg  *pgx.ConnConfig
}

func NewDB() *DB {
	cfg := config.GetDBConfig()
	rootCertPool := x509.NewCertPool()
	pem, err := os.ReadFile(cfg.CrtPath)
	if err != nil {
		panic(err)
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		panic("Failed to append PEM.")
	}

	connstring := fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=verify-full target_session_attrs=read-write",
		cfg.Host, cfg.Port, cfg.DBName, cfg.User, cfg.Pswd)

	connConfig, err := pgx.ParseConfig(connstring)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	connConfig.TLSConfig = &tls.Config{
		RootCAs:            rootCertPool,
		InsecureSkipVerify: true,
	}

	return &DB{
		dbCfg:   *cfg,
		connCfg: connConfig,
	}
}

func (db *DB) GetRecords() string {
	conn, err := pgx.ConnectConfig(context.Background(), db.connCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM records ORDER BY record DESC LIMIT 10")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
	defer rows.Close()

	sss := dbrequests.RecordList{List: make([]dbrequests.Record, 0)}
	for rows.Next() {
		ans := dbrequests.Record{}
		rows.Scan(&ans.Player, &ans.Val)
		sss.List = append(sss.List, ans)
		fmt.Printf("%+v", ans)
	}

	b, err := json.Marshal(sss)
	if err != nil {
		fmt.Fprintf(os.Stderr, "json failed: %v\n", err)
	}

	return string(b)
}

func (db *DB) AddPlayer(data []byte) {
	conn, err := pgx.ConnectConfig(context.Background(), db.connCfg)
	if err != nil {
		log.Default().Printf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	stru := dbrequests.AddRequest{}
	err = json.Unmarshal(data, &stru)
	if err != nil {
		log.Default().Printf("Unable to unmarshall json: %v\n", err)
	}

	log.Default().Printf("adding: %v\n", stru)
	rows, err := conn.Query(context.Background(), fmt.Sprintf("INSERT INTO records (player, record) VALUES ('%s', 0)", stru.SanitizedPlayer()))
	if err != nil {
		log.Default().Printf("QueryRow failed: %v\n", err)
	}
	defer rows.Close()
}

func (db *DB) ChangeRecordForPlayer(data []byte) {
	conn, err := pgx.ConnectConfig(context.Background(), db.connCfg)
	if err != nil {
		log.Default().Printf("Unable to connect to database: %v\n", err)
	}
	defer conn.Close(context.Background())

	stru := dbrequests.ChangePlayerRequest{}
	err = json.Unmarshal(data, &stru)
	if err != nil {
		log.Default().Printf("Unable to unmarshall json: %v\n", err)
	}

	rows, err := conn.Query(context.Background(), fmt.Sprintf("UPDATE records SET record = %d WHERE player = '%s'", stru.Val, stru.SanitizedPlayer()))
	if err != nil {
		log.Default().Printf("QueryRow failed: %v\n", err)
	}
	defer rows.Close()
}
