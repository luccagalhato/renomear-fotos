package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strconv"

	_ "github.com/denisenkom/go-mssqldb"
)

type sqlStruct struct {
	url *url.URL
	db  *sql.DB
}

func (s *sqlStruct) init() error {
	s.url = &url.URL{
		Scheme: "sqlserver",
		User:   url.UserPassword("felipe.restum", "z12xcvab@#$"),
		// Host:     "10.58.192.51:1433",
		Host:     "179.183.30.186:3215",
		RawQuery: url.Values{}.Encode(),
	}
	var err error
	s.db, err = sql.Open("sqlserver", s.url.String())
	if err != nil {
		return err
	}

	ctx := context.Background()

	return s.db.PingContext(ctx)
}

func (s *sqlStruct) Close() error {
	return s.db.Close()
}

// GetGtinCode ...
func (s *sqlStruct) GetGtinCode(code string, cor string) (filenames []string, err error) {
	corInt, err := strconv.Atoi(cor)
	if err != nil {
		return
	}
	rows, err := s.db.QueryContext(context.Background(), fmt.Sprintf(`SELECT CODIGO_BARRA FROM LINX_TBFG..PRODUTOS_BARRA
	WHERE TIPO_COD_BAR = '1' AND PRODUTO = '%s' AND COR_PRODUTO = %d`, code, corInt))
	if err != nil {
		return
	}
	filenames = make([]string, 0)
	for rows.Next() {
		var code string
		if err = rows.Scan(&code); err != nil {
			continue
		}
		filenames = append(filenames, code)
	}

	if len(filenames) == 0 {
		err = fmt.Errorf("não encontrado")
	}

	return
}

var sqlVar *sqlStruct
