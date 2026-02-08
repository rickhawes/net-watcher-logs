-- This SQL script initializes the PostgreSQL database for the net_watcher_log application.
-- 
-- It is meant to be run automatically when the PostgreSQL container is first started, 
-- as specified in the docker-compose.yml file.
CREATE DATABASE net_watcher_logs;

GRANT ALL PRIVILEGES ON DATABASE net_watcher_logs TO postgres;

ALTER DATABASE net_watcher_logs SET timezone TO 'America/Los_Angeles';

\c net_watcher_logs

GRANT ALL ON SCHEMA public TO postgres;
