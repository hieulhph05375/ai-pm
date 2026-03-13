import { test, expect } from '@playwright/test';

test.describe('Project Management', () => {
    test('should create a new project successfully', async ({ page }) => {
        await page.goto('/projects');

        // Click Create New Project
        await page.click('button:has-text("Create New Project")');

        // Fill the form
        const projectID = `E2E-${Date.now()}`;
        const projectName = `E2E Test Project ${Date.now()}`;

        await page.fill('input[placeholder="Enter project name..."]', projectName);
        await page.fill('input[placeholder="Example: PRJ-001"]', projectID);
        await page.fill('input[placeholder="Manager name..."]', 'E2E Manager');

        // Save Project
        await page.click('button:has-text("Save Project")');

        // Check if project appears in the list
        await expect(page.locator(`text=${projectName}`)).toBeVisible();
        await expect(page.locator(`text=${projectID}`)).toBeVisible();
    });

    test('should search for a project', async ({ page }) => {
        await page.goto('/projects');

        // Search for the project we might have created or a known one
        // Let's assume there is at least one project from seed data or previous test
        const searchInput = page.locator('input[placeholder="Search project name, ID..."]');
        await searchInput.fill('E2E');
        await searchInput.press('Enter');

        // Verify result appears
        await expect(page.locator('tbody tr').first()).toBeVisible();
    });
});
