SET FOREIGN_KEY_CHECKS = 0;

-- ====================================================
-- 1. Drop FK constraints from `projects` table
-- ====================================================
ALTER TABLE `projects`
DROP FOREIGN KEY IF EXISTS `fk_projects_main_contractor`;

ALTER TABLE `projects`
DROP FOREIGN KEY IF EXISTS `fk_projects_offsite_fabricator`;

ALTER TABLE `projects`
DROP FOREIGN KEY IF EXISTS `fk_projects_worker_company`;

ALTER TABLE `projects`
DROP FOREIGN KEY IF EXISTS `fk_projects_worker_company_client`;

-- 2. Drop old ID columns and add inline text columns
ALTER TABLE `projects`
DROP COLUMN IF EXISTS `main_contractor_id`,
DROP COLUMN IF EXISTS `offsite_fabricator_id`,
DROP COLUMN IF EXISTS `worker_company_id`,
DROP COLUMN IF EXISTS `worker_company_client_id`,
ADD COLUMN `main_contractor_name` varchar(255) DEFAULT NULL,
ADD COLUMN `main_contractor_uen` varchar(50) DEFAULT NULL,
ADD COLUMN `offsite_fabricator_name` varchar(255) DEFAULT NULL,
ADD COLUMN `offsite_fabricator_uen` varchar(50) DEFAULT NULL,
ADD COLUMN `offsite_fabricator_location` varchar(255) DEFAULT NULL,
ADD COLUMN `worker_company_name` varchar(255) DEFAULT NULL,
ADD COLUMN `worker_company_uen` varchar(50) DEFAULT NULL,
ADD COLUMN `worker_company_client_name` varchar(255) DEFAULT NULL,
ADD COLUMN `worker_company_client_uen` varchar(50) DEFAULT NULL;

-- ====================================================
-- 3. Drop FK constraint from `users` table
-- ====================================================
ALTER TABLE `users` DROP FOREIGN KEY IF EXISTS `users_ibfk_2`;

ALTER TABLE `users` DROP KEY IF EXISTS `company_id`;

ALTER TABLE `users`
DROP COLUMN IF EXISTS `company_id`,
ADD COLUMN `company_name` varchar(255) DEFAULT NULL;

-- ====================================================
-- 4. Drop the companies table
-- ====================================================
DROP TABLE IF EXISTS `companies`;

SET FOREIGN_KEY_CHECKS = 1;