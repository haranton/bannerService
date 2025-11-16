CREATE TABLE IF NOT EXISTS banners (
     id SERIAL PRIMARY KEY,
     content jsonb,
     is_active boolean,
     created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
     updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE feature_tag_banner (
    tag_id BIGINT NOT NULL,
    feature_id BIGINT NOT NULL,
    banner_id BIGINT NOT NULL REFERENCES banners(id) ON DELETE CASCADE,
    PRIMARY KEY (tag_id, feature_id)
);

