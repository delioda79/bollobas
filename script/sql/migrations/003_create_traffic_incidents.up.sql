CREATE TABLE IF NOT EXISTS `traffic_incidents` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `date` datetime NOT NULL,
  `type` int(11) NOT NULL COMMENT 'Type of traffic event that was reported',
  `plates` varchar(60) NOT NULL DEFAULT '' COMMENT 'Vehicle’s Plate Number',
  `licence` varchar(60) NOT NULL DEFAULT '' COMMENT 'Driver’s License Number',
  `travel_distance` varchar(60) NOT NULL DEFAULT '' COMMENT 'Travel distance in meters (3km ranges)',
  `travel_time` varchar(60) NOT NULL DEFAULT '' COMMENT 'Travel time in minutes (5 minute ranges)',
  `coordinates` varchar(40) NOT NULL DEFAULT '' COMMENT 'Coordinates (longitude and latitude) where the incident occurred',
  PRIMARY KEY (`id`),
  KEY `date` (`date`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;