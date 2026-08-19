CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE tags (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    slug TEXT NOT NULL UNIQUE
);

CREATE TABLE wallpapers (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    slug TEXT NOT NULL UNIQUE,
    image_url TEXT NOT NULL,
    thumbnail_url TEXT NOT NULL,
    width INTEGER NOT NULL CHECK (width > 0),
    height INTEGER NOT NULL CHECK (height > 0),
    file_size BIGINT NOT NULL CHECK (file_size >= 0),
    downloads BIGINT NOT NULL DEFAULT 0 CHECK (downloads >= 0),
    views BIGINT NOT NULL DEFAULT 0 CHECK (views >= 0),
    category_id BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE wallpaper_tags (
    wallpaper_id BIGINT NOT NULL REFERENCES wallpapers(id) ON DELETE CASCADE,
    tag_id BIGINT NOT NULL REFERENCES tags(id) ON DELETE CASCADE,
    PRIMARY KEY (wallpaper_id, tag_id)
);

CREATE INDEX wallpapers_category_id_idx ON wallpapers(category_id);
CREATE INDEX wallpapers_created_at_idx ON wallpapers(created_at DESC);
CREATE INDEX wallpaper_tags_tag_id_idx ON wallpaper_tags(tag_id);

