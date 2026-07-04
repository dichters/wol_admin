<template>
  <div class="container">
    <div class="header">
      <h1>{{ t('title') }}</h1>
      <div class="lang-switch" @click="toggleLang">
        <span :class="{ active: locale === 'zh' }">中文</span>
        <span class="divider">|</span>
        <span :class="{ active: locale === 'en' }">EN</span>
      </div>
    </div>

    <div class="cards">
      <div class="card">
        <button class="btn btn-wake" :disabled="waking" @click="doWOL">
          {{ waking ? t('waking') : t('wakeUp') }}
        </button>
      </div>
      <div class="card">
        <button class="btn btn-shutdown" :disabled="shutting" @click="doShutdown">
          {{ shutting ? t('shutting') : t('shutdown') }}
        </button>
      </div>
    </div>

    <div class="version-info" v-if="versionInfo">
      <span>{{ t('version') }}: {{ versionInfo.version }}</span>
      <span>{{ t('arch') }}: {{ versionInfo.arch }}</span>
      <span>{{ t('buildTime') }}: {{ versionInfo.build_time }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { fetchVersion, fetchWOL, fetchShutdown } from '../api'

const { t, locale } = useI18n()
const router = useRouter()

const waking = ref(false)
const shutting = ref(false)

interface VersionInfo {
  version: string
  arch: string
  build_time: string
}

const versionInfo = ref<VersionInfo | null>(null)

onMounted(async () => {
  try {
    const data = await fetchVersion()
    versionInfo.value = data
  } catch {
    // silently ignore version fetch errors
  }
})

function toggleLang() {
  const newLocale = locale.value === 'zh' ? 'en' : 'zh'
  locale.value = newLocale
  if (newLocale === 'en') {
    router.push('/en')
  } else {
    router.push('/')
  }
}

async function doWOL() {
  waking.value = true
  try {
    const res = await fetchWOL()
    alert(res.code === 0 ? t('wolSuccess') : t('error') + ': ' + res.message)
  } catch {
    alert(t('error'))
  } finally {
    waking.value = false
  }
}

async function doShutdown() {
  if (!confirm(t('confirmShutdown'))) return
  shutting.value = true
  try {
    const res = await fetchShutdown()
    alert(res.code === 0 ? t('shutdownSuccess') : t('error') + ': ' + res.message)
  } catch {
    alert(t('error'))
  } finally {
    shutting.value = false
  }
}
</script>

<style scoped>
.container {
  max-width: 480px;
  margin: 0 auto;
  padding: 24px 16px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
}

.header h1 {
  font-size: 1.5rem;
  margin: 0;
  color: #e0e0e0;
}

.lang-switch {
  cursor: pointer;
  user-select: none;
  font-size: 0.9rem;
  color: #888;
}

.lang-switch .active {
  color: #fff;
  font-weight: bold;
}

.lang-switch .divider {
  margin: 0 4px;
  color: #555;
}

.lang-switch:hover {
  color: #ccc;
}

.cards {
  display: flex;
  gap: 16px;
  margin-bottom: 32px;
}

.card {
  flex: 1;
}

.btn {
  width: 100%;
  padding: 16px 0;
  border: none;
  border-radius: 8px;
  font-size: 1.1rem;
  font-weight: 600;
  cursor: pointer;
  transition: opacity 0.2s;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-wake {
  background: #4caf50;
  color: #fff;
}

.btn-wake:hover:not(:disabled) {
  background: #43a047;
}

.btn-shutdown {
  background: #f44336;
  color: #fff;
}

.btn-shutdown:hover:not(:disabled) {
  background: #e53935;
}

.version-info {
  display: flex;
  gap: 16px;
  font-size: 0.8rem;
  color: #666;
  flex-wrap: wrap;
}
</style>
