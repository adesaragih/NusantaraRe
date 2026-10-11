-- 258 - unggah CSV MARINE CARGO (tiket 74 P6, keputusan work owner 11-10-2026 K74-14..17) atas infrastruktur tiket 75. Rentang KEDUA
-- nbfacin 255-259. Ditulis, TIDAK dijalankan agent; work owner / DBA yang menjalankan. Nol COMMIT (ADR-U-0029).
--
-- T_NB_CSV_UNGGAH.JENIS - importer pekerjaan: 'fire' (template FIRE tiket 66 / 75) atau 'cargo' (InsertUploadMarineCargo_act);
-- pekerjaan lama (sebelum 258) = 'fire' lewat DEFAULT. KOLOM - judul kolom yang dikenali, urut berkas, dipisah ";" (cargo dibaca
-- per NAMA judul, urutan bebas - View Upload menampilkan urutan berkas); NULL = urutan template FIRE.
-- T_NB_CSV_BARIS - penampung kolom CSV cargo: satu kolom VARCHAR2(4000) per judul cargo, awalan C_ (pola 257). CURRENCY dan
-- TOP_RISK sudah ada (kolom FIRE yang sama namanya, dipakai bersama).
ALTER TABLE {skema}.T_NB_CSV_UNGGAH ADD (
  JENIS VARCHAR2(10) DEFAULT 'fire' NOT NULL,
  KOLOM VARCHAR2(4000)
)
/
ALTER TABLE {skema}.T_NB_CSV_BARIS ADD (
  C_TRADING VARCHAR2(4000),
  C_DEPARTURE VARCHAR2(4000),
  C_DESTINATION VARCHAR2(4000),
  C_GOODS VARCHAR2(4000),
  C_GOODS_NAME VARCHAR2(4000),
  C_PACKING VARCHAR2(4000),
  C_CONVEYENCE VARCHAR2(4000),
  C_CONVEYENCE_NAME VARCHAR2(4000),
  C_COVERAGE_IDX VARCHAR2(4000),
  C_COVERAGE_NAME VARCHAR2(4000),
  C_TSI VARCHAR2(4000),
  C_AGENT VARCHAR2(4000),
  C_DATASHIP VARCHAR2(4000),
  C_INVOICE VARCHAR2(4000),
  C_ETD VARCHAR2(4000),
  C_BLNUMBER VARCHAR2(4000)
)
/
