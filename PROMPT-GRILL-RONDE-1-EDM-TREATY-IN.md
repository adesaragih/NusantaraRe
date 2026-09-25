# PROMPT GRILLING — RONDE 1 · EDM Treaty In

> Salin seluruh isi berkas ini sebagai prompt ke sesi eksekutor.
>
> ⚠️ **Ronde ini berbeda sifatnya dari empat ronde NB Treaty In.** Sebagian besar aturan modul ini
> **sudah dikenali** lewat NB Treaty In. Yang dicari adalah **selisihnya**, bukan seluruhnya.

---

## 0. LINGKUP — DIKUNCI

Modul utama: **`D:\XML\RNM_BRD\EDM Treaty In`**.

`D:\XML\RNM_BRD\NB Treaty In` boleh dibaca **hanya sebagai pembanding**. Temuan dari sana tidak
masuk sensus modul ini.

`D:\XML\RNM_BRD\` adalah korpus **READ-ONLY**. Menulis hanya ke `OUTPUT_HASIL_RNM\`.
`D:\XML\nusantara-re\` **terlarang**.

---

## 1. YANG SUDAH DIKETAHUI — JANGAN DITEMUKAN ULANG

EDM Treaty In **163 berkas, 17.639.246 B**. Dibandingkan NB Treaty In:

| | Berkas |
| --- | ---: |
| Bernama sama dengan NB Treaty In | **73** |
| ⇒ di antaranya **isinya sama**, beda hanya metadata ekspor | **69** |
| ⇒ kelas berbeda | **3** |
| ⇒ jumlah langkah berbeda | **1** |
| **Hanya ada di EDM** | **90** |

⭐ **Enam puluh sembilan aturan sudah dikenali lewat NB Treaty In.** Jangan diurai ulang. Bila
sebuah aturan bernama sama dan isinya sama, cukup sebut *"sama dengan NB Treaty In"*.

Berkas keadaan NB — `.scratch\nb-treaty-in\KEADAAN-NB-TREATY-IN.md` — memuat 12 ketetapan yang
sudah diputuskan work owner. **Berlaku juga di sini kecuali terbukti sebaliknya.** Bila EDM
berperilaku berbeda pada salah satunya, **itu temuan**, dan wajib dilaporkan.

### Empat aturan yang benar-benar berbeda — periksa keempatnya

| Aturan | NB Treaty In | EDM Treaty In |
| --- | --- | --- |
| `serviceInsertArasapas_act` | kelas `Data-PolicyTreatyIn` | kelas `Work` |
| `SetCategoryAttach` | kelas `Data-OfferFacIn` | kelas `Work` |
| `IsUW` | kelas `Work` | kelas `Data-OfferFacIn-LocationReinsurance` |
| `CountSpreading_Act` | 7 langkah | **8 langkah** |

Yang terakhir paling menarik: satu langkah lebih banyak, di aturan penghitung spreading.
**Cari langkah mana, dan apa yang dikerjakannya.**

---

## 2. SASARAN RONDE INI

### Sasaran 1 — DATA LAMA, DATA BARU, DAN SELISIHNYA (utama)

Inilah inti endorsemen, dan NB Treaty In tidak punya apa pun yang menyerupainya. Tujuh Section
besar, seluruhnya hanya ada di EDM:

| Ukuran | Berkas |
| ---: | --- |
| 947.721 | `Section\DetailPolicyTreatyInPropNewData2` |
| 746.314 | `Section\DetailPolicyTreatyInPropOldData2` |
| 745.389 | `Section\DetailPolicyTreatyInPropValueDifference` |
| 745.191 | `Section\DetailPolicyTreatyInPropOldData` |
| 674.337 | `Section\DetailPolicyTreatyInPropNewData` |

Ditambah `Activity\CalculateDifferenceEDM_act` (216.239 B).

**Yang dicari:**

- medan apa saja yang dibandingkan lama-lawan-baru, dan **medan apa yang tidak**
- bagaimana selisih dihitung — per medan, atau hanya untuk medan uang
- apa arti akhiran `2` pada `PropNewData2` dan `PropOldData2` — dua tingkat, atau dua tampilan
- apakah data lama **disalin** ke berkas endorsemen atau **dibaca** dari polis induk
- apa yang terjadi bila polis induk berubah sesudah endorsemen dibuat

### Sasaran 2 — ADENDUM DAN PREMI TAMBAHAN

`Section\DetailPolicyTreatyInAddendum` (913.122 B) · `Activity\InsetTreatyInProdAddendum_Act`
(307.657 B) · `Section\DetailPolicyTreatyInAddPremi` (688.183 B) ·
`Section\DetailPolicyAddPremiDetail` (249.693 B).

**Yang dicari:** apakah adendum dan premi tambahan dua hal berbeda atau dua nama untuk satu hal;
kapan masing-masing terbit; dan apakah keduanya menghasilkan nomor tersendiri.

### Sasaran 3 — ⛔ PANGGILAN KE SISTEM LUAR

**NB Treaty In tidak punya ini sama sekali.** EDM punya satu `ConnectREST`:

```
ConnectREST\convertJsonNusareToProduction
   pyServiceName          convertJsonNusareToProduction
   pyBaseURLSelectionType SETTING
   pyBaseURLSetting       LinkService!LinkService
   pyUseAuthentication    false
   pyResponseTimeout      30000
   dikirim dari clipboard .OfferFacIn.PolicyData.PolicyNo · .pzInsKey · .pxCreateDateTime
   balasan dipetakan ke   .StatusService
```

Dan setelan yang memberinya alamat:

```
SystemSettings\LinkService     pySetting = "=ResponLink.URL"
```

Alamatnya **bukan nilai tetap**, melainkan ungkapan yang dibaca dari properti saat berjalan.

**Yang dicari:**

- aturan mana yang **memanggil** layanan ini, dan pada tahap apa
- apa yang terjadi bila panggilan **gagal** atau **habis waktu** — 30 detik itu lama
- dari mana `ResponLink.URL` terisi
- ⛔ **`pyUseAuthentication` = false** — apakah benar tidak ada autentikasi sama sekali
- apakah `.StatusService` dibaca sesudahnya, dan apa akibat nilainya

⚠️ Ini menyentuh **ADR-0013** *(resolusi endpoint via LinkService)* dan **ADR-0004** *(endpoint
sebagai env var, superseded)*. **Jangan membuat ADR baru** — laporkan saja bila temuannya
bertentangan.

### Sasaran 4 — ALUR

Satu `Flow`, 202.269 B. Bandingkan dengan `InputRealizationTreatyIn` milik NB — 31 kotak,
32 sambungan, empat cabang putusan.

**Yang dicari:** tahap apa yang ada di EDM tetapi tidak di NB, tangga persetujuannya sama atau
berbeda, dan di mana panggilan ke sistem luar duduk di dalam alur.

### Sasaran 5 — SEMBILAN PULUH BERKAS EDM-SAJA

Sensus: kelompokkan ke-90 menurut perannya. Sebutkan berapa yang **terjangkau** dari titik masuk
nyata dan berapa **yatim** — sama seperti cara yang dipakai NB Treaty In, yang menemukan 28 aturan
yatim milik modul lain.

⚠️ Folder modul adalah **dependency closure**. Jangan menganggap isi folder sebagai daftar
pekerjaan EDM.

---

## 3. DISIPLIN

Setiap pernyataan membawa bukti: **jalur berkas + tipe rule + nama rule**.

Penanda: `[terverifikasi]` · `[dugaan]` · `[terbuka]` · `[data DBA]` · `[keputusan work owner]` ·
`[penyimpangan sadar]`.

Setiap angka membawa **perintah yang menghasilkannya**. Sensus dihitung **dua cara yang
benar-benar berbeda**.

**Jangan menutup pertanyaan terbuka sendiri.** **Jangan menyalin nilai berupa nama orang** — catat
nama rule, nama medan, dan jumlahnya saja. Nomor polis tidak ditulis apa adanya.

**Nol kode, nol DDL, nol usulan daftar kolom.** Ronde ini membaca.

### Lima jebakan yang sudah pernah menjerat

1. `<rowdata REPEATINGINDEX="n"/>` yang menutup sendiri tidak tertangkap pola
   `<rowdata...>(.*?)</rowdata>` — sel kosong hilang, pasangan bergeser.
2. `pyRowNum` pada `pyOrConditions` berbasis **nol**; baris tabel berbasis satu.
3. `pyStepsPage`, `pyCriteriaValue`, `pyResult`, `pyShapeName`, `pyTaskLabel` **tidak ada** di
   ekspor ini. Periksa daftar tag lebih dulu.
4. Menyatakan sesuatu nihil tanpa menyebut lingkup penelusuran.
5. ⭐ **Naskah SQL ADA** — di `<pyBrowseSQL>` pada `RDBList`. EDM punya **36 RDBList**.
   Baris tabel keputusan ada di `pyCondition`, `pyOrConditions`, `pyResults`.

Jalur korpus memuat spasi — pakai Python, bukan loop shell.

---

## 4. KELUARAN — DUA BERKAS

### `OUTPUT_HASIL_RNM\.scratch\edm-treaty-in\grilling-ronde-1.md`

Temuan lengkap, bab A sampai E mengikuti kelima sasaran.

Bab pembuka wajib: **apa yang sama dengan NB Treaty In**, supaya pembaca tahu apa yang
**tidak** perlu dibaca ulang.

### `OUTPUT_HASIL_RNM\.scratch\edm-treaty-in\PERTANYAAN-RONDE-1.md`

**Hanya pertanyaan.** Bentuknya sama seperti modul NB:

```
## Pn — <judul dalam bahasa bisnis>

<satu atau dua paragraf, bahasa bisnis>

**Konteks:** <mengapa ini tidak bisa dijawab dari berkas>

**Bentuk jawaban yang diharapkan:** <pilihan ganda / butuh berkas / butuh contoh>

**Dampak bila salah:** <akibat nyata>

rujukan: <bukti teknis>

**Jawaban:**

> *(tulis di sini)*
```

⭐ **Nomor pertanyaan mulai dari P50.** P1–P49 dipakai NB Treaty In; P41 dan P48 sudah ditarik.
Penomoran **berlanjut**, tidak mengulang dari satu — supaya satu nomor berarti satu pertanyaan di
seluruh proyek.

Kelompokkan menurut pemilik: pengembang Pega lama · Product & Underwriting · DBA · Finance · IAM ·
pemilik export Pega.

⚠️ Bila sebuah pertanyaan **sudah dijawab** untuk NB Treaty In dan jawabannya berlaku di sini,
**jangan ditanyakan ulang** — rujuk nomornya.

---

## 5. BAB WAJIB — TELEMETRI EKSEKUSI

Berkas grilling ditutup dengan `## TELEMETRI EKSEKUSI`.

Angka token sejati tidak terlihat dari dalam sesi. Ukur dari luar:

```
claude --print --output-format json "<prompt>" > hasil.json
```

Catat token keluaran · cache-read · panggilan alat · durasi · biaya · byte dibaca.
Bila tidak diukur dari luar, **katakan begitu**; bila diperoleh dengan mengurangkan baseline,
**katakan itu juga**.

---

## 6. YANG MENANDAKAN RONDE INI BERHASIL

- Ke-90 berkas EDM-saja **dibuka**, bukan hanya didaftar. Sebutkan byte yang benar-benar dibaca.
- Pola **data lama / data baru / selisih** dijelaskan sampai tingkat medan.
- Panggilan ke sistem luar tertelusur: siapa memanggil, kapan, apa akibat gagal, ada autentikasi
  atau tidak.
- Keempat aturan yang berbeda dari NB diperiksa satu per satu.
- Sensus terjangkau lawan yatim dihitung dua cara.
- Pertanyaan baru ada di berkasnya sendiri, bernomor mulai **P50**, dan **tidak mengulang** yang
  sudah dijawab di NB Treaty In.
- Bab telemetri jujur tentang cara pengukurannya.

---

*Disusun 22 September 2026, sesudah NB Treaty In selesai sampai tahap tiket. Seluruh angka di dalam
prompt ini diverifikasi dari korpus dan boleh diuji ulang.*
