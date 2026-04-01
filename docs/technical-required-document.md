TRD — IAM System Core

1. Document Metadata
TRD ID: TRD-IAM-Core
Status: Draft
Owner: Backend Team
Last Updated: 2025-12-20
Related Playbook: ai/playbooks/build_feature_from_trd.md
Purpose: Establish traceability for IAM features including User Management, Role Assignment, and Master Data.

2. Purpose & Context
2.1 Background
The system requires a local identity and access management (IAM) layer that integrates with a Central authority while supporting local role assignments and group structures. This feature consolidates user synchronization, role management, and context-aware data access (institutions/groups).

2.2 Goals
Goals:
- Enable user synchronization from Central to local `auth.users`.
- Support comprehensive Role-Based Access Control (RBAC) via `iam.roles` and `iam.user_roles`.
- Provide master data endpoints for Permissions, Institutions, and Groups to support UI flows.

2.3 Non-Goals (Critical)
Non-Goals:
- No local authentication (password management handled by Central/IDP).
- No creation of Institutions (read-only from Central/Upstream).
- No complexity in Group hierarchies beyond what is defined in `iam.groups`.

3. Scope Definition
3.1 In Scope
In Scope:
- User listing, syncing (POST), and status updates.
- Role creation, update, deletion, and assignment to users.
- Permission listing and Group listing.
- Integration with existing `auth` and `iam` schemas.

3.2 Out of Scope
Out of Scope:
- User profile updates (must be done in Central).
- Complex multi-tenancy logic beyond `institution_id` filtering.

4. Domain Model
4.1 Entities

Entity: User
Fields:
- id (uuid, PK)
- external_subject (string, unique, sync from Central)
- institution_id (uuid, FK)
- status (string)
- metadata (jsonb, contains profile info)
- created_at (timestamptz)

Entity: Role
Fields:
- id (bigint, PK)
- institution_id (uuid, FK)
- name (string)
- description (string)
- is_active (boolean)

Entity: UserRole
Fields:
- id (bigint, PK)
- user_id (uuid, FK)
- role_id (bigint, FK)
- group_id (bigint, FK, NOT NULL)
- institution_id (uuid, FK)

4.2 Domain Rules
Rules:
- User `external_subject` must be unique.
- Role names must be unique per `institution_id`.
- `group_id` is mandatory for Role Assignment; default group must be resolved if missing.
- Users cannot be updated locally except for `status`.

5. Usecases

5.1 Usecase List
Usecases:
- ListUsers
- SyncUser
- GetUser
- UpdateUserStatus
- ListRoles
- CreateRole
- UpdateRole
- DeleteRole
- AssignUserRole
- RevokeUserRole
- ListUserRoles
- ListGroups
- ListPermissions
- ListInstitutions

5.2 Usecase Specifications

Usecase: SyncUser
Input:
- code (string, from Central)
Steps:
1. Fetch user details from Central API using `code`.
2. Map Central response to `auth.users` schema (extract metadata).
3. Upsert user into `auth.users` using `external_subject`.
Output:
- user_id

Usecase: AssignUserRole
Input:
- user_id
- role_id
- institution_id
- group_id (optional, logic to resolve default)
Steps:
1. Validate `user_id` and `role_id` exist.
2. Resolve `group_id` if missing (Logic: User's default group or Institution's default).
3. Insert into `iam.user_roles`.
Output:
- assignment_id

6. Interfaces

6.1 HTTP APIs

POST /users
Request:
{
  "code": "EMP001"
}
Response: 201
{
  "id": "uuid..."
}

GET /users
Query: institution_id, status, search
Response: 200
[ { "id": "...", "metadata": { "full_name": "..." } } ]

PATCH /users/{id}/status
Request:
{
  "status": "inactive"
}

POST /users/{id}/roles
Request:
{
  "role_id": 123,
  "institution_id": "uuid...",
  "group_id": 456
}

GET /roles
Query: institution_id
Response: 200
[ { "id": 123, "name": "Admin" } ]

POST /roles
Request:
{
  "institution_id": "uuid...",
  "name": "Supervisor",
  "permissions": ["module.read", "module.write"]
}

6.2 Events (Publishers)
Event: local.user.synced
Exchange: local.user
Routing Key: user.synced
Payload: { "user_id": "..." }

6.3 Consumers (Optional)
None specified in source contract.

6.4 Cron Jobs (Optional)
None specified.

7. Data & Persistence
7.1 Tables
Table: auth.users
PK: id
Indexes: (institution_id, status), (external_subject)

Table: iam.roles
PK: id
Indexes: (institution_id)

Table: iam.user_roles
PK: id
Indexes: (user_id), (role_id)

Table: iam.role_permissions
PK: (role_id, permission_code)

7.2 Migrations
Migration Required: NO (Assumes existing schema based on contract reference to "Existing database schema")

8. Non-Functional Requirements
Performance:
- User list must support pagination and efficient JSONB searching.
Reliability:
- Central sync must handle timeouts gracefully.
Security:
- All writes require proper RBAC checks (e.g., `user.write`, `role.write`).

9. Observability
Logging:
- Log syncing events with Central.
- Log Role assignment changes.

10. Testing Requirements
Unit Tests:
- Metadata extraction logic.
- Role assignment validation (group_id resolution).
Integration Tests:
- Full SyncUser flow (mocking Central).
- Role CRUD operations.

11. Open Questions / Decisions
Open:
- How to resolve `group_id` if not provided in UI? (Proposed: default to a known "General" group item).
- Is `generate-username` endpoint purely a proxy or does it need local state? (Assumed: Proxy only).

12. TRD → AI Mapping
Feature Parts:
- HTTP: SyncUser → add_http_endpoint
- HTTP: RoleCRUD → add_http_endpoint
- HTTP: UserRoleAssignment → add_http_endpoint

Required Rules:
- ai/rules/20_usecase_service.md
- ai/rules/30_handler_http.md
- ai/rules/10_repository.md

Templates Expected:
- ai/templates/usecase.md
- ai/templates/handler_http.md
- ai/templates/repository.md

Checklist:
- ai/checklist/feature_done.md
