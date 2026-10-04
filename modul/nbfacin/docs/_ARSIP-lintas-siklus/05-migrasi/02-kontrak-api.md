# Kontrak API — Facultative Inward

> Sumber: sintesis dokumen discovery di `D:\migrasi\RNM\OUTPUT\`, dibangun dari korpus
> `D:\migrasi\RNM\{NB FacIn, RNW Fac In, Endorsment Fac In}\`.
> Label mengikuti `CLAUDE.md` §3.
>
> **Status dokumen: [usulan].** Korpus **tidak memuat** kontrak API — Pega tidak memisahkan
> backend dan frontend seperti ini. Yang direkam korpus adalah *perilaku* dan *bentuk data*;
> pemetaannya ke HTTP adalah rancangan baru. Setiap endpoint di bawah menyebut temuan korpus yang
> membentuknya, sehingga yang mana rekaman dan yang mana rancangan tetap dapat dibedakan.

---

## 1. Aturan lintas-endpoint yang **mengikat**

### 1.1 Uang selalu string desimal

`[terverifikasi]` (`04-aturan/02-formula-dan-status.md` §7) `number` JSON adalah IEEE-754 double di
sebagian besar parser; produk antara rumus premi mencapai orde 10²³. Karena itu:

```json
{ "tsi": { "amount": "25000000000.0000", "currency": "IDR" } }
```

**Bukan** `"tsi": 25000000000.0`. Aturan ini berlaku untuk uang, TSI, rate, persen, dan limit —
tanpa kecuali.

Format: titik sebagai pemisah desimal, tanpa pemisah ribuan. Format tampilan Indonesia (koma) adalah
urusan lapisan presentasi React, bukan kontrak.

⚠️ `[terverifikasi]` Sistem lama menerima nilai sebagai **string berkoma desimal** di 1.057 titik.
Parser di batas input wajib menerima keduanya; normalisasi terjadi sekali, di batas.

### 1.2 Nama field mengikuti `DATA_JSON`

`[terverifikasi]` (`02-data/01-model-domain.md` §7) Agregat case diserialkan ke
`JSON_POLIS.DATA_JSON` dengan nama properti identik. Payload API yang memetakan langsung ke agregat
**mempertahankan nama field apa adanya**, termasuk yang namanya menyesatkan.

Pengecualian yang disengaja — hanya untuk field yang **tidak** ikut serialisasi: `LetterNo`
dipaparkan sebagai `next_approver_position` (`CONTEXT.md` §4).

### 1.3 Diskriminator siklus adalah atribut, bukan path

`[terverifikasi]` (§1.1 `05-migrasi/01-arsitektur-target.md`) Satu basis rule melayani tiga siklus.
Karena itu **tidak ada** `/api/nb/...`, `/api/rnw/...`, `/api/edm/...`. Yang ada satu himpunan
endpoint dengan `status_business` sebagai atribut case.

### 1.4 Tangga akseptasi = satu transisi per panggilan

`[terverifikasi]` (`01-flow/04-mesin-akseptasi.md` §0) Satu keputusan manusia = satu putaran. API
**tidak** menyediakan endpoint yang menghitung seluruh rantai approver sekaligus — bentuk itu tidak
ada di sistem lama dan akan berperilaku berbeda saat rantai terputus.

### 1.5 Endpoint luar tidak pernah literal

`[terverifikasi]` `CLAUDE.md` §4.4. Alamat servis luar diambil dari tabel Oracle `M_LINK_SERVICE`
(kunci `KATEGORI_1` + `KATEGORI_2`) saat runtime. **Isi tabel itu tidak ada di korpus.**

---

## 2. Sumber daya utama

### 2.1 Case (`OfferFacIn`)

```
GET    /api/v1/cases/{id}
POST   /api/v1/cases                       buat case baru
PATCH  /api/v1/cases/{id}                  simpan perubahan parsial
GET    /api/v1/cases?queue=&position=      daftar kerja per antrean
```

Bentuk respons — ringkas; agregat penuh mengikuti `02-data/01-model-domain.md`:

```json
{
  "id": "NB-000123",
  "status_business": 1,
  "state": {
    "current_queue": "ReasFacInUnderwriting",
    "next_approver_position": "SENIORUW",
    "last_decision": 4
  },
  "offer_fac_in": { }
}
```

⚠️ `[pertanyaan terbuka]` **`last_decision` (`ProposalAcceptStatus`) dipaparkan sebagai angka mentah**
karena artinya belum diketahui. `DecisionTable/IsUWAccepted` mendeklarasikan enam hasil (`confirm`,
`reject`, `ask`, `banding`, `revise`, `decline`) tetapi **baris hasilnya tidak ikut terekspor**, jadi
pemetaan angka→arti tidak ada di korpus. Memberi nama pada nilai-nilai ini sekarang berarti menebak
arah **setiap** keputusan underwriting. Enum bernama ditambahkan setelah bisnis menjawab
(`02-data/03-batas-pengetahuan.md`).

⚠️ `[terverifikasi]` Bukti bahwa menebak berbahaya: nilai `4` ditulis oleh rule bernama
`SetBandingProposal_DT` (banding) tetapi **diuji** oleh `When/IsFacout.xml` (fac out) —
`RNW Fac In\When\IsFacout.xml:148` → `pyWorkPage.ProposalAcceptStatus = 4`. Satu nilai, dua nama
pemakai yang temanya berbeda.

⚠️ `[terverifikasi]` `ProposalAcceptStatus` ditulis sebagai angka `4` tetapi dibaca sebagai string
`"4"` di `SetBanding_ACT`. Kontrak harus memilih satu representasi dan mendokumentasikannya;
perbedaan ini adalah kandidat perbaikan, bukan sesuatu yang diperbaiki diam-diam.

### 2.2 Keputusan akseptasi — satu transisi

```
POST   /api/v1/cases/{id}/decisions
```

```json
{ "decision": 4, "note": "…" }
```

Respons mengembalikan keadaan **setelah** satu transisi:

```json
{
  "state": {
    "current_queue": "ReasFacInDivHead",
    "next_approver_position": "KADIVTEKNIK",
    "last_decision": 4
  },
  "ladder_finished": false
}
```

`[terverifikasi]` `ladder_finished: true` ketika tidak ada `To*` yang cocok dengan
`next_approver_position` — cabang `Else` gerbang. **Itu penyelesaian normal, bukan galat**; API
mengembalikan `200`, bukan `4xx`.

### 2.3 Banding

```
POST   /api/v1/cases/{id}/appeals
```

`[terverifikasi]` Jalur terpisah dari keputusan biasa: `SetBanding_ACT` memetakan
`next_approver_position` / `current_queue` → nama tiket Pega, lalu melompatkan alur. Di target,
lompatan tiket menjadi transisi bernama di mesin keadaan.

⚠️ `[pertanyaan terbuka]` `GetAksepBanding_SQL` selalu memakai tabel limit Property untuk mencari
jabatan atasan, **apa pun lini bisnisnya**. Bug atau memang hierarki jabatan hanya ada di satu
tabel — belum terjawab. Kontrak banding tidak dapat difinalkan sebelum ini dijawab.

### 2.4 Premi

```
POST   /api/v1/cases/{id}/premium:calculate
```

⚠️ `[terverifikasi]` Perhitungan premi **bercabang menurut siklus** dan presisinya berbeda antar
rule: aritmetika yang sama dibulatkan 4 desimal di `HitungPremi_FacInDT` dan 20 desimal di
`CountPremi_ACT`. Respons wajib menyertakan asal rule agar rekonsiliasi paralel run mungkin:

```json
{
  "premium": { "amount": "1234.5678", "currency": "IDR" },
  "rounding": { "decimals": 4, "source_rule": "HitungPremi_FacInDT" }
}
```

⚠️ `[pertanyaan terbuka]` **MEMBLOKIR** — satuan `.Rate` tidak konsisten: satu rule membagi dengan
`1e9` (konsisten dengan rate per mille), rule lain dengan `1e4` (konsisten dengan persen). Satu ordo
besaran 10 berbeda. Endpoint ini tidak dapat diimplementasikan sebelum bisnis menjawab.

### 2.5 Produksi

```
POST   /api/v1/cases/{id}/production
```

`[terverifikasi]` Konversi ke produksi menulis pasangan **nilai-sesudah** dan **delta** berdampingan
(akhiran `_MENJADI` / `_SELISIH`). Payload membawa pasangan itu, **bukan** satu nilai — memilih satu
saja menghilangkan informasi yang dipakai endorsement.

`[terverifikasi]` Penomoran versi polis bersifat *append-only* (`PRODKE = COUNT(NOPOLIS) − 1`).
Tidak ada endpoint yang memperbarui baris polis lama.

---

## 3. Registry predikat — bukan endpoint publik

`[terverifikasi]` 226 nama rule `When` unik dievaluasi **di dalam** `services`, tidak dipaparkan
sebagai API. Alasannya: predikat adalah detail implementasi tangga, dan memaparkannya membekukan
bentuk internal ke dalam kontrak.

⚠️ Jumlah rule yang kondisinya tidak terbaca sama sekali adalah **0** — mengoreksi `CLAUDE.md` §4.5.
Lihat `05-migrasi/01-arsitektur-target.md` §2.1.

---

## 4. Penanganan galat

| Situasi | HTTP | Catatan |
| --- | --- | --- |
| Tangga selesai (tidak ada approver berikutnya) | `200` | **Bukan galat** — `[terverifikasi]` |
| Predikat yang kondisinya tidak diketahui dievaluasi | `500` | Gagal keras (`CLAUDE.md` §4.5). **Jangan** `return false` — aplikasi yang berjalan dan salah diam-diam lebih berbahaya. |
| Nilai enumerasi di luar yang terbaca korpus | `422` | Celah enumerasi adalah temuan, bukan nilai valid. |
| Nilai uang dikirim sebagai `number` JSON | `400` | Menjaga §1.1 di batas, bukan di dalam. |

---

## 5. Yang **belum** dapat dikontrakkan

| Area | Alasan |
| --- | --- |
| Endpoint spreading / capacity / scoring | Alurnya belum didiscovery (`05-migrasi/01` §5). |
| Payload servis luar | Isi `M_LINK_SERVICE` dan aktivitas servis tidak ada di korpus. |
| Enum bernama untuk `last_decision`, `Type`, `EdmType` | Baris hasil DecisionTable tidak terekspor. **Memblokir.** |
| Kontrak render per-layar | 2.306 `pyPreDataTransform` — memindahkannya ke frontend mengubah perilaku; keputusan arsitektur yang belum diambil. |

---

*Daftar lengkap pertanyaan yang memblokir ada di `02-data/03-batas-pengetahuan.md` dan
`05-migrasi/03-risiko-dan-pertanyaan.md`.*
