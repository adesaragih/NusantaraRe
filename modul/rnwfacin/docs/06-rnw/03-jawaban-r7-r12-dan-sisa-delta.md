# Discovery Renewal — Putaran 3: jawaban R-7…R-12 dan sisa berkas delta

> **Sumber:** `D:\migrasi\RNM\RNW Fac In\` dan `NB FacIn\`, diurai `[System.Xml.XmlDocument]`.
> `Endorsment Fac In\` dan korpus Treaty **tidak dibaca** (K-030, K-005). Label `CLAUDE.md` §3.
>
> Metode kunci putaran ini: **periksa `<pxRuleObjClass>` dan tag pembawa, bukan nama** — pelajaran
> yang dipaksakan oleh koreksi R-1 pada putaran 2.

---

## Ringkasan temuan

### ⛔ T-9 — **`SumTreatyCapacity_Act` adalah rujukan menggantung sungguhan di RNW. Menyentuh K-006.**

`[terverifikasi]` Dua activity yang **ada** di RNW memanggilnya sebagai **langkah activity sungguhan**:

```
RNW\Activity\CountASMShareTotal_ACT.xml  L444
    <pyStepsActivityName>Call SumTreatyCapacity_Act</pyStepsActivityName>
RNW\Activity\SetValidateDate_Act.xml     L1129
    <pyStepsActivityName>Call SumTreatyCapacity_Act</pyStepsActivityName>
```

`[terverifikasi]` Keberadaan berkasnya:

| Rule | NB | RNW |
| --- | :-: | :-: |
| `SumTreatyCapacity_Act` | **ada** | ⛔ **TIDAK ADA** |
| `CountTreatyCapacity_Act` (yang dipanggilnya) | tidak ada | tidak ada — tiba lewat `DDL\` (amandemen K-006) |

**Rantai cabang 4 K-006 ("kapasitas treaty") putus SATU HOP LEBIH AWAL di renewal:**

```
NB :  CountASMShareTotal_ACT / SetValidateDate_Act → SumTreatyCapacity_Act (ADA) → CountTreatyCapacity_Act (tiba lewat DDL\)
RNW:  CountASMShareTotal_ACT / SetValidateDate_Act → SumTreatyCapacity_Act (HILANG) ✗
```

⚠️ **Ini bukan vonis.** Aturan `CLAUDE.md` §4.5 berlaku penuh: **"hilang dari ekspor" ≠ "usang"**.
Dan di sini bobotnya lebih besar dari biasanya, karena T-1 membuktikan kedua folder diekspor dari
keadaan yang sama — satu berkas yang hadir di satu sisi saja menjadi lebih mencurigakan, bukan kurang.

📌 **DITANDAI UNTUK KEPUTUSAN WORK OWNER** — menyentuh **K-006**.

### ⚠️ T-10 — **Koreksi kedua atas laporan saya sendiri: RNW TIDAK mengisi celah K-004**

Putaran 2 menulis bahwa `serviceInsertArasapasRNW_act` adalah *"implementasi kelas `Work` yang ekspor
NB tidak punya"*. **Rumusan itu menyesatkan.**

`[terverifikasi]` Pemeriksaan ulang seluruh activity bernuansa *arasapas*:

| Activity | NB | RNW | Kelas | Langkah |
| --- | :-: | :-: | --- | ---: |
| `serviceInsertArasapas_act` | ada | ada | `Data-PolicyTreatyIn` | **1** (stub) |
| `serviceInsertArasapasEDM_act` | **ada** | **ada** | **`ASM-FW-GISFW-Work`** | **20** |
| `serviceInsertArasapasRNW_act` | tidak | **ada** | **`ASM-FW-GISFW-Work`** | **10** |

**NB sudah punya activity konversi kelas `Work` sejak awal** — `serviceInsertArasapasEDM_act`,
20 langkah. Jadi bukan benar bahwa "kelas `Work` tidak ada di NB".

`[terverifikasi]` Yang benar-benar hilang lebih spesifik. Stub NB menargetkan kelas `Work`:

```
NB\Activity\serviceInsertArasapas_act.xml
  L230  <pyStepsParentClass>ASM-FW-GISFW-Data-PolicyTreatyIn</pyStepsParentClass>
  L236  <pyStepsClassName>ASM-FW-GISFW-Work</pyStepsClassName>       ← sasaran panggilan
  L221  <pyStepsActivityName>Call serviceInsertArasapas_act</pyStepsActivityName>
```

Yang hilang adalah **`serviceInsertArasapas_act` kelas `Work` — nama tanpa sufiks** — dan ia **tetap
tidak ada di NB maupun RNW**.

**Kesimpulan R-11: celah K-004 TIDAK terisi.** `serviceInsertArasapasRNW_act` adalah **saudara**
dengan nama berbeda, bukan rule yang dicari.

**Yang tetap bernilai:** kini ada **dua saudara terbaca** (EDM 20 langkah, RNW 10 langkah) yang
memperlihatkan **polanya**, dan itu berguna untuk merancang jalur konversi:

```
1 Call serviceInsertArasapas_act   6 Connect-REST
2 Property-Set                     7 Property-Set
3 RDB-List                         8 Call InsertLogServiceProd
4 Property-Set                     9 Page-Remove
5 Call ASM-FW-GISFW-Int-M_LINK_SERVICE.GetLinkService    10 RDB-List
```

`[terverifikasi]` Pola **endpoint diambil dari `M_LINK_SERVICE` → `Connect-REST` → catat log layanan**
menguatkan `CLAUDE.md` §4.4 langsung dari korpus.

📌 **DITANDAI UNTUK KEPUTUSAN WORK OWNER** — apakah pola saudara ini boleh dipakai merancang jalur
konversi NB meski bukan rule aslinya. **Tidak saya putuskan.**

### ⚠️ T-11 — `Section\EditMarketing` mengikat properti **`.Password`**

`[terverifikasi]` Empat properti terikat: `.ID` · **`.Password`** · `.QuotationData.MOID` ·
`.pyTemplateButton`, pada kelas `ASM-FW-GISFW-Data-OfferFacIn`.

Sebuah layar penyuntingan marketing officer yang meminta **ID dan kata sandi**. `[pertanyaan terbuka]`
Ini mekanisme otorisasi tertanam di lapisan data; apa yang divalidasi dan terhadap apa **belum
ditelusuri**. Perlu tinjauan keamanan sebelum diport (`CLAUDE.md` §6).

⛔ Tidak ada nilai kredensial yang disalin ke dokumen ini — hanya nama propertinya.

---

## R-7 · Tipe rujukan 42 nama — **TERJAWAB untuk dua yang paling berdampak**

Uji: baca **tag pembawa** dan `<pxRuleObjClass>`, bukan nama.

| Nama | Muncul di | Tag pembawa | Vonis |
| --- | --- | --- | --- |
| `GetLimitAkseptasi_ActFlow` | `Activity\GetLimitAkseptasi_JUW_UW` L61 | **`<pzOriginalInstanceKey>`** | ✅ **BUKAN rujukan** — metadata asal-usul ekspor |
| `SumTreatyCapacity_Act` | `CountASMShareTotal_ACT` L444 · `SetValidateDate_Act` L1129 | **`<pyStepsActivityName>`** | ⛔ **RUJUKAN SUNGGUHAN** — lihat T-9 |

Isi baris metadata itu `[terverifikasi]`:

```
<pzOriginalInstanceKey>RULE-OBJ-ACTIVITY ASM-FW-GISFW-WORK GETLIMITAKSEPTASI_ACTFLOW #20250703T063546.429 GMT</pzOriginalInstanceKey>
```

📌 Ia mencatat bahwa `GetLimitAkseptasi_JUW_UW` **berasal dari salinan** `GetLimitAkseptasi_ActFlow` —
jejak *save-as*, bukan panggilan. Dan ia bertag `pz*`, yaitu **persis metadata yang dibuang normalisasi
putaran 1** — karena itu ia tidak pernah memengaruhi temuan identik byte-per-byte.

**Sisa 40 nama belum diperiksa satu per satu** — tetap `[pertanyaan terbuka]`. Dua yang paling
berdampak sudah selesai, dan keduanya jatuh ke sisi yang berlawanan: satu tabrakan, satu sungguhan.
Itu sendiri pelajaran — **tidak ada aturan umum yang bisa menggantikan pemeriksaan per nama.**

---

## R-8 · Gerbang masuk renewal — **TERJAWAB: di luar korpus**

`[terverifikasi]` `Flow\InputRenewalFacultativeIn` **tidak dirujuk berkas mana pun** selain dirinya
sendiri, dan `IsOfferFacIn` **dirujuk nol kali** di RNW.

**Penentu "kapan sebuah kasus masuk alur renewal" berada di konfigurasi work type / portal Pega yang
tidak ikut terekspor.** Ini bukan celah pembacaan — memang bukan rule.

📌 **Konsekuensi: sistem baru harus MENETAPKAN gerbang masuk renewal secara eksplisit.** Tidak ada
perilaku terekam yang dapat direproduksi, sehingga ini **keputusan bisnis, bukan porting**.

---

## R-9 / R-5 · Di mana "hitung ulang" renewal — **TERJAWAB: tidak ada hitung ulang**

`[terverifikasi]` Seluruh Activity RNW yang menggerbangi `StatusBusiness = 2` — hanya **empat**:

| Activity | Kondisi | Yang digerbangi |
| --- | --- | --- |
| `GetDateValidity_ACT` L1109 | `StatusBusiness==2` | masa berlaku tanggal |
| `SetValidateDate_PostAct` L3148 | `StatusBusiness==2 \|\| ==3` | validasi tanggal |
| `CheckSpreadingProtect_ACT` L3621 | `=="2" \|\| =="3"` | proteksi spreading |
| `InputDtlPayment_PreAct` L5107 | `==2 \|\| ProposalPosition=="1"` | input pembayaran |

**Tidak satu pun perhitungan premi.** Keempatnya berkas **bersama** yang identik byte-per-byte dengan
NB.

**Kesimpulan:** renewal **tidak punya mesin hitung ulang**. Ia menjalankan **perhitungan NB yang sama
persis, tanpa perubahan**, atas data kasus baru. Kaitan ke polis lama adalah **rujukan**
(`OldPolicyNo`), bukan pembawaan data.

📌 Ini **menjelaskan K-015 dari sisi mekanisme**: renewal dinilai atas nilai pertanggungan **penuh**
bukan karena ada kebijakan khusus, melainkan karena **tidak ada perhitungan selisih sama sekali** —
yang ada hanya perhitungan baru di atas kasus baru. Keputusan bisnisnya (K-015) dan mekanismenya
sejalan.

⚠️ **Ketidakkonsistenan yang dicatat sebagai kandidat perbaikan:** `GetDateValidity_ACT` membandingkan
`StatusBusiness==2` sebagai **angka**, sementara `CheckSpreadingProtect_ACT` membandingkan `=="2"`
sebagai **string**. Pola yang sama sudah tercatat di `CLAUDE.md` §4.1. Direproduksi apa adanya.

---

## R-10 · `PROSESCOPY` — **TERJAWAB apa adanya, artinya tidak disimpulkan**

`[terverifikasi]` Isi penuh `RDBList\CariBusinessGID` `<pySaveSQL>`:

```sql
DECLARE
  vTHN_TREATY VARCHAR2(10);  vTOP_ID VARCHAR2(10);
  vIDTreatyYear VARCHAR2(10); errmsg VARCHAR2(4000);
BEGIN
  vTHN_TREATY   := '2007';
  vTOP_ID       := '10018';
  vIDTreatyYear := '1000134';
  POOLDATA.PROSESCOPY(vTHN_TREATY, vTOP_ID, vIDTreatyYear, errmsg);
  dbms_output.put_line(errmsg);
END;
```

Empat sifat terbaca, `[terverifikasi]` semuanya:

1. Seluruh argumennya **literal keras** — tahun `'2007'`, dua ID tetap. Tidak ada parameter masuk.
2. Memanggil `POOLDATA.PROSESCOPY`, yang **dikonfirmasi DBA tidak ada di basis data**.
3. Keluarannya ke `dbms_output` — kanal yang tidak dibaca aplikasi.
4. **Tidak berhubungan dengan `<pyBrowseSQL>` di rule yang sama**, yang mencari kelompok bisnis.

`[dugaan]` Bentuknya menyerupai **potongan uji coba yang tertinggal**, bukan logika produksi. ⛔ **Tidak
disimpulkan sebagai kode mati** — pemanggilnya (`GetBusinessGroup_Act`) memakai `RDB-List`, dan apakah
jalur Save pernah tereksekusi **belum dibuktikan**. `[pertanyaan terbuka]`, diport apa adanya
(`CLAUDE.md` §1).

---

## R-12 · `Work-Renewal` — **TERJAWAB: irisan pelaporan, bukan kelas kerja**

`[terverifikasi]`

| Pemakai | Kelas |
| --- | --- |
| `Flow\InputRenewalFacultativeIn` → `pyWorkClass` | **`ASM-FW-GISFW-Work`** — sama dengan NB |
| `ReportDefinition\RenewalList_RD` → `pyClassName` | `ASM-FW-GISFW-Work-Renewal` |

`Work-Renewal` **tidak dipakai** oleh flow, activity, section, maupun rule lain mana pun — hanya satu
ReportDefinition. **Kasus renewal dibuat pada kelas kerja yang sama dengan NB.**

`[pertanyaan terbuka]` Definisi kelas `Work-Renewal` sendiri tidak ada di korpus, sehingga apakah ia
subkelas nyata dengan instance tersendiri tidak dapat dipastikan dari sini.

---

## Sisa berkas delta — terpetakan

| Berkas | Temuan |
| --- | --- |
| `FlowAction\EditMarketing` | kelas `Data-OfferFacIn`; `pySectionReference` = `EditMarketing` |
| `Section\EditMarketing` | 4 properti: `.ID` · **`.Password`** · `.QuotationData.MOID` · template → **T-11** |
| `Harness\ChooseInsured` | kelas `Data-Quotation`; hanya 3 properti (`.Name`, `.pyID`, template) — **sangat tipis** |
| `Section\ChooseInsuredDtl` | 3 properti yang sama — pemilih daftar, bukan penarik data polis lama |

**Kedua puluh berkas delta kini terpetakan.**

### `Struktur_InputRenewalFacultativeIn.xlsx` — sah, belum dibaca isinya

`[terverifikasi]` Berkas **OOXML sah** (byte awal `50 4B 03 04`), **7,02 MB**, satu sheet bernama
`Struktur`, **81.379 baris × 28 kolom**.

Baris pertama hanya memuat `No` di kolom pertama — judul kolom kemungkinan tidak di baris 1.
**Isinya belum dibaca**: volume 81 ribu baris memerlukan putaran tersendiri. Format **tidak rusak**.

---

## Status seluruh pertanyaan

| # | Status |
| --- | --- |
| R-1 | ✅ ditutup — properti, bukan rujukan menggantung |
| R-2 | ✅ ditutup → menjadi R-8 |
| R-3 | ✅ ditutup → dikonfirmasi ulang R-12 |
| R-4 | ✅ ditutup — `BusinessCode` → `business` → `BusinessGroupID` |
| R-5 | ✅ ditutup → dijawab bersama R-9 |
| R-6 | 🟡 132 bersih; **40 nama** belum diperiksa tipenya |
| R-7 | 🟡 dua terpenting selesai; 40 sisa terbuka |
| R-8 | ✅ ditutup — **keputusan bisnis**, bukan porting |
| R-9 | ✅ ditutup — tidak ada hitung ulang |
| R-10 | 🟡 isi terbaca; **arti tidak disimpulkan** |
| R-11 | ✅ ditutup — **celah K-004 tidak terisi** |
| R-12 | ✅ ditutup — irisan pelaporan |

## Pertanyaan terbuka baru

| # | Pertanyaan |
| ---: | --- |
| R-13 | ⛔ `SumTreatyCapacity_Act` menggantung di RNW — apakah cabang kapasitas treaty memang tidak berlaku di renewal, atau ekspor RNW tidak lengkap? **Menyentuh K-006** |
| R-14 | `Section\EditMarketing` meminta **ID + kata sandi** — apa yang divalidasi, terhadap apa? Perlu tinjauan keamanan |
| R-15 | Bolehkah pola `serviceInsertArasapasEDM/RNW_act` dipakai merancang jalur konversi NB, meski bukan rule aslinya? **Menyentuh K-004** |
| R-16 | Isi `Struktur_InputRenewalFacultativeIn.xlsx` (81.379 baris) — putaran tersendiri |
| R-17 | 40 nama sisa R-6 — tipe rujukan per nama |

---

*Tanpa nama orang, tanpa alamat email, tanpa data pelanggan, tanpa nilai kredensial.*
