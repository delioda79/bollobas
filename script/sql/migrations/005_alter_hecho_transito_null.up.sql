ALTER TABLE traffic_incidents MODIFY COLUMN `type` int(11) NULL COMMENT 'Type of traffic event that was reported';
ALTER TABLE traffic_incidents MODIFY COLUMN plates varchar(60) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' NULL COMMENT 'Vehicle’s Plate Number';
ALTER TABLE traffic_incidents MODIFY COLUMN licence varchar(60) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' NULL COMMENT 'Driver’s License Number';
ALTER TABLE traffic_incidents MODIFY COLUMN travel_distance varchar(60) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' NULL COMMENT 'Travel distance in meters (3km ranges)';
ALTER TABLE traffic_incidents MODIFY COLUMN travel_time varchar(60) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' NULL COMMENT 'Travel time in minutes (5 minute ranges)';
ALTER TABLE traffic_incidents MODIFY COLUMN coordinates varchar(40) CHARACTER SET utf8 COLLATE utf8_general_ci DEFAULT '' NULL COMMENT 'Coordinates (longitude and latitude) where the incident occurred';
