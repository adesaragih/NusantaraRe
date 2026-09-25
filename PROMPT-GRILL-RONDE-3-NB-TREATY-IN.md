# PROMPT — Grilling Ronde 3: NB Treaty In

> Ketik sendiri di sesi kerja baru. `grill-with-docs` ber-`disable-model-invocation: true`.
>
> Ronde 1: 868 baris, 24 butir. Ronde 2: 1.048 baris, 13 butir baru, P18–P28.

---

## Kenapa ronde 3 ada, dan kenapa isinya berbeda dari ronde 1 dan 2

Ronde 2 membuktikan **942 langkah `Property-Set` tidak terekspor isinya** — nol dari 942. Rantai
perhitungan tidak dapat direkonstruksi dari langkah `Activity`.

Justru karena itu, **sumber penetapan nilai yang MASIH terbaca menjadi jauh lebih berharga.**
Ronde 3 memburu sisa itu. Terukur, dan belum satu pun disentuh di ronde 1 maupun 2:

| Sumber | Jumlah | Keadaan |
| --- | ---: | --- |
| `DataTransform` | 12 berkas | **terbaca penuh** — 80 nama properti, 62 nilai, 12 aksi |
| Langkah `Java` ber-source | 8 | **source terekspor** — 343 di seluruh korpus, nol pernah dibaca |
| `pyStepsObjectName` di `Activity` | 457 berisi | sasaran langkah terbaca walau nilainya tidak |
| `ReportDefinition` | 15 berkas | nol dibuka |
| `FlowAction` | 9 berkas | nol dibuka |
| `Flow` | 1 berkas | ditelusur di D2, belum dibaca ulang |

Perintah audit:

```
grep -ohE '<pyPropertiesName>[^<]+</pyPropertiesName>' "NB Treaty In"/DataTransform/*.xml | wc -l   # 80
grep -ohE '<pyStepsJavaSource>[^<]+'                   "NB Treaty In"/Activity/*.xml      | wc -l   # 8
grep -ohE '<pyStepsObjectName>[^<]+</pyStepsObjectName>' "NB Treaty In"/Activity/*.xml    | wc -l   # 457
```

---

```
/mattpocock-skills:grill-with-docs

## Lingkup

Modul ronde ini: HANYA "NB Treaty In" di D:\XML\RNM_BRD\.

⛔ JANGAN membuka "EDM Treaty In". OQ-025 tetap berdiri sebagai pertanyaan.
⛔ JANGAN membuka folder korpus modul lain. ⛔ JANGAN membuka D:\XML\nusantara-re\.
⛔ JANGAN membaca .scratch\ konteks lain. Bila sambungan sesi menyuntikkannya tanpa diminta,
   catat di Bab E dan nyatakan apa yang tidak merembes.

## Baca lebih dulu

- .scratch/nb-treaty-in/grilling-ronde-1.md dan grilling-ronde-2.md — SELURUHNYA
- discovery/flows/NB Treaty In.md §5 dan §7
- discovery/open-questions.md untuk OQ yang disebut kedua ronde
- CONTEXT.md, docs/adr/

## Keluaran

Tulis HANYA ke .scratch/nb-treaty-in/grilling-ronde-3.md. SATU berkas.
⛔ JANGAN menyunting berkas ronde 1 atau ronde 2. Bila keduanya keliru, tulis RALAT di berkas
   ronde 3, kutip kalimat lamanya utuh, sebut cara mana yang keliru dan kenapa.
⛔ JANGAN membuat spec.md. ⛔ JANGAN membuat tiket.

## Yang SUDAH selesai — jangan diulang

Ronde 1: 5 blok PL/SQL dibaca utuh · 41 RDBList terurai · 3 procedure, nol badannya ada ·
31 objek Oracle, 2 ditulis lewat SQL terbaca · sensus 278 berkas, 10 tipe rule.
Ronde 2: 25 Section + 6 Harness disisir (733 syarat, 229 bermakna, 11 murni tata letak) ·
75 dari 75 When dibaca · 2 DecisionTable, baris keputusan tidak terekspor · 30 rule hitung uang ·
942 Property-Set nol isinya · 20 nomor polis produksi di 3 aturan hidup · 91 elemen mati ·
18 dari 52 pyConditionString tidak cocok dengan aturan yang dieksekusi.

## EMPAT SASARAN RONDE 3

### Sasaran 1 — 12 DataTransform ⭐ prioritas tertinggi

Ini satu-satunya tipe rule di modul ini yang menetapkan nilai properti DAN isinya terekspor.
Terukur: 80 nama properti, 62 nilai, 12 aksi.

  AddToListCommentsPolicyTreatyIn_DT · DeptHeadTreatyInUW_preDT · DeptHeadTreatyIn_UW_postDT
  InboxPolicyTreatyIn_postDT · InputPolicyTreatyIn_preDT · SearchHierarkiSourceBizAgent_PostDT
  SystemSetOneYear_DT · TestTreatyToFacStatus · TreatyEnableDisableInput · btnCedingCO_DT
  btnSOB_DT · setCategoryAttachment_DT

Untuk tiap rule: properti apa yang ditetapkan · nilai atau ekspresi apa · aksi apa
(Set / Append / Remove / When) · dipanggil dari mana.

⭐ Nyatakan secara eksplisit: **berapa banyak dari rantai penetapan nilai modul ini yang dapat
dipulihkan dari 12 berkas ini**, dan berapa yang tetap hilang bersama 942 langkah. Itu angka yang
menentukan apakah spec perhitungan mungkin ditulis sebagian.

### Sasaran 2 — 8 langkah Java ber-source ⭐ belum pernah dibaca siapa pun

Tag `<pyStepsJavaSource>` berisi di 8 tempat. Berkasnya:

  FetchMasterTreatyIn · GetTreatyName · InputPolicyTreatyInDetail_NonProp
  InputPolicyTreatyInDetail_preACT · InputPolicyTreatyOutDetail_NonProp
  InputPolicyTreatyOutDetail_preACT · SetTreatyIn_Act

Untuk tiap langkah: apa yang dilakukan source-nya · properti apa yang dibaca dan ditulis ·
apakah ia menghitung uang · apakah ia memanggil sesuatu di luar Pega.

⚠️ Kutip source-nya seperlunya untuk membuktikan klaim, jangan seluruhnya.
⛔ JANGAN menulis ulang sebagai Go. ⛔ JANGAN menebak maksud kode yang tidak jelas — catat
sebagai [terbuka].

### Sasaran 3 — 37 berkas yang masih nol dibuka

  15 ReportDefinition · 9 FlowAction · 12 DataTransform (sasaran 1) · 1 Flow

ReportDefinition: kolom apa yang ditampilkan, dari tabel apa, dengan saringan apa.
FlowAction: validasi saat pindah tahap, dan siapa yang boleh menekannya.
Flow: konfirmasi terhadap telusur D2, bukan telusur ulang.

### Sasaran 4 — 62 Activity non-uang, dipetakan bentuknya

Nilainya hilang, tetapi bentuknya tidak. Yang masih terbaca: jenis langkah, urutannya, dan
`pyStepsObjectName` (457 berisi).

Sensus jenis langkah di 92 Activity, sudah terukur — konfirmasi lalu perdalam:
  Property-Set 942 · Page-Set-Messages 105 · RDB-List 66 · Page-New 33 · Property-Remove 25 ·
  Page-Remove 25 · Property-Set-Messages 17 · Page-Clear-Messages 17 · Java 8 · Page-Copy 6 ·
  Obj-Save 6 · Obj-Browse 6 · langkah Call bernama

⭐ Petakan rantai panggil antar-Activity: siapa memanggil siapa. Itu terbaca dari langkah `Call`
walau parameternya tidak. Hasilnya: kerangka alur perhitungan tanpa isinya — berguna saat ekspor
ulang datang, karena isi tinggal ditempelkan ke kerangka yang sudah jadi.

Bila anggaran habis, sasaran 4 yang pertama dipotong, DAN pemotongannya dilaporkan.

## Aturan yang mengikat

CLAUDE.md §4 dan §4a berlaku penuh — sama persis seperti ronde 1 dan 2:

1. Bukti wajib path + class + nama rule.
2. Identitas rule dari <pxInsName>. <pzOriginalInstanceKey> BUKAN identitas.
3. Label wajib: [terverifikasi] / [dugaan] / [terbuka] / [data DBA].
4. Setiap angka disertai perintah auditnya.
5. Sensus DUA cara berbeda, jendelanya disebut. Berselisih ⇒ "belum punya data", cantumkan keduanya.
6. Properti dicari dengan awalan halaman.
7. Grep dulu, baca rentang seperlunya. Catat berkas yang belum habis dibaca.
8. Jangan menebak skema, tipe kolom, arti kode, singkatan, rumus, isi procedure, struktur JSON.
9. Jangan menyatukan identitas berkonflik.
10. Nilai nama orang tidak disalin.
11. Nomor OQ dari register. 20 butir baru ronde 1 dan 2 BELUM bernomor register — rujuk sebagai
    "ronde 1 baru #n" dan "ronde 2 baru #n". Nomor bebas berikutnya OQ-073.
12. ⛔ Agent TIDAK BOLEH menutup butir terbuka sendiri.
13. ⛔ Kode NOL. DDL NOL. CREATE TABLE NOL. Daftar kolom usulan NOL.

⭐ Uji instrumen WAJIB. Ronde 1 gagal 2 kali, ronde 2 gagal 4 kali, keenam-enamnya tertangkap
sebelum publikasi. Uji tiap alat penghitung atas butir yang jawabannya sudah diketahui, dan
laporkan hasil ujinya bersama sensusnya.

⚠️ Satu galat yang SUDAH TERULANG dua ronde: ejaan pemilik dicampur — [Product+UW] dan
[Product+Underwriting]. Pakai [Product+Underwriting] saja, dan cek ejaannya sebelum merekap.

## Bab yang WAJIB ada

Bab A — Tambahan stored procedure dan SQL mentah. Boleh nihil asal dibuktikan.
Bab B — Tambahan objek Oracle. Boleh nihil asal dibuktikan.
Bab C — Register [terbuka]: butir lama yang berubah status + butir baru. Rekap per pemilik,
        ejaan diseragamkan.
Bab D — PERTANYAAN SIAP KIRIM, HANYA yang baru, lanjut dari P29. Bentuk sama seperti ronde 1
        dan 2. ⛔ Nol dijawab sendiri.
Bab E — Lampiran bukti isolasi, dengan perintah audit tiap baris.

⭐ Bab F — PEMULIHAN RANTAI NILAI. Bab baru, hanya di ronde ini.
   Jawab satu pertanyaan dengan angka: dari seluruh penetapan nilai di modul ini, berapa persen
   yang dapat dipulihkan dari DataTransform + Java + pyStepsObjectName, dan berapa persen yang
   tetap hilang bersama 942 langkah Property-Set.
   Ini yang menentukan apakah spec perhitungan dapat ditulis sebagian sambil menunggu ekspor ulang,
   atau tidak sama sekali.

## Tutup dengan

- Cakupan sesudah ronde 3: berapa dibaca utuh · berapa disisir · berapa MASIH nol dibuka
- Berapa butir lahir · berapa ditutup (harus NOL)
- Apakah frontier masih terbuka
- ⭐ Vonis: apakah NB Treaty In sudah habis digrilling dengan bahan yang ada, atau masih ada ronde 4
```

---

## Menjalankan sambil mengukur biaya

```bash
cd /d/XML/RNM_BRD/OUTPUT_HASIL_RNM
export PYTHONIOENCODING=utf-8

run_json=".run-nb-treaty-in-ronde-3.json"
claude --print --output-format json -- "$(cat PROMPT-GRILL-RONDE-3-NB-TREATY-IN.md)" > "$run_json"

py -c "
import json
d = json.load(open('$run_json', encoding='utf-8'))
print('biaya USD     :', d.get('total_cost_usd'))
print('jumlah giliran:', d.get('num_turns'))
for k, v in (d.get('usage') or {}).items():
    print(f'  {k}: {v}')
"
```

## Perkiraan

| | Ronde 1 | Ronde 2 | Ronde 3 (perkiraan) |
| --- | ---: | ---: | ---: |
| keluaran | 192.021 | 225.898 | **160.000–200.000** |
| cache baca | 13.532.360 | 22.688.258 | **20.000.000–26.000.000** |
| panggilan | 69 | 72 | **60–75** |

Lebih ringan dari ronde 2 pada keluaran — bahan bacanya jauh lebih kecil (12 DataTransform 12 KB,
8 langkah Java, 37 berkas menengah) dibanding 16,6 MB Section dan Harness. Tetapi cache baca naik
karena dua berkas ronde sebelumnya ikut dibaca tiap giliran.
