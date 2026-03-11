<script setup lang="ts">
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { useTaskStore } from '@/stores/tasks'
import AppModal from '@/components/ui/AppModal.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppButton from '@/components/ui/AppButton.vue'
import type { Task, TaskStatus } from '@/types/task'

const { t } = useI18n()

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

watch(
  () => props.task,
  (task) => {
    if (task) {
      title.value = task.title
      description.value = task.description
      status.value = task.status
      deadline.value = task.deadline ? task.deadline.slice(0, 16) : ''
    }
  },
)

const statuses: { value: TaskStatus; labelKey: string }[] = [
  { value: 'new', labelKey: 'task.new' },
  { value: 'in_progress', labelKey: 'task.inProgress' },
  { value: 'done', labelKey: 'task.done' },
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
    toast.success(t('task.taskUpdated'))
    emit('close')
  } catch {
    toast.error(t('task.taskUpdateFailed'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AppModal :open="open" :title="t('task.editTitle')" @close="$emit('close')">
    <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
      <AppInput v-model="title" :label="t('task.title')" />
      <div class="flex flex-col gap-1.5">
        <label class="text-sm font-medium text-text-secondary">{{ t('task.description') }}</label>
        <textarea
          v-model="description"
          rows="3"
          class="w-full px-4 py-2.5 rounded-xl bg-white/5 border border-border text-sm text-text-primary placeholder-text-muted resize-none focus:outline-none focus:ring-2 focus:ring-accent/50 focus:border-accent transition-all"
        />
      </div>
      <div class="flex flex-col gap-1.5">
        <label class="text-sm font-medium text-text-secondary">{{ t('task.status') }}</label>
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
            {{ t(s.labelKey) }}
          </button>
        </div>
      </div>
      <AppInput v-model="deadline" :label="t('task.deadline')" type="datetime-local" />
      <div class="flex justify-end gap-3 mt-2">
        <AppButton variant="secondary" type="button" @click="$emit('close')">
          {{ t('common.cancel') }}
        </AppButton>
        <AppButton type="submit" :loading="loading">{{ t('task.saveButton') }}</AppButton>
      </div>
    </form>
  </AppModal>
</template>
