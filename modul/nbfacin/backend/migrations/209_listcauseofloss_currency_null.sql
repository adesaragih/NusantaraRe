-- 209 - T_LISTCAUSEOFLOSS.CURRENCY boleh NULL (tiket 42b, sub-tab Loss Record tab Object).
--
-- Keputusan work owner 07-10-2026 (verbatim: "di tab object, sub menu loss record, Currency ini tidak wajib di isi, ganti") =
-- Pega `NB FacIn\Section\InputCauseOfLoss_FacIn.xml` sel `.Currency` pyRequired=false. Meralat N-9 / A133 (migrasi 191:
-- `DEFAULT 'UNKNOWN' NOT NULL` dari rancangan K-069). Catatan tanpa mata uang disimpan NULL - aplikasi tetap TIDAK menulis
-- 'UNKNOWN' (K-012); DEFAULT 'UNKNOWN' DIBIARKAN (aplikasi selalu menyebut kolomnya di INSERT, jadi DEFAULT tidak pernah
-- terpakai olehnya; membuang DEFAULT = perubahan kedua yang tidak diminta). Tipe ditulis ulang apa adanya supaya pembaca
-- tipe penjaga STRUKTUR tetap VARCHAR2(50). 191 sudah jalan di DEV - tidak diubah. Ditulis, TIDAK dijalankan agent. Nol
-- COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_LISTCAUSEOFLOSS MODIFY (
  CURRENCY VARCHAR2(50) NULL
)
/
