import { test, expect } from '@playwright/test';

test.describe('Work Breakdown Structure (WBS)', () => {
    test.beforeEach(async ({ page }) => {
        page.on('pageerror', err => console.log('PAGE ERROR:', err.message, err.stack));
        page.on('console', msg => {
            if (msg.type() === 'error') console.log('CONSOLE ERROR:', msg.text());
        });

        await page.goto('/projects');
        await page.click('text=Test Project Alpha');
        await expect(page).toHaveURL(/\/projects\/\d+/);

        // Navigate to WBS page
        await page.click('text=Open WBS');
        await expect(page).toHaveURL(/\/projects\/\d+\/wbs/);
    });

    test('should display WBS tree and allow creating a node', async ({ page }) => {
        await expect(page.locator('text=Project WBS & Timeline')).toBeVisible();

        // Check if "Phase 1" exists (from seed)
        await expect(page.locator('text=Phase 1')).toBeVisible();

        // Create a new task with unique name
        const uniqueTaskName = `E2E Task ${Date.now()}`;
        await page.click('button:has-text("New Task")');

        await page.fill('#title', uniqueTaskName);
        await page.selectOption('#type', 'Task');
        await page.fill('#progress', '10');
        await page.fill('#effort', '16');

        // Dates (ensure within project range)
        await page.fill('#start', '2026-04-01');
        await page.fill('#end', '2026-04-05');

        await page.click('button:has-text("Create")');

        // Verify toast or list update
        await expect(page.locator('text=Added successfully')).toBeVisible();
        await expect(page.locator(`text=${uniqueTaskName}`)).toBeVisible();
    });

    test('should allow editing a node', async ({ page }) => {
        // Find "Phase 1" and click edit 
        const phaseRow = page.locator('div[role="row"]:has-text("Phase 1")');
        await phaseRow.hover(); // Trigger actions visibility
        await phaseRow.locator('button[title="Edit Task"]').click();

        const uniquePhaseName = `Phase 1 - Upd ${Date.now()}`;
        await page.fill('#title', uniquePhaseName);

        // Fill dates to pass validation
        await page.fill('#start', '2026-05-01');
        await page.fill('#end', '2026-05-10');

        await page.click('button:has-text("Save Changes")');

        await expect(page.locator('text=Updated successfully')).toBeVisible();
        await expect(page.locator(`text=${uniquePhaseName}`)).toBeVisible();
    });
});
