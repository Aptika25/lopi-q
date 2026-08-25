-- Seed single Super Admin profile in users table (user-service)
INSERT INTO users (id, nip, email, name, role, jabatan, unit_kerja, password_hash, is_active)
VALUES (1, '199708192025061003', 'aswan@bulukumbakab.go.id', 'Muhammad Aswan, S.T.', 'superadmin', 'HEAD OF DISKOMINFO', 'Diskominfo Kab. Bulukumba', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true)
ON CONFLICT (email) DO UPDATE SET nip = '199708192025061003', name = 'Muhammad Aswan, S.T.';

SELECT setval(pg_get_serial_sequence('users', 'id'), (SELECT COALESCE(MAX(id), 1) FROM users));
