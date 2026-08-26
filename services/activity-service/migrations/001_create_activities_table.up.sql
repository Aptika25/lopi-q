CREATE TABLE IF NOT EXISTS tb_lopiq_daily_activities (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    user_nip VARCHAR(50) NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    activity_date VARCHAR(50) NOT NULL,
    photo_url TEXT DEFAULT '',
    status VARCHAR(20) DEFAULT 'PENDING',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_activities_user_id ON tb_lopiq_daily_activities(user_id);
CREATE INDEX IF NOT EXISTS idx_activities_status ON tb_lopiq_daily_activities(status);
