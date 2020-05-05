CREATE TABLE IF NOT EXISTS `operator_stats` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `date` datetime NOT NULL,
  `operator_id` varchar(45) NOT NULL DEFAULT '' COMMENT 'Unique identifier for each driver, generated randomly each month',
  `gender` int(11) NOT NULL COMMENT 'Driver gender: 1 = male; 2 = female; 3 = other; 4 = prefer not to answer',
  `completed_trips` int(11) NOT NULL COMMENT 'Number of rides completed in the reported month',
  `days_since` int(11) NOT NULL COMMENT 'Time in days, since the driver had the possibility to accept a ride for the first time',
  `age_range` varchar(60) NOT NULL DEFAULT '' COMMENT 'Age (per five-year group)',
  `hours_connected` varchar(60) NOT NULL DEFAULT '' COMMENT 'Hours that the driver was connected to the application during the reported month (in 25-hour ranges)',
  `trip_hours` varchar(60) NOT NULL DEFAULT '' COMMENT 'Hours the driver spent with passengers on board and with ride requests during the reported month (in 25-hour ranges)',
  `tot_revenue` varchar(60) NOT NULL DEFAULT '' COMMENT 'Total revenue for rides, including bonuses and tips (in ranges of $ 1,000)',
  PRIMARY KEY (`id`),
  KEY `date` (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;