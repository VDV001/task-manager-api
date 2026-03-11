import { describe, it, expect } from 'vitest'
import { createI18n } from 'vue-i18n'
import en from '@/i18n/en'
import ru from '@/i18n/ru'

function createTestI18n(locale: 'en' | 'ru' = 'en') {
  return createI18n({
    legacy: false,
    locale,
    fallbackLocale: 'en',
    messages: { en, ru },
  })
}

describe('i18n', () => {
  it('has matching keys between en and ru', () => {
    function getKeys(obj: Record<string, unknown>, prefix = ''): string[] {
      return Object.entries(obj).flatMap(([key, value]) => {
        const fullKey = prefix ? `${prefix}.${key}` : key
        if (typeof value === 'object' && value !== null) {
          return getKeys(value as Record<string, unknown>, fullKey)
        }
        return [fullKey]
      })
    }

    const enKeys = getKeys(en).sort()
    const ruKeys = getKeys(ru).sort()
    expect(enKeys).toEqual(ruKeys)
  })

  it('translates common keys in English', () => {
    const i18n = createTestI18n('en')
    const { t } = i18n.global

    expect(t('common.save')).toBe('Save')
    expect(t('common.cancel')).toBe('Cancel')
    expect(t('common.delete')).toBe('Delete')
  })

  it('translates common keys in Russian', () => {
    const i18n = createTestI18n('ru')
    const { t } = i18n.global

    expect(t('common.save')).toBe('Сохранить')
    expect(t('common.cancel')).toBe('Отмена')
    expect(t('common.delete')).toBe('Удалить')
  })

  it('handles interpolation in English', () => {
    const i18n = createTestI18n('en')
    const { t } = i18n.global

    expect(t('dashboard.totalTasks', { count: 42 })).toBe('42 tasks total')
    expect(t('validation.required', { field: 'Email' })).toBe('Email is required')
  })

  it('handles interpolation in Russian', () => {
    const i18n = createTestI18n('ru')
    const { t } = i18n.global

    expect(t('dashboard.totalTasks', { count: 42 })).toBe('42 задач всего')
    expect(t('validation.required', { field: 'Email' })).toBe('Поле Email обязательно')
  })

  it('all task statuses have translations', () => {
    const i18n = createTestI18n('en')
    const { t } = i18n.global

    expect(t('task.new')).toBe('New')
    expect(t('task.inProgress')).toBe('In Progress')
    expect(t('task.done')).toBe('Done')
  })

  it('no translation returns key (fallback check)', () => {
    const i18n = createTestI18n('en')
    const { t } = i18n.global

    // Non-existent key should return the key path
    expect(t('nonexistent.key')).toBe('nonexistent.key')
  })

  it('switches locale dynamically', () => {
    const i18n = createTestI18n('en')
    const { t } = i18n.global

    expect(t('common.save')).toBe('Save')

    i18n.global.locale.value = 'ru'
    expect(t('common.save')).toBe('Сохранить')
  })
})
