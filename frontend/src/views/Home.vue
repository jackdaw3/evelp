<template>
  <div class="home">
    <Header />
    <div style="margin-top: -13px">
      <el-divider></el-divider>
    </div>
    <div style="display: flex; margin-top: -8px">
      <div>
        <Dialog @form-change="formChange" />
      </div>
      <div style="width: 100%; margin-left: 5px">
        <el-cascader
          v-model="corporation.value"
          style="width: 25%"
          v-loading="corporation.loading"
          :placeholder="corporation.placeholder"
          :options="corporation.lists"
          element-loading-background="rgba(0, 0, 0, 0.8)"
          clearable
          filterable
          @change="loadTable"
          ref="cascader"
        ></el-cascader>
      </div>
      <el-button
        icon="el-icon-download"
        circle
        size="medium"
        @click="exportExcel"
        style="height: 50%; float: right; margin-right: 15px; cursor: pointer"
      ></el-button>
    </div>
    <Table ref="Table" style="margin-top: -10px" />
  </div>
</template>

<script>
import Dialog from "@/components/Dialog.vue";
import Header from "@/components/Header.vue";
import Table from "@/components/Table.vue";

const backend = "https://eve-lp.com/api/";
const the_forge = "10000002";

export default {
  name: "Home",
  components: {
    Dialog,
    Header,
    Table,
  },
  created() {
    if (localStorage.lang == null) {
      localStorage.lang = "en";
    }
    if (localStorage.lang) {
      this.$i18n.locale = localStorage.lang;
    }
    if (localStorage.form) {
      this.$store.dispatch("setForm", JSON.parse(localStorage.form));
    }
    this.loadFactions();
  },
  data() {
    return {
      corporation: {
        value: "",
        loading: "",
        placeholder: this.$t("message.corporation.placeholder"),
        lists: [],
      },
    };
  },
  computed: {
    form: function () {
      return this.$store.state.form;
    },
  },
  methods: {
    loadFactions() {
      this.corporation.loading = true;
      this.axios
        .get(backend + "faction", {
          params: {
            lang: this.$i18n.locale,
          },
        })
        .then((response) => {
          var factions = response.data;
          this.loadCorporations(factions);
          this.corporation.loading = false;
        });
    },
    loadCorporations(list) {
      var factions = new Array();
      for (let i = 0; i < list.length; ++i) {
        var faction = new Object();
        faction.value = list[i].FactionId;
        faction.label = list[i].FactionName;
        if (list.length > 0) {
          var corporations = new Array();
          var corporationList = list[i].Corporations;
          for (let j = 0; j < corporationList.length; ++j) {
            var corporation = new Object();
            corporation.value = corporationList[j].CorporationId;
            corporation.label = corporationList[j].CorporationName;
            corporations.push(corporation);
          }
        }
        faction.children = corporations;
        factions.push(faction);
      }
      this.corporation.lists = factions;
    },
    loadTable(value) {
      if (value == "" || value == null) {
        this.$store.dispatch("setTableData", []);
        return;
      }
      const corporationId = parseInt(value[1]);
      const corporationName =
        this.$refs["cascader"].getCheckedNodes()[0].pathLabels[1];
      this.$store.dispatch("setCorporationId", corporationId);
      this.$store.dispatch("setCorporationName", corporationName);
      this.corporation.loading = true;
      this.axios
        .get(backend + "offer", {
          params: {
            regionId: the_forge,
            scope: this.form.scope,
            corporationId: corporationId,
            lang: this.$i18n.locale,
            days: this.form.days,
            tax: this.form.tax,
            productPrice: this.form.productPrice,
            materialPrice: this.form.materialPrice,
          },
        })
        .then((response) => {
          var data = response.data;
          for (let i = 0; i < data.length; ++i) {
            var matertials = data[i].Matertials;
            if (matertials == null) {
              continue;
            }
            let count = 0;
            for (let j = 0; j + count < matertials.length; ) {
              if (count == 0) {
                matertials[j].length = 1;
                ++count;
                continue;
              }
              if (
                matertials[j].IsBluePrint == matertials[j + count].IsBluePrint
              ) {
                matertials[j].length += 1;
                matertials[j + count].length = 0;
                ++count;
              } else {
                j += count;
                count = 0;
              }
            }
          }
          this.$store.dispatch("setTableData", data);
          this.corporation.loading = false;
        })
        .catch(() => {
          this.$store.dispatch("setTableData", []);
          this.corporation.loading = false;
        });
    },
    formChange() {
      this.reloadTable();
    },
    exportExcel() {
      this.$refs.Table.exportExcel();
    },
    reloadTable() {
      var cascaderValue = this.$refs["cascader"].getCheckedNodes()[0];
      if (cascaderValue != null) {
        var facAndcorp = new Array();
        facAndcorp[0] = cascaderValue.parent.value;
        facAndcorp[1] = cascaderValue.value;
        this.loadTable(facAndcorp);
      }
    },
  },
  watch: {
    "$i18n.locale"() {
      this.corporation.placeholder = this.$t("message.corporation.placeholder");
      this.loadFactions(this.facList);
      this.reloadTable();
    },
  },
};
</script>

<style>
hr {
  display: block;
  height: 1px;
  border: 0;
  border-top: 0px solid #ccc;
  margin: 0.618em 0;
  padding: 0;
}
</style>
