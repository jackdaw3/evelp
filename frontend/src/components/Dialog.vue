<template>
  <div class="Dialog">
    <el-button icon="el-icon-setting" circle size="medium" @click="dialogVisible = true"></el-button>&nbsp;
    <el-dialog
      :title="dialog.title"
      :visible.sync="dialogVisible"
      height="60%"
      width="32%"
      :before-close="handleClose"
    >
      <el-form label-width="20%" style="margin-top: -3%" :model="formLabelAlign">
        <el-form-item :label="dialog.materialPrice">
          <el-select :placeholder="dialog.materialPlaceholder" style="width:90%">
            <el-option :label="dialog.buyPrice" value="buy"></el-option>
            <el-option :label="dialog.sellPrice" value="sell"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="dialog.offerPrice">
          <el-select :placeholder="dialog.offerPlaceholder" style="width:90%">
            <el-option :label="dialog.buyPrice" value="buy"></el-option>
            <el-option :label="dialog.sellPrice" value="sell"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="dialog.scope">
          <el-tooltip class="item" effect="dark" :content="dialog.scopeContent" placement="right">
            <el-select :placeholder="dialog.scopePlaceholder" style="width:90%">
              <el-option label="1%" value="0.01"></el-option>
              <el-option label="5%" value="0.05"></el-option>
              <el-option label="10%" value="0.1"></el-option>
              <el-option label="25%" value="0.25"></el-option>
            </el-select>
          </el-tooltip>
        </el-form-item>
        <el-form-item :label="dialog.days">
          <el-tooltip class="item" effect="dark" :content="dialog.daysContent" placement="right">
            <el-select :placeholder="dialog.daysPlaceholder" style="width:90%">
              <el-option :label="dialog.week" value="7"></el-option>
              <el-option :label="dialog.month" value="30"></el-option>
              <el-option :label="dialog.season" value="90"></el-option>
            </el-select>
          </el-tooltip>
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>
<script>
export default {
  data() {
    return {
      dialogVisible: false,
      dialog: this.$t("message.dialog"),
    };
  },
  watch: {
    "$i18n.locale"() {
      this.dialog = this.$t("message.dialog");
    },
  },
};
</script>
