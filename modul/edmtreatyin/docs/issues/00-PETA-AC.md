# 00 · PETA AC — penyimpanan relasional EDM Treaty In

> ## ⛔ KOREKSI 06-10-2026 — rekonsiliasi dengan keadaan nyata (log: `../KOREKSI-DOKUMEN-2026-10-06.md`)
>
> | # | Bunyi lama (dikutip) | Bunyi baru | Bukti |
> | ---: | --- | --- | --- |
> | 1 | tiket 11 *"⛔ tetap blocked"*, tiket 12 *"⛔ blocked · tiket baru"*, *"⇒ 10 siap · 2 tertahan"* | ⭐ **11 dan 12 `ready-for-agent`** ⇒ **12 siap · 0 tertahan** | `11-…md` baris 3–14, `12-…md` baris 3–14 (penahan gugur 23-09 sore); `../spec-penyimpanan-relasional.md` bab *KEPUTUSAN 23-09-2026 sore* |
> | 2 | `ID-27b` *"tidak tersentuh satu tiket pun"* · `ID-27c` *"tidak terbawa ke tiket mana pun"* | ⭐ keduanya **tiket 12** (bagian 1 = `ID-27b`, bagian 2 = `ID-27c`); **AC 57** (isi = `REMARK` 128, `ID-27b`) dipindah **11 → 12** | `12-…md` baris 18, 37–53; spec-penyimpanan AC 57 |
> | 3 | tiket 08 *"AC 35–38 · AC 56 (5 AC)"* | **AC 35–38 (4 AC)** — AC 56 ditarik dari lingkup | spec-penyimpanan AC 56; `08-…md` baris 60 |
> | 4 | *"Sisanya bergerbang NB 18 atau NB 19, yang keduanya menunggu DBA"* | gerbang NB **16 · 17 · 18 · 19 · 20 · 22 berstatus selesai** (sisa penahan pihak luar K11) | `modul/nbtreatyin/docs/issues/16…` baris 3, `17…` 3, `18…` 20, `19…` 14, `20…` 3, `22…` 14 |
> | 5 | *"Sepuluh tabel dasar dibuat tiket NB"* | tabel dasar nyata: `T_WORK_POLIS` (premiumlistlife) + `T_GENERAL_POLIS_TREATY` + delapan anak `T_POLIS_*` migrasi NB **321–328** (`T_POLIS_SURVEY` 328 ikut) | `modul/nbtreatyin/backend/migrations/320…328` |
> | 6 | peringatan *"tumpang tindih … wajib diputuskan"* | ✅ sudah diputuskan: NB 24–28 bertanda **DIGANTIKAN**. ⚠️ Kelimanya masih menunjuk jalur lama `.scratch\edm-treaty-in\issues\…` — **tidak disunting** (berkas NB); jalur kini `modul/edmtreatyin/docs/issues/` | `modul/nbtreatyin/docs/issues/24…28` baris 3–6, 32 |
> | 7 | `KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md` (tanpa jalur) | `modul/nbtreatyin/docs/KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md` | berkas ada di folder itu |
>
> ⛔ Nol butir `[terbuka]` ditutup oleh koreksi ini.

---

> ## ⭐⭐ DIPERBARUI 23 September 2026 sore — tujuh butir tertahan ditutup
>
> Ronde keputusan dua belas butir menutup **seluruh** penahan `[work owner]`. Enam tiket di berkas
> ini berubah label.
>
> | Tiket | Semula | Sekarang | Sebab |
> | ---: | --- | --- | --- |
> | **02** | ⛔ blocked | ⭐ `ready-for-agent` | endorsemen **tanpa batas**, kolom generasi dilebarkan |
> | **05** | ⛔ blocked | ⭐ `ready-for-agent` | boleh endorse sesudah batal, cukup **peringatan** |
> | **07** | ⛔ blocked | ⭐ `ready-for-agent` | **kedua** anak tabel proyeksi dibuat |
> | **08** | ⛔ blocked | ⭐ `ready-for-agent` | ⚠️ **lingkup menyusut** — lihat di bawah |
> | **09** | ⛔ blocked | ⭐ `ready-for-agent` | dokumen lama dipindah **seluruhnya** |
> | **10** | ⛔ blocked | ⭐ `ready-for-agent` | idem |
> | **11** | ⛔ blocked | ~~⛔ **tetap blocked**~~ ⭐ `ready-for-agent` *(koreksi 06-10)* | ~~`[data DBA]` presisi fisik kolom uang~~ gugur 23-09 sore — `[keputusan work owner]` *"Selesaikan, jangan jadi permasalahan."* |
> | **12** | — | ~~⛔ **blocked** · **tiket baru**~~ ⭐ `ready-for-agent` · tiket baru *(koreksi 06-10)* | ~~`[data DBA]` panduan bentuk dokumen basi~~ gugur 23-09 sore — `[keputusan work owner]` *"Abaikan `JSON_DATAGUIDE`…"* |
>
> ⇒ ~~**10 siap · 2 tertahan**, keduanya menunggu **DBA**, bukan work owner.~~ ⭐ **12 siap · 0 tertahan** *(koreksi 06-10, butir 1)*.
>
> ### ⛔ Tabel sebaran tambahan DIBATALKAN
>
> `~~T_POLIS_BREAKDOWN_SPREAD~~` **tidak ada.** Keempat medannya turunan: `TREATY_TYPE` dan
> `SHARE_PERCENTAGE` dari master `POOLDATA.PROPORTIONALARRG`, dua medan uangnya dihitung dari
> `TotalPremium` dan `TotalClaim` yang **sudah** tersimpan pada polis.
>
> Akibatnya: ~~**AC 56 diganti isinya**~~ ⛔ ⭐ **AC 56 DITARIK DARI LINGKUP.** Sebaran tambahan
> tidak dimigrasi sama sekali — nol tabel, **nol perhitungan**, nol layar. Cacah AC berlaku **58 → 57**.
> `ID-6b` **dicabut**, sehingga cacah ID turun **50 → 49**.
>
> Cacah tabel kembali ke semula: **EDM Prop 10 · EDM NonProp 14**.
>
> ### ⭐ Gerbang NB tidak berubah
>
> Label `ready-for-agent` **bukan** berarti dapat dimulai. Setiap tiket di sini tetap menunggu
> gerbang NB-nya. Lima tiket punya gerbang yang sudah siap dikerjakan — **02 · 03 · 04 · 06 · 07**
> *(gerbang NB 16, 17, 20)*. ~~Sisanya bergerbang **NB 18** atau **NB 19**, yang keduanya menunggu DBA.~~
> ⭐ *(koreksi 06-10, butir 4)* Gerbang NB 16–20 dan 22 kini **selesai**; sisa penahannya hanya K11
> (skema uji Oracle) — `modul/nbtreatyin/docs/issues/16…22` baris *Status*.
>
> Rinciannya di `modul/nbtreatyin/docs/KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`.

---


> ## ⭐ TAMBAHAN 23 September 2026 — kesebelas tiket ini **menggantikan** tiket NB 24–28
>
> Ronde tiket sebelumnya menutup **dua** spec sekaligus, sehingga tiket endorsemen jatuh ke folder
> polis baru sebagai nomor **24–28**. Kelimanya sekarang **ditandai digantikan** — tidak dihapus,
> tidak dikerjakan. Penggantinya ada di berkas ini.
>
> | Tiket lama | Digantikan |
> | --- | --- |
> | NB 24 | **01** · **02** · **03** · **05** |
> | NB 25 | **06** |
> | NB 26 | **07** |
> | NB 27 | **08** |
> | NB 28 | **09** · **10** · **11** |
>
> ⭐ Alasannya satu: **tiket di berkas ini menyebut gate NB-nya**, tiket 24–28 tidak. Tanpa gate
> itu pekerjaan endorsemen bisa dimulai sebelum tabelnya ada.
>
> ⚠️ Tumpang tindih itu **kekeliruan penyusun brief**, bukan pelaksana ronde mana pun.
>
> ---
>
> ### ⛔ Dua butir yang belum tertutup — hasil audit sesudah ronde, bukan koreksi tiket
>
> Spec memuat **50 ID unik**: 47 bernomor murni **+ `ID-6b` · `ID-27b` · `ID-27c`**. Ketiganya
> disisipkan sesudah spec pertama ditulis, jadi wajar tak terbaca ronde tiket. Cacah **47** di
> kepala berkas ini benar untuk yang bernomor murni.
>
> | ID | Isi | Keadaan |
> | --- | --- | --- |
> | `ID-6b` | tabel sebaran tambahan | ✅ tertutup **isi**nya oleh tiket **08**, label ID-nya saja yang tak dikutip |
> | `ID-27b` | medan `REMARK` *(panjang 128)* ikut dimigrasi | ~~⛔ **tidak tersentuh satu tiket pun**~~ ✅ **tiket 12 bagian 1** + **AC 57** *(koreksi 06-10)* |
> | `ID-27c` | data guide berkedudukan pelengkap — satu dokumen memuat 95 jalur yang tak ada di dalamnya | ~~⛔ butir tertahan `[data DBA]` ini **tidak terbawa** ke tiket mana pun~~ ✅ **tiket 12 bagian 2** — tanpa nomor AC di spec; penahannya gugur 23-09 sore *(koreksi 06-10)* |
>
> ⭐ Akibat praktisnya: tiket yang menetapkan **daftar kolom** berjalan seolah data guide lengkap.
> Ia tidak lengkap. ~~Butir itu perlu ditambahkan sebagai `blocked` sebelum daftar kolom dikunci.~~
> ✅ Sudah ditambahkan sebagai **tiket 12**, lalu penahannya gugur 23-09 sore (`12-…md` baris 3–9).
>
> ⚠️ Catatan ini **tidak menyunting satu pun tiket**. Penutupannya keputusan work owner.

---


> **Sumber:** `spec-penyimpanan-relasional.md` — **58 AC · ~~47~~ ~~50~~ ⭐ 49 ID** *(47 bernomor murni + `ID-27b` · `ID-27c`; ~~`ID-6b`~~ ⛔ dicabut 23-09 sore)*.
> ⛔ Spec **tidak disunting** ronde ini. ⛔ Tiket NB **tidak disentuh**.

| Ukuran | Nilai |
| --- | ---: |
| AC pada spec | ~~**58**~~ ⭐ **57** *(56 ditarik)* |
| ⭐ tercakup tiket | ⭐ **57** |
| ⛔ yatim | ✅ **0** |
| ganda | ✅ **0** |
| tiket | ~~**11** — **4** siap · **7** tertahan~~ ~~**12** — **10** siap · **2** tertahan~~ ⭐ **12** — **12** siap · **0** tertahan *(koreksi 06-10, butir 1)* |

---

## ⚠️ EDM tidak membuat satu pun tabel dasar

~~Sepuluh tabel dasar dibuat **tiket NB**.~~ ⭐ *(koreksi 06-10, butir 5)* Tabel dasar dibuat **migrasi NB**:
`T_GENERAL_POLIS_TREATY` (320) + delapan anak `T_POLIS_*` (321–328, termasuk `T_POLIS_SURVEY`), di atas
`T_WORK_POLIS` milik premiumlistlife. EDM hanya menambah empat tabel **proyeksi** (migrasi 360–363,
`../STRUKTUR-TABEL-EDM-TREATY-IN.md`). Setiap tiket di bawah menyebut tiket NB yang
menggatenya — dikerjakan **sesudahnya**, bukan bersamaan.

| # | Tiket | Status | Blocked by | Gate NB | Menutup |
| ---: | --- | --- | --- | --- | --- |
| **01** | Kolom khas endorsemen pada tabel yang sudah ada | ⭐ `ready-for-agent` | — *(dapat mulai sesudah tiket NB)* | **16** *(kerangka penyimpanan)* · **19** *(pemecah dokumen)* | AC **32–34** · AC **54–55** *(5 AC)* |
| **02** | Rantai generasi dan larangan percabangan | ⭐ `ready-for-agent` | ~~`[work owner] batas berapa kali satu polis boleh di-endorse~~ ✅ **gugur 23-09 sore** — batas teknis 99, batas dagang belum ditetapkan | **16** *(kerangka penyimpanan dan kunci generasi)* | AC **1–6** *(6 AC)* |
| **03** | Keutuhan nomor urut antar generasi | ⭐ `ready-for-agent` | **02** | **17** *(nomor urut baris anak)* | AC **7–9** *(3 AC)* |
| **04** | Nomor urut di endorsemen — tanpa penghapusan | ⭐ `ready-for-agent` | **03** | **17** *(nomor urut baris anak)* | AC **10–13** *(4 AC)* |
| **05** | Pembatalan sebagai generasi bernilai nol | ⭐ `ready-for-agent` | ~~`[work owner] sesudah dibatalkan, polis masih boleh di-endorse lagi atau tidak~~ ✅ **gugur 23-09 sore** | **19** *(pemecah dokumen — kolom jenis berkas)* | AC **14–15** *(2 AC)* |
| **06** | Perhitungan selisih di lapisan layanan | ⭐ `ready-for-agent` | **02** · **03** | **16** *(kerangka penyimpanan)* | AC **16–22** *(7 AC)* |
| **07** | Proyeksi selisih yang dapat dibaca langsung | ⭐ `ready-for-agent` | **06** · ~~`[work owner] anak proyeksi dibutuhkan pembaca SQL atau tidak~~ ✅ **gugur 23-09 sore** | **20** *(transaksi tunggal dan skema eksplisit)* | AC **23–31** *(9 AC)* |
| **08** | Perbedaan perhitungan sebaran dan rincian angsuran | ⭐ `ready-for-agent` | ~~`[work owner] beda dagang tabel sebaran tambahan~~ ✅ **gugur 23-09 sore** dari sebaran risiko belum dijelaskan | **19** *(pemecah dokumen menjadi baris)* | AC **35–38** ~~· AC **56** *(5 AC)*~~ *(4 AC — AC 56 ditarik; koreksi 06-10)* |
| **09** | Dua penanda migrasi | ⭐ `ready-for-agent` | **04** · **07** · ~~`[work owner] lingkup pemindahan dokumen lama~~ ✅ **gugur 23-09 sore** | **22** *(pemuat dokumen lama)* | AC **39–43** *(5 AC)* |
| **10** | Pemuat migrasi endorsemen | ⭐ `ready-for-agent` | **09** · ~~`[work owner] lingkup pemindahan dokumen lama~~ ✅ **gugur 23-09 sore** | **22** *(pemuat dokumen lama)* | AC **44** *(1 AC)* |
| **11** | Kepatuhan lapisan dan tipe kolom | ⭐ `ready-for-agent` | ~~`[data DBA] presisi fisik kolom uang~~ ✅ **gugur 23-09 sore** — dua belas digit di depan koma belum diuji terhadap nilai terbesar | **18** *(tipe kolom dan presisi uang)* · **20** *(transaksi tunggal dan skema eksplisit)* | AC **45–53** · AC ~~**57–58** *(11 AC)*~~ **58** *(10 AC — AC 57 ke tiket 12; koreksi 06-10)* |
| **12** | Medan catatan dan penguncian daftar kolom | ⭐ `ready-for-agent` | ~~`[data DBA]` panduan bentuk dokumen basi~~ ✅ **gugur 23-09 sore** | **19** *(pemecah dokumen menjadi baris)* | AC **57** · `ID-27b` · `ID-27c` *(1 AC; baris ditambah koreksi 06-10)* |

---

## AC → tiket

| Tiket | AC yang ditutupnya |
| ---: | --- |
| **01** | 32–34 · 54–55 |
| **02** | 1–6 |
| **03** | 7–9 |
| **04** | 10–13 |
| **05** | 14–15 |
| **06** | 16–22 |
| **07** | 23–31 |
| **08** | 35–38 · ~~56~~ *(ditarik)* |
| **09** | 39–43 |
| **10** | 44 |
| **11** | 45–53 · ~~57–58~~ 58 |
| **12** | 57 *(+ `ID-27c`, tanpa AC)* — *baris ditambah koreksi 06-10* |

---

## ⛔ ~~Tujuh tiket `blocked`~~, dan butir yang menahannya — ⭐ *(koreksi 06-10)* **nol tiket `blocked`**

| Tiket | Tertahan | Pemilik |
| ---: | --- | --- |
| ~~**02**~~ ✅ | ~~batas berapa kali satu polis boleh di-endorse — batas teknis 99~~ **gugur 23-09 sore** | ~~`[work owner]`~~ |
| ~~**05**~~ ✅ | ~~sesudah dibatalkan, polis masih boleh di-endorse lagi atau tidak~~ **gugur 23-09 sore** | ~~`[work owner]`~~ |
| ~~**07**~~ ✅ | ~~anak tabel proyeksi dibutuhkan pembaca SQL atau tidak~~ **gugur 23-09 sore** | ~~`[work owner]`~~ |
| ~~**08**~~ ✅ | ~~beda dagang tabel sebaran tambahan dari sebaran risiko~~ **gugur 23-09 sore** | ~~`[work owner]`~~ |
| ~~**09**~~ ✅ | ~~lingkup pemindahan dokumen lama~~ **gugur 23-09 sore** | ~~`[work owner]`~~ |
| ~~**10**~~ ✅ | ~~idem~~ **gugur 23-09 sore** | ~~`[work owner]`~~ |
| ~~**11**~~ ✅ | ~~presisi fisik — dua belas digit di depan koma belum diuji~~ **gugur 23-09 sore** | ~~`[data DBA]`~~ |

~~⭐ **Empat siap dikerjakan** — **01**, **03**, **04**, **06** — begitu tiket NB yang menggatenya
selesai.~~ ⭐ *(koreksi 06-10)* **Dua belas siap**; tiket NB penggatenya selesai (lihat blok KOREKSI butir 4).

⚠️ **Tiket `blocked` tetap ditulis lengkap.** ⛔ Jangan ditandai siap sebelum butirnya dijawab.

---

## ⛔⛔ Peringatan: tiket ini TUMPANG TINDIH dengan tiket NB 24–28

> ✅ *(koreksi 06-10, butir 6)* **Sudah diputuskan** — NB 24–28 bertanda **DIGANTIKAN**
> (`modul/nbtreatyin/docs/issues/24…28` baris 3–6, 32). Bunyi di bawah dipertahankan sebagai sejarah.

~~⚠️ **Wajib diputuskan sebelum pekerjaan dimulai.**~~

Ronde tiket sebelumnya menulis **tiket NB 24–28**, yang menutup **AC 1–58 spec yang sama ini**, ke
dalam folder tiket NB. Ronde ini menulis **sebelas tiket** yang menutup AC yang sama, ke dalam
folder ini.

⇒ ⛔ **Ke-58 AC kini tercakup dua kali, di dua tempat.**

| | folder tiket NB | folder ini |
| --- | --- | --- |
| tiket | **24–28** *(5)* | **01–11** *(11)* |
| AC EDM tercakup | 58 | 58 |
| kehalusan | lima potongan besar | sebelas potongan |
| menyebut gate NB | ⛔ tidak | ⭐ **ya, 11 dari 11** |

⭐ **Yang disarankan:** pertahankan **sebelas tiket di folder ini** — lebih halus, dan tiap tiket
menyatakan ketergantungan pada tiket NB yang membuat tabelnya. Lalu **tandai 24–28 sebagai
digantikan**.

⛔ **Tidak dikerjakan ronde ini** — brief melarang menyentuh tiket NB. Keputusannya milik
`[work owner]`.

---

## ⚠️ Dua penyimpangan dari pembagian yang disarankan brief

Brief menyarankan **sepuluh** potongan. Ditulis **sebelas**, dan sebabnya:

1. ⭐ **Tiket 11 ditambahkan** — AC **45–53** *(kepatuhan lapisan)* dan **57–58** *(tipe kolom)*
   **tidak tercakup** pembagian brief. Tanpa tiket itu, **sebelas AC menjadi yatim**.
   ⭐ *(koreksi 06-10, butir 2)* AC 57 berisi medan `REMARK` (`ID-27b`), bukan tipe kolom umum →
   kini milik **tiket 12**; tiket 11 menutup **10 AC**.
2. ⭐ **Tiga AC ditempatkan ulang** — AC **54–55** *(bentuk tabel tidak berubah)* masuk tiket **01**,
   sebab di sanalah penjaganya berada; ~~AC **56** *(tabel sebaran tambahan)* masuk tiket **08**~~ ⛔ **AC 56 ditarik dari lingkup 23-09 sore**,
   sebab ia soal sebaran, bukan soal bentuk.

⚠️ **Dan dua angka brief yang tidak cocok dengan yang terukur:**

| | Brief | Terukur |
| --- | ---: | ---: |
| Implementation decision | 50 | **47** |
| rentang ID untuk kolom khas endorsemen | ID-15..ID-19 | ⛔ **ID-20 · ID-21 · ID-22** — ID-15..ID-19 sebenarnya berisi aturan keutuhan, nomor urut, dan pembatalan |

⭐ Peta ini memakai **ID yang sebenarnya ada di spec**.

---

## Disiplin ronde ini

| ⛔ Larangan | Dipatuhi? | Bukti |
| --- | :---: | --- |
| nol nama tabel di **judul** | ✅ | diperiksa pola — **0** dari 11 |
| spec tidak disunting | ✅ | **0** tulis |
| tiket NB tidak disentuh | ✅ | **0** tulis |
| nol `CREATE TABLE` di tiket | ✅ | **0** — presisi fisik dicocokkan DBA **di dalam** tiket |
| nol nama orang · nomor polis harfiah · cuplikan produksi | ✅ | **0** · **0** · **0** |
| jangan menutup butir `[terbuka]` | ✅ | **0** ditutup; tujuh butir disebut sebagai penahan |
| setiap tiket menyebut ketergantungan NB | ✅ | **11 / 11** |

---

## TELEMETRI EKSEKUSI

### ⛔ Yang TIDAK diukur

| Ukuran | Keadaan |
| --- | --- |
| pengukuran **dari luar** | ⛔ **tidak dilakukan** — perintah itu memulai sesi **baru dan terpisah**, yang ongkosnya bukan ongkos ronde ini |
| durasi · biaya · panggilan alat | ⛔ **tidak diukur** — baseline-nya tidak diambil |

⛔ **Tidak satu pun ditaksir lalu disajikan sebagai angka terukur.**

### ⭐ Yang TERUKUR — selisih terhadap baseline sesi

| Ukuran | Nilai |
| --- | ---: |
| panggilan model | **15** |
| token keluar | **89.716** |
| token cache ditulis | **85.553** |
| token cache dibaca | **12.081.207** |
| tiket ditulis | **11** |
| AC dipetakan | **58** |
| berkas spec disunting | **0** |

⚠️ Angka token sejati tidak terlihat dari dalam sesi; yang di atas **pengurangan terhadap baseline**
yang diambil pada awal ronde.

---

*Disusun 23 September 2026 dari `spec-penyimpanan-relasional.md` EDM Treaty In.*