-- Migration 013: Phase out trade_code from workers table
-- trade_code is superseded by person_trade (SGBuildex compliance field)
-- The attendance table retains its own trade_code column (unrelated)

ALTER TABLE `workers` DROP COLUMN `trade_code`;