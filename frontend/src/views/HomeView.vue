<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { fetchVersion, fetchWOL, fetchShutdown } from '../api'
import { debounce } from '../utils/debounce'

const { t, locale } = useI18n()
const router = useRouter()

const wolLoading = ref(false)
const shutdownLoading = ref(false)
const message = ref('')
const messageType = ref<'success' | 'error' | ''>('')

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
    // silently ignore
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

async function handleWOL() {
  if (wolLoading.value) return
  wolLoading.value = true
  message.value = ''
  try {
    const res = await fetchWOL()
    if (res.code === 0) {
      messageType.value = 'success'
      message.value = t('wolSuccess')
    } else {
      messageType.value = 'error'
      message.value = res.message || t('error')
    }
  } catch (e: any) {
    messageType.value = 'error'
    message.value = t('error') + '：' + (e.message || String(e))
  } finally {
    setTimeout(() => {
      wolLoading.value = false
    }, 3000)
  }
}

async function handleShutdown() {
  if (shutdownLoading.value) return
  shutdownLoading.value = true
  message.value = ''
  try {
    const res = await fetchShutdown()
    if (res.code === 0) {
      messageType.value = 'success'
      message.value = t('shutdownSuccess')
    } else {
      messageType.value = 'error'
      message.value = res.message || t('error')
    }
  } catch (e: any) {
    messageType.value = 'error'
    message.value = t('error') + '：' + (e.message || String(e))
  } finally {
    setTimeout(() => {
      shutdownLoading.value = false
    }, 3000)
  }
}

const debouncedWOL = debounce(handleWOL, 500)
const debouncedShutdown = debounce(handleShutdown, 500)
</script>

<template>
  <div class="lang-switch" @click="toggleLang">
    <span :class="{ active: locale === 'zh' }">中文</span>
    <span class="divider">|</span>
    <span :class="{ active: locale === 'en' }">EN</span>
  </div>

  <div class="container">
    <h1>🖥️ {{ t('title') }}</h1>
    <p class="subtitle">{{ t('subtitle') }}</p>

    <div class="button-group">
      <button
        class="btn btn-wol"
        :disabled="wolLoading"
        @click="debouncedWOL"
      >
        <span v-if="wolLoading" class="spinner"></span>
        {{ wolLoading ? t('waking') : '⚡ ' + t('wakeUp') }}
      </button>

      <button
        class="btn btn-shutdown"
        :disabled="shutdownLoading"
        @click="debouncedShutdown"
      >
        <span v-if="shutdownLoading" class="spinner"></span>
        {{ shutdownLoading ? t('shutting') : '🔴 ' + t('shutdown') }}
      </button>
    </div>

    <div v-if="message" :class="['message', messageType]">
      {{ message }}
    </div>

    <div class="version-info" v-if="versionInfo">
      v{{ versionInfo.version }} · {{ versionInfo.arch }} · {{ versionInfo.build_time }}
    </div>
  </div>
</template>

<style scoped>
.container {
  text-align: center;
  padding: 2rem;
}

h1 {
  font-size: 2rem;
  margin-bottom: 0.5rem;
  color: #ffffff;
}

.lang-switch {
  position: fixed;
  top: 1rem;
  right: 1rem;
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

.subtitle {
  color: #888;
  margin-bottom: 2rem;
  font-size: 0.95rem;
}

.button-group {
  display: flex;
  gap: 1.5rem;
  justify-content: center;
  flex-wrap: wrap;
}

.btn {
  padding: 1rem 2.5rem;
  font-size: 1.1rem;
  font-weight: 600;
  border: none;
  border-radius: 12px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  transition: all 0.2s ease;
  min-width: 180px;
  justify-content: center;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.btn-wol {
  background: linear-gradient(135deg, #00b894, #00cec9);
  color: #fff;
}

.btn-wol:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(0, 184, 148, 0.4);
}

.btn-shutdown {
  background: linear-gradient(135deg, #d63031, #e17055);
  color: #fff;
}

.btn-shutdown:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(214, 48, 49, 0.4);
}

.message {
  margin-top: 1.5rem;
  padding: 0.8rem 1.5rem;
  border-radius: 8px;
  font-size: 0.95rem;
  max-width: 400px;
  margin-left: auto;
  margin-right: auto;
}

.message.success {
  background: rgba(0, 184, 148, 0.15);
  border: 1px solid rgba(0, 184, 148, 0.3);
  color: #00b894;
}

.message.error {
  background: rgba(214, 48, 49, 0.15);
  border: 1px solid rgba(214, 48, 49, 0.3);
  color: #d63031;
}

.spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.version-info {
  margin-top: 2rem;
  font-size: 0.75rem;
  color: #555;
}
</style>
