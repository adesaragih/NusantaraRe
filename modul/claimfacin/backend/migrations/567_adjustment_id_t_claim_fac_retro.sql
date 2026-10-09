-- 567 - T_CLAIM_FAC_RETRO (Claim Prop 527): di FAC retro melekat ke adjustment (keputusan K5 STRUKTUR FAC;
-- GenerateDLAFacin_Act 15.2.1.2.1 `Adjustment.FacRetroList`). CLAIM_ID tetap diisi; baris tingkat klaim (pengganti
-- FacRetroList polis, CheckLimit_Act1 9.1.1.4) ber-ADJUSTMENT_ID kosong.
-- NOL COMMIT (ADR-U-0029), nol MODIFY / DROP kolom yang sudah ada. -migrate dijalankan work owner.
ALTER TABLE {skema}.T_CLAIM_FAC_RETRO ADD (
  ADJUSTMENT_ID VARCHAR2(32)
)
/
ALTER TABLE {skema}.T_CLAIM_FAC_RETRO ADD CONSTRAINT FK_CLAIM_FAC_RETRO_ADJ FOREIGN KEY (ADJUSTMENT_ID)
  REFERENCES {skema}.T_CLAIM_ADJUSTMENT (ID) ON DELETE CASCADE
/
CREATE INDEX {skema}.IX_CLAIM_FAC_RETRO_ADJ ON {skema}.T_CLAIM_FAC_RETRO (ADJUSTMENT_ID)
/
