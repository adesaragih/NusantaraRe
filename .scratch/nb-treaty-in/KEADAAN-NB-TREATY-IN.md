# Keadaan NB Treaty In — 22 September 2026

## Apa berkas ini

Berkas ini menyatakan **keadaan modul NB Treaty In yang berlaku sekarang**, sesudah empat ronde
grilling, 47 jawaban work owner, dan satu perbaikan berkas korpus.

Ia ada karena berkas `grilling-ronde-1..4.md` **tersegel** — isinya tidak disunting setelah ronde
selesai. Sebagian angkanya kini basi, dan sebagian temuannya keliru. Berkas ini memuat yang benar.

### Urutan kewenangan

Bila dua sumber bertentangan, yang di atas menang:

| # | Sumber | Sifat |
| ---: | --- | --- |
| 1 | `PERTANYAAN-untuk-*.md` — lembar jawaban per pemilik | keputusan work owner, mengikat |
| 2 | **berkas ini** | keadaan terukur, dapat diuji ulang |
| 3 | `grilling-ronde-1..4.md` | latar; **tidak berlaku** di mana bertentangan dengan 1 atau 2 |

---

## 1 · Angka yang berlaku

Seluruhnya diukur ulang 22 September 2026 sesudah `Protection_Act` diperbaiki.

| Hal | Berlaku | Pernah tertulis | Di mana |
| --- | ---: | ---: | --- |
| Berkas dalam modul | **278** | 278 | — |
| Ukuran modul | **35.492.317 B** | 35.678.284 B | ronde 2, 4 |
| Langkah `Property-Set` | **938** | 942 | ronde 2, 3, 4 |
| Langkah `Page-Clear-Messages` | **16** | 17 | ronde 3 |
| Baris `DecisionTable\BusinessType_DeT` | **36** | 9 | ronde 4 |

### Sebaran berkas

| Tipe | Berkas | Byte |
| --- | ---: | ---: |
| `Activity` | 92 | 14.690.375 |
| `When` | 75 | 1.899.674 |
| `RDBList` | 41 | 268.811 |
| `Section` | 25 | 14.229.383 |
| `ReportDefinition` | 15 | 1.084.901 |
| `DataTransform` | 12 | 286.884 |
| `FlowAction` | 9 | 284.415 |
| `Harness` | 6 | 2.406.451 |
| `DecisionTable` | 2 | 53.678 |
| `Flow` | 1 | 287.745 |

### Sebaran langkah Activity — 1.342 seluruhnya

`Property-Set` 938 · `Page-Set-Messages` 97 · `Call` 74 · `RDB-List` 66 · `Page-New` 33 ·
`Page-Remove` 24 · `Property-Remove` 24 · `Property-Set-Messages` 20 · `Page-Clear-Messages` 16

---

## 2 · Lingkup — folder ini bukan daftar pekerjaan Treaty

**Folder `NB Treaty In` adalah *dependency closure*.** Ia memuat setiap aturan yang disentuh modul
ini, termasuk milik modul lain, karena pewarisan kelas Pega.

Dari 92 `Activity`:

| | Rule | Langkah `Property-Set` |
| --- | ---: | ---: |
| **Terjangkau** dari titik masuk nyata | 64 | 377 |
| **Yatim** — tidak dipanggil siapa pun | 28 | 561 |

Titik masuk nyata = dirujuk dari `Flow`, `Section`, `Harness`, `FlowAction`, atau `DataTransform`,
lalu diikuti pemanggilannya sampai tidak ada aturan baru yang tercapai.

### 28 aturan yatim — tidak dimigrasi

Seluruhnya keluarga spreading/premi milik Fac In:

| Langkah | Aturan | | Langkah | Aturan |
| ---: | --- | --- | ---: | --- |
| 100 | `SumTSIPremiSpreadedRNM_FIRE_Act` | | 17 | `GetTreatyName` |
| 62 | `SumTSIPremiSpreadedRNM_ANEKA_Act` | | 13 | `SumTSIPremiSpreadedRNM_TRAVEL_Act` |
| 49 | `ProtectFIREMBUPA_Act` | | 13 | `ProtectPremiPolicy_Act` |
| 35 | `CheckSpreadingProtectAnekaGolf_ACT` | | 13 | `CheckSpreadingProtect_ACT` |
| 33 | `SumTSIPremiSpreadedRNM_Act` | | 12 | `SumTSIPremiSpreadedRNM_PA_Act` |
| 32 | `CheckSpreadingProtectFire_ACT` | | 11 | `cekSpreadingFactIn` |
| 25 | `CheckSpreadingProtectMCargoMBU_ACT` | | 10 | `SetFlagOccupation_ACT` |
| 25 | `CheckSpreadingProtectPATravel_ACT` | | 8 | `ProtectShareCedant_Act` |
| 21 | `SumTSIPremiSpreadedRNM_MBU_Act` | | 4 | `ProtectCurrencyTSI_Act` |
| 21 | `ProtectCoverage_Act` | | 4 | `GetDateValidity_ACT` |
| 21 | `SumTSIPremiSpreadedRNM_GOLF_Act` | | 4 | `CheckRISLIP` |
| 17 | `SumTSIPremiSpreadedRNM_MARINECARGO_Act` | | 3 | `CheckRISLIP_EDM` |

Ditambah `ProtectShipData_Act` (3), `ProtectRenewal_Act` (2), `ProtectCedingCo` (2),
`GetCurrencyMaster` (1).

### Sebabnya — satu berkas yang salah dimasukkan

Sampai 22 September 2026 pukul 15.45, folder ini memuat `Protection_Act` **versi Fac**:
kelas `ASM-FW-GISFW-Work`, 343.271 B, 37 langkah, 10 pemanggilan — identik byte dengan versi di
NB FacIn dan RNW Fac In.

Versi Treaty yang benar: kelas `ASM-FW-GISFW-Data-PolicyTreatyIn`, 157.304 B, **16 langkah**,
satu pemanggilan (`ProtectionNonProp_Act`), dan **tidak memanggil satu pun** aturan pada rantai
spreading.

Penggantian berkas itu memutus rantai Fac di akarnya dan menjadikan 28 aturan di atas yatim.

Penelusuran berkas lain dengan pola yang sama — modul Treaty memakai varian berkelas berbagi
padahal varian berkelas Treaty ada di modul Treaty lain — menemukan **nol kasus lain**.

---

## 3 · Sumber data — JSON ditinggalkan

Pembacaan **dan** penulisan tidak lagi memakai kolom dokumen `JSONDATA`. Sumbernya pindah ke view
relasional yang sudah ada:

```
POOLDATA.TREATYINDETAILJOINEDM
   = pooldata.treatyindetail  UNION ALL  pooldata.treatyindetailedm
```

**39 kolom.** Uji kecukupan: `ReportDefinition\BrowseTreatyInDetail.xml` merujuk 33 medan;
ke-33-nya tersedia di view; **nol yang hilang**. Dari 31 berkas layar, 29 kolom view dipakai
langsung.

Enam kolom tidak dipakai laporan: `BROKERAGE` `COMMENCEMENT` `TERMINATION` `RIOGR` `RIONR`
`RNM_SHARE`. Lima di antaranya tetap dipakai di tempat lain; **`RNM_SHARE` nol dipakai di seluruh
modul**.

Enam aturan pembongkar JSON **tidak dimigrasi**: `FetchMasterTreatyIn` · `SetTreatyIn_Act` ·
`InputPolicyTreatyInDetail_preACT` · `InputPolicyTreatyInDetail_NonProp` ·
`InputPolicyTreatyOutDetail_preACT` · `InputPolicyTreatyOutDetail_NonProp`.

---

## 4 · Keputusan yang mengikat implementasi

| # | Ketetapan | Butir |
| ---: | --- | --- |
| 1 | `IsApproved`: `0` = ditolak, selain itu = disetujui. Aturan hidup ada di `DecisionTable`, bukan `When`. Dibandingkan sebagai **teks** | P24, P6 |
| 2 | Admin menolak → berkas diselesaikan sebagai ditolak. Atasan menolak → kembali ke admin | P24 |
| 3 | Nilai uang berpresisi penuh; pembulatan hanya di titik penyajian, tidak pernah di repository | P29 |
| 4 | `DEDUCTION1` `DEDUCTION2` `BROKERAGE` `RNM_SHARE` adalah **persentase**, bukan uang. `12.5` berarti 12,5 persen | P29 |
| 5 | `COMMENCEMENT` `TERMINATION` dibaca apa adanya; **P32 tetap berlaku** untuk penulisan tanggal baru | P29 |
| 6 | Setiap query menulis skema `POOLDATA.` eksplisit | P3 |
| 7 | Seluruh urutan penyimpanan dibungkus **satu transaksi** — lebih ketat dari sistem lama | P2 |
| 8 | `OPERATORID` diisi dari identitas akses login; `PIC` dari nama tampilan | P4, P33 |
| 9 | Pencarian berdasarkan **nomor urut** antrean tidak dimigrasi; diganti pemeriksaan keanggotaan | P25 |
| 10 | Bila keterangan rule `When` berbeda dari syarat yang dijalankan, **yang dijalankan benar** | P23 |
| 11 | 38 medan tetap terkunci permanen; 27 medan wajib tetap wajib, berbeda menurut tingkat | P45, P47 |
| 12 | `TypeTax` dan potongan 2,2 % ditiru apa adanya, termasuk ketidakseragaman presisi 4 lawan 8 | P46 |

---

## 5 · Koreksi atas ronde 1–4

| Bunyi lama | Yang benar | Asal |
| --- | --- | --- |
| *"942 langkah Property-Set"* | 938 di folder; **268** yang diminta sesudah penyaringan | ronde 2, 3, 4 |
| *"17 langkah Page-Clear-Messages"* | **16** | ronde 3 |
| *"lapisan layar nol dibuka"* | ronde 2 sudah menyisir ke-31 berkas | ronde 3 |
| *"naskah SQL tidak ikut terekspor"* | ada di `<pyBrowseSQL>` pada `RDBList` — 41 di modul ini, **1.151** di korpus | sesi asisten |
| *"baris tabel keputusan tidak terkirim"* | ada di `pyCondition`, `pyOrConditions`, `pyResults` | ronde 2 |
| *"BusinessType_DeT tinggal 9 baris dari 34"* | **36 baris utuh**; `pyRowNum` berbasis nol menandai baris ber-daftar-OR, bukan nomor baris | ronde 4 |
| *"160 setelan wajib tidak menempel medan mana pun"* | seluruhnya dapat dipasangkan — **27 medan wajib di 6 layar** | ronde 4 |
| *"seluruh 38 medan terkunci ada di layar Kepala Departemen"* | **36 di sana, 2 di layar biasa** | ronde 2 |
| ⭐⭐ *"isi 268 langkah penetapan nilai tidak ikut terekspor"* | ⭐ **ADA di ekspor**, di tag `PropertiesName`/`PropertiesValue` — **tanpa awalan `py`**. **938 dari 939** langkah `Property-Set` membawa isinya sendiri; **2.481** pasangan nama=nilai terisi; **707** pada ke-51 aturan yang diminta; **nol** tanda terpotong | ronde 2, 3, 4 |
| ⭐ *"layar Kepala Departemen tidak dapat disimpan sama sekali"* | ⭐ **dapat disimpan** — dari **34** medan wajib-sekaligus-terkunci, **28** punya langkah pengisi; 3 sisanya *(`ExcessLoss`, `OutstandingClaim`, `SalvageValue`)* **diketik underwriter** di `DetailPolicyTreatyIn`/`GeneralPolicyTreatyIn` | ronde 4 |
| ⚠️ *"938 langkah Property-Set"* | **939** seluruhnya; **938** yang berisi. Angka 938 adalah cacah yang **berisi**, bukan cacah seluruhnya | `VERIFIKASI-P18.md` |

⭐⭐ **Tiga pertanyaan ditarik karena premisnya keliru: **P41**, **P48**, dan — 2026-09-22 — **P18**.**

> ⚠️ **P18 adalah yang termahal dari ketiganya.** Ia disebut penahan terberat seluruh proyek
> selama tiga ronde, dan ia **tidak pernah ada**. ⭐ **Sebabnya dua huruf: `py`.** Tim migrasi
> memeriksa `pyPropRef`, `pyPropertiesName`, `pyPropertiesValue` — **dengan** awalan — menemukan
> nol, dan mempercayainya. Tag yang sebenarnya berisi tidak memakai awalan itu.
> ⚠️ Di `DataTransform` justru **kebalikannya**. Rinciannya di `VERIFIKASI-P18.md`.

---

## 6 · Yang masih terbuka

| # | Menahan apa | Pemilik |
| --- | --- | --- |
| ✅ ~~**P1**~~ | ~~penulisan penyimpanan data — naskah tiga stored procedure~~ | ✅ **TERJAWAB 2026-09-22** — **empat** naskah diterima |

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — tabel ini semula memuat satu baris lagi:
> > *"| **P18** | perhitungan uang **dan** alur persetujuan Kepala Departemen — 268 langkah |
> > pemilik export Pega |"*
> >
> > *"Surat keduanya sudah siap: `SURAT-P1-KE-DBA.md`, `SURAT-P18-KE-PEMILIK-EXPORT-PEGA.md`."*
>
> ⭐ **P18 ditarik.** Suratnya **dibatalkan sebelum dikirim**; berkasnya diganti blok RALAT.

⭐⭐ **Penahan proyek NIHIL.** `[terverifikasi]` 2026-09-22

> ⛔⛔ **RALAT.** `[penyimpangan sadar]` 2026-09-22 — kalimat ini semula berbunyi:
> > *"⭐ **Tinggal satu yang menahan.** Suratnya sudah siap: `SURAT-P1-KE-DBA.md`."*
>
> **P1 terjawab work owner pada hari yang sama.** Naskah **empat** stored procedure diterima
> *(`PEGA_TREATY_IN` · `PEGA_JSON_POLIS_TREATYIN` · `PROC_GENERATE_SEQUENCE_NUMBER` ·
> `PEGA_DELETE_ERROR_KONVERSI`)*, ditambah **dua contoh** `DATA_JSON`. Suratnya tidak perlu dikirim.

⚠️ **Yang tersisa bukan penahan, melainkan pekerjaan tim migrasi sendiri:** **sensus properti**
`PolicyTreatyIn` — 94 nama unik dan 9 simpul bersarang baru **batas bawah**, dan rancangan tabel
penyimpanan bergantung padanya. Lihat tiket `00`.

Butir terbuka yang **tidak menahan**: `RNM_SHARE` persentase dari apa · tipe Oracle kolom uang ·
empat wadah layar berisi 104 medan (P44) ·
ketidakseragaman presisi 4 lawan 8 desimal.

> ⚠️ **RALAT** 2026-09-22 — daftar ini semula memuat ⛔ *"sembilan medan wajib-tapi-terkunci
> menunggu P18"*. **Tidak lagi menunggu apa pun.**

---

## 7 · Cara menguji ulang berkas ini

Seluruh angka dapat diperiksa ulang dari korpus. Pakai Python, bukan loop shell — jalur korpus
memuat spasi.

```python
import glob, os, re, html
base = r"D:/XML/RNM_BRD/NB Treaty In"

# berkas dan byte
allf = glob.glob(os.path.join(base, "**", "*.xml"), recursive=True)
print(len(allf), sum(os.path.getsize(f) for f in allf))

# langkah Property-Set
n = 0
for f in glob.glob(os.path.join(base, "Activity", "*.xml")):
    u = html.unescape(html.unescape(open(f, encoding="utf-8", errors="replace").read()))
    n += sum(1 for s in re.findall(r"<pyStepsActivityName>([^<]*)</pyStepsActivityName>", u)
             if s == "Property-Set")
print(n)
```

Jangkauan aturan tersimpan di `jangkauan.json`; daftar 51 aturan terjangkau yang semula diminta
pada P18 di `lingkup-p18-diminta.json` — ⚠️ **permintaannya ditarik, daftarnya tetap sahih.**

### ⭐ Lima jebakan yang sudah pernah menjerat  — ⚠️ *semula empat*

0. ⭐⭐ **Nama tag berbeda menurut JENIS aturan.** Di `Activity` yang berisi adalah
   `PropertiesName`/`PropertiesValue` **tanpa** awalan `py`; di `DataTransform` **kebalikannya**.
   Memeriksa satu jenis lalu menyimpulkan untuk jenis lain melahirkan **nol palsu**.
   ⛔ Jebakan inilah yang melahirkan P18. **Periksa daftar tag yang benar-benar ada lebih dulu.**

1. `<rowdata REPEATINGINDEX="n"/>` yang menutup sendiri tidak tertangkap pola
   `<rowdata...>(.*?)</rowdata>` — sel kosong hilang, pasangan bergeser.
2. `pyRowNum` pada `pyOrConditions` berbasis **nol**; baris tabel berbasis satu.
3. `pyStepsPage`, `pyCriteriaValue`, `pyResult`, `pyShapeName`, `pyTaskLabel` **tidak ada** di
   ekspor ini. Periksa daftar tag lebih dulu.
4. Menyatakan sesuatu nihil tanpa menyebut lingkup penelusuran. Sebutkan selalu apa yang dicari
   dan di mana.

---

*Disusun oleh sesi asisten, 22 September 2026. Belum diverifikasi ulang oleh sesi eksekutor.
Setiap angka di berkas ini dapat diuji ulang dari korpus dengan cara di bab 7.*
