package dbquery

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"hleb_flip/internal/config"
	dbrequests "hleb_flip/internal/db_requests"
	"os"

	"github.com/jackc/pgx/v5"
)

func Wtf() string {
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

	conn, err := pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())

	rows, err := conn.Query(context.Background(), "SELECT * FROM records")
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
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
