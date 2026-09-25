# -*- coding: utf-8 -*-
"""
sapu-nama-elemen.py — inventaris SELURUH nama elemen XML di korpus, beserta
jenis aturan tempat ia muncul.

KENAPA ADA:  CONTEXT.md 2.0-i. Sebuah sapuan yang bersandar pada SATU nama tag
             buta terhadap jenis aturan yang memakai nama lain untuk gagasan
             yang sama (`PropertiesName` vs `pyPropertiesName`).

BATAS KEMAMPUAN PERKAKAS INI:
  - Ia TIDAK dapat menyatakan sebuah sapuan benar. Ia hanya dapat menunjukkan
    nama elemen yang ADA dan TIDAK dipakai sapuan itu.
  - Ia membaca nama elemen, bukan artinya. Penyaringan menghasilkan CALON;
    setiap calon harus dibuka berkasnya.
  - Ia buta terhadap apa pun yang tidak ada di ekspor (L-10).

Keluaran:
  datar-nama-elemen.csv        nama, jumlah, jumlah berkas, jenis aturan
  datar-calon-rujukan.csv      calon berbentuk RUJUKAN (memanggil/merujuk aturan lain)
"""
import os, io, sys, csv, re
import xml.etree.ElementTree as ET
from collections import defaultdict

AKAR = [r'D:\XML_NURE\Treaty In', r'D:\XML_NURE\Treaty In Adjustment']
KELUAR = r'D:\XML_NURE\_migration-docs\treaty-in\alat'

jumlah   = defaultdict(int)
berkas_n = defaultdict(set)
jenis_n  = defaultdict(set)

n_berkas = 0
tolak_parse = []

for akar in AKAR:
    for dp, dn, fn in os.walk(akar):
        for f in fn:
            if not f.lower().endswith('.xml'):
                continue
            p = os.path.join(dp, f)
            try:
                root = ET.parse(p).getroot()
            except Exception as e:
                tolak_parse.append((p, str(e)))
                continue
            n_berkas += 1
            jenis = root.findtext('pxObjClass') or root.tag
            for el in root.iter():
                jumlah[el.tag]   += 1
                berkas_n[el.tag].add(p)
                jenis_n[el.tag].add(jenis)

# ---- keluaran 1: seluruh nama elemen
with io.open(os.path.join(KELUAR, 'datar-nama-elemen.csv'), 'w', encoding='utf-8', newline='') as fh:
    w = csv.writer(fh)
    w.writerow(['nama_elemen', 'jumlah', 'jumlah_berkas', 'jenis_aturan'])
    for k in sorted(jumlah, key=lambda x: (-jumlah[x], x)):
        w.writerow([k, jumlah[k], len(berkas_n[k]), ' | '.join(sorted(jenis_n[k]))])

# ---- keluaran 2: calon RUJUKAN
# bentuk rujukan = nama yang menyebut JENIS ATURAN LAIN, atau menyebut nama/kelas
# sebuah aturan yang dipanggil/dimuat.
KATA = ['activity', 'rule', 'when', 'transform', 'section', 'harness', 'flow',
        'model', 'report', 'list', 'strategy', 'decision', 'map', 'declare',
        'validate', 'service', 'connect', 'html', 'stream', 'utility', 'function',
        'library', 'class', 'applies', 'called', 'call', 'invoke', 'ref', 'include']
EKOR = ['name', 'class', 'ref', 'key', 'id', 'reference', 'references', 'label']

def calon_rujukan(nama):
    low = nama.lower()
    for kata in KATA:
        if kata in low:
            # harus BERPASANGAN: menyebut nama/kelas/rujukan sesuatu
            for ek in EKOR:
                if ek in low:
                    return kata
            # atau namanya sendiri sudah rujukan langsung
            if low.endswith(kata) and kata in ('ref', 'call', 'invoke'):
                return kata
    return None

with io.open(os.path.join(KELUAR, 'datar-calon-rujukan.csv'), 'w', encoding='utf-8', newline='') as fh:
    w = csv.writer(fh)
    w.writerow(['nama_elemen', 'kata_pemicu', 'jumlah', 'jumlah_berkas', 'jenis_aturan'])
    n = 0
    for k in sorted(jumlah, key=lambda x: (-jumlah[x], x)):
        kata = calon_rujukan(k)
        if kata:
            n += 1
            w.writerow([k, kata, jumlah[k], len(berkas_n[k]), ' | '.join(sorted(jenis_n[k]))])

print('berkas dibaca          : %d' % n_berkas)
print('GAGAL DIURAI (ditolak) : %d' % len(tolak_parse))
for p, e in tolak_parse[:10]:
    print('   - %s : %s' % (p, e))
print('nama elemen berbeda    : %d' % len(jumlah))
print('calon berbentuk rujukan: %d' % n)
