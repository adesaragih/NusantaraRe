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

## A. Work owner

| OQ | Pertanyaan | Rekomendasi (yang dibangun) | Bukti |
| --- | --- | --- | --- |
| **OQ-KCFI-02** | `KomitePost_Adjustment` S25 (kasir) bergerbang `Local.FacRetro == 1 [T=3]` berkomentar "JIKA RETRO KONVESI WAKTU PRINT DLA", tetapi `pre=false` — gerbangnya tidak dinilai, kasir **tetap** berjalan untuk Fac Retro. Maksudnya retro dilewati? | Ikut XML: muatan kasir retro ikut diantre (produksi saja); konversi S17 dan log S19 tetap dilewati untuk retro | `Activity/KomitePost_Adjustment.xml` S17, S19, S25 |
| **OQ-KCFI-03** | Pop-up Reject Claim (ClaimComiteeReject LS4) dan Close Claim (PreventRejectClaim LS23) punya isian **Chronology / Extent Of Loss / Policy Liability** (`TempCommiteClaim.*`); `SendRejectClaimToKomite2` / `SendCloseClaimToKomite` menyalinnya ke `childPageKomite.Komite.*` dan `ClaimData.ClaimComitee(<LAST>)`, lalu ShowTransfer LS39 menampilkannya. Kasus TT3 / TT4 **tanpa adjustment**, jadi tidak ada kolom penampung (TT2 menyimpannya di baris adjustment `DataCommitteFacin.*`). Simpan di 3 kolom baru `T_GENERAL_KOMITE` (ADD, di luar KCF-03)? | Tidak disimpan: isian tetap dapat diketik (tampil di pop-up sampai ditutup), layar komite TT3 / TT4 hanya menampilkan **Remarks** (`ClaimData.Remark` klaim induk, langkah 2). Nol DDL tanpa izin | `SendRejectClaimToKomite2.xml` / `SendCloseClaimToKomite.xml` (salinan `TempCommiteClaim.*`), Section `ShowTransfer` LS39 |
| **OQ-KCFI-06** | Akun yang memegang **`ReasClaimSPVA` dan `ReasClaimSPVB` sekaligus**: XML memakai jabatan profil (`pyPosition == "SPV B"`, satu nilai). Sekarang pita SPV B berlaku bila pemutus anggota `ReasClaimSPVB`, walau ia juga anggota SPVA | Keanggotaan SPVB menang (pita berlaku). Usul: satu akun hanya di salah satu dari dua workbasket | `ApprovalKomite_Act` L5, `models.AnggotaSPVB` |
| **OQ-KCFI-07** | ADR-0038 / 0039 / 0040 (`OUTPUT_HASIL_RNM/docs/bersama/adr/`, di luar repo) bersumber dari spec September dan kini bertentangan dengan RALAT 10-10-2026 (ADR-0040 premis gugur; ADR-0038 "bukan karena pencacah" padahal XML memakai `KomiteCount`; ADR-0039 "data kutipan disalin utuh") | Tidak disunting (di luar berkas yang diizinkan). Usul: ADR diperbarui atau diberi catatan "digantikan KCF-01–04" | `docs/spec.md` Lampiran 3 (RALAT) |
| **OQ-KCFI-08** | Penjaga Claim Life `repository/migrasibatas_test.go` `TestKolomTakDibawaHanyaAdaDiKatalog` (sudah gagal di HEAD) kini juga menyebut `backend/repository/warisan.go` modul ini: `STS_KONVERSI` termasuk "kolom yang sengaja tidak dibawa". Kolom itu kolom procedure `PEGA_JSON_OS_AKSEP_KLAIM` yang ditulis `SaveReject_ACT_KMT` S17 (TT3, `1` bila Value kosong); preseden HEAD: `claimnonprop/backend/repository/tulisan.go` | Kolom tetap ditulis (ikut XML). Butuh pengecualian di berkas penjaga Claim Life (modul lain — tidak disunting) | `SaveReject_ACT_KMT.xml` S17; penjaga Claim Life |

## B. Pemilik ekspor Pega

| OQ | Rule yang dibutuhkan | Akibat sekarang |
| --- | --- | --- |
| **OQ-KCFI-01** | Stream HTML `AcceptanceNotePDF*` (`PrintPDFAccep_MultiAksep_KMT`), `CommitteReject_CC`, `CommitteCloseClaim`, badan email `EmailKlaim_HTML_KMT` | Nomor, penanda (`IsPrintAccept`), nama berkas VERBATIM disiapkan; **berkas PDF tidak dibuat** dan badan email tidak dikarang. Submit yang di Pega menerbitkan PDF menjawab `info` OQ-KCFI-01. Bila stream diberikan: digambar go-pdf/fpdf (pola Komite Claim Prop) |

## C. DBA

| OQ | Pertanyaan | Akibat sekarang |
| --- | --- | --- |
| — | `-migrate` 640–642 menunggu work owner. Sebelum 642: setiap kueri yang menyebut `T_GENERAL_KOMITE.TRANSFER_TYPE` → `ORA-00904` di DEV (wajar; dicoba baca-saja 10-10-2026: varian tanpa kolom itu lolos urai / SELECT) | — |
| — | Mundur 641 = `MODIFY (ADJUSTMENT_ID NOT NULL NOVALIDATE)`: baris TT3 / TT4 ber-`ADJUSTMENT_ID` kosong yang sudah ada dibiarkan (kendala tidak memvalidasi baris lama) | — |

## Catatan (tanpa keputusan)

| Butir | Catatan |
| --- | --- |
| **OQ-KCFI-04** | **Ditutup 10-10-2026 tanpa keputusan** — gerbang kasir `StatusKonversi = 1` ADA di XML (`HitServiceToKasirKMT_Act` / `HitServiceToKasir_Act` langkah 3, transisi pasca-langkah yang tidak tercetak alat dump awal); dibangun di kedua modul (prompt §2 butir 6 benar) |
| **OQ-KCFI-05** | `CurrencyDol` (kurs tukar adjustment) tidak disimpan Claim Fac In; "Total in IDR" memakai `POOLDATA.GETCURRENCYSTANDARD` — sumber yang sama dengan `SetNilaiResikoSendiri` Claim Fac In tahap 1 |
| K7 | Asal "K7 · 2026-09-20" tidak tertelusur (`keputusan-sebelum-to-spec.md` menulis Q3 belum diputuskan); isinya kini dipastikan XML + KCF-02 — tidak menahan. Nomor K7 juga dipakai STRUKTUR Claim Fac In §5d untuk keputusan lain |
| MBU / Travel | Di Pega, `IsMBU` sisi komite menguji `pyWorkPage.Quotation.BusinessCode` yang tidak disalin `SetValueKomite` S12; sistem baru membaca halaman klaim utuh lewat kontrak, jadi cacat itu tidak terbawa |
| OQ-CFI-26 / 28 | Dedupe kasir saat muatan masih antre dan satu adjustment per KMT — tetap terbuka di Claim Fac In |
