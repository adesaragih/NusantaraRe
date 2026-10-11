-- 208 - T_GENERAL_POLIS.NB_STATUS (tiket 57, Submit tahap marketing): teks status kasus `pyWorkPage.NBStatus`
-- ("NB IS IN TEAM LEADER'S INBOX", "NB WAS DECLINED BY  <OPERATOR>", ...; keputusan 9d 07-10-2026, SATU kolom).
--
-- NBStatus milik halaman KERJA, bukan OfferFacIn, jadi tidak ada di rancangan (`loader/skema_gen.go`) - amandemen loader
-- (turunan: tidak diisi dari dokumen lama). Ditaruh di T_GENERAL_POLIS (milik nbfacin, berbagi PK dengan T_WORK_POLIS)
-- karena T_WORK_POLIS milik premiumlistlife (K-064). VARCHAR2(1000): teks = "NB IS IN " + MARKETINGOFFICER.CLIENTNAME
-- (VARCHAR2(100)) / M_LOGIN_GO.NAME (150) + akhiran - lebar sumber Pega tidak tertulis. Kolom BARU lewat ALTER biasa
-- (pola 195 / 203). Ditulis, TIDAK dijalankan agent. Nol COMMIT (ADR-U-0029).
ALTER TABLE {skema}.T_GENERAL_POLIS ADD (
  NB_STATUS VARCHAR2(1000)
)
/
