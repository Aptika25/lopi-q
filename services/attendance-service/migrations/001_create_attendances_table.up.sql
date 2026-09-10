CREATE TABLE IF NOT EXISTS tb_lopiq_attendances (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    user_nip VARCHAR(50) NOT NULL,
    user_name VARCHAR(150) NOT NULL,
    type VARCHAR(20) NOT NULL, -- 'MASUK' or 'PULANG'
    timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    distance_meters DOUBLE PRECISION NOT NULL DEFAULT 0,
    within_radius BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_attendances_user_id ON tb_lopiq_attendances(user_id);
CREATE INDEX IF NOT EXISTS idx_attendances_timestamp ON tb_lopiq_attendances(timestamp);
