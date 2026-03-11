import { test, expect } from '@playwright/test'

// Helper to set up fake auth
async function authenticateUser(page: import('@playwright/test').Page) {
  await page.goto('/login')
  await page.evaluate(() => {
    localStorage.setItem('access_token', 'fake-token-for-e2e')
    localStorage.setItem('refresh_token', 'fake-refresh-for-e2e')
  })
}

test.describe('Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    await authenticateUser(page)
  })

  test('dashboard page loads with header', async ({ page }) => {
    await page.goto('/')
    await expect(page.locator('h1').first()).toContainText(/Tasks|Задачи/, {
      timeout: 10000,
    })
  })

  test('demo mode toggle works', async ({ page }) => {
    await page.goto('/')

    const demoButton = page.getByRole('button', { name: /Demo|Демо/ })
    await expect(demoButton).toBeVisible({ timeout: 10000 })

    await demoButton.click()

    // Should show 10K demo indicator
    await expect(page.getByText(/10K|10 000/)).toBeVisible({ timeout: 5000 })
  })

  test('view mode switches between cards and table', async ({ page }) => {
    await page.goto('/')

    // Activate demo mode first to have data
    const demoButton = page.getByRole('button', { name: /Demo|Демо/ })
    await expect(demoButton).toBeVisible({ timeout: 10000 })
    await demoButton.click()

    // Wait for demo data to appear
    await expect(page.getByText(/10K|10 000/)).toBeVisible({ timeout: 5000 })

    // Find the table view toggle (second button in the view switcher)
    const tableButton = page.locator(
      'button:has(svg.lucide-table2), button:has(svg.lucide-table-2)',
    )
    if (await tableButton.count()) {
      await tableButton.first().click()

      // Table view should show search in table and column headers
      await expect(page.getByPlaceholder(/Search in table|Поиск в таблице/)).toBeVisible({
        timeout: 5000,
      })
    }
  })

  test('charts are visible on dashboard', async ({ page }) => {
    await page.goto('/')

    // Activate demo mode
    const demoButton = page.getByRole('button', { name: /Demo|Демо/ })
    await expect(demoButton).toBeVisible({ timeout: 10000 })
    await demoButton.click()

    // Chart sections should be visible
    await expect(page.getByText(/Tasks Created Over Time|Создание задач по времени/)).toBeVisible()
    await expect(page.getByText(/Tasks by Status|Задачи по статусу/)).toBeVisible()

    // Canvas elements for charts should exist
    const canvases = page.locator('canvas')
    await expect(canvases.first()).toBeVisible({ timeout: 5000 })
  })

  test('new task button is visible when not in demo mode', async ({ page }) => {
    await page.goto('/')

    await expect(page.getByRole('button', { name: /New Task|Новая задача/ })).toBeVisible({
      timeout: 10000,
    })
  })

  test('new task button is hidden in demo mode', async ({ page }) => {
    await page.goto('/')

    const demoButton = page.getByRole('button', { name: /Demo|Демо/ })
    await expect(demoButton).toBeVisible({ timeout: 10000 })
    await demoButton.click()

    await expect(page.getByRole('button', { name: /New Task|Новая задача/ })).not.toBeVisible()
  })
})
