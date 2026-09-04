-- Seed default intern accounts for direct access and testing
INSERT INTO auth_users (id, nip, email, name, role, jabatan, unit_kerja, password, is_active)
VALUES 
(2, '0091755987', 'adeanisa150299@gmail.com', 'adhe anisa', 'intern', 'SMKS TI Bulukumba', 'TKJ', '$2a$10$A.t5X3JcRyEohmF/VNEbyuCf3URanCEk0dO3ViXuLUtBSZMoAryme', true),
(4, '3084633444', 'apriliahikma45@gmail.com', 'aprilia hikma', 'intern', 'SMKS Muhammadiyah Bulukumba', 'TKJ', '$2a$10$pFM0wbg/InF/VZqdPPczxO8VXhYZpsxLhGRCzqcwaAbtmS9WXad3i', true),
(6, '0009804858', 'nurhidayah18032009@gmail.com', 'Nurhidayah', 'intern', 'SMKS muhammadiyah Bulukumba', 'TKJ', '$2a$10$acYyTR4pfSq/A7gCO2Aj2.pOaOmDfXXhWb6QopJ4/XxYtEWo8BR1K', true)
ON CONFLICT (email) DO UPDATE SET nip = EXCLUDED.nip, name = EXCLUDED.name, role = EXCLUDED.role, jabatan = EXCLUDED.jabatan, unit_kerja = EXCLUDED.unit_kerja;

SELECT setval(pg_get_serial_sequence('auth_users', 'id'), (SELECT COALESCE(MAX(id), 1) FROM auth_users));
