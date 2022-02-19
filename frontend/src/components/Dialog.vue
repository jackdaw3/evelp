<template>
  <div class="Dialog">
    <el-button icon="el-icon-setting" circle size="medium" @click="dialogVisible = true"></el-button>&nbsp;
    <el-dialog
      :title="dialogLabel.title"
      :visible.sync="dialogVisible"
      height="60%"
      width="32%"
    >
      <el-form label-width="20%" style="margin-top: -3%" v-model="formData">
        <el-form-item :label="dialogLabel.materialPrice">
          <el-select
            :placeholder="dialogLabel.materialPlaceholder"
            v-model="formData.materialPrice"
            style="width:90%"
          >
            <el-option :label="dialogLabel.buyPrice" value="buy"></el-option>
            <el-option :label="dialogLabel.sellPrice" value="sell"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="dialogLabel.productPrice">
          <el-select
            :placeholder="dialogLabel.productPlaceholder"
            v-model="formData.productPrice"
            style="width:90%"
          >
            <el-option :label="dialogLabel.buyPrice" value="buy"></el-option>
            <el-option :label="dialogLabel.sellPrice" value="sell"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="dialogLabel.scope">
          <el-select
            :placeholder="dialogLabel.scopePlaceholder"
            v-model="formData.scope"
            style="width:90%"
          >
            <el-option label="1%" value="0.01"></el-option>
            <el-option label="5%" value="0.05"></el-option>
            <el-option label="10%" value="0.1"></el-option>
            <el-option label="25%" value="0.25"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="dialogLabel.days">
          <el-select
            :placeholder="dialogLabel.daysPlaceholder"
            v-model="formData.days"
            style="width:90%"
          >
            <el-option :label="dialogLabel.week" value="7"></el-option>
            <el-option :label="dialogLabel.month" value="30"></el-option>
            <el-option :label="dialogLabel.season" value="90"></el-option>
          </el-select>
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>
<script>
export default {
  props: {
    form: Object,
  },
  data() {
    return {
      dialogVisible: false,
      dialogLabel: this.$t("message.dialog"),
      formData: this.form,
    };
  },
  watch: {
    formData(value) {
      this.$emit("formChange", value);
    },
    "$i18n.locale"() {
      this.dialogLabel = this.$t("message.dialog");
    },
  },
};
</script>
