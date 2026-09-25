# 10: Jembatan predikat lini bisnis → resolver skala

**What to build:** Lini bisnis sebuah kasus **ditentukan oleh predikat yang sesungguhnya dipakai
sistem lama**, lalu mengalir ke resolver skala — bukan diisi pemanggil sebagai parameter yang
diasumsikan benar.

Ini titik integrasi antara registry predikat dan perhitungan premi: dua bagian yang sudah bekerja
sendiri-sendiri, disambungkan di satu tempat yang dapat diuji. Predikat yang terlibat mengenali
FIRE, PA, MBU, ANEKA, MARINE CARGO, GOLF, dan BONDING.

⚠️ **Kehati-hatian yang diwarisi dari discovery.** Lini bisnis ditentukan satu properti, **tetapi dua
predikat membacanya lewat jalur berbeda di dalam agregat** — sebagian memakai dua sampai tiga jalur.
Apakah ketiga jalur selalu sinkron **belum terjawab**. Sampai terjawab, jalur-jalur itu
**dipertahankan apa adanya** dan perbedaannya dicatat saat runtime; properti itu **tidak boleh
dinormalisasi** menjadi satu field.

⚠️ Sebagian nilai lini bisnis muncul **dengan spasi di depan atau belakang**. Dipangkas di batas
input, dan dicatat sebagai kandidat perbaikan.

**Blocked by:** 08, 03

**Status:** ready-for-agent

- [ ] Lini bisnis diturunkan lewat registry predikat, bukan diterima mentah dari pemanggil
- [ ] Hasilnya mengalir ke resolver skala; **tidak ada** jalur yang melewati resolver
- [ ] Jalur baca ganda pada properti penentu **dipertahankan**, tidak dinormalisasi
- [ ] Ketidaksinkronan antar jalur **tercatat saat runtime** alih-alih dipilih diam-diam salah satunya
- [ ] Spasi di ujung nilai dipangkas di batas input, dicatat sebagai kandidat perbaikan
- [ ] Kasus yang lini bisnisnya tidak dikenali **gagal keras** lewat resolver (tiket 03), bukan memakai default
