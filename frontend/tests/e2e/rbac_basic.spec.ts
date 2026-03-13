import { test, expect } from '@playwright/test';

test.describe('RBAC - Sidebar Menu Visibility', () => {
    test('Menu visibility based on role', async ({ page }, testInfo) => {
        const role = testInfo.project.name;
        if (role === 'setup') return; // Skip setup project

        await page.goto('/');

        // Everyone should see Projects (basic access)
        await expect(page.getByRole('link', { name: 'Projects' })).toBeVisible();

        if (role === 'admin' || role === 'pmo') {
            await expect(page.getByRole('link', { name: 'Portfolio' })).toBeVisible();
        }

        if (role === 'admin') {
            await expect(page.getByRole('link', { name: 'Users' })).toBeVisible();
            await expect(page.getByRole('link', { name: 'Settings' })).toBeVisible();
        } else {
            // Non-admin roles should NOT see these
            await expect(page.getByRole('link', { name: 'Users' })).not.toBeVisible();
            await expect(page.getByRole('link', { name: 'Settings' })).not.toBeVisible();
        }
    });
});
