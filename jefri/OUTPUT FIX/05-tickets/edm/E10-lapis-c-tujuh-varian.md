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

**Status:** blocked

- [ ] Ketujuh varian menandai simpul sesuai **kedalaman daftar** lini masing-masing
- [ ] Nilai penanda sesuai yang tersimpan di korpus, **berkutip**
- [ ] **Hanya varian kebakaran** yang menyetel penanda prorata tingkat agregat; enam lainnya tidak
- [ ] Penanda prorata itu **tidak** termasuk field yang disalin lapis A
- [ ] **K-046** `K046_Life_DuaPenandaOldData` — varian jiwa menyetel **dua** penanda berbeda; enam lainnya satu
- [ ] ⚠️ Lapis C menangani **jiwa tetapi bukan travel** — himpunannya berbeda dari lapis B (E08)
- [ ] Menyebut rule Pega asalnya dalam komentar (§4.6)
