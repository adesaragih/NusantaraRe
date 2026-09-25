# -*- coding: utf-8 -*-
"""
sapu-tingkat-paket-uang.py — langkah 1 sesi to-spec.

APA YANG DAPAT IA BUKTIKAN
  Ia dapat menunjukkan bahwa sebuah paket uang DITULIS oleh sebuah ekspresi, dan
  memperlihatkan ekspresinya apa adanya.
APA YANG TIDAK
  Ia TIDAK dapat menyatakan sebuah paket uang bertingkat TREATY_100_PERSEN.
  Ketiadaan perkalian Share/Pct BUKAN bukti tingkat 100% -- ia bisa berarti
  penulisnya tidak terekspor (L-10: Declare Expression / Declare Trigger),
  atau nilainya diketik langsung oleh pengguna dan tingkatnya adalah KESEPAKATAN,
  bukan perhitungan. Yang begitu dilaporkan sebagai BELUM, bukan sebagai 100%.
PENANDA MATI diperiksa lebih dulu: pyStepsBlockName == "//", pyStepsPreCondition == "false".
  Langkah mati dicetak TERPISAH dan tidak pernah menjadi dasar putusan.
PENOLAKAN dilaporkan: sasaran yang memuat operator (ekspresi when) dihitung.
"""
import os, io, re, sys
import xml.etree.ElementTree as ET
from collections import defaultdict

AKAR = [r'D:\XML_NURE\Treaty In', r'D:\XML_NURE\Treaty In Adjustment']
OP = re.compile(r'==|!=|>=|<=|\|\||&&')
KALI = re.compile(r'(Share|Pct|Percent|Proportion|RNMShare|QSPct)', re.I)

SASARAN = {
 'NILAI_EGNPI':              ['EGNPI'],
 'BATAS_MAKSIMUM_KELOMPOK':  ['MaxCoGroup'],
 'BATAS_MAKSIMUM_NON_KELOMPOK':['MaxCoNonGroup'],
 'BATAS_PILIHAN':            ['OptionLimit'],
 'LIMIT_AGREGAT':            ['AgregateLimit','AggregateLimit'],
 'MDP':                      ['MDPList'],
 'MDP_MINIMUM':              ['MDPMinList'],
 'NILAI_RETENSI':            ['RetentionList'],
 'NILAI_TERMIN':             ['Installment'],
 'NILAI_BATAS':              ['Earthquake','FloodJab','FloodNation','RSMDLimit'],
 'PREMI_BRUTO':              ['GrossPremiumList'],
 'PREMI_BRUTO_MINIMUM':      ['GrossPremiumMinList'],
 'CADANGAN_PREMI':           ['ReserveList'],
 'KAPASITAS_SURPLUS':        ['IOOLimitList'],
}
# KALIBRASI -- satu kasus positif yang SUDAH DIKETAHUI jawabannya (CONTEXT 2.0-i / TA-04).
# NILAI_PENYEBARAN bertingkat BAGIAN_NURE menurut STRUKTUR-DATA.md 1.4.
# Bila sapuan ini tidak menemukan perkalian Share/Pct padanya, sapuannya yang salah.
KALIBRASI = {'__KALIBRASI_NILAI_PENYEBARAN': ['BreakDownSprdList','SpreadingList']}

tolak_ekspresi=0; tolak_parse=0
hit=defaultdict(list)   # paket -> (hidup, berkas, bentuk, sasaran, nilai, sebab_mati)

def cocok(tgt, kunci):
    return any(re.search(r'(^|\.)'+re.escape(k)+r'(\b|\.|\()', tgt) for k in kunci)

def rekam(paket, hidup, sebab, p, bentuk, tgt, val):
    hit[paket].append((hidup, os.path.basename(p), bentuk, tgt, (val or '').strip()[:160], sebab))

berkas=[]
for a in AKAR:
    for dp,dn,fn in os.walk(a):
        for f in fn:
            if f.lower().endswith('.xml'): berkas.append(os.path.join(dp,f))

SEMUA = dict(SASARAN); SEMUA.update(KALIBRASI)

for p in berkas:
    try: root=ET.parse(p).getroot()
    except Exception: tolak_parse+=1; continue
    for par in root.iter():
        for c in list(par):
            if c.tag in ('pyIncludedRuleXML','pyRuleVersionsList'): par.remove(c)
    # peta: rowdata langkah -> status hidup
    mati_ids=set()
    for steps in root.iter('pySteps'):
        for row in steps.findall('rowdata'):
            b=row.findtext('pyStepsBlockName') or ''
            pc=row.findtext('pyStepsPreCondition') or ''
            sebab=[]
            if b.strip()=='//': sebab.append('blok //')
            if pc.strip()=='false': sebab.append('preCondition=false')
            if sebab: mati_ids.add((id(row), ' + '.join(sebab)))
    mati_map={k:v for k,v in mati_ids}
    def status(node):
        # telusuri ke atas: apakah node berada di dalam rowdata langkah yang mati
        return None
    for steps in root.iter('pySteps'):
        for row in steps.findall('rowdata'):
            sebab = mati_map.get(id(row))
            hidup = sebab is None
            for rd in row.iter('rowdata'):
                d={c.tag:(c.text or '') for c in rd}
                for nm,vl,bentuk in (('PropertiesName','PropertiesValue','Property-Set'),
                                     ('pyPropertiesName','pyPropertiesValue','DataTransform')):
                    if nm in d:
                        tgt=d[nm].strip()
                        if not tgt: continue
                        if OP.search(tgt): tolak_ekspresi+=1; continue
                        for paket,kunci in SEMUA.items():
                            if cocok(tgt,kunci): rekam(paket,hidup,sebab,p,bentuk,tgt,d.get(vl))
    # data transform di luar pySteps
    for rd in root.iter('rowdata'):
        d={c.tag:(c.text or '') for c in rd}
        if 'pyPropertiesName' in d:
            tgt=d['pyPropertiesName'].strip()
            if tgt and not OP.search(tgt):
                dis=(rd.findtext('pyDisabled') or '').strip()=='true'
                for paket,kunci in SEMUA.items():
                    if cocok(tgt,kunci):
                        rekam(paket, not dis, 'pyDisabled=true' if dis else None, p,'DataTransform(DT)',tgt,d.get('pyPropertiesValue'))

print("DITOLAK  sasaran memuat operator (ekspresi when, bukan penugasan) : %d" % tolak_ekspresi)
print("DITOLAK  berkas tidak terurai                                     : %d" % tolak_parse)
print()
for paket in list(SASARAN)+list(KALIBRASI):
    rows=hit.get(paket,[])
    # buang duplikat
    uniq=[]; seen=set()
    for r in rows:
        k=(r[1],r[3],r[4])
        if k in seen: continue
        seen.add(k); uniq.append(r)
    hidup=[r for r in uniq if r[0]]; mati=[r for r in uniq if not r[0]]
    print("="*100)
    print("%s   penulis HIDUP: %d   penulis MATI: %d" % (paket, len(hidup), len(mati)))
    for r in hidup[:14]:
        tanda = "  <== PERKALIAN Share/Pct" if KALI.search(r[4]) else ""
        print("   HIDUP %-40s %-16s %-44s = %s%s" % (r[1][:40], r[2][:16], r[3][:44], r[4][:80], tanda))
    if len(hidup)>14: print("   ... (%d lagi)" % (len(hidup)-14))
    for r in mati[:5]:
        print("   MATI  %-40s [%s] %-40s = %s" % (r[1][:40], r[5], r[3][:40], r[4][:60]))
