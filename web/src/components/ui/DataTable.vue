<script lang="ts">
export interface Column<T> {
  key: keyof T & string
  label: string
  sortable?: boolean
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  render?: (value: any, row: T) => string
  width?: string
}
</script>

<script setup lang="ts" generic="T extends Record<string, any>">
import { ref, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useVirtualizer } from '@tanstack/vue-virtual'
import { ArrowUp, ArrowDown, ArrowUpDown, Search } from 'lucide-vue-next'

const { t } = useI18n()

const props = withDefaults(
  defineProps<{
    columns: Column<T>[]
    data: T[]
    rowHeight?: number
    searchable?: boolean
  }>(),
  {
    rowHeight: 48,
    searchable: true,
  },
)

const emit = defineEmits<{
  rowClick: [row: T]
  'update:filteredData': [data: T[]]
}>()

// Search
const searchQuery = ref('')

// Sorting
const sortKey = ref<string>('')
const sortOrder = ref<'asc' | 'desc'>('asc')

function toggleSort(key: string) {
  if (sortKey.value === key) {
    sortOrder.value = sortOrder.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortKey.value = key
    sortOrder.value = 'asc'
  }
}

// Filtered + sorted data
const processedData = computed(() => {
  let result = [...props.data]

  // Global search
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter((row) =>
      props.columns.some((col) => {
        const val = row[col.key]
        return val != null && String(val).toLowerCase().includes(query)
      }),
    )
  }

  // Sort
  if (sortKey.value) {
    const key = sortKey.value
    const order = sortOrder.value === 'asc' ? 1 : -1
    result.sort((a, b) => {
      const aVal = a[key]
      const bVal = b[key]
      if (aVal == null) return 1
      if (bVal == null) return -1
      if (aVal < bVal) return -1 * order
      if (aVal > bVal) return 1 * order
      return 0
    })
  }

  return result
})

// Emit filtered data to parent for chart sync
watch(
  processedData,
  (data) => {
    emit('update:filteredData', data)
  },
  { immediate: true },
)

// Virtual scroll
const parentRef = ref<HTMLElement | null>(null)

const virtualizer = useVirtualizer({
  get count() {
    return processedData.value.length
  },
  getScrollElement: () => parentRef.value,
  estimateSize: () => props.rowHeight,
  overscan: 20,
})
</script>

<template>
  <div
    role="table"
    class="flex flex-col border border-border rounded-2xl bg-bg-card/60 backdrop-blur-sm overflow-hidden"
  >
    <!-- Search -->
    <div v-if="searchable" class="p-4 border-b border-border">
      <div class="relative">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-text-muted" />
        <input
          v-model="searchQuery"
          type="text"
          :aria-label="t('dashboard.searchInTable')"
          :placeholder="t('dashboard.searchInTable')"
          class="w-full pl-10 pr-4 py-2 rounded-xl bg-white/5 border border-border text-sm text-text-primary placeholder-text-muted focus:outline-none focus:ring-2 focus:ring-accent/50 focus:border-accent transition-all"
        />
      </div>
      <p class="text-xs text-text-muted mt-2">
        {{ t('dashboard.rows', { count: processedData.length.toLocaleString() }) }}
      </p>
    </div>

    <!-- Header -->
    <div class="flex bg-white/[0.03] border-b border-border">
      <div
        v-for="col in columns"
        :key="col.key"
        :style="{ width: col.width ?? '1fr', flex: col.width ? `0 0 ${col.width}` : '1' }"
        class="px-4 py-3 text-xs font-semibold text-text-secondary uppercase tracking-wider select-none"
        :class="
          col.sortable !== false ? 'cursor-pointer hover:text-text-primary transition-colors' : ''
        "
        @click="col.sortable !== false && toggleSort(col.key)"
      >
        <div class="flex items-center gap-1.5">
          <span>{{ col.label }}</span>
          <template v-if="col.sortable !== false">
            <ArrowUp
              v-if="sortKey === col.key && sortOrder === 'asc'"
              class="w-3.5 h-3.5 text-accent"
            />
            <ArrowDown
              v-else-if="sortKey === col.key && sortOrder === 'desc'"
              class="w-3.5 h-3.5 text-accent"
            />
            <ArrowUpDown v-else class="w-3.5 h-3.5 opacity-30" />
          </template>
        </div>
      </div>
    </div>

    <!-- Virtual scroll body -->
    <div ref="parentRef" class="overflow-auto" :style="{ maxHeight: '600px' }">
      <div
        :style="{ height: `${virtualizer.getTotalSize()}px`, width: '100%', position: 'relative' }"
      >
        <div
          v-for="virtualRow in virtualizer.getVirtualItems()"
          :key="virtualRow.index"
          class="flex absolute w-full border-b border-border/50 hover:bg-white/[0.03] transition-colors cursor-pointer"
          :style="{
            height: `${virtualRow.size}px`,
            transform: `translateY(${virtualRow.start}px)`,
          }"
          @click="emit('rowClick', processedData[virtualRow.index]!)"
        >
          <div
            v-for="col in columns"
            :key="col.key"
            :style="{ width: col.width ?? '1fr', flex: col.width ? `0 0 ${col.width}` : '1' }"
            class="px-4 flex items-center text-sm text-text-primary truncate"
          >
            {{
              col.render
                ? col.render(
                    processedData[virtualRow.index]![col.key],
                    processedData[virtualRow.index]!,
                  )
                : processedData[virtualRow.index]![col.key]
            }}
          </div>
        </div>
      </div>
    </div>
  </div>
</template>
