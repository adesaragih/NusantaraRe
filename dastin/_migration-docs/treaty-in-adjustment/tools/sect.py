# -*- coding: utf-8 -*-
"""Ringkasan Section/Harness Pega: aktivitas yang dipanggil, seksi anak, properti,
   dan kondisi tampil. pyVisible=ALWAYS MENGABAIKAN pyCondition di sebelahnya."""
import os, sys, xml.etree.ElementTree as ET, collections
BASE=r"D:\XML_NURE\Treaty In Adjustment"
def find(n):
    if os.path.exists(n): return n
    for dp,dn,fn in os.walk(BASE):
        for f in fn:
            if f.lower() in (n.lower(), n.lower()+'.xml'): return os.path.join(dp,f)
for name in sys.argv[1:]:
    p=find(name)
    if not p: print("TIDAK KETEMU",name); continue
    r=ET.parse(p).getroot()
    print("### %s | %s | %s %s | %s"%((r.findtext('pyRuleName') or '').strip(),(r.findtext('pyClassName') or '').strip(),
        (r.findtext('pyRuleSet') or '').strip(),(r.findtext('pyRuleSetVersion') or '').strip(), os.path.basename(os.path.dirname(p))))
    acts=collections.Counter(); secs=collections.Counter(); props=collections.Counter()
    vis=[]; dt=collections.Counter()
    for e in r.iter():
        t=(e.text or '').strip()
        if not t: continue
        if e.tag in ('pyActivityName','pyActionActivity'): acts[t]+=1
        elif e.tag in ('pySectionName','pyEmbedSectionName','pySectionNameValue'): secs[t]+=1
        elif e.tag in ('pyDataTransformName','DataTransform'): dt[t]+=1
        elif e.tag=='pyPropertyName': props[t]+=1
        elif e.tag in ('pyCondition','pyVisibleWhen','pyWhen'): vis.append((e.tag,t))
        elif e.tag=='pyVisible' and t!='ALWAYS': vis.append((e.tag,t))
    def show(lbl,c):
        if c: print("  %s: %s"%(lbl,", ".join("%s"%k for k,_ in c.most_common())))
    show("aktivitas",acts); show("data transform",dt); show("seksi anak",secs)
    if props: print("  properti (%d): %s"%(len(props),", ".join(sorted(props))[:1200]))
    if vis:
        print("  kondisi:")
        seen=set()
        for k,v in vis:
            if (k,v) in seen: continue
            seen.add((k,v)); print("    %-14s %s"%(k,v[:130]))
    print()
