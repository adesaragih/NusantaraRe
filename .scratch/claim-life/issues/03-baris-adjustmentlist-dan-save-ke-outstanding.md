# 03: Baris `AdjustmentList` + Save ke Outstanding

**Status:** ready-for-agent

**Blocked by:** 02 (register klaim + penomoran), **14 (skema relasional klaim — PREFACTOR)**

> ⚠️ **Diselaraskan 2026-09-16 — revisi penyimpanan.** Dua perubahan mengikat:
> **(1)** baris adjustment melekat pada **PESERTA** (`T_CLAIMLF_PREMIUMLIST_DETAIL`), bukan pada
> klaim; **(2)** dokumen **per peserta** menjadi **gerbang simpan** ke Outstanding.
> Lihat blok AC "Penyimpanan relasional" dan "Dokumen per peserta" di bawah.

## Hasil & nilai pengguna

Sebagai **ReasLifeAdmin**, saya dapat menginput baris **AdjustmentList** pertama pada sebuah klaim
dan menyimpannya ke Outstanding — sehingga klaim itu tercatat sebagai sedang berjalan dan menjadi
antrean kerja yang terlihat. *(User story 3 di spec)*

Ini tiket yang **memperkenalkan unit keputusan** sistem ini: barisnya, bukan klaimnya.

## Area codebase

`internal/models` (baris adjustment sebagai entitas dengan statusnya sendiri), `internal/repository`
(penulisan baris), `internal/services` (aturan penyimpanan ke Outstanding), `internal/handlers`
(endpoint input baris), `frontend/` (formulir dan daftar baris adjustment).

## Rule Pega sumber

| Rule | Identitas | Perilaku yang ditiru |
| --- | --- | --- |
| `Claim Life/Activity/SaveOutStandingLife_Act.xml` | `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` **satu-satunya penulis nilai `STS_REJECT = 0`** di seluruh korpus Claim Life |
| `Claim Life/Section/AdjustmentDetail_Section.xml` | `ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `ADJUSTMENTDETAIL_SECTION` / `RULE-OBJ-HTML-SECTION` | `[terverifikasi]` tombol "Save to Outstanding", dan bentuk baris adjustment |
| `Claim Life/Activity/SetIndexAdjustmentList.xml` | `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `SETINDEXADJUSTMENTLIST` / `RULE-OBJ-ACTIVITY` | `[terverifikasi]` baris baru mewarisi **8 kolom** dari `AdjustmentList(1)`: `SHARE_NUSANTARA_RE`, `CEDING_RETENTION`, `SUM_REASURED`, `SUM_INSURED`, `SHARE_RETRO`, `RETROCEDED_SHARE`, `CURRENCYID`, `CURRENCY` — **tanpa** `STS_REJECT` |

## ADR terkait

**ADR-0011** (unit status = baris `AdjustmentList`), **ADR-0003** (8 kolom uang + `EM_PERCENT`
persen; nama kolom tidak dapat dipakai menebak sifatnya).

## Acceptance criteria

- [ ] Baris `AdjustmentList` yang baru disimpan ke Outstanding berstatus **Outstanding** (`0`).
      *(AC 1 spec)*
- [ ] Sebuah klaim dapat memuat **beberapa** baris adjustment sekaligus, masing-masing dengan
      statusnya sendiri. *(AC 25 spec)*
- [ ] Baris yang ditambahkan setelah baris pertama **mewarisi delapan kolom** di atas dari baris
      pertama, dan **tidak** mewarisi status. *(AC 5 spec)*
- [ ] Kedelapan kolom uang diperlakukan sebagai uang; `EM_PERCENT` **tidak** diperlakukan sebagai
      uang. *(AC 23 spec)*
- [ ] Tidak ada nilai uang sebagai *binary floating point* di lapisan mana pun maupun di kontrak API.
      *(AC 22 spec)*

### Penyimpanan relasional ⚠️ BARU 2026-09-16 — spec §2b

- [ ] ⚠️ Setiap baris adjustment menunjuk **satu peserta** lewat `PREMIUM_LIST_DETAIL_ID`. Test yang
      menemukan FK adjustment menunjuk **header klaim** **gagal**. *(AC 33 spec; penyimpangan
      sadar 2 — perbaikan relasi, bukan peniruan)*
- [ ] ⚠️ Dua peserta dengan masing-masing dua putaran adjustment menghasilkan **empat baris yang
      seluruhnya dapat ditelusuri ke peserta yang benar**. *(AC 34 spec)*
- [ ] ⚠️ Peserta menyimpan **penanda dipilih-untuk-diklaim** (`IS_CHECK`) — inilah penyimpan aturan
      "hanya peserta yang diklaim". *(AC 39 spec; penyimpangan sadar 3)*
- [ ] ⚠️ Tanggal **diterima**, **konfirmasi**, dan **penyelesaian** tersimpan **per peserta**, bukan
      di header. *(AC 40 spec)*
- [ ] Peserta menyimpan `STATUS`, `RECOMMENDATION`, `SOURCE_ID`, dan `CEDING_RETENTION`.
      *(AC 41 spec)*
- [ ] ⚠️ `STS_REJECT` baris adjustment diisi **nilai sebenarnya menurut aksi** — Admin insert
      Outstanding → `0`; SPV tambah Outstanding → `0`. Test yang menemukan nilai di-hardcode
      **gagal**. *(AC 42 spec; penyimpangan sadar 4)*
- [ ] ⚠️ `ACCEPTATION_DATE` **tidak** distempel saat insert; ia diisi **tanggal akseptasi
      sebenarnya** saat baris benar-benar diaksep. *(AC 43 spec; penyimpangan sadar 4)*
- [ ] Peserta beserta seluruh baris adjustment-nya tersimpan dalam **satu transaksi**. *(AC 49 spec)*
- [ ] Baris adjustment menyimpan **nama bank**, **id bank**, dan **nomor rekening**, dan ketiganya
      dapat diisi dari layar rincian adjustment. *(AC 56 spec; `[terverifikasi]` class
      `ASM-FW-GISFW-Data-AdjustmentLife`, tampil di `Claim Life/Section/AdjustmentDetail_Section.xml`)*
- [ ] Ketiga field bank **boleh kosong saat Save ke Outstanding** — ia baru menjadi gerbang pada
      **penyerahan ke Komite** (tiket 10). *(AC 57 spec)*

### Dokumen per peserta — **gerbang simpan** ⚠️ BARU 2026-09-16

`[terverifikasi]` `Claim Life/Activity/SaveOutStandingLife_Act.xml`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`) beriterasi atas
`.DocumentList` **per peserta** dan menolak simpan dengan
`"The document hasn't been uploaded person number "+<nomor>` serta
`"Documents are incomplete, please complete the documents"`.

- [ ] ⚠️ Dokumen tersimpan **per peserta** di `DOCUMENT_CLAIM` dan dapat dibaca dengan `SELECT`
      biasa — **bukan** lewat mekanisme lampiran bawaan. *(AC 44 spec; penyimpangan sadar 5)*
- [ ] ⚠️ Menyimpan ke Outstanding **ditolak** bila ada peserta yang dokumennya belum lengkap, dengan
      pesan yang **menyebut peserta mana**. *(AC 45 spec; penyimpangan sadar 5)*
- [ ] Kolom isian `DOCUMENT_CLAIM` **diturunkan dari sensus `.DocumentList`** pada activity di atas,
      dan **keputusannya dicatat** — **jangan tebak dari nama tabel**. *(tiket 14 §Catatan)*
- [ ] Halaman React menampilkan daftar baris adjustment dengan status masing-masing sebagai kata,
      bukan angka. *(AC 26 spec)*

### Spreading adjustment ⚠️ BARU 2026-09-16

`[terverifikasi]` `Claim Life/Activity/SpreadingClaimLife_Act.xml`
(`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` / `SPREADINGCLAIMLIFE_ACT` / `RULE-OBJ-ACTIVITY`) menghitung dan
mengisi `.SpreadingList` (per treaty-year, class `ASM-FW-GISFW-Int-TREATYYEAR_LIFE`) beserta
`.RetroLifeList` (per reinsurer, class `ASM-FW-GISFW-Int-RETROCESSIONLIFE`) pada **baris
adjustment**.

- [ ] ⚠️ Menyimpan baris adjustment **menghitung dan menyimpan** spreading-nya: satu baris per
      treaty-year, dan di bawahnya satu baris per reinsurer. *(AC 59 spec; tiket 14; penyimpangan sadar —
      spreading dibekukan)*
- [ ] ⚠️ Setiap baris spreading menunjuk **satu baris adjustment**; setiap baris spreading retro
      menunjuk **satu baris spreading**. Test yang menemukan keduanya menggantung pada peserta atau
      pada header klaim **gagal**. *(AC 58 spec; tiket 14)*
- [ ] Nilai turunan tersimpan sesuai perhitungan yang terbukti: `AMOUNT = RetrocadedShare ×
      PERCENT_SHARE ÷ 100`, dan `PREMIUM_SPREADED_GROSS = RATE × (1 + EM_PERCENT) × AMOUNT`.
      *(`[terverifikasi]` `SpreadingClaimLife_Act`)*
- [ ] ⚠️ `[terbuka]` **`PREMIUM_SPREADED_NET` tidak dinyatakan selesai** sebelum Product + UW
      menetapkan rumusnya. `[terverifikasi]` korpus memuat **dua** cabang di rule yang sama —
      `GROSS − Comm` dan `GROSS − Discount − Comm` — dan **mana yang berlaku tidak terbaca**.
      **Jangan tebak.** *(tiket 14 §Blocker)*
- [ ] ⚠️ Nilai spreading **dibekukan**: perubahan master treaty sesudahnya **tidak mengubah** angka
      yang sudah tersimpan pada adjustment itu. *(AC 59 spec)*
- [ ] Adjustment tanpa retrosesi tersimpan dengan **nol baris** spreading — **bukan** kegagalan.
- [ ] Baris adjustment beserta seluruh spreading dan spreading retro-nya tersimpan dalam **satu
      transaksi**. *(AC 49 spec)*

## Catatan

`[terverifikasi]` Empat kolom bernama `SHARE_*` / `*_SHARE` ternyata **uang**, bukan rasio — nama
kolom di korpus ini terbukti menipu (bandingkan `STS_REJECT`, yang nilai `1`-nya berarti *diaksep*).
Klasifikasi mengikat ada di **ADR-0003**; jangan menyimpulkan dari nama.

`[terbuka]` **OQ-060** (pemilik **Product+UW**) — apakah seluruh baris satu klaim wajib bermata uang
sama. `[terverifikasi]` dalam praktiknya seragam karena `CURRENCY` disalin dari baris 1, tetapi
strukturnya membolehkan campur. **Tidak memblokir tiket ini**; memblokir bentuk akhir tipe uang.

## Perintah verifikasi

```
go test ./internal/...
cd frontend && npm test
make check
```
