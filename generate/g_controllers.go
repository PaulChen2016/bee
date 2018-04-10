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
	"fmt"
	"os"
	"path"
	"strings"

	beeLogger "github.com/PaulChen2016/bee/logger"
	"github.com/PaulChen2016/bee/logger/colors"
	"github.com/PaulChen2016/bee/utils"
)

func GenerateController(cname, currpath string) {
	w := colors.NewColorWriter(os.Stdout)

	p, f := path.Split(cname)
	controllerName := strings.Title(f)
	packageName := "controllers"

	if p != "" {
		i := strings.LastIndex(p[:len(p)-1], "/")
		packageName = p[i+1 : len(p)-1]
	}

	beeLogger.Log.Infof("Using '%s' as controller name", controllerName)
	beeLogger.Log.Infof("Using '%s' as package name", packageName)

	fp := path.Join(currpath, "controllers", p)
	if _, err := os.Stat(fp); os.IsNotExist(err) {
		// Create the controller's directory
		if err := os.MkdirAll(fp, 0777); err != nil {
			beeLogger.Log.Fatalf("Could not create controllers directory: %s", err)
		}
	}

	fpath := path.Join(fp, strings.ToLower(controllerName)+".go")
	if f, err := os.OpenFile(fpath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0666); err == nil {
		defer utils.CloseFile(f)

		modelPath := path.Join(currpath, "models", strings.ToLower(controllerName)+".go")

		var content string
		if _, err := os.Stat(modelPath); err == nil {
			beeLogger.Log.Infof("Using matching model '%s'", controllerName)
			content = strings.Replace(controllerModelTpl, "{{packageName}}", packageName, -1)
			pkgPath := getPackagePath(currpath)
			content = strings.Replace(content, "{{pkgPath}}", pkgPath, -1)
		} else {
			content = strings.Replace(controllerTpl, "{{packageName}}", packageName, -1)
		}

		content = strings.Replace(content, "{{controllerName}}", controllerName, -1)
		f.WriteString(content)

		// Run 'gofmt' on the generated source code
		utils.FormatSourceCode(fpath)
		fmt.Fprintf(w, "\t%s%screate%s\t %s%s\n", "\x1b[32m", "\x1b[1m", "\x1b[21m", fpath, "\x1b[0m")
	} else {
		beeLogger.Log.Fatalf("Could not create controller file: %s", err)
	}
}

var controllerTpl = `package {{packageName}}

import (
	"github.com/astaxie/beego"
)

// {{controllerName}}Controller operations for {{controllerName}}
type {{controllerName}}Controller struct {
	beego.Controller
}

// URLMapping ...
func (c *{{controllerName}}Controller) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Create
// @Description create {{controllerName}}
// @Param	body		body 	models.{{controllerName}}	true		"body for {{controllerName}} content"
// @Success 201 {object} models.{{controllerName}}
// @Failure 403 body is empty
// @router / [post]
func (c *{{controllerName}}Controller) Post() {

}

// GetOne ...
// @Title GetOne
// @Description get {{controllerName}} by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.{{controllerName}}
// @Failure 403 :id is empty
// @router /:id [get]
func (c *{{controllerName}}Controller) GetOne() {

}

// GetAll ...
// @Title GetAll
// @Description get {{controllerName}}
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.{{controllerName}}
// @Failure 403
// @router / [get]
func (c *{{controllerName}}Controller) GetAll() {

}

// Put ...
// @Title Put
// @Description update the {{controllerName}}
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.{{controllerName}}	true		"body for {{controllerName}} content"
// @Success 200 {object} models.{{controllerName}}
// @Failure 403 :id is not int
// @router /:id [put]
func (c *{{controllerName}}Controller) Put() {

}

// Delete ...
// @Title Delete
// @Description delete the {{controllerName}}
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *{{controllerName}}Controller) Delete() {

}
`

var controllerModelTpl = `package {{packageName}}

import (
	"{{pkgPath}}/models"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/PaulChen2016/common"
	"github.com/astaxie/beego/logs"
	"github.com/tealeg/xlsx"
)

// {{controllerName}}Controller operations for {{controllerName}}
type {{controllerName}}Controller struct {
	common.BaseController
}

// URLMapping ...
func (c *{{controllerName}}Controller) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
	c.Mapping("DeleteList", c.DeleteList)
	c.Mapping("Import", c.Import)
}

// Post ...
// @Title Post
// @Description create {{controllerName}}
// @Param	body		body 	models.{{controllerName}}	true		"body for {{controllerName}} content"
// @Success 201 {int} models.{{controllerName}}
// @Failure 403 body is empty
// @router / [post]
func (c *{{controllerName}}Controller) Post() {
	var v models.{{controllerName}}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if _, err := models.Add{{controllerName}}(&v); err == nil {
		c.Ctx.Output.SetStatus(201)
		c.Data["json"] = common.RestfulResult{Code: 0, Msg: v}
	} else {
		c.Data["json"] = common.RestfulResult{Code: -1, Msg: err.Error()}
	}
	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get {{controllerName}} by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.{{controllerName}}
// @Failure 403 :id is empty
// @router /:id [get]
func (c *{{controllerName}}Controller) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.Get{{controllerName}}ById(id)
	if err != nil {
		c.Data["json"] = common.RestfulResult{Code: -1, Msg: err.Error()}
	} else {
		c.Data["json"] = common.RestfulResult{Code: 0, Msg: v}
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get {{controllerName}}
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.{{controllerName}}
// @Failure 403
// @router / [get]
func (c *{{controllerName}}Controller) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query []*common.QueryConditon
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k|type:v|v|v,k|type:v|v|v  其中Type可以没有,默认值是 MultiText
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") { // 分割多个查询key
			qcondtion := new(common.QueryConditon)
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key:value pair," + cond)
				c.ServeJSON()
				return
			}
			k_init, v_init := kv[0], kv[1]         // 初始分割查询key和value（备注，value是多个用|分割）
			key_type := strings.Split(k_init, "|") // 解析key中的type信息
			if len(key_type) == 2 {
				qcondtion.QueryKey = key_type[0]
				qcondtion.QueryType = key_type[1]
			} else if len(key_type) == 1 {
				qcondtion.QueryKey = key_type[0]
				qcondtion.QueryType = common.MultiText
			} else {
				c.Data["json"] = errors.New("Error: invalid query key|type format," + k_init)
				c.ServeJSON()
				return
			}
			qcondtion.QueryValues = strings.Split(v_init, "|") // 解析出values信息
			qcondtion.QueryKey = strings.Replace(qcondtion.QueryKey, ".", "__", -1)
			query = append(query, qcondtion)
		}
	}
	l, count, err := models.GetAll{{controllerName}}(query, fields, sortby, order, offset, limit)
	if err != nil {
		c.Data["json"] = common.RestfulResult{Code: -1, Msg: err.Error()}
	} else {
		c.Data["json"] = common.RestfulResult{Code: 0, Msg: struct {
			Items interface{}
			Total int64
		}{
			Items: l,
			Total: count,
		}}
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the {{controllerName}}
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.{{controllerName}}	true		"body for {{controllerName}} content"
// @Success 200 {object} models.{{controllerName}}
// @Failure 403 :id is not int
// @router /:id [put]
func (c *{{controllerName}}Controller) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.{{controllerName}}{Id: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.Update{{controllerName}}ById(&v); err == nil {
		c.Data["json"] = common.RestfulResult{Code: 0, Msg: "OK"}
	} else {
		c.Data["json"] = common.RestfulResult{Code: -1, Msg: err.Error()}
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the {{controllerName}}
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *{{controllerName}}Controller) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.Delete{{controllerName}}(id); err == nil {
		c.Data["json"] = common.RestfulResult{Code: 0, Msg: "OK"}
	} else {
		c.Data["json"] = common.RestfulResult{Code: -1, Msg: err.Error()}
	}
	c.ServeJSON()
}


// DeleteList ...
// @Title multi-Delete
// @Description delete multi {{controllerName}}s
// @Param	ids	 	string	true		"The ids you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /deletelist [delete]
func (c *{{controllerName}}Controller) DeleteList() {
	// fields: col1,col2,entity.col3
	idslice := []interface{}{}
	if v := c.GetString("ids"); v != "" {
		s := strings.Split(v, ",")
		for _, id := range s {
			idslice = append(idslice, id)
		}
	}
	if err := models.MultDelete{{controllerName}}ByIDs(idslice); err != nil {
		c.Data["json"] = common.RestfulResult{Code: -1, Msg: err.Error()}
	} else {
		c.Data["json"] = common.RestfulResult{Code: 0, Msg: "OK"}
	}
	c.ServeJSON()
}


// Import ...
// @Title 批量导入 {{controllerName}}
// @Description import {{controllerName}}
// @Param   excel file
// @Success 200 {object} import success!
// @Failure 403 file context is incorrect
// @router /import [post]
func (c *{{controllerName}}Controller) Import() {
	var err error
	fpath := filepath.Join(common.GetUploadPath(), "{{controllerName}}"+time.Now().Format("2006-01-02-15-04-05")+".xlsx")
	err = c.SaveToFile("file", fpath)
	if err != nil {
		logs.Error("Upload {{controllerName}} excel file err,", err.Error())
	}
	if err = import{{controllerName}}(fpath); err == nil {
		c.Data["json"] = common.RestfulResult{Code: 0, Msg: ""}
		//删除上传文件
		os.Remove(fpath)
	} else {
		c.Data["json"] = common.RestfulResult{Code: -1, Msg: err.Error()}
	}
	c.ServeJSON()
}

//import 具体方法
func import{{controllerName}}(fpath string) (err error) {
	var xlFile *xlsx.File
	xlFile, err = xlsx.OpenFile(fpath)
	for _, sheet := range xlFile.Sheets {
		if len(sheet.Rows) <= 2 {
			logs.Info("sheet context is null of import file,", fpath, sheet.Name)
			continue
		}
		headRow := sheet.Rows[1]
		for _, row := range sheet.Rows[2:] {
			m := models.{{controllerName}}{}
			s := reflect.ValueOf(&m).Elem()
			for col, cell := range row.Cells {
				// v, _ := cell.String()
				var v interface{}
				attr, _ := headRow.Cells[col].String()
				if !s.FieldByName(attr).IsValid() { //不识别的属性，继续
					logs.Warning("unknow attr:", attr)
					continue
				}
				field := s.FieldByName(attr).Type().Kind()
				switch field {
				case reflect.Float64:
					v, err = cell.Float()
				case reflect.String:
					v, err = cell.String()
				case reflect.Int:
					v, err = cell.Int()
				case reflect.Bool:
					v = cell.Bool()
				case reflect.Int64:
					v, err = cell.Int64()
				case reflect.Ptr:
					// switch attr {
					// case "BetaUser":
					// 	var accountID string
					// 	if accountID, err = cell.String(); err == nil {
					// 		v, err = models.GetBetaUserByAccountID(accountID)
					// 	}
					// default:
					// 	err = fmt.Errorf("Unkown BetaUserScore ptr attr type: %s", attr)
					// 	logs.Error(err.Error())
					// 	return err
					// }
				default:
					err = fmt.Errorf("Unkown {{controllerName}} attr type: %s", field)
					logs.Error(err.Error())
					return err
				}
				if err != nil {
					logs.Error("parse cell value failed when import {{controllerName}}s,", err.Error())
					return err
				}
				refvalue := reflect.ValueOf(v)
				s.FieldByName(attr).Set(refvalue)
			}
			if _, err = models.Add{{controllerName}}(&m); err != nil {
				logs.Error("models.Add{{controllerName}} failed,", err.Error())
				return err
			}
		}
	}
	return
}
`
