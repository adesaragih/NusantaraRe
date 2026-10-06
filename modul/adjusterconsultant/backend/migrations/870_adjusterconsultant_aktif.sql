-- 870 - ADJUSTERCONSULTANT.IS_ACTIVE: flag nonaktif pengganti hapus (keputusan work owner 05-10-2026: "tambah flag
-- non aktif"). Delete di layar Pega `MstAdjusterConsultant` membuka konfirmasi hapus TREATY GROUP (salin-tempel), jadi
-- adjuster/consultant tidak pernah terhapus; di sini tidak ada hapus permanen - Deactivate mengisi '0'.
--
-- Baris yang sudah ada menjadi aktif lewat DEFAULT. Layar klaim Pega lama tidak mengenal kolom ini.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TABLE {skema}.ADJUSTERCONSULTANT ADD (
  IS_ACTIVE VARCHAR2(1) DEFAULT '1' NOT NULL
)
/
ALTER TABLE {skema}.ADJUSTERCONSULTANT ADD CONSTRAINT CK_ADJUSTERCONSULTANT_ACTIVE CHECK (IS_ACTIVE IN ('0', '1'))
/
