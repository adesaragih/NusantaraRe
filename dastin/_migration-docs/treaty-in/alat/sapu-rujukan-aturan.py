# -*- coding: utf-8 -*-
"""
sapu-rujukan-aturan.py — SIAPA MERUJUK SIAPA, dari indeks rujukan Pega sendiri.

KENAPA ADA:  sapuan pemanggil sebelumnya bersandar pada pencarian teks atau pada
             satu nama tag. Keduanya tidak menyatakan JENIS aturan yang dirujuk,
             dan sapuan satu-nama-tag buta terhadap jenis aturan yang memakai
             nama lain (`PropertiesName` vs `pyPropertiesName`).

SUMBERNYA:   <pxRuleReferences> — indeks yang DIHITUNG PEGA SENDIRI saat aturan
             disimpan, berisi setiap aturan yang dirujuk berkas itu beserta
             jenisnya. Ia tidak bergantung pada nama tag mana pun yang kita tebak.

BATAS KEMAMPUAN PERKAKAS INI — dibaca sebelum hasilnya dipercaya:
  1. Indeks ini DIHITUNG SAAT SIMPAN. Aturan yang dirujuk lewat nama yang
     DIRAKIT SAAT JALAN (`Call @(…)`, nama dari parameter) TIDAK akan muncul.
  2. 86 dari 708 berkas TIDAK memuat indeks ini sama sekali (lihat keluaran).
     Untuk ke-86 itu perkakas ini DIAM, dan diamnya BUKAN nol rujukan.
  3. Ia hanya meliputi korpus yang ada. Jenis aturan yang tidak terekspor (L-10)
     tidak dapat menjadi perujuk di sini walaupun di sistem ia merujuk.
  => Perkakas ini dapat menyatakan sebuah aturan DIRUJUK. Ia TIDAK dapat
     menyatakan sebuah aturan TIDAK DIRUJUK; untuk itu ketiga batas di atas
     harus ditutup lebih dulu.

Keluaran:
  datar-rujukan-aturan.csv        perujuk -> yang dirujuk (jenis, kelas, nama)
  datar-aturan-tanpa-perujuk.csv  aturan di korpus yang nol kali dirujuk
"""
import os, io, csv
import xml.etree.ElementTree as ET
from collections import defaultdict

AKAR = [r'D:\XML_NURE\Treaty In', r'D:\XML_NURE\Treaty In Adjustment']
KELUAR = r'D:\XML_NURE\_migration-docs\treaty-in\alat'

# aturan yang ADA di korpus: (jenis, nama-huruf-besar) -> berkas
ada = {}
# rujukan
rujukan = []            # (berkas_perujuk, jenis_dirujuk, kelas_dirujuk, nama_dirujuk)
dirujuk_oleh = defaultdict(set)

tanpa_indeks = []
n = 0

berkas = []
for akar in AKAR:
    for dp, dn, fn in os.walk(akar):
        for f in fn:
            if f.lower().endswith('.xml'):
                berkas.append(os.path.join(dp, f))

for p in berkas:
    try:
        root = ET.parse(p).getroot()
    except Exception:
        continue
    n += 1
    jenis = (root.findtext('pxObjClass') or '').strip()
    nama  = (root.findtext('pyRuleName') or '').strip()
    if jenis and nama:
        ada.setdefault((jenis, nama.upper()), []).append(p)

for p in berkas:
    try:
        root = ET.parse(p).getroot()
    except Exception:
        continue
    rr = root.find('.//pxRuleReferences')
    rows = rr.findall('rowdata') if rr is not None else []
    if not rows:
        tanpa_indeks.append(p)
        continue
    for row in rows:
        d = {c.tag: (c.text or '').strip() for c in row}
        j = d.get('pxRuleObjClass', '')
        k = d.get('pxRuleClassName', '')
        nm = d.get('pyRuleName', '')
        if not nm:
            continue
        rujukan.append((p, j, k, nm))
        dirujuk_oleh[(j, nm.upper())].add(p)

with io.open(os.path.join(KELUAR, 'datar-rujukan-aturan.csv'), 'w', encoding='utf-8', newline='') as fh:
    w = csv.writer(fh); w.writerow(['berkas_perujuk', 'jenis_dirujuk', 'kelas_dirujuk', 'nama_dirujuk'])
    for r in rujukan:
        w.writerow([os.path.basename(r[0]), r[1], r[2], r[3]])

nol = []
with io.open(os.path.join(KELUAR, 'datar-aturan-tanpa-perujuk.csv'), 'w', encoding='utf-8', newline='') as fh:
    w = csv.writer(fh); w.writerow(['jenis', 'nama', 'berkas', 'jumlah_perujuk'])
    for (j, NM), paths in sorted(ada.items()):
        # perujuk selain dirinya sendiri
        pj = dirujuk_oleh.get((j, NM), set()) - set(paths)
        if not pj:
            nol.append((j, NM))
            w.writerow([j, NM, os.path.basename(paths[0]), 0])

print('berkas dibaca                         : %d' % n)
print('DITOLAK — tanpa <pxRuleReferences>    : %d  (indeks tidak ada; diamnya BUKAN nol rujukan)' % len(tanpa_indeks))
from collections import Counter
c = Counter(os.path.basename(os.path.dirname(p)) for p in tanpa_indeks)
for k, v in c.most_common():
    print('      %-22s %d' % (k, v))
print('aturan di korpus                      : %d' % len(ada))
print('baris rujukan                         : %d' % len(rujukan))
print('aturan korpus NOL kali dirujuk        : %d' % len(nol))
cj = Counter(j for j, _ in nol)
for k, v in cj.most_common():
    print('      %-28s %d' % (k, v))
