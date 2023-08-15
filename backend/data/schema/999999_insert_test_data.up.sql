
-- Users

-- email = 'test@example.com', password = 'pass'
INSERT INTO users(id, email, password, created_at, updated_at)
VALUES('test', 'test@example.com', '$2a$10$YfHxWNfL8Ba2ltl6TRHMVuN0WPXxAuB5L7w1Y0jqaFcn2bUDoUq9W', NOW(), NOW());

-- email = 'dev@example.com', password = 'pass'
INSERT INTO users(id, email, password, created_at, updated_at)
VALUES('dev', 'dev@example.com', '$2a$10$YfHxWNfL8Ba2ltl6TRHMVuN0WPXxAuB5L7w1Y0jqaFcn2bUDoUq9W', NOW(), NOW());

-- email = 'admin@example.com', password = 'pass'
INSERT INTO users(id, email, password, created_at, updated_at)
VALUES('admin', 'admin@example.com', '$2a$10$YfHxWNfL8Ba2ltl6TRHMVuN0WPXxAuB5L7w1Y0jqaFcn2bUDoUq9W', NOW(), NOW());

-- Tasks

INSERT INTO tasks(id, user_id, name, is_completed, created_at, updated_at)
VALUES('t1', 'test', 'Test Task 1', false, NOW(), NOW());

INSERT INTO tasks(id, user_id, name, is_completed, created_at, updated_at)
VALUES('t2', 'test', 'Test Task 2', true, NOW(), NOW());

INSERT INTO tasks(id, user_id, name, is_completed, created_at, updated_at)
VALUES('t3', 'dev', 'Dev Task 1', false, NOW(), NOW());

INSERT INTO tasks(id, user_id, name, is_completed, created_at, updated_at)
VALUES('t4', 'admin', 'Admin Task 1', false, NOW(), NOW());
