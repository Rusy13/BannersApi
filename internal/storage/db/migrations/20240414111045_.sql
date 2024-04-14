-- +goose Up
-- +goose StatementBegin


CREATE TABLE banner_versions (
 id SERIAL PRIMARY KEY,
 banner_id INT NOT NULL,
 version_number INT NOT NULL,
 content JSONB NOT NULL,
 is_active BOOLEAN NOT NULL,
 created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table banners;
drop table featuretag;
-- +goose StatementEnd

