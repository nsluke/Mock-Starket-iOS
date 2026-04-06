import { test, expect } from '@playwright/test';
import { mockAPI, loginAsGuest } from './helpers';

test.describe('Navigation', () => {
  test.beforeEach(async ({ page }) => {
    await mockAPI(page);
    await loginAsGuest(page);
  });

  test('sidebar shows all nav items', async ({ page }) => {
    await expect(page.getByRole('link', { name: 'Market' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Portfolio' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Orders' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Leaderboard' })).toBeVisible();
  });

  test('sidebar shows extra nav items', async ({ page }) => {
    await expect(page.getByRole('link', { name: 'Watchlist' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Alerts' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Challenges' })).toBeVisible();
    await expect(page.getByRole('link', { name: 'Achievements' })).toBeVisible();
  });

  test('navigates to each main page', async ({ page }) => {
    await page.getByRole('link', { name: 'Portfolio' }).click();
    await expect(page).toHaveURL('/portfolio');

    await page.getByRole('link', { name: 'Orders' }).click();
    await expect(page).toHaveURL('/orders');

    await page.getByRole('link', { name: 'Leaderboard' }).click();
    await expect(page).toHaveURL('/leaderboard');

    await page.getByRole('link', { name: 'Settings' }).click();
    await expect(page).toHaveURL('/settings');
  });

  test('shows connection status in sidebar', async ({ page }) => {
    await expect(page.getByText(/Connected|Reconnecting|Disconnected/)).toBeVisible();
  });
});
