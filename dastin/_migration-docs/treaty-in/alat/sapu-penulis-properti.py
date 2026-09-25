# -*- coding: utf-8 -*-
"""
sapu-penulis-properti.py — SIAPA MENULIS properti kontrak, atas KELIMA bentuk.

KENAPA ADA:  sapuan penulis sebelumnya hanya mengenali <PropertiesName> (Activity).
             Data transform memakai <pyPropertiesName> — beda awalan `py` — dan
             73 berkas tidak pernah ikut tersapu selama sebelas sesi.

NAMA TAG YANG DIPAKAI SAPUAN INI, dan ia wajib dicocokkan ke
`datar-nama-elemen.csv` sebelum hasilnya dipercaya (CONTEXT.md 2.0-i):
    PropertiesName    / PropertiesValue      Rule-Obj-Activity   Property-Set
    pyPropertiesName  / pyPropertiesValue    Rule-Obj-Model      (data transform)
    pyPropertyTarget  / pySetValueOnSelect   Harness/Section     pemetaan pencarian
    pyTargetProperty                         Report-Definition   BUKAN penulis —
                                             ia kolom keluaran laporan, relatif
                                             terhadap kelas laporannya sendiri.
    CopyFrom / CopyInto  (pyStepsCallParams) Rule-Obj-Activity   SALIN HALAMAN.
                                             Bentuk KELIMA, ditemukan 24 Sep 2026
                                             lewat "nol yang mustahil": OLDDATA
                                             punya empat layar dan nol penulis.
                                             Sasarannya duduk di PARAMETER LANGKAH,
                                             bukan di pasangan nama/nilai -- jadi
                                             ia tak terlihat oleh bentuk mana pun
                                             di atasnya.

BATAS: bentuk yang tidak terekspor (Declare Expression, Declare Trigger — L-10)
       tidak terlihat sapuan ini sama sekali.
"""
import os, io, csv, re
import xml.etree.ElementTree as ET
from collections import defaultdict

AKAR = [r'D:\XML_NURE\Treaty In', r'D:\XML_NURE\Treaty In Adjustment']
KELUAR = r'D:\XML_NURE\_migration-docs\treaty-in\alat'

hasil = []          # (bentuk, berkas, jenis_aturan, sasaran, nilai)
tolak_ekspresi = 0
OP = re.compile(r'==|!=|>=|<=|\|\||&&')

def rekam(bentuk, p, jenis, tgt, val):
    global tolak_ekspresi
    tgt = (tgt or '').strip()
    if not tgt:
        return
    if OP.search(tgt):
        tolak_ekspresi += 1
        return
    hasil.append((bentuk, os.path.basename(p), jenis, tgt, (val or '').strip()[:120]))

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
    jenis = (root.findtext('pxObjClass') or '').strip()
    for row in root.iter('rowdata'):
        d = {c.tag: (c.text or '') for c in row}
        if 'PropertiesName' in d:
            rekam('Activity/Property-Set', p, jenis, d['PropertiesName'], d.get('PropertiesValue'))
        if 'pyPropertiesName' in d:
            rekam('DataTransform', p, jenis, d['pyPropertiesName'], d.get('pyPropertiesValue'))
        if 'pyPropertyTarget' in d:
            rekam('Harness/Section-pemetaan', p, jenis, d['pyPropertyTarget'], d.get('pyDisplayProperty'))
    for cp in root.iter('pyStepsCallParams'):
        into = (cp.findtext('CopyInto') or '').strip()
        frm  = (cp.findtext('CopyFrom') or '').strip()
        if into:
            rekam('Salin-halaman', p, jenis, into, frm)

with io.open(os.path.join(KELUAR, 'datar-penulis-properti.csv'), 'w', encoding='utf-8', newline='') as fh:
    w = csv.writer(fh); w.writerow(['bentuk', 'berkas', 'jenis_aturan', 'sasaran', 'nilai'])
    for r in sorted(set(hasil)):
        w.writerow(r)

from collections import Counter
c = Counter(r[0] for r in hasil)
print('baris penugasan ditemukan:')
for k, v in c.most_common():
    print('   %-26s %d' % (k, v))
print('DITOLAK — sasaran memuat operator (ia ekspresi when, bukan penugasan): %d' % tolak_ekspresi)
print('sasaran berbeda           : %d' % len(set(r[3] for r in hasil)))
ti = sorted(set(r[3] for r in hasil if r[3].startswith('TreatyIn.')))
print('sasaran berawalan TreatyIn.: %d' % len(ti))
