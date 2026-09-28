# Paritas aksi — PremiumList Life

> Berkas ini lahir 28-09-2026 untuk sensus remark (GILIRAN-12 paket 2); modul ini sebelumnya tanpa
> berkas paritas. Paritas layar per tiket tetap di `issues/`.

## Sensus remark 28-09-2026 (GILIRAN-12 paket 2)

Setiap activity PremiumList Life dicetak beserta `pyStepsBlockName`; `//` = ter-remark (beserta
anaknya). 38 activity, **9** ber-remark, **36** langkah `//` (label `SHOWFAILURE`, `AA`, `Err`, `REM`,
`COMMITSTEP` dll. bukan remark).

**Cara hitung** `[terverifikasi]` — dua cara, jendela = seluruh `PremiumList Life/Activity/*.xml`:

| Cara | Perintah | activity | ber-remark | langkah `//` |
| --- | --- | ---: | ---: | ---: |
| A — pengurai XML (expat), per `rowdata` ber-`pyStepPageReference` yang `pyStepsBlockName` tepat `//` | skrip sensus sesi (cetak activity, nomor langkah, bNNN) | 38 | 9 | 36 |
| B — grep mentah per berkas | `grep -o "<pyStepsBlockName>//</pyStepsBlockName>" "<berkas>" \| wc -l`, dijumlah | 38 | 9 | 36 |

Alat diuji lebih dulu pada item yang sudah diketahui: `Claim Life/Activity/SaveOutStandingLife_Act.xml`
— 8 langkah `//` (dibaca ulang pemeriksa independen GILIRAN-11) — kedua cara menjawab 8. Label
lompatan lain (`EXIT`, `send`, `AA`, …) bukan remark dan tidak terhitung.

⚠️ **Cakupan celah langkah hidup DINYATAKAN:** sensus ini mendaftar langkah hidup yang tidak ditiru
untuk activity **ber-remark** saja. Activity tanpa remark tidak diaudit ulang langkah demi langkah di
giliran ini; statusnya tetap seperti baris paritasnya (bila ada) — **belum** audit penuh.


| Activity | Langkah `//` (b) | Ditiru? | Keputusan |
| --- | --- | --- | --- |
| `ProtectAccept` | 4.1 b720 (NoOffer) · 6 b2340 (+6.1–6.5; 6.3 b2756) | **ya — dibuang** (`models/polis_validasi.go`) | tiket 01 ralat; gerbang posisi dibetulkan |
| `ValidasiUploadPL_act` | 9.2 b8025 · 9.3 b8198 · 9.4 b8368 · 9.5 b8517 · 9.11 b9667 · 10 b13018 · 11 b13171 · 12 b13322 · 13 b13475 · 19 b14393 · 33 b16535 | **ya — dibuang** (sertifikat terpakai/ganda, `MEDICAL_STATUS`, lampiran; + duplikat nama+DOB) | tiket 04 ralat, AC disunting di tempat |
| `SubmitPremiumList_Act` | 15 b3398 (geser `MM.YYYY.` bila di atas tanggal closing) | tidak — tetapi `models/polis_periode.go` mengutip b3491 | kutipan diralat; premis "dua versi hidup" → OQ-PL-13 |
| `InsertJsonPolisLife_Act` | 16 b5294 `Commit` · 17 b5383 `Connect-REST InsertLifePremiumDetail` | tidak (Go memegang transaksi; salinan peserta lewat `repository/polis_warisan.go`, keputusan pl2) | nama layanan step 17 diralat di spec, 05b, 06 → OQ-PL-14 |
| `SavePremiumList_Act` | 8.3 b1891 · 8.5 b3704 · 8.6 b3863 · 8.7 b4050 · 8.9 b5887 · 11 b6411 · 13 b6642 | tidak | tercatat tiket 03/05a |
| `UploadCSVLifePremium_Act` | 5 b889 · 6 b1042 · 7 b1188 (+7.2 b2517, 7.4 b2797) · 8 b2922 | tidak (tinjauan tidak menulis) | — |
| `InputOfferLife_ACT` | 1 b354 (`ASMForceCaseClose`) | tidak — Decline menutup lewat flow | tercatat tiket 01 |
| `LoadAttachmentData` | 1 b276 · 2 b493 · 3 b661 · 4 b771 | tidak | — |
| `serviceInsertArasapasLife_act` | 9 b1597 (`UpdateErrorNoteJsonPolis`) | tidak (efek stub) | — |

**Langkah HIDUP yang tidak ditiru dan belum tercatat** (kini dicatat sebagai celah):

- `ProtectAccept` 3 (b624, `EmailTypePL==1` atau KELUAR) dan 8 (b3792, daftar mata uang kosong di
  Premium) — tiket 01. `ValidasiPenawaran` sendiri belum tersambung ke rute mana pun.
- `ValidasiUploadPL_act` 2 — mengisi 0 kolom uang yang kosong (sub-langkah 2.1–2.32) → **OQ-PL-12**;
  3–8 (b6874–b7767) — keberadaan medan header; pesannya ada di `models/polis_unggah.go`, belum dipakai.
- `SavePremiumList_Act` 8.1 (b1635, batas umur) dan 8.2 (b1833, batas sum insured; dilewati TP/TR
  berkode RNML-FL b1803), masing-masing dengan keluar di langkah 9 (b6146).
- `Calculate1_Act` 7.5–7.9 — `CekDoubleInsured` untuk akumulasi retensi ceding (`SumCRetensi` lawan
  `MaxCRetensi`).

**OQ baru (untuk work owner):**

- **OQ-PL-12** — langkah 2 `ValidasiUploadPL_act` menjadikan kolom uang yang kosong **0** sebelum
  pemeriksaan "HARUS ADA"; aplikasi menolaknya (ADR-U-0027: kosong bukan nol). Ikuti Pega (terima
  sebagai 0) atau pertahankan penolakan sebagai penyimpangan sadar?
- **OQ-PL-13** — "ikuti yang dari DB" berdiri di atas premis bahwa `SubmitPremiumList_Act` membaca
  ambang secara hidup; nyatanya ia hanya mengalir ke langkah 15 yang ter-remark. Pembaca hidup tabel
  itu adalah `PROC_GENERATE_SEQUENCE_NUMBER` (nomor PL); `ProdDateTime` di Pega memakai `>25`
  tertanam. Apakah ambang DB juga berlaku untuk `ProdDateTime`?
- **OQ-PL-14** — cerita spec 34 membuang `ConvertJsonNusareToProduction` dengan alasan step 17
  ter-remark; step 17 memanggil `InsertLifePremiumDetail`, dan `convertJsonNusareToProduction` hidup
  di `serviceInsertArasapasLife_act` langkah 5 (b963). Tetap dibuang?
