CREATE TABLE IF NOT EXISTS "songs"  (
    "id" SERIAL NOT NULL UNIQUE,
    "song" VARCHAR(255) NOT NULL,
    "group" VARCHAR(255) NOT NULL,
    "release_date" DATE,
    "text" TEXT,
    "link" VARCHAR(255),
    "created_at" TIMESTAMP NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP NOT NULL DEFAULT NOW()
)