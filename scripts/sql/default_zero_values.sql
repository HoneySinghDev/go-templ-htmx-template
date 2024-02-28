CREATE OR REPLACE FUNCTION report_default_go_sql_zero_values_mismatch()
    RETURNS SETOF information_schema.columns
    AS $$
BEGIN
    RETURN QUERY
    SELECT
        *
    FROM
        information_schema.columns
    WHERE (table_schema = 'public'
        AND column_default IS NOT NULL)
        AND (
            (data_type = 'boolean' AND column_default <> 'false'::text AND column_default <> '''false'''::text)
            OR (data_type IN ('char', 'character', 'varchar', 'character varying', 'text')
                AND column_default NOT LIKE '''%''')
            OR (data_type IN ('smallint', 'integer', 'bigint', 'smallserial', 'serial', 'bigserial')
                AND column_default <> '0'::text AND column_default NOT LIKE 'nextval(%'::text)
            OR (data_type IN ('decimal', 'numeric', 'real', 'double precision')
                AND column_default <> '0.0'::text)
        );
END
$$
LANGUAGE plpgsql
SECURITY DEFINER;
