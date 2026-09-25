# F14: Rekonsiliasi eksak Fac Out

**What to build:** Pembanding **nol-selisih** untuk baris produksi Fac Out, memakai **kerangka
rekonsiliasi New Business** — tanpa pembanding baru.

`[terverifikasi]` Sasaran pembandingnya sudah ada dan sudah flat: `DDL\FACOUTPRODUCTION.txt` — **65
kolom**, dengan pasangan `_MENJADI`/`_SELISIH`. Pembanding membaca baris yang ditulis **F10**, bukan
struktur staging.

`[terverifikasi]` Uang di tabel itu bertipe `NUMBER` dengan **tiga skala berbeda**: `NUMBER(20,4)`
(nilai pokok), `NUMBER(20,8)` (premi & komisi coverage), `NUMBER(25,20)` (`RATE`, `PRORATE`).
Pembanding **menghormati ketiganya apa adanya** dan tidak membulatkan ke satu skala bersama.

⛔ **Seluruh kejanggalan yang diport apa adanya harus muncul sebagai selisih NOL**, bukan sebagai
selisih yang dimaafkan. Bila sistem baru "memperbaiki" salah satunya diam-diam, pembanding inilah yang
menangkapnya.

Kejanggalan Fac Out yang **wajib** menghasilkan selisih nol:

| Penanda | Isi | Tiket |
| --- | --- | --- |
| ~~`K046_RIComIN_TidakPernahDisetel`~~ | ⛔ **DICABUT (K-060)** — keanehan tidak ada pada salinan `DDL\` yang berlaku; ia akibat salinan korpus basi. **Bukan** kandidat perbaikan, **tidak** diuji sebagai keanehan | — |
| `K046_Prorate_EnamPenugasanSalingMenimpa` | satu properti ditimpa berurutan | F06 |
| `K046_BackUpStatus_DuaGayaPembanding` | `==3` telanjang vs `=="3"` berkutip | F05 |
| `K046_LocalFacout_AngkaSebagaiString` | `local.Facout=="0"` | F05 |
| `K046_IsFacRetroOffer_NolBerkutip` | `"0"` berkutip vs `0` telanjang | F05 |
| `K046_OkupasiFacOut_GayaUjiTidakSeragam` | awalan dan kesetaraan bercampur | F03 |
| `K046_GenerateFacRetroNo_TanpaElse` | `FacType` tak dikenal → kode kosong | F09 |
| `K046_GuardFacOut_HanyaType7` | gerbang idempotensi hanya aktif satu jenis penyesuaian | F11 |
| `K046_MenjadiSelisih_TidakSetangkup` | pasangan kolom tidak setangkup | F11 |

⛔ **Tidak pernah menyentuh berkas mentah ber-PII.** Fixture yang dipakai **sudah
ter-de-identifikasi**. Kolom `FACOUTPRODUCTION` yang ber-PII — `DOB`, `LICENSEPLATE`, `ENGINENUMBER`,
`ADDRESS`, `OBJECTNAME`, `REINSURER_NAME` — **tidak pernah** muncul di laporan pembanding, hanya
dihitung ada/tidaknya.

⚠️ **Uji dengan fixture, bukan dengan mengandalkan data produksi.** K-019 mencatat fitur fac out
**usang secara bisnis**, sehingga volume kasus yang melewati jalur ini **mungkin nol di data mutakhir**.
Bila paralel run tidak pernah menyentuh jalur Fac Out, itu **bukan bukti port-nya benar** — hanya
bukti jalurnya tidak terpakai. `[pertanyaan terbuka]` berapa banyak kasus Fac Out yang benar-benar ada
di data mutakhir — korpus tidak memuat data produksi, jadi angkanya hanya dapat datang dari DBA.

⚠️ **Dua presisi pembanding di F06 tidak diseragamkan** — rate pada 10 desimal, premi pada 4. Pembanding
menghormati keduanya apa adanya (ADR-0005).

⚠️ Selisih tipe dari **A.5** (string `"0"` → angka `0`) adalah **perbaikan sadar** yang sudah
dijelaskan K-046; pembanding **menormalkan** sebelum membandingkan dan **tidak** menandainya cacat.

**Asal (Pega).** — (pembanding, bukan port rule)

**Keputusan.** K-046 · K-019 · K-007 · K-025 · ADR-0001 · ADR-0005

**Blocked by:** F10 · `..\15-de-identifikasi-berkas-kasus.md` · `..\16-rekonsiliasi-eksak-tahap-1.md`
⛔ Kedua tiket NB itu **belum dikerjakan** — lihat §0 indeks.

**Status:** blocked

- [ ] Memakai **kerangka rekonsiliasi New Business**, tanpa pembanding baru
- [ ] ⛔ **Tidak pernah** membaca berkas mentah ber-PII; hanya fixture ter-de-identifikasi
- [ ] **Delapan** penanda K-046 yang masih berlaku menghasilkan **selisih nol** (semula sembilan; `K046_RIComIN_TidakPernahDisetel` dicabut oleh K-060)
- [ ] ⛔ Jalur Fac Out diuji dengan **fixture**, tidak mengandalkan ada-tidaknya kasus di data produksi
- [ ] Dua presisi pembanding (10 dan 4 desimal) dihormati apa adanya
- [ ] Setiap selisih yang tersisa **dapat ditelusuri** ke rule Pega asalnya (§4.6)
- [ ] Nilai `DOB`, plat nomor, nomor mesin, alamat, nama orang, dan nomor polis produksi **tidak pernah** muncul di laporan pembanding
