ALTER TABLE aggregated_trips ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW() AFTER eod_empty_time;
ALTER TABLE aggregated_trips ADD COLUMN produced_at DATETIME NOT NULL DEFAULT NOW() AFTER created_at;

ALTER TABLE operator_stats ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW() AFTER tot_revenue;
ALTER TABLE operator_stats ADD COLUMN produced_at DATETIME NOT NULL DEFAULT NOW() AFTER created_at;

ALTER TABLE traffic_incidents ADD COLUMN created_at TIMESTAMP NOT NULL DEFAULT NOW() AFTER coordinates;
ALTER TABLE traffic_incidents ADD COLUMN produced_at DATETIME NOT NULL DEFAULT NOW() AFTER created_at;