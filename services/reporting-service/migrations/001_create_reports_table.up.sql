CREATE TABLE IF NOT EXISTS tb_lopiq_attendance_reports (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    user_name VARCHAR(255) NOT NULL,
    user_nip VARCHAR(50) NOT NULL,
    department VARCHAR(255) NOT NULL,
    date VARCHAR(50) NOT NULL,
    clock_in VARCHAR(50) DEFAULT '--:--',
    clock_out VARCHAR(50) DEFAULT '--:--',
    total_hours VARCHAR(50) DEFAULT '0h 0m',
    status VARCHAR(20) DEFAULT 'HADIR',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_reports_date ON tb_lopiq_attendance_reports(date);
CREATE INDEX IF NOT EXISTS idx_reports_department ON tb_lopiq_attendance_reports(department);
