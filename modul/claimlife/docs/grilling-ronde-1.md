# Grilling — Claim — Life — Ronde 1

Status: answered (work owner, 2026-09-14) — menunggu persetujuan frontier Ronde 2
Konteks: `claim-life` (Claim — Life)
Tanggal: 2026-09-14
Skill: `/mattpocock-skills:grill-with-docs` (grilling + domain-modeling)

> **Berkas ini adalah ARSIP PERTANYAAN, bukan jawaban dan bukan keputusan.**
> Tidak ada yang ditulis ke `CONTEXT.md` maupun `docs/adr/` atas dasar berkas ini.
> Keputusan baru dicatat setelah work owner mengisi kolom **Jawaban work owner**.
> Yang tidak dijawab **tetap menjadi OQ terbuka** di `discovery/open-questions.md` — tidak ditebak.

## Bahan bukti (READ-ONLY)

`discovery/modules/Claim Life.md`, `discovery/flows/Claim Life.md`,
`discovery/flows/_SUMMARY-claim.md`, `discovery/context-map.md` §2.6,
`discovery/understanding-report.md` §2.4, `discovery/glossary.md`,
`discovery/open-questions.md`, `discovery/inventory/_oq011-konflik-isi.md`.

Korpus `D:\XML\RNM_BRD\` dan `D:\XML\nusantara-re\` READ-ONLY. Tulis hanya ke `OUTPUT_HASIL_RNM\`.

## Lima OQ pemblokir konteks ini

Dari register `discovery/open-questions.md` (nomor diambil dari register, bukan dari ingatan):
**OQ-002**, **OQ-018**, **OQ-020**, **OQ-021**, **OQ-029**.
Sumber: `discovery/D3-D4-CLOSING-REPORT.md` §4.1 — Claim — Life adalah konteks dengan pemblokir
paling sedikit (5) dan cakupan bukti `full`.

---

## Temuan pencarian fakta sebelum ronde ini

Dua temuan **mengoreksi** catatan FASE A untuk konteks ini. Keduanya fakta terukur, bukan jawaban
atas Q di bawah.

### T-1 `[terverifikasi]` — Claim Life tidak memuat identitas orang ter-hardcode

`pyPosition` dibandingkan terhadap **tiga kode peran** di **17 berkas**:
`'ReasLifeAdmin'`, `'ReasLifeSPV'`, `'ReasLifeMedicalAdvisor'`.
`OperatorID.pyUserIdentifier` (2 berkas) dan `OperatorID.pyUserName` (4 berkas) dipakai sebagai
**data** — bukan guard terhadap literal nama orang.

```
grep -rhoE "pyPosition[ ]*[=!]+[ ]*'[^']*'|pyPosition[ ]*[=!]+[ ]*\"[^\"]*\"" "Claim Life" --include="*.xml" \
  | grep -oE "'[^']*'|\"[^\"]*\"" | sort -u
grep -rhoE "(pyUserIdentifier|pyUserName)[^<]{0,45}" "Claim Life" --include="*.xml" \
  | sed 's/&amp;#61;/=/g;s/\]\[/ /g;s/[][]//g' | grep -E "[=!]"      # -> kosong
```

Berkas yang memuat `pyUserIdentifier`/`pyUserName`: `Activity/RejectOSClaimLife_Act.xml`,
`Activity/SendEmailKlaimLF.xml`, `Harness/Committe_Life.xml`, `RDBList/GetTokenStorage_SQL.xml`,
`Section/ClaimComite.xml`, `Section/RejectOSClaimLife_Sec.xml`.

**Ini model otorisasi terbersih di korpus.** Nilai `IT Developer` yang tercatat di D1 berasal dari
sapuan korpus-wide — **tidak ada di modul ini**.

### T-2 `[terverifikasi]` — Claim Life bersih dari hostname DEV, tetapi berasal dari 4 server

```
grep -rlio "appdev\.nusantarare" "Claim Life" --include="*.xml"     # -> kosong
grep -rho "<pxHostId>[^<]*" "Claim Life" --include="*.xml" | sed 's/<[^>]*>//' | sort | uniq -c
```
→ `pega-nusre` **73**, `jboss1073` **64**, `jboss117` **1**, `1d407e1106c501b92d737506992c6d06` **1**.

Ekspor modul ini **gabungan beberapa server**, bukan snapshot satu lingkungan.

---

## Q1 — Batas konteks: Claim Life sendiri, atau bersama Komite Claim Life?

**Bukti `[terverifikasi]`:**

| Fakta | Path + rule |
| --- | --- |
| Satu rule tulis dipakai dua sisi (hash ternormalisasi `c50bfd9a12`, **tidak** di register OQ-011) | `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` dan `Komite Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml` — `ASM-FW-GISFW-INT-LIFE_PREMIUM_DETAIL` / `RNM!UPDATEOSAKSEPTASICLAIMLIFE_SQL` / `RULE-CONNECT-SQL` |
| Activity Arasapas hanya ada di Claim Life tetapi dipanggil Komite (OQ-035) | `Claim Life/Activity/serviceInsertArasapasClaimLife_act.xml` |
| Penyerahan ke Komite ada sebagai Activity **di luar graf flow** (OQ-039) | `Claim Life/Activity/CreateKMTLife_Act.xml`, `Claim Life/Activity/GetListKomiteLife.xml` |
| Komite membawa 13 OQ pemblokir, hampir seluruhnya RBAC | `discovery/D3-D4-CLOSING-REPORT.md` §4.2 |

**Opsi:**

- **(a)** Claim Life saja; Komite Claim Life diperlakukan sebagai sistem luar dengan kontrak eksplisit atas `OS_AKSEPTASI_KLAIM_LIFE`
- **(b)** Claim Life + Komite Claim Life sebagai satu bounded context
- **(c)** Claim Life saja, tetapi rule tulis bersama ikut dispesifikasikan di dalam konteks ini

➡️ **Usulan: (a)**, dengan kontrak tertulis atas `OS_AKSEPTASI_KLAIM_LIFE` sebagai batas.
Alasan: menarik Komite masuk akan mengubah konteks tersiap (5 pemblokir, cakupan `full`) menjadi
yang paling terblokir. **Risiko yang diakui terbuka:** OQ-039 membuat pemicu penyerahan ke Komite
tidak diketahui **di opsi mana pun** — memilih (b) tidak menyelesaikannya.

**Pemilik:** (DBA / Product+UW / IAM / IT-infra) → **Product+UW**

**Jawaban work owner:** `[dijawab 2026-09-14]`

**(a)** — Claim Life sendiri; **Komite Claim Life = konteks/sistem luar** dengan kontrak atas
`OS_AKSEPTASI_KLAIM_LIFE`. Isi & alur internal Komite **tidak dispesifikasikan** di konteks ini.

**Titik penyerahan kini TERVERIFIKASI** (sebelumnya OQ-039): komite dibuat oleh
`Claim Life/Activity/CreateKMTLife_Act.xml` (`ASM-FW-GISFW-DATA-ADJUSTMENTLIFE` /
`CREATEKMTLIFE_ACT` / `RULE-OBJ-ACTIVITY`) yang men-*spawn* child work berkelas
`ASM-FW-GCNMFW-Work-KomiteLife` (`Call pxAddChildWork` + `Obj-Save`; `pyWorkPage` =
`ASM-FW-GCNMFW-Work-ClaimLife`). **Dipicu dari UI, bukan dari Flow** —
`<pyActivity>CreateKMTLife_Act</pyActivity>` di `Claim Life/Section/ClaimComite.xml` dan
`Claim Life/Harness/Committe_Life.xml`. Itu sebabnya tidak terlihat di graf flow.

**Diperiksa ulang ke korpus — cocok.** Tambahan yang ditemukan saat verifikasi: activity ini juga
memanggil `pxRetrieveReportData`, merujuk class roster `ASM-FW-GCNMFW-Int-EMAILKOMITE`, dan
memanggil `SendEmailKlaimLF`.

**MASIH TERBUKA:** kondisi/aturan bisnis **kapan** penyerahan dipicu (otomatis di atas nilai
tertentu, atau manual) → **OQ-039 dipersempit**, pemilik **Product+UW**.


---

## Q2 — Parity perilaku atau perbaikan? Apa non-goal-nya?

**Bukti `[terverifikasi]`:** tiga hal tidak dapat paritas karena memang tidak ada di Pega:

| Hal | Bukti |
| --- | --- |
| Model peran | **OQ-007** — nol rule identitas/otorisasi di 17 tipe rule korpus |
| Representasi uang | aturan proyek: uang material **jangan `float`** (`CLAUDE.md` §7) |
| Endpoint sebagai literal | **OQ-047** — daftar endpoint di tabel `M_LINK_SERVICE`, isinya tidak ada di korpus |

➡️ **Usulan: paritas perilaku sebagai baseline, dengan tiga penyimpangan eksplisit**, masing-masing
menjadi ADR: (1) RBAC dirancang dari tiga kode peran yang sudah ada (lihat T-1); (2) uang non-float;
(3) endpoint menjadi env var.

**Non-goal yang diusulkan:** tidak memindahkan tombol/jalur ber-penanda `(dev)`; tidak memindahkan
nilai ter-hardcode sebagai konstanta kode; **tidak merapikan alur** — jalur balik `IsSendtoAdmin`
dari tiga titik dipertahankan apa adanya (`Claim Life/Flow/Register_Flow.xml`,
`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `REGISTER_FLOW` / `RULE-OBJ-FLOW`).

**Pemilik:** (DBA / Product+UW / IAM / IT-infra) → **Product+UW**

**Jawaban work owner:** `[dijawab 2026-09-14]`

**Meniru sistem lama DAN memperbaikinya**, kecuali **3 hal** yang masing-masing menjadi ADR:
(1) RBAC dari 3 kode peran existing; (2) uang non-float; (3) endpoint jadi env var.

**Non-goal disetujui:** tidak memindahkan tombol `(dev)`; tidak menyalin nilai ter-hardcode menjadi
konstanta; tidak merapikan alur; **tidak menambah produk/lini baru**.


---

## Q3 — OQ-018: korpus ini cerminan production atau bukan, untuk Claim Life spesifik?

**Bukti `[terverifikasi]`:** lihat **T-2** di atas — nol hostname DEV di modul ini, tetapi
**4 `pxHostId` berbeda** (`pega-nusre` 73, `jboss1073` 64, `jboss117` 1, satu id hash 1).

Konteks register: **OQ-018** mencatat hostname DEV `appdev.nusantarare.com` (7×) di **81 berkas /
15 modul** korpus — Claim Life **bukan** salah satunya.

➡️ **Usulan: perlakukan sebagai "mendekati production, belum terkonfirmasi".** Yang diminta dari
pemilik export adalah satu hal konkret: **`pxHostId` mana yang production**. Sampai itu dijawab,
setiap perilaku yang berbeda antar `pxHostId` ditandai risiko.
**Tidak** disarankan menganggapnya production hanya karena tidak ada string DEV — 4 server adalah
sinyal yang lebih kuat daripada ketiadaan satu string.

**Pemilik:** (DBA / Product+UW / IAM / IT-infra) → **IT-infra / pemilik export Pega**

**Jawaban work owner:** `[dijawab 2026-09-14]`

**`jboss1073` = production, `jboss117` = dev.** Sistem memakai **mirroring**, sehingga perbedaan
antar `pxHostId` **tidak dianggap risiko perilaku**.

**OQ-018 untuk Claim Life → ditutup** (korpus dianggap production).

_Catatan cakupan:_ lingkungan `pega-nusre` (73 berkas) dan `1d407e1106c5…` belum dinyatakan;
hostname DEV di 15 modul lain tetap terbuka.


---

## Q4 — OQ-021: tiga kode peran — lengkap, dan apa artinya?

**Bukti `[terverifikasi]`:** lihat **T-1**. Pola gerbang yang terbaca, antara lain:

```
pyPosition != 'ReasLifeAdmin' || pyWorkPage.ClaimData.…      (8 kemunculan)
pyPosition == 'ReasLifeMedicalAdvisor'                        (5 kemunculan)
pyPosition == 'ReasLifeSPV'                                   (6 kemunculan)
```

Tahapan siklus yang perlu dipetakan ke peran
(`Claim Life/Flow/Register_Flow.xml`, `<pyStartActivity>Start2`):
**Register → Outstanding → Medical Check → Claim Analis**.

**Yang dibutuhkan:**
1. Apakah ketiganya daftar **lengkap** peran di siklus klaim Life?
2. Siapa melakukan apa di keempat tahap?
3. Apakah satu orang boleh memegang lebih dari satu peran?

➡️ **Usulan pemetaan — `[dugaan]`, bukan bukti:** Register + Outstanding → `ReasLifeAdmin`;
Medical Check → `ReasLifeMedicalAdvisor`; Claim Analis → `ReasLifeSPV`.
Dasarnya posisi guard, **bukan** pernyataan korpus — graf tidak menyatakan pemetaan shape → peran
(**OQ-024**: `<pyWorkBasket>` kosong, `<pyRouteTo>` = `Custom`).
Bila work owner tidak yakin, **OQ-021 tetap terbuka** dan pemetaan ini tetap ditandai `[dugaan]`.

**Pemilik:** (DBA / Product+UW / IAM / IT-infra) → **IAM + Product+UW**

**Jawaban work owner:** `[dijawab 2026-09-14]`

**Tiga peran SUDAH LENGKAP.** Pemetaan tahap → peran **dikonfirmasi work owner, bukan lagi dugaan**:

| Tahap | Peran |
| --- | --- |
| Register + Outstanding | `ReasLifeAdmin` |
| Medical Check | `ReasLifeMedicalAdvisor` |
| Claim Analis | `ReasLifeSPV` |

**Rangkap peran: TIDAK boleh**, kecuali akses ditambahkan eksplisit di role akun.
**OQ-021 → ditutup untuk Claim Life.** RBAC ADR memakai ketiga peran ini.


---

## Q5 — OQ-020: arti kode status di Claim Life

**Bukti `[terverifikasi]`**, seluruhnya **tanpa penjelasan di korpus**:

| Kode | Nilai literal | Pemakaian | Bukti |
| --- | --- | --- | --- |
| `STS_REJECT` | `'0'`, `'1'`, `'2'` | 50 kemunculan; `=='1' \|\| =='2'` dikelompokkan 7×; `=="0"&&…Type` 2×; `!=1` 1× | sapuan `Claim Life` |
| `pyWorkPage.SendtoAdmin` | `1` | guard `IsSendtoAdmin` — jalur balik dari 3 titik | `Claim Life/When/IsSendtoAdmin.xml` (`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOADMIN` / `RULE-OBJ-WHEN`) |
| `AcceptStatus` | `2` | 1 kemunculan | sapuan `Claim Life` |
| `ContentNote` | `"DEATH"` | gerbang `RDB-List` + `Property-Set` | `Claim Life/Activity/SaveOutStandingLife_Act.xml` |
| `BusinessCode` | `L1`…`L11` | satu precondition menguji 11 kode berderet (**OQ-038**) | idem |

`[terverifikasi]` **`PaymentType` tidak dipakai sama sekali** di modul ini — berbeda dari tiga modul
Claim lain (`discovery/flows/_SUMMARY-claim.md` §4).

Perintah audit:
```
grep -rhoE "(AcceptStatus|STS_REJECT|SendtoAdmin)[ ]*[=!]+[ ]*[\"']?[0-9A-Za-z]{1,10}" "Claim Life" \
  --include="*.xml" | sort | uniq -c
```

➡️ **Usulan prioritas: dua yang paling menentukan.**
**`STS_REJECT`** — karena `1` dan `2` selalu dikelompokkan, `[dugaan]` `0` = aktif dan `1`/`2` = dua
bentuk penolakan yang berbeda; **perlu konfirmasi apa bedanya**.
**`ContentNote = "DEATH"`** — apakah ia sebab klaim, dan adakah nilai lain.
`BusinessCode L1–L11` dapat menyusul sebagai daftar referensi, tetapi **jangan ditebak** karena
ter-hardcode dalam satu precondition (OQ-038).

**Pemilik:** (DBA / Product+UW / IAM / IT-infra) → **Product+UW** (+ Finance untuk kode bernilai uang)

**Jawaban work owner:** `[dijawab 2026-09-14]`

| Kode | Nilai | Arti |
| --- | --- | --- |
| `STS_REJECT` | `0` / `1` / `2` | **Outstanding** / **Aksep** / **Reject** |
| `AcceptStatus` | `1` / `2` | diaksep / reject |
| `SendtoAdmin` | `1` | dari `ReasLifeMedicalAdvisor` atau `ReasLifeSPV` → kembali ke `ReasLifeAdmin` |
| `ContentNote` | `DEATH`/`HEALTH`/`CI`/`TPD`/`TI` | **jenis klaim**, diturunkan dari `BusinessCode` |

⚠️ **Nama `STS_REJECT` menyesatkan — nilai `1` = diaksep, bukan ditolak.** Peringatan direkam di
`CONTEXT.md` dan di register OQ-020.

**`BusinessCode` `L1`–`L21`** (nama produk + jenis klaim) diberikan lengkap → tabel penuh di
`CONTEXT.md`. Modul ini hanya menguji `L1`–`L11` (semua `DEATH`); `L12`–`L21` ada di master produk.

**OQ-020 → ditutup untuk Claim Life. OQ-038 → ditutup.**


---

## Q6 — OQ-029: `IsPEGAPROD` dan `IsSendtoMedical` — keduanya tidak terbaca

**Bukti `[terverifikasi]`:** `<pyLabel>` kedua rule hanya berisi template kosong
`[first value][relation][second value]`.

| Rule | Menggerbangi | Path |
| --- | --- | --- |
| `@BASECLASS` / `ISPEGAPROD` / `RULE-OBJ-WHEN` (**berkonflik, OQ-011 #305**) | **tiga efek keluar**: unggah berkas, simpan utama, email | `Claim Life/When/IsPEGAPROD.xml`; dipakai di `Activity/InsertGoogleStorage_Act.xml`, `Activity/SaveOutStandingLife_Act.xml` (642.787 byte), `Activity/SendEmailKlaimLF.xml` |
| `ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `ISSENDTOMEDICAL` / `RULE-OBJ-WHEN` | pengembalian kasus dari **Claim Analis → Medical Check** | `Claim Life/When/IsSendtoMedical.xml` |

Pembanding: `When/IsSendtoAdmin.xml` **terbaca** (`pyWorkPage.SendtoAdmin = 1`).

➡️ **Usulan terbelah, karena keduanya berbeda sifat:**
**`IsPEGAPROD` — jangan dimigrasikan sebagai rule.** Ganti dengan flag lingkungan eksplisit
(mis. `ENV=production`) yang menggerbangi ketiga efek keluar yang sama. Perilakunya menjadi terbaca,
bukan tersembunyi di rule berkonflik. **Kandidat ADR.**
**`IsSendtoMedical` — tetap OQ terbuka.** Kapan analis mengembalikan kasus ke medis adalah **aturan
bisnis**; tidak boleh dikarang.

**Pemilik:** (DBA / Product+UW / IAM / IT-infra) → `IsPEGAPROD`: **IT-infra**; `IsSendtoMedical`: **Product+UW**

**Jawaban work owner:** `[dijawab 2026-09-14]`

**`IsPEGAPROD`:** setuju — **jangan ditiru sebagai rule**; ganti flag lingkungan
(`ENV=production`) yang menggerbangi 3 efek keluar (unggah berkas, simpan utama, email).
**Jadikan ADR.**

**`IsSendtoMedical = 1`:** dari **`ReasLifeSPV`**, kasus dikembalikan ke
**`ReasLifeMedicalAdvisor`**.

**OQ-029 → ditutup untuk Claim Life (keduanya).**


---

## Q7 — OQ-002: penomoran klaim dan `PROC_GENERATE_SEQUENCE_NUMBER`

**Bukti `[terverifikasi]`:** `Claim Life/Activity/SaveOutStandingLife_Act.xml`
(`ASM-FW-GCNMFW-WORK-CLAIMLIFE` / `SAVEOUTSTANDINGLIFE_ACT` / `RULE-OBJ-ACTIVITY`, 642.787 byte —
**16 langkah pertama terbaca, sisanya belum**) memanggil lewat `RequestType`:

| RequestType | Class |
| --- | --- |
| `Generate_NoKlaim_Life` | `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL` |
| `Generate_NoKlaim_LifeRetro` | `ASM-FW-GISFW-Int-LIFE_PREMIUM_DETAIL` |
| `GetSequenceNumber_SQL` | `ASM-FW-GISFW-Int-policyjson` → memanggil `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` |

**Isi ketiganya tidak ada di korpus** (OQ-002). Penomoran **bercabang retro / non-retro** — pola
sejajar dengan penomoran akseptasi di Komite Claim Life (`Generate_NoAccept_KMT_Life` vs
`…_LifeRetro`).

➡️ **Usulan: yang diminta bukan body-nya, melainkan kontraknya** — (i) format nomor klaim Life
retro dan non-retro, (ii) apa yang membedakan keduanya, (iii) apakah sequence per-tahun, per-lini,
atau global.
Pembanding yang menunjukkan pertanyaan ini bisa dijawab singkat: format nomor polis treaty **sudah
terbaca penuh** dari SQL — `'RNM-' || {InputData.CARI20} || '.T' || {…BusinessOldId} || '.' ||
to_char(sysdate,'MM.yyyy') || '.' || LPAD(POOLDATA.JSON_POLIS_TREATYIN_SEQ.NEXTVAL,5,'0')`
(`NB Treaty In/RDBList/GenerateNoPolicy.xml`, `ASM-FW-GISFW-INT-POLISTREATYIN` /
`ASM!GENERATENOPOLICY` / `RULE-CONNECT-SQL`).
Bila tidak terjawab, **OQ-002 tetap terbuka** dan spec menyebut penomoran sebagai **dependensi
database**, bukan mereplikasinya.

**Pemilik:** (DBA / Product+UW / IAM / IT-infra) → **DBA** (+ Product+UW untuk arti format)

**Jawaban work owner:** `[dijawab 2026-09-14]`

**Dua SQL lama `Generate_NoKlaim_Life` dan `Generate_NoKlaim_LifeRetro` SUDAH TIDAK DIPAKAI
(di-remark).** Penomoran sekarang lewat `GetSequenceNumber_SQL` →
`POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER`. **Jangan migrasikan dua SQL lama itu.**

**Diperiksa ulang ke korpus — cocok.** Asimetri indeks rujukan di `SaveOutStandingLife_Act.xml`:
`Generate_NoKlaim_Life` dan `…LifeRetro` muncul sebagai `<RequestType>` tetapi **0 kali** di
`pxRuleReferences`; `GetSequenceNumber_SQL` muncul **2 kali**. Ini **mengoreksi**
`discovery/flows/Claim Life.md` §3 yang mendaftar keduanya sebagai dipanggil.

**MASIH TERBUKA:** kontrak `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` — format nomor klaim dan
reset per tahun/lini/global → **OQ-002 dipersempit**, pemilik **DBA**.

---

## Catatan frontier

Ketujuh pertanyaan di atas adalah **seluruh frontier ronde 1** — yaitu keputusan yang prasyaratnya
sudah tuntas dan dapat dijawab sekarang tanpa menebak jawaban yang belum terdengar.

**Menunggu ronde berikutnya** (bergantung pada jawaban di atas):

| Menunggu | Bergantung pada |
| --- | --- |
| Mesin status & guard lengkap | Q5, Q6 |
| Penanganan galat & jejak audit | Q3, Q5 |
| Representasi uang (ADR non-float) | Q1 (batas), + skema `OS_AKSEPTASI_KLAIM_LIFE` |
| Bentuk API & UI | Q1, Q2 |
| Strategi migrasi data & cutover | Q3 |

## Status tindak lanjut

- [ ] Work owner mengisi kolom **Jawaban work owner** di tiap Q
- [ ] Jawaban dicatat sebagai keputusan: istilah → `CONTEXT.md`; keputusan sulit dibalik → `docs/adr/`
- [ ] OQ yang tidak terjawab **tetap terbuka** di `discovery/open-questions.md`
- [ ] Ronde 2 dihitung ulang dari frontier yang baru

**Belum dikerjakan dan tidak boleh dikerjakan sebelum kolom di atas terisi:** `spec.md`,
`issues/`, `CONTEXT.md`, `docs/adr/`.
