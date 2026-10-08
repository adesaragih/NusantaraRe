# Spec — Komite Claim Prop (migrasi Pega → Go + React + Oracle)

**Tanggal:** 2026-09-18 · **Lingkup:** Komite Claim Prop saja `[keputusan work owner]` 2026-09-18
**Sumber:** `grilling-ronde-1.md` … `grilling-ronde-8.md` di folder yang sama.
**Korpus:** `D:\XML\RNM_BRD\Komite Claim Prop\` — 80 berkas rule, READ-ONLY.

---

## RALAT 08-10-2026 — implementasi satu modul

> Prompt `_brief/PROMPT-IMPLEMENTASI-MODUL-KOMITE-CLAIM-PROP.md` §3 dan §7. Kalimat lama di bab-bab di bawah **tetap
> tertulis** (dikutip di sini); yang berlaku adalah butir berikut.

1. **Kontrak** — keputusan 18 dan `urutan-tiket.md` (*"seam memakai ulang milik Claim Prop, tidak menambah seam baru"*):
   seam UJI tetap milik Claim Prop; batas MODUL kini kontrak baru `kontrak.KlaimTreatyKomite`
   (`inti/backend/kontrak/klaimtreaty.go`), disediakan `claimprop`, dipakai modul ini. Tulisan ke tabel klaim Prop
   hanya lewat kontrak itu, di dalam transaksi Submit.
2. **Kasus komite lahir** — opsi B 07-10-2026: Claim Prop melahirkan kasus `TKMT-` beserta tangga awal; modul ini
   tidak melahirkan kasus.
3. **JSON / procedure / COMMIT** — `InsertJsonClaimTreaty_act` (`PEGA_JSON_KLAIM_PNC`) dan `SaveOSClaim_SQL`
   (`PEGA_JSON_OS_AKSEP_KLAIM`) ditulis ulang sebagai SQL langsung; procedure tidak dipanggil, nol `COMMIT`.
   `JSON_KLAIM` tanpa DATA_JSON; `OS_AKSEPTASI_KLAIM` **dengan** DATA_JSON halaman `TempOSAkseptasi` (keputusan work
   owner 08-10-2026 sore, *"isi json nya khusus os_akseptasi_klaim"*).
4. **AC 56** (*"Kesembilan rule yang di Pega menyimpan sendiri tetap menyimpan sendiri"*) — diganti SATU transaksi
   aplikasi per Submit (prompt §6 butir 7); efek keluar diantre outbox di transaksi itu. Penyimpangan sadar
   (`docs/PARITAS.md`).
5. **Tiga jalur** — TT 2 (ADJUSTMENT) dibangun; TT 3 (`KomitePost_Reject`) tanpa penulis `TransferType = 3` di
   korpus Claim Prop (hanya dibaca `SendEmailKlaimRejectClose`) — tak terjangkau; TT 4 (`KomitePost_Close`) ditunda
   **OQ-CP-06** (keputusan work owner 08-10-2026).
6. **Nomor akseptasi** — rujukan ke `GenerateNoAcceptTreaty` / `GENERATE_NOACCEPTTREATYIN` adalah jalur ter-remark
   (S16.1-S16.4); yang hidup S16.5-S16.9 lewat `inti/backend/penomor`.
7. **Ketelitian kolom** — *"20 digit seluruhnya, 8 di belakang koma"* → DDL yang ada `NUMBER(38,10)` (keputusan
   work owner 07-10-2026); modul ini tidak membuat kolom uang.
8. **Indeks `ADJUSTMENT_ID`** — *"index UNIK"* → indeks biasa (migrasi 681, keputusan work owner 08-10-2026);
   keunikan dijaga `KOMITE_ID UNIQUE` di kedua tabel adjustment.
9. **AC 81-86** ditempel ke tiket: 81-82 → 04; 83 → 13; 84 → 03 (aksi tak ada di wajah TT 2); 85-86 → 11.
10. **Jawaban work owner 08-10-2026 (laporan implementasi):** OQ-KCP-01 "a" — isian Subjectivity tingkat 1 disimpan di
    `T_GENERAL_KOMITE` (migrasi 682); OQ-KCP-06 "a" — baris subjectivity dapat diserahkan ulang ke komite (Claim Prop +
    komite, S7); OQ-KCP-02 "b" — **AC 83** (*"Seluruh data kasus komite dipindahkan"*) DIRALAT: kasus komite lama tidak
    dimigrasi, klaimnya dimuat Claim Prop; OQ-KCP-03 "ikuti" — bacaan lintas skema bergerbang IsPEGAPROD.
11. **Ekspor tambahan work owner 08-10-2026 (OQ-KCP-04):** prompt values `AcceptStatus` / `AcceptanceStatus` /
    `KomiteAproval` / `Payable` / `SubjectivityNote` → label layar dan dropdown "Subjectivity Note"; harness
    `ViewClaimFormKomite` → "View more details" aktif (berkas Claim Prop hanya-baca); stream `EmailKlaim_HTML_KMT` dan
    `FILEAcceptanceNote` → isi email dan PDF dokumen akseptasi dirakit saat efek dikirim (MUATAN outbox hanya
    pengenal). Harness `Confirm` = bawaan platform Pega, tidak dibutuhkan. OQ-KCP-05 diabaikan work owner. Konversi
    PDF (`HTMLToPDF`): OQ-KCP-07 "A" — pustaka Go `github.com/go-pdf/fpdf` (go.mod), PDF digambar dari halaman
    TempAcceptedNo.

## Cara membaca berkas ini

Pembacanya **pengembang Go dan React yang tidak tahu Pega**. Istilah Pega hanya dipakai bila perlu,
dan dijelaskan sekali di tempat pertama muncul.

**Berkas ini menerangkan APA YANG ADA di sistem lama, bukan apa yang sebaiknya dibuat.** Kalimat
yang menetapkan rancangan hanya sah bila bertanda `[keputusan work owner]`.

| Tanda | Artinya |
| --- | --- |
| `[terverifikasi]` | dibaca sendiri dari korpus dengan parser XML; buktinya disebut |
| `[data work owner]` | fakta yang diberikan work owner, bukan dari korpus |
| `[keputusan work owner]` | keputusan; **hanya ini yang mengikat rancangan** |
| `[data DBA]` | DDL atau isi basis data produksi dari DBA |
| `[terbuka]` | **belum terjawab** — jangan ditebak, jangan dianggap selesai |
| ⚠️ | risiko atau jebakan yang mudah terlewat |
| ⭐ | temuan yang mengubah gambaran |

**Istilah Pega yang dipakai berulang:**

| Istilah | Artinya di sini |
| --- | --- |
| **rule** | satu berkas aturan; identitasnya `pxInsName`, berbentuk `KELAS!NAMA` |
| **activity** | urutan langkah bernomor — mirip fungsi |
| **step / langkah** | satu baris kerja di dalam activity; bernomor, boleh bersarang (`16.5`) |
| **flow** | diagram daur hidup kasus |
| **assignment** | tahap yang menunggu manusia; masuk ke kotak kerja seseorang |
| **flow action** | tombol/aksi yang dikerjakan orang saat memegang assignment |
| **section** | satu layar atau bagian layar |
| **work class / kelas kerja** | jenis kasus; menentukan properti apa yang dimiliki kasus |

⛔ **Bukti perilaku selalu berupa: path berkas + nama rule + nomor langkah Pega.** Tidak ada nomor
baris XML di berkas ini — nomor baris tidak stabil dan tidak bisa diverifikasi ulang.

#### Berkas `Struktur_*.xlsx` di dalam folder korpus — **bukan sumber kebenaran**

`[keputusan work owner]` 2026-09-19 — **berkas `Struktur_*.xlsx` di dalam folder korpus dibuat oleh
tim sendiri untuk memetakan XML.** Ia **artefak turunan**, bukan sumber kebenaran: ⛔ **tidak boleh
dikutip sebagai `[terverifikasi]`**, dan bukti tetap wajib **path berkas + nama rule + nomor langkah
Pega**.

⚠️ **Karena itu setiap sensus korpus wajib menyebut penyebutnya "berkas `.xml`", bukan "berkas".**
Folder `Komite Claim Prop\` berisi **80 berkas `.xml`** ditambah **satu `.xlsx` pemetaan** —
**81 berkas seluruhnya**. Angka **80** yang dipakai spec ini **SAH**, karena ia menghitung **rule**.

---

## Problem Statement

Nusantara Re menjalankan persetujuan komite untuk klaim reasuransi **proporsional** di atas Pega.
Setiap kali sebuah baris penyesuaian klaim (*adjustment*) perlu disetujui, sistem membuat satu
kasus komite tersendiri, lalu mengedarkannya ke para penyetuju **satu per satu, berurutan**. Ketika
penyetuju terakhir menyetujui, sistem menerbitkan nomor akseptasi, membuat dokumen PDF, mengirim
data pembayaran ke Kasir, memanggil layanan arasapas, menulis beberapa tabel riwayat, dan
mengirim email.

Pega akan ditinggalkan. Persoalannya:

1. **Perilakunya tidak tertulis di mana pun.** Yang ada hanya 80 berkas rule hasil ekspor, dan
   aturan bisnisnya tersebar di gerbang langkah, gerbang layar, dan potongan Java.
2. **Sebagian besar aturannya tidak terlihat dari luar.** Delapan ronde penyisiran menemukan
   **empat keluarga wadah berbeda** yang masing-masing menyimpan aturan — tiga di antaranya baru
   ketahuan di ronde 4, 6, dan 7. Setiap kali satu wadah terlewat, kesimpulan "tidak ada" menjadi
   salah.
3. **Ada jalur mati yang menyamar sebagai jalur hidup.** Lima perintah lompat menunjuk tanda yang
   tidak ada; satu properti hanya tampil di balik gerbang yang tidak pernah benar; satu titik lompat
   darurat sudah tidak dipakai. Semuanya masih ada di ekspor.
4. **Angka uangnya dihitung dengan presisi yang tidak seragam** — persen dibagi seratus dengan
   tiga cara berbeda di dalam satu modul, dan **tidak ada pembulatan sama sekali**.
5. **Semua efek keluar berjalan sebelum penyimpanan.** Kalau penyimpanan gagal, uang sudah
   dikirim dan email sudah sampai.

Tanpa keterangan yang bisa diverifikasi ulang, pemindahan akan menebak — dan tebakan pada jalur
pembayaran berarti uang salah kirim.

---

## Solution

Berkas ini menyusun **apa yang benar-benar dikerjakan sistem lama**, dengan setiap pernyataan
menyebut buktinya, sehingga orang berikutnya bisa memeriksa ulang tanpa mengulang delapan ronde.

Isinya **sebelas** bab: lingkup dan identitas · daur hidup kasus · tangga penyetuju · layar
komite · tiga jalur · nomor akseptasi · angka uang · efek keluar · tabel yang disentuh ·
kumpulan keputusan work owner · titik yang **sengaja diubah**. Ditambah lampiran aturan baca.

> ⛔ **RALAT 2026-09-19** — angka babnya keliru sejak terbitan pertama; kalimat lamanya **dikutip,
> tidak dihapus**:
>
> > *"Isinya **sembilan** bab: lingkup dan identitas · daur hidup kasus · tangga penyetuju · layar
> > komite · tiga jalur · nomor akseptasi · angka uang · efek keluar · tabel yang disentuh.
> > Ditambah kumpulan keputusan work owner, daftar titik yang **sengaja diubah**, dan lampiran
> > aturan baca."*
>
> Yang keliru: **kumpulan keputusan work owner** dan **titik yang sengaja diubah** sudah menjadi
> **Bab 10** dan **Bab 11** bernomor, bukan lampiran tambahan di luar hitungan. Babnya **sebelas**.

**Tiga hal yang membuat berkas ini berbeda dari catatan biasa:**

- **Butir yang belum terjawab ditulis sebagai `[terbuka]`, bukan diisi dengan tebakan.** Ada
  **tiga belas** butir semacam itu, dan semuanya ikut terbawa ke bab masing-masing.
  > ⛔ **RALAT 2026-09-19** — kalimat lamanya **dikutip, tidak dihapus**: *"Ada **sembilan** butir
  > semacam itu, dan semuanya ikut terbawa ke bab masing-masing."* Register bertambah sesudah
  > terbitan pertama; jumlah yang berlaku adalah jumlah baris di bab **Butir `[terbuka]` — daftar
  > penuh**, yaitu **13**.
- **Titik yang sengaja berbeda dari Pega dipisahkan** ke satu daftar pendek, supaya tidak tercampur
  dengan perilaku yang ditiru.
- **Aturan membaca ekspor Pega ditulis di lampiran**, supaya angka mana pun di berkas ini bisa
  dihitung ulang oleh orang lain.

### Batas dengan Claim Prop

`[keputusan work owner]` 2026-09-18 — **"intinya ini aku mau supaya ini hanya bahas Komite Claim
Prop."**

Kasus komite selalu punya **kasus klaim induk**. Batasnya ditetapkan lewat **aturan awalan**
(bab 1). Singkatnya: data yang awalannya titik atau `pyWorkPage.` **milik kasus komite**; yang
awalannya `pyWorkCover.` **milik kasus klaim** dan **hanya dibaca**, tidak digarap modul ini.

---

## User Stories

### Menerima kasus dan giliran

1. Sebagai **penyetuju komite**, saya ingin kasus komite yang menunggu keputusan saya muncul di
   kotak kerja saya, supaya saya tahu ada yang perlu saya tinjau.
2. Sebagai **penyetuju komite**, saya ingin hanya menerima kasus **saat giliran saya**, supaya
   saya tidak menilai sesuatu yang belum dilihat penyetuju sebelumnya.
3. Sebagai **penyetuju komite**, saya ingin melihat siapa saja penyetuju lain dan urutannya,
   supaya saya tahu posisi saya dalam rangkaian.
4. Sebagai **penyetuju komite**, saya ingin melihat keputusan dan catatan penyetuju sebelum saya,
   supaya saya menilai dengan konteks yang sama.
5. Sebagai **administrator**, saya ingin daftar penyetuju ditetapkan saat kasus komite dibuat,
   supaya rangkaiannya tidak berubah di tengah jalan.

### Melihat kasus

6. Sebagai **penyetuju komite**, saya ingin melihat data klaim induknya — tertanggung, nomor
   polis, tanggal kejadian, penyebab, lokasi — supaya saya paham apa yang sedang diputuskan.
7. Sebagai **penyetuju komite**, saya ingin melihat nilai penyesuaian yang diajukan berikut mata
   uangnya, supaya saya tahu besaran yang saya setujui.
8. Sebagai **penyetuju komite**, saya ingin melihat total dalam rupiah, supaya saya bisa
   membandingkan antar mata uang.
9. Sebagai **penyetuju komite**, saya ingin melihat rincian objek pertanggungan, supaya saya bisa
   memeriksa dasar perhitungannya.
10. Sebagai **penyetuju komite**, saya ingin layar menyesuaikan diri dengan jenis pengajuan, supaya
    saya tidak dibingungkan kolom yang tidak relevan.
11. Sebagai **penyetuju komite**, saya ingin melihat siapa yang membuat kasus dan kapan, supaya
    saya tahu umur pengajuan.
12. Sebagai **penyetuju komite**, saya ingin kolom yang tidak boleh saya ubah tampil sebagai
    bacaan saja, supaya saya tidak mengubah sesuatu yang bukan wewenang saya.

### Memutuskan

13. Sebagai **penyetuju komite**, saya ingin menyatakan setuju atau tolak, supaya keputusan saya
    tercatat.
14. Sebagai **penyetuju komite**, saya ingin menulis catatan pada keputusan saya, supaya alasan
    saya terbaca orang berikutnya.
15. Sebagai **penyetuju komite**, saya ingin menandai bahwa persetujuan saya bersyarat
    (*subjectivity*), supaya persetujuan itu tidak dianggap final.
16. Sebagai **penyetuju komite**, saya ingin menulis isi syarat itu ketika saya menandainya,
    supaya syaratnya jelas.
17. Sebagai **penyetuju komite**, saya ingin mengusulkan penutupan atau pencadangan klaim, supaya
    usul itu terbawa ke proses berikutnya.
18. Sebagai **penyetuju komite**, saya ingin tanggal dan identitas saya tercatat otomatis saat
    saya memutuskan, supaya saya tidak perlu mengisinya.

### Tangga persetujuan

19. Sebagai **pengelola proses**, saya ingin kasus otomatis berpindah ke penyetuju berikutnya
    setelah satu penyetuju memutuskan, supaya tidak ada langkah manual.
20. Sebagai **pengelola proses**, saya ingin kasus selesai ketika penyetuju terakhir menyetujui,
    supaya tidak ada tahap menggantung.
21. Sebagai **pengelola proses**, saya ingin satu penolakan menghentikan rangkaian, supaya
    penyetuju berikutnya tidak diminta menilai sesuatu yang sudah ditolak.
22. Sebagai **pengelola proses**, saya ingin sisa penyetuju ditandai otomatis ketika rangkaian
    berhenti karena penolakan, supaya daftar penyetuju tidak tertinggal setengah terisi.

### Akseptasi dan hasilnya

23. Sebagai **bagian akseptasi**, saya ingin nomor akseptasi terbit **hanya sekali**, saat
    penyetuju terakhir menyetujui, supaya tidak ada nomor ganda.
24. Sebagai **bagian akseptasi**, saya ingin nomor akseptasi tidak terbit ketika persetujuan masih
    bersyarat, supaya nomor hanya melekat pada persetujuan penuh.
25. Sebagai **bagian akseptasi**, saya ingin baris akseptasi tercatat di daftar akseptasi klaim,
    supaya angkanya bisa ditelusuri.
26. Sebagai **bagian akseptasi**, saya ingin dokumen akseptasi PDF dibuat otomatis, supaya tidak
    disusun manual.
27. Sebagai **bagian keuangan**, saya ingin data pembayaran terkirim ke Kasir setelah akseptasi,
    supaya pembayaran bisa diproses.
28. Sebagai **bagian keuangan**, saya ingin pengiriman ke Kasir hanya terjadi bila nomor akseptasi
    sudah tercatat, supaya tidak ada pembayaran tanpa dasar.
29. Sebagai **pihak terkait**, saya ingin menerima email pemberitahuan hasil akseptasi, supaya saya
    tahu tanpa membuka sistem.
30. Sebagai **auditor**, saya ingin setiap keputusan komite tercatat di riwayat akseptasi, supaya
    bisa diperiksa belakangan.
31. Sebagai **auditor**, saya ingin penolakan klaim tercatat tersendiri, supaya alasan penolakan
    bisa dilacak.

### Kejelasan dan kepercayaan angka

32. Sebagai **pengguna angka**, saya ingin kurs yang dipakai tersimpan bersama transaksinya, supaya
    angka lama tidak berubah ketika kurs berubah.
33. Sebagai **pengguna angka**, saya ingin hasil hitungan di layar ditampilkan dengan jumlah angka
    desimal yang tetap, supaya angka yang sama tidak tampil berbeda-beda.
34. Sebagai **bagian keuangan**, saya ingin nilai yang dikirim ke Kasir dihitung dengan cara yang
    sama di semua jalur, supaya tidak ada selisih yang tidak bisa diterangkan.
35. Sebagai **auditor**, saya ingin tahu angka mana yang disimpan dan angka mana yang dihitung
    ulang, supaya saya tahu mana yang bisa berubah.

### Yang sengaja tidak dibawa

36. Sebagai **pengembang**, saya ingin jalur mati tidak ikut dipindahkan, supaya saya tidak
    menghabiskan waktu memindahkan kode yang tidak pernah berjalan.
37. Sebagai **pengembang**, saya ingin tahu titik mana yang **sengaja diubah** dari sistem lama,
    supaya saya tidak menganggapnya cacat pemindahan.
38. Sebagai **pengembang**, saya ingin tahu pertanyaan mana yang **belum terjawab**, supaya saya
    tidak menebak di tempat yang berbahaya.
39. Sebagai **pemeriksa**, saya ingin bisa menghitung ulang sendiri angka apa pun di berkas ini,
    supaya saya tidak perlu percaya begitu saja.

---

## Implementation Decisions

> ⛔ Sembilan bab berikut **menerangkan sistem lama**. Kalimat yang menetapkan rancangan sistem
> baru **hanya** yang bertanda `[keputusan work owner]`.

---

### Bab 1 — Lingkup dan identitas

*Bab ini menjawab: apa saja yang termasuk modul ini, bagaimana mengenali satu aturan, dan di mana
batasnya dengan modul klaim.*

**Kelas kasus komite** `[terverifikasi]` adalah **`ASM-FW-GCNMFW-Work-KomiteTreaty`**. Nama
"Treaty" adalah sebutan Pega untuk lini **proporsional** — lihat ronde 1 §3.6.

**Isi modul** `[terverifikasi]` — 80 berkas rule (ronde 1 §1): Activity 35 · RDBList 23 · When 6 ·
ReportDefinition 4 · ConnectREST 3 · DataTransform 2 · FlowAction 2 · Section 2 · Flow 1 ·
DecisionTable 1 · SystemSettings 1. Folder `Harness` ada dan kosong.

#### Identitas satu aturan

⭐ **Identitas rule = `pxInsName`, berbentuk `KELAS!NAMA`.** `[terverifikasi]` 1023 dari 1024
berkas di enam modul ber-`pxInsName` sama dengan nama berkasnya; satu-satunya pengecualian adalah
artefak ekspor ganda Windows (ronde 1, blok RALAT di kepala berkas).

⚠️ **Kecocokan nama berkas BUKAN bukti rule yang sama.** Pelajaran ini muncul **dua kali**:

- `[keputusan work owner]` 2026-09-18 — **"itu hanya nama yang sama. Ada nama yang sama tapi tetap
  menjalankan sesuai modulnya masing-masing."** Temuan lama *"`KomitePost*` dipanggil dari empat
  lini"* **dibatalkan** karena lahir dari mencocokkan nama berkas (ronde 3 §3.5).
- `[terverifikasi]` Dua berkas bernama `ViewDetailInterest.xml` di modul ini ternyata **rule
  berbeda pada kelas berbeda** — `Section/` berkelas `Work-KomiteTreaty`, `FlowAction/` berkelas
  `ASM-FW-GISFW-Data-TreatyInTotal` (ronde 8 §C1).

#### ⭐ Aturan awalan — pemilah data komite dan data klaim

`[keputusan work owner]` 2026-09-18:

| Awalan | Pemilik |
| --- | --- |
| **`.Xxx`** | kasus **KOMITE** |
| **`pyWorkPage.Xxx`** | kasus **KOMITE** *(ditulis lengkap)* |
| **`pyWorkCover.Xxx`** | kasus **KLAIM induk** — **dibaca saja, tidak digarap modul ini** |

Komite **membaca** pengajuan baris penyesuaian yang dikirim kasus klaim induknya. Isi baris itu
sampai ke kasus komite lewat **penyalinan saat kasus dibuat**, oleh
`Claim Prop/Activity/AddKomiteTreatyChild_ACT` (`ASM-FW-GCNMFW-DATA-ADJUSTMENT!ADDKOMITETREATYCHILD_ACT`),
yang memakai metode `Page-Copy` dan menulis 36 penugasan bernama `Adjustment` (ronde 5 §C3)
`[terverifikasi]`. **Tidak ada tulis-menulis lintas modul saat berjalan.**

#### Jejak Save-As

`[terverifikasi]` 20 rule di modul ini dibuat dengan menyalin rule lain (ronde 1 §2a). Gunanya
untuk migrasi: rule hasil salin **mewarisi kolom, parameter, dan gerbang** dari leluhurnya, jadi
kolom yang tampak tanpa penulis sering berasal dari sana. ⛔ Silsilah **bukan** bukti dua rule
berperilaku sama.

#### `[terbuka]` bab ini

- **Pewarisan kelas kerja tidak ada di ekspor ini.** `[terverifikasi]` Tidak ada satu pun berkas
  `Rule-Obj-Class`, dan penyisiran seluruh wadah teks (`pyMemo`, `pyUsage`, `pyXMLSignature`,
  `pyNotes`, `pyDescription`) di 80 berkas **tidak menemukan** penyebutan kelas induk (ronde 7 §D3).
  **Yang perlu diminta ke work owner:** ekspor `Rule-Obj-Class` untuk `Work-KomiteTreaty` ·
  `Work-ClaimTreaty` · `Data-Adjustment` · `Data-ClaimData` · `Data-Comitee`, **atau** tangkapan
  layar tab *Class* masing-masing yang menampilkan kelas induk dan properti yang dideklarasikan.

---

### Bab 2 — Daur hidup kasus

*Bab ini menjawab: kasus komite melewati tahap apa saja, apa yang memindahkannya, dan kapan
selesai.*

`[terverifikasi]` Modul ini punya **satu** berkas flow: `Flow/KomiteTreaty_Flow.xml`
(`ASM-FW-GCNMFW-WORK-KOMITETREATY!KOMITETREATY_FLOW`). Bentuknya **empat**, penghubungnya **empat**
(ronde 6 §D).

> ⭐ **TAMBAHAN 2026-09-19 — cara menghitungnya, supaya sensus berikutnya tidak salah lagi.**
> ⛔ **Kalimat di atas TIDAK diubah; ia benar.**
>
> `[terverifikasi]` Angka **empat/empat** dihitung dari wadah **`pyShapes`** dan **`pyConnectors`**
> di dalam `pyModelProcess`. Berkas alurnya **juga** memuat **satu pengubah** *(`pyModifiers`)* dan
> **tiga tiket** *(`pyTicketShapes`)*, ditambah **perute** dan **pemberitahu** yang menggantung di
> dalam bentuk — ⛔ **tak satu pun dari semua itu adalah bentuk.**
>
> ⚠️ **Sensus yang menjumlahkan seluruh entri `Data-MO-*` tanpa memisahkan wadahnya akan memberi
> 12, dan angka 12 itu BUKAN jumlah bentuk.** Dari tiga tiket itu, **dua cangkang kosong** tanpa
> nama dan tanpa pengenal; yang bernama **`komiteAccept_ticket`** — sudah diputuskan **tidak
> dibuat** *(AC 7)*. Rinciannya di `utang-lintas-modul.md` §1.

| Bentuk | Jenis | Nama | Siapa yang mengerjakan |
| --- | --- | --- | --- |
| `Start1` | mulai | — | sistem, saat kasus komite dibuat |
| **`ASSIGNMENT63`** | **assignment** | **`KomiteRouter`** | ⭐ **penyetuju komite** — masuk kotak kerjanya |
| **`Decision1`** | gerbang keputusan | **`KomiteLoop`** | sistem |
| `END52` | selesai | — | sistem · status akhir **`Resolved-Completed`** |

| Penghubung | Dari → ke | Pemicu |
| --- | --- | --- |
| `Transition1` | `Start1` → `ASSIGNMENT63` | selalu |
| **`TRANSITION54`** | `ASSIGNMENT63` → `Decision1` | ⭐ **aksi pengguna `ViewTransferDtl`** |
| **`Transition2`** | `Decision1` → **`ASSIGNMENT63`** | ⭐ **bila `IsKomiteLoop`** — **jalur balik** |
| `Transition3` | `Decision1` → `END52` | selain itu, bernama `NoLoop` |

**Ceritanya:** kasus dibuat → langsung masuk kotak kerja penyetuju, dan `KomiteRouter` yang
menentukan ke siapa → penyetuju membuka layar dan menekan kirim → sistem memeriksa apakah masih
ada penyetuju berikutnya → bila ya, kasus **kembali ke kotak kerja**; bila tidak, kasus **selesai**.

⭐ **Hanya ada satu jalur balik**, dan tidak ada yang kembali ke awal.

#### Tiga aksi lain pada assignment

`[terverifikasi]` Assignment `KomiteRouter` membawa **tiga aksi lokal**, dan ketiganya **bawaan
Pega**, bukan buatan Nusantara Re: memindahkan tugas ke orang lain, membuat kasus sampingan, dan
melibatkan pihak luar (ronde 7 §D1). Siapa yang **berhak** memakainya tidak terbaca — lihat
bab 10, keputusan tentang hak akses.

#### Titik lompat darurat

`[keputusan work owner]` 2026-09-18 — **`komiteAccept_ticket` sudah tidak dipakai. Jalur mati —
dibuang, tidak dipindahkan, dan tidak ditelusuri siapa pemicunya.** Perlakuannya sama dengan lima
lompatan menggantung (bab 10).

---

### Bab 3 — Tangga penyetuju

*Bab ini menjawab: bagaimana sistem tahu giliran siapa, kapan tangga berhenti, dan apa yang terjadi
pada sisa penyetuju.*

Tiga hal yang menjalankannya, dan **tidak satu pun tampil di layar** `[terverifikasi]` (ronde 2 §B4):

| Nama | Isinya |
| --- | --- |
| **`KomiteLoop`** | **berapa penyetuju yang dibutuhkan** |
| **`KomiteCount`** | **penyetuju ke berapa yang sedang berjalan** |
| **`KomiteList`** | daftar penyetuju berikut keputusan, catatan, dan tanggalnya |

#### Dari mana `KomiteLoop` diisi

`[terverifikasi]` = **jumlah baris daftar penyetuju**, ditetapkan **sekali saat kasus komite
dibuat**; pada jalur tutup/tolak yang memakai satu penyetuju tetap, nilainya `1` (ronde 1 §9.2).
**Nol penulis di modul ini** — jadi di jalur modul ini ia tidak berubah di tengah jalan.

`[keputusan work owner]` 2026-09-18 — **"ikuti activity yang berjalan."** Ia **kolom yang
disimpan**, diisi ulang di titik yang sama dengan Pega mengisinya: saat kasus komite dibuat, dan —
untuk lini Fac In — di dalam modul komitenya. **Bukan dihitung ulang setiap kali dibaca.**

#### Bagaimana giliran ditentukan

`[terverifikasi]` `Activity/KomiteRouter` punya 8 langkah, **6 di-remark**. Yang hidup: langkah
**6** dan anaknya **6.1**. Langkah 6.1 bergerbang *"keputusan penyetuju masih kosong"* dan
menetapkan **kepada siapa assignment dirutekan**, diambil dari akun operator penyetuju yang sedang
giliran (ronde 1 §5, dikoreksi ronde 4 §B5). Tangga naik di langkah **5**.

⭐ **Tangga berputar DI DALAM satu tahap.** `[terverifikasi]` Flow hanya punya satu assignment, dan
jalur baliknya menuju assignment yang sama. **Setiap penyetuju adalah kunjungan baru ke tahap yang
sama**, bukan tahap baru (ronde 6 §D5, dikuatkan ronde 7 §C4).

#### Kapan tangga berhenti

`[terverifikasi]` Tiga langkah penentu di `Activity/KomitePostAdjustment` — langkah **17** (menulis
baris akseptasi), **21** (membuat PDF), **34** (mengirim ke Kasir) — masing-masing punya **dua**
baris syarat yang **harus dua-duanya lolos**:

```
baris 1   penyetuju terakhir  DAN  keputusannya setuju
baris 2   persetujuan BUKAN bersyarat
```

⭐ Artinya ketiganya hanya berjalan pada **penyetuju terakhir, yang menyetujui, tanpa syarat**
(ronde 1 §3.2; diuji ulang ronde 4 §B1 dan tetap berdiri).

⚠️ Pada langkah **34** urutan kedua barisnya terbalik dari 17 dan 21. **Hasilnya sama.**

#### Penolakan

`[terverifikasi]` Bila keputusan **tolak**, langkah **12** melompat ke tanda `EXT`, yaitu langkah
**25**, yang memaksa pencacah ke nilai akhir sehingga tangga berhenti. Langkah **26** kemudian
menyisir daftar penyetuju dan anaknya **26.1** menandai sisa penyetuju sebagai tertolak otomatis
(ronde 1 §9.1, ronde 3 §A2).

#### `[terbuka]` bab ini

- **Arti `REPEAT` dan `EMBEDDED`, dan halaman apa yang diulang.** `[terverifikasi]` 85 langkah di
  modul ini berulang — 65 `EMBEDDED`, 20 `REPEAT`. Bedanya terbaca di struktur (`REPEAT` selalu
  membawa iterasi/awal/batas bernilai satu), **tetapi artinya tidak dinyatakan di mana pun**, dan
  **tidak ada elemen yang menyebut halaman mana yang diulang** (ronde 7 §C2).
- **Berapa kali perulangan berputar.** Menentukan berapa kali angka berubah dan berapa kali
  pengiriman ke Kasir terjadi — lihat bab 7 dan bab 8.
- **Urutan pemeriksaan bila dua keluarga gerbang sama-sama terisi** tidak terbaca dari struktur
  (ronde 4 §A3).

---

### Bab 4 — Layar komite

*Bab ini menjawab: apa yang dilihat penyetuju, apa yang boleh diisinya, dan apa yang menentukan
sebuah kolom tampil atau tidak.*

#### Bentuknya: pembungkus, badan, penyimpan

`[terverifikasi]` `FlowAction/ViewTransferDtl` (`ASM-FW-GCNMFW-WORK-KOMITETREATY!VIEWTRANSFERDTL`)
**tidak punya kolom sendiri** — ia pembungkus. Kontraknya tiga baris (ronde 2 §B2):

```
SetKomiteList_Act   menyiapkan   ->   ShowTransfer   menampilkan   ->   KomitePost   menyimpan
```

Tombolnya **Submit** dan **Cancel**.

⭐ `SetKomiteList_Act` **tidak mengisi daftar penyetuju sama sekali**, meski namanya begitu.
`[terverifikasi]` Yang dikerjakannya adalah **menjumlahkan penyesuaian per mata uang**, dan
langkah **1**-nya **menghentikan activity** bila jenis pengajuan bukan penyesuaian (ronde 5 §D).

#### ✅ Famili gerbang ketiga — **DITUTUP**

> `[keputusan work owner]` 2026-09-18 — **butir kolom layar DITUTUP. Buktinya sudah cukup.**
> **Angka 93 · komite 62 · klaim induk 31 SAH.**
>
> Dasarnya dua: wadah famili ketiga di berkas ini **ada satu dan tidak hidup**, dan hitungan 15
> gerbang famili A **sudah menyapu kedua kelas**.

`[terverifikasi]` Diperiksa langsung di berkas layar utama modul ini:

| Wadah | Isinya | Yang **HIDUP** |
| --- | --- | --- |
| **`pyUserData`** | **354 blok** — **162** berkelas *HeaderElements*, **192** berkelas *UserData* | **15** |
| **`pyDefaultUserData`** — famili ketiga | **1 blok**, berkelas *UserData* | ⭐ **0** |

**Dua hal yang dibuktikannya:**

1. ⭐ **Hitungan 15 gerbang famili A sudah MENYAPU KEDUA KELAS.** Kelas *UserData* **tidak pernah
   terlewat** — ia hanya tidak pernah disebut namanya di ronde 1–8.
2. ⭐ **Famili ketiga di berkas ini ADA SATU dan TIDAK HIDUP.**

⛔ **Karena itu angka di bawah tidak berubah.**

#### Berapa kolomnya

`[terverifikasi]` `Section/ShowTransfer` memuat 348 sel; yang benar-benar menampilkan medan
menghasilkan **93 properti berbeda** (ronde 6 §C):

| Ember | Properti berbeda |
| --- | --- |
| **KOMITE** — awalan titik dan `pyWorkPage.` | **62** |
| **KLAIM induk** — awalan `pyWorkCover.` | **31** |
| **TOTAL** | **93** |

**Definisi yang dipakai:** hanya sel yang benar-benar menjadi kotak di layar. Label, teks bantu,
dan properti yang hanya muncul di dalam teks syarat **tidak dihitung** — semuanya menentukan
tampilan, tetapi bukan kotak (ronde 6 §C1).

⚠️ **Dua yang berada di luar hitungan tetapi DATA NYATA dan wajib ikut terbawa:** **"dibuat oleh"**
dan **"tanggal dibuat"**. Keduanya tampil di layar; keduanya di luar hitungan hanya karena aturan
pencatatan memisahkan properti bawaan Pega (ronde 7 §A3).

#### Apa yang menentukan sebuah kolom tampil

`[terverifikasi]` **Dua famili gerbang, bukan satu** (ronde 3 §B):

| Famili | Letaknya | Hidup |
| --- | --- | --- |
| **A** | pada **sel** — di sub-halaman milik sel itu | **15** |
| **B** | pada **layout** yang membungkus sel | **9** |

**Penanda gerbang mati: `NEVER` dan `1=2`.** `[terverifikasi]` **12 sel** duduk di baliknya, jadi
**tidak pernah tampil** (ronde 3 §B4).

⭐ **Saklar besar layar ini adalah jenis pengajuan.** Dari 62 properti komite, **24** hanya tampil
pada jalur **penyesuaian**. Tiga judul bagian layar — ADJUSTMENT, REJECT, CLOSE — juga bergerbang
jenis pengajuan, jadi layar **berganti wajah** (ronde 3 §B1).

#### Apa yang boleh diisi penyetuju

`[terverifikasi]` (ronde 3 §B3):

| Isian | Kapan bisa diisi |
| --- | --- |
| **keputusan setuju/tolak** | ✅ **selalu** · wajib isi |
| **catatan** | ✅ **selalu** · wajib isi |
| penanda persetujuan bersyarat | bila **menyetujui** dan **jalur penyesuaian** |
| isi syarat | bila penanda bersyarat **dicentang** |
| usul tutup klaim | bila **jalur penyesuaian** |
| usul cadangkan klaim | bila **jalur penyesuaian** |

⭐ Rantainya bertingkat dan masuk akal: isi syarat baru muncul **setelah** penandanya dicentang.

#### Layar kedua

`[terverifikasi]` `Section/ViewDetailInterest` (kelas komite) memuat **15 sel**, hanya **3**
medan — nama objek, mata uang, dan nilai pertanggungan per objek — **ketiganya hanya-baca dan
wajib isi**. **Ketiganya sudah ada di layar utama**, jadi **nol kolom tambahan** (ronde 8 §C).

`[terverifikasi]` **Tidak ada layar ketiga.** Modul ini punya 2 Section dan 2 FlowAction,
**keempatnya sudah dibaca** (ronde 8 §C6).

#### ⭐ Daftar penyetuju DITAMPILKAN — **titik yang sengaja diubah**

⚠️ `[terverifikasi]` **Di Pega, penyetuju memutuskan tanpa melihat siapa pun sebelum dia.** Bab 3
sudah mencatatnya: ketiga hal yang menjalankan tangga — **`KomiteLoop`**, **`KomiteCount`**, dan
**`KomiteList`** — **tidak satu pun tampil di layar**, padahal `KomiteList` justru berisi *"daftar
penyetuju berikut keputusan, catatan, dan tanggalnya"*. Penyaringan ke-93 medan layar
**tidak menemukan satu pun** yang menampilkannya.

⚠️ **Dua user story menuntut sebaliknya** — **US 3** *"melihat siapa saja penyetuju lain dan
urutannya"* dan **US 4** *"melihat keputusan dan catatan penyetuju sebelum saya, supaya saya menilai
dengan konteks yang sama"*. Keduanya ditulis di terbitan pertama, **sebelum** diketahui bahwa daftar
itu tidak tampil, dan **tidak pernah dicabut**.

`[keputusan work owner — atas rekomendasi asisten]` 2026-09-19 — **daftar penyetuju DITAMPILKAN di
layar komite.** Yang tampil: **urutan**, **siapa**, dan — untuk penyetuju yang **sudah** memutuskan
— **keputusan, catatan, dan tanggalnya**. Penyetuju yang **belum** memutuskan tampil tanpa isi.

**Dasarnya tiga, dan tidak satu pun berupa tebakan atas Pega:**

1. **Datanya sudah ada.** Kesembilan kolom baris penyetuju sudah dikunci di bab 9 — `KOMITE_URUT` ·
   `KOMITE_OPERATORID` · `KOMITE_JABATAN` · `KOMITE_APPROVAL` · `KOMITE_COMMENT` · `DATE_APPROVE`.
   Menampilkannya **tidak menuntut satu kolom pun yang baru**.
2. **AC 53 sudah mewajibkan keputusan tersimpan per baris penyetuju**, bukan di header. Yang kurang
   hanya **menampilkannya kembali**.
3. **Mencabut US 3 dan US 4 berarti membuang niat yang tertulis** dan tidak pernah ditarik —
   perbuatan yang lebih besar daripada mencatat satu penyimpangan yang mudah dibatalkan.

⛔ **Ia penyimpangan sadar, bukan paritas** — lihat **bab 11 titik 5**, **AC 79** dan **AC 80**.
⛔ **Tandanya sengaja `[keputusan work owner — atas rekomendasi asisten]`:** pilihannya diserahkan
kepada asisten, dan **pencabutannya murah** — menyentuh dua AC, satu baris bab 11, satu keputusan
bab 10, dan tiket **03**. ⛔ **Nol kolom basis data tersentuh bila dicabut.**

⚠️ **Yang TIDAK diputuskan di sini:** apakah penyetuju boleh melihat keputusan penyetuju
**sesudahnya** *(tidak ada isinya saat ia memutuskan)*, dan apakah daftar itu tampil juga pada jalur
**REJECT** dan **CLOSE**. Keduanya mengikuti aturan layar yang sudah ada di bab ini.

#### `[terbuka]` bab ini

- ✅ **Angka 93 · 62 · 31 — DITUTUP** `[keputusan work owner]` 2026-09-18. **Bukan lagi butir
  terbuka.** Dasarnya di blok di atas.
- **Dari mana layar kedua dibuka belum terbaca.** `[terverifikasi]` Section komite
  `ViewDetailInterest` tidak dirujuk oleh layar utama, pembungkusnya, maupun berkas flow; satu-
  satunya flow action bernama sama merujuk section **pada kelas lain** (ronde 8 §C4).

---

### Bab 5 — Tiga jalur

*Bab ini menjawab: kasus komite bisa berupa apa saja, dan apa bedanya di layar dan di mesin.*

`[keputusan work owner]` 2026-09-18 — **nama resmi ketiga jalur:**

| Kode | Nama |
| --- | --- |
| **2** | **TRANSFER ADJUSTMENT** — persetujuan penyesuaian nilai klaim |
| **3** | **REJECT** — penolakan |
| **4** | **CLOSE** — penutupan klaim |

`[terverifikasi]` Rule pemilah `KomitePost` memilah **tepat tiga** arah itu (ronde 1 §3.5).

#### Kode keempat yang tidak ada

`[terverifikasi]` **Tidak ada satu pun penulis yang menetapkan nilai `1`.** Yang ada hanya pembaca:
teks syarat `TransferType=="1"` muncul tujuh kali, dan **ketiga langkah yang memakainya bergerbang
mati** (ronde 1 §9.8). **Kesimpulan: `1` adalah sisa jalur lama.**

#### Apa bedanya di layar

`[keputusan work owner]` 2026-09-18 — pada jalur **REJECT** dan **CLOSE**, penyetuju memang hanya
mengisi **dua** hal: keputusan dan catatan. Empat isian lain hanya ada di jalur **TRANSFER
ADJUSTMENT**. **Ditiru apa adanya.**

⚠️ Dikuatkan mesinnya: penyiap layar **berhenti di langkah pertama** bila jenis pengajuan bukan
penyesuaian, dan catatan pengembangnya sendiri berbunyi *"skip selain TransferType = 2"*
`[terverifikasi]` (ronde 5 §D, ronde 7 §B5).

---

### Bab 6 — Nomor akseptasi

*Bab ini menjawab: kapan nomor akseptasi terbit, dari mana angkanya, dan mengapa hanya sekali.*

`[terverifikasi]` Penerbitan nomor ada di satu blok langkah bersarang di
`Activity/KomitePostAdjustment` — langkah **16** dan anak-anaknya **16.1** sampai **16.9**
(ronde 1 §3.4, ronde 2 §C).

| Langkah | Yang dikerjakan |
| --- | --- |
| **16** | **induk blok**; bergerbang **bukan persetujuan bersyarat**; merupakan langkah **berulang** |
| 16.1 · 16.2 · 16.3 · 16.4 | ⛔ **di-remark** — tidak berjalan |
| **16.5** | mengambil kode produksi |
| **16.6** | menyiapkan nilai |
| **16.7** | membangkitkan bulan-tahun dan nomor urut |
| **16.8** | merangkai potongan nomor |
| **16.9** | **menuliskan nomor akseptasi dan menandai status akseptasi** |

⭐ **Gerbangnya berlapis dua:** di tingkat induk *"bukan persetujuan bersyarat"*, di tingkat anak
*"penyetuju terakhir dan disetujui"*. Itulah yang membuat nomor terbit **sekali saja**, di ujung
rangkaian (ronde 1 §3.2).

#### Rule yang dipanggil

⚠️ Nama rule pada langkah pengambilan data **tidak tertulis di satu tempat** — ia terpecah tiga
bagian yang harus dibaca bersama (lampiran, aturan 4). `[terverifikasi]` Hasilnya (ronde 2 §C):

| Langkah | Rule |
| --- | --- |
| 16.1 *(remark)* | `ASM-FW-GISFW-INT-POLICYJSON!RNM!GETTANGGALCLOSING_SQL` |
| 16.4 *(remark)* | `ASM-FW-GCNMFW-INT-V_POLIS!GCNM!GENERATENOACCEPTTREATY` |
| **16.5** | `ASM-FW-GISFW-INT-POLICYJSON!RNM!GETKODEPRODNONLIFE_SQL` |
| **16.7** | `ASM-FW-GISFW-INT-POLICYJSON!RNM!GETSEQUENCENUMBER_SQL` |

⚠️ Dua rule yang langkahnya di-remark **tidak ada di ekspor modul ini** — konsisten: yang dimatikan
tidak ikut diekspor `[terverifikasi]`.

---

### Bab 7 — Angka uang

*Bab ini menjawab: dari mana kurs diambil, bagaimana total dihitung, dan seberapa teliti angkanya.*

#### Kurs

`[keputusan work owner]` 2026-09-18 — **"ikuti aja query di xml CurrencyStandard."**

`[terverifikasi]` Kurs diambil oleh rule Connect-SQL `ASM!CURRENCYSTANDARD`, lewat sebuah **fungsi
tersimpan Oracle** dengan **dua masukan saja: kode mata uang dan tanggal server saat query
dijalankan** (ronde 5 §9.6 / §A7 ronde 5). **Tidak ada nomor kasus, tidak ada lini, tidak ada
tanggal kasus.**

⭐ **Artinya kurs yang terkunci adalah kurs TANGGAL PENCARIAN DIJALANKAN** — bukan tanggal kasus
dibuat, bukan tanggal komite menyetujui.

⚠️ **Isi fungsi tersimpannya tidak ada di korpus**, dan `[keputusan work owner]` tidak diminta ke
DBA.

`[terverifikasi]` Penguncian nilainya terjadi **di luar modul ini** — pola *ambil-hanya-bila-kosong*
ada di `Claim Prop/Activity/SetNameCurrency_Act`. **Di dalam modul ini nol pola penguncian untuk
uang** (ronde 8 §B3 · §B5).

#### Yang dihitung modul ini

`[terverifikasi]` **Satu-satunya penghitung uang berkelas komite adalah `SetKomiteList_Act`**,
langkah **7.1.1**, di dalam perulangan (ronde 8 §B2):

```
total kotor        := total kotor        + nilai kotor baris
total bersih       := total bersih       + nilai penyesuaian baris
total kotor  IDR   := total kotor  IDR   + ( nilai kotor baris      x kurs )
total bersih IDR   := total bersih IDR   + ( nilai penyesuaian baris x kurs )
```

lalu langkah **7.2** menaruh kedua total pertama ke properti kasus.

⭐ **Gerbang penjumlahannya penting:** baris yang **berstatus tolak tidak ikut dijumlahkan**, dan
penjumlahan dilakukan **per mata uang**.

#### Yang dirujuk dari Claim Prop

⚠️ `[terverifikasi]` **Lima activity penghitung lain yang ada di folder ekspor modul ini berkelas
`ASM-FW-GCNMFW-Work-ClaimTreaty` — kelas kasus KLAIM, bukan komite** (ronde 8 §B1). Menurut aturan
awalan (bab 1), **pemiliknya Claim Prop**. Rumusnya dicatat di ronde 8 §B2 sebagai **bahan
rujukan**, bukan sebagai garapan modul ini.

Rantai ketergantungannya, dalam bahasa biasa (ronde 8 §B6): kurs dicari lebih dulu → persentase
diambil dari master → nilai per baris dihitung → nilai rupiah per baris dihitung → total per mata
uang ditumpuk → total kasus ditulis → spreading dihitung dari total. **Baru sesudah itu komite
masuk**, menjumlahkan baris penyesuaian yang tidak ditolak.

#### ⚠️ Seberapa teliti angkanya di Pega

`[terverifikasi]` Sensus 80 berkas, tidak peka huruf besar-kecil (ronde 8 §B4):

⛔ **TIDAK ADA PEMBULATAN SAMA SEKALI.** Tidak ada pembulatan, pemotongan, pemformatan angka,
maupun penetapan skala di mana pun.

⚠️ **Satu-satunya kendali ketelitian adalah argumen ketiga pada pembagian — dan ia tidak seragam.**
Persen dibagi seratus dengan **tiga cara berbeda**, dua di antaranya **di activity yang sama untuk
menghitung hal yang sama**: dengan ketelitian sepuluh angka, dengan ketelitian empat angka, dan
**tanpa kendali ketelitian sama sekali**. Ditambah satu pemotongan tersembunyi ke nol angka di
rule pengirim email.

#### ⭐ Ketelitian angka di sistem baru — **TITIK YANG SENGAJA DIUBAH**

`[keputusan work owner]` 2026-09-18:

| Di mana | Ketetapan |
| --- | --- |
| **Hitungan di Go** | **ketelitian penuh, nol pembulatan di tengah jalan** |
| **Tampilan layar** | **empat angka di belakang koma** |
| **Kolom angka di basis data** | ⭐ **MENGIKUTI DDL YANG SUDAH ADA** — **20 digit seluruhnya, 8 di antaranya di belakang koma.** Berlaku untuk **uang, persen, dan kurs** `[keputusan work owner]` |
| **Kolom pencacah** | **bilangan bulat biasa** — `KOMITE_LOOP` · `KOMITE_COUNT` · `KOMITE_URUT`. **Tidak ikut aturan di atas** |
| **Tipe angka di Go** | **tipe desimal, BUKAN bilangan pecahan biner** `[keputusan work owner — atas rekomendasi asisten]` |

> ## ⛔ RALAT — rekomendasi asisten DICABUT
>
> `[keputusan work owner]` 2026-09-18 — **ralat mengikuti konvensi proyek: teks dan angka lama
> DIKUTIP sebagai jejak, tidak dihilangkan.**
>
> **Yang dicabut, dikutip apa adanya:**
>
> > *"bilangan desimal berskala tetap — **18 digit di depan koma, 20 di belakang**. Berlaku untuk
> > uang, persen, dan kurs"* — bertanda `[keputusan work owner — atas rekomendasi asisten]`.
>
> **Penggantinya:** work owner memeriksa basis data dan memutuskan **mengikuti bentuk kolom yang
> sudah ada di produksi — 20 digit seluruhnya, 8 di antaranya di belakang koma**. Tandanya kini
> `[keputusan work owner]` polos; **rekomendasi asisten sudah tidak dipakai**.
>
> ⚠️ Bentuk lama **lebih lebar** dari penggantinya di kedua sisi koma. Selisih itulah yang
> melahirkan dua butir terbuka di bawah.
>
> Yang **tetap** bertanda *atas rekomendasi asisten* hanya **tipe desimal, bukan bilangan pecahan
> biner**.

> ⚠️ **PEMBULATAN TERJADI DI BATAS PENYIMPANAN, DAN ITU DITERIMA SADAR.**
>
> Hitungan di Go berjalan **sampai 20 angka di belakang koma tanpa pembulatan di tengah jalan**,
> tetapi kolomnya hanya menampung **8 angka di belakang koma**. Jadi **angka persen yang di Pega
> dihitung sampai 10 desimal akan menjadi 8 saat disimpan.** Itu **titik pembulatan satu-satunya**,
> dan letaknya **di batas penyimpanan**, bukan di tengah perhitungan.

#### `[terbuka]` bab ini

- ✅ **Batas 12 digit di depan koma — DITUTUP** `[keputusan work owner]` 2026-09-19: **ikuti
  bentuk kolom yang sudah ada apa adanya, tidak ditanyakan ke DBA.** ⚠️ **Risiko diterima
  sadar:** bila nilai melewati batas, basis data **menolak menyimpan**, bukan membulatkan.
- ⭐ **BARU — seberapa besar selisih akibat persen 10 desimal menjadi 8 saat disimpan.**
  Selisihnya **menumpuk lewat perulangan** sebelum sampai ke Kasir (lihat butir berikutnya).
  **Besarnya belum diukur.** ⛔ Tidak ditebak.
- **Mana ketelitian Pega yang sah** — tiga cara berbeda untuk operasi yang sama, ditambah
  pemotongan tersembunyi. Keputusan di atas menetapkan **apa yang dikerjakan sistem baru**;
  pertanyaan **mana yang benar di Pega** tetap terbuka, dan ia yang menentukan bagaimana selisih
  terhadap data lama diterangkan.
- **Berapa kali perulangan berputar.** `[terverifikasi]` **91 dari 121** penugasan beraritmetika di
  modul ini berada **di dalam langkah berulang** (ronde 8 §B7). Berapa kali angkanya berubah
  bergantung pada jawaban butir terbuka di bab 3.

⭐ Kabar baiknya `[terverifikasi]`: **nol perhitungan uang di dalam potongan Java.** Kedua puluh
potongan Java hanya membersihkan duplikat, merakit teks, dan menangani berkas — **rumus uang
seluruhnya terbaca dari penugasan properti biasa** (ronde 7 §B3).

---

### Bab 8 — Efek keluar

*Bab ini menjawab: apa yang keluar dari sistem saat komite menyetujui, kapan, dan apa yang terjadi
bila gagal.*

`[terverifikasi]` Delapan hal keluar, seluruhnya dari `Activity/KomitePostAdjustment` (ronde 5 §B2):

| Langkah | Rule | Yang keluar | Tujuan |
| --- | --- | --- | --- |
| **17** | `SaveAcceptation_Act` | baris akseptasi | daftar akseptasi klaim |
| **21** | `SaveAcceptationTreaty_TKMT` | **dokumen PDF akseptasi** | penyimpanan dokumen |
| **28** | `InsertJsonClaimTreaty_act` | data klaim bentuk JSON | tabel JSON klaim |
| **29** | `KonversiKlaim_Act` | panggilan layanan | ⭐ **arasapas** |
| **31** | `InsertLogServiceClaim` | baris log | tabel pemantauan klaim |
| **33** | `InsertHistoryAkseptasiPega_Sql` | 6 kolom riwayat | tabel riwayat akseptasi |
| **34** | `HitServiceToKasirKMT_Act` | data pembayaran | ⭐ **Kasir** |
| **35** | `SendEmailKlaim_KMT` | **email** | penerima akseptasi |

#### Syarat sebelum mengirim ke Kasir

`[data work owner]` 2026-09-18 — pengiriman ke Kasir **hanya boleh jalan bila nomor akseptasi sudah
tercatat** di daftar akseptasi klaim. **Itu syarat resmi yang dipertahankan.** Buktinya: rule
pemeriksa mengambil nomor akseptasi, membuang titiknya, mencari di tabel akseptasi, dan menyalakan
penanda hanya bila ketemu (ronde 5 §A7).

#### Alamat tujuan layanan

`[terverifikasi]` Alamat diambil dari **tabel alamat layanan**, disaring **dua kunci** — kategori
dan sub-kategori — sehingga hasilnya satu baris. Empat pemanggilnya mengirim pasangan kunci yang
berbeda-beda: Kasir, klaim, dan dua untuk penyimpanan berkas (ronde 6 §B1). **Nol cacat** — dugaan
lama bahwa alamat diambil tanpa saringan **sudah dicabut**.

#### ⚠️ Dua sikap berbeda terhadap kegagalan

`[terverifikasi]` (ronde 5 §B4, diuji ulang ronde 7 §B4):

| Efek keluar | Bila gagal |
| --- | --- |
| **Kasir** | ⭐ **melompat ke langkah pengiriman email, lalu selesai** — alur **tidak berhenti** |
| **arasapas** | melompat ke langkah pencatatan log — alur **tidak berhenti** |
| pengiriman ke Kasir di tingkat dalam | ada **email galat tersendiri** |
| **dokumen PDF** | ⭐ **melempar galat dan berhenti** |
| **unggah berkas** | ⛔ **tidak ditangani** |
| **email** | ⛔ **tidak ditangani** |

⭐ **Modul ini punya dua sikap yang berbeda, dan keduanya disengaja:** yang satu **melompat dan
melanjutkan**, yang lain **berhenti dengan galat**.

`[keputusan work owner]` 2026-09-18 — **kiriman ke Kasir gagal → email TETAP dikirim, kasus tetap
jalan, nol percobaan ulang. DITIRU APA ADANYA.**

> ⚠️ **Risiko yang diterima sadar:** kegagalan **tidak terlihat** sampai ada orang yang memeriksa
> belakangan. Tidak ada pemberitahuan bahwa transfer gagal; yang sampai ke penerima justru **email
> keberhasilan akseptasi**. Selisih antara "email terkirim" dan "uang terkirim" **hanya ketahuan
> dari pemeriksaan manual**.

#### ⭐ Urutan terhadap penyimpanan — **TITIK YANG SENGAJA DIUBAH**

`[terverifikasi]` Di Pega, **kedelapan efek keluar berjalan SEBELUM penyimpanan**. Penyimpanan ada
di langkah **41**, sedangkan Kasir di **34** dan email di **35**. **Nol yang sesudah** (ronde 5 §B5).

⚠️ Akibatnya di Pega: bila penyimpanan gagal, **uang sudah dikirim, email sudah sampai, dokumen
sudah dibuat, dan tiga tabel log sudah terisi** — sementara kasusnya sendiri tidak tersimpan.

> `[keputusan work owner]` 2026-09-18 — **URUTAN INI TIDAK DITIRU.** Fakta di atas dicatat sebagai
> **fakta**, tetapi **tidak mengikat rancangan**. Urutannya **akan disesuaikan nanti**.
>
> ⛔ **Jangan mengunci urutan Pega sebagai syarat.** ⛔ Urutan penggantinya **belum diputuskan** dan
> tidak dirancang di berkas ini.

#### `[terbuka]` bab ini

- **Unggah berkas dan email tanpa penanganan gagal** — kesimpulan dari **empat keluarga wadah yang
  seluruhnya sudah dibaca**, jadi cukup kuat, tetapi belum dibawa ke work owner.
- **Catatan pengembang yang tidak berjejak.** `[terverifikasi]` Rule pengunggah berkas bercatat
  **"FIX ERROR HANDLING"**, tetapi **tidak terlihat penanganan apa pun** di keempat keluarga wadah
  (ronde 7 §B5). Entah perbaikannya di tempat lain, entah catatannya menyusul perubahan yang
  kemudian dicabut.

---

### Bab 9 — Tabel yang disentuh

*Bab ini menjawab: tabel apa saja yang ditulis, siapa yang menyimpan sendiri, dan bentuk tabel yang
sudah pasti.*

#### Sembilan rule yang menyimpan sendiri

`[terverifikasi]` Dari 23 rule basis data di modul ini, **sembilan memuat perintah simpan
sendiri** — artinya mereka **tidak menunggu penyimpanan akhir** (ronde 1 §3.7):

| Rule | Tabel |
| --- | --- |
| `GCNM!INSERTCLAIMREJECTED_SQL` | tabel penolakan klaim |
| `ASM!INSERTHISTORYAKSEPTASIPEGA_SQL` | tabel riwayat akseptasi |
| `RNM!INSERTLOGMOP_SQL` | log kiriman ke Kasir |
| `RNM!INSERTLOGSERVICECLAIM` | log pemantauan klaim |
| `RNM!GENERATEIMAGEID_SQL` | tabel penyimpanan berkas |
| `RNM!GETTANGGALCLOSING_SQL` | blok prosedur |
| `RNM!INSERT_T_STORAGE_SQL` | tabel penyimpanan berkas |
| `GCNM!INSERTCLAIMPNC` | blok prosedur |
| `GCNM!SAVEDATATOOSAKSEPTASI` · `ASM!SAVEOSCLAIM_SQL` | daftar akseptasi |

`[keputusan work owner]` 2026-09-18 — **"ditiru apa adanya."** Kesembilannya **tetap menyimpan
sendiri**.

> ⚠️ **Risiko yang diterima sadar:** bila langkah sesudahnya gagal, **kesembilan tabel tetap
> terisi**, dan **tidak ada rule pembatal di korpus** — tidak ada penghapusan, pembatalan, maupun
> penanda batal di satu pun dari 23 rule yang disisir. **Data separuh jadi adalah perilaku yang
> disengaja ditiru, bukan cacat yang terlewat.**

#### Dua tabel yang bentuknya sudah pasti

`[data DBA]` 2026-09-18 — keduanya di skema yang sama, dan penulisan tanpa awalan skema di rule
terselesaikan lewat skema bawaan koneksi (ronde 3 §12).

**Tabel riwayat akseptasi — 7 kolom, jalur komite mengisi 6:** pengenal kasus · tanggal transfer ·
status · nama pengguna · kotak kerja · pengenal komite — dan **satu kolom akun operator yang
TIDAK pernah diisi jalur komite**.

⭐ **Pengenal kasus bukan nomor polos** — ia **kunci instance penuh berikut awalan kelas**, panjang
150 karakter. Dua dari lima indeksnya memotong awalan kelas, yang membuktikan ada pihak lain yang
mencari dengan bentuk terpotong. ⚠️ **Fakta untuk pengembang, bukan usulan:** ia **teks**, bukan
angka, dan **bukan kunci asing basis data**.

**Tabel penyimpanan berkas — 8 kolom**, kunci utamanya pengenal berkas, **nol kunci asing**.
⚠️ Dua kolomnya menunjukkan tabel ini **dipakai bersama banyak aplikasi**, bukan milik satu modul.

⛔ **Bentuk, kolom, tipe, indeks, dan batasan kedua tabel mengikuti apa yang sudah ada di produksi,
dan tidak ditetapkan ulang oleh rancangan mana pun di modul ini.**

#### Sembilan kolom tabel penyetuju yang sudah dikunci

`[keputusan work owner]` — kesembilan kolom berikut sudah ditetapkan, ditulis apa adanya
(ronde 5 §C4):

| Kolom | Isinya |
| --- | --- |
| `ID` | kunci baris |
| `DATA_KOMITE_ID` | penunjuk ke kasus komitenya |
| `KOMITE_URUT` | penyetuju ke berapa |
| `KOMITE_OPERATORID` | akun operator penyetuju |
| `KOMITE_JABATAN` | jabatan penyetuju |
| `KOMITE_EMAIL` | email penyetuju |
| `KOMITE_APPROVAL` | keputusan |
| `KOMITE_COMMENT` | catatan |
| `DATE_APPROVE` | tanggal diputuskan |

`[keputusan work owner]` 2026-09-18 — **satu kolom tanggal saja, `DATE_APPROVE`.** Di Pega ada dua
properti tanggal yang diisi **langkah yang sama** dengan nilai identik; yang kedua **tidak pernah
ditampilkan**. Yang kedua **sengaja tidak dijadikan kolom** — jangan menambahkannya kembali karena
"ada di korpus".

> ⚠️ **19 properti baris penyesuaian milik kasus komite BELUM tertampung di kesembilan kolom itu.**
> `[terverifikasi]` (ronde 5 §C2). Kesembilan kolom di atas adalah **satu baris penyetuju**;
> properti baris penyesuaian dan isi kasus komite **tempatnya di tabel lain**.
>
> ⛔ **Didaftarkan saja. Tidak ada usulan kolom baru di berkas ini — itu keputusan work owner.**

#### ⭐ Dua kolom penanda usul di header kasus komite — nilainya sudah ditetapkan

`[keputusan work owner]` 2026-09-19 — **kedua penanda usul disimpan sebagai teks satu huruf:
`'1'` bila diusulkan, `'0'` bila tidak.** Kolomnya **`CHAR(1)`** dan **wajib isi** — kotak centang
yang **tidak pernah disentuh** tersimpan **`'0'`**, bukan kosong.

⚠️ **Migrasi data lama:** nilai yang di Pega bermakna **benar/salah** dikonversi otomatis menjadi
**`'1'`/`'0'`**, dan **baris lama yang nilainya kosong menjadi `'0'`**.

⛔ **Butir `[terbuka]` tentang nilai kedua penanda usul dengan demikian TERJAWAB.** Ia tidak pernah
masuk register 13 butir, jadi jumlah register **tidak berubah karenanya**. Nama dan tempat kolomnya
tetap mengikuti **acuan tunggal** `STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md` — berkas ini hanya
menetapkan **nilai dan tipenya**.

#### `[terbuka]` bab ini

- **Kolom akun operator di tabel riwayat akseptasi tidak pernah diisi jalur komite.** Adakah
  penulis lain di luar modul ini, atau kolom itu memang selalu kosong. ⛔ Tidak ditebak, dan
  **tidak diusulkan dibuang**.

---

### Bab 10 — Kumpulan keputusan work owner

*Bab ini mengumpulkan seluruh keputusan yang diucapkan work owner selama ronde 1–8, dengan kalimat
aslinya bila ada, supaya tidak perlu mencarinya satu per satu di delapan berkas.*

Keputusan **1–20** bertanggal **2026-09-18**; keputusan **21–28** bertanggal **2026-09-19**.

> ⛔ **RALAT 2026-09-19 — ketiga.** Kalimat lamanya **dikutip, tidak dihapus**: *"Keputusan
> **1–20** bertanggal **2026-09-18**; keputusan **21–25** bertanggal **2026-09-19**."* Tiga
> keputusan lagi masuk pada hari yang sama, dan **satu keputusan lama dicabut**.

> ⛔ **RALAT 2026-09-19 — kedua**. Kalimat lamanya **dikutip berurutan, tidak satu pun dihapus**:
>
> > *"Seluruhnya bertanggal **2026-09-18**."*
> >
> > *"Seluruhnya bertanggal **2026-09-18**, **kecuali keputusan 21 dan 22 yang bertanggal
> > 2026-09-19**."*
> >
> > *(catatan ralat ronde sebelumnya, ikut dikutip)* — *"⛔ **RALAT 2026-09-19** — kalimat
> > lamanya **dikutip, tidak dihapus**: "Seluruhnya bertanggal **2026-09-18**." Ia benar
> > sampai keputusan 20; dua keputusan berikutnya bertanggal lebih baru."*
>
> Yang pertama benar sampai keputusan 20; yang kedua benar sampai keputusan 22. Tiga keputusan
> lagi masuk pada hari yang sama.
>
> ⚠️ **Kalimat pengantar bab ini juga sudah tidak tepat.** Ia berbunyi *"seluruh keputusan yang
> diucapkan work owner selama **ronde 1–8**"* — padahal keputusan **21–25** lahir **sesudah**
> delapan ronde itu, saat spec dan tiket sedang dirapikan. ⛔ Kalimatnya **tidak diubah**;
> perluasannya dicatat di sini saja.

| # | Keputusan | Kalimat asli / inti |
| --- | --- | --- |
| **1** | **Lingkup satu modul** | *"intinya ini aku mau supaya ini hanya bahas Komite Claim Prop."* Komite Life / Non Prop / Fac In digarap sendiri-sendiri nanti. ⚠️ Risiko rancangan berselisih antar lini **diterima sadar**. |
| **2** | **Nama yang sama bukan rule yang sama** | *"itu hanya nama yang sama. Ada nama yang sama tapi tetap menjalankan sesuai modulnya masing-masing."* |
| **3** | **Sembilan rule tetap menyimpan sendiri** | *"ditiru apa adanya."* Termasuk tabel penolakan klaim dan tabel riwayat akseptasi. |
| **4** | **Pencacah penyetuju adalah kolom yang disimpan** | *"ikuti activity yang berjalan."* Diisi ulang di titik yang sama dengan Pega, **bukan** dihitung ulang tiap dibaca. |
| **5** | **Kurs mengikuti query yang ada** | *"ikuti aja query di xml CurrencyStandard."* Dipanggil apa adanya dengan dua masukan yang sama; isi fungsinya **tidak diminta ke DBA**. |
| **6** | **Satu kolom tanggal persetujuan** | `DATE_APPROVE` saja; yang kedua **sengaja tidak dijadikan kolom**. |
| **7** | **Satu properti jenis treaty dibuang** | Ia hanya ada di balik gerbang mati, tidak pernah tampil. **Tidak dibuat di sistem baru.** |
| **8** | **Nomor klaim dan dua kotak total estimasi di balik `NEVER` dibuang** | **Tidak dibuat di layar baru.** |
| **9** | **Lima perintah lompat menggantung dibuang** | Kelimanya menunjuk tanda yang tidak ada, dan kelimanya bergerbang mati. ⛔ **Tidak ditebak ke mana seharusnya melompat. Tidak ditandai cacat. Dibuang.** |
| **10** | **Aturan awalan** | titik dan `pyWorkPage.` = komite · `pyWorkCover.` = klaim induk, **dibaca saja**. |
| **11** | **Jalur REJECT dan CLOSE: dua isian** | keputusan dan catatan. **Ditiru apa adanya.** |
| **12** | **Gagal kirim ke Kasir: email tetap jalan** | kasus tetap jalan, nol percobaan ulang. **Ditiru apa adanya**, dengan risiko diterima sadar. |
| **13** | **Penanda konversi adalah syarat resmi** | kiriman ke Kasir hanya boleh jalan bila nomor akseptasi sudah tercatat. **Dipertahankan.** |
| **14** | ⭐ **Urutan efek keluar tidak ditiru** | dicatat sebagai fakta, **tidak mengikat rancangan**; disesuaikan nanti. |
| **15** | **Hak akses per activity diabaikan** | tidak dipindahkan dan tidak disisir. ⚠️ Didukung temuan: dari 35 baris hak akses, yang **kelasnya** terisi 34, yang **nama haknya** terisi **nol** — daftar itu menyebut kelas tanpa menyebut hak `[terverifikasi]` (ronde 7 §A2). |
| **16** | **Titik lompat darurat dibuang** | `komiteAccept_ticket` sudah tidak dipakai; jalur mati. |
| **17** | ⭐ **Ketelitian angka** — **diralat** | hitungan **sampai 20 angka di belakang koma tanpa pembulatan di tengah jalan** · **kolom basis data MENGIKUTI DDL YANG SUDAH ADA — 20 digit, 8 di belakang koma** · tampilan **4 angka** · pencacah tetap bilangan bulat · tipe desimal di Go. ⚠️ **Pembulatan di batas penyimpanan, diterima sadar.** ⛔ **Yang dicabut, dikutip sebagai jejak:** *"18 digit di depan koma, 20 di belakang"*, bertanda *atas rekomendasi asisten*. **Lihat bab 7.** |
| **18** | **Seam memakai ulang Claim Prop** | modul komite **tidak membangun seam sendiri**. ⚠️ Akibatnya pengujian modul ini **bergantung pada seam Claim Prop sudah ada lebih dulu**. |
| **19** | **Efek keluar diuji dengan layanan sungguhan** | bukan pengganti tiruan, dan dijalankan di **lingkungan uji terpisah, bukan produksi**. |
| **20** | **Pengujian urutan efek keluar menunggu** | selama urutan barunya belum ditetapkan, bagian itu **tidak punya sasaran uji**. |
| **21** | ⭐ **Nilai dua penanda usul** *(2026-09-19)* | disimpan sebagai **teks satu huruf** — `'1'` bila diusulkan, `'0'` bila tidak. Kolomnya **`CHAR(1)`**, **wajib isi**; kotak yang tidak disentuh tersimpan **`'0'`**, bukan kosong. ⚠️ Migrasi: benar/salah → `'1'`/`'0'`, kosong → `'0'`. **Lihat bab 9.** |
| **22** | ⚠️ **Baris penyesuaian yang sudah diserahkan ke komite menjadi BEKU** *(2026-09-19)* | tidak dapat diubah dan tidak dapat dihapus, **selamanya**; **penolakan komite tidak mencairkannya**. Perbaikan lewat **baris penyesuaian baru**. ⛔ **Perilaku BARU, bukan paritas** — Pega tidak punya kunci semacam ini. Lingkup: **Claim Prop** *(Claim Non Prop menyusul; Life tidak termasuk)*. Penegakannya **milik modul Claim Prop**. **Lihat bab 11.** ⛔ **DIPERSEMPIT keputusan 24** — kata *"selamanya"* di baris ini hanya berlaku untuk **sunting**, tidak untuk **hapus klaim**. Kalimat di atas **dibiarkan apa adanya sebagai jejak**. |
| **23** | ⭐ **Berkas `Struktur_*.xlsx` di folder korpus bukan sumber kebenaran** *(2026-09-19)* | dibuat **tim sendiri** untuk memetakan XML — **artefak turunan**. ⛔ **Tidak boleh dikutip `[terverifikasi]`.** ⚠️ Setiap sensus korpus menyebut penyebutnya **"berkas `.xml`"**, bukan "berkas". **Sudah diuraikan di bab "Cara membaca berkas ini".** |
| **24** | ⛔ **DICABUT keputusan 28, 2026-09-19** — *sebabnya: larangan hapus klaim ternyata **permanen**, bukan sementara, sehingga premis keputusan ini gugur seluruhnya.* ⚠️ Teksnya **tidak dihapus**, dikutip apa adanya sebagai jejak: ⚠️ **Kunci beku berlaku sampai komite selesai, bukan selamanya** *(2026-09-19)* | baris beku **tidak dapat disunting selamanya** dan **tidak dapat dihapus satu per satu**; **klaim induknya terkunci hanya selama kasus komitenya berjalan**. Sesudah komite selesai atau menolak, hapus klaim **boleh** dan **mengkaskade seperti biasa**. ⛔ **Mempersempit keputusan 22.** **Lihat AC 76 · 77 · 78 dan bab 11.** |
| **25** | ⚠️ **Penegakan kunci beku dibagi di modul Claim Prop** *(2026-09-19)* | tiket **11 MEMASANG** — menandai baris menjadi beku saat penyerahan ke komite tercatat; tiket **08 MENEGAKKAN** — menjadikannya **sifat baris** yang menahan **setiap** jalur sunting; **jalur kaskade hapus di tiket 00** modul itu **disesuaikan**. ⛔ **Modul ini hanya pemicunya — nol penegakan di sini.** ⚠️ **Kata "memasang" dipersempit keputusan 26** — yang dipasang bukan kolom penanda, melainkan **kasus komitenya sendiri**. |
| **26** | ⭐ **Baris beku dikenali sebagai TURUNAN — nol kolom penanda baru** *(2026-09-19)* | sebuah baris penyesuaian **dianggap beku bila ada kasus komite yang menunjuknya** — **berjalan maupun selesai**. ⛔ **Nol kolom baru pada baris penyesuaian.** Penegakannya di modul **Claim Prop, tiket 08**. |
| **27** | ⚠️ **Kolom `FLAG_ON_GOING_COMMITTEE` pada `T_GENERAL_CLAIM` DIBUANG** *(2026-09-19)* | tidak dibawa ke sistem baru; **layar akseptasi menghitung sendiri** *"komite sedang berjalan"* dari kasus komitenya. ⚠️ **Penyimpangan sadar:** `[terverifikasi]` kolom itu **ada di Pega** — ditulis `Activity/AddKomiteTreatyChild_ACT.xml` langkah **3**, juga `KomitePostAdjustment` langkah **26.1** dan **27**, dan **dibaca layar akseptasi** — tetapi **sengaja tidak dibawa** karena nilainya **dapat diturunkan**. Dibuang di modul **Claim Prop, tiket 00**. **Lihat bab 11, titik 4.** |
| **28** | ⚠️ **Klaim yang PERNAH punya kasus komite tidak dapat dihapus — SELAMANYA** *(2026-09-19)* | walau komitenya **sudah selesai atau menolak**. Klaim yang **belum pernah** ke komite **tetap dapat dihapus** dan **mengkaskade**. ⛔ **MENCABUT keputusan 24.** **Lihat AC 77.** |
| ⭐⭐ **30** | ⚠️ **Keputusan penyetuju hanya boleh disimpan pemilik `KomiteID` tingkat berjalan** *(2026-09-19)* | **ADR-0014** mewajibkannya, dan modul ini **belum pernah menyebutnya** — ke-15 ADR baru diadu dengan modul ini hari ini. `[terverifikasi]` Di Pega **nol pemeriksaan semacam itu**: `KomiteID` hanya muncul di `KomiteRouter`, dan di sana ia **menetapkan tujuan rute**, bukan menguji penyimpan. `AcceptStatus` ditulis **tepat satu** penugasan properti di seluruh 80 berkas. ⭐ **Alasan pokoknya bukan "ADR bilang begitu":** di Pega kotak kerja membatasi **secara kebetulan** — orang lain tidak melihat kasusnya, jadi jarang menyimpannya. **Di sistem baru kebetulan itu tidak ada**, sehingga tanpa pemeriksaan eksplisit lubangnya **MELEBAR, bukan pindah**. ⚠️ **Penyimpangan sadar** — preseden **Komite Claim Life** sudah ada pada ADR yang sama. ⭐⭐ **Dan preseden yang lebih dekat lagi: modul CLAIM PROP sudah menegakkannya** — **AC 57** modul itu berbunyi *"Seluruh gerbang wewenang **ditegakkan di lapisan layanan** (**ADR-0014**)"* dan **AC 58** mengujinya lewat endpoint tanpa UI. ⚠️ **Modul komite-lah yang tertinggal**, bukan ADR-nya yang asing bagi keluarga modul ini. **Lihat AC 81 · 82, AC 67, bab 11 titik 6.** |
| **29** | ⭐ **Daftar penyetuju DITAMPILKAN di layar komite** *(2026-09-19)* | urutan · siapa · dan untuk yang **sudah** memutuskan: **keputusan, catatan, tanggal**. ⚠️ **Penyimpangan sadar** — `[terverifikasi]` di Pega **tidak satu pun dari ketiganya tampil di layar**, sehingga penyetuju memutuskan **tanpa konteks penyetuju sebelumnya**. Menutup **US 3** dan **US 4**, yang sebelumnya **nol AC**. ⛔ Tandanya **`[keputusan work owner — atas rekomendasi asisten]`** — pilihannya diserahkan kepada asisten, **mudah dicabut**. **Lihat bab 4, bab 11 titik 5, AC 79 · 80.** |

---

### Bab 11 — Titik yang **SENGAJA DIUBAH** dari Pega

*Bab ini memisahkan hal yang sengaja berbeda dari hal yang ditiru, supaya tidak tercampur dan tidak
disalahpahami sebagai cacat pemindahan.*

⭐ **Ada enam.**

> ⛔ **RALAT 2026-09-19 — keempat.** Keempat kalimat lamanya **dikutip berurutan, tidak satu pun
> dihapus**:
>
> > *"⭐ **Ada dua, dan keduanya besar.**"* · *"⭐ **Ada tiga.**"* · *"⭐ **Ada empat.**"* ·
> > *"⭐ **Ada lima.**"*
>
> Titik **keenam** — **wewenang menyimpan keputusan ditegakkan** — ditambahkan
> `[keputusan work owner]` 2026-09-19, menutup lubang **ADR-0014**.

> ⛔ **RALAT 2026-09-19 — ketiga.** Ketiga kalimat lamanya **dikutip berurutan, tidak satu pun
> dihapus**:
>
> > *"⭐ **Ada dua, dan keduanya besar.**"* · *"⭐ **Ada tiga.**"* · *"⭐ **Ada empat.**"*
>
> Titik **kelima** — **daftar penyetuju yang ditampilkan** — ditambahkan
> `[keputusan work owner — atas rekomendasi asisten]` 2026-09-19, menutup **US 3** dan **US 4**.

> ⛔ **RALAT 2026-09-19 — kedua.** Kedua kalimat lamanya **dikutip berurutan, tidak satu pun
> dihapus**:
>
> > *"⭐ **Ada dua, dan keduanya besar.**"*
> >
> > *"⭐ **Ada tiga.**"*
>
> Titik **ketiga** — **baris penyesuaian yang beku** — dan titik **keempat** — **satu kolom Pega
> yang sengaja dibuang** — keduanya ditambahkan `[keputusan work owner]` 2026-09-19.
>
> *(anotasi ralat ronde sebelumnya, ikut dikutip)* — *"⛔ **RALAT 2026-09-19** — kalimat lamanya
> **dikutip, tidak dihapus**: "⭐ **Ada dua, dan keduanya besar.**" Titik ketiga — **baris
> penyesuaian yang beku** — ditambahkan `[keputusan work owner]` 2026-09-19."*

| # | Titik | Di Pega | Di sistem baru | Bab |
| --- | --- | --- | --- | --- |
| **1** | **Urutan efek keluar terhadap penyimpanan** | kedelapan efek keluar berjalan **sebelum** penyimpanan | **akan disesuaikan** — urutan penggantinya **belum diputuskan** | bab 8 |
| **2** | **Ketelitian angka** | **empat rupa berbeda**, nol pembulatan di mana pun | hitungan sampai **20 angka** tanpa pembulatan di tengah jalan · **satu titik pembulatan di batas penyimpanan** · tampilan **4 angka** | bab 7 |
| **3** | ⭐ **Baris penyesuaian yang sudah diserahkan ke komite menjadi BEKU** | ⛔ **tidak ada kunci semacam ini** — baris boleh berubah atau hilang sesudah diserahkan, tanpa satu pun pemeriksaan | **tidak dapat disunting selamanya**; **tidak dapat dihapus satu per satu**; **klaim induknya tidak dapat dihapus sama sekali, selamanya**. Penolakan komite **tidak mencairkan apa pun**. Perbaikan lewat **baris penyesuaian baru** | bab 2 · bab 11 |
| **4** | ⭐ **Satu kolom Pega sengaja DIBUANG — `FLAG_ON_GOING_COMMITTEE`** | `[terverifikasi]` kolomnya **ada dan diisi** — `AddKomiteTreatyChild_ACT` langkah **3**, `KomitePostAdjustment` langkah **26.1** dan **27** — serta **dibaca layar akseptasi** | **tidak dibawa**; nilainya **diturunkan** dari kasus komitenya saat dibutuhkan | bab 9 · bab 11 |
| ⭐⭐ **6** | ⚠️ **Wewenang MENYIMPAN keputusan ditegakkan** | ⛔ **tidak ditegakkan sama sekali** — Pega hanya **menempatkan** kasus di kotak kerja pemilik `KomiteID`; siapa pun yang dapat membuka layar **dapat menyimpan keputusan** untuk tingkat mana pun | **hanya pemilik `KomiteID` tingkat berjalan** yang dapat menyimpan; yang lain **ditolak di lapisan layanan** dengan galat yang **terlihat**. ⛔ Bukan disembunyikan layarnya — **ditolak** | bab 3 · bab 10 · **ADR-0014** |
| **5** | ⭐ **Daftar penyetuju DITAMPILKAN di layar** | ⛔ **tidak tampil sama sekali** — `KomiteLoop`, `KomiteCount`, dan `KomiteList` ketiganya hanya mesin; penyetuju memutuskan **tanpa melihat siapa pun sebelum dia** | **tampil**: urutan · siapa · dan untuk yang sudah memutuskan **keputusan, catatan, tanggal**. ⛔ **Nol kolom baru** — datanya sudah ada di baris penyetuju | bab 3 · bab 4 |

#### ⚠️ Titik 3 — baris penyesuaian yang beku

`[keputusan work owner]` 2026-09-19 — ⚠️ **baris penyesuaian yang sudah dikirim meminta persetujuan
komite menjadi BEKU.** Bekunya **satu umur saja — selamanya**, pada ketiga sisinya:

| Sisi | Berlaku sampai kapan |
| --- | --- |
| **tidak dapat disunting** | **selamanya** — **bekunya tidak mencair meskipun komite menolak** |
| **tidak dapat dihapus satu per satu** | **selamanya** |
| **klaim induknya tidak dapat dihapus** | **selamanya** — walau komitenya sudah selesai maupun menolak |

Klaim yang **belum pernah** punya kasus komite **tetap dapat dihapus** dan **mengkaskade seperti
biasa**. Perbaikan atas isi baris tetap dilakukan dengan **baris penyesuaian baru**, bukan dengan
menyunting baris lama.

⭐ **Bekunya TURUNAN, bukan kolom.** `[keputusan work owner]` 2026-09-19 — sebuah baris penyesuaian
**dianggap beku bila ada kasus komite yang menunjuknya**, **berjalan maupun selesai**. ⛔ **Nol
kolom penanda baru** pada baris penyesuaian; yang menjadi buktinya adalah **keberadaan kasus
komitenya sendiri**.

⚠️ **Sunting dan hapus diperlakukan SAMA:** yang dijaga adalah **jejak baris yang pernah dinilai
komite**, dan jejak itu tidak boleh hilang lewat **jalan mana pun** — baik dengan menyunting
barisnya, menghapus barisnya, maupun menghapus klaim induknya.

> ⛔ **RALAT 2026-09-19 — kedua.** Uraian ini semula menggambarkan kunci **selamanya di kedua
> sisi**, lalu **dua sisi yang umurnya berbeda**. Keduanya salah. Kalimat-kalimat lamanya
> **dikutip utuh, tidak satu pun dihapus**:
>
> > *"Bekunya punya **dua sisi yang umurnya berbeda**"* — dengan baris tabel
> > *"**klaim induknya tidak dapat dihapus** | **hanya selama kasus komitenya berjalan**"*.
> >
> > *"Sesudah komite **selesai atau menolak**, menghapus klaim **boleh** dan **mengkaskade seperti
> > biasa** — termasuk **menyapu baris beku itu**."*
> >
> > *"⚠️ **Sunting dan hapus sengaja diperlakukan berbeda:** yang dijaga selamanya adalah **isi
> > baris yang sudah pernah dinilai komite**; yang tidak dijaga selamanya adalah **hak menghapus
> > klaimnya**, karena mematikan hapus-klaim akan mematikan fitur yang sah."*
> >
> > *(dari sel tabel titik 3)* — *"Penolakan komite **tidak mencairkan hak sunting**, tetapi
> > **membuka hak hapus klaim**."*
>
> `[keputusan work owner]` 2026-09-19: **klaim yang PERNAH punya kasus komite tidak dapat dihapus,
> selamanya.** Lihat keputusan **28** di bab 10, yang **mencabut keputusan 24**.

> ⛔ **RALAT 2026-09-19** — uraian ini semula menggambarkan kunci yang berlaku **selamanya di
> kedua sisi**. Kedua kalimat lamanya **dikutip utuh, tidak dihapus**:
>
> > *"`[keputusan work owner]` 2026-09-19 — ⚠️ **baris penyesuaian yang sudah dikirim meminta
> > persetujuan komite menjadi BEKU: tidak dapat diubah dan tidak dapat dihapus, selamanya.
> > Bekunya tidak mencair meskipun komite menolak** — perbaikan dilakukan dengan **baris
> > penyesuaian baru**, bukan dengan menyunting baris lama."*
> >
> > *(dari baris tabel titik 3)* — *"**beku selamanya**: tidak dapat diubah, tidak dapat dihapus;
> > **penolakan komite tidak mencairkannya**. Perbaikan lewat **baris penyesuaian baru**"*
>
> ⚠️ Bentuk lama **menutup kaskade hapus klaim sepenuhnya**, dan itu bertabrakan dengan
> `claim-prop/spec.md` **AC 5** — *"Menghapus klaim **mengkaskade sampai tingkat terdalam**;
> test wajib memeriksa cicit."* `[keputusan work owner]` 2026-09-19: **kunci beku tidak mematikan
> fitur hapus klaim; yang dikunci adalah saat komitenya masih berjalan.** **Lihat keputusan 24 di
> bab 10.**

⚠️ **Pega tidak punya kunci semacam ini.** `[terverifikasi]` Di seluruh ekspor modul ini:
**nol pemeriksaan, nol pesan kesalahan, nol penanganan** untuk baris penyesuaian yang berubah atau
hilang sesudah diserahkan. **Karena itu aturan ini dibangun, bukan dimigrasikan.**

**Lingkup keputusan:** Claim Prop dan Claim Non Prop. `[keputusan work owner]` **Claim Non Prop
dikerjakan menyusul — fokus sekarang Claim Prop saja. Claim — Life tidak termasuk.**

⚠️ **Kuncinya dibangun di modul Claim Prop**, karena daftar penyesuaian miliknya. **Modul ini
hanya pemicunya:** saat penyerahan ke komite tercatat, baris induknya menjadi beku.

`[keputusan work owner]` 2026-09-19 — **penegakannya dibagi dua di modul itu:** tiket **11
MEMASANG** *(menandai baris menjadi beku saat penyerahan tercatat)* · tiket **08 MENEGAKKAN**
*(menjadikannya sifat baris yang menahan setiap jalur sunting)*, dan **jalur kaskade hapus di
tiket 00** modul itu **disesuaikan**. ⛔ **Nol penegakan di modul ini.**

⛔ **Ia tidak berlaku surut.** Kasus komite **lama hasil migrasi** tidak terlindungi olehnya —
lihat butir **14** di register `[terbuka]`.

#### Kenapa ketelitian angka **TETAP di sini**, meski kolomnya mengikuti DDL yang sudah ada

⛔ **Ia tetap titik yang sengaja diubah, bukan "ditiru apa adanya".** Alasannya satu kalimat:
**yang mengikuti DDL lama hanyalah BENTUK KOLOMNYA; CARA MENGHITUNGNYA berubah** — Pega memakai
empat rupa ketelitian yang berbeda-beda tanpa pembulatan sama sekali, sedangkan sistem baru
memakai **satu** aturan dengan **satu titik pembulatan di batas penyimpanan**.

⚠️ **Akibat titik 2:** hasil hitungan sistem baru **akan berbeda dari Pega pada angka di belakang
koma**. Itu **disengaja**, bukan cacat. Khususnya: **angka persen yang di Pega dihitung sampai 10
desimal akan menjadi 8 saat disimpan.**

#### ⚠️ Titik 4 — kolom `FLAG_ON_GOING_COMMITTEE` yang sengaja dibuang

`[keputusan work owner]` 2026-09-19 — ⚠️ **kolom `FLAG_ON_GOING_COMMITTEE` pada `T_GENERAL_CLAIM`
tidak dibawa ke sistem baru.** **Layar akseptasi menghitung sendiri** *"komite sedang berjalan"*
dari **kasus komitenya**, bukan dari kolom tersimpan.

`[terverifikasi]` **Kolom itu ada di Pega dan benar-benar dipakai** — ditulis
`Activity/AddKomiteTreatyChild_ACT.xml` langkah **3**, ditulis juga `KomitePostAdjustment` langkah
**26.1** dan **27**, dan **dibaca layar akseptasi**. ⛔ **Karena itu ia penyimpangan sadar, bukan
jalur mati yang dibuang** — ia hidup, dan tetap tidak dibawa.

**Alasannya satu kalimat:** nilainya **dapat diturunkan** dari kasus komitenya, jadi menyimpannya
berarti memelihara dua sumber kebenaran untuk satu fakta. Sejalan dengan **titik 3**, yang juga
mengenali beku sebagai **turunan**, bukan sebagai kolom.

⚠️ **Dibuangnya di modul Claim Prop, tiket 00** — tabelnya milik modul itu. ⛔ **Nol pekerjaan
di modul ini**; ia dicatat di sini karena **pemicunya** ada di sini.

⚠️ **Bentuk kolom dan nama tabel tetap mengikuti acuan tunggalnya**,
`claim-prop/STRUKTUR-TABEL-CLAIM-PROP.md`. Berkas ini **hanya mencatat bahwa kolom itu tidak
dibawa**.

#### ⭐ Titik 5 — daftar penyetuju yang ditampilkan

`[keputusan work owner — atas rekomendasi asisten]` 2026-09-19 — **layar komite menampilkan daftar
penyetuju**: **urutan**, **siapa**, dan untuk yang **sudah** memutuskan — **keputusan, catatan, dan
tanggalnya**.

⚠️ **Di Pega ia tidak tampil sama sekali.** `[terverifikasi]` Ketiga penggerak tangga —
`KomiteLoop`, `KomiteCount`, `KomiteList` — **nol di antaranya tampil di layar** (bab 3), dan
penyaringan ke-93 medan layar **tidak menemukan satu pun** yang menampilkan daftar itu. **Akibatnya
di Pega: setiap penyetuju memutuskan tanpa konteks penyetuju sebelumnya.**

⛔ **Nol kolom basis data baru.** Yang ditampilkan sudah tersimpan di baris penyetuju — kesembilan
kolom yang dikunci di bab 9. Ini **murni penambahan tampilan**.

**Alasannya ditulis lengkap di bab 4**, berikut sebab tandanya **`atas rekomendasi asisten`** dan
apa saja yang tersentuh bila kamu mencabutnya. **Lihat AC 79 dan AC 80.**

⛔ **Ia menutup US 3 dan US 4**, dua user story yang sebelumnya **nol AC**.

**Tidak ada titik keenam.** Seluruh keputusan lain di bab 10 berbunyi *"ditiru apa adanya"*,
*"dibuang karena jalur mati"*, atau menyangkut **cara menguji**, bukan perilaku — tidak satu pun
mengubah perilaku yang hidup.

> ⛔ **RALAT 2026-09-19** — kalimat lamanya **dikutip utuh, tidak dihapus**:
>
> > *"`[terverifikasi]` **Tidak ada titik ketiga.** Seluruh keputusan lain di bab 10 berbunyi
> > *"ditiru apa adanya"*, *"dibuang karena jalur mati"*, atau menyangkut **cara menguji**,
> > bukan perilaku — tidak satu pun mengubah perilaku yang hidup."*
>
> Titik ketiga kini **ada** — baris penyesuaian yang beku, dan **titik keempat** menyusul pada
> hari yang sama — kolom `FLAG_ON_GOING_COMMITTEE` yang sengaja dibuang. Kalimat *"Tidak ada
> titik **keempat**"* karenanya **juga sudah usang**, dan diganti *"Tidak ada titik **kelima**"*
> — bentuk lamanya dikutip di sini.
>
> ⛔ **Dan sekali lagi, pada hari yang sama:** titik **kelima** — **daftar penyetuju yang
> ditampilkan** — membuat *"Tidak ada titik **kelima**"* ikut usang; ia diganti *"Tidak ada titik
> **keenam**"*, dan bentuk lamanya **dikutip di sini**. ⚠️ Rantai ralat ini panjang justru karena
> **nol kalimat lama dihapus**. Tanda `[terverifikasi]` juga dilepas dari kalimat ini:
> ketiadaan titik berikutnya adalah **hasil penyisiran bab 10**, bukan pembacaan korpus.
>
> ✅ **AC 72 SUDAH ikut diralat** `[keputusan work owner]` 2026-09-19 — ia kini berbunyi
> *"Tidak ada penyimpangan di luar yang terdaftar di Bab 11"*, **tanpa angka** dan **tanpa tanda
> golongan**. Nomornya **tetap 72**.
>
> > ⛔ **Catatan lamanya dikutip, tidak dihapus:** *"⚠️ **AC 72 belum ikut diralat.** Ia masih
> > berbunyi *"Tidak ada penyimpangan ketiga"*. ⛔ Nomor dan makna AC 1–73 dikunci oleh perintah
> > kerja, jadi AC 72 **tidak disentuh di sini** — ralatnya **keputusan work owner**, bukan
> > keputusan berkas ini."*
>
> Ralatnya memang datang dari work owner, pada hari yang sama.

---

## Testing Decisions

### Apa yang membuat test baik di sini

Menguji **perilaku yang terlihat dari luar**, bukan cara kerja di dalamnya. Untuk modul ini artinya:
*"penyetuju terakhir menyetujui → nomor akseptasi terbit sekali, dan data pembayaran terkirim"* —
**bukan** *"activity X memanggil rule Y di langkah 34"*. Nomor langkah Pega adalah **bukti asal
usul**, bukan sasaran uji.

### Seam — **memakai ulang seam Claim Prop**, tidak menambah

`[keputusan work owner]` 2026-09-18 — **modul komite MEMAKAI ULANG seam modul induknya, Claim
Prop. Tidak membangun seam sendiri.** Sama dengan pilihan yang sudah diambil modul Komite Claim
Life.

> ⚠️ **Akibatnya, ditulis apa adanya:** pengujian Komite Claim Prop **bergantung pada seam Claim
> Prop sudah ada lebih dulu**. Selama seam itu belum berdiri, modul ini **belum bisa diuji**
> — bukan karena bahannya kurang, melainkan karena pintunya belum dibuat.

Dari sisi perilaku lama, bentuknya memang cocok: `[terverifikasi]` seluruh perilaku modul ini
masuk lewat **satu pintu** — penyetuju menekan kirim pada layar komite — dan keluar lewat
**delapan efek keluar** yang sudah terdaftar di bab 8.

### Cara menguji delapan efek keluar

`[keputusan work owner]` 2026-09-18 — **kedelapan efek keluar diuji dengan LAYANAN SUNGGUHAN,
bukan pengganti tiruan** — dan dijalankan di **LINGKUNGAN UJI TERPISAH, bukan produksi**.

⚠️ Alamat lingkungan uji itu ditemukan lewat **tabel alamat layanan** (bab 8) dengan **pasangan
kategori sendiri**; **barisnya mana BELUM ditetapkan** → `[terbuka]`.

### Menguji urutan efek keluar

`[keputusan work owner]` 2026-09-18 — **pengujian urutan efek keluar MENUNGGU urutan barunya
ditetapkan.** Selama belum, bagian itu **tidak punya sasaran uji** → `[terbuka]`.

⛔ Ini bukan kekurangan berkas ini: urutan barunya memang **belum diputuskan** (bab 11).

### Modul yang diuji

**Mengikuti modul induk Claim Prop** — konsekuensi langsung dari keputusan seam di atas. Modul
komite tidak menambah modul uji tersendiri.

### Prior art

Modul `komite-claim-life` sudah mengambil pilihan yang sama — **memakai ulang seam modul
induknya, tidak menambah**. `claim-prop` dan `premiumlist-life` punya bagian yang sama dan bisa
dijadikan pembanding bentuk.

### `[terbuka]` bab ini — **dua**

1. **Baris mana di tabel alamat layanan yang menunjuk lingkungan uji** — pasangan kategorinya
   belum ditetapkan.
2. **Sasaran uji untuk urutan efek keluar** — menunggu urutan barunya diputuskan.

---

## Acceptance Criteria

*Bab ini **tidak memutuskan apa pun.** Ia menyatakan ulang — dalam bentuk yang bisa diuji dari
luar — keputusan yang **sudah** ada di Bab 1–11, di tiket, dan di register butir `[terbuka]`.
Bila sebuah butir di bawah terasa seperti keputusan baru, ia salah tulis: laporkan, jangan
dilaksanakan.*

**Cara membacanya.** Setiap butir dibuka **tanda golongan** dan ditutup **Bab asalnya**:

| Tanda golongan | Artinya |
| --- | --- |
| `[terverifikasi]` | perilaku **ditiru** dari Pega; buktinya ada di Bab yang disebut |
| `[keputusan work owner]` | perilaku **diputuskan**, bukan ditiru |
| `[data DBA]` | bersandar pada DDL atau isi basis data produksi |
| `[terbuka]` | **belum terjawab** — menyebut nomor butirnya di register 12 butir |
| ⚠️ | menguji **penyimpangan sadar** dari Pega — **daftar lengkapnya di Bab 11** |

> ⛔ **RALAT 2026-09-19** — baris di atas semula menyebut angka; kalimat lamanya **dikutip, tidak
> dihapus**: *"menguji **penyimpangan sadar** dari Pega — hanya ada **dua**, keduanya di Bab 11"*.
> Penyimpangannya sudah bertambah, dan seperti **AC 72**, bentuk **tak-berangka** dipakai supaya
> baris ini tidak usang lagi tiap kali daftarnya bertambah.

⛔ Butir ber-`[terbuka]` **bukan sasaran uji**. Ia menandai tempat yang **belum punya sasaran uji**,
supaya ketiadaannya terlihat dan tidak ditambal dengan tebakan.

---

### Lingkup dan identitas

1. `[keputusan work owner]` Sebuah kasus komite menunjuk **tepat satu** baris penyesuaian pada kasus
   klaim induk. Penunjuk itu **wajib isi** dan **unik** — dua kasus komite **tidak dapat** menunjuk
   baris penyesuaian yang sama, dan kasus komite **tidak dapat** dibuat tanpa penunjuk itu. *(Bab 1)*
2. `[keputusan work owner]` **Layar** komite tidak menulis satu pun data milik kasus klaim induk; di
   layar, data berawalan halaman induk **hanya dibaca**. Penulisan ke kasus klaim induk hanya terjadi
   lewat langkah penyimpanan — lihat AC 55 *(dua penanda usul)* dan AC 13 *(status penolakan)*.
   *(Bab 1)*
3. `[terbuka]` **butir 2** — daftar properti yang **diwarisi** kasus komite dari kelas induknya belum
   dapat diperiksa: ekspor Pega tidak memuat satu pun berkas kelas kerja. Selama ekspor kelas itu
   belum ada, **kelengkapan** kolom kasus komite **tidak punya sasaran uji**. ⛔ Jangan menutupnya
   dengan menganggap daftar kolom yang ada sudah lengkap. *(Bab 1)*

### Daur hidup kasus

4. `[terverifikasi]` Kasus komite yang baru dibuat **langsung** masuk kotak kerja seorang penyetuju.
   Test yang menemukan tahap antara — persetujuan pembuatan, antrean umum, atau tahap tunggu —
   **gagal**. *(Bab 2)*
5. `[terverifikasi]` Kasus yang tidak lagi punya penyetuju berikutnya **selesai**, berstatus akhir
   *selesai-tuntas*, dan **tidak muncul lagi** di kotak kerja siapa pun. Tidak ada tahap menggantung.
   *(Bab 2)*
6. `[terverifikasi]` Satu-satunya jalan kembali adalah kembali ke **tahap penyetuju yang sama**.
   Test yang menemukan kasus komite kembali ke tahap awal, atau berpindah ke tahap baru per
   penyetuju, **gagal**. *(Bab 2)*
7. `[keputusan work owner]` **+ `[terverifikasi]` 2026-09-19** Titik lompat darurat
   `komiteAccept_ticket` **tidak dibuat**. ⭐ **Dan tidak membangunnya adalah PARITAS, bukan
   penyimpangan:** sapuan **seluruh 409 berkas** dua modul menemukan **nol** pemakaian metode
   `Obj-Set-Tickets` dan **nol** pemakaian activity `SetTicket`, sehingga di Pega pun tiket itu
   **tidak pernah dilempar oleh siapa pun**. Test yang menemukan jalan masuk ke alur komite selain
   penyerahan dari kasus klaim **gagal**. *(Bab 2)*

   > ⛔ **TAMBAHAN 2026-09-19** — tanda golongan lamanya **dikutip, tidak dihapus**:
   > *"`[keputusan work owner]` Titik lompat darurat `komiteAccept_ticket` **tidak dibuat**."*
   > Yang ditambahkan **hanya** `[terverifikasi]`: keputusan ronde 7 diambil **tanpa bukti korpus**
   > *(butir 8 register ronde 6 berbunyi "siapa pemicunya")*, dan bukti itu **sekarang ada**.
   > ⛔ Isinya tidak berubah sedikit pun.

### Tangga penyetuju

8. `[terverifikasi]` Daftar penyetuju ditetapkan **sekali, saat kasus dibuat**, dan tidak berubah di
   tengah jalan. Penambahan atau pengurangan penyetuju sesudah kasus lahir **ditolak**. *(Bab 3)*
9. `[keputusan work owner]` Pencacah **jumlah penyetuju yang dibutuhkan** adalah **kolom yang
   disimpan**, diisi dari jumlah baris daftar penyetuju pada saat kasus dibuat. Test yang
   menemukannya **dihitung ulang setiap kali dibaca** **gagal**. *(Bab 3)*
10. `[terverifikasi]` Kotak kerja seorang penyetuju **hanya** memuat kasus yang sedang **gilirannya**.
    Penyetuju yang belum tiba gilirannya tidak melihatnya; penyetuju yang sudah memutuskan tidak
    melihatnya lagi. *(Bab 3)*
11. `[terverifikasi]` Tujuan rute diambil dari **akun operator** penyetuju yang sedang giliran —
    bukan jabatannya, dan bukan nama antrean. *(Bab 3)*
12. `[terverifikasi]` Pencacah **penyetuju keberapa** naik **tepat satu** tiap putaran, dan jumlah
    putaran **tidak melebihi** jumlah penyetuju yang dibutuhkan. Test yang menemukannya melompat
    **gagal**. *(Bab 3)*
13. `[terverifikasi]` Satu keputusan **tolak** menghentikan rangkaian: penyetuju berikutnya
    **tidak menerima** kasus, **seluruh sisa** penyetuju ditandai tertolak otomatis, dan status
    penolakan sampai ke **baris penyesuaian induk**. Tidak ada baris penyetuju yang tertinggal
    berkeputusan kosong. *(Bab 3)*
14. `[terverifikasi]` Antrean berjenjang lama — empat antrean bertingkat yang keenam langkahnya
    di-remark di Pega — **tidak dibuat**. Test yang menemukannya **gagal**. *(Bab 3)*
15. `[terbuka]` **butir 3** — arti kedua jenis langkah berulang Pega, dan **halaman mana** yang
    diulang, tidak dinyatakan di mana pun dalam ekspor. Selama itu, **apa** yang diputari di tiap
    langkah berulang **tidak punya sasaran uji**. *(Bab 3)*
16. `[terbuka]` **butir 4** — **berapa kali** sebuah langkah berulang berputar belum terbaca. Karena
    itu jumlah penandaan sisa penyetuju, jumlah penulisan nilai, dan jumlah pengiriman ke Kasir
    **belum punya angka sasaran**. ⛔ Jangan menebak "sekali". *(Bab 3)*
17. `[terbuka]` **butir 5** — urutan pemeriksaan ketika dua keluarga gerbang pada satu langkah
    **sama-sama terisi** tidak terbaca dari struktur ekspor. Selama itu, perilaku langkah semacam
    itu **tidak punya sasaran uji yang pasti**. *(Bab 3)*

### Layar komite

18. `[keputusan work owner]` Layar utama komite menampilkan **93** medan — **62** milik kasus komite,
    **31** milik kasus klaim induk. Angka ini **SAH** dan tidak diturunkan ulang. *(Bab 4)*
19. `[terverifikasi]` **"Dibuat oleh"** dan **"tanggal dibuat"** **tampil di layar**, meski keduanya
    berada di luar hitungan 93. Test yang tidak menemukannya **gagal**. *(Bab 4)*
20. `[terverifikasi]` Medan yang di Pega hanya-baca tampil sebagai **bacaan saja** dan **tidak dapat
    diubah** penyetuju. *(Bab 4)*
21. `[keputusan work owner]` Properti jenis treaty, nomor klaim, dan dua kotak total estimasi — yang
    di Pega duduk di balik gerbang mati dan **tidak pernah tampil** — **tidak dibuat sama sekali**.
    Test yang menemukan salah satunya **gagal**. *(Bab 4)*
22. `[terverifikasi]` Layar kedua — rincian objek pertanggungan — menampilkan **tiga** medan: nama
    objek, mata uang, dan nilai pertanggungan per objek. Ketiganya **hanya-baca**, dan ketiganya
    **sudah ada di layar utama**: layar kedua **tidak menambah satu kolom pun**. *(Bab 4)*
23. `[terverifikasi]` **Keputusan setuju/tolak** dan **catatan** **selalu** dapat disunting dan
    **wajib isi**, pada ketiga jalur. Pengiriman tanpa salah satunya **ditolak**. *(Bab 4)*
24. `[terverifikasi]` Keempat isian bersyarat muncul **bertingkat**: penanda persetujuan bersyarat
    hanya bila penyetuju **menyetujui** **dan** jalurnya **penyesuaian** · isi syarat hanya
    **setelah** penandanya dicentang · kedua penanda usul hanya pada jalur **penyesuaian**.
    Keempatnya **tidak wajib isi** — pengiriman tanpa mengisinya tetap diterima. *(Bab 4)*
25. `[terverifikasi]` Tanggal dan identitas penyetuju **terisi otomatis** saat ia memutuskan, dan
    penyetuju **tidak dapat mengubahnya**. *(Bab 4)*
26. `[terbuka]` **butir 9** — **dari mana layar kedua dibuka** belum terbaca: layar itu tidak dirujuk
    layar utama, pembungkusnya, maupun berkas alur. Jalan masuknya **belum punya sasaran uji**.
    ⛔ Jangan mengarang tombolnya. *(Bab 4)*

79. `[keputusan work owner]` ⚠️ **Penyimpangan sadar 5 — daftar penyetuju tampil.** Layar komite
    menampilkan **daftar penyetuju**: **urutan** dan **siapa**, untuk seluruh tingkat — bukan hanya
    yang sudah lewat. Test yang tidak menemukannya **gagal**. **⚠️ penyimpangan sadar**; di Pega
    daftar itu **tidak tampil sama sekali**. *(Bab 4 · menutup US 3)*
80. `[keputusan work owner]` ⚠️ **Penyimpangan sadar 5 — keputusan penyetuju sebelumnya tampil.**
    Untuk setiap penyetuju yang **sudah memutuskan**, layar menampilkan **keputusannya, catatannya,
    dan tanggalnya**; yang **belum** memutuskan tampil **tanpa isi**. ⛔ **Nol kolom basis data
    baru** — ketiganya sudah tersimpan di baris penyetuju *(AC 53)*. **⚠️ penyimpangan sadar.**
    *(Bab 4 · menutup US 4)*

### Tiga jalur

27. `[keputusan work owner]` Kasus komite berjalan pada **tepat tiga** jalur — **TRANSFER
    ADJUSTMENT**, **REJECT**, **CLOSE**. Jalur keempat **tidak dibuat**: di Pega nilainya tidak
    pernah ditulis siapa pun, dan ketiga langkah yang membacanya bergerbang mati. *(Bab 5)*
28. `[terverifikasi]` Pada jalur **penyesuaian**, ke-**24** medan bergerbang **tampil**; pada jalur
    **tolak** dan **tutup** ke-24 medan itu **tidak tampil**, dan **ketiga judul bagian layar**
    berganti mengikuti jalur. *(Bab 5)*
29. `[keputusan work owner]` Pada jalur **REJECT** dan **CLOSE**, penyetuju hanya mengisi **dua**
    hal — keputusan dan catatan. **Ditiru apa adanya.** *(Bab 5)*
30. `[terverifikasi]` Pada jalur **bukan penyesuaian**, penghitungan total **tidak berjalan sama
    sekali**. *(Bab 5)*

### Nomor akseptasi

31. `[terverifikasi]` Nomor akseptasi terbit **tepat sekali**: pada penyetuju **terakhir**, yang
    **menyetujui**, **tanpa syarat**. *(Bab 6)*
32. `[terverifikasi]` Nomor **tidak terbit** ketika persetujuannya **bersyarat**, meski penyetuju
    terakhir menyetujui; dan **tidak terbit** pada penyetuju selain yang terakhir. *(Bab 6)*
33. `[terverifikasi]` Menjalankan ulang alur pada kasus yang **sudah bernomor** **tidak menerbitkan
    nomor kedua**. *(Bab 6)*
34. `[terverifikasi]` Status akseptasi pada baris penyesuaian induk terisi **bersama** nomornya —
    bukan sesudahnya dan bukan terpisah. *(Bab 6)*
35. `[terverifikasi]` Kedua rule yang di Pega langkahnya **di-remark** **tidak dipanggil**; keduanya
    juga tidak ada di ekspor modul ini. Test yang menemukan salah satunya dipanggil **gagal**.
    *(Bab 6)*

### Angka uang

36. `[terverifikasi]` Total dijumlahkan **per mata uang**, dan baris **berstatus tolak tidak ikut
    dijumlahkan**. *(Bab 7)*
37. `[terverifikasi]` Baris bermata-uang-sama **digabung** lebih dulu, sebelum dijumlahkan. *(Bab 7)*
38. `[terverifikasi]` Total dalam rupiah = total mata uang **dikali kurs** yang tersimpan pada baris
    itu — bukan kurs yang dicari ulang saat menampilkan. *(Bab 7)*
39. `[keputusan work owner]` Kurs diambil dengan **dua masukan saja**: kode mata uang dan **tanggal
    saat pencarian dijalankan**. Test yang menemukan nomor kasus, lini, atau tanggal kasus ikut
    dikirim **gagal**. *(Bab 7)*
40. `[keputusan work owner]` Nilai uang melewati aplikasi sebagai **tipe desimal**; **nol** bilangan
    pecahan biner di lapisan mana pun maupun di kontrak antarmuka. *(Bab 7)*
41. `[keputusan work owner]` Angka uang **di layar** tampil dengan **4 angka di belakang koma** —
    seragam, di semua tempat. *(Bab 7)*
42. `[keputusan work owner]` Nilai yang **melewati batas digit** kolom **ditolak dengan galat yang
    terlihat** — bukan dibulatkan diam-diam, dan bukan ditelan. *(Bab 7)*
43. `[terbuka]` **butir 8** — **mana** ketelitian pembagian persen Pega yang sah belum ditetapkan;
    di Pega ada tiga cara berbeda untuk operasi yang sama, ditambah satu pemotongan tersembunyi.
    Selama itu, **selisih terhadap data lama belum dapat diterangkan** dan bagian itu belum punya
    sasaran uji. *(Bab 7)*
44. `[terbuka]` **butir 11** — **berapa besar** selisih akibat persen 10 desimal menjadi 8 saat
    disimpan **belum diukur**, dan selisih itu **menumpuk lewat perulangan** sebelum sampai ke
    Kasir. ⛔ Jangan menebak "kecil, boleh diabaikan". *(Bab 7)*

### Efek keluar

45. `[terverifikasi]` Nomor akseptasi terbit → **delapan** efek keluar berjalan: baris akseptasi ·
    dokumen PDF akseptasi · data klaim bentuk JSON · panggilan arasapas · baris log pemantauan ·
    riwayat akseptasi · data pembayaran ke Kasir · email. Test yang menemukan salah satunya tidak
    berjalan **gagal**. *(Bab 8)*
46. `[keputusan work owner]` Data pembayaran terkirim ke Kasir **hanya bila** nomor akseptasi
    **sudah tercatat** di daftar akseptasi klaim. Bila belum tercatat, pengiriman **tidak
    berjalan**. Fakta asalnya bertanda `[data work owner]`; yang mengikat adalah keputusan **13**
    di Bab 10 — *"syarat resmi yang dipertahankan"*. *(Bab 8)*
47. `[keputusan work owner]` Kegagalan pengiriman ke Kasir **tidak menghentikan kasus**, **email
    tetap terkirim**, dan **tidak ada percobaan ulang otomatis**. Kegagalannya **tercatat** di jalur
    galatnya sendiri. **Ditiru apa adanya**, dengan risiko diterima sadar. *(Bab 8)*
48. `[terverifikasi]` Kegagalan pembuatan **dokumen PDF** **menghentikan dengan galat yang terlihat**
    — sikap yang **berbeda** dari pengiriman ke Kasir, dan perbedaan itu **disengaja**. *(Bab 8)*
49. `[terverifikasi]` Alamat tujuan layanan diambil dari tabel alamat layanan **disaring dua kunci**
    — kategori dan sub-kategori — sehingga hasilnya **satu baris**, dan pasangan kuncinya **berbeda
    per tujuan**. Test yang menemukan baris pertama diambil tanpa saringan **gagal**. *(Bab 8)*
50. `[terbuka]` **butir 6** — **unggah berkas** dan **email** **tanpa penanganan gagal** belum dibawa
    ke work owner. Perilaku saat gagal **belum punya sasaran uji**, dan ⛔ **tidak boleh ditambal
    diam-diam**: bila memang tidak tertangani, itu **didokumentasikan**. *(Bab 8)*
51. `[terbuka]` **butir 7** — catatan pengembang **"FIX ERROR HANDLING"** pada rule pengunggah berkas
    **tidak berjejak** di keempat keluarga wadah yang sudah disisir habis. Apakah penanganannya ada
    di tempat lain **belum terjawab**. *(Bab 8)*
52. `[terbuka]` **butir 12** — **baris mana** di tabel alamat layanan yang menunjuk **lingkungan
    uji** belum ditetapkan. Selama itu, kedelapan efek keluar **belum dapat dijalankan** dengan
    layanan sungguhan sebagaimana diputuskan di AC 69. *(Bab 8)*

### Tabel yang disentuh

53. `[terverifikasi]` Keputusan dan catatan tersimpan pada **baris penyetuju yang bersangkutan**,
    bukan di header kasus komite. *(Bab 9)*
54. `[keputusan work owner]` Hanya ada **satu** kolom tanggal persetujuan. Di Pega ada dua properti
    tanggal yang diisi langkah yang sama dengan nilai identik; yang kedua **sengaja tidak dijadikan
    kolom**. Test yang menemukan kolom tanggal kedua **gagal** — jangan menambahkannya kembali karena
    "ada di korpus". *(Bab 9)*
55. `[keputusan work owner]` Kedua **penanda usul** — usul tutup klaim dan usul cadangkan klaim —
    tersimpan di **header kasus komite**, dan nilai yang sama **juga sampai** ke kasus klaim induk
    dengan **nama penanda yang berbeda**; langkah penyalinnya di Pega **tanpa gerbang**, jadi selalu
    berjalan. *(Bab 9)*
56. `[keputusan work owner]` Kesembilan rule yang di Pega **menyimpan sendiri** **tetap menyimpan
    sendiri**. Bila langkah sesudahnya gagal, tabelnya **tetap terisi** dan **tidak ada pembatalan,
    penghapusan, maupun penanda batal**. Test yang menemukan pembatalan **gagal** — data separuh jadi
    adalah perilaku yang **disengaja ditiru**. *(Bab 9)*
57. `[data DBA]` Riwayat akseptasi terisi **6 kolom** dari jalur komite: pengenal kasus · tanggal
    transfer · status · nama pengguna · kotak kerja · pengenal komite. Kolom ketujuh — akun
    operator — **tetap kosong**. *(Bab 9)*
58. `[terverifikasi]` Kolom **kotak kerja** pada riwayat akseptasi terisi **literal tetap**, selalu,
    tanpa cabang. *(Bab 9)*
59. `[data DBA]` **Pengenal kasus** pada riwayat akseptasi adalah **teks** berisi kunci instance
    penuh **berikut awalan kelas** — **bukan angka**, dan **bukan kunci asing basis data**. *(Bab 9)*
60. `[keputusan work owner]` Baris komite ber-`ID` berawalan **`TKMT-`**, ber-`COVER_KEY` berisi `ID`
    baris klaim induk berawalan **`CLMP-`**, dan ber-`LINI` **`PROP`**. Header kasus komite memakai
    **kunci yang sama persis** — test yang menemukan kolom penyambung terpisah **gagal**. *(Bab 9)*
61. `[keputusan work owner]` Ketiga pencacah — jumlah penyetuju · penyetuju keberapa · urutan
    penyetuju — tersimpan sebagai **bilangan bulat**, bukan angka desimal. Ketiganya **tidak ikut**
    aturan ketelitian angka. *(Bab 9)*
62. `[terverifikasi]` Penolakan klaim tercatat **tersendiri**, di tabel penolakan klaim, terpisah
    dari riwayat akseptasi. *(Bab 9)*
63. `[terbuka]` **butir 1** — kolom **akun operator** di riwayat akseptasi **tidak pernah diisi jalur
    komite**. Adakah penulis lain di luar modul ini, atau kolom itu memang selalu kosong, **belum
    terjawab**. ⛔ Tidak ditebak, dan **tidak diusulkan dibuang**. *(Bab 9)*
64. `[keputusan work owner]` Kasus komite lama **terbaca di sistem baru** lengkap dengan daftar
    penyetuju, keputusan, catatan, dan tanggalnya; kedua pencacah tangga dipindahkan **apa adanya**.
    *(Bab 9 · tiket 13)*
65. `[keputusan work owner]` Penunjuk **posisional** lama diterjemahkan menjadi penunjuk baris
    penyesuaian **apa adanya** — **termasuk yang sudah salah alamat di Pega**. Migrasi **tidak
    menolak** dan **tidak memperbaiki** baris yang mencurigakan; test yang menemukan baris ditolak
    atau diperbaiki **gagal**. *(Bab 9 · tiket 13)*
66. `[keputusan work owner]` Nilai uang **lama** yang **melewati batas digit** kolom menyebabkan
    **kegagalan yang terlihat** saat migrasi — bukan pembulatan diam-diam. *(Bab 7 · tiket 13)*

74. `[keputusan work owner]` Kedua **penanda usul** tersimpan **hanya** sebagai `'1'` atau `'0'`;
    nilai lain **ditolak**. Kotak yang **tidak disentuh** tersimpan `'0'`, bukan kosong. *(Bab 9)*
75. `[keputusan work owner]` Migrasi data lama mengubah nilai **benar/salah** menjadi `'1'`/`'0'`,
    dan nilai **kosong** menjadi `'0'`; test yang menemukan kolom usul **kosong** sesudah migrasi
    **gagal**. *(Bab 9)*

### Kumpulan keputusan work owner

67. `[keputusan work owner]` Daftar hak akses Pega **tidak dipindahkan** — ia menyebut kelas **tanpa
    menyebut hak**, jadi tidak menegakkan apa pun. Wewenang di sistem baru ditentukan **aturan peran
    sistem baru**; penegakan **giliran** bukan penegakan **wewenang**. *(Bab 10)*

    > ⭐ **LANJUTAN 2026-09-19 — kalimat di atas TETAP BERLAKU dan tidak dicabut.** Ia benar:
    > giliran memang bukan wewenang. ⚠️ **Yang kurang hanyalah kelanjutannya** — sampai hari ini
    > *"aturan peran sistem baru"* **belum ditulis di mana pun**, sehingga wewenang memutuskan
    > menggantung tanpa sasaran uji.
    >
    > ⭐ **Sekarang ia punya AC-nya sendiri: AC 81 dan AC 82.** Wewenang **menyimpan keputusan**
    > tidak lagi menunggu aturan peran yang belum ada — ia ditegakkan **per `KomiteID` pada tingkat
    > yang sedang berjalan**, sesuai **ADR-0014**. *(Bab 10 keputusan 30 · Bab 11 titik 6)*
68. `[keputusan work owner]` Kelima **perintah lompat menggantung** **tidak dibuat**; test yang
    menemukan lompatan ke tanda yang tidak ada **gagal**. Lompatan yang **hidup dan tandanya ada**
    — lompatan pada jalur penolakan — **tetap dipindahkan**. *(Bab 10)*
69. `[keputusan work owner]` Kedelapan efek keluar diuji dengan **layanan sungguhan**, bukan pengganti
    tiruan, dan dijalankan di **lingkungan uji terpisah, bukan produksi**. *(Bab 10)*
83. `[keputusan work owner]` **2026-09-19** · **ADR-0009** **Seluruh data kasus komite dipindahkan;
    tidak ada koeksistensi dua penulis.** Test yang menemukan kasus komite berjalan diselesaikan di
    sistem lama sementara yang baru sudah menerima kasus baru **gagal**. ⚠️ **Alasan khas modul
    ini, dan ia lebih kuat daripada di Claim Life:** **keputusan 28** menetapkan klaim yang pernah
    punya kasus komite **tidak dapat dihapus selamanya**, dan **AC 78** menetapkan bekunya **tidak
    mencair**. Dua penulis atas klaim yang sama, dengan aturan beku yang **hanya dipahami satu
    sisi**, **akan** melahirkan selisih yang tidak dapat diterangkan. *(Bab 10)*
86. `[keputusan work owner]` **+ `[terverifikasi]` 2026-09-19** Daftar **jenis berkas yang diterima**
    mengikuti daftar lama **apa adanya — tidak ditambah**. Test yang menemukan jenis berkas di luar
    daftar lama diterima **gagal**. ⚠️ `[terverifikasi]` Salinan rule penentu jenis berkas yang
    lebih baru di modul saudara memuat **48** baris lawan **42** di sini *(catatan pengembangnya
    berbunyi "add avi")*; ⛔ **menambah jenis yang diterima adalah PERUBAHAN PERILAKU**, dan tidak
    ada yang memintanya. *(Bab 9)*

### Titik yang **SENGAJA DIUBAH** dari Pega

70. `[keputusan work owner]` ⚠️ **Penyimpangan sadar 1 — urutan efek keluar.** Di Pega kedelapan efek
    keluar berjalan **sebelum** penyimpanan; urutan itu **TIDAK DITIRU**. ⛔ Urutan penggantinya
    **belum diputuskan**, jadi AC ini **tidak mengunci urutan mana pun** — dan test yang mengunci
    **urutan Pega** sebagai syarat **gagal**. *(Bab 11)*
71. `[keputusan work owner]` ⚠️ **Penyimpangan sadar 2 — ketelitian angka.** Hitungan berjalan sampai
    **20 angka di belakang koma tanpa pembulatan di tengah jalan**, dan **satu-satunya** titik
    pembulatan ada di **batas penyimpanan**. Hasil sistem baru **akan berbeda dari Pega di angka
    belakang koma** — khususnya persen yang di Pega dihitung sampai 10 desimal menjadi **8** saat
    disimpan. Perbedaan itu **disengaja**: test yang menuntut hasil **identik** dengan Pega **gagal**.
    *(Bab 11)*
72. `[keputusan work owner]` **Tidak ada penyimpangan di luar yang terdaftar di Bab 11.** Perilaku yang berbeda dari Pega
    **hanya** yang tercatat sebagai **penyimpangan sadar** di bab itu; sisanya **ditiru apa adanya**
    atau **dibuang sebagai jalur mati**. Test atau kode yang memperkenalkan penyimpangan yang
    **tidak terdaftar di Bab 11** **gagal**. *(Bab 11)*
    > ⛔ **RALAT 2026-09-19** `[keputusan work owner]` — **AC 72 dibuat tidak berangka**, supaya
    > tidak perlu diralat lagi setiap ada penyimpangan baru. Kalimat lamanya **dikutip utuh, tidak
    > dihapus**:
    >
    > > *"72. `[terverifikasi]` **Tidak ada penyimpangan ketiga.** Perilaku yang berbeda dari Pega
    > > di luar AC 70 dan AC 71 **tidak diperkenalkan**; sisanya **ditiru apa adanya** atau
    > > **dibuang sebagai jalur mati**. Test atau kode yang memperkenalkan perbedaan ketiga
    > > **gagal**."*
    >
    > ⚠️ Kalimat lama menyebut **angka** — *"tidak ada penyimpangan ketiga"* — dan menjadi salah
    > begitu penyimpangan ke-3 diputuskan. Bentuk **tak-berangka** ini tidak akan usang lagi.
    >
    > ⚠️ **Tanda `[terverifikasi]` DILEPAS, dan ⛔ tidak diganti tanda lain.** Alasannya:
    > *"tidak ada penyimpangan lain"* **bukan fakta yang dibaca dari korpus**, melainkan pernyataan
    > tentang **isi berkas ini sendiri**. Kalimatnya berdiri sebagai **aturan bab**, bukan sebagai
    > temuan — karena itu ia **satu-satunya AC tanpa tanda golongan**, dan itu **disengaja**.
    >
    > ✅ **SUSULAN 2026-09-19** `[keputusan work owner]` — **AC 72 tetap tak-berangka, tetapi
    > dibuka tanda `[keputusan work owner]`.** Alasan pelepasan tanda di atas **tetap dicatat
    > sebagai jejak**, tetapi **sudah tidak berlaku**: kalimat ini memang bukan temuan korpus,
    > melainkan **keputusan** tentang bagaimana berkas ini dibaca — dan itu **punya tandanya
    > sendiri**. ⛔ **Nol AC tanpa tanda golongan lagi.**
73. `[terbuka]` **butir 13** — **sasaran uji untuk urutan efek keluar** menunggu urutan barunya
    ditetapkan. Selama belum, bagian itu **tidak punya sasaran uji**, dan itu **bukan kekurangan
    berkas ini**. *(Bab 11)*
76. `[keputusan work owner]` ⚠️ **Penyimpangan sadar 3 — baris beku, sisi ubah.** Baris penyesuaian
    yang sudah **diserahkan ke komite** **tidak dapat diubah**; upaya menyuntingnya **ditolak di
    lapisan layanan**. *(Bab 11)*
77. `[keputusan work owner]` ⚠️ **Penyimpangan sadar 3 — sisi hapus.** Baris penyesuaian yang sudah
    diserahkan ke komite **tidak dapat dihapus satu per satu**; upaya menghapusnya **ditolak**.
    **Klaim induknya tidak dapat dihapus sama sekali — selamanya**, walau komitenya **sudah
    selesai maupun menolak**. Klaim yang **belum pernah** punya kasus komite **tetap dapat
    dihapus** dan **mengkaskade seperti biasa**. *(Bab 11)*
    > ⛔ **RALAT 2026-09-19 — kedua** `[keputusan work owner]`: **larangannya PERMANEN, bukan
    > sementara.** Kedua kalimat lamanya **dikutip utuh, tidak satu pun dihapus**:
    >
    > > *"77. `[keputusan work owner]` ⚠️ **Penyimpangan sadar 3 — sisi hapus.** Baris penyesuaian
    > > yang sudah diserahkan ke komite **tidak dapat dihapus**; upaya menghapusnya **ditolak**,
    > > **termasuk lewat kaskade hapus klaim**."*
    > >
    > > *"77. `[keputusan work owner]` ⚠️ **Penyimpangan sadar 3 — sisi hapus.** Baris penyesuaian
    > > yang sudah diserahkan ke komite **tidak dapat dihapus satu per satu**; upaya menghapusnya
    > > **ditolak**. **Klaim induknya tidak dapat dihapus selama masih ada kasus komite yang
    > > berjalan.** Sesudah komite **selesai atau menolak**, menghapus klaim **boleh** dan
    > > **mengkaskade seperti biasa** — termasuk **menyapu baris beku itu**."*
    >
    > *(anotasi ralat ronde sebelumnya, dikutip utuh)* — *"⛔ **RALAT 2026-09-19**
    > `[keputusan work owner]` — **kunci beku tidak mematikan fitur hapus klaim; yang dikunci
    > adalah saat komitenya masih berjalan.** Kalimat lamanya **dikutip utuh, tidak dihapus**:"*
    > dan *"⚠️ Kalimat lama **menutup kaskade hapus klaim sepenuhnya**, dan itu **bertabrakan**
    > dengan `claim-prop/spec.md` **AC 5** — *"Menghapus klaim **mengkaskade sampai tingkat
    > terdalam**; test wajib memeriksa cicit."* Keputusan work owner: **kuncinya berlaku sampai
    > komite selesai, bukan selamanya.**"*
    >
    > ⚠️ Bentuk **pertama** menutup kaskade hapus klaim sepenuhnya; bentuk **kedua** membukanya
    > kembali sesudah komite selesai. Keduanya salah. `[keputusan work owner]` 2026-09-19:
    > **klaim yang PERNAH punya kasus komite tidak dapat dihapus, selamanya** — lihat keputusan
    > **28** di bab 10, yang **mencabut keputusan 24**.
    >
    > ⚠️ **Akibatnya terhadap `claim-prop/spec.md` AC 5** — *"Menghapus klaim **mengkaskade sampai
    > tingkat terdalam**; test wajib memeriksa cicit."* — kaskade itu kini punya **pengecualian
    > permanen**, bukan pengecualian sementara. ⛔ Berkas modul itu **tidak disunting dari sini**.
78. `[keputusan work owner]` ⚠️ **Penyimpangan sadar 3 — bekunya tidak mencair.** **Penolakan oleh
    komite tidak mencairkan** bekunya; test yang **berhasil menyunting** baris sesudah penolakan
    **gagal**. *(Bab 11)*

⚠️ **AC 77 dan AC 78 saling menguatkan:** **keduanya permanen**. Penolakan komite **tidak
mengembalikan hak menyunting** baris **(AC 78)** dan **tidak membuka hak menghapus klaimnya**
**(AC 77)**. **Sekali sebuah baris pernah dinilai komite, jejaknya tidak dapat dihapus lewat jalan
mana pun.**

> ⛔ **RALAT 2026-09-19** — kalimat lamanya **dikutip utuh, tidak dihapus**: *"⚠️ **AC 77 dan
> AC 78 tidak bertentangan:** penolakan komite **tidak mengembalikan hak menyunting** baris
> **(AC 78)**, tetapi **membuka hak menghapus klaimnya** **(AC 77)**. **Sunting dan hapus
> diperlakukan berbeda dengan sengaja.**"*
>
> ⚠️ Ia menerangkan dua sisi yang **umurnya berbeda**. Sesudah `[keputusan work owner]`
> 2026-09-19, **umurnya sama** — keduanya selamanya — jadi tidak ada lagi beda yang perlu
> diterangkan.

81. `[keputusan work owner]` ⚠️ **Penyimpangan sadar 6 — wewenang menyimpan keputusan.**
    **Hanya pemilik `KomiteID` pada tingkat yang sedang berjalan** yang dapat **menyimpan
    keputusan** penyetuju. Test yang **berhasil menyimpan** keputusan sebagai pengguna lain
    **gagal** — termasuk pengguna yang **sudah** memutuskan di tingkat sebelumnya dan pengguna yang
    **belum** tiba gilirannya. ⛔ **Menyembunyikan layarnya bukan pemenuhan AC ini**; yang diuji
    adalah **penolakan di lapisan layanan**. *(Bab 10 keputusan 30 · Bab 11 titik 6 · ADR-0014)*
82. `[keputusan work owner]` ⚠️ **Penyimpangan sadar 6 — penolakannya terlihat.** Penolakan pada
    AC 81 menghasilkan **galat yang terbaca pengguna**. Test yang menemukan penyimpanan
    **diam-diam diabaikan** — tampak berhasil di layar tetapi tidak tersimpan — **gagal**.
    *(Bab 10 keputusan 30 · Bab 11 titik 6)*
84. `[keputusan work owner]` ⚠️ **Penyimpangan sadar — wewenang menghapus catatan kronologi.**
    Wewenang itu adalah **izin eksplisit yang diberikan lewat peran**, **bukan** perbandingan
    terhadap **teks jabatan** operator. Test yang menemukan wewenang bergantung pada **isi medan
    jabatan** **gagal**. ⚠️ `[terverifikasi]` Di Pega ia berupa perbandingan terhadap teks
    jabatan — satu-satunya pemeriksaan wewenang sejati di modul ini, dan **rule yang sama persis
    ada di Claim Prop**. ⛔ Teks jabatan adalah **data kepegawaian**: begitu satu huruf diganti,
    wewenangnya **berubah diam-diam**. *(Bab 10)*
85. `[keputusan work owner]` ⚠️ **Penyimpangan sadar — pengenal berkas yang sudah ditambal.**
    Pengenal berkas dibuat dengan cara berketelitian **nanodetik** **dan** disertai **nilai unik
    sejagat** dari basis data. Test yang menemukan **dua unggahan berdekatan menghasilkan pengenal
    yang sama** **gagal**. ⚠️ `[terverifikasi]` Versi yang dipakai modul ini **belum ditambal** —
    ia bersandar pada ketelitian **milidetik tanpa nilai unik**, sehingga dua unggahan dalam
    milidetik yang sama menghasilkan **pengenal kembar**. ⭐ **Tambalannya bukan rancangan baru:**
    ia **sudah berjalan** pada salinan rule yang sama di modul saudara. *(Bab 9)*

---

**Jumlah AC: 86.** **AC tanpa tanda golongan: 0.** **Butir register yang terwakili: 12 dari 13**
— butir **14** *(kasus lama hasil migrasi yang penunjuknya salah alamat)* **sengaja belum punya AC**,
karena perilakunya belum diputuskan.

> ⛔ **RALAT 2026-09-19 — kedua.** Kalimat lamanya **dikutip, tidak dihapus**:
> *"**Jumlah AC: 80.** **AC tanpa tanda golongan: 0.** **Butir register yang terwakili: 12 dari
> 13**"*. **Enam AC ditambahkan** dari lima keputusan work owner 2026-09-19 — **81 · 82**
> *(wewenang menyimpan keputusan, **ADR-0014**)* · **83** *(migrasi penuh, **ADR-0009**)* ·
> **84** *(izin hapus catatan kronologi)* · **85** *(pengenal berkas yang sudah ditambal)* ·
> **86** *(jenis berkas, paritas)*. ⛔ **Nol AC lama dinomori ulang.**
>
> ⚠️ **Empat di antaranya penyimpangan sadar** — 81 · 82 · 84 · 85 — sehingga **Bab 11
> bertambah satu titik, menjadi ENAM**.

> ⛔ **RALAT 2026-09-19** — kalimat lamanya **dikutip, tidak dihapus**: *"**Jumlah AC: 73.**
> **AC tanpa tanda golongan: 0.** **Butir register yang terwakili: 12 dari 12**"*. Lima AC
> ditambahkan — **74 · 75** *(nilai penanda usul)* dan **76 · 77 · 78** *(baris beku)* — dan
> register bertambah satu butir. ⛔ **Nol AC lama dinomori ulang.**
>
> ⚠️ **Akibatnya nomor AC TIDAK lagi menaik terus di dalam berkas.** AC **74 · 75** duduk di
> akhir sub-judul *Tabel yang disentuh*, jadi letaknya **sebelum** AC 67. Itu **disengaja**:
> sub-judul menentukan tempat, nomor menentukan identitas — dan **empat belas tiket merujuk AC
> lewat nomor**, jadi nomor tidak boleh digeser demi kerapian urutan.

**Daftar lama, dikutip apa adanya:** **Butir register yang terwakili: 12 dari 12**
— butir 1 *(AC 63)* · 2 *(AC 3)* · 3 *(AC 15)* · 4 *(AC 16)* · 5 *(AC 17)* · 6 *(AC 50)* ·
7 *(AC 51)* · 8 *(AC 43)* · 9 *(AC 26)* · 11 *(AC 44)* · 12 *(AC 52)* · 13 *(AC 73)*.

⛔ **Nol butir `[terbuka]` ditutup di bab ini.** Bab AC hanya menandai **di mana** sasaran uji belum
ada; yang menutup butir adalah work owner.

---

### ✅ Celah cakupan — peta 39 user story terhadap AC

*Bab ini sebelumnya tidak punya peta semacam ini, sehingga klaim "spec selesai" belum pernah diadu
dengan ke-39 user story-nya sendiri. Peta ini menutup celah itu.* **Dihitung 2026-09-19.**

| US | Ditutup AC | Catatan |
| --- | --- | --- |
| 1 · 2 | 4 · 10 · 11 | kotak kerja dan giliran |
| **3** | **79** | ⭐ **sebelumnya NOL** — ditutup penyimpangan sadar 5 |
| **4** | **80** | ⭐ **sebelumnya NOL** — ditutup penyimpangan sadar 5 |
| 5 | 8 · 9 | daftar penyetuju ditetapkan sekali |
| 6 · 7 | 18 | 93 medan layar |
| 8 | 38 | total rupiah |
| 9 | 22 | layar kedua — ⚠️ jalan masuknya masih butir `[terbuka]` **9** |
| 10 | 28 | tiga wajah menurut jalur |
| 11 | 19 | "dibuat oleh" dan "tanggal dibuat" |
| 12 | 20 | medan hanya-baca |
| 13 · 14 | 23 | keputusan dan catatan, wajib isi |
| 15 · 16 | 24 | rantai isian bersyarat |
| 17 | 24 · 55 | dua penanda usul |
| 18 | 25 | tanggal dan identitas otomatis |
| 19 | 12 | tangga maju |
| 20 | 5 | selesai di penyetuju terakhir |
| 21 · 22 | 13 | penolakan menghentikan, sisa ditandai |
| 23 | 31 · 33 | nomor terbit sekali |
| 24 | 32 | tidak terbit bila bersyarat |
| 25 | 45 | baris akseptasi |
| 26 | 45 · 48 | dokumen PDF |
| 27 | 45 · 46 | kiriman ke Kasir |
| 28 | 46 | syarat nomor akseptasi |
| 29 | 45 · 47 | email |
| 30 | 57 | riwayat akseptasi |
| 31 | 62 | penolakan klaim tersendiri |
| 32 | 38 | kurs tersimpan bersama transaksinya |
| 33 | 41 | empat angka di belakang koma |
| **34** | 36 · 40 · 42 **— sebagian** | ⚠️ *"dihitung dengan cara yang sama di semua jalur"* **belum punya AC sendiri**; nilai yang dikirim ke Kasir dihitung **di modul Claim Prop**, dan butir `[terbuka]` **8** *(ketelitian Pega mana yang sah)* masih terbuka |
| **35** | 9 · 38 **— sebagian** | ⚠️ *"tahu angka mana disimpan dan mana dihitung ulang"* tertutup **per kasus**, bukan sebagai **daftar menyeluruh** |
| 36 | 7 · 14 · 21 · 35 · 68 | jalur mati tidak dipindahkan |
| 37 | 70 · 71 · 72 · 76 · 77 · 78 · 79 · 80 | titik yang sengaja diubah |
| 38 | 3 · 15 · 16 · 17 · 26 · 43 · 44 · 50 · 51 · 52 · 63 · 73 | dua belas AC ber-`[terbuka]` |
| **39** | — | ⭐ **meta, sengaja tanpa AC**: *"bisa menghitung ulang sendiri angka apa pun"* dipenuhi **bab Lampiran — Aturan baca ekspor Pega**, bukan oleh perilaku yang diuji |

**Hasilnya:** **US tanpa AC sama sekali: NOL** *(turun dari **dua** — US 3 dan US 4)* ·
**tertutup sebagian: dua** — **US 34** dan **US 35** · **meta tanpa AC: satu** — **US 39**.

⚠️ **US 34 dan US 35 sengaja TIDAK dipaksa tertutup.** Menambah AC untuk keduanya berarti
memutuskan hal yang belum diputuskan — untuk US 34, **ketelitian Pega mana yang sah** *(butir
`[terbuka]` 8)*; untuk US 35, **bentuk daftar "disimpan versus dihitung"** yang belum diminta siapa
pun. ⛔ Keduanya dicatat di sini apa adanya.

---

## Out of Scope

Di luar lingkup berkas ini, dan **bukan** karena tidak penting:

1. **Komite untuk lini lain** — Life, Non Prop, dan Fac In. `[keputusan work owner]` lingkup satu
   modul.
2. **Modul Claim Prop.** Termasuk **lima activity penghitung** yang berkelas kasus klaim meski
   berkasnya ada di folder ekspor modul ini (bab 7).
3. **Rule pembaca tabel riwayat akseptasi** yang ada di modul Fac In dan Treaty In. Sisi **tulis**
   milik modul ini sudah selesai dan benar.
4. **Isi fungsi tersimpan Oracle untuk kurs.** `[keputusan work owner]` tidak diminta ke DBA.
5. **Hak akses per activity.** `[keputusan work owner]` diabaikan.
6. **Lima lompatan menggantung, properti jenis treaty, nomor klaim di balik gerbang mati, dan titik
   lompat darurat.** `[keputusan work owner]` dibuang, tidak dipindahkan.
7. **Bentuk tabel baru selain sembilan kolom yang sudah dikunci.** Termasuk tempat bagi 19 properti
   baris penyesuaian — itu keputusan work owner yang belum diambil.
8. **Urutan efek keluar yang baru.** Sengaja diubah, tetapi **belum diputuskan**.
9. **Sensus kode arah dan wadah gerbang di seluruh korpus.** Terverifikasi tetapi milik lintas
   modul.

---

## Butir `[terbuka]` — daftar penuh

⛔ **Tidak satu pun ditutup di berkas ini.** Setiap butir **ikut terbawa** ke bab masing-masing.

**Jumlah akhir: 13.** *(12 → 13: butir **14** ditambahkan 2026-09-19 — kasus komite lama hasil
migrasi yang penunjuknya salah alamat.)* ⛔ **Nol butir lama ditutup.**

> ⛔ **RALAT 2026-09-19** — kalimat lamanya **dikutip, tidak dihapus**: *"**Jumlah akhir: 12.**
> *(13 → 12: butir "12 digit" **DITUTUP** `[keputusan work owner]` 2026-09-19 — ikuti bentuk kolom
> yang ada apa adanya, risiko gagal simpan diterima sadar.)*"* Penutupan butir "12 digit" itu
> **tetap berlaku**; yang berubah hanyalah **jumlahnya**, karena satu butir baru masuk.
>
> ⚠️ **Penomoran tetap melompat:** nomor **10** memang tidak pernah dipakai, dan nomor butir baru
> adalah **14**, melanjutkan nomor tertinggi — bukan mengisi lubang nomor 10.

| # | Butir | Bab | Menunggu |
| --- | --- | --- | --- |
| 1 | Kolom akun operator di tabel riwayat akseptasi tidak pernah diisi jalur komite | 9 | korpus / DBA |
| 2 | Pewarisan kelas kerja — tidak ada di ekspor ini, permintaan sudah dirumuskan | 1 | work owner |
| 3 | Arti dua jenis perulangan, dan halaman apa yang diulang | 3 | korpus |
| 4 | **Berapa kali** perulangan berputar — menentukan berapa kali angka berubah dan berapa kali kiriman Kasir terjadi | 3 · 7 · 8 | korpus |
| 5 | Urutan pemeriksaan bila dua keluarga gerbang sama-sama terisi | 3 | korpus |
| 6 | Unggah berkas dan email **tanpa penanganan gagal** | 8 | work owner |
| 7 | Catatan pengembang **"FIX ERROR HANDLING"** yang tidak berjejak | 8 | korpus |
| 8 | Ketelitian pembagian persen di Pega tidak seragam — mana yang sah | 7 | work owner |
| 9 | Dari mana layar kedua dibuka | 4 | korpus |
| **11** | ⭐ **BARU** — seberapa besar **selisih akibat persen 10 desimal menjadi 8** saat disimpan; menumpuk lewat perulangan sebelum sampai ke Kasir | 7 | korpus / work owner |
| **12** | ⭐ **BARU** — **baris mana di tabel alamat layanan** yang menunjuk lingkungan uji | Testing | work owner |
| **13** | ⭐ **BARU** — **sasaran uji untuk urutan efek keluar**, menunggu urutan barunya ditetapkan | Testing | work owner |
| **14** | ⭐ **BARU 2026-09-19** — kasus komite **LAMA hasil migrasi** yang penunjuk baris penyesuaiannya **sudah salah alamat atau menunjuk baris yang tidak ada** — apa yang dilakukan aplikasi saat kasus itu dibuka: **menolak terang-terangan**, atau **diam seperti Pega**. ⚠️ **Kunci beku (bab 11) TIDAK menutup butir ini karena tidak berlaku surut.** | 2 · 11 · Migrasi | work owner |

> **RALAT 08-10-2026 — penilaian ulang tiga butir** (prompt §7 butir 9). Yang ditutup hanya yang terbukti:
>
> - **Butir 4 DITUTUP.** Perulangan berputar sekali per tingkat tangga: setiap Submit menjalankan
>   `KomitePostAdjustment` sekali, S40 menaikkan `KomiteCount` (prakondisi nonaktif), dan Decision `KomiteLoop`
>   (`IsKomiteLoop`: `.AcceptStatus="1" AND .KomiteCount <= .KomiteLoop`) mengembalikan assignment ke KomiteRouter
>   sampai tingkat akhir. Langkah uang / Kasir (S14-S34) bergerbang `KomiteCount == TotalKomite && AcceptStatus == 1`
>   - terjadi **sekali** per kasus komite. Bukti: `Activity/KomitePostAdjustment.xml` S13-S40, `Flow/KomiteTreaty_Flow.xml`.
> - **Butir 11 DITUTUP.** Kolom yang ada `NUMBER(38,10)` menyimpan 10 desimal; persen 10 desimal tidak terpangkas
>   menjadi 8.
> - **Butir 3 TETAP.** Arti blok `ULANG[1]` (S16, HitServiceToKasirKMT_Act S14) tidak dinyatakan ekspor; perilakunya
>   ditiru (sekali jalan) tanpa menafsirkan nama.
> - **Butir 9 TETAP** (pintu masuk `ViewDetailInterest`).
>
> Jumlah butir terbuka: 13 → **11**.

#### ⭐ Apakah ketiga belas butir ini MENGHAMBAT — penilaian 2026-09-19

⛔ **Penilaian, bukan penutupan.** **Nol butir ditutup di sini**; yang berubah hanyalah bahwa
statusnya kini **dinyatakan**, bukan dibiarkan tidak diketahui.

| Golongan | Butir | Artinya |
| --- | --- | --- |
| ⛔ **MENGHAMBAT PEMBANGUNAN** — **satu** | **9** | **dari mana layar kedua dibuka** — tiket **03** tidak dapat membangun jalan masuknya, dan ⛔ mengarang tombol dilarang |
| ⚠️ **MENGHAMBAT PENGUJIAN** — **lima** | **4 · 6 · 12 · 13 · 14** | dapat dibangun, **belum dapat diuji tuntas**: berapa kali perulangan berputar *(4)* · perilaku saat unggah berkas dan email gagal *(6)* · baris lingkungan uji di tabel alamat layanan *(12)* · sasaran uji urutan efek keluar *(13)* · perilaku kasus lama hasil migrasi yang penunjuknya salah *(14)* |
| ✅ **TIDAK MENGHAMBAT** — **tujuh** | **1 · 2 · 3 · 5 · 7 · 8 · 11** | catatan dan risiko; perilaku yang ditiru adalah perilaku yang **berjalan**, dan tiketnya tetap dapat dikerjakan |

⚠️ **Modul ini BERBEDA dari Claim Prop.** Di sana pernyataannya *"`[terbuka]` yang memblokir:
NIHIL"*. Di sini **enam butir menghambat** — satu menghambat pembangunan, lima menghambat
pengujian. ⛔ **Jangan menyalin kalimat "memblokir NIHIL" ke modul ini.**

⚠️ **Akibat langsungnya pada urutan kerja:** tiket **03** tidak dapat dinyatakan tuntas sebelum
butir **9** terjawab, dan tiket **10 · 11 · 12** tidak dapat diuji dengan layanan sungguhan
sebelum butir **12** terjawab — padahal `[keputusan work owner]` 2026-09-18 mewajibkan justru itu.

#### ✅ Yang DITUTUP — satu

**Peringatan famili gerbang ketiga.** Di draf pertama ia menggantung angka kolom layar. Pemeriksaan
sudah dilakukan di berkas modul ini, `[terverifikasi]` angkanya tidak berubah, dan
`[keputusan work owner]` 2026-09-18 **menutupnya**. **Angka 93 · komite 62 · klaim induk 31 SAH**
(bab 4). Ia **tidak pernah masuk tabel 13 di atas**, jadi jumlahnya tidak berubah.

---

## Further Notes

### Langkah berikutnya sesudah spec

`[keputusan work owner]` 2026-09-18 — **tiket ditulis sebagai berkas `issues/*.md`**, seperti
modul-modul sebelumnya di proyek ini. ⛔ **Bukan dipublikasi ke issue tracker.**

### Ke mana 19 properti baris penyesuaian pergi

`[keputusan work owner]` 2026-09-18 — **dipilah dua**:

| Berapa | Perlakuan |
| --- | --- |
| **17** yang **hanya dibaca** | **DIBACA dari baris penyesuaian milik Claim Prop** — **tidak disalin** ke tabel komite |
| **2** yang **disunting penyetuju** — usul tutup klaim dan usul cadangkan klaim | **DISIMPAN di tabel komite** |

⛔ **Tempat persisnya kedua kolom itu ditetapkan di berkas struktur tabel**, bukan di sini.
Lihat `STRUKTUR-TABEL-KOMITE-CLAIM-PROP.md`.

⚠️ **Catatan yang tidak boleh hilang** `[terverifikasi]`: di Pega kedua penanda itu **punya dua
tempat**. Sesudah penyetuju mengisinya di layar, langkah **11** `KomitePostAdjustment`
**menyalinnya ke kasus klaim induk** dengan nama lain — penanda *tutup berkas* dan penanda
*klaim dicadangkan*. Langkah itu **tanpa gerbang**, jadi selalu berjalan.

### Mengapa delapan ronde

Tiga kali berturut-turut kesimpulan **"tidak ada"** ternyata salah, karena ditarik dari wadah yang
belum dibaca:

| Ronde | Yang dikira tidak ada | Ternyata |
| --- | --- | --- |
| 2 → 3 | *"sel tidak pernah bergerbang sendiri"* | ada **15** gerbang di sub-halaman milik sel |
| 3 → 4 | *"566 baris syarat, itu semuanya"* | ada **552 baris lagi** di keluarga gerbang kedua |
| 5 → 6 | *"alamat layanan diambil tanpa saringan"* | saringannya ada, di wadah yang tidak diparse |

⭐ **Pelajarannya, dan ia berlaku untuk siapa pun yang memeriksa ulang berkas ini:**
**"tidak ketemu" hanya sah bila seluruh kedalaman sudah dicari.** Empat keluarga wadah aturan di
tingkat langkah kini sudah dibaca habis (lampiran).

### Nama yang menyesatkan

`[terverifikasi]` Tiga hal yang namanya tidak menggambarkan isinya — jangan percaya nama:

- Penyiap layar yang namanya menyebut daftar penyetuju **tidak menyentuh daftar penyetuju sama
  sekali**; ia menjumlahkan penyesuaian per mata uang.
- Sebuah kolom bernama seperti mata uang asing **isinya hasil kali dengan kurs**, bukan mata uang
  itu.
- Dua berkas bernama sama di modul ini adalah **dua rule berbeda pada kelas berbeda**.

### Nama orang yang tertulis di dalam aturan

⚠️ `[terverifikasi]` Dua tempat menuliskan **nama orang dan jabatan langsung di dalam aturan** —
satu pada jalur penutupan komite, satu pada pencarian batas wewenang. Bila orangnya pindah, jalur
itu **diam-diam salah alamat** (ronde 1 §9.4). **Dicatat sebagai temuan; perbaikannya meja work
owner.**

---

## Lampiran — Aturan baca ekspor Pega

*Lampiran ini ada supaya orang berikutnya bisa memverifikasi sendiri setiap angka di berkas ini,
tanpa mengulang delapan ronde. Semuanya `[terverifikasi]`.*

### 1. Struktur langkah dan jebakan bersarang

Langkah bersarang berada **di dalam** elemen induknya. ⚠️ Gerbang milik induk **terserialisasi
sesudah** anak-anaknya, jadi **pembacaan urutan baris pasti salah atribusi**. **Hanya parser XML
yang aman.** Di layar, sarangnya mencapai **lima tingkat**.

### 2. Empat keluarga wadah aturan pada satu langkah

| Keluarga | Isinya | Diuji kapan |
| --- | --- | --- |
| **1 — prasyarat** | gerbang **sebelum** langkah jalan | sebelum |
| **2 — transisi** | ke mana lanjut **sesudah** langkah jalan; syaratnya sering *"langkah barusan gagal"* | sesudah |
| **3 — perulangan** | apakah langkah berulang, dan jenisnya | — |
| **4 — parameter metode** | isi panggilan; **bermuka dua** (lihat aturan 5) | — |

Setiap langkah **selalu membawa keempat wadah**; yang membedakan hanya terisi atau tidak.
⚠️ Keluarga 1 dan 2 punya **saklar hidup-mati masing-masing**, dengan logika tiga nilai yang sama:
`true` berlaku · `false` dimatikan · kosong tidak pernah ada.

### 3. Enam kode arah `[data work owner]`

| Kode | Nama di Pega | Arti |
| --- | --- | --- |
| **1** | Jump to Later Step | lompat ke langkah bertanda; sasarannya di kolom parameter |
| **2** | Continue Whens | periksa baris syarat berikutnya; habis baris → **langkah dijalankan** |
| **3** | Skip Step | langkah **tidak** dijalankan |
| **4** | Exit Iteration | putus perulangan |
| **5** | Skip Whens | **berhenti memeriksa**, langkah **langsung dijalankan** |
| **6** | Exit Activity | hentikan activity |
| *(kosong)* | — | tidak ada tindakan dipilih |

⭐ **Baris syarat adalah RANTAI, dibaca dari atas.** *Continue Whens* + *Skip Step* menyusun **DAN**;
satu *Skip Whens* mengubahnya jadi **ATAU**.
⭐ **Tanda langkah** dibaca dari kolom paling kiri; nilai `//` berarti **langkah di-remark**.

### 4. Letak nama rule basis data

⚠️ Langkah pengambilan data **tidak menyebut nama rule di satu tempat**. Nama rule terpecah **tiga
bagian** di parameter panggilan — **kelas**, **jalur akses**, dan **jenis permintaan** — yang bila
dirangkai dengan tanda seru menghasilkan `pxInsName` persis.

### 5. Letak saringan pencarian data

⚠️ Saringan **bukan satu tag**. Satu baris saringan adalah **satu baris data** yang membawa
penanda-pilih, nama kolom, pembanding, dan nilai **bersama-sama** — dan ia duduk di **wadah
keluarga 4**, wadah yang **sama** dengan tempat penugasan properti biasa disimpan.

⭐ **Wadah itu bermuka dua.** Untuk penugasan ia daftar pasangan nama-nilai; untuk pencarian data ia
daftar kolom dan saringan. **Membacanya sebagai satu jenis saja adalah cara paling mudah kehilangan
separuh isinya.**
**Membedakannya:** baris ber-pembanding **dan** bernilai = **saringan**; baris berpenanda-pilih
tanpa keduanya = **kolom yang diambil**.

### 6. Letak parameter sebuah activity

⚠️ Parameter dideklarasikan di sebuah wadah yang isinya **dokumen XML ter-escape di dalam satu
simpul teks**. Parser XML biasa membacanya sebagai teks polos; ia harus **di-unescape lalu diurai
ulang**. **18 dari 35** activity modul ini memakainya.

### 7. Famili gerbang layar — **dua yang hidup, satu yang tidak**

| Famili | Tingkat | Saklarnya | Syaratnya di | Hidup di modul ini |
| --- | --- | --- | --- | --- |
| **A** | **sel** | penanda tampil bernilai *"lainnya"* | sub-halaman milik sel | **15** |
| **B** | **layout** | pilihan tampil bernilai *bersyarat* | properti syarat layout | **9** |
| **C** — famili ketiga | **sel** | idem famili A, tetapi di **wadah bawaan** yang terpisah | wadah itu sendiri | ⭐ **0** |

**Penanda gerbang mati: `NEVER` dan `1=2`.**

⭐ **Famili C ditemukan modul Claim Prop, dan sudah diperiksa di berkas modul ini.**
`[terverifikasi]` Di layar utama: wadah famili A memuat **354 blok** — **162** berkelas
*HeaderElements* dan **192** berkelas *UserData* — dengan **15 hidup**; wadah famili C memuat
**1 blok** berkelas *UserData*, **nol hidup**.

⛔ **Dua akibat, dan keduanya penting bagi siapa pun yang menghitung ulang:**

1. **Hitungan 15 itu sudah menyapu KEDUA kelas.** Kelas *UserData* **tidak pernah terlewat**, ia
   hanya tidak pernah disebut namanya di ronde 1–8.
2. **Famili C ada satu dan tidak hidup**, jadi **angka kolom layar tidak berubah**.

⚠️ Butirnya **tetap `[terbuka]`** (bab 4) — yang belum adalah pemeriksaan yang sama di berkas
milik Claim Prop, dan itu **di luar lingkup modul ini**.

### 8. Huruf besar-kecil

⭐ **Ekspor Pega tidak konsisten kapitalnya.** `[terverifikasi]` Menyaring nama metode dengan huruf
persis memberi **17** langkah Java; tidak peka huruf memberi **20** — ketiga selisihnya ditulis
huruf kecil.

> **Aturan: setiap penyaringan nilai harus TIDAK PEKA huruf besar-kecil. Dan lebih baik menyaring
> dari ISI daripada dari LABEL.**

`[terverifikasi]` Dua belas penyaring lain di ronde 1–8 sudah diperiksa ulang: **nol yang
terpengaruh**.

### 9. Identitas rule

**`pxInsName`, berbentuk `KELAS!NAMA`.** ⛔ **Kecocokan nama berkas bukan bukti rule yang sama** —
kesalahan ini sudah terjadi dua kali di proyek ini.
Jejak salinan rule ada di properti terpisah; ia **hanya** untuk membaca asal-usul, **bukan**
identitas.
