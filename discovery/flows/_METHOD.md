# Konvensi Telusur STEP D2

Mengikat untuk **seluruh** konteks D2. Ditetapkan 2026-09-13, setelah D1 selesai (20/20 modul).

## 0. Batas fase

D2 **memahami dan merekam perilaku existing apa adanya**. D2 **tidak**:

- menilai benar/salah, efisien/tidak;
- merancang Go, React, atau skema database;
- membuat ADR, BRD, spec, atau tiket (itu FASE B, setelah gate D4);
- menyimpulkan *mengapa* sesuatu dilakukan — hanya *apa* yang dilakukan.

## 1. Empat aturan anti-halusinasi (dari `../D1-CLOSING-REPORT.md` §4)

### 1.1 OQ-011 — cek register sebelum membuka rule

Sebelum menelusur rule apa pun, cek **`../inventory/_oq011-konflik-isi.md`**. Bila identitas rule
itu terdaftar, ia punya **lebih dari satu isi** di korpus.

- Telusur **varian milik modul konteks yang sedang dikerjakan**.
- Catat eksplisit bahwa modul lain punya versi berbeda.
- **Jangan** menulis "rule X melakukan Y" tanpa menyebut varian mana.
- **Jangan** memilih satu varian sebagai "yang benar".
- Bila konteks bergantung pada rule berkonflik yang variannya **tidak ada** di modul itu sendiri:
  **berhenti**, catat sebagai blocker. Jangan meminjam varian dari modul lain.

### 1.2 OQ-018 — korpus belum tentu production

Korpus memuat hostname DEV (`appdev.nusantarare.com`, 7×) dan dirakit dari lebih dari satu server
Pega (dua `<pxHostId>`). **Setiap pernyataan perilaku menyebut FILE mana yang dibaca**, bukan hanya
nama rule.

### 1.3 OQ-002 / OQ-017 — logika di luar korpus = batas pengetahuan

67 stored procedure dan 8 objek database link isinya **tidak ada di korpus**. Pemanggilannya
dicatat dengan pola tetap:

> Rule `<class/nama/tipe>` (`<path>`) memanggil `POOLDATA.Y` dengan parameter `a`, `b`, `c`;
> hasilnya ditaruh di `<properti>` dan dipakai untuk `<langkah berikutnya>`.
> **Isi `POOLDATA.Y` tidak diketahui** (OQ-002).

**Dilarang** menuliskan procedure sebagai langkah yang dijelaskan ("menghitung premi", "memvalidasi
polis") — namanya bukan spesifikasinya.

### 1.4 Kode, status, dan identitas orang

- Nilai kode/status/enumerasi dicatat **literal apa adanya**. Artinya ditulis
  `arti belum terverifikasi` bila korpus tidak menjelaskannya (OQ-020).
- Guard berbasis identitas orang (OQ-021) dicatat sebagai **temuan RBAC**: rule, tag, dan efeknya.
  **Nilai nama orang tidak disalin** ke artefak D2.
- Nama workbasket/worklist **dicatat** — itu justru bukti peran yang paling konkret (OQ-007).

## 2. Cara menelusur

### 2.1 Titik masuk

Mulai dari rule `Flow` → `<pyStartActivity>`. Untuk 5 modul tanpa `Flow` (OQ-005), titik masuk
alternatif ditetapkan lebih dulu per konteks — bukan diimprovisasi.

### 2.2 Membaca graf Flow `[terverifikasi]`

Struktur Flow Pega di korpus ini terbaca dari tag berikut:

| Tag | Isi | Guna |
| --- | --- | --- |
| `<pyShapeType>` | `Data-MO-Activity-Assignment`, `Data-MO-Gateway-Decision`, `Data-MO-Activity-Utility`, `Data-MO-Connector-Transition` | jenis shape |
| `<pxSubscript>` | `Assignment2`, `Decision8`, `Utility1`, … | **ID shape** |
| `<pyMOName>` | mis. `IS SPV TREATY 1`, `HIT SERVICE ARASAPAS` | label shape (label, **bukan** perilaku) |
| `<pyFrom>` / `<pyTo>` | ID shape | **arah connector** |
| `<pyConditionType>` | `When`, `ALWAYS`, … | jenis guard connector |
| `<pyExpression>` / `<pyTaskStatusOrWhen>` | nama rule `When` atau `ALWAYS` | **guard** |
| `<pyImplementation>` | `WorkBasket` / `Worklist` | cara assignment dirutekan |
| `<pyRouteTo>` | mis. `Custom` | mekanisme routing |
| `<pyRuleName>` | nama Activity / FlowAction / When / workbasket yang dirujuk | **rujukan rule** |
| `<pyClassName>` | class konteks shape | class |

Perintah audit graf:

```
grep -o "<pyShapeType>[^<]*" "<flow>.xml" | sort | uniq -c
paste <(grep -o "<pyFrom>[^<]*" "<flow>.xml" | sed 's/<pyFrom>//') \
      <(grep -o "<pyTo>[^<]*"   "<flow>.xml" | sed 's/<pyTo>//')
grep -o "<pyRuleName>[^<]*" "<flow>.xml" | sed 's/<pyRuleName>//' | sort -u
```

**Catatan penting:** `<pyMOName>` adalah **label yang diketik orang**, bukan perilaku. Label
`HIT SERVICE ARASAPAS` tidak membuktikan apa pun tentang apa yang dilakukan shape itu — yang
membuktikan adalah rule yang dipanggilnya. Perlakukan label seperti nama rule: **tidak dapat
dipercaya sendirian** (D1 membuktikan nama menipu berulang kali).

### 2.3 Dari shape ke rule

| Dari | Baca | Catat |
| --- | --- | --- |
| **Activity** | `<pyStepsActivityName>` (method tiap step), `<pyStepsPreCondParamsWhen>` (precondition/guard), `<pyStepsJavaSource>` (Java tertanam) | **apa** yang dilakukan tiap step, berurutan. Java tertanam dirangkum, bukan disalin utuh |
| **When** | kondisi literal | kondisi apa adanya, termasuk nilai yang diuji |
| **RDBList** | `<pyBrowseSQL>` | tabel + kolom + operasi (SELECT/INSERT/UPDATE/DELETE). **Pertahankan `@`** (db-link). Tipe kolom **tidak** disimpulkan (OQ-001) |
| **Sub-flow** | Flow yang dipanggil | telusur terpisah, jangan diratakan |
| **ConnectREST** | `<pyServiceName>`, `<pyBaseURLSelectionType>`, `<pyBaseURLSetting>` | target sebagai **konfigurasi**; URL literal = temuan (OQ-018) |

### 2.4 UI (Harness / Section)

- `<pyInclude>` = **embed**, dependency nyata → dicatat.
- `<pyType>` bernilai `FIELD` / `LAYOUT` / `SUB_SECTION` → struktur layar.
- **`<pySection>` di aksi Refresh BUKAN dependency** → jangan dihitung.

### 2.5 Hemat konteks

File Flow/Harness/Section besar (terbesar 8,2 MB). **Grep dulu untuk menemukan posisi, lalu baca
rentang baris seperlunya.** Jangan membaca file utuh.

## 3. Format keluaran per konteks

Setiap `flows/<konteks>.md` memakai **7 bagian** ini, berurutan:

| § | Isi |
| --- | --- |
| **1** | **Diagram alur tekstual** — langkah → langkah, dengan guard di tiap cabang. Termasuk assignment ke workbasket/worklist |
| **2** | **Status/state yang berubah** — tabel: properti, nilai literal, rule yang mengubah, + `arti belum terverifikasi` bila perlu |
| **3** | **Objek Oracle yang disentuh** — tabel: objek, operasi, rule sumber, path |
| **4** | **Integrasi eksternal** — ConnectREST/service sebagai konfigurasi |
| **5** | **Batas pengetahuan** — SP, db-link, kolom JSON yang isinya tidak diketahui |
| **6** | **Daftar rule terlibat** — `class / nama / tipe` + path, ditandai bila ada di register OQ-011 |
| **7** | **Pertanyaan terbuka baru** — masuk juga ke `../open-questions.md` |

Istilah domain baru yang ditemukan → tambahkan ke `../glossary-draft.md` dengan kolom bukti.

## 4. Label wajib

Setiap pernyataan diberi `[terverifikasi]` (didukung tag yang dikutip), `[dugaan]` (dari pola/nama,
belum dikonfirmasi), atau `[pertanyaan terbuka]`.

Setiap angka disertai perintah yang menghasilkannya.
