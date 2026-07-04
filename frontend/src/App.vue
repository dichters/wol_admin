<template>
  <router-view />
</template>

<script setup lang="ts">
import { watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute } from 'vue-router'

const { locale } = useI18n()
const route = useRoute()

// Sync locale from URL path
watch(
  () => route.path,
  (path) => {
    if (path.includes('/en')) {
      locale.value = 'en'
    } else {
      locale.value = 'zh'
    }
  },
  { immediate: true }
)
</script>

<style>
body {
  margin: 0;
  background: #1a1a2e;
  color: #e0e0e0;
}
</style>
