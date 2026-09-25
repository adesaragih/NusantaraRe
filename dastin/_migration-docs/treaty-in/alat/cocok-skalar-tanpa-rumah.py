# -*- coding: utf-8 -*-
"""
cocok-skalar-tanpa-rumah.py - mencocokkan SKALAR KEDALAMAN-1 pohon clipboard
`TreatyIn` terhadap seluruh kolom *Asal* yang sudah diberi rumah di §10
`SPEC-MODEL-DATA.md`, KE DUA ARAH.

KENAPA SUMBER INI, DAN BUKAN §3.2
---------------------------------
Blok koreksi §10.2 menyatakan kelima belas jalur §3.2 yang menganggur sudah
diadili satu per satu dan TIDAK menjelaskan empat atribut yang hilang. Maka
sumber yang dipakai di sini harus MANDIRI dari §3.2.

`datar-treatyin-lama.csv` dibangkitkan `buat-pohon-treatyin.py` dari sapuan 329
berkas ekspor - bukan dari §3.2, bukan dari §10. F-8 sudah membuktikan semesta
§3.2 bolong: `BrokeragePct` tidak ada di dalamnya sama sekali.

Skalar kedalaman-1 dipakai karena ia setara dengan "atribut tingkat versi":
anak langsung halaman `TreatyIn` yang bukan Page maupun Page List.

DUA ARAH, DAN KEDUANYA DICETAK
------------------------------
  ARAH 1  skalar pohon yang TIDAK muncul sebagai Asal mana pun di §10
  ARAH 2  Asal di §10 yang TIDAK ada sebagai skalar kedalaman-1 di pohon

Arah 2 hampir pasti tidak kosong dan itu WAJAR - banyak Asal menunjuk ruas
bersarang atau bertanda *baru*. Ia dicetak supaya hasil bersih di arah 1 tidak
terbaca sebagai bukti bahwa pencocokannya sempurna.

BATAS KEMAMPUAN - dinyatakan di dalam keluarannya sendiri
---------------------------------------------------------
Pencocokannya LEKSIKAL atas nama ruas terakhir. Ia dapat menyatakan sebuah nama
TIDAK MUNCUL di kolom Asal mana pun; ia TIDAK dapat menyatakan bahwa ruas itu
karena itu tidak punya rumah - ia mungkin punya rumah dengan nama Asal yang
ditulis berbeda, atau sengaja dibuang dengan alasan yang tertulis di prosa.
Setiap baris keluarannya adalah CALON, bukan putusan.
"""
import io, os, re, sys

AKAR = r"D:\XML_NURE\_migration-docs\treaty-in"
POHON = os.path.join(AKAR, "4-erd-dan-tabel-datar", "datar-treatyin-lama.csv")
SPEC = os.path.join(AKAR, "SPEC-MODEL-DATA.md")

RE_BACKTIK = re.compile(r"`([^`]+)`")
RE_KEPALA = re.compile(r"^\|\s*Nama\s*\|")
RE_MULAI10 = re.compile(r"^### 10\.")
RE_AKHIR = re.compile(r"^## 11\.")

# ruas yang memang bukan atribut versi - disebut supaya penolakannya terbaca
BUKAN_ATRIBUT = {"ActualValue", "ValueDifference", "OLDDATA", "ValueBeforeProrate"}


def skalar_kedalaman1():
    hasil, ditolak = [], {}
    with io.open(POHON, encoding="utf-8-sig") as f:
        kepala = f.readline().rstrip("\n").split(",")
        i_jalur, i_dalam, i_nama, i_jenis = (kepala.index(k) for k in
                                             ("JALUR", "KEDALAMAN", "NAMA", "JENIS"))
        for ln in f:
            sel = ln.rstrip("\n").split(",")
            if len(sel) <= i_jenis:
                ditolak["baris tidak terurai"] = ditolak.get("baris tidak terurai", 0) + 1
                continue
            if sel[i_dalam] != "1":
                continue
            jenis, nama = sel[i_jenis], sel[i_nama]
            if jenis != "Skalar":
                ditolak["bukan Skalar (Page / Page List)"] = ditolak.get(
                    "bukan Skalar (Page / Page List)", 0) + 1
                continue
            if nama in BUKAN_ATRIBUT:
                ditolak["pohon cermin"] = ditolak.get("pohon cermin", 0) + 1
                continue
            hasil.append(nama)
    return hasil, ditolak


def asal_di_spec():
    asal, dalam10 = set(), False
    with io.open(SPEC, encoding="utf-8") as f:
        baris = f.read().splitlines()
    i = 0
    while i < len(baris):
        ln = baris[i]
        if RE_MULAI10.match(ln):
            dalam10 = True
        if RE_AKHIR.match(ln):
            break
        if dalam10 and RE_KEPALA.match(ln):
            i += 2
            while i < len(baris) and baris[i].lstrip().startswith("|"):
                sel = [s.strip() for s in baris[i].strip().strip("|").split("|")]
                if len(sel) >= 2:
                    for nm in RE_BACKTIK.findall(sel[1]):
                        for potong in re.split(r"[\[\]\.\s/]+", nm):
                            if potong:
                                asal.add(potong)
                i += 1
            continue
        i += 1
    return asal


def main():
    skalar, ditolak = skalar_kedalaman1()
    asal = asal_di_spec()

    tanpa_rumah = sorted(n for n in skalar if n not in asal)

    print("=" * 96)
    print("SKALAR KEDALAMAN-1 `TreatyIn` YANG TIDAK MUNCUL SEBAGAI *Asal* MANA PUN DI §10")
    print("=" * 96)
    print("skalar kedalaman-1 diperiksa : %d" % len(skalar))
    print("nama *Asal* terkumpul dari §10: %d" % len(asal))
    print("TANPA RUMAH (calon)           : %d" % len(tanpa_rumah))
    print("-" * 96)
    for n in tanpa_rumah:
        print("  %s" % n)
    print("-" * 96)
    print("YANG TIDAK MASUK PEMERIKSAAN, dan atas dasar apa:")
    for sebab, n in sorted(ditolak.items()):
        print("  %-36s %d" % (sebab, n))
    print()
    print("BATAS PERKAKAS:")
    print("  Pencocokan LEKSIKAL atas nama ruas. Ia menyatakan sebuah nama TIDAK MUNCUL di kolom")
    print("  Asal mana pun - BUKAN bahwa ruas itu tidak punya rumah. Ia mungkin berumah dengan")
    print("  nama Asal yang ditulis berbeda, atau dibuang dengan alasan yang tertulis di prosa.")
    print("  Setiap baris di atas adalah CALON, bukan putusan.")
    return 0


sys.exit(main())
