-- +goose Up
-- +goose StatementBegin

CREATE TABLE banners (
     id SERIAL PRIMARY KEY,
     content JSONB NOT NULL,
     is_active BOOLEAN NOT NULL,
     created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE featuretag (
     feature_id INT NOT NULL,
     tag_id INT NOT NULL,
     banner_id INT NOT NULL,
     PRIMARY KEY (tag_id, feature_id),
     FOREIGN KEY (banner_id) REFERENCES banners(id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table banners;
drop table featuretag;
-- +goose StatementEnd

