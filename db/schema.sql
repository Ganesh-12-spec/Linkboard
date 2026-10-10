CREATE TABLE bookmarks (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    title VARCHAR(255) NOT NULL
        CHECK (length(trim(title)) > 0),

    url TEXT NOT NULL
        CHECK (length(trim(url)) > 0),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
