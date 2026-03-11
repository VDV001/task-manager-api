<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { toast } from 'vue-sonner'
import { Plus, LayoutGrid, Table2, Database } from 'lucide-vue-next'
import { useTaskStore } from '@/stores/tasks'
import AppButton from '@/components/ui/AppButton.vue'
import AppCard from '@/components/ui/AppCard.vue'
import AppSkeleton from '@/components/ui/AppSkeleton.vue'
import TaskCard from '@/components/tasks/TaskCard.vue'
import TaskFilters from '@/components/tasks/TaskFilters.vue'
import TaskCreateModal from '@/components/tasks/TaskCreateModal.vue'
import TaskEditModal from '@/components/tasks/TaskEditModal.vue'
import TaskPagination from '@/components/tasks/TaskPagination.vue'
import DataTable from '@/components/ui/DataTable.vue'
import TasksLineChart from '@/components/charts/TasksLineChart.vue'
import TasksPieChart from '@/components/charts/TasksPieChart.vue'
import { generateMockTasks } from '@/lib/mock-data'
import type { Task, TaskFilter } from '@/types/task'
import { formatDistanceToNow, parseISO } from 'date-fns'

const taskStore = useTaskStore()

const showCreate = ref(false)
const showEdit = ref(false)
const editingTask = ref<Task | null>(null)
const viewMode = ref<'cards' | 'table'>('cards')
const demoMode = ref(false)
const mockTasks = ref<Task[]>([])

// Table columns definition
const tableColumns = [
  { key: 'title' as const, label: 'Title', width: '30%' },
  {
    key: 'status' as const,
    label: 'Status',
    width: '120px',
    render: (val: string) => {
      const map: Record<string, string> = {
        new: '● New',
        in_progress: '◐ In Progress',
        done: '✓ Done',
      }
      return map[val] ?? val
    },
  },
  { key: 'description' as const, label: 'Description' },
  {
    key: 'deadline' as const,
    label: 'Deadline',
    width: '160px',
    render: (val: string) => {
      if (!val) return '—'
      return formatDistanceToNow(parseISO(val), { addSuffix: true })
    },
  },
  {
    key: 'created_at' as const,
    label: 'Created',
    width: '160px',
    render: (val: string) => formatDistanceToNow(parseISO(val), { addSuffix: true }),
  },
]

// Data for charts — use mock data in demo mode, real data otherwise
const chartTasks = computed(() => (demoMode.value ? mockTasks.value : taskStore.tasks))

onMounted(() => {
  taskStore.fetchTasks()
  taskStore.fetchStats()
})

function toggleDemo() {
  demoMode.value = !demoMode.value
  if (demoMode.value && mockTasks.value.length === 0) {
    mockTasks.value = generateMockTasks(10000)
  }
}

function openEdit(task: Task) {
  editingTask.value = task
  showEdit.value = true
}

function handleTableRowClick(row: Task) {
  if (!demoMode.value) {
    openEdit(row)
  }
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

const displayTasks = computed(() => (demoMode.value ? mockTasks.value : taskStore.tasks))
</script>

<template>
  <div>
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <div>
        <h1 class="text-2xl font-bold text-text-primary">Tasks</h1>
        <p class="text-sm text-text-secondary mt-1">
          {{
            demoMode
              ? `${mockTasks.length.toLocaleString()} demo rows`
              : `${taskStore.meta?.total ?? 0} tasks total`
          }}
        </p>
      </div>
      <div class="flex items-center gap-2">
        <!-- Demo mode toggle -->
        <AppButton :variant="demoMode ? 'primary' : 'secondary'" size="sm" @click="toggleDemo">
          <Database class="w-4 h-4" />
          {{ demoMode ? '10K Demo' : 'Demo' }}
        </AppButton>
        <!-- View toggle -->
        <div class="flex bg-white/5 rounded-xl p-1 border border-border">
          <button
            :class="[
              'p-2 rounded-lg transition-all cursor-pointer',
              viewMode === 'cards'
                ? 'bg-accent/15 text-accent'
                : 'text-text-secondary hover:text-text-primary',
            ]"
            @click="viewMode = 'cards'"
          >
            <LayoutGrid class="w-4 h-4" />
          </button>
          <button
            :class="[
              'p-2 rounded-lg transition-all cursor-pointer',
              viewMode === 'table'
                ? 'bg-accent/15 text-accent'
                : 'text-text-secondary hover:text-text-primary',
            ]"
            @click="viewMode = 'table'"
          >
            <Table2 class="w-4 h-4" />
          </button>
        </div>
        <AppButton v-if="!demoMode" @click="showCreate = true">
          <Plus class="w-4 h-4" />
          New Task
        </AppButton>
      </div>
    </div>

    <!-- Charts -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-5 mb-8">
      <AppCard>
        <h3 class="text-sm font-semibold text-text-secondary mb-4">Tasks Created Over Time</h3>
        <TasksLineChart :tasks="chartTasks" />
      </AppCard>
      <AppCard>
        <h3 class="text-sm font-semibold text-text-secondary mb-4">Tasks by Status</h3>
        <TasksPieChart :stats="taskStore.stats" />
      </AppCard>
    </div>

    <!-- Filters (non-demo only) -->
    <TaskFilters
      v-if="!demoMode"
      :filter="taskStore.filter"
      class="mb-6"
      @update="handleFilterUpdate"
    />

    <!-- Table view -->
    <template v-if="viewMode === 'table'">
      <DataTable
        :columns="tableColumns"
        :data="displayTasks"
        :row-height="48"
        @row-click="handleTableRowClick"
      />
    </template>

    <!-- Card view -->
    <template v-else>
      <div
        v-if="taskStore.loading && !demoMode"
        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
      >
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
        v-else-if="displayTasks.length"
        class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4"
      >
        <TaskCard
          v-for="task in demoMode ? displayTasks.slice(0, 12) : displayTasks"
          :key="task.id"
          :task="task"
          @edit="openEdit"
          @delete="handleDelete"
        />
        <p v-if="demoMode" class="col-span-full text-center text-text-muted text-sm py-4">
          Showing 12 of {{ mockTasks.length.toLocaleString() }} — switch to table view for virtual
          scroll
        </p>
      </div>

      <div v-else class="text-center py-20">
        <p class="text-text-muted text-lg mb-4">No tasks yet</p>
        <AppButton @click="showCreate = true">
          <Plus class="w-4 h-4" />
          Create your first task
        </AppButton>
      </div>

      <!-- Pagination (non-demo only) -->
      <TaskPagination
        v-if="!demoMode"
        :page="taskStore.filter.page ?? 1"
        :total-pages="taskStore.totalPages"
        @update:page="taskStore.setPage"
      />
    </template>

    <!-- Modals -->
    <TaskCreateModal :open="showCreate" @close="showCreate = false" />
    <TaskEditModal :open="showEdit" :task="editingTask" @close="showEdit = false" />
  </div>
</template>
