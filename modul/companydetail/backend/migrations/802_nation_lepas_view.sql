-- 802 - NATION berhenti menjadi view (perintah work owner 04-10-2026: "untuk negara jangan dari enumeration, ambil
-- dari select * from NATION; sebelum itu ubah view itu jadi flat table").
--
-- Sebelumnya: NATION = view di atas M_NATION (57 baris DEV), kolom ID, OLDID, NOTE, NATIONINITIAL. Langkah ini
-- membuang view-nya; 803 membuat tabel datar bernama SAMA dengan kolom yang SAMA, 804 menyalin isinya - pembaca lama
-- (rule Pega BrowseNation_RD dan lainnya) tetap membaca nama dan kolom yang sama.
--
-- Tiga langkah terpisah supaya pelari yang gagal di tengah dapat diulang: view yang sudah dibuang tidak dibuang dua
-- kali. Jalur mundur 802 membuat ulang view aslinya.
--
-- NOL COMMIT (ADR-U-0029). -migrate dijalankan work owner.
DROP VIEW {skema}.NATION
/
