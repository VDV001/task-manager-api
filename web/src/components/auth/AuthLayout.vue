<script setup lang="ts">
import { useI18n } from 'vue-i18n'
import { CheckSquare, Languages } from 'lucide-vue-next'
import { setLocale, locale } from '@/i18n'

const { t } = useI18n()

function toggleLocale() {
  setLocale(locale.value === 'en' ? 'ru' : 'en')
}
</script>

<template>
  <div class="min-h-screen flex flex-col items-center justify-center p-4 relative overflow-hidden">
    <!-- Language switcher -->
    <div class="fixed top-6 right-6 z-50">
      <button
        class="flex items-center gap-1.5 px-3 py-1.5 rounded-full text-xs font-medium text-text-secondary hover:text-text-primary hover:bg-white/[0.06] transition-all duration-200 cursor-pointer border border-white/[0.08]"
        @click="toggleLocale"
      >
        <Languages class="w-3.5 h-3.5" />
        <span class="uppercase">{{ locale === 'en' ? 'RU' : 'EN' }}</span>
      </button>
    </div>

    <!-- Background effects -->
    <div
      class="absolute top-[-20%] left-[-10%] w-[500px] h-[500px] rounded-full bg-accent/[0.06] blur-[100px] pointer-events-none"
    />
    <div
      class="absolute bottom-[-20%] right-[-10%] w-[400px] h-[400px] rounded-full bg-info/[0.04] blur-[100px] pointer-events-none"
    />

    <!-- Logo -->
    <div class="flex items-center gap-2.5 mb-8 relative z-10">
      <CheckSquare class="w-7 h-7 text-accent" />
      <span class="text-xl font-bold text-text-primary tracking-tight">TaskFlow</span>
    </div>

    <!-- Auth Card -->
    <div class="w-full max-w-md relative z-10">
      <div class="bg-bg-secondary border border-white/[0.08] rounded-2xl shadow-xl p-8">
        <slot />
      </div>
    </div>

    <!-- Footer -->
    <p class="mt-8 text-center text-xs text-text-muted relative z-10">
      {{ t('common.footer') }}
    </p>
  </div>
</template>
