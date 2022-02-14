<template>
  <div class="home">
    <Header />
    <hr />
    <div style="display: flex">
      <div style="width: 100%; margin-left: 5px">
        <el-cascader
          v-model="corporation.value"
          style="width: 25%"
          v-loading="corporation.loading"
          :placeholder="corporation.label"
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
import Header from "@/components/Header.vue";

export default {
  name: "Home",
  components: {
    Header,
  },
  mounted() {
    if (localStorage.lang == null) {
      localStorage.lang = "en";
    }
    if (localStorage.lang) {
      this.$i18n.locale = localStorage.lang;
    }
  },
  data() {
    return {
      corporation: {
        value: "",
        loading: "",
        label: "",
        lists: [],
      },
    };
  },
  methods: {
    loadFactions() {
      this.selectLoading = true;
      this.axios.get(backendUrl + "factionList").then((response) => {
        const list = response.data;
        this.facList = list;
        this.initFacList(list);
        this.selectLoading = false;
      });
    },
  },
  watch: {
    "$i18n.locale"() {
      localStorage.lang = this.$i18n.locale;
    },
  },
};
</script>
