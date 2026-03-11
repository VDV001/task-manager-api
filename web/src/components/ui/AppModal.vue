<script setup lang="ts">
import { X } from 'lucide-vue-next'

defineProps<{
  open: boolean
  title: string
}>()

defineEmits<{
  close: []
}>()
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="$emit('close')" />
        <div
          class="relative w-full max-w-lg rounded-2xl border border-border bg-bg-secondary shadow-2xl"
        >
          <div class="flex items-center justify-between p-5 border-b border-border">
            <h2 class="text-lg font-semibold text-text-primary">{{ title }}</h2>
            <button
              class="p-1.5 rounded-lg hover:bg-white/5 text-text-muted transition-colors cursor-pointer"
              @click="$emit('close')"
            >
              <X class="w-5 h-5" />
            </button>
          </div>
          <div class="p-5">
            <slot />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
</style>
