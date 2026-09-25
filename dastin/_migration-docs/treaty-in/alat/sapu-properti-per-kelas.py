# -*- coding: utf-8 -*-
"""
sapu-properti-per-kelas.py — 24 September 2026

KENAPA PERKAKAS INI ADA
-----------------------
`buat-pohon-treatyin.py` menyusun pohon clipboard dengan satu penjaga yang BENAR:

    `Primary` hanya berarti TreatyIn bila ATURANNYA SENDIRI applies-to kelas akar.

Tanpa penjaga itu, simpul palsu masuk ke tingkat akar. Tetapi penjaga itu punya
akibat yang tidak pernah diperiksa: properti yang ditulis aturan ber-applies-to
KELAS ANAK **tidak dipindahkan ke bawah kelas anaknya** — ia dibuang sama sekali.

Akibatnya terbukti: `BreakDownSprdList` dan `BreakDownSprdListXOL` muncul 54 kali
di ekspor (6 berkas Treaty In) dan **NOL kali** di pohon 985 simpul. Keduanya
adalah sumber entitas `RINCIAN_PENYEBARAN`, yang disebut `SPEC-MODEL-DATA.md` §3.5.

Perkakas ini memakai SUMBER YANG MANDIRI dari penyisiran jalur: entri indeks
rujukan aturan, yang memasangkan

    <pxRuleClassName>  kelas pemilik properti
    <pyRuleName>       nama propertinya
    <pxRuleObjClass>   = Rule-Obj-Property

sehingga daftar properti per kelas Pega diperoleh TANPA melewati awalan jalur
`TreatyIn.` sama sekali. Bandingkan hasilnya dengan pohon, dan selisihnya adalah
titik buta pohon — terukur, bukan dugaan.

KELUARAN
    datar-properti-per-kelas.csv    seluruh properti per kelas
    datar-titik-buta-pohon.csv      yang ada di kelas tetapi TIDAK ada di pohon
"""
import csv
import io
import os
import re
import sys

AKAR_XML = r'D:\XML_NURE'
FOLDER = ['Treaty In', 'Treaty In Adjustment']
KELUARAN = r'D:\XML_NURE\_migration-docs\treaty-in\4-erd-dan-tabel-datar'
POHON = os.path.join(KELUARAN, 'datar-treatyin-lama.csv')

AWALAN_KELAS = 'ASM-FW-GISFW-'
# Properti bawaan Pega — bukan model bisnis, tidak dihitung.
BAWAAN = re.compile(r'^(px|py|pz)')

RE_ROW = re.compile(r'<rowdata\b[^>]*>(.*?)</rowdata>', re.S)
RE_KELAS = re.compile(r'<pxRuleClassName>(.*?)</pxRuleClassName>', re.S)
RE_NAMA = re.compile(r'<pyRuleName>(.*?)</pyRuleName>', re.S)
RE_JENIS = re.compile(r'<pxRuleObjClass>(.*?)</pxRuleObjClass>', re.S)


def sapu_berkas(path):
    """Kembalikan himpunan (kelas, properti) dari entri indeks berkas ini."""
    hasil = set()
    try:
        s = io.open(path, encoding='utf-8', errors='replace').read()
    except OSError:
        return hasil
    if 'Rule-Obj-Property' not in s:
        return hasil
    for m in RE_ROW.finditer(s):
        blok = m.group(1)
        j = RE_JENIS.search(blok)
        if not j or j.group(1).strip() != 'Rule-Obj-Property':
            continue
        k = RE_KELAS.search(blok)
        n = RE_NAMA.search(blok)
        if not k or not n:
            continue
        kelas = k.group(1).strip()
        nama = n.group(1).strip()
        if not kelas.startswith(AWALAN_KELAS):
            continue
        if BAWAAN.match(nama):
            continue
        hasil.add((kelas, nama))
    return hasil


def baca_pohon():
    """
    Dari pohon: untuk tiap kelas, himpunan NAMA ANAK yang pernah muncul di
    bawah simpul berkelas itu. Itulah yang pohon "tahu" tentang kelas tersebut.
    """
    if not os.path.exists(POHON):
        sys.exit('pohon tidak ada: ' + POHON)
    baris = list(csv.DictReader(io.open(POHON, encoding='utf-8-sig')))
    kelas_jalur = {}
    for r in baris:
        kelas_jalur[r['JALUR']] = r['KELAS_PEGA'].strip()
    anak = {}
    semua_nama = set()
    for r in baris:
        jalur = r['JALUR']
        semua_nama.add(r['NAMA'])
        induk = jalur.rsplit('.', 1)[0]
        kelas_induk = kelas_jalur.get(induk, '')
        if kelas_induk and kelas_induk != '(tidak dideklarasikan)':
            anak.setdefault(kelas_induk, set()).add(r['NAMA'])
    return anak, semua_nama


def main():
    pasangan = set()
    cacah_berkas = 0
    for f in FOLDER:
        akar = os.path.join(AKAR_XML, f)
        if not os.path.isdir(akar):
            continue
        for d, _, berkas in os.walk(akar):
            for b in berkas:
                if not b.lower().endswith('.xml'):
                    continue
                cacah_berkas += 1
                pasangan |= sapu_berkas(os.path.join(d, b))

    anak_pohon, semua_nama_pohon = baca_pohon()

    per_kelas = {}
    for kelas, nama in pasangan:
        per_kelas.setdefault(kelas, set()).add(nama)

    # 1. seluruh properti per kelas
    p1 = os.path.join(KELUARAN, 'datar-properti-per-kelas.csv')
    with io.open(p1, 'w', encoding='utf-8-sig', newline='') as fh:
        w = csv.writer(fh)
        w.writerow(['KELAS_PEGA', 'CACAH_PROPERTI', 'PROPERTI'])
        for kelas in sorted(per_kelas):
            w.writerow([kelas, len(per_kelas[kelas]),
                        ' | '.join(sorted(per_kelas[kelas]))])

    # 2. titik buta: properti kelas yang tidak pernah muncul sebagai anak
    #    simpul berkelas itu di pohon.
    p2 = os.path.join(KELUARAN, 'datar-titik-buta-pohon.csv')
    buta = []
    for kelas in sorted(per_kelas):
        tahu = anak_pohon.get(kelas, set())
        for nama in sorted(per_kelas[kelas] - tahu):
            # ADA_DI_POHON_DI_TEMPAT_LAIN membedakan dua hal yang sangat berbeda:
            #   tidak ada sama sekali  -> benar-benar hilang
            #   ada di tempat lain     -> hanya tidak tergantung di kelas ini
            buta.append([kelas, nama,
                         'ya' if nama in semua_nama_pohon else 'TIDAK'])
    with io.open(p2, 'w', encoding='utf-8-sig', newline='') as fh:
        w = csv.writer(fh)
        w.writerow(['KELAS_PEGA', 'PROPERTI', 'ADA_DI_POHON_DI_TEMPAT_LAIN'])
        w.writerows(buta)

    hilang_total = [b for b in buta if b[2] == 'TIDAK']
    print('berkas disapu              :', cacah_berkas)
    print('kelas dengan properti      :', len(per_kelas))
    print('pasangan (kelas, properti) :', len(pasangan))
    print('titik buta (semua)         :', len(buta))
    print('titik buta HILANG TOTAL    :', len(hilang_total))
    print()
    print('-> ' + p1)
    print('-> ' + p2)


if __name__ == '__main__':
    main()
