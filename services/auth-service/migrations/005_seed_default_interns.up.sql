-- Seed default intern accounts for direct access and testing
INSERT INTO auth_users (id, nip, email, name, role, jabatan, unit_kerja, password, is_active)
VALUES 
(101, '0051234567', 'admin@example.com', 'Sarah Jenkins', 'intern', 'SMK Negeri 1 Bulukumba', 'Rekayasa Perangkat Lunak', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true),
(102, '2024001', 'hikma@gmail.com', 'Hikma', 'intern', 'Universitas Negeri Makassar', 'Product Design', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true),
(103, '2024002', 'budi@gmail.com', 'Budi Santoso', 'intern', 'SMK Negeri 1 Bulukumba', 'Frontend Dev', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true)
ON CONFLICT (email) DO UPDATE SET nip = EXCLUDED.nip, name = EXCLUDED.name, role = EXCLUDED.role, jabatan = EXCLUDED.jabatan, unit_kerja = EXCLUDED.unit_kerja;

SELECT setval(pg_get_serial_sequence('auth_users', 'id'), (SELECT COALESCE(MAX(id), 1) FROM auth_users));
