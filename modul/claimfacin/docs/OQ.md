# OQ — pertanyaan terbuka modul Claim Fac In

> 10-10-2026, tahap 1 dibangun (prompt "IMPLEMENTASI CLAIM FAC IN, TAHAP 1 DARI 2"). OQ-CFI-01–04 diputuskan work owner
> 09-10-2026 (`MODUL.md` bab "Keputusan"); bawaan a–g prompt §3 dipakai sampai work owner mengubahnya. OQ lain terbit
> saat membangun; kolom **Rekomendasi** = yang dibangun sekarang. Status setiap tombol: [`PARITAS.md`](PARITAS.md).
>
> Pengelompokan per pemilik jawaban (pola `modul/claimprop/docs/OQ.md`):
> - **work owner:** keputusan.
> - **pemilik ekspor Pega:** rule yang tidak ada di korpus.
> - **DBA:** objek basis data dan data warisan.
>
> Nomor OQ-CFI-05 s.d. 10 tidak terbit (tidak dipakai). Akun orang, alamat surel, kata sandi, dan alamat yang tertulis
> mati di XML **tidak dicetak** di dokumen ini.

## Diputuskan

| OQ | Keputusan | Tanggal |
| --- | --- | --- |
| **OQ-CFI-01** | Tabel baru `T_CLAIM_OBJECT` / `T_CLAIM_OBJECT_ITEM`; tabel `T_CLAIM_*` Claim Prop dipakai ulang dengan `OBJECT_ITEM_ID` (ALTER ADD); baris FAC mengisi `CLAIM_ID` dan `OBJECT_ITEM_ID`; nol MODIFY | 09-10-2026 |
| **OQ-CFI-02** | Awalan `CLM-` / `KMT-` sama dengan Pega, nomor `SEQ_WORK_CLAIM`; lini disaring lewat `T_WORK_CLAIM.LINI`, tidak pernah lewat awalan | 09-10-2026 |
| **OQ-CFI-03** | Kelainan XML diperbaiki semua (`PARITAS.md` §7–§8); maksud tidak pasti → OQ | 09-10-2026 |
| **OQ-CFI-04** | Komite pola Claim Prop, tanpa menu; tahap 1 hanya melahirkan kasus `KMT-` TT2 | 09-10-2026 |
| **OQ-CFI-27** | **KCF-03**: TT3 Reject / TT4 Close Without Payment diizinkan untuk Fac In — migrasi `komiteclaimfacin` 641 (`T_GENERAL_KOMITE.ADJUSTMENT_ID` boleh kosong) + 642 (`TRANSFER_TYPE` 2 / 3 / 4); komite satu tingkat `ReasClaimDeptHead`. Tombol **Yes** kedua pop-up dinyalakan (tahap 2, `PARITAS.md` §5) | 10-10-2026 |

## A. Work owner

| OQ | Pertanyaan | Rekomendasi (yang dibangun) | Bukti |
| --- | --- | --- | --- |
| **OQ-CFI-12** | `CheckDate_Act` 9 / 13 `Param.StartDate := Policy.StartDateTime+1`: tanggal polis sesudah `CopyNB_Act` 11 berbentuk teks `yyyyMMdd`, `+1` menjadi penggabungan teks, `@CompareDates` 11 selalu salah — pesan "Date of Loss should not be more than Start Date Policy" tidak pernah tampil. Maksudnya +1 hari, tanpa +1, atau memang dimatikan? | Ikut hasil XML (pesan awal polis tidak dipasang); langkah 15 (akhir polis) dibangun | `Activity/CheckDate_Act.xml` 9, 11, 13, 15 |
| **OQ-CFI-14** | Layanan REST luar `getPremiumPaidOn` / `getPremiumPaidOnMarine` (View Status Payment Premi), `getPaymentClaim` / `getPayAttachment` (View Status Payment Claim / View Payment Attachment) belum disetujui untuk dipanggil | Tombol tampil nonaktif dengan `title` OQ. Jalur `getPremiumPaidOnMarine` khusus satu polis (bawaan c) dibuang | `ConnectREST/*`, Section `InputRegister`, `Adjusment_SC` |
| **OQ-CFI-17** | **Save Spreading** (kata sandi tertulis mati di `SaveSpreadingSP_act`) ada di container `LS42` section `Estimasi` bergerbang `ClaimData.ExGratia = 1`; penulis satu-satunya `InsertObjects_dt` 11 menulis `0`. Tombol tidak pernah tampil | Tidak dibangun; kata sandi tidak ditulis di mana pun. Bila fitur ini memang dipakai di Pega hidup (screenshot), penulis `ExGratia` tingkat klaim perlu diekspor | `Section/Estimasi.xml` LS42, `DataTransform/InsertObjects_dt.xml` 11 |
| **OQ-CFI-18** | Prompt values `associated` tidak diekspor: label Payment Type 1–7 (6 / 7 tidak terbaca), Acceptance Status, Report Type, Reporter Status, Share Retro, Tipe Analisis | Kode DB tampil apa adanya; label ReportType / IndividualRiskType / Payable dari screenshot Claim Prop (kelas properti sama, `[inferensi]`). Mohon screenshot | Section `InputAdjustment` LS15, `InputRegisterDetail` |
| **OQ-CFI-19** | Kronologi (`ClaimData.Chronology`) disimpan di `T_VIEW_SUGGEST` (tabel bersama): `ASMDateTimeChronology` → `DATE_SUGGEST`, `ASMUser` → `PIC_SUGGEST`, `ASMUserID` (jabatan, kolom "User") → `IS_CEDING_CONFIRM`, `pyNote` → `COMMENT_SUGGEST`, `ASMNoteType` (kolom "Status") → `INITIAL_SUGGEST`. Pemetaan ke kolom bersama yang ada (tanpa kolom baru) disetujui? | Seperti di kiri (pola Claim Prop) | `repository/halaman.go` `sqlSisipKronologi` |
| **OQ-CFI-23** | Tombol **+** grid Adjustment tidak ber-`addRow`; `CountTotalEstimasi_Act` menulis `.Adjustment(<LAST>)`. Di Pega hidup, apakah + menambah baris atau mengisi ulang baris terakhir? | `[inferensi]` baris baru; pemeriksaan 20–21 atas baris terakhir yang sudah ada | `Section/Adjusment_SC.xml`, `Activity/CountTotalEstimasi_Act.xml` 20–21 |
| **OQ-CFI-24** | `CekLimit.CARI1` (halaman requestor) **tidak pernah ditulis** rule mana pun, sehingga PLA treaty **G** (`GeneratePLATreaty_Act`) dan DLA treaty **S** (`DLAFacintoTreaty_Act`) tidak pernah jalan; `CheckLimitSpreadingTreaty_Act` membaca baris objek, dan prioritas operator langkah 20–21 (`A && B \|\| C`) tidak jelas | Ditiru apa adanya (jalur treaty dibangun, gerbangnya tidak tercapai). Perlu penulis `CekLimit.CARI1` dari ekspor / screenshot | `Activity/GeneratePLA.xml`, `ChooseDla_Act.xml` 3.1.4, `CheckLimitSpreadingTreaty_Act.xml` 20–21 |
| **OQ-CFI-26** | Panggilan kasir nyata (`SendAcceptationToKasir`) belum disetujui; `DIRECTTOKASIR_LOG` ditulis bersama tanggapan kasir. Outbox tanpa dedupe: `StatusKasir` baru terisi dari `DIRECTTOKASIR_LOG` sesudah pelaksana menjalankan muatan, sehingga klik Acceptation ulang ketika kirim pertama masih antre mengantre muatan kedua ber-`AcceptedNo` sama (Pega memanggil kasir sinkron, `StatusKasir` langsung terisi). Dedupe "masih antre" butuh pembacaan `T_LOG_SERVICE_RNM` milik inti (kontrak baru) | Muatan kasir masuk outbox hanya di produksi; modul ini belum punya pelaksana outbox (`modul.go` `TanpaPekerja`), jadi panggilan nyata dan `DIRECTTOKASIR_LOG` belum terjadi. Dedupe menunggu keputusan | `Activity/HitServiceToKasir_Act.xml` 13.2–13.6 |
| **OQ-CFI-28** | `CreateKMTNo_Act` 6.8 menyertakan adjustment PT 2 lain yang belum berkomite ke KMT yang sama, tetapi indeks unik `UQ_CLAIM_ADJUSTMENT_KOMITE` menolak dua adjustment bertaut ke satu KMT | Satu adjustment per KMT; bundel PT 2 tidak dibangun (tanpa MODIFY indeks) | `Activity/CreateKMTNo_Act.xml` 6.8 |
| **OQ-CFI-29** | Harness `New` / `NewSample` tidak diekspor; jalur B2B (`When/IsSPK`, `OfferFacIn.IsB2B = "SPK"` → langsung Choose Surveyor) hanya diputuskan di Start, saat polis belum dipilih. Dari mana kasus B2B lahir? | Add Claim = Start → Decision B2B salah → Input Register; jalur B2B tidak terjangkau dari Add Claim | `Flow/Register_Flow.xml` Decision4, `When/IsSPK.xml` |
| **OQ-CFI-31** | **Claim Life** `GET /api/klaim-life/{id}` membaca kepala `T_GENERAL_CLAIM` menurut ID **tanpa saringan LINI**; kasus FACIN berawalan `CLM-` yang sama, sehingga kepalanya terbaca bila ID diketik di URL Claim Life (hanya-baca; penulis per-ID Life menolak karena tahapnya tidak dikenal). Keadaan yang sama sudah berlaku untuk `CLMP-` / `CLMNP-` | Tidak disunting (modul lain). Usul: pembaca header Life menyaring `NVL(w.LINI, 'LIFE') = 'LIFE'` lewat join `T_WORK_CLAIM` | `PARITAS.md` §10; `claimlife/backend/repository/klaimlife.go` `AmbilHeader` |
| **OQ-CFI-32** | Ikon **+** Consultant / Adjuster memakai rute pinjaman `POST /api/adjuster-consultant`, yang di `cmd/api/rakit.go` (`ruteDipinjam`) hanya dipinjamkan ke `claimprop` (memanggilnya dari modul ini juga gagal penjaga `TestPanggilanLintasModulTerdaftar`) | Ikon tampil nonaktif dengan `title` OQ; memilih adjuster / consultant yang ada tetap jalan. Usul: tambah `claimfacin` ke `ruteDipinjam` (berkas inti, tidak disunting) | `cmd/api/rakit.go` `ruteDipinjam` |
| **OQ-CFI-33** | `CloseClaim` 7 `UpdateTotalJob_sql` (`mst_user_teknis.TOTAL_JOB`) membaca `ClaimData.UserTeknis` yang tanpa penulis (medan LS177 `InputRegisterDetail` VIS `1=2`) — langkah tidak pernah berjalan. Prompt §2 butir 6 menyebutnya | Tidak dibangun (ikut XML). Bila beban kerja teknis memang dihitung, sumber `UserTeknis` perlu ditetapkan | `Activity/CloseClaim.xml` 7, `Section/InputRegisterDetail.xml` LS177 |
| **OQ-CFI-34** | `CheckPeriodePolicy` langkah 1 `Exit-Activity` pre=false keluar tanpa syarat — cek periode polis dan endorsement (2–6) tidak pernah jalan. Sengaja? | Tidak dibangun (ikut XML) | `Activity/CheckPeriodePolicy.xml` 1–6 |

## B. Pemilik ekspor Pega

| OQ | Rule yang dibutuhkan | Akibat sekarang |
| --- | --- | --- |
| **OQ-CFI-13** | Validate FlowAction `ValidateDate` (InputRegister) | Tidak dibangun; `CheckDate*` dan `ProteksiDataRegister_Act` tetap memeriksa tanggal |
| **OQ-CFI-15** | FlowAction `InputSubProgressClaim` tanpa post-activity; sumber Progres1 / Progres2 (`AddDelProgress` hanya mengubah clipboard) | Tombol "Input Progres Claim" nonaktif |
| **OQ-CFI-20** | Aliran HTML dokumen (`Property-Set-HTML` + `HTMLToPDF`): Claim Face Sheet per lini, PLA, Draft DLA, DLA fac / treaty, nota akseptasi `PrintPDFAccep_MultiAksep` | Nomor, data, dan penanda (`IsCFS`, `PrintFaceClaim`, `PlaStatus`, `DLAStatus`) ditulis seperti XML; berkas PDF belum dibuat (pesan info di layar). **Sebagian terjawab 10-10-2026**: work owner memberi `Claim Fac In/AcceptanceNotePDF.xml` -> nota akseptasi tombol Acceptation lini Fire / Aneka / Golf / Marine Cargo dibangun (PARITAS baris Acceptation, §8 butir 24); stream `AcceptanceNotePDFMBU` / `AcceptanceNotePDFTravel` / `AcceptanceNotePDFPA` dan aliran CFS / PLA / DLA masih tidak diekspor |
| **OQ-CFI-21** | Activity `ASM-FW-GCNMFW-Data-ObjectItem.GeneratePLA` (PLA fac retro) | Nomor PLA fac belum terbit; `IsPicTransfer = 1` disetel (`[inferensi]` pasangan `GeneratePLATreaty_Act` 19) supaya Send to PIC Claim kasus retro dapat jalan |
| **OQ-CFI-22** | Grid penerima `Email.RecipientList` + subjek PLA / DLA, isi surel (`SendEmailDLA_ACT` memuat kredensial tertulis mati — tidak disalin) | Penerima tidak dibangun; muatan outbox surel hanya pengenal (pola Komite Claim Prop) |
| **OQ-CFI-25** | Harness `ClaimCommittee` (grid item PA / Travel "Send / Sent Claim to committee"), local action `ShowRetro` kelas SpreadingRisk (spreading adjustment); `DLAList` / `RetroList.TotalClaim = "123"` tanpa pembaca di korpus | Tombol nonaktif; `DLAList` tidak dibangun |
| — | Prompt values `associated` (lihat OQ-CFI-18) | — |

## C. DBA

| OQ | Pertanyaan | Akibat sekarang |
| --- | --- | --- |
| **OQ-CFI-11** | Polis `FACINPRODUCTION` yang tidak punya `JSON_POLIS.DATA_JSON` tidak dapat disalin (`CopyNB_Act` membaca halaman polis dari JSON) | Ditolak terang di Choose Polis (`OQ-CFI-11` di pesan) |
| **OQ-CFI-30** | **NB Fac In sistem baru** (`modul/nbfacin`, menyala) tidak menulis `FACINPRODUCTION` / `JSON_POLIS` — polis yang lahir di sistem baru tidak ditemukan Choose Polis (prompt §4) | Choose Polis hanya polis warisan. Perlu sumber polis bersama (keputusan lintas modul; nbfacin tidak disunting) |
| — | Migrasi `560`–`567` + `978` menunggu `-migrate` work owner; kueri baru ke objek itu dicoba read-only di DEV sesudahnya (`ORA-00942` / `ORA-00904` sekarang wajar) | — |

## Ditutup tanpa keputusan (catatan)

| OQ | Catatan |
| --- | --- |
| **OQ-CFI-16** | View Retro List menampilkan `OfferFacIn.FacRetroList` polis yang dibaca ulang dari `JSON_POLIS`; pengganti treaty `CheckLimit_Act1` 9.1.1.4 disimpan klaim (baris retro tingkat klaim, `ClaimData.FacRetroTreaty`) dan dipasang ulang setiap muat (`TerapkanRetroTreaty`). Tidak perlu keputusan |
