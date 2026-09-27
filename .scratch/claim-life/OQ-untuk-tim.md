# Pertanyaan untuk Tim — pembuka tiket Claim — Life

Dokumen bantu untuk mengumpulkan jawaban dari DBA, Product+UW, Finance, dan IT-infra.
Tujuannya membuka 6 tiket berstatus `needs-info` menjadi `ready-for-agent`.

**Cara pakai:** bawa ke tiap pemilik peran, isi kolom "Jawaban", lalu jawaban ini akan dicatat ke
CONTEXT.md/ADR/register lewat skill (bukan diketik langsung ke artefak). Yang tak terjawab tetap
OQ terbuka — jangan ditebak.

Sumber pertanyaan: `discovery/open-questions.md` (nomor OQ otoritatif di sana). Tiket yang
diblokir disebut di tiap butir.

---

## Untuk DBA

### OQ-002 — Kontrak `POOLDATA.PROC_GENERATE_SEQUENCE_NUMBER` (memblokir tiket 02, 12)
Penomoran klaim Life memanggil procedure ini; badannya tidak ada di korpus.
Yang dibutuhkan:
1. Format nomor klaim Life yang dihasilkan (retro dan non-retro) — bentuk string persisnya.
2. Apa yang membedakan nomor retro vs non-retro.
3. Sequence di-reset per tahun, per lini bisnis, atau berjalan terus (global)?
4. (untuk tiket 12) Kontrak `GET_TOKEN_STORAGE` — parameter masuk/keluar, dan apakah commit sendiri.

> **Jawaban:**

### OQ-001 — DDL Oracle produksi (memblokir tiket 13; TIDAK memblokir tiket 01)
Tidak ada DDL di korpus. Untuk migrasi data nyata (tiket 13) dibutuhkan definisi tabel:
`OS_AKSEPTASI_KLAIM_LIFE`, `JSON_KLAIM`, `M_LIFE_PREMIUM_DETAIL`, dan tabel klaim Life terkait —
nama kolom, tipe, presisi, nullability, PK/FK, index.
Catatan: tiket 01 (skema uji provisional) TIDAK menunggu ini — nama kolom sudah terbaca dari SQL,
tipe longgar. Yang menunggu hanya migrasi produksi.

> **Jawaban:**

### OQ-013 — Batas transaksi & identitas pengguna di dalam stored procedure (memblokir tiket 13)
Procedure `POOLDATA.*` melakukan `COMMIT` di dalam blok PL/SQL, dan `{OperatorID.pyUserIdentifier}`
dikirim sebagai parameter. Yang dibutuhkan:
1. Apakah semua procedure `POOLDATA.*` commit sendiri?
2. Identitas pengguna dipakai untuk apa di dalam procedure — audit trail, otorisasi, atau keduanya?

> **Jawaban:**

### OQ-018 — pxHostId production (memblokir tiket 13 untuk cutover; sebagian sudah dijawab)
Sudah dijawab untuk Claim Life: jboss1073 = production, jboss117 = dev (mirroring).
Sisa untuk cutover: konfirmasi pxHostId `pega-nusre` dan satu id-hash — lingkungan apa?

> **Jawaban:**

### OQ-047 — Daftar endpoint di tabel `M_LINK_SERVICE` (memblokir tiket 12)
Alamat integrasi (Google Storage, email, Arasapas, konversi) dibaca dari tabel `M_LINK_SERVICE`,
bukan hanya SystemSettings. Yang dibutuhkan: daftar endpoint aktual (untuk jadi env var), atau
konfirmasi bahwa semuanya via tabel itu dan strukturnya.

> **Jawaban:**

---

## Untuk Product + Underwriting

### OQ-032 — Sumber nilai `KomiteLoop` (jumlah tingkat tangga komite) (memblokir tiket 10)
`KomiteLoop` menentukan berapa tingkat persetujuan komite. Nilainya ditentukan Claim Life, tapi
sumber/aturan penentuannya tidak ada di korpus.
Pertanyaan: apa yang menentukan jumlah tingkat komite untuk sebuah klaim? (nilai klaim, jenis, dll)

> **Jawaban:**

### OQ-037 — Ambang nominal roster komite (memblokir tiket 10; Finance + Product+UW)
Komposisi roster komite ditentukan ambang nominal ter-hardcode (mis. batas 30jt/50jt di modul lain).
Pertanyaan: apa aturan ambang nominal yang menentukan siapa masuk roster komite untuk klaim Life?
Apakah mata uangnya IDR? Apakah dapat dikonfigurasi?

> **Jawaban:**

### Tinjauan aturan turunan "klaim selesai" (memblokir tiket 04)
Spec §3 menetapkan (bukan temuan korpus, tapi konsekuensi logis jawaban Anda):
> Klaim SELESAI bila tidak ada lagi baris AdjustmentList bernilai 0 DAN ada minimal satu baris
> bernilai 1 (aksep). Bila tak ada 0 dan tak ada 1 → ditolak seluruhnya, tapi masih bisa
> dilanjutkan dengan baris baru.
Pertanyaan: apakah aturan ini benar menurut proses bisnis Anda?

> **Jawaban:**

### Satuan pergeseran tanggal DOL (memblokir tiket 06)
Di ValidasiDOL_Act, cabang TP/TR memakai argumen `@addCalendar(...,1,...)` — [dugaan] +1 hari.
Definisi fungsi tidak ada di korpus.
Pertanyaan: untuk tipe TP/TR, apakah jendela validasi Date of Loss digeser +1 hari? Bila ya,
mengapa (aturan bisnis) — atau ini kekeliruan lama?

> **Jawaban:**

### OQ-060 — Cakupan kolom CURRENCY pada rekam akseptasi Life (memblokir bentuk tipe uang, tiket terkait)
CURRENCY ada di tingkat baris dan disalin dari baris 1. Pertanyaan: apakah satu klaim Life selalu
satu mata uang, atau boleh campur antar baris?

> **Jawaban:**

---

## Ringkasan pemetaan OQ → tiket

| OQ | Pemilik | Memblokir tiket |
| --- | --- | --- |
| OQ-002 | DBA | 02, 12 |
| OQ-001 | DBA | 13 (bukan 01) |
| OQ-013 | DBA | 13 |
| OQ-018 | IT-infra/DBA | 13 (cutover) |
| OQ-047 | DBA/Platform | 12 |
| OQ-032 | Product+UW | 10 |
| OQ-037 | Finance + Product+UW | 10 |
| aturan "klaim selesai" | Product+UW (work owner) | 04 |
| satuan DOL | Product+UW | 06 |
| OQ-060 | Product+UW + DBA | bentuk tipe uang |

**Setelah terjawab:** jawaban dibawa kembali ke sesi Claude, dicatat ke CONTEXT.md/ADR/register
lewat alur skill, lalu tiket yang bersangkutan dinaikkan `needs-info` → `ready-for-agent`.
Yang tetap tak terjawab: tiket tetap `needs-info`, jangan ditebak.
```

## 27 September 2026 — tiga pertanyaan untuk pemilik ekspor / pengembang Pega

**OQ-A.** Apakah ada **harness portal/navigasi** untuk Claim Life *(padanan `SFAPortalOpportunities`
di NB Treaty In)* beserta tombol Create-nya? Bila ada, mohon diekspor. Sampai terjawab, menu sidebar
Claim Life adalah **minimum berbukti**: `Inbox Claim Life` *(label `[tidak ada di korpus]`)* dan
`Register` *(`Register_Flow.xml:155`)*.

**OQ-B.** `Activity\setDetailClaim_act.xml` tampak **residu uji pengembang**: `Obj-Open-By-Handle`
pada satu handle literal *(b284)*, `Property-Set` tanggal literal Januari–Februari 2026 *(b478,
b711)*, prasyarat yang membandingkan `.NAME_OF_INSURED` dengan satu nama literal *(b872)*, lalu
`Obj-Save` *(b960)*. Ia terpasang pada tombol `Choose` popup pencarian polis
*(`SearchPolicy_Section.xml` 3337–3356, 3468)*. **Mohon konfirmasi** apakah ia memang tidak
dipakai di produksi; kami **tidak menirunya**.

**OQ-C.** `SendtoAdmin_Act` dan `SendtoAdmin_Act1` keduanya berprasyarat
`pyWorkPage.pyPosition=="ReasLifeMedicalAdvisor"` *(WhenTrue=2, WhenFalse=3 = lewati)*, padahal
tombol pemanggilnya — `Send Back to Register` dan `Send to Medical Check` — berdiri di layar
**Outstanding**, yang `pyPosition`-nya `ReasLifeAdmin`. Akibatnya menurut XML apa adanya: kedua
tombol itu **tidak menulis apa pun** pada posisi Admin. **Mohon konfirmasi** apakah prasyaratnya
memang demikian di produksi, atau salah tempel. Keputusan kami *(butir aw)* mengikuti **maksud**
yang terang dari label tombol dan penyambung alurnya, dan cacatnya dilaporkan di sini.

**OQ-D** *(urutan modul, untuk work owner)*. Sebelas medan layar Register terikat
`.PolicyDataLife.*` dan terisi dari kasus **PremiumList Life**. Register Claim Life karena itu
**menunggu modul PremiumList Life** untuk lengkap. Mohon konfirmasi urutan pengerjaan modul.

## 27 September 2026 — dua pertanyaan tambahan (paket Register(2))

**OQ-E** *(cacat halus, untuk pengembang Pega)*. `RDBList/GetPesertaClaim_sql1.xml:85` memasang
**kedua** `LIKE` tanpa syarat:

> `AND CERTIFICATE_NO LIKE '%'||{SearchPolicyHolder.CARI2}||'%'`
> `AND UPPER(NAME_OF_INSURED) LIKE '%'||{SearchPolicyHolder.CARI3}||'%'`

Di Oracle, `X LIKE '%'` bernilai **FALSE** ketika `X` NULL. Akibatnya kotak pencarian `Find
Insured` yang dibiarkan **kosong** pun membuang setiap peserta yang `NAME_OF_INSURED`- atau
`CERTIFICATE_NO`-nya NULL — baris yang hilang tanpa seorang pun memintanya, dan tanpa pesan.

**Mohon konfirmasi** apakah itu memang dikehendaki. Kami **tidak menirunya**: penyaring hanya
terpasang untuk kotak yang terisi, sehingga kotak kosong berarti *"jangan saring"* — yang memang
dibaca orang dari kotak kosong.

**OQ-F** *(untuk pemilik ekspor)*. `Activity/UploadCSVClaimLife_Act.xml` hanya **tiga** langkah:
`Page-Remove TempWorkPage` *(b250)*, `Page-New TempWorkPage` *(b340)*, lalu
`Call pxUploadCSVResults` *(b488)* — yaitu mesin unggah **bawaan platform**.

Pemetaan kolom CSV ke medan peserta karena itu **tidak ada di korpus**: ia tersimpan di
konfigurasi gadget, bukan di rule yang diekspor. **Mohon kirimkan** definisi pemetaan kolomnya
bila fitur unggah CSV memang dipakai; tanpa itu fitur ini tidak dapat ditiru tanpa mengarang, dan
kami tidak mengarang.
