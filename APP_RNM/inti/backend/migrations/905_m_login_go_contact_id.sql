-- M_LOGIN_GO: CONTACT_ID, username unik, dan email unik (keputusan work owner 03-10-2026, pilihan V1 rancangan
-- Kelola Marketing Officer: "M_LOGIN_GO ID nya pake CON-xxx"; "tambahkan proteksi format email, atau email sudah
-- terdaftar"; "tambahkan juga proteksi username sudah ada"). Login lewat email DIBATALKAN work owner ("LOGIN LEWAT
-- EMAIL TIDAK JADI!!"): login tetap dengan username saja.
--
-- LOGIN_ID TETAP username dan tetap PK - sesi, pelaku, menu, dan workbasket tidak berubah. Kolom baru CONTACT_ID
-- adalah kode orang yang TIDAK PERNAH berubah: `CON-n`, n dari M_LOGIN_GO_CONTACT_SEQ. Nomornya mulai 1001, di atas
-- nomor kontak SFAGIS terbesar (`ASM-SFAGIS-WORK-CONTACT CON-108` di DEV), supaya satu nomor CON tidak pernah
-- menunjuk dua orang. Akun yang sudah ada diberi nomor di sini, lalu kolomnya NOT NULL dan unik.
--
-- LOGIN_ID unik TANPA BEDA HURUF (`LOWER(LOGIN_ID)`, permintaan work owner "tambahkan juga proteksi username sudah
-- ada"): `uji` dan `UJI` tidak boleh dua akun. Aplikasi memeriksanya lebih dulu; index ini pengaman terakhir untuk
-- dua simpan serentak.
--
-- EMAIL disimpan huruf kecil dan UNIK tanpa beda huruf (`LOWER(EMAIL)`, "proteksi email sudah terdaftar"): satu
-- email hanya milik satu akun. Email kosong (NULL) tidak ikut index, jadi akun tanpa email tetap boleh banyak.
-- ⚠️ Index username atau email GAGAL (ORA-01452) bila dua akun sudah berbagi username tanpa beda huruf atau berbagi
-- email - bersihkan dulu, lalu ulangi. DEV 03-10-2026: 3 akun, nol ganda, nol email.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
CREATE SEQUENCE {skema}.M_LOGIN_GO_CONTACT_SEQ START WITH 1001 INCREMENT BY 1 NOCACHE
/
ALTER TABLE {skema}.M_LOGIN_GO ADD (
  CONTACT_ID VARCHAR2(20)
)
/
UPDATE {skema}.M_LOGIN_GO SET CONTACT_ID = 'CON-' || {skema}.M_LOGIN_GO_CONTACT_SEQ.NEXTVAL
 WHERE CONTACT_ID IS NULL
/
ALTER TABLE {skema}.M_LOGIN_GO MODIFY (CONTACT_ID NOT NULL)
/
CREATE UNIQUE INDEX {skema}.UX_M_LOGIN_GO_CONTACT_ID ON {skema}.M_LOGIN_GO (CONTACT_ID)
/
CREATE UNIQUE INDEX {skema}.UX_M_LOGIN_GO_LOGIN_ID ON {skema}.M_LOGIN_GO (LOWER(LOGIN_ID))
/
UPDATE {skema}.M_LOGIN_GO SET EMAIL = LOWER(EMAIL) WHERE EMAIL <> LOWER(EMAIL)
/
CREATE UNIQUE INDEX {skema}.UX_M_LOGIN_GO_EMAIL ON {skema}.M_LOGIN_GO (LOWER(EMAIL))
/
