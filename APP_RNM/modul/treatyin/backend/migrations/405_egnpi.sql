-- EGNPI - EGNPI per kelompok treaty per mata uang
--
-- Migrasi tiket 23 Treaty In. Satu berkas = satu langkah migrasi.
-- Pernyataan dipisahkan oleh baris yang hanya berisi tanda garis miring.
--
-- Asal: `5-tiket/issues/23-*.md`, `2-to-spec/KAMUS-KOLOM.md`, `ddl-usulan/13_EGNPI.sql`.
--
-- Penyelarasan presisi dan skema: `docs/KEPUTUSAN-PENYELARASAN-REPO.md`.
-- INV-18 - TANPA "ON DELETE": bawaan Oracle MENOLAK, dan menolak yang
-- dikehendaki. Baris anak yang masih ada membuat penghapusan induk gagal
-- dengan ORA-02292, bukan diam-diam ikut terhapus.
--
-- INV-09: satu EGNPI per (versi, kelompok treaty, mata uang).
--
-- CONSTRAINT YANG SENGAJA TIDAK DIPASANG: NILAI_EGNPI tidak dipasangkan
-- dengan PERSEN_BAGIAN_DIPAKAI (INV-39, INV-40) - tingkat pencatatannya belum
-- ditentukan. DITAGIH: T-4 dari teknik treaty.
CREATE TABLE {skema}.EGNPI (
  ID_EGNPI            NUMBER(19)          NOT NULL,
  ID_VERSI_KONTRAK    NUMBER(19)          NOT NULL,
  ID_KELOMPOK_TREATY  NUMBER(19)          NOT NULL,
  ID_KELAS_BISNIS     NUMBER(19),
  NILAI_EGNPI         NUMBER(38,8)        NOT NULL,
  KODE_MATA_UANG      VARCHAR2(1000 CHAR) NOT NULL,
  TANGGAL_BERLAKU     DATE,
  PROPORSI            NUMBER(38,8),
  CATATAN             VARCHAR2(1000 CHAR),
  CONSTRAINT PK_EGNPI PRIMARY KEY (ID_EGNPI),
  CONSTRAINT UQ_EGNPI UNIQUE (ID_VERSI_KONTRAK, ID_KELOMPOK_TREATY, KODE_MATA_UANG),
  CONSTRAINT FK_EGNPI_1 FOREIGN KEY (ID_VERSI_KONTRAK)
    REFERENCES {skema}.VERSI_KONTRAK (ID_VERSI_KONTRAK),
  CONSTRAINT FK_EGNPI_2 FOREIGN KEY (ID_KELOMPOK_TREATY)
    REFERENCES {skema}.KELOMPOK_TREATY (ID_KELOMPOK_TREATY),
  CONSTRAINT FK_EGNPI_3 FOREIGN KEY (ID_KELAS_BISNIS)
    REFERENCES {skema}.KELAS_BISNIS (ID_KELAS_BISNIS)
)
/
CREATE INDEX {skema}.IX_EGNPI_KELOMPOK ON {skema}.EGNPI (ID_KELOMPOK_TREATY)
/
CREATE INDEX {skema}.IX_EGNPI_KELAS_BISNIS ON {skema}.EGNPI (ID_KELAS_BISNIS)
/
