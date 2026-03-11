import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import { createI18n } from 'vue-i18n'
import AppBadge from '@/components/ui/AppBadge.vue'
import en from '@/i18n/en'
import ru from '@/i18n/ru'

import type { TaskStatus } from '@/types/task'

function createWrapper(status: TaskStatus) {
  const i18n = createI18n({
    legacy: false,
    locale: 'en',
    messages: { en, ru },
  })

  return mount(AppBadge, {
    props: { status },
    global: {
      plugins: [i18n],
    },
  })
}

describe('AppBadge', () => {
  it('renders "New" for new status', () => {
    const wrapper = createWrapper('new')
    expect(wrapper.text()).toBe('New')
  })

  it('renders "In Progress" for in_progress status', () => {
    const wrapper = createWrapper('in_progress')
    expect(wrapper.text()).toBe('In Progress')
  })

  it('renders "Done" for done status', () => {
    const wrapper = createWrapper('done')
    expect(wrapper.text()).toBe('Done')
  })

  it('applies info styling for new status', () => {
    const wrapper = createWrapper('new')
    expect(wrapper.html()).toContain('bg-info')
  })

  it('applies warning styling for in_progress status', () => {
    const wrapper = createWrapper('in_progress')
    expect(wrapper.html()).toContain('bg-warning')
  })

  it('applies success styling for done status', () => {
    const wrapper = createWrapper('done')
    expect(wrapper.html()).toContain('bg-success')
  })

  it('renders in Russian when locale is ru', () => {
    const i18n = createI18n({
      legacy: false,
      locale: 'ru',
      messages: { en, ru },
    })

    const wrapper = mount(AppBadge, {
      props: { status: 'done' as const },
      global: { plugins: [i18n] },
    })

    expect(wrapper.text()).toBe('Готово')
  })
})
