<template>
  <div class="Table">
    <el-table :data="tableData2" id="table" :cell-style="tableStyle" style="width: 100%">
      <el-table-column prop="date" min-width="50%" label="日期"></el-table-column>
      <el-table-column prop="name" min-width="50%" label="姓名"></el-table-column>
      <el-table-column prop="address" min-width="50%" label="地址"></el-table-column>
    </el-table>
  </div>
</template>

  <script>
export default {
  props: {
    table: Array,
    form: Object,
  },
  data() {
    return {
      tableData2: [
        {
          date: "2016-05-02",
          name: "王小虎",
          address: "上海市普陀区金沙江路 1518 弄",
        },
        {
          date: "2016-05-04",
          name: "王小虎",
          address: "上海市普陀区金沙江路 1517 弄",
        },
        {
          date: "2016-05-01",
          name: "王小虎",
          address: "上海市普陀区金沙江路 1519 弄",
        },
        {
          date: "2016-05-03",
          name: "王小虎",
          address: "上海市普陀区金沙江路 1516 弄",
        },
      ],
    };
  },
  methods: {
    tableStyle() {
      return {
        padding: "0",
        "font-size": "14px",
      };
    },
  },
};
</script>