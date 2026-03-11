import { test, expect } from '@playwright/test'

test.describe('Navigation & Route Guards', () => {
  test('unauthenticated user is redirected to login', async ({ page }) => {
    await page.goto('/')
    await expect(page).toHaveURL(/\/login/)
  })

  test('unauthenticated user cannot access stats', async ({ page }) => {
    await page.goto('/stats')
    await expect(page).toHaveURL(/\/login/)
  })

  test('unknown routes redirect to root', async ({ page }) => {
    await page.goto('/some-random-page')
    await expect(page).toHaveURL(/\/login/)
  })

  test('authenticated user sees dashboard', async ({ page }) => {
    await page.goto('/login')
    await page.evaluate(() => {
      localStorage.setItem('access_token', 'fake-token-for-e2e')
      localStorage.setItem('refresh_token', 'fake-refresh-for-e2e')
    })
    await page.goto('/')

    // Dashboard heading should be visible (API fails but page renders)
    await expect(page.locator('h1')).toContainText(/Tasks|Задачи/, { timeout: 10000 })
  })

  test('authenticated user is redirected from login to dashboard', async ({ page }) => {
    await page.goto('/login')
    await page.evaluate(() => {
      localStorage.setItem('access_token', 'fake-token-for-e2e')
      localStorage.setItem('refresh_token', 'fake-refresh-for-e2e')
    })
    await page.goto('/login')

    await expect(page).not.toHaveURL(/\/login/)
  })
})
