<template>
  <div class="Dialog">
    <el-button icon="el-icon-setting" circle size="medium" @click="dialogVisible = true"></el-button>&nbsp;
    <el-dialog :title="dialogLabel.title" :visible.sync="dialogVisible" @close="closeDialog" width="38%"
      >
      <el-tabs v-model="activeName" type="card">
        <el-tab-pane :label="dialogLabel.dataTitle" name="first">
          <el-form label-width="20%" style="margin-top: 2.5%" v-model="form">
            <el-form-item :label="dialogLabel.materialPrice">
              <el-select :placeholder="dialogLabel.materialPlaceholder" v-model="form.materialPrice" style="width: 90%">
                <el-option :label="dialogLabel.buyPrice" value="buy"></el-option>
                <el-option :label="dialogLabel.sellPrice" value="sell"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item :label="dialogLabel.productPrice">
              <el-select :placeholder="dialogLabel.productPlaceholder" v-model="form.productPrice" style="width: 90%">
                <el-option :label="dialogLabel.buyPrice" value="buy"></el-option>
                <el-option :label="dialogLabel.sellPrice" value="sell"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item :label="dialogLabel.scope">
              <el-select :placeholder="dialogLabel.scopePlaceholder" v-model="form.scope" style="width: 90%">
                <el-option label="1%" value="0.01"></el-option>
                <el-option label="5%" value="0.05"></el-option>
                <el-option label="10%" value="0.1"></el-option>
                <el-option label="25%" value="0.25"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item :label="dialogLabel.days">
              <el-select :placeholder="dialogLabel.daysPlaceholder" v-model="form.days" style="width: 90%">
                <el-option :label="dialogLabel.week" value="7"></el-option>
                <el-option :label="dialogLabel.month" value="30"></el-option>
                <el-option :label="dialogLabel.season" value="90"></el-option>
              </el-select>
            </el-form-item>
            <el-form-item :label="dialogLabel.tax">
              <el-slider v-model="form.tax" :max="20" show-input :format-tooltip="taxFormat"
                style="width: 90%"></el-slider>
            </el-form-item>
          </el-form>
        </el-tab-pane>
        <el-tab-pane :label="dialogLabel.quickbarTitle" name="second">
          <el-transfer style="display: flex; align-items: center; justify-content: center; margin-top: 2.5%"
            :render-content="renderFunc" filterable v-model="form.corporations"
            :titles="[dialogLabel.sourceList, dialogLabel.targetList]" :data="corporationData">
          </el-transfer>
        </el-tab-pane>
        <el-tab-pane :label="dialogLabel.desc.title" name="third">
          <el-collapse accordion>
            <el-collapse-item :title="dialogLabel.desc.dataDesc" name="1">
              <div style="word-break: keep-all;">
                <span style="font-weight: bold; margin-right: 10px;">{{ dialogLabel.materialPrice }}</span>
                {{ dialogLabel.desc.materialContent }}
                <br>
                <span style="font-weight: bold; margin-right: 10px;">{{ dialogLabel.productPrice }}</span>
                {{ dialogLabel.desc.productPriceContent }}
                <br>
                <span style="font-weight: bold; margin-right: 10px;">{{ dialogLabel.scope }}</span>
                {{ dialogLabel.desc.scopeContent }}
                <br>
                <span style="font-weight: bold; margin-right: 10px;">{{ dialogLabel.days }}</span>
                {{ dialogLabel.desc.daysContent }}
                <br>
                <span style="font-weight: bold; margin-right: 10px;">{{ dialogLabel.tax }}</span>
                {{ dialogLabel.desc.taxContent }}
              </div>
            </el-collapse-item>
            <el-collapse-item :title="dialogLabel.desc.tableDesc" name="2">       
              <div style="word-break: keep-all;">
                <span style="font-weight: bold; margin-right: 10px;">{{ tableLabel.cost }}</span>
                {{ dialogLabel.desc.costContent }}
                <br>
                <span style="font-weight: bold; margin-right: 10px;">{{ tableLabel.income }}</span>
                {{ dialogLabel.desc.incomeContent }}
                <br>
                <span style="font-weight: bold; margin-right: 10px;">{{ tableLabel.volume }}</span>
                {{ dialogLabel.desc.volumeContent }}
                <br>
                <span style="font-weight: bold; margin-right: 10px;">{{ tableLabel.saleIndex }}</span>
                {{ dialogLabel.desc.saleIndexContent }}
                <br>
                <span style="font-weight: bold; margin-right: 10px;">{{ tableLabel.unitProfit }}</span>
                {{ dialogLabel.desc.unitProfitContent }}
              </div>

            </el-collapse-item>
            <el-collapse-item :title="dialogLabel.desc.claim" name="3">
              <div style="word-break: keep-all;">
                {{ dialogLabel.desc.claimContent }}
              </div>
            </el-collapse-item>
          </el-collapse>
        </el-tab-pane>
      </el-tabs>
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
    corporationData: function () {
      return this.$store.state.corporationData;
    },
  },
  data() {
    return {
      dialogVisible: false,
      dialogLabel: this.$t("message.dialog"),
      tableLabel: this.$t("message.table"),
      activeName: 'first',
      renderFunc(h, option) {
          return <span title={option.label}>{option.label}</span>;
      },
    };
  },
  methods: {
    reset() {
      this.form.materialPrice = "sell";
      this.form.productPrice = "buy";
      this.form.days = "7";
      this.form.scope = "0.05";
      this.form.tax = 0;
      this.form.corporations = [];
    },
    closeDialog() {
      if (!localStorage.form) {
        this.$store.dispatch("setForm", this.form);
        localStorage.form = JSON.stringify(this.form);
        this.$emit("form-change");
      } else {
        const old = JSON.parse(localStorage.form);
        var is_same = (this.form.corporations.length == old.corporations.length) && this.form.corporations.every(el => old.corporations.includes(el));
        if (
          this.form.materialPrice != old.materialPrice ||
          this.form.productPrice != old.productPrice ||
          this.form.days != old.days ||
          this.form.tax != old.tax || !is_same
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
<style>
.v-modal {
  background-color: rgba(0, 0, 0, 0.9) !important;
}
</style>