import { defineConfig, devices } from '@playwright/test';

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
    testDir: './tests/e2e',
    /* Maximum time one test can run for. */
    timeout: 30 * 1000,
    expect: {
        /**
         * Maximum time expect() should wait for the condition to be met.
         * For example in `await expect(locator).toBeVisible();`
         */
        timeout: 5000
    },
    /* Run tests in files in parallel */
    fullyParallel: true,
    /* Fail the build on CI if you accidentally left test.only in the source code. */
    forbidOnly: !!process.env.CI,
    /* Retry on CI only */
    retries: process.env.CI ? 2 : 0,
    /* Opt out of parallel tests on local. */
    workers: process.env.CI ? 1 : undefined,
    /* Reporter to use. See https://playwright.dev/docs/test-reporters */
    reporter: 'html',
    /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
    use: {
        /* Base URL to use in actions like `await page.goto('/')`. */
        baseURL: 'http://localhost:5173',

        /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
        trace: 'on-first-retry',

        screenshot: 'only-on-failure',
    },

    /* Configure projects for major browsers */
    projects: [
        { name: 'setup', testMatch: /.*\.setup\.ts/ },
        {
            name: 'admin',
            use: {
                ...devices['Desktop Chrome'],
                storageState: 'tests/e2e/.auth/admin.json',
            },
            dependencies: ['setup'],
        },
        {
            name: 'pmo',
            use: {
                ...devices['Desktop Chrome'],
                storageState: 'tests/e2e/.auth/pmo.json',
            },
            dependencies: ['setup'],
        },
        {
            name: 'pm',
            use: {
                ...devices['Desktop Chrome'],
                storageState: 'tests/e2e/.auth/pm.json',
            },
            dependencies: ['setup'],
        },
        {
            name: 'teamlead',
            use: {
                ...devices['Desktop Chrome'],
                storageState: 'tests/e2e/.auth/teamlead.json',
            },
            dependencies: ['setup'],
        },
        {
            name: 'member',
            use: {
                ...devices['Desktop Chrome'],
                storageState: 'tests/e2e/.auth/member.json',
            },
            dependencies: ['setup'],
        },
        {
            name: 'viewer',
            use: {
                ...devices['Desktop Chrome'],
                storageState: 'tests/e2e/.auth/viewer.json',
            },
            dependencies: ['setup'],
        },
    ],

    /* Run your local dev server before starting the tests */
    webServer: [
        {
            command: 'cd ../backend && PORT=8081 DATABASE_URL=postgres://postgres:Caikeo@1234@localhost:5432/projectmgmt_test?sslmode=disable go run cmd/server/main.go',
            url: 'http://localhost:8081/api/v1/roles',
            reuseExistingServer: false,
            timeout: 120 * 1000,
            ignoreHTTPSErrors: true,
        },
        {
            command: 'pnpm dev',
            url: 'http://localhost:5173',
            reuseExistingServer: !process.env.CI,
            timeout: 120 * 1000,
        }
    ],
});
