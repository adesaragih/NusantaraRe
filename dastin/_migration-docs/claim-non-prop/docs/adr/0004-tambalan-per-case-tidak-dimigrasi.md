---
status: accepted
label: DECIDED
---

# Tambalan per-case tidak ikut dimigrasi; angkanya dipindahkan sebagai data

<!-- STEMPEL ASAL -->
> **Dasar bukti**: ekspor XML `D:\XML_NURE\Claim Non Prop`, **279 berkas**, berkas tertanggal 2026-09-08 s/d 2026-09-09, rule termutakhir di dalamnya `pxUpdateDateTime = 2026-08-30`.
> Dokumen ini hanya berlaku untuk keadaan sistem pada ekspor tersebut. Tambalan yang ditambahkan sesudahnya tidak tercermin di sini; deteksinya lewat sapuan ulang, bukan lewat register.

Sapuan seluruh folder menemukan **9 identitas klaim atau master treaty yang ditulis langsung di dalam kode kalkulasi**, tersebar di 4 rule dan 14 langkah — antara lain `CLMNP-975` yang menimpa Premi Pemulihan dengan angka mati `881928.966808370` (IDR) dan `806851161.1895` (USD), serta `CLMNP-232` yang memetakan Biaya Penilaian secara manual per indeks. Kami memutuskan tidak satu pun dari cabang ini dibawa ke sistem baru: kodenya dibuang, tetapi nilai akhirnya dipertahankan sebagai data sehingga saldo klaim yang bersangkutan tetap benar.

## Consequences

Daftar lengkap ke-9 identitas beserta nilai yang ditimpa dan rumus normal yang di-bypass ada di `BLUEPRINT.md` §7, dan menjadi instruksi bagi tim migrasi data — bukan bagi tim pembangun aplikasi.

Salah satu temuan memperkuat bahwa ini memang tambalan, bukan aturan: pada `InputOutStandingClmTNP_PreAct` langkah 5, deskripsi langkah berbunyi `pyWorkPage.pyID=="CLMNP-50"` sementara kondisi eksekusinya `pyWorkPage.pyID=="CLMNP-232"` — deskripsi dan kondisi tidak sinkron, ciri khas salin-tempel.

Pemeriksaan ulang terhadap `MEMORI_PEMAHAMAN.MD` tidak menemukan keterangan bahwa treaty `1000393` punya perlakuan bisnis khusus. Karena itu ia diperlakukan sebagai tambalan, bukan sebagai dimensi yang hilang dari model data. Bila pemilik proses kemudian menyatakan sebaliknya, keputusan ini harus ditinjau ulang dan model data perlu menambah atribut yang menjelaskan perbedaan perlakuan tersebut.

## Tambahan 18 September 2026 — persetujuan eksplisit dan mekanisme pengganti

Persetujuan diberikan tanpa menunggu status terbuka/tutup tiap klaim: **tidak satu pun dari 29 langkah tambalan dimigrasi.**

Penggantinya satu mekanisme tunggal: **koreksi bernilai tercatat** — `nilai_sebelum`, `nilai_sesudah`, `alasan`, `pelaku`, `waktu`. Satu bentuk untuk semua kasus, bukan cabang per klaim di dalam kode.

Status tiap klaim yang disebut di dalam tambalan (`CLMNP-232`, `CLMNP-975`, `CLMNP-861`, `CLMNP-50`, `CLMNP-367`, `CLMNP-382`, `IDMaster 1000393`, `1001130`) diukur lewat REQ-015, dan itu **bukan penghalang** keputusan ini.

Satu hal yang dapat disimpulkan tanpa menunggu query: `CLMNP-975` ditambal **2026-07-16**, dua bulan sebelum ekspor ini. Tambalan tidak ditulis untuk klaim yang sudah tutup, jadi klaim itu hampir pasti masih terbuka dan tambalannya masih aktif.

**Karena itu `CLMNP-975` dijadikan kasus uji utama shadow-run.** Bila sistem baru menghasilkan angka yang benar untuk klaim itu **tanpa** tambalan apa pun, tesis "ini tambalan data, bukan aturan bisnis" terbukti untuk seluruh kelompoknya sekaligus — dan 29 langkah itu gugur bersama-sama, bukan satu per satu.
