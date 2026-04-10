import { test, expect } from '@playwright/test';
import { mockAPI, loginAsGuest } from './helpers';

test.describe('Market Page', () => {
  test.beforeEach(async ({ page }) => {
    await mockAPI(page);
    await loginAsGuest(page);
  });

  test('displays market heading and summary section', async ({ page }) => {
    await expect(page.getByRole('heading', { name: 'Market' })).toBeVisible();
    await expect(page.getByText('Market Index')).toBeVisible();
  });

  test('displays search input and asset filters', async ({ page }) => {
    await expect(page.getByPlaceholder('Search stocks...')).toBeVisible();
    await expect(page.getByRole('button', { name: 'All', exact: true })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Stocks' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'ETFs' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Crypto' })).toBeVisible();
  });

  test('displays stock table with headers', async ({ page }) => {
    await expect(page.getByRole('table')).toBeVisible();
    await expect(page.getByRole('columnheader', { name: /Stock/ })).toBeVisible();
    await expect(page.getByRole('columnheader', { name: /Price/ })).toBeVisible();
  });

  test('stock table contains rows', async ({ page }) => {
    // Wait for stock data to load and render into the table
    const firstRow = page.locator('tbody tr').first();
    await expect(firstRow).toBeVisible({ timeout: 10_000 });
    expect(await page.locator('tbody tr').count()).toBeGreaterThan(0);
  });

  test('search clears input and table still has rows', async ({ page }) => {
    // Wait for data to load
    await expect(page.locator('tbody tr').first()).toBeVisible({ timeout: 10_000 });

    // Type and clear search
    const searchInput = page.getByPlaceholder('Search stocks...');
    await searchInput.fill('test');
    await page.waitForTimeout(400);
    await searchInput.clear();
    await page.waitForTimeout(400);

    // Table should still have rows after clearing
    expect(await page.locator('tbody tr').count()).toBeGreaterThan(0);
  });

  test('asset filter buttons are clickable', async ({ page }) => {
    await page.getByRole('button', { name: 'ETFs' }).click();
    // Should filter — either shows ETF rows or an empty table
    await expect(page.getByRole('table')).toBeVisible();
  });
});
