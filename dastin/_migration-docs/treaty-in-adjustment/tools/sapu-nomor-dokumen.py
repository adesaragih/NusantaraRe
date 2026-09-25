# -*- coding: utf-8 -*-
"""
sapu-nomor-dokumen.py — GRILL-D §1: adakah properti yang menyimpan NOMOR DOKUMEN
addendum eksternal?

APA YANG DAPAT DIBUKTIKAN PERKAKAS INI
    Bahwa sebuah nama properti MUNCUL, dan di berkas mana, dan sebagai apa
    (ditulis, dibaca, atau hanya dideklarasikan).
APA YANG TIDAK
    Ia TIDAK dapat menyatakan sebuah properti TIDAK ADA di sistem lama. Dua bentuk
    berada di luar jangkauannya dan dinyatakan di sini, bukan disembunyikan:
      - Rule-Declare-Expression dan Rule-Declare-Trigger TIDAK TEREKSPOR (L-10);
      - langkah Java dan langkah SQL/REST tidak terurai sebagai penugasan properti.
    Maka nol di sini berarti NOL DI ANTARA BENTUK YANG TERURAI, dan kalimat itu
    ikut dicetak bersama hasilnya.

KALIBRASI -- wajib sebelum klaim negatif apa pun (METODE; CONTEXT 2.0-i).
    Tiga properti yang SUDAH DIKETAHUI ADA dimasukkan ke sapuan yang sama.
    Bila salah satu tidak ditemukan, penyapunya rusak dan nolnya tidak bermakna.

NAMA TAG YANG DIPAKAI, dan ia diperiksa terhadap ../treaty-in/alat/datar-nama-elemen.csv:
    PropertiesName / PropertiesValue      Rule-Obj-Activity     Property-Set
    pyPropertiesName / pyPropertiesValue  Rule-Obj-Model        data transform
    pyPropertyTarget                      Harness / Section     pemetaan pencarian
    pyTargetProperty                      Report-Definition     kolom keluaran laporan
    CopyFrom / CopyInto                   Rule-Obj-Activity     salin halaman (bentuk kelima)
"""
import io, os, re, csv, json
from collections import defaultdict, Counter

AKAR = os.path.dirname(os.path.abspath(__file__))
EKSPOR = [(r'D:\XML_NURE\Treaty In', 'induk'),
          (r'D:\XML_NURE\Treaty In Adjustment', 'adjustment')]
NAMA_ELEMEN = os.path.join(AKAR, '..', '..', 'treaty-in', 'alat', 'datar-nama-elemen.csv')

TAG = ['PropertiesName', 'pyPropertiesName', 'pyPropertyTarget', 'pyTargetProperty',
       'CopyFrom', 'CopyInto', 'PropertiesValue', 'pyPropertiesValue']

# ── DAFTAR NAMA YANG DISAPU, dicatat di dalam perkakasnya ────────────────────
CALON = [
    'DocNo', 'DocumentNo', 'DocNumber', 'NoDoc', 'NoDokumen', 'NomorDokumen',
    'AddendumNo', 'NoAddendum', 'AddendumNumber', 'AddendumRef', 'AddendumID',
    'EndorseNo', 'EndorsementNo', 'NoEndorse',
    'RefNo', 'NoRef', 'ReferenceNo', 'ContractRefNo', 'RefDoc', 'DocRef',
    'SlipNo', 'NoSlip', 'LetterNo', 'NoSurat', 'SuratNo',
    'AmendmentNo', 'NoAmendment', 'AddNo',
]
# Tiga kasus positif yang SUDAH DIKETAHUI ADA -- kalibrasi.
KALIBRASI = ['EDMState', 'EDMMaterialType', 'OLDID']

pola = re.compile(r'(?i)(' + '|'.join(re.escape(c) for c in CALON + KALIBRASI) + r')')

# ── 0. periksa nama tag terhadap daftar nama elemen ──────────────────────────
punya = set()
try:
    with io.open(NAMA_ELEMEN, encoding='utf-8') as f:
        for r in csv.DictReader(f):
            punya.add(r['nama_elemen'])
except Exception as e:
    punya = None

print('=== sapu-nomor-dokumen.py ===')
if punya is None:
    print('PERIKSA TAG : daftar nama elemen TIDAK TERBACA -- sapuan tetap jalan, tetapi')
    print('              namanya belum diperiksa. Itu batas, bukan hasil.')
else:
    tak = [t for t in TAG if t not in punya]
    print('PERIKSA TAG : %d dari %d nama tag ADA di datar-nama-elemen.csv' % (len(TAG) - len(tak), len(TAG)))
    if tak:
        print('              TIDAK ADA di daftar, dan karena itu tidak akan pernah cocok: %s' % ', '.join(tak))

# ── 1. sapu ──────────────────────────────────────────────────────────────────
hit = defaultdict(list)       # nama -> [(sisi, berkas, bentuk, konteks)]
tolak = Counter()
berkas_n = 0
for akar, sisi_default in EKSPOR:
    for dp, dn, fn in os.walk(akar):
        jenis = os.path.basename(dp)
        for f in fn:
            if not f.lower().endswith('.xml'):
                tolak['bukan .xml'] += 1
                continue
            p = os.path.join(dp, f)
            try:
                s = io.open(p, encoding='utf-8', errors='ignore').read()
            except Exception:
                tolak['berkas tidak terbaca'] += 1
                continue
            berkas_n += 1
            # buang salinan terbungkus -- ia bukan badan aturan
            for t in ('pyIncludedRuleXML', 'pyRuleVersionsList'):
                s = re.sub(r'<%s\b[^>]*>.*?</%s>' % (t, t), '', s, flags=re.S)
            if not pola.search(s):
                continue
            # SISI dibaca dari pzInsKey, BUKAN dari nama berkas
            ins = re.findall(r'<pzInsKey>([^<]*)</pzInsKey>', s)
            sisi = 'TIDAK TERBACA'
            if ins:
                k = ' '.join(ins).upper()
                ada56 = 'TREATY_IN_EDM' in k or '56' in k
                sisi = '56-KHAS' if ada56 else 'IRISAN'
            for tag in TAG:
                for m in re.finditer(r'<%s>([^<]*)</%s>' % (tag, tag), s):
                    v = m.group(1)
                    mm = pola.search(v)
                    if mm:
                        hit[mm.group(1)].append((sisi, os.path.basename(p), jenis, tag, v[:70]))
            # kemunculan di luar keempat bentuk penulis -- deklarasi, kondisi, layar
            for mm in pola.finditer(s):
                nm = mm.group(1)
                if not any(h[1] == os.path.basename(p) and h[3] in TAG for h in hit[nm]):
                    hit[nm].append((sisi, os.path.basename(p), jenis, '(bukan penulis)', ''))
                    break

print('BERKAS disapu : %d' % berkas_n)
print('DITOLAK -- satu baris per sebab:')
for k, v in tolak.most_common():
    print('   %-34s %d' % (k, v))
if not tolak:
    print('   (tidak ada)')
print()
print('KALIBRASI -- tiga kasus positif yang sudah diketahui ADA:')
lulus = True
for k in KALIBRASI:
    n = len(hit.get(k, []))
    print('   %-20s %s (%d kemunculan)' % (k, 'DITEMUKAN' if n else 'TIDAK DITEMUKAN  <-- PENYAPU RUSAK', n))
    if not n:
        lulus = False
print('   -> %s' % ('LULUS. Nol pada calon di bawah BERMAKNA.' if lulus else 'GAGAL. Jangan pakai hasilnya.'))
print()
print('CALON NOMOR DOKUMEN -- %d nama disapu:' % len(CALON))
ada = [(c, hit[c]) for c in CALON if hit.get(c)]
if not ada:
    print('   NOL. Tidak satu pun dari %d nama muncul di kedua ekspor.' % len(CALON))
else:
    for c, rows in ada:
        per_sisi = Counter(r[0] for r in rows)
        penulis = [r for r in rows if r[3] in TAG]
        print('   %-18s %2d kemunculan | sisi: %s | sebagai PENULIS: %d'
              % (c, len(rows), dict(per_sisi), len(penulis)))
        for r in rows[:4]:
            print('        %-10s %-34s %-22s %s' % (r[0], r[1][:34], r[3], r[4][:44]))
print()
print('BATAS YANG DINYATAKAN BERSAMA HASIL INI:')
print('   Rule-Declare-Expression dan Rule-Declare-Trigger TIDAK TEREKSPOR (L-10).')
print('   Langkah Java dan langkah SQL/REST tidak terurai sebagai penugasan properti.')
print('   Maka nol di atas berarti NOL DI ANTARA BENTUK YANG TERURAI.')

json.dump({c: hit.get(c, []) for c in CALON + KALIBRASI},
          io.open(os.path.join(AKAR, 'sapu-nomor-dokumen.json'), 'w', encoding='utf-8'),
          ensure_ascii=False, indent=1)
