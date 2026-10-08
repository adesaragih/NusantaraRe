-- ⛔⛔ BERKAS INI MENGUBAH TABEL MILIK MODUL `treatyin`, TETAPI TINGGAL DI
-- FOLDER MIGRASI `treatyinadjustment` — tempat parkir yang sama dengan
-- `444`–`451`: rentang `treatyin` (`400-439`) PENUH.
--
-- ---------------------------------------------------------------------
-- `ValueDifference.TotalLimitsROL` — selisih yang nol punya kolom
-- ---------------------------------------------------------------------
--
-- Audit jalur simpan 8 Oktober 2026 (permintaan pemilik proses *"perbaiki
-- alur simpan"*) mengadu 113 kunci layar ekspor dengan peta pendaratan.
-- Empat belas nol punya rumah; sesudah disaring, satu di antaranya cacat
-- sungguhan dan BERKOLOM — yaitu berkas ini.
--
-- ⭐ `services/hitung_selisih.go` baris 232 menulisnya pada SETIAP
-- perhitungan selisih:
--
--   vd.TotalLimitsROL = ActualValue.TotalLimitsROL − OLDDATA.TotalLimitsROL
--
-- Dua saudaranya SUDAH punya kolom sejak awal — `VALUEDIFF_RNMSHARE` dan
-- `VALUEDIFF_BROKERAGEPCT` — dan yang ketiga ini terlewat. Akibatnya ia
-- dilaporkan "TIDAK tersimpan" dan hilang setiap Save.
--
-- ⚠️ Nama kolomnya mengikuti pola kedua saudaranya (`VALUEDIFF_` + ruas),
-- BUKAN huruf besar kunci penuh: kunci bertitik di peta pendaratan memang
-- dipetakan ke nama kolom yang diringkas, dan `TestPetaPendaratanCocokDenganDDL`
-- mengadu peta terhadap DDL, bukan terhadap aturan penamaan.
--
-- ⚠️ `VARCHAR2(4000 CHAR)` seperti SELURUH kolom pendaratan.

ALTER TABLE {skema}.T_TREATY_REVISION ADD (
  VALUEDIFF_TOTALLIMITSROL         VARCHAR2(4000 CHAR)
)
/
