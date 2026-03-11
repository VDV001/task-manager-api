<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import { toast } from 'vue-sonner'
import AuthLayout from '@/components/auth/AuthLayout.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppButton from '@/components/ui/AppButton.vue'
import type { ApiError } from '@/types/api'

const { t } = useI18n()
const router = useRouter()
const auth = useAuthStore()

const name = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)
const errors = ref<Record<string, string>>({})

async function handleSubmit() {
  errors.value = {}
  if (!name.value) errors.value.name = t('validation.required', { field: t('auth.name') })
  if (!email.value) errors.value.email = t('validation.required', { field: t('auth.email') })
  if (!password.value)
    errors.value.password = t('validation.required', { field: t('auth.password') })
  else if (password.value.length < 6)
    errors.value.password = t('validation.minLength', { count: 6 })
  if (Object.keys(errors.value).length) return

  loading.value = true
  try {
    await auth.register({ name: name.value, email: email.value, password: password.value })
    toast.success(t('auth.accountCreated'))
    router.push({ name: 'dashboard' })
  } catch (err: unknown) {
    const apiErr = (err as Record<string, unknown>)?.data as ApiError | undefined
    if (apiErr?.error?.code === 'CONFLICT') {
      errors.value.email = t('auth.emailTaken')
    } else if (apiErr?.error?.details?.length) {
      apiErr.error.details.forEach((d) => {
        errors.value[d.field] = d.message
      })
    } else {
      toast.error(apiErr?.error?.message ?? t('auth.somethingWrong'))
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AuthLayout>
    <h2 class="text-xl font-semibold text-text-primary mb-1">{{ t('auth.register') }}</h2>
    <p class="text-sm text-text-secondary mb-6">{{ t('auth.createAccountSubtitle') }}</p>

    <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
      <AppInput
        v-model="name"
        :label="t('auth.name')"
        :placeholder="t('auth.namePlaceholder')"
        :error="errors.name"
      />
      <AppInput
        v-model="email"
        :label="t('auth.email')"
        type="email"
        :placeholder="t('auth.emailPlaceholder')"
        :error="errors.email"
      />
      <AppInput
        v-model="password"
        :label="t('auth.password')"
        type="password"
        :placeholder="t('auth.passwordMinLength')"
        :error="errors.password"
      />
      <AppButton type="submit" :loading="loading" class="mt-2 w-full">
        {{ t('auth.register') }}
      </AppButton>
    </form>

    <p class="text-sm text-text-muted text-center mt-6">
      {{ t('auth.hasAccount') }}
      <RouterLink to="/login" class="text-accent hover:text-accent-hover transition-colors">
        {{ t('auth.signIn') }}
      </RouterLink>
    </p>
  </AuthLayout>
</template>
