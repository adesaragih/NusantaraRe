# PROMPT — GILIRAN 14 *(sesi baru atau sesi yang sama, folder `OUTPUT_HASIL_RNM`, cabang `main` @ `665d7e8` atau lebih baru)*: **bp baris adjustment lahir saat Submit Register → bq Decision3 dirutekan dari bendera → br sunting sel adjustment → status tiket + panduan**

> Hanya konteks **Claim Life, PremiumList Life, Komite Claim Life**. GILIRAN-3 *(A/B/C)*, 4–13 tetap rujukan; aturan berhenti
> GILIRAN-6; **setiap pembacaan activity mencetak `pyStepsBlockName`**. Berkas `tco_*` tidak disentuh.

## 0. KEADAAN AWAL — DIVERIFIKASI ASISTEN 29-09-2026

| Klaim laporan GILIRAN-13 | Diperiksa ulang | Hasil |
| --- | --- | --- |
| 5 commit `a291a20` → `665d7e8`; migrasi `057` | ada; `057_seq_work_polis_dan_flag_ongoing.sql` + down | ✅ |
| Go 659 · 0 · 40 SKIP *(db)*; vitest 408 | dijalankan ulang di `main`: **659 PASS · 0 FAIL**, dengan tag `db` **40 SKIP**; vitest **408**; vet, gofmt, tsc bersih | ✅ |
| `Add` tampil hanya bila `ReasLifeSPV` | `ClaimLifeDetailGCNM.xml` b18160 `pyWorkPage.pyPosition =='ReasLifeSPV'` | ✅ |
| Decision3 memakai bendera | `InputPolicyHolder.xml` b274 `IsFlagOnGoingPolicy`, b276 `FlagOnGoingPolicy`, b283 `Decision3` | ✅ → **bq** |
| ⛔ **OQ-N9 "Pega melahirkan baris pertama saat pendaftaran"** | **Terbukti XML, bukan pertanyaan**: `Activity/SavePesertaClaim.xml` *(tombol `Submit` `InputRegisterClaimLife.xml` b27369 → b27393)* langkah **7.8** *(`RH_1.pySteps(7).pySteps(8)` b3672, **hidup** — `pyStepsBlockName` kosong b3682)*, prasyarat `.IsCheck=="true"` b3631 *(WhenTrue=2 lanjut, WhenFalse=3 lewati)*: "Set property AdjustmentList" b3685 mengisi `PremiumListDetail(<LAST>).AdjustmentList(<LAST>)` dengan **delapan** medan dari peserta — `CEDING_RETENTION` b3696, `SHARE_NUSANTARA_RE` b3742, `SUM_INSURED` b3768, `SUM_REASURED` b3788, `SHARE_RETRO` b3808, `CLAIM_AMOUNT` b3828, `RETROCEDED_SHARE` b3848, `CURRENCY` b3868 *(nilai dari b3697–b3849: `@divide(@toDecimal(@replaceAll(x,",",".")),1,4)`)*. `SaveInsuredClaim_Act` b1341–b1532 melakukan hal yang sama pada halaman sementara `TempDetail`. Pendaftaran Go **tidak** menyisipkan baris adjustment | ⛔ **bp**; ralat atas bo |
| DEV `T_MIGRASI` | masih `056` / 28 — **`057` belum dijalankan** | ⚠️ work owner |

## 1. KEPUTUSAN

| Butir | Isi |
| --- | --- |
| **bp** `[DIPUTUSKAN — bukti §0; veto work owner]` | `POST /api/klaim-life` *(Submit Register)* menyisipkan **satu baris `T_CLAIMLF_ADJUSTMENT` per peserta ber-`IsCheck = "true"`** dengan delapan medan di atas, rumus VERBATIM *(pembulatan 4 angka lewat pembagian — keanehan warisan yang sudah dikenal)*, dalam **transaksi pendaftaran yang sama**. Baris ini adalah `.AdjustmentList(1)` yang diwarisi `putaran` *(`SetIndexAdjustmentList` langkah 3)*. **bo diralat**: `Add` pada grid kosong tidak lagi jalan lahir baris pertama *(grid tidak pernah kosong untuk peserta terpilih)*; `Add` tetap = `putaran` bergerbang `ReasLifeSPV` b18160. Klaim lama tanpa baris *(DEV: nol klaim terlihat)* → A4 migrasi data |
| **bq** `[DIPUTUSKAN — XML; veto work owner]` | PremiumList **Decision3** dirutekan **otomatis** dari `FLAG_ONGOING_POLICY` *(`IsFlagOnGoingPolicy`, baca decision table utuh: baris, nilai, konektor tujuan)* — layar **tidak** menanyakannya lagi; OQ-PL-16 ditutup. Kedua tombol kini berbeda perilaku seperti Pega |
| **br** `[DIPUTUSKAN — XML; veto work owner]` | sunting sel baris adjustment: **hanya** kolom yang `ClaimLifeDetailGCNM.xml` grid b17126 tampilkan **dapat disunting** *(baca `pyEditOptions`/`pyReadOnly`/`pyCondition` tiap kolom; tujuh gerbang `.STS_REJECT` berlaku)*; `PUT /api/klaim-life/{id}/peserta/{pesertaId}/adjustment/{adjId}` bergerbang tahap + pemegang; kolom read-only ditolak 422; uang `apd.Decimal`; OQ-N8 ditutup |
| OQ-N7 *(Delete: penanda hapus atau tidak berlaku)*, OQ-PL-15 *(awal `SEQ_WORK_POLIS`)*, dan OQ lain | **tetap untuk work owner** |

## 2. URUTAN — satu commit per paket

| # | Paket | Commit |
| ---: | --- | --- |
| 1 | **bp** — Submit Register melahirkan baris adjustment; uji: peserta `IsCheck` true → satu baris dengan delapan nilai dari contoh literal; false → nol baris; tiket 02/03 diralat bertanggal *(bo, OQ-N9)*; panduan uji Claim Life §1.3 diperbarui *(data sintetis tidak lagi perlu baris adjustment manual)* | `claim-life: bp — Submit Register melahirkan baris adjustment (SavePesertaClaim 7.8)` |
| 2 | **bq** — Decision3 otomatis | `premiumlist-life: bq — Decision3 dirutekan dari FLAG_ONGOING_POLICY` |
| 3 | **br** — sunting sel adjustment + layar grid dapat disunting | `claim-life: br — sunting sel adjustment menurut grid XML` |
| 4 | **status tiket + panduan + data uji sintetis** *(sesuaikan `data_uji_tiga_modul.sql` dengan bp)* | `docs: status tiket, panduan, data uji — giliran 14` |

## 3. LAPORAN · TELEMETRI

Baris pertama alasan berhenti; tabel **paket → commit → rule XML → rute/komponen**; ralat tiket; OQ; angka uji tiap commit dengan dan
tanpa tag `db`; bab **TELEMETRI EKSEKUSI**. `-migrate` tidak dijalankan executor; tulis ke DEV tidak dilakukan.

---

*Disusun 29 September 2026 sesudah verifikasi `9784049..665d7e8` (uji dijalankan ulang), pembacaan `SavePesertaClaim.xml` langkah 7.8
(b3600–b3868, prasyarat dan blok remark), `SaveInsuredClaim_Act.xml` b1341–b1532, `ClaimLifeDetailGCNM.xml` b18160,
`InputPolicyHolder.xml` b274–b283, dan `T_MIGRASI` DEV.*
