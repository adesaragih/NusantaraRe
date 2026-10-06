// ⛔ BERKAS BANGKITAN — JANGAN DISUNTING TANGAN.
// Dibangkitkan `alat/ekstrak_kerangka.py` dari Section ekspor `Treaty In Adjustment`,
// sesudah `<pyIncludedRuleXML>` dibuang dengan menghitung kedalaman sarang.
// Offset `at` = posisi bita di berkas Section SESUDAH pembuangan itu.

import type { ButirKerangka, GridKerangka, Kerangka, Terbuang } from './jenis'

export const KERANGKA_TAB: Readonly<Record<string, Kerangka>> = {
 "TreatyInTabsNonProportional#Maximum Retention": {
  "at": 13608,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 23476,
    "judul": "Maximum Retention",
    "syarat": [],
    "anak": [
     {
      "t": "grid",
      "at": 40873,
      "prop": "TreatyIn.Retention",
      "dari": "sisi",
      "larik": "Retention",
      "syarat": [],
      "kolom": [
       "Treaty Group",
       "Currency",
       "Amount"
      ],
      "kunci": [
       "TreatyGroup",
       "Currency",
       "Amount"
      ],
      "lebar": [
       245,
       127,
       379
      ],
      "desimal": [
       null,
       null,
       null
      ],
      "format": [
       "pxAutoComplete",
       "",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null,
       null
      ],
      "atSel": [
       66870,
       74649,
       79118
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInRetention!pyGridRowDetails"
     }
    ]
   },
   {
    "t": "grid",
    "at": 147233,
    "prop": "TreatyIn.TotalRetentionAmountNP",
    "dari": "sisi",
    "larik": "TotalRetentionAmountNP",
    "syarat": [],
    "kolom": [
     "Total Retention Amount",
     "Value"
    ],
    "kunci": [
     "Currency",
     "Value"
    ],
    "lebar": [
     193,
     349
    ],
    "desimal": [
     null,
     2
    ],
    "format": [
     "pxNumber",
     "pxNumber"
    ],
    "syaratSel": [
     null,
     null
    ],
    "atSel": [
     159212,
     164168
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
   },
   {
    "t": "tombol",
    "at": 198918,
    "label": "Update Total",
    "syarat": [
     "TreatyIn.IsEditData !='1'"
    ]
   }
  ]
 },
 "TreatyInTabsNonProportional#Event Limits": {
  "at": 220908,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 231300,
    "judul": "Event Limits",
    "syarat": [],
    "anak": [
     {
      "t": "medan",
      "at": 260828,
      "label": "RSMD Limit",
      "dari": "sisi",
      "kunci": "CurrencyRSMD",
      "format": "pxAutoComplete",
      "desimal": null,
      "syarat": []
     },
     {
      "t": "medan",
      "at": 286560,
      "label": "",
      "dari": "sisi",
      "kunci": "RSMDLimit",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": []
     },
     {
      "t": "medan",
      "at": 317490,
      "label": "Earthquake Limit",
      "dari": "sisi",
      "kunci": "CurrencyEarthquake",
      "format": "pxAutoComplete",
      "desimal": null,
      "syarat": []
     },
     {
      "t": "medan",
      "at": 343882,
      "label": "",
      "dari": "sisi",
      "kunci": "Earthquake",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": []
     },
     {
      "t": "medan",
      "at": 374771,
      "label": "Flood Limit (Jabodetabek)",
      "dari": "sisi",
      "kunci": "CurrencyFloodJab",
      "format": "pxAutoComplete",
      "desimal": null,
      "syarat": []
     },
     {
      "t": "medan",
      "at": 400573,
      "label": "",
      "dari": "sisi",
      "kunci": "FloodJab",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": []
     },
     {
      "t": "medan",
      "at": 431464,
      "label": "Flood Limit (Nationwide)",
      "dari": "sisi",
      "kunci": "CurrencyFloodNat",
      "format": "pxAutoComplete",
      "desimal": null,
      "syarat": []
     },
     {
      "t": "medan",
      "at": 457287,
      "label": "",
      "dari": "sisi",
      "kunci": "FloodNation",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": []
     }
    ]
   }
  ]
 },
 "TreatyInTabsNonProportional#EGNPI": {
  "at": 488216,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 506990,
    "judul": "Estimate Gross Net Premium Income",
    "syarat": [],
    "anak": [
     {
      "t": "grid",
      "at": 524268,
      "prop": "TreatyIn.EGNPI",
      "dari": "sisi",
      "larik": "EGNPI",
      "syarat": [],
      "kolom": [
       "Treaty Group",
       "As Date",
       "Proportion %",
       "Currency",
       "Amount",
       "Amount in IDR"
      ],
      "kunci": [
       "TreatyGroup",
       "AsDate",
       "Proportion",
       "Currency",
       "Amount",
       "AmountIDR"
      ],
      "lebar": [
       287,
       197,
       142,
       137,
       222,
       223
      ],
      "desimal": [
       null,
       null,
       2,
       null,
       2,
       null
      ],
      "format": [
       "pxAutoComplete",
       "",
       "pxNumber",
       "",
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null,
       null,
       null,
       null,
       null
      ],
      "atSel": [
       562421,
       567864,
       572297,
       577381,
       581815,
       586894
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInEGNPI!pyGridRowDetails"
     }
    ]
   },
   {
    "t": "grid",
    "at": 655147,
    "prop": "TreatyIn.TotalEgnpiAmountNP",
    "dari": "sisi",
    "larik": "TotalEgnpiAmountNP",
    "syarat": [],
    "kolom": [
     "Total EGNPI Amount",
     "Value"
    ],
    "kunci": [
     "Currency",
     "Value"
    ],
    "lebar": [
     194,
     349
    ],
    "desimal": [
     null,
     2
    ],
    "format": [
     "pxNumber",
     "pxNumber"
    ],
    "syaratSel": [
     null,
     null
    ],
    "atSel": [
     667054,
     671998
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
   },
   {
    "t": "teks",
    "at": 733571,
    "teks": "Total Amount in IDR",
    "syarat": []
   },
   {
    "t": "teks",
    "at": 753273,
    "teks": "IDR",
    "syarat": []
   },
   {
    "t": "medan",
    "at": 757649,
    "label": "",
    "dari": "sisi",
    "kunci": "TotalEgnpiAmount",
    "format": "pxNumber",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "teks",
    "at": 790523,
    "teks": "Total Proportion %",
    "syarat": []
   },
   {
    "t": "medan",
    "at": 801254,
    "label": "",
    "dari": "sisi",
    "kunci": "TotalEgnpiProportion",
    "format": "pxNumber",
    "desimal": 2,
    "syarat": []
   },
   {
    "t": "tombol",
    "at": 840132,
    "label": "Update Total",
    "syarat": [
     "TreatyIn.ViewState !='1'"
    ]
   },
   {
    "t": "tombol",
    "at": 849212,
    "label": "Update EGNPI Value",
    "syarat": [
     "TreatyIn.ViewState !='1' && TreatyMasterInEDM"
    ]
   }
  ]
 },
 "TreatyInTabsNonProportional#Limits": {
  "at": 883398,
  "syarat": [],
  "isi": [
   {
    "t": "grid",
    "at": 919339,
    "prop": "TreatyIn.Limits",
    "dari": "sisi",
    "larik": "Limits",
    "syarat": [],
    "kolom": [
     "Layers",
     "",
     "",
     "",
     "",
     "100% Limits ( IDR )",
     "Deductible ( IDR )",
     "100% Limits ( USD )",
     "Deductible ( USD )"
    ],
    "kunci": [
     "LayerType",
     "Layer",
     "pyTemplateRichTextEditor",
     "LayerPartType",
     "LayerPart",
     "Limit",
     "Deductible",
     "Limit2",
     "Deductible2"
    ],
    "lebar": [
     102,
     102,
     102,
     102,
     105,
     185,
     187,
     185,
     187
    ],
    "desimal": [
     null,
     null,
     null,
     null,
     null,
     2,
     2,
     2,
     2
    ],
    "format": [
     "pxTextInput",
     "",
     "pxTextInput",
     "pxTextInput",
     "",
     "pxNumber",
     "pxNumber",
     "pxNumber",
     "pxNumber"
    ],
    "syaratSel": [
     null,
     null,
     null,
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "atSel": [
     971285,
     977204,
     981637,
     987644,
     993471,
     997908,
     1002988,
     1008073,
     1013782
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails"
   },
   {
    "t": "blok",
    "at": 1058846,
    "judul": "Summary of Limit",
    "syarat": [],
    "anak": [
     {
      "t": "grid",
      "at": 1085096,
      "prop": "TreatyIn.LimitSummaryList",
      "dari": "sisi",
      "larik": "LimitSummaryList",
      "syarat": [],
      "kolom": [
       "Note",
       "100% Limit (IDR)",
       "100% Limit (USD)",
       "MDP (IDR)",
       "MDP (USD)",
       "Agregate Limit (IDR)",
       "Agregate Limit (USD)",
       "Deductible (IDR)",
       "Deductible (USD)"
      ],
      "kunci": [
       "Note",
       "Limit",
       "Limit2",
       "MDP",
       "MDP2",
       "AggregateLimit",
       "AggregateLimit2",
       "Deductible",
       "Deductible2"
      ],
      "lebar": [
       252,
       170,
       172,
       100,
       100,
       170,
       170,
       170,
       170
      ],
      "desimal": [
       null,
       null,
       null,
       null,
       null,
       null,
       null,
       null,
       null
      ],
      "format": [
       "",
       "",
       "",
       "",
       "",
       "",
       "",
       "",
       ""
      ],
      "syaratSel": [
       null,
       null,
       null,
       null,
       null,
       null,
       null,
       null,
       null
      ],
      "atSel": [
       1125421,
       1129850,
       1134288,
       1138727,
       1143162,
       1147598,
       1152044,
       1156492,
       1160935
      ],
      "templatBaris": "ASM-FW-GISFW-Data-LimitSummaryList!pyGridModalTemplate"
     }
    ]
   },
   {
    "t": "blok",
    "at": 1289784,
    "judul": "Total All Layers",
    "syarat": [],
    "anak": [
     {
      "t": "grid",
      "at": 1316034,
      "prop": "TreatyIn.TotalLimitIOONP",
      "dari": "sisi",
      "larik": "TotalLimitIOONP",
      "syarat": [],
      "kolom": [
       "Total 100% Limit",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       193,
       349
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       1327936,
       1332880
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "grid",
      "at": 1380685,
      "prop": "TreatyIn.TotalLimitDeductblNP",
      "dari": "sisi",
      "larik": "TotalLimitDeductblNP",
      "syarat": [],
      "kolom": [
       "Total Deductible",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       194,
       349
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       1392592,
       1397536
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "grid",
      "at": 1445341,
      "prop": "TreatyIn.TotalLimitPremiEarnNP",
      "dari": "sisi",
      "larik": "TotalLimitPremiEarnNP",
      "syarat": [],
      "kolom": [
       "Total Premium Earned",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       193,
       349
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       1457253,
       1462197
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "grid",
      "at": 1510002,
      "prop": "TreatyIn.TotalLimitMDPNP",
      "dari": "sisi",
      "larik": "TotalLimitMDPNP",
      "syarat": [],
      "kolom": [
       "Total MDP",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       193,
       349
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       1521897,
       1526841
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "teks",
      "at": 1588332,
      "teks": "Total ROL",
      "syarat": []
     },
     {
      "t": "medan",
      "at": 1599054,
      "label": "",
      "dari": "sisi",
      "kunci": "TotalLimitsROL",
      "format": "pxNumber",
      "desimal": 2,
      "syarat": []
     },
     {
      "t": "tombol",
      "at": 1638126,
      "label": "Update Total",
      "syarat": [
       "TreatyIn.ViewState !='1'"
      ]
     },
     {
      "t": "tombol",
      "at": 1652282,
      "label": "Update Value in List",
      "syarat": [
       "TreatyIn.ViewState !='1' && TreatyMasterInEDM"
      ]
     }
    ]
   }
  ]
 },
 "TreatyInTabsNonProportional#Share": {
  "at": 1693064,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 1703355,
    "judul": "Share",
    "syarat": [],
    "anak": [
     {
      "t": "teks",
      "at": 1710157,
      "teks": "Share Non Pro Rate (actual)",
      "syarat": [
       "TreatyIn.IsProRate = True"
      ]
     },
     {
      "t": "medan",
      "at": 1741623,
      "label": "% RNM Share",
      "dari": "sisi",
      "kunci": "RNMShare",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": []
     },
     {
      "t": "medan",
      "at": 1750505,
      "label": "% Brokerage",
      "dari": "sisi",
      "kunci": "BrokeragePercent",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": []
     },
     {
      "t": "medan",
      "at": 1768828,
      "label": "",
      "dari": "sisi",
      "kunci": "RNMShareAcrossTheBoard",
      "format": "pxCheckbox",
      "desimal": null,
      "syarat": [],
      "caption": "Share Across The Board"
     },
     {
      "t": "medan",
      "at": 1797115,
      "label": "Share to Other Retro",
      "dari": "sisi",
      "kunci": "FacultativeShare",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": []
     },
     {
      "t": "medan",
      "at": 1806086,
      "label": "Brokerage From Other Retro",
      "dari": "sisi",
      "kunci": "FacultativeShareBrokerage",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": [
       "TreatyIn.FacultativeShare >0"
      ]
     },
     {
      "t": "tombol",
      "at": 1828787,
      "label": "Update Summary",
      "syarat": [
       "TreatyIn.ViewState !='1'"
      ]
     },
     {
      "t": "grid",
      "at": 1878628,
      "prop": "TreatyIn.ShareReins",
      "dari": "sisi",
      "larik": "ShareReins",
      "syarat": [],
      "kolom": [
       "Reinsurer Name",
       "Layer",
       "% Share"
      ],
      "kunci": [
       "ReinsName",
       "Layer",
       "SharePct"
      ],
      "lebar": [
       424,
       180,
       188
      ],
      "desimal": [
       null,
       null,
       2
      ],
      "format": [
       "pxAutoComplete",
       "pxTextInput",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null,
       null
      ],
      "atSel": [
       1904195,
       1912259,
       1918452
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInShareReins!pyGridModalTemplate"
     },
     {
      "t": "blok",
      "at": 1958098,
      "judul": "",
      "syarat": [
       "TreatyIn.FacultativeShare >0"
      ],
      "anak": [
       {
        "t": "grid",
        "at": 1975456,
        "prop": "TreatyIn.ShareFacultativeReinsurers",
        "dari": "sisi",
        "larik": "ShareFacultativeReinsurers",
        "syarat": [],
        "kolom": [
         "Facultative Reinsurers",
         "Layer",
         "% Share"
        ],
        "kunci": [
         "ReinsName",
         "Layer",
         "SharePct"
        ],
        "lebar": [
         428,
         182,
         189
        ],
        "desimal": [
         null,
         null,
         2
        ],
        "format": [
         "pxAutoComplete",
         "pxTextInput",
         "pxNumber"
        ],
        "syaratSel": [
         null,
         null,
         null
        ],
        "atSel": [
         2001074,
         2009138,
         2015331
        ],
        "templatBaris": "ASM-FW-GISFW-Data-TreatyInShareReins!pyGridModalTemplate"
       }
      ]
     },
     {
      "t": "blok",
      "at": 2091194,
      "judul": "RNM Share",
      "syarat": [],
      "anak": [
       {
        "t": "blok",
        "at": 2100300,
        "judul": "",
        "syarat": [
         "TreatyIn.FacultativeShare >0"
        ],
        "anak": [
         {
          "t": "medan",
          "at": 2107030,
          "label": "Share to RNM :",
          "dari": "sisi",
          "kunci": "RnmShareDeducted",
          "format": "pxTextInput",
          "desimal": null,
          "syarat": []
         },
         {
          "t": "teks",
          "at": 2118388,
          "teks": "%",
          "syarat": []
         }
        ]
       },
       {
        "t": "blok",
        "at": 2131350,
        "judul": "RNM Share",
        "syarat": [],
        "anak": [
         {
          "t": "grid",
          "at": 2148707,
          "prop": "TreatyIn.Share",
          "dari": "sisi",
          "larik": "Share",
          "syarat": [],
          "kolom": [
           "",
           "",
           "",
           "",
           "",
           "",
           "100% Limit",
           "",
           "100% Limit",
           "",
           "MDP",
           "",
           "MDP",
           "% Share"
          ],
          "kunci": [
           "LayerType",
           "Layer",
           "pyTemplateInputBox",
           "LayerPartType",
           "LayerPart",
           "RnmLimitListDisplay(1).Currency",
           "RnmLimitListDisplay(1).Value",
           "RnmLimitListDisplay(2).Currency",
           "RnmLimitListDisplay(2).Value",
           "RnmGrossPremiDisplay(1).Currency",
           "RnmGrossPremiDisplay(1).Value",
           "RnmGrossPremiDisplay(2).Currency",
           "RnmGrossPremiDisplay(2).Value",
           "RNMShare"
          ],
          "lebar": [
           138,
           60,
           92,
           138,
           103,
           75,
           163,
           75,
           163,
           80,
           120,
           80,
           120,
           100
          ],
          "desimal": [
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null
          ],
          "format": [
           "pxTextInput",
           "",
           "",
           "pxTextInput",
           "",
           "pxTextInput",
           "",
           "pxTextInput",
           "",
           "pxTextInput",
           "",
           "pxTextInput",
           "",
           "pxNumber"
          ],
          "syaratSel": [
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null
          ],
          "atSel": [
           2207544,
           2213478,
           2217921,
           2222321,
           2228160,
           2232609,
           2238497,
           2243056,
           2248944,
           2253503,
           2259393,
           2263954,
           2269844,
           2274405
          ],
          "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridRowDetails"
         }
        ]
       },
       {
        "t": "blok",
        "at": 2321865,
        "judul": "Summarry of RNM Share",
        "syarat": [],
        "anak": [
         {
          "t": "grid",
          "at": 2339255,
          "prop": "TreatyIn.LimitShareSummaryList",
          "dari": "sisi",
          "larik": "LimitShareSummaryList",
          "syarat": [],
          "kolom": [
           "Note",
           "100% Limit (IDR)",
           "100% Limit (USD)",
           "MDP (IDR)",
           "MDP (USD)",
           "Deduction (IDR)",
           "Deduction (USD)",
           "Net Premi (IDR)",
           "Net Premi (USD)"
          ],
          "kunci": [
           "Note",
           "Limit",
           "Limit2",
           "MDP",
           "MDP2",
           "Deductible",
           "Deductible2",
           "NetPremi",
           "NetPremi2"
          ],
          "lebar": [
           293,
           199,
           204,
           123,
           119,
           205,
           210,
           150,
           144
          ],
          "desimal": [
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null
          ],
          "format": [
           "pxNumber",
           "pxNumber",
           "pxNumber",
           "pxNumber",
           "pxNumber",
           "pxNumber",
           "pxNumber",
           "pxNumber",
           "pxNumber"
          ],
          "syaratSel": [
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null,
           null
          ],
          "atSel": [
           2379731,
           2384686,
           2389747,
           2394809,
           2399867,
           2404926,
           2409992,
           2415059,
           2420122
          ],
          "templatBaris": "ASM-FW-GISFW-Data-LimitSummaryList!pyGridModalTemplate"
         }
        ]
       }
      ]
     }
    ]
   },
   {
    "t": "blok",
    "at": 2483759,
    "judul": "Total All Layers RNM Share",
    "syarat": [],
    "anak": [
     {
      "t": "grid",
      "at": 2519171,
      "prop": "TreatyIn.TotalShareRnmNP",
      "dari": "sisi",
      "larik": "TotalShareRnmNP",
      "syarat": [],
      "kolom": [
       "Total RNM Limit (RNM Share)",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       194,
       352
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       2531152,
       2536109
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "grid",
      "at": 2584196,
      "prop": "TreatyIn.TotalSpreadedRnmProp",
      "dari": "sisi",
      "larik": "TotalSpreadedRnmProp",
      "syarat": [],
      "kolom": [
       "Total OR Limit",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       192,
       350
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       2596169,
       2601126
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "grid",
      "at": 2649233,
      "prop": "TreatyIn.TotalSpreadedRnmRIProp",
      "dari": "sisi",
      "larik": "TotalSpreadedRnmRIProp",
      "syarat": [],
      "kolom": [
       "Total R/I Limit",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       192,
       350
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       2661210,
       2666167
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "grid",
      "at": 2729593,
      "prop": "TreatyIn.TotalShareGrossMinNP",
      "dari": "sisi",
      "larik": "TotalShareGrossMinNP",
      "syarat": [],
      "kolom": [
       "Total Gross Min Premium",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       194,
       350
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       2741575,
       2746532
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "grid",
      "at": 2809958,
      "prop": "TreatyIn.TotalShareGrossNP",
      "dari": "sisi",
      "larik": "TotalShareGrossNP",
      "syarat": [],
      "kolom": [
       "Total Gross Premium (MDP)",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       194,
       350
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       2821939,
       2826896
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "blok",
      "at": 2873031,
      "judul": "hiddden",
      "syarat": [],
      "anak": [
       {
        "t": "grid",
        "at": 2890324,
        "prop": "TreatyIn.TotalShareDeductionNP",
        "dari": "sisi",
        "larik": "TotalShareDeductionNP",
        "syarat": [],
        "kolom": [
         "Total Deduction",
         "Value"
        ],
        "kunci": [
         "Currency",
         "Value"
        ],
        "lebar": [
         194,
         352
        ],
        "desimal": [
         null,
         2
        ],
        "format": [
         "pxNumber",
         "pxNumber"
        ],
        "syaratSel": [
         null,
         null
        ],
        "atSel": [
         2902299,
         2907256
        ],
        "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
       }
      ]
     },
     {
      "t": "grid",
      "at": 2970722,
      "prop": "TreatyIn.TotalShareNetNP",
      "dari": "sisi",
      "larik": "TotalShareNetNP",
      "syarat": [],
      "kolom": [
       "Total Net Premium",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       193,
       349
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       2982693,
       2987650
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "grid",
      "at": 3035737,
      "prop": "TreatyIn.TotalSpreadedNetPremi",
      "dari": "sisi",
      "larik": "TotalSpreadedNetPremi",
      "syarat": [],
      "kolom": [
       "Total OR Net Premium",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       192,
       348
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       3047717,
       3052674
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "grid",
      "at": 3099627,
      "prop": "TreatyIn.TotalSpreadedNetPremiRI",
      "dari": "sisi",
      "larik": "TotalSpreadedNetPremiRI",
      "syarat": [],
      "kolom": [
       "Total R/I Net Premium",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       192,
       348
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       3111610,
       3116567
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "tombol",
      "at": 3184690,
      "label": "Update Total",
      "syarat": [
       "TreatyIn.ViewState !='1'"
      ]
     },
     {
      "t": "tombol",
      "at": 3198699,
      "label": "Update Value in Share",
      "syarat": [
       "TreatyIn.ViewState !='1' && TreatyMasterInEDM"
      ]
     },
     {
      "t": "tombol",
      "at": 3236226,
      "label": "Hide Facultative Share (unused)",
      "syarat": [
       "TreatyIn.FacultativeShare>0 && facsharedisp.CARI1 = '1'"
      ]
     }
    ]
   }
  ]
 },
 "TreatyInTabsNonProportional#Retro": {
  "at": 3289387,
  "syarat": [],
  "isi": [
   {
    "t": "medan",
    "at": 3305948,
    "label": "Has share to retro:",
    "dari": "sisi",
    "kunci": "IsMultipleRetro",
    "format": "pxCheckbox",
    "desimal": null,
    "syarat": [],
    "caption": "Has Share To Retro:"
   },
   {
    "t": "blok",
    "at": 3419521,
    "judul": "",
    "syarat": [
     "TreatyIn.IsMultipleRetro = 'true'"
    ],
    "anak": [
     {
      "t": "include",
      "at": 3426220,
      "nama": "TreatyInFacultativeShareCalculation",
      "syarat": []
     }
    ]
   }
  ]
 },
 "TreatyInTabsNonProportional#Installment": {
  "at": 3463245,
  "syarat": [],
  "isi": [
   {
    "t": "medan",
    "at": 3489051,
    "label": "Installment",
    "dari": "sisi",
    "kunci": "InstallmentNo",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "tombol",
    "at": 3516400,
    "label": "Update Value",
    "syarat": [
     "TreatyIn.ViewState !='1'"
    ]
   },
   {
    "t": "grid",
    "at": 3567479,
    "prop": "TreatyIn.Installment",
    "dari": "sisi",
    "larik": "Installment",
    "syarat": [],
    "kolom": [
     "Currency"
    ],
    "kunci": [
     "Currency"
    ],
    "lebar": [
     1300
    ],
    "desimal": [
     null
    ],
    "format": [
     ""
    ],
    "syaratSel": [
     null
    ],
    "atSel": [
     3575411
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInInstallment!pyGridRowDetails"
   },
   {
    "t": "grid",
    "at": 3773231,
    "prop": "TreatyIn.TotalInstallmentNP",
    "dari": "sisi",
    "larik": "TotalInstallmentNP",
    "syarat": [],
    "kolom": [
     "Total Installment Amount",
     "Value"
    ],
    "kunci": [
     "Currency",
     "Value"
    ],
    "lebar": [
     194,
     349
    ],
    "desimal": [
     null,
     2
    ],
    "format": [
     "pxNumber",
     "pxNumber"
    ],
    "syaratSel": [
     null,
     null
    ],
    "atSel": [
     3785144,
     3790088
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
   },
   {
    "t": "tombol",
    "at": 3824681,
    "label": "Update Total",
    "syarat": [
     "TreatyIn.ViewState !='1'"
    ]
   }
  ]
 },
 "TreatyInTabsNonProportional#Value Difference": {
  "at": 3846190,
  "syarat": [
   "TreatyIn.EDMState != 3 && TreatyIn.EDMMaterialType == 1"
  ],
  "isi": [
   {
    "t": "include",
    "at": 3867690,
    "nama": "TreatyInTabsNonProportionalValueDifference",
    "syarat": []
   }
  ]
 },
 "TreatyInTabsNonProportional#Exclusions": {
  "at": 3904723,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 3915122,
    "judul": "Exclusions",
    "syarat": [],
    "anak": [
     {
      "t": "medan",
      "at": 3921877,
      "label": "",
      "dari": "sisi",
      "kunci": "Exclusions",
      "format": "pxTextArea",
      "desimal": null,
      "syarat": []
     }
    ]
   }
  ]
 },
 "TreatyInTabsNonProportional#Special Conditions": {
  "at": 3940127,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 3950103,
    "judul": "Special Conditions",
    "syarat": [],
    "anak": [
     {
      "t": "medan",
      "at": 3956866,
      "label": "",
      "dari": "sisi",
      "kunci": "SpecialConditions",
      "format": "pxTextArea",
      "desimal": null,
      "syarat": []
     }
    ]
   }
  ]
 },
 "TreatyInTabsNonProportional#Information & Submit": {
  "at": 3975084,
  "syarat": [],
  "isi": [
   {
    "t": "include",
    "at": 3982960,
    "nama": "TreatyInfoSubmit",
    "syarat": []
   }
  ]
 },
 "TreatyInTabsNonProportionalOldData#Maximum Retention": {
  "at": 20697,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 31226,
    "judul": "Maximum Retention",
    "syarat": [],
    "anak": [
     {
      "t": "grid",
      "at": 48622,
      "prop": "TreatyIn.OLDDATA.Retention",
      "dari": "sisi",
      "larik": "Retention",
      "syarat": [],
      "kolom": [
       "Treaty Group",
       "Currency",
       "Amount"
      ],
      "kunci": [
       "TreatyGroup",
       "Currency",
       "Amount"
      ],
      "lebar": [
       243,
       126,
       376
      ],
      "desimal": [
       null,
       null,
       null
      ],
      "format": [
       "pxAutoComplete",
       "",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null,
       null
      ],
      "atSel": [
       65013,
       72792,
       77261
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInRetention!pyGridRowDetails"
     }
    ]
   },
   {
    "t": "grid",
    "at": 136305,
    "prop": "TreatyIn.OLDDATA.TotalRetentionAmountNP",
    "dari": "sisi",
    "larik": "TotalRetentionAmountNP",
    "syarat": [],
    "kolom": [
     "Total Retention Amount",
     "Value"
    ],
    "kunci": [
     "Currency",
     "Value"
    ],
    "lebar": [
     193,
     349
    ],
    "desimal": [
     null,
     2
    ],
    "format": [
     "pxNumber",
     "pxNumber"
    ],
    "syaratSel": [
     null,
     null
    ],
    "atSel": [
     148292,
     153248
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
   }
  ]
 },
 "TreatyInTabsNonProportionalOldData#EGNPI": {
  "at": 203767,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 222812,
    "judul": "Estimate Gross Net Premium Income",
    "syarat": [],
    "anak": [
     {
      "t": "grid",
      "at": 240206,
      "prop": "TreatyIn.OLDDATA.EGNPI",
      "dari": "sisi",
      "larik": "EGNPI",
      "syarat": [],
      "kolom": [
       "Treaty Group",
       "As Date",
       "Proportion %",
       "Currency",
       "Amount",
       "Amount in IDR"
      ],
      "kunci": [
       "TreatyGroup",
       "AsDate",
       "Proportion",
       "Currency",
       "Amount",
       "AmountIDR"
      ],
      "lebar": [
       284,
       195,
       140,
       136,
       219,
       220
      ],
      "desimal": [
       null,
       null,
       2,
       null,
       2,
       null
      ],
      "format": [
       "pxAutoComplete",
       "",
       "pxNumber",
       "",
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null,
       null,
       null,
       null,
       null
      ],
      "atSel": [
       269147,
       274603,
       279049,
       284146,
       288593,
       293685
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInEGNPI!pyGridRowDetails"
     }
    ]
   },
   {
    "t": "grid",
    "at": 353999,
    "prop": "TreatyIn.OLDDATA.TotalEgnpiAmountNP",
    "dari": "sisi",
    "larik": "TotalEgnpiAmountNP",
    "syarat": [],
    "kolom": [
     "Total EGNPI Amount",
     "Value"
    ],
    "kunci": [
     "Currency",
     "Value"
    ],
    "lebar": [
     194,
     349
    ],
    "desimal": [
     null,
     2
    ],
    "format": [
     "pxNumber",
     "pxNumber"
    ],
    "syaratSel": [
     null,
     null
    ],
    "atSel": [
     365982,
     370939
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
   },
   {
    "t": "teks",
    "at": 432765,
    "teks": "Total Amount in IDR",
    "syarat": []
   },
   {
    "t": "teks",
    "at": 452509,
    "teks": "IDR",
    "syarat": []
   },
   {
    "t": "medan",
    "at": 456898,
    "label": "",
    "dari": "sisi",
    "kunci": "TotalEgnpiAmount",
    "format": "pxNumber",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "teks",
    "at": 489848,
    "teks": "Total Proportion %",
    "syarat": []
   },
   {
    "t": "medan",
    "at": 500592,
    "label": "",
    "dari": "sisi",
    "kunci": "TotalEgnpiProportion",
    "format": "pxNumber",
    "desimal": 2,
    "syarat": []
   }
  ]
 },
 "TreatyInTabsNonProportionalOldData#Limits": {
  "at": 552942,
  "syarat": [],
  "isi": [
   {
    "t": "grid",
    "at": 589279,
    "prop": "TreatyIn.OLDDATA.Limits",
    "dari": "sisi",
    "larik": "Limits",
    "syarat": [],
    "kolom": [
     "Layers",
     "",
     "",
     "",
     "",
     "Currency",
     "100% Limits",
     "Deductible"
    ],
    "kunci": [
     "LayerType",
     "Layer",
     "pyTemplateRichTextEditor",
     "LayerPartType",
     "LayerPart",
     "Currency",
     "Limit",
     "Deductible"
    ],
    "lebar": [
     101,
     101,
     101,
     101,
     104,
     122,
     182,
     182
    ],
    "desimal": [
     null,
     null,
     null,
     null,
     null,
     null,
     2,
     2
    ],
    "format": [
     "pxTextInput",
     "",
     "pxTextInput",
     "pxTextInput",
     "",
     "pxNumber",
     "pxNumber",
     "pxNumber"
    ],
    "syaratSel": [
     null,
     null,
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "atSel": [
     626067,
     631999,
     636445,
     642550,
     648390,
     652840,
     657987,
     663080
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails"
   },
   {
    "t": "blok",
    "at": 698474,
    "judul": "Summary of Limit",
    "syarat": [],
    "anak": [
     {
      "t": "grid",
      "at": 724876,
      "prop": "TreatyIn.OLDDATA.LimitSummaryList",
      "dari": "sisi",
      "larik": "LimitSummaryList",
      "syarat": [],
      "kolom": [
       "Note",
       "100% Limit (IDR)",
       "100% Limit (USD)",
       "MDP (IDR)",
       "MDP (USD)",
       "Agregate Limit (IDR)",
       "Agregate Limit (USD)",
       "Deductible (IDR)",
       "Deductible (USD)"
      ],
      "kunci": [
       "Note",
       "Limit",
       "Limit2",
       "MDP",
       "MDP2",
       "AggregateLimit",
       "AggregateLimit2",
       "Deductible",
       "Deductible2"
      ],
      "lebar": [
       252,
       170,
       172,
       100,
       100,
       170,
       170,
       170,
       170
      ],
      "desimal": [
       null,
       null,
       null,
       null,
       null,
       null,
       null,
       null,
       null
      ],
      "format": [
       "",
       "",
       "",
       "",
       "",
       "",
       "",
       "",
       ""
      ],
      "syaratSel": [
       null,
       null,
       null,
       null,
       null,
       null,
       null,
       null,
       null
      ],
      "atSel": [
       765368,
       769810,
       774261,
       778713,
       783161,
       787610,
       792069,
       796530,
       800986
      ],
      "templatBaris": "ASM-FW-GISFW-Data-LimitSummaryList!pyGridModalTemplate"
     }
    ]
   },
   {
    "t": "blok",
    "at": 930602,
    "judul": "Total All Layers",
    "syarat": [],
    "anak": [
     {
      "t": "grid",
      "at": 957005,
      "prop": "TreatyIn.OLDDATA.TotalLimitIOONP",
      "dari": "sisi",
      "larik": "TotalLimitIOONP",
      "syarat": [],
      "kolom": [
       "Total 100% Limit",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       193,
       349
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       968983,
       973940
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "grid",
      "at": 1022027,
      "prop": "TreatyIn.OLDDATA.TotalLimitDeductblNP",
      "dari": "sisi",
      "larik": "TotalLimitDeductblNP",
      "syarat": [],
      "kolom": [
       "Total Deductible",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       194,
       349
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       1034010,
       1038967
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "grid",
      "at": 1087054,
      "prop": "TreatyIn.OLDDATA.TotalLimitPremiEarnNP",
      "dari": "sisi",
      "larik": "TotalLimitPremiEarnNP",
      "syarat": [],
      "kolom": [
       "Total Premium Earned",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       193,
       349
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       1099042,
       1103999
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "grid",
      "at": 1152086,
      "prop": "TreatyIn.OLDDATA.TotalLimitMDPNP",
      "dari": "sisi",
      "larik": "TotalLimitMDPNP",
      "syarat": [],
      "kolom": [
       "Total MDP",
       "Value"
      ],
      "kunci": [
       "Currency",
       "Value"
      ],
      "lebar": [
       193,
       349
      ],
      "desimal": [
       null,
       2
      ],
      "format": [
       "pxNumber",
       "pxNumber"
      ],
      "syaratSel": [
       null,
       null
      ],
      "atSel": [
       1164057,
       1169014
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "teks",
      "at": 1230751,
      "teks": "Total ROL",
      "syarat": []
     },
     {
      "t": "medan",
      "at": 1241486,
      "label": "",
      "dari": "sisi",
      "kunci": "TotalLimitsROL",
      "format": "pxNumber",
      "desimal": 2,
      "syarat": []
     }
    ]
   }
  ]
 },
 "TreatyInTabsNonProportionalOldData#Share": {
  "at": 1287504,
  "syarat": [],
  "isi": [
   {
    "t": "include",
    "at": 1304635,
    "nama": "TreatyInTabsNonProportionalOldDataShare",
    "syarat": []
   }
  ]
 },
 "TreatyInTabsNonProportionalOldData#Retro": {
  "at": 1344913,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 1355371,
    "judul": "",
    "syarat": [
     "TreatyIn.IsMultipleRetro = 'true'"
    ],
    "anak": [
     {
      "t": "include",
      "at": 1362092,
      "nama": "TreatyInFacultativeShareCalculationOldData",
      "syarat": []
     }
    ]
   }
  ]
 },
 "TreatyInTabsNonProportionalOldData#Installment": {
  "at": 1404109,
  "syarat": [],
  "isi": [
   {
    "t": "medan",
    "at": 1430204,
    "label": "Installment",
    "dari": "sisi",
    "kunci": "InstallmentNo",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "grid",
    "at": 1484350,
    "prop": "TreatyIn.OLDDATA.Installment",
    "dari": "sisi",
    "larik": "Installment",
    "syarat": [],
    "kolom": [
     "Currency"
    ],
    "kunci": [
     "Currency"
    ],
    "lebar": [
     1300
    ],
    "desimal": [
     null
    ],
    "format": [
     ""
    ],
    "syaratSel": [
     null
    ],
    "atSel": [
     1492345
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInInstallment!pyGridRowDetails"
   },
   {
    "t": "grid",
    "at": 1690755,
    "prop": "TreatyIn.OLDDATA.TotalInstallmentNP",
    "dari": "sisi",
    "larik": "TotalInstallmentNP",
    "syarat": [],
    "kolom": [
     "Total Installment Amount",
     "Value"
    ],
    "kunci": [
     "Currency",
     "Value"
    ],
    "lebar": [
     194,
     349
    ],
    "desimal": [
     null,
     2
    ],
    "format": [
     "pxNumber",
     "pxNumber"
    ],
    "syaratSel": [
     null,
     null
    ],
    "atSel": [
     1702744,
     1707701
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
   }
  ]
 },
 "TreatyInTabsNonProportionalOldData#Exclusions": {
  "at": 1758152,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 1768688,
    "judul": "Exclusions",
    "syarat": [],
    "anak": [
     {
      "t": "medan",
      "at": 1775459,
      "label": "",
      "dari": "akar",
      "kunci": "Exclusions",
      "format": "pxTextArea",
      "desimal": null,
      "syarat": []
     }
    ]
   }
  ]
 },
 "TreatyInTabsNonProportionalOldData#Special Conditions": {
  "at": 1796524,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 1806635,
    "judul": "Special Conditions",
    "syarat": [],
    "anak": [
     {
      "t": "medan",
      "at": 1813414,
      "label": "",
      "dari": "akar",
      "kunci": "SpecialConditions",
      "format": "pxTextArea",
      "desimal": null,
      "syarat": []
     }
    ]
   }
  ]
 },
 "TreatyInTabsProportional#Reporting Period": {
  "at": 23811,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 43323,
    "judul": "Account Reporting Period",
    "syarat": [],
    "anak": []
   },
   {
    "t": "medan",
    "at": 98997,
    "label": "Start Date",
    "dari": "sisi",
    "kunci": "ReportingStart",
    "format": "pxDateTime",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "medan",
    "at": 120405,
    "label": "End Date",
    "dari": "sisi",
    "kunci": "ReportingEnd",
    "format": "pxDateTime",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "medan",
    "at": 139204,
    "label": "Period",
    "dari": "sisi",
    "kunci": "ReportingPeriod",
    "format": "pxDropdown",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "medan",
    "at": 144723,
    "label": "Interval",
    "dari": "sisi",
    "kunci": "ReportingInterval",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": [
     "TreatyIn.ReportingPeriod='other'"
    ]
   },
   {
    "t": "medan",
    "at": 182295,
    "label": "Submission",
    "dari": "sisi",
    "kunci": "ReportingSubmission",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "medan",
    "at": 204274,
    "label": "Confirmation",
    "dari": "sisi",
    "kunci": "ReportingConfirmation",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "medan",
    "at": 226262,
    "label": "Settlement",
    "dari": "sisi",
    "kunci": "ReportingSettlement",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "blok",
    "at": 247909,
    "judul": "Account Reporting Period",
    "syarat": [],
    "anak": [
     {
      "t": "tombol",
      "at": 254615,
      "label": "Apply",
      "syarat": [
       "TreatyIn.IsEditData !='1'"
      ]
     },
     {
      "t": "teks",
      "at": 270162,
      "teks": "Start Date, Due",
      "syarat": []
     },
     {
      "t": "teks",
      "at": 274479,
      "teks": ", and Interval",
      "syarat": [
       "TreatyIn.ReportingPeriod='other'"
      ]
     },
     {
      "t": "teks",
      "at": 279526,
      "teks": ". Must Not Be Empty",
      "syarat": []
     }
    ]
   },
   {
    "t": "grid",
    "at": 304113,
    "prop": "TreatyIn.ReportingPeriodList",
    "dari": "sisi",
    "larik": "ReportingPeriodList",
    "syarat": [],
    "kolom": [
     "Period",
     "Auto Calculate",
     "Initial Date",
     "Submission Due",
     "Confirmation Due",
     "Settlement Due"
    ],
    "kunci": [
     "Period",
     "AutoCalculate",
     "InitialDate",
     "SubmissionDue",
     "ConfirmationDue",
     "SettlementDue"
    ],
    "lebar": [
     220,
     104,
     217,
     228,
     228,
     218
    ],
    "desimal": [
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "format": [
     "pxTextInput",
     "pxCheckbox",
     "pxDateTime",
     "pxDateTime",
     "pxDateTime",
     "pxDateTime"
    ],
    "syaratSel": [
     null,
     "TreatyIn.ViewState != 1",
     null,
     null,
     null,
     null
    ],
    "atSel": [
     333359,
     339333,
     345959,
     355661,
     361774,
     367889
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInAccountReport!pyGridRowDetails"
   }
  ]
 },
 "TreatyInTabsProportional#Portfolio": {
  "at": 414663,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 424764,
    "judul": "Portfolio",
    "syarat": [],
    "anak": [
     {
      "t": "grid",
      "at": 442135,
      "prop": "TreatyIn.Portfolio",
      "dari": "sisi",
      "larik": "Portfolio",
      "syarat": [],
      "kolom": [
       "Portfolio Type",
       "Premium / Loss Type",
       "Description"
      ],
      "kunci": [
       "TypePortfolio",
       "Type",
       "Description"
      ],
      "lebar": [
       111,
       157,
       648
      ],
      "desimal": [
       null,
       null,
       null
      ],
      "format": [
       "",
       "",
       "pxTextArea"
      ],
      "syaratSel": [
       null,
       null,
       null
      ],
      "atSel": [
       467551,
       472622,
       477569
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInPortfolio!pyGridModalTemplate"
     }
    ]
   }
  ]
 },
 "TreatyInTabsProportional#Limits": {
  "at": 524276,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 534799,
    "judul": "Limits",
    "syarat": [],
    "anak": [
     {
      "t": "grid",
      "at": 567629,
      "prop": "TreatyIn.Limits",
      "dari": "sisi",
      "larik": "Limits",
      "syarat": [],
      "kolom": [
       "Kind of Treaty"
      ],
      "kunci": [
       "TreatyType"
      ],
      "lebar": [
       1315
      ],
      "desimal": [
       null
      ],
      "format": [
       "pxAutoComplete"
      ],
      "syaratSel": [
       null
      ],
      "atSel": [
       584772
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails"
     }
    ]
   }
  ]
 },
 "TreatyInTabsProportional#Share": {
  "at": 649111,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 659633,
    "judul": "Total Share",
    "syarat": [],
    "anak": [
     {
      "t": "tombol",
      "at": 684428,
      "label": "Refresh",
      "syarat": [
       "TreatyIn.ViewState !='1'"
      ]
     },
     {
      "t": "medan",
      "at": 711242,
      "label": "% RNM Share",
      "dari": "sisi",
      "kunci": "RNMShareP",
      "format": "pxNumber",
      "desimal": 2,
      "syarat": []
     },
     {
      "t": "medan",
      "at": 717944,
      "label": "% Brokerage",
      "dari": "sisi",
      "kunci": "BrokeragePercentP",
      "format": "pxNumber",
      "desimal": 2,
      "syarat": []
     },
     {
      "t": "medan",
      "at": 723387,
      "label": "Option",
      "dari": "sisi",
      "kunci": "OptionLimit",
      "format": "pxDropdown",
      "desimal": null,
      "syarat": []
     },
     {
      "t": "medan",
      "at": 780884,
      "label": "Share to Other Retro",
      "dari": "sisi",
      "kunci": "FacShare",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": [
       "TreatyIn.IsMultipleRetro"
      ]
     },
     {
      "t": "medan",
      "at": 790889,
      "label": "Brokerage From Other Retro",
      "dari": "sisi",
      "kunci": "FacShareBrokerage",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": [
       "TreatyIn.IsMultipleRetro"
      ]
     },
     {
      "t": "blok",
      "at": 842608,
      "judul": "",
      "syarat": [
       "TreatyIn.IsMultipleRetro"
      ],
      "anak": [
       {
        "t": "grid",
        "at": 860338,
        "prop": "TreatyIn.ShareFacultativeReinsurers",
        "dari": "sisi",
        "larik": "ShareFacultativeReinsurers",
        "syarat": [],
        "kolom": [
         "Facultative Reinsurers",
         "Broker",
         "% Share"
        ],
        "kunci": [
         "ReinsName",
         "BrokerName",
         "SharePct"
        ],
        "lebar": [
         323,
         204,
         148
        ],
        "desimal": [
         null,
         null,
         2
        ],
        "format": [
         "pxLink",
         "pxLink",
         "pxNumber"
        ],
        "syaratSel": [
         null,
         null,
         null
        ],
        "atSel": [
         886836,
         904928,
         922680
        ],
        "templatBaris": "ASM-FW-GISFW-Data-TreatyInShareReins!pyGridModalTemplate"
       }
      ]
     }
    ]
   },
   {
    "t": "include",
    "at": 989614,
    "nama": "TreatyInShareProp",
    "syarat": []
   }
  ]
 },
 "TreatyInTabsProportional#Retro": {
  "at": 1028337,
  "syarat": [
   "TreatyIn.IsMultipleRetro"
  ],
  "isi": [
   {
    "t": "include",
    "at": 1036283,
    "nama": "TreatyInFacultativeRetro",
    "syarat": []
   }
  ]
 },
 "TreatyInTabsProportional#Co-Ins Scale": {
  "at": 1075235,
  "syarat": [],
  "isi": [
   {
    "t": "medan",
    "at": 1190828,
    "label": "Max Co-Insurance Panel (Non Group)",
    "dari": "sisi",
    "kunci": "MaxCoNonGroup",
    "format": "pxNumber",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "medan",
    "at": 1198436,
    "label": "Max Co-Insurance Panel (Group)",
    "dari": "sisi",
    "kunci": "MaxCoGroup",
    "format": "pxNumber",
    "desimal": null,
    "syarat": []
   }
  ]
 },
 "TreatyInTabsProportional#Accumulation": {
  "at": 1228049,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 1247565,
    "judul": "Accumulation Control",
    "syarat": [],
    "anak": []
   },
   {
    "t": "blok",
    "at": 1269571,
    "judul": "Account Reporting Period",
    "syarat": [],
    "anak": [
     {
      "t": "medan",
      "at": 1276277,
      "label": "Period",
      "dari": "sisi",
      "kunci": "AccumulationPeriod",
      "format": "pxDropdown",
      "desimal": null,
      "syarat": []
     }
    ]
   },
   {
    "t": "grid",
    "at": 1344630,
    "prop": "TreatyIn.AccumulationList",
    "dari": "sisi",
    "larik": "AccumulationList",
    "syarat": [],
    "kolom": [
     "Period",
     "Reporting Date",
     "Submission Days",
     "Submission Due"
    ],
    "kunci": [
     "Period",
     "ReportDate",
     "SubDays",
     "SubDueDate"
    ],
    "lebar": [
     98,
     202,
     199,
     215
    ],
    "desimal": [
     null,
     null,
     null,
     null
    ],
    "format": [
     "pxTextInput",
     "pxDateTime",
     "pxTextInput",
     "pxDateTime"
    ],
    "syaratSel": [
     null,
     null,
     null,
     null
    ],
    "atSel": [
     1374097,
     1380069,
     1389247,
     1398507
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInAccumulation!pyGridRowDetails"
   }
  ]
 },
 "TreatyInTabsProportional#Exclusions": {
  "at": 1452568,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 1463104,
    "judul": "Exclusions",
    "syarat": [],
    "anak": [
     {
      "t": "medan",
      "at": 1469875,
      "label": "",
      "dari": "sisi",
      "kunci": "ExclusionsP",
      "format": "pxTextArea",
      "desimal": null,
      "syarat": []
     }
    ]
   }
  ]
 },
 "TreatyInTabsProportional#Special Conditions": {
  "at": 1494461,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 1505005,
    "judul": "Special Conditions",
    "syarat": [],
    "anak": [
     {
      "t": "medan",
      "at": 1511784,
      "label": "",
      "dari": "sisi",
      "kunci": "SpecialConditionsP",
      "format": "pxTextArea",
      "desimal": null,
      "syarat": []
     }
    ]
   }
  ]
 },
 "TreatyInTabsProportional#Information & Submit": {
  "at": 1536386,
  "syarat": [],
  "isi": [
   {
    "t": "include",
    "at": 1544383,
    "nama": "TreatyInfoSubmit",
    "syarat": []
   }
  ]
 },
 "TreatyInTabsProportional#Achievement In IDR": {
  "at": 1578427,
  "syarat": [],
  "isi": [
   {
    "t": "include",
    "at": 1595116,
    "nama": "TreatyInTabsAchievement",
    "syarat": []
   },
   {
    "t": "grid",
    "at": 1639375,
    "prop": "TreatyIn.Limits",
    "dari": "sisi",
    "larik": "Limits",
    "syarat": [],
    "kolom": [
     "Kind of Treaty"
    ],
    "kunci": [
     "TreatyType"
    ],
    "lebar": [
     1315
    ],
    "desimal": [
     null
    ],
    "format": [
     "pxAutoComplete"
    ],
    "syaratSel": [
     null
    ],
    "atSel": [
     1647348
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails"
   }
  ]
 },
 "TreatyInTabsProportionalOldData#Reporting Period": {
  "at": 20392,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 39911,
    "judul": "Account Reporting Period",
    "syarat": [],
    "anak": []
   },
   {
    "t": "medan",
    "at": 95609,
    "label": "Start Date",
    "dari": "sisi",
    "kunci": "ReportingStart",
    "format": "pxDateTime",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "medan",
    "at": 116989,
    "label": "End Date",
    "dari": "sisi",
    "kunci": "ReportingEnd",
    "format": "pxDateTime",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "medan",
    "at": 135711,
    "label": "Period",
    "dari": "sisi",
    "kunci": "ReportingPeriod",
    "format": "pxDropdown",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "medan",
    "at": 141196,
    "label": "Interval",
    "dari": "sisi",
    "kunci": "ReportingInterval",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": [
     "TreatyIn.ReportingPeriod='other'"
    ]
   },
   {
    "t": "medan",
    "at": 178702,
    "label": "Submission",
    "dari": "sisi",
    "kunci": "ReportingSubmission",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "medan",
    "at": 200222,
    "label": "Confirmation",
    "dari": "sisi",
    "kunci": "ReportingConfirmation",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "medan",
    "at": 221750,
    "label": "Settlement",
    "dari": "sisi",
    "kunci": "ReportingSettlement",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": []
   },
   {
    "t": "grid",
    "at": 254396,
    "prop": "TreatyIn.OLDDATA.ReportingPeriodList",
    "dari": "sisi",
    "larik": "ReportingPeriodList",
    "syarat": [],
    "kolom": [
     "Period",
     "Auto Calculate",
     "Initial Date",
     "Submission Due",
     "Confirmation Due",
     "Settlement Due"
    ],
    "kunci": [
     "Period",
     "AutoCalculate",
     "InitialDate",
     "SubmissionDue",
     "ConfirmationDue",
     "SettlementDue"
    ],
    "lebar": [
     219,
     107,
     217,
     228,
     228,
     218
    ],
    "desimal": [
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "format": [
     "pxTextInput",
     "pxCheckbox",
     "pxDateTime",
     "pxDateTime",
     "pxDateTime",
     "pxDateTime"
    ],
    "syaratSel": [
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "atSel": [
     283659,
     289641,
     296214,
     305783,
     311763,
     317745
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInAccountReport!pyGridRowDetails"
   }
  ]
 },
 "TreatyInTabsProportionalOldData#Portfolio": {
  "at": 365061,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 375164,
    "judul": "Portfolio",
    "syarat": [],
    "anak": [
     {
      "t": "grid",
      "at": 392543,
      "prop": "TreatyIn.OLDDATA.Portfolio",
      "dari": "sisi",
      "larik": "Portfolio",
      "syarat": [],
      "kolom": [
       "Portfolio Type",
       "Premium / Loss Type",
       "Description"
      ],
      "kunci": [
       "TypePortfolio",
       "Type",
       "Description"
      ],
      "lebar": [
       110,
       157,
       647
      ],
      "desimal": [
       null,
       null,
       null
      ],
      "format": [
       "",
       "",
       "pxTextArea"
      ],
      "syaratSel": [
       null,
       null,
       null
      ],
      "atSel": [
       408675,
       413206,
       417613
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInPortfolio!pyGridModalTemplate"
     }
    ]
   }
  ]
 },
 "TreatyInTabsProportionalOldData#Limits": {
  "at": 456157,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 466681,
    "judul": "Limits",
    "syarat": [],
    "anak": [
     {
      "t": "grid",
      "at": 499520,
      "prop": "TreatyIn.OLDDATA.Limits",
      "dari": "sisi",
      "larik": "Limits",
      "syarat": [],
      "kolom": [
       "Kind of Treaty"
      ],
      "kunci": [
       "TreatyType"
      ],
      "lebar": [
       1313
      ],
      "desimal": [
       null
      ],
      "format": [
       "pxAutoComplete"
      ],
      "syaratSel": [
       null
      ],
      "atSel": [
       507483
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails"
     }
    ]
   }
  ]
 },
 "TreatyInTabsProportionalOldData#Share": {
  "at": 554245,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 564768,
    "judul": "Total Share",
    "syarat": [],
    "anak": [
     {
      "t": "medan",
      "at": 589556,
      "label": "% RNM Share",
      "dari": "sisi",
      "kunci": "RNMShareP",
      "format": "pxNumber",
      "desimal": 2,
      "syarat": []
     },
     {
      "t": "medan",
      "at": 599000,
      "label": "% Brokerage",
      "dari": "sisi",
      "kunci": "BrokeragePercentP",
      "format": "pxNumber",
      "desimal": 2,
      "syarat": []
     }
    ]
   },
   {
    "t": "grid",
    "at": 645459,
    "prop": "TreatyIn.OLDDATA.Limits",
    "dari": "sisi",
    "larik": "Limits",
    "syarat": [],
    "kolom": [
     "Kind of Treaty"
    ],
    "kunci": [
     "TreatyType"
    ],
    "lebar": [
     800
    ],
    "desimal": [
     null
    ],
    "format": [
     "pxAutoComplete"
    ],
    "syaratSel": [
     null
    ],
    "atSel": [
     653422
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails"
   },
   {
    "t": "grid",
    "at": 713764,
    "prop": "TreatyIn.OLDDATA.TotalShareRnmProp",
    "dari": "sisi",
    "larik": "TotalShareRnmProp",
    "syarat": [],
    "kolom": [
     "Total Share RNM Limit",
     "Value"
    ],
    "kunci": [
     "Currency",
     "Value"
    ],
    "lebar": [
     192,
     347
    ],
    "desimal": [
     null,
     2
    ],
    "format": [
     "pxNumber",
     "pxNumber"
    ],
    "syaratSel": [
     null,
     null
    ],
    "atSel": [
     725749,
     730706
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
   },
   {
    "t": "grid",
    "at": 784556,
    "prop": "TreatyIn.OLDDATA.TotalSpreadedRnmProp",
    "dari": "sisi",
    "larik": "TotalSpreadedRnmProp",
    "syarat": [],
    "kolom": [
     "Total Value Spreading OR",
     "Value"
    ],
    "kunci": [
     "Currency",
     "Value"
    ],
    "lebar": [
     192,
     347
    ],
    "desimal": [
     null,
     2
    ],
    "format": [
     "pxNumber",
     "pxNumber"
    ],
    "syaratSel": [
     null,
     null
    ],
    "atSel": [
     796547,
     801504
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
   },
   {
    "t": "grid",
    "at": 855354,
    "prop": "TreatyIn.OLDDATA.TotalSpreadedRnmRIProp",
    "dari": "sisi",
    "larik": "TotalSpreadedRnmRIProp",
    "syarat": [],
    "kolom": [
     "Total Value Spreading R/I",
     "Value"
    ],
    "kunci": [
     "Currency",
     "Value"
    ],
    "lebar": [
     192,
     347
    ],
    "desimal": [
     null,
     2
    ],
    "format": [
     "pxNumber",
     "pxNumber"
    ],
    "syaratSel": [
     null,
     null
    ],
    "atSel": [
     867348,
     872305
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
   }
  ]
 },
 "TreatyInTabsProportionalOldData#Accumulation": {
  "at": 916409,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 935926,
    "judul": "Accumulation Control",
    "syarat": [],
    "anak": []
   },
   {
    "t": "blok",
    "at": 957940,
    "judul": "Account Reporting Period",
    "syarat": [],
    "anak": [
     {
      "t": "medan",
      "at": 964647,
      "label": "Period",
      "dari": "akar",
      "kunci": "AccumulationPeriod",
      "format": "pxDropdown",
      "desimal": null,
      "syarat": []
     }
    ]
   },
   {
    "t": "grid",
    "at": 1032718,
    "prop": "TreatyIn.OLDDATA.AccumulationList",
    "dari": "sisi",
    "larik": "AccumulationList",
    "syarat": [],
    "kolom": [
     "Period",
     "Reporting Date",
     "Submission Days",
     "Submission Due"
    ],
    "kunci": [
     "Period",
     "ReportDate",
     "SubDays",
     "SubDueDate"
    ],
    "lebar": [
     97,
     202,
     199,
     218
    ],
    "desimal": [
     null,
     null,
     null,
     null
    ],
    "format": [
     "pxTextInput",
     "pxDateTime",
     "pxTextInput",
     "pxDateTime"
    ],
    "syaratSel": [
     null,
     null,
     null,
     null
    ],
    "atSel": [
     1053165,
     1059145,
     1068190,
     1077324
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInAccumulation!pyGridRowDetails"
   }
  ]
 },
 "TreatyInTabsProportionalOldData#Exclusions": {
  "at": 1123409,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 1133947,
    "judul": "Exclusions",
    "syarat": [],
    "anak": [
     {
      "t": "medan",
      "at": 1140719,
      "label": "",
      "dari": "sisi",
      "kunci": "ExclusionsP",
      "format": "pxTextArea",
      "desimal": null,
      "syarat": []
     }
    ]
   }
  ]
 },
 "TreatyInTabsProportionalOldData#Special Conditions": {
  "at": 1164915,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 1175461,
    "judul": "Special Conditions",
    "syarat": [],
    "anak": [
     {
      "t": "medan",
      "at": 1182241,
      "label": "",
      "dari": "sisi",
      "kunci": "SpecialConditionsP",
      "format": "pxTextArea",
      "desimal": null,
      "syarat": []
     }
    ]
   }
  ]
 }
}

export const KERANGKA_INCLUDE: Readonly<Record<string, readonly ButirKerangka[]>> = {
 "TreatyInTabsNonProportionalOldDataShare": [
  {
   "t": "blok",
   "at": 10154,
   "judul": "Share",
   "syarat": [],
   "anak": [
    {
     "t": "medan",
     "at": 35337,
     "label": "% RNM Share",
     "dari": "sisi",
     "kunci": "RNMShare",
     "format": "pxTextInput",
     "desimal": null,
     "syarat": []
    },
    {
     "t": "medan",
     "at": 43965,
     "label": "% Brokerage",
     "dari": "sisi",
     "kunci": "BrokeragePercent",
     "format": "pxTextInput",
     "desimal": null,
     "syarat": []
    },
    {
     "t": "medan",
     "at": 70496,
     "label": "Facultative Share",
     "dari": "sisi",
     "kunci": "FacultativeShare",
     "format": "pxTextInput",
     "desimal": null,
     "syarat": []
    },
    {
     "t": "medan",
     "at": 79208,
     "label": "Fakultative Brokerage",
     "dari": "sisi",
     "kunci": "FacultativeShareBrokerage",
     "format": "pxTextInput",
     "desimal": null,
     "syarat": []
    },
    {
     "t": "grid",
     "at": 120312,
     "prop": "TreatyIn.OLDDATA.ShareReins",
     "dari": "sisi",
     "larik": "ShareReins",
     "syarat": [],
     "kolom": [
      "Reinsurer Name",
      "Layer",
      "% Share"
     ],
     "kunci": [
      "ReinsName",
      "Layer",
      "SharePct"
     ],
     "lebar": [
      427,
      182,
      189
     ],
     "desimal": [
      null,
      null,
      2
     ],
     "format": [
      "pxAutoComplete",
      "pxTextInput",
      "pxNumber"
     ],
     "syaratSel": [
      null,
      null,
      null
     ],
     "atSel": [
      135906,
      143710,
      149131
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInShareReins!pyGridModalTemplate"
    },
    {
     "t": "blok",
     "at": 209454,
     "judul": "RNM Share",
     "syarat": [],
     "anak": [
      {
       "t": "blok",
       "at": 218121,
       "judul": "",
       "syarat": [
        "TreatyIn.FacultativeShare >0"
       ],
       "anak": [
        {
         "t": "medan",
         "at": 224542,
         "label": "Share to RNM :",
         "dari": "sisi",
         "kunci": "RnmShareDeducted",
         "format": "pxTextInput",
         "desimal": null,
         "syarat": []
        },
        {
         "t": "teks",
         "at": 235201,
         "teks": "%",
         "syarat": []
        }
       ]
      },
      {
       "t": "blok",
       "at": 248087,
       "judul": "RNM Share",
       "syarat": [],
       "anak": [
        {
         "t": "grid",
         "at": 264853,
         "prop": "TreatyIn.OLDDATA.Share",
         "dari": "sisi",
         "larik": "Share",
         "syarat": [],
         "kolom": [
          "",
          "",
          "",
          "",
          "",
          "",
          "100% Limit",
          "",
          "100% Limit",
          "",
          "MDP",
          "",
          "MDP"
         ],
         "kunci": [
          "LayerType",
          "Layer",
          "pyTemplateInputBox",
          "LayerPartType",
          "LayerPart",
          "RnmLimitListDisplay(1).Currency",
          "RnmLimitListDisplay(1).Value",
          "RnmLimitListDisplay(2).Currency",
          "RnmLimitListDisplay(2).Value",
          "RnmGrossPremiDisplay(1).Currency",
          "RnmGrossPremiDisplay(1).Value",
          "RnmGrossPremiDisplay(2).Currency",
          "RnmGrossPremiDisplay(2).Value"
         ],
         "lebar": [
          138,
          60,
          92,
          138,
          103,
          75,
          163,
          75,
          163,
          80,
          120,
          80,
          120
         ],
         "desimal": [
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null
         ],
         "format": [
          "pxTextInput",
          "",
          "",
          "pxTextInput",
          "",
          "pxTextInput",
          "",
          "pxTextInput",
          "",
          "pxTextInput",
          "",
          "pxTextInput",
          ""
         ],
         "syaratSel": [
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null
         ],
         "atSel": [
          309200,
          314890,
          319089,
          323246,
          328842,
          333048,
          338693,
          343009,
          348654,
          352970,
          358617,
          362935,
          368582
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridRowDetails"
        }
       ]
      },
      {
       "t": "blok",
       "at": 406230,
       "judul": "Summarry of RNM Share",
       "syarat": [],
       "anak": [
        {
         "t": "grid",
         "at": 422945,
         "prop": "TreatyIn.OLDDATA.LimitShareSummaryList",
         "dari": "sisi",
         "larik": "LimitShareSummaryList",
         "syarat": [],
         "kolom": [
          "Note",
          "100% Limit (IDR)",
          "100% Limit (USD)",
          "MDP (IDR)",
          "MDP (USD)",
          "Deduction (IDR)",
          "Deduction (USD)",
          "Net Premi (IDR)",
          "Net Premi (USD)"
         ],
         "kunci": [
          "Note",
          "Limit",
          "Limit2",
          "MDP",
          "MDP2",
          "Deductible",
          "Deductible2",
          "NetPremi",
          "NetPremi2"
         ],
         "lebar": [
          293,
          199,
          204,
          123,
          119,
          205,
          210,
          150,
          144
         ],
         "desimal": [
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null
         ],
         "format": [
          "",
          "pxNumber",
          "pxNumber",
          "pxNumber",
          "pxNumber",
          "pxNumber",
          "pxNumber",
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null,
          null
         ],
         "atSel": [
          461557,
          465756,
          470574,
          475393,
          480208,
          485024,
          489847,
          494671,
          499491
         ],
         "templatBaris": "ASM-FW-GISFW-Data-LimitSummaryList!pyGridModalTemplate"
        }
       ]
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 550672,
   "judul": "Total All Layers RNM Share",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 586202,
     "prop": "TreatyIn.OLDDATA.TotalShareRnmNP",
     "dari": "sisi",
     "larik": "TotalShareRnmNP",
     "syarat": [],
     "kolom": [
      "Total RNM Limit (RNM Share)",
      "Value"
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      194,
      352
     ],
     "desimal": [
      null,
      2
     ],
     "format": [
      "pxNumber",
      "pxNumber"
     ],
     "syaratSel": [
      null,
      null
     ],
     "atSel": [
      597774,
      602488
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 649772,
     "prop": "TreatyIn.OLDDATA.TotalSpreadedRnmProp",
     "dari": "sisi",
     "larik": "TotalSpreadedRnmProp",
     "syarat": [],
     "kolom": [
      "Total OR Limit",
      "Value"
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      192,
      350
     ],
     "desimal": [
      null,
      2
     ],
     "format": [
      "pxNumber",
      "pxNumber"
     ],
     "syaratSel": [
      null,
      null
     ],
     "atSel": [
      661336,
      666050
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 713334,
     "prop": "TreatyIn.OLDDATA.TotalSpreadedRnmRIProp",
     "dari": "sisi",
     "larik": "TotalSpreadedRnmRIProp",
     "syarat": [],
     "kolom": [
      "Total R/I Limit",
      "Value"
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      192,
      350
     ],
     "desimal": [
      null,
      2
     ],
     "format": [
      "pxNumber",
      "pxNumber"
     ],
     "syaratSel": [
      null,
      null
     ],
     "atSel": [
      724902,
      729616
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 792221,
     "prop": "TreatyIn.OLDDATA.TotalShareGrossNP",
     "dari": "sisi",
     "larik": "TotalShareGrossNP",
     "syarat": [],
     "kolom": [
      "Total Gross Premium (MDP)",
      "Value"
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      194,
      350
     ],
     "desimal": [
      null,
      2
     ],
     "format": [
      "pxNumber",
      "pxNumber"
     ],
     "syaratSel": [
      null,
      null
     ],
     "atSel": [
      803794,
      808508
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "blok",
     "at": 854379,
     "judul": "hiddden",
     "syarat": [],
     "anak": [
      {
       "t": "grid",
       "at": 871088,
       "prop": "TreatyIn.OLDDATA.TotalShareDeductionNP",
       "dari": "sisi",
       "larik": "TotalShareDeductionNP",
       "syarat": [],
       "kolom": [
        "Total Deduction",
        "Value"
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        194,
        352
       ],
       "desimal": [
        null,
        2
       ],
       "format": [
        "pxNumber",
        "pxNumber"
       ],
       "syaratSel": [
        null,
        null
       ],
       "atSel": [
        882655,
        887369
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    },
    {
     "t": "grid",
     "at": 949966,
     "prop": "TreatyIn.OLDDATA.TotalShareNetNP",
     "dari": "sisi",
     "larik": "TotalShareNetNP",
     "syarat": [],
     "kolom": [
      "Total Net Premium",
      "Value"
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      193,
      349
     ],
     "desimal": [
      null,
      2
     ],
     "format": [
      "pxNumber",
      "pxNumber"
     ],
     "syaratSel": [
      null,
      null
     ],
     "atSel": [
      961528,
      966242
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 1013526,
     "prop": "TreatyIn.OLDDATA.TotalSpreadedNetPremi",
     "dari": "sisi",
     "larik": "TotalSpreadedNetPremi",
     "syarat": [],
     "kolom": [
      "Total OR Net Premium",
      "Value"
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      192,
      348
     ],
     "desimal": [
      null,
      2
     ],
     "format": [
      "pxNumber",
      "pxNumber"
     ],
     "syaratSel": [
      null,
      null
     ],
     "atSel": [
      1025097,
      1029811
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 1075941,
     "prop": "TreatyIn.OLDDATA.TotalSpreadedNetPremiRI",
     "dari": "sisi",
     "larik": "TotalSpreadedNetPremiRI",
     "syarat": [],
     "kolom": [
      "Total R/I Net Premium",
      "Value"
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      192,
      348
     ],
     "desimal": [
      null,
      2
     ],
     "format": [
      "pxNumber",
      "pxNumber"
     ],
     "syaratSel": [
      null,
      null
     ],
     "atSel": [
      1087515,
      1092229
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "tombol",
     "at": 1171408,
     "label": "Hide Facultative Share (unused)",
     "syarat": [
      "TreatyIn.FacultativeShare>0 && facsharedisp.CARI1 = '1'"
     ]
    },
    {
     "t": "tombol",
     "at": 1180222,
     "label": "Show Facultative Share",
     "syarat": [
      "TreatyIn.FacultativeShare>0"
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1211455,
   "judul": "",
   "syarat": [
    "facsharedisp.CARI1 = 1"
   ],
   "anak": [
    {
     "t": "blok",
     "at": 1221737,
     "judul": "Facultative Share Calculation",
     "syarat": [],
     "anak": [
      {
       "t": "blok",
       "at": 1230506,
       "judul": "Facultative Share",
       "syarat": [],
       "anak": [
        {
         "t": "blok",
         "at": 1239181,
         "judul": "Facultative Share",
         "syarat": [],
         "anak": [
          {
           "t": "grid",
           "at": 1255963,
           "prop": "TreatyIn.FacultativeShareList",
           "dari": "akar",
           "larik": "FacultativeShareList",
           "syarat": [],
           "kolom": [
            "",
            "",
            "",
            "",
            "",
            "",
            "100% Limit",
            "",
            "100% Limit",
            "",
            "MDP",
            "",
            "MDP"
           ],
           "kunci": [
            "LayerType",
            "Layer",
            "pyTemplateInputBox",
            "LayerPartType",
            "LayerPart",
            "RnmLimitListDisplay(1).Currency",
            "RnmLimitListDisplay(1).Value",
            "RnmLimitListDisplay(2).Currency",
            "RnmLimitListDisplay(2).Value",
            "RnmGrossPremiDisplay(1).Currency",
            "RnmGrossPremiDisplay(1).Value",
            "RnmGrossPremiDisplay(2).Currency",
            "RnmGrossPremiDisplay(2).Value"
           ],
           "lebar": [
            138,
            60,
            92,
            138,
            103,
            75,
            163,
            75,
            163,
            80,
            120,
            80,
            120
           ],
           "desimal": [
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "format": [
            "pxTextInput",
            "",
            "",
            "pxTextInput",
            "",
            "pxTextInput",
            "",
            "pxTextInput",
            "",
            "pxTextInput",
            "",
            "pxTextInput",
            ""
           ],
           "syaratSel": [
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "atSel": [
            1309365,
            1315056,
            1319256,
            1323413,
            1329009,
            1333215,
            1338860,
            1343176,
            1348821,
            1353137,
            1358784,
            1363102,
            1368749
           ],
           "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridRowDetails"
          }
         ]
        },
        {
         "t": "blok",
         "at": 1414294,
         "judul": "Summarry of Facultative Share",
         "syarat": [],
         "anak": [
          {
           "t": "grid",
           "at": 1431018,
           "prop": "TreatyIn.LimitFacShareSummaryList",
           "dari": "akar",
           "larik": "LimitFacShareSummaryList",
           "syarat": [],
           "kolom": [
            "Note",
            "100% Limit (IDR)",
            "100% Limit (USD)",
            "MDP (IDR)",
            "MDP (USD)",
            "Deduction (IDR)",
            "Deduction (USD)",
            "Net Premi (IDR)",
            "Net Premi (USD)"
           ],
           "kunci": [
            "Note",
            "Limit",
            "Limit2",
            "MDP",
            "MDP2",
            "Deductible",
            "Deductible2",
            "NetPremi",
            "NetPremi2"
           ],
           "lebar": [
            293,
            199,
            204,
            123,
            119,
            205,
            210,
            150,
            144
           ],
           "desimal": [
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "format": [
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            ""
           ],
           "syaratSel": [
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "atSel": [
            1469626,
            1473825,
            1478033,
            1482242,
            1486447,
            1490653,
            1494866,
            1499080,
            1503290
           ],
           "templatBaris": "ASM-FW-GISFW-Data-LimitSummaryList!pyGridModalTemplate"
          }
         ]
        }
       ]
      },
      {
       "t": "blok",
       "at": 1543945,
       "judul": "Total All Layers Facultative Share",
       "syarat": [],
       "anak": [
        {
         "t": "grid",
         "at": 1578012,
         "prop": "TreatyIn.TotalFacShareRnmNP",
         "dari": "akar",
         "larik": "TotalFacShareRnmNP",
         "syarat": [],
         "kolom": [
          "Total RNM Limit (RNM Share)",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          194,
          352
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          1589580,
          1594294
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "grid",
         "at": 1656531,
         "prop": "TreatyIn.TotalFacShareGrossNP",
         "dari": "akar",
         "larik": "TotalFacShareGrossNP",
         "syarat": [],
         "kolom": [
          "Total Gross Premium (MDP)",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          194,
          350
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          1668099,
          1672813
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "blok",
         "at": 1718323,
         "judul": "hiddden",
         "syarat": [],
         "anak": [
          {
           "t": "grid",
           "at": 1735032,
           "prop": "TreatyIn.TotalFacShareDeductionNP",
           "dari": "akar",
           "larik": "TotalFacShareDeductionNP",
           "syarat": [],
           "kolom": [
            "Total Deduction",
            "Value"
           ],
           "kunci": [
            "Currency",
            "Value"
           ],
           "lebar": [
            194,
            352
           ],
           "desimal": [
            null,
            2
           ],
           "format": [
            "pxNumber",
            "pxNumber"
           ],
           "syaratSel": [
            null,
            null
           ],
           "atSel": [
            1746594,
            1751308
           ],
           "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
          }
         ]
        },
        {
         "t": "grid",
         "at": 1813565,
         "prop": "TreatyIn.TotalFacShareNetNP",
         "dari": "akar",
         "larik": "TotalFacShareNetNP",
         "syarat": [],
         "kolom": [
          "Total Net Premium",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          193,
          349
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          1825123,
          1829837
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        }
       ]
      }
     ]
    }
   ]
  }
 ],
 "TreatyInShareProp": [
  {
   "t": "blok",
   "at": 18288,
   "judul": "RNM Share",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 28683,
     "judul": "",
     "syarat": [
      "TreatyIn.FacultativeShare >0"
     ],
     "anak": [
      {
       "t": "medan",
       "at": 35412,
       "label": "Share to RNM :",
       "dari": "sisi",
       "kunci": "RnmShareDeducted",
       "format": "pxNumber",
       "desimal": 2,
       "syarat": []
      },
      {
       "t": "teks",
       "at": 47860,
       "teks": "%",
       "syarat": []
      }
     ]
    },
    {
     "t": "grid",
     "at": 78804,
     "prop": "TreatyIn.Limits",
     "dari": "sisi",
     "larik": "Limits",
     "syarat": [],
     "kolom": [
      "Kind of Treaty"
     ],
     "kunci": [
      "TreatyType"
     ],
     "lebar": [
      800
     ],
     "desimal": [
      null
     ],
     "format": [
      "pxAutoComplete"
     ],
     "syaratSel": [
      null
     ],
     "atSel": [
      86754
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails"
    },
    {
     "t": "grid",
     "at": 147065,
     "prop": "TreatyIn.TotalShareRnmProp",
     "dari": "sisi",
     "larik": "TotalShareRnmProp",
     "syarat": [],
     "kolom": [
      "Total Share RNM Limit",
      "Value"
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      192,
      347
     ],
     "desimal": [
      null,
      2
     ],
     "format": [
      "pxNumber",
      "pxNumber"
     ],
     "syaratSel": [
      null,
      null
     ],
     "atSel": [
      159037,
      163993
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 217821,
     "prop": "TreatyIn.TotalSpreadedRnmProp",
     "dari": "sisi",
     "larik": "TotalSpreadedRnmProp",
     "syarat": [],
     "kolom": [
      "Total Value Spreading OR",
      "Value"
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      192,
      347
     ],
     "desimal": [
      null,
      2
     ],
     "format": [
      "pxNumber",
      "pxNumber"
     ],
     "syaratSel": [
      null,
      null
     ],
     "atSel": [
      229799,
      234755
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 288583,
     "prop": "TreatyIn.TotalSpreadedRnmRIProp",
     "dari": "sisi",
     "larik": "TotalSpreadedRnmRIProp",
     "syarat": [],
     "kolom": [
      "Total Value Spreading R/I",
      "Value"
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      192,
      347
     ],
     "desimal": [
      null,
      2
     ],
     "format": [
      "pxNumber",
      "pxNumber"
     ],
     "syaratSel": [
      null,
      null
     ],
     "atSel": [
      300568,
      305525
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    }
   ]
  }
 ],
 "TreatyInfoSubmit": [
  {
   "t": "medan",
   "at": 36657,
   "label": "Additional Information",
   "dari": "sisi",
   "kunci": "Information",
   "format": "pxTextArea",
   "desimal": null,
   "syarat": []
  },
  {
   "t": "medan",
   "at": 64606,
   "label": "Comment",
   "dari": "sisi",
   "kunci": "Comment",
   "format": "pxTextArea",
   "desimal": null,
   "syarat": [
    "TreatyIn.ViewState != 1 ||TreatyIn.RevisionState=1"
   ]
  },
  {
   "t": "blok",
   "at": 148957,
   "judul": "",
   "syarat": [
    "!TreatyMasterInEDM"
   ],
   "anak": [
    {
     "t": "tombol",
     "at": 164660,
     "label": "Submit",
     "syarat": [
      "(TreatyIn.ViewState !='1' && TreatyIn.StatusAkseptasi != 'Resolve Complete')"
     ]
    },
    {
     "t": "tombol",
     "at": 179405,
     "label": "Submit",
     "syarat": [
      "TreatyIn.RevisionState='1'"
     ]
    },
    {
     "t": "tombol",
     "at": 200534,
     "label": "Decline offer",
     "syarat": [
      "TreatyIn.ViewState != '1'"
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 223240,
   "judul": "",
   "syarat": [
    "TreatyMasterInEDM"
   ],
   "anak": [
    {
     "t": "tombol",
     "at": 229956,
     "label": "Submit",
     "syarat": [
      "TreatyIn.ViewState !='1' && TreatyIn.StatusAkseptasi != 'Resolve Complete'"
     ]
    },
    {
     "t": "tombol",
     "at": 241510,
     "label": "Decline offer",
     "syarat": [
      "TreatyIn.ViewState != '1'"
     ]
    }
   ]
  }
 ],
 "TreatyInTabsAchievement": [],
 "TreatyInTabsNonProportionalValueDifference": [
  {
   "t": "blok",
   "at": 19968,
   "judul": "Value Difference",
   "syarat": [],
   "anak": [
    {
     "t": "tombol",
     "at": 26349,
     "label": "Update Value",
     "syarat": [
      "TreatyIn.ViewState != 1"
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 47293,
   "judul": "After Pro Rate Calculation",
   "syarat": [
    "TreatyIn.IsProRate != true"
   ],
   "anak": [
    {
     "t": "include",
     "at": 55000,
     "nama": "TreatyInTabsNPValueDifferenceProRate",
     "syarat": []
    }
   ]
  },
  {
   "t": "blok",
   "at": 85258,
   "judul": "",
   "syarat": [
    "TreatyIn.IsProRate = true"
   ],
   "anak": [
    {
     "t": "blok",
     "at": 95392,
     "judul": "Before Pro Rate Calculation",
     "syarat": [],
     "anak": [
      {
       "t": "include",
       "at": 101742,
       "nama": "TreatyInTabsNPValueDifference_NoProRate",
       "syarat": []
      }
     ]
    },
    {
     "t": "blok",
     "at": 136631,
     "judul": "After Pro Rate Calculation",
     "syarat": [],
     "anak": [
      {
       "t": "include",
       "at": 142980,
       "nama": "TreatyInTabsNPValueDifferenceProRate",
       "syarat": []
      }
     ]
    }
   ]
  }
 ],
 "TreatyInTabsNPValueDifferenceProRate": [
  {
   "t": "blok",
   "at": 9755,
   "judul": "After Pro Rate",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 19866,
     "judul": "",
     "syarat": [
      "TreatyIn.IsProRate == true"
     ],
     "anak": [
      {
       "t": "teks",
       "at": 34946,
       "teks": "Pro Rate:",
       "syarat": []
      },
      {
       "t": "medan",
       "at": 43386,
       "label": "",
       "dari": "sisi",
       "kunci": "ProRateDays",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": []
      },
      {
       "t": "teks",
       "at": 53144,
       "teks": "/",
       "syarat": []
      },
      {
       "t": "medan",
       "at": 61576,
       "label": "",
       "dari": "sisi",
       "kunci": "ProRateTotalDays",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": []
      },
      {
       "t": "medan",
       "at": 81824,
       "label": "Pro Rate Percentage",
       "dari": "sisi",
       "kunci": "ProRatePercent",
       "format": "pxTextInput",
       "desimal": 4,
       "syarat": []
      },
      {
       "t": "teks",
       "at": 92027,
       "teks": "%",
       "syarat": []
      }
     ]
    },
    {
     "t": "blok",
     "at": 111289,
     "judul": "Share Difference",
     "syarat": [],
     "anak": [
      {
       "t": "medan",
       "at": 135023,
       "label": "% RNM Share",
       "dari": "sisi",
       "kunci": "ValueDifference.RNMShare",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": []
      },
      {
       "t": "medan",
       "at": 143662,
       "label": "% Brokerage",
       "dari": "sisi",
       "kunci": "ValueDifference.BrokeragePercent",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": []
      },
      {
       "t": "medan",
       "at": 179127,
       "label": "Brokerage From Other Retro",
       "dari": "sisi",
       "kunci": "FacultativeShareBrokerage",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [
        "TreatyIn.FacultativeShare >0"
       ]
      },
      {
       "t": "blok",
       "at": 450695,
       "judul": "RNM Share",
       "syarat": [],
       "anak": [
        {
         "t": "blok",
         "at": 459365,
         "judul": "",
         "syarat": [
          "TreatyIn.FacultativeShare >0"
         ],
         "anak": [
          {
           "t": "medan",
           "at": 465817,
           "label": "Share to RNM :",
           "dari": "sisi",
           "kunci": "RnmShareDeducted",
           "format": "pxTextInput",
           "desimal": null,
           "syarat": []
          },
          {
           "t": "teks",
           "at": 476472,
           "teks": "%",
           "syarat": []
          }
         ]
        },
        {
         "t": "blok",
         "at": 489362,
         "judul": "RNM Share",
         "syarat": [],
         "anak": [
          {
           "t": "grid",
           "at": 506144,
           "prop": "TreatyIn.ValueDifference.Share",
           "dari": "sisi",
           "larik": "ValueDifference.Share",
           "syarat": [],
           "kolom": [
            "",
            "",
            "",
            "",
            "",
            "",
            "100% Limit",
            "",
            "100% Limit",
            "",
            "MDP",
            "",
            "MDP"
           ],
           "kunci": [
            "LayerType",
            "Layer",
            "pyTemplateInputBox",
            "LayerPartType",
            "LayerPart",
            "RnmLimitListDisplay(1).Currency",
            "RnmLimitListDisplay(1).Value",
            "RnmLimitListDisplay(2).Currency",
            "RnmLimitListDisplay(2).Value",
            "RnmGrossPremiDisplay(1).Currency",
            "RnmGrossPremiDisplay(1).Value",
            "RnmGrossPremiDisplay(2).Currency",
            "RnmGrossPremiDisplay(2).Value"
           ],
           "lebar": [
            138,
            60,
            92,
            138,
            103,
            75,
            163,
            75,
            163,
            80,
            120,
            80,
            120
           ],
           "desimal": [
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "format": [
            "pxTextInput",
            "",
            "",
            "pxTextInput",
            "",
            "pxTextInput",
            "",
            "pxTextInput",
            "",
            "pxTextInput",
            "",
            "pxTextInput",
            ""
           ],
           "syaratSel": [
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "atSel": [
            559563,
            565255,
            569456,
            573614,
            579211,
            583418,
            589064,
            593381,
            599027,
            603344,
            608992,
            613311,
            618959
           ],
           "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridRowDetails"
          }
         ]
        },
        {
         "t": "blok",
         "at": 664555,
         "judul": "Summarry of RNM Share",
         "syarat": [],
         "anak": [
          {
           "t": "grid",
           "at": 681279,
           "prop": "TreatyIn.ValueDifference.LimitShareSummaryList",
           "dari": "sisi",
           "larik": "ValueDifference.LimitShareSummaryList",
           "syarat": [],
           "kolom": [
            "Note",
            "100% Limit (IDR)",
            "100% Limit (USD)",
            "MDP (IDR)",
            "MDP (USD)",
            "Deduction (IDR)",
            "Deduction (USD)",
            "Net Premi (IDR)",
            "Net Premi (USD)"
           ],
           "kunci": [
            "Note",
            "Limit",
            "Limit2",
            "MDP",
            "MDP2",
            "Deductible",
            "Deductible2",
            "NetPremi",
            "NetPremi2"
           ],
           "lebar": [
            293,
            199,
            204,
            123,
            119,
            205,
            210,
            150,
            144
           ],
           "desimal": [
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "format": [
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            "",
            ""
           ],
           "syaratSel": [
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "atSel": [
            719911,
            724111,
            728320,
            732530,
            736736,
            740943,
            745157,
            749372,
            753583
           ],
           "templatBaris": "ASM-FW-GISFW-Data-LimitSummaryList!pyGridModalTemplate"
          }
         ]
        }
       ]
      },
      {
       "t": "blok",
       "at": 800576,
       "judul": "Total All Layers RNM Share",
       "syarat": [],
       "anak": [
        {
         "t": "grid",
         "at": 834649,
         "prop": "TreatyIn.ValueDifference.TotalShareRnmNP",
         "dari": "sisi",
         "larik": "ValueDifference.TotalShareRnmNP",
         "syarat": [],
         "kolom": [
          "Total RNM Limit (RNM Share)",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          194,
          352
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          846235,
          850950
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "grid",
         "at": 898209,
         "prop": "TreatyIn.ValueDifference.TotalSpreadedRnmProp",
         "dari": "sisi",
         "larik": "ValueDifference.TotalSpreadedRnmProp",
         "syarat": [],
         "kolom": [
          "Total OR Limit",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          192,
          350
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          909787,
          914502
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "grid",
         "at": 961761,
         "prop": "TreatyIn.ValueDifference.TotalSpreadedRnmRIProp",
         "dari": "sisi",
         "larik": "ValueDifference.TotalSpreadedRnmRIProp",
         "syarat": [],
         "kolom": [
          "Total R/I Limit",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          192,
          350
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          973343,
          978058
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "grid",
         "at": 1040298,
         "prop": "TreatyIn.ValueDifference.TotalShareGrossNP",
         "dari": "sisi",
         "larik": "ValueDifference.TotalShareGrossNP",
         "syarat": [],
         "kolom": [
          "Total Gross Premium (MDP)",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          194,
          350
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          1051884,
          1056599
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "blok",
         "at": 1102143,
         "judul": "hiddden",
         "syarat": [],
         "anak": [
          {
           "t": "grid",
           "at": 1118861,
           "prop": "TreatyIn.ValueDifference.TotalShareDeductionNP",
           "dari": "sisi",
           "larik": "ValueDifference.TotalShareDeductionNP",
           "syarat": [],
           "kolom": [
            "Total Deduction",
            "Value"
           ],
           "kunci": [
            "Currency",
            "Value"
           ],
           "lebar": [
            194,
            352
           ],
           "desimal": [
            null,
            2
           ],
           "format": [
            "pxNumber",
            "pxNumber"
           ],
           "syaratSel": [
            null,
            null
           ],
           "atSel": [
            1130441,
            1135156
           ],
           "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
          }
         ]
        },
        {
         "t": "grid",
         "at": 1197436,
         "prop": "TreatyIn.ValueDifference.TotalShareNetNP",
         "dari": "sisi",
         "larik": "ValueDifference.TotalShareNetNP",
         "syarat": [],
         "kolom": [
          "Total Net Premium",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          193,
          349
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          1209012,
          1213727
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "grid",
         "at": 1261006,
         "prop": "TreatyIn.ValueDifference.TotalSpreadedNetPremi",
         "dari": "sisi",
         "larik": "ValueDifference.TotalSpreadedNetPremi",
         "syarat": [],
         "kolom": [
          "Total OR Net Premium",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          192,
          348
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          1272591,
          1277306
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "grid",
         "at": 1323411,
         "prop": "TreatyIn.ValueDifference.TotalSpreadedNetPremiRI",
         "dari": "sisi",
         "larik": "ValueDifference.TotalSpreadedNetPremiRI",
         "syarat": [],
         "kolom": [
          "Total R/I Net Premium",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          192,
          348
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          1334999,
          1339714
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        }
       ]
      }
     ]
    },
    {
     "t": "blok",
     "at": 1396182,
     "judul": "Installment",
     "syarat": [],
     "anak": [
      {
       "t": "medan",
       "at": 1419847,
       "label": "Installment",
       "dari": "sisi",
       "kunci": "InstallmentNo",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": []
      },
      {
       "t": "grid",
       "at": 1492192,
       "prop": "TreatyIn.ValueDifference.Installment",
       "dari": "sisi",
       "larik": "ValueDifference.Installment",
       "syarat": [],
       "kolom": [
        "Currency"
       ],
       "kunci": [
        "Currency"
       ],
       "lebar": [
        1300
       ],
       "desimal": [
        null
       ],
       "format": [
        ""
       ],
       "syaratSel": [
        null
       ],
       "atSel": [
        1499991
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInInstallment!pyGridRowDetails"
      },
      {
       "t": "grid",
       "at": 1567513,
       "prop": "TreatyIn.ValueDifference.TotalInstallmentNP",
       "dari": "sisi",
       "larik": "ValueDifference.TotalInstallmentNP",
       "syarat": [],
       "kolom": [
        "Total Installment Amount",
        "Value"
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        194,
        349
       ],
       "desimal": [
        null,
        2
       ],
       "format": [
        "pxNumber",
        "pxNumber"
       ],
       "syaratSel": [
        null,
        null
       ],
       "atSel": [
        1579236,
        1583951
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    }
   ]
  }
 ],
 "TreatyInTabsNPValueDifference_NoProRate": [
  {
   "t": "blok",
   "at": 10016,
   "judul": "After Pro Rate",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 20185,
     "judul": "Share Difference",
     "syarat": [],
     "anak": [
      {
       "t": "medan",
       "at": 43907,
       "label": "% RNM Share",
       "dari": "sisi",
       "kunci": "ValueBeforeProrate.RNMShare",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": []
      },
      {
       "t": "medan",
       "at": 52547,
       "label": "% Brokerage",
       "dari": "sisi",
       "kunci": "ValueBeforeProrate.BrokeragePercent",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": []
      },
      {
       "t": "medan",
       "at": 88007,
       "label": "Brokerage From Other Retro",
       "dari": "sisi",
       "kunci": "FacultativeShareBrokerage",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [
        "TreatyIn.FacultativeShare >0"
       ]
      },
      {
       "t": "blok",
       "at": 359430,
       "judul": "RNM Share",
       "syarat": [],
       "anak": [
        {
         "t": "blok",
         "at": 368097,
         "judul": "",
         "syarat": [
          "TreatyIn.FacultativeShare >0"
         ],
         "anak": [
          {
           "t": "medan",
           "at": 374548,
           "label": "Share to RNM :",
           "dari": "sisi",
           "kunci": "RnmShareDeducted",
           "format": "pxTextInput",
           "desimal": null,
           "syarat": []
          },
          {
           "t": "teks",
           "at": 385199,
           "teks": "%",
           "syarat": []
          }
         ]
        },
        {
         "t": "blok",
         "at": 398086,
         "judul": "RNM Share",
         "syarat": [],
         "anak": [
          {
           "t": "grid",
           "at": 414859,
           "prop": "TreatyIn.ValueBeforeProrate.Share",
           "dari": "sisi",
           "larik": "ValueBeforeProrate.Share",
           "syarat": [],
           "kolom": [
            "",
            "",
            "",
            "",
            "",
            "",
            "100% Limit",
            "",
            "100% Limit",
            "",
            "MDP",
            "",
            "MDP"
           ],
           "kunci": [
            "LayerType",
            "Layer",
            "pyTemplateInputBox",
            "LayerPartType",
            "LayerPart",
            "RnmLimitListDisplay(1).Currency",
            "RnmLimitListDisplay(1).Value",
            "RnmLimitListDisplay(2).Currency",
            "RnmLimitListDisplay(2).Value",
            "RnmGrossPremiDisplay(1).Currency",
            "RnmGrossPremiDisplay(1).Value",
            "RnmGrossPremiDisplay(2).Currency",
            "RnmGrossPremiDisplay(2).Value"
           ],
           "lebar": [
            138,
            60,
            92,
            138,
            103,
            75,
            163,
            75,
            163,
            80,
            120,
            80,
            120
           ],
           "desimal": [
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "format": [
            "pxTextInput",
            "",
            "",
            "pxTextInput",
            "",
            "pxTextInput",
            "",
            "pxTextInput",
            "",
            "pxTextInput",
            "",
            "pxTextInput",
            ""
           ],
           "syaratSel": [
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "atSel": [
            468264,
            473955,
            478155,
            482312,
            487908,
            492114,
            497759,
            502075,
            507720,
            512036,
            517683,
            522001,
            527648
           ],
           "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridRowDetails"
          }
         ]
        },
        {
         "t": "blok",
         "at": 573179,
         "judul": "Summarry of RNM Share",
         "syarat": [],
         "anak": [
          {
           "t": "grid",
           "at": 589894,
           "prop": "TreatyIn.ValueBeforeProrate.LimitShareSummaryList",
           "dari": "sisi",
           "larik": "ValueBeforeProrate.LimitShareSummaryList",
           "syarat": [],
           "kolom": [
            "Note",
            "100% Limit (IDR)",
            "100% Limit (USD)",
            "MDP (IDR)",
            "MDP (USD)",
            "Deduction (IDR)",
            "Deduction (USD)",
            "Net Premi (IDR)",
            "Net Premi (USD)"
           ],
           "kunci": [
            "Note",
            "Limit",
            "Limit2",
            "MDP",
            "MDP2",
            "Deductible",
            "Deductible2",
            "NetPremi",
            "NetPremi2"
           ],
           "lebar": [
            293,
            199,
            204,
            123,
            119,
            205,
            210,
            150,
            144
           ],
           "desimal": [
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "format": [
            "",
            "pxNumber",
            "pxNumber",
            "pxNumber",
            "pxNumber",
            "pxNumber",
            "pxNumber",
            "pxNumber",
            "pxNumber"
           ],
           "syaratSel": [
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "atSel": [
            628517,
            632716,
            637534,
            642353,
            647168,
            651984,
            656807,
            661631,
            666451
           ],
           "templatBaris": "ASM-FW-GISFW-Data-LimitSummaryList!pyGridModalTemplate"
          }
         ]
        }
       ]
      },
      {
       "t": "blok",
       "at": 714032,
       "judul": "Total All Layers RNM Share",
       "syarat": [],
       "anak": [
        {
         "t": "grid",
         "at": 748092,
         "prop": "TreatyIn.ValueBeforeProrate.TotalShareRnmNP",
         "dari": "sisi",
         "larik": "ValueBeforeProrate.TotalShareRnmNP",
         "syarat": [],
         "kolom": [
          "Total RNM Limit (RNM Share)",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          194,
          352
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          759676,
          764390
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "grid",
         "at": 811628,
         "prop": "TreatyIn.ValueBeforeProrate.TotalSpreadedRnmProp",
         "dari": "sisi",
         "larik": "ValueBeforeProrate.TotalSpreadedRnmProp",
         "syarat": [],
         "kolom": [
          "Total OR Limit",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          192,
          350
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          823204,
          827918
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "grid",
         "at": 875156,
         "prop": "TreatyIn.ValueBeforeProrate.TotalSpreadedRnmRIProp",
         "dari": "sisi",
         "larik": "ValueBeforeProrate.TotalSpreadedRnmRIProp",
         "syarat": [],
         "kolom": [
          "Total R/I Limit",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          192,
          350
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          886736,
          891450
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "grid",
         "at": 953667,
         "prop": "TreatyIn.ValueBeforeProrate.TotalShareGrossNP",
         "dari": "sisi",
         "larik": "ValueBeforeProrate.TotalShareGrossNP",
         "syarat": [],
         "kolom": [
          "Total Gross Premium (MDP)",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          194,
          350
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          965251,
          969965
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "blok",
         "at": 1015495,
         "judul": "hiddden",
         "syarat": [],
         "anak": [
          {
           "t": "grid",
           "at": 1032204,
           "prop": "TreatyIn.ValueBeforeProrate.TotalShareDeductionNP",
           "dari": "sisi",
           "larik": "ValueBeforeProrate.TotalShareDeductionNP",
           "syarat": [],
           "kolom": [
            "Total Deduction",
            "Value"
           ],
           "kunci": [
            "Currency",
            "Value"
           ],
           "lebar": [
            194,
            352
           ],
           "desimal": [
            null,
            2
           ],
           "format": [
            "pxNumber",
            "pxNumber"
           ],
           "syaratSel": [
            null,
            null
           ],
           "atSel": [
            1043782,
            1048496
           ],
           "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
          }
         ]
        },
        {
         "t": "grid",
         "at": 1110733,
         "prop": "TreatyIn.ValueBeforeProrate.TotalShareNetNP",
         "dari": "sisi",
         "larik": "ValueBeforeProrate.TotalShareNetNP",
         "syarat": [],
         "kolom": [
          "Total Net Premium",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          193,
          349
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          1122307,
          1127021
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "grid",
         "at": 1174259,
         "prop": "TreatyIn.ValueBeforeProrate.TotalSpreadedNetPremi",
         "dari": "sisi",
         "larik": "ValueBeforeProrate.TotalSpreadedNetPremi",
         "syarat": [],
         "kolom": [
          "Total OR Net Premium",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          192,
          348
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          1185842,
          1190556
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "grid",
         "at": 1236640,
         "prop": "TreatyIn.ValueBeforeProrate.TotalSpreadedNetPremiRI",
         "dari": "sisi",
         "larik": "ValueBeforeProrate.TotalSpreadedNetPremiRI",
         "syarat": [],
         "kolom": [
          "Total R/I Net Premium",
          "Value"
         ],
         "kunci": [
          "Currency",
          "Value"
         ],
         "lebar": [
          192,
          348
         ],
         "desimal": [
          null,
          2
         ],
         "format": [
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          1248226,
          1252940
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        }
       ]
      }
     ]
    },
    {
     "t": "blok",
     "at": 1309391,
     "judul": "Installment",
     "syarat": [],
     "anak": [
      {
       "t": "medan",
       "at": 1333051,
       "label": "Installment",
       "dari": "sisi",
       "kunci": "InstallmentNo",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": []
      },
      {
       "t": "grid",
       "at": 1405380,
       "prop": "TreatyIn.ValueBeforeProrate.Installment",
       "dari": "sisi",
       "larik": "ValueBeforeProrate.Installment",
       "syarat": [],
       "kolom": [
        "Currency"
       ],
       "kunci": [
        "Currency"
       ],
       "lebar": [
        1300
       ],
       "desimal": [
        null
       ],
       "format": [
        ""
       ],
       "syaratSel": [
        null
       ],
       "atSel": [
        1413178
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInInstallment!pyGridRowDetails"
      },
      {
       "t": "grid",
       "at": 1480651,
       "prop": "TreatyIn.ValueBeforeProrate.TotalInstallmentNP",
       "dari": "sisi",
       "larik": "ValueBeforeProrate.TotalInstallmentNP",
       "syarat": [],
       "kolom": [
        "Total Installment Amount",
        "Value"
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        194,
        349
       ],
       "desimal": [
        null,
        2
       ],
       "format": [
        "pxNumber",
        "pxNumber"
       ],
       "syaratSel": [
        null,
        null
       ],
       "atSel": [
        1492372,
        1497086
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    }
   ]
  }
 ]
}

export const GRID_KURS: Readonly<Record<'lama' | 'baru', GridKerangka>> = {
 "lama": {
  "t": "grid",
  "at": 173882,
  "prop": "TreatyIn.OLDDATA.CurrencyList",
  "dari": "sisi",
  "larik": "CurrencyList",
  "syarat": [],
  "kolom": [
   "Currency",
   "Value to IDR",
   "Valid From",
   "Valid Until"
  ],
  "kunci": [
   "Currency",
   "Conversion",
   "PeriodStart",
   "PeriodEnd"
  ],
  "lebar": [
   188,
   208,
   193,
   198
  ],
  "desimal": [
   null,
   null,
   null,
   null
  ],
  "format": [
   "pxAutoComplete",
   "pxTextInput",
   "pxDateTime",
   "pxDateTime"
  ],
  "syaratSel": [
   null,
   null,
   null,
   null
  ],
  "atSel": [
   195570,
   206033,
   212117,
   218258
  ],
  "templatBaris": "ASM-FW-GISFW-Data-TreatyInCurrencyList!pyGridModalTemplate"
 },
 "baru": {
  "t": "grid",
  "at": 407552,
  "prop": "TreatyIn.CurrencyList",
  "dari": "sisi",
  "larik": "CurrencyList",
  "syarat": [],
  "kolom": [
   "Currency",
   "Value to IDR",
   "Valid From",
   "Valid Until"
  ],
  "kunci": [
   "CurrencyID",
   "Conversion",
   "PeriodStart",
   "PeriodEnd"
  ],
  "lebar": [
   191,
   208,
   194,
   198
  ],
  "desimal": [
   null,
   null,
   null,
   null
  ],
  "format": [
   "pxDropdown",
   "pxTextInput",
   "pxDateTime",
   "pxDateTime"
  ],
  "syaratSel": [
   null,
   null,
   null,
   null
  ],
  "atSel": [
   439485,
   454068,
   460497,
   466983
  ],
  "templatBaris": "ASM-FW-GISFW-Data-TreatyInCurrencyList!pyGridModalTemplate"
 }
}

export const DIBUANG: readonly Terbuang[] = [
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 84489,
  "jenis": "kolom",
  "alasan": "kolom tombol (jalur tulis)"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 192646,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 238070,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 298880,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 303702,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 356161,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 360983,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 412848,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 417672,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 591940,
  "jenis": "kolom",
  "alasan": "kolom tombol (jalur tulis)"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 769045,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 775304,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 812396,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 818655,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 824914,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 1019496,
  "jenis": "kolom",
  "alasan": "kolom tombol (jalur tulis)"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 1201663,
  "jenis": "blok",
  "alasan": "penjaga mati: Never"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 1555175,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 1610389,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 1616648,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 1622907,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 1923992,
  "jenis": "kolom",
  "alasan": "kolom tombol (jalur tulis)"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 2020871,
  "jenis": "kolom",
  "alasan": "kolom tombol (jalur tulis)"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 2067996,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 2113110,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 2279656,
  "jenis": "kolom",
  "alasan": "kolom tombol (jalur tulis)"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 3151407,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 3220454,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 3227224,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 3245214,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 3324438,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 3611911,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 3818422,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportional",
  "at": 3862701,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldData",
  "at": 181738,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldData",
  "at": 468315,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldData",
  "at": 474587,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldData",
  "at": 512108,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldData",
  "at": 518380,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldData",
  "at": 524652,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldData",
  "at": 841977,
  "jenis": "blok",
  "alasan": "penjaga mati: Never"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldData",
  "at": 1197494,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldData",
  "at": 1259228,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldData",
  "at": 1528971,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldData",
  "at": 1736201,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsProportional",
  "at": 50128,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsProportional",
  "at": 264437,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsProportional",
  "at": 483459,
  "jenis": "kolom",
  "alasan": "kolom tombol (jalur tulis)"
 },
 {
  "berkas": "TreatyInTabsProportional",
  "at": 541585,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsProportional",
  "at": 592801,
  "jenis": "kolom",
  "alasan": "kolom tombol (jalur tulis)"
 },
 {
  "berkas": "TreatyInTabsProportional",
  "at": 624574,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER && TreatyMasterInEDM"
 },
 {
  "berkas": "TreatyInTabsProportional",
  "at": 738460,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsProportional",
  "at": 769543,
  "jenis": "sel",
  "alasan": "penjaga mati: never"
 },
 {
  "berkas": "TreatyInTabsProportional",
  "at": 941807,
  "jenis": "kolom",
  "alasan": "kolom tombol (jalur tulis)"
 },
 {
  "berkas": "TreatyInTabsProportional",
  "at": 1094235,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsProportional",
  "at": 1254367,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsProportional",
  "at": 1285075,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsProportional",
  "at": 1299072,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsProportional",
  "at": 1404482,
  "jenis": "kolom",
  "alasan": "kolom tombol (jalur tulis)"
 },
 {
  "berkas": "TreatyInTabsProportionalOldData",
  "at": 46718,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsProportionalOldData",
  "at": 473467,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsProportionalOldData",
  "at": 759233,
  "jenis": "sel",
  "alasan": "penjaga mati: 3=4"
 },
 {
  "berkas": "TreatyInTabsProportionalOldData",
  "at": 830031,
  "jenis": "sel",
  "alasan": "penjaga mati: 3=4"
 },
 {
  "berkas": "TreatyInTabsProportionalOldData",
  "at": 942729,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsProportionalOldData",
  "at": 973311,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsProportionalOldData",
  "at": 987206,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldDataShare",
  "at": 186861,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldDataShare",
  "at": 229982,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldDataShare",
  "at": 1126992,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldDataShare",
  "at": 1142424,
  "jenis": "sel",
  "alasan": "penjaga mati: TreatyIn.ViewState !='1' && NEVER"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldDataShare",
  "at": 1155984,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldDataShare",
  "at": 1162580,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldDataShare",
  "at": 1373067,
  "jenis": "kolom",
  "alasan": "kolom tombol (jalur tulis)"
 },
 {
  "berkas": "TreatyInTabsNonProportionalOldDataShare",
  "at": 1864405,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInShareProp",
  "at": 41717,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "TreatyInShareProp",
  "at": 192509,
  "jenis": "sel",
  "alasan": "penjaga mati: 3=4"
 },
 {
  "berkas": "TreatyInShareProp",
  "at": 263271,
  "jenis": "sel",
  "alasan": "penjaga mati: 3=4"
 },
 {
  "berkas": "TreatyInfoSubmit",
  "at": 50147,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInfoSubmit",
  "at": 81615,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInfoSubmit",
  "at": 86779,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInfoSubmit",
  "at": 94258,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInfoSubmit",
  "at": 134677,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInfoSubmit",
  "at": 209160,
  "jenis": "sel",
  "alasan": "syarat identitas operator (blok dev)"
 },
 {
  "berkas": "TreatyInfoSubmit",
  "at": 253707,
  "jenis": "sel",
  "alasan": "syarat identitas operator (blok dev)"
 },
 {
  "berkas": "TreatyInTabsNPValueDifferenceProRate",
  "at": 39320,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "TreatyInTabsNPValueDifferenceProRate",
  "at": 48653,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "TreatyInTabsNPValueDifferenceProRate",
  "at": 57510,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "TreatyInTabsNPValueDifferenceProRate",
  "at": 87536,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "TreatyInTabsNPValueDifferenceProRate",
  "at": 170207,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNPValueDifferenceProRate",
  "at": 200969,
  "jenis": "sel",
  "alasan": "penjaga mati: FALSE && TreatyIn.ViewState !='1'"
 },
 {
  "berkas": "TreatyInTabsNPValueDifferenceProRate",
  "at": 219857,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNPValueDifferenceProRate",
  "at": 428436,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNPValueDifferenceProRate",
  "at": 471251,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "TreatyInTabsNPValueDifferenceProRate",
  "at": 623278,
  "jenis": "kolom",
  "alasan": "kolom tombol (jalur tulis)"
 },
 {
  "berkas": "TreatyInTabsNPValueDifferenceProRate",
  "at": 1374273,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNPValueDifferenceProRate",
  "at": 1452727,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNPValueDifferenceProRate",
  "at": 1542403,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNPValueDifferenceProRate",
  "at": 1612219,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNPValueDifference_NoProRate",
  "at": 79089,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNPValueDifference_NoProRate",
  "at": 109847,
  "jenis": "sel",
  "alasan": "penjaga mati: FALSE && TreatyIn.ViewState !='1'"
 },
 {
  "berkas": "TreatyInTabsNPValueDifference_NoProRate",
  "at": 128731,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNPValueDifference_NoProRate",
  "at": 337178,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNPValueDifference_NoProRate",
  "at": 379980,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "TreatyInTabsNPValueDifference_NoProRate",
  "at": 531966,
  "jenis": "kolom",
  "alasan": "kolom tombol (jalur tulis)"
 },
 {
  "berkas": "TreatyInTabsNPValueDifference_NoProRate",
  "at": 1287488,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNPValueDifference_NoProRate",
  "at": 1365928,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNPValueDifference_NoProRate",
  "at": 1455552,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNPValueDifference_NoProRate",
  "at": 1525343,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInNONProportional",
  "at": 473425,
  "jenis": "kolom",
  "alasan": "kolom tombol (jalur tulis)"
 }
]
