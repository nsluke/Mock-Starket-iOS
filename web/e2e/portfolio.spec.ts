import { test, expect } from '@playwright/test';
import { mockAPI, loginAsGuest } from './helpers';

test.describe('Portfolio Page', () => {
  test.beforeEach(async ({ page }) => {
    await mockAPI(page);
    await loginAsGuest(page);
    await page.getByRole('link', { name: 'Portfolio' }).click();
    await expect(page).toHaveURL('/portfolio');
  });

  test('displays portfolio summary cards', async ({ page }) => {
    await expect(page.getByText('Net Worth')).toBeVisible();
    await expect(page.getByText('$105,000.00')).toBeVisible();
    // Use the main content area to avoid matching sidebar balance
    const main = page.locator('main');
    await expect(main.getByText('Cash')).toBeVisible();
    await expect(main.getByText('$85,000.00')).toBeVisible();
  });

  test('displays holdings with position data', async ({ page }) => {
    const main = page.locator('main');
    await expect(main.getByText('PIPE')).toBeVisible();
    await expect(main.getByText('$15,550.00')).toBeVisible();
  });

  test('switches between holdings and trade history tabs', async ({ page }) => {
    await expect(page.getByText('Holdings')).toBeVisible();
    await page.getByRole('button', { name: 'Trade History' }).click();
    await expect(page.getByText('No trades yet.')).toBeVisible();
  });
});
