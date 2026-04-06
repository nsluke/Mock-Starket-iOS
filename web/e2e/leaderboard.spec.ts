import { test, expect } from '@playwright/test';
import { mockAPI, loginAsGuest } from './helpers';

test.describe('Leaderboard Page', () => {
  test.beforeEach(async ({ page }) => {
    await mockAPI(page);
    await loginAsGuest(page);
    await page.getByRole('link', { name: 'Leaderboard' }).click();
  });

  test('displays leaderboard with rankings', async ({ page }) => {
    await expect(page.getByText('TraderJoe')).toBeVisible();
    await expect(page.getByText('StockPro')).toBeVisible();
    await expect(page.getByText('$120,000.00')).toBeVisible();
  });

  test('shows period toggle buttons', async ({ page }) => {
    await expect(page.getByRole('button', { name: 'All Time' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Weekly' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Daily' })).toBeVisible();
  });
});
