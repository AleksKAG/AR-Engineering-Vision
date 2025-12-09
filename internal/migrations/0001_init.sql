-- USERS
CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    email text UNIQUE NOT NULL,
    password_hash text NOT NULL,
    created_at timestamptz DEFAULT now()
);

-- PROJECTS
CREATE TABLE projects (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id uuid REFERENCES users(id),
    name text NOT NULL,
    description text,
    created_at timestamptz DEFAULT now()
);

-- MODELS
CREATE TABLE models (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id uuid REFERENCES projects(id),
    filename text,
    content_type text,
    s3_key text,
    created_at timestamptz DEFAULT now()
);

-- ROOMS
CREATE TABLE rooms (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id uuid REFERENCES projects(id),
    name text,
    bbox jsonb,
    created_at timestamptz DEFAULT now()
);

-- ELEMENTS
CREATE TABLE elements (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id uuid REFERENCES rooms(id),
    type text,
    world_coords jsonb,
    properties jsonb,
    created_at timestamptz DEFAULT now()
);
