-- 211 - pembetulan nama workbasket posisi case NB Fac In: ReasFacInJuniorUnderwriting -> ReasFacInJuniorUnderwritingB.
--
-- Keputusan work owner 07-10-2026 (verbatim: "jika ada yang menggunakan ReasFacInJuniorUnderwriting ganti menjadi
-- ReasFacInJuniorUnderwritingB, karena kedepannya akan menggunakan itu"); ADR-0032 (pembetulan nama + tabel pemetaan lama ->
-- baru di docs/issues/58 bab "Pemetaan workbasket", penyimpangan sadar). Hanya baris case Fac In (LINI 'FAC' = LiniFacIn)
-- yang dibuat / dipindah modul ini (Accept Team Leader tiket 58, Banding JUW_B tiket 57). Riwayat HISTORYAKSEPTASIPEGA lama
-- TIDAK diubah (catatan sejarah). TGL_UPDATE tidak disentuh (bukan perubahan oleh pengguna). Ditulis, TIDAK dijalankan agent.
-- Nol COMMIT (ADR-U-0029).
UPDATE {skema}.T_WORK_POLIS SET POSITION = 'ReasFacInJuniorUnderwritingB'
 WHERE LINI = 'FAC' AND POSITION = 'ReasFacInJuniorUnderwriting'
/
