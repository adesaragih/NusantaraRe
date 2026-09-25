import os, xml.etree.ElementTree as ET, collections
BASE=r"D:\XML_NURE"
VOL={'pyRuleFormStatusTime','pyShowJavaWindowName','pyJavaGenerateTime','pyJavaClassName','pyCellId','pySectionId'}
def items(p):
    r=ET.parse(p).getroot(); out=[]
    def walk(n,path):
        for c in n:
            if c.tag.startswith('px') or c.tag.startswith('pz') or c.tag in VOL: continue
            out.append((path+"/"+c.tag,(c.text or '').strip())); walk(c,path+"/"+c.tag)
    walk(r,r.tag); return out
def rels(root):
    s=set()
    for dp,dn,fn in os.walk(root):
        for f in fn:
            if f.lower().endswith(".xml"): s.add(os.path.relpath(os.path.join(dp,f),root))
    return s
A=os.path.join(BASE,"Treaty In Adjustment"); T=os.path.join(BASE,"Treaty In")
ra,rt=rels(A),rels(T)
print("Adjustment:",len(ra),"TreatyIn:",len(rt),"irisan:",len(ra&rt),"hanya-Adj:",len(ra-rt),"hanya-TI:",len(rt-ra))
diff=[]
for rel in sorted(ra&rt):
    ca=collections.Counter(items(os.path.join(A,rel))); ct=collections.Counter(items(os.path.join(T,rel)))
    d=(ca-ct)+(ct-ca)
    if d: diff.append((rel,sum(d.values()),sorted({p.split("/")[-1] for p,v in d})))
print("\nirisan yang ISINYA BEDA (setelah buang artefak ekspor/generated):",len(diff))
for rel,n,tags in diff: print("  %-55s %4d  %s"%(rel,n,",".join(tags)[:110]))
