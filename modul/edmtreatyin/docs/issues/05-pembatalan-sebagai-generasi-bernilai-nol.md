# 05: Pembatalan sebagai generasi bernilai nol

> ## ⭐ PENAHAN GUGUR — 23 September 2026
>
> `[keputusan work owner]` Polis yang sudah dibatalkan **masih boleh di-endorse lagi**. Nol gerbang peran — layar cukup memberi **peringatan**, bukan penghalang.
>
> ⭐ **Status berubah `blocked` → `ready-for-agent`.** Dua baris di kepala tiket dicoret, bunyinya tidak dihapus.
>
> Rinciannya di `modul/nbtreatyin/docs/KEPUTUSAN-RONDE-12-BUTIR-2026-09-23.md`, butir 4.

> ⛔ **KOREKSI 06-10-2026** (log `../KOREKSI-DOKUMEN-2026-10-06.md`). Bunyi lama AC 14 *"seluruh kolom
> uangnya **bernilai nol**"* (dan spec-penyimpanan ID-19 *"menolkan seluruh kolom uang"*) **tidak persis
> XML**. `Activity/SetEDMTCancel.xml`: langkah **1.1** (proporsional) menolkan **16** medan —
> `PremiOgp RiCommOgp ResultOgp1 OveriddingCommOgp ResultOgp2 PremiOnp RiCommOnp ResultOnp1
> OveriddingCommOnp ResultOnp2 Claim OutstandingClaim SalvageValue ExcessLoss Deduction1 Deduction2`;
> langkah **1.2** menyalin `Installment/ListInstallment/SpreadingRiskList/TreatyXOLList` dari data lama;
> **1.3** `PaymentTotal = Premium = 0`; **1.4** `PremiumSpreaded = SharePercentage = 0`. ⛔ **Tidak dinolkan**
> rule ini: `GrossPremium`, `NetPremium`, `BalanceDueTo`, `PPNValue`, `PPHValue`, `BalanceBeforeTax/PPH`,
> `ClaimSpreaded`, `ClaimPercentage`. Langkah **2.x** (NonProp) menolkan seluruh uang lapisan; mata uang,
> `DueTo`, lapisan disalin dari data lama. Panggilan yang efektif: `EDMChooseBusiness_Act` langkah **5**
> (`EDMType=='4'`); panggilan `CreateEDMT` langkah 16 berjalan atas halaman portal **sebelum** kasus
> dibuat (langkah 18) — `[dugaan]` tanpa efek. ⚠️ Tiru 16 medan XML atau nolkan seluruh uang — **butir WO**.
> Bab *Kenapa tiket ini `blocked`* di bawah dicoret: penahannya gugur 23-09.

---


**Status:** ⭐ **ready-for-agent** *(semula ~~blocked~~ — 23-09-2026)*
~~**Blocked by:** ⛔ `[work owner]` **sesudah dibatalkan, polis masih boleh di-endorse lagi atau tidak**~~ ⛔ **penahan gugur 23-09-2026**
**Bergantung pada tiket NB:** **19** *(pemecah dokumen — kolom jenis berkas)*
**Menutup:** AC **14–15** *(2 AC)*
**Sumber:** `spec-penyimpanan-relasional.md` ID-19 · ID-41

## Hasil & nilai pengguna

Pembatalan polis adalah **jenis endorsemen**, dipilih di awal saat berkas dibuat — bukan alur
tersendiri dan bukan tombol terpisah.

⛔⛔ **Pembatalan BUKAN penghapusan.** Ia menghasilkan **generasi baru berisi nol**, sehingga
jejaknya tetap ada dan dapat dipertanggungjawabkan.
⛔ **Tiket yang membuat jalur hapus untuk pembatalan salah.**

## Yang dibangun

Jenis endorsemen tersimpan di kolomnya sendiri, dan salah satu nilainya berarti pembatalan.
Memilihnya menghasilkan generasi baru yang seluruh kolom uangnya **bernilai nol**, lewat jalur
penyimpanan yang **sama** dengan endorsemen lain.

⛔ **Nol jalur hapus dibangun untuk pembatalan.**

## Batas — yang TIDAK termasuk

⛔ Aturan nomor urut — tiket **04**; pembatalan **tidak** menghapus baris, jadi aturannya berlaku apa adanya.
⛔ Perhitungan selisih terhadap generasi yang dibatalkan — tiket **06**.

## Cara mengujinya

Lewat seam `repository`. ⭐ **Uji utama:** polis dibatalkan, lalu **cacah baris rincian dihitung**
— jumlahnya harus **sama** dengan generasi sebelumnya, dan seluruh kolom uangnya nol.
⛔ Baris yang berkurang berarti gagal.

## Acceptance criteria

- [ ] **AC 14** — pembatalan menghasilkan generasi baru dengan kolom uang **bernilai nol**, bukan penghapusan *(koreksi 06-10: XML menolkan 16 medan + angsuran + spreading, bukan seluruh — lihat blok KOREKSI)*
- [ ] **AC 15** — jenis endorsemen tersimpan di kolomnya sendiri dan dipilih di awal

## ⛔ ~~Kenapa tiket ini `blocked`~~ — ✅ penahan gugur 23-09 *(dicoret koreksi 06-10)*

~~`[work owner]` **Belum dijawab apakah polis yang sudah dibatalkan masih boleh di-endorse lagi.**
Jawabannya menentukan apakah generasi pembatalan menjadi **ujung rantai** — sehingga penunjuk
generasi berikutnya harus ditolak — atau sekadar generasi biasa yang nilainya nol.~~

~~⭐ **Kedua AC dapat dibangun tanpa menunggu**; yang menunggu hanya **penjaga rantai sesudahnya**.~~
⭐ `[keputusan work owner]` 23-09: **boleh** di-endorse sesudah batal, layar cukup memberi peringatan
(spec-penyimpanan bab *KEPUTUSAN 23-09-2026* butir 3).