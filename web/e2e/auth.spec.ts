import { test, expect } from '@playwright/test';
import { mockAPI, loginAsGuest } from './helpers';

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
    await expect(page.getByText('Real stocks. Virtual money. Zero risk.')).toBeVisible();
  });

  test('authenticated user is redirected from login to market', async ({ page }) => {
    // Middleware redirects logged-in users away from /login
    await page.context().addCookies([{
      name: 'mockstarket_token',
      value: 'guest_e2e_test_token',
      domain: 'localhost',
      path: '/',
    }]);
    await page.goto('/login');
    await expect(page).toHaveURL('/market');
  });

  test('sets token in localStorage after login', async ({ page }) => {
    await loginAsGuest(page);

    const token = await page.evaluate(() => localStorage.getItem('mockstarket_token'));
    expect(token).toBeTruthy();
    expect(token).toContain('guest_');
  });
});
