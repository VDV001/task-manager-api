import { test, expect } from '@playwright/test'

test.describe('Authentication', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/login')
    await page.evaluate(() => {
      localStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
    })
  })

  test('login page renders correctly', async ({ page }) => {
    await page.goto('/login')

    await expect(page.getByText('Welcome back')).toBeVisible()
    await expect(page.getByText('Sign in to your account')).toBeVisible()
    await expect(page.getByPlaceholder('you@example.com')).toBeVisible()
    await expect(page.getByPlaceholder('Enter your password')).toBeVisible()
    await expect(page.getByRole('button', { name: 'Sign in' })).toBeVisible()
    await expect(page.getByText("Don't have an account?")).toBeVisible()
  })

  test('login form validates required fields', async ({ page }) => {
    await page.goto('/login')

    await page.getByRole('button', { name: 'Sign in' }).click()

    await expect(page.getByText(/Email is required/i)).toBeVisible()
    await expect(page.getByText(/Password is required/i)).toBeVisible()
  })

  test('register page renders correctly', async ({ page }) => {
    await page.goto('/register')

    await expect(page.getByRole('heading', { name: 'Create account' })).toBeVisible()
    await expect(page.getByPlaceholder('John Doe')).toBeVisible()
    await expect(page.getByPlaceholder('you@example.com')).toBeVisible()
    await expect(page.getByRole('button', { name: 'Create account' })).toBeVisible()
    await expect(page.getByText('Already have an account?')).toBeVisible()
  })

  test('register form validates required fields', async ({ page }) => {
    await page.goto('/register')

    await page.getByRole('button', { name: 'Create account' }).click()

    await expect(page.getByText(/Name is required/i)).toBeVisible()
    await expect(page.getByText(/Email is required/i)).toBeVisible()
    await expect(page.getByText(/Password is required/i)).toBeVisible()
  })

  test('register form validates password min length', async ({ page }) => {
    await page.goto('/register')

    await page.getByPlaceholder('John Doe').fill('Test')
    await page.getByPlaceholder('you@example.com').fill('test@test.com')
    await page.getByPlaceholder('Min. 6 characters').fill('12345')

    await page.getByRole('button', { name: 'Create account' }).click()

    await expect(page.getByText(/Minimum 6 characters/i)).toBeVisible()
  })

  test('navigate from login to register', async ({ page }) => {
    await page.goto('/login')

    await page.getByText('Sign up').click()
    await expect(page).toHaveURL(/\/register/)
  })

  test('navigate from register to login', async ({ page }) => {
    await page.goto('/register')

    await page.getByText('Sign in', { exact: true }).click()
    await expect(page).toHaveURL(/\/login/)
  })
})
