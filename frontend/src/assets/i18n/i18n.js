import Vue from 'vue'

import VueI18n from 'vue-i18n'
import de from './lang/de'
import en from './lang/en'
import fr from './lang/fr'
import ja from './lang/ja'
import ru from './lang/ru'
import zh from './lang/zh'
Vue.use(VueI18n)

const messages = {
    de: de,
    en: en,
    fr: fr,
    ja: ja,
    ru: ru,
    zh: zh
}

const i18n = new VueI18n({
    locale: 'en', 
    messages
})


export default i18n