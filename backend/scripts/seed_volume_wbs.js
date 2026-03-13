const { Client } = require('pg');
const crypto = require('crypto');

// Target the local pm_db database
const connectionString = 'postgresql://postgres:Caikeo@1234@localhost:5432/pm_db?sslmode=disable';

async function seedVolume() {
    console.log('🚀 Starting Volume WBS Seeder (10,000 Nodes)...');

    const client = new Client({ connectionString });

    try {
        await client.connect();

        // 1. Create a dummy project for load testing
        console.log('📦 Creating dummy project...');
        const projectInsert = `
            INSERT INTO projects (project_id, project_name, description, project_status, current_phase, planned_start_date, planned_end_date)
            VALUES ($1, $2, $3, $4, $5, $6, $7)
            RETURNING id;
        `;
        const projectValues = [
            'PRJ-10K-' + Date.now(),
            'Load Test Project 10K',
            'A generated project containing 10,000 WBS nodes to stress test the Ltree indexes and hierarchy retrieval.',
            'Running',
            'Execution',
            '2026-01-01',
            '2026-12-31'
        ];
        const resProject = await client.query(projectInsert, projectValues);
        const projectId = resProject.rows[0].id;
        console.log(`✅ Dummy Project created with ID: ${projectId}`);

        // 2. Clear existing nodes for this project just in case
        await client.query('DELETE FROM wbs_nodes WHERE project_id = $1', [projectId]);

        // 3. Get the Max ID to avoid duplicate key issues
        const maxIdRes = await client.query('SELECT COALESCE(MAX(id), 0) as max_id FROM wbs_nodes');
        let nextId = parseInt(maxIdRes.rows[0].max_id, 10) + 1;

        console.log(`🌳 Generating 10,000 WBS Node hierarchies starting from ID ${nextId}...`);

        let nodes = [];
        let paths = [];

        // Level 1: 10 main phases
        for (let i = 1; i <= 10; i++) {
            const phaseId = nextId++;
            const phasePath = phaseId.toString();
            paths.push(phasePath);
            nodes.push({ id: phaseId, path: phasePath, title: `Phase ${i}`, type: 'Phase' });

            // Level 2: 10 sub-phases per phase (10 * 10 = 100)
            for (let j = 1; j <= 10; j++) {
                const subPhaseId = nextId++;
                const subPhasePath = `${phasePath}.${subPhaseId}`;
                paths.push(subPhasePath);
                nodes.push({ id: subPhaseId, path: subPhasePath, title: `Sub-Phase ${i}.${j}`, type: 'Phase' });

                // Level 3: 10 tasks per sub-phase (100 * 10 = 1,000)
                for (let k = 1; k <= 10; k++) {
                    const taskId = nextId++;
                    const taskPath = `${subPhasePath}.${taskId}`;
                    paths.push(taskPath);
                    nodes.push({ id: taskId, path: taskPath, title: `Task ${i}.${j}.${k}`, type: 'Task' });

                    // Level 4: 9 sub-tasks per task (1,000 * 9 = 9,000)
                    for (let l = 1; l <= 9; l++) {
                        const subTaskId = nextId++;
                        const subTaskPath = `${taskPath}.${subTaskId}`;
                        paths.push(subTaskPath);
                        nodes.push({ id: subTaskId, path: subTaskPath, title: `Sub-Task ${i}.${j}.${k}.${l}`, type: 'Task' });
                    }
                }
            }
        }

        // Total Nodes: 10 (L1) + 100 (L2) + 1000 (L3) + 9000 (L4) = 10,110 nodes

        console.log(`💾 Inserting ${nodes.length} nodes using bulk insert...`);

        // We will insert in chunks of 1000 to avoid exceeding parameter limits
        const CHUNK_SIZE = 1000;
        for (let i = 0; i < nodes.length; i += CHUNK_SIZE) {
            const chunk = nodes.slice(i, i + CHUNK_SIZE);

            const values = [];
            const placeholders = [];
            let paramIndex = 1;

            chunk.forEach((node) => {
                placeholders.push(`($${paramIndex}, $${paramIndex + 1}, $${paramIndex + 2}, $${paramIndex + 3}, $${paramIndex + 4})`);
                values.push(projectId, node.id, node.title, node.type, node.path);
                paramIndex += 5;
            });

            const query = `
                INSERT INTO wbs_nodes (project_id, id, title, type, path)
                VALUES ${placeholders.join(', ')};
            `;

            await client.query(query, values);
            process.stdout.write(`\rInserted chunk ${Math.floor(i / CHUNK_SIZE) + 1}/${Math.ceil(nodes.length / CHUNK_SIZE)}`);
        }

        console.log('\n✨ Volume Seeding Complete! Project ID:', projectId);

    } catch (err) {
        console.error('\n❌ Error during seeding:', err);
    } finally {
        await client.end();
    }
}

seedVolume();
