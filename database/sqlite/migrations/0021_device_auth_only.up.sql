-- Device deployment schema: remove web user/role/menu authorization tables.
DROP TABLE IF EXISTS app_role_access;
DROP TABLE IF EXISTS app_user_role;
DROP TABLE IF EXISTS app_menu_permission;
DROP TABLE IF EXISTS app_menu;
DROP TABLE IF EXISTS app_role;
DROP TABLE IF EXISTS app_login;
DROP TABLE IF EXISTS app_login_attempt;
DROP TABLE IF EXISTS sessions;
DROP TABLE IF EXISTS app_user;

CREATE TABLE sessions (
    id TEXT PRIMARY KEY NOT NULL,
    user_id TEXT NOT NULL,
    access_token TEXT NOT NULL,
    access_token_expired_at DATETIME NOT NULL,
    refresh_token TEXT NOT NULL,
    refresh_token_expired_at DATETIME NOT NULL,
    ip_address TEXT NOT NULL,
    user_agent TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME
);

CREATE TABLE app_login (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    access_token TEXT,
    status TEXT CHECK (status IN ('login', 'logout')),
    ip_address TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    type TEXT DEFAULT 'api' CHECK (type IN ('api', 'web'))
);

CREATE TABLE app_login_attempt (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    pin TEXT,
    ip_address TEXT,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
