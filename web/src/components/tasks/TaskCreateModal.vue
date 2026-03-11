<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import { toast } from 'vue-sonner'
import { useTaskStore } from '@/stores/tasks'
import AppModal from '@/components/ui/AppModal.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppButton from '@/components/ui/AppButton.vue'

const { t } = useI18n()

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
    errors.value.title = t('validation.required', { field: t('task.title') })
    return
  }
  loading.value = true
  try {
    await taskStore.createTask({
      title: title.value,
      description: description.value || undefined,
      deadline: deadline.value ? new Date(deadline.value).toISOString() : undefined,
    })
    toast.success(t('task.taskCreated'))
    title.value = ''
    description.value = ''
    deadline.value = ''
    emit('close')
  } catch {
    toast.error(t('task.taskCreateFailed'))
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AppModal :open="open" :title="t('task.createTitle')" @close="$emit('close')">
    <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
      <AppInput
        v-model="title"
        :label="t('task.title')"
        :placeholder="t('task.titlePlaceholder')"
        :error="errors.title"
      />
      <div class="flex flex-col gap-1.5">
        <label class="text-sm font-medium text-text-secondary">{{ t('task.description') }}</label>
        <textarea
          v-model="description"
          rows="3"
          :placeholder="t('task.descriptionPlaceholder')"
          class="w-full px-4 py-2.5 rounded-xl bg-white/5 border border-border text-sm text-text-primary placeholder-text-muted resize-none focus:outline-none focus:ring-2 focus:ring-accent/50 focus:border-accent transition-all"
        />
      </div>
      <AppInput v-model="deadline" :label="t('task.deadline')" type="datetime-local" />
      <div class="flex justify-end gap-3 mt-2">
        <AppButton variant="secondary" type="button" @click="$emit('close')">
          {{ t('common.cancel') }}
        </AppButton>
        <AppButton type="submit" :loading="loading">{{ t('task.createButton') }}</AppButton>
      </div>
    </form>
  </AppModal>
</template>
