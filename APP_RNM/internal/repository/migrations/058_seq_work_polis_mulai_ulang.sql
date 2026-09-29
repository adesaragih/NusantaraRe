-- SEQ_WORK_POLIS dimulai ulang DI ATAS nomor lama - OQ-PL-15 (GILIRAN-15).
--
-- `[DIPUTUSKAN; veto work owner]` jawaban work owner 29-09-2026 ("ikuti
-- rekomendasi"): awal SEQ_WORK_POLIS dimajukan melewati nomor kasus lama.
-- 057 memulainya dari 1, sedangkan pengenal warisan `NBLF-<n>` sudah jauh di
-- atasnya - kasus baru akan lahir dengan nomor yang SAMA dengan kasus lama di
-- Pega.
--
-- `[data DEV - brief GILIRAN-15 §0, agregat, dibaca asisten 29-09-2026;
-- perintah auditnya tidak disertakan brief dan BELUM diverifikasi executor]`
-- nomor `NBLF-` tertinggi yang TERLIHAT adalah 22373 (`JSON_POLIS` /
-- `POLICYJSONLIFE`, 33 baris ber-`NBLF-`); 22374 = tertinggi terlihat + 1.
--
-- ⚠️ Label:
-- [sementara — DBA memastikan pyLastReservedID awalan NBLF- di PC_DATA_UNIQUEID sebelum data nyata]
-- Penghitung Pega yang sebenarnya (`PC_DATA_UNIQUEID`) TIDAK terlihat dari akun
-- `POOLDATA`; "tertinggi yang terlihat" belum tentu tertinggi yang pernah
-- dicadangkan Pega - OQ-PL-17.
--
-- Penjaga kata cadangan (`katacadangan_test.go`) TIDAK disesuaikan: ia
-- memeriksa nama KOLOM, dan langkah ini tidak membuat kolom apa pun.
--
-- ⛔ DROP lalu CREATE, bukan ALTER: bentuk yang brief tetapkan, dan satu-satunya
-- yang tidak bergantung versi Oracle (`ALTER SEQUENCE … RESTART` baru ada
-- sejak 18c).
--
-- ⚠️ BILA LANGKAH INI GAGAL DI TENGAH (DROP sudah jalan, CREATE belum): DDL
-- Oracle menetapkan dirinya sendiri, dan percobaan ulang akan berhenti di
-- ORA-02289 pada DROP - jalur maju hanya memaafkan CREATE yang menjawab
-- ORA-00955. Pulihkan dengan menjalankan pernyataan CREATE di bawah secara
-- manual (DBA), lalu ulangi `-migrate`.
--
-- ⚠️ Dua kasus uji `NBLF-2` dan `NBLF-3` di `T_WORK_POLIS` DEV TIDAK disentuh
-- langkah ini (brief GILIRAN-15 §3) - penghapusannya keputusan work owner.
--
-- ⛔ NOL `COMMIT` (ADR-U-0029). `-migrate` dijalankan work owner.
DROP SEQUENCE {skema}.SEQ_WORK_POLIS
/
CREATE SEQUENCE {skema}.SEQ_WORK_POLIS START WITH 22374 INCREMENT BY 1 NOCACHE NOCYCLE
/
