# -*- coding: utf-8 -*-
"""
periksa-cacah-atribut.py - mencocokkan CACAH YANG DIUMUMKAN di judul setiap §10.x
`SPEC-MODEL-DATA.md` terhadap CACAH YANG BENAR-BENAR TERDAFTAR di tabel atributnya.

Lahir dari F-2 dan F-3: §10.2 mengaku 48 memuat 44, §10.21 mengaku 6 memuat 5.
Tiga instans di satu berkas bukan kebetulan, jadi polanya disapu, bukan ditambal
satu per satu.

DUA CACAH, DAN KEDUANYA DICETAK
-------------------------------
  BARIS  - jumlah baris data di tabel atribut
  NAMA   - jumlah nama kolom berbacktick di kolom pertama

Keduanya berbeda ketika satu baris memuat lebih dari satu nama - misalnya
`NILAI_SEBELUM` / `NILAI_SESUDAH` di §10.21. Perkakas ini TIDAK memilih salah
satunya; ia mencetak keduanya, sebab yang satu menjawab "berapa baris ditulis"
dan yang lain menjawab "berapa kolom akan berdiri di DDL".

BATAS KEMAMPUAN - dinyatakan di dalam keluarannya sendiri
---------------------------------------------------------
Perkakas ini TIDAK dapat menyatakan sebuah cacah BENAR. Ia hanya dapat
menyatakan bahwa dua pernyataan di dalam berkas yang sama COCOK, atau TIDAK
COCOK. Bila judul dan tabel sama-sama salah, ia diam - dan diamnya bukan
pernyataan bahwa keduanya benar.

Ia juga tidak membaca angka yang ditulis di luar judul, misalnya di dalam blok
kutip koreksi. Angka semacam itu dilaporkan sebagai "angka lain di judul" hanya
bila ia berada di judul.
"""
import io, re, sys, collections

BERKAS = r"D:\XML_NURE\_migration-docs\treaty-in\SPEC-MODEL-DATA.md"

RE_JUDUL   = re.compile(r"^### (10\.\d+[a-z]?)\s+(.*)$")
RE_AKHIR   = re.compile(r"^## 11\.")
RE_ATRIBUT = re.compile(r"(\d+)\s*(?:\*\*)?\s*atribut", re.I)
RE_ANGKA   = re.compile(r"(?<![\w.])(\d{1,3})(?![\w.])")
RE_KEPALA  = re.compile(r"^\|\s*Nama\s*\|")
RE_BACKTIK = re.compile(r"`([^`]+)`")


def baca_bagian(baris):
    """Pecah menjadi (kode, judul, isi[]) per ### 10.x, berhenti di ## 11."""
    bagian, kini = [], None
    for ln in baris:
        if RE_AKHIR.match(ln):
            break
        m = RE_JUDUL.match(ln)
        if m:
            if kini:
                bagian.append(kini)
            kini = [m.group(1), m.group(2).rstrip(), []]
        elif kini:
            kini[2].append(ln)
    if kini:
        bagian.append(kini)
    return bagian


def tabel_atribut(isi):
    """Kembalikan (jumlah_baris, daftar_nama) dari SELURUH tabel ber-kepala '| Nama |'."""
    baris_total, nama_total, n_tabel = 0, [], 0
    i = 0
    while i < len(isi):
        if RE_KEPALA.match(isi[i]):
            n_tabel += 1
            i += 2                      # lewati kepala + garis pemisah
            while i < len(isi) and isi[i].lstrip().startswith("|"):
                sel = [s.strip() for s in isi[i].strip().strip("|").split("|")]
                if sel and sel[0]:
                    baris_total += 1
                    nama_total.extend(RE_BACKTIK.findall(sel[0]))
                i += 1
        else:
            i += 1
    return baris_total, nama_total, n_tabel


def main():
    with io.open(BERKAS, encoding="utf-8") as f:
        baris = f.read().splitlines()

    bagian = baca_bagian(baris)
    hasil, ditolak = [], collections.OrderedDict()

    for kode, judul, isi in bagian:
        m = RE_ATRIBUT.search(judul)
        diumumkan = int(m.group(1)) if m else None
        lain = [int(x) for x in RE_ANGKA.findall(judul)
                if diumumkan is None or int(x) != diumumkan]
        lain = [x for x in lain if not (10 <= x <= 10 and kode.startswith("10."))]
        n_baris, nama, n_tabel = tabel_atribut(isi)

        if n_tabel == 0:
            ditolak.setdefault("tanpa tabel atribut", []).append(kode)
            continue
        if diumumkan is None:
            ditolak.setdefault("judul tidak menyebut 'N atribut'", []).append(kode)

        hasil.append((kode, judul, diumumkan, lain, n_baris, len(nama), nama))

    lebar = max(len(h[1]) for h in hasil)
    lebar = min(lebar, 58)
    print("=" * 110)
    print("COCOK-SILANG CACAH ATRIBUT  ·  SPEC-MODEL-DATA.md §10")
    print("=" * 110)
    print("%-8s %-*s %6s %6s %6s  %s" % ("§", lebar, "judul", "UMUM", "BARIS", "NAMA", "hasil"))
    print("-" * 110)

    n_cocok = n_beda = 0
    for kode, judul, diumumkan, lain, n_baris, n_nama, nama in hasil:
        j = judul if len(judul) <= lebar else judul[:lebar - 1] + "~"
        if diumumkan is None:
            tanda = "? judul tanpa cacah"
        elif diumumkan == n_baris == n_nama:
            tanda = "cocok"; n_cocok += 1
        else:
            bagian_beda = []
            if diumumkan != n_baris:
                bagian_beda.append("umum%+d baris" % (n_baris - diumumkan))
            if n_baris != n_nama:
                bagian_beda.append("baris%+d nama" % (n_nama - n_baris))
            tanda = "TIDAK COCOK -> " + ", ".join(bagian_beda)
            n_beda += 1
        ekstra = ("  [angka lain di judul: %s]" % ",".join(map(str, lain))) if lain else ""
        print("%-8s %-*s %6s %6d %6d  %s%s" % (
            kode, lebar, j,
            "-" if diumumkan is None else diumumkan,
            n_baris, n_nama, tanda, ekstra))

    print("-" * 110)
    print("RINGKASAN : %d bagian bertabel · %d cocok · %d TIDAK COCOK" % (len(hasil), n_cocok, n_beda))
    print()
    print("YANG TIDAK MASUK HASIL, dan atas dasar apa:")
    if not ditolak:
        print("  (tidak ada)")
    for sebab, daftar in ditolak.items():
        print("  %-34s %d  -> %s" % (sebab, len(daftar), ", ".join(daftar)))
    print()
    print("BATAS PERKAKAS:")
    print("  Ia TIDAK dapat menyatakan sebuah cacah BENAR. Ia hanya menyatakan dua pernyataan")
    print("  di dalam berkas yang sama COCOK atau TIDAK COCOK. Bila judul dan tabel sama-sama")
    print("  salah, ia DIAM - dan diamnya bukan pernyataan bahwa keduanya benar.")
    return 0


sys.exit(main())
