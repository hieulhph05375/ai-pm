import { test, expect } from '@playwright/test';

test.describe('Authentication', () => {
    test('should show error with invalid credentials', async ({ page }) => {
        await page.goto('/login');

        await page.fill('#email', 'wrong@example.com');
        await page.fill('#password', 'wrongpassword');
        await page.click('button[type="submit"]');

        const errorAlert = page.locator('[role="alert"]');
        await expect(errorAlert).toBeVisible();
        await expect(errorAlert).toContainText('Login failed');
    });

    test('should login successfully with admin credentials', async ({ page }) => {
        await page.goto('/login');

        // Using default admin credentials from seed data
        await page.fill('#email', 'admin@example.com');
        await page.fill('#password', 'password');
        await page.click('button[type="submit"]');

        // Should redirect to dashboard/home
        await expect(page).toHaveURL('/');

        // Check for some dashboard element (e.g., "Running Projects")
        await expect(page.locator('text=Running Projects')).toBeVisible();
    });
});
