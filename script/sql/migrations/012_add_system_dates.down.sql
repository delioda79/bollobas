ALTER TABLE operator_stats DROP COLUMN created_at;
ALTER TABLE operator_stats DROP COLUMN produced_at;

ALTER TABLE aggregated_trips DROP COLUMN created_at;
ALTER TABLE aggregated_trips DROP COLUMN produced_at;

ALTER TABLE traffic_incidents DROP COLUMN created_at;
ALTER TABLE traffic_incidents DROP COLUMN produced_at;