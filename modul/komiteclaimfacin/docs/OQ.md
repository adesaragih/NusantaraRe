# OQ — pertanyaan terbuka modul Komite Claim Fac In

> 10-10-2026, tahap 2 dibangun (prompt "IMPLEMENTASI KOMITE CLAIM FAC IN, TAHAP 2 DARI 2"). Keputusan work owner
> KCF-01–04 di `MODUL.md` bab "Keputusan"; OQ lain terbit saat membangun. Kolom **Rekomendasi** = yang dibangun sekarang.
> Status setiap langkah: [`PARITAS.md`](PARITAS.md). Pengelompokan per pemilik jawaban (pola `modul/claimfacin/docs/OQ.md`).
> Akun orang, alamat surel, dan alamat yang tertulis mati di XML **tidak dicetak** di dokumen ini.

## Diputuskan

| Keputusan | Isi | Tanggal |
| --- | --- | --- |
| **KCF-01** | Roster `EMAILKOMITE` FACIN → workbasket (migrasi 640); anggota `ReasClaimSPVB` juga memutus tingkat 1; pita SPV B 30.000.000 < total ≤ 57.750.000 → satu jenjang atas | 10-10-2026 |
| **KCF-02** | Perluasan tangga ikut XML (`ApprovalKomite_Act` L4–L8), ditampilkan saat tingkat 1 membuka, disimpan saat tingkat 1 Submit; Fac Retro melewatinya | 10-10-2026 |
| **KCF-03** | TT3 / TT4 untuk Fac In: `ADJUSTMENT_ID` boleh kosong (641) + `TRANSFER_TYPE` (642); satu tingkat `ReasClaimDeptHead` | 10-10-2026 |
| **KCF-04** | Kasus komite lama Pega `KMT-` tidak dimigrasi | 10-10-2026 |
| **OQ-KCFI-02** | Kasir kasus Fac Retro **ikut XML**: tetap dikirim bila gerbangnya terpenuhi (DirectToKasir, StatusKasir kosong, status konversi 1); tanpa pengecualian retro (retro biasanya baru terkonversi sesudah Print DLA di klaim) | 10-10-2026 |
| **OQ-KCFI-03** | **Boleh**: 3 kolom nullable `T_GENERAL_KOMITE` (`KOMITE_CIRCUM_CAUSE_OF_LOSS` / `KOMITE_EXTENT_OF_LOSS` / `KOMITE_LEGAL_LIABILITY`, ADD saja, migrasi 643); isian pop-up `SendRejectClaimToKomite2` / `SendCloseClaimToKomite` disimpan dan tampil di LS39; STRUKTUR pemilik tabel diperbarui | 10-10-2026 |
| **OQ-KCFI-06** | **Anggap SPV A**: pemutus tingkat 1 anggota `ReasClaimSPVA` + `ReasClaimSPVB` sekaligus -> pita SPV B tidak berlaku; pita hanya untuk anggota `ReasClaimSPVB` saja (rangkap di produksi dijaga Kelola User) | 10-10-2026 |
| **OQ-KCFI-08** | **Boleh, sempit**: `STS_KONVERSI` kosong -> kolom tidak ikut INSERT OS (preseden Claim Prop); pengecualian penjaga Claim Life `repository/migrasibatas_test.go` hanya untuk `backend/repository/osakseptasi.go` | 10-10-2026 |

## A. Work owner

| OQ | Pertanyaan | Rekomendasi (yang dibangun) | Bukti |
| --- | --- | --- | --- |
| **OQ-KCFI-07** | ADR-0038 / 0039 / 0040 (`OUTPUT_HASIL_RNM/docs/bersama/adr/`, di luar repo) bersumber dari spec September dan kini bertentangan dengan RALAT 10-10-2026 (ADR-0040 premis gugur; ADR-0038 "bukan karena pencacah" padahal XML memakai `KomiteCount`; ADR-0039 "data kutipan disalin utuh") | Tidak disunting (di luar berkas yang diizinkan). Usul: ADR diperbarui atau diberi catatan "digantikan KCF-01–04" | `docs/spec.md` Lampiran 3 (RALAT) |

## B. Pemilik ekspor Pega

| OQ | Rule yang dibutuhkan | Akibat sekarang |
| --- | --- | --- |
| **OQ-KCFI-01** | Stream HTML `AcceptanceNotePDFMBU` / `AcceptanceNotePDFTravel` / `AcceptanceNotePDFPA` (`PrintPDFAccep_MultiAksep_KMT` S19-S21), `CommitteReject_CC`, `CommitteCloseClaim`, badan email `EmailKlaim_HTML_KMT` | **Sebagian terjawab 10-10-2026**: work owner memberi `Claim Fac In/AcceptanceNotePDF.xml` -> PDF akseptasi lini Fire / Aneka / Golf / Marine Cargo digambar go-pdf/fpdf (pola Komite Claim Prop), diunggah sesudah Submit tersimpan, dicatat T_STORAGE_IMAGE + DOCUMENT_CLAIM (PARITAS S9). Stream lain masih tidak diekspor: Submit lini MBU / Travel / PA dan TT3 / TT4 menjawab `info` yang menyebut stream-nya (`services.InfoTanpaStream`); nama berkas VERBATIM tetap disiapkan; badan email tidak dikarang |

## C. DBA

| OQ | Pertanyaan | Akibat sekarang |
| --- | --- | --- |
| — | `-migrate` 640–642 **sudah dijalankan** work owner 10-10-2026 15:02 (diverifikasi baca-saja: ADJUSTMENT_ID nullable, TRANSFER_TYPE + CHECK, roster FACIN workbasket). **643 menunggu** `-migrate`: sebelum itu kueri kepala kasus komite (tiga kolom teks) → `ORA-00904` di DEV | — |
| — | Mundur 641 = `MODIFY (ADJUSTMENT_ID NOT NULL NOVALIDATE)`: baris TT3 / TT4 ber-`ADJUSTMENT_ID` kosong yang sudah ada dibiarkan (kendala tidak memvalidasi baris lama) | — |

## Catatan (tanpa keputusan)

| Butir | Catatan |
| --- | --- |
| **OQ-KCFI-04** | **Ditutup 10-10-2026 tanpa keputusan** — gerbang kasir `StatusKonversi = 1` ADA di XML (`HitServiceToKasirKMT_Act` / `HitServiceToKasir_Act` langkah 3, transisi pasca-langkah yang tidak tercetak alat dump awal); dibangun di kedua modul (prompt §2 butir 6 benar) |
| **OQ-KCFI-05** | `CurrencyDol` (kurs tukar adjustment) tidak disimpan Claim Fac In; "Total in IDR" memakai `POOLDATA.GETCURRENCYSTANDARD` — sumber yang sama dengan `SetNilaiResikoSendiri` Claim Fac In tahap 1 |
| K7 | Asal "K7 · 2026-09-20" tidak tertelusur (`keputusan-sebelum-to-spec.md` menulis Q3 belum diputuskan); isinya kini dipastikan XML + KCF-02 — tidak menahan. Nomor K7 juga dipakai STRUKTUR Claim Fac In §5d untuk keputusan lain |
| MBU / Travel | Di Pega, `IsMBU` sisi komite menguji `pyWorkPage.Quotation.BusinessCode` yang tidak disalin `SetValueKomite` S12; sistem baru membaca halaman klaim utuh lewat kontrak, jadi cacat itu tidak terbawa |
| OQ-CFI-26 / 28 | Dedupe kasir saat muatan masih antre dan satu adjustment per KMT — tetap terbuka di Claim Fac In |
