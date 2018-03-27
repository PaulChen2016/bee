<template>
    <div class="app-container calendar-list-container">
        <div class="filter-container">
            <el-button-group style='margin-left:0px'>
                <el-button @click="addBtnClick" class="filter-item" type="primary" icon="el-icon-plus">新增</el-button>
                <el-button @click="deleteBtnClick" class="filter-item" type="primary" icon="el-icon-delete" :disabled="multi_select.length <= 0">删除</el-button>
            </el-button-group>

            <el-button-group style='margin-left:15px'>
                <el-button class="filter-item" type="primary" @click="importScoreDialog.visible = true">导入</el-button>
                <el-button class="filter-item" type="primary" @click="handleExport">导出</el-button>
            </el-button-group>

            <div style="float:right">
                <el-select class="filter-item" v-model="query.key" placeholder="请选择">
                    <el-option label="XXXX" value="{{modelName}}.XXXX"> </el-option>

                </el-select>
                <el-input v-model="query.value" style="width: 200px;" class="filter-item" placeholder="请输入内容"></el-input>
                <el-button @click="queryaction" class="filter-item" type="primary" icon="el-icon-search">搜索</el-button>
            </div>
        </div>

        <el-table header-row-class-name="tblheader" :data="result.tabledata" stripe style="width: 100%" @selection-change="handleSelectionChange" @sort-change="sortTable" v-loading="listLoading" @cell-dblclick="celldblclick" @filter-change="filterchange">
            <el-table-column type="selection" width="45" fixed="left"></el-table-column>
            <el-table-column prop="BetaUser" label="用户名" width="150" sortable="custom">
                <template slot-scope="scope">
                    {{scope.row.BetaUser.Name}}
                </template>
            </el-table-column>
            <el-table-column prop="XXXX" label="XXXX" width="140" sortable="custom" column-key="XXXX" :filters="filterOptions.XXXX">
                <template slot-scope="scope">
                    {{{'-2':"兑换积分", '-1':"惩罚积分", 1:"奖励积分", 2:"问卷积分"}[scope.row.Type]}}
                </template>
            </el-table-column>
            <el-table-column prop="XXXX" label="XXXX" width="120" sortable="custom"></el-table-column>
            <el-table-column prop="XXXX" label="XXXX" sortable="custom" show-overflow-tooltip></el-table-column>

            <el-table-column prop="XXXX" label="操作" width="160" fixed="right">
                <template slot-scope="scope">
                    <el-tooltip class="item" effect="dark" content='编辑 XXXX' placement="top-start">
                        <el-button @click="handleEditEvt(scope.row)" icon="el-icon-edit" size="mini" type="primary" :loading="listLoading"></el-button>
                    </el-tooltip>
                    <el-tooltip class="item" effect="dark" content='删除 XXXX' placement="top-start">
                        <el-button @click="handleDeleteOne(scope.row.Id)" size="mini" type="danger" icon="el-icon-delete" :loading="listLoading"></el-button>
                    </el-tooltip>
                </template>
            </el-table-column>
        </el-table>

        <div v-show="!listLoading" class="pagination-container">
            <el-pagination @size-change="handleSizeChange" @current-change="handleCurrentChange" :current-page.sync="result.page" :page-sizes="[10,20,30, 50]" :page-size="result.limit" layout="total, sizes, prev, pager, next, jumper" :total="result.total">
            </el-pagination>
        </div>

        <el-dialog :title="newDialog.add?'新增':'修改'" :visible.sync="newDialog.formVisible">
            <el-form ref="itemForm" :model="newDialog.itemForm" :rules="newDialog.rules">
                <el-form-item prop="XXXX.Id" label="Beta用户" label-width="120px">
                    <el-select v-model="newDialog.itemForm.XXXX.Id" placeholder="Beta用户姓名(自动查询)" filterable clearable remote :remote-method="getBetaUsers" size="1000px">
                        <el-option v-for="user in newDialog.XXXXs" :label="user.Name + '(' + user.AccountID + ')'" :key="user.Id" :value="user.Id"> </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item prop="XXXX" label="XXXX" label-width="120px">
                    <!-- <el-input v-model="newDialog.itemForm.Product"></el-input> -->
                    <el-select v-model="newDialog.itemForm.XXXX" placeholder="选择XXXX" filterable clearable size="1000px">
                        <el-option v-for="XXXX in newDialog.XXXX" :label="XXXX" :key="XXXX" :value="XXXX"> </el-option>
                    </el-select>
                </el-form-item>
                <el-form-item prop="XXXX" label="积分类型" label-width="120px">
                    <el-select v-model="newDialog.itemForm.XXXX" placeholder="选择积分类型" filterable clearable size="1000px">
                        <el-option label="兑换积分" :value=-2> </el-option>
                    </el-select>
                </el-form-item>

                <el-form-item prop="Desc" label="备注" label-width="120px">
                    <el-input v-model="newDialog.itemForm.Desc" type="textarea"> </el-input>
                </el-form-item>

                <el-form-item label-width="120px">
                    <el-button @click="clickAddorModBtn" type="primary">{{newDialog.add?"立即创建":"保存修改"}}</el-button>
                    <el-button @click="newDialog.formVisible = false">取消</el-button>
                </el-form-item>
            </el-form>
        </el-dialog>

        <el-dialog title="详情" :visible.sync="cellDetailDialog.Visible">
            <pre>{{cellDetailDialog.detail}}</pre>
            <span slot="footer" class="dialog-footer">
                <el-button type="primary" @click="() => {cellDetailDialog.Visible = false; cellDetailDialog.detail =''}">确 定</el-button>
            </span>
        </el-dialog>

        <el-dialog title="导入XXXX" :visible.sync="importScoreDialog.visible">
            <el-upload class="upload-demo" :action="uploadAction" :on-preview="importPreview" :file-list="importScoreDialog.fileList" :on-success="importSuccess" :on-error="importError">
                <el-button size="small" type="primary">点击上传</el-button>
                <div slot="tip" class="el-upload__tip">只能上传 *.xlsx 文件，且不超过500kb</div>
            </el-upload>
            <span slot="footer" class="dialog-footer">
                <el-button type="primary" @click="importScoreDialog.visible = false">确 定</el-button>
            </span>
        </el-dialog>

    </div>
</template>
  
  <script>
import {
  getItemList,
  addItem,
  updateItem,
  deleteItems,
  deleteItem
} from "api/XXXX/{{modelName}}.js";

export default {
  data() {
    return {
      listLoading: false,
      exportLoading: false,
      uploadAction: "",
      importScoreDialog: {
        visible: false
      },
      query: {
        key: "",
        value: ""
      },
      filterOptions: {
        XXXX: [{ text: "兑换积分", value: -2 }]
      },
      filters: {
        // begin 过滤查询项的值
      },
      result: {
        tabledata: [],
        total: 0,
        limit: 10,
        page: 1,
        sortby: ["Id"],
        order: ["desc"]
      },
      importDialog: {
        visible: false
      },
      newDialog: {
        add: false,
        formVisible: false,
        BetaUsers: [],
        products: "",
        itemForm: {
          XXXX: {
            Id: ""
          },
          Name: "",
          Desc: ""
        },
        rules: {
          "XXXX.Id": [
            {
              type: "number",
              required: true,
              message: "xxxxxx",
              trigger: "change"
            }
          ],
          XXXX: [{ required: true, message: "请选择xxx", trigger: "change" }]
        }
      },
      multi_select: [],
      cellDetailDialog: {
        Visible: false,
        detail: ""
      }
    };
  },
  methods: {
    getList() {
      this.listLoading = true;
      const vm = this;
      let queryTmp = "";
      if (this.query.key !== "" && this.query.value !== "") {
        queryTmp = this.query.key + ":" + this.query.value;
      }

      // 增加表头过滤条件
      for (const key in this.filters) {
        // 如果过滤条件为空，下一个
        if (this.filters[key].length <= 0) {
          continue;
        }
        if (queryTmp !== "") {
          queryTmp =
            queryTmp +
            "," +
            key +
            "|multi-select" +
            ":" +
            this.filters[key].join("|");
        } else {
          queryTmp = key + "|multi-select" + ":" + this.filters[key].join("|");
        }
      }
      const query = {
        limit: this.result.limit,
        query: queryTmp,
        offset: (this.result.page - 1) * this.result.limit,
        sortby: this.result.sortby.join(","),
        order: this.result.order.join(",")
      };

      getItemList(query).then(response => {
        vm.result.total = response.data.Msg.Total;
        vm.result.tabledata = response.data.Msg.Items;
        console.log(vm.result.tabledata);
        vm.listLoading = false;
      });
    },
    queryaction() {
      this.getList();
    },
    handleSizeChange(val) {
      this.result.limit = val;
      this.getList();
    },
    handleCurrentChange(val) {
      this.result.page = val;
      this.getList();
    },
    handleEditEvt(row) {
      this.listLoading = true;
      this.newDialog.BetaUsers = [row.BetaUser]; // 修改是，下来菜单显示用户名
      for (const p in row) {
        this.newDialog.itemForm[p] = row[p];
      }
      this.listLoading = false;
      this.newDialog.add = false;
      this.newDialog.formVisible = true;
    },
    clickAddorModBtn() {
      const _this = this;
      this.$refs.itemForm.validate(valid => {
        if (valid) {
          // 根据积分类型，计算积分正负数
          if (_this.newDialog.itemForm.Type < 0) {
            _this.newDialog.itemForm.Score = -Math.abs(
              _this.newDialog.itemForm.Score
            );
          } else {
            _this.newDialog.itemForm.Score = Math.abs(
              _this.newDialog.itemForm.Score
            );
          }
          if (_this.newDialog.add) {
            _this.newDialog.itemForm.Id = 0; // 清除修改时记录的Id信息，避免造成添加重复
            addItem(_this.newDialog.itemForm).then(() => {
              _this.$notify({
                title: "消息",
                message: "积分流水添加成功",
                offset: 100,
                type: "success"
              });
              _this.getList();
              _this.newDialog.formVisible = false;
            });
          } else {
            updateItem(_this.newDialog.itemForm).then(() => {
              _this.$notify({
                title: "消息",
                message: "积分流水修改成功",
                offset: 100,
                type: "success"
              });
              _this.getList();
              _this.newDialog.formVisible = false;
            });
          }
        }
      });
    },
    handleSelectionChange(val) {
      this.multi_select = val.reduce((prev, next) => {
        prev.push(next.Id);
        return prev;
      }, []);
    },
    addBtnClick() {
      this.newDialog.add = true;
      this.newDialog.formVisible = true;
    },
    deleteBtnClick() {
      this.$confirm("此操作将永久删除选中的积分流水, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
        .then(() => {
          deleteItems(this.multi_select.join(",")).then(
            () => {
              this.$notify({
                title: "消息",
                message: "删除成功",
                offset: 100,
                type: "success"
              });
              this.getList();
            },
            errRespnse => {
              this.$alert(errRespnse, "提示", {
                confirmButtonText: "确定"
              });
            }
          );
        })
        .catch(() => {
          console.log("取消删除");
        });
    },
    handleDeleteOne(id) {
      this.$confirm("此操作将永久删除该积分流水, 是否继续?", "提示", {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning"
      })
        .then(() => {
          deleteItem(id).then(
            () => {
              this.$notify({
                title: "消息",
                message: "删除成功",
                offset: 100,
                type: "success"
              });
              this.getList();
            },
            errRespnse => {
              this.$alert(errRespnse, "提示", {
                confirmButtonText: "确定"
              });
            }
          );
        })
        .catch(() => {
          console.log("取消删除");
        });
    },
    sortTable(col) {
      this.result.sortby = [];
      this.result.order = [];
      if (col.prop) {
        this.result.sortby.push(col.prop.trim());
        this.result.order.push(col.order.replace("ending", ""));
      }
      this.getList();
    },
    celldblclick(row, column) {
      // this.cellDetailDialog.detail = row[column.property]
      this.cellDetailDialog.detail = column.property
        .split(".")
        .reduce((prev, next) => prev[next], row);
      this.cellDetailDialog.Visible = true;
    },
    filterchange(filter) {
      // 存储过滤条件到filter对象中
      for (const key in filter) {
        this.filters[key] = filter[key];
      }
      this.getList();
    },
    importPreview(file) {
      console.log(file.type);
      const isExcel = file.type === "xlsx";
      if (!isExcel) {
        this.$message.error("上传文件只能是xlsx格式!");
      }
    },
    importSuccess(response, file, fileList) {
      console.log(response, file, fileList);
      this.getList();
    },
    importError(err, file, fileList) {
      console.log(err);
    },
    handleExport() {
      this.exportLoading = true;
      const vm = this;
      let queryTmp = "";
      if (this.query.key !== "" && this.query.value !== "") {
        queryTmp = this.query.key + ":" + this.query.value;
      }

      // 增加表头过滤条件
      for (const key in this.filters) {
        // 如果过滤条件为空，下一个
        if (this.filters[key].length <= 0) {
          continue;
        }
        if (queryTmp !== "") {
          queryTmp =
            queryTmp +
            "," +
            key +
            "|multi-select" +
            ":" +
            this.filters[key].join("|");
        } else {
          queryTmp = key + "|multi-select" + ":" + this.filters[key].join("|");
        }
      }

      const query = {
        limit: -1,
        query: queryTmp,
        offset: (this.result.page - 1) * this.result.limit,
        sortby: this.result.sortby.join(","),
        order: this.result.order.join(",")
      };
      getItemList(query).then(response => {
        require.ensure([], () => {
          const { export_json_to_excel } = require("vendor/Export2Excel");
          const tHeader = [
            "用户名",
            "用户账号ID",
            "产品",
            "积分类型",
            "积分值",
            "备注"
          ];
          const filterVal = [
            "Name",
            "AccountID",
            "Product",
            "Type",
            "Score",
            "Desc"
          ];
          let data = response.data.Msg.Items;
          if (vm.multi_select.length > 0) {
            data = data.filter(item => {
              if (vm.multi_select.some(i => item.Id === i)) {
                return true;
              }
              return false;
            });
          }
          data = data.map(v => {
            return filterVal.map(j => {
              if (j === "Name") {
                if (!v["BetaUser"]) {
                  return null;
                } else {
                  return v["BetaUser"].Name;
                }
              } else if (j === "AccountID") {
                if (!v["BetaUser"]) {
                  return null;
                } else {
                  return v["BetaUser"].AccountID;
                }
              }
              return v[j];
            });
          });
          data.unshift(filterVal);
          export_json_to_excel(tHeader, data, "userScore-export");
          this.exportLoading = false;
        });
      });
    }
  },
  created() {
    this.getList();
    getProducts().then(response => {
      this.newDialog.products = response.data.Msg;
    });
    if (process.env.NODE_ENV === "production") {
      this.uploadAction =
        "http://" + document.location.host + "/entiy/userscore/import";
    } else {
      this.uploadAction = " http://127.0.0.1:20000/entiy/userscore/import";
    }
  }
};
</script>
  
  <style rel="stylesheet/scss" lang="scss" scoped>
.app-container {
  margin-top: 5px;
  text-align: left;
  .filter-container {
    margin-bottom: 15px;
  }
}

.el-button-group {
  margin-left: 10px;
}
</style>
<style>
.tblheader > th {
  background-color: rgb(241, 246, 253);
}
</style>