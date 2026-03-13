import { test, expect } from '@playwright/test';

test.describe('Visual Regression - Baselines', () => {
    test('Login Page - Visual Baseline', async ({ page }) => {
        await page.goto('/login');
        await expect(page).toHaveScreenshot('login-page.png');
    });

    test('Dashboard - Visual Baseline', async ({ page }) => {
        // This will use the storageState from the project (admin by default if run in admin project)
        await page.goto('/');
        await expect(page).toHaveScreenshot('dashboard-home.png');
    });

    test('Sidebar - Expanded Visual Baseline', async ({ page }) => {
        await page.goto('/');
        // Sidebar is usually expanded by default
        await expect(page.locator('aside')).toHaveScreenshot('sidebar-expanded.png');
    });
});
