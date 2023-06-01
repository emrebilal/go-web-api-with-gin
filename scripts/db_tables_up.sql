CREATE TABLE IF NOT EXISTS ratings
(
    id           bigserial
        CONSTRAINT ratings_pk
        PRIMARY KEY,
    username     varchar(36) NOT NULL,
    provider_id  varchar(32) NOT NULL,
    service_id   varchar(32) NOT NULL,
    rate         int         NOT NULL,
    created_date timestamp
);

CREATE UNIQUE INDEX IF NOT EXISTS uix_ratings_service_id
    ON ratings (service_id);
