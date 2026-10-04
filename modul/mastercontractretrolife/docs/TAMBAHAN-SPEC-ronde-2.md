# TAMBAHAN SPEC — hasil Grilling Ronde 2

**Tanggal:** 2026-09-19 · **Modul:** Master Contract Retro Life
**Menunjuk balik ke:** `spec.md` *(885 baris, tidak disunting)*

> ⛔ **Ini DAFTAR TAMBAHAN, bukan spec baru.** `spec.md` **tidak disunting, satu byte pun**.
> ⛔ Kalimat lama yang diralat **dikutip apa adanya** di sebelah ralatnya — nol dihapus.
> ⛔ Nol butir `[terbuka]` / OQ dinyatakan tertutup. ⛔ Nol revisi ADR diusulkan.
> ⛔ `CREATE TABLE` NOL · DDL NOL · kode NOL.

---

## T1 · ⭐ RALAT — §8 "Gerbang konsistensi tahun": gerbangnya **MATI**, bukan keliru cara

**Bab yang disentuh:** `spec.md` **§8 Gerbang konsistensi tahun**

**Kalimat lama yang diralat** *(dikutip dari dasar yang dipakai §8 dan tiket 03)*:

> *Pega memeriksa konsistensi tahun dengan membandingkan `@substring(...,6,10)` terhadap
> `TREATYYEAR_LIFE` — cara yang keliru karena kolomnya bertipe `DATE`.*

**Yang DITAMBAHKAN, kalimat utuh siap salin:**

> `[terverifikasi]` **Pemeriksaan konsistensi tahun tidak pernah berjalan.** Ketiga baris syarat
> yang memuatnya — `Activity/SaveSecurityLife_Act.xml` langkah **4 · 5 · 6**, dan tiga baris
> kembarannya di `Activity/SaveSecurityReinsurerLife_Act.xml` langkah **4 · 5 · 6** — seluruhnya
> berflag prakondisi **`false`**, sehingga gerbangnya **tersimpan tetapi dimatikan** dan
> langkahnya berjalan **tanpa saringan**. ⭐ Jadi dasar faktualnya bukan *"Pega memeriksanya dengan
> cara yang keliru"*, melainkan ***"Pega tidak memeriksanya sama sekali"***.

**Menambah atau meralat:** ⭐ **MERALAT** — dasar faktualnya, **bukan** keputusannya.
⚠️ Keputusan `[keputusan work owner]` Q8 *(gerbang tahun **dibangun**, memakai nilai `DATE`)*
**tidak berubah** — ia justru **menguat**, karena kini terbukti tidak ada apa pun yang ditiru.

**Bukti:** `Master Contract Retro Life/Activity/SaveSecurityLife_Act.xml` langkah 4 · 5 · 6 ·
`.../SaveSecurityReinsurerLife_Act.xml` langkah 4 · 5 · 6 — medan `pyStepsPreCondition` = `false`.

---

## T2 · ⭐ TAMBAHAN — §7 "Total share dan validasi": **nol penjaga di KEDUA sisi**

**Bab yang disentuh:** `spec.md` **§7 Total share dan validasi**

**Yang DITAMBAHKAN:**

> `[terverifikasi]` **Tidak ada penjaga total share di sisi Pega maupun di sisi basis data.**
> Di Pega, `Activity/CountingPercentShare_Act.xml` **menjumlahkan** `PCTSHARE` *(langkah 3.1)* dan
> menuliskan hasilnya ke medan tampilan *(langkah 4)* — **nol perbandingan terhadap 100**, nol
> penolakan. `[data DBA]` Kelima stored procedure penulis juga **nol validasi bisnis**: tidak ada
> cek 100%, tidak ada cek anak, tidak ada unique selain kunci utama. ⭐ Karena itu aturan total
> share di sistem baru **dibangun sepenuhnya**, bukan dimigrasikan.

⚠️ **Catatan penamaan yang wajib dibawa:** hasil penjumlahan ditulis ke medan bernama
**`STDRATING`** *(standard rating)* — nama yang **tidak menyebut share sama sekali**. Bila kolom
serupa muncul saat migrasi, jangan mengira ia rating.

**Menambah atau meralat:** **MENAMBAH.**

**Bukti:** `Activity/CountingPercentShare_Act.xml` langkah **3.1** dan **4** ·
`procedure-bodies-from-dba.md` baris temuan *"Validasi bisnis — NOL"*.

---

## T3 · ⭐ RALAT — §15 "Tiga activity yang ternyata bukan rumus": jumlah langkahnya

**Bab yang disentuh:** `spec.md` **§15**

**Kalimat lama** *(dari `grilling-ronde-1.md`, yang §15 rujuk)*:

> *`CountingPercentShare_Act` … **4 langkah***
> *`TreatyLimit_TypeProtect` … (**4 langkah**)*

**Yang DITAMBAHKAN:**

> `[terverifikasi]` Dihitung ulang dengan pengurai XML yang membaca langkah **bersarang**:
> `CountingPercentShare_Act` **5 langkah** *(bukan 4)* dan `TreatyLimit_TypeProtect` **5 langkah**
> *(bukan 4)*; `SetValueRetroLimit_TreatyYearLife` **2 langkah** *(cocok)*. Selisihnya adalah
> langkah anak **3.1** dan **4.1** yang tidak terhitung pada ronde 1.
> ⭐ **Isi dan kesimpulan §15 tidak berubah** — ketiganya **nol** langkah ber-remark, sehingga
> pembacaan dari isi langkah tidak dirusak oleh langkah mati.

**Menambah atau meralat:** ⭐ **MERALAT angkanya saja.** ⛔ Kesimpulan §15 **BERTAHAN**.

**Bukti:** ketiga berkas di `Master Contract Retro Life/Activity/`, dihitung dengan jangkar
`Embed-ActivitySteps`.

---

## T4 · ⭐ TAMBAHAN — §12 "Penegakan `HASIL1`": kini dapat dinyatakan **apa** yang hilang

**Bab yang disentuh:** `spec.md` **§12 Penegakan `HASIL1` / `o_message`**

**Kalimat lama yang diralat** *(dari `grilling-ronde-1.md` §C2)*:

> *hanya `DeleteRowBusiness.xml` yang menyebut `HASIL1`; tidak satu pun activity `Save*`
> membacanya.*

**Yang DITAMBAHKAN:**

> `[terverifikasi]` **Tujuh berkas dari 66 menyebut `HASIL1`**, bukan satu — kelima rule SQL jalur
> simpan *(`SaveMasterTreatyYear_Life_SQL` · `SaveMasterTreatyContract_Life_SQL` ·
> `SaveMasterTreatyReinsurer_Life_SQL` · `SaveMasterTreatySecurityReinsurer_Life_SQL` ·
> `SaveMasterTreatyBusiness_Life_SQL`)* **mendeklarasikan** penampung galat itu di dalam blok
> anonimnya.
>
> ⭐ **Tetapi nol activity membacanya kembali.** Satu-satunya activity yang menyentuh penampung
> serupa, `Activity/DeleteRowBusiness.xml`, memakai **`HASIL12`** — penampung dengan nama berbeda.
>
> `[data DBA]` Dengan body procedure di tangan, kini dapat dinyatakan **apa yang jatuh diam-diam**:
> saat gagal, procedure menaruh pesan galat **berikut `SQLERRM`** ke `o_message`, menjalankan
> `ROLLBACK`, lalu mengembalikan kendali. ⭐ **Pega tidak pernah melihat pesan itu**, sehingga
> pengguna melihat penyimpanan yang **seolah berhasil** padahal seluruh perubahannya sudah
> dibatalkan.

⚠️ **Tambahan bentuk:** `o_message` berisi **HTML** *(`<span style="color:red">…</span>`)*. Di
sistem baru pesan itu **tidak boleh diteruskan mentah** ke antarmuka.

**Menambah atau meralat:** ⭐ **MERALAT** *(jumlah berkas)* **dan MENAMBAH** *(apa yang hilang)*.

**Bukti:** kelima berkas di `Master Contract Retro Life/RDBList/` · `Activity/DeleteRowBusiness.xml` ·
`procedure-bodies-from-dba.md`.

---

## T5 · ⭐ TAMBAHAN — §9 "Menghapus": basis data **kini mengaskade sendiri**

**Bab yang disentuh:** `spec.md` **§9 Menghapus — kaskade dengan konfirmasi**

**Kalimat lama** *(dari `grilling-ronde-1.md` §C)*:

> *nol kaskade* — dinyatakan dari sisi Pega.

**Yang DITAMBAHKAN:**

> `[terverifikasi]` Pernyataan *"nol kaskade"* **tetap benar untuk sisi Pega**: nol rule Pega yang
> menghapus baris anak saat induknya dihapus.
>
> ⭐ `[data DBA]` **Tetapi basis data kini melakukannya sendiri.** Empat kunci tamu dengan mode
> **`ON DELETE CASCADE`** sudah dipasang: kontrak→tahun · reinsurer→kontrak · security→reinsurer ·
> business→kontrak. Kelima tabel juga sudah punya kunci utama pada `ID`.
>
> ⚠️ **Artinya risikonya terbalik.** Yang lama dicatat sebagai *"anak menjadi yatim"* kini menjadi
> ***"menghapus induk menghapus seluruh anaknya, dan Pega tidak tahu itu terjadi"*** — karena
> Pega tidak pernah membaca hasil penghapusan.

⚠️ Keputusan `[keputusan work owner]` Q5 — **kaskade + popup konfirmasi Ya/Batal** — **tidak
berubah**, dan `[data DBA]` menyatakannya **selaras** dengan kunci tamu itu. ⭐ Yang berubah hanya
**siapa yang melakukan kaskade**. ⛔ Butir ini **tidak dinyatakan tertutup**.

**Menambah atau meralat:** **MENAMBAH** — memperluas, tidak membatalkan.

**Bukti:** `ddl-tables-from-dba.md`, bagian *"KEADAAN FINAL"* dan tabel kunci tamu.

---

## T6 · ⭐ TAMBAHAN — §11 "Jalur simpan": lima perilaku procedure yang **tidak ada di Pega**

**Bab yang disentuh:** `spec.md` **§11 Jalur simpan — lima procedure, dipanggil apa adanya**

**Yang DITAMBAHKAN:**

> `[data DBA]` Kelima procedure berpola **identik**, dan lima perilakunya **tidak terbaca sama
> sekali dari sisi Pega**:
>
> 1. ⭐ **Simpan berbentuk *upsert* dikunci `ID`** — baris dicari lebih dulu; ada = diperbarui,
>    tidak ada = disisipkan. Pega memanggil satu rule "Save" tanpa membedakan keduanya.
> 2. ⭐ **Pengenal baris diciptakan basis data**, bukan aplikasi — `'1'` diikuti nomor urut
>    enam digit dari *sequence* milik tiap tabel. Pengenal yang dikirim Pega **diabaikan** saat
>    baris baru disisipkan.
> 3. ⭐ **Cap waktu selalu waktu basis data.** Parameter cap waktu yang dikirim Pega **diabaikan**.
> 4. ⭐ **Setiap procedure menutup transaksinya sendiri**, dan membatalkannya sendiri saat gagal.
>    Tidak ada transaksi yang melintasi beberapa baris.
> 5. ⭐ **Pesan galat berisi HTML** berikut teks galat basis data mentah.

**Menambah atau meralat:** **MENAMBAH.**

**Bukti:** `procedure-bodies-from-dba.md`, tabel *"Temuan lintas kelima procedure"*.

---

## T7 · ⭐ TAMBAHAN — §14 "Skema Oracle target": batas angka **tidak dijaga basis data**

**Bab yang disentuh:** `spec.md` **§14 Skema Oracle target**

**Yang DITAMBAHKAN:**

> `[data DBA]` Seluruh kolom uang *(enam kolom batas pada kontrak)* dan seluruh kolom share dan
> komisi bertipe **angka tanpa batas digit yang ditetapkan**. ⭐ Artinya basis data **tidak
> membatasi apa pun**, dan aturan ketelitian angka yang berlaku — **hitung penuh, nol pembulatan
> di tengah jalan, tampilkan empat angka di belakang koma** — sepenuhnya **keputusan aplikasi**.
>
> ⚠️ **Akibatnya arah kegagalan terbalik dari modul lain:** nilai yang lebih panjang dari batas
> aplikasi **akan diterima Oracle** tetapi **ditolak aplikasi**. Penolakan itu harus terlihat,
> bukan dibulatkan diam-diam.

⭐ **Satu pengecualian yang wajib dicatat:** kolom **`RIRATE`** pada tabel business bertipe **teks**,
bukan angka. ⛔ Aturan ketelitian angka **tidak berlaku** padanya, dan ia **tidak boleh dipaksa
menjadi angka** tanpa keputusan work owner — lihat **Pertanyaan A** di `grilling-ronde-2.md` §G.

**Menambah atau meralat:** **MENAMBAH.**

**Bukti:** `ddl-tables-from-dba.md`, temuan 1 · 2 · 3.

---

## T8 · ⭐ TAMBAHAN — §3 "Entitas dan hierarki": empat layar utama **belum pernah dibaca**

**Bab yang disentuh:** `spec.md` **§3 Entitas dan hierarki** · **§4 Empat grid**

**Yang DITAMBAHKAN:**

> ⚠️ `[terbuka]` **Empat layar milik modul ini belum pernah dibaca isinya** — total **1,69 MB**,
> lebih besar dari keempat layar yang sudah dibaca *(1,10 MB)*:
> `Section/InputBusinessLifeReinsurers.xml` · `Section/InputDtlRetrocessionLife.xml` ·
> `Section/InputRetroLimitReinsurers.xml` · `Section/InputSecurityLifeReinsurers.xml`.
>
> `[terverifikasi]` Ketiganya *(kecuali `InputDtlRetrocessionLife`)* adalah layar yang disambung
> keempat pembungkus "Inbox" modul ini — jadi **layar utamanya justru yang belum dibaca**.
>
> ⚠️ **Sebelas dari dua belas laporan juga belum pernah dibaca**, termasuk dua yang namanya
> menyebut *rate* — kemungkinan sumber data tampilan rate, tetapi **belum terbukti**.

**Menambah atau meralat:** **MENAMBAH** — butir `[terbuka]` baru, ⛔ **bukan** pembatalan §3/§4.

**Bukti:** sensus `Section/` dan `ReportDefinition/` di `grilling-ronde-2.md` §B4 dan §B5.

---

## T9 · ⭐ TAMBAHAN — pernyataan lingkup ekspor

**Bab yang disentuh:** `spec.md` — **bab pembuka / Further Notes**

**Yang DITAMBAHKAN:**

> ⚠️ `[terbuka]` **Ekspor korpus modul ini tidak seragam lingkungannya.** Dari 66 berkas:
> **53** tercap `pegadevnusare2` · **12** `pega` · **1** `pegaprdnusare`. Penanda mesin punya
> **lima** nilai berbeda, dua di antaranya **belum pernah dinyatakan**.
>
> ⛔ **Tidak disimpulkan bahwa modul ini dev atau production.** Dicatat supaya pembaca spec tahu
> bahwa seluruh pembacaan korpus modul ini berdiri di atas ekspor dengan sebaran itu.

**Menambah atau meralat:** **MENAMBAH.**

**Bukti:** `grilling-ronde-2.md` §A6.

---

## Rekap

| | Jumlah |
| --- | --- |
| **MENAMBAH** | **6** — T2 · T5 · T6 · T7 · T8 · T9 |
| **MERALAT** | **3** — T1 · T3 · T4 |
| **TOTAL butir** | **9** |

**Bab `spec.md` yang tersentuh:** **§3** · **§4** · **§7** · **§8** · **§9** · **§11** · **§12** ·
**§14** · **§15** · bab pembuka — ⭐ **sepuluh tempat**.

⚠️ **Catatan penomoran AC:** `spec.md` memuat bab *Acceptance Criteria* **tanpa nomor "AC N"** yang
dapat dirujuk — penyaringan `AC [0-9]+` atas seluruh berkas menghasilkan **nol**. ⛔ Karena itu
butir di atas menunjuk **nomor bab**, bukan nomor AC. **Itu batas dari berkasnya, bukan kelalaian
ronde ini.**

⛔ **Nol butir `[terbuka]` / OQ dinyatakan tertutup.**
⛔ **`spec.md` tidak disunting.**
