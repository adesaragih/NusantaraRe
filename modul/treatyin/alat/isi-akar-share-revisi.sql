-- Mengisi SEKALI kolom akar Share di `T_TREATY_REVISION` untuk kontrak lama
-- (Pega) — sesudah migrasi `445_kolom_share_revisi` terpasang.
--
-- ⛔ BUKAN berkas migrasi, dan sengaja TIDAK di `backend/migrations/`: ia
-- MENULIS DATA ke POOLDATA, jadi dijalankan dengan tangan oleh pemilik
-- proses, sesudah diperiksa. Ganti `{skema}` dengan skema sasaran
-- (`POOLDATA`) sebelum menjalankannya.
--
-- ⛔ NOL `COMMIT` di berkas ini. Periksa hasilnya dengan kueri di ekor
-- berkas, lalu `COMMIT` — atau `ROLLBACK` bila tidak sesuai.
--
-- ⭐ SUMBER — keputusan pemakai 7 Oktober 2026. Hanya nilai yang PEGA TULIS
-- SENDIRI; nol JSON, nol `M_TREATY_IN`, nol tebakan dari rasio Gross/MDP:
--
--   RNMSHARE          1. salinan baris Share  `T_TREATY_SHARE.RNMSHARE`
--                        (`TreatyInNonAddItem` [14.1] `.RNMShare := TreatyIn.RNMShare`),
--                        baris berurutan terkecil yang berisi
--                     2. `TREATYINDETAIL.RNM_SHARE`
--                        (`SaveTreatyInDetail_Act` [5.1] `RNM_SHARE := TreatyIn.RNMShare`)
--   BROKERAGEPERCENT  1. `DEDUCTIONPCT` deduksi `Brokerage fee` baris Share
--                        PERTAMA (`TreatyInSetBrokerage` [2])
--                     2. `TREATYINDETAIL.BROKERAGE`
--
--   RNMSHAREACROSSTHEBOARD dan RNMSHAREDEDUCTED TIDAK diisi: yang pertama
--   tidak punya sumber di mana pun (layar memakai bawaan Pega `true`), yang
--   kedua turunan (`RNMShare − FacultativeShare`) yang layar hitung sendiri.
--
-- Terukur 7 Oktober 2026 (baca saja): dari 1.079 kontrak Non-Prop yang
-- punya baris Share, 314 tertutup sumber 1, +480 tertutup `TREATYINDETAIL`
-- (283 kontrak yang punya keduanya: nilainya sama 283/283), 285 tetap
-- kosong — diisi pemakai lalu Save.
--
-- ⚠️ Hanya kolom yang MASIH KOSONG yang diisi: nilai yang sudah disimpan
-- aplikasi (Save) tidak pernah ditimpa, dan berkas ini aman diulang.

UPDATE {skema}.T_TREATY_REVISION r
   SET r.RNMSHARE = COALESCE(
         (SELECT MAX(s.RNMSHARE) KEEP (DENSE_RANK FIRST ORDER BY s.URUTAN)
            FROM {skema}.T_TREATY_SHARE s
           WHERE s.MASTERID = r.MASTERID
             AND TRIM(s.RNMSHARE) IS NOT NULL),
         (SELECT RTRIM(TO_CHAR(MAX(d.RNM_SHARE), 'FM99999999990.99999999'), '.')
            FROM {skema}.TREATYINDETAIL d
           WHERE d.TREATYID = r.MASTERID
             AND d.PROPORTIONTYPE = 'NonProportional'))
 WHERE r.PROPORTIONTYPE = 'NonProportional'
   AND TRIM(r.RNMSHARE) IS NULL
/

UPDATE {skema}.T_TREATY_REVISION r
   SET r.BROKERAGEPERCENT = COALESCE(
         (SELECT MAX(dd.DEDUCTIONPCT) KEEP (DENSE_RANK FIRST ORDER BY dd.URUTAN)
            FROM {skema}.T_TREATY_SHARE s
            JOIN {skema}.T_TREATY_SHARE_DEDUCTION dd ON dd.IDINDUK = s.ID
           WHERE s.MASTERID = r.MASTERID
             AND s.URUTAN = (SELECT MIN(s2.URUTAN) FROM {skema}.T_TREATY_SHARE s2 WHERE s2.MASTERID = r.MASTERID)
             AND dd.COMMENT_ = 'Brokerage fee'
             AND TRIM(dd.DEDUCTIONPCT) IS NOT NULL),
         (SELECT RTRIM(TO_CHAR(MAX(d.BROKERAGE), 'FM99999999990.99999999'), '.')
            FROM {skema}.TREATYINDETAIL d
           WHERE d.TREATYID = r.MASTERID
             AND d.PROPORTIONTYPE = 'NonProportional'))
 WHERE r.PROPORTIONTYPE = 'NonProportional'
   AND TRIM(r.BROKERAGEPERCENT) IS NULL
/

-- Periksa sebelum COMMIT:
--
--   SELECT COUNT(*) kontrak,
--          SUM(CASE WHEN TRIM(RNMSHARE) IS NOT NULL THEN 1 ELSE 0 END) rnm_terisi,
--          SUM(CASE WHEN TRIM(BROKERAGEPERCENT) IS NOT NULL THEN 1 ELSE 0 END) brokerage_terisi
--     FROM {skema}.T_TREATY_REVISION
--    WHERE PROPORTIONTYPE = 'NonProportional';
