import { test, expect } from '@playwright/test';

test.describe('Analytics & Portfolio', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/');
    });

    test('should display Portfolio Dashboard with stats', async ({ page }) => {
        await page.goto('/portfolio');
        await expect(page.locator('h1.text-3xl')).toContainText('Portfolio Dashboard');

        // Check for stat cards
        await expect(page.locator('text=Total Projects')).toBeVisible();
        await expect(page.locator('text=Total Budget')).toBeVisible();

        // Wait for data to load (skeleton should disappear)
        await expect(page.locator('text=Portfolio Health')).toBeVisible();
    });

    test('should display Resource Workload Heatmap', async ({ page }) => {
        await page.goto('/portfolio');
        // The heatmap is a component in the portfolio page
        await expect(page.locator('text=Resource Workload Analysis')).toBeVisible();
    });
});
