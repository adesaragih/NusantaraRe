# Kesiapan Claim — Life untuk `to-spec`

Tanggal: 2026-09-14 (perhitungan ke-4, setelah Ronde 4)
Dasar: grilling Ronde 1–4, 12 ADR, `CONTEXT.md`, register `discovery/open-questions.md`
(OQ-001…OQ-063)

---

## Putusan: **MATANG untuk `to-spec`**

**Tidak ada lagi pemblokir sejati. Frontier Claim — Life kosong.**

Empat ronde grilling menutup seluruh pertanyaan yang, bila tidak dijawab, akan memaksa spesifikasi
menebak. Yang tersisa adalah **10 OQ non-pemblokir** — pertanyaan yang jawabannya memengaruhi
*detail* implementasi, bukan bentuk spesifikasi. Semuanya akan **ditandai "terbuka" di dalam spec**
pada tempatnya masing-masing.

Perjalanan empat perhitungan:

| Blok | Ke-1 | Ke-2 | Ke-3 | **Ke-4 (sekarang)** |
| --- | --- | --- | --- | --- |
| **A** — pemblokir sejati | 3 | 1 | 1 | **0** ✅ |
| **B** — fakta hilang, non-pemblokir | 6 | 8 | 8 | **7** |
| **C** — implementasi | 3 | 3 | 3 | 3 |
| **D** — keputusan desain | 1 | 1 | 3 | **1** (sisa: aturan "klaim selesai") |
| ADR | 10 (1 draft) | 11 | 11 | **12, seluruhnya accepted** |

---

## Keputusan yang sudah lengkap — 12 ADR

| ADR | Keputusan |
| --- | --- |
| 0001 | Batas konteks Claim Life ↔ Komite Life; **tiga kontrak** (penyerahan, jalur balik, tabel bersama) |
| 0002 | RBAC tiga peran |
| 0003 | Uang non-float; 8 kolom uang + `EM_PERCENT` persen |
| 0004 | Endpoint sebagai env var |
| 0005 | `IsPEGAPROD` → flag lingkungan |
| 0006 | Penomoran via `PROC_GENERATE_SEQUENCE_NUMBER`; 2 SQL lama tidak dimigrasi |
| 0007 | Jejak audit siapa+kapan setiap transisi dan jalur balik |
| 0008 | Efek keluar asinkron, tidak memblokir, dengan antre-ulang |
| 0009 | Migrasi penuh; koeksistensi ditolak |
| 0010 | Penyimpanan berkas tetap Google Storage |
| 0011 | **Unit status = baris `AdjustmentList`**; dua sumber reject; nilai `2` selalu berarti "baris ditolak"; `SaveAdjustment_Act` dead |
| 0012 | **Wewenang kirim ke Komite bergantung `Type`** — paritas, risiko RBAC tercatat |

**Mesin status lengkap dan dapat ditulis utuh:** enam langkah, tiga peran, dua sumber penolakan
dengan makna setara pada tingkat baris, terminal per baris, klaim tidak terminal, revisi = baris
baru, dan pengecualian `TP`/`TR` pada wewenang.

**Kosakata lengkap:** `CONTEXT.md` memuat tahapan, tiga peran, `STS_REJECT` beserta mesin statusnya,
`AdjustmentList`, `PremiumListDetail`, `AcceptStatus`, `KomiteLoop`, `Type`/`TP`/`TR`, `SendtoAdmin`,
`SendtoMedical`, `ContentNote`, `BusinessCode` `L1`–`L21`, dan batas konteks — seluruhnya dengan
baris `_Bukti_`.

---

## Yang akan ditandai "terbuka" di dalam spec — 10 OQ non-pemblokir

Tidak satu pun menghalangi penulisan spec. Masing-masing ditandai pada bagian yang menyentuhnya.

### Fakta bisnis yang hilang — Product+UW

| OQ | Bagian spec | Akibat bila tetap terbuka |
| --- | --- | --- |
| **OQ-020** (sisa) | Arti `QP` / `QR` | Perilaku keduanya **sudah terbaca** (jendela gross, pergeseran nol); hanya namanya yang belum. `[dugaan]` satu konfirmasi dapat menutupnya |
| **OQ-032** | Penyerahan ke Komite | Apa yang menentukan `KomiteLoop` — jumlah putaran sebelum keputusan final |
| **OQ-037** | Pemilihan roster | Ambang `Param.LIMIT_BOTTOM` pada `FilterEmailKomiteWithLimit` |
| **OQ-060** | Tipe uang | Apakah keseragaman mata uang per klaim adalah aturan atau kebetulan |
| **OQ-035** | Efek keluar | Kepemilikan `serviceInsertArasapasClaimLife_act` — satu salinan, dua konteks |
| **OQ-039** (sisa) | — | **Tertutup untuk Claim — Life**; tersisa untuk 3 modul Claim lain, di luar konteks ini |

### Batas pengetahuan database — DBA

| OQ | Bagian spec | Akibat bila tetap terbuka |
| --- | --- | --- |
| **OQ-001** | Model data | Tidak ada DDL — tipe, presisi, skala kolom; bentuk penyimpanan baris adjustment |
| **OQ-002** | Penomoran + unggah berkas | Kontrak `PROC_GENERATE_SEQUENCE_NUMBER` dan `GET_TOKEN_STORAGE` |
| **OQ-013** | Transaksi | `COMMIT` berada di dalam blok PL/SQL |
| **OQ-047** | Efek keluar | Isi tabel `M_LINK_SERVICE` — daftar endpoint sebenarnya |

### Lingkungan — IT-infra

| OQ | Bagian spec |
| --- | --- |
| **OQ-018** (sisa) | Lingkungan `pega-nusre`; bucket production vs dev |

---

## Satu keputusan desain yang tersisa — bukan pemblokir, tetapi harus diambil saat spec ditulis

### Kapan sebuah klaim dianggap "selesai"?

Konsekuensi langsung dari Ronde 4: karena `STS_REJECT = 2` **selalu** berarti "baris ini ditolak"
dan **tidak pernah** "klaim selesai", maka **tidak ada nilai status yang menandakan klaim selesai.**
Selesainya klaim adalah **keadaan turunan** dari kumpulan barisnya — misalnya "ada satu baris
bernilai `1`", atau "tidak ada lagi baris bernilai `0`".

Ini **bukan fakta yang hilang dari korpus** (sistem lama memang tidak menyimpannya), melainkan
aturan turunan yang perlu dinyatakan saat spec ditulis. Karena itu ia **tidak** menjadi OQ.

Dua hal lain yang perlu dinyatakan di dalam spec, keduanya sudah berbukti dan berkeputusan:
- **Penegakan peran di lapisan layanan** — sistem lama tidak seragam (**ADR-0002**, **ADR-0012**).
- **Validasi dibawa apa adanya** — 16/16 FlowAction `Claim Life` membawa rule validasi; termasuk
  pergeseran **+1 hari** pada jendela DOL untuk `TP`/`TR` (`ValidasiDOL_Act`). Dibawa sebagai
  paritas, sejalan dengan ADR-0012; bila kelak dianggap keliru, itu keputusan terpisah.

---

## Kesimpulan

**Claim — Life MATANG untuk `to-spec`.** Frontier kosong: tidak ada pertanyaan tersisa yang, bila
tidak dijawab, akan memaksa spec menebak.

`spec.md` dan `issues/` **belum disentuh** — menunggu **persetujuan eksplisit Anda** untuk
menjalankan `to-spec`.
