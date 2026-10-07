-- 181 - SEQ_WORK_POLIS_NB: penghitung nomor case NB (tiket 29), butir 74.2 dan 76.
--
-- ⛔⛔ WAJIB DIISI SEBELUM DIJALANKAN. Ganti {NB_MULAI} di bawah dengan
-- (nomor NB Pega terakhir + 1) - keputusan work owner 03-10-2026 "Lanjut dari nomor
-- terakhir Pega" (butir 74.2). Angka itu TIDAK ada di repo; diisi work owner/DBA.
-- Selama belum diganti, Oracle MENOLAK pernyataan ini (penanda bukan angka) dan
-- `-migrate` seluruh aplikasi berhenti di langkah ini - disengaja, pilihan work owner
-- (butir 76.5), supaya nomor tidak pernah diam-diam mulai dari angka yang salah.
--
-- Hanya ANGKA; awalan `NB-` dan bentuknya (tanpa nol di depan) dirakit di Go
-- (`repository/casenb.go`), pola SEQ_WORK_POLIS premiumlistlife (057/058).
-- NOCACHE: nomor tidak loncat saat instance dimulai ulang (pola 057).
CREATE SEQUENCE {skema}.SEQ_WORK_POLIS_NB START WITH 151028 INCREMENT BY 1 NOCACHE NOCYCLE
/
