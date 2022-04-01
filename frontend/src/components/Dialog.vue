<template>
  <div class="Dialog">
    <el-button
      icon="el-icon-setting"
      circle
      size="medium"
      @click="dialogVisible = true"
    ></el-button
    >&nbsp;
    <el-dialog
      :title="dialogLabel.title"
      :visible.sync="dialogVisible"
      @close="closeDialog"
      width="32%"
      style="background-color:rgba(0, 0, 0, 0.7)"
    >
      <el-form label-width="20%" style="margin-top: -3%" v-model="form">
        <el-form-item :label="dialogLabel.materialPrice">
          <el-select
            :placeholder="dialogLabel.materialPlaceholder"
            v-model="form.materialPrice"
            style="width: 90%"
          >
            <el-option :label="dialogLabel.buyPrice" value="buy"></el-option>
            <el-option :label="dialogLabel.sellPrice" value="sell"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="dialogLabel.productPrice">
          <el-select
            :placeholder="dialogLabel.productPlaceholder"
            v-model="form.productPrice"
            style="width: 90%"
          >
            <el-option :label="dialogLabel.buyPrice" value="buy"></el-option>
            <el-option :label="dialogLabel.sellPrice" value="sell"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="dialogLabel.scope">
          <el-select
            :placeholder="dialogLabel.scopePlaceholder"
            v-model="form.scope"
            style="width: 90%"
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
            v-model="form.days"
            style="width: 90%"
          >
            <el-option :label="dialogLabel.week" value="7"></el-option>
            <el-option :label="dialogLabel.month" value="30"></el-option>
            <el-option :label="dialogLabel.season" value="90"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item :label="dialogLabel.tax">
          <el-slider
            v-model="form.tax"
            :max="20"
            show-input
            :format-tooltip="taxFormat"
            style="width: 90%"
          ></el-slider>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="reset" size="medium">{{
          dialogLabel.reset
        }}</el-button>
        <el-button @click="dialogVisible = false" size="medium">{{
          dialogLabel.close
        }}</el-button>
      </span>
    </el-dialog>
  </div>
</template>
<script>
export default {
  computed: {
    form: function () {
      return this.$store.state.form;
    },
  },
  data() {
    return {
      dialogVisible: false,
      dialogLabel: this.$t("message.dialog"),
    };
  },
  methods: {
    reset() {
      this.form.materialPrice = "sell";
      this.form.productPrice = "buy";
      this.form.days = "7";
      this.form.scope = "0.05";
      this.form.tax = 0;
    },
    closeDialog() {
      if (!localStorage.form) {
        this.$store.dispatch("setForm", this.form);
        localStorage.form = JSON.stringify(this.form);
        this.$emit("form-change");
      } else {
        const old = JSON.parse(localStorage.form);
        if (
          this.form.materialPrice != old.materialPrice ||
          this.form.productPrice != old.productPrice ||
          this.form.days != old.days ||
          this.form.tax != old.tax
        ) {
          this.$store.dispatch("setForm", this.form);
          localStorage.form = JSON.stringify(this.form);
          this.$emit("form-change");
        }
      }
    },
    taxFormat(e) {
      return e + "%";
    },
  },
  watch: {
    "$i18n.locale"() {
      this.dialogLabel = this.$t("message.dialog");
    },
  },
};
</script>
