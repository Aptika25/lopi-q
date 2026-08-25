-- Seed single superadmin auth account (aswan@bulukumbakab.go.id, NIP: 199708192025061003, password: Asw&a198)
INSERT INTO auth_users (id, nip, email, name, role, jabatan, unit_kerja, password, is_active)
VALUES (1, '199708192025061003', 'aswan@bulukumbakab.go.id', 'Muhammad Aswan, S.T.', 'superadmin', 'HEAD OF DISKOMINFO', 'Diskominfo Kab. Bulukumba', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true)
ON CONFLICT (email) DO UPDATE SET nip = '199708192025061003', name = 'Muhammad Aswan, S.T.';

SELECT setval(pg_get_serial_sequence('auth_users', 'id'), (SELECT COALESCE(MAX(id), 1) FROM auth_users));
