# KAMUS KOLOM — modul Claim Non Prop

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop` (279 berkas), `pengetahuan/ddl/` (49 objek), dan `MEMORI_PEMAHAMAN.MD`.
> Dibangkitkan 18 September 2026 dari definisi yang sama dengan `ddl-usulan/`. Keduanya tidak dapat berbeda.

> **Cacah berkas ini, supaya tidak dicacah keliru lagi:** **396** kolom di 22 bagian tabel, **56** kolom di 9 bagian view, **5** baris catatan — **457** baris seluruhnya. Yang dimaksud "kolom tabel" adalah **396**; angka **452** adalah 396 + 56 dan dicabut sebagai sebutan jumlah kolom tabel. Perintahnya di `sampah/keluaran/RINGKASAN-TABEL-DATAR.md` bagian 2.1.

**Kolom tanpa sumber tidak ada di sini.** Presisi merujuk kelompok domain AK-1 (`REGISTER-RATIFIKASI.md`), berstatus **provisional**.

| Kelompok | Tipe | Isi |
|---|---|---|
| U1 | `NUMBER(38,20)` | nilai uang antara |
| P1 | `NUMBER(38,20)` | persentase dan porsi |
| K1 | `NUMBER(38,20)` | kurs |
| L1 | `NUMBER(38,20)` | limit layer dan premi deposit |
| P2 | `NUMBER(11,8)` | tarif pajak dan brokerage |
| C1 | `NUMBER(9)` | cacah dan nomor urut |

---

## 1. `KLAIM`

Aggregate root. Satu klaim satu kejadian kerugian.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | ADR-0021 |
| `NOMOR_KLAIM` | `VARCHAR2(32 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 (.NoClaim/.ClaimNo); ADR-0021 |
| `NOMOR_KLAIM_CEDANT` | `VARCHAR2(64 CHAR)` | — | Y | cedant belum memberi nomornya | MEMORI_PEMAHAMAN.MD §4.2 (.CNPClmNoCedant) |
| `NOMOR_POLIS` | `VARCHAR2(64 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 (.PolicyNo); ADR-0021 |
| `NOMOR_POLIS_CEDANT` | `VARCHAR2(64 CHAR)` | — | Y | tidak dicatat saat registrasi | ddl TABLE_PC_ASM_FW_GCNMFW_WORK (POLICY_CEDING); ADR-0021 |
| `ID_TREATY` | `VARCHAR2(64 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 (.IDMaster) |
| `NAMA_TERTANGGUNG` | `VARCHAR2(255 CHAR)` | — | Y | belum diisi | MEMORI_PEMAHAMAN.MD §4.2 (.InsuredName) |
| `LINI_USAHA` | `VARCHAR2(40 CHAR)` | — | Y | **belum ditetapkan** — bukan "konvensional". Penentu lamanya (node yang mengeksekusi, nama akun notifikasi) tidak ikut pindah | SPEC §21.3; G1; D33, D42 — **tanpa CHECK domain**, nilainya belum terbaca |
| `TANGGAL_KEJADIAN` | `DATE` | — | N | — | BLUEPRINT.md §2.2; ADR-0001, ADR-0022 |
| `TANGGAL_LAPOR` | `DATE` | — | Y | belum diisi | MEMORI_PEMAHAMAN.MD §4.2 (.ReportDate) |
| `TANGGAL_TERIMA` | `DATE` | — | Y | belum diisi | MEMORI_PEMAHAMAN.MD §4.2 (.DateReceived) |
| `TREATY_MULAI` | `DATE` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 (.StartDateTreaty); ADR-0022 |
| `TREATY_AKHIR` | `DATE` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 (.EndDateTreaty); ADR-0022, FINDING-005 |
| `PENYEBAB_KERUGIAN_ID` | `VARCHAR2(50 CHAR)` | — | Y | belum dipilih | MEMORI_PEMAHAMAN.MD §4.2 (.CauseOfLossID) |
| `PORSI_CEDANT` | `NUMBER(38,20)` | P1 | Y | belum diisi | MEMORI_PEMAHAMAN.MD §4.2 (.ShareCeding) |
| `PORSI_REASURADUR` | `NUMBER(38,20)` | P1 | Y | belum diisi | MEMORI_PEMAHAMAN.MD §4.5 (RNMShare) |
| `STATUS_KLAIM` | `VARCHAR2(40 CHAR)` | — | N | — | BLUEPRINT.md §3.1 |
| `PENUTUPAN_LAMA` | `VARCHAR2(64 CHAR)` | — | Y | **tidak ditutup di sistem lama** — bukan tidak diketahui | SPEC §21.6; B9. Menandai klaim yang tertutup lewat jalur yang tidak dimigrasi, supaya asal-usulnya terbaca tanpa jalurnya ikut hidup |
| `PENANDA_BERSYARAT` | `NUMBER(1)` | — | N | — | BLUEPRINT.md §8.3 — field pendaratan, penulis belum dikonfirmasi. **Domain biner ditetapkan di sini, bukan diwarisi**: SPEC §21.1 melarang memberi domain biner kepada penanda lama tanpa bukti, dan nilai lamanya masuk lewat pemetaan, bukan penyalinan mentah |
| `CATATAN_BERSYARAT` | `VARCHAR2(2000 CHAR)` | — | Y | tidak ada syarat | BLUEPRINT.md §8.3 |
| `KEADAAN_BARIS` | `VARCHAR2(16 CHAR)` | — | N | — | ADR-0019 lapis 2; ADR-0014 penularan |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 2. `NILAI_KLAIM_MATA_UANG`

Nilai kerugian per mata uang. Masukan, diisi pengguna.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_NILAI_KLAIM` | `NUMBER(19)` | — | N | — | BLUEPRINT.md §2.1 |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | BLUEPRINT.md §2.1 |
| `MATA_UANG` | `VARCHAR2(3 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2; ADR-0007 |
| `KURS` | `NUMBER(38,20)` | K1 | Y | kurs belum ada — nilai IDR tidak dihitung | FINDING-006; ADR-0014 |
| `KURS_TANGGAL` | `DATE` | — | Y | idem | ADR-0007 |
| `KURS_SUMBER` | `VARCHAR2(32 CHAR)` | — | Y | idem | ADR-0007; FINDING-006 |
| `NILAI_KERUGIAN` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | MEMORI_PEMAHAMAN.MD §4.2 ListClaimAmount.ClaimAmount |
| `NILAI_KERUGIAN_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | MEMORI_PEMAHAMAN.MD §4.2 ListClaimAmount.ClaimAmount; ADR-0007 |
| `BIAYA_PENILAIAN` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | MEMORI_PEMAHAMAN.MD §4.2 ListClaimAmount.AdjusterFee |
| `BIAYA_PENILAIAN_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | MEMORI_PEMAHAMAN.MD §4.2 ListClaimAmount.AdjusterFee; ADR-0007 |
| `SALVAGE` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | MEMORI_PEMAHAMAN.MD §4.2 ListClaimAmount.Salvage |
| `SALVAGE_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | MEMORI_PEMAHAMAN.MD §4.2 ListClaimAmount.Salvage; ADR-0007 |
| `BIAYA_LAIN` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | MEMORI_PEMAHAMAN.MD §4.2 ListClaimAmount.CNPOthersFee |
| `BIAYA_LAIN_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | MEMORI_PEMAHAMAN.MD §4.2 ListClaimAmount.CNPOthersFee; ADR-0007 |
| `PRORATA_KLAIM` | `NUMBER(38,20)` | P1 | Y | belum diisi | MEMORI_PEMAHAMAN.MD §4.2 ListClaimAmount.PctProrateClaim |
| `NILAI_RETENSI_CEDANT` | `NUMBER(38,20)` | U1 | Y | belum diisi | MEMORI_PEMAHAMAN.MD §4.2 ListClaimAmount.CNPDeductible |
| `KEADAAN_BARIS` | `VARCHAR2(16 CHAR)` | — | N | — | ADR-0019 lapis 2; ADR-0014 penularan |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 3. `ALOKASI_LAYER`

Hasil alokasi kerugian per layer. Tanpa baris Retensi Cedant.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_ALOKASI` | `NUMBER(19)` | — | N | — | BLUEPRINT.md §2.4 |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | BLUEPRINT.md §2.4 |
| `LAYER` | `VARCHAR2(50 CHAR)` | — | N | — | BLUEPRINT.md §2.4; ddl TABLE_PROPORTIONALARRG |
| `LAYER_JENIS` | `VARCHAR2(50 CHAR)` | — | N | — | BLUEPRINT.md §2.4 |
| `LAYER_BAGIAN` | `VARCHAR2(50 CHAR)` | — | N | — | BLUEPRINT.md §2.4 |
| `LAYER_BAGIAN_JENIS` | `VARCHAR2(50 CHAR)` | — | N | — | BLUEPRINT.md §2.4 |
| `JENIS_REASURANSI` | `VARCHAR2(100 CHAR)` | — | Y | lookup REINSURANCETYPE tidak menemukan baris — belum ditelusuri apakah itu memicu galat | SPEC bagian 0; ddl TABLE_REINSURANCETYPE (ID, NOTE) |
| `JENIS_REASURANSI_ID` | `VARCHAR2(50 CHAR)` | — | Y | idem | SPEC bagian 0 |
| `MATA_UANG` | `VARCHAR2(3 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2; ADR-0007 |
| `KURS` | `NUMBER(38,20)` | K1 | Y | kurs belum ada — nilai IDR tidak dihitung | FINDING-006; ADR-0014 |
| `KURS_TANGGAL` | `DATE` | — | Y | idem | ADR-0007 |
| `KURS_SUMBER` | `VARCHAR2(32 CHAR)` | — | Y | idem | ADR-0007; FINDING-006 |
| `NILAI_TERALOKASI_HITUNG` | `NUMBER(38,20)` | U1 | Y | belum pernah dihitung | MEMORI_PEMAHAMAN.MD §4.2 SpreadingRisk |
| `NILAI_TERALOKASI_SUNTING` | `NUMBER(38,20)` | U1 | Y | tidak pernah disunting — keadaan normal | ADR-0008 |
| `NILAI_TERALOKASI_DIPAKAI` | `NUMBER(38,20) GENERATED ALWAYS AS (COALESCE(NILAI_TERALOKASI_SUNTING, NILAI_TERALOKASI_HITUNG)) VIRTUAL` | U1 | Y | kedua-duanya kosong | ADR-0008 — kolom turunan, tidak dapat ditulis |
| `BIAYA_PENILAIAN_HITUNG` | `NUMBER(38,20)` | U1 | Y | belum pernah dihitung | MEMORI_PEMAHAMAN.MD §4.2 SpreadingRisk |
| `BIAYA_PENILAIAN_SUNTING` | `NUMBER(38,20)` | U1 | Y | tidak pernah disunting — keadaan normal | ADR-0008 |
| `BIAYA_PENILAIAN_DIPAKAI` | `NUMBER(38,20) GENERATED ALWAYS AS (COALESCE(BIAYA_PENILAIAN_SUNTING, BIAYA_PENILAIAN_HITUNG)) VIRTUAL` | U1 | Y | kedua-duanya kosong | ADR-0008 — kolom turunan, tidak dapat ditulis |
| `SALVAGE_HITUNG` | `NUMBER(38,20)` | U1 | Y | belum pernah dihitung | MEMORI_PEMAHAMAN.MD §4.2 SpreadingRisk |
| `SALVAGE_SUNTING` | `NUMBER(38,20)` | U1 | Y | tidak pernah disunting — keadaan normal | ADR-0008 |
| `SALVAGE_DIPAKAI` | `NUMBER(38,20) GENERATED ALWAYS AS (COALESCE(SALVAGE_SUNTING, SALVAGE_HITUNG)) VIRTUAL` | U1 | Y | kedua-duanya kosong | ADR-0008 — kolom turunan, tidak dapat ditulis |
| `BIAYA_LAIN_HITUNG` | `NUMBER(38,20)` | U1 | Y | belum pernah dihitung | MEMORI_PEMAHAMAN.MD §4.2 SpreadingRisk |
| `BIAYA_LAIN_SUNTING` | `NUMBER(38,20)` | U1 | Y | tidak pernah disunting — keadaan normal | ADR-0008 |
| `BIAYA_LAIN_DIPAKAI` | `NUMBER(38,20) GENERATED ALWAYS AS (COALESCE(BIAYA_LAIN_SUNTING, BIAYA_LAIN_HITUNG)) VIRTUAL` | U1 | Y | kedua-duanya kosong | ADR-0008 — kolom turunan, tidak dapat ditulis |
| `DISUNTING_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah disunting | ADR-0008 |
| `DISUNTING_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DISUNTING_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah disunting | ADR-0008 |
| `ALASAN_SUNTING` | `VARCHAR2(500 CHAR)` | — | Y | idem — tidak pernah dibaca mesin | ADR-0008, ADR-0009 |
| `NILAI_TERALOKASI_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi | ADR-0007 |
| `PORSI_REASURADUR` | `NUMBER(38,20)` | P1 | Y | belum di-snapshot | MEMORI_PEMAHAMAN.MD §4.5 (RNMShare) |
| `PORSI_CEDANT` | `NUMBER(38,20)` | P1 | Y | belum di-snapshot | MEMORI_PEMAHAMAN.MD §4.2 (.ShareCeding) |
| `PRORATA_KLAIM` | `NUMBER(38,20)` | P1 | Y | belum di-snapshot | MEMORI_PEMAHAMAN.MD §4.2 |
| `LIMIT_LAYER` | `NUMBER(38,20)` | L1 | Y | belum di-snapshot | MEMORI_PEMAHAMAN.MD §4.5 Limits[].Limit |
| `PREMI_DEPOSIT` | `NUMBER(38,20)` | L1 | Y | belum di-snapshot | MEMORI_PEMAHAMAN.MD §4.5 Limits[].MDPList |
| `PERSEN_PREMI_PEMULIHAN` | `NUMBER(38,20)` | P1 | Y | belum di-snapshot | MEMORI_PEMAHAMAN.MD §4.5 Limits[].ReinstatementPct |
| `URUTAN_PENGISIAN` | `NUMBER(9)` | C1 | N | — | MEMORI_PEMAHAMAN.MD §6.2 Mode 3; ADR-0011 |
| `KEADAAN_BARIS` | `VARCHAR2(16 CHAR)` | — | N | — | ADR-0019 lapis 2; ADR-0014 penularan |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 4. `RETENSI_CEDANT`

Bagian pertama kerugian yang ditanggung cedant. Tanpa kolom jenis reasuransi.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_RETENSI` | `NUMBER(19)` | — | N | — | ADR-0010 |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | ADR-0010 |
| `MATA_UANG` | `VARCHAR2(3 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2; ADR-0007 |
| `KURS` | `NUMBER(38,20)` | K1 | Y | kurs belum ada — nilai IDR tidak dihitung | FINDING-006; ADR-0014 |
| `KURS_TANGGAL` | `DATE` | — | Y | idem | ADR-0007 |
| `KURS_SUMBER` | `VARCHAR2(32 CHAR)` | — | Y | idem | ADR-0007; FINDING-006 |
| `NILAI_RETENSI_HITUNG` | `NUMBER(38,20)` | U1 | Y | belum pernah dihitung | BLUEPRINT.md §2.4; ADR-0010 |
| `NILAI_RETENSI_SUNTING` | `NUMBER(38,20)` | U1 | Y | tidak pernah disunting — keadaan normal | ADR-0008 |
| `NILAI_RETENSI_DIPAKAI` | `NUMBER(38,20) GENERATED ALWAYS AS (COALESCE(NILAI_RETENSI_SUNTING, NILAI_RETENSI_HITUNG)) VIRTUAL` | U1 | Y | kedua-duanya kosong | ADR-0008 — kolom turunan, tidak dapat ditulis |
| `DISUNTING_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah disunting | ADR-0008 |
| `DISUNTING_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DISUNTING_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah disunting | ADR-0008 |
| `ALASAN_SUNTING` | `VARCHAR2(500 CHAR)` | — | Y | idem — tidak pernah dibaca mesin | ADR-0008, ADR-0009 |
| `NILAI_RETENSI_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi | ADR-0007 |
| `PRORATA_KLAIM` | `NUMBER(38,20)` | P1 | Y | belum di-snapshot | MEMORI_PEMAHAMAN.MD §6.2 |
| `KEADAAN_BARIS` | `VARCHAR2(16 CHAR)` | — | N | — | ADR-0019 lapis 2; ADR-0014 penularan |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 5. `AKSEPTASI`

Satu entitas dengan keadaan. Satu per klaim per layer per mata uang.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_AKSEPTASI` | `NUMBER(19)` | — | N | — | ADR-0015 |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | ADR-0024 |
| `LAYER` | `VARCHAR2(50 CHAR)` | — | N | — | BLUEPRINT.md §2.4; ddl TABLE_PROPORTIONALARRG |
| `LAYER_JENIS` | `VARCHAR2(50 CHAR)` | — | N | — | BLUEPRINT.md §2.4 |
| `LAYER_BAGIAN` | `VARCHAR2(50 CHAR)` | — | N | — | BLUEPRINT.md §2.4 |
| `LAYER_BAGIAN_JENIS` | `VARCHAR2(50 CHAR)` | — | N | — | BLUEPRINT.md §2.4 |
| `JENIS_REASURANSI` | `VARCHAR2(100 CHAR)` | — | Y | lookup REINSURANCETYPE tidak menemukan baris — belum ditelusuri apakah itu memicu galat | SPEC bagian 0; ddl TABLE_REINSURANCETYPE (ID, NOTE) |
| `JENIS_REASURANSI_ID` | `VARCHAR2(50 CHAR)` | — | Y | idem | SPEC bagian 0 |
| `MATA_UANG` | `VARCHAR2(3 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2; ADR-0007 |
| `KURS` | `NUMBER(38,20)` | K1 | Y | kurs belum ada — nilai IDR tidak dihitung | FINDING-006; ADR-0014 |
| `KURS_TANGGAL` | `DATE` | — | Y | idem | ADR-0007 |
| `KURS_SUMBER` | `VARCHAR2(32 CHAR)` | — | Y | idem | ADR-0007; FINDING-006 |
| `NILAI_DISETUJUI` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | ddl TABLE_OS_AKSEPTASI_KLAIM; ADR-0024 |
| `NILAI_DISETUJUI_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | ddl TABLE_OS_AKSEPTASI_KLAIM; ADR-0024; ADR-0007 |
| `KEADAAN_AKSEPTASI` | `VARCHAR2(24 CHAR)` | — | N | — | ADR-0015, ADR-0009 |
| `NOMOR_AKSEPTASI` | `VARCHAR2(64 CHAR)` | — | Y | belum terbit — belum lolos seluruh jenjang | MEMORI_PEMAHAMAN.MD §6.8; CONTEXT.md |
| `TANGGAL_AKSEPTASI` | `DATE` | — | Y | belum terbit | BLUEPRINT.md §8.3 |
| `KEADAAN_BARIS` | `VARCHAR2(16 CHAR)` | — | N | — | ADR-0019 lapis 2; ADR-0014 penularan |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 6. `ADJUSTMENT`

Unit transaksi pembayaran klaim. Satu klaim banyak Adjustment.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_ADJUSTMENT` | `NUMBER(19)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.3 |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.3 |
| `NOMOR_URUT` | `NUMBER(9)` | C1 | N | — | BLUEPRINT.md §8.4 — pengganti IndexObject |
| `ID_REKENING` | `NUMBER(19)` | — | Y | penerima belum ditentukan | MEMORI_PEMAHAMAN.MD §4.2 ReceiverClaim |
| `JENIS_PEMBAYARAN` | `VARCHAR2(24 CHAR)` | — | Y | belum ditentukan. Arti tiap nilainya EXTERNAL — A6b, REQ-027 | MEMORI_PEMAHAMAN.MD §4.3; §7.4 |
| `MATA_UANG` | `VARCHAR2(3 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2; ADR-0007 |
| `KURS` | `NUMBER(38,20)` | K1 | Y | kurs belum ada — nilai IDR tidak dihitung | FINDING-006; ADR-0014 |
| `KURS_TANGGAL` | `DATE` | — | Y | idem | ADR-0007 |
| `KURS_SUMBER` | `VARCHAR2(32 CHAR)` | — | Y | idem | ADR-0007; FINDING-006 |
| `NILAI_ADJUSTMENT` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | MEMORI_PEMAHAMAN.MD §4.3 (.ValueAdjustment) — mata uangnya EXTERNAL: H1/H2 |
| `NILAI_ADJUSTMENT_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | MEMORI_PEMAHAMAN.MD §4.3 (.ValueAdjustment) — mata uangnya EXTERNAL: H1/H2; ADR-0007 |
| `KEPUTUSAN_KOMITE` | `VARCHAR2(16 CHAR)` | — | **Y** | **belum pernah dikirim ke komite** — berbeda dari `'0'` yang berarti sudah dikirim dan belum diputus. Keadaan **keempat**, dan justru yang lolos penjaga `CloseClaimMD` (D2) | BLUEPRINT.md §8.3 — field pendaratan. Tanpa CHECK: domainnya EXTERNAL (I1/I2) |
| `NOMOR_AKSEPTASI` | `VARCHAR2(64 CHAR)` | — | Y | belum terbit | BLUEPRINT.md §8.3 |
| `TANGGAL_AKSEPTASI` | `DATE` | — | Y | belum terbit | BLUEPRINT.md §8.3 |
| `PENANDA_BERSYARAT` | `NUMBER(1)` | — | N | — | BLUEPRINT.md §8.3 |
| `CATATAN_BERSYARAT` | `VARCHAR2(2000 CHAR)` | — | Y | tidak ada syarat | BLUEPRINT.md §8.3 |
| `USUL_TUTUP` | `NUMBER(1)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.3 (.IsProposeClose) |
| `LANGSUNG_KE_KASIR` | `NUMBER(1)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.3 (.DirectToKasir) |
| `KEADAAN_BARIS` | `VARCHAR2(16 CHAR)` | — | N | — | ADR-0019 lapis 2; ADR-0014 penularan |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 7. `PREMI_PEMULIHAN`

Premi untuk memulihkan kapasitas layer, beserta seluruh masukan rumusnya.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_PREMI_PEMULIHAN` | `NUMBER(19)` | — | N | — | MEMORI_PEMAHAMAN.MD §6.3 |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | MEMORI_PEMAHAMAN.MD §6.3 |
| `LAYER` | `VARCHAR2(50 CHAR)` | — | N | — | BLUEPRINT.md §2.4; ddl TABLE_PROPORTIONALARRG |
| `LAYER_JENIS` | `VARCHAR2(50 CHAR)` | — | N | — | BLUEPRINT.md §2.4 |
| `LAYER_BAGIAN` | `VARCHAR2(50 CHAR)` | — | N | — | BLUEPRINT.md §2.4 |
| `LAYER_BAGIAN_JENIS` | `VARCHAR2(50 CHAR)` | — | N | — | BLUEPRINT.md §2.4 |
| `MATA_UANG` | `VARCHAR2(3 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2; ADR-0007 |
| `KURS` | `NUMBER(38,20)` | K1 | Y | kurs belum ada — nilai IDR tidak dihitung | FINDING-006; ADR-0014 |
| `KURS_TANGGAL` | `DATE` | — | Y | idem | ADR-0007 |
| `KURS_SUMBER` | `VARCHAR2(32 CHAR)` | — | Y | idem | ADR-0007; FINDING-006 |
| `NILAI_DASAR` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | MEMORI_PEMAHAMAN.MD §6.3 ReinstatementList.IDR |
| `NILAI_DASAR_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | MEMORI_PEMAHAMAN.MD §6.3 ReinstatementList.IDR; ADR-0007 |
| `BIAYA_PENILAIAN` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | MEMORI_PEMAHAMAN.MD §6.3 |
| `BIAYA_PENILAIAN_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | MEMORI_PEMAHAMAN.MD §6.3; ADR-0007 |
| `SALVAGE_DISESUAIKAN` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | MEMORI_PEMAHAMAN.MD §6.3 (Salvage x 100 / RNMShare) |
| `SALVAGE_DISESUAIKAN_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | MEMORI_PEMAHAMAN.MD §6.3 (Salvage x 100 / RNMShare); ADR-0007 |
| `LIMIT_LAYER` | `NUMBER(38,20)` | L1 | N | — | MEMORI_PEMAHAMAN.MD §6.3 (CNPLimit) |
| `PREMI_DEPOSIT` | `NUMBER(38,20)` | L1 | N | — | MEMORI_PEMAHAMAN.MD §6.3 (CNPMDP) |
| `PERSEN_PREMI_PEMULIHAN` | `NUMBER(38,20)` | P1 | N | — | MEMORI_PEMAHAMAN.MD §6.3 (CNPPctReinstate) |
| `PORSI_REASURADUR` | `NUMBER(38,20)` | P1 | N | — | MEMORI_PEMAHAMAN.MD §6.3 (CNPShareRNM) |
| `PREMI_PEMULIHAN` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | MEMORI_PEMAHAMAN.MD §6.3 (CNPReinstatePremium); FINDING-007 |
| `PREMI_PEMULIHAN_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | MEMORI_PEMAHAMAN.MD §6.3 (CNPReinstatePremium); FINDING-007; ADR-0007 |
| `PREMI_PEMULIHAN_PORSI` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | MEMORI_PEMAHAMAN.MD §6.3 (CNPReinsPremiRNM) |
| `PREMI_PEMULIHAN_PORSI_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | MEMORI_PEMAHAMAN.MD §6.3 (CNPReinsPremiRNM); ADR-0007 |
| `KEADAAN_BARIS` | `VARCHAR2(16 CHAR)` | — | N | — | ADR-0019 lapis 2; ADR-0014 penularan |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 8. `REKENING_PENERIMA`

Salinan keadaan rekening pada saat pembayaran. Pengecualian sah ADR-0023.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_REKENING` | `NUMBER(19)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 ReceiverClaim |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 |
| `NAMA_PENERIMA` | `VARCHAR2(255 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 (.Name) |
| `NAMA_BANK` | `VARCHAR2(128 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 (.NameOfBank) |
| `CABANG_BANK` | `VARCHAR2(128 CHAR)` | — | Y | tidak dicatat | MEMORI_PEMAHAMAN.MD §4.2 (.BranchOfBank) |
| `NOMOR_REKENING` | `VARCHAR2(64 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 (.NoAccount) |
| `KODE_SWIFT` | `VARCHAR2(16 CHAR)` | — | Y | bank dalam negeri | MEMORI_PEMAHAMAN.MD §4.2 (.SwiftCode) |
| `MATA_UANG` | `VARCHAR2(3 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 9. `OBJEK_PERTANGGUNGAN`

Objek yang dipertanggungkan. Hanya InterestList; ObjectList masuk daftar lubang.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_OBJEK` | `NUMBER(19)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 InterestList |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 |
| `NOMOR_URUT` | `NUMBER(9)` | C1 | N | — | MEMORI_PEMAHAMAN.MD §4.2 |
| `NAMA_OBJEK` | `VARCHAR2(255 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 (.ObjectName) |
| `MATA_UANG` | `VARCHAR2(3 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2; ADR-0007 |
| `KURS` | `NUMBER(38,20)` | K1 | Y | kurs belum ada — nilai IDR tidak dihitung | FINDING-006; ADR-0014 |
| `KURS_TANGGAL` | `DATE` | — | Y | idem | ADR-0007 |
| `KURS_SUMBER` | `VARCHAR2(32 CHAR)` | — | Y | idem | ADR-0007; FINDING-006 |
| `NILAI_PERTANGGUNGAN` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | MEMORI_PEMAHAMAN.MD §4.2 (.TSIPerObject) |
| `NILAI_PERTANGGUNGAN_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | MEMORI_PEMAHAMAN.MD §4.2 (.TSIPerObject); ADR-0007 |
| `KEADAAN_BARIS` | `VARCHAR2(16 CHAR)` | — | N | — | ADR-0019 lapis 2; ADR-0014 penularan |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 10. `PEMBAGIAN_KERUGIAN`

Masukan mesin alokasi, bukan keluarannya.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_PEMBAGIAN` | `NUMBER(19)` | — | N | — | MEMORI_PEMAHAMAN.MD §6.2 Langkah 2 |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | BLUEPRINT.md §2.1 CNPSpreadLoss |
| `NOMOR_URUT` | `NUMBER(9)` | C1 | N | — | BLUEPRINT.md §2.1 |
| `PERSEN_KERUGIAN` | `NUMBER(38,20)` | P1 | N | — | BLUEPRINT.md §2.1 (.ClaimPercentage) |
| `MATA_UANG` | `VARCHAR2(3 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2; ADR-0007 |
| `KURS` | `NUMBER(38,20)` | K1 | Y | kurs belum ada — nilai IDR tidak dihitung | FINDING-006; ADR-0014 |
| `KURS_TANGGAL` | `DATE` | — | Y | idem | ADR-0007 |
| `KURS_SUMBER` | `VARCHAR2(32 CHAR)` | — | Y | idem | ADR-0007; FINDING-006 |
| `NILAI_KERUGIAN` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | BLUEPRINT.md §2.1 (.ClaimAmountAdjust) |
| `NILAI_KERUGIAN_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | BLUEPRINT.md §2.1 (.ClaimAmountAdjust); ADR-0007 |
| `BIAYA_PENILAIAN` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | BLUEPRINT.md §2.1 |
| `BIAYA_PENILAIAN_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | BLUEPRINT.md §2.1; ADR-0007 |
| `SALVAGE` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | BLUEPRINT.md §2.1 |
| `SALVAGE_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | BLUEPRINT.md §2.1; ADR-0007 |
| `BIAYA_LAIN` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | BLUEPRINT.md §2.1 |
| `BIAYA_LAIN_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | BLUEPRINT.md §2.1; ADR-0007 |
| `KEADAAN_BARIS` | `VARCHAR2(16 CHAR)` | — | N | — | ADR-0019 lapis 2; ADR-0014 penularan |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 11. `PENYEBARAN`

Pembagian bagian reasuradur kepada pihak penerima berikutnya.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_PENYEBARAN` | `NUMBER(19)` | — | N | — | CONTEXT.md Penyebaran |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | BLUEPRINT.md §2.1 SpreadingClaim/SpreadingBreakQS |
| `NOMOR_URUT` | `NUMBER(9)` | C1 | N | — | BLUEPRINT.md §2.1 |
| `JENIS_REASURANSI` | `VARCHAR2(100 CHAR)` | — | N | — | CountLossAllocation_act Langkah 10 |
| `JENIS_REASURANSI_ID` | `VARCHAR2(50 CHAR)` | — | Y | lookup tidak menemukan baris | idem |
| `PERSEN_PENYEBARAN` | `NUMBER(38,20)` | P1 | N | — | CountLossAllocation_act Langkah 10 (Pct) |
| `MATA_UANG` | `VARCHAR2(3 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2; ADR-0007 |
| `KURS` | `NUMBER(38,20)` | K1 | Y | kurs belum ada — nilai IDR tidak dihitung | FINDING-006; ADR-0014 |
| `KURS_TANGGAL` | `DATE` | — | Y | idem | ADR-0007 |
| `KURS_SUMBER` | `VARCHAR2(32 CHAR)` | — | Y | idem | ADR-0007; FINDING-006 |
| `NILAI_TERSEBAR` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | BLUEPRINT.md §2.1 (.ClaimSpreaded) |
| `NILAI_TERSEBAR_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | BLUEPRINT.md §2.1 (.ClaimSpreaded); ADR-0007 |
| `BIAYA_PENILAIAN` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | BLUEPRINT.md §2.1 |
| `BIAYA_PENILAIAN_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | BLUEPRINT.md §2.1; ADR-0007 |
| `SALVAGE` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | BLUEPRINT.md §2.1 |
| `SALVAGE_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | BLUEPRINT.md §2.1; ADR-0007 |
| `BIAYA_LAIN` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | BLUEPRINT.md §2.1 |
| `BIAYA_LAIN_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | BLUEPRINT.md §2.1; ADR-0007 |
| `KEADAAN_BARIS` | `VARCHAR2(16 CHAR)` | — | N | — | ADR-0019 lapis 2; ADR-0014 penularan |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 12. `ESTIMASI_AWAL`

Estimasi kerugian awal per jenis reasuransi per mata uang.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_ESTIMASI` | `NUMBER(19)` | — | N | — | BLUEPRINT.md §2.1 EstimationList |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | BLUEPRINT.md §2.1 |
| `JENIS_REASURANSI` | `VARCHAR2(100 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 (.TypeLoss); SPEC bagian 0 |
| `MATA_UANG` | `VARCHAR2(3 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2; ADR-0007 |
| `KURS` | `NUMBER(38,20)` | K1 | Y | kurs belum ada — nilai IDR tidak dihitung | FINDING-006; ADR-0014 |
| `KURS_TANGGAL` | `DATE` | — | Y | idem | ADR-0007 |
| `KURS_SUMBER` | `VARCHAR2(32 CHAR)` | — | Y | idem | ADR-0007; FINDING-006 |
| `NILAI_ESTIMASI` | `NUMBER(38,20)` | U1 | Y | belum dihitung — bukan nol | MEMORI_PEMAHAMAN.MD §4.2 (.TotalEstimasi) |
| `NILAI_ESTIMASI_IDR` | `NUMBER(38,20)` | U1 | Y | belum dikonversi; lihat KEADAAN_BARIS | MEMORI_PEMAHAMAN.MD §4.2 (.TotalEstimasi); ADR-0007 |
| `KEADAAN_BARIS` | `VARCHAR2(16 CHAR)` | — | N | — | ADR-0019 lapis 2; ADR-0014 penularan |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 13. `KRONOLOGI_KLAIM`

Catatan berurutan tindakan penting. Tidak pernah dibaca mesin.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_KRONOLOGI` | `NUMBER(19)` | — | N | — | CONTEXT.md Kronologi |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | MEMORI_PEMAHAMAN.MD §4.2 SuggestList |
| `NOMOR_URUT` | `NUMBER(9)` | C1 | N | — | MEMORI_PEMAHAMAN.MD §4.2 |
| `JENIS_TINDAKAN` | `VARCHAR2(40 CHAR)` | — | N | — | ADR-0009 — terstruktur, bukan teks bebas |
| `KETERANGAN` | `VARCHAR2(2000 CHAR)` | — | Y | tidak ada keterangan | ADR-0009 — tidak pernah dibaca mesin |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 14. `DOKUMEN_KLAIM`

Dokumen dirujuk, tidak disimpan. Berkasnya di Google Cloud Storage.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_DOKUMEN` | `NUMBER(19)` | — | N | — | ADR-0027 |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | ADR-0027 |
| `NOMOR_URUT` | `NUMBER(9)` | C1 | N | — | BLUEPRINT.md §2.1 Attachment |
| `NAMA_BERKAS` | `VARCHAR2(255 CHAR)` | — | N | — | MEMORI_PEMAHAMAN.MD §5.7 |
| `JENIS_MEDIA` | `VARCHAR2(128 CHAR)` | — | Y | tidak terdeteksi | InsertGoogleStorage_Act |
| `URL` | `VARCHAR2(2000 CHAR)` | — | Y | URL belum pernah diterbitkan | ADR-0027; T_STORAGE_IMAGE |
| `URL_KEDALUWARSA` | `TIMESTAMP(6)` | — | Y | tidak ada URL | ADR-0027 (EXPDATE) |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 15. `TUTUP_BUKU`

Tanggal tutup buku bertanggal berlaku. Satu sumber kebenaran.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_TUTUP_BUKU` | `NUMBER(19)` | — | N | — | ADR-0025 |
| `LINGKUP` | `VARCHAR2(40 CHAR)` | — | N | — | ADR-0025 — disediakan sekarang meski belum diperlukan |
| `BERLAKU_SEJAK` | `DATE` | — | N | — | ADR-0025 |
| `HARI_TUTUP_BUKU` | `NUMBER(9)` | C1 | N | — | AK-6.3 — HitServiceToKasir_Act (25); PROC_GENERATE_SEQUENCE_NUMBER |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 16. `TARIF_BERLAKU`

Tarif pajak dan brokerage sebagai data bertanggal berlaku.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_TARIF` | `NUMBER(19)` | — | N | — | AK-5; ADR-0004 |
| `LINGKUP` | `VARCHAR2(40 CHAR)` | — | N | — | AK-5 — baris pertama bernilai global |
| `JENIS_TARIF` | `VARCHAR2(24 CHAR)` | — | N | — | BLUEPRINT.md §6.4 SetPPNPPH |
| `BERLAKU_SEJAK` | `DATE` | — | N | — | AK-5; preseden ADR-0025 |
| `NILAI_TARIF_PERSEN` | `NUMBER(11,8)` | P2 | N | — | BLUEPRINT.md §6.4 — skala 8 |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 17. `KOREKSI_NILAI`

Pengganti 29 langkah tambalan per-case. Koreksi bernilai tercatat.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_KOREKSI` | `NUMBER(19)` | — | N | — | ADR-0004 |
| `TABEL_SASARAN` | `VARCHAR2(30 CHAR)` | — | N | — | ADR-0004 |
| `ID_SASARAN` | `NUMBER(19)` | — | N | — | ADR-0004 |
| `KOLOM_SASARAN` | `VARCHAR2(30 CHAR)` | — | N | — | ADR-0004 |
| `NILAI_SEBELUM` | `NUMBER(38,20)` | U1 | Y | tidak ada nilai sebelumnya — koreksi mengisi yang semula kosong | ADR-0004 |
| `NILAI_SESUDAH` | `NUMBER(38,20)` | U1 | Y | koreksi membatalkan nilai, bukan menggantinya | ADR-0004 |
| `ALASAN` | `VARCHAR2(2000 CHAR)` | — | N | — | ADR-0004, ADR-0009 |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 18. `KLAIM_PENJAGA_TANGGAL`

Klaim terdampak penjaga tanggal yang salah. Daftar menetap.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_PENJAGA_TANGGAL` | `NUMBER(19)` | — | N | — | ADR-0018 |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | ADR-0018 |
| `SEBAB` | `VARCHAR2(255 CHAR)` | — | N | — | FINDING-005 |
| `DIBUAT_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0006; FINDING-002 §9 |
| `DIBUAT_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan — bukan tidak diketahui | FINDING-002 §9 |
| `DIBUAT_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0006 |
| `DIUBAH_OLEH` | `VARCHAR2(128 CHAR)` | — | Y | belum pernah diubah sejak dibuat | ADR-0006 |
| `DIUBAH_ATAS_NAMA` | `VARCHAR2(128 CHAR)` | — | Y | tidak ada perwakilan | FINDING-002 §9 |
| `DIUBAH_PADA` | `TIMESTAMP(6)` | — | Y | belum pernah diubah | ADR-0006 |

## 19. `MIGRASI_KORELASI`

Jembatan ke sistem lama. Dijatuhkan setelah paritas diterima (ADR-0005).

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_KORELASI` | `NUMBER(19)` | — | N | — | ADR-0021 |
| `TABEL_TUJUAN` | `VARCHAR2(30 CHAR)` | — | N | — | polimorfik — tanpa FK; panjangnya mengikuti batas nama objek, SPEC §16 |
| `ID_TUJUAN` | `NUMBER(19)` | — | N | — | ADR-0021 |
| `PZINSKEY_LAMA` | `VARCHAR2(255 CHAR)` | — | Y | **bukan berasal dari work object Pega** — bukan tidak diketahui | BLUEPRINT §8.4a; ddl `TABLE_PC_ASM_FW_GCNMFW_WORK` |
| `PYID_LAMA` | `VARCHAR2(32 CHAR)` | — | Y | tidak punya nomor klaim lama | ddl tabel work; **tanpa unique di sistem lama** (D13) |
| `CASEID_LAMA` | `VARCHAR2(255 CHAR)` | — | Y | tidak tercatat di sisi basis data lama | BLUEPRINT §13.2; salinan `pzInsKey` |
| `INDEX_OBJECT_LAMA` | `NUMBER(9)` | — | Y | bukan baris Adjustment | BLUEPRINT §8.4 — **posisi, bukan pengenal** (D17) |
| `URUTAN_BARIS_LAMA` | `NUMBER(9)` | — | Y | bukan baris di dalam PageList tertanam | BLUEPRINT §8.4a |
| `SUMBER_LAMA` | `VARCHAR2(64 CHAR)` | — | N | — | ADR-0023 |
| `DIMIGRASI_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0005 |

> Kelima kolom `*_LAMA` boleh kosong sendiri-sendiri; **kosong semua ditolak** (`CK_MIGRASI_KORELASI_1`). Istilah `_Avoid_` (`CASEID`, `pyID`) muncul di sini sebagai **nama lama bersufiks `_LAMA`** — pengecualian ADR-0021, sama dengan lapisan view kompatibilitas.

## 20. `MIGRASI_NILAI_DITOLAK`

Nilai lama yang tidak lolos aturan penguraian. Tidak dibulatkan, tidak dibuang (ADR-0003, ADR-0014, AK-4).

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_DITOLAK` | `NUMBER(19)` | — | N | — | ADR-0021 |
| `ID_PENDARATAN` | `NUMBER(19)` | — | Y | **bukan berasal dari muatan** — bukan tidak diketahui (ADR-0019). Nilai ini datang langsung dari kolom tabel lama | SPEC §21.8 — seam tiket `13`↔`14`; `FK_MIGRASI_NILAI_DITOLAK_1` |
| `SUMBER_LAMA` | `VARCHAR2(64 CHAR)` | — | N | — | ADR-0014 |
| `PENGENAL_LAMA` | `VARCHAR2(500 CHAR)` | — | N | — | ADR-0014; `pzInsKey` sampai 255 byte, ditambah penunjuk baris |
| `KOLOM_LAMA` | `VARCHAR2(64 CHAR)` | — | N | — | ADR-0014 |
| `NILAI_MENTAH` | `VARCHAR2(4000 CHAR)` | — | N | — | ADR-0003 — **teks, bukan angka**: mengubahnya jadi angka adalah persis hal yang gagal dilakukan. D39 (uang melintas sebagai `String`) |
| `SEBAB_DITOLAK` | `VARCHAR2(255 CHAR)` | — | N | — | AK-4 — **tanpa CHECK domain**: sebab yang belum pernah terjadi tidak dapat didaftar lebih dulu |
| `DITEMUKAN_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0014 |

## 21. `MIGRASI_PENDARATAN`

Satu-satunya tempat bentuk lama boleh masuk utuh, termasuk JSON. Dijatuhkan setelah paritas diterima.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_PENDARATAN` | `NUMBER(19)` | — | N | — | ADR-0021 |
| `SUMBER_LAMA` | `VARCHAR2(64 CHAR)` | — | N | — | ADR-0023 |
| `PENGENAL_LAMA` | `VARCHAR2(500 CHAR)` | — | N | — | ADR-0023 |
| `MUATAN` | `CLOB` | — | N | — | BLUEPRINT.md §13.2, §13.5 — **tanpa `IS JSON`**: bentuknya tidak dibatasi di sistem lama (D43), dan sistem lama memang menghasilkan JSON rusak (D41) |
| `BENTUK_TERURAI` | `NUMBER(1)` | — | N | — | SPEC §21.8 — hasil penguraian dicatat **sebagai fakta**, bukan dipakai sebagai syarat masuk |
| `SEBAB_GAGAL_URAI` | `VARCHAR2(255 CHAR)` | — | Y | **terurai** — `CK_MIGRASI_PENDARATAN_2` menolak "gagal tanpa sebab" dan "berhasil tapi bersebab" | SPEC §21.8 |
| `DITERIMA_PADA` | `TIMESTAMP(6)` | — | N | — | ADR-0014 |

## 22. `ARSIP_MUATAN_KELUAR`

Apa yang BENAR-BENAR dikirim ke hilir, bukan apa yang seharusnya. Tulis-sekali, tidak pernah diubah.

| Kolom | Tipe | Kel. | Null | Arti NULL | Sumber |
|---|---|---|---|---|---|
| `ID_ARSIP` | `NUMBER(19)` | — | N | — | AK-1 — "arsipnya tabel tersendiri, bukan kolom di tabel bisnis" |
| `TUJUAN` | `VARCHAR2(24 CHAR)` | — | N | — | S1: SaveDataToOSAksep_Act (akseptasi) / HitServiceToKasir_Act (kasir) |
| `ID_KLAIM` | `NUMBER(19)` | — | N | — | ADR-0024 |
| `JENIS_REASURANSI` | `VARCHAR2(100 CHAR)` | — | N | — | S1: InputParamOs.TypeLoss <- .TreatyName; RDBList/GetDataOS.xml penyaring a.data_json.TypeLoss |
| `MATA_UANG` | `VARCHAR2(3 CHAR)` | — | N | — | RDBList/GetDataOS.xml penyaring a.data_json.Currency |
| `DITOLAK` | `NUMBER(1)` | — | N | — | RDBList/GetDataOS.xml — penyaringnya memuat STS_REJECT = 0; muatan ditolak TIDAK ikut dijumlahkan |
| `NILAI_TERSEBAR` | `NUMBER(38,2)` | tepi | N | — | S1: InputParamOs.Value — dibulatkan 2 di tepi (AK-1) |
| `NILAI_BRUTO` | `NUMBER(38,2)` | tepi | N | — | S1: InputParamOs.GrossValue |
| `BIAYA_PENILAIAN` | `NUMBER(38,2)` | tepi | N | — | S1: InputParamOs.Adjusterfee |
| `SALVAGE` | `NUMBER(38,2)` | tepi | N | — | S1: InputParamOs.Salvage |
| `BIAYA_LAIN` | `NUMBER(38,2)` | tepi | N | — | S1: InputParamOs.CNPOthersFee |
| `KURS` | `NUMBER(38,20)` | K1 | Y | muatan tidak membawa kurs | S1: InputParamOs.KursValue |
| `DIKIRIM_OLEH` | `VARCHAR2(128 CHAR)` | — | N | — | ADR-0017 — lewat akun aplikasi |
| `DIKIRIM_PADA` | `TIMESTAMP(6)` | — | N | — | catatan, bukan turunan |

---

# VIEW

Kolom view adalah **kontrak ke Arasapas dan kasir** — lebih mengikat daripada kolom tabel kanonik, karena pihak di seberangnya tidak dapat diubah. Daftar kolomnya **dibaca dari XML** lewat sapuan S1, bukan dirancang.

## V1. `V_AKSEPTASI_KOMPATIBEL`

Kontrak keluar ke Arasapas dan kasir. Bentuk lama dipertahankan, termasuk CASEID (pengecualian ADR-0021).

| Kolom | Isi | Sumber |
|---|---|---|
| `CASEID` | pengenal lama; hanya hidup di lapisan view | MIGRASI_KORELASI.CASEID_LAMA; ADR-0021 |
| `NOCLAIM` | Nomor Klaim | S1: InputParamOs.NoClaim |
| `TYPELOSS` | jenis reasuransi — satuan kasar | S1: InputParamOs.TypeLoss <- .TreatyName |
| `TYPELOSSID` | id jenis reasuransi | S1: InputParamOs.TypeLossID <- .TreatyType |
| `CURRENCY` | mata uang | S1: InputParamOs.Currency |
| `KURSVALUE` | kurs | S1: InputParamOs.KursValue <- .Kurs |
| `VALUE` | nilai tersebar, dijumlahkan, dibulatkan 2 | S1: InputParamOs.Value <- .ClaimSpreaded |
| `GROSSVALUE` | nilai bruto, dijumlahkan, dibulatkan 2 | S1: InputParamOs.GrossValue |
| `ADJUSTERFEE` | biaya penilaian, dijumlahkan | S1: InputParamOs.Adjusterfee |
| `SALVAGE` | salvage, dijumlahkan | S1: InputParamOs.Salvage |
| `CNPOTHERSFEE` | biaya lain, dijumlahkan | S1: InputParamOs.CNPOthersFee |
| `PERSENRNM` | porsi reasuradur — tidak dijumlahkan | S1: InputParamOs.PersenRNM |
| `IDMASTERTREATY` | pengenal treaty | S1: InputParamOs.IDMasterTreaty |
| `CAUSEOFLOSS` | penyebab kerugian | S1: InputParamOs.CauseOfLoss |
| `CAUSEOFLOSSID` | id penyebab kerugian | S1: InputParamOs.CauseOfLossID |

## V2. `V_AKSEPTASI_DITOLAK`

Kelompok yang ditolak V_AKSEPTASI_KOMPATIBEL karena besaran per-layer tidak sepakat. Tidak ada yang hilang diam-diam.

| Kolom | Isi | Sumber |
|---|---|---|
| `ID_KLAIM` | klaim | SPEC bagian 0 |
| `JENIS_REASURANSI` | satuan kasar | SPEC bagian 0 |
| `MATA_UANG` | mata uang | SPEC bagian 0 |
| `CACAH_BARIS` | berapa baris halus dilebur | SPEC bagian 0 |
| `SEBAB` | besaran mana yang tidak sepakat | SPEC bagian 0 |

## V3. `V_REKAP_KLAIM_MATA_UANG`

Rekap per klaim per mata uang. Retensi Cedant sebagai kolom TERPISAH, tidak pernah tercampur.

| Kolom | Isi | Sumber |
|---|---|---|
| `ID_KLAIM` | klaim | ADR-0010 |
| `MATA_UANG` | mata uang | ADR-0007 |
| `NILAI_TERALOKASI_LAYER` | jumlah alokasi seluruh layer | ADR-0010 |
| `NILAI_RETENSI_CEDANT` | Retensi Cedant — kolom terpisah | ADR-0010 |
| `PREMI_PEMULIHAN` | jumlah premi pemulihan | MEMORI_PEMAHAMAN.MD §6.3 |

## V4. `V_PARITAS_SHADOW`

Perbandingan baris per baris terhadap sistem lama. ADR-0005; toleransi di AK-3.

| Kolom | Isi | Sumber |
|---|---|---|
| `TABEL_TUJUAN` | tabel kanonik | SPEC bagian 8 |
| `ID_TUJUAN` | baris kanonik | SPEC bagian 8 |
| `PZINSKEY_LAMA` | pengenal lama | BLUEPRINT.md §8.4b |
| `PYID_LAMA` | pengenal lama | ddl TABLE_PC_ASM_FW_GCNMFW_WORK |
| `CASEID_LAMA` | pengenal lama | BLUEPRINT.md §13.2 |
| `INDEX_OBJECT_LAMA` | posisi numerik lama | BLUEPRINT.md §8.4 |
| `KEADAAN_BARIS` | baris belum lengkap dikeluarkan dari perbandingan | ADR-0019 |

## V5. `V_TOTAL_ESTIMASI`

Pengganti PageList ListTotalEstimation. Agregasi murni.

| Kolom | Isi | Sumber |
|---|---|---|
| `ID_KLAIM` | klaim | BLUEPRINT.md §2.1 |
| `MATA_UANG` | mata uang | BLUEPRINT.md §2.1 |
| `NILAI_KERUGIAN` | jumlah | MEMORI_PEMAHAMAN.MD §4.2 |
| `BIAYA_PENILAIAN` | jumlah | MEMORI_PEMAHAMAN.MD §4.2 |
| `SALVAGE` | jumlah | MEMORI_PEMAHAMAN.MD §4.2 |
| `BIAYA_LAIN` | jumlah | MEMORI_PEMAHAMAN.MD §4.2 |

## V6. `V_PENYEBARAN_AGREGAT`

Pengganti PageList SpreadingAdjustment. Agregasi per jenis reasuransi per mata uang.

| Kolom | Isi | Sumber |
|---|---|---|
| `ID_KLAIM` | klaim | MEMORI_PEMAHAMAN.MD §6.4 |
| `JENIS_REASURANSI` | kunci agregasi lama (Currency, TreatyName) | MEMORI_PEMAHAMAN.MD §6.4 |
| `MATA_UANG` | idem | MEMORI_PEMAHAMAN.MD §6.4 |
| `NILAI_TERSEBAR` | jumlah | MEMORI_PEMAHAMAN.MD §6.4 |
| `BIAYA_PENILAIAN` | jumlah | MEMORI_PEMAHAMAN.MD §6.4 |
| `SALVAGE` | jumlah | MEMORI_PEMAHAMAN.MD §6.4 |
| `BIAYA_LAIN` | jumlah | MEMORI_PEMAHAMAN.MD §6.4 |

## V7. `V_PENYEBARAN_AGREGAT_QS`

Pengganti PageList SpreadingAdjustmentQS. Sama, disaring ke jenis quota share.

| Kolom | Isi | Sumber |
|---|---|---|
| `ID_KLAIM` | klaim | MEMORI_PEMAHAMAN.MD §6.4 |
| `JENIS_REASURANSI` | jenis quota share | CountLossAllocation_act Langkah 10 |
| `MATA_UANG` | mata uang | MEMORI_PEMAHAMAN.MD §6.4 |
| `NILAI_TERSEBAR` | jumlah | MEMORI_PEMAHAMAN.MD §6.4 |

## V8. `V_TOTAL_NILAI_PERTANGGUNGAN`

Pengganti PageList TotalInterestInsured.

| Kolom | Isi | Sumber |
|---|---|---|
| `ID_KLAIM` | klaim | BLUEPRINT.md §2.1 |
| `MATA_UANG` | mata uang | MEMORI_PEMAHAMAN.MD §4.2 |
| `NILAI_PERTANGGUNGAN` | jumlah | MEMORI_PEMAHAMAN.MD §4.2 |

## V9. `V_REKAP_AKSEPTASI`

Pengganti PageList ListClaimAcceptation.

| Kolom | Isi | Sumber |
|---|---|---|
| `ID_KLAIM` | klaim | BLUEPRINT.md §2.1 |
| `MATA_UANG` | mata uang | ADR-0024 |
| `CACAH_AKSEPTASI` | berapa baris akseptasi | ADR-0024 |
| `NILAI_DISETUJUI` | jumlah | ADR-0015 |


---

## Kolom yang bertambah 19 September 2026

Lima kolom tabel dan satu kolom view, seluruhnya lahir dari tiket yang ditutup hari ini.

| Tabel / view | Kolom | Tipe | Null | Arti NULL | Dari |
|---|---|---|---|---|---|
| `KLAIM` | `LINI_USAHA` | `VARCHAR2(40 CHAR)` | Y | belum ditetapkan — bukan "konvensional" | tiket `12`, SPEC §21.3 |
| `KLAIM` | `PENUTUPAN_LAMA` | `VARCHAR2(64 CHAR)` | Y | tidak ditutup di sistem lama | tiket `12`, SPEC §21.6 |
| `MIGRASI_PENDARATAN` | `BENTUK_TERURAI` | `NUMBER(1)` | N | — | tiket `13`, SPEC §21.8 |
| `MIGRASI_PENDARATAN` | `SEBAB_GAGAL_URAI` | `VARCHAR2(255 CHAR)` | Y | terurai | tiket `13` |
| `MIGRASI_NILAI_DITOLAK` | `ID_PENDARATAN` | `NUMBER(19)` | Y | bukan berasal dari muatan | tiket `14`, seam `13`↔`14` |
| enam view agregat | `CACAH_BELUM_LENGKAP` | `NUMBER` | N | — | tiket `30`, `31` |

Dan satu kolom **dilonggarkan**: `ADJUSTMENT.KEPUTUSAN_KOMITE` dari `NOT NULL` menjadi nullable — tanpa itu, keadaan keempat tidak dapat disimpan sama sekali (tiket `18`).

Satu kolom **diganti nama** bersama tabelnya: `ID_DAFTAR` → `ID_PENJAGA_TANGGAL` pada `KLAIM_PENJAGA_TANGGAL`.
