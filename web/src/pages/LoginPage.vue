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

const email = ref('')
const password = ref('')
const loading = ref(false)
const errors = ref<Record<string, string>>({})

async function handleSubmit() {
  errors.value = {}
  if (!email.value) errors.value.email = t('validation.required', { field: t('auth.email') })
  if (!password.value)
    errors.value.password = t('validation.required', { field: t('auth.password') })
  if (Object.keys(errors.value).length) return

  loading.value = true
  try {
    await auth.login({ email: email.value, password: password.value })
    toast.success(t('auth.welcomeToast'))
    router.push({ name: 'dashboard' })
  } catch (err: unknown) {
    const apiErr = (err as Record<string, unknown>)?.data as ApiError | undefined
    if (apiErr?.error?.code === 'UNAUTHORIZED') {
      toast.error(t('auth.invalidCredentials'))
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
    <h2 class="text-xl font-semibold text-text-primary mb-1">{{ t('auth.welcomeBack') }}</h2>
    <p class="text-sm text-text-secondary mb-6">{{ t('auth.signInSubtitle') }}</p>

    <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
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
        :placeholder="t('auth.passwordPlaceholder')"
        :error="errors.password"
      />
      <AppButton type="submit" :loading="loading" class="mt-2 w-full">
        {{ t('auth.login') }}
      </AppButton>
    </form>

    <p class="text-sm text-text-muted text-center mt-6">
      {{ t('auth.noAccount') }}
      <RouterLink to="/register" class="text-accent hover:text-accent-hover transition-colors">
        {{ t('auth.signUp') }}
      </RouterLink>
    </p>
  </AuthLayout>
</template>
