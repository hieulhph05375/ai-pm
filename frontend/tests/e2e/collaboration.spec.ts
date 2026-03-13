import { test, expect } from '@playwright/test';

test.describe('Collaboration & Notifications', () => {
    test.beforeEach(async ({ page }) => {
        await page.goto('/');
    });

    test('should open notification dropdown', async ({ page }) => {
        const bell = page.getByLabel('Notifications');
        await expect(bell).toBeVisible();

        await bell.click();

        // Check for dropdown header
        await expect(page.locator('text=Notifications').first()).toBeVisible();
    });

    // Note: To test actual receipt of notifications, we'd need to trigger an event
    // or seed a notification for the user.
    test('should show empty state if no notifications', async ({ page }) => {
        await page.click('button[aria-label="Notifications"]');

        // Since we just seeded the DB and haven't run cron yet, it might be empty
        // or contain the sample notifications if seeded.
        // SeedAll in fixtures_test.go doesn't seed notifications yet.

        await expect(page.locator('text=No notifications')).toBeVisible();
    });
});
