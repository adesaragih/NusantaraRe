---
status: tertahan
---

# 11: Lima port hilir berdiri kosong dengan niat tercatat lebih dulu

*Asal: `SPEC-KOMITE-01.md` bagian Persyaratan, 70 persyaratan `S-xxx`. Aliran A-1a, A-1b, A-2, A-3 saja — nol persyaratan untuk A-4, A-5, maupun isi A-6.*

**What to build:** Tidak satu pun jalur berakhir di efek hilir tanpa port yang **dinamai**.
Niat memanggil tiap port dicatat **sebelum** panggilan, dengan kunci idempotensi per efek,
dan kegagalan port tidak membatalkan keputusan yang sudah tercatat — ia meninggalkan niat yang
belum terpenuhi dan **terbaca**.

Lima port: pencatatan akseptasi, pembalikan akseptasi, pengiriman kasir, penerbitan surat,
pengunggahan dokumen. Kelimanya **kosong**.

**Persyaratan:** `S-037`.

**Tidak termasuk:** **isi kelima port**, seluruhnya berpagar. PAGAR-02 dan PAGAR-03 —
`PEGA_JSON_OS_AKSEP_KLAIM`, `PEGA_JSON_OS_AKSEP_SUBJECTIVITY`, `HISTORYAKSEPTASIPEGA`
(`INVENTARIS-BUKTI.md` §2.1) dan nol baris `OS_AKSEPTASI_KLAIM.DATA_JSON` (§2.5 baris 2).
PAGAR-04 — `XOL2_AKSEP_KLAIM` (§2.1) dan `SetProtectionEstimation` (§2.3). PAGAR-05 dan
PAGAR-06 — nol baris `POOLDATA.DIRECTTOKASIR_LOG` (§2.5 baris 1). PAGAR-07 —
`PostEmailKomiteCNP` (§2.3). PAGAR-08 — `SendEmailWithAttachments` (§2.3).

**Jalur gagal:** port dipanggil tanpa niat tercatat lebih dulu → uji gagal · port gagal →
keputusan yang sudah tercatat **tidak** dibatalkan; niat tetap ada dan terbaca sebagai belum
terpenuhi.

**Uji:** BARU — port diuji sebagai **port kosong**: uji memeriksa bahwa port dipanggil dengan
niat yang tercatat lebih dulu dan kunci idempotensi, **bukan** memeriksa apa yang dikirimnya.

**Menggantikan:** `KomitePostAdjustment`·19.1 → `PostEmailKomiteCNP` · `InsertOSKlaimCNP`,
`InsertOSSubjectivityCNP` · `InsertXOLKlaimCNP` · `SaveRejectOSKomiteCNP` ·
`HitServiceToKasirKMT_Act`·10.3, 10.7, 10.9 · `GenerateAccCNP_act` → `InsertDocument_Act` ·
`SendEmailKlaim_KMT`, `SendEmailKlaimRejectClose_KMT` · `InsertHistoryAkseptasiPega_Sql`.
Seluruhnya **DIPAGARI**, bukan digantikan — yang dibangun hanyalah pintunya.

**Blocked by:**

- `05` — Akibat pada klaim dihitung di satu tempat dan ditulis dalam transaksi yang sama
- `08` — Nomor akseptasi terbit sekali, di dalam transaksi keputusan


**Dasar:** DECIDED(keputusan beku no. 6, `ADR-0031`, `ADR-0032`). Kelima port berpagar
PAGAR-02 sampai PAGAR-08 pada `SPEC-KOMITE-01.md` bagian *Out of Scope*.

- [ ] Kelima port **dinamai** dan **kosong**; nol jalur berakhir di efek hilir tanpa melewati
      salah satunya.
- [ ] Niat dicatat **sebelum** panggilan, dengan kunci idempotensi per efek.
- [ ] Efek ke luar berada **di luar** batas transaksi dan dijalankan sesudahnya.
- [ ] Kegagalan port tidak membatalkan keputusan; niat yang belum terpenuhi terbaca.
- [ ] Nol persyaratan, nol muatan, dan nol penjaga ditulis untuk isi port mana pun —
      **termasuk pada alasannya**, bukan hanya pada kesimpulannya.

**Ketidakpastian:** **`KonversiKlaim_Act`** dan **`InsertJsonClaimTreatyNonProp_act`** adalah
efek ke luar yang dipanggil dari jalur keputusan sistem lama dan **tidak terpeta ke satu pun
dari kelima port**. Berkas keempat rule-nya **ada di repo** (`INVENTARIS-BUKTI.md` §2.4),
sehingga keduanya **lubang rancangan, bukan lubang bukti**, dan tidak boleh dipagari.
Keduanya berkeadaan `BELUM DIPUTUSKAN`; keputusannya milik pemilik proses.
