-- 883 - CITYINPUT.PROVINCENAME, isi awal dari RW, dan kota RW yang belum ada di CITYINPUT.
--
-- Keputusan work owner 04-10-2026 (diteruskan sesi nusantarare-9d): "di cityinput tambahin PROVINCENAME, data nya diambil
-- dari tabel rw . rw.CITYNAME=cityinput.note baru di isi cityinput .PROVINCENAME". Lebar = RW.PROVINCENAME VARCHAR2(4000)
-- `[terverifikasi]` DDL `D:\migrasi\RNM\DDL\RW.txt`. Kolom biasa (diubah menu City), bukan turunan PROVINCE.NOTE.
--
-- Isi awal (MD-11, keputusan agent sesi c3 - menunggu konfirmasi):
--   - Pencocokan PERSIS kata work owner: `r.CITYNAME = c.NOTE` - tanpa UPPER / TRIM.
--   - Satu CITYNAME dapat muncul di banyak baris RW (satu baris per RW / kode pos). Diisi HANYA bila baris RW yang cocok
--     memuat TEPAT SATU nilai PROVINCENAME berbeda (COUNT(DISTINCT) = 1); lebih dari satu = dibiarkan NULL, diisi
--     manusia lewat menu City. PROVINCENAME RW yang NULL tidak dihitung COUNT(DISTINCT): satu nilai + NULL = diisi.
--   - RW.STS_AKTIF TIDAK disaring: work owner tidak menyebutnya, dan nama provinsi baris RW nonaktif tetap sebuah
--     pemetaan kota -> provinsi. Akibatnya: nama lama di baris nonaktif yang berbeda membuat kota itu NULL (ambigu),
--     bukan salah isi. Alternatifnya (`AND r.STS_AKTIF = '1'`, filter RD BrowseRW) menunggu keputusan work owner.
--   - `c.PROVINCENAME IS NULL`: isi yang sudah ada tidak ditimpa bila pernyataan ini dijalankan lagi.
--
-- Kota baru dari RW (keputusan work owner 04-10-2026: "cek di tabel rw, ada kolom cityname ..., bandingkan ke tabel
-- city, jika belum ada di tabel city, km tambahin, untuk id nya +1 dari id terakhir. cityname=note"; "tidak usah
-- zipcode, provincename aja"). MD-12, keputusan agent sesi c3 - menunggu konfirmasi:
--   - Satu baris per CITYNAME RW berbeda yang TRIM-nya tidak kosong (NULL / spasi saja dilewati) dan tidak sama PERSIS
--     dengan NOTE mana pun di CITYINPUT - kunci banding sama dengan UPDATE (tanpa UPPER / TRIM; "Jakarta" dan
--     "JAKARTA" = dua kota). RW.STS_AKTIF tidak disaring (sama dengan UPDATE, MD-11).
--   - NOTE = CITYNAME; PROVINCENAME = aturan UPDATE (tepat satu nilai berbeda, selain itu NULL); STS_AKTIF = DEFAULT
--     '1'; PROVINCEID / BRANCHID / EMAIL / MOID / JABODETABEKSTATUS NULL (tidak diminta).
--   - ID = "+1 dari id terakhir": MAX ID CITYINPUT yang SELURUHNYA digit (`REGEXP_LIKE '^[0-9]+$'`; ID lain tidak ikut
--     dihitung) + ROW_NUMBER() urut CITYNAME, disimpan teks tanpa nol di depan (TO_CHAR angka). Tanpa ID digit: mulai 1.
--     Satu pernyataan: nomornya berurutan, tanpa tabrakan dengan sesama baris baru.
--   - Penanda CREATE_OP = 'MIGRASI-883', TGL_CREATE = SYSDATE (kolom jejak 882): jalur mundur membuang TEPAT baris ini.
-- Pernyataan ini SESUDAH UPDATE: baris baru sudah membawa PROVINCENAME-nya sendiri.
--
-- Bila 883 gagal SESUDAH ALTER, pengulangan berhenti di ORA-01430 (kolom sudah ada) dan T_MIGRASI belum mencatat 883 -
-- pola 059: DBA menjalankan UPDATE di bawah manual, lalu mencatat `883_cityinput_provincename` di T_MIGRASI.
-- Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.CITYINPUT ADD (
  PROVINCENAME VARCHAR2(4000)
)
/
UPDATE {skema}.CITYINPUT c
   SET c.PROVINCENAME = (SELECT MAX(r.PROVINCENAME) FROM {skema}.RW r WHERE r.CITYNAME = c.NOTE)
 WHERE c.PROVINCENAME IS NULL
   AND (SELECT COUNT(DISTINCT r.PROVINCENAME) FROM {skema}.RW r WHERE r.CITYNAME = c.NOTE) = 1
/
INSERT INTO {skema}.CITYINPUT (ID, NOTE, PROVINCENAME, CREATE_OP, TGL_CREATE)
SELECT TO_CHAR(m.MAKS + ROW_NUMBER() OVER (ORDER BY b.CITYNAME)), b.CITYNAME,
       CASE WHEN b.NPROV = 1 THEN b.PROV END, 'MIGRASI-883', SYSDATE
  FROM (SELECT r.CITYNAME, COUNT(DISTINCT r.PROVINCENAME) AS NPROV, MAX(r.PROVINCENAME) AS PROV
          FROM {skema}.RW r
         WHERE TRIM(r.CITYNAME) IS NOT NULL
         GROUP BY r.CITYNAME) b
 CROSS JOIN (SELECT NVL(MAX(CASE WHEN REGEXP_LIKE(c.ID, '^[0-9]+$') THEN TO_NUMBER(c.ID) END), 0) AS MAKS
               FROM {skema}.CITYINPUT c) m
 WHERE NOT EXISTS (SELECT 1 FROM {skema}.CITYINPUT x WHERE x.NOTE = b.CITYNAME)
/
