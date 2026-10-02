# FASE A DISCOVERY — LAPORAN PENUTUP (D3 + D4)

Migrasi Pega → Go + React, Nusantara Re. Ditutup **2026-09-14**.
Korpus `D:\XML\RNM_BRD\` — **9.369 file, 20 modul**, READ-ONLY sepanjang fase.

> **Catatan pelaksanaan.** Fase ini **tidak** dijalankan lewat skill `/wayfinder`. Skill Matt Pocock
> yang terpasang (`mattpocock-skills@1.2.3`) ber-frontmatter `disable-model-invocation: true`
> sehingga **hanya dapat dipicu manusia**. Atas keputusan manusia, D3+D4 dikerjakan sebagai
> **sintesis metodologi FASE A proyek ini** — format `modules/*.md`, `glossary.md`,
> `context-map.md`, `understanding-report.md`, sama seperti D0–D2. Alur internal skill wayfinder
> (peta decision-ticket di issue tracker) **tidak direplikasi**. Cara menjalankan skill FASE B
> secara manual ada di **§6**.

---

## 1. Apa yang dihasilkan FASE A

| Step | Keluaran | Status |
| --- | --- | --- |
| **D0** | kerangka discovery + fakta korpus terukur | selesai |
| **D1** | `inventory/<modul>.md` ×20 + `_summary.md` + `_oq011-konflik-isi.md` | selesai — 20/20 modul, 9.369/9.369 file (100 %) |
| **D2** | `flows/<konteks>.md` ×20 + 4 sintesis + 2 dokumen metode | selesai — **20 konteks / 428 rule**, 5 tahap |
| **D3** | `modules/<modul>.md` ×20 + **`glossary.md`** (169 entri) | **selesai** |
| **D4** | **`context-map.md`** + **`understanding-report.md`** | **selesai** |

Penutup per step: `D1-CLOSING-REPORT.md`, `D2-CLOSING-REPORT.md`, dan dokumen ini.

`[terverifikasi]` Angka kunci:

| Ukuran | Nilai | Perintah audit |
| --- | ---: | --- |
| File korpus | 9.369 | `find . -path ./OUTPUT_HASIL_RNM -prune -o -name '*.xml' -print \| wc -l` |
| Konteks ditelusur | 20 | `ls discovery/flows/*.md \| grep -vc "^.*_"` |
| Rule ditelusur | 428 | jumlahkan `**N rule ditelusur**` di `flows/*.md` |
| Catatan modul | 20 | `ls discovery/modules/*.md \| grep -vc "_README"` |
| Entri glossary | 169 | lihat `glossary.md` §Cakupan terukur |
| Bounded context | 9 + 1 absent | `context-map.md` §1 |
| Pertanyaan terbuka | **57 terbuka / 2 terjawab** | `awk` di `D2-CLOSING-REPORT.md` §4 |

---

## 2. GATE — yang harus manusia review dan putuskan

**FASE B (ADR/BRD/spec/tiket) tidak boleh dimulai sebelum gate ini lewat.**

### 2.1 Yang harus di-review

| Dokumen | Apa yang diperiksa |
| --- | --- |
| **`understanding-report.md`** | Apakah narasi "bagaimana aplikasi bekerja" cocok dengan yang Anda ketahui? Yang **salah** lebih penting dilaporkan daripada yang benar. |
| **`context-map.md`** | Apakah sembilan bounded context masuk akal secara bisnis? Apakah label `full`/`partial` dapat diterima? |
| **`glossary.md`** | Apakah ada istilah yang **salah dicatat**? (Bukan: apakah ada yang belum lengkap — itu memang disengaja.) |
| **`open-questions.md`** | 57 pertanyaan — mana yang sudah bisa Anda jawab sekarang? |

### 2.2 Tujuh keputusan yang **tidak dapat diturunkan dari korpus**

`[terverifikasi]` Bukti sudah lengkap untuk kelima yang pertama; keputusannya milik manusia.

| # | Keputusan | Bukti yang tersedia | OQ |
| ---: | --- | --- | --- |
| 1 | Facultative NB/RNW/EDM → **satu** konteks atau **tiga**? | satu ruleset; **99,0 %** identitas `RNW` ⊂ `NB`; diskriminator `IsNB`/`IsRenewal`/`IsEDM` ada di ketiganya | **OQ-015** |
| 2 | `Treaty In` + `Adjustment` → **satu** atau **dua**? | **277 dari 323** file bernama sama identik (85,8 %); mesin status `Akseptasi_DT` satu rule (`58b8e650`) | **OQ-010** |
| 3 | Ke mana **`Treaty Contract Out`** (303 rule) ditempatkan? | terbukti **bukan outward** — master arrangement inward (5 uji) | **OQ-022** |
| 4 | Trio Claim non-life → **tiga** konteks atau satu bervarian? | overlap 12–40 %; **tanpa** diskriminator runtime | **OQ-019** |
| 5 | Siapa **pemilik** data lintas konteks? | Claim membaca master Treaty; Master Life membawa penulis `M_TREATY_IN` | **OQ-042**, **OQ-056** |
| 6 | **Model RBAC** — harus dirancang dari nol | nol rule peran di korpus; 4 pola otorisasi tersebar | **OQ-007** dkk |
| 7 | **Versi rule mana yang berlaku di production** | 533 identitas berkonflik; 429 tetap berbeda setelah normalisasi 21 tag | **OQ-011** |

### 2.3 Risiko yang harus diakui sebelum FASE B

Rincian berbukti di `understanding-report.md` §5:

1. **Korpus belum tentu production** — hostname DEV (81 file/15 modul), 2 `pxHostId`, **10 tombol
   `(dev)`** termasuk `Force Resolve Complete`, `DBMS_OUTPUT.PUT_LINE` di SQL, folder
   `Claude outputs`.
2. **Rumus bisnis di luar korpus** — 67 SP tanpa body; pada 2 modul **seluruh jalur tulis master
   lewat SP**; 8 tabel rate di database lain.
3. **RBAC harus dirancang** — nol rule peran; 5 identitas orang ter-hardcode **sebagai tujuan rute**.
4. **Uang material tanpa mata uang** — `0.45`/`0.05`, `"3000000000"` **dibandingkan sebagai string**,
   `30000000.00`/`50000000.00`/`57750000.00`. Aturan proyek: **jangan `float`**.
5. **Identitas & host ter-hardcode** — endpoint sesungguhnya di tabel `M_LINK_SERVICE`.
6. **Model AI pihak ketiga di alur underwriting** — `GeminiAIGoogle_Act`, isinya belum dibaca.
7. **Ketergantungan terputus** — `SERVICEINSERTARASAPAS_ACT` dipanggil 3 modul, tidak ada di
   ketiganya.

---

## 3. Seluruh 57 pertanyaan terbuka, per pemilik peran

Tanda **⛔** = **pemblokir FASE B** (38 dari 57). Tanda ▫ = tidak memblokir FASE B (memblokir D4
atau bersifat higienis).

Pemilik utama (yang pertama disebut di register): **P+UW 28, DBA 11, IAM 10, Finance 4,
pemilik export Pega 3, Arsitektur Pega 1**.

### 3.1 DBA — 11 sebagai pemilik utama, 14 tersentuh

| | OQ | Ringkas |
| --- | --- | --- |
| ⛔ | **OQ-001** | Tidak ada DDL Oracle maupun definisi properti di korpus |
| ⛔ | **OQ-002** | Body **67 stored procedure** tidak ada di korpus |
| ⛔ | **OQ-008** | Arti prefix `ASM` / `RNM` / `GCNM` pada kunci RDBList |
| ⛔ | **OQ-012** | Struktur JSON di kolom `JSONDATA` / `DATA_JSON` (bersama P+UW) |
| ⛔ | **OQ-013** | Batas transaksi & identitas pengguna di dalam stored procedure |
| ▫ | OQ-016 | Peran 10 skema Oracle; mana dalam ruang lingkup (bersama P+UW) |
| ⛔ | **OQ-017** | Db-link `ASMD.SINARMAS.CO.ID` — 8 tabel rate & kurs di luar sistem (bersama Actuarial, P+UW) |
| ⛔ | **OQ-018** | Host/URL ter-hardcode termasuk hostname DEV (bersama IAM) |
| ⛔ | **OQ-029** | `IsPEGAPROD` dkk — 5 keluarga `When` tak terbaca (bersama IAM, P+UW) |
| ⛔ | **OQ-047** | Daftar endpoint di tabel `M_LINK_SERVICE` (bersama Platform) |
| ⛔ | **OQ-055** | Nomor revisi di dalam string ID, ter-hardcode di rule **dan** SQL (bersama P+UW) |
| ⛔ | **OQ-058** | Tabel internal Pega `DATAPEGA.PC_*` diakses langsung lewat SQL (bersama Arsitektur Pega) |
| ⛔ | **OQ-059** | Slot parameter generik `CARI1`…`CARI30` menuju SQL (bersama P+UW) |
| ▫ | OQ-056 | Master produk life membawa jalur tulis `M_TREATY_IN` (bersama P+UW) |

### 3.2 Product + Underwriting — 28 sebagai pemilik utama, 44 tersentuh

| | OQ | Ringkas |
| --- | --- | --- |
| ▫ | OQ-003 | Folder `excludeXML` — aktif atau tidak |
| ▫ | OQ-005 | Lima modul tanpa `Flow` — **terjawab untuk keperluan D2** |
| ▫ | OQ-009 | Rule di `@BASECLASS` (202/303 di `Treaty Contract Out`) |
| ▫ | OQ-010 | `Treaty In` ↔ `Adjustment`: satu atau dua |
| ⛔ | **OQ-011** | **533 identitas berkonflik** — versi mana yang berlaku di production |
| ⛔ | **OQ-014** | Arti singkatan `EDM` |
| ▫ | OQ-015 | Facultative: satu ruleset atau tiga — **terjawab sebagian** |
| ▫ | OQ-019 | Trio Claim tidak berbagi basis kode — apa pembedanya |
| ⛔ | **OQ-020** | **Arti seluruh kode/enumerasi** (bersama Finance) |
| ▫ | OQ-022 | `Treaty Contract Out` bukan outward — **terjawab sebagian** |
| ▫ | OQ-023 | Shape/tombol tanpa jalur masuk; guard `FALSE &&` |
| ⛔ | **OQ-025** | `SERVICEINSERTARASAPAS_ACT` dipanggil 3 modul, tidak ada di ketiganya |
| ⛔ | **OQ-026** | Satu nama rule, dua tipe — mana yang dipakai flow |
| ▫ | OQ-028 | `WorkList` vs `WorkBasket` (bersama IAM) |
| ⛔ | **OQ-030** | Ambang tanggal 25 (bersama Finance) |
| ▫ | OQ-031 | Empat ID `1000032`–`1000035` |
| ⛔ | **OQ-032** | Sumber nilai `.KomiteLoop` |
| ▫ | OQ-033 | `KomiteRouter` hanya menetapkan sasaran bila `TransferType == '2'` |
| ▫ | OQ-034 | Tanggal cutover ter-hardcode `20250207` |
| ⛔ | **OQ-036** | Empat sasaran routing komite ter-hardcode (bersama IAM) |
| ⛔ | **OQ-037** | Ambang nominal menentukan komposisi roster komite (bersama Finance) |
| ▫ | OQ-038 | Sebelas kode lini bisnis `L1`…`L11` |
| ▫ | OQ-039 | Adjustment / Close / Reject tidak muncul di graf |
| ⛔ | **OQ-040** | Limit wewenang: hardcode vs database (bersama Finance) |
| ⛔ | **OQ-041** | `ISCLM*` berkonflik **dan** tidak terbaca |
| ▫ | OQ-042 | Klaim membaca master treaty outward |
| ⛔ | **OQ-043** | **Baris tabel keputusan tidak terekspor** (49 `DecisionTable`) |
| ⛔ | **OQ-044** | `ProposalAcceptStatus = 4` ambigu |
| ⛔ | **OQ-046** | Pangsa retro 0,45/0,05 & ambang 3 miliar (bersama Actuarial) |
| ⛔ | **OQ-048** | `GeminiAIGoogle_Act` — AI pihak ketiga di alur underwriting (bersama Security) |
| ⛔ | **OQ-050** | 10 tombol `(dev)` termasuk `Force Resolve Complete` |
| ⛔ | **OQ-052** | `Reject` vs `Decline` — bedanya tidak dijelaskan |
| ⛔ | **OQ-057** | Kode `OR` & `REINSTYPEID` tanpa nilai literal |
| ⛔ | **OQ-021**, **OQ-024**, **OQ-045**, **OQ-053** | lihat §3.4 (pemilik utama IAM) |

### 3.3 Actuarial — 2 tersentuh

| | OQ | Ringkas |
| --- | --- | --- |
| ⛔ | **OQ-017** | Tabel rate & kurs di luar sistem (bersama DBA, P+UW) |
| ⛔ | **OQ-046** | Pangsa retro & ambang ter-hardcode (bersama P+UW) |

### 3.4 IAM — 10 sebagai pemilik utama, 11 tersentuh

| | OQ | Ringkas |
| --- | --- | --- |
| ⛔ | **OQ-007** | **Tidak ada rule identitas/otorisasi di korpus** |
| ⛔ | **OQ-018** | Host/URL ter-hardcode (bersama DBA) |
| ⛔ | **OQ-021** | Identitas orang ter-hardcode sebagai **guard** (bersama P+UW) |
| ⛔ | **OQ-024** | Pemetaan Assignment → workbasket tidak terbaca (bersama P+UW) |
| ⛔ | **OQ-027** | `OperatorID.pyTelephone` menyimpan kode peran — **4 nilai** |
| ▫ | OQ-028 | `WorkList` vs `WorkBasket` (bersama P+UW) |
| ⛔ | **OQ-029** | 5 keluarga `When` tak terbaca (bersama DBA, P+UW) |
| ⛔ | **OQ-036** | Empat sasaran routing komite ter-hardcode (bersama P+UW) |
| ⛔ | **OQ-045** | `LetterNo` sebagai token routing persetujuan (bersama P+UW) |
| ⛔ | **OQ-051** | Otorisasi lewat **indeks tetap** `pyWorkBasketList(2)` |
| ⛔ | **OQ-053** | **5 identitas orang ditetapkan sebagai pemilik tugas berikutnya** (bersama P+UW) |

### 3.5 Finance — 4 tersentuh

| | OQ | Ringkas |
| --- | --- | --- |
| ⛔ | **OQ-020** | Arti kode `PaymentType`, `TransferType`, dan seluruh enumerasi (bersama P+UW) |
| ⛔ | **OQ-030** | Ambang tanggal 25 (bersama P+UW) |
| ⛔ | **OQ-037** | Ambang nominal menentukan roster komite (bersama P+UW) |
| ⛔ | **OQ-040** | Limit wewenang: hardcode vs database (bersama P+UW) |

### 3.6 Pemilik export Pega / Arsitektur Pega / Platform / Security

| | OQ | Ringkas | Pemilik |
| --- | --- | --- | --- |
| ⛔ | **OQ-035** | Activity dipanggil tetapi salinannya hanya di modul lain | pemilik export |
| ⛔ | **OQ-043** | Baris tabel keputusan tidak ikut terekspor | pemilik export + P+UW |
| ▫ | OQ-049 | Empat class bersufiks `ENDORSEMENT`, satu berprefix `ASM-SFAGIS-` | Arsitektur Pega + pemilik export |
| ▫ | OQ-054 | Folder `Claude outputs` berisi berkas non-Pega di dalam ekspor | pemilik export |
| ⛔ | **OQ-047** | `M_LINK_SERVICE` | Platform + DBA |
| ⛔ | **OQ-048** | Gemini AI di alur underwriting | Security + P+UW |
| ⛔ | **OQ-058** | Tabel internal Pega diakses langsung | Arsitektur Pega + DBA |

---

## 4. Usulan urutan FASE B — berbasis bukti

`[terverifikasi]` Diukur per bounded context: jumlah OQ yang menyentuhnya, berapa di antaranya
**pemblokir FASE B**, dan cakupan bukti dari `context-map.md`.

| Urutan | Bounded context | File | Cakupan | OQ total | **Pemblokir** |
| ---: | --- | ---: | --- | ---: | ---: |
| **1** | **Claim — Life** | 136 | **full** | 9 | **5** |
| **2** | **Life — Penawaran & Premium List** | 199 | **full** | 13 | **9** |
| **3** | **Treaty Arrangement** (`Treaty Contract Out`) | 303 | partial | 11 | **5** |
| 4 | Claim — Non-Life | 1.031 | partial | 14 | 8 |
| 5 | Life — Master | 180 | partial | 15 | 10 |
| 6 | Komite | 300 | **full** | 16 | 13 |
| 7 | Treaty Inward — Realisasi & Endorsement | 441 | **full** | 20 | 15 |
| 8 | Treaty Inward — Master & Akseptasi | 708 | partial | 24 | 18 |
| **9** | **Facultative Inward** | **6.071** | partial | 26 | **20** |

Perintah audit (dapat dijalankan ulang):
```
for m in <modul konteks>; do awk '/^## 7\./{f=1} f' "modules/$m.md" | grep -oE "OQ-[0-9]{3}"; done \
 | sort -u | tee /dev/stderr | grep -cFf blockers.txt
```

### 4.1 Yang paling siap dimulai

**`Claim — Life`** (rekomendasi pertama). Alasannya berbukti:

- **cakupan `full`** — alur Register → Outstanding → Medical Check → Claim Analis terekam penuh;
- **pemblokir paling sedikit (5)**: OQ-002, OQ-018, OQ-020, OQ-021, OQ-029;
- **modul terkecil dan paling terpisah** — 136 file, overlap identitas hanya 12–14 % terhadap modul
  Claim lain, sehingga keputusan di sini **tidak menyeret konteks lain**;
- hanya **3 stored procedure** (paling sedikit di domain Claim);
- **tidak memakai `PaymentType`**, sehingga OQ-020 berdampak lebih sempit daripada di modul lain.

**`Life — Penawaran & Premium List`** (rekomendasi kedua) — juga `full`, 199 file, dan berbagi mesin
dengan konteks pertama (`InsertJsonPolisLife_Act`, `IsLifeAccepted`), sehingga jawaban OQ yang sama
dipakai dua kali.

**`Treaty Arrangement`** (rekomendasi ketiga) — hanya 5 pemblokir dan batas konteksnya sudah jelas
(OQ-022 terjawab: bukan outward). Peringatan: cakupan `partial` karena **seluruh jalur tulis lewat
stored procedure** dan **202/303 rule di `@BASECLASS`**.

### 4.2 Yang harus menunggu jawaban bisnis

**`Facultative Inward`** — **jangan dimulai lebih dulu** meski ia domain terbesar (6.071 file,
64,8 % korpus). Alasannya:

- **20 pemblokir** — terbanyak;
- **rumus spreading/capacity/scoring belum terbaca** (34–38 activity, terbesar 996 KB);
- **pemetaan `ProposalAcceptStatus` → hasil tidak ada di korpus** (OQ-043) — inti aturan akseptasi;
- **input rating di 8 tabel db-link di luar sistem** (OQ-017);
- keputusan satu-vs-tiga konteks (OQ-015) belum diambil.

**`Treaty Inward — Master & Akseptasi`** (18 pemblokir) — menunggu OQ-010, OQ-052 (beda
`Reject`/`Decline`), OQ-053 (5 identitas ter-hardcode), OQ-055.

**`Komite`** (13 pemblokir) — cakupan `full`, tetapi hampir seluruh pemblokirnya adalah **RBAC dan
ambang nominal** (OQ-021, OQ-024, OQ-036, OQ-037, OQ-040). Konteks ini baru produktif **setelah
model peran diputuskan**.

### 4.3 Empat OQ yang membuka paling banyak konteks sekaligus

`[terverifikasi]` Bila hanya sedikit waktu narasumber yang tersedia, empat ini memberi hasil
terbesar:

| OQ | Menyentuh | Mengapa |
| --- | ---: | --- |
| **OQ-020** (arti kode/enumerasi) | **9 dari 9 konteks** | mengunci 51 entri glossary Bagian 2 |
| **OQ-002** (body stored procedure) | **9 dari 9 konteks** | 67 SP; pada 2 modul seluruh jalur tulis |
| **OQ-021** (+ OQ-007, OQ-024, OQ-045, OQ-051, OQ-053) | 8 konteks | model RBAC — tidak ada yang bisa disalin |
| **OQ-011** (versi rule production) | 7 konteks | 533 identitas berkonflik; risiko halusinasi paling konkret |

---

## 5. Kesiapan FASE B per keluaran

| Keluaran FASE B | Bahan yang sudah ada | Yang masih kurang |
| --- | --- | --- |
| **ADR** (keputusan arsitektur) | `context-map.md` §5 — 7 keputusan dengan buktinya | jawaban manusia atas ketujuhnya |
| **BRD** (kebutuhan bisnis) | `understanding-report.md` §2 — narasi per konteks | arti kode (OQ-020), aturan persetujuan (OQ-043) |
| **Spec** | `flows/*.md` — alur, guard, objek Oracle per konteks | rumus (OQ-002, OQ-017), struktur JSON (OQ-012), tipe kolom (OQ-001) |
| **Tiket** | `modules/*.md` — cakupan per modul | penetapan bounded context final (OQ-010/015/019/022) |
| **`CONTEXT.md`** | **`glossary.md`** — 169 entri berkolom Bukti | OQ-008, OQ-014, OQ-020, OQ-057 |

---

## 6. CARA MENJALANKAN FASE B SECARA MANUAL

`[terverifikasi]` Skill FASE B **tidak dapat dipanggil oleh agent**. Frontmatter-nya
`disable-model-invocation: true`, dan percobaan pemanggilan ditolak dengan pesan:
*"Ask the user to run /mattpocock-skills:<nama> themselves — it cannot be invoked via the Skill
tool."*

Audit:
```
find "C:/Users/Administrator/.claude/plugins/cache/claude-plugins-official/mattpocock-skills/1.2.3/skills" \
  -name SKILL.md | while read f; do echo "$(basename $(dirname $f)): $(grep -i '^disable-model-invocation' $f)"; done
```

### 6.1 Command yang harus **Anda** ketik sendiri

Plugin sudah terpasang (`mattpocock-skills@claude-plugins-official` v1.2.3, scope user,
`installedAt 2026-08-31`) — **tidak perlu setup ulang**.

| Tahap FASE B | Command |
| --- | --- |
| Grilling konteks terpilih | `/mattpocock-skills:grill-with-docs` |
| Menyusun spec | `/mattpocock-skills:to-spec` |
| Memecah jadi tiket | `/mattpocock-skills:to-tickets` |

Bila autocomplete tidak memunculkannya, coba nama pendek (`/grill-with-docs`). Ketik `/` lalu gulir
daftar untuk memastikan ejaannya. Bila daftar kosong: `claude plugin list` untuk memastikan plugin
aktif, lalu restart sesi.

### 6.2 Prasyarat yang perlu Anda putuskan lebih dulu

**(a) Issue tracker.** Skill keluarga ini bekerja dengan tiket di issue tracker repo.
`D:\XML\RNM_BRD` **bukan repository git** dan tidak punya issue tracker.

Audit: `git -C "D:/XML/RNM_BRD" status` → bukan repo.

Opsi untuk Anda pilih — **saya tidak memutuskan**:

| Opsi | Bentuk | Catatan |
| --- | --- | --- |
| **A** | **Markdown lokal** di `OUTPUT_HASIL_RNM/.scratch/<konteks>/` | paling dekat dengan cara kerja proyek ini sejauh ini; tetap READ-ONLY terhadap korpus |
| B | Inisialisasi git lokal + issue tracker ringan | menambah perkakas baru di tengah fase |
| C | Tracker perusahaan (mis. GitLab sesuai runbook `CLAUDE-REST-API-GITLAB-HEADLESS-RUNBOOK-ID.md`) | perlu kredensial & persetujuan; **jangan taruh dump production atau data pelanggan** |

**(b) Jawaban OQ pemblokir** untuk konteks yang dipilih. Untuk rekomendasi pertama
(**Claim — Life**) yang diperlukan hanya **5**: OQ-002, OQ-018, OQ-020, OQ-021, OQ-029.

**(c) Bahan yang akan diberikan ke skill.** Untuk `Claim — Life`:
`modules/Claim Life.md`, `flows/Claim Life.md`, `flows/_SUMMARY-claim.md`,
`context-map.md` §2.6, `understanding-report.md` §2.4, `glossary.md`, `open-questions.md`.

### 6.3 Batas yang tetap berlaku

- Korpus `D:\XML\RNM_BRD\` dan `D:\XML\nusantara-re\` tetap **READ-ONLY**; tulis hanya ke
  `OUTPUT_HASIL_RNM\`.
- **Jangan** memasukkan secret/token/data pelanggan/dump production ke prompt, kode, atau test.
- Endpoint & host internal = **konfigurasi/env var**, bukan literal.
- Uang material **jangan `float`**.
- Butuh persetujuan manusia sebelum: operasi destruktif, migration DB, perubahan authn/authz,
  akses credential/regulated data, perubahan sistem eksternal.

---

## 7. Pernyataan penutup FASE A

`[terverifikasi]` **FASE A Discovery (D0–D4) SELESAI.** Dua puluh modul terinventaris, dua puluh
konteks tertelusur end-to-end (428 rule), dua puluh catatan modul, glossary 169 entri,
context map sembilan konteks + satu absent, dan understanding report — seluruhnya berbukti
`path + rule`, dengan setiap angka disertai perintah audit.

**Fase ini sekarang menunggu GATE manusia.** FASE B tidak dimulai, dan tidak akan dimulai oleh
agent — skill-nya hanya dapat Anda picu sendiri (§6).

`[terverifikasi]` Korpus tidak tersentuh sepanjang D0–D4:
```
find . -path ./OUTPUT_HASIL_RNM -prune -o -type f -newermt "-1 day" -print   # kosong
find "D:/XML/nusantara-re" -type f -newermt "-1 day" -print                  # kosong
```
