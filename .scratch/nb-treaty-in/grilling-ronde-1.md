# NB Treaty In — Grilling Ronde 1
## Tiga procedure tanpa badan, dan satu halaman parameter yang artinya berpindah per rule

Konteks: **Treaty Inward — Realisasi & Endorsement**. Modul yang digrill: **hanya `NB Treaty In`**.
Korpus READ-ONLY `D:\XML\RNM_BRD\NB Treaty In\`. Aturan: `CLAUDE.md` §4 dan §4a.

⛔ **Ronde ini MENEMUKAN dan BERTANYA.** Nol keputusan rancangan · nol kode · nol DDL ·
`CREATE TABLE` nol · nol daftar kolom usulan · ⛔ **nol butir terbuka yang saya tutup sendiri.**

> **SENSUS BERKAS INI**
>
> ⭐ **Jendelanya: berkas ini MINUS blok sensus ini sendiri.** ⛔ Disebut terang-terangan supaya
> blok ini tidak mengubah angka yang ia klaim.
>
> **842 baris di jendela** · **10 bab** `## ` ·
> ⚠️ *(berkas utuh: **866** menurut pemisah baris, **865** menurut `wc -l` — ⭐ selisih satu,
> sebabnya pergantian baris penutup; blok ini **24** baris. Jendela **842** sama pada kedua cara
> bila alat yang sama dipakai di kedua sisi pengurangan.)*
> **23 sub-bab** `### ` · **17 sub-sub** `#### ` · **21 tabel**
>
> **17 pernyataan berpenanda terverifikasi** · **3 berpenanda dugaan** · **5 berpenanda data DBA**
>
> ⭐ **Register Bab C: 24 butir — dihitung DUA CARA.** *(a)* baris tabel ⇒ **17 terdaftar + 7 baru**;
> *(b)* nomor di kolom pertama ⇒ OQ **3-digit 17 buah**, butir baru bernomor **1–7 tanpa nomor
> hilang**. ✅ **Sepakat.** ⛔ **Nol butir ditutup.**
>
> ⭐ **Bab D: 17 pertanyaan, P1–P17 utuh**, terbagi 4+6+3+4 menurut pemilik. ✅ **Cocok.**
>
> ⚠️ **Dua ungkapan terlarang muncul di jendela ini** — `CREATE TABLE` **2×** dan
> `CREATE OR REPLACE` **1×**. ⭐ Keduanya **disebut, bukan dipakai**: yang pertama di kalimat yang
> menyatakannya nol, yang kedua **di dalam perintah audit** yang membuktikan nol badan procedure
> ada. ⛔ Dinyatakan di sini apa adanya supaya pemeriksa berikutnya tidak menyangka sensusnya
> berbohong.
>
> ⛔ **Nol keputusan rancangan · nol kode · nol DDL · nol daftar kolom usulan · nol nomor baris
> XML · nol butir ditutup sendiri · nol nilai nama orang.**

---

## Label dan cara baca

| Label | Artinya |
| --- | --- |
| `[terverifikasi]` | ada **path + class + nama rule**, dan perintah audit yang menghasilkannya |
| `[dugaan]` | pembacaan saya, ⛔ **belum terbukti** |
| `[terbuka]` | belum terjawab — wajib masuk Bab C |

**Pemilik:** `[work owner]` · `[DBA]` · `[Product+Underwriting]` · `[Finance]` · `[IAM]` ·
`[pengembang Pega lama]` · `[pemilik export Pega]`

⛔ **Dua artefak turunan TIDAK dipakai sebagai sumber** — `Struktur_MenuNBTreatyIn.xlsx` di akar
korpus, dan **salinannya** di folder `Claude outputs\`. Keduanya buatan tim, bukan ekspor Pega.
⛔ Tidak ada satu pun `[terverifikasi]` di berkas ini yang bersandar padanya.

---

## §0 — Sensus, dan ujian atas instrumennya

### 0.1 Sensus tipe rule — dihitung DUA CARA, keduanya SEPAKAT

⭐ **Jendela: seluruh 278 berkas `.xml` di bawah `NB Treaty In\`.**

| Tipe rule | Cara A — **nama subfolder** | Cara B — **field 1 `pzOriginalInstanceKey` pertama per berkas** | Sepakat? |
| --- | ---: | ---: | :---: |
| Activity | 92 | `RULE-OBJ-ACTIVITY` 92 | ✅ |
| When | 75 | `RULE-OBJ-WHEN` 75 | ✅ |
| RDBList | 41 | `RULE-CONNECT-SQL` 41 | ✅ |
| Section | 25 | `RULE-HTML-SECTION` 25 | ✅ |
| ReportDefinition | 15 | `RULE-OBJ-REPORT-DEFINITION` 15 | ✅ |
| DataTransform | 12 | `RULE-OBJ-MODEL` 12 | ✅ |
| FlowAction | 9 | `RULE-OBJ-FLOWACTION` 9 | ✅ |
| Harness | 6 | `RULE-HTML-HARNESS` 6 | ✅ |
| ⭐ **DecisionTable** | **2** | `RULE-DECLARE-DECISIONTABLE` **2** | ✅ |
| Flow | 1 | `RULE-OBJ-FLOW` 1 | ✅ |
| **TOTAL** | **278** | **278** | ✅ |

Perintah audit cara A:
`for d in */; do find "$d" -name '*.xml' | wc -l; done`
Perintah audit cara B: untuk tiap berkas,
`grep -o '<pzOriginalInstanceKey>[A-Z-]*' "$f" | head -1`, lalu `sort | uniq -c`.
`[terverifikasi]` **Nol berkas tanpa `pzOriginalInstanceKey`** — jadi cara B meliput seluruh 278.

⚠️ **Selisih terhadap brief, dan ia bukan selisih korpus:** daftar D1 di brief menyebut sembilan
tipe *(Activity 92 · When 75 · RDBList 41 · Section 25 · ReportDefinition 15 · DataTransform 12 ·
FlowAction 9 · Harness 6 · Flow 1)* yang berjumlah **276**, ⛔ **bukan 278** — ⭐ **`DecisionTable`
2 luput dari daftar itu.** ⚠️ Tetapi `discovery/modules/NB Treaty In.md` §3 **sudah mencatatnya
benar** *("2 DecisionTable")*. ⭐ Jadi yang tidak lengkap adalah kutipan di brief, bukan
pencatatan sebelumnya. ⛔ **Bukan alasan berhenti.**

⚠️ **`DecisionTable` bukan tipe sepele di modul ini** — salah satu dari dua isinya adalah
`isApproved`, yang menggerbangi enam shape Decision *(lihat §0.4)*.

### 0.2 Sensus berkas dan objek lain

| Yang dicacah | Angka | Cara kedua |
| --- | ---: | --- |
| berkas `.xml` | **278** | `find`-count vs jumlah kolom cara B ⇒ ✅ |
| subfolder di akar | **11** | ⚠️ tetapi **hanya 10 berisi `.xml`** — lihat di bawah |
| `.xlsx` di akar | **1** | `Struktur_MenuNBTreatyIn.xlsx` |
| ⚠️ `.xlsx` **total di korpus** | ⚠️ **2** | ⛔ **satu lagi di `Claude outputs\`** |
| berkas non-`.xml` | **2** | keduanya `.xlsx` di atas |

⚠️ **Subfolder ke-11 bernama `Claude outputs`** dan ⛔ **bukan tipe rule** — ia memuat
`Struktur_MenuNBTreatyIn.xlsx` **275.860 B**, berkas non-Pega **di dalam** folder korpus.
⭐ `[terverifikasi]` Ini **sudah terdaftar sebagai OQ-054**, bukan temuan baru. ⛔ **Tidak dibuka.**

**Lima berkas terbesar** `[terverifikasi]` — perintah:
`find . -name '*.xml' -printf '%s %p\n' | sort -rn | head -5`

| Berkas | Ukuran |
| --- | ---: |
| `Section\GeneralPolicyTreatyIn.xml` | 1.881 KB |
| `Section\DetailPolicyTreatyIn.xml` | 1.871 KB |
| `Section\DetailPolicyTreatyInNonProportional.xml` | 1.624 KB |
| `Section\GeneralDeptHeadTreatyIn_UW.xml` | 1.547 KB |
| `Section\DetailDeptHeadTreatyIn_UW.xml` | 1.530 KB |

⭐ **Kelimanya `Section`** — ⚠️ artinya berat modul ini ada di **layar**, dan layarnya
**tidak saya urai** *(lihat §0.3)*.

### 0.3 ⚠️ Cakupan — apa yang benar-benar saya baca

⛔ **Ini bagian yang paling menentukan nilai ronde ini, dan saya menaruhnya di depan.**

| Golongan | Berkas | Cara |
| --- | ---: | --- |
| ⭐ **Dibaca UTUH, isi SQL-nya dikutip** | **5** | kelima blok PL/SQL — `GetSequenceNumber_SQL` · `InsertHistoryAkseptasiPega_Sql` · `InsertViewSuggest_SQL` · `SavePolisTreatyIn_SQL` · `SaveTreatyIn` |
| ⭐ **Disisir pengurai penuh** *(SQL diekstrak dan diurai)* | **41** | seluruh `RDBList\` |
| **Disisir pola atas seluruh teks** | **278** | slot `CARI` berawalan halaman · `pzOriginalInstanceKey` · `pyTelephone` |
| ⛔ **TIDAK dibuka sama sekali** | **25** | ⛔ **seluruh `Section\`** — termasuk kelima berkas terbesar |
| ⛔ **TIDAK dibuka sama sekali** | **6** | ⛔ **seluruh `Harness\`** |
| ⛔ **TIDAK dibuka** | **2** | kedua `.xlsx` — ⭐ **sengaja**, artefak turunan |
| ⭐ **Dibaca dari catatan lama, bukan dari korpus** | **15** | rule jalur inti — sudah ditelusur di D2, ronde ini **mengkonfirmasi**, bukan membaca ulang |

⚠️⚠️ **Akibat yang wajib disadari:** `Section` + `Harness` = **31 berkas**, dan ⛔ **nol dari
keduanya saya buka.** ⭐ Kelima berkas terbesar modul ini ada di situ. ⛔ **Apa pun yang hanya hidup
di dalam layar — validasi di layar, medan tersembunyi, tombol yang menyaring jabatan — belum
terbaca sama sekali di ronde ini.**

⚠️ **92 Activity dan 75 When juga belum dibuka satu per satu.** Yang saya lakukan atas keduanya
hanyalah **sisiran pola**, bukan pembacaan langkah.

### 0.4 ⛔ Ujian instrumen — dan instrumen saya GAGAL DUA KALI

⛔⛔ `CLAUDE.md` §4a menuntut instrumen diuji atas butir yang jawabannya **sudah diketahui**.
⭐ **Saya mengujinya, dan itu menyelamatkan Bab A dari dua angka palsu.**

**Butir uji:** `discovery/modules/NB Treaty In.md` §6 mencatat `[terverifikasi]` **tiga** stored
procedure. Instrumen saya melaporkan **empat**.

| # | Yang instrumen saya katakan | Yang sebenarnya, sesudah SQL-nya dibaca | Sebab galatnya |
| --- | --- | --- | --- |
| ⛔ **1** | `POOLDATA.HISTORYAKSEPTASIPRODUCTION` adalah **procedure berparameter 15** | ⛔ **SALAH.** Ia **TABEL**, disasar `INSERT INTO` dengan **15 KOLOM** | pola `schema.NAMA(` tidak membedakan **panggilan procedure** dari **`INSERT INTO tabel (daftar kolom)`** |
| ⛔ **2** | `PROC_GENERATE_SEQUENCE_NUMBER` berparameter **6** | ⛔ **SALAH — 5.** *(dua di antaranya `out`)* | pemisah koma memotong **di dalam** `TO_DATE({…},'DD/MM/YYYY')`, sehingga satu argumen terhitung dua |

⭐ **Cara kedua yang dipakai untuk memperbaikinya**, dan ia benar-benar berbeda: buang lebih dulu
seluruh pernyataan `INSERT INTO … ;` dari blok, **baru** cari `schema.NAMA(`. Hasilnya
**tepat 3 procedure**, dan dua blok sisanya dinyatakan **DML langsung**. ✅ **Sepakat dengan
catatan D3.**

⚠️ **Kalimat CLAUDE.md §4a berlaku harfiah di sini** — *"Verifikator yang instrumennya sendiri
belum diverifikasi memproduksi tuduhan palsu."* ⛔ Tanpa ujian ini, Bab A akan menyodorkan
**procedure keempat yang tidak ada** ke rapat keputusan work owner.

---

## §0.5 — Konfirmasi empat batas pengetahuan yang sudah tercatat

⭐ Brief meminta **mengkonfirmasi dan memperdalam**, bukan menemukan ulang. Hasilnya:

| § | Butir | Masih berlaku? | Pendalaman ronde ini |
| --- | --- | :---: | --- |
| **5.1** | `SavePolisTreatyIn_SQL` → `POOLDATA.PEGA_JSON_POLIS_TREATYIN`, **8 parameter**, `COMMIT` di dalam blok PL/SQL | ✅ **YA, persis** | ⭐ SQL dikutip utuh di Bab A; ⭐ **`COMMIT` itu milik blok Pega, bukan terbukti milik procedure** — perbedaan yang belum pernah dipisahkan |
| **5.2** | **BLOCKER** `serviceInsertArasapas_act` memanggil identitas yang variannya tidak ada di modul ini → OQ-025 | ✅ **YA** | ⛔ **tetap blocker**; ⛔ varian modul lain **tidak dipinjam**; ⭐ telusur tetap **berhenti** di situ |
| **5.3** | `@ASM.GetPageJSONString()` menghasilkan JSON yang ditulis ke basis data; source tidak ada di korpus → OQ-012 | ✅ **YA** | ⭐ **dan sekarang terlihat lebih buruk:** nilainya masuk lewat **`InputData.CARI3`**, slot yang **artinya berpindah antar rule** — lihat §0.6 |
| **5.4** | `isApproved` ada sebagai **`RULE-OBJ-WHEN` DAN `RULE-DECLARE-DECISIONTABLE`** berclass+nama sama → OQ-026 | ✅ **YA** | ⭐ `[terverifikasi]` keduanya terbaca di sensus cara B, dan ⭐ **`pxUpdateDateTime`-nya BERBEDA** — `#20170608T022121.683` *(When)* lawan `#20170608T024733.462` *(DecisionTable)*, ⚠️ **selisih ±26 menit di hari yang sama** |

⭐ **Yang ditambahkan §5.4 oleh ronde ini:** karena identitas rule berempat bagian, kedua berkas itu
**bukan satu rule yang terbaca dua kali** — ⛔ **ia dua rule berbeda** yang berebut nama yang sama.
⚠️ Selisih 26 menit `[dugaan]` menandakan yang kedua dibuat **segera sesudah** yang pertama, ⛔ tetapi
**mana yang menang saat flow berjalan tetap tidak terbaca.**

### §0.6 ⭐⭐ Temuan baru — halaman parameter generik dipakai bersama, artinya berpindah per rule

⛔⛔ **Ini temuan terberat ronde ini, dan ia memperdalam OQ-059 sampai ke titik yang menggigit.**

`[terverifikasi]` Slot `CARI` di modul ini **selalu** berawalan halaman *(aturan §4a dipatuhi:
dicari sebagai `<halaman>.CARI<n>`, ⛔ bukan `CARI<n>` telanjang)*. Halamannya **14 buah**:

| Halaman pembawa | Kemunculan slot |
| --- | ---: |
| ⭐ `InputData` | **23** |
| `SearchCurrencyValueIn` | 8 |
| `ParamInProdCase` | 7 |
| `InsertHistory` | 5 |
| `Track` · `ParamSeq` | 3 · 3 |
| `ParamData` · `OldID` · `InputDataCredit` | 2 masing-masing |
| `TempOccupation` · `TempObjItemIn` · `TempNopolis` · `SearchClient` · `DataSearch` | 1 masing-masing |

Perintah audit: `grep -rhoE '\{[A-Za-z][A-Za-z0-9_.]*\.CARI[0-9]+' … | sed 's/\.CARI[0-9]*$//' | sort | uniq -c`

⭐ **Dan inilah yang berbahaya** `[terverifikasi]`: halaman `InputData` dipakai **18 rule berbeda**
dengan **15 slot berbeda** *(`CARI1`–`CARI10`, `CARI20`, `CARI21`, `CARI23`, `CARI24`, `CARI44`)*.
⛔ **Satu slot yang sama membawa muatan berbeda di rule berbeda.** Buktinya sepasang:

| Rule | Slot | Muatannya di situ |
| --- | --- | --- |
| `RDBList\SavePolisTreatyIn_SQL.xml` | ⭐ `InputData.CARI3` | **JSON** — parameter ke-7 procedure |
| `RDBList\InsertViewSuggest_SQL.xml` | ⭐ `InputData.CARI3` | **kolom `POSISI`** pada tabel riwayat |

`[terverifikasi]` `InputData.CARI3` muncul di **7 rule**; `InputData.CARI1` di **14 rule**.
Perintah: `grep -rl 'InputData\.CARI3[^0-9]' . --include='*.xml'`

⚠️ **Akibatnya pada migrasi, dan ia bukan soal gaya:** ⛔ **awalan halaman TIDAK CUKUP** untuk
mengunci arti sebuah slot. ⭐ Artinya ditentukan oleh **rule yang sedang berjalan**, dan
⛔ **rule mana yang mengisi slot itu sebelum SQL dipanggil belum saya telusur.** ⚠️ Memindahkan ini
ke Go dengan memetakan `InputData.CARI3` ke **satu** kolom akan **salah di enam rule lainnya**.

---

## Bab A — Stored procedure dan SQL mentah

⭐ **Kenapa bab ini ada:** work owner sedang memutuskan apakah pemanggilan stored procedure dihapus
seluruhnya. ⛔ Tanpa badan procedure, memindahkan logikanya ke Go berarti **menebak** — dan menebak
dilarang.

### A.1 Tiga stored procedure — dan ⛔ NOL badannya ada

`[terverifikasi]` **41 dari 41** rule `RDBList` punya `<pyBrowseSQL>` berisi
*(`ls RDBList | wc -l` = 41; blok SQL terekstrak = 41)*. **5** di antaranya blok PL/SQL
`BEGIN…END`; **3** memanggil procedure, **2** melakukan DML langsung.

| Procedure | Rule pemanggil — class / nama / tipe | Path | Param | `COMMIT`? | Badan ada di artefak `[data DBA]`? |
| --- | --- | --- | ---: | --- | --- |
| ⛔ **`POOLDATA.PEGA_TREATY_IN`** | `ASM-FW-GISFW-Int-TREATY_IN` / `ASM-FW-GISFW-INT-TREATY_IN!ASM!SAVETREATYIN` / `RULE-CONNECT-SQL` | `RDBList\SaveTreatyIn.xml` | ⭐ **24** *(3 `out`)* | ⚠️ **ya, di blok Pega** | ⛔ **BELUM** |
| ⛔ **`POOLDATA.PEGA_JSON_POLIS_TREATYIN`** | `ASM-FW-GISFW-Int-POLISTREATYIN` / `ASM-FW-GISFW-INT-POLISTREATYIN!ASM!SAVEPOLISTREATYIN_SQL` / `RULE-CONNECT-SQL` | `RDBList\SavePolisTreatyIn_SQL.xml` | **8** *(1 `out`)* | ⚠️ **ya, di blok Pega** | ⛔ **BELUM** |
| ⚠️ **`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`** | `ASM-FW-GISFW-Int-policyjson` / `ASM-FW-GISFW-INT-POLICYJSON!RNM!GETSEQUENCENUMBER_SQL` / `RULE-CONNECT-SQL` | `RDBList\GetSequenceNumber_SQL.xml` | **5** *(2 `out`)* | ⚠️ **ya, di blok Pega** — ⭐ **tetapi procedure-nya sendiri TIDAK commit** `[data DBA]` | ⚠️ **SEBAGIAN** — perilaku commit-nya tercatat, ⛔ **badannya tidak** |

⭐⭐ **Pembedaan yang belum pernah dipisahkan, dan ia mengikat rancangan transaksi Go:**
⚠️ `COMMIT` pada ketiga baris di atas **berada di dalam blok PL/SQL yang ditulis Pega**, ⛔ **bukan
terbukti berada di dalam procedure.** `[terverifikasi]` Untuk `PROC_GENERATE_SEQUENCE_NUMBER`
keduanya **terbukti berbeda**: register OQ-002 mencatat `[data DBA]` *"commit sendiri: tidak"*,
padahal blok Pega yang memanggilnya **ber-`COMMIT`**. ⛔ **Untuk dua procedure lainnya, mana yang
commit belum diketahui** — dan itu menentukan apakah Go boleh membungkus keduanya dalam satu
transaksi.

### A.2 Dua blok PL/SQL yang menulis LANGSUNG, bukan lewat procedure

| Rule | Path | Sasaran | Kolom diisi | `COMMIT` |
| --- | --- | --- | ---: | --- |
| `ASM-FW-GISFW-Work` / `…!RNM!INSERTVIEWSUGGEST_SQL` | `RDBList\InsertViewSuggest_SQL.xml` | ⭐ `POOLDATA.historyakseptasiproduction` | **15** | **ya, di blok** |
| *(RDBList)* | `RDBList\InsertHistoryAkseptasiPega_Sql.xml` | ⭐ `HISTORYAKSEPTASIPEGA` — ⚠️ **tanpa awalan schema** | **6** | **ya, di blok** |

⭐ **`InsertHistoryAkseptasiPega_Sql` adalah satu-satunya tempat di modul ini yang DDL-nya sudah
diketahui** `[data DBA]`: tabelnya **`POOLDATA.HISTORYAKSEPTASIPEGA`**, tablespace `POOLMASTER`,
**7 kolom**, lima index semuanya bertumpu `ID_PEGA`.
⚠️ **Rule ini mengisi 6 dari 7** — ⛔ kolom `OPERATORID` **tidak diisi olehnya**.
⚠️ Dan ia menulis **tanpa awalan schema**, sehingga **schema tujuannya bergantung pada pengguna
koneksi** — ⛔ tidak terbaca dari rule.

### A.3 Satu sequence, dan satu fungsi Pega tanpa source

| Objek | Dipakai di | Keadaan |
| --- | --- | --- |
| `POOLDATA.JSON_POLIS_TREATYIN_SEQ.NEXTVAL` | `RDBList\GenerateNoPolicy.xml` | ⭐ membentuk nomor polis dari `dual`; ⛔ **isi `InputData.CARI20` tidak terlihat diisi** → OQ-059 |
| ⛔ `@ASM.GetPageJSONString()` | `Activity\SaveJsonPolisTreatyIn_Act.xml` langkah 6 | ⛔ **source tidak ada di korpus** — nol tipe rule fungsi/library → OQ-012 |

### A.4 ⭐ Rekap Bab A dalam satu kalimat untuk keputusan work owner

⛔ **Tiga procedure dipanggil; NOL badannya ada di `OUTPUT_HASIL_RNM`.**
`[terverifikasi]` Perintah audit: `grep -rli 'CREATE OR REPLACE\|CREATE PROCEDURE\|CREATE FUNCTION'`
atas seluruh `OUTPUT_HASIL_RNM` **di luar `.scratch\`** ⇒ ⛔ **nol berkas.**
⭐ Yang sudah diserahkan DBA *(register OQ-002, 2026-09-15)* adalah **tiga procedure lain** —
`PEGA_M_LIFE_PREMIUM_SUMMARY`, `INSERTJSONPOLISLIFE`, `INSERTJSONOFFERLIFE` — ⛔ **ketiganya jalur
Life, bukan salah satu dari tiga yang dipanggil modul ini.**

⚠️ **Jadi keputusan "hapus pemanggilan stored procedure" belum dapat dinilai untuk modul ini.**
⛔ **24 parameter** `PEGA_TREATY_IN` adalah permukaan terlebar, dan apa yang dilakukannya
**tidak terbaca sama sekali.**

---

## Bab B — Tabel Oracle yang disentuh

`[terverifikasi]` **31 objek berbeda**, dari **41 blok SQL** di **41 rule**.

⭐ **Dihitung DUA CARA yang berbeda, dan hasilnya sepakat:**
*(a)* pola `FROM` / `JOIN` / `INSERT INTO` / `UPDATE` / `DELETE FROM` / `MERGE INTO` ⇒ **31**;
*(b)* daftar tabel **berkoma** sesudah `FROM` sampai kata kunci penutup ⇒ **30**, ⭐ selisihnya
**tepat** `POOLDATA.HISTORYAKSEPTASIPRODUCTION` — ⭐ **wajar, sebab ia hanya disasar `INSERT`, tidak
pernah muncul sesudah `FROM`**. ✅ **Nol objek baru ditemukan cara (b).**

⚠️ **Kenapa cara (b) perlu:** `JOIN` = **0** kemunculan sementara `FROM` = **52**, ⛔ jadi
gabungannya **gaya koma lama**, dan pola `FROM` polos hanya menangkap tabel **pertama**.
`[terverifikasi]` Hanya **1** blok berdaftar dua tabel — `RDBList\GetDataByID_SQL.xml` — dan
tabel keduanya **sudah tercakup**.

**Uji kelengkapan dua arah atas penulis:** berkas-menulis menurut pola = **2**; berkas ber-DML
menurut kata kunci = **2**; ⛔ **A\B kosong, B\A kosong.** ✅

### B.1 Yang DITULIS — hanya dua, dan keduanya tabel riwayat

| Objek | Schema | Operasi | Rule |
| --- | --- | --- | --- |
| ⭐ `POOLDATA.HISTORYAKSEPTASIPRODUCTION` | `POOLDATA` | **TULIS** | `RDBList\InsertViewSuggest_SQL.xml` |
| ⭐ `HISTORYAKSEPTASIPEGA` | ⚠️ **tanpa awalan** | **BACA + TULIS** | tulis: `RDBList\InsertHistoryAkseptasiPega_Sql.xml` · baca: `RDBList\GetFlagReject_SQL.xml` |

⛔⛔ **Dan inilah bentuk modul ini yang paling penting untuk dicatat:**
`[terverifikasi]` **Dari 31 objek, hanya 2 ditulis lewat SQL yang terbaca.** ⭐ **29 sisanya
hanya DIBACA.** ⚠️ **Seluruh penyimpanan data treaty yang sebenarnya lewat tiga procedure di
Bab A** — dan ⛔ **sasarannya tidak terbaca.**

### B.2 Yang hanya DIBACA — 29 objek

| Schema | Objek |
| --- | --- |
| `POOLDATA` | `M_TREATY_IN` *(2 rule)* · `M_TREATY_IN_EDM` *(2)* · `M_TREATY_IN_DETAIL_EDM` · ⭐ `M_TREATY_OUT` *(2)* · `REINSURANCETYPE` *(2)* · `PROPORTIONALARRG` · `TREATYGROUP` · `TREATYINPRODUCTION` · `TANGGAL_CLOSING` · `KODE_PRODUKSI` · `CURRENCY` |
| `DATAPEGA` | `PC_ASM_FW_GCNMFW_WORK` — ⚠️ **tabel kerja Pega dibaca sebagai SQL biasa** *(`RDBList\GetCountClaim.xml`)* |
| ⚠️ **tanpa awalan** | `AGENT` · `BUSINESS` · `CATEGORY_ATTACH_REAS` *(2)* · `CLIENT` · `CLIENTEMAIL` · `CURRENCY` *(3)* · `FACINPRODUCTION` · `JSON_POLIS` · `M_REINSURANCETYPE` · `OCCUPATION` · `PROPORTIONALARRG` *(2)* · `REINSURANCETYPE` *(2)* · `RW` · `TREATYBUSINESS` *(2)* · `TREATYCONTRACT` · `TREATYEXCHANGEYEARLY` · `V_JN_OBJ_ITEM` |

⚠️⚠️ **Empat nama muncul DUA KALI — sekali berawalan `POOLDATA`, sekali tanpa awalan:**
`CURRENCY` · `REINSURANCETYPE` · `PROPORTIONALARRG` · *(dan `HISTORYAKSEPTASI*` dalam dua ejaan
berbeda)*. ⛔ **Apakah itu satu tabel yang dieja dua cara, atau dua tabel berbeda, TIDAK TERBACA
dari rule** — awalan yang hilang diisi oleh **schema pengguna koneksi**. ⭐ Ini persis jebakan
sensus nomor 4 di `CLAUDE.md` §4a *(satu objek dieja beberapa cara)*, dan di sini ia **nyata**.

### B.3 ⛔ Objek yang struktur kolomnya TIDAK DAPAT DINYATAKAN

⭐ Brief meminta ini ditandai. Daftarnya:

| Objek | Kenapa tidak dapat dinyatakan |
| --- | --- |
| ⛔ **seluruh sasaran `PEGA_TREATY_IN`** | ⛔ **tidak diketahui sama sekali** — 24 parameter masuk, badan procedure tidak ada |
| ⛔ **seluruh sasaran `PEGA_JSON_POLIS_TREATYIN`** | ⛔ 8 parameter masuk; ⚠️ salah satunya **JSON utuh** → ⛔ struktur JSON juga tidak terbaca *(OQ-012)* |
| ⛔ **seluruh sasaran `PROC_GENERATE_SEQUENCE_NUMBER`** | ⛔ badan tidak ada; ⚠️ hanya perilaku commit-nya tercatat |
| ⚠️ `JSON_POLIS` | ⭐ **dibaca** di sini *(`RDBList\GetTgl_InputJsonPolis_SQL.xml`)*, ⛔ **tetapi apa yang menulisnya di modul ini tidak terbaca** — `[dugaan]` lewat `PEGA_JSON_POLIS_TREATYIN`, ⛔ **belum terverifikasi** |
| ⚠️ 17 objek **tanpa awalan schema** | ⛔ schema-nya ditentukan pengguna koneksi, bukan rule |

⭐ **Satu-satunya pengecualian, dan ia satu-satunya:** `POOLDATA.HISTORYAKSEPTASIPEGA` —
**7 kolom + 5 index terbaca** `[data DBA]`.

---

## Bab C — Register `[terbuka]` berpemilik

⭐ **Nomor OQ diambil dari register `discovery/open-questions.md`**, ⛔ bukan dari ingatan.
`[terverifikasi]` Register memuat **72 OQ** — **70 terbuka**, **2 terjawab**. ⭐ **Nomor bebas
berikutnya: OQ-073.**

⛔⛔ **NOL butir saya tutup di ronde ini.** Penutupan hanya oleh pemilik peran.

### C.1 Butir yang SUDAH terdaftar dan dikonfirmasi masih berlaku

| OQ | Butir | Pemilik | Kenapa memblokir | Akibat bila dijawab salah |
| --- | --- | --- | --- | --- |
| ⛔ **002** | Badan **3 procedure** yang dipanggil modul ini tidak ada | `[DBA]` | ⛔ **seluruh penyimpanan treaty lewat sini** — 24 parameter tanpa badan | ⛔ logika tersalin separuh; data tersimpan berbeda dari Pega tanpa ada yang tahu |
| ⛔ **013** | Batas transaksi ada di sisi basis data; `COMMIT` di blok **atau** di procedure | `[DBA]` | ⛔ menentukan apakah Go boleh satu transaksi | ⛔ commit di tengah ⇒ **data separuh jadi permanen** saat langkah berikutnya gagal |
| ⛔ **025** | `SERVICEINSERTARASAPAS_ACT` dipanggil, variannya **tidak ada** di modul ini — **BLOCKER** | `[Product+Underwriting]` · `[pemilik export Pega]` | ⛔ efek keluar terakhir jalur realisasi **tidak dapat dinyatakan** | ⛔ integrasi keuangan hilang atau terkirim dua kali |
| ⛔ **026** | `isApproved` ada sebagai **When** dan **DecisionTable**, identitas sama | `[Product+Underwriting]` | ⛔ **gerbang pertama seluruh alur** — enam shape Decision memakainya | ⛔ kasus yang seharusnya berhenti **lolos**, atau sebaliknya |
| ⛔ **012** | Struktur JSON yang ditulis; `@ASM.GetPageJSONString()` tanpa source | `[DBA]` · `[pengembang Pega lama]` | ⛔ isi kolom JSON adalah **muatan utama** yang disimpan | ⛔ data lama tidak terbaca sistem baru |
| ⚠️ **059** | Slot `CARI1`…`CARI30` generik menuju SQL | `[pengembang Pega lama]` | ⛔ **diperdalam ronde ini** — lihat C.2 butir baru #1 | ⛔ kolom terpetakan ke muatan yang salah |
| ⚠️ **023** | Sub-graf komite klaim + Direktur **tanpa connector masuk** | `[Product+Underwriting]` | ⛔ apakah persetujuan Direktur **ada atau tidak** di alur ini | ⛔ satu jenjang persetujuan hilang, atau dibangun padahal mati |
| ⚠️ **024** | Pemetaan Assignment → workbasket tidak terbaca; routing `Custom` | `[IAM]` · `[Product+Underwriting]` | ⛔ siapa memegang giliran tiap tahap | ⛔ pekerjaan masuk ke kotak yang salah |
| ⛔ **027** | `OperatorID.pyTelephone` menyimpan **kode peran** | `[IAM]` | ⛔ satu-satunya penentu peran yang terbaca | ⛔ wewenang salah orang |
| ⚠️ **021** | Identitas orang ter-hardcode sebagai guard | `[IAM]` | ⛔ 12 berkas ber-`pyUserIdentifier` | ⛔ orang yang sudah pindah tetap berwenang |
| ⚠️ **028** | **Keenam** Assignment `WorkBasket`, nol `WorkList` | `[Product+Underwriting]` | ⛔ model penugasan menentukan bentuk kotak masuk | ⛔ antrean pribadi dibangun padahal semuanya antrean bersama |
| ⚠️ **020** | Arti kode penggerbang *(`IsApproved=1`, `LetterNo="TREATYINDEPTHEAD"`, `BusinessCode!="40"`, `isFOR`=`EDM`/`POLICY`)* | `[Product+Underwriting]` | ⛔ percabangan bergantung kode tanpa arti | ⛔ cabang salah dipilih |
| ⚠️ **022** | `POOLDATA.M_TREATY_OUT` disentuh **8 berkas** modul ini, nol di modul bernama outward | `[Product+Underwriting]` | ⚠️ batas konteks | ⚠️ tanggung jawab treaty outward salah tempat |
| ⚠️ **011** | 533 identitas rule berisi berbeda antar modul | `[pemilik export Pega]` | ⛔ mana yang berlaku di production | ⛔ perilaku disalin dari varian yang tidak dipakai |
| ⚠️ **018** | Korpus memuat hostname **DEV**, dirakit dari >1 server | `[pemilik export Pega]` | ⚠️ **seluruh berkas ini** belum tentu cerminan production | ⚠️ spesifikasi dibangun atas lingkungan yang salah |
| **009** | Rule tinggal di class bawaan Pega, bukan class aplikasi | `[pengembang Pega lama]` | ⚠️ cakupan migrasi | ⚠️ rule terlewat atau terbawa berlebih |
| **054** | Folder `Claude outputs` berisi berkas non-Pega di dalam korpus | `[pemilik export Pega]` | ⚠️ kebersihan korpus | ⚠️ artefak turunan disalahbaca sebagai sumber |

### C.2 ⭐ Butir yang LAHIR di ronde ini — belum bernomor register

⛔ **Saya tidak memberi nomor OQ sendiri** — register ada di berkas lain dan ronde ini hanya boleh
menulis satu berkas. ⭐ **Nomor bebas berikutnya OQ-073**; penomoran resmi menunggu pemilik register.

| # | Butir baru | Pemilik | Kenapa memblokir | Akibat bila dijawab salah |
| --- | --- | --- | --- | --- |
| ⭐⭐ **1** | **Halaman parameter generik dipakai bersama 18 rule, dan satu slot berarti berbeda per rule** — `InputData.CARI3` = **JSON** di satu rule, kolom **`POSISI`** di rule lain | `[pengembang Pega lama]` | ⛔ **awalan halaman tidak cukup** mengunci arti; pemetaan slot→kolom **tidak dapat dibuat sekali untuk semua** | ⛔⛔ satu pemetaan salah **merusak enam rule lain** yang memakai slot sama |
| ⭐ **2** | **`COMMIT` di blok Pega lawan `COMMIT` di dalam procedure** — untuk `PROC_GENERATE_SEQUENCE_NUMBER` keduanya **terbukti berbeda**; untuk dua procedure lain **belum diketahui** | `[DBA]` | ⛔ menentukan apakah dua penyimpanan boleh satu transaksi | ⛔ separuh data permanen saat sisanya gagal |
| ⭐ **3** | **Empat nama tabel muncul dalam dua ejaan** — berawalan `POOLDATA` dan tanpa awalan *(`CURRENCY`, `REINSURANCETYPE`, `PROPORTIONALARRG`, dan dua ejaan tabel riwayat)* | `[DBA]` | ⛔ satu tabel atau dua? schema diisi pengguna koneksi | ⛔ menulis ke schema yang salah, atau membaca tabel yang salah |
| ⭐ **4** | **`HISTORYAKSEPTASIPEGA` ditulis tanpa awalan schema**, dan rule ini mengisi **6 dari 7 kolom** — `OPERATORID` tidak diisi | `[DBA]` · `[pengembang Pega lama]` | ⚠️ kolom yang tidak pernah diisi **dari modul ini** — apakah modul lain mengisinya | ⚠️ jejak audit tidak lengkap, atau kolom mati ikut termigrasi |
| ⚠️ **5** | **`DATAPEGA.PC_ASM_FW_GCNMFW_WORK` dibaca sebagai SQL biasa** — tabel kerja internal Pega diakses langsung *(`RDBList\GetCountClaim.xml`)* | `[pengembang Pega lama]` | ⛔ tabel ini **hilang** begitu Pega ditinggalkan | ⛔ perhitungan yang bertumpu padanya **mati tanpa suara** |
| ⚠️ **6** | **Nol dari 25 `Section` dan nol dari 6 `Harness` dibuka** — kelima berkas terbesar modul ada di situ | `[pengembang Pega lama]` | ⚠️ validasi layar dan penyaring tombol **belum terbaca sama sekali** | ⚠️ aturan yang hanya hidup di layar hilang dalam migrasi |
| ⚠️ **7** | **`isApproved` versi When dan DecisionTable dibuat berselisih ±26 menit** di hari yang sama | `[pengembang Pega lama]` | ⚠️ petunjuk mana yang menggantikan mana | ⚠️ gerbang pertama alur memakai aturan yang salah |

### C.3 Rekap jumlah per pemilik

⭐ **24 butir** — **17 terdaftar** *(C.1)* + **7 baru** *(C.2)*. ✅ **Dihitung dua cara**
*(baris tabel yang memuat pemilik, dan kemunculan literal nama pemilik)*; **keduanya sepakat**.

> ⛔⛔ **RALAT — ditangkap oleh sensus berkas ini sendiri, sebelum dilaporkan.**
> Tabel di bawah ini sebelumnya memuat **tiga angka yang salah**, dan kalimatnya dikutip utuh di
> sini alih-alih dihapus:
>
> *"`[DBA]` **5** · `[pengembang Pega lama]` **7** · `[Product+UW]` **6** …
> Jumlah sebutan pemilik **25**, jumlah butir 24."*
>
> ⭐ **Yang benar: `[DBA]` 6 · `[pengembang Pega lama]` 8 · `[Product+Underwriting]` 7 ·
> sebutan pemilik 28.**
>
> ⚠️ **Dan sebab galatnya persis jebakan sensus nomor 4 di `CLAUDE.md` §4a** — *"nama
> dicocokkan peka huruf besar-kecil; satu objek bisa dieja beberapa cara."* ⛔ **Saya sendiri
> mengeja pemilik yang sama dua cara** — `[Product+UW]` di **4 tempat** dan
> `[Product+Underwriting]` di tempat lain — sehingga cacahan pertama kehilangan empat baris.
> ⭐ **Ejaannya kini diseragamkan menjadi `[Product+Underwriting]` di seluruh berkas.**

| Pemilik | Butir | Nomornya |
| --- | ---: | --- |
| ⛔ `[pengembang Pega lama]` | ⭐ **8** | OQ-059 · OQ-009 · OQ-012 *(bersama)* · baru #1 · #2 *(bersama)* · #4 *(bersama)* · #5 · #6 · #7 |
| ⛔ `[Product+Underwriting]` | ⭐ **7** | OQ-025 *(bersama)* · OQ-026 · OQ-023 · OQ-024 *(bersama)* · OQ-028 · OQ-020 · OQ-022 |
| ⛔ `[DBA]` | ⭐ **6** | OQ-002 · OQ-013 · OQ-012 *(bersama)* · baru #2 · #3 · #4 *(bersama)* |
| ⚠️ `[pemilik export Pega]` | **4** | OQ-011 · OQ-018 · OQ-054 · OQ-025 *(bersama)* |
| ⚠️ `[IAM]` | **3** | OQ-027 · OQ-021 · OQ-024 *(bersama)* |
| `[Finance]` | ⭐ **0** | ⛔ **tidak ada butir uang yang lahir di ronde ini** — ⚠️ sebab rantai hitung modul ini **belum saya sentuh**, bukan sebab bersih |
| `[work owner]` | ⭐ **0** | ⭐ seluruh butir jatuh ke pemilik yang lebih spesifik |

⚠️ **Jumlah sebutan pemilik 28, jumlah butir 24** — ⭐ sebab **empat butir berpemilik ganda**.
⛔ Bukan salah hitung.

⭐ **Bab D memuat 17 pertanyaan, dan pembagiannya cocok:** `[DBA]` **P1–P4** *(4)* ·
`[Product+Underwriting]` **P5–P10** *(6)* · `[IAM]` **P11–P13** *(3)* ·
`[pengembang Pega lama]` **P14–P17** *(4)*. ✅ **17 = 4+6+3+4**, bernomor **P1–P17 tanpa nomor
hilang**.

---

## Bab D — PERTANYAAN SIAP KIRIM

⛔⛔ **Bab ini dapat dibawa ke rapat apa adanya.** Penjawabnya **tidak membaca XML**.
⛔ **Nol pertanyaan saya jawab sendiri.** Diurutkan dari yang **paling memblokir**.

---

### Untuk DBA

#### **P1 — Kami butuh isi tiga program penyimpan data di basis data**

Ada tiga program yang tersimpan **di dalam basis data** — bukan di aplikasi — dan seluruh
penyimpanan data treaty inward melewatinya. Salah satunya menerima **24 keterangan sekaligus**.
Kami dapat melihat keterangan apa yang **dikirim masuk**, ⛔ **tetapi tidak dapat melihat apa yang
dilakukan program itu** — tabel mana yang diisi, kolom mana, pemeriksaan apa yang dijalankan.

**Konteks:** sistem baru harus menyimpan data yang sama, dan tanpa isi program ini kami akan
**menebak** ke mana data pergi.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — naskah ketiga program berikut namanya:
`PEGA_TREATY_IN`, `PEGA_JSON_POLIS_TREATYIN`, `PROC_GENERATE_SEQUENCE_NUMBER`.
⚠️ Untuk yang ketiga, sebagian keterangannya sudah kami terima pada **15 September 2026**; yang
belum adalah naskah lengkapnya.

**Dampak bila salah:** data treaty tersimpan di tempat atau bentuk yang berbeda dari sekarang,
dan **tidak seorang pun akan tahu sampai laporan tidak cocok**.

rujukan: `RDBList\SaveTreatyIn.xml` · `RDBList\SavePolisTreatyIn_SQL.xml` ·
`RDBList\GetSequenceNumber_SQL.xml` — OQ-002

---

#### **P2 — Apakah ketiga program itu menyelesaikan penyimpanannya sendiri, atau menunggu aplikasi?**

Ketika data disimpan, ada satu titik di mana penyimpanan menjadi **permanen dan tidak bisa
dibatalkan**. Kami menemukan bahwa perintah "jadikan permanen" itu ditulis **di sisi aplikasi**,
⛔ tetapi untuk salah satu dari tiga program itu keterangan Anda sebelumnya menyebut bahwa
**program itu sendiri tidak menjadikannya permanen**. ⚠️ Untuk dua program lainnya kami **belum
tahu**.

**Konteks:** bila sebuah program menjadikan datanya permanen sendiri **di tengah** pekerjaan,
maka kegagalan pada langkah berikutnya **tidak dapat membatalkan** yang sudah tersimpan.

**Bentuk jawaban yang diharapkan:** pilihan ganda, **untuk masing-masing dari tiga program**:
**(a)** program menjadikan permanen sendiri · **(b)** program menyerahkannya ke aplikasi ·
**(c)** tergantung jalur, jelaskan.

**Dampac bila salah:** ⛔ separuh data menjadi permanen sementara separuh lainnya batal — dan
selisihnya hanya ketahuan berbulan-bulan kemudian.

rujukan: `COMMIT` di dalam blok PL/SQL Pega vs di dalam procedure — OQ-013, butir baru #2

---

#### **P3 — Empat nama tabel muncul dalam dua ejaan. Satu tabel atau dua?**

Empat nama tabel — mata uang, jenis reasuransi, pengaturan proporsional, dan tabel riwayat —
muncul **dua kali** di dalam aturan: sekali **dengan nama gudang data di depannya**, sekali
**tanpa**. Kami tidak dapat memastikan dari aturan itu apakah keduanya menunjuk **tabel yang sama**.

**Konteks:** bila tanpa awalan berarti "gudang milik pengguna yang sedang menyambung", maka
jawabannya bisa berbeda antara lingkungan uji dan lingkungan sebenarnya.

**Bentuk jawaban yang diharapkan:** untuk tiap nama, satu baris — **nama gudang yang benar**, dan
apakah ada tabel lain bernama sama di gudang berbeda. Contoh nilai: `POOLDATA.CURRENCY` dan
`CURRENCY` ⇒ *"sama"* atau *"berbeda, yang kedua ada di gudang X"*.

**Dampak bila salah:** sistem baru **membaca tabel yang salah** atau **menulis ke gudang yang
salah**, dan keduanya tidak menimbulkan pesan galat.

rujukan: `CURRENCY` · `REINSURANCETYPE` · `PROPORTIONALARRG` · `HISTORYAKSEPTASIPEGA` /
`HISTORYAKSEPTASIPRODUCTION` — butir baru #3

---

#### **P4 — Satu kolom pada tabel riwayat tidak pernah diisi dari modul ini. Siapa mengisinya?**

Tabel riwayat akseptasi punya **tujuh kolom**. Bagian sistem yang kami periksa mengisi **enam**.
Kolom ketujuh — penampung identitas operator — ⛔ **tidak diisi sama sekali dari sini**.

**Konteks:** kami perlu tahu apakah kolom itu diisi bagian lain sistem, atau memang tidak pernah
dipakai — sebab yang tidak pernah dipakai tidak perlu dibawa ke sistem baru.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** diisi bagian lain, sebutkan mana ·
**(b)** tidak pernah dipakai, boleh ditinggalkan · **(c)** dipakai untuk data lama saja.
⭐ Bila memungkinkan: **berapa baris yang kolom itu tidak kosong** di basis data sekarang.

**Dampak bila salah:** jejak audit sistem baru **kehilangan satu keterangan** yang selama ini ada,
atau sebaliknya kami membawa kolom mati selamanya.

rujukan: `POOLDATA.HISTORYAKSEPTASIPEGA` kolom `OPERATORID` ·
`RDBList\InsertHistoryAkseptasiPega_Sql.xml` — butir baru #4

---

### Untuk Product + Underwriting

#### **P5 — Apakah persetujuan Direktur dan jalur klaim benar-benar berjalan?**

Di dalam gambar alur kerja ada **lima kotak** yang menggambarkan sebuah jalur: "Klaim" →
"Manajer Klaim" → "Persetujuan Direktur". ⛔ **Tidak ada satu pun garis yang masuk ke jalur itu.**
⚠️ Jadi entah jalur itu **sudah lama tidak dipakai**, atau **garisnya hilang** ketika sistem lama
diekspor.

**Konteks:** kalau jalur itu hidup, sistem baru harus punya jenjang persetujuan Direktur; kalau
mati, membangunnya berarti menambah jenjang yang tidak pernah ada.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** jalur ini masih dipakai, jelaskan
kapan · **(b)** sudah lama tidak dipakai · **(c)** tidak tahu, perlu diperiksa di sistem berjalan.
⭐ Bila (a): **berapa kali dipakai dalam setahun terakhir.**

**Dampak bila salah:** ⛔ satu jenjang persetujuan **hilang** dari sistem baru, atau sebuah jenjang
**dibangun dan diwajibkan** padahal tidak pernah ada.

rujukan: `Flow\InputRealizationTreatyIn.xml` — `Decision5`/`Assignment5`/`Assignment1`, sub-graf
tanpa connector masuk — OQ-023

---

#### **P6 — Ada dua aturan berbeda dengan nama sama untuk memutuskan "sudah disetujui atau belum". Mana yang berlaku?**

Gerbang **pertama** seluruh alur adalah pemeriksaan "apakah ini sudah benar/disetujui". Di sistem
lama ada **dua aturan terpisah dengan nama yang persis sama** untuk itu, dibuat **berselisih
sekitar 26 menit** pada hari yang sama, dan ⛔ **kami tidak dapat menentukan mana yang sebenarnya
dipakai.**

**Konteks:** enam titik percabangan memakai pemeriksaan ini; bila kami memilih yang salah, arah
seluruh alur bisa berubah.

**Bentuk jawaban yang diharapkan:** **butuh berkas atau demonstrasi** — satu contoh kasus nyata
beserta keterangan apakah ia lolos gerbang pertama, sehingga kami dapat mencocokkan. Atau, bila
diketahui: **aturan mana yang berlaku**, berbentuk daftar syarat.

**Dampak bila salah:** ⛔ pengajuan yang seharusnya **berhenti** malah **diteruskan**, atau
sebaliknya pengajuan sah **tertahan**.

rujukan: `When\isApproved.xml` dan `DecisionTable\isApproved.xml`, keduanya
`ASM-FW-GISFW-WORK!ISAPPROVED` — OQ-026, butir baru #7

---

#### **P7 — Apa arti kode-kode yang menentukan percabangan?**

Beberapa percabangan alur diputuskan oleh **kode tanpa keterangan**. Yang kami temukan:
angka **1** berarti "lolos"; teks **`TREATYINDEPTHEAD`** disimpan **di kolom nomor surat** dan
dipakai sebagai penanda arah; ⛔ kode bisnis **`40`** menghentikan satu langkah penyimpanan; dan
sebuah penanda bernilai **`EDM`** atau **`POLICY`** membedakan dua cara penyimpanan.

**Konteks:** kode-kode ini menggerakkan alur, jadi artinya menentukan perilaku, bukan tampilan.

**Bentuk jawaban yang diharapkan:** satu baris per kode — **apa artinya**, dan **nilai lain apa
yang mungkin muncul**. Contoh: `40` ⇒ *"lini usaha …, dikecualikan karena …"*.

**Dampak bila salah:** ⛔ pengajuan diarahkan ke jenjang yang salah, atau **dilewatkan dari
penyimpanan** tanpa ada yang menyadarinya.

rujukan: `.IsApproved`, `LetterNo`, `.Quotation.BusinessCode != "40"`, `param.isFOR` — OQ-020

---

#### **P8 — Apa yang sebenarnya dikirim ke layanan "ARASAPAS" di akhir proses?**

Langkah **terakhir** jalur realisasi memanggil sesuatu bernama "HIT SERVICE ARASAPAS".
⛔ **Isi langkah itu tidak ada di dalam ekspor yang kami terima untuk modul ini** — hanya
pembungkusnya. Salinan isinya ada di **dua modul lain**, dan ⛔ **kami sengaja tidak meminjamnya**,
sebab keduanya terdaftar berisi berbeda.

**Konteks:** ini efek keluar terakhir — bila ia mengirim data ke sistem keuangan, kegagalannya
berarti transaksi tidak tercatat di sana.

**Bentuk jawaban yang diharapkan:** pilihan ganda + keterangan — **(a)** kirim ke sistem
faktur/keuangan, sebutkan sistemnya · **(b)** kirim ke sistem lain · **(c)** sudah tidak dipakai.
⭐ Dan: **apa yang terjadi sekarang bila pengiriman itu gagal** — dicoba ulang, atau diabaikan?

**Dampak bila salah:** ⛔ transaksi **tidak sampai** ke sistem tujuan, atau ⛔ **terkirim dua kali**.

rujukan: `Flow\InputRealizationTreatyIn.xml` `Utility2` →
`Activity\serviceInsertArasapas_act.xml` — OQ-025 *(BLOCKER)*

---

#### **P9 — Mengapa modul realisasi treaty masuk menyentuh data treaty KELUAR?**

Modul yang menangani **treaty masuk** ternyata menyentuh tabel **treaty keluar** di **delapan
berkas** — sementara modul yang **namanya** treaty keluar tidak menyentuhnya sama sekali.

**Konteks:** kami perlu tahu apakah ini memang alur bisnis *(misalnya bagian treaty masuk
diteruskan keluar)*, atau peninggalan sejarah.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** memang satu proses, jelaskan ·
**(b)** peninggalan, tidak dipakai lagi · **(c)** dipakai untuk laporan saja.

**Dampak bila salah:** tanggung jawab treaty keluar **dibangun di tempat yang salah**, dan
perubahan di satu sisi merusak sisi lain.

rujukan: `POOLDATA.M_TREATY_OUT`, class `…INT-TREATYOUTDETAIL` — OQ-022

---

#### **P10 — Apakah semua pekerjaan masuk ke antrean BERSAMA, tidak pernah ke antrean pribadi?**

Keenam titik penugasan di alur ini memakai **antrean bersama** — ⛔ **tidak satu pun** menugaskan
ke **kotak masuk pribadi** seseorang.

**Konteks:** ini menentukan bentuk kotak masuk di sistem baru: daftar bersama yang siapa pun dalam
peran itu boleh ambil, atau tugas milik satu orang.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** benar, semua antrean bersama ·
**(b)** seharusnya ada yang pribadi, sebutkan tahap mana · **(c)** dulu pribadi, diubah.

**Dampak bila salah:** pekerjaan **menumpuk tanpa pemilik**, atau sebaliknya **terkunci pada satu
orang** yang sedang tidak ada.

rujukan: keenam Assignment `<pyImplementation>WorkBasket` / `<pyRouteTo>Custom` — OQ-028

---

### Untuk IAM

#### **P11 — Peran seseorang disimpan di kolom NOMOR TELEPON. Di mana peran yang sebenarnya?**

⛔ Satu-satunya penentu peran yang dapat kami baca adalah **kolom nomor telepon** pada data
pengguna. Isinya bukan nomor telepon, melainkan **kode seperti `TREATY1` dan `SPVTREATY1`**.
Percabangan alur diputuskan dari kode itu.

**Konteks:** sistem baru butuh daftar peran yang sebenarnya. Kami tidak dapat menurunkannya dari
kolom yang dipakai untuk hal lain.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh daftar** — peran apa saja yang ada di proses treaty
inward, **siapa saja pemegangnya sekarang**, dan **apa arti tiap kode** yang muncul di kolom itu
*(`TREATY1`, `TREATY2`, `SPVTREATY1`, `SPVTREATY2`)*. ⚠️ Serta: **siapa yang berwenang mengubah
isi kolom itu hari ini.**

**Dampak bila salah:** ⛔ orang yang tidak berwenang **dapat menyetujui**, atau orang yang
berwenang **terkunci di luar**.

rujukan: `OperatorID.pyTelephone` — `When\IsTreaty1.xml`, `When\IsSPVTreaty1.xml` — OQ-027

---

#### **P12 — Ada nama orang tertentu tertulis di dalam aturan. Apakah itu masih benar?**

Di beberapa aturan, wewenang diputuskan dengan **mencocokkan nama orang tertentu** yang tertulis
langsung di dalam aturan itu. ⛔ **Nilainya tidak kami salin ke berkas mana pun.**

**Konteks:** aturan seperti ini berhenti bekerja saat orangnya pindah jabatan atau keluar, dan
tidak seorang pun mendapat pemberitahuan.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** ganti dengan peran, sebutkan perannya ·
**(b)** memang harus orang tertentu, jelaskan alasannya · **(c)** sudah tidak relevan, hapus.

**Dampak bila salah:** ⛔ wewenang tetap menempel pada orang yang sudah **tidak berhak**.

rujukan: `pyWorkPage.pxCreateOperator` di `When\IsSPVCreate.xml`;
`OperatorID.pyUserIdentifier` di 12 berkas — OQ-021

---

#### **P13 — Siapa yang memegang pekerjaan di tiap tahap?**

Alur menyebut **lima nama antrean** — Admin, Kepala Seksi, Ketua Kelompok, Kepala Departemen, dan
Direktur. ⛔ Tetapi **hubungan antara tahap dan antrean tidak tertulis** — penentuan antreannya
dilakukan secara khusus, dan kolom yang biasanya menyimpannya **kosong**.

**Konteks:** tanpa pemetaan ini, kami tidak dapat menyatakan siapa memegang giliran di tahap mana.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh tabel sederhana** — untuk setiap tahap
*("Input Realisasi", "Akseptasi Kepala Treaty", "Akseptasi Kepala Departemen")*, **antrean mana**
yang menerimanya.

**Dampak bila salah:** pekerjaan masuk ke **kotak yang salah** dan tidak ada yang merasa
bertanggung jawab.

rujukan: `<pyWorkBasket>` kosong, `<pyRouteTo>Custom`; nama antrean di `<pyRuleName>` — OQ-024

---

### Untuk pengembang Pega lama

#### **P14 — Satu wadah keterangan dipakai bersama oleh 18 aturan, dan isinya berbeda-beda. Bagaimana cara membacanya?**

Ada sebuah **wadah keterangan bernomor** — slot 1, slot 2, slot 3, dan seterusnya — yang dipakai
**bersama oleh 18 aturan berbeda**. ⛔ **Slot yang sama membawa hal yang berbeda di aturan yang
berbeda**: slot ke-3 membawa **seluruh muatan data dalam bentuk teks** di satu aturan, dan membawa
**posisi/jabatan** di aturan lain.

**Konteks:** kami hendak memetakan setiap slot ke kolom yang tepat. Bila pemetaannya dibuat satu
kali untuk semua, ⛔ **ia akan salah di sebagian besar aturan**.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** memang per-aturan, tidak ada makna
tetap · **(b)** ada dokumen pemetaan, **butuh berkas** · **(c)** ada aturan tak tertulis,
jelaskan. ⭐ Bila (a): **aturan mana yang mengisi slot itu sebelum penyimpanan dijalankan.**

**Dampak bila salah:** ⛔⛔ satu pemetaan yang salah **merusak enam aturan lain** yang memakai slot
yang sama — dan kesalahannya **tidak menimbulkan pesan galat**, hanya data yang salah tempat.

rujukan: halaman `InputData` 15 slot / 18 rule; `InputData.CARI3` di 7 rule —
OQ-059, butir baru #1

---

#### **P15 — Sebuah fungsi yang membentuk seluruh muatan data tidak ada naskahnya. Apa isinya?**

Ada satu fungsi buatan sendiri yang **mengubah data satu kasus menjadi satu teks besar**, dan teks
itulah yang disimpan ke basis data sebagai muatan utama. ⛔ **Naskah fungsi itu tidak ada di dalam
ekspor mana pun.**

**Konteks:** data lama tersimpan dalam bentuk yang dihasilkan fungsi ini. Sistem baru harus dapat
**membacanya kembali**.

**Bentuk jawaban yang diharapkan:** ⭐ **butuh berkas** — naskah fungsi tersebut, **atau** beberapa
**contoh nyata** hasilnya dari basis data *(3–5 contoh sudah cukup untuk menurunkan bentuknya)*.

**Dampak bila salah:** ⛔ **seluruh data lama tidak terbaca** oleh sistem baru.

rujukan: `@ASM.GetPageJSONString()` di `Activity\SaveJsonPolisTreatyIn_Act.xml` langkah 6 — OQ-012

---

#### **P16 — Satu perhitungan membaca tabel internal Pega secara langsung. Apa penggantinya?**

Sebuah perhitungan membaca **tabel kerja internal sistem lama** secara langsung, seolah tabel
biasa. ⛔ **Tabel itu akan hilang** begitu sistem lama ditinggalkan.

**Konteks:** apa pun yang bergantung pada perhitungan itu akan **berhenti bekerja tanpa pesan
galat**.

**Bentuk jawaban yang diharapkan:** pilihan ganda — **(a)** perhitungan ini masih dipakai,
jelaskan untuk apa · **(b)** sudah tidak dipakai · **(c)** tidak tahu. ⭐ Bila (a): **angka apa
yang sebenarnya dicari** dari situ.

**Dampak bila salah:** sebuah angka di layar menjadi **selalu nol** dan tidak ada yang menyadarinya.

rujukan: `DATAPEGA.PC_ASM_FW_GCNMFW_WORK` di `RDBList\GetCountClaim.xml` — butir baru #5

---

#### **P17 — Aturan yang hanya hidup di layar belum kami baca sama sekali. Apa yang ada di sana?**

⛔ **Kami belum membuka satu pun dari 31 berkas layar** modul ini — dan **kelima berkas terbesar
seluruh modul ada di antaranya**. Artinya pemeriksaan yang dilakukan **di layar** — medan wajib,
tombol yang hanya muncul untuk jabatan tertentu, perhitungan yang berjalan saat mengetik — ⛔ **belum
terbaca.**

**Konteks:** kami menyatakan ini terang-terangan agar tidak ada yang menyangka ronde ini sudah
meliput modulnya. ⭐ Ronde berikutnya perlu arahan **layar mana yang paling penting**.

**Bentuk jawaban yang diharapkan:** ⭐ **urutan prioritas** — dari daftar layar utama
*(input polis treaty, rincian polis, rincian non-proporsional, layar Kepala Departemen)*, mana yang
**paling banyak memuat aturan bisnis**, bukan hanya tampilan.

**Dampak bila salah:** ronde berikutnya membaca **8 MB layar** dan menemukan sebagian besarnya
tata letak, sementara aturan yang penting tetap terlewat.

rujukan: 25 `Section` + 6 `Harness`, **nol dibuka**; lima terbesar semuanya `Section` — butir baru #6

---

## Bab E — Lampiran bukti isolasi

| Yang dibuktikan | Keadaan | Perintah audit |
| --- | --- | --- |
| ⛔ folder korpus selain `NB Treaty In` | ✅ **NOL dibuka** | seluruh pembacaan korpus bergantung pada `ROOT = D:\XML\RNM_BRD\NB Treaty In`; ⛔ **`EDM Treaty In` TIDAK dibuka** walau satu bounded context |
| ⛔ `D:\XML\nusantara-re\` | ✅ **NOL disentuh** | tidak pernah muncul sebagai argumen perintah apa pun |
| ⚠️ isi `.scratch\` konteks lain | ⚠️ **lihat catatan jujur di bawah** | `ls -la .scratch` — mtime tiap folder |
| korpus `NB Treaty In` tidak berubah | ✅ **md5 `62a3735ebb2eabb788c8cd8abdda78ca`** | `find . -name '*.xml' \| sort \| xargs cat \| md5sum` |
| berkas dibuat | ✅ **tepat SATU** | `ls .scratch\nb-treaty-in\` |
| spec / tiket | ✅ **NOL** | ⛔ tidak dibuat, sesuai larangan brief |
| kode · DDL · `CREATE TABLE` · daftar kolom usulan | ✅ **NOL** | ⚠️ ungkapan terlarangnya **disebut** di kalimat yang menyatakannya nol — disebut, bukan dipakai |
| butir terbuka yang saya tutup | ✅ ⛔ **NOL** | Bab C — setiap butir tetap berpemilik |
| nilai nama orang | ✅ **NOL disalin** | P12 menyebut **pola dan jumlahnya**, ⛔ bukan namanya |

### ⚠️ Satu baris yang TIDAK bersih, dan saya menuliskannya apa adanya

⛔ **Invarian `.scratch` konteks lain = NOL dibaca tidak dapat saya klaim bersih.** Sebabnya:
⭐ **ronde ini menyambung sesi yang sudah panjang**, dan ketika sesi disambung, harness
**menyuntikkan ulang dua bacaan lama** ke dalam konteks saya — keduanya `urutan-tiket.md`
dari dua folder `.scratch` konteks lain. ⛔ **Bukan saya yang memanggilnya di ronde ini**, tetapi
isinya **ada di depan mata saya**, dan menyembunyikan itu akan membuat invarian ini bohong.

⭐ **Bukti bahwa isinya tidak merembes ke berkas ini:** ⛔ **nol nama tabel, nol nama kolom, dan nol
pola rancangan** dari konteks lain dipakai di sini. Seluruh **31 objek Oracle** di Bab B berasal
dari sisiran SQL korpus `NB Treaty In` sendiri. ⭐ Rujukan lintas-konteks yang **memang** saya
pakai hanyalah **register bersama** yang brief perintahkan dibaca —
`discovery/`, `CONTEXT.md`, `docs/adr/` — ⛔ **bukan `.scratch` modul lain.**

⚠️ **Dan satu lagi:** Bab A menyebut tiga procedure jalur Life yang badannya sudah diserahkan
*(`PEGA_M_LIFE_PREMIUM_SUMMARY`, `INSERTJSONPOLISLIFE`, `INSERTJSONOFFERLIFE`)*. ⭐ Itu dibaca dari
**`discovery/open-questions.md`**, register yang brief perintahkan dibaca, ⛔ **bukan dari korpus
modul lain.** ⭐ Dan penyebutannya justru untuk **membuktikan bahwa ketiganya BUKAN** procedure
yang dipanggil modul ini.

---

## Penutup — rekap

| | Jumlah |
| --- | ---: |
| butir `[terbuka]` yang **lahir** di ronde ini | ⭐ **7** *(belum bernomor register; nomor bebas berikutnya OQ-073)* |
| butir terdaftar yang **dikonfirmasi masih berlaku** | **17** |
| ⛔ butir yang **saya tutup** | ⛔ **0** — ⭐ **sesuai larangan; penutupan hanya oleh pemilik peran** |
| empat batas pengetahuan yang diminta dikonfirmasi | ✅ **4 dari 4 masih berlaku**, ⭐ **3 diperdalam** |
| ⛔ galat instrumen saya sendiri yang tertangkap dan diralat | ⛔ **2** |

### Apakah frontier masih terbuka?

⛔⛔ **MASIH TERBUKA, dan lebar.** ⭐ Tiga sebab, disebut berdasarkan besarnya:

1. ⛔ **31 dari 278 berkas tidak dibuka sama sekali** — seluruh `Section` dan seluruh `Harness`,
   termasuk kelima berkas terbesar modul. ⭐ Aturan yang hidup di layar **belum tersentuh.**
2. ⛔ **92 Activity dan 75 When belum dibaca langkah demi langkah** — hanya disisir pola.
   ⚠️ **75 When** adalah jumlah yang tidak biasa: ⭐ modul ini memutuskan banyak hal lewat
   penggolong, dan **isi penggolongnya belum dibaca.**
3. ⛔ **Rantai hitung uang belum disentuh sama sekali** — ⭐ itulah sebabnya `[Finance]` punya
   **nol butir** di Bab C, dan ⚠️ **nol itu bukan kabar baik, ia lubang.**

### Syarat agar ronde 2 produktif

| # | Syarat | Sifat |
| --- | --- | --- |
| ⭐ **1** | **P17 dijawab** — mana dari empat layar besar yang memuat aturan bisnis, bukan tata letak | ⭐ **tanpa ini ronde 2 membakar 8 MB untuk menemukan tata letak** |
| ⭐ **2** | **P14 dijawab** — cara membaca wadah slot generik | ⛔ **menahan**: tanpa ini tiap temuan SQL berhenti di "slot ini entah apa" |
| ⭐ **3** | **P1 dijawab** — naskah tiga procedure | ⛔ **menahan Bab B**: 29 dari 31 objek hanya dibaca; **penyimpanan sejatinya tak terlihat** |
| **4** | **P6 dijawab** — `isApproved` mana yang berlaku | ⛔ menahan pembacaan alur: gerbang pertama |
| **5** | Izin membuka **`EDM Treaty In`** | ⚠️ ⛔ **keputusan work owner**, bukan saya — lihat catatan di bawah |

### ⚠️ Satu modul yang saya TIDAK buka, dan kenapa saya tetap menyebutnya

⛔ **`EDM Treaty In` tidak saya buka**, sesuai larangan brief. ⭐ Tetapi brief juga meminta:
bila ada modul lain yang relevan, **tulis sebagai pertanyaan, sebutkan modul apa dan kenapa, lalu
berhenti.** Jadi:

⭐ **Modul: `EDM Treaty In`.** ⛔ **Kenapa: ia memuat satu-satunya implementasi langkah terakhir
jalur realisasi modul ini** *(`SERVICEINSERTARASAPAS_ACT`, ±232 KB)*, dan **selama ia tidak dibuka,
efek keluar terakhir NB Treaty In tidak dapat dinyatakan sama sekali.** ⚠️ Catatan lama sudah
merekam bahwa dua varian yang ada **berbeda hanya pada 3 tag metadata**, ⛔ tetapi blocker
prosedural tetap berdiri.

⛔ **Saya berhenti di sini.** Kapan modul lain boleh disentuh adalah **keputusan work owner**.
