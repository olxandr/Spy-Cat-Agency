CREATE TABLE targets (
    id SERIAL PRIMARY KEY,
    mission_id INT NOT NULL REFERENCES missions(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    country VARCHAR(100) NOT NULL,
    notes TEXT,
    is_completed BOOLEAN NOT NULL DEFAULT FALSE,
    CONSTRAINT target_unique_per_mission UNIQUE (mission_id, name)
);

