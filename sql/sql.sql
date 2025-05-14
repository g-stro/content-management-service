SET TIMEZONE = 'UTC';

CREATE TABLE "content"
(
    "id"                 SERIAL PRIMARY KEY,
    "name" VARCHAR(255),
    "description"        TEXT,
    "creation_date"      TIMESTAMP,
    "last_modified_date" TIMESTAMP
);

CREATE TABLE "content_type"
(
    "id"   SERIAL PRIMARY KEY,
    "name" VARCHAR(50)
);

CREATE TABLE "content_details"
(
    "id"              SERIAL PRIMARY KEY,
    "content_id"      INTEGER,
    "content_type_id" INTEGER,
    "value"           TEXT,
    FOREIGN KEY ("content_id") REFERENCES "content" ("id"),
    FOREIGN KEY ("content_type_id") REFERENCES "content_type" ("id")
);

INSERT INTO content_type (id, name)
VALUES
    (1, 'text'),
    (2, 'image'),
    (3, 'video')
ON CONFLICT DO NOTHING;
