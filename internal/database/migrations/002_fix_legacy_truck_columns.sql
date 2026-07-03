DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'trips'
          AND column_name = 'truck_id'
    ) THEN
        ALTER TABLE trips ALTER COLUMN truck_id DROP NOT NULL;
    END IF;
END $$;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'trip_compartments'
          AND column_name = 'truck_compartment_id'
    ) THEN
        ALTER TABLE trip_compartments ALTER COLUMN truck_compartment_id DROP NOT NULL;
    END IF;
END $$;