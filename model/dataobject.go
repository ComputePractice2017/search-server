package model

import (
	"strings"

	r "gopkg.in/gorethink/gorethink.v3"
)

const (
	searchDbName       string = "search"
	documentsTableName string = "docs"
)

type Document struct {
	ID   string `json:"id"`
	Url  string `json:"url"`
	Text string `json:"text"`
}

var session *r.Session

func InitSession() error {
	var err error
	session, err = r.Connect(r.ConnectOpts{
		Address: "localhost",
	})
	if err != nil {
		return err
	}

	err = createDbIfNotExists()
	if err != nil {
		return err
	}

	err = createTableIfNotExists()
	if err != nil {
		return err
	}

	return nil
}

func createDbIfNotExists() error {
	_, err := r.DBCreate(searchDbName).Run(session)
	if err != nil && strings.Contains(err.Error(), "exists") {
		return nil
	}
	return err
}

func createTableIfNotExists() error {
	_, err := r.DB(searchDbName).TableCreate(documentsTableName).Run(session)
	if err != nil && strings.Contains(err.Error(), "exists") {
		return nil
	}
	return err
}

func FindDocs(query string) ([]Document, error) {
	res, err := r.DB(searchDbName).Table(documentsTableName).Filter(func(row r.Term) r.Term {
		return row.Field("text").Match("(?i)" + query)
	}).Run(session)
	if err != nil {
		return nil, err
	}

	var response []Document
	err = res.All(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
