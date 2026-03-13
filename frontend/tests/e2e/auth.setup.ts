import { test as setup, expect } from '@playwright/test';
import fs from 'fs';
import path from 'path';

const authDir = 'tests/e2e/.auth';

const users = [
    { role: 'admin', email: 'admin@example.com', password: 'password' },
    { role: 'pmo', email: 'pmo@example.com', password: 'password' },
    { role: 'pm', email: 'projectmanager@example.com', password: 'password' },
    { role: 'teamlead', email: 'teamlead@example.com', password: 'password' },
    { role: 'member', email: 'member@example.com', password: 'password' },
    { role: 'viewer', email: 'viewer@example.com', password: 'password' },
];

setup('authenticate all roles', async ({ browser }) => {
    // Ensure auth directory exists
    if (!fs.existsSync(authDir)) {
        fs.mkdirSync(authDir, { recursive: true });
    }

    for (const user of users) {
        console.log(`Authenticating role: ${user.role} (${user.email})...`);
        const context = await browser.newContext();
        const page = await context.newPage();
        const authFile = path.join(authDir, `${user.role}.json`);

        await page.goto('/login');
        await page.fill('#email', user.email);
        await page.fill('#password', user.password);
        await page.click('button[type="submit"]');

        // Wait for login to complete (redirect to home or visible element)
        try {
            await page.waitForURL(url => url.pathname === '/', { timeout: 15000 });
            console.log(`Successfully authenticated role: ${user.role}`);
        } catch (err) {
            console.error(`Failed to authenticate role: ${user.role}. URL: ${page.url()}`);
            // Take a debug screenshot
            await page.screenshot({ path: path.join(authDir, `fail_${user.role}.png`) });
            throw err;
        }

        // Save storage state for this role
        await context.storageState({ path: authFile });
        await context.close();
    }
});
