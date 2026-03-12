<script setup lang="ts">
import { watch, onUnmounted } from 'vue'
import { X } from 'lucide-vue-next'

let counter = 0
const modalTitleId = `modal-title-${++counter}`

const props = defineProps<{
  open: boolean
  title: string
}>()

const emit = defineEmits<{
  close: []
}>()

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape') emit('close')
}

watch(
  () => props.open,
  (open) => {
    if (open) {
      document.addEventListener('keydown', onKeydown)
    } else {
      document.removeEventListener('keydown', onKeydown)
    }
  },
)

onUnmounted(() => document.removeEventListener('keydown', onKeydown))
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="open"
        role="dialog"
        aria-modal="true"
        :aria-labelledby="modalTitleId"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
      >
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm" @click="$emit('close')" />
        <div
          class="relative w-full max-w-lg rounded-2xl border border-border bg-bg-secondary shadow-2xl"
        >
          <div class="flex items-center justify-between p-5 border-b border-border">
            <h2 :id="modalTitleId" class="text-lg font-semibold text-text-primary">{{ title }}</h2>
            <button
              aria-label="Close"
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
