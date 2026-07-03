DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'trip_compartments'
          AND column_name = 'cliente'
    ) THEN
        ALTER TABLE trip_compartments ALTER COLUMN cliente DROP NOT NULL;
    END IF;
END $$;

DO $$
BEGIN
    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_name = 'trip_compartments'
          AND column_name = 'produto'
    ) THEN
        ALTER TABLE trip_compartments ALTER COLUMN produto DROP NOT NULL;
    END IF;
END $$;