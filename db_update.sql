ALTER TABLE `meters` ADD `current_reading` DECIMAL(10,2) NOT NULL DEFAULT '0.00' AFTER `valve_status`;
ALTER TABLE `meters` ADD `previous_reading` DECIMAL(10,2) NOT NULL DEFAULT '0.00' AFTER `valve_status`;
ALTER TABLE `meters` ADD `meter_type` VARCHAR(50) NOT NULL DEFAULT 'Prepaid Meter' AFTER `meter_description`;