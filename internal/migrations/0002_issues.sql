CREATE TABLE issues (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    element_id uuid REFERENCES elements(id),
    photo_url text NOT NULL,
    comment text,
    deviation_mm float, -- Смещение в мм (например, 72)
    status text DEFAULT 'open', -- open, fixed
    created_at timestamptz DEFAULT now()
);

ALTER TABLE issues 
ADD COLUMN ai_detected_type text,
ADD COLUMN ai_confidence float,
ADD COLUMN is_match boolean;

-- Таблица для сгенерированных PDF-отчётов
CREATE TABLE reports (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id uuid REFERENCES projects(id),
    s3_key text NOT NULL,
    generated_at timestamptz DEFAULT now()
);
