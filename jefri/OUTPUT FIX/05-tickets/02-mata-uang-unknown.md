# 02: Mata uang yang tidak diketahui sebagai keadaan eksplisit

**What to build:** Sistem baru dapat membaca dan menampilkan nilai uang yang **tidak membawa mata
uang**, sebagaimana sistem lama melakukannya hari ini — tanpa menebak, dan tanpa berhenti.

Dasarnya terukur: korpus memang kehilangan informasi ini (112 layar menampilkan uang tanpa field mata
uang mana pun), dan pengukuran produksi D1 menemukan **4 dari 691.925 baris** tiba tanpa mata uang.
Langka, tetapi nyata — dan justru kelangkaannya yang berbahaya: cacat yang muncul 4 kali dari 691.925
**tidak akan tertangkap pengujian sampel**.

Karena itu keadaan tidak-diketahui dibuat **eksplisit dan terlihat**, lalu gagal keras **hanya di
titik yang benar-benar memerlukan jawabannya**. Default diam-diam ke satu mata uang ditolak: itu
menebak, dan tebakannya akan salah persis pada kasus paling mahal.

**Blocked by:** 01

**Status:** ready-for-agent

- [ ] Keadaan mata uang tidak diketahui adalah **nilai sah yang eksplisit**, bukan kosong dan bukan kegagalan
- [ ] Baca, tampilkan, dan simpan di memori: **diizinkan**
- [ ] Aritmetika antar nilai yang sama-sama tidak diketahui mata uangnya: **diizinkan**
- [ ] Aritmetika **lintas mata uang** yang melibatkannya: **gagal keras**
- [ ] **Tidak ada** default diam-diam ke mata uang mana pun di kode
- [ ] Implementasi menghasilkan **hitungan** berapa banyak nilai tiba tanpa mata uang — sebagai ukuran, bukan dugaan
- [ ] Kewajiban **"gagal keras saat menulis ke Oracle"** dicatat sebagai TODO yang terikat pada tiket jalur tulis produksi; **tidak** diimplementasikan sekarang karena seam-nya belum ada, dan **tidak** dihapus dari daftar kewajiban
