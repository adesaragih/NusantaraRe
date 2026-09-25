# PROMPT — Grilling Ronde 1: NB Treaty In

> Ketik sendiri di sesi kerja baru. `grill-with-docs` ber-`disable-model-invocation: true`,
> agent tidak dapat memanggilnya (CLAUDE.md §8).

---

```
/mattpocock-skills:grill-with-docs

## Lingkup

Konteks: Treaty Inward — Realisasi & Endorsement.
Modul yang digrilling ronde ini: HANYA "NB Treaty In" (278 berkas .xml) di D:\XML\RNM_BRD\.

⛔ JANGAN membuka "EDM Treaty In" di ronde ini, walau ia satu bounded context dengan NB Treaty In.
Work owner memutuskan modul dikerjakan satu per satu. EDM Treaty In menyusul di ronde terpisah.
⛔ JANGAN membuka folder korpus modul lain mana pun.
⛔ JANGAN membuka D:\XML\nusantara-re\ — di-blacklist total.
⛔ JANGAN membaca isi .scratch\ milik konteks lain.

## Baca lebih dulu

- discovery/context-map.md §2.3 (konteks ini, cakupan `full`)
- discovery/flows/NB Treaty In.md — TERUTAMA §5 Batas pengetahuan dan §7 Pertanyaan terbuka baru
- discovery/modules/NB Treaty In.md
- discovery/inventory/NB Treaty In.md
- discovery/open-questions.md — hanya OQ yang menyentuh modul ini
- CONTEXT.md
- docs/adr/ — ADR-0001 sampai ADR-0015

Titik masuk sudah terverifikasi, jangan dicari ulang:
NB Treaty In/Flow/InputRealizationTreatyIn.xml
identitas ASM-FW-GISFW-WORK!INPUTREALIZATIONTREATYIN, pyStartActivity = Start1.

Isi modul per tipe rule (D1): Activity 92 · When 75 · RDBList 41 · Section 25 · ReportDefinition 15
· DataTransform 12 · FlowAction 9 · Harness 6 · Flow 1.

## Keluaran

Tulis HANYA ke .scratch/nb-treaty-in/grilling-ronde-1.md. SATU berkas. Tidak ada berkas lain.
⛔ JANGAN membuat spec.md di ronde ini. ⛔ JANGAN membuat tiket.

## Aturan yang mengikat

CLAUDE.md §4 dan §4a berlaku penuh:

1. Setiap klaim perilaku wajib bukti path + class + nama rule. Sebut berkas yang dibaca.
2. Identitas rule = class / nama dari <pxInsName>. <pzOriginalInstanceKey> adalah asal salinan,
   BUKAN identitas.
3. Label wajib di setiap pernyataan: [terverifikasi] / [dugaan] / [terbuka].
4. Setiap angka disertai perintah audit yang menghasilkannya.
5. Tiap sensus dihitung DUA cara yang benar-benar berbeda, dan jendelanya disebut. Bila kedua cara
   berselisih, tulis "belum punya data" dan cantumkan keduanya — jangan memilih salah satu.
6. Properti korpus dicari dengan awalan halaman: InputData.CARI21, bukan CARI21 telanjang.
7. Grep dulu, baca rentang seperlunya. Catat bila sebuah berkas belum habis dibaca.
8. Jangan menebak skema, tipe kolom, arti kode/status/enumerasi, kepanjangan singkatan, rumus,
   isi stored procedure, atau struktur JSON.
9. Jangan menyatukan identitas yang berkonflik. Cek _oq011-konflik-isi.md lebih dulu.
10. Nilai nama orang tidak disalin. Catat rule, tag, dan jumlahnya saja.
11. Nomor OQ diambil dari register discovery/open-questions.md, bukan dari ingatan.
12. ⛔ Agent TIDAK BOLEH menutup butir terbuka sendiri. Penutupan hanya oleh pemilik peran
    (work owner / DBA / Product+Underwriting / Actuarial / Finance / IAM).
13. ⛔ Kode NOL. DDL NOL. CREATE TABLE NOL. Daftar kolom usulan NOL.

## Empat batas pengetahuan yang SUDAH diketahui — konfirmasi, jangan temukan ulang

discovery/flows/NB Treaty In.md sudah merekam empat hal berikut. Tugas ronde ini adalah
mengkonfirmasi apakah masih berlaku, dan memperdalamnya — bukan menemukannya lagi dari nol:

- §5.1 — RDBList/SavePolisTreatyIn_SQL.xml memanggil POOLDATA.PEGA_JSON_POLIS_TREATYIN dengan
  8 parameter, diikuti COMMIT di dalam blok PL/SQL. Isi procedure tidak diketahui (OQ-002),
  batas transaksi di sisi database (OQ-013).
- §5.2 — BLOCKER. Activity/serviceInsertArasapas_act.xml memanggil identitas yang variannya
  TIDAK ADA di NB Treaty In (register OQ-011 entri #415, 2 varian, 2 isi berbeda). Telusur
  BERHENTI di sini per _METHOD.md §1.1 aturan 4. Varian dari modul lain TIDAK DIPINJAM. → OQ-025.
- §5.3 — @ASM.GetPageJSONString() menghasilkan JSON yang ditulis ke database. Source fungsi ini
  tidak ada di korpus. → memperkuat OQ-012.
- §5.4 — isApproved ada sebagai DUA tipe rule dengan class+nama sama (RULE-OBJ-WHEN dan
  RULE-DECLARE-DECISIONTABLE). Enam shape Decision merujuknya; mana yang dipakai tidak dapat
  ditentukan dari flow. → OQ-026.

## Bab yang WAJIB ada di keluaran

### Bab A — Stored procedure dan SQL mentah

Daftar lengkap stored procedure / function Oracle yang dipanggil NB Treaty In. Per baris:
nama procedure · rule pemanggil (class + nama + path) · jumlah parameter · apakah ia COMMIT sendiri
· apakah body-nya sudah ada di artefak [data DBA] mana pun di OUTPUT_HASIL_RNM, atau belum.

Alasan bab ini wajib: work owner sedang memutuskan apakah pemanggilan stored procedure dihapus
seluruhnya. Keputusan itu butuh daftar lengkap per modul, dan butuh tahu body mana yang sudah
diserahkan DBA. Tanpa body, memindahkan logikanya ke Go berarti menebak — dan menebak dilarang.

### Bab B — Tabel Oracle yang disentuh

Per tabel: nama · schema · dibaca atau ditulis atau keduanya · rule yang menyentuhnya.
Tandai tabel yang HANYA disentuh lewat stored procedure, karena untuk tabel itu struktur kolomnya
tidak dapat dinyatakan.

### Bab C — Register [terbuka] berpemilik

Tabel: nomor · butir · pemilik · kenapa ia memblokir · apa akibatnya kalau dijawab salah.
Pemilik dipilih dari: [work owner] · [DBA] · [Product+Underwriting] · [Finance] · [IAM] ·
[pengembang Pega lama] · [pemilik export Pega].
Tutup bab ini dengan rekap jumlah per pemilik.

### Bab D — PERTANYAAN SIAP KIRIM ⭐ wajib, ini yang paling penting

Ubah Bab C menjadi dokumen yang bisa LANGSUNG dibawa ke rapat dan dijawab tanpa penjelasan
tambahan dari agent. Kelompokkan per pemilik. Untuk setiap pertanyaan, tulis:

1. Pertanyaannya dalam bahasa bisnis — TANPA istilah Pega, TANPA nama rule, TANPA path.
   Orang yang menjawab tidak membaca XML.
2. Konteks satu kalimat: kenapa pertanyaan ini muncul.
3. Bentuk jawaban yang diharapkan: pilihan ganda bila ada, contoh nilai, atau "butuh berkas X".
4. Dampak bila salah dijawab, dalam satu kalimat.
5. Rujukan teknis diletakkan di baris terpisah paling bawah, ditandai "rujukan:", supaya penjawab
   boleh mengabaikannya.

Urutkan dari yang paling memblokir. Beri nomor P1, P2, P3, dan seterusnya.
⛔ JANGAN menjawab sendiri satu pun pertanyaan di bab ini.

### Bab E — Lampiran bukti isolasi

Tabel yang menunjukkan: folder korpus selain NB Treaty In = NOL dibuka · isi .scratch selain
nb-treaty-in = NOL dibaca · D:\XML\nusantara-re\ = NOL disentuh · berkas dibuat = tepat SATU ·
korpus tidak berubah. Sertakan perintah audit untuk tiap baris.

## Tutup dengan

Rekap: berapa butir [terbuka] lahir di ronde ini · berapa ditutup (harus NOL, agent tidak boleh
menutup) · apakah frontier masih terbuka · apa syarat ronde 2 bisa produktif.
```

---

## Cara menjalankan sambil mengukur token dan biaya

`jq` TIDAK terpasang di mesin ini (`command -v jq` → kosong). `node` juga tidak. `py` ada
(Python 3.14.7). Jadi pakai `py`, bukan `jq`.

```bash
cd /d/XML/RNM_BRD/OUTPUT_HASIL_RNM
export PYTHONIOENCODING=utf-8

prompt_file="PROMPT-GRILL-RONDE-1-NB-TREATY-IN.md"
run_json=".run-nb-treaty-in-ronde-1.json"

claude --print --output-format json -- "$(cat "$prompt_file")" > "$run_json"

py -c "
import json
d = json.load(open('$run_json', encoding='utf-8'))
print('biaya USD   :', d.get('total_cost_usd'))
print('durasi ms   :', d.get('duration_ms'))
print('jumlah giliran:', d.get('num_turns'))
print('session_id  :', d.get('session_id'))
print('token:')
for k, v in (d.get('usage') or {}).items():
    print(f'  {k}: {v}')
"
```

`total_cost_usd` dan `usage` adalah angka sebenarnya dari server, bukan proksi counter dalam sesi.

## ⚠️ Uji dulu sebelum menjalankan yang panjang

Ronde ini diperkirakan 250.000–400.000 token. Jangan tembak langsung. Uji dulu apakah
`--print` menerima slash command sebagai prompt — `grill-with-docs` ber-`disable-model-invocation:
true`, dan belum terbukti apakah larangan itu juga berlaku saat slash command dikirim lewat
`--print`:

```bash
claude --print --output-format json -- "/mattpocock-skills:grill-with-docs sebutkan saja apakah skill ini aktif. jangan membaca berkas apa pun. jangan menulis berkas apa pun." > uji.json
py -c "import json;d=json.load(open('uji.json',encoding='utf-8'));print(d.get('total_cost_usd'));print(d.get('result','')[:400])"
```

Bila uji itu menolak, jalur `--print` gugur untuk skill ini. Ketik slash command-nya sendiri di
sesi interaktif, dan ambil angka tokennya dengan `/caveman-stats` atau laporan pemakaian bawaan.

## Catatan operasional

- `--print` berjalan headless satu tembakan. Tidak bisa disela di tengah jalan. Untuk pekerjaan
  250–400 ribu token, itu risiko: kalau arahnya melenceng di menit kelima, Anda baru tahu di akhir.
  Sesi interaktif memberi kesempatan menghentikan.
- Berkas `.run-*.json` berisi `session_id`. Simpan — ia memudahkan melacak ulang jalannya sesi.
