# -*- coding: utf-8 -*-
"""
sapu-tetapan-di-kode.py — 24 September 2026

ATURAN YANG MELAHIRKANNYA
    "Setiap tetapan yang ditulis langsung di kode adalah CALON."

Pola "cacat bersembunyi di balik nilai bawaan" sudah muncul empat kali di modul ini:
deposit premium = minimum premium, persentase master berjumlah 100, rasio limit
menimpa tarif premi pemulihan, dan `ReinstatementPct = "100"` yang membuat ketentuan
pemulihan bertingkat MUSTAHIL dinyatakan.

Keempatnya berbentuk sama: sebuah nilai TETAP ditulis ke properti bisnis, sehingga dua
hal yang berbeda menghasilkan angka yang sama dan perbedaannya tidak pernah terlihat.

Perkakas ini mengubah pola itu dari pengamatan menjadi ALAT CARI: sapu seluruh
`Property-Set` yang nilainya literal, buang yang jelas bukan besaran bisnis, dan
laporkan sisanya untuk diadili satu per satu.

ATURAN PERKAKAS (berlaku untuk seluruh perkakas di proyek ini):
    yang ditolak ikut dilaporkan beserta jumlahnya — penolakan yang diam membuat
    yang tersisa terbaca sebagai keseluruhan.
"""
import collections
import io
import os
import re

AKAR = r'D:\XML_NURE'
FOLDER = ['Treaty In', 'Treaty In Adjustment']

RE_SET = re.compile(
    r'<PropertiesName>([^<]*)</PropertiesName>\s*<PropertiesValue>([^<]*)</PropertiesValue>')

# Literal: angka telanjang, atau teks dalam tanda kutip. Bukan rujukan properti,
# bukan ekspresi, bukan pemanggilan fungsi.
RE_LITERAL_ANGKA = re.compile(r'^-?\d+(?:\.\d+)?$')
RE_LITERAL_TEKS = re.compile(r'^"([^"]*)"$')

# Sasaran yang BUKAN properti bisnis — ditolak, dan dihitung.
TOLAK_SASARAN = re.compile(r'^(Local\.|Param\.|Temp|pg|py|px|pz)', re.I)
# Nilai yang literal tetapi tidak menarik: kosong, nol, penanda benar/salah.
TOLAK_NILAI = {'', '""', 'true', 'false', 'True', 'False', 'Y', 'N', '"Y"', '"N"',
               '"true"', '"false"'}


def literal(v):
    v = v.strip()
    if v in TOLAK_NILAI:
        return None
    if RE_LITERAL_ANGKA.match(v):
        return v
    m = RE_LITERAL_TEKS.match(v)
    if m and m.group(1).strip():
        return v
    return None


def main():
    temuan = []
    tolak_sasaran = 0
    tolak_nilai = 0
    tolak_bukan_literal = 0
    n_berkas = 0

    for f in FOLDER:
        akar = os.path.join(AKAR, f)
        if not os.path.isdir(akar):
            continue
        for d, _, berkas in os.walk(akar):
            for b in berkas:
                if not b.lower().endswith('.xml'):
                    continue
                n_berkas += 1
                p = os.path.join(d, b)
                s = io.open(p, encoding='utf-8', errors='replace').read()
                for m in RE_SET.finditer(s):
                    nama = m.group(1).replace('&lt;', '<').replace('&gt;', '>').strip()
                    nilai = m.group(2).replace('&quot;', '"').strip()
                    lit = literal(nilai)
                    if lit is None:
                        if nilai in TOLAK_NILAI:
                            tolak_nilai += 1
                        else:
                            tolak_bukan_literal += 1
                        continue
                    if TOLAK_SASARAN.match(nama):
                        tolak_sasaran += 1
                        continue
                    temuan.append((os.path.relpath(p, AKAR), nama, lit))

    # nilai yang paling sering ditulis tetap
    per_nilai = collections.Counter(t[2] for t in temuan)
    per_ruas = collections.Counter(t[1].split('.')[-1] for t in temuan)

    print('berkas disapu                 :', n_berkas)
    print('Property-Set bernilai LITERAL ke properti bisnis :', len(temuan))
    print('--- yang DITOLAK, dan sebabnya ---')
    print('  sasaran bukan properti bisnis (Local/Param/Temp/py…):', tolak_sasaran)
    print('  nilai kosong / benar-salah                          :', tolak_nilai)
    print('  nilainya bukan literal (rujukan atau ekspresi)      :', tolak_bukan_literal)
    print()
    print('--- nilai tetap yang paling sering ditulis ---')
    for v, n in per_nilai.most_common(12):
        print('   %-14s %d' % (v, n))
    print()
    print('--- ruas yang paling sering diisi tetapan ---')
    for v, n in per_ruas.most_common(15):
        print('   %-28s %d' % (v, n))
    print()
    print('--- seluruh temuan, per berkas ---')
    for berkas, nama, lit in sorted(temuan):
        print('   %-58s %-46s = %s' % (berkas[:58], nama[-46:], lit))


if __name__ == '__main__':
    main()
