<template>
  <div class="header">
    <el-dropdown>
      <el-button type="primary" size="small" icon="el-icon-menu">{{ headerLabel.menu }}</el-button>
      <el-dropdown-menu slot="dropdown">
        <router-link to="/">
          <el-dropdown-item>
            <i class="el-icon-s-home"></i>
            <span>{{ headerLabel.home }}</span>
          </el-dropdown-item>
        </router-link>
      </el-dropdown-menu>
    </el-dropdown>

    <el-divider direction="vertical"></el-divider>
    <span>EVE-LP</span>
    <el-divider direction="vertical"></el-divider>
    <span>{{ headerLabel.serverName }}</span>
    <el-divider direction="vertical"></el-divider>
    <span>{{ headerLabel.market }}</span>

    <el-dropdown
      @command="langChange"
      style="float: right;margin-right: 15px;margin-top:8px;cursor: pointer"
    >
      <span class="el-dropdown-link">
        {{ language }}
        <i class="el-icon-arrow-down el-icon--right"></i>
      </span>
      <el-dropdown-menu slot="dropdown">
        <el-dropdown-item command="de">Deutsch</el-dropdown-item>
        <el-dropdown-item command="en">English</el-dropdown-item>
        <el-dropdown-item command="fr">Français</el-dropdown-item>
        <el-dropdown-item command="ja">日本語</el-dropdown-item>
        <el-dropdown-item command="ru">Pусский</el-dropdown-item>
        <el-dropdown-item command="zh">中文</el-dropdown-item>
      </el-dropdown-menu>
    </el-dropdown>
  </div>
</template>

<script>
export default {
  mounted() {
    if (localStorage.lang == null) {
      localStorage.lang = "en";
    }
    if (localStorage.lang) {
      this.$i18n.locale = localStorage.lang;
    }
    this.language = this.langLabel(this.$i18n.locale);
  },
  data() {
    return {
      headerLabel: this.$t("message.header"),
      language: "",
    };
  },
  methods: {
    langChange(lang) {
      this.$i18n.locale = lang;
    },
    langLabel(lang) {
      switch (lang) {
        case "de":
          return "Deutsch";
        case "en":
          return "English";
        case "fr":
          return "Français";
        case "ja":
          return "日本語";
        case "ru":
          return "Pусский";
        case "zh":
          return "中文";
      }
    },
  },
  watch: {
    "$i18n.locale"() {
      localStorage.lang = this.$i18n.locale;
      this.headerLabel = this.$t("message.header");
      this.language = this.langLabel(this.$i18n.locale);
    },
  },
};
</script>

<style>
A:link {
  color: black;
  text-decoration: none;
}
.el-dropdown-link {
  cursor: pointer;
  color: #409eff;
}
.el-icon-arrow-down {
  font-size: 12px;
}
</style>