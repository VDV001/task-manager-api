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

const name = ref('')
const email = ref('')
const password = ref('')
const loading = ref(false)
const errors = ref<Record<string, string>>({})

async function handleSubmit() {
  errors.value = {}
  if (!name.value) errors.value.name = 'Name is required'
  if (!email.value) errors.value.email = 'Email is required'
  if (!password.value) errors.value.password = 'Password is required'
  else if (password.value.length < 6) errors.value.password = 'Minimum 6 characters'
  if (Object.keys(errors.value).length) return

  loading.value = true
  try {
    await auth.register({ name: name.value, email: email.value, password: password.value })
    toast.success('Account created!')
    router.push({ name: 'dashboard' })
  } catch (err: any) {
    const apiErr = err.data as ApiError | undefined
    if (apiErr?.error?.code === 'CONFLICT') {
      errors.value.email = 'Email already registered'
    } else if (apiErr?.error?.details?.length) {
      apiErr.error.details.forEach((d) => {
        errors.value[d.field] = d.message
      })
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
    <h2 class="text-xl font-semibold text-text-primary mb-1">Create account</h2>
    <p class="text-sm text-text-secondary mb-6">Get started with TaskFlow</p>

    <form class="flex flex-col gap-4" @submit.prevent="handleSubmit">
      <AppInput v-model="name" label="Name" placeholder="John Doe" :error="errors.name" />
      <AppInput v-model="email" label="Email" type="email" placeholder="you@example.com" :error="errors.email" />
      <AppInput v-model="password" label="Password" type="password" placeholder="Min. 6 characters" :error="errors.password" />
      <AppButton type="submit" :loading="loading" class="mt-2 w-full">
        Create account
      </AppButton>
    </form>

    <p class="text-sm text-text-muted text-center mt-6">
      Already have an account?
      <RouterLink to="/login" class="text-accent hover:text-accent-hover transition-colors">
        Sign in
      </RouterLink>
    </p>
  </AuthLayout>
</template>
