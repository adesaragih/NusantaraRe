-- 808 - matikan dua trigger warisan M_CLIENT (perintah work owner 04-10-2026: "TRIGGER KAMU YG TUTUP";
-- "triger dimatikan ok").
--
-- Kedua trigger (AFTER INSERT OR UPDATE ON M_CLIENT, FOR EACH ROW) menulis ulang CLIENT_ADDRESS dan CLIENT_PICLIST
-- dari dokumen M_CLIENT setiap kali Pega menyimpannya - MERGE lalu DELETE baris yang tidak ada di dokumen. Sesudah
-- Company Detail pindah ke aplikasi Go, kedua tabel itu sumber datanya sendiri: satu simpan dari Pega
-- (RDBINSERTCLIENT lewat GetInsuredID NB FacIn/RNW/Endorsment) akan menghapus PIC dan alamat buatan Go.
-- Trigger TIDAK dibuang - hanya dinonaktifkan; jalur mundur menghidupkannya lagi.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
ALTER TRIGGER {skema}.TRG_M_CLIENT DISABLE
/
ALTER TRIGGER {skema}.TRG_M_CLIENT_PIC DISABLE
/
