import Vue from "vue";

import VueI18n from "vue-i18n";
import de from "./lang/de";
import en from "./lang/en";
import fr from "./lang/fr";
import ja from "./lang/ja";
import ru from "./lang/ru";
import zh from "./lang/zh";

import locale from "element-ui/lib/locale";
import elementDe from "element-ui/lib/locale/lang/de"; 
import elementEn from "element-ui/lib/locale/lang/en"; 
import elementFr from "element-ui/lib/locale/lang/fr"; 
import elementJa from "element-ui/lib/locale/lang/ja"; 
import elementRu from "element-ui/lib/locale/lang/ru-RU"; 
import elementZh from "element-ui/lib/locale/lang/zh-CN"; 


Vue.use(VueI18n);

const messages = {
  de: { ...de, ...elementDe },
  en: { ...en, ...elementEn },
  fr: { ...fr, ...elementFr },
  ja: { ...ja, ...elementJa },
  ru: { ...ru, ...elementRu },
  zh: { ...zh, ...elementZh },
};

const i18n = new VueI18n({
  locale: "en",
  messages,
});

locale.i18n((key, value) => i18n.t(key, value));
export default i18n;
