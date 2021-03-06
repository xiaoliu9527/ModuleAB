package models

import (
	"fmt"

	"github.com/ModuleAB/ModuleAB/server/common"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/validation"
	"github.com/pborman/uuid"
)

type Oas struct {
	Id         string        `orm:"pk;size(36)" json:"id" valid:"Match(/^[A-Fa-f0-9]{8}-([A-Fa-f0-9]{4}-){3}[A-Fa-f0-9]{12}$/)"`
	Endpoint   string        `json:"endpoint" valid:"Required"`
	VaultName  string        `orm:"size(32) json:"vaultName" valid:"Required"`
	VaultId    string        `orm:"size(32) json:"vaultId" valid:"Required"`
	BackupSets []*BackupSets `orm:"reverse(many)"`
	Jobs       []*OasJobs    `orm:"reverse(many)"`
}

func init() {
	if prefix := beego.AppConfig.String("database::mysqlprefex"); prefix != "" {
		orm.RegisterModelWithPrefix(prefix, new(Oas))
	} else {
		orm.RegisterModel(new(Oas))
	}

}

func AddOas(a *Oas) (string, error) {
	beego.Debug("[M] Got data:", a)
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return "", err
	}

	a.Id = uuid.New()
	beego.Debug("[M] Got new id:", a.Id)
	validator := new(validation.Validation)
	valid, err := validator.Valid(a)
	if err != nil {
		o.Rollback()
		return "", err
	}
	if !valid {
		o.Rollback()
		var errS string
		for _, err := range validator.Errors {
			errS = fmt.Sprintf("%s, %s:%s", errS, err.Key, err.Message)
		}
		return "", fmt.Errorf("Bad info: %s", errS)
	}
	beego.Debug("[M] Got new data:", a)
	_, err = o.Insert(a)
	if err != nil {
		o.Rollback()
		return "", err
	}
	beego.Debug("[M] Oas info saved")
	o.Commit()
	return a.Id, nil
}

func DeleteOas(a *Oas) error {
	beego.Debug("[M] Got data:", a)
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return err
	}
	validator := new(validation.Validation)
	valid, err := validator.Valid(a)
	if err != nil {
		o.Rollback()
		return err
	}
	if !valid {
		o.Rollback()
		var errS string
		for _, err := range validator.Errors {
			errS = fmt.Sprintf("%s, %s:%s", errS, err.Key, err.Message)
		}
		return fmt.Errorf("Bad info: %s", errS)
	}
	_, err = o.Delete(a)
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

func UpdateOas(a *Oas) error {
	beego.Debug("[M] Got data:", a)
	o := orm.NewOrm()
	err := o.Begin()
	if err != nil {
		return err
	}
	validator := new(validation.Validation)
	valid, err := validator.Valid(a)
	if err != nil {
		o.Rollback()
		return err
	}
	if !valid {
		o.Rollback()
		var errS string
		for _, err := range validator.Errors {
			errS = fmt.Sprintf("%s, %s:%s", errS, err.Key, err.Message)
		}
		return fmt.Errorf("Bad info: %s", errS)
	}
	_, err = o.Update(a)
	if err != nil {
		o.Rollback()
		return err
	}
	o.Commit()
	return nil
}

// If get all, just use &Oas{}
func GetOas(cond *Oas, limit, index int) ([]*Oas, error) {
	r := make([]*Oas, 0)
	o := orm.NewOrm()
	q := o.QueryTable("oas")
	if cond.Id != "" {
		q = q.Filter("id", cond.Id)
	}
	if cond.Endpoint != "" {
		q = q.Filter("endpoint", cond.Endpoint)
	}
	if cond.VaultName != "" {
		q = q.Filter("vault_name", cond.VaultName)
	}
	if limit > 0 {
		q = q.Limit(limit)
	}
	if index > 0 {
		q = q.Offset(index)
	}
	_, err := q.All(&r)

	if err != nil {
		return nil, err
	}
	for _, v := range r {
		o.LoadRelated(v, "BackupSets", common.RelDepth)
		o.LoadRelated(v, "Jobs", common.RelDepth)
	}
	return r, nil
}
