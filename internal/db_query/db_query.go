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

	rows, err := conn.Query(context.Background(), "SELECT player, record FROM records ORDER BY record DESC LIMIT 10")
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

func (db *DB) AddPlayer(data []byte) int {
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
	log.Default().Printf("adding: %s\n", fmt.Sprintf("VALUES ('%s', 0)", stru.SanitizedPlayer()))
	rows, err := conn.Query(context.Background(), fmt.Sprintf("INSERT INTO records (player, record) VALUES ('%s', 0) RETURNING id", stru.SanitizedPlayer()))
	if err != nil {
		log.Default().Printf("QueryRow failed: %v\n", err)
	}
	defer rows.Close()

	var newId int
	for rows.Next() {
		err = rows.Scan(&newId)
		if err != nil {
			log.Default().Println("failed to retrieve new player id")
		}
	}
	return newId
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

	log.Default().Printf("QueryRow failed: %s\n", fmt.Sprintf("record = %d WHERE id = %d", stru.Val, stru.Player))
	rows, err := conn.Query(context.Background(), fmt.Sprintf("UPDATE records SET record = %d WHERE id = %d", stru.Val, stru.Player))
	if err != nil {
		log.Default().Printf("QueryRow failed: %v\n", err)
	}
	defer rows.Close()
}

func (db *DB) GetPlayerRecord(id int) string {
	conn, err := pgx.ConnectConfig(context.Background(), db.connCfg)
	if err != nil {
		log.Default().Printf("QueryRow failed: %v\n", err)
	}
	defer conn.Close(context.Background())

	log.Default().Printf("SELECT player, record FROM records WHERE id = %d", id)
	rows, err := conn.Query(context.Background(), fmt.Sprintf("SELECT player, record FROM records WHERE id = %d", id))
	if err != nil {
		log.Default().Printf("QueryRow failed: %v\n", err)
	}
	defer rows.Close()
	ans := dbrequests.Record{}
	for rows.Next() {
		rows.Scan(&ans.Player, &ans.Val)
		fmt.Printf("%+v", ans)
	}

	b, err := json.Marshal(ans)
	if err != nil {
		log.Default().Printf("QueryRow failed: %v\n", err)
	}
	log.Default().Printf("json: %s", string(b))
	return string(b)
}
