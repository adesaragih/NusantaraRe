# 19: Pemecah dokumen menjadi baris

> ## ⭐ PENAHAN GUGUR — 23 September 2026 sore
>
> `[keputusan work owner]` *"Abaikan `JSON_DATAGUIDE`, ikuti dari data yang digunakan di Activity dan Section."* ⭐ Daftar medan disusun ulang dari **302 berkas aturan** — **394 medan unik**, 232 di antaranya tampil di layar. Hasilnya di **`DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md`**. Data guide turun derajat jadi **penambal**, sah hanya untuk panjang maksimum per medan.
>
> ⭐ **`blocked` → `ready-for-agent`.** ⛔ **Nol butir `[data DBA]` tersisa di tiket ini.**
>
> Rinciannya: `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`.

---


**Status:** sebagian *(implementasi 2026-10-03, cabang `modul/nbtreatyin/implementasi`; semula: ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026 sore)*)*
~~**Blocked by:** **16** · **17** · **18** · ⛔ `[data DBA]` **daftar kolom lengkap** — panduan bentuk dokumen dari DBA terbukti **basi**~~ ⛔ **gugur 23-09-2026 sore**
**Menutup:** NB AC **26–44** *(19 AC)*
**Sumber:** `nb-treaty-in\spec-penyimpanan-relasional.md` ID-21..ID-31

## Hasil & nilai pengguna

Dokumen polis lama dipecah menjadi baris tabel: data umum, kuotasi, ceding, angsuran, penyebaran,
lapisan, dan riwayat usulan. Sesudah tiket ini, bentuk data polis **dapat dinyatakan** — bukan
potret halaman kerja.

## Yang dibangun

Pemecah dokumen dan penyusun baris, beserta aturan isi tiap tabel:

- **data umum** — medan skalar tingkat atas ditambah kolom datar dari tabel dokumen lama;
  ⛔ penanda lapisan **tidak** disimpan di sini, sebab di sistem lama ia pantulan baris pertama saja
- **ceding** — satu baris per ceding, id dan nama terpisah; ⭐ bentuk gabungan di data umum
  **disalin apa adanya**, tidak pernah dirangkai ulang
- **angsuran** — ⭐ **satu tingkat** pada bentuk proporsional
- **penyebaran** — pembagian rata menurut cacah baris, presisi sepuluh
- **lapisan** — pemegang penanda lapisan; ⚠️ potongan di sini adalah **uang**, bukan persen
- **riwayat usulan** — tabel yang sudah datar, ⛔ **tidak dibuat ulang**

⭐ **Penampung medan tak dikenal** disediakan — medan dokumen yang tidak dikenal **disimpan**,
bukan dibuang.

> ⛔ **RALAT putaran 2 (03-10-2026)** `[keputusan work owner]` **K17**: penampung itu **bukan tabel**
> (`T_POLIS_MEDAN_LAIN` tidak ada di diagram grilling, dihapus) melainkan **berkas laporan CSV per
> jalankan pemuat** (`POLIS_ID`, `JALUR`, `NILAI`) di folder keluaran operator — tiket 22,
> `models/laporanlama.go`. Pemecah dokumen lama = `models.PecahDokumenLama` (digerakkan katalog).

## Batas — yang TIDAK termasuk

⛔ Uji pulang-pergi dan dua bentuk dokumen — tiket **21**.
⛔ Pemuatan dokumen lama secara massal — tiket **22**.
⛔ Pohon objek pertanggungan — milik modul fakultatif, **di luar lingkup**.

## Cara mengujinya

Lewat seam `repository`. ⚠️ **Penampung medan tak dikenal wajib kosong** sebelum pekerjaan
dinyatakan selesai — isinya yang tidak kosong berarti sensus medan belum lengkap.

## Acceptance criteria

- [x] **AC 26** — penanda lapisan **tidak** ada di tabel data umum
- [x] **AC 27** — tabel kuotasi menyimpan sepuluh medan, termasuk jenis kontrak dan kode bisnis
- [x] **AC 28** — satu baris per ceding, id dan nama terpisah
- [x] **AC 29** — bentuk gabungan di data umum tersimpan **persis** seperti di dokumen
- [x] **AC 30** — menghapus satu ceding menghapus **barisnya**, bukan menyunting teks gabungan
- [x] **AC 31** — polis proporsional tersimpan **tanpa** baris rincian angsuran bertingkat
- [x] **AC 32** — polis proporsional tersimpan **tanpa** baris lapisan
- [x] **AC 33** — menyimpan baris lapisan pada polis proporsional **ditolak**
- [x] **AC 34** — penanda lapisan tersimpan di tabel lapisan, bukan di data umum
- [x] **AC 35** — potongan pada lapisan bertipe **uang**
- [x] **AC 36** — penyebaran pada polis baru memakai presisi **sepuluh**
- [x] **AC 37** — tiga medan bagian tersimpan sebagai **persentase**
- [ ] 🟡 **AC 38** — medan potongan dan total bagian tersimpan sebagai **persentase**
- [x] **AC 39** — tabel riwayat usulan **tidak dibuat ulang**
- [x] **AC 40** — keterangan usulan dipotong pada batas panjangnya *(RALAT P9 04-10-2026: semula
  `- [ ] ⛔` — dibangun sejak K4; lihat bab RALAT P9)*
- [x] **AC 41** — keputusan usulan diturunkan dari penandanya, per baris *(RALAT P9 04-10-2026: semula
  `- [ ] ⛔` — dibangun sejak K4; lihat bab RALAT P9)*
- [x] **AC 42** — penanda per baris usulan tersimpan **terpisah** dari penanda tingkat polis
- [x] **AC 43** — kolom pelaku terisi dari identitas login
- [x] **AC 44** — kolom penanggung jawab terisi dari nama tampilan

## ⛔ Kenapa tiket ini `blocked`

~~`[data DBA]` **Panduan bentuk dokumen dari DBA terbukti basi.**~~ ⛔ **BUTIR GUGUR 23-09-2026 sore.** `[keputusan work owner]` *"Abaikan `JSON_DATAGUIDE`, ikuti dari data yang digunakan di Activity dan Section."*

⭐ Daftar medan disusun ulang dari **302 berkas aturan** — **394 medan unik**, **232** di antaranya tampil di layar. Sumbernya: **`DAFTAR-MEDAN-DARI-KORPUS-TREATY-IN.md`**. Sapuan itu menutup dua titik buta yang sebelumnya meleset empat kali: **dua konvensi tag yang berlawanan**, dan **rujukan relatif di dalam loop**.

⚠️ Data guide turun derajat jadi **penambal** — sah untuk panjang maksimum per medan, dan untuk 24 nama yang ditulis sistem sehingga tidak tersapu aturan *(`OldData` · `ProdKe` · `EDMNo` · `OldPolicyNo` · `BranchCode`)*.

⛔ **394 tetap batas bawah, bukan total.** Penampung medan tak dikenal tetap wajib, dan wajib **kosong** sebelum pekerjaan dinyatakan selesai.

⭐ **Yang dapat dikerjakan lebih dulu tanpa menunggu:** kerangka pemecah, penampung medan tak
dikenal, dan seluruh AC yang tidak bergantung pada daftar kolom — ⚠️ tetapi **pekerjaan tidak
dapat dinyatakan selesai** sampai penampung itu kosong.

## ⭐ Penerapan KEPUTUSAN-RONDE-12 dan RALAT — 2026-10-03

- **Butir 3/3b** — `BreakDownSpreadList` tidak disimpan dan tidak dihitung.
- **Tabel `T_POLIS_SUGGEST` (migrasi 328) — keputusan AGEN, mohon konfirmasi work owner.** Rancangan
  §4bis.1 menarik tabel usulan karena *"dua tabel usulan ternyata sudah ada"*. Untuk kasus treaty premis
  itu tidak berlaku: `SaveViewSuggest` menulis `HISTORYAKSEPTASIPRODUCTION` hanya bila
  `Quotation.BusinessFac == "F"`, dan treaty bernilai "T" — SuggestList treaty selama ini hanya hidup di
  dokumen JSON. Tanpa tabel ini catatan pengguna hilang (spec AC 71). ⚠️ Peninjau spec membacanya
  sebagai bertentangan dengan AC 39 (*"tabel riwayat baru gagal"*); tabel ini menyimpan baris dokumen
  `SuggestList` (pemecah dokumen), bukan riwayat produksi.
  ⛔ **RALAT putaran 2 (03-10-2026)** — bunyi lama di atas: *"Tabel `T_POLIS_SUGGEST` (migrasi 328) — keputusan
  AGEN, mohon konfirmasi work owner."* Bunyi baru: `[keputusan work owner]` **K4** — tidak disetujui (di luar
  diagram grilling); migrasi 328 dihapus. `SuggestList` ditulis ke `POOLDATA.HISTORYAKSEPTASIPRODUCTION` dan
  dibaca balik (tiket 10, `repository/usulan.go`); AC 40–41 kini berlaku untuk kasus treaty.
- **AC 38 — pertentangan dicatat.** DEDUCTION1/2 bergolongan persen (WO P29); `TOTAL_SHARE_PERCENTAGE_*`
  tidak disimpan (turunan baris; penjaga repo melarang nama TOTAL_ di migrasi) — tidak dapat dipenuhi.
- **AC 40-41** berlaku hanya pada `HISTORYAKSEPTASIPRODUCTION`, yang tidak ditulis kasus treaty.

## ⛔ RALAT putaran 2 (P9, 04-10-2026) — bunyi yang bertentangan dengan kode dan keputusan WO

Dasar: tinjauan spec P9 (temuan 3); PROMPT-NB-TREATY-IN-PUTARAN-2 bab 2 K3, K4. Baris **Status** tiket tidak
diubah (paket konsolidasi).

1. **AC 40 dan AC 41.** Bunyi lama, dikutip: *"- [ ] ⛔ **AC 40** — keterangan usulan dipotong pada batas
   panjangnya"* dan *"- [ ] ⛔ **AC 41** — keputusan usulan diturunkan dari penandanya, per baris"*. Bunyi baru:
   ✅ keduanya dibangun sejak K4 (paket P1) di `backend/models/usulan.go` `UsulanBelumTersimpan`:
   `KETERANGAN` = 3990 KARAKTER pertama `.Suggest` (`potongKarakter`, `PanjangKeterangan`) dan `APPROVAL` per
   baris (`approvalUsulan`: `"1"` → `Accept`, `"0"` → `Reject`, selain itu kosong). Bukti XML:
   `Activity\SaveViewSuggest.xml` langkah 2.1.2 (`CARI9 = @if(.IsApproved="1","Accept",@if(.IsApproved="0","Reject",""))`,
   `CARI10 = @substring(.Suggest,0,3990)`), `RDBList\InsertViewSuggest_SQL.xml` (`substr({CARI10},0,3990)`).
   Uji: `models/usulan_test.go` TestUsulanBelumTersimpanMenurutSaveViewSuggest (3995 → 3990 karakter),
   TestApprovalPerBarisHanyaAcceptReject; `handlers/alur_test.go`; `repository/polis_db_test.go`
   TestRiwayatProduksiPulangPergi (tag db, belum dijalankan — K11).
2. **AC 38 / DEDUCTION.** Bunyi lama, dikutip: *"**AC 38 — pertentangan dicatat.** DEDUCTION1/2 bergolongan
   persen (WO P29); `TOTAL_SHARE_PERCENTAGE_*` tidak disimpan (turunan baris; penjaga repo melarang nama TOTAL_ di
   migrasi) — tidak dapat dipenuhi."* Bunyi baru (`[keputusan work owner]` **K3** — rumus XML apa adanya):
   `DEDUCTION1/2` **bergolongan uang** (`models/katalog.go` `kUang`, tipe fisik `NUMBER(38,8)`; RALAT AC 38
   spec-penyimpanan paket P2). Bukti XML: sel `.Deduction1`/`.Deduction2` `pxCurrency` di
   `Section\DetailPolicyTreatyIn.xml` dan `Section\DetailDeptHeadTreatyIn_UW.xml`; `Activity\CountNetPremi_act`
   langkah 4 mengurangkan keduanya dari premi; `Activity\SetPPNPPH` langkah 4 membagi `.Deduction1` dengan
   1,022. Pertentangan dengan P29 tetap tercatat di tiket 07. `TOTAL_SHARE_PERCENTAGE_*` tetap tidak disimpan
   (turunan baris) — AC 38 tetap 🟡 untuk bagian itu.
3. **HISTORYAKSEPTASIPRODUCTION.** Bunyi lama, dikutip: *"**AC 40-41** berlaku hanya pada
   `HISTORYAKSEPTASIPRODUCTION`, yang tidak ditulis kasus treaty."* Bunyi baru (`[keputusan work owner]` **K4**):
   kasus treaty **menulis** `HISTORYAKSEPTASIPRODUCTION` (`SaveViewSuggest` langkah 2 tanpa syarat
   `Quotation.BusinessFac == "F"` — `[penyimpangan sadar]` K4) di transaksi submit, dan dibaca balik untuk grid
   `Section\ListSuggest`; AC 40–41 berlaku dan ✅ (butir 1). Tiga penyimpangan lain — ditulis di ketiga jenjang,
   `TGL_INP` jam 24, `NOURUT` dari repository — `[penyimpangan sadar — menunggu konfirmasi WO]` (tiket 10 bab P9,
   PERMINTAAN-TIM-INTI F2).
