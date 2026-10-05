-- M_LOGIN_GO_MENU.HAK: tingkat akses satu menu untuk satu akun (keputusan work owner 04-10-2026: "menu bisa dilihat
-- orang, tapi orang itu belum tentu bisa menambah isinya"; tandanya di M_LOGIN_GO_MENU, per menu).
--
--   PENUH - membaca dan menulis (perilaku sebelum migrasi ini);
--   LIHAT - hanya membaca: gerbang `cmd/api` menolak POST/PUT/DELETE modul itu, layar menyembunyikan tombolnya.
--
-- Baris yang sudah ada menjadi PENUH lewat DEFAULT, jadi tidak ada akun yang kehilangan akses. LIHAT hanya boleh untuk
-- menu modul yang MENDAFTAR (`inti.Pendaftaran.HakLihat`) - diperiksa Kelola User, bukan oleh CHECK di sini, sebab
-- daftarnya hidup di kode modul. LIHAT berlaku juga untuk superadmin (pemegang `kelolauser`, keputusan work owner
-- 05-10-2026); Kelola User sendiri tidak pernah LIHAT, jadi superadmin selalu dapat mengembalikan aksesnya.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
ALTER TABLE {skema}.M_LOGIN_GO_MENU ADD (
  HAK VARCHAR2(10) DEFAULT 'PENUH' NOT NULL
)
/
ALTER TABLE {skema}.M_LOGIN_GO_MENU ADD CONSTRAINT CK_M_LOGIN_GO_MENU_HAK CHECK (HAK IN ('LIHAT', 'PENUH'))
/
