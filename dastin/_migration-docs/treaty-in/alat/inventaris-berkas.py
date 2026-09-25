# -*- coding: utf-8 -*-
"""
INVENTARIS BERKAS INDUK — treaty-in (bentuk 2)
BATAS KEMAMPUAN: perkakas ini TIDAK dapat menyatakan sebuah berkas SUDAH DIBACA.
Ia menyatakan DISEBUT / TIDAK PERNAH DISEBUT. Penyebutan bukan pembacaan.
BENTUK PENYEBUTAN yang dicari (tiga, bukan satu — §2.0-j):
  B1 nama berkas utuh          contoh: SPEC-INVARIAN.md
  B2 penanda ADR              contoh: ADR-0055   (untuk docs/adr/*.md)
  B3 nama objek tanpa ekstensi contoh: TREATYINDETAIL, PEGA_TREATY_IN (untuk Table/, PROCEDURE/)
"""
import io, os, re
ROOT = r"D:\XML_NURE\_migration-docs\treaty-in"
EXTRA = [r"D:\XML_NURE\_migration-docs\treaty-in-adjustment", r"D:\XML_NURE\_migration-docs\claim-non-prop"]

files=[]
for dp,dn,fn in os.walk(ROOT):
    for f in fn:
        full=os.path.join(dp,f); rel=os.path.relpath(full,ROOT).replace("\\","/")
        files.append((rel,f,os.path.getsize(full)))
files.sort()

corpus=[]; tolak=[]
for base in [ROOT]+EXTRA:
    for dp,dn,fn in os.walk(base):
        for f in fn:
            full=os.path.join(dp,f)
            if f.rsplit(".",1)[-1].lower() in ("xlsx","png","jpg","pdf","zip"): tolak.append(f); continue
            try: corpus.append((os.path.abspath(full), io.open(full,encoding="utf-8",errors="ignore").read()))
            except Exception: tolak.append(f)

def bentuk(rel, base):
    b=[base]
    m=re.match(r"docs/adr/(\d{4})-", rel)
    if m: b.append("ADR-"+m.group(1))
    if rel.startswith(("Table/","PROCEDURE/")): b.append(base.rsplit(".",1)[0])
    return b

hasil=[]
for rel,base,size in files:
    forms=bentuk(rel,base); n=0; via=set()
    self_abs=os.path.abspath(os.path.join(ROOT,rel))
    for full,t in corpus:
        if full==self_abs: continue
        for k,fm in enumerate(forms):
            c=t.count(fm)
            if c: n+=c; via.add("B%d"%(k+1))
    hasil.append((rel,size,n,",".join(sorted(via)) or "-"))

yatim=[h for h in hasil if h[2]==0]
print("JUMLAH BERKAS: %d | KORPUS: %d berkas teks | DITOLAK biner: %d (%s)" % (len(files),len(corpus),len(tolak),", ".join(sorted(set(tolak)))))
print()
print("=== TIDAK PERNAH DISEBUT, dalam SATU PUN dari tiga bentuk ===")
for rel,size,n,via in yatim: print("  %-56s %8d B" % (rel,size))
print("  (jumlah: %d)" % len(yatim))
print()
print("=== DISEBUT 1-3 KALI (calon 'belum pernah dipakai langkah mana pun') ===")
for rel,size,n,via in sorted([h for h in hasil if 0<h[2]<=3], key=lambda x:x[2]):
    print("  %-56s %3d sebutan  via %s" % (rel,n,via))
