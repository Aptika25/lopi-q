-- SQL Script untuk menambahkan akun Peserta Magang baru secara langsung ke PostgreSQL DB
-- Password default: password123 (Bcrypt hash: $2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO)

-- 1. Insert ke Database db_lopiq_auth (Tabel auth_users)
INSERT INTO auth_users (nip, email, name, role, jabatan, unit_kerja, password, is_active)
VALUES (
    '0091755897', 
    'adeanisa150299@gmail.com', 
    'ade anisa', 
    'intern', 
    'SMK TI Bulukumba', 
    'Rekayasa Perangkat Lunak', 
    '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', 
    true
)
ON CONFLICT (email) DO UPDATE SET 
    nip = EXCLUDED.nip, 
    name = EXCLUDED.name, 
    role = EXCLUDED.role, 
    jabatan = EXCLUDED.jabatan, 
    unit_kerja = EXCLUDED.unit_kerja, 
    password = EXCLUDED.password, 
    is_active = true;

-- 2. Insert ke Database db_lopiq_user (Tabel users)
INSERT INTO users (nip, email, name, role, jabatan, unit_kerja, password_hash, is_active)
VALUES (
    '0091755897', 
    'adeanisa150299@gmail.com', 
    'ade anisa', 
    'intern', 
    'SMK TI Bulukumba', 
    'Rekayasa Perangkat Lunak', 
    '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', 
    true
)
ON CONFLICT (email) DO UPDATE SET 
    nip = EXCLUDED.nip, 
    name = EXCLUDED.name, 
    role = EXCLUDED.role, 
    jabatan = EXCLUDED.jabatan, 
    unit_kerja = EXCLUDED.unit_kerja, 
    password_hash = EXCLUDED.password_hash, 
    is_active = true;
