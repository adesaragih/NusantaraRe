# Grilling — Claim — Life — Ronde 2

Status: answered (work owner, 2026-09-14) — keputusan sudah dicatat ke `CONTEXT.md` dan `docs/adr/`
Konteks: `claim-life` (Claim — Life)
Tanggal: 2026-09-14
Skill: `/mattpocock-skills:grill-with-docs` (grilling + domain-modeling)
Ronde sebelumnya: `grilling-ronde-1.md` (Q1–Q7)

> Berkas ini **arsip ronde**: pertanyaan apa adanya, jawaban work owner apa adanya, lalu ke mana
> tiap jawaban dicatat. Ditulis ulang dari transkrip sesi agar rujukan `sumber:` di `docs/adr/`
> menunjuk berkas yang benar-benar ada. **Tidak ada jawaban yang ditebak atau diparafrase.**

---

## Q8 — `STS_REJECT` dan `AcceptStatus` sama-sama mengkodekan aksep/reject. Redundan, atau beda peran?

Dari jawaban Ronde 1: `STS_REJECT` `1` = Aksep, `2` = Reject; `AcceptStatus` `1` = diaksep,
`2` = reject. **Nilai dan artinya identik.** Bedanya: `STS_REJECT` adalah **kolom tabel**
`OS_AKSEPTASI_KLAIM_LIFE` (terbaca di `Claim Life/RDBList/UpdateOsAkseptasiClaimLife_sql.xml`) dan
muncul 50× di modul; `AcceptStatus` hanya **1 kemunculan** di `Claim Life` — tetapi ia kode utama di
keempat modul **Komite** (`.AcceptStatus = "1"` menggerbangi `IsKomiteLoop`).

➡️ **Dugaan: `AcceptStatus` milik konteks Komite, `STS_REJECT` milik Claim Life**, dan satu
kemunculan di Claim Life adalah titik sentuh. Bila benar, sistem baru menyimpan **satu** status klaim
dan memetakannya di batas kontrak.

**Jawaban work owner:**
> Setuju dugaan Claude. Satu status klaim di sistem baru; STS_REJECT (Claim Life) dan AcceptStatus
> (Komite) dipetakan di batas kontrak. Bukti mendukung (AcceptStatus memang kode utama Komite).

**Pemilik:** Product+UW
**Dicatat ke:** `CONTEXT.md` (entri `AcceptStatus`), **ADR-0001** Consequences.
**Menyisakan:** **OQ-061** — `STS_REJECT` hidup di tiga tingkat (klaim, `PremiumListDetail`,
`AdjustmentList`); "satu status klaim" belum menyatakan bagaimana tingkat baris diperlakukan.

---

## Q9 — Mesin status: di titik mana `STS_REJECT` ditulis, dan apakah `1`/`2` terminal?

Graf (`Claim Life/Flow/Register_Flow.xml`) punya 4 tahap dan 2 jalur balik. `STS_REJECT = 0` =
Outstanding. Tiga hal belum diketahui: **(i)** apakah `0` ditulis saat Register atau saat masuk
Outstanding; **(ii)** apakah `1`/`2` terminal, atau klaim bisa dibuka ulang; **(iii)** status apa yang
berlaku selagi kasus di Medical Check.

➡️ **Dugaan: `0` sejak Register, bertahan sepanjang Medical Check, berubah ke `1`/`2` hanya di
Claim Analis, dan keduanya terminal.**

**Jawaban work owner:**
> Saat di admin, adjustmentlist insert dengan STS_REJECT = 0
> lalu di saat di klaim analisis kirim komite, jika komite di aksep AcceptStatus =1 maka
> STS_REJECT = 1, jika AcceptStatus =2 maka STS_REJECT = 2

**Pemilik:** Product+UW
**Dicatat ke:** `CONTEXT.md` (blok **Perpindahan nilai** pada entri `STS_REJECT`, entri baru
`AdjustmentList`), **ADR-0001** Kontrak 2 (jalur balik).

**Dikuatkan bukti korpus `[terverifikasi]`:** `Komite Claim Life/Activity/KomitePostAdjustment.xml`
(`ASM-FW-GCNMFW-WORK-KOMITELIFE!KOMITEPOSTADJUSTMENT` / `RULE-OBJ-ACTIVITY`) memuat 6 `Property-Set`
→ `STS_REJECT = 1`, 2 → `= 2`, dan precondition `pyWorkPage.AcceptStatus = 1` (4×) / `==2` (1×)
bersama `pyWorkPage.KomiteCount == pyWorkPage.KomiteLoop`.

**Menyisakan:** (i) dan (ii) tidak terjawab langsung — apakah `1`/`2` **terminal** dan apakah klaim
dapat dibuka ulang masih terbuka; serta **OQ-039**, karena
`Claim Life/Activity/SaveAdjustment_Act.xml` menulis `STS_REJECT = 1` **tanpa** melalui Komite.

---

## Q10 — Dari sembilan kolom numerik, mana uang dan mana rasio/persen? *(mengisi ADR-0003)*

➡️ Nama menyarankan `EM_PERCENT` persen dan `SHARE_*` rasio — tetapi **nama bukan bukti**, dan di
korpus ini nama sudah terbukti menipu (`STS_REJECT`). Tidak ditebak.

**Jawaban work owner:**
> SUM_INSURED uang · CEDING_RETENTION uang · SUM_REASURED uang · SHARE_NUSANTARA_RE uang ·
> CLAIM_AMOUNT uang · SHARE_RETRO uang · CLAIM_RETRO uang · RETROCEDED_SHARE uang ·
> EM_PERCENT *(dikosongkan)*

Susulan work owner pada pesan berikutnya:
> EM_PERCENT		   PERSEN

**Pemilik:** Product+UW (+ Finance)
**Dicatat ke:** **ADR-0003** — dinaikkan dari `draft` ke `accepted`.
**Menyisakan:** **OQ-060** — pertanyaan susulan pada Q10 ("`CURRENCY` berlaku untuk seluruh rekam,
atau ada kolom bermata uang berbeda?") **tidak dijawab**.

---

## Q11 — Jejak audit: apa yang wajib terekam?

Bukti yang ada: `OperatorID.pyUserIdentifier` / `pyUserName` dipakai sebagai data di
`Activity/RejectOSClaimLife_Act.xml` dan `Activity/SendEmailKlaimLF.xml`. Tabel
`OS_AKSEPTASI_KLAIM_LIFE` punya kolom `CREATEOPNAME`, `ACCEPTATION_DATE`, `CONFIRMATION_DATE`,
`CLAIM_RECEIVED_DATE`, `COMPLETE_DATE`. Ada juga `RDBList/InsertLogServiceClaim.xml` →
`INSERT INTO pooldata.monitoring_klaim_log`.

➡️ **Usulan: rekam siapa+kapan untuk setiap transisi status dan setiap jalur balik.**

**Jawaban work owner:**
> Audit: rekam SIAPA + KAPAN untuk setiap transisi status DAN setiap jalur balik (SendtoAdmin,
> SendtoMedical) — bukan hanya CREATEOPNAME + tanggal seperti sekarang. Ini penyimpangan sadar dari
> sistem lama (perbaikan), catat sebagai keputusan/ADR bila perlu.

**Pemilik:** Product+UW
**Dicatat ke:** **ADR-0007** (`accepted`).

---

## Q12 — Apa yang terjadi bila efek keluar gagal?

Empat efek keluar: unggah Google Storage, email, `serviceInsertArasapasClaimLife_act`, dan
`convertJsonNusareToProductionClaimLife`. Di konteks facultative ada pola gerbang
`IsSuccessHitService`; di Claim Life hanya ditemukan **pencatatan** (`monitoring_klaim_log`),
**bukan** gerbang retry.

➡️ **Dugaan: kegagalan dicatat tetapi tidak memblokir** alur klaim.

**Jawaban work owner:**
> Efek keluar (Google Storage, email, Arasapas, konversi) bersifat ASINKRON: catat kegagalan,
> JANGAN blokir alur klaim, sediakan antre-ulang (retry). Pertahankan perilaku lama "tidak
> memblokir", tambah keandalan. Catat sebagai keputusan.

**Pemilik:** Product+UW + IT-infra
**Dicatat ke:** **ADR-0008** (`accepted`).
**Menyisakan:** apakah klaim boleh mencapai status akhir selagi efek keluar masih tertunda —
dicatat sebagai konsekuensi terbuka di ADR-0008, bukan OQ tersendiri.

---

## Q13 — Migrasi data: klaim yang sedang berjalan ikut pindah?

➡️ **Usulan: hanya klaim baru; klaim berjalan diselesaikan di Pega** — keduanya menulis ke
`OS_AKSEPTASI_KLAIM_LIFE` yang sama, sehingga koeksistensi mungkin tanpa migrasi data.

**Jawaban work owner:**
> Semua data akan dipindah

**Pemilik:** Product+UW + DBA
**Dicatat ke:** **ADR-0009** (`accepted`) — usulan koeksistensi **ditolak**, penolakannya dicatat
supaya tidak diusulkan ulang.

---

## Q14 — Penyimpanan lampiran tetap Google Storage?

Bukti: `Activity/InsertGoogleStorage_Act.xml`, `GetUrlGoogleStorage_Act.xml`, tabel
`T_STORAGE_IMAGE` + `POOLDATA.T_FOLDER_IMAGE`, procedure `POOLDATA.GET_TOKEN_STORAGE` (isinya tidak
ada di korpus, OQ-002), dan `RDBList/Insert_T_Storage_SQL.xml` / `Update_T_Storage_SQL.xml`.

➡️ **Usulan: pertahankan Google Storage.**

**Jawaban work owner:**
> pertahankan Google Storage — mengganti penyimpanan berarti memigrasikan berkas yang sudah ada, dan
> itu di luar cakupan konteks ini. Tetapi GET_TOKEN_STORAGE tetap batas pengetahuan; bila
> dipertahankan, kontraknya perlu DBA.

**Pemilik:** IT-infra + DBA
**Dicatat ke:** **ADR-0010** — ditulis `draft`, **menunggu persetujuan**, karena Q14 tidak termasuk
daftar ADR yang diminta pada instruksi tindak lanjut (Q11, Q12, Q13).

---

## Q15 — Muatan kontrak ke Komite — apakah yang saya temukan sudah lengkap?

➡️ **Pertanyaannya bukan "apa muatannya"** — itu sudah terbaca — melainkan apakah ini muatan yang
**benar**, dan apakah `KomiteLoop` memang ditentukan Claim Life.

**Jawaban work owner:**
> Muatan kontrak ke Komite yang terbaca sudah benar, TAMBAHKAN yang seharusnya ikut namun belum:
> nilai klaim, mata uang (CURRENCY), dan STS_REJECT saat penyerahan — supaya Komite punya konteks
> tanpa membaca balik ke Claim Life. KomiteLoop (jumlah tingkat tangga) DITENTUKAN Claim Life;
> sumber nilainya tetap OQ-032 terbuka sampai dipastikan.

**Pemilik:** Product+UW
**Dicatat ke:** **ADR-0001** §"Tambahan pada muatan", `CONTEXT.md` (entri baru `KomiteLoop`),
register **OQ-032** (dipersempit: kepemilikan terjawab, penentu nilai terbuka).

---

## Instruksi tindak lanjut dari work owner

> Setelah mencatat: (1) selesaikan ADR yang prasyaratnya kini lengkap (Q11 audit, Q12 efek keluar,
> Q13 migrasi/koeksistensi bila memenuhi kriteria ADR); ADR-0003 TETAP draft menunggu Q10.
> (2) Perbarui CONTEXT.md + open-questions.md. (3) Hitung apakah grilling Claim Life sudah MATANG
> untuk to-spec — daftar apa yang masih menghalangi (harusnya tinggal Q10/ADR-0003 + OQ-032

**Catatan atas butir (1):** instruksi "ADR-0003 TETAP draft menunggu Q10" ditulis **sebelum**
work owner mengirim `EM_PERCENT  PERSEN`. Dengan susulan itu Q10 lengkap kecuali cakupan `CURRENCY`,
sehingga ADR-0003 dinaikkan ke `accepted` dan sisa pertanyaannya dipindah ke **OQ-060**.

## Status tindak lanjut

- ADR ditulis/diperbarui: **0001** (diperluas), **0003** (draft → accepted), **0007**, **0008**,
  **0009** (baru, accepted), **0010** (baru, draft — menunggu persetujuan)
- `CONTEXT.md`: entri `STS_REJECT` diperluas; entri baru `AdjustmentList`, `KomiteLoop`; entri
  `AcceptStatus` dipertajam sebagai kosakata konteks luar
- `discovery/open-questions.md`: **OQ-032** dan **OQ-039** dipersempit; **OQ-060** dan **OQ-061**
  dibuat
- `spec.md` dan `issues/` **tidak disentuh** — menunggu persetujuan eksplisit
