# -*- coding: utf-8 -*-
"""
urai-sepuluh.py — membaca §10 SPEC-MODEL-DATA.md sebagai SATU SUMBER definisi.
BATAS: ia hanya melihat baris tabel beratribut lima kolom di dalam §10.x.
       Entitas yang atributnya ditulis dalam bentuk lain TIDAK terlihat, dan
       ketidakterlihatan itu DILAPORKAN, bukan didiamkan.
"""
import io, re, json, sys
P = r"D:\XML_NURE\_migration-docs\treaty-in\SPEC-MODEL-DATA.md"
baris = io.open(P, encoding="utf-8").read().split("\n")

ent=None; hasil={}; urut=[]
tolak_tanpa_tabel=[]; tolak_baris=0
judul=re.compile(r"^### (10\.\d+[a-z]?)\s+`([A-Z_]{2,})`")
for i,l in enumerate(baris):
    m=judul.match(l)
    if m:
        ent=m.group(2); 
        if ent not in hasil: hasil[ent]=[]; urut.append((m.group(1),ent))
        continue
    if l.startswith("### ") and not judul.match(l): ent=None; continue
    if ent and l.startswith("| ") and l.count("|")>=5:
        sel=[c.strip() for c in l.strip().strip("|").split("|")]
        if len(sel)<4: tolak_baris+=1; continue
        nama=sel[0].strip("`* ")
        if nama in ("Nama","---") or set(nama)<=set("-: "): continue
        if not re.match(r"^[A-Z][A-Z0-9_]*$", nama): tolak_baris+=1; continue
        asal=sel[1].strip("`* "); tipe=sel[2].strip("* "); kosong=sel[3].strip("* ")
        cat = sel[4] if len(sel)>4 else ""
        hasil[ent].append(dict(kolom=nama, asal=asal, tipe=tipe, kosong=kosong, catatan=cat))
for e,v in hasil.items():
    if not v: tolak_tanpa_tabel.append(e)
print("ENTITAS dengan tabel atribut terbaca : %d" % len([e for e in hasil if hasil[e]]))
print("ENTITAS tanpa tabel atribut (DITOLAK): %d -> %s" % (len(tolak_tanpa_tabel), ", ".join(tolak_tanpa_tabel) or "-"))
print("BARIS tabel ditolak (bukan atribut)  : %d" % tolak_baris)
print()
for sec,e in urut:
    if hasil[e]: print("  %-7s %-24s %2d atribut" % (sec,e,len(hasil[e])))
io.open(r"D:\XML_NURE\_migration-docs\treaty-in\alat\definisi-sepuluh.json","w",encoding="utf-8").write(
    json.dumps({"urut":urut,"entitas":hasil}, ensure_ascii=False, indent=1))
