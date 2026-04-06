import { test, expect } from '@playwright/test';
import { mockAPI } from './helpers';

test.describe('Authentication', () => {
  test.beforeEach(async ({ page }) => {
    await mockAPI(page);
  });

  test('redirects unauthenticated user to login', async ({ page }) => {
    await page.goto('/market');
    await expect(page).toHaveURL('/login');
  });

  test('login page shows app branding', async ({ page }) => {
    await page.goto('/login');
    await expect(page.getByText('Mock Starket')).toBeVisible();
    await expect(page.getByText('Learn to trade. Risk nothing.')).toBeVisible();
  });

  test('guest login redirects to market', async ({ page }) => {
    await page.goto('/login');
    await page.getByRole('button', { name: 'Continue as Guest' }).click();
    await expect(page).toHaveURL('/market');
  });

  test('sets token in localStorage after login', async ({ page }) => {
    await page.goto('/login');
    await page.getByRole('button', { name: 'Continue as Guest' }).click();
    await page.waitForURL('/market');

    const token = await page.evaluate(() => localStorage.getItem('mockstarket_token'));
    expect(token).toBeTruthy();
    expect(token).toContain('guest_');
  });
});
