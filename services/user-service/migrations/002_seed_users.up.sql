-- Seed: Profil Super Admin dan Call Taker di tabel users (user-service)

INSERT INTO users (nip, email, name, role, jabatan, unit_kerja, password_hash, is_active)
VALUES
  -- Super Admin
  ('199501012020011000', 'aswan@bulukumbakab.go.id', 'Muhammad Aswan', 'superadmin', 'HEAD OF DISKOMINFO', 'Diskominfo Kab. Bulukumba', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true),

  -- Call Takers (10 petugas)
  ('19940503202521 1 138', 'amappalua@bulukumbakab.go.id', 'A.Mappalua, S.Pd', 'call_taker', 'PENATA LAYANAN OPERASIONAL', 'Dinas Sosial', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true),
  ('19870304202521 1 061', 'suherman@bulukumbakab.go.id', 'Suherman, S.Pd', 'call_taker', 'PENATA LAYANAN OPERASIONAL', 'Badan Penanggulangan Bencana Daerah', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true),
  ('20000206202521 1 166', 'riswandirisman@bulukumbakab.go.id', 'Riswandi Risman', 'call_taker', 'OPERATOR LAYANAN OPERASIONAL', 'Dinas Kesehatan', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true),
  ('19900215202521 1 114', 'abilkizri@bulukumbakab.go.id', 'Abil Kizri', 'call_taker', 'OPERATOR LAYANAN OPERASIONAL', 'Dinas Perhubungan', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true),
  ('19911005202521 1 087', 'imamardiyansah@bulukumbakab.go.id', 'Imam Ardiyansah', 'call_taker', 'OPERATOR LAYANAN OPERASIONAL', 'Satpol, Pemadam Kebakaran dan Penyelamatan', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true),
  ('19861130202521 1 101', 'abdrahim@bulukumbakab.go.id', 'Abd.Rahim', 'call_taker', 'OPERATOR LAYANAN OPERASIONAL', 'Dinas Sosial', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true),
  ('19860304202521 1 147', 'munawir@bulukumbakab.go.id', 'Munawir Syadzali', 'call_taker', 'PENATA LAYANAN OPERASIONAL', 'Badan Penanggulangan Bencana Daerah', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true),
  ('19760802200604 1 017', 'abdullah@bulukumbakab.go.id', 'Abdullah, S.Kep., Ns', 'call_taker', 'PERENCANA', 'Dinas Kesehatan', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true),
  ('19860712202521 1 089', 'ismail@bulukumbakab.go.id', 'Ismail, S.Sos', 'call_taker', 'PENATA LAYANAN OPERASIONAL', 'Dinas Perhubungan', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true),
  ('19960328202521 1 050', 'aldiafdal@bulukumbakab.go.id', 'Aldi Afdali Saputra', 'call_taker', 'OPERATOR LAYANAN OPERASIONAL', 'Satpol, Pemadam Kebakaran dan Penyelamatan', '$2a$10$EwQk2ADnVXXIVSSSueM4sOnO9Py1TQB0l5Bynadgn1Ke7TXT6W/vO', true)
ON CONFLICT (email) DO NOTHING;

SELECT setval(pg_get_serial_sequence('users', 'id'), (SELECT COALESCE(MAX(id), 1) FROM users));
