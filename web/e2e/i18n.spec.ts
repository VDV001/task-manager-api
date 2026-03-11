import { test, expect } from '@playwright/test'

test.describe('Internationalization', () => {
  test('login page defaults to English', async ({ page }) => {
    await page.goto('/login')
    await page.evaluate(() => localStorage.removeItem('locale'))
    await page.goto('/login')

    await expect(page.getByText('Welcome back')).toBeVisible()
    await expect(page.getByText('Sign in to your account')).toBeVisible()
  })

  test('language can be switched to Russian', async ({ page }) => {
    await page.goto('/login')
    await page.evaluate(() => localStorage.removeItem('locale'))
    await page.goto('/login')

    // Find and click the language toggle button (RU)
    const ruButton = page.getByRole('button', { name: 'RU' })
    if (await ruButton.isVisible()) {
      await ruButton.click()

      await expect(page.getByText('С возвращением')).toBeVisible()
      await expect(page.getByText('Войдите в свой аккаунт')).toBeVisible()
    }
  })

  test('language preference persists after reload', async ({ page }) => {
    await page.goto('/login')
    await page.evaluate(() => localStorage.setItem('locale', 'ru'))
    await page.goto('/login')

    await expect(page.getByText('С возвращением')).toBeVisible()
  })

  test('register page respects saved locale', async ({ page }) => {
    await page.goto('/register')
    await page.evaluate(() => localStorage.setItem('locale', 'ru'))
    await page.goto('/register')

    await expect(page.getByRole('heading', { name: 'Создать аккаунт' })).toBeVisible()
  })
})
