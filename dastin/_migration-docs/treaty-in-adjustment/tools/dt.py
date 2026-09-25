# -*- coding: utf-8 -*-
import os, sys, xml.etree.ElementTree as ET
BASE=r"D:\XML_NURE"
def find(name, root=os.path.join(BASE,"Treaty In Adjustment")):
    for dp,dn,fn in os.walk(root):
        for f in fn:
            if f.lower() in (name.lower(), name.lower()+".xml"): return os.path.join(dp,f)
def txt(n,t):
    e=n.find(t); return (e.text or '').strip() if e is not None and e.text else ''
for name in sys.argv[1:]:
    p=find(name) or name
    r=ET.parse(p).getroot()
    print("### %s | %s | %s %s"%(txt(r,'pyRuleName'),txt(r,'pyClassName'),txt(r,'pyRuleSet'),txt(r,'pyRuleSetVersion')))
    m=txt(r,'pyNote') or txt(r,'pyMemo')
    if m: print("memo:",m[:200])
    def walk(node,pref):
        i=0
        for row in node.findall('rowdata'):
            i+=1
            act=txt(row,'pyActionName'); nm=txt(row,'pyPropertiesName'); vl=txt(row,'pyPropertiesValue')
            dis=txt(row,'pyDisabled')
            print("%s%-6s %-3s %-16s %-48s %s"%(pref,i,"MATI" if dis=="true" else "",act,nm[:48],vl[:110]))
            sub=row.find('pyProperties')
            if sub is not None: walk(sub,pref+str(i)+".")
    a=r.find('pyProperties')
    if a is None: print("(kosong)")
    else: walk(a,"")
    print()
