# TAMBAHAN TIKET — hasil Grilling Ronde 2

**Tanggal:** 2026-09-19 · **Modul:** Master Contract Retro Life
**Menunjuk balik ke:** `issues/01-…` sampai `issues/12-…` *(dua belas tiket, tidak disunting)*

> ⛔ **NOL berkas tiket lama disunting.** ⛔ `CREATE TABLE` NOL · DDL NOL · kode NOL.
> ⛔ Nol butir `[terbuka]` / OQ dinyatakan tertutup.
> Bentuk tiket baru mengikuti `issues/05-reinsurer-share-dan-total.md` — tidak dikarang baru.

---

# BAGIAN (a) — TAMBAHAN untuk tiket 01–12 yang sudah ada

## Tiket **03** — Gerbang tahun lewat nilai tanggal

| | |
| --- | --- |
| **Yang ditambahkan** | ⭐ **Gerbang konsistensi tahun di Pega TIDAK PERNAH BERJALAN.** Ketiga baris syaratnya di `Activity/SaveSecurityLife_Act.xml` langkah **4 · 5 · 6** dan tiga kembarannya di `Activity/SaveSecurityReinsurerLife_Act.xml` langkah **4 · 5 · 6** berflag prakondisi **`false`** — tersimpan tetapi dimatikan. Langkahnya berjalan **tanpa saringan**. |
| **Bukti** | kedua berkas di `Master Contract Retro Life/Activity/`, langkah 4 · 5 · 6, medan `pyStepsPreCondition` |
| **Menambah / meralat** | ⭐ **MERALAT dasar faktualnya** — dari *"Pega memeriksanya dengan cara keliru"* menjadi *"Pega tidak memeriksanya sama sekali"*. ⛔ Keputusan tiketnya **tidak berubah**: gerbangnya tetap **dibangun** memakai nilai tanggal |
| **Akibat untuk pengujian** | Test tidak boleh berbunyi *"meniru perilaku Pega"* — tidak ada perilaku untuk ditiru. Test menguji **aturan baru** |

## Tiket **05** — Reinsurer, share, dan total

| | |
| --- | --- |
| **Yang ditambahkan** | ⭐ **Nol penjaga total share di kedua sisi.** Pega menjumlah lalu menampilkan *(`CountingPercentShare_Act` langkah **3.1** dan **4**)* — nol perbandingan terhadap 100. `[data DBA]` procedure juga **nol validasi bisnis**. ⚠️ Hasil penjumlahan ditulis ke medan bernama **`STDRATING`** |
| **Bukti** | `Activity/CountingPercentShare_Act.xml` langkah 3.1 · 4 · `procedure-bodies-from-dba.md` |
| **Menambah / meralat** | **MENAMBAH** — menguatkan keputusan Q4 *(tampilkan, jangan blokir)* dengan bukti dua sisi |
| **Peringatan migrasi** | Bila kolom bernama *rating* ditemui saat migrasi, **periksa isinya** — bisa jadi ia total share |

## Tiket **07** — Business dan tampilan rate

| | |
| --- | --- |
| **Yang ditambahkan** | ⭐ ⚠️ **`RIRATE` disimpan sebagai TEKS sepanjang 1000 karakter**, bukan angka. `RIRATEID` berdampingan dengannya. Keduanya punya **4 penulis** *(`NewInputBusinessLife_Act` · `SetBusinessListLife_Act` · `SetParamRate`)* dan **8–9 pembaca**, jadi keduanya **dipakai**, bukan sisa |
| **Bukti** | `ddl-tables-from-dba.md` temuan 3 · sensus 66 berkas di `grilling-ronde-2.md` §B3 |
| **Menambah / meralat** | ⭐ **MERALAT** — ronde 1 menulis `RIRATE` *"belum saya telusur"*; kini terbaca tipenya, tetapi **isinya belum** |
| **⛔ MEMBLOKIR** | ⭐ **Tiket ini MACET** sampai **Pertanyaan A** dijawab: bentuk kolom, validasi, dan migrasi nilainya ketiganya bergantung apakah `RIRATE` satu angka atau teks berstruktur |

## Tiket **08** — Terapkan ke semua dengan pratinjau

| | |
| --- | --- |
| **Yang ditambahkan** | ⚠️ **Gerbangnya terbelah.** `Activity/SaveBusinessToAllLife_Act.xml` langkah **3** berflag **`false`** — **mati**; tetapi langkah **3.1** dan **3.2** berflag **`true`** — **berlaku**. Syaratnya sama persis: `.REINSTYPEID == Param.REINSTYPEID` |
| **Bukti** | `Activity/SaveBusinessToAllLife_Act.xml` langkah 3 · 3.1 · 3.2 |
| **Menambah / meralat** | **MENAMBAH** — ronde 1 mengutip gerbang ini tanpa membaca arahnya |
| **Akibat** | Penyaringan *"hanya baris berjenis sama"* **berjalan di tingkat anak**, bukan di tingkat induk. Pratinjau harus menunjukkan baris mana yang benar-benar tersentuh |

## Tiket **09** — Kaskade hapus dengan popup konfirmasi

| | |
| --- | --- |
| **Yang ditambahkan** | ⭐ `[data DBA]` **Basis data kini mengaskade sendiri** — **empat** kunci tamu bermode `ON DELETE CASCADE`: kontrak→tahun · reinsurer→kontrak · security→reinsurer · business→kontrak. Kelima tabel juga sudah punya kunci utama |
| **Bukti** | `ddl-tables-from-dba.md`, bagian *KEADAAN FINAL* |
| **Menambah / meralat** | **MENAMBAH** — *"nol kaskade"* **tetap benar untuk Pega**; yang berubah adalah **siapa** yang mengaskade |
| **Akibat untuk pengujian** | ⚠️ Test kaskade **tidak boleh** mengandaikan aplikasi yang menghapus anak. Test harus membuktikan **hasil akhirnya**, dan memastikan popup konfirmasi muncul **sebelum** basis data bertindak |

## Tiket **10** — Penegakan `HASIL1` di lima jalur simpan

| | |
| --- | --- |
| **Yang ditambahkan** | ⭐ **Tujuh berkas dari 66 menyebut `HASIL1`**, bukan satu — kelima rule SQL simpan **mendeklarasikannya**. ⭐ Tetapi **nol activity membacanya**. Satu-satunya activity yang menyentuh penampung serupa memakai **`HASIL12`**, nama yang **berbeda**. `[data DBA]` Saat gagal, procedure menaruh pesan galat berikut teks galat basis data ke penampung itu, **menjalankan pembatalan**, lalu selesai — sehingga **seluruh perubahan dibatalkan tanpa ada yang melihatnya** |
| **Bukti** | lima berkas di `RDBList/` · `Activity/DeleteRowBusiness.xml` · `procedure-bodies-from-dba.md` |
| **Menambah / meralat** | ⭐ **MERALAT** *(jumlah berkas: 1 → 7)* **dan MENAMBAH** *(apa yang hilang)* |
| **Tambahan bentuk** | Pesan galat berisi **HTML**; sistem baru **tidak boleh meneruskannya mentah** |
| **`[terbuka]` yang menyentuhnya** | **Pertanyaan D** — apakah `HASIL12` salah ketik atau penampung lain |

## Tiket **11** — Laporan kontrak dengan total share tidak 100%

| | |
| --- | --- |
| **Yang ditambahkan** | ⚠️ **Sebelas dari dua belas laporan modul ini belum pernah dibaca**, termasuk `BrowseRateLife_RD` dan `BrowseRateLifeSummary`. ⛔ Kelas sumber, kolom, dan saringannya **tidak berhasil dibaca** di ronde ini — disisir tiga medan, ketiganya nol |
| **Bukti** | sensus `ReportDefinition/` di `grilling-ronde-2.md` §B5 |
| **Menambah / meralat** | **MENAMBAH** — butir `[terbuka]`, ⛔ bukan pembatalan |
| **Akibat** | Bila salah satu laporan itu **sudah** menghitung total share, tiket ini berubah dari *"bangun laporan"* menjadi *"tiru laporan"* |

## Tiket **12** — Migrasi skema lima tabel

| | |
| --- | --- |
| **Yang ditambahkan** | ⭐ `[data DBA]` **Kunci utama kini ada di kelima tabel** · **empat kunci tamu bermode kaskade** sudah dipasang · **wajib-isi tetap nol** di basis data. ⚠️ Kolom angka **tidak dibatasi digitnya** oleh basis data — batas ketelitian sepenuhnya keputusan aplikasi. ⭐ Satu kolom, `RIRATE`, bertipe **teks** |
| **Bukti** | `ddl-tables-from-dba.md` temuan 1 · 2 · 3 · 4 · 6, dan bagian *KEADAAN FINAL* |
| **Menambah / meralat** | **MENAMBAH** |
| **Peringatan** | ⚠️ Kunci tamu kaskade **mungkin dipasang setelah data lama dibuat** — data yatim lama bisa masih ada. Migrasi wajib **memverifikasi keadaan nyata**, bukan berasumsi |

---

# BAGIAN (b) — TIKET BARU

⭐ **Dua tiket baru**, diberi nomor mulai **13** supaya tidak bertabrakan.

---

# 13: Validasi simpan yang di Pega tertulis tetapi tidak berjalan

**Status:** ⚠️ **menunggu keputusan work owner** *(Pertanyaan B)* — bukan `ready-for-agent`

**Blocked by:** 02 (kontrak lahir), 05 (reinsurer lahir), 06 (security reinsurer lahir)

## Hasil & nilai pengguna

Sebagai **admin master**, saya diberi tahu **saat menyimpan** bila nama reinsurer atau tanggal mulai
belum saya isi — bukan menemukan barisnya kosong berhari-hari kemudian saat laporan dibuat.

## Area codebase

| Lapisan | Isi |
| --- | --- |
| `internal/services` | Validasi wajib-isi sebelum memanggil jalur simpan |
| `internal/handlers` | Pesan galat per medan pada respons simpan |
| `frontend/` | Penanda medan wajib; galat ditampilkan di medannya |

## Perilaku Pega yang ditiru, berikut buktinya

⭐ **Yang ditiru adalah ketiadaannya** — dan itu disengaja untuk dicatat, bukan untuk diteruskan.

| Rule | Langkah | Syarat yang tertulis | Flag |
| --- | --- | --- | --- |
| `Activity/SaveSecurityLife_Act.xml` | **4 · 5 · 6** | nama reinsurer kosong · tanggal mulai kosong · konsistensi tahun | ⛔ **`false` — mati** |
| `Activity/SaveSecurityReinsurerLife_Act.xml` | **4 · 5 · 6** | tiga syarat yang **sama persis** | ⛔ **`false` — mati** |
| `Activity/SaveBusinessLife_Act.xml` | **6** | tiga syarat yang **sama persis** | ⛔ **`false` — mati** |
| `Activity/SaveBusinessToAllLife_Act.xml` | **3** | jenis reasuransi sama | ⛔ **`false` — mati** |

⭐ **21 dari 42 baris syarat di jalur simpan modul ini bergerbang mati.**
`[data DBA]` Dan kelima procedure penulis **nol validasi bisnis**. ⭐ **Tidak ada penjaga di kedua
sisi.**

## Keputusan work owner yang mengikat

⚠️ **BELUM ADA.** Tiket ini **menunggu Pertanyaan B** di `grilling-ronde-2.md` §G:

> *Apakah sistem baru **membangun** pemeriksaan itu, atau **meniru keadaan sekarang** yang tanpa
> pemeriksaan?*

⭐ **Rekomendasi asisten — bukan keputusan:** **bangun**, dengan dua syarat: *(i)* dinyatakan
sebagai **penyimpangan sadar** di spec, dan *(ii)* migrasi data lama **tidak menolak** baris yang
sudah terlanjur kosong. Alasannya: pemeriksaan itu **ditulis sendiri oleh pengembangnya** lalu
dimatikan — bentuk yang menyerupai keputusan sementara, bukan keputusan rancangan.

## Apa yang harus diuji

- [ ] simpan reinsurer **tanpa nama** → **ditolak**, dan pesannya menyebut medan mana
- [ ] simpan **tanpa tanggal mulai** → **ditolak**
- [ ] simpan dengan tahun kontrak **tidak sesuai tahun treaty** → **ditolak**
- [ ] baris **lama** yang sudah terlanjur kosong → **masih dapat dibuka dan diperbaiki**, tidak
      dikunci oleh validasi baru
- [ ] bila keputusan jatuh ke **meniru**: seluruh test di atas **dibalik** menjadi *"diterima"*

## Butir `[terbuka]` yang menyentuhnya

- ⭐ **Pertanyaan B** — dibangun atau ditiru *(memblokir tiket ini)*
- ⚠️ **H2 `grilling-ronde-2.md`** — arti flag `false` **belum diuji di layar Pega** untuk modul ini.
  Bila artinya ternyata lain, **21 baris berbalik menjadi hidup** dan tiket ini **gugur seluruhnya**

---

# 14: Empat layar dan sebelas laporan yang belum pernah dibaca

**Status:** ⚠️ **belum dapat dikerjakan** — pekerjaan pembacaan korpus, bukan pembangunan

**Blocked by:** tidak ada

## Hasil & nilai pengguna

Sebagai **pembangun**, saya tidak menemukan di tengah pekerjaan bahwa layar yang saya bangun
ternyata **sudah ada bentuknya** di sistem lama dan saya membangunnya berbeda tanpa alasan.

## Area codebase

⛔ **Belum dapat ditentukan** — tiket ini menghasilkan **bacaan**, bukan kode.

## Perilaku Pega yang ditiru, berikut buktinya

⭐ **Belum terbaca.** Yang terbaca hanyalah **keberadaan dan ukurannya**:

| Section belum dibaca | Byte | Disambung oleh |
| --- | --- | --- |
| `InputRetroLimitReinsurers` | **541 845** | `Harness/InboxRetroLimitReinsurers` |
| `InputSecurityLifeReinsurers` | **496 437** | `Harness/InboxRetroLifeReinsurersList` |
| `InputBusinessLifeReinsurers` | **449 805** | `Harness/InboxBusinessLifeReinsurers` |
| `InputDtlRetrocessionLife` | **203 911** | ⛔ belum terbaca |

⭐ **Total 1,69 MB — lebih besar dari keempat layar yang SUDAH dibaca (1,10 MB).**

**Sebelas laporan belum dibaca:** `BrowseBusinessLife_RD` · `BrowseCedingCoLife_RD` ·
`BrowseDetailTreatyReisurerLife_RD` · **`BrowseRateLifeSummary`** · **`BrowseRateLife_RD`** ·
`BrowseRetrocessionLife_RD` · `BrowseSecurityReinsurer_Life_RD` · `BrowseTreatyBusiness_Life_RD` ·
`BrowseTreatyContract_Life_RD` · `BrowseTreatyYear_Life_RD` · `BrowseTreatyYear_RD`.

⚠️ **Gerbang layar yang sudah terbaca dari keempat pembungkusnya:** **40 hidup · 15 mati**.
Penanda mati di sini berarti **niat menyembunyikan**, bukan gerbang rusak.

## Keputusan work owner yang mengikat

⛔ **Belum ada** — dan **belum diperlukan**. Tiket ini membaca, bukan memutuskan.

## Apa yang harus diuji

⛔ **Belum ada test** — hasilnya berupa bacaan yang menambal **tiket 04 · 05 · 06 · 07 · 11**.

**Yang harus dihasilkan tiket ini:**

- [ ] untuk tiap **4 Section**: kolom yang ditampilkan, gerbang layarnya, dan medan yang dapat disunting
- [ ] untuk tiap **11 laporan**: kelas sumbernya, saringannya, kolom yang ditampilkan
- [ ] jawaban apakah **`BrowseRateLife_RD` / `BrowseRateLifeSummary`** adalah sumber tampilan rate
- [ ] jawaban apakah salah satu laporan **sudah menghitung total share** — bila ya, **tiket 11**
      berubah bentuk

## Butir `[terbuka]` yang menyentuhnya

- ⚠️ medan kolom Harness **tidak ketemu** di `pyPropertyName` — medan yang benar belum diketahui
- ⚠️ kelas/kolom/saringan laporan **tidak ketemu** di tiga medan yang disisir
- ⭐ **Pertanyaan A** — bila `BrowseRateLife_RD` menerangkan bentuk `RIRATE`, Pertanyaan A bisa
  terjawab dari korpus tanpa menunggu work owner

---

## Urutan ketergantungan

```
14 (baca 4 layar + 11 laporan)   ──┐
   tidak diblokir siapa pun        │  hasilnya menambal 04 · 05 · 06 · 07 · 11
                                   │  dan BISA menjawab Pertanyaan A
                                   ▼
Pertanyaan A  (RIRATE teks?)  ───► 07  (business dan tampilan rate)   ⛔ MACET tanpa ini
                                   ▼
Pertanyaan B  (validasi?)     ───► 13  (validasi simpan)              ⛔ MACET tanpa ini
                                   │
                                   └─► menyentuh 05 · 06 · 07 · 10
```

⭐ **Tiket 14 didahulukan** — ia satu-satunya yang **tidak diblokir apa pun**, dan hasilnya
**berpeluang menjawab Pertanyaan A dari korpus** sehingga tiket 07 tidak perlu menunggu work owner.

⛔ **Tiket 13 tidak boleh dimulai** sebelum Pertanyaan B dijawab — dan ⚠️ **bisa gugur seluruhnya**
bila arti flag `false` ternyata berbeda *(lihat `grilling-ronde-2.md` §H2)*.

---

## Rekap

| | Jumlah |
| --- | --- |
| Tiket lama **ditambal** | ⭐ **8** — 03 · 05 · 07 · 08 · 09 · 10 · 11 · 12 |
| di antaranya **MERALAT** | **3** — 03 · 07 · 10 |
| di antaranya **MENAMBAH** | **5** — 05 · 08 · 09 · 11 · 12 |
| Tiket lama **tidak tersentuh** | **4** — 01 · 02 · 04 · 06 |
| ⭐ **Tiket BARU** | **2** — **13** *(validasi simpan yang tidak berjalan)* · **14** *(empat layar dan sebelas laporan yang belum dibaca)* |

⛔ **NOL berkas tiket lama disunting.**
