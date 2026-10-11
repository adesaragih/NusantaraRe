-- 532 - T_VIEW_SUGGEST (tabel bersama milik PremiumList Life) mendapat induk kedua: klaim Claim Prop.
-- Keputusan work owner 07-10-2026 "1 tabel aja gabung life dan non life": riwayat .ClaimData.SuggestList ("Claim
-- History") satu tabel dengan riwayat penawaran PremiumList Life. Diagram sheet Claim Prop F38-F40: FK CLAIM_ID ->
-- T_GENERAL_CLAIM.ID, CASCADE. CHECK: tepat satu induk terisi. Katalog DEV 07-10-2026: seluruh baris ber-PREMIUM_LIST_ID,
-- jadi CHECK dapat dipasang tanpa membersihkan data. Baris kolomnya dicatat di STRUKTUR PremiumList Life.
-- Jejak audit Claim Prop (AC 4, 77-82): PIC_SUGGEST = pelaku, IS_CEDING_CONFIRM = tingkat wewenang, kolom terpisah.
ALTER TABLE {skema}.T_VIEW_SUGGEST ADD (
  CLAIM_ID VARCHAR2(32)
)
/
ALTER TABLE {skema}.T_VIEW_SUGGEST ADD CONSTRAINT FK_VS_CLAIM FOREIGN KEY (CLAIM_ID) REFERENCES {skema}.T_GENERAL_CLAIM (ID) ON DELETE CASCADE
/
ALTER TABLE {skema}.T_VIEW_SUGGEST ADD CONSTRAINT CK_VS_SATU_INDUK CHECK ((PREMIUM_LIST_ID IS NULL AND CLAIM_ID IS NOT NULL) OR (PREMIUM_LIST_ID IS NOT NULL AND CLAIM_ID IS NULL))
/
CREATE INDEX {skema}.IDX_VS_CLAIM ON {skema}.T_VIEW_SUGGEST (CLAIM_ID)
/
