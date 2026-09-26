# 03: Baris `AdjustmentList` + Save ke Outstanding

**Status:** claimed

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

---

## Implementasi — 26 September 2026 malam (prasyarat, sesi batch 02–06)

**Status: `claimed`** — **0 dari 26 AC tertutup.** Yang dikerjakan commit ini adalah **dua prasyarat
§10 brief**, bukan AC tiket ini. Dinyatakan terus terang supaya tidak terbaca sebagai kemajuan yang
tidak ada.

Test Go **111 → 112**, nol FAIL; test bertag `db` **23 → 24**, seluruhnya SKIP; 88 modul.

### ⛔ Cacat tersembunyi yang ditutup

`services/pendaftaran_db_test.go` (tiket 02) memanggil `AmbilUntukKlaim`, yang membaca
`M_LIFE_PREMIUM_DETAIL` — tabel yang **skema uji tidak pernah buat**. Hari ini ia SKIP bersama yang
lain, jadi cacatnya tak terlihat; **dengan Oracle ia akan GAGAL di pembacaan peserta**, bukan
menguji pendaftaran. Tiga test db tiket 02 karena itu selama ini kosong isinya.

Sekarang skema uji membuat tiruannya: **24 kolom `kolomSalin` + `EDMSTATUS`**, tipe `[data DBA]`
dari katalog — bukan diturunkan dari nama. Fixture dua peserta: satu ber-`EDMSTATUS` **NULL**
(new business, harus muncul) dan satu **`Batal`** (harus disaring keluar), sehingga penyaring hidup
akhirnya teruji terhadap Oracle dan bukan hanya terhadap dirinya sendiri.

⛔ Kolom `KTP` **tidak ikut ditiru**: ia tidak pernah dibaca, dan tabel uji yang menyediakan tempat
untuk nomor identitas adalah undangan.

### ⛔ Satu cacat saya sendiri, ditemukan saat memverifikasi tipe kolom

`kolomSalin` membaca **`STNC` apa adanya**, padahal katalog menyebutnya **`DATE`** (kolom 29).
Membaca kolom tanggal tanpa `TO_CHAR` membuat bentuknya bergantung `NLS_DATE_FORMAT` sesi — jebakan
yang sama yang diperangi sepanjang tiket 14, dan saya sendiri yang memasangnya di tiket 02 lanjutan.
Sudah dibungkus `TO_CHAR`.

⚠️ `[terbuka]` Kolom tujuannya di `003` bernama `STNC_TREATY` dan bertipe `VARCHAR2(64)`; sumbernya
`DATE`. Selisih tipe itu belum diputuskan siapa pun, dan executor tidak mengubah DDL tanpa keputusan.

### Penjaga posisi `kolomSalin`

`salinKePeserta` membaca hasil `SELECT` lewat **indeks tetap 0–23**. Satu kolom yang disisipkan di
tengah menggeser seluruh sisanya **tanpa satu pun galat** — nilai hanya mendarat di medan yang
salah, dan itu baru terlihat jauh di hilir kalau pernah terlihat. `TestUrutanKolomSalinDikunci`
mengunci cacahnya (24) dan nama kolom pada tiap posisi.

### Tiruan tabel treaty

`RETROCESSIONLIFE` ditiru **dengan tipe aslinya** — di instance pengembangan ia **VIEW** ber-13
kolom yang seluruhnya `VARCHAR2(4000)`, termasuk `PERCENTSHARE`, `RATE`, `COMMISION`, `OVR_COMM`.
⭐ Tiruannya sengaja **tidak** dibuat `NUMBER` yang "lebih benar": justru jalur *teks → `ParseDecimal`
→ laporkan yang gagal* itulah yang perlu diuji, dan tiruan bertipe `NUMBER` membuat Oracle mengurai
angkanya lebih dulu sehingga pembacanya tidak pernah menemui teks. `TREATYYEAR_LIFE` ikut, 7 kolom.

### ⛔ `RATE_LIFE` tidak ditiru — dan itu memblokir satu masukan rumus spreading

Katalog baru memuat **enam kolom pertamanya** (`ID`, `IDUSEDBY`, `USEDBY`, `TYPE`, `GENDER`,
`CONTRACT`), dan **tidak satu pun di antaranya kolom rate**. Menirunya berarti mengarang bentuk, dan
membaca `RATE` dari tabel yang bentuknya dikarang berarti mengarang angkanya.

Akibatnya pada AC 22: rumus `PREMIUM_SPREADED_GROSS = RATE × (1 + EM_PERCENT) × AMOUNT` kehilangan
sumber `RATE`-nya. ⚠️ `[data DBA]` — daftar kolom `RATE_LIFE` beserta tipenya diperlukan sebelum
pembacanya ditulis.

### Penjaga AC 29 dipertajam, bukan dilonggarkan

Tiruan tabel peserta **mendeklarasikan** kolom `EDMSTATUS`, dan penjaga lama menuduhnya sebagai
penyaring kedua. Aturannya kini membedakan **menyaring** dari **menyebut**: yang dilarang adalah
bentuk `EDMSTATUS IS`/`NOT IN`/`IN`/`=` di luar pembacanya. ⚠️ Percobaan pertama saya memakai pola
`"EDMSTATUS)"`, yang cocok dengan **daftar kolom INSERT** — penjaga yang menuduh hal yang bukan
aturan. Diperbaiki, lalu dibuktikan masih menangkap penyaring kedua yang sungguhan.

### Yang BELUM dikerjakan — seluruh AC tiket ini

26 AC masih terbuka. Yang belum ada sama sekali: `Adjustment.Tambah` dan pewarisan delapan kolom,
`Spreading.Hitung`, pembaca treaty, gerbang dokumen, langkah migrasi `010` (butir **ad**), kedua
pintu HTTP, dan seluruh frontend-nya. Sensus `.DocumentList` untuk **ad** baru sampai pada daftar
berkas korpusnya, belum pada kolomnya.
