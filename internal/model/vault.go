// Package model содержит модели доменной области
package model

import "time"

type VaultItem struct {
	ParentId        *int64          `json:"ParentId"`
	Id              int64           `json:"Id"`
	MasterId        int64           `json:"MasterId"`
	Title           string          `json:"Title"`
	Detail          string          `json:"Detail"`
	PartNumber      string          `json:"PartNumber"`
	Comm            string          `json:"Comm"`
	LastModDate     time.Time       `json:"LastModDate"`
	LastModUserId   int64           `json:"LastModUserId"`
	LastModUserName string          `json:"LastModUserName"`
	LfCycStateId    int             `json:"LfCycStateId"`
	NumSchmId       int             `json:"NumSchmId"`
	RevId           int64           `json:"RevId"`
	RevNum          string          `json:"RevNum"`
	VerNum          int             `json:"VerNum"`
	CadBOMStruct    string          `json:"CadBOMStruct"`
	CatName         string          `json:"CatName"`
	CatSystemName   string          `json:"CatSystemName"`
	CatId           int             `json:"CatId"`
	Quant           *float64        `json:"Quant"`
	PositionNum     *string         `json:"PositionNum"`
	UnitID          *int            `json:"UnitID"`
	Units           *string         `json:"Units"`
	Properties      []VaultProperty `json:"Properties"`
	Files           []VaultFile     `json:"Files"`
}

type VaultProperty struct {
	SysName  string      `json:"SysName"`
	DispName string      `json:"DispName"`
	Val      interface{} `json:"Val"`
}

type VaultFile struct {
	FileName    string    `json:"FileName"`
	Id          int64     `json:"Id"`
	MasterId    int64     `json:"MasterId"`
	VerNum      int       `json:"VerNum"`
	LastModDate time.Time `json:"LastModDate"`
	LinkType    string    `json:"LinkType"`
}
