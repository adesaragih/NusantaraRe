-- =====================================================================
-- V01_V_AKSEPTASI_KOMPATIBEL.sql
-- USULAN. BELUM PERNAH DIJALANKAN. UNTUK DIBACA, BUKAN UNTUK DIJALANKAN.
--
-- Tiket 29. Kontrak keluar ke Arasapas dan kasir. Bentuk lama dipertahankan,
-- termasuk CASEID (pengecualian ADR-0021). Daftar kolom DIBACA dari XML
-- (sapuan S1), bukan dirancang.
-- =====================================================================


-- ---------------------------------------------------------------------
-- View ini menghasilkan SELISIH, bukan nilai mutlak
-- ---------------------------------------------------------------------
-- Ini yang membedakannya dari sekadar agregat, dan versi sebelumnya
-- MELEWATKANNYA sama sekali — ia mengeluarkan nilai mutlak.
--
-- Sistem hilir tidak menerima keadaan; ia menerima TAMBAHAN. Buktinya
-- RDBList\GetDataOS.xml:
--
--   select sum(nvl(a.data_json.Value,0)) ...
--     from OS_AKSEPTASI_KLAIM a
--    WHERE CASEID = {pyWorkPage.pzInsKey}
--      AND a.data_json.TypeLoss = {InputParamOs.TypeLoss}
--      AND a.data_json.Currency = {InputParamOs.Currency}
--      AND STS_REJECT = 0
--
-- lalu SaveDataToOSAksep_Act MENGURANGKAN hasilnya. Prosedur di seberang
-- SELALU INSERT, tidak pernah UPDATE — dan kedua hal itu pasangan wajib.
--
-- Maka yang dikirim = nilai sekarang DIKURANGI apa yang sudah pernah
-- dikirim. Yang sudah pernah dikirim ada di ARSIP_MUATAN_KELUAR, dan
-- HANYA baris yang TIDAK ditolak yang terhitung — STS_REJECT = 0 di sisi
-- sana, DITOLAK = 0 di sisi sini.
--
-- Mengeluarkan nilai mutlak ke jalur yang menjumlahkan tambahan berarti
-- MELIPATGANDAKAN setiap nilai pada pengiriman kedua.


-- ---------------------------------------------------------------------
-- Penjaga: kelompok yang tidak sepakat DITOLAK, tidak diam-diam dirata-rata
-- ---------------------------------------------------------------------
-- Agregasi halus (per bagian layer) ke kasar (per jenis reasuransi) melebur
-- beberapa baris jadi satu. Besaran ADITIF boleh dijumlahkan; besaran
-- PER-LAYER tidak — Limit Layer, Premi Deposit, persen premi pemulihan,
-- kurs, dan porsi reasuradur. Bila baris yang dilebur TIDAK SEPAKAT pada
-- salah satunya, tidak ada satu nilai pun yang benar untuk dikirim.
--
-- View ini MENOLAK menghasilkan baris untuk kelompok itu, dan kelompoknya
-- muncul di V_AKSEPTASI_DITOLAK beserta sebabnya. Tidak ada yang hilang
-- diam-diam — itu syaratnya, dan keduanya HARUS memakai pengelompokan dan
-- penyaring yang SAMA PERSIS, kalau tidak yang ditolak di sini dapat lenyap
-- dari sana juga. Versi sebelumnya TIDAK sama: V02 menyambung tanpa
-- LAYER_BAGIAN dan tanpa menyaring KEADAAN_BARIS.
--
-- REQ-033 mengukur apakah ketidaksepakatan itu pernah terjadi di data lama.
-- Ia MENUNGGUI, tidak menahan: jawabannya mengubah BERAPA BANYAK kelompok
-- yang ditolak, bukan apakah penolakannya ada.


CREATE OR REPLACE VIEW KLAIMNP.V_AKSEPTASI_KOMPATIBEL AS
WITH SEKARANG AS (
  SELECT a.ID_KLAIM,
         a.JENIS_REASURANSI,
         a.JENIS_REASURANSI_ID,
         a.MATA_UANG,
         MIN(a.KURS)                                AS KURS,
         MIN(al.PORSI_REASURADUR)                   AS PORSI_REASURADUR,
         SUM(al.NILAI_TERALOKASI_DIPAKAI)           AS NILAI_TERSEBAR,
         SUM(a.NILAI_DISETUJUI)                     AS NILAI_BRUTO,
         SUM(al.BIAYA_PENILAIAN_DIPAKAI)            AS BIAYA_PENILAIAN,
         SUM(al.SALVAGE_DIPAKAI)                    AS SALVAGE,
         SUM(al.BIAYA_LAIN_DIPAKAI)                 AS BIAYA_LAIN
    FROM KLAIMNP.AKSEPTASI a
    LEFT JOIN KLAIMNP.ALOKASI_LAYER al
           ON al.ID_KLAIM = a.ID_KLAIM
          AND al.LAYER = a.LAYER
          AND al.LAYER_JENIS = a.LAYER_JENIS
          AND al.LAYER_BAGIAN = a.LAYER_BAGIAN
          AND al.LAYER_BAGIAN_JENIS = a.LAYER_BAGIAN_JENIS
          AND al.MATA_UANG = a.MATA_UANG
   WHERE a.KEADAAN_BARIS = 'LENGKAP'
     AND (al.ID_ALOKASI IS NULL OR al.KEADAAN_BARIS = 'LENGKAP')
   GROUP BY a.ID_KLAIM, a.JENIS_REASURANSI, a.JENIS_REASURANSI_ID, a.MATA_UANG
  HAVING MIN(al.LIMIT_LAYER)            = MAX(al.LIMIT_LAYER)
     AND MIN(al.PREMI_DEPOSIT)          = MAX(al.PREMI_DEPOSIT)
     AND MIN(al.PERSEN_PREMI_PEMULIHAN) = MAX(al.PERSEN_PREMI_PEMULIHAN)
     AND MIN(a.KURS)                    = MAX(a.KURS)
     AND MIN(al.PORSI_REASURADUR)       = MAX(al.PORSI_REASURADUR)
),
SUDAH_DIKIRIM AS (
  SELECT ID_KLAIM, JENIS_REASURANSI, MATA_UANG,
         SUM(NILAI_TERSEBAR)   AS NILAI_TERSEBAR,
         SUM(NILAI_BRUTO)      AS NILAI_BRUTO,
         SUM(BIAYA_PENILAIAN)  AS BIAYA_PENILAIAN,
         SUM(SALVAGE)          AS SALVAGE,
         SUM(BIAYA_LAIN)       AS BIAYA_LAIN
    FROM KLAIMNP.ARSIP_MUATAN_KELUAR
   WHERE DITOLAK = 0
   GROUP BY ID_KLAIM, JENIS_REASURANSI, MATA_UANG
)
SELECT mk.CASEID_LAMA                AS CASEID,
       kl.NOMOR_KLAIM                AS NOCLAIM,
       s.JENIS_REASURANSI            AS TYPELOSS,
       s.JENIS_REASURANSI_ID         AS TYPELOSSID,
       s.MATA_UANG                   AS CURRENCY,
       s.KURS                        AS KURSVALUE,
       ROUND(s.NILAI_TERSEBAR   - NVL(d.NILAI_TERSEBAR,   0), 2) AS VALUE,
       ROUND(s.NILAI_BRUTO      - NVL(d.NILAI_BRUTO,      0), 2) AS GROSSVALUE,
       ROUND(s.BIAYA_PENILAIAN  - NVL(d.BIAYA_PENILAIAN,  0), 2) AS ADJUSTERFEE,
       ROUND(s.SALVAGE          - NVL(d.SALVAGE,          0), 2) AS SALVAGE,
       ROUND(s.BIAYA_LAIN       - NVL(d.BIAYA_LAIN,       0), 2) AS CNPOTHERSFEE,
       s.PORSI_REASURADUR            AS PERSENRNM,
       kl.ID_TREATY                  AS IDMASTERTREATY,
       kl.PENYEBAB_KERUGIAN_ID       AS CAUSEOFLOSSID,
       kl.PENYEBAB_KERUGIAN_ID       AS CAUSEOFLOSS
  FROM SEKARANG s
  JOIN KLAIMNP.KLAIM kl ON kl.ID_KLAIM = s.ID_KLAIM
  LEFT JOIN SUDAH_DIKIRIM d
         ON d.ID_KLAIM = s.ID_KLAIM
        AND d.JENIS_REASURANSI = s.JENIS_REASURANSI
        AND d.MATA_UANG = s.MATA_UANG
  LEFT JOIN KLAIMNP.MIGRASI_KORELASI mk
         ON mk.TABEL_TUJUAN = 'KLAIM' AND mk.ID_TUJUAN = kl.ID_KLAIM;


-- ---------------------------------------------------------------------
-- Apa yang view ini TIDAK lakukan
-- ---------------------------------------------------------------------
--   1. Ia TIDAK menyaring selisih nol. Kelompok yang sudah terkirim penuh
--      muncul dengan nilai 0 — dan itu BENAR: nol berarti "tidak ada yang
--      perlu ditambahkan", dan itu berbeda dari kelompoknya tidak muncul.
--      Penyaringannya urusan pengirim, bukan urusan kontrak.
--   2. Ia TIDAK mengirim. Pengiriman lapisan aplikasi; view ini menyatakan
--      APA yang akan dikirim bila pengiriman terjadi.
--   3. Ia TIDAK menyimpan. ADR-0023 utuh: ia membaca arsip, bukan
--      menyalinnya. ADR-0017 utuh: yang menulis arsip adalah proses
--      pengiriman lewat akun aplikasi.
--   4. LayerPart dan LayerPartType TIDAK muncul di keluaran — keduanya
--      satuan halus, dan kontrak lama tidak mengenalnya.
-- ---------------------------------------------------------------------

-- ADR-0017/0028: hilir membaca HANYA lewat view, tanpa hak atas tabel kanonik.
GRANT SELECT ON KLAIMNP.V_AKSEPTASI_KOMPATIBEL TO POOLDATA;
GRANT SELECT ON KLAIMNP.V_AKSEPTASI_KOMPATIBEL TO KLAIMNP_HILIR;
