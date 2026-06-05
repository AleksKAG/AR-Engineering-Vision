CREATE TABLE issues (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    element_id uuid REFERENCES elements(id),
    photo_url text NOT NULL,
    comment text,
    deviation_mm float, -- Смещение в мм (например, 72)
    status text DEFAULT 'open', -- open, fixed
    created_at timestamptz DEFAULT now()
);
