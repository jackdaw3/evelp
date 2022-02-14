<template>
  <div class="home">
    <Header />
    <hr />
    <div style="display: flex">
      <div>
        <Dialog />
      </div>
      <div style="width: 100%; margin-left: 5px">
        <el-cascader
          v-model="corporation.value"
          style="width: 25%"
          v-loading="corporation.loading"
          :placeholder="corporation.placeholder"
          :options="corporation.lists"
          clearable
          filterable
          @change="corporationChange"
        >
        </el-cascader>
      </div>
    </div>
  </div>
</template>

<script>
import Dialog from "@/components/Dialog.vue";
import Header from "@/components/Header.vue";

const backend = "http://localhost:9000/";

export default {
  name: "Home",
  components: {
    Dialog,
    Header,
  },
  mounted() {
    if (localStorage.lang == null) {
      localStorage.lang = "en";
    }
    if (localStorage.lang) {
      this.$i18n.locale = localStorage.lang;
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
    corporationChange() {},
  },
  watch: {
    "$i18n.locale"() {
      this.corporation.placeholder = this.$t("message.corporation.placeholder");
      this.loadFactions(this.facList);
      localStorage.lang = this.$i18n.locale;
    },
  },
};
</script>
