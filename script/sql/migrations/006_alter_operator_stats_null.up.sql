ALTER TABLE operator_stats MODIFY COLUMN operator_id varchar(45) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' NULL COMMENT 'Unique identifier for each driver, generated randomly each month';
ALTER TABLE operator_stats MODIFY COLUMN gender int(11) NULL COMMENT 'Driver gender: 1 = male; 2 = female; 3 = other; 4 = prefer not to answer';
ALTER TABLE operator_stats MODIFY COLUMN completed_trips int(11) NULL COMMENT 'Number of rides completed in the reported month';
ALTER TABLE operator_stats MODIFY COLUMN days_since int(11) NULL COMMENT 'Time in days, since the driver had the possibility to accept a ride for the first time';
ALTER TABLE operator_stats MODIFY COLUMN age_range varchar(60) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' NULL COMMENT 'Age (per five-year group)';
ALTER TABLE operator_stats MODIFY COLUMN hours_connected varchar(60) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' NULL COMMENT 'Hours that the driver was connected to the application during the reported month (in 25-hour ranges)';
ALTER TABLE operator_stats MODIFY COLUMN trip_hours varchar(60) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' NULL COMMENT 'Hours the driver spent with passengers on board and with ride requests during the reported month (in 25-hour ranges)';
ALTER TABLE operator_stats MODIFY COLUMN tot_revenue varchar(60) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' NULL COMMENT 'Total revenue for rides, including bonuses and tips (in ranges of $ 1,000)';
