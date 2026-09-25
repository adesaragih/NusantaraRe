# PROMPT — Grilling Ronde 2: NB Treaty In

> Ketik sendiri di sesi kerja baru. `grill-with-docs` ber-`disable-model-invocation: true`,
> agent tidak dapat memanggilnya (CLAUDE.md §8).
>
> Ronde 1 selesai 2026-09-22: `.scratch/nb-treaty-in/grilling-ronde-1.md`, 868 baris, 24 butir
> terbuka, nol ditutup.

---

```
/mattpocock-skills:grill-with-docs

## Lingkup

Konteks: Treaty Inward — Realisasi & Endorsement.
Modul ronde ini: HANYA "NB Treaty In" (278 berkas .xml) di D:\XML\RNM_BRD\.

⛔ JANGAN membuka "EDM Treaty In". Blocker OQ-025 tetap berdiri dan tetap dicatat sebagai
pertanyaan, bukan dipecahkan dengan meminjam varian dari sana. Izin membukanya adalah keputusan
work owner yang belum diberikan.
⛔ JANGAN membuka folder korpus modul lain mana pun.
⛔ JANGAN membuka D:\XML\nusantara-re\ — di-blacklist total.
⛔ JANGAN membaca isi .scratch\ milik konteks lain. Bila sambungan sesi menyuntikkannya tanpa
   diminta, catat kejadiannya di Bab E dan nyatakan apa yang tidak merembes.

## Baca lebih dulu

- .scratch/nb-treaty-in/grilling-ronde-1.md — SELURUHNYA. Ini dasar ronde 2.
- discovery/flows/NB Treaty In.md §5 dan §7
- discovery/open-questions.md — hanya OQ yang disebut di Bab C ronde 1
- CONTEXT.md, docs/adr/

## Keluaran

Tulis HANYA ke .scratch/nb-treaty-in/grilling-ronde-2.md. SATU berkas.
⛔ JANGAN menyunting grilling-ronde-1.md. Bila ronde 1 keliru, tulis RALAT di berkas ronde 2,
   kutip kalimat lamanya utuh, dan sebut cara mana yang keliru dan kenapa.
⛔ JANGAN membuat spec.md. ⛔ JANGAN membuat tiket.

## Yang SUDAH dikerjakan ronde 1 — jangan diulang

- 5 blok PL/SQL dibaca utuh, SQL-nya dikutip: GetSequenceNumber_SQL · InsertHistoryAkseptasiPega_Sql
  · InsertViewSuggest_SQL · SavePolisTreatyIn_SQL · SaveTreatyIn
- seluruh 41 RDBList disisir dengan pengurai penuh
- sensus tipe rule, dua cara, sepakat: Activity 92 · When 75 · RDBList 41 · Section 25 ·
  ReportDefinition 15 · DataTransform 12 · FlowAction 9 · Harness 6 · DecisionTable 2 · Flow 1 = 278
- Bab A: 3 procedure, nol badannya ada — PEGA_TREATY_IN 24 param · PEGA_JSON_POLIS_TREATYIN 8 ·
  PROC_GENERATE_SEQUENCE_NUMBER 5
- Bab B: 31 objek Oracle, hanya 2 ditulis lewat SQL terbaca, 29 hanya dibaca
- Bab C: 24 butir (17 terdaftar + 7 baru), nol ditutup
- Bab D: 17 pertanyaan siap kirim, P1–P17

Ronde 2 MEMPERDALAM ini, tidak menghitung ulang. Bila sebuah angka ronde 1 ternyata salah,
RALAT dengan aturan di atas.

## EMPAT SASARAN RONDE 2 — urut menurut prioritas

### Sasaran 1 — 31 berkas yang nol dibuka: seluruh Section dan seluruh Harness

Ronde 1 tidak membuka satu pun dari 25 Section dan 6 Harness. Terukur: 16.635.834 byte,
46,6 % modul menurut ukuran. Kelima berkas terbesar modul ada di situ:

  Section/GeneralPolicyTreatyIn.xml                   1.926.378 B
  Section/DetailPolicyTreatyIn.xml                    1.916.240 B
  Section/DetailPolicyTreatyInNonProportional.xml     1.662.941 B
  Section/GeneralDeptHeadTreatyIn_UW.xml              1.584.498 B
  Section/DetailDeptHeadTreatyIn_UW.xml               1.567.128 B

⛔ JANGAN membaca berkas-berkas ini utuh. Aturan §4 butir 4 berlaku: grep dulu, baca rentang
seperlunya, dan CATAT bila sebuah berkas belum habis dibaca.

Yang dicari di layar — dan hanya ini:
- validasi yang hanya hidup di layar: medan wajib isi, batas nilai, format
- medan tersembunyi dan medan hanya-baca
- tombol dan bagian yang muncul atau hilang menurut jabatan atau status
- rujukan ke When atau DecisionTable sebagai penyaring tampilan
⛔ BUKAN tata letak, BUKAN urutan kolom, BUKAN gaya tampilan. Bila sebuah Section ternyata murni
tata letak, tulis satu baris "murni tata letak" dan lanjut. Itu temuan yang sah dan murah.

### Sasaran 2 — rantai hitung uang, yang ronde 1 tidak sentuh sama sekali

Ronde 1 menutup dengan [Finance] = 0 butir, dan menyebut sendiri bahwa nol itu lubang, bukan
kebersihan. Tutup lubang itu.

Terukur: 30 berkas di Activity\ dan When\ bernama pola hitung uang. Perintah:
  find Activity When -name '*.xml' | grep -icE 'count|sum|total|calc|premi|rate|limit|share|comm|brokerage|tax'

Yang terbesar, dan belum pernah dibaca siapa pun di proyek ini:
  Activity/SumTSIPremiSpreadedRNM_FIRE_Act.xml    996.052 B
  Activity/SumTSIPremiSpreadedRNM_ANEKA_Act.xml   748.216 B
  Activity/SumTSIPremiSpreadedRNM_Act.xml         516.395 B
  Activity/ProtectFIREMBUPA_Act.xml               875.119 B
  Activity/CheckSpreadingProtectAnekaGolf_ACT.xml 398.150 B

Untuk tiap rule hitung uang, catat: apa yang masuk · apa yang keluar · disimpan atau dihitung
ulang · apakah nilainya bermata uang atau telanjang · apakah ada ambang ter-hardcode.

⚠️ Peringatan yang sudah terbukti di modul lain: korpus ini memuat nilai uang ter-hardcode TANPA
mata uang, dan setidaknya satu ambang dibandingkan sebagai STRING, bukan angka. Bila menemukan
pola itu di sini, catat persis bentuknya dan kutip barisnya.
⛔ JANGAN menulis ulang rumusnya sebagai kode. ⛔ JANGAN menebak arti singkatan.

### Sasaran 3 — dua DecisionTable, keduanya belum pernah dibaca

  DecisionTable/isApproved.xml     ← OQ-026, gerbang pertama seluruh alur
  DecisionTable/BusinessType_DeT.xml

Untuk isApproved: baca KEDUA varian — When\isApproved.xml dan DecisionTable\isApproved.xml —
lalu nyatakan apa bedanya. Ronde 1 mencatat keduanya dibuat berselisih ±26 menit di hari sama.
⚠️ Bila baris keputusan DecisionTable tidak ikut terekspor, katakan begitu. Itu sudah terjadi pada
49 DecisionTable lain di korpus. ⛔ JANGAN menebak isinya.

### Sasaran 4 — 92 Activity dan 75 When, dibaca langkah demi langkah

Ronde 1 hanya menyisir pola. 75 When adalah jumlah tidak biasa: modul ini memutuskan banyak hal
lewat penggolong, dan isi penggolongnya belum dibaca.

Kerjakan setelah sasaran 1–3. Bila anggaran habis, BERHENTI dan katakan berapa yang terbaca dari
berapa. ⛔ Sasaran yang dikerjakan separuh lalu dilaporkan utuh adalah kegagalan yang lebih buruk
daripada sasaran yang tidak dikerjakan.

## Aturan yang mengikat — sama seperti ronde 1

CLAUDE.md §4 dan §4a berlaku penuh:

1. Setiap klaim perilaku wajib bukti path + class + nama rule.
2. Identitas rule = class / nama dari <pxInsName>. <pzOriginalInstanceKey> BUKAN identitas.
3. Label wajib: [terverifikasi] / [dugaan] / [terbuka] / [data DBA].
4. Setiap angka disertai perintah audit yang menghasilkannya.
5. Tiap sensus dihitung DUA cara yang benar-benar berbeda, dan jendelanya disebut. Berselisih ⇒
   tulis "belum punya data" dan cantumkan keduanya.
6. Properti korpus dicari dengan awalan halaman: InputData.CARI21, bukan CARI21 telanjang.
7. Grep dulu, baca rentang seperlunya. Catat berkas yang belum habis dibaca.
8. Jangan menebak skema, tipe kolom, arti kode, kepanjangan singkatan, rumus, isi stored
   procedure, atau struktur JSON.
9. Jangan menyatukan identitas berkonflik. Cek _oq011-konflik-isi.md lebih dulu.
10. Nilai nama orang tidak disalin. Catat rule, tag, dan jumlahnya saja.
11. Nomor OQ dari register discovery/open-questions.md. Tujuh butir baru ronde 1 BELUM bernomor
    register; rujuk sebagai "ronde 1 baru #1" sampai #7. Nomor bebas berikutnya OQ-073.
12. ⛔ Agent TIDAK BOLEH menutup butir terbuka sendiri.
13. ⛔ Kode NOL. DDL NOL. CREATE TABLE NOL. Daftar kolom usulan NOL.

⭐ Ujian instrumen WAJIB, dan ronde 1 membuktikan kenapa: instrumennya gagal DUA KALI dan ujian
menangkapnya — sebuah tabel terbaca sebagai procedure, dan koma di dalam TO_DATE membuat jumlah
parameter salah. Uji tiap alat penghitung atas butir yang jawabannya sudah diketahui, dan laporkan
hasil ujinya bersama sensusnya. Bila instrumen gagal, tulis kegagalannya — jangan diam-diam
diperbaiki.

## Bab yang WAJIB ada di keluaran

Bab A — Tambahan stored procedure dan SQL mentah yang ditemukan di Section, Harness, Activity,
        When, dan DecisionTable. Bila tidak ada tambahan, katakan begitu dengan perintah auditnya.
Bab B — Tambahan objek Oracle. Sama, boleh nihil asal dibuktikan.
Bab C — Register [terbuka] berpemilik: butir ronde 1 yang berubah statusnya, ditambah butir baru
        ronde 2. Rekap per pemilik. ⭐ [Finance] wajib disebut eksplisit: berapa butir uang lahir
        di ronde ini, dan bila tetap nol, kenapa.
Bab D — PERTANYAAN SIAP KIRIM, HANYA yang baru. Lanjutkan penomoran dari P18. Bentuknya sama
        persis seperti ronde 1: pertanyaan dalam bahasa bisnis tanpa istilah Pega tanpa nama rule
        tanpa path · konteks satu kalimat · bentuk jawaban yang diharapkan · dampak bila salah
        dijawab · rujukan teknis di baris terpisah paling bawah.
        ⛔ JANGAN menjawab sendiri satu pun.
Bab E — Lampiran bukti isolasi, dengan perintah audit tiap baris.

## Tutup dengan

- Cakupan sesudah ronde 2: berapa berkas dibaca utuh · berapa disisir · berapa MASIH nol dibuka,
  dan dari golongan mana
- Berapa butir lahir · berapa ditutup (harus NOL)
- Apakah frontier masih terbuka
- Apa syarat ronde 3 produktif, atau nyatakan bila modul ini sudah siap masuk to-spec begitu
  OQ-025 dan naskah tiga procedure tersedia
```

---

## Menjalankan sambil mengukur biaya

`jq` tidak terpasang. `node` tidak terpasang. `py` ada (Python 3.14.7).

```bash
cd /d/XML/RNM_BRD/OUTPUT_HASIL_RNM
export PYTHONIOENCODING=utf-8

run_json=".run-nb-treaty-in-ronde-2.json"
claude --print --output-format json -- "$(cat PROMPT-GRILL-RONDE-2-NB-TREATY-IN.md)" > "$run_json"

py -c "
import json
d = json.load(open('$run_json', encoding='utf-8'))
print('biaya USD     :', d.get('total_cost_usd'))
print('durasi ms     :', d.get('duration_ms'))
print('jumlah giliran:', d.get('num_turns'))
print('session_id    :', d.get('session_id'))
for k, v in (d.get('usage') or {}).items():
    print(f'  {k}: {v}')
"
```

## Perkiraan biaya

Dasar: angka nyata ronde 1 — keluaran 192.021 · cache tulis 363.795 · cache baca 13.532.360 ·
69 panggilan.

Ronde 2 menyentuh 46,6 % modul yang belum dibaca, ditambah 30 rule hitung uang dan lima berkas
Activity terbesar. Perkiraan **1,5 sampai 2 kali ronde 1**. Bila angka sesungguhnya jauh melampaui
itu, sasaran 4 adalah yang pertama dipotong — dan pemotongannya harus dilaporkan, bukan disembunyikan.
