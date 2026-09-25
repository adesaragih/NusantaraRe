# E22: Rekonsiliasi eksak siklus endorsement

**What to build:** Pembanding **nol-selisih** untuk kasus endorsement, memakai kerangka yang sudah
dibangun untuk New Business — **tanpa pembanding baru**.

⛔ **Seluruh kejanggalan yang diport apa adanya harus muncul sebagai selisih NOL**, bukan sebagai
selisih yang dimaafkan. Bila sistem baru "memperbaiki" salah satunya diam-diam, pembanding inilah
yang menangkapnya.

⚠️ **Satu pengecualian yang sah:** perbaikan A.5 mengubah tipe nilai cadangan premi Nusantara Re dari
string menjadi angka. Pembandingan mentah akan menunjukkan **beda tipe meski nilainya sama**.
Pembanding **menormalkan** sebelum membandingkan, dan selisih semacam ini **dijelaskan oleh A.5** —
bukan ditandai cacat.

⛔ **Tidak pernah menyentuh berkas mentah ber-PII.** Fixture yang dipakai adalah yang sudah
ter-de-identifikasi.

⚠️ Tiket ini **tidak bersandar** pada klaim arsip bahwa nilai dasar akseptasi adalah selisih TSI —
klaim itu tetap **`[belum diuji]`**.

**Asal (Pega).** — (pembanding, bukan port rule)

**Keputusan.** K-046 · A.5 · ADR-0001 · K-025

**Blocked by:** `..\15-de-identifikasi-berkas-kasus.md` · `..\16-rekonsiliasi-eksak-tahap-1.md` ·
E09 · E14 · E15 · E19
⛔ Kedua tiket NB itu **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] Memakai **kerangka rekonsiliasi New Business**, tanpa pembanding baru
- [ ] ⛔ **Tidak pernah** membaca berkas mentah ber-PII; hanya fixture ter-de-identifikasi
- [ ] Seluruh kejanggalan K-046 menghasilkan **selisih nol**
- [ ] Pembanding **menormalkan** string dan angka; selisih tipe dari **A.5** dijelaskan, bukan ditandai cacat
- [ ] Mencakup ketiga lapis before-image, selisih pembayaran, selisih baris spreading, dan premi jiwa
- [ ] Setiap selisih yang tersisa **dapat ditelusuri** ke rule Pega asalnya (§4.6)
- [ ] Nilai medis, nomor polis, dan nama orang **tidak pernah** muncul di laporan pembanding
