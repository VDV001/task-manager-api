<script setup lang="ts">
import { computed } from 'vue'
import { cn } from '@/lib/utils'

let counter = 0

const props = defineProps<{
  id?: string
  label?: string
  error?: string
  modelValue?: string
  type?: string
  placeholder?: string
}>()

defineEmits<{
  'update:modelValue': [value: string]
}>()

const autoId = `app-input-${++counter}`
const inputId = computed(() => props.id || autoId)
const errorId = computed(() => `${inputId.value}-error`)
</script>

<template>
  <div class="flex flex-col gap-1.5">
    <label v-if="label" :for="inputId" class="text-sm font-medium text-text-secondary">
      {{ label }}
    </label>
    <input
      :id="inputId"
      :type="type ?? 'text'"
      :value="modelValue"
      :placeholder="placeholder"
      :aria-invalid="error ? true : undefined"
      :aria-describedby="error ? errorId : undefined"
      :class="
        cn(
          'w-full h-11 px-4 py-2.5 rounded-xl bg-white/[0.06] border text-text-primary placeholder-text-muted text-sm transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-accent/50',
          error ? 'border-danger focus:ring-danger/50' : 'border-white/[0.06] focus:border-accent',
        )
      "
      @input="$emit('update:modelValue', ($event.target as HTMLInputElement).value)"
    />
    <p v-if="error" :id="errorId" class="text-xs text-danger" role="alert">{{ error }}</p>
  </div>
</template>
