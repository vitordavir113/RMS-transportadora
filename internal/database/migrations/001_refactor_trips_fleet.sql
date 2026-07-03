ALTER TABLE trips
ALTER COLUMN truck_id DROP NOT NULL;

ALTER TABLE trip_compartments
ALTER COLUMN truck_compartment_id DROP NOT NULL;

ALTER TABLE trips
ADD COLUMN IF NOT EXISTS tractor_id BIGINT;

ALTER TABLE trips
ADD COLUMN IF NOT EXISTS trailer_id BIGINT;

ALTER TABLE trips
ADD COLUMN IF NOT EXISTS driver_id BIGINT;

ALTER TABLE trips
ADD COLUMN IF NOT EXISTS tractor_plate_snapshot VARCHAR(10);

ALTER TABLE trips
ADD COLUMN IF NOT EXISTS trailer_plate_snapshot VARCHAR(10);

ALTER TABLE trips
ADD COLUMN IF NOT EXISTS driver_name_snapshot VARCHAR(150);

ALTER TABLE trip_compartments
ADD COLUMN IF NOT EXISTS trailer_compartment_id BIGINT;

ALTER TABLE trip_compartments
ADD COLUMN IF NOT EXISTS client_id BIGINT;

ALTER TABLE trip_compartments
ADD COLUMN IF NOT EXISTS client_name_snapshot VARCHAR(150);

ALTER TABLE trip_compartments
ADD COLUMN IF NOT EXISTS freight_value_snapshot DOUBLE PRECISION DEFAULT 0;

ALTER TABLE trip_compartments
ADD COLUMN IF NOT EXISTS freight_type_snapshot VARCHAR(20);

ALTER TABLE trip_compartments
ADD COLUMN IF NOT EXISTS freight_total DOUBLE PRECISION DEFAULT 0;

CREATE INDEX IF NOT EXISTS idx_trips_tractor_id ON trips (tractor_id);
CREATE INDEX IF NOT EXISTS idx_trips_trailer_id ON trips (trailer_id);
CREATE INDEX IF NOT EXISTS idx_trips_driver_id ON trips (driver_id);

CREATE INDEX IF NOT EXISTS idx_trip_compartments_trailer_compartment_id ON trip_compartments (trailer_compartment_id);
CREATE INDEX IF NOT EXISTS idx_trip_compartments_client_id ON trip_compartments (client_id);