-- Pembalikan `444_kolom_revisi_yang_layar_baca.sql`.
--
-- ⚠ SATU pernyataan, bukan enam belas. `ALTER TABLE … DROP (a, b, c)`
-- membuang seluruhnya dalam satu kali jalan; enam belas pernyataan terpisah
-- membuat kegagalan di yang kesembilan meninggalkan tabel setengah jalan,
-- dan DDL Oracle mengikat sendiri sehingga setengah jalan itu PERMANEN.
--
-- ⛔ DATA HILANG, dan itu memang artinya: keenam belas kolom ini diisi
-- pemuat dari dokumen, jadi nilainya dapat dilahirkan kembali dengan
-- menjalankan ulang `pemuat -semua -ikat`. Nol nilai yang hanya hidup di
-- sini.

ALTER TABLE {skema}.T_TREATY_REVISION DROP (
  CEDINGID,
  LEADINGREINSSOURCEID,
  LEADINGREINSID,
  INFORMATION,
  POSITION,
  POSITIONUSERNAME,
  CHOOSESTATUSAKSEPTASI,
  ACCUMULATIONPERIOD,
  COMMENTTEKS,
  REPORTINGSTART,
  REPORTINGEND,
  REPORTINGPERIOD,
  REPORTINGINTERVAL,
  REPORTINGSUBMISSION,
  REPORTINGCONFIRMATION,
  REPORTINGSETTLEMENT
)
/
