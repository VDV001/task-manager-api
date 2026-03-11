<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { toast } from 'vue-sonner'
import AuthLayout from '@/components/auth/AuthLayout.vue'
import AppInput from '@/components/ui/AppInput.vue'
import AppButton from '@/components/ui/AppButton.vue'
import type { ApiError } from '@/types/api'

const router = useRouter()
const auth = useAuthStore()

const email = ref('')
const password = ref('')
const loading = ref(false)
const errors = ref<Record<string, string>>({})

async function handleSubmit() {
  errors.value = {}
  if (!email.value) errors.value.email = 'Email is required'
  if (!password.value) errors.value.password = 'Password is required'
  if (Object.keys(errors.value).length) return

  loading.value = true
  try {
    await auth.login({ email: email.value, password: password.value })
    toast.success('Welcome back!')
    router.push({ name: 'dashboard' })
  } catch (err: any) {
    const apiErr = err.data as ApiError | undefined
    if (apiErr?.error?.code === 'UNAUTHORIZED') {
      toast.error('Invalid email or password')
    } else {
      toast.error(apiErr?.error?.message ?? 'Something went wrong')
    }
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AuthLayout>
    <h2 class="text-xl font-semibold text-text-primary mb-1">Welcome back</h2>
    <p class="text-sm text-text-secondary mb-6">Sign in to your account</p>

    <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
      <AppInput
        v-model="email"
        label="Email"
        type="email"
        placeholder="you@example.com"
        :error="errors.email"
      />
      <AppInput
        v-model="password"
        label="Password"
        type="password"
        placeholder="Enter your password"
        :error="errors.password"
      />
      <AppButton type="submit" :loading="loading" class="mt-2 w-full">
        Sign in
      </AppButton>
    </form>

    <p class="text-sm text-text-muted text-center mt-6">
      Don't have an account?
      <RouterLink to="/register" class="text-accent hover:text-accent-hover transition-colors">
        Sign up
      </RouterLink>
    </p>
  </AuthLayout>
</template>
