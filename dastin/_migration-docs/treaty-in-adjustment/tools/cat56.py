# -*- coding: utf-8 -*-
import os, xml.etree.ElementTree as ET
BASE=r"D:\XML_NURE"; A=os.path.join(BASE,"Treaty In Adjustment"); T=os.path.join(BASE,"Treaty In")
def rels(root):
    s=set()
    for dp,dn,fn in os.walk(root):
        for f in fn:
            if f.lower().endswith(".xml"): s.add(os.path.relpath(os.path.join(dp,f),root))
    return s
def txt(n,t):
    e=n.find(t); return (e.text or '').strip() if e is not None and e.text else ''
only=sorted(rels(A)-rels(T))
print("%-12s %-42s %-46s %-12s %s"%("TIPE","RULE","CLASS","RULESET","MEMO"))
for rel in only:
    p=os.path.join(A,rel); r=ET.parse(p).getroot()
    memo=txt(r,'pyNote') or txt(r,'pyMemo') or txt(r,'pyDescription') or txt(r,'pyRuleDescription')
    print("%-12s %-42s %-46s %-12s %s"%(rel.split(os.sep)[0], txt(r,'pyRuleName')[:42], txt(r,'pyClassName')[:46],
          txt(r,'pyRuleSet')+" "+txt(r,'pyRuleSetVersion'), memo[:70].replace("\n"," ")))
