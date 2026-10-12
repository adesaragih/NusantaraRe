-- 254 - T_WORK_POLIS.POSITION_NOTE (tiket 77, keputusan work owner 12-10-2026: "kalo tambahin kolom baru untuk positionnote
-- gimana? jadi di position itu tetap offer, policy, binding" -> "Ya, kerjakan"; opsi 1: kolom di T_WORK_POLIS milik premiumlistlife,
-- STRUKTUR premiumlistlife ditambah satu baris atas izin work owner 12-10-2026).
--
-- Pega punya dua properti: PositionNote (workbasket = kotak masuk) dan IsCedingConfirm (tahap). Sesudah 254, untuk case NB FacIn
-- (LINI FAC) POSITION_NOTE = workbasket pemegang case (NULL sesudah Decline) dan POSITION = tahap Offer / Binding / Policy dari
-- FLAG_ONGOING_POLICY (repository/posisinote.go). Kolom ADITIF, boleh NULL: baris lini lain tidak disentuh; baris NB lama dibaca
-- lewat posisi efektif sampai dikonversi (docs/alat/konversi_position_note_254*.sql).
-- T_WORK_POLIS lahir di premiumlistlife 050 -> 254 aman menurut urutan nama. Rentang KETIGA nbfacin 250-254 HABIS sesudah 254.
-- TIDAK wajib sebelum biner: biner mendeteksi kolomnya; tanpa 254 perilaku sama dengan sebelum tiket 77. Ditulis, TIDAK dijalankan
-- agent; DBA yang menjalankan. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_WORK_POLIS ADD (
  POSITION_NOTE VARCHAR2(255)
)
/
