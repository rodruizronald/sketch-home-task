DROP TABLE IF EXISTS canvas;

CREATE TABLE canvas (
    canvas_id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE,
    width INT,
    height INT,
    drawings JSONB
);

ALTER TABLE canvas
    OWNER TO postgres;