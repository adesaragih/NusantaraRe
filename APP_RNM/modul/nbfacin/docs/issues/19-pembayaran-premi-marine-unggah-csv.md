# 19: Pembayaran premi MARINE CARGO di jalur unggah CSV anggota polis

> ⚠️ **Disusun agent dari XML atas perintah work owner — bukan hasil `/to-tickets`.** Perintah:
> `PROMPT-LANJUT-IMPLEMENT-TANPA-TANYA-NB-RNW.md` (root repo, aturan 5), diteruskan sesi `nusantarare-0f`
> 01-10-2026 (butir 2 pesannya).

**What to build:** layanan yang menghitung `pyWorkPage.Policy.Payment.Premium` dan `.Diskon` kasus MARINE CARGO
New Business persis seperti `PremiPaymentMarine` — jumlah premi dan diskon seluruh coverage kargo.

**Asal** `[terverifikasi]` — semua `D:\migrasi\RNM\NB FacIn\Activity\`:

| Rule | Langkah | Isi |
| --- | --- | --- |
| `UploadCSVPolicyMember_PostAct.xml` (ASM-FW-GISFW-DATA-POLICY!UPLOADCSVPOLICYMEMBER_POSTACT) | 22 dari 22 langkah akar, L16763 | `call InputDtlPayment_PreAct`, tanpa prakondisi |
| `UploadCSVPolicyMemberEDM_PostAct.xml` (…!UPLOADCSVPOLICYMEMBEREDM_POSTACT) | 26 dari 42, L27688 | sama, tanpa prakondisi |
| `InputDtlPayment_PreAct.xml` | 45, L10293 | `Call PremiPaymentMarine`, prakondisi `IsMarineCargo` |
| `PremiPaymentMarine.xml` (ASM-FW-GISFW-WORK!PREMIPAYMENTMARINE) | 1 (`IsNB`) | 1.1 `.Policy.Payment.Premium = @sum(.CargoList().CoverageList().Premium)`, `.Policy.Payment.Diskon = @sum(.CargoList().CoverageList().Discount)`; 1.2–1.4 berlabel `//` (remark, butir 43) |
| | 2 (`IsEDM`) | selisih terhadap `.OldCoverage(1)`, komisi, total cicilan — **cabang EDM, tidak di tiket ini** |

⚠️ `IsNB` membaca `pyWorkPage.Quotation.StatusBusiness == 1`, sedangkan `IsEDM` membaca
`pyWorkPage.OfferFacIn.QuotationData.StatusBusiness == 3` — **dua halaman berbeda**, jadi keduanya bisa benar sekaligus.

**Di luar tiket ini:** isi lain kedua `UploadCSV*_PostAct` (812 KB dan 1,7 MB; parsing CSV dan anggota polis — belum habis
dibaca) dan 53 langkah akar lain `InputDtlPayment_PreAct` (54 langkah akar: pembongkar dan urai `pySteps`); cabang EDM
`PremiPaymentMarine`.

**Blocked by:** 18

**Status:** ready-for-human — layanan diport 01-10-2026; nilai lama tidak terekam di fixture (belum tercakup); A45 menunggu konfirmasi

- [x] Gerbang `IsMarineCargo` (InputDtlPayment_PreAct langkah 45) dan `IsNB` / `IsEDM` (PremiPaymentMarine 1, 2) lewat
      registry predikat NB
- [x] Premium dan Diskon = jumlah eksak seluruh coverage seluruh kargo
- [x] Cabang EDM → galat "belum diport", tidak dihitung diam-diam
- [x] Nilai kosong / mata uang campuran ditolak, bukan ditebak
- [x] Rekonsiliasi: nilai lama `pyWorkPage.Policy.Payment` dicari di fixture; bila tidak ada, dinyatakan belum tercakup

## Comments

### 2026-10-01 — layanan diport (agent)

**Kode:** `backend/services/pembayaran/marine.go` (`Marine`, `Coverage`, `Hasil`, `ErrCabangEDM`, `ErrTanpaCabang`,
`ErrMasukan`); baris rekonsiliasi `PremiPaymentMarine langkah 1.1 Policy.Payment.Premium` di
`rekonsiliasi.barisPembayaran`.

**Rekonsiliasi:** `nb-marinecargo-1` → sistem baru `26395202.0000` (104 coverage IDR, diskon nol). Jumlah dihitung dua
cara (Python `decimal` dan `fractions`): 26.395.202. Status **belum tercakup**: `.Policy.Payment` di rule ini relatif
terhadap `pyWorkPage`, sedangkan akar fixture adalah halaman OfferFacIn; `CurrencyList(1).Policy.Payment.Premium`
(`1319760.1`) adalah halaman lain dan tidak setara.

**Temuan:** `IsNB` = `pyWorkPage.Quotation.StatusBusiness == 1` → kasus **renewal** MARINE (StatusBusiness 2) tidak
melewati cabang mana pun di rule NB ini (`ErrTanpaCabang`). `[pertanyaan terbuka]` apakah renewal memakai salinan rule di
folder RNW atau memang tidak menghitung Payment.

**Keputusan agent:** A45 (lintas mata uang ditolak). **Uji mutasi:** termasuk di
`../alat/mutasi_tiket18_19.py` (log `mutasi_tiket18_19.log`): 7 mutan pembayaran, 6 tertangkap, 1 ekuivalen.

### 2026-10-02 — batas bukti (agent, register "Ralat kedua" butir 3 dan "Ralat ketiga" butir 4)

- `IsNB` dan `IsEDM` membaca kolom yang sama di fixture (adaptor `kasusJSON`) — "keduanya benar" hanya teruji buatan.
- `.CargoList()` di rule ini = `pyWorkPage.CargoList`; fixture memberi `OfferFacIn.CargoList` — kesetaraan `[dugaan]`.
- Tes mata uang campur kini `errors.Is(…, uang.ErrMataUangBerbeda)`; mata uang kosong ditolak (A46).

