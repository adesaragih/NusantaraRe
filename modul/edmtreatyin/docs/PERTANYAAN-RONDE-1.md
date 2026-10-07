# Pertanyaan — EDM Treaty In, ronde 1
## Migrasi Treaty Inward — Realisasi & Endorsement · modul *EDM Treaty In*

⭐ **Sebelas pertanyaan baru, P50–P60.** ⭐ Penomoran **berlanjut** dari NB Treaty In *(P1–P49;
P41, P48, dan — 2026-09-22 — **P18** ditarik)* — ⭐ satu nomor berarti satu pertanyaan **di seluruh
proyek**.

> ⭐⭐ **Keadaan 2026-09-22 sore: kesebelasnya TERJAWAB.** P50–P60 tertutup seluruhnya — sepuluh
> dijawab pemiliknya, satu *(**P60**)* **ditutup dari korpus** sesudah isi langkah terbukti terbaca.
>
> ⭐ **Keadaan terukur modul ini kini ada di `KEADAAN-EDM-TREATY-IN.md`** — itu sumber angka,
> bukan berkas ini.
>
> ⛔⛔ **EMPAT PERNYATAAN RONDE 1 TERBUKTI KELIRU.** `grilling-ronde-1.md` **tersegel dan tidak
> disunting**; koreksinya hidup di `KEADAAN-EDM-TREATY-IN.md` Bab 6, dan ringkasnya:
>
> | Bunyi ronde 1 | Yang benar |
> | --- | --- |
> | ⛔ *"isi langkah penetapan nilai tidak terekspor"* | ⭐ **ada** — di `PropertiesName`/`PropertiesValue` **tanpa awalan `py`**. **379 dari 380** langkah membawa isinya *(99,7 %)*, **1.309** pasangan terisi |
> | ⛔ *"pemeriksaan keberhasilan hanya memeriksa FacIn"* | ⭐ memeriksa **keduanya**, tiga cabang — sudah diralat di **P50** |
> | ⛔ *"`BusinessType_DeT` tinggal 9 baris dari 34"* | ⭐ **36 baris utuh**; terverifikasi ulang dari **data produksi** |
> | ⚠️ *"beda hanya metadata ekspor"* untuk 69 aturan | ⭐ **bertahan untuk ke-69-nya** — identik perilaku **dan** catatan. Yang beda hanya **4** dari 73, dan ketiga berkas yang catatannya berbeda seluruhnya ada di dalam keempat itu |
>
> ⚠️ **Angka 378 langkah dan 1.062 pasangan yang pernah beredar DITARIK** — enam rekonstruksi
> dicoba, tidak satu pun menghasilkannya. Yang berlaku **380** dan **1.309**, sepakat dua cara.
⛔ **Berkas ini hanya berisi pertanyaan** — temuan dan buktinya ada di `grilling-ronde-1.md`.

**Cara menjawab:** tulis di bawah tiap pertanyaan. Baris `rujukan:` boleh diabaikan. ⚠️ *"Tidak
tahu"* adalah jawaban yang sah. ⛔ Nomor jangan diubah.

| Pemilik | Pertanyaan |
| --- | --- |
| ⛔ **Product & Underwriting** | **P50 · P51 · P54 · P55 · P56 · P58** |
| ⛔ **pengembang Pega lama** | **P52 · P57** — ⭐ ~~P60~~ **ditutup dari korpus 2026-09-22**, tidak perlu dijawab |
| ⛔ **IAM** | **P53** |
| ⛔ **DBA** | **P59** |
| Finance · pemilik export Pega | ⭐ **tidak ada ronde ini** |

⚠️ **Yang TIDAK ditanyakan ulang**, karena sudah dijawab untuk NB Treaty In dan jawabannya berlaku
di sini: **P23** *(keterangan aturan lawan syarat yang dijalankan — yang dijalankan benar)* ·
**P25** *(penunjukan antrean menurut nomor urut — tidak dimigrasi)* · **P13** *(tangga tiga
jenjang)* · **P29** *(dokumen JSON dibuang)*.

---

# Untuk Product & Underwriting

## P50 — Data dikirim DUA KALI ke sistem produksi. Kapan yang kedua terkirim?

Ketika sebuah endorsemen selesai, sistem lama mengirimkan datanya ke sistem produksi
**dua kali berturut-turut**. Yang pertama untuk **bisnis masuk**, yang kedua untuk **bisnis
keluar**. ⭐ Yang pertama **selalu** dikirim; ⚠️ yang kedua **hanya bila suatu syarat terpenuhi** —
dan syarat itu tidak dapat kami baca.

Lebih jauh lagi: pemeriksaan "apakah pengirimannya berhasil" **hanya memeriksa yang pertama**.
⛔ Kami tidak menemukan satu pun pemeriksaan untuk yang kedua.

**Konteks:** bila pengiriman kedua gagal tanpa diperiksa, sebagian data endorsemen tidak sampai
ke sistem produksi dan **tidak ada yang tahu**.

**Bentuk jawaban yang diharapkan:** **(1)** kapan pengiriman kedua perlu dilakukan — endorsemen
jenis apa · **(2)** apakah kegagalan pengiriman kedua **harus** ditangani seperti yang pertama,
atau memang boleh diabaikan.

**Dampak bila salah:** ⛔ data penempatan keluar **tidak sampai** ke sistem produksi, atau
sebaliknya terkirim untuk endorsemen yang tidak seharusnya.

rujukan: `Activity\serviceInsertArasapas_act.xml` langkah 6 *(tanpa gerbang)* dan 7 *(bergerbang)*;
`When\IsSuccessHitService.xml` hanya menguji status **FacIn**;
`RDBList\INSERTJSON_JSONPOLISMONITORING_FACIN.xml` — namanya pun menyebut FacIn saja

**Jawaban:**

> **Untuk EDM Treaty In tidak ada FacOut. Pengiriman kedua TIDAK dijalankan.**
> `[keputusan work owner]` 2026-09-22
>
> Langkah 7 — `Connect-REST` *"Hit service Arasapas 2 - FACOUT"* — **tidak pernah berjalan** pada
> berkas EDM Treaty In. Konversi FacOut bukan bagian modul ini. **Tidak dimigrasi.**
>
> **Bukti dari korpus yang bersesuaian.** `[terverifikasi]` Bentuk syarat yang dijalankan
> `When\IsSuccessHitService.xml`:
>
> ```
> [StatusService.StsKonversiFacIn ][=][1]
> [StatusService.StsKonversiFacOut][=][""]
> [StatusService.StsKonversiFacOut][=][1]
> ```
>
> Berhasil = FacIn = 1 **dan** (FacOut kosong **atau** FacOut = 1). Cabang **"FacOut kosong"**
> itulah yang menampung keadaan EDM Treaty In: tidak pernah dicoba, karena itu kosong, dan itu
> tetap dihitung berhasil. Keputusan work owner dan bentuk rule-nya bersesuaian.
>
> Penguat: di seluruh modul EDM Treaty In, `FacOut` hanya muncul di **dua berkas** — aktivitas yang
> memuat langkahnya, dan rule pemeriksa di atas. Nol layar, nol tabel pemantauan, nol rule lain.
>
> ⛔ **KOREKSI 06-10-2026 (XML) — penguat ini tidak lengkap; keputusan WO di atas tidak diubah.** Gerbang
> langkah 7 adalah `When/IsFacRetro` (`OfferFacIn.IsFacRetro = 1`), dan penanda itu **diset di jalur EDM**:
> `Activity/InputPolicyTreatyEDMDetail_NP.xml` dan `…_NP_AdjPremi.xml` langkah 3 — `OfferFacIn.IsFacRetro = "1"`
> bila `TreatyIn.FacultativeShare > 0` (dipanggil `SetValueEDM_Act` 9/10 ← `EDMChooseBusiness_Act` 1). Sapuan
> kata *"FacOut"* melewatkannya. ⇒ *"tidak pernah berjalan"* tidak terbukti; **butir WO**: apakah P50 tetap
> berlaku untuk master ber-`FacultativeShare > 0`. Log: `KOREKSI-DOKUMEN-2026-10-06.md`.
>
> ---
>
> ⛔ **RALAT atas laporan ronde 1.** `[penyimpangan sadar]`
>
> Ronde 1 menulis: *"pemeriksaan apakah pengirimannya berhasil hanya memeriksa yang pertama; kami
> tidak menemukan satu pun pemeriksaan untuk yang kedua."*
>
> **Keliru.** `IsSuccessHitService` **memeriksa keduanya**, dengan tiga cabang seperti di atas.
> Yang terbaca ronde 1 adalah `pyConditionString` — keterangan yang ditulis manusia — yang memang
> hanya berbunyi `StsKonversiFacIn = 1`.
>
> Ini **persis jebakan yang sudah diputuskan di P23**: bila keterangan berbeda dari syarat yang
> dijalankan, **yang dijalankan benar**. Keputusan itu berlaku juga di modul ini.
>
> ---
>
> **Yang mengikat sistem baru:**
>
> | | |
> | --- | --- |
> | Langkah FacOut | **tidak dibangun** untuk EDM Treaty In |
> | Pemeriksaan keberhasilan | cukup menguji FacIn |
> | Baris pemantauan | cukup FacIn — tidak perlu padanan FacOut |
>
> `[terbuka]` Aktivitas `serviceInsertArasapas_act` **dipakai bersama** empat modul —
> `EDM Treaty In`, `Endorsment Fac In`, `NB FacIn`, `RNW Fac In`. Langkah FacOut yang tidak dipakai
> di sini **kemungkinan dipakai modul Fac**. Ketika modul Fac digarap, pertanyaan "kapan FacOut
> dikirim" **wajib dibuka lagi di sana** — di sini ia tertutup, bukan terjawab untuk semua.

---

## P51 — Bila pengiriman gagal, sistem MENGHAPUS data produksi. Sengaja?

⛔ Di dalam jalur penanganan galat terdapat satu langkah yang **menghapus data produksi**, dan
catatannya menyatakan itu apa adanya: *"hapus produksi jika error"*. Penghapusannya dijalankan
lewat sebuah program di basis data.

**Konteks:** menghapus adalah tindakan yang **tidak dapat dibatalkan**. Kami perlu tahu apa yang
dihapus — satu baris, satu endorsemen, atau lebih luas — dan apakah itu memang yang dimaksudkan
ketika pengiriman gagal.

**Bentuk jawaban yang diharapkan:** pilihan ganda + keterangan — **(a)** benar, data produksi yang
sudah tersisip harus dibatalkan bila konversi gagal · **(b)** seharusnya ditandai gagal, bukan
dihapus · **(c)** tidak tahu, perlu diperiksa.
⭐ Dan yang paling membantu: **seberapa sering ini terjadi** dalam setahun terakhir.

**Dampak bila salah:** ⛔⛔ data produksi **hilang permanen** karena kegagalan sementara di
jaringan, atau sebaliknya data yang salah **tertinggal** di sistem produksi.

rujukan: `Activity\serviceInsertArasapas_act.xml` langkah 12, `RDB-List`, catatan
*"DELETE PRODUKSI JIKA ERROR (PROCEDURE)"*, bergerbang — ronde 1 §C.2

**Jawaban:**

> **(a) Benar — data produksi HARUS dihapus bila konversi gagal.** `[keputusan work owner]` 2026-09-22
>
> Dan satu batasan yang menyertainya: **penghapusan hanya boleh lewat satu jalur itu.** Tidak ada
> kode lain di sistem baru yang boleh menghapus data produksi.
>
> **Jejak yang sudah terverifikasi.** `[terverifikasi]`
>
> ```
> Activity\serviceInsertArasapas_act  langkah 12  RDB-List
>    -> RDBList\DeleteDataProduction
>    -> BEGIN POOLDATA.PEGA_DELETE_ERROR_KONVERSI(
>            {pyWorkPage.pzInsKey},                 -- kunci kasus
>            {pyWorkPage.Quotation.BusinessFac},    -- lini bisnis
>            {OutputData.HASIL10 out} ); END;
> ```
>
> Lingkup penghapusan ditentukan **dua parameter**: kunci kasus dan lini bisnis. Bukan satu baris,
> melainkan satu kasus untuk satu lini. Berapa baris yang berarti hanya dapat dijawab naskah
> procedure-nya.
>
> Dipakai **empat modul**: `EDM Treaty In`, `Endorsment Fac In`, `NB FacIn`, `RNW Fac In` —
> pola yang sama di seluruh jalur konversi ke produksi.
>
> ---
>
> **Yang mengikat sistem baru.**
>
> | Ketetapan | |
> | --- | --- |
> | Kegagalan konversi **menghapus** data produksi yang sudah tersisip | ditiru apa adanya |
> | Penghapusan hanya lewat **satu jalur** | satu fungsi repository, tidak dipanggil dari tempat lain |
> | Tidak ada kode lain yang boleh menghapus data produksi | ditegakkan lewat uji |
>
> `[penyimpangan sadar]` Mandat proyek melarang pemanggilan stored procedure, sehingga logika
> `PEGA_DELETE_ERROR_KONVERSI` **ditulis ulang di Go**. Yang dipertahankan adalah **aturannya** —
> satu jalur penghapusan, dipakai hanya oleh penanganan gagal-konversi — bukan pemanggilan
> procedure-nya. Bila yang dimaksud work owner justru mempertahankan procedure itu sendiri,
> ketetapan ini wajib diralat.
>
> ~~`[terbuka]`~~ ✅ **`PEGA_DELETE_ERROR_KONVERSI` — naskahnya DITERIMA 2026-09-22.** Semula
> tercatat sebagai *"procedure KEEMPAT yang belum diminta"*; surat **P1** kini memuat keempatnya dan
> **sudah terjawab**. Lihat `modul/nbtreatyin/docs/PERTANYAAN-untuk-DBA.md` §P1 nomor 4.
> *(Koreksi 06-10: jalur semula `..\nb-treaty-in\…` terpotong menjadi dua baris karena `\n` terbaca sebagai baris baru.)*
>
> ~~`[terbuka]`~~ ✅ **Pertanyaan `COMMIT` terjawab: procedure ini meng-commit DI DALAM**, sesudah
> seluruh penghapusan berhasil. Galat pada tabel mana pun menghasilkan `ROLLBACK` dan penghentian.
> ⭐ **Penghapusannya utuh atau tidak sama sekali**, dan **tidak** bergantung pada transaksi
> pemanggil.

---

## ⭐⭐ PENAJAMAN 2026-09-22 — penghapusannya BERSYARAT

> ⚠️ **Ini penajaman, BUKAN pembatalan.** `[terverifikasi]` `[data DBA]`
> Keputusan work owner *"(a) Benar — data produksi HARUS dihapus bila konversi gagal"*
> **tetap berlaku utuh.** Yang bertambah adalah **penjaga** yang tidak terlihat dari sisi Pega,
> sebab ia ada di dalam naskah procedure.

### Yang baru diketahui

⭐⭐ **Penghapusan tidak menyentuh data yang sudah berhasil dikonversi.** Penjaganya:

```
hanya menghapus bila   STS_KONVERSI IS NULL   atau   STS_KONVERSI <> '1'
lingkup selalu         WHERE IDPEGA = <satu kasus>
```

| `Bisnis` | Tabel yang dihapus |
| --- | --- |
| `'T'` — treaty | `JSON_POLIS` · `TREATYINPRODUCTION` · `TREATYINPRODUCTION_BACKUP` |
| `'F'` — fakultatif | `JSON_POLIS` · `FACINPRODUCTION` · `FACINPRODUCTION_BACKUP` · `FACOUTPRODUCTION` |

⭐ **Pertanyaan lama *"berapa baris yang terhapus"* kini terjawab:** baris milik **satu kasus**
*(`IDPEGA`)*, pada **tiga tabel** untuk treaty dan **empat** untuk fakultatif — dan **nol baris**
bila kasus itu sudah bertanda konversi berhasil.

### Yang berubah bagi sistem baru

| Sebelum penajaman | Sesudah |
| --- | --- |
| *"kegagalan konversi menghapus data produksi yang sudah tersisip"* | ⭐ **ditambah syarat**: hanya bila `STS_KONVERSI` belum bernilai `1`. Penulisan ulang di Go **wajib membawa penjaga ini**; menghilangkannya membuat kegagalan berulang dapat menghapus konversi yang sudah sah |
| lingkup *"satu kasus untuk satu lini"* | ✅ **ditegaskan tepat** — `WHERE IDPEGA =`, satu kasus |
| ketergantungan transaksi belum pasti | ✅ **meng-commit di dalam**, utuh atau tidak sama sekali |

⚠️ `[penyimpangan sadar]` yang sudah tercatat **tetap berlaku**: logikanya ditulis ulang di Go,
bukan memanggil procedure. Penjaga `STS_KONVERSI` menjadi **bagian dari logika yang ditulis ulang
itu**, bukan alasan mempertahankan procedure-nya.

⛔ **Satu butir `[terbuka]` BARU lahir dari naskah ini** — tabel `_BACKUP` ikut terhapus.
Dicatat pada lembar `[work owner]`, **tidak ditutup di sini**.
>
> `[terbuka]` **Berapa kali jalur ini berjalan dalam setahun terakhir** — dihitung dari baris
> ber-`sts_konversi = 9` pada `JSON_POLIS_MONITORING`. Tidak menahan, tetapi menentukan seberapa
> mendesak uji regresinya.
>
> **Kegagalan WAJIB diberitahukan kepada pengguna lebih dulu, baru penghapusan dijalankan.**
> `[keputusan work owner]` 2026-09-22
>
> Urutannya mengikat:
>
> ```
> 1. konversi gagal
> 2. pesan galat DITAMPILKAN kepada pengguna     <- wajib, tidak boleh dilewati
> 3. jalur penghapusan dijalankan
> 4. kegagalan dicatat di tabel pemantauan dan surel
> ```
>
> `[penyimpangan sadar]` Ini **berbeda dari sistem lama**. Ronde 1 melaporkan dua langkah pada jalur
> galat tidak aktif — pemasang pesan galat (langkah 10) dan penutup paksa kasus (langkah 18) —
> sehingga galatnya hanya tercatat di tabel pemantauan dan surel, tidak tampak di layar.
>
> ⚠️ Laporan itu **belum dapat dibuktikan ulang**: isi gerbang langkah
> (`pyStepsPreCondParamsWhen`) kosong di ekspor untuk kedelapan langkah bergerbang, sehingga hanya
> terlihat *bahwa* ada gerbang, bukan bunyinya. Status temuannya `[dugaan]`.
>
> **Keputusan ini berlaku terlepas dari benar tidaknya temuan itu.** Bila ternyata sistem lama
> memang sudah menampilkan pesannya, sistem baru hanya meniru. Bila tidak, sistem baru
> memperbaikinya. Keduanya menghasilkan perilaku yang sama, dan itu yang dikehendaki.
>
> **Yang mengikat pengujian:** test yang menemukan penghapusan berjalan **tanpa** pesan galat
> tampil lebih dulu **gagal**. Pencatatan ke tabel pemantauan dan surel **tidak menggantikan**
> pesan di layar — ketiganya berjalan, bukan salah satu.

---

## P54 — Endorsemen yang dibuat sesudah tanggal 25 diberi nomor bulan BERIKUTNYA. Kenapa?

Sistem lama memeriksa tanggal hari ini ketika membentuk nomor endorsemen: bila sudah **lewat
tanggal 25**, nomornya dibuat untuk **bulan berikutnya**, bukan bulan berjalan.

**Konteks:** ini aturan periode produksi, dan ⭐ **tidak ada padanannya** pada penomoran realisasi
polis baru. Kami perlu tahu apakah ia masih berlaku, dan apakah tanggalnya memang 25.

**Bentuk jawaban yang diharapkan:** **(1)** apakah aturan ini masih berlaku · **(2)** apakah
batasnya memang tanggal **25**, atau mengikuti tanggal tutup buku yang dapat berubah ·
**(3)** apakah ia berlaku untuk **semua** jenis endorsemen.

**Dampak bila salah:** ⛔ endorsemen masuk **periode produksi yang salah**, dan itu menggeser
pengakuan premi antar bulan.

rujukan: `Activity\GeneratePolicyNoTreatyAddendum_Act.xml`, catatan langkah
*"when not above 25 in that month"* dan *"Determine if it's past 25 on current month, if true move
to next month"* — ronde 1 §B.3

**Jawaban:**

> **Tanggal tutup buku diambil dari `POOLDATA.TANGGAL_CLOSING`, bukan angka yang ditulis mati.**
> `[keputusan work owner]` 2026-09-22
>
> Bila tanggalnya berubah di tabel itu, **seluruh modul ikut berubah**. Tidak ada satu pun nilai
> tanggal tutup buku yang ditanam di dalam kode.
>
> Angka **25** pada catatan langkah — *"when not above 25 in that month"* dan *"Determine if it's
> past 25 on current month, if true move to next month"* — adalah **keterangan**, bukan nilai yang
> dijalankan. Sejalan dengan keputusan **P23**: bila keterangan berbeda dari yang dijalankan, yang
> dijalankan benar.
>
> **Bukti dari korpus.** `[terverifikasi]` `Activity\GeneratePolicyNoTreatyAddendum_Act.xml`
> memanggil `RDBList\GETTanggalClosing_SQL` — `SELECT * FROM POOLDATA.TANGGAL_CLOSING` — di dalam
> aktivitas yang sama dengan pembentukan nomor endorsemen.
>
> Tabel itu **dipakai 13 modul**: `Claim Fac In` · `Claim Life` · `Claim Prop` · `Claim Non Prop` ·
> `Komite Claim` (empat modul) · `NB FacIn` · `NB Treaty In` · `PremiumList Life` ·
> `RNW Fac In` · `Endorsment Fac In` · `EDM Treaty In`. Ia **sumber tunggal tanggal tutup buku
> seluruh sistem**.
>
> ---
>
> **Yang mengikat sistem baru:**
>
> | | |
> | --- | --- |
> | Batas periode produksi | **dibaca dari `POOLDATA.TANGGAL_CLOSING`** setiap kali dibutuhkan |
> | Nilai tanggal di dalam kode | **dilarang** — termasuk sebagai nilai bawaan |
> | Perubahan tanggal | berlaku serentak untuk seluruh modul, tanpa rilis perangkat lunak |
>
> **Yang mengikat pengujian:** test yang menemukan batas periode berasal dari nilai di dalam kode,
> bukan dari tabel, **gagal**. Dan test yang mengubah isi tabel harus melihat perilaku penomoran
> ikut berubah **tanpa perubahan kode**.
>
> `[terbuka]` **Bentuk isi tabel belum terlihat** — apakah satu baris berlaku global, atau satu
> baris per periode, atau per lini bisnis. Itu menentukan bagaimana pembacaannya ditulis.
> Cukup dijawab dengan `SELECT * FROM POOLDATA.TANGGAL_CLOSING` beserta beberapa baris contohnya.
> **Tidak menahan** penulisan spec; menahan penulisan kode pembacanya.
>
> `[terbuka]` **Apakah aturan geser-bulan berlaku untuk semua jenis endorsemen** belum dipastikan.
> Pertanyaan asal butir ini menyebut tiga bagian; dua sudah terjawab oleh keputusan di atas.
> Bagian ketiga tetap terbuka, **tidak menahan**.

---

## P55 — Adendum dan premi tambahan: kapan masing-masing dipakai?

Kami menemukan bahwa **adendum** dan **premi tambahan** adalah **dua hal berbeda**, bukan dua nama
untuk satu hal. Adendum memuat **medan kontrak lengkap** — nama bisnis, ceding, mata uang,
tanggal. Premi tambahan **hanya memuat angka uang**, dan rinciannya dipecah **per lapisan**.

**Konteks:** kami dapat melihat bedanya dari isinya, ⛔ tetapi tidak **kapan** masing-masing
dipakai, dan apakah satu endorsemen dapat memuat keduanya sekaligus.

**Bentuk jawaban yang diharapkan:** **(1)** perubahan seperti apa yang menjadi **adendum** ·
**(2)** perubahan seperti apa yang menjadi **premi tambahan** · **(3)** dapatkah keduanya terjadi
pada satu endorsemen · **(4)** apakah keduanya menghasilkan **nomor tersendiri**, atau berbagi satu
nomor.

**Dampak bila salah:** ⛔ pengguna dipaksa memakai jenis endorsemen yang salah, dan penomoran
serta pengakuan preminya ikut salah.

rujukan: `Section\DetailPolicyTreatyInAddendum` **34** properti lawan
`Section\DetailPolicyTreatyInAddPremi` **14**; `Section\DetailPolicyAddPremiDetail` berkelas
`…Data-XOLRealisasiData` memuat lapisan; `RDBList\GenerateNoEDMTreaty.xml` — ronde 1 §B

**Jawaban:**

> **DITUTUP. Penomoran satu jalur; pembedaan adendum dan premi tambahan diikuti apa adanya.**
> `[keputusan work owner]` 2026-09-22
>
> **Bagian (4) — penomoran — terjawab dari korpus.** `[terverifikasi]`
>
> Nomor endorsemen ditetapkan **satu kali di awal**, di `Activity\SetEDMTNoPolis.xml`
> (kelas `Data-Portal`, dipanggil `Activity\CreateEDMT.xml`):
>
> ```
> PolicyTreatyIn.EDMNo = PolicyTreatyIn.OldData.PolicyNo + "/E"
>                      + @If(PolicyTreatyIn.ProdKe < 10, "0" + ProdKe, ProdKe)
> ```
>
> Contoh: polis `UJI-POL-0001` menghasilkan endorsemen
> `UJI-POL-0001/E01`. *(contoh disamarkan 06-10-2026 — semula nomor polis berformat nyata; aturan prompt eksekusi: nol nomor polis di dokumen.)*
>
> Sumbernya **nomor polis induk** (`OldData.PolicyNo`), bukan deret tersendiri. Adendum dan premi
> tambahan karena itu **berbagi satu nomor** — tidak ada penomoran terpisah untuk masing-masing.
>
> ---
>
> ⛔ **RALAT atas jawaban sementara sesi asisten.** `[penyimpangan sadar]`
>
> Jawaban sementara sebelumnya menyebut bahwa nomor endorsemen dibentuk
> `RDBList\GenerateNoEDMTreaty.xml`:
>
> > *"`SELECT 'RNM-E' || <BusinessOldId> || '.' || <periode> || '.' ||
> > LPAD(TO_CHAR(POOLDATA.ENDORSEMENT_SEQ.NEXTVAL), 5, '0')`"*
>
> **Keliru.** Rule itu memang ada, tetapi **tidak pernah tercapai**. Ia dipanggil
> `Activity\GeneratePolicyNoTreatyAddendum_Act.xml`, yang memuat **tiga langkah bersyarat
> `pyWorkPage.PolicyTreatyIn.EDMNo==""`**. `[terverifikasi]` Karena `SetEDMTNoPolis` sudah mengisi
> `EDMNo` sejak berkas dibuat, syarat itu tidak pernah benar.
>
> **`RNM-E…` dan `POOLDATA.ENDORSEMENT_SEQ` tidak dimigrasi.**
>
> Dua langkah lain pada aktivitas yang sama juga mati, ber-`pyStepsBlockName` = `//`:
> langkah 3 *(copy quotation)* dan langkah 4 *(Fetch Current Date)*. Matinya langkah 4 menguatkan
> **P54**: tanggal batas periode diambil dari `POOLDATA.TANGGAL_CLOSING` lewat langkah 5, bukan dari
> tanggal hari ini.
>
> ---
>
> **Bagian (1), (2), (3) — pembedaan adendum dan premi tambahan — diikuti apa adanya.**
>
> Keduanya tetap dua hal berbeda sebagaimana di sistem lama: adendum memuat **medan kontrak
> lengkap** (34 properti di `Section\DetailPolicyTreatyInAddendum`), premi tambahan hanya memuat
> **angka uang** dan dirinci per lapisan (14 properti di `Section\DetailPolicyTreatyInAddPremi`,
> rincian di `Section\DetailPolicyAddPremiDetail` berkelas `…Data-XOLRealisasiData`).
>
> Sebaran pemakaian: `Addendum` disebut di 16 berkas, `AddPremi` di 4 — premi tambahan tampak
> sebagai **variasi di dalam** alur adendum, bukan alur tersendiri. `[dugaan]`
>
> `[terbuka]` **Aturan bisnis kapan masing-masing dipakai belum tertulis di mana pun**, dan tidak
> dapat disimpulkan dari korpus. **Tidak menahan** penulisan spec — keduanya dibangun sebagaimana
> adanya. Menahan hanya bila kelak diperlukan validasi yang menolak pilihan jenis yang salah.
> Pemiliknya **Product & Underwriting**.

---

## P56 — Sebuah endorsemen dapat DIBATALKAN. Apa yang terjadi pada selisihnya?

Di dalam aturan penghitung selisih terdapat cabang yang ditandai **"endorsemen dibatalkan"**.
⛔ Kami tidak dapat membaca apa yang dikerjakan cabang itu.

**Konteks:** endorsemen mengubah nilai kontrak. Bila ia dibatalkan, nilai kontrak harus kembali —
tetapi kami tidak tahu apakah ia **kembali ke nilai lama**, **dinolkan**, atau **dibiarkan** dengan
penanda batal.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** nilai kontrak kembali ke keadaan
sebelum endorsemen · **(b)** endorsemen tetap tercatat tetapi selisihnya dinolkan · **(c)** cara
lain, jelaskan.
⭐ Dan: **siapa yang boleh membatalkan**, dan **sampai tahap apa** pembatalan masih mungkin.

**Dampak bila salah:** ⛔ nilai kontrak **tidak kembali** sesudah pembatalan, dan selisihnya ikut
terbawa ke laporan.

rujukan: `Activity\CalculateDifferenceEDM_act.xml`, catatan langkah *"EDM BATAL"*; ditambah
`Activity\SetEDMTCancel.xml` di antara 31 berkas EDM inti — ronde 1 §A.4

**Jawaban:**

> **Mekanismenya terjawab dari korpus: data BARU dinolkan, data LAMA tetap, endorsemen tetap
> tercatat.** `[terverifikasi]` `[keputusan work owner]` 2026-09-22
>
> Bukan **(a)** *(kembali ke keadaan sebelum endorsemen)* dan bukan **(b)** *(selisih dinolkan)*,
> melainkan bentuk ketiga:
>
> ```
> selisih = nilai baru - nilai lama
> saat batal:  nilai baru := 0
> maka         selisih   := -(nilai lama)
> ```
>
> Endorsemen **tidak dihapus**. Ia tetap tercatat dengan nomornya sendiri, membawa selisih yang
> **membalikkan seluruh nilai kontrak**.
>
> **Bukti.** `[terverifikasi]` `Activity\SetEDMTCancel.xml` — kelas `ASM-FW-GISFW-Work`,
> enam langkah, seluruhnya `Property-Set`:
>
> | # | Catatan langkah |
> | ---: | --- |
> | 1 | *"Set Property to zero"* |
> | 2 | *"set installment &"* |
> | 5 | *"Set Total Properties"* |
> | 6 | *"Set Per Layer Properties"* |
>
> Dan catatan di `Activity\EDMChooseBusiness_Act.xml` menyatakannya terang-terangan:
> **"When EDM cancel, set new value to 0"**.
>
> Penolkan berlaku sampai **tingkat lapisan** — langkah 6.
>
> Dipanggil dari dua tempat: `Activity\CreateEDMT.xml` dan `Activity\EDMChooseBusiness_Act.xml`.
>
> ---
>
> `[terbuka]` **Siapa yang boleh membatalkan, dan sampai tahap apa.** Tidak dapat dibaca dari
> korpus. `[dugaan]` Kedua pemanggilnya berada di tahap **awal** — saat berkas dibuat dan saat
> bisnis dipilih — sehingga pembatalan mungkin hanya dapat dilakukan **sebelum endorsemen
> diajukan**. Belum terbukti: pemanggilan dari layar tidak terekspor. Pemiliknya
> **Product & Underwriting**.
>
> `[terbuka]` **Apakah endorsemen yang dibatalkan tetap masuk laporan.** Karena selisihnya menjadi
> negatif penuh dan berkasnya tetap tercatat, ia akan muncul pada pengakuan premi periode itu
> **sebagai pengurang**. Bila yang dikehendaki adalah *"seolah tidak pernah ada"*, perlakuannya
> harus berbeda — dan itu keputusan **Finance**, bukan keputusan teknis.
>
> **Cukup untuk menulis spec.** Kedua butir terbuka di atas **tidak menahan**: mekanismenya sudah
> pasti, yang belum pasti hanya wewenang dan perlakuan pelaporannya.

---

## P58 — Data lama disalin ke berkas endorsemen, atau dibaca dari polis induknya?

Layar endorsemen menampilkan **nilai lama** berdampingan dengan **nilai baru**. Nilai lama itu
disimpan dengan penanda tersendiri, yang `[dugaan]` berarti ia **disalin** ke berkas endorsemen
saat endorsemen dibuat — ⛔ tetapi kami belum dapat memastikannya.

**Konteks:** bedanya besar. Bila **disalin**, nilai lama **beku** pada saat endorsemen dibuat.
Bila **dibaca**, ia ikut berubah bila polis induknya berubah — dan selisihnya berubah bersamanya.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** disalin dan beku · **(b)** dibaca dari
polis induk · **(c)** tidak tahu.
⭐ Dan pertanyaan pendampingnya: **apa yang terjadi bila polis induk berubah sesudah endorsemen
dibuat tetapi belum disetujui?**

**Dampak bila salah:** ⛔ selisih yang ditampilkan **berbeda** dari selisih yang tersimpan, dan
keduanya tidak dapat direkonsiliasi.

rujukan: layar `…PropOldData` memakai penanda halaman tersendiri untuk nilai lama; layar
`…PropValueDifference` memakai penanda ketiga — ronde 1 §A.1 dan §A.4

**Jawaban:**

> **(a) Disalin dan beku.** `[keputusan work owner]` `[terverifikasi]` 2026-09-22
>
> Nilai lama **disalin** ke berkas endorsemen saat berkas dibuat, lalu **tidak pernah dibaca ulang**
> dari polis induk. Perubahan pada polis induk sesudah endorsemen dibuat **tidak berpengaruh**, dan
> selisihnya tetap.
>
> **Bukti — `Activity\CreateEDMT.xml`, 20 langkah.** `[terverifikasi]` Lima yang menentukan:
>
> | # | Langkah | Catatan |
> | ---: | --- | --- |
> | 7 | `RDB-List` | *"Call data from json_polis"* — baca polis induk dari basis data |
> | 8 | `Java` | *"Map oracle column to clipboard"* |
> | 10 | `Page-Copy` | *"Copy from TempPage.PolicyTreatyIn to MergePage"* |
> | 11 | `Page-Copy` | *"Copy from MergePage.PolicyTreatyIn.OldData.Quo…"* |
> | 13 | `Page-Remove` | *"Remove TempPage"* |
>
> Halaman sumbernya **dibuang di langkah 13**. Sesudah itu tidak ada lagi sambungan ke polis induk.
>
> Perilaku ini benar untuk dokumen endorsemen: ia mencatat perubahan terhadap **keadaan pada saat
> itu**, bukan terhadap keadaan sekarang.
>
> **Pertanyaan pendamping terjawab sekaligus:** bila polis induk berubah sesudah endorsemen dibuat
> tetapi belum disetujui, **tidak terjadi apa-apa** pada endorsemennya.
>
> ---
>
> **Dua hal lain yang ikut terbukti dari aktivitas yang sama.**
>
> `[terverifikasi]` **Pembatalan adalah jenis endorsemen, bukan tindakan terpisah.** Langkah 16
> memanggil `SetEDMTCancel` dengan catatan *"when Param.edmtype = 4, run this activity to z[ero]"*.
> Melengkapi **P56**: pembatalan dipilih di awal, saat berkas dibuat.
>
> `[terbuka]` **Tiga pemeriksaan pencegah dimatikan** — ber-`pyStepsBlockName` = `//`:
>
> | # | Pemeriksaan | Catatan |
> | ---: | --- | --- |
> | 4 | `CheckNopolisAvailability` | *"Check Nopolis Availability"* |
> | 5 | `CheckDupeAdjPremi` | *"Check Duplicate Adjustment Premi"* |
> | 6 | `CheckOngoingTreatyEDM` | *"Check Ongoing for same nopolis"* |
>
> Yang ketiga paling berat akibatnya. Bila dua endorsemen dapat berjalan **bersamaan** atas polis
> yang sama, keduanya menyalin nilai lama yang sama pada langkah 7-11, lalu masing-masing menyimpan
> selisih terhadap dasar yang sama. Hasilnya dapat salah **tanpa ada yang menyadarinya** — dan
> keputusan **P57** bahwa endorsemen berlapis menghitung selisih terhadap keadaan tepat sebelumnya
> justru mengandaikan tidak ada dua yang berjalan serentak.
>
> Sifatnya sama dengan **P52**: perlu dipastikan apakah dimatikan sengaja atau tertinggal.
> **Belum diangkat menjadi pertanyaan bernomor** — menunggu keputusan work owner apakah perlu.
>
> ---
>
> ⭐⭐ **DIPERIKSA ULANG 2026-09-22 terhadap data produksi — jawaban ini TETAP BERDIRI.**
> `[terverifikasi]`
>
> Contoh `DATA_JSON` sungguhan memperlihatkan **`OldData` ada di dalam dokumen POLIS BARU**, bukan
> hanya di endorsemen. Itu tidak membatalkan jawaban ini: `CreateEDMT` langkah 10 tetap menyalin
> *"from TempPage.PolicyTreatyIn to MergePage.PolicyTreatyIn.OldData"*, dan langkah 13 tetap
> membuang halaman sumbernya. ⭐ **Penyalinan dan pembekuannya tidak berubah.**
>
> ⭐ **Sebab `OldData` sudah ada di polis baru pun terbaca dari korpus:**
> `Activity\TreatyRealizationCheckXOLList` **di NB Treaty In** menetapkan
> `pyWorkPage.PolicyTreatyIn.OldData.TreatyXOLList = pyWorkPage.PolicyTreatyIn.TreatyXOLList`
> saat realisasi polis baru.
>
> ⛔⛔ `[terbuka]` **BARU — `OldData` di dalam `OldData`.** Karena halaman yang disalin langkah 10
> **sudah** memuat `OldData` miliknya sendiri, hasilnya bersarang. ⚠️ Dan korpus **tidak pernah
> membacanya**: rujukan berbentuk `OldData…OldData` = **0**. ⇒ Susunan itu **dibuat tetapi tidak
> dipakai**, menumpuk pada tiap endorsemen berlapis. **Perlu dipastikan apakah sengaja.**
> Pemiliknya `[work owner]`. ⛔ **Tidak ditutup di sini.**

---

# Untuk pengembang Pega lama

## P52 — Dua langkah penanganan galat DIMATIKAN. Sengaja?

Di dalam jalur penanganan kegagalan pengiriman, **dua langkah dimatikan** — ditandai sebagai
nonaktif di dalam aturannya sendiri:

- langkah yang **memasang pesan galat untuk pengguna**;
- langkah yang **menutup paksa berkasnya**.

⛔ Akibatnya, ketika pengiriman gagal, **pengguna tidak melihat apa pun di layar**. Yang tersisa
hanyalah catatan di tabel pemantauan dan **sebuah surel berlampiran**.

**Konteks:** kami perlu tahu apakah pengguna memang **tidak seharusnya** diberi tahu, atau kedua
langkah itu dimatikan sementara lalu terlupakan.

**Bentuk jawaban yang diharapkan:** untuk masing-masing — **(a)** memang sengaja dimatikan,
jelaskan alasannya · **(b)** dimatikan sementara, seharusnya kembali · **(c)** tidak tahu.
⭐ Dan: **siapa penerima surel** itu sekarang.

**Dampak bila salah:** ⛔ pengguna mengira endorsemennya berhasil padahal tidak, dan baru tahu
ketika ada yang membaca surel atau memeriksa tabel pemantauan.

rujukan: `Activity\serviceInsertArasapas_act.xml` — langkah **10** *(pemasang pesan)* dan **18**
*(penutup paksa)*, keduanya ber-`pyStepsBlockName` berawalan `//`; langkah 20
*"SEND EMAIL JIKA ERROR KONVERSI"* — ronde 1 §C.2

**Jawaban:**

> **Langkah 10 dihidupkan. Langkah 18 tetap mati.** `[keputusan work owner]` 2026-09-22
>
> | Langkah | Sistem lama | Sistem baru |
> | ---: | --- | --- |
> | **10** — pasang pesan galat untuk pengguna | dimatikan | ⭐ **dibangun dan aktif** |
> | **18** — tutup paksa berkas | dimatikan | **tetap tidak dibangun** |
>
> **Bukti cara keduanya dimatikan.** `[terverifikasi]` Penandanya `pyStepsBlockName` bernilai `//`
> — dikomentari, bukan bergerbang. Diperiksa langsung di
> `Activity\serviceInsertArasapas_act.xml`: langkah 10 dan 18 bernilai `//`; langkah 8 bernilai
> `JMP`; sisanya kosong.
>
> ---
>
> **Langkah 10 — (b), seharusnya kembali.** Sudah ditetapkan di **P51**: pesan galat wajib
> ditampilkan kepada pengguna **sebelum** jalur penghapusan dijalankan. Butir ini hanya menegaskan
> bahwa sistem lama memang mematikannya, dan sistem baru tidak menirunya.
> `[penyimpangan sadar]`
>
> **Langkah 18 — (a), memang tidak dipakai, dan tidak dihidupkan.**
>
> Alasannya: menutup paksa berkas sesudah kegagalan berarti pengguna **kehilangan kesempatan
> mencoba lagi**, padahal data produksinya sudah dihapus. Tidak ada yang tersisa, dan tidak ada
> yang dapat diperbaiki.
>
> Membiarkan berkas **tetap terbuka** justru benar: konversi gagal, data produksi bersih kembali,
> berkas tetap di tangan pengguna untuk dicoba ulang sesudah penyebabnya diatasi. Penyebab tersering
> adalah kegagalan jaringan — hal yang memang layak dicoba lagi.
>
> **Yang mengikat pengujian:** test yang menemukan berkas **tertutup otomatis** sesudah kegagalan
> konversi **gagal**. Berkas tetap terbuka dan dapat dikirim ulang oleh pengguna yang sama.
>
> ---
>
> **Keadaan akhir jalur gagal-konversi di sistem baru:**
>
> ```
> 1. konversi gagal
> 2. pesan galat DITAMPILKAN kepada pengguna        <- langkah 10, dihidupkan
> 3. jalur penghapusan data produksi dijalankan     <- langkah 12
> 4. dicatat di tabel pemantauan, status gagal      <- langkah 15, 16
> 5. surel pemberitahuan dikirim                    <- langkah 20
> 6. berkas TETAP TERBUKA, dapat dicoba ulang       <- langkah 18 tidak dibangun
> ```
>
> `[terbuka]` **Siapa penerima surel langkah 20 sekarang?** Satu-satunya bagian butir ini yang tidak
> dapat dijawab dari korpus. Perannya berubah sesudah keputusan di atas: dulu surel itu
> **satu-satunya** yang mengetahui kegagalan; kini pesan layar yang utama dan surel menjadi
> cadangan. **Tidak menahan.** Tetapi bila penerimanya ternyata satu orang yang sudah berpindah
> bagian, itu perlu diketahui sebelum go-live.

---

## P57 — Ada penanda selisih DI DALAM data lama. Berarti endorsemen berlapis?

Layar endorsemen punya dua pasang: satu pasang menampilkan nilai lama dan baru, ⭐ pasangan kedua
menampilkan **selisih yang tersimpan di dalam data lama** — penanda bersarang dua tingkat.

**Konteks:** `[dugaan]` itu berarti selisih dari **endorsemen sebelumnya** ikut dibawa, sehingga
sebuah polis dapat diendorse **berkali-kali** dan tiap endorsemen menyimpan selisihnya sendiri.
⛔ Kami belum dapat memastikannya.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** benar, satu polis dapat diendorse
berkali-kali dan tiap kali menyimpan selisihnya · **(b)** hanya satu endorsemen per polis; pasangan
kedua untuk hal lain, jelaskan · **(c)** tidak tahu.
⭐ Bila **(a)**: **selisih yang mana yang dipakai laporan** — yang terakhir, atau jumlah seluruhnya?

**Dampak bila salah:** ⛔ endorsemen kedua **menimpa** yang pertama, atau sebaliknya selisih
**dijumlahkan dua kali** dalam laporan.

rujukan: `Section\DetailPolicyTreatyInPropOldData2.xml` memakai penanda bersarang dua tingkat;
irisan dengan `…PropOldData` **15** properti, hanya di yang ber-akhiran 2 **27** — ronde 1 §A.3

**Jawaban:**

> **(a) Benar — satu polis dapat diendorse berkali-kali, dan tiap endorsemen menyimpan selisihnya
> sendiri.** `[keputusan work owner]` `[terverifikasi]` 2026-09-22
>
> Selisih dihitung terhadap **keadaan tepat sebelumnya**, bukan terhadap polis asli.
>
> **Bukti pertama — penomoran menghitung naik.** `Activity\SetEDMTNoPolis.xml`:
>
> ```
> MergePage.PolicyTreatyIn.ProdKe = @toInt(OutputData.pxResults(1).CARI1) + 1
> ```
>
> Langkah 2 membaca nilai dari basis data lalu **menambah satu**. Itulah yang menghasilkan
> `/E01`, `/E02`, `/E03`. Bila endorsemen hanya boleh sekali, penghitungan itu tidak diperlukan.
>
> **Bukti kedua — layar memilih menurut apakah data lama sendiri sudah pernah diendorse.**
> `[terverifikasi]`
>
> | Berkas | Syarat tampil | Keadaan |
> | --- | --- | --- |
> | `Section\DetailPolicyTreatyInPropOldData` | `.OldData.EDMNo = ''` | data lama **belum pernah** diendorse |
> | `Section\DetailPolicyTreatyInPropOldData2` | `.OldData.EDMNo != ''` | data lama **sudah** berupa endorsemen |
>
> `Activity\EDMTCalculateTreatyDifference.xml` memakai syarat yang sama untuk memilih cara
> menghitung selisih.
>
> ⭐ **Arti akhiran `2` terjawab pasti**: bukan dua tingkat tampilan, melainkan **dua keadaan** —
> layar biasa untuk endorsemen pertama, layar ber-akhiran 2 untuk endorsemen atas endorsemen.
> Ini mengoreksi dugaan ronde 1 yang menyebutnya *"selisih dari endorsemen sebelumnya"* sebagai
> `[dugaan]`; kini `[terverifikasi]`.
>
> ---
>
> `[terbuka]` **Selisih mana yang dipakai laporan — yang terakhir, atau jumlah seluruhnya?**
> Tidak dapat dijawab dari korpus, dan **bedanya besar**: bila sebuah polis diendorse tiga kali,
> laporan dapat menampilkan selisih ketiga saja atau jumlah ketiganya. Salah pilih berarti premi
> diakui dua kali, atau satu putaran hilang. Pemiliknya **Finance**, bukan Product — ini soal
> pengakuan premi.
>
> `[terbuka]` **Adakah batas berapa kali satu polis boleh diendorse?** `ProdKe` dibentuk dua digit
> (`@If(ProdKe < 10, "0" + ProdKe, ProdKe)`), sehingga batas teknisnya **99**. Bila ada batas bisnis
> yang lebih rendah, itu validasi yang perlu dibangun; bila tidak ada, angka 99 tetap batas yang
> perlu diketahui sebelum sistem baru memakainya. **Tidak menahan.**
>
> ⛔ **KOREKSI 06-10-2026.** (XML) Rumus itu hanya membubuhkan nol di bawah sepuluh; ≥ 100 menjadi tiga digit —
> **tidak ada batas 99**. (Keputusan) ✅ ditutup `[keputusan work owner]` 23-09: **tanpa batas**
> (`spec-penyimpanan-relasional.md` bab *KEPUTUSAN 23-09-2026* butir 1).
>
> ---
>
> ⭐⭐ **DIPERIKSA ULANG 2026-09-22 terhadap data produksi — jawaban ini TETAP BERDIRI.**
> `[terverifikasi]`
>
> Contoh `DATA_JSON` memperlihatkan `OldData` **ada di dokumen polis baru**. Pemilih layar di atas
> menguji **isi `.OldData.EDMNo`** *(kosong lawan terisi)*, **bukan keberadaan `OldData`** — maka
> ia tetap sahih.
>
> ⭐ **Temuan itu justru menjelaskan kenapa ujinya ditulis begitu:** bila `OldData` ada bahkan pada
> polis baru, uji keberadaan akan **selalu benar** dan tidak berguna. Pengembang lama memilih medan
> yang membedakan, bukan simpulnya.
>
> ⚠️ **Satu butir `[terbuka]` baru yang menyentuh jawaban ini** — `OldData` bersarang di dalam
> `OldData` pada endorsemen berlapis; lihat catatan di **P58**.

---

## P60 — Satu langkah tambahan pada penghitung penyebaran. Apa yang dikerjakannya?  ⭐ **DITUTUP DARI KORPUS**

Aturan penghitung penyebaran di modul endorsemen punya **satu langkah lebih banyak** daripada
versi di modul polis baru. Catatan langkahnya menyebut keadaan yang ditanganinya: **daftar
penyebaran berisi satu baris, dan persentase pembagiannya kosong atau nol**.

⛔ Yang **dikerjakan** langkah itu tidak dapat kami baca — nilainya tidak ikut terkirim, sama
seperti seluruh langkah penetapan nilai lainnya.

**Konteks:** keadaan "satu baris, persentase kosong" tampak seperti **kasus khusus** yang
ditambahkan belakangan. Kami perlu tahu apa yang seharusnya terjadi.

**Bentuk jawaban yang diharapkan:** satu paragraf — **apa yang ditetapkan** ketika keadaan itu
terjadi, dan **kenapa** endorsemen memerlukannya sementara polis baru tidak.

**Dampak bila salah:** ⛔ penyebaran bernilai salah pada endorsemen yang hanya punya satu baris —
dan itu justru bentuk endorsemen yang paling umum.

rujukan: `Activity\CountSpreading_Act.xml` — **8 langkah** di EDM lawan **7** di modul polis baru;
catatan langkah *"JIKA SPREADINGLIST 1 DAN SplitRNMSharePct NULL || 0"*; ⛔ nilai `Property-Set`
tidak terekspor *(sama seperti **P18**)* — ronde 1 §F.1

> ⛔⛔ **RALAT atas baris `rujukan:` di atas.** `[penyimpangan sadar]` 2026-09-22
> Anak kalimat ⛔ *"nilai `Property-Set` tidak terekspor (sama seperti **P18**)"* **salah**.
> Nilainya **terekspor**, di tag `PropertiesName`/`PropertiesValue` — **tanpa awalan `py`**.
> **P18 sendiri sudah ditarik.** Lihat `modul/nbtreatyin/docs/VERIFIKASI-P18.md`.
> *(Koreksi 06-10: jalur semula `..\nb-treaty-in\…` terpotong menjadi dua baris.)*

**Jawaban:**

> ⭐⭐ **DIJAWAB DARI KORPUS, bukan ditanyakan.** `[terverifikasi]` 2026-09-22
>
> ⚠️ **Tidak perlu jawaban dari `[pengembang Pega lama]`.** Butir ini ditutup dari bahan yang
> sudah ada, sesudah isi langkah terbukti terbaca.
>
> ### a. Apa yang dikerjakan langkah tambahan itu
>
> `Activity\CountSpreading_Act.xml` di EDM punya **8 langkah bermetode**, di NB **7** — cocok
> dengan rujukan. Langkah tambahannya adalah **langkah 4**, berketerangan
> *"JIKA SPREADINGLIST 1 DAN SplitRNMSharePct NULL || 0"*, dan isinya **satu penetapan tunggal**:
>
> ```
> .SpreadingRiskList(1).SplitRNMSharePct = 100
> ```
>
> ### b. Kenapa nilai itu diperlukan — rumus di bawahnya berbeda
>
> Langkah berikutnya menghitung `SharePercentage` bila ia kosong, dan **kedua modul memakai rumus
> yang berlainan**:
>
> | | Rumus bila `SharePercentage` kosong | Presisi |
> | --- | --- | ---: |
> | **NB** *(langkah 4.1)* | `@if(.SharePercentage=="", 100/@toDecimal(@LengthOfPageList(Primary.SpreadingRiskList)), .SharePercentage)` | **10** |
> | **EDM** *(langkah 5.1)* | `@if(.SharePercentage=="", @divide(.SplitRNMSharePct, Primary.RNMShare, 20), .SharePercentage)` | **20** |
>
> ⭐ **Di situlah sebabnya.** Rumus EDM **membagi dengan `SplitRNMSharePct`**; rumus NB tidak
> menyentuh medan itu sama sekali — ia memakai **cacah baris**, yang tidak pernah kosong selama
> daftarnya ada. Maka hanya EDM yang dapat kejatuhan medan kosong, dan hanya EDM yang butuh
> nilai bawaan.
>
> ⭐ Presisi juga berbeda: **10** di NB lawan **20** di EDM, termasuk pada `PremiumSpreaded` dan
> `ClaimSpreaded`. ⚠️ Ini **ketidakseragaman presisi yang sama** dengan yang sudah tercatat di
> butir `[terbuka]` P46 — disebut di sini supaya tidak dicari ulang, **tidak ditutup di sini**.
>
> ### c. Kenapa endorsemen memerlukannya sementara polis baru tidak
>
> ⭐⭐ **Karena aturan yang MENGISI medan itu hanya ada di NB.** `[terverifikasi]`
>
> | | `TreatyNonPropSetSpreading` | Yang mengisi `SplitRNMSharePct` |
> | --- | --- | --- |
> | **NB Treaty In** | ⭐ **ADA** | rumusnya sendiri: `SpreadingRiskList(<LAST>).SplitRNMSharePct = .Pct`, dan untuk kasus satu baris `SharePercentage = "100"` langsung |
> | **EDM Treaty In** | ⛔ **TIDAK ADA di modul** | ⛔ **tidak ada satu pun** — satu-satunya berkas EDM yang menyebut medan itu adalah `CountSpreading_Act` sendiri |
>
> Di **polis baru**, daftar penyebaran **dibangun dari nol** oleh `TreatyNonPropSetSpreading`, yang
> mengisi `SplitRNMSharePct` **dan** `SharePercentage` pada saat pembuatan. Ketika
> `CountSpreading_Act` berjalan, medannya sudah terisi.
>
> Di **endorsemen**, daftar itu **datang dari polis yang sedang diendorse**, bukan dibangun ulang —
> dan `SplitRNMSharePct` dapat tiba **kosong atau nol**. Karena rumus EDM membaginya, langkah
> tambahan itu memberinya nilai bawaan **100** untuk kasus satu baris, yang membuat pembagiannya
> menjadi `100 / RNMShare`.
>
> ⭐ **Terjemahan dagangnya:** bila endorsemen hanya punya satu baris penyebaran, seluruh bagian
> dianggap **milik RNM sepenuhnya**, dan persentasenya dihitung terhadap bagian RNM.
>
> ### d. Satu perbedaan struktural lagi, dicatat apa adanya
>
> `[terverifikasi]` Langkah 5.1 EDM membawa `pyStepsPreCondition = "false"`; langkah 4.1 NB tidak
> membawa tag itu sama sekali *(self-closing)*. Menurut aturan baca yang sudah mengikat sejak
> ronde 3, `pyStepsPreCondition = "false"` **mematikan GERBANGNYA, bukan langkahnya** — sehingga
> **kedua langkah sama-sama berjalan**. Dicatat supaya tidak salah dibaca kemudian.
>
> ### ⚠️ Yang TETAP `[terbuka]`, dan tidak ditutup di sini
>
> **Kenapa RNM memilih dasar `SplitRNMSharePct / RNMShare` untuk endorsemen** sementara polis baru
> memakai pembagian rata menurut cacah baris — itu **pertanyaan maksud dagang**, bukan pertanyaan
> bahan, dan **tidak terbaca dari korpus**. Pemiliknya `[Product+Underwriting]`.
> ⛔ **Tidak dikarang di sini.**

---

# Untuk IAM

## P53 — Pengiriman ke sistem produksi berjalan TANPA autentikasi. Benar?

Panggilan dari sistem lama ke sistem produksi disetel **tanpa autentikasi sama sekali** — tidak
ada kunci, tidak ada kata sandi, tidak ada token. Batas waktunya **tiga puluh detik**.

**Konteks:** yang dikirim adalah **nomor polis, pengenal berkas, dan waktu pembuatan** — bukan data
nasabah, tetapi cukup untuk menyisipkan data ke sistem produksi.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** benar, jalur ini tertutup di jaringan
internal sehingga autentikasi tidak diperlukan · **(b)** seharusnya ada autentikasi, dan ini
kelemahan yang diketahui · **(c)** ada autentikasi di lapisan lain, jelaskan di mana.
⭐ Dan: **apakah sistem baru boleh menambahkan autentikasi**, atau sistem penerimanya belum siap.

**Dampak bila salah:** ⛔ siapa pun yang dapat menjangkau alamat itu **dapat menyisipkan data** ke
sistem produksi; atau sebaliknya sistem baru menambahkan autentikasi yang **ditolak** penerimanya.

rujukan: `ConnectREST\convertJsonNusareToProduction.xml` — `pyUseAuthentication` = `false`,
`pyResponseTimeout` = `30000`; dikirim `PolicyNo` · `pzInsKey` · `pxCreateDateTime` — ronde 1 §C.3

**Jawaban:**

> **(a) Benar — jalur ini tertutup di jaringan internal, autentikasi tidak diperlukan.**
> `[keputusan work owner]` 2026-09-22
>
> Dan satu keterangan yang menentukan arah: **autentikasi, bila kelak ada, datang dari sistem
> penerima** — sebab kitalah yang memanggil mereka, bukan sebaliknya. Sistem baru **tidak
> menambahkan autentikasi atas kehendak sendiri**.
>
> **Keadaan sekarang, terverifikasi.** `[terverifikasi]`
> `ConnectREST\convertJsonNusareToProduction.xml` — `pyUseAuthentication` = `false`,
> `pyResponseTimeout` = `30000`. Yang dikirim: `PolicyNo`, `pzInsKey`, `pxCreateDateTime` —
> bukan data nasabah.
>
> ---
>
> **Yang mengikat rancangan, meskipun jawabannya "tidak perlu".**
>
> Autentikasi dibangun sebagai **setelan yang sengaja dimatikan**, bukan sebagai fitur yang tidak
> ada. Satu baris konfigurasi bernilai mati — bukan kode yang tidak ditulis.
>
> Sebabnya praktis: bila kelak sistem penerima mewajibkan token, menghidupkannya menjadi perubahan
> **konfigurasi**. Bila mekanismenya tidak dibangun sama sekali, itu menjadi perubahan **kode**.
>
> Dan perubahan kode itu tidak berhenti di satu tempat: mekanisme alamat yang sama —
> `GetLinkService` pada kelas `ASM-FW-GISFW-Int-M_LINK_SERVICE` — dipakai **18 modul**, untuk
> Google Storage, layanan Kasir, Gemini AI, konversi klaim, dan konversi produksi ini.
> `[terverifikasi]` Keputusan di sini karena itu menjadi **preseden bagi seluruh integrasi keluar**,
> bukan bagi satu panggilan.
>
> Sejalan dengan **ADR-0013** *(resolusi endpoint via LinkService)*. Tidak ada ADR baru.
>
> `[terbuka]` **Batas waktu 30 detik belum ditinjau.** Itu lama untuk panggilan yang menahan
> pengguna di layar. Bila penerima biasanya menjawab dalam hitungan detik, batas selama itu hanya
> memperpanjang waktu tunggu saat ada gangguan — tanpa menambah peluang berhasil. **Tidak menahan**;
> layak ditetapkan bersama pemilik sistem penerima ketika jalur ini diuji.
>
> `[terbuka]` **Alamat dibaca dari tabel `M_LINK_SERVICE` saat berjalan**, bukan dari berkas
> konfigurasi. Artinya siapa pun yang dapat menulis ke tabel itu dapat mengalihkan panggilan.
> Karena jaringannya tertutup, ini **tidak mendesak** — tetapi wewenang tulis ke tabel itu layak
> diperiksa bersama IAM ketika modul lain yang memakai mekanisme sama digarap.

---

# Untuk DBA

## P59 — Alamat sistem produksi dibaca dari sebuah tabel saat berjalan. Isinya apa?

Alamat sistem produksi **tidak ditulis di dalam aturan**. Ia dibaca dari **sebuah tabel layanan**
setiap kali dibutuhkan, lewat sebuah ungkapan yang menunjuk properti, bukan nilai tetap.

**Konteks:** kami perlu tahu isi tabel itu sebelum menulis lapisan integrasi — berapa baris, apa
yang membedakan satu baris dari yang lain, dan bagaimana baris yang tepat dipilih.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh isi tabel** — seluruh barisnya, atau sekurangnya
kolom pengenal dan kolom alamatnya. ⭐ Ditambah satu keterangan: apakah alamatnya **berbeda antar
lingkungan** *(uji dan produksi)*, dan bagaimana pembedaannya.

**Dampak bila salah:** ⛔ sistem baru memanggil **alamat yang salah** — dan bila itu alamat
produksi dari lingkungan uji, data uji masuk ke sistem produksi.

rujukan: `SystemSettings\LinkService.xml` — `pySetting` berupa ungkapan, bukan nilai tetap;
`Activity\GetLinkService.xml` berkelas `ASM-FW-GISFW-Int-M_LINK_SERVICE`, **4 langkah**, memakai
`Obj-Browse`; sejalan **ADR-0013** — ronde 1 §C.3

**Jawaban:**

> **Lingkungan uji dan produksi memakai basis data berbeda, sehingga isi tabelnya berbeda dengan
> sendirinya.** `[keputusan work owner]` 2026-09-22
>
> Tidak diperlukan penjaga tambahan di lapisan aplikasi untuk mencegah lingkungan uji memanggil
> alamat produksi — pemisahannya sudah terjadi di tingkat basis data.
>
> **Struktur tabel terbaca dari korpus.** `[terverifikasi]` `Activity\GetLinkService.xml`
> berkelas `ASM-FW-GISFW-Int-M_LINK_SERVICE`, empat langkah — `Page-New`, `Obj-Browse`,
> `Property-Set`, `Page-Remove` — merujuk **tiga properti**:
>
> | Kolom | Peran |
> | --- | --- |
> | `KATEGORI_1` | penggolong tingkat satu |
> | `KATEGORI_2` | penggolong tingkat dua |
> | `URL` | alamat layanan |
>
> Pemilihan baris memakai **dua kategori**, bukan satu. Itu menjelaskan bagaimana satu tabel
> melayani Google Storage, layanan Kasir, Gemini AI, konversi klaim, dan konversi produksi
> sekaligus — masing-masing memakai pasangan kategori tersendiri.
>
> Hasilnya masuk ke `ResponLink.URL`, yang lalu dibaca setelan `SystemSettings\LinkService`
> (`pySetting` = `=ResponLink.URL`). Sejalan **ADR-0013**; tidak ada ADR baru.
>
> ---
>
> **Yang mengikat sistem baru:**
>
> | | |
> | --- | --- |
> | Alamat layanan | **dibaca dari `POOLDATA.M_LINK_SERVICE`**, dipilih dengan dua kategori |
> | Alamat di dalam kode atau berkas konfigurasi aplikasi | **dilarang** — termasuk sebagai nilai bawaan |
> | Pemisahan lingkungan | dijamin oleh **basis data yang berbeda**, bukan oleh kode |
>
> **Yang mengikat pengujian:** test yang menemukan alamat layanan berasal dari nilai di dalam kode,
> bukan dari tabel, **gagal**.
>
> `[terbuka]` **Isi tabelnya belum dilihat** — berapa baris, dan pasangan kategori mana yang
> menunjuk konversi produksi. Diperlukan sebelum lapisan integrasi ditulis, bukan sebelum spec.
> Satu perintah cukup:
> `SELECT kategori_1, kategori_2, url FROM POOLDATA.M_LINK_SERVICE;`
> Ditumpangkan ke surat DBA yang sudah ada. **Tidak menahan.**

---

## Sesudah Anda menjawab

⭐ Kembalikan lembar ini apa adanya — **tidak perlu lengkap**. Jawaban sebagian tetap berguna.

⚠️ **Dua di antaranya menahan pekerjaan:** **P51** *(penghapusan data produksi)* dan **P50**
*(pengiriman kedua)* — keduanya menyentuh **data yang sudah masuk sistem produksi**, dan salah
menebaknya tidak dapat dibatalkan.

---

*Disusun 22 September 2026, ronde 1 modul EDM Treaty In, sesudah NB Treaty In selesai sampai tahap
tiket. Tidak ada pertanyaan di lembar ini yang dijawab sendiri oleh tim migrasi.*
