-- DUMMY SEEDER
-- Populates the database with initial data for testing.
-- Usage: psql -d <database> -f seeder.sql

DO $$
DECLARE
    -- IDs
    uni_tech_id UUID;
    uni_health_id UUID;
    
    admin_user_id UUID;
    staff_user_id UUID;
    student_user_id UUID;
    
    role_admin_id UUID;
    role_staff_id UUID;
    
    perm_user_view_id UUID;
    perm_user_edit_id UUID;
    perm_roles_view_id UUID;
    
    group_it_id UUID;
    group_med_id UUID;
BEGIN
    -- =================================================================================================
    -- 1. INSTITUTIONS
    -- =================================================================================================
    RAISE NOTICE 'Seeding Institutions...';

    INSERT INTO auth.institutions (code, name, description, is_active)
    VALUES ('TECH-UNI', 'Tech University', 'Top technical university', true)
    RETURNING id INTO uni_tech_id;
    
    INSERT INTO auth.institutions (code, name, description, is_active)
    VALUES ('HEALTH-INS', 'Health Institute', 'Leading health academy', true)
    RETURNING id INTO uni_health_id;

    -- =================================================================================================
    -- 2. USERS (Tech Uni)
    -- =================================================================================================
    RAISE NOTICE 'Seeding Users...';

    -- Admin User
    INSERT INTO auth.users (institution_id, external_subject, identity_provider, metadata, status)
    VALUES 
    (uni_tech_id, 'sub_admin_tech_001', 'central', '{"full_name": "Tech Admin", "email": "admin@tech.edu"}'::JSONB, 'active')
    RETURNING id INTO admin_user_id;

    -- Staff User
    INSERT INTO auth.users (institution_id, external_subject, identity_provider, metadata, status)
    VALUES 
    (uni_tech_id, 'sub_staff_tech_001', 'central', '{"full_name": "Tech Staff", "email": "staff@tech.edu"}'::JSONB, 'active')
    RETURNING id INTO staff_user_id;

    -- Student User
    INSERT INTO auth.users (institution_id, external_subject, identity_provider, metadata, status)
    VALUES 
    (uni_tech_id, 'sub_student_tech_001', 'central', '{"full_name": "Tech Student", "email": "student@tech.edu"}'::JSONB, 'active')
    RETURNING id INTO student_user_id;

    -- =================================================================================================
    -- 3. ROLES (Tech Uni)
    -- =================================================================================================
    RAISE NOTICE 'Seeding Roles...';

    -- Super Admin Role
    INSERT INTO iam.roles (institution_id, name, display_name, description, is_system_role)
    VALUES (uni_tech_id, 'super_admin', 'Super Administrator', 'Full system access', true)
    RETURNING id INTO role_admin_id;

    -- Academic Staff Role
    INSERT INTO iam.roles (institution_id, name, display_name, description, is_system_role)
    VALUES (uni_tech_id, 'academic_staff', 'Academic Staff', 'Standard staff permissions', false)
    RETURNING id INTO role_staff_id;

    -- =================================================================================================
    -- 4. PERMISSIONS (Tech Uni)
    -- =================================================================================================
    RAISE NOTICE 'Seeding Permissions...';

    -- Format: module.sub_module.page.action
    INSERT INTO iam.permissions (institution_id, code, module, sub_module, page, action, scope_type)
    VALUES 
    (uni_tech_id, 'users.manage.all.view', 'users', 'manage', 'all', 'view', 'ui')
    RETURNING id INTO perm_user_view_id;

    INSERT INTO iam.permissions (institution_id, code, module, sub_module, page, action, scope_type)
    VALUES 
    (uni_tech_id, 'users.manage.all.edit', 'users', 'manage', 'all', 'edit', 'api')
    RETURNING id INTO perm_user_edit_id;

    INSERT INTO iam.permissions (institution_id, code, module, sub_module, page, action, scope_type)
    VALUES 
    (uni_tech_id, 'roles.manage.all.view', 'roles', 'manage', 'all', 'view', 'ui')
    RETURNING id INTO perm_roles_view_id;

    -- IDs are now captured directly via RETURNING clauses above

    -- =================================================================================================
    -- 5. ROLE PERMISSIONS
    -- =================================================================================================
    RAISE NOTICE 'Assigning Permissions to Roles...';

    -- Admin gets all
    INSERT INTO iam.role_permissions (role_id, permission_id)
    VALUES 
    (role_admin_id, perm_user_view_id),
    (role_admin_id, perm_user_edit_id),
    (role_admin_id, perm_roles_view_id);

    -- Staff gets view only
    INSERT INTO iam.role_permissions (role_id, permission_id)
    VALUES 
    (role_staff_id, perm_user_view_id);

    -- =================================================================================================
    -- 6. GROUPS (Tech Uni)
    -- =================================================================================================
    RAISE NOTICE 'Seeding Groups...';

    INSERT INTO iam.groups (institution_id, name, group_type, description)
    VALUES (uni_tech_id, 'IT Department', 'department', 'IT Services and Support')
    RETURNING id INTO group_it_id;

    -- =================================================================================================
    -- 7. USER ROLES (Tech Uni)
    -- =================================================================================================
    RAISE NOTICE 'Assigning Roles to Users...';

    -- Admin -> Super Admin Role -> IT Group
    INSERT INTO iam.user_roles (institution_id, user_id, role_id, group_id)
    VALUES (uni_tech_id, admin_user_id, role_admin_id, group_it_id);

    -- Staff -> Academic Staff Role -> IT Group
    INSERT INTO iam.user_roles (institution_id, user_id, role_id, group_id)
    VALUES (uni_tech_id, staff_user_id, role_staff_id, group_it_id);

    RAISE NOTICE 'Seeding Completed Successfully!';
    RAISE NOTICE 'Institution ID: %', uni_tech_id;
    RAISE NOTICE 'Admin User ID: %', admin_user_id;

END $$;
