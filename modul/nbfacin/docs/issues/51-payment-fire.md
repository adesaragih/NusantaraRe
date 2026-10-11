# 51: Tab Payment kasus FIRE

> ⚠️ **Disusun agent sesi `nusantarare-0f` — bukan hasil `/to-tickets`.** Pemicu: work owner 06-10-2026 "lanjut ke tab
> payment" (tab Loss Record dilewati: di Pega hanya tampil untuk `IsIT` = `OperatorID.pyPosition = "IT Developer"`).
> Survei korpus oleh subagen (pengurai XML); butir kunci diperiksa ulang sesi 0f. Acuan tata letak = gambar layar Pega
> `D:\migrasi\RNM\DDL\Blok Currency, Payment, Installment.JPG`.

**Status:** frontend dibangun 06-10-2026 (sesi 0f); backend selesai 06-10-2026 (sesi c3, uji hijau) dengan ketiga
keputusan work owner diterapkan. Migrasi 201 ditulis, menunggu dijalankan work owner — bersama 200, sebelum biner baru.

## Bukti `[terverifikasi]` — `D:\migrasi\RNM\NB FacIn\` (F) dan `D:\migrasi\RNM\Endorsment Fac In\` (S)

- `F\Section\InputInwardFacultativeDtl.xml` tab "Payment" non-Life (±93952–96543): label "Payment Information", section
  `PaymentCurrencyList`, tombol Save → `IsThereAnyObjectLocation_Act(IsSave="Save")` lalu `ProtectPremiPolicy_Act`.
- `F\Section\PaymentCurrencyList.xml`: "Installment" (`.Policy.Payment.Installment`, wajib, `pyMax=2`) change/enter →
  `FillPaymentInstallment(Installment)`; grid `pyWorkPage.OfferFacIn.CurrencyList`: Currency (`.Name`) · Premium RNM
  (`.Policy.Payment.Premium`) · Total Payment (`.SumTotalPayment`), master-detail → `InputDtlPayment_FacIn`.
- `F\Section\InputDtlPayment_FacIn.xml` "Installment Data Information", grid `.Policy.Payment.ListInstallment`: No
  (`.InstallmentNo`) · Due Date* (`.DueDate`, dd/MM/yyyy, dapat diisi) · % Installment (`.InstallmentPercentage`, dapat
  diisi; change → `SetValidateInstallment_Act(index)`) · Premium RNM (`.Premium`) · Deduction (`.RICommision`) ·
  Deduction 2 (`.Deduction2`) · Total Brokerage Fee (`.BrokerageFee`) · PPH n % / PPN n % (`.PPh` / `.PPN`, kepala =
  `PPH_Note` / `PPN_Note`) · Discount (`.Discount`) · Policy Fee (`.AdminFee`, hanya `IsHealthIndividu` → tersembunyi
  untuk FIRE) · Total Payment (`.PaymentTotal`); pesan `InputParam.HASIL12` di atas grid.
- `FillPaymentInstallment` / `SetValidateInstallment_Act` kelas yang dipanggil layar **tidak ada di F** (F hanya versi
  `DATA-POLICYTREATYIN`); versi benar di S: `ASM-FW-GISFW-WORK!FILLPAYMENTINSTALLMENT` (bangun ulang daftar per mata uang:
  % = 100/n, terakhir = sisa; `.DueDate` = StartDateTime + 3 bulan, sama untuk semua baris — S @2693) dan
  `…-CURRENCY!SETVALIDATEINSTALLMENT_ACT` ("Total % Installment Can't Be More Than 100%" / "Total % Installment Must Be
  100%").
- `F\Activity\SaveFillPaymentInstallment.xml` + `CountPaymentInstallment_Act.xml` (Save): per cicilan jumlah = nilai mata
  uang × % / 100; `SumTotalPayment` = `NetPremium` = Premium − Commision − BrokerageFee + PPh + PPN
  (`CountPaymentInstallment_Act`, versi 2026-06-23 — Deduction 2 dan Discount TIDAK dikurangkan).
- `F\Activity\ProtectPremiPolicy_Act.xml`: "List Payment can't be null!", "List Installment can't be null!", "Premi on
  paymnet can't be 0", "Total premium RNM on installment didn't match with premium RNM!", "Premium Master Policy must be 0".
- Discount cicilan = `Policy.Payment.Diskon` mata uang, tidak pernah diisi di F → 0 (diskon coverage sudah dikurangkan
  di Premium RNM).

## Keputusan work owner 06-10-2026

1. **Total Payment dikurangi Deduction 2** — "ikut dikurangkan, mungkin kelewat untuk activity itu ... jangan halu".
   Dasar korpus `[terverifikasi]`: `SaveFillPaymentInstallment.xml` langkah 3.1.2 (@1064–1065, memo "tambah deduction 2",
   2026-03-16) `NetPremium = Premium − Commision − Deduction2 − BrokerageFee + PPh + PPN`; `CountPaymentInstallment_Act.xml`
   (memo "ubah …", 2026-06-23) menimpanya tanpa Deduction 2. Dipakai rumus 3.1.2 untuk `NetPremium`, `SumTotalPayment` dan
   `PaymentTotal` cicilan. Discount tidak dikurangkan (tidak ada di rumus Pega mana pun; nilainya 0 untuk NB).
2. `FillPaymentInstallment` / `SetValidateInstallment_Act` versi folder Endorsment Fac In dipakai — "setuju".
3. Due Date bawaan = mulai polis + 3 bulan, sama untuk semua cicilan — "iya ikut pega".

## Backend (sesi c3, 06-10-2026)

Kontrak final = kontrak usulan. `GET /api/nbfacin/kasus/{caseId}/payment` → `TampilanPayment` persis `api.ts` (larik selalu
`[]`, uang teks 8 desimal, `installmentPercentage` teks APA ADANYA sampai 20 desimal, Due Date `DD-MM-YYYY`, `""` bila
kosong). `POST …/payment/hitung` {aksi, currency, idx, installment, mataUang[{currency, cicilan[{installmentNo, dueDate,
installmentPercentage}]}]} → tampilan TANPA simpan. `PUT …/payment` (badan sama TANPA aksi / currency / idx) → hitung +
proteksi + simpan → bentuk GET. Kunci asing → 400. Galat: 400 isian / aksi / proteksi (`ErrMasukanPayment`, pesan verbatim
Pega), 401 PUT tanpa identitas, 404, 409 mata uang badan ≠ objek, 503.

- `backend/migrations/201_payment.sql` (+ `_down`): `T_GENERAL_POLIS.PAY_INSTALLMENT` NUMBER(5), `T_CURRENCYLIST.SUM_TOTAL_PAYMENT`
  / `PAY_NET_PREMIUM` NUMBER(38,8), tabel `T_LISTINSTALLMENT` (rancangan, sebagian) + indeks + sequence. **Ditulis, belum
  dijalankan.** ⚠️ Urutan deploy: **200 DAN 201 sebelum biner baru** — Save tab Premium Deduction kini juga menyentuh
  `T_LISTINSTALLMENT` (K51-7). Penjaga rancangan membaca 201 (272 kolom); `docs/STRUKTUR-TABEL-NB-FACIN.md`.
- `models/payment.go`, `repository/payment.go` + `repository/currencylist.go` (bersama tiket 50), `services/payment.go`,
  `handlers/payment.go` (+ uji ketiganya).

## Keputusan backend (sesi c3)

- **K51-1** (keputusan work owner 06-10-2026, "ikut dikurangkan … jangan halu") NetPremium = SumTotalPayment = Premium −
  Commision − **Deduction2** − BrokerageFee + PPh + PPN, lalu @Math.divide(x, 1, 20) — `NB FacIn\Activity\
  SaveFillPaymentInstallment.xml` 3.1.2 `[terverifikasi]`; PaymentTotal cicilan = NetPremium × % / 100. MENYIMPANG sadar dari
  `CountPaymentInstallment_Act` 3.2 (dan FillPaymentInstallment 5.1.2 / 5.1.4 untuk SumTotalPayment) yang membuang Deduction2.
  Discount tidak dikurangkan (0 untuk NB). Proteksi total memakai NetPremium ini.
- **K51-2** (keputusan work owner "setuju") Aksi installment / persen memakai `Endorsment Fac In\Activity\
  FillPaymentInstallment.xml` (ASM-FW-GISFW-WORK!FILLPAYMENTINSTALLMENT) dan `SetValidateInstallment_Act.xml`
  (ASM-FW-GISFW-DATA-OFFERFACIN-CURRENCY!SETVALIDATEINSTALLMENT_ACT) — varian modul lain, dipinjam atas persetujuan (CLAUDE.md
  §4 butir 6). DueDate bawaan = Begin date WIB + 3 bulan, SAMA untuk semua cicilan (keputusan "iya ikut pega"); tambah bulan
  mengikuti Calendar.add Java (31 Jan + 1 bulan = 28 Feb). Nilai uang cicilan hasil kedua aksi dihitung dengan rumus Save
  (Premium / RICommision / PaymentTotal skala 20, Deduction2 skala 4) — bukan skala 4 FillPaymentInstallment, dan bukan
  `PaymentTotal = Premium − RICommision` SetValidateInstallment langkah 6 — supaya layar menampilkan nilai yang akan disimpan.
- **K51-3** `INSTALLMENT_PERCENTAGE` disimpan teks desimal apa adanya (VARCHAR2(50), pola PCT_LIMIT): persen Pega berskala 20
  (`[terverifikasi]` contoh `DDL\CONTOH` "33.33333333333333333333"); NUMBER(38,8) membuat Σ ≠ 100 sehingga Save meratakan ulang
  persen pengguna. Due Date tersimpan `YYYYMMDD` (`[terverifikasi]` 200 dari 200 contoh, dihitung dua cara: pengurai XML = grep
  tag); kabel `DD-MM-YYYY`.
- **K51-4** GET menghitung ulang nilai uang cicilan dari persen tersimpan dan nilai Payment mata uang terkini (tab Premium
  Deduction) — setara Pega yang menjalankan SaveFillPaymentInstallment + CountPaymentInstallment_Act di setiap Save tab lain.
- **K51-5** (keputusan work owner "Reject, don't save") Proteksi `ProtectPremiPolicy_Act` → 400, TIDAK disimpan (Pega
  menyimpan dulu). Pesan verbatim, urut Page-Set-Messages: "List Payment can't be null!", "List Installment can't be null!",
  "Total premium RNM on installment didn't match with premium RNM!", "Premi on paymnet can't be 0" (ejaan Pega). "Premium
  Master Policy must be 0" (MOP / PolicyType 1) dan cabang Life / EDM di luar cakupan; syarat `PolicyType == 1` pesan premi 0
  tidak dibawa (PolicyType tersimpan sebagai label layar, butir 78.2).
- **K51-6** Proteksi total dibawa VERBATIM: When 2 langkah 5.5 berantai OR — pesan hanya bila KETIGA selisih (premi, deduction,
  total payment; 4 desimal) ≥ 0.001 (kemungkinan maksudnya AND; `[pertanyaan terbuka]` ke work owner). Karena Save meratakan
  persen sampai Σ = 100, praktis tidak pernah menyala.
- **K51-7** Save PUT: jumlah cicilan per mata uang wajib = Installment (400, pesan agent; Pega mengiterasi baris yang ada
  dengan Param.Installment), Due Date wajib (pyRequiredWhen !IsEDM), Installment 1–99 (pyMax 2). Persen dipertahankan bila Σ
  = 100 tepat, selain itu diratakan 100/n (SaveFill 3.1.5.3/3.1.5.4) dan No dinomori ulang. NetPremium negatif → 400 (pesan
  agent; pengikat angka Oracle modul ini tidak menerima tanda minus `[dugaan]`). CurrencyList ditulis menurut nama mata uang
  (`simpanCurrencyList`, bersama tiket 50): Save tab Premium Deduction memperbarui baris yang ada sehingga cicilan tidak ikut
  terhapus; mata uang yang hilang dibuang beserta cicilannya. Nilai Payment mata uang dibaca di luar transaksi tulis (dari
  tampilan tab Premium Deduction); setiap GET menghitung ulang (K51-4). Save tab Payment TIDAK menimpa PAY_* baris mata
  uang yang sudah ada (milik tab Premium Deduction; hanya baris baru diisi nilai tampilan) — temuan review standar. Kode mata
  uang > 10 bita → 400 (T_LISTINSTALLMENT.CURRENCY_CODE VARCHAR2(10) rancangan). Save tab Premium Deduction menyegarkan
  cicilan tersimpan (K50-12). Nama mata uang ganda di data lama tidak ditangani (baris kedua tidak tersentuh) —
  `[pertanyaan terbuka]` bila loader memuatnya.

## Frontend (sesi 0f)

- `components/TabPayment.tsx` (`pesanPersenCicilan`, `installmentSah`) + uji; `api.ts` (`ambilPayment`, `hitungPayment`,
  `simpanPayment`, `TampilanPayment`, `MataUangPayment`, `CicilanPayment`, `AksiPayment`); `labels.ts` (`PAYMENT`,
  `TEKS_PAYMENT`); `nbfacin.css`; `pages/InwardFacultative.tsx`.
- Keputusan agent: **P-1** hitung saat lepas fokus / Enter; **P-2** Due Date diubah tanpa memanggil server (Pega tanpa
  aksi), tersimpan lewat Save.
