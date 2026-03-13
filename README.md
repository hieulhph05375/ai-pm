# Project Management System (PM-System)

Enterprise-grade Project Management System built with Go (Backend) and SvelteKit (Frontend).

## Key Features

- **Hierarchical WBS**: Manage complex project structures with LTree (Depth up to 5 levels).
- **Interactive Gantt Chart**: Built-in scheduling and visual progress tracking.
- **Dynamic RBAC**: Flexible role-based access control at both system and project levels.
- **Performance Optimized**: Handles 10k+ task nodes with sub-500ms latency.
- **Production Ready**: Fully dockerized stack for seamless deployment.

## Repository

GitHub: [https://github.com/caikeoboompro/project_task_mgmt](https://github.com/caikeoboompro/project_task_mgmt)

## Installation & Deployment

For production setup using Docker, please refer to the detailed guide in the `installation` directory:

[Installation Guide (INSTALL.md)](installation/INSTALL.md)

### Quick Start (Development)

1. **Backend**:

   ```bash
   cd backend
   go run cmd/server/main.go
   ```

2. **Frontend**:

   ```bash
   cd frontend
   npm run dev
   ```

## Documentation

- [Roadmap](.gsd/ROADMAP.md) - Project progress and upcoming features.
- [Decisions](.gsd/DECISIONS.md) - Technical architecture and design choices.

---
*Developed under GSD Methodology.*
