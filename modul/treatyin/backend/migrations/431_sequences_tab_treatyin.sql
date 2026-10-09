-- Sequence identitas untuk delapan tabel tab Treaty In
--
-- Migrasi tiket tab Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Berdiri TERPISAH dari tabelnya, mengikuti 402, 415, 419 dan 421.
--
-- INV-02: pengenal datang dari SEQUENCE, tidak pernah dari cap waktu, teks,
-- maupun dari `MASTERID` ditambah urutan. INV-03: `NOCYCLE` - satu nomor
-- tidak pernah terpakai dua kali.
--
-- ⚠️ Awalan `SEQ_MTI_`, bukan `SEQ_TRIN_` yang dipakai 415/419/421. Sebabnya
-- batas 30 bita, dan ia menggigit di sini: `SEQ_TRIN_M_TREATYIN_REPORTINGPERIOD`
-- 35 bita. Awalan `MTI` dipakai KONSISTEN untuk seluruh objek kedelapan tabel
-- ini - PK, UQ, FK, index, sequence - supaya tidak ada yang perlu menebak
-- kapan singkatan berlaku dan kapan tidak. Yang terpanjang di bawah,
-- SEQ_MTI_INSTALLMENTITEM, 23 bita; `TestNamaObjekDiBawahTigaPuluhBita`
-- membacanya dari teks DDL, bukan dari daftar ini.
--
-- ⛔ `START WITH 1` aman walau tabelnya akan dimuat ulang berkali-kali:
-- pengenalnya tidak pernah dipakai ulang, dan pemuat yang memuat ulang
-- MENGOSONGKAN tabelnya lalu mengambil nomor baru. Yang tidak boleh terjadi
-- adalah dua baris berbeda bernomor sama, dan NOCYCLE-lah yang mencegahnya.
CREATE SEQUENCE {skema}.SEQ_MTI_REPORTINGPERIOD START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_MTI_PORTFOLIO START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_MTI_ACCUMULATION START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_MTI_EGNPI START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_MTI_RETENTION START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_MTI_INSTALLMENT START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_MTI_INSTALLMENTITEM START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
CREATE SEQUENCE {skema}.SEQ_MTI_COMMENT START WITH 1 INCREMENT BY 1 NOCACHE NOCYCLE
/
