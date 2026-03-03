SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `pitstop_authorisations`;

CREATE TABLE IF NOT EXISTS pitstop_authorisations (
    pitstop_auth_id VARCHAR(50) PRIMARY KEY,
    dataset_id VARCHAR(50) NOT NULL,
    dataset_name VARCHAR(255) NOT NULL,
    user_id VARCHAR(50),
    regulator_id CHAR(36) NOT NULL,
    regulator_name VARCHAR(255) NOT NULL,
    maincon_id CHAR(36) NOT NULL,
    maincon_name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
    last_synced_at DATETIME DEFAULT NULL
);