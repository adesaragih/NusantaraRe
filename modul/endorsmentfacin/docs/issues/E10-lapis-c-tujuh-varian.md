# E10: Lapis C — penanda baris warisan, tujuh varian

**What to build:** Setiap simpul daftar yang berasal dari polis lama ditandai, sehingga baris warisan
dapat dibedakan dari baris yang ditambahkan saat endorsement.

`[terverifikasi]` Tujuh varian, satu per lini bisnis. **Cacah penanda mengikuti kedalaman model data
lini itu**, bukan kelengkapan implementasi — lini dengan empat tingkat daftar bersarang menandai
lebih banyak simpul daripada lini dengan satu tingkat.

⚠️ Nilai penanda tersimpan **berkutip**. Pola pencarian tanpa kutip memberi **nol hasil** dan pernah
menyesatkan.

**Asal (Pega).** `Activity\SetOLDValueToEDMWork_{FIRE,MC,Aneka,MBU,LIFE,PA,GOLF}` ·
dipanggil `Activity\SetValueToEDMWork` langkah 14.7–14.13

**Keputusan.** K-047 · K-046 · `CLAUDE.md` §4.6

**Blocked by:** E06

**Status:** sebagian — 01-10-2026, `backend/services/lapisc.go`; test `lapisc_test.go`. ⛔ Korpus
menyimpang dari tiket di empat titik (laporan §2): tiap varian juga **menyalin** daftarnya dari
OldData dan **menggabungkan spreading** per `TreatyType`; `FlagOldData` dipasang **MC, MBU, PA, dan
LIFE**; FIRE berakhir di `SpreadingList`, bukan `LayerList`. Kode mengikuti korpus. Menunggu work
owner untuk butir K-046 jiwa.

- [x] Ketujuh varian menandai simpul sesuai **kedalaman daftar** lini masing-masing — kedalaman menurut korpus, lihat laporan §2.3
- [x] Nilai penanda sesuai yang tersimpan di korpus, **berkutip** — isi string `old`
- [x] **Hanya varian kebakaran** yang menyetel penanda prorata tingkat agregat; enam lainnya tidak
- [x] Penanda prorata itu **tidak** termasuk field yang disalin lapis A
- [ ] **K-046** `K046_Life_DuaPenandaOldData` — ⛔ **premis dibantah korpus**: `FlagOldData` + `IsOldData` dipasang juga oleh MC, MBU, dan PA (`TestLapisCFlagOldDataBukanHanyaLife`). Menunggu work owner: ralat tiket atau tunjukkan bukti lain
- [x] ⚠️ Lapis C menangani **jiwa tetapi bukan travel** — himpunannya berbeda dari lapis B (E08)
- [x] Menyebut rule Pega asalnya dalam komentar (§4.6)
