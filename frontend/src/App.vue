<script setup lang="ts">
import { ref } from 'vue'
import { wol, shutdown } from './api'
import { debounce } from './utils/debounce'

const wolLoading = ref(false)
const shutdownLoading = ref(false)
const message = ref('')
const messageType = ref<'success' | 'error' | ''>('')

async function handleWOL() {
  if (wolLoading.value) return
  wolLoading.value = true
  message.value = ''
  try {
    const res = await wol()
    if (res.code === 0) {
      messageType.value = 'success'
      message.value = 'WOL 开机包已发送，请等待 NAS 启动...'
    } else {
      messageType.value = 'error'
      message.value = res.message || '操作失败'
    }
  } catch (e: any) {
    messageType.value = 'error'
    message.value = '网络错误：' + (e.message || String(e))
  } finally {
    // Keep button disabled for 3s as client-side anti-shake
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
    const res = await shutdown()
    if (res.code === 0) {
      messageType.value = 'success'
      message.value = '关机命令已发送，NAS 将安全关闭...'
    } else {
      messageType.value = 'error'
      message.value = res.message || '操作失败'
    }
  } catch (e: any) {
    messageType.value = 'error'
    message.value = '网络错误：' + (e.message || String(e))
  } finally {
    setTimeout(() => {
      shutdownLoading.value = false
    }, 3000)
  }
}

// Debounced versions (not really needed with button lock, but demonstrates the utility)
const debouncedWOL = debounce(handleWOL, 500)
const debouncedShutdown = debounce(handleShutdown, 500)
</script>

<template>
  <div class="container">
    <h1>🖥️ NAS 远程控制</h1>
    <p class="subtitle">Wake-on-LAN 网络唤醒 / SSH 远程关机</p>

    <div class="button-group">
      <button
        class="btn btn-wol"
        :disabled="wolLoading"
        @click="debouncedWOL"
      >
        <span v-if="wolLoading" class="spinner"></span>
        {{ wolLoading ? '发送中...' : '⚡ WOL 开机' }}
      </button>

      <button
        class="btn btn-shutdown"
        :disabled="shutdownLoading"
        @click="debouncedShutdown"
      >
        <span v-if="shutdownLoading" class="spinner"></span>
        {{ shutdownLoading ? '发送中...' : '🔴 关机' }}
      </button>
    </div>

    <div v-if="message" :class="['message', messageType]">
      {{ message }}
    </div>
  </div>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  background: #1a1a2e;
  color: #e0e0e0;
  min-height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
}

.container {
  text-align: center;
  padding: 2rem;
}

h1 {
  font-size: 2rem;
  margin-bottom: 0.5rem;
  color: #ffffff;
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
</style>
