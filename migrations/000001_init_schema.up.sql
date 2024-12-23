-- Create Food Table
CREATE TABLE foods
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL,
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

-- Create Group Table
CREATE TABLE groups
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       TEXT NOT NULL,
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

-- Create FoodGroup Table
CREATE TABLE foods_groups
(
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    food_id    UUID REFERENCES foods (id),
    group_id   UUID REFERENCES groups (id),
    max_size   FLOAT NOT NULL,
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);
