CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    url VARCHAR(255) NOT NULL
);

CREATE INDEX url_idx on links (url);
