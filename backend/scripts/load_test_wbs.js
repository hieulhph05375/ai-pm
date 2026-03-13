import http from 'k6/http';
import { check, sleep } from 'k6';

// Read target Project ID from ENV, default to 34 (the one we seeded)
const PROJECT_ID = __ENV.PROJECT_ID || 34;
const BASE_URL = 'http://localhost:8080/api/v1';

export const options = {
    // Ramp up to 50 virtual users over 10s, hold for 20s, ramp down over 5s
    stages: [
        { duration: '10s', target: 50 },
        { duration: '20s', target: 50 },
        { duration: '5s', target: 0 },
    ],
    thresholds: {
        // We want the 95th percentile response time to be under 500ms initially
        http_req_duration: ['p(95)<500'],
    },
};

// Login once per VU to get token
export function setup() {
    const loginRes = http.post(`${BASE_URL}/auth/login`, JSON.stringify({
        email: 'admin@example.com',
        password: 'password'
    }), {
        headers: { 'Content-Type': 'application/json' }
    });

    if (loginRes.status !== 200) {
        throw new Error('Login failed: ' + loginRes.body);
    }

    return { token: loginRes.json('access_token') };
}

export default function (data) {
    const params = {
        headers: {
            'Authorization': `Bearer ${data.token}`,
            'Content-Type': 'application/json',
        },
    };

    // Hit the WBS endpoint which loads the tree
    const res = http.get(`${BASE_URL}/projects/${PROJECT_ID}/wbs`, params);

    check(res, {
        'status is 200': (r) => r.status === 200,
        'returns wbs data': (r) => r.json('data') !== undefined,
    });

    sleep(1); // Standard 1s sleep to simulate regular browser fetching
}
