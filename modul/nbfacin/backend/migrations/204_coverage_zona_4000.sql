-- 204 - T_COVERAGELIST.ZONE VARCHAR2(1000) -> VARCHAR2(4000) (tiket 52, bug DEV 06-10-2026, diteruskan sesi `nusantarare-9d`).
--
-- Save tab Object gagal "coverages[k].zone paling banyak 1000 byte": nilai Zone (`ZONES.DESCRIPTION` terpilih; view atas
-- JSON m_zones, lebar kolom tidak tertulis) di DEV melampaui 1000 byte. Perintah work owner: "ini juga perbaiki, buat
-- 4000". 203 kemungkinan SUDAH dijalankan di DEV: tidak diubah, migrasi baru (pola 192). Pelebaran VARCHAR2 aman walau
-- kolom berisi. Amandemen loader `amandemenLebar`. Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_COVERAGELIST MODIFY (ZONE VARCHAR2(4000))
/
