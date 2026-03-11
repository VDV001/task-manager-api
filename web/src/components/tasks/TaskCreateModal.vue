<script setup lang="ts">
import { ref } from 'vue'
import { toast } from 'vue-sonner'
import { useTaskStore } from '@/stores/tasks'
import AppModal from '@/components/ui/AppModal.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppButton from '@/components/ui/AppButton.vue'

defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: [] }>()

const taskStore = useTaskStore()
const title = ref('')
const description = ref('')
const deadline = ref('')
const loading = ref(false)
const errors = ref<Record<string, string>>({})

async function handleSubmit() {
  errors.value = {}
  if (!title.value.trim()) {
    errors.value.title = 'Title is required'
    return
  }
  loading.value = true
  try {
    await taskStore.createTask({
      title: title.value,
      description: description.value || undefined,
      deadline: deadline.value ? new Date(deadline.value).toISOString() : undefined,
    })
    toast.success('Task created')
    title.value = ''
    description.value = ''
    deadline.value = ''
    emit('close')
  } catch {
    toast.error('Failed to create task')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AppModal :open="open" title="Create Task" @close="$emit('close')">
    <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
      <AppInput v-model="title" label="Title" placeholder="Task title" :error="errors.title" />
      <div class="flex flex-col gap-1.5">
        <label class="text-sm font-medium text-text-secondary">Description</label>
        <textarea
          v-model="description"
          rows="3"
          placeholder="Optional description"
          class="w-full px-4 py-2.5 rounded-xl bg-white/5 border border-border text-sm text-text-primary placeholder-text-muted resize-none focus:outline-none focus:ring-2 focus:ring-accent/50 focus:border-accent transition-all"
        />
      </div>
      <AppInput v-model="deadline" label="Deadline" type="datetime-local" />
      <div class="flex justify-end gap-3 mt-2">
        <AppButton variant="secondary" type="button" @click="$emit('close')">Cancel</AppButton>
        <AppButton type="submit" :loading="loading">Create</AppButton>
      </div>
    </form>
  </AppModal>
</template>
