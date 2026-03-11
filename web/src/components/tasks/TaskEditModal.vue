<script setup lang="ts">
import { ref, watch } from 'vue'
import { toast } from 'vue-sonner'
import { useTaskStore } from '@/stores/tasks'
import AppModal from '@/components/ui/AppModal.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppButton from '@/components/ui/AppButton.vue'
import type { Task, TaskStatus } from '@/types/task'

const props = defineProps<{
  open: boolean
  task: Task | null
}>()
const emit = defineEmits<{ close: [] }>()

const taskStore = useTaskStore()
const title = ref('')
const description = ref('')
const status = ref<TaskStatus>('new')
const deadline = ref('')
const loading = ref(false)

watch(() => props.task, (t) => {
  if (t) {
    title.value = t.title
    description.value = t.description
    status.value = t.status
    deadline.value = t.deadline ? t.deadline.slice(0, 16) : ''
  }
})

const statuses: { value: TaskStatus; label: string }[] = [
  { value: 'new', label: 'New' },
  { value: 'in_progress', label: 'In Progress' },
  { value: 'done', label: 'Done' },
]

async function handleSubmit() {
  if (!props.task) return
  loading.value = true
  try {
    await taskStore.updateTask(props.task.id, {
      title: title.value,
      description: description.value,
      status: status.value,
      deadline: deadline.value ? new Date(deadline.value).toISOString() : undefined,
    })
    toast.success('Task updated')
    emit('close')
  } catch {
    toast.error('Failed to update task')
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AppModal :open="open" title="Edit Task" @close="$emit('close')">
    <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
      <AppInput v-model="title" label="Title" />
      <div class="flex flex-col gap-1.5">
        <label class="text-sm font-medium text-text-secondary">Description</label>
        <textarea
          v-model="description"
          rows="3"
          class="w-full px-4 py-2.5 rounded-xl bg-white/5 border border-border text-sm text-text-primary placeholder-text-muted resize-none focus:outline-none focus:ring-2 focus:ring-accent/50 focus:border-accent transition-all"
        />
      </div>
      <div class="flex flex-col gap-1.5">
        <label class="text-sm font-medium text-text-secondary">Status</label>
        <div class="flex gap-2">
          <button
            v-for="s in statuses"
            :key="s.value"
            type="button"
            :class="[
              'px-3 py-2 rounded-xl text-xs font-medium border transition-all cursor-pointer',
              status === s.value
                ? 'bg-accent/15 text-accent border-accent/30'
                : 'bg-white/5 text-text-secondary border-border hover:bg-white/10',
            ]"
            @click="status = s.value"
          >
            {{ s.label }}
          </button>
        </div>
      </div>
      <AppInput v-model="deadline" label="Deadline" type="datetime-local" />
      <div class="flex justify-end gap-3 mt-2">
        <AppButton variant="secondary" type="button" @click="$emit('close')">Cancel</AppButton>
        <AppButton type="submit" :loading="loading">Save</AppButton>
      </div>
    </form>
  </AppModal>
</template>
