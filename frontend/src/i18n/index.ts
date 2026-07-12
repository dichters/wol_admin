import { createI18n } from 'vue-i18n'
import zh from './zh'
import en from './en'

const i18n = createI18n({
  legacy: false,
  locale: detectLocale(),
  fallbackLocale: 'zh',
  messages: { zh, en },
})

function detectLocale(): string {
  // Check URL path first: /wol/en means English
  const path = window.location.pathname
  if (path.includes('/en')) {
    return 'en'
  }
  return 'zh'
}

export default i18n
