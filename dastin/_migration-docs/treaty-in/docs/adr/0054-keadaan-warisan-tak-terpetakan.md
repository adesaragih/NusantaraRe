# ADR-0054 — Keadaan warisan yang tidak ada padanannya

**Status:** diterima, 23 September 2026
**Berlaku untuk:** Treaty In
**Menyelesaikan:** tabrakan antara ADR-0042 dan ADR-0046

## Konteks

Dua keputusan yang sudah diambil saling bertabrakan pada satu titik:

| ADR | Isinya |
|---|---|
| ADR-0042 | data warisan dipindahkan **apa adanya**, termasuk cacatnya; sejarah tidak dibersihkan |
| ADR-0046 | keadaan siklus hidup adalah **satu nilai bertipe tegas dengan himpunan tertutup** |

Titik tabrakannya nyata, bukan hipotetis. Aturan `TreatyInSetValue` di sistem lama menyetel
`StatusAkseptasi` ke nilai `"test"`, dan aturan itu terpasang di enam layar. Nilai `"test"` tidak
punya padanan di antara keadaan sah mana pun.

Memindahkannya apa adanya melanggar himpunan tertutup. Memetakannya ke salah satu keadaan sah
berarti membersihkan sejarah, yang ADR-0042 larang. Keduanya tidak bisa sekaligus dipenuhi tanpa
aturan tambahan.

Ini tabrakan pertama di antara keduanya dan hampir pasti bukan yang terakhir, karena `"test"`
bukan satu-satunya nilai liar yang mungkin ada di dua puluh tahun data. Karena itu ia diselesaikan
sebagai **aturan**, bukan sebagai penanganan satu kasus.

## Keputusan

Himpunan keadaan memuat **satu anggota tambahan yang menyatakan dirinya sendiri sebagai warisan tak
terpetakan**: `WARISAN_TAK_TERPETAKAN`.

Ia bukan "lainnya" dan bukan "tidak diketahui". Artinya tepat satu hal: *nilai warisan yang tidak
ada padanannya di antara keadaan sah, dan nilai aslinya tercatat.*

Tiga ketentuan yang menyertainya:

1. **Nilai asli disimpan.** Atribut `KEADAAN_WARISAN_ASLI` memuat teks keadaan lama apa adanya.
   Ia terisi **hanya** pada baris berkeadaan `WARISAN_TAK_TERPETAKAN`, dan tidak pernah dibaca oleh
   perhitungan mana pun — ia catatan, bukan masukan.

2. **Tidak ada perpindahan masuk selain migrasi.** Sistem berjalan tidak pernah bisa menghasilkan
   keadaan ini. Satu-satunya pintu masuknya adalah pemindahan data lama.

3. **Satu-satunya perpindahan keluar adalah `PERBAIKAN_WARISAN`.** Perpindahan ini menetapkan
   keadaan sah yang **dipilih secara eksplisit** oleh orang yang memperbaikinya — bukan ditebak
   sistem — dan tercatat di jejak perubahan (ADR-0045) beserta siapa dan kapan. Sesudah itu baris
   tersebut berjalan seperti baris lain.

## Konsekuensi

- Himpunan keadaan **tetap tertutup**: tidak ada nilai di luar daftar, karena yang di luar daftar
  punya satu tempat bernama.
- Sejarah **tetap tidak dibersihkan**: nilai aslinya ada, dapat dibaca, dan tidak ada yang
  menyulapnya menjadi keadaan yang tidak pernah dicapainya.
- Baris bermasalah **tidak bisa menyelinap ke alur kerja biasa**, karena ia tidak punya perpindahan
  keluar selain perbaikan yang sadar.
- Ini instans dari aturan sentuh-perbaiki pada ADR-0042: invarian baru ditegakkan saat sebuah
  kontrak disentuh, bukan lewat pembersihan massal.

## Berapa banyak, dan apakah itu mengubah keputusan

Jumlahnya belum diketahui; ia diukur oleh **Uji S-7** pada berkas permintaan DBA.

Hasilnya **tidak mengubah keputusan ini**. Bila nol, aturannya tetap ditulis, karena ia menutup
kelas persoalan dan bukan satu nilai. Yang berubah hanya perkiraan berapa banyak pekerjaan
perbaikan yang menunggu di masa depan.
