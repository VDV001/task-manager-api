<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import { Plus } from 'lucide-vue-next'
import { useTaskStore } from '@/stores/tasks'
import AppButton from '@/components/ui/AppButton.vue'
import AppSkeleton from '@/components/ui/AppSkeleton.vue'
import TaskCard from '@/components/tasks/TaskCard.vue'
import TaskFilters from '@/components/tasks/TaskFilters.vue'
import TaskCreateModal from '@/components/tasks/TaskCreateModal.vue'
import TaskEditModal from '@/components/tasks/TaskEditModal.vue'
import TaskPagination from '@/components/tasks/TaskPagination.vue'
import type { Task, TaskFilter } from '@/types/task'

const taskStore = useTaskStore()

const showCreate = ref(false)
const showEdit = ref(false)
const editingTask = ref<Task | null>(null)

onMounted(() => {
  taskStore.fetchTasks()
  taskStore.fetchStats()
})

function openEdit(task: Task) {
  editingTask.value = task
  showEdit.value = true
}

async function handleDelete(task: Task) {
  if (!confirm(`Delete "${task.title}"?`)) return
  try {
    await taskStore.deleteTask(task.id)
    toast.success('Task deleted')
  } catch {
    toast.error('Failed to delete task')
  }
}

function handleFilterUpdate(update: Partial<TaskFilter>) {
  taskStore.setFilter(update)
}
</script>

<template>
  <div>
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-2xl font-bold text-text-primary">Tasks</h1>
        <p class="text-sm text-text-secondary mt-1">{{ taskStore.meta?.total ?? 0 }} tasks total</p>
      </div>
      <AppButton @click="showCreate = true">
        <Plus class="w-4 h-4" />
        New Task
      </AppButton>
    </div>

    <!-- Filters -->
    <TaskFilters :filter="taskStore.filter" class="mb-6" @update="handleFilterUpdate" />

    <!-- Task grid -->
    <div v-if="taskStore.loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="i in 6"
        :key="i"
        class="rounded-2xl border border-border bg-bg-card/60 p-5 space-y-3"
      >
        <AppSkeleton class="w-20" />
        <AppSkeleton class="w-3/4 h-5" />
        <AppSkeleton class="w-full" />
        <AppSkeleton class="w-1/3" />
      </div>
    </div>

    <div
      v-else-if="taskStore.tasks.length"
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
    >
      <TaskCard
        v-for="task in taskStore.tasks"
        :key="task.id"
        :task="task"
        @edit="openEdit"
        @delete="handleDelete"
      />
    </div>

    <div v-else class="text-center py-20">
      <p class="text-text-muted text-lg mb-4">No tasks yet</p>
      <AppButton @click="showCreate = true">
        <Plus class="w-4 h-4" />
        Create your first task
      </AppButton>
    </div>

    <!-- Pagination -->
    <TaskPagination
      :page="taskStore.filter.page ?? 1"
      :total-pages="taskStore.totalPages"
      @update:page="taskStore.setPage"
    />

    <!-- Modals -->
    <TaskCreateModal :open="showCreate" @close="showCreate = false" />
    <TaskEditModal :open="showEdit" :task="editingTask" @close="showEdit = false" />
  </div>
</template>
