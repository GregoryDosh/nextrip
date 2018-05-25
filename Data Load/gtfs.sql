BEGIN;

CREATE SCHEMA mt;

CREATE TABLE mt.stops
(
    stop_id int primary key,
    stop_code text,
    stop_name text,
    stop_desc text,
    stop_lat double precision,
    stop_lon double precision,
    zone_id text,
    stop_url text,
    location_type int,
    wheelchair_boarding int
);

COPY mt.stops FROM PROGRAM 'wget ftp://ftp.gisdata.mn.gov/pub/gdrs/data/pub/us_mn_state_metc/trans_transit_schedule_google_fd/csv_trans_transit_schedule_google_fd.zip 2> /dev/null 1>&2 && unzip csv_trans_transit_schedule_google_fd.zip 2>/dev/null 1>&2 && cat stops.txt' CSV HEADER;

END;