// Copyright 2013 bee authors
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package generate

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	beeLogger "github.com/PaulChen2016/bee/logger"
	"github.com/PaulChen2016/bee/logger/colors"
	"github.com/PaulChen2016/bee/utils"
)

func GenerateModel(mname, fields, currpath string) {
	w := colors.NewColorWriter(os.Stdout)

	p, f := path.Split(mname)
	modelName := strings.Title(f)
	packageName := "models"
	if p != "" {
		i := strings.LastIndex(p[:len(p)-1], "/")
		packageName = p[i+1 : len(p)-1]
	}

	modelStruct, hastime, err := getStruct(modelName, fields)
	if err != nil {
		beeLogger.Log.Fatalf("Could not generate the model struct: %s", err)
	}

	beeLogger.Log.Infof("Using '%s' as model name", modelName)
	beeLogger.Log.Infof("Using '%s' as package name", packageName)

	fp := path.Join(currpath, "models", p)
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		// Create the model's directory
		if err := os.MkdirAll(fp, 0777); err != nil {
			beeLogger.Log.Fatalf("Could not create the model directory: %s", err)
		}
	}

	fpath := path.Join(fp, strings.ToLower(modelName)+".go")
	if f, err := os.OpenFile(fpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
		defer utils.CloseFile(f)
		content := strings.Replace(modelTpl, "{{packageName}}", packageName, -1)
		content = strings.Replace(content, "{{modelName}}", modelName, -1)
		content = strings.Replace(content, "{{modelStruct}}", modelStruct, -1)
		if hastime {
			content = strings.Replace(content, "{{timePkg}}", `"time"`, -1)
		} else {
			content = strings.Replace(content, "{{timePkg}}", "", -1)
		}
		f.WriteString(content)
		// Run 'gofmt' on the generated source code
		utils.FormatSourceCode(fpath)
		fmt.Fprintf(w, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", fpath, "\x1b[0m")
	} else {
		beeLogger.Log.Fatalf("Could not create model file: %s", err)
	}
}

func getStruct(structname, fields string) (string, bool, error) {
	if fields == "" {
		return "", false, errors.New("fields cannot be empty")
	}

	hastime := false
	structStr := "type " + structname + " struct{\n"
	fds := strings.Split(fields, ",")
	for i, v := range fds {
		kv := strings.SplitN(v, ":", 2)
		if len(kv) != 2 {
			return "", false, errors.New("the fields format is wrong. Should be key:type,key:type " + v)
		}

		typ, tag, hastimeinner := getType(kv[1])
		if typ == "" {
			return "", false, errors.New("the fields format is wrong. Should be key:type,key:type " + v)
		}

		if i == 0 && strings.ToLower(kv[0]) != "id" {
			structStr = structStr + "ID     int64     `orm:\"auto\"`\n"
		}

		if hastimeinner {
			hastime = true
		}
		structStr = structStr + utils.CamelString(kv[0]) + "       " + typ + "     " + tag + "\n"
	}
	structStr += "}\n"
	return structStr, hastime, nil
}

// fields support type
// http://beego.me/docs/mvc/model/models.md#mysql
func getType(ktype string) (kt, tag string, hasTime bool) {
	kv := strings.SplitN(ktype, ":", 2)
	switch kv[0] {
	case "string":
		if len(kv) == 2 {
			return "string", "`orm:\"size(" + kv[1] + ")\"`", false
		}
		return "string", "`orm:\"size(128)\"`", false
	case "text":
		return "string", "`orm:\"type(longtext)\"`", false
	case "auto":
		return "int64", "`orm:\"auto\"`", false
	case "pk":
		return "int64", "`orm:\"pk\"`", false
	case "datetime":
		return "time.Time", "`orm:\"type(datetime)\"`", true
	case "int", "int8", "int16", "int32", "int64":
		fallthrough
	case "uint", "uint8", "uint16", "uint32", "uint64":
		fallthrough
	case "bool":
		fallthrough
	case "float32", "float64":
		return kv[0], "", false
	case "float":
		return "float64", "", false
	}
	return "", "", false
}

var modelTpl = `package {{packageName}}

import (
	"errors"
	"reflect"
	"strconv"
	{{timePkg}}
	"github.com/PaulChen2016/tools/curd/common"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

// {{modelName}} ...
{{modelStruct}}

func init() {
	orm.RegisterModel(new({{modelName}}))
}

// Add{{modelName}} insert a new {{modelName}} into database and returns
// last inserted ID on success.
func Add{{modelName}}(m *{{modelName}}) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// Get{{modelName}}ByID retrieves {{modelName}} by ID. Returns error if
// ID doesn't exist
func Get{{modelName}}ByID(id int64) (v *{{modelName}}, err error) {
	o := orm.NewOrm()
	v = &{{modelName}}{ID: id}
	if err = o.QueryTable(new({{modelName}})).Filter("ID", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAll{{modelName}} retrieves all {{modelName}} matches certain condition. Returns empty list if
// no records exist
func GetAll{{modelName}}(querys []*common.QueryConditon, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, totalcount int64, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new({{modelName}}))
	// query QueryCondition
	cond := orm.NewCondition()
	for _, query := range querys {
		var k string
		cond1 := orm.NewCondition()
		switch query.QueryType {
		case common.MultiSelect:
			k = query.QueryKey // + "__iexact"
			for _, v := range query.QueryValues {
				cond1 = cond1.Or(k, v)
			}
			cond = cond.AndCond(cond1)
		case common.MultiText:
			k = query.QueryKey + "__icontains"
			for _, v := range query.QueryValues {
				cond1 = cond1.Or(k, v)
			}
			cond = cond.AndCond(cond1)
		case common.NumRange:
			if len(query.QueryValues) == 2 {
				var from, to float64
				if from, err = strconv.ParseFloat(query.QueryValues[0], 64); err != nil {
					logs.Error(err.Error())
					return
				}
				if to, err = strconv.ParseFloat(query.QueryValues[1], 64); err != nil {
					logs.Error(err.Error())
					return
				}
				k = query.QueryKey + "__gte"
				cond1 = cond1.Or(k, from)
				k = query.QueryKey + "__lte"
				cond1 = cond1.And(k, to)
				cond = cond.AndCond(cond1)
			} else {
				k = query.QueryKey + "__icontains"
				for _, v := range query.QueryValues {
					cond1 = cond1.Or(k, v)
				}
				cond = cond.AndCond(cond1)
			}
		default:
			k = query.QueryKey + "__icontains"
			for _, v := range query.QueryValues {
				cond1 = cond1.Or(k, v)
			}
			cond = cond.AndCond(cond1)
		}
	}
	qs = qs.SetCond(cond)
	//获取查询条件过滤的总记录数
	if totalcount, err = qs.Count(); err != nil {
		return
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, 0, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, 0, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, 0, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, 0, errors.New("Error: unused 'order' fields")
		}
	}

	var l []{{modelName}}
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, totalcount, nil
	}
	return nil, 0, err
}

// Update{{modelName}}ByID and returns error if
// the record to be updated doesn't exist
func Update{{modelName}}ByID(m *{{modelName}}, fields []string) (err error) {
	o := orm.NewOrm()
	v := {{modelName}}{ID: m.ID}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m, fields...); err == nil {
			logs.Debug("Number of {{modelName}} update in database:", num)
		}
	}
	return
}

// Delete{{modelName}} deletes {{modelName}} by ID and returns error if
// the record to be deleted doesn't exist
func Delete{{modelName}}(id int64) (err error) {
	o := orm.NewOrm()
	v := {{modelName}}{ID: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&{{modelName}}{ID: id}); err == nil {
			logs.Debug("Number of {{modelName}} deleted in database:", num)
		}
	}
	return
}

// MultDelete{{modelName}}ByIDs and returns error if
// the delete not success
func MultDelete{{modelName}}ByIDs(ids []interface{}) (err error) {
	o := orm.NewOrm()
	if num, err := o.QueryTable(new({{modelName}})).Filter("id__in", ids...).Delete(); err == nil {
		logs.Debug("delete {{modelName}} from database, num:", num)
		return nil
	}
	return
}
`
