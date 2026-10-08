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
       "Amount",
       ""
      ],
      "kunci": [
       "TreatyGroup",
       "Currency",
       "Amount",
       ""
      ],
      "lebar": [
       245,
       127,
       379,
       153
      ],
      "desimal": [
       null,
       null,
       null,
       null
      ],
      "format": [
       "pxAutoComplete",
       "",
       "pxNumber",
       "pxButton"
      ],
      "syaratSel": [
       null,
       null,
       null,
       "TreatyIn.IsEditData !='1'"
      ],
      "atSel": [
       66870,
       74649,
       79118,
       84489
      ],
      "baca": [
       null,
       null,
       null,
       "selalu"
      ],
      "tombol": [
       null,
       null,
       null,
       {
        "t": "tombol",
        "at": 84489,
        "label": "Delete",
        "syarat": [
         "TreatyIn.IsEditData !='1'"
        ],
        "aksi": [
         {
          "aksi": "deleteRow"
         }
        ],
        "nonaktif": [
         "TreatyIn.EDMMaterialType = 2"
        ]
       }
      ],
      "tombolKepala": [
       null,
       null,
       null,
       {
        "t": "tombol",
        "at": 56897,
        "label": "Add",
        "syarat": [
         "TreatyIn.IsEditData !='1'"
        ],
        "aksi": [
         {
          "aksi": "refresh",
          "aktivitas": "TreatyInNonAddItem",
          "param": {
           "Type": "retention"
          }
         }
        ],
        "ikon": "rpadd.gif",
        "nonaktif": [
         "TreatyIn.EDMMaterialType = 2"
        ]
       }
      ],
      "pilihan": [
       {
        "sumber": "reportdefinition",
        "rd": "BrowseTreatyGroup_RD",
        "nilai": "TreatyGroupName"
       },
       null,
       null,
       null
      ],
      "aksiUbah": [
       null,
       null,
       null,
       null
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInRetention!pyGridRowDetails",
      "rincian": "MaxRetention"
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
    "baca": [
     null,
     null
    ],
    "tombol": [
     null,
     null
    ],
    "tombolKepala": [
     null,
     null
    ],
    "pilihan": [
     null,
     null
    ],
    "aksiUbah": [
     null,
     null
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
   },
   {
    "t": "tombol",
    "at": 198918,
    "label": "Update Total",
    "syarat": [
     "TreatyIn.IsEditData !='1'"
    ],
    "aksi": [
     {
      "aksi": "refresh",
      "aktivitas": "TreatyInNPSetTotal",
      "param": {
       "type": "retention"
      }
     }
    ],
    "nonaktif": [
     "TreatyIn.EDMMaterialType = 2"
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
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState = 1",
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ],
      "pilihan": {
       "sumber": "reportdefinition",
       "rd": "BrowseCurrencyTreatyIn_RD",
       "nilai": "Currency"
      },
      "aksiUbah": [
       {
        "aksi": "postValue"
       }
      ]
     },
     {
      "t": "medan",
      "at": 286560,
      "label": "",
      "dari": "sisi",
      "kunci": "RSMDLimit",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState = 1",
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ]
     },
     {
      "t": "medan",
      "at": 317490,
      "label": "Earthquake Limit",
      "dari": "sisi",
      "kunci": "CurrencyEarthquake",
      "format": "pxAutoComplete",
      "desimal": null,
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ],
      "pilihan": {
       "sumber": "reportdefinition",
       "rd": "BrowseCurrencyTreatyIn_RD",
       "nilai": "Currency"
      },
      "aksiUbah": [
       {
        "aksi": "postValue"
       }
      ]
     },
     {
      "t": "medan",
      "at": 343882,
      "label": "",
      "dari": "sisi",
      "kunci": "Earthquake",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ]
     },
     {
      "t": "medan",
      "at": 374771,
      "label": "Flood Limit (Jabodetabek)",
      "dari": "sisi",
      "kunci": "CurrencyFloodJab",
      "format": "pxAutoComplete",
      "desimal": null,
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ],
      "pilihan": {
       "sumber": "reportdefinition",
       "rd": "BrowseCurrencyTreatyIn_RD",
       "nilai": "Currency"
      },
      "aksiUbah": [
       {
        "aksi": "postValue"
       }
      ]
     },
     {
      "t": "medan",
      "at": 400573,
      "label": "",
      "dari": "sisi",
      "kunci": "FloodJab",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ]
     },
     {
      "t": "medan",
      "at": 431464,
      "label": "Flood Limit (Nationwide)",
      "dari": "sisi",
      "kunci": "CurrencyFloodNat",
      "format": "pxAutoComplete",
      "desimal": null,
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ],
      "pilihan": {
       "sumber": "reportdefinition",
       "rd": "BrowseCurrencyTreatyIn_RD",
       "nilai": "Currency"
      },
      "aksiUbah": [
       {
        "aksi": "postValue"
       }
      ]
     },
     {
      "t": "medan",
      "at": 457287,
      "label": "",
      "dari": "sisi",
      "kunci": "FloodNation",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ]
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
       "Amount in IDR",
       ""
      ],
      "kunci": [
       "TreatyGroup",
       "AsDate",
       "Proportion",
       "Currency",
       "Amount",
       "AmountIDR",
       ""
      ],
      "lebar": [
       287,
       197,
       142,
       137,
       222,
       223,
       125
      ],
      "desimal": [
       null,
       null,
       2,
       null,
       2,
       null,
       null
      ],
      "format": [
       "pxAutoComplete",
       "",
       "pxNumber",
       "",
       "pxNumber",
       "pxNumber",
       "pxButton"
      ],
      "syaratSel": [
       null,
       null,
       null,
       null,
       null,
       null,
       "TreatyIn.ViewState !='1'"
      ],
      "atSel": [
       562421,
       567864,
       572297,
       577381,
       581815,
       586894,
       591940
      ],
      "baca": [
       [
        "1=1"
       ],
       null,
       null,
       null,
       null,
       null,
       "selalu"
      ],
      "tombol": [
       null,
       null,
       null,
       null,
       null,
       null,
       {
        "t": "tombol",
        "at": 591940,
        "label": "Delete",
        "syarat": [
         "TreatyIn.ViewState !='1'"
        ],
        "aksi": [
         {
          "aksi": "deleteRow"
         }
        ],
        "nonaktif": [
         "TreatyIn.EDMMaterialType = 2"
        ]
       }
      ],
      "tombolKepala": [
       null,
       null,
       null,
       null,
       null,
       null,
       {
        "t": "tombol",
        "at": 552737,
        "label": "Add",
        "syarat": [
         "TreatyIn.ViewState !='1'"
        ],
        "aksi": [
         {
          "aksi": "refresh",
          "aktivitas": "TreatyInNonAddItem",
          "param": {
           "Type": "egnpi"
          }
         }
        ],
        "ikon": "rpadd.gif",
        "nonaktif": [
         "TreatyIn.EDMMaterialType = 2"
        ]
       }
      ],
      "pilihan": [
       {
        "sumber": "associated"
       },
       null,
       null,
       null,
       null,
       null,
       null
      ],
      "aksiUbah": [
       null,
       null,
       null,
       null,
       null,
       null,
       null
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInEGNPI!pyGridRowDetails",
      "rincian": "DetailEGNPI"
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
    "baca": [
     null,
     null
    ],
    "tombol": [
     null,
     null
    ],
    "tombolKepala": [
     null,
     null
    ],
    "pilihan": [
     null,
     null
    ],
    "aksiUbah": [
     null,
     null
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
    "syarat": [],
    "baca": "selalu"
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
    "syarat": [],
    "baca": "selalu"
   },
   {
    "t": "tombol",
    "at": 840132,
    "label": "Update Total",
    "syarat": [
     "TreatyIn.ViewState !='1'"
    ],
    "aksi": [
     {
      "aksi": "refresh",
      "aktivitas": "TreatyInNPSetTotal",
      "param": {
       "type": "egnpi"
      }
     }
    ],
    "nonaktif": [
     "TreatyIn.EDMMaterialType = 2"
    ]
   },
   {
    "t": "tombol",
    "at": 849212,
    "label": "Update EGNPI Value",
    "syarat": [
     "TreatyIn.ViewState !='1' && TreatyMasterInEDM"
    ],
    "aksi": [
     {
      "aksi": "refresh",
      "aktivitas": "TreatyInEGNPIListValue"
     }
    ],
    "nonaktif": [
     "TreatyIn.EDMMaterialType = 2"
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
     "Deductible ( USD )",
     ""
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
     "Deductible2",
     ""
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
     187,
     108
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
     2,
     null
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
     "pxNumber",
     "pxButton"
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
     "TreatyIn.ViewState !='1'"
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
     1013782,
     1019496
    ],
    "baca": [
     null,
     null,
     "selalu",
     null,
     null,
     null,
     null,
     null,
     null,
     "selalu"
    ],
    "tombol": [
     null,
     null,
     null,
     null,
     null,
     null,
     null,
     null,
     null,
     {
      "t": "tombol",
      "at": 1019496,
      "label": "Delete",
      "syarat": [
       "TreatyIn.ViewState !='1'"
      ],
      "aksi": [
       {
        "aksi": "deleteRow"
       }
      ],
      "nonaktif": [
       "TreatyIn.EDMMaterialType = 2"
      ]
     }
    ],
    "tombolKepala": [
     null,
     null,
     null,
     null,
     null,
     null,
     null,
     null,
     null,
     {
      "t": "tombol",
      "at": 961609,
      "label": "add Layer",
      "syarat": [
       "TreatyIn.ViewState !='1'"
      ],
      "aksi": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInNonAddItem",
        "param": {
         "Type": "limits"
        }
       }
      ],
      "nonaktif": [
       "TreatyIn.EDMMaterialType = 2"
      ]
     }
    ],
    "pilihan": [
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
    "aksiUbah": [
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
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails",
    "rincian": "Layers"
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
      "baca": [
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
      "tombol": [
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
      "tombolKepala": [
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
      "pilihan": [
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
      "aksiUbah": [
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
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
      "syarat": [],
      "baca": "selalu"
     },
     {
      "t": "tombol",
      "at": 1638126,
      "label": "Update Total",
      "syarat": [
       "TreatyIn.ViewState !='1'"
      ],
      "aksi": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInNPSetTotal",
        "param": {
         "type": "limits"
        }
       },
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInSummaryMDP"
       },
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInSummaryLimit"
       }
      ],
      "nonaktif": [
       "TreatyIn.EDMMaterialType = 2"
      ]
     },
     {
      "t": "tombol",
      "at": 1652282,
      "label": "Update Value in List",
      "syarat": [
       "TreatyIn.ViewState !='1' && TreatyMasterInEDM"
      ],
      "aksi": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInLimitsListValue"
       }
      ],
      "nonaktif": [
       "TreatyIn.EDMMaterialType = 2"
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
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ],
      "aksiUbah": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInXOLAddSpreading"
       }
      ]
     },
     {
      "t": "medan",
      "at": 1750505,
      "label": "% Brokerage",
      "dari": "sisi",
      "kunci": "BrokeragePercent",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ],
      "aksiUbah": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInSetBrokerage"
       },
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInXOLAddSpreading"
       }
      ]
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
      "caption": "Share Across The Board",
      "baca": [
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ],
      "aksiUbah": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInXOLAddSpreading",
        "syarat": "TreatyIn.RNMShareAcrossTheBoard = 'true'"
       },
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInSetBrokerage",
        "syarat": "TreatyIn.BrokeragePercent > 0"
       }
      ]
     },
     {
      "t": "medan",
      "at": 1797115,
      "label": "Share to Other Retro",
      "dari": "sisi",
      "kunci": "FacultativeShare",
      "format": "pxTextInput",
      "desimal": null,
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ],
      "aksiUbah": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInXOLAddSpreading"
       }
      ]
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
      ],
      "baca": [
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ],
      "aksiUbah": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInXOLAddSpreading"
       }
      ]
     },
     {
      "t": "tombol",
      "at": 1828787,
      "label": "Update Summary",
      "syarat": [
       "TreatyIn.ViewState !='1'"
      ],
      "aksi": [
       {
        "aksi": "setValue",
        "param": {
         "TreatyIn.FacultativeShare": "0"
        },
        "syarat": "TreatyIn.FacultativeShare = ''"
       },
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInNonAddItem",
        "param": {
         "Type": "share"
        }
       },
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInXOLAddSpreading"
       },
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInSetBrokerage"
       },
       {
        "aksi": "refresh",
        "aktivitas": "TreatyEDMCalculateDifference",
        "syarat": "TreatyIn.EDMMaterialType = 7897987 && TreatyIn.EDMState != 3"
       }
      ],
      "nonaktif": [
       "TreatyIn.EDMMaterialType = 2"
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
       "% Share",
       ""
      ],
      "kunci": [
       "ReinsName",
       "Layer",
       "SharePct",
       ""
      ],
      "lebar": [
       424,
       180,
       188,
       159
      ],
      "desimal": [
       null,
       null,
       2,
       null
      ],
      "format": [
       "pxAutoComplete",
       "pxTextInput",
       "pxNumber",
       "pxButton"
      ],
      "syaratSel": [
       null,
       null,
       null,
       "TreatyIn.ViewState !='1'"
      ],
      "atSel": [
       1904195,
       1912259,
       1918452,
       1923992
      ],
      "baca": [
       [
        "TreatyIn.ViewState = 1"
       ],
       [
        "TreatyIn.ViewState = 1"
       ],
       [
        "TreatyIn.ViewState = 1"
       ],
       "selalu"
      ],
      "tombol": [
       null,
       null,
       null,
       {
        "t": "tombol",
        "at": 1923992,
        "label": "Delete",
        "syarat": [
         "TreatyIn.ViewState !='1'"
        ],
        "aksi": [
         {
          "aksi": "deleteRow"
         }
        ]
       }
      ],
      "tombolKepala": [
       null,
       null,
       null,
       {
        "t": "tombol",
        "at": 1894484,
        "label": "Add",
        "syarat": [
         "TreatyIn.ViewState !='1'"
        ],
        "aksi": [
         {
          "aksi": "refresh",
          "aktivitas": "TreatyInNonAddItem",
          "param": {
           "Type": "sharereins"
          }
         }
        ],
        "ikon": "rpadd.gif"
       }
      ],
      "pilihan": [
       {
        "sumber": "reportdefinition",
        "rd": "BrowseAgentNusaRe_RD",
        "nilai": "ClientName"
       },
       null,
       null,
       null
      ],
      "aksiUbah": [
       null,
       null,
       null,
       null
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
         "% Share",
         ""
        ],
        "kunci": [
         "ReinsName",
         "Layer",
         "SharePct",
         ""
        ],
        "lebar": [
         428,
         182,
         189,
         160
        ],
        "desimal": [
         null,
         null,
         2,
         null
        ],
        "format": [
         "pxAutoComplete",
         "pxTextInput",
         "pxNumber",
         "pxButton"
        ],
        "syaratSel": [
         null,
         null,
         null,
         "TreatyIn.ViewState !='1'"
        ],
        "atSel": [
         2001074,
         2009138,
         2015331,
         2020871
        ],
        "baca": [
         [
          "TreatyIn.ViewState = 1"
         ],
         [
          "TreatyIn.ViewState = 1"
         ],
         [
          "TreatyIn.ViewState = 1"
         ],
         "selalu"
        ],
        "tombol": [
         null,
         null,
         null,
         {
          "t": "tombol",
          "at": 2020871,
          "label": "Delete",
          "syarat": [
           "TreatyIn.ViewState !='1'"
          ],
          "aksi": [
           {
            "aksi": "deleteRow"
           }
          ]
         }
        ],
        "tombolKepala": [
         null,
         null,
         null,
         {
          "t": "tombol",
          "at": 1991359,
          "label": "Add",
          "syarat": [
           "TreatyIn.ViewState !='1'"
          ],
          "aksi": [
           {
            "aksi": "refresh",
            "aktivitas": "TreatyInNonAddItem",
            "param": {
             "Type": "sharefacname"
            }
           }
          ],
          "ikon": "rpadd.gif"
         }
        ],
        "pilihan": [
         {
          "sumber": "reportdefinition",
          "rd": "BrowseAgentNusaRe_RD",
          "nilai": "ClientName"
         },
         null,
         null,
         null
        ],
        "aksiUbah": [
         null,
         null,
         null,
         null
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
          "syarat": [],
          "baca": "selalu"
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
        "judul": "",
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
          "baca": [
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
           "selalu"
          ],
          "tombol": [
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
          "tombolKepala": [
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
          "pilihan": [
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
          "aksiUbah": [
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
          "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridRowDetails",
          "rincian": "Share"
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
          "baca": [
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
          "tombol": [
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
          "tombolKepala": [
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
          "pilihan": [
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
          "aksiUbah": [
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "blok",
      "at": 2873031,
      "judul": "",
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
        "baca": [
         null,
         null
        ],
        "tombol": [
         null,
         null
        ],
        "tombolKepala": [
         null,
         null
        ],
        "pilihan": [
         null,
         null
        ],
        "aksiUbah": [
         null,
         null
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
     },
     {
      "t": "tombol",
      "at": 3184690,
      "label": "Update Total",
      "syarat": [
       "TreatyIn.ViewState !='1'"
      ],
      "aksi": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInNPSetTotal",
        "param": {
         "type": "share"
        }
       },
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInSummaryLimitShare"
       },
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInSummaryLimitFacShare"
       }
      ],
      "nonaktif": [
       "TreatyIn.EDMMaterialType = 2"
      ]
     },
     {
      "t": "tombol",
      "at": 3198699,
      "label": "Update Value in Share",
      "syarat": [
       "TreatyIn.ViewState !='1' && TreatyMasterInEDM"
      ],
      "aksi": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInShareListValue"
       }
      ],
      "nonaktif": [
       "TreatyIn.EDMMaterialType = 2"
      ]
     },
     {
      "t": "tombol",
      "at": 3236226,
      "label": "Hide Facultative Share (unused)",
      "syarat": [
       "TreatyIn.FacultativeShare>0 && facsharedisp.CARI1 = '1'"
      ],
      "aksi": [
       {
        "aksi": "refresh",
        "transformasi": "TreatyInShowHideFacShare",
        "paramDT": {
         "type": "0"
        }
       }
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
    "caption": "Has Share To Retro:",
    "baca": [
     "TreatyIn.ProportionType = \"Proportional\"",
     "TreatyIn.ViewState = 1"
    ],
    "aksiUbah": [
     {
      "aksi": "refresh"
     }
    ]
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
    "syarat": [],
    "baca": [
     "TreatyIn.ViewState =='1' || TreatyIn.EDMMaterialType = 2"
    ],
    "aksiUbah": [
     {
      "aksi": "postValue"
     },
     {
      "aksi": "refresh",
      "aktivitas": "TreatyInSetValueInstallment",
      "param": {
       "Installment": "TreatyIn.InstallmentNo"
      }
     }
    ]
   },
   {
    "t": "tombol",
    "at": 3516400,
    "label": "Update Value",
    "syarat": [
     "TreatyIn.ViewState !='1'"
    ],
    "aksi": [
     {
      "aksi": "refresh",
      "aktivitas": "TreatyInSetValueInstallment",
      "param": {
       "Installment": "TreatyIn.InstallmentNo",
       "status": "update"
      }
     }
    ],
    "nonaktif": [
     "TreatyIn.EDMMaterialType = 2"
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
    "baca": [
     null
    ],
    "tombol": [
     null
    ],
    "tombolKepala": [
     null
    ],
    "pilihan": [
     null
    ],
    "aksiUbah": [
     null
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInInstallment!pyGridRowDetails",
    "rincian": "Installments"
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
    "baca": [
     null,
     null
    ],
    "tombol": [
     null,
     null
    ],
    "tombolKepala": [
     null,
     null
    ],
    "pilihan": [
     null,
     null
    ],
    "aksiUbah": [
     null,
     null
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
   },
   {
    "t": "tombol",
    "at": 3824681,
    "label": "Update Total",
    "syarat": [
     "TreatyIn.ViewState !='1'"
    ],
    "aksi": [
     {
      "aksi": "refresh",
      "aktivitas": "TreatyInNPSetTotal",
      "param": {
       "type": "installment"
      }
     }
    ],
    "nonaktif": [
     "TreatyIn.EDMMaterialType = 2"
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
      "syarat": [],
      "baca": [
       "TreatyIn.IsEditData= 1 || TreatyIn.EDMMaterialType = 1"
      ]
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
      "syarat": [],
      "baca": [
       "TreatyIn.IsEditData= 1 || TreatyIn.EDMMaterialType = 1"
      ]
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
      "baca": [
       null,
       null,
       null
      ],
      "tombol": [
       null,
       null,
       null
      ],
      "tombolKepala": [
       null,
       null,
       null
      ],
      "pilihan": [
       {
        "sumber": "reportdefinition",
        "rd": "BrowseTreatyGroup_RD",
        "nilai": "TreatyGroupName"
       },
       null,
       null
      ],
      "aksiUbah": [
       null,
       null,
       null
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInRetention!pyGridRowDetails",
      "rincian": "MaxRetentionOldData"
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
    "baca": [
     null,
     null
    ],
    "tombol": [
     null,
     null
    ],
    "tombolKepala": [
     null,
     null
    ],
    "pilihan": [
     null,
     null
    ],
    "aksiUbah": [
     null,
     null
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
      "baca": [
       [
        "1=1"
       ],
       null,
       null,
       null,
       null,
       null
      ],
      "tombol": [
       null,
       null,
       null,
       null,
       null,
       null
      ],
      "tombolKepala": [
       null,
       null,
       null,
       null,
       null,
       null
      ],
      "pilihan": [
       {
        "sumber": "associated"
       },
       null,
       null,
       null,
       null,
       null
      ],
      "aksiUbah": [
       null,
       null,
       null,
       null,
       null,
       null
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInEGNPI!pyGridRowDetails",
      "rincian": "DetailEGNPIOldData"
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
    "baca": [
     null,
     null
    ],
    "tombol": [
     null,
     null
    ],
    "tombolKepala": [
     null,
     null
    ],
    "pilihan": [
     null,
     null
    ],
    "aksiUbah": [
     null,
     null
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
    "syarat": [],
    "baca": "selalu"
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
    "syarat": [],
    "baca": "selalu"
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
    "baca": [
     null,
     null,
     "selalu",
     null,
     null,
     null,
     null,
     null
    ],
    "tombol": [
     null,
     null,
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "tombolKepala": [
     null,
     null,
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "pilihan": [
     null,
     null,
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "aksiUbah": [
     null,
     null,
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails",
    "rincian": "LayersOldData"
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
      "baca": [
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
      "tombol": [
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
      "tombolKepala": [
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
      "pilihan": [
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
      "aksiUbah": [
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
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
      "baca": [
       null,
       null
      ],
      "tombol": [
       null,
       null
      ],
      "tombolKepala": [
       null,
       null
      ],
      "pilihan": [
       null,
       null
      ],
      "aksiUbah": [
       null,
       null
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
      "syarat": [],
      "baca": "selalu"
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
    "syarat": [],
    "baca": "selalu",
    "aksiUbah": [
     {
      "aksi": "postValue"
     },
     {
      "aksi": "refresh",
      "aktivitas": "TreatyInSetValueInstallment",
      "param": {
       "Installment": "TreatyIn.InstallmentNo"
      }
     }
    ]
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
    "baca": [
     null
    ],
    "tombol": [
     null
    ],
    "tombolKepala": [
     null
    ],
    "pilihan": [
     null
    ],
    "aksiUbah": [
     null
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInInstallment!pyGridRowDetails",
    "rincian": "Installments_ReadOnly"
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
    "baca": [
     null,
     null
    ],
    "tombol": [
     null,
     null
    ],
    "tombolKepala": [
     null,
     null
    ],
    "pilihan": [
     null,
     null
    ],
    "aksiUbah": [
     null,
     null
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
      "syarat": [],
      "baca": "selalu"
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
      "syarat": [],
      "baca": "selalu"
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
    "syarat": [],
    "baca": [
     "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
    ]
   },
   {
    "t": "medan",
    "at": 120405,
    "label": "End Date",
    "dari": "sisi",
    "kunci": "ReportingEnd",
    "format": "pxDateTime",
    "desimal": null,
    "syarat": [],
    "baca": [
     "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
    ]
   },
   {
    "t": "medan",
    "at": 139204,
    "label": "Period",
    "dari": "sisi",
    "kunci": "ReportingPeriod",
    "format": "pxDropdown",
    "desimal": null,
    "syarat": [],
    "baca": [
     "TreatyIn.ViewState = 1"
    ],
    "pilihan": {
     "sumber": "associated"
    }
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
    ],
    "baca": [
     "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
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
    "syarat": [],
    "baca": [
     "TreatyIn.IsEditData= 1 || TreatyIn.EDMMaterialType = 2"
    ]
   },
   {
    "t": "medan",
    "at": 204274,
    "label": "Confirmation",
    "dari": "sisi",
    "kunci": "ReportingConfirmation",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": [],
    "baca": [
     "TreatyIn.IsEditData= 1 || TreatyIn.EDMMaterialType = 2"
    ]
   },
   {
    "t": "medan",
    "at": 226262,
    "label": "Settlement",
    "dari": "sisi",
    "kunci": "ReportingSettlement",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": [],
    "baca": [
     "TreatyIn.IsEditData= 1 || TreatyIn.EDMMaterialType = 2"
    ]
   },
   {
    "t": "blok",
    "at": 247909,
    "judul": "",
    "syarat": [],
    "anak": [
     {
      "t": "tombol",
      "at": 254615,
      "label": "Apply",
      "syarat": [
       "TreatyIn.IsEditData !='1'"
      ],
      "aksi": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInSetReport",
        "param": {
         "startdate": "TreatyIn.ReportingStart",
         "autocalculate": "true"
        }
       }
      ],
      "nonaktif": [
       "TreatyIn.EDMMaterialType = 2"
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
    "baca": [
     "selalu",
     [
      "TreatyIn.EDMMaterialType = 2"
     ],
     [
      "TreatyIn.ViewState = 1",
      "TreatyIn.EDMMaterialType = 2"
     ],
     [
      "TreatyIn.ViewState = 1",
      "TreatyIn.EDMMaterialType = 2"
     ],
     [
      "TreatyIn.ViewState = 1",
      "TreatyIn.EDMMaterialType = 2"
     ],
     [
      "TreatyIn.ViewState = 1",
      "TreatyIn.EDMMaterialType = 2"
     ]
    ],
    "tombol": [
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "tombolKepala": [
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "pilihan": [
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "aksiUbah": [
     null,
     null,
     [
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInSetReport",
       "param": {
        "startdate": ".InitialDate",
        "autocalculate": ".AutoCalculate"
       }
      }
     ],
     null,
     null,
     null
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
       "Description",
       ""
      ],
      "kunci": [
       "TypePortfolio",
       "Type",
       "Description",
       ""
      ],
      "lebar": [
       111,
       157,
       648,
       102
      ],
      "desimal": [
       null,
       null,
       null,
       null
      ],
      "format": [
       "",
       "",
       "pxTextArea",
       "pxButton"
      ],
      "syaratSel": [
       null,
       null,
       null,
       "TreatyIn.IsEditData!='1'"
      ],
      "atSel": [
       467551,
       472622,
       477569,
       483459
      ],
      "baca": [
       [
        "TreatyIn.IsEditData= 1",
        "TreatyIn.EDMMaterialType = 2"
       ],
       [
        "TreatyIn.IsEditData= 1",
        "TreatyIn.EDMMaterialType = 2"
       ],
       [
        "TreatyIn.IsEditData= 1",
        "TreatyIn.EDMMaterialType = 2"
       ],
       "selalu"
      ],
      "tombol": [
       null,
       null,
       null,
       {
        "t": "tombol",
        "at": 483459,
        "label": "Delete",
        "syarat": [
         "TreatyIn.IsEditData!='1'"
        ],
        "aksi": [
         {
          "aksi": "deleteRow"
         }
        ],
        "nonaktif": [
         "TreatyIn.EDMMaterialType = 2"
        ]
       }
      ],
      "tombolKepala": [
       null,
       null,
       null,
       {
        "t": "tombol",
        "at": 457899,
        "label": "Add",
        "syarat": [
         "TreatyIn.IsEditData!='1'"
        ],
        "aksi": [
         {
          "aksi": "refresh",
          "aktivitas": "TreatyInPropAdd",
          "param": {
           "Type": "portfolio"
          }
         }
        ],
        "ikon": "pyWorkActionsAddWork.png",
        "nonaktif": [
         "TreatyIn.EDMMaterialType = 2"
        ]
       }
      ],
      "pilihan": [
       null,
       null,
       null,
       null
      ],
      "aksiUbah": [
       null,
       null,
       null,
       null
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
       "Kind of Treaty",
       ""
      ],
      "kunci": [
       "TreatyType",
       ""
      ],
      "lebar": [
       1315,
       80
      ],
      "desimal": [
       null,
       null
      ],
      "format": [
       "pxAutoComplete",
       "pxButton"
      ],
      "syaratSel": [
       null,
       "TreatyIn.ViewState !='1'"
      ],
      "atSel": [
       584772,
       592801
      ],
      "baca": [
       null,
       "selalu"
      ],
      "tombol": [
       null,
       {
        "t": "tombol",
        "at": 592801,
        "label": "Delete",
        "syarat": [
         "TreatyIn.ViewState !='1'"
        ],
        "aksi": [
         {
          "aksi": "deleteRow"
         }
        ],
        "nonaktif": [
         "TreatyIn.EDMMaterialType = 2"
        ]
       }
      ],
      "tombolKepala": [
       null,
       {
        "t": "tombol",
        "at": 575223,
        "label": "Add",
        "syarat": [
         "TreatyIn.ViewState !='1'"
        ],
        "aksi": [
         {
          "aksi": "refresh",
          "aktivitas": "TreatyInPropAdd",
          "param": {
           "Type": "limits"
          }
         }
        ],
        "ikon": "pyWorkActionsAddWork.png",
        "nonaktif": [
         "TreatyIn.EDMMaterialType = 2"
        ]
       }
      ],
      "pilihan": [
       {
        "sumber": "reportdefinition",
        "rd": "BrowseReinsuranceType_RD",
        "nilai": "Note"
       },
       null
      ],
      "aksiUbah": [
       null,
       null
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails",
      "rincian": "LimitProportional"
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
      ],
      "aksi": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInPropshare"
       }
      ],
      "nonaktif": [
       "TreatyIn.EDMMaterialType = 2"
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
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState =='1' || TreatyIn.EDMMaterialType = 2"
      ],
      "aksiUbah": [
       {
        "aksi": "postValue"
       },
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInPropshare"
       }
      ]
     },
     {
      "t": "medan",
      "at": 717944,
      "label": "% Brokerage",
      "dari": "sisi",
      "kunci": "BrokeragePercentP",
      "format": "pxNumber",
      "desimal": 2,
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState =='1' || TreatyIn.EDMMaterialType = 2"
      ],
      "aksiUbah": [
       {
        "aksi": "postValue"
       }
      ]
     },
     {
      "t": "medan",
      "at": 723387,
      "label": "Option",
      "dari": "sisi",
      "kunci": "OptionLimit",
      "format": "pxDropdown",
      "desimal": null,
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState =='1' || TreatyIn.EDMMaterialType = 2"
      ],
      "pilihan": {
       "sumber": "associated"
      },
      "aksiUbah": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInPropshare"
       }
      ]
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
      ],
      "baca": [
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ],
      "aksiUbah": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInPropshare"
       }
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
      ],
      "baca": [
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ],
      "aksiUbah": [
       {
        "aksi": "refresh"
       }
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
         "% Share",
         ""
        ],
        "kunci": [
         "ReinsName",
         "BrokerName",
         "SharePct",
         ""
        ],
        "lebar": [
         323,
         204,
         148,
         80
        ],
        "desimal": [
         null,
         null,
         2,
         null
        ],
        "format": [
         "pxLink",
         "pxLink",
         "pxNumber",
         "pxButton"
        ],
        "syaratSel": [
         null,
         null,
         null,
         "TreatyIn.ViewState !='1'"
        ],
        "atSel": [
         886836,
         904928,
         922680,
         941807
        ],
        "baca": [
         [
          "TreatyIn.ViewState = 1"
         ],
         [
          "TreatyIn.ViewState = 1"
         ],
         [
          "TreatyIn.ViewState = 1"
         ],
         "selalu"
        ],
        "tombol": [
         null,
         null,
         null,
         {
          "t": "tombol",
          "at": 941807,
          "label": "Delete",
          "syarat": [
           "TreatyIn.ViewState !='1'"
          ],
          "aksi": [
           {
            "aksi": "deleteRow"
           },
           {
            "aksi": "refresh"
           }
          ]
         }
        ],
        "tombolKepala": [
         null,
         null,
         null,
         {
          "t": "tombol",
          "at": 876741,
          "label": "Add",
          "syarat": [
           "TreatyIn.ViewState !='1'"
          ],
          "aksi": [
           {
            "aksi": "refresh",
            "aktivitas": "AddFacRetroProp"
           }
          ],
          "ikon": "rpadd.gif"
         }
        ],
        "pilihan": [
         null,
         null,
         null,
         null
        ],
        "aksiUbah": [
         null,
         null,
         [
          {
           "aksi": "postValue"
          },
          {
           "aksi": "runActivity",
           "aktivitas": "CountRetroShare_Act"
          },
          {
           "aksi": "refreshRowItem"
          },
          {
           "aksi": "refresh"
          },
          {
           "aksi": "refresh"
          },
          {
           "aksi": "refresh"
          }
         ],
         null
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
    "t": "grid",
    "at": 1111572,
    "prop": "TreatyIn.CoInScale",
    "dari": "sisi",
    "larik": "CoInScale",
    "syarat": [],
    "kolom": [
     "Co-Insurance Share",
     "% Treaty Limit",
     ""
    ],
    "kunci": [
     "CoInShare",
     "PctLimit",
     ""
    ],
    "lebar": [
     206,
     206,
     101
    ],
    "desimal": [
     null,
     null,
     null
    ],
    "format": [
     "pxTextInput",
     "pxNumber",
     "pxButton"
    ],
    "syaratSel": [
     null,
     null,
     "TreatyIn.IsEditData!='1'"
    ],
    "atSel": [
     1134997,
     1143339,
     1150344
    ],
    "baca": [
     [
      "TreatyIn.IsEditData='1'"
     ],
     [
      "TreatyIn.IsEditData='1'"
     ],
     "selalu"
    ],
    "tombol": [
     null,
     null,
     {
      "t": "tombol",
      "at": 1150344,
      "label": "",
      "syarat": [
       "TreatyIn.IsEditData!='1'"
      ],
      "aksi": [
       {
        "aksi": "deleteRow"
       }
      ],
      "ikon": "IconTrash.png"
     }
    ],
    "tombolKepala": [
     null,
     null,
     {
      "t": "tombol",
      "at": 1123559,
      "label": "",
      "syarat": [
       "TreatyIn.IsEditData!='1'"
      ],
      "aksi": [
       {
        "aksi": "addRow"
       },
       {
        "aksi": "refresh"
       }
      ],
      "ikon": "IconAdd.png"
     }
    ],
    "pilihan": [
     null,
     null,
     null
    ],
    "aksiUbah": [
     [
      {
       "aksi": "postValue"
      }
     ],
     [
      {
       "aksi": "postValue"
      }
     ],
     null
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInCoInScaleDetails!pyGridModalTemplate"
   },
   {
    "t": "medan",
    "at": 1190828,
    "label": "Max Co-Insurance Panel (Non Group)",
    "dari": "sisi",
    "kunci": "MaxCoNonGroup",
    "format": "pxNumber",
    "desimal": null,
    "syarat": [],
    "baca": [
     "TreatyIn.ViewState ='1'",
     "TreatyIn.IsEditData ='1'"
    ],
    "aksiUbah": [
     {
      "aksi": "postValue"
     }
    ]
   },
   {
    "t": "medan",
    "at": 1198436,
    "label": "Max Co-Insurance Panel (Group)",
    "dari": "sisi",
    "kunci": "MaxCoGroup",
    "format": "pxNumber",
    "desimal": null,
    "syarat": [],
    "baca": [
     "TreatyIn.ViewState ='1'",
     "TreatyIn.IsEditData ='1'"
    ],
    "aksiUbah": [
     {
      "aksi": "postValue"
     }
    ]
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
    "judul": "",
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
      "syarat": [],
      "baca": [
       "TreatyIn.ViewState = 1",
       "TreatyIn.EDMMaterialType = 2"
      ],
      "pilihan": {
       "sumber": "associated"
      },
      "aksiUbah": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInSetAccountReport",
        "transformasi": "TreatyInDeleteAccumulationLists"
       }
      ]
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
     "Submission Due",
     ""
    ],
    "kunci": [
     "Period",
     "ReportDate",
     "SubDays",
     "SubDueDate",
     ""
    ],
    "lebar": [
     98,
     202,
     199,
     215,
     87
    ],
    "desimal": [
     null,
     null,
     null,
     null,
     null
    ],
    "format": [
     "pxTextInput",
     "pxDateTime",
     "pxTextInput",
     "pxDateTime",
     "pxButton"
    ],
    "syaratSel": [
     null,
     null,
     null,
     null,
     "TreatyIn.ViewState !='1'"
    ],
    "atSel": [
     1374097,
     1380069,
     1389247,
     1398507,
     1404482
    ],
    "baca": [
     "selalu",
     [
      "TreatyIn.ViewState = 1",
      "TreatyIn.EDMMaterialType = 2"
     ],
     [
      "TreatyIn.ViewState = 1",
      "TreatyIn.EDMMaterialType = 2"
     ],
     "selalu",
     "selalu"
    ],
    "tombol": [
     null,
     null,
     null,
     null,
     {
      "t": "tombol",
      "at": 1404482,
      "label": "Delete",
      "syarat": [
       "TreatyIn.ViewState !='1'"
      ],
      "aksi": [
       {
        "aksi": "deleteRow"
       }
      ],
      "nonaktif": [
       "TreatyIn.EDMMaterialType = 2"
      ]
     }
    ],
    "tombolKepala": [
     null,
     null,
     null,
     null,
     {
      "t": "tombol",
      "at": 1364703,
      "label": "Add",
      "syarat": [
       "TreatyIn.ViewState !='1'"
      ],
      "aksi": [
       {
        "aksi": "refresh",
        "transformasi": "TreatyInAddAccumulation"
       }
      ],
      "ikon": "pyWorkActionsAddWork.png",
      "nonaktif": [
       "TreatyIn.EDMMaterialType = 2"
      ]
     }
    ],
    "pilihan": [
     null,
     null,
     null,
     null,
     null
    ],
    "aksiUbah": [
     null,
     [
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInAccumulationSetSubDue"
      }
     ],
     [
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInAccumulationSetSubDue"
      }
     ],
     null,
     null
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
      "syarat": [],
      "baca": [
       "TreatyIn.IsEditData= 1 || TreatyIn.EDMMaterialType = 1"
      ],
      "aksiUbah": [
       {
        "aksi": "refresh",
        "transformasi": "TreatyInCopyConditions"
       }
      ]
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
      "syarat": [],
      "baca": [
       "TreatyIn.IsEditData= 1 || TreatyIn.EDMMaterialType = 1"
      ],
      "aksiUbah": [
       {
        "aksi": "refresh",
        "transformasi": "TreatyInCopyConditions"
       }
      ]
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
    "baca": [
     null
    ],
    "tombol": [
     null
    ],
    "tombolKepala": [
     null
    ],
    "pilihan": [
     {
      "sumber": "reportdefinition",
      "rd": "BrowseReinsuranceType_RD",
      "nilai": "Note"
     }
    ],
    "aksiUbah": [
     null
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails",
    "rincian": "AchievementCombine"
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
    "syarat": [],
    "baca": "selalu"
   },
   {
    "t": "medan",
    "at": 116989,
    "label": "End Date",
    "dari": "sisi",
    "kunci": "ReportingEnd",
    "format": "pxDateTime",
    "desimal": null,
    "syarat": [],
    "baca": "selalu"
   },
   {
    "t": "medan",
    "at": 135711,
    "label": "Period",
    "dari": "sisi",
    "kunci": "ReportingPeriod",
    "format": "pxDropdown",
    "desimal": null,
    "syarat": [],
    "baca": "selalu",
    "pilihan": {
     "sumber": "associated"
    }
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
    ],
    "baca": "selalu"
   },
   {
    "t": "medan",
    "at": 178702,
    "label": "Submission",
    "dari": "sisi",
    "kunci": "ReportingSubmission",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": [],
    "baca": "selalu"
   },
   {
    "t": "medan",
    "at": 200222,
    "label": "Confirmation",
    "dari": "sisi",
    "kunci": "ReportingConfirmation",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": [],
    "baca": "selalu"
   },
   {
    "t": "medan",
    "at": 221750,
    "label": "Settlement",
    "dari": "sisi",
    "kunci": "ReportingSettlement",
    "format": "pxTextInput",
    "desimal": null,
    "syarat": [],
    "baca": "selalu"
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
    "baca": [
     "selalu",
     "selalu",
     "selalu",
     "selalu",
     "selalu",
     "selalu"
    ],
    "tombol": [
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "tombolKepala": [
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "pilihan": [
     null,
     null,
     null,
     null,
     null,
     null
    ],
    "aksiUbah": [
     null,
     null,
     [
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInSetReport",
       "param": {
        "startdate": ".InitialDate",
        "autocalculate": ".AutoCalculate"
       }
      }
     ],
     null,
     null,
     null
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
      "baca": [
       "selalu",
       "selalu",
       "selalu"
      ],
      "tombol": [
       null,
       null,
       null
      ],
      "tombolKepala": [
       null,
       null,
       null
      ],
      "pilihan": [
       null,
       null,
       null
      ],
      "aksiUbah": [
       null,
       null,
       null
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
      "baca": [
       null
      ],
      "tombol": [
       null
      ],
      "tombolKepala": [
       null
      ],
      "pilihan": [
       {
        "sumber": "reportdefinition",
        "rd": "BrowseReinsuranceType_RD",
        "nilai": "Note"
       }
      ],
      "aksiUbah": [
       null
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails",
      "rincian": "LimitProportionalOldData"
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
      "syarat": [],
      "baca": "selalu",
      "aksiUbah": [
       {
        "aksi": "postValue"
       },
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInPropshare"
       }
      ]
     },
     {
      "t": "medan",
      "at": 599000,
      "label": "% Brokerage",
      "dari": "sisi",
      "kunci": "BrokeragePercentP",
      "format": "pxNumber",
      "desimal": 2,
      "syarat": [],
      "baca": "selalu",
      "aksiUbah": [
       {
        "aksi": "postValue"
       }
      ]
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
    "baca": [
     null
    ],
    "tombol": [
     null
    ],
    "tombolKepala": [
     null
    ],
    "pilihan": [
     {
      "sumber": "reportdefinition",
      "rd": "BrowseReinsuranceType_RD",
      "nilai": "Note"
     }
    ],
    "aksiUbah": [
     null
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails",
    "rincian": "TotalLimitsOldData"
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
    "baca": [
     null,
     null
    ],
    "tombol": [
     null,
     null
    ],
    "tombolKepala": [
     null,
     null
    ],
    "pilihan": [
     null,
     null
    ],
    "aksiUbah": [
     null,
     null
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
    "baca": [
     null,
     null
    ],
    "tombol": [
     null,
     null
    ],
    "tombolKepala": [
     null,
     null
    ],
    "pilihan": [
     null,
     null
    ],
    "aksiUbah": [
     null,
     null
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
    "baca": [
     null,
     null
    ],
    "tombol": [
     null,
     null
    ],
    "tombolKepala": [
     null,
     null
    ],
    "pilihan": [
     null,
     null
    ],
    "aksiUbah": [
     null,
     null
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
    "judul": "",
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
      "syarat": [],
      "baca": "selalu",
      "pilihan": {
       "sumber": "associated"
      },
      "aksiUbah": [
       {
        "aksi": "refresh",
        "aktivitas": "TreatyInSetAccountReport",
        "transformasi": "TreatyInDeleteAccumulationLists"
       }
      ]
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
    "baca": [
     "selalu",
     "selalu",
     "selalu",
     "selalu"
    ],
    "tombol": [
     null,
     null,
     null,
     null
    ],
    "tombolKepala": [
     null,
     null,
     null,
     null
    ],
    "pilihan": [
     null,
     null,
     null,
     null
    ],
    "aksiUbah": [
     null,
     [
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInAccumulationSetSubDue"
      }
     ],
     [
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInAccumulationSetSubDue"
      }
     ],
     null
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
      "syarat": [],
      "baca": "selalu",
      "aksiUbah": [
       {
        "aksi": "refresh",
        "transformasi": "TreatyInCopyConditions"
       }
      ]
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
      "syarat": [],
      "baca": "selalu",
      "aksiUbah": [
       {
        "aksi": "refresh",
        "transformasi": "TreatyInCopyConditions"
       }
      ]
     }
    ]
   }
  ]
 },
 "TreatyInTabsNonProportionalAdjustPremi#Actual GNPI": {
  "at": 20240,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 39307,
    "judul": "Actual Premium",
    "syarat": [],
    "anak": [
     {
      "t": "grid",
      "at": 56680,
      "prop": "TreatyIn.ActualValue.EGNPI",
      "dari": "sisi",
      "larik": "ActualValue.EGNPI",
      "syarat": [],
      "kolom": [
       "Treaty Group",
       "As Date",
       "Proportion %",
       "Currency",
       "Amount",
       "Amount in IDR",
       ""
      ],
      "kunci": [
       "TreatyGroup",
       "AsDate",
       "Proportion",
       "Currency",
       "Amount",
       "AmountIDR",
       ""
      ],
      "lebar": [
       292,
       202,
       144,
       139,
       228,
       230,
       128
      ],
      "desimal": [
       null,
       null,
       2,
       null,
       2,
       null,
       null
      ],
      "format": [
       "pxAutoComplete",
       "",
       "pxNumber",
       "",
       "pxNumber",
       "pxNumber",
       "pxButton"
      ],
      "syaratSel": [
       null,
       null,
       null,
       null,
       null,
       null,
       "TreatyIn.ViewState !='1'"
      ],
      "atSel": [
       94893,
       100348,
       104793,
       109889,
       114335,
       119426,
       124484
      ],
      "baca": [
       [
        "1=1"
       ],
       null,
       null,
       null,
       null,
       null,
       "selalu"
      ],
      "tombol": [
       null,
       null,
       null,
       null,
       null,
       null,
       {
        "t": "tombol",
        "at": 124484,
        "label": "Delete",
        "syarat": [
         "TreatyIn.ViewState !='1'"
        ],
        "aksi": [
         {
          "aksi": "deleteRow"
         }
        ]
       }
      ],
      "tombolKepala": [
       null,
       null,
       null,
       null,
       null,
       null,
       {
        "t": "tombol",
        "at": 85260,
        "label": "Add",
        "syarat": [
         "TreatyIn.ViewState !='1'"
        ],
        "aksi": [
         {
          "aksi": "refresh",
          "aktivitas": "TreatyInNonAddItem",
          "param": {
           "Type": "actualpremium"
          }
         }
        ],
        "ikon": "rpadd.gif"
       }
      ],
      "pilihan": [
       {
        "sumber": "associated"
       },
       null,
       null,
       null,
       null,
       null,
       null
      ],
      "aksiUbah": [
       null,
       null,
       null,
       null,
       null,
       null,
       null
      ],
      "templatBaris": "ASM-FW-GISFW-Data-TreatyInEGNPI!pyGridRowDetails",
      "rincian": "DetailEGNPI"
     }
    ]
   },
   {
    "t": "grid",
    "at": 187941,
    "prop": "TreatyIn.ActualValue.TotalEgnpiAmountNP",
    "dari": "sisi",
    "larik": "ActualValue.TotalEgnpiAmountNP",
    "syarat": [],
    "kolom": [
     "Total Actual Premium Amount",
     "Value"
    ],
    "kunci": [
     "Currency",
     "Value"
    ],
    "lebar": [
     195,
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
     199933,
     204889
    ],
    "baca": [
     null,
     null
    ],
    "tombol": [
     null,
     null
    ],
    "tombolKepala": [
     null,
     null
    ],
    "pilihan": [
     null,
     null
    ],
    "aksiUbah": [
     null,
     null
    ],
    "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
   },
   {
    "t": "teks",
    "at": 266681,
    "teks": "Total Amount in IDR",
    "syarat": []
   },
   {
    "t": "teks",
    "at": 286423,
    "teks": "IDR",
    "syarat": []
   },
   {
    "t": "medan",
    "at": 290811,
    "label": "",
    "dari": "sisi",
    "kunci": "ActualValue.TotalEgnpiAmount",
    "format": "pxNumber",
    "desimal": null,
    "syarat": [],
    "baca": "selalu"
   },
   {
    "t": "teks",
    "at": 323761,
    "teks": "Total Proportion %",
    "syarat": []
   },
   {
    "t": "medan",
    "at": 334504,
    "label": "",
    "dari": "sisi",
    "kunci": "ActualValue.TotalEgnpiProportion",
    "format": "pxNumber",
    "desimal": 2,
    "syarat": [],
    "baca": "selalu"
   },
   {
    "t": "tombol",
    "at": 364838,
    "label": "Update Total",
    "syarat": [
     "TreatyIn.ViewState !='1'"
    ],
    "aksi": [
     {
      "aksi": "refresh",
      "aktivitas": "TreatyInActualUpdateValue",
      "param": {
       "type": ""
      }
     }
    ]
   }
  ]
 },
 "TreatyInTabsNonProportionalAdjustPremi#Actual Limits": {
  "at": 395887,
  "syarat": [],
  "isi": [
   {
    "t": "include",
    "at": 403864,
    "nama": "TreatyInActualLimits",
    "syarat": []
   }
  ]
 },
 "TreatyInTabsNonProportionalAdjustPremi#Actual Share": {
  "at": 437918,
  "syarat": [],
  "isi": [
   {
    "t": "include",
    "at": 445907,
    "nama": "TreatyInActualShare",
    "syarat": []
   }
  ]
 },
 "TreatyInTabsNonProportionalAdjustPremi#Actual Retro": {
  "at": 479958,
  "syarat": [],
  "isi": [
   {
    "t": "blok",
    "at": 490433,
    "judul": "",
    "syarat": [
     "TreatyIn.IsMultipleRetro = 'true'"
    ],
    "anak": [
     {
      "t": "include",
      "at": 497133,
      "nama": "TreatyInActualFacultativeShareCalculation",
      "syarat": []
     }
    ]
   }
  ]
 },
 "TreatyInTabsNonProportionalAdjustPremi#Premium Adjustment": {
  "at": 539146,
  "syarat": [],
  "isi": [
   {
    "t": "include",
    "at": 555839,
    "nama": "TreatyInActualSumary",
    "syarat": []
   }
  ]
 },
 "TreatyInTabsNonProportionalAdjustPremi#Information & Submit": {
  "at": 596060,
  "syarat": [],
  "isi": [
   {
    "t": "include",
    "at": 603615,
    "nama": "TreatyInfoSubmit",
    "syarat": []
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
     "syarat": [],
     "baca": "selalu",
     "aksiUbah": [
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInXOLAddSpreading"
      }
     ]
    },
    {
     "t": "medan",
     "at": 43965,
     "label": "% Brokerage",
     "dari": "sisi",
     "kunci": "BrokeragePercent",
     "format": "pxTextInput",
     "desimal": null,
     "syarat": [],
     "baca": "selalu",
     "aksiUbah": [
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInSetBrokerage"
      },
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInXOLAddSpreading"
      }
     ]
    },
    {
     "t": "medan",
     "at": 70496,
     "label": "Facultative Share",
     "dari": "sisi",
     "kunci": "FacultativeShare",
     "format": "pxTextInput",
     "desimal": null,
     "syarat": [],
     "baca": "selalu",
     "aksiUbah": [
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInXOLAddSpreading"
      }
     ]
    },
    {
     "t": "medan",
     "at": 79208,
     "label": "Fakultative Brokerage",
     "dari": "sisi",
     "kunci": "FacultativeShareBrokerage",
     "format": "pxTextInput",
     "desimal": null,
     "syarat": [],
     "baca": "selalu",
     "aksiUbah": [
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInXOLAddSpreading"
      }
     ]
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
     "baca": [
      "selalu",
      "selalu",
      "selalu"
     ],
     "tombol": [
      null,
      null,
      null
     ],
     "tombolKepala": [
      null,
      null,
      null
     ],
     "pilihan": [
      {
       "sumber": "reportdefinition",
       "rd": "BrowseAgentNusaRe_RD",
       "nilai": "ClientName"
      },
      null,
      null
     ],
     "aksiUbah": [
      null,
      null,
      null
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
         "syarat": [],
         "baca": "selalu"
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
       "judul": "",
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
         "baca": [
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
         "tombol": [
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
         "tombolKepala": [
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
         "pilihan": [
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
         "aksiUbah": [
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
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridRowDetails",
         "rincian": "ShareOldData"
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
         "baca": [
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
         "tombol": [
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
         "tombolKepala": [
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
         "pilihan": [
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
         "aksiUbah": [
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
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
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
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
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
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
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
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "blok",
     "at": 854379,
     "judul": "",
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
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
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
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
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
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
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
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "tombol",
     "at": 1171408,
     "label": "Hide Facultative Share (unused)",
     "syarat": [
      "TreatyIn.FacultativeShare>0 && facsharedisp.CARI1 = '1'"
     ],
     "aksi": [
      {
       "aksi": "refresh",
       "transformasi": "TreatyInShowHideFacShare",
       "paramDT": {
        "type": "0"
       }
      }
     ]
    },
    {
     "t": "tombol",
     "at": 1180222,
     "label": "Show Facultative Share",
     "syarat": [
      "TreatyIn.FacultativeShare>0"
     ],
     "aksi": [
      {
       "aksi": "showHarness"
      }
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
         "judul": "",
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
           "baca": [
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
           "tombol": [
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
           "tombolKepala": [
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
           "pilihan": [
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
           "aksiUbah": [
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
           "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridRowDetails",
           "rincian": "Share"
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
           "baca": [
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
           "tombol": [
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
           "tombolKepala": [
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
           "pilihan": [
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
           "aksiUbah": [
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "blok",
         "at": 1718323,
         "judul": "",
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
           "baca": [
            null,
            null
           ],
           "tombol": [
            null,
            null
           ],
           "tombolKepala": [
            null,
            null
           ],
           "pilihan": [
            null,
            null
           ],
           "aksiUbah": [
            null,
            null
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
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
       "syarat": [],
       "baca": "selalu"
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
     "baca": [
      null
     ],
     "tombol": [
      null
     ],
     "tombolKepala": [
      null
     ],
     "pilihan": [
      {
       "sumber": "reportdefinition",
       "rd": "BrowseReinsuranceType_RD",
       "nilai": "Note"
      }
     ],
     "aksiUbah": [
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails",
     "rincian": "TotalLimits"
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
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
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
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
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
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
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
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
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
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh"
    }
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
     ],
     "aksi": [
      {
       "aksi": "runActivity",
       "aktivitas": "TreatyInSubmit"
      },
      {
       "aksi": "refresh"
      },
      {
       "aksi": "refresh"
      }
     ]
    },
    {
     "t": "tombol",
     "at": 179405,
     "label": "Submit",
     "syarat": [
      "TreatyIn.RevisionState='1'"
     ],
     "aksi": [
      {
       "aksi": "runActivity",
       "aktivitas": "TreatyInSubmit"
      },
      {
       "aksi": "refresh"
      },
      {
       "aksi": "refresh"
      }
     ]
    },
    {
     "t": "tombol",
     "at": 200534,
     "label": "Decline offer",
     "syarat": [
      "TreatyIn.ViewState != '1'"
     ],
     "aksi": [
      {
       "aksi": "localAction",
       "aktivitas": "TreatyInDeclineConfirmation"
      }
     ],
     "ikon": "pi pi-mobile-phone-solid",
     "nonaktif": [
      "TreatyIn.ID = ''"
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
     ],
     "aksi": [
      {
       "aksi": "refresh"
      },
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInSubmitEDM"
      }
     ],
     "nonaktif": [
      "TreatyIn.EDMEffective = ''"
     ]
    },
    {
     "t": "tombol",
     "at": 241510,
     "label": "Decline offer",
     "syarat": [
      "TreatyIn.ViewState != '1'"
     ],
     "aksi": [
      {
       "aksi": "localAction",
       "aktivitas": "TreatyInDeclineConfirmationEDM"
      }
     ],
     "ikon": "pi pi-mobile-phone-solid",
     "nonaktif": [
      "TreatyIn.ID = ''"
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
     ],
     "aksi": [
      {
       "aksi": "refresh",
       "aktivitas": "TreatyEDMCalculateDifference"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 47293,
   "judul": "",
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
   "judul": "",
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
       "syarat": [],
       "baca": "selalu"
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
       "syarat": [],
       "baca": "selalu"
      },
      {
       "t": "medan",
       "at": 81824,
       "label": "Pro Rate Percentage",
       "dari": "sisi",
       "kunci": "ProRatePercent",
       "format": "pxTextInput",
       "desimal": 4,
       "syarat": [],
       "baca": "selalu"
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
       "syarat": [],
       "baca": "selalu",
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInXOLAddSpreading"
        }
       ]
      },
      {
       "t": "medan",
       "at": 143662,
       "label": "% Brokerage",
       "dari": "sisi",
       "kunci": "ValueDifference.BrokeragePercent",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu",
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInSetBrokerage"
        },
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInXOLAddSpreading"
        }
       ]
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
       ],
       "baca": [
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ],
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInXOLAddSpreading"
        }
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
           "syarat": [],
           "baca": "selalu"
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
         "judul": "",
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
           "baca": [
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
           "tombol": [
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
           "tombolKepala": [
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
           "pilihan": [
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
           "aksiUbah": [
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
           "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridRowDetails",
           "rincian": "ShareOldData"
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
           "baca": [
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
           "tombol": [
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
           "tombolKepala": [
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
           "pilihan": [
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
           "aksiUbah": [
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "blok",
         "at": 1102143,
         "judul": "",
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
           "baca": [
            null,
            null
           ],
           "tombol": [
            null,
            null
           ],
           "tombolKepala": [
            null,
            null
           ],
           "pilihan": [
            null,
            null
           ],
           "aksiUbah": [
            null,
            null
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
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
       "syarat": [],
       "baca": "selalu",
       "aksiUbah": [
        {
         "aksi": "postValue"
        },
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInSetValueInstallment",
         "param": {
          "Installment": "TreatyIn.InstallmentNo"
         }
        }
       ]
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
       "baca": [
        null
       ],
       "tombol": [
        null
       ],
       "tombolKepala": [
        null
       ],
       "pilihan": [
        null
       ],
       "aksiUbah": [
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInInstallment!pyGridRowDetails",
       "rincian": "Installments"
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
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
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
   "judul": "",
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
       "syarat": [],
       "baca": "selalu",
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInXOLAddSpreading"
        }
       ]
      },
      {
       "t": "medan",
       "at": 52547,
       "label": "% Brokerage",
       "dari": "sisi",
       "kunci": "ValueBeforeProrate.BrokeragePercent",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu",
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInSetBrokerage"
        },
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInXOLAddSpreading"
        }
       ]
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
       ],
       "baca": [
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ],
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInXOLAddSpreading"
        }
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
           "syarat": [],
           "baca": "selalu"
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
         "judul": "",
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
           "baca": [
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
           "tombol": [
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
           "tombolKepala": [
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
           "pilihan": [
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
           "aksiUbah": [
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
           "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridRowDetails",
           "rincian": "ShareOldData"
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
           "baca": [
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
           "tombol": [
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
           "tombolKepala": [
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
           "pilihan": [
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
           "aksiUbah": [
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        },
        {
         "t": "blok",
         "at": 1015495,
         "judul": "",
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
           "baca": [
            null,
            null
           ],
           "tombol": [
            null,
            null
           ],
           "tombolKepala": [
            null,
            null
           ],
           "pilihan": [
            null,
            null
           ],
           "aksiUbah": [
            null,
            null
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
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
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
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
       "syarat": [],
       "baca": "selalu",
       "aksiUbah": [
        {
         "aksi": "postValue"
        },
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInSetValueInstallment",
         "param": {
          "Installment": "TreatyIn.InstallmentNo"
         }
        }
       ]
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
       "baca": [
        null
       ],
       "tombol": [
        null
       ],
       "tombolKepala": [
        null
       ],
       "pilihan": [
        null
       ],
       "aksiUbah": [
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInInstallment!pyGridRowDetails",
       "rincian": "Installments"
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
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    }
   ]
  }
 ],
 "TreatyInActualLimits": [
  {
   "t": "blok",
   "at": 8781,
   "judul": "",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 24247,
     "judul": "",
     "syarat": [],
     "anak": [
      {
       "t": "grid",
       "at": 40960,
       "prop": "TreatyIn.ActualValue.Limits",
       "dari": "sisi",
       "larik": "ActualValue.Limits",
       "syarat": [],
       "kolom": [
        "Layers",
        "",
        "",
        "",
        ""
       ],
       "kunci": [
        "LayerType",
        "Layer",
        "pyTemplateRichTextEditor",
        "LayerPartType",
        "LayerPart"
       ],
       "lebar": [
        102,
        102,
        102,
        102,
        105
       ],
       "desimal": [
        null,
        null,
        null,
        null,
        null
       ],
       "format": [
        "pxTextInput",
        "",
        "pxTextInput",
        "pxTextInput",
        ""
       ],
       "syaratSel": [
        null,
        null,
        null,
        null,
        null
       ],
       "atSel": [
        64549,
        70238,
        74441,
        79784,
        85381
       ],
       "baca": [
        null,
        null,
        "selalu",
        null,
        null
       ],
       "tombol": [
        null,
        null,
        null,
        null,
        null
       ],
       "tombolKepala": [
        null,
        null,
        null,
        null,
        null
       ],
       "pilihan": [
        null,
        null,
        null,
        null,
        null
       ],
       "aksiUbah": [
        null,
        null,
        null,
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails",
       "rincian": "LayersEDM"
      }
     ]
    },
    {
     "t": "blok",
     "at": 118335,
     "judul": "Summary of Actual Limits",
     "syarat": [],
     "anak": [
      {
       "t": "grid",
       "at": 143723,
       "prop": "TreatyIn.ActualValue.LimitSummaryList",
       "dari": "sisi",
       "larik": "ActualValue.LimitSummaryList",
       "syarat": [],
       "kolom": [
        "Note",
        "MDP (IDR)",
        "MDP (USD)"
       ],
       "kunci": [
        "Note",
        "MDP",
        "MDP2"
       ],
       "lebar": [
        252,
        100,
        100
       ],
       "desimal": [
        null,
        null,
        null
       ],
       "format": [
        "",
        "",
        ""
       ],
       "syaratSel": [
        null,
        null,
        null
       ],
       "atSel": [
        159203,
        163402,
        167607
       ],
       "baca": [
        null,
        null,
        null
       ],
       "tombol": [
        null,
        null,
        null
       ],
       "tombolKepala": [
        null,
        null,
        null
       ],
       "pilihan": [
        null,
        null,
        null
       ],
       "aksiUbah": [
        null,
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-LimitSummaryList!pyGridModalTemplate"
      }
     ]
    },
    {
     "t": "blok",
     "at": 204923,
     "judul": "Total All Layers",
     "syarat": [],
     "anak": [
      {
       "t": "grid",
       "at": 230304,
       "prop": "TreatyIn.ActualValue.TotalLimitPremiEarnNP",
       "dari": "sisi",
       "larik": "ActualValue.TotalLimitPremiEarnNP",
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
        241880,
        246594
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 293852,
       "prop": "TreatyIn.ActualValue.TotalLimitMDPNP",
       "dari": "sisi",
       "larik": "ActualValue.TotalLimitMDPNP",
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
        305411,
        310125
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "teks",
       "at": 370473,
       "teks": "Total ROL",
       "syarat": []
      },
      {
       "t": "medan",
       "at": 381156,
       "label": "",
       "dari": "sisi",
       "kunci": "ActualValue.TotalLimitsROL",
       "format": "pxNumber",
       "desimal": 2,
       "syarat": [],
       "baca": "selalu"
      },
      {
       "t": "tombol",
       "at": 404884,
       "label": "Update Total",
       "syarat": [
        "TreatyIn.ViewState !='1'"
       ],
       "aksi": [
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInNPSetTotalActual",
         "param": {
          "type": "limits"
         }
        },
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInSummaryLimitActual"
        },
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInSummaryLimitShareActual"
        }
       ]
      }
     ]
    }
   ]
  }
 ],
 "TreatyInActualShare": [
  {
   "t": "blok",
   "at": 8749,
   "judul": "Actual Share",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 38157,
     "judul": "",
     "syarat": [],
     "anak": [
      {
       "t": "medan",
       "at": 61884,
       "label": "% RNM Share",
       "dari": "sisi",
       "kunci": "ActualValue.RNMShare",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu",
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInXOLAddSpreading"
        }
       ]
      },
      {
       "t": "medan",
       "at": 70213,
       "label": "% Brokerage",
       "dari": "sisi",
       "kunci": "ActualValue.BrokeragePercent",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu",
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInSetBrokerage"
        },
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInXOLAddSpreading"
        }
       ]
      },
      {
       "t": "medan",
       "at": 96748,
       "label": "Facultative Share",
       "dari": "sisi",
       "kunci": "FacultativeShare",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu",
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInXOLAddSpreading"
        }
       ]
      },
      {
       "t": "medan",
       "at": 105452,
       "label": "Fakultative Brokerage",
       "dari": "sisi",
       "kunci": "FacultativeShareBrokerage",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu",
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInXOLAddSpreading"
        }
       ]
      },
      {
       "t": "tombol",
       "at": 127216,
       "label": "Update Summary",
       "syarat": [
        "TreatyIn.ViewState !='1'"
       ],
       "aksi": [
        {
         "aksi": "setValue",
         "param": {
          "TreatyIn.FacultativeShare": "0"
         },
         "syarat": "TreatyIn.FacultativeShare = ''"
        },
        {
         "aksi": "refresh",
         "aktivitas": "TreatyInActualUpdateValueShare",
         "param": {
          "type": "share"
         }
        },
        {
         "aksi": "refresh"
        },
        {
         "aksi": "refresh"
        }
       ]
      },
      {
       "t": "grid",
       "at": 162556,
       "prop": "TreatyIn.ActualValue.ShareReins",
       "dari": "sisi",
       "larik": "ActualValue.ShareReins",
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
        178154,
        185927,
        191317
       ],
       "baca": [
        null,
        null,
        null
       ],
       "tombol": [
        null,
        null,
        null
       ],
       "tombolKepala": [
        null,
        null,
        null
       ],
       "pilihan": [
        {
         "sumber": "reportdefinition",
         "rd": "BrowseAgentNusaRe_RD",
         "nilai": "ClientName"
        },
        null,
        null
       ],
       "aksiUbah": [
        null,
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInShareReins!pyGridModalTemplate"
      },
      {
       "t": "blok",
       "at": 250929,
       "judul": "RNM Share",
       "syarat": [],
       "anak": [
        {
         "t": "blok",
         "at": 259596,
         "judul": "",
         "syarat": [],
         "anak": [
          {
           "t": "grid",
           "at": 276362,
           "prop": "TreatyIn.ActualValue.Share",
           "dari": "sisi",
           "larik": "ActualValue.Share",
           "syarat": [],
           "kolom": [
            "",
            "",
            "",
            "",
            "",
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
            307716,
            313406,
            317605,
            321761,
            327356,
            331561,
            337206,
            341522,
            347167
           ],
           "baca": [
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
           "tombol": [
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
           "tombolKepala": [
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
           "pilihan": [
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
           "aksiUbah": [
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
           "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridRowDetails",
           "rincian": "Share"
          }
         ]
        },
        {
         "t": "blok",
         "at": 382309,
         "judul": "Summarry of RNM Share",
         "syarat": [],
         "anak": [
          {
           "t": "grid",
           "at": 399024,
           "prop": "TreatyIn.ActualValue.LimitShareSummaryList",
           "dari": "sisi",
           "larik": "ActualValue.LimitShareSummaryList",
           "syarat": [],
           "kolom": [
            "Note",
            "MDP (IDR)",
            "MDP (USD)",
            "Deduction (IDR)",
            "Deduction (USD)",
            "Net Premi (IDR)",
            "Net Premi (USD)"
           ],
           "kunci": [
            "Note",
            "MDP",
            "MDP2",
            "Deductible",
            "Deductible2",
            "NetPremi",
            "NetPremi2"
           ],
           "lebar": [
            293,
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
            null
           ],
           "format": [
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
            null
           ],
           "atSel": [
            429928,
            434127,
            438332,
            442538,
            446751,
            450965,
            455175
           ],
           "baca": [
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "tombol": [
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "tombolKepala": [
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "pilihan": [
            null,
            null,
            null,
            null,
            null,
            null,
            null
           ],
           "aksiUbah": [
            null,
            null,
            null,
            null,
            null,
            null,
            null
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
     "at": 507509,
     "judul": "Total All Layers RNM Share",
     "syarat": [],
     "anak": [
      {
       "t": "grid",
       "at": 541568,
       "prop": "TreatyIn.ActualValue.TotalShareGrossNP",
       "dari": "sisi",
       "larik": "ActualValue.TotalShareGrossNP",
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
        553145,
        557859
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "blok",
       "at": 603389,
       "judul": "",
       "syarat": [],
       "anak": [
        {
         "t": "grid",
         "at": 620098,
         "prop": "TreatyIn.ActualValue.TotalShareDeductionNP",
         "dari": "sisi",
         "larik": "ActualValue.TotalShareDeductionNP",
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
          631669,
          636383
         ],
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        }
       ]
      },
      {
       "t": "grid",
       "at": 698620,
       "prop": "TreatyIn.ActualValue.TotalShareNetNP",
       "dari": "sisi",
       "larik": "ActualValue.TotalShareNetNP",
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
        710187,
        714901
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 762159,
       "prop": "TreatyIn.ActualValue.TotalSpreadedNetPremi",
       "dari": "sisi",
       "larik": "ActualValue.TotalSpreadedNetPremi",
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
        773735,
        778449
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 824553,
       "prop": "TreatyIn.ActualValue.TotalSpreadedNetPremiRI",
       "dari": "sisi",
       "larik": "ActualValue.TotalSpreadedNetPremiRI",
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
        836132,
        840846
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "tombol",
       "at": 890194,
       "label": "Show Facultative Share",
       "syarat": [
        "TreatyIn.FacultativeShare>0"
       ],
       "aksi": [
        {
         "aksi": "showHarness"
        }
       ]
      }
     ]
    }
   ]
  }
 ],
 "TreatyInActualSumary": [
  {
   "t": "blok",
   "at": 18926,
   "judul": "Premium Difference",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 36300,
     "prop": "TreatyIn.ValueDifference.EGNPI",
     "dari": "sisi",
     "larik": "ValueDifference.EGNPI",
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
      65243,
      70698,
      75143,
      80239,
      84685,
      89776
     ],
     "baca": [
      [
       "1=1"
      ],
      null,
      null,
      null,
      null,
      null
     ],
     "tombol": [
      null,
      null,
      null,
      null,
      null,
      null
     ],
     "tombolKepala": [
      null,
      null,
      null,
      null,
      null,
      null
     ],
     "pilihan": [
      {
       "sumber": "associated"
      },
      null,
      null,
      null,
      null,
      null
     ],
     "aksiUbah": [
      null,
      null,
      null,
      null,
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInEGNPI!pyGridRowDetails",
     "rincian": "DetailEGNPI"
    }
   ]
  },
  {
   "t": "grid",
   "at": 150382,
   "prop": "TreatyIn.ValueDifference.TotalEgnpiAmountNP",
   "dari": "sisi",
   "larik": "ValueDifference.TotalEgnpiAmountNP",
   "syarat": [],
   "kolom": [
    "Total Premium Difference",
    "Value"
   ],
   "kunci": [
    "Currency",
    "Value"
   ],
   "lebar": [
    195,
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
    162374,
    167330
   ],
   "baca": [
    null,
    null
   ],
   "tombol": [
    null,
    null
   ],
   "tombolKepala": [
    null,
    null
   ],
   "pilihan": [
    null,
    null
   ],
   "aksiUbah": [
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "blok",
   "at": 218150,
   "judul": "Limits Difference",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 235528,
     "prop": "TreatyIn.ValueDifference.Limits",
     "dari": "sisi",
     "larik": "ValueDifference.Limits",
     "syarat": [],
     "kolom": [
      "Layers",
      "",
      "",
      "",
      ""
     ],
     "kunci": [
      "LayerType",
      "Layer",
      "pyTemplateRichTextEditor",
      "LayerPartType",
      "LayerPart"
     ],
     "lebar": [
      101,
      101,
      101,
      101,
      104
     ],
     "desimal": [
      null,
      null,
      null,
      null,
      null
     ],
     "format": [
      "pxTextInput",
      "",
      "pxTextInput",
      "pxTextInput",
      ""
     ],
     "syaratSel": [
      null,
      null,
      null,
      null,
      null
     ],
     "atSel": [
      260154,
      266085,
      270530,
      276483,
      282323
     ],
     "baca": [
      null,
      null,
      "selalu",
      null,
      null
     ],
     "tombol": [
      null,
      null,
      null,
      null,
      null
     ],
     "tombolKepala": [
      null,
      null,
      null,
      null,
      null
     ],
     "pilihan": [
      null,
      null,
      null,
      null,
      null
     ],
     "aksiUbah": [
      null,
      null,
      null,
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails",
     "rincian": "Layers"
    }
   ]
  },
  {
   "t": "blok",
   "at": 315504,
   "judul": "Summary of Limit",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 341906,
     "prop": "TreatyIn.ValueDifference.LimitSummaryList",
     "dari": "sisi",
     "larik": "ValueDifference.LimitSummaryList",
     "syarat": [],
     "kolom": [
      "Note",
      "MDP (IDR)",
      "MDP (USD)",
      "Agregate Limit (IDR)",
      "Agregate Limit (USD)",
      "Deductible (IDR)",
      "Deductible (USD)"
     ],
     "kunci": [
      "Note",
      "MDP",
      "MDP2",
      "AggregateLimit",
      "AggregateLimit2",
      "Deductible",
      "Deductible2"
     ],
     "lebar": [
      252,
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
      null
     ],
     "format": [
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
      null
     ],
     "atSel": [
      374278,
      378720,
      383168,
      387617,
      392076,
      396537,
      400993
     ],
     "baca": [
      null,
      null,
      null,
      null,
      null,
      null,
      null
     ],
     "tombol": [
      null,
      null,
      null,
      null,
      null,
      null,
      null
     ],
     "tombolKepala": [
      null,
      null,
      null,
      null,
      null,
      null,
      null
     ],
     "pilihan": [
      null,
      null,
      null,
      null,
      null,
      null,
      null
     ],
     "aksiUbah": [
      null,
      null,
      null,
      null,
      null,
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-LimitSummaryList!pyGridModalTemplate"
    }
   ]
  },
  {
   "t": "blok",
   "at": 441087,
   "judul": "Total All Layers",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 467490,
     "prop": "TreatyIn.ValueDifference.TotalLimitPremiEarnNP",
     "dari": "sisi",
     "larik": "ValueDifference.TotalLimitPremiEarnNP",
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
      479486,
      484443
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 532550,
     "prop": "TreatyIn.ValueDifference.TotalLimitMDPNP",
     "dari": "sisi",
     "larik": "ValueDifference.TotalLimitMDPNP",
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
      544529,
      549486
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "teks",
     "at": 611244,
     "teks": "Total ROL",
     "syarat": []
    },
    {
     "t": "medan",
     "at": 621979,
     "label": "",
     "dari": "sisi",
     "kunci": "ValueDifference.TotalLimitsROL",
     "format": "pxNumber",
     "desimal": 2,
     "syarat": [],
     "baca": "selalu"
    }
   ]
  },
  {
   "t": "blok",
   "at": 672669,
   "judul": "RNM Share",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 681775,
     "judul": "",
     "syarat": [],
     "anak": [
      {
       "t": "grid",
       "at": 699132,
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
        731125,
        737059,
        741502,
        745902,
        751741,
        756190,
        762079,
        766604,
        772493
       ],
       "baca": [
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
       "tombol": [
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
       "tombolKepala": [
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
       "pilihan": [
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
       "aksiUbah": [
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
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridRowDetails",
       "rincian": "Share"
      }
     ]
    },
    {
     "t": "blok",
     "at": 807895,
     "judul": "Summarry of RNM Share",
     "syarat": [],
     "anak": [
      {
       "t": "grid",
       "at": 825285,
       "prop": "TreatyIn.ValueDifference.LimitShareSummaryList",
       "dari": "sisi",
       "larik": "ValueDifference.LimitShareSummaryList",
       "syarat": [],
       "kolom": [
        "Note",
        "MDP (IDR)",
        "MDP (USD)",
        "Deduction (IDR)",
        "Deduction (USD)",
        "Net Premi (IDR)",
        "Net Premi (USD)"
       ],
       "kunci": [
        "Note",
        "MDP",
        "MDP2",
        "Deductible",
        "Deductible2",
        "NetPremi",
        "NetPremi2"
       ],
       "lebar": [
        293,
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
        null
       ],
       "format": [
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
        null
       ],
       "atSel": [
        857649,
        862091,
        866539,
        870988,
        875444,
        879901,
        884354
       ],
       "baca": [
        null,
        null,
        null,
        null,
        null,
        null,
        null
       ],
       "tombol": [
        null,
        null,
        null,
        null,
        null,
        null,
        null
       ],
       "tombolKepala": [
        null,
        null,
        null,
        null,
        null,
        null,
        null
       ],
       "pilihan": [
        null,
        null,
        null,
        null,
        null,
        null,
        null
       ],
       "aksiUbah": [
        null,
        null,
        null,
        null,
        null,
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-LimitSummaryList!pyGridModalTemplate"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 930411,
   "judul": "Total All Layers RNM Share",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 965823,
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
      977820,
      982777
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "blok",
     "at": 1028911,
     "judul": "",
     "syarat": [],
     "anak": [
      {
       "t": "grid",
       "at": 1046204,
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
        1058195,
        1063152
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    },
    {
     "t": "grid",
     "at": 1126618,
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
      1138605,
      1143562
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 1191649,
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
      1203645,
      1208602
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 1255535,
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
      1267534,
      1272491
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "tombol",
     "at": 1323059,
     "label": "Update Total",
     "syarat": [
      "TreatyIn.ViewState !='1'"
     ],
     "aksi": [
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInNPSetTotalActual",
       "param": {
        "type": "limits"
       }
      },
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInSummaryLimitActual"
      },
      {
       "aksi": "refresh",
       "aktivitas": "TreatyInSummaryLimitShareActual"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1374675,
   "judul": "",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1383814,
     "judul": "Other Treaty Retro",
     "syarat": [],
     "anak": [
      {
       "t": "blok",
       "at": 1392928,
       "judul": "",
       "syarat": [],
       "anak": [
        {
         "t": "grid",
         "at": 1410293,
         "prop": "TreatyIn.ValueDifference.FacultativeShareList",
         "dari": "sisi",
         "larik": "ValueDifference.FacultativeShareList",
         "syarat": [],
         "kolom": [
          "",
          "",
          "",
          "",
          "",
          "",
          "",
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
          "RnmLimitListDisplay(2).Currency",
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
          75,
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
          null
         ],
         "format": [
          "pxTextInput",
          "",
          "",
          "pxTextInput",
          "",
          "pxTextInput",
          "pxTextInput",
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
          null
         ],
         "atSel": [
          1447528,
          1453462,
          1457905,
          1462305,
          1468144,
          1472593,
          1478481,
          1484369,
          1490258,
          1494783,
          1500673
         ],
         "baca": [
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
         "tombol": [
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
         "tombolKepala": [
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
         "pilihan": [
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
         "aksiUbah": [
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
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridRowDetails",
         "rincian": "ShareRetro"
        }
       ]
      },
      {
       "t": "blok",
       "at": 1537301,
       "judul": "Summarry of Other Treaty Retro",
       "syarat": [],
       "anak": [
        {
         "t": "grid",
         "at": 1554700,
         "prop": "TreatyIn.ValueDifference.LimitFacShareSummaryList",
         "dari": "sisi",
         "larik": "ValueDifference.LimitFacShareSummaryList",
         "syarat": [],
         "kolom": [
          "Note",
          "MDP (IDR)",
          "MDP (USD)",
          "Deduction (IDR)",
          "Deduction (USD)",
          "Net Premi (IDR)",
          "Net Premi (USD)"
         ],
         "kunci": [
          "Note",
          "MDP",
          "MDP2",
          "Deductible",
          "Deductible2",
          "NetPremi",
          "NetPremi2"
         ],
         "lebar": [
          293,
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
          null
         ],
         "format": [
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
          null
         ],
         "atSel": [
          1587067,
          1591509,
          1595957,
          1600406,
          1604862,
          1609319,
          1613772
         ],
         "baca": [
          null,
          null,
          null,
          null,
          null,
          null,
          null
         ],
         "tombol": [
          null,
          null,
          null,
          null,
          null,
          null,
          null
         ],
         "tombolKepala": [
          null,
          null,
          null,
          null,
          null,
          null,
          null
         ],
         "pilihan": [
          null,
          null,
          null,
          null,
          null,
          null,
          null
         ],
         "aksiUbah": [
          null,
          null,
          null,
          null,
          null,
          null,
          null
         ],
         "templatBaris": "ASM-FW-GISFW-Data-LimitSummaryList!pyGridModalTemplate"
        }
       ]
      }
     ]
    },
    {
     "t": "blok",
     "at": 1653507,
     "judul": "Total All Layers Other Treaty Retro",
     "syarat": [],
     "anak": [
      {
       "t": "grid",
       "at": 1688928,
       "prop": "TreatyIn.ValueDifference.TotalFacShareGrossNP",
       "dari": "sisi",
       "larik": "ValueDifference.TotalFacShareGrossNP",
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
        1700928,
        1705885
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "blok",
       "at": 1752020,
       "judul": "",
       "syarat": [],
       "anak": [
        {
         "t": "grid",
         "at": 1769313,
         "prop": "TreatyIn.ValueDifference.TotalFacShareDeductionNP",
         "dari": "sisi",
         "larik": "ValueDifference.TotalFacShareDeductionNP",
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
          1781307,
          1786264
         ],
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
        }
       ]
      },
      {
       "t": "grid",
       "at": 1849710,
       "prop": "TreatyIn.ValueDifference.TotalFacShareNetNP",
       "dari": "sisi",
       "larik": "ValueDifference.TotalFacShareNetNP",
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
        1861700,
        1866657
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
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
  "baca": [
   "selalu",
   "selalu",
   "selalu",
   "selalu"
  ],
  "tombol": [
   null,
   null,
   null,
   null
  ],
  "tombolKepala": [
   null,
   null,
   null,
   null
  ],
  "pilihan": [
   {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrency_RD",
    "nilai": "Currency"
   },
   null,
   null,
   null
  ],
  "aksiUbah": [
   [
    {
     "aksi": "refresh"
    }
   ],
   null,
   null,
   null
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
   "Valid Until",
   ""
  ],
  "kunci": [
   "CurrencyID",
   "Conversion",
   "PeriodStart",
   "PeriodEnd",
   ""
  ],
  "lebar": [
   191,
   208,
   194,
   198,
   89
  ],
  "desimal": [
   null,
   null,
   null,
   null,
   null
  ],
  "format": [
   "pxDropdown",
   "pxTextInput",
   "pxDateTime",
   "pxDateTime",
   "pxButton"
  ],
  "syaratSel": [
   null,
   null,
   null,
   null,
   "TreatyIn.ViewState !='1'"
  ],
  "atSel": [
   439485,
   454068,
   460497,
   466983,
   473425
  ],
  "baca": [
   [
    "TreatyIn.ViewState = 1"
   ],
   [
    "TreatyIn.ViewState = 1",
    "TreatyIn.EDMMaterialType = 2"
   ],
   [
    "TreatyIn.ViewState = 1",
    "TreatyIn.EDMMaterialType = 2"
   ],
   [
    "TreatyIn.ViewState = 1",
    "TreatyIn.EDMMaterialType = 2"
   ],
   "selalu"
  ],
  "tombol": [
   null,
   null,
   null,
   null,
   {
    "t": "tombol",
    "at": 473425,
    "label": "Delete",
    "syarat": [
     "TreatyIn.ViewState !='1'"
    ],
    "aksi": [
     {
      "aksi": "deleteRow"
     }
    ],
    "nonaktif": [
     "TreatyIn.EDMMaterialType = 2"
    ]
   }
  ],
  "tombolKepala": [
   null,
   null,
   null,
   null,
   {
    "t": "tombol",
    "at": 429700,
    "label": "Add",
    "syarat": [
     "TreatyIn.ViewState !='1'"
    ],
    "aksi": [
     {
      "aksi": "refresh",
      "aktivitas": "TreatyInAddCurrency"
     }
    ],
    "ikon": "pyWorkActionsAddWork.png",
    "nonaktif": [
     "TreatyIn.EDMMaterialType = 2"
    ]
   }
  ],
  "pilihan": [
   {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrency_RD",
    "nilai": "ID",
    "tampil": "Currency"
   },
   null,
   null,
   null,
   null
  ],
  "aksiUbah": [
   [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "runActivity",
     "aktivitas": "SetCurrNameMasterTreaty_Act"
    },
    {
     "aksi": "refresh"
    }
   ],
   null,
   null,
   null,
   null
  ],
  "templatBaris": "ASM-FW-GISFW-Data-TreatyInCurrencyList!pyGridModalTemplate"
 }
}

// Section rincian baris (`expandPane`) — ikatan RELATIF atas halaman baris.
export const KERANGKA_RINCIAN: Readonly<Record<string, readonly ButirKerangka[]>> = {
 "AchievementCombine": [
  {
   "t": "grid",
   "at": 35121,
   "prop": ".Detail",
   "dari": "sisi",
   "larik": "Detail",
   "syarat": [],
   "kolom": [
    "Treaty Group",
    "Reins Type",
    "Gross Premium Before Claim in IDR",
    "Net Premium Before Claim in IDR",
    "Incured Claim in IDR",
    "Net Loss Ratio"
   ],
   "kunci": [
    "TreatyGroup",
    "TreatyType",
    "TotalAchPremium",
    "TotalAchNetPremium",
    "TotalAchIncured",
    "LossRatio"
   ],
   "lebar": [
    128,
    132,
    152,
    150,
    166,
    135
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
    "pxTextInput",
    "pxTextInput",
    "pxTextInput",
    "pxTextInput",
    "pxTextInput"
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
    64948,
    70451,
    75854,
    81645,
    86957,
    92266
   ],
   "baca": [
    "selalu",
    "selalu",
    "selalu",
    "selalu",
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    null,
    null,
    null,
    null,
    null
   ],
   "tombolKepala": [
    null,
    null,
    null,
    null,
    null,
    null
   ],
   "pilihan": [
    null,
    null,
    null,
    null,
    null,
    null
   ],
   "aksiUbah": [
    null,
    null,
    null,
    null,
    null,
    null
   ],
   "templatBaris": "Data-!pyGridModalTemplate"
  }
 ],
 "CoBList": [
  {
   "t": "medan",
   "at": 15558,
   "label": "Treaty Group",
   "dari": "sisi",
   "kunci": "TreatyGroup",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1",
    "TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseTreatyGroup_RD",
    "nilai": "TreatyGroupName"
   },
   "aksiUbah": [
    {
     "aksi": "refresh",
     "aktivitas": "TotalEgnpi",
     "param": {
      "idxLimit": ".ID"
     }
    },
    {
     "aksi": "refresh",
     "aktivitas": "TotalEgnpi",
     "param": {
      "idxLimit": ""
     }
    }
   ]
  },
  {
   "t": "grid",
   "at": 48813,
   "prop": ".ClassOfBusinessList",
   "dari": "sisi",
   "larik": "ClassOfBusinessList",
   "syarat": [],
   "kolom": [
    "Class of Business"
   ],
   "kunci": [
    "ClassOfBusiness"
   ],
   "lebar": [
    722
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
    65903
   ],
   "baca": [
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ]
   ],
   "tombol": [
    null
   ],
   "tombolKepala": [
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseTreatyBusinessWOType_RD",
     "nilai": "BIZNAME"
    }
   ],
   "aksiUbah": [
    [
     {
      "aksi": "runDataTransform",
      "transformasi": "SetCoB",
      "paramDT": {
       "businessID": ".ClassOfBusinessID",
       "business": ".ClassOfBusiness"
      }
     },
     {
      "aksi": "postValue"
     }
    ]
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitLayer!pyGridRowDetails"
  }
 ],
 "CoBListOldData": [
  {
   "t": "medan",
   "at": 15245,
   "label": "Treaty Group",
   "dari": "sisi",
   "kunci": "TreatyGroup",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseTreatyGroup_RD",
    "nilai": "TreatyGroupName"
   },
   "aksiUbah": [
    {
     "aksi": "refresh",
     "aktivitas": "TotalEgnpi",
     "param": {
      "idxLimit": ".ID"
     }
    },
    {
     "aksi": "refresh",
     "aktivitas": "TotalEgnpi",
     "param": {
      "idxLimit": ""
     }
    }
   ]
  },
  {
   "t": "grid",
   "at": 48417,
   "prop": ".ClassOfBusinessList",
   "dari": "sisi",
   "larik": "ClassOfBusinessList",
   "syarat": [],
   "kolom": [
    "Class of Business"
   ],
   "kunci": [
    "ClassOfBusiness"
   ],
   "lebar": [
    719
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
    56336
   ],
   "baca": [
    "selalu"
   ],
   "tombol": [
    null
   ],
   "tombolKepala": [
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseTreatyBusinessWOType_RD",
     "nilai": "BIZNAME"
    }
   ],
   "aksiUbah": [
    [
     {
      "aksi": "runDataTransform",
      "transformasi": "SetCoB",
      "paramDT": {
       "businessID": ".ClassOfBusinessID",
       "business": ".ClassOfBusiness"
      }
     },
     {
      "aksi": "postValue"
     }
    ]
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitLayer!pyGridRowDetails"
  }
 ],
 "CoBListReadOnly": [
  {
   "t": "medan",
   "at": 15187,
   "label": "Treaty Group",
   "dari": "sisi",
   "kunci": "TreatyGroup",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseTreatyGroup_RD",
    "nilai": "TreatyGroupName"
   },
   "aksiUbah": [
    {
     "aksi": "refresh"
    }
   ]
  },
  {
   "t": "grid",
   "at": 45602,
   "prop": ".ClassOfBusinessList",
   "dari": "sisi",
   "larik": "ClassOfBusinessList",
   "syarat": [],
   "kolom": [
    "Class of Business"
   ],
   "kunci": [
    "ClassOfBusiness"
   ],
   "lebar": [
    196
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
    53521
   ],
   "baca": [
    "selalu"
   ],
   "tombol": [
    null
   ],
   "tombolKepala": [
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseTreatyBusinessWOType_RD",
     "nilai": "BIZNAME"
    }
   ],
   "aksiUbah": [
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitLayer!pyGridRowDetails"
  }
 ],
 "DetailEGNPI": [
  {
   "t": "medan",
   "at": 16166,
   "label": "Treaty Group",
   "dari": "sisi",
   "kunci": "TreatyGroup",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1",
    "TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseTreatyGroup_RD",
    "nilai": "TreatyGroupName"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 26530,
   "label": "As At",
   "dari": "sisi",
   "kunci": "AsDate",
   "format": "pxDateTime",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1",
    "TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 52955,
   "label": "Amount",
   "dari": "sisi",
   "kunci": "Currency",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh",
     "aktivitas": "SetAmountConversion"
    }
   ]
  },
  {
   "t": "medan",
   "at": 90365,
   "label": "",
   "dari": "sisi",
   "kunci": "Amount",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh",
     "aktivitas": "SetAmountConversion"
    }
   ]
  },
  {
   "t": "medan",
   "at": 102409,
   "label": "",
   "dari": "sisi",
   "kunci": "AmountIDR",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 124120,
   "label": "Proportion %",
   "dari": "sisi",
   "kunci": "Proportion",
   "format": "pxTextInput",
   "desimal": 10,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 132880,
   "label": "Text Area",
   "dari": "sisi",
   "kunci": "Note",
   "format": "pxTextArea",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  }
 ],
 "DetailEGNPIOldData": [
  {
   "t": "medan",
   "at": 15294,
   "label": "Treaty Group",
   "dari": "sisi",
   "kunci": "TreatyGroup",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseTreatyGroup_RD",
    "nilai": "TreatyGroupName"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 25524,
   "label": "As At",
   "dari": "sisi",
   "kunci": "AsDate",
   "format": "pxDateTime",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 51795,
   "label": "Amount",
   "dari": "sisi",
   "kunci": "Currency",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh",
     "aktivitas": "SetAmountConversion"
    }
   ]
  },
  {
   "t": "medan",
   "at": 89119,
   "label": "",
   "dari": "sisi",
   "kunci": "Amount",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh",
     "aktivitas": "SetAmountConversion"
    }
   ]
  },
  {
   "t": "medan",
   "at": 101077,
   "label": "",
   "dari": "sisi",
   "kunci": "AmountIDR",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 122788,
   "label": "Proportion %",
   "dari": "sisi",
   "kunci": "Proportion",
   "format": "pxTextInput",
   "desimal": 10,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 131548,
   "label": "Text Area",
   "dari": "sisi",
   "kunci": "Note",
   "format": "pxTextArea",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  }
 ],
 "DetailLimits": [
  {
   "t": "medan",
   "at": 33452,
   "label": "Treaty Group",
   "dari": "sisi",
   "kunci": "TreatyGroupID",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1",
    "TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseTreatyGroup_RD",
    "nilai": "ID",
    "tampil": "TreatyGroupName"
   },
   "aksiUbah": [
    {
     "aksi": "refresh",
     "aktivitas": "SetTreatyGroupName_Act"
    },
    {
     "aksi": "refresh"
    },
    {
     "aksi": "refresh",
     "aktivitas": "FetchQSfromMaster",
     "transformasi": "SetDetailsID",
     "param": {
      "ParentReinsTypeID": ""
     },
     "syarat": ".TreatyType = 'QUOTA SHARE'"
    },
    {
     "aksi": "refresh",
     "aktivitas": "LimitCalculation",
     "param": {
      "kindoftreaty": "surplus",
      "add": "",
      "autocalculate": "true"
     },
     "syarat": ".TreatyType = 'SURPLUS'"
    }
   ]
  },
  {
   "t": "grid",
   "at": 86887,
   "prop": ".COBList",
   "dari": "sisi",
   "larik": "COBList",
   "syarat": [],
   "kolom": [
    "Class of Business"
   ],
   "kunci": [
    "ClassOfBusiness"
   ],
   "lebar": [
    272
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
    104472
   ],
   "baca": [
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ]
   ],
   "tombol": [
    null
   ],
   "tombolKepala": [
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseTreatyBusinessWOType_RD",
     "nilai": "BIZNAME"
    }
   ],
   "aksiUbah": [
    [
     {
      "aksi": "runDataTransform",
      "transformasi": "SetCoBID",
      "paramDT": {
       "id": ".ClassOfBusinessID"
      }
     },
     {
      "aksi": "postValue"
     }
    ]
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails"
  },
  {
   "t": "blok",
   "at": 170234,
   "judul": "",
   "syarat": [
    ".TreatyType = 'QUOTA SHARE'"
   ],
   "anak": [
    {
     "t": "medan",
     "at": 185962,
     "label": "QS %",
     "dari": "sisi",
     "kunci": "QSPct",
     "format": "pxNumber",
     "desimal": 2,
     "syarat": [],
     "baca": [
      "TreatyIn.ViewState = 1",
      "TreatyIn.EDMMaterialType = 2"
     ],
     "aksiUbah": [
      {
       "aksi": "postValue"
      },
      {
       "aksi": "refresh",
       "aktivitas": "LimitCalculation",
       "param": {
        "kindoftreaty": "qs",
        "add": "",
        "autocalculate": "true"
       }
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 212267,
   "judul": "",
   "syarat": [
    ".TreatyType = 'SURPLUS' || .TreatyType = '2ND SURPLUS' || .TreatyType = '3RD SURPLUS' || .TreatyType = 'SPECIAL SURPLUS'"
   ],
   "anak": [
    {
     "t": "medan",
     "at": 228088,
     "label": "Lines",
     "dari": "sisi",
     "kunci": "Surplus",
     "format": "pxNumber",
     "desimal": 2,
     "syarat": [],
     "baca": [
      "TreatyIn.ViewState = 1",
      "TreatyIn.EDMMaterialType = 2"
     ],
     "aksiUbah": [
      {
       "aksi": "postValue"
      },
      {
       "aksi": "refresh",
       "aktivitas": "LimitCalculation",
       "param": {
        "kindoftreaty": "surplus",
        "add": "",
        "autocalculate": "true"
       }
      }
     ]
    }
   ]
  },
  {
   "t": "teks",
   "at": 283497,
   "teks": "100% Limit",
   "syarat": []
  },
  {
   "t": "teks",
   "at": 297544,
   "teks": "100",
   "syarat": [
    ".TreatyType = 'QUOTA SHARE'"
   ]
  },
  {
   "t": "teks",
   "at": 302320,
   "teks": "%",
   "syarat": [
    ".TreatyType = 'QUOTA SHARE'"
   ]
  },
  {
   "t": "grid",
   "at": 345241,
   "prop": ".IOOLimitList",
   "dari": "sisi",
   "larik": "IOOLimitList",
   "syarat": [],
   "kolom": [
    "",
    "",
    "",
    ""
   ],
   "kunci": [
    "Currency",
    "Value",
    "",
    "Layer"
   ],
   "lebar": [
    194,
    361,
    102,
    128
   ],
   "desimal": [
    null,
    2,
    null,
    null
   ],
   "format": [
    "pxDropdown",
    "pxNumber",
    "pxButton",
    "pxCheckbox"
   ],
   "syaratSel": [
    null,
    null,
    "TreatyIn.ViewState !='1'",
    ".Note = 'QUOTA SHARE'"
   ],
   "atSel": [
    369627,
    385320,
    395493,
    407837
   ],
   "baca": [
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ],
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ],
    "selalu",
    [
     "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
    ]
   ],
   "tombol": [
    null,
    null,
    {
     "t": "tombol",
     "at": 395493,
     "label": "Remove",
     "syarat": [
      "TreatyIn.ViewState !='1'"
     ],
     "aksi": [
      {
       "aksi": "deleteRow"
      },
      {
       "aksi": "refresh",
       "aktivitas": "LimitCalculation",
       "param": {
        "kindoftreaty": "qs",
        "add": "",
        "autocalculate": ".Layer"
       },
       "syarat": ".Note = 'QUOTA SHARE'"
      }
     ],
     "nonaktif": [
      "TreatyIn.EDMMaterialType = 2"
     ]
    },
    null
   ],
   "tombolKepala": [
    null,
    null,
    {
     "t": "tombol",
     "at": 356966,
     "label": "Add",
     "syarat": [
      "TreatyIn.ViewState !='1'"
     ],
     "aksi": [
      {
       "aksi": "refresh",
       "transformasi": "AddLimitRetentionCession",
       "paramDT": {
        "type": "limit"
       }
      }
     ],
     "nonaktif": [
      "TreatyIn.EDMMaterialType = 2"
     ]
    },
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseCurrencyTreatyIn_RD",
     "nilai": "Currency",
     "tampil": "Currency"
    },
    null,
    null,
    null
   ],
   "aksiUbah": [
    [
     {
      "aksi": "postValue"
     },
     {
      "aksi": "runActivity",
      "aktivitas": "SetCurrName_Act"
     },
     {
      "aksi": "refresh",
      "aktivitas": "LimitCalculation",
      "param": {
       "kindoftreaty": "QS",
       "add": "",
       "autocalculate": ".Layer"
      },
      "syarat": ".Note = 'QUOTA SHARE'"
     }
    ],
    [
     {
      "aksi": "refresh",
      "aktivitas": "LimitCalculation",
      "param": {
       "kindoftreaty": "qs",
       "add": "",
       "autocalculate": ".Layer"
      },
      "syarat": ".Note = 'QUOTA SHARE'"
     }
    ],
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "teks",
   "at": 460405,
   "teks": "Retention",
   "syarat": []
  },
  {
   "t": "medan",
   "at": 474451,
   "label": "",
   "dari": "sisi",
   "kunci": "RetentionPct",
   "format": "pxNumber",
   "desimal": null,
   "syarat": [
    ".TreatyType = 'QUOTA SHARE'"
   ],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh",
     "aktivitas": "LimitCalculation",
     "param": {
      "kindoftreaty": "qs",
      "ioo": ".IOOLimit"
     }
    }
   ]
  },
  {
   "t": "teks",
   "at": 481479,
   "teks": "%",
   "syarat": [
    ".TreatyType = 'QUOTA SHARE'"
   ]
  },
  {
   "t": "grid",
   "at": 524406,
   "prop": ".RetentionList",
   "dari": "sisi",
   "larik": "RetentionList",
   "syarat": [],
   "kolom": [
    "",
    "",
    "",
    ""
   ],
   "kunci": [
    "Currency",
    "Value",
    "",
    "Layer"
   ],
   "lebar": [
    196,
    363,
    104,
    128
   ],
   "desimal": [
    null,
    2,
    null,
    null
   ],
   "format": [
    "pxDropdown",
    "pxNumber",
    "pxButton",
    "pxCheckbox"
   ],
   "syaratSel": [
    null,
    null,
    "TreatyIn.ViewState !='1'",
    ".Note = 'SURPLUS' || .Note = '2ND SURPLUS' || .Note = '3RD SURPLUS' || .Note = 'SPECIAL SURPLUS'"
   ],
   "atSel": [
    548441,
    566006,
    578978,
    590393
   ],
   "baca": [
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ],
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ],
    "selalu",
    [
     "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
    ]
   ],
   "tombol": [
    null,
    null,
    {
     "t": "tombol",
     "at": 578978,
     "label": "Remove",
     "syarat": [
      "TreatyIn.ViewState !='1'"
     ],
     "aksi": [
      {
       "aksi": "deleteRow"
      },
      {
       "aksi": "refresh",
       "aktivitas": "LimitCalculation",
       "param": {
        "kindoftreaty": "surplus",
        "add": "man",
        "autocalculate": ".Layer"
       }
      }
     ],
     "nonaktif": [
      "TreatyIn.EDMMaterialType = 2"
     ]
    },
    null
   ],
   "tombolKepala": [
    null,
    null,
    {
     "t": "tombol",
     "at": 536135,
     "label": "Add",
     "syarat": [
      "TreatyIn.ViewState !='1'"
     ],
     "aksi": [
      {
       "aksi": "refresh",
       "transformasi": "AddLimitRetentionCession",
       "paramDT": {
        "type": "retention"
       }
      }
     ],
     "nonaktif": [
      "TreatyIn.EDMMaterialType = 2"
     ]
    },
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseCurrencyTreatyIn_RD",
     "nilai": "Currency",
     "tampil": "Currency"
    },
    null,
    null,
    null
   ],
   "aksiUbah": [
    [
     {
      "aksi": "postValue"
     },
     {
      "aksi": "runActivity",
      "aktivitas": "SetCurrName_Act"
     },
     {
      "aksi": "refresh",
      "aktivitas": "LimitCalculation",
      "param": {
       "kindoftreaty": "surplus",
       "add": "man",
       "autocalculate": ".Layer"
      },
      "syarat": ".Note = 'SURPLUS' || .Note = '2ND SURPLUS' || .Note = '3RD SURPLUS'"
     }
    ],
    [
     {
      "aksi": "refresh",
      "aktivitas": "LimitCalculation",
      "param": {
       "kindoftreaty": "surplus",
       "add": "man",
       "autocalculate": ".Layer"
      },
      "syarat": ".Note = 'SURPLUS' || .Note = '2ND SURPLUS' || .Note = '3RD SURPLUS' || .Note = 'SPECIAL SURPLUS'"
     }
    ],
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "teks",
   "at": 642766,
   "teks": "Cession to R/I",
   "syarat": []
  },
  {
   "t": "medan",
   "at": 656819,
   "label": "",
   "dari": "sisi",
   "kunci": "CessionPct",
   "format": "pxNumber",
   "desimal": null,
   "syarat": [
    ".TreatyType = 'QUOTA SHARE'"
   ],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh",
     "aktivitas": "LimitCalculation",
     "param": {
      "kindoftreaty": "qs",
      "ioo": ".IOOLimit"
     }
    }
   ]
  },
  {
   "t": "teks",
   "at": 663861,
   "teks": "%",
   "syarat": [
    ".TreatyType = 'QUOTA SHARE'"
   ]
  },
  {
   "t": "grid",
   "at": 706791,
   "prop": ".CessionList",
   "dari": "sisi",
   "larik": "CessionList",
   "syarat": [],
   "kolom": [
    "",
    "",
    ""
   ],
   "kunci": [
    "Currency",
    "Value",
    ""
   ],
   "lebar": [
    193,
    352,
    102
   ],
   "desimal": [
    null,
    2,
    null
   ],
   "format": [
    "pxDropdown",
    "pxNumber",
    "pxButton"
   ],
   "syaratSel": [
    null,
    null,
    "TreatyIn.ViewState !='1'"
   ],
   "atSel": [
    728136,
    741842,
    747116
   ],
   "baca": [
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ],
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ],
    "selalu"
   ],
   "tombol": [
    null,
    null,
    {
     "t": "tombol",
     "at": 747116,
     "label": "Remove",
     "syarat": [
      "TreatyIn.ViewState !='1'"
     ],
     "aksi": [
      {
       "aksi": "deleteRow"
      }
     ],
     "nonaktif": [
      "TreatyIn.EDMMaterialType = 2"
     ]
    }
   ],
   "tombolKepala": [
    null,
    null,
    {
     "t": "tombol",
     "at": 718518,
     "label": "Add",
     "syarat": [
      "TreatyIn.ViewState !='1'"
     ],
     "aksi": [
      {
       "aksi": "refresh",
       "transformasi": "AddLimitRetentionCession",
       "paramDT": {
        "type": "cession"
       }
      }
     ],
     "nonaktif": [
      "TreatyIn.EDMMaterialType = 2"
     ]
    }
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseCurrencyTreatyIn_RD",
     "nilai": "Currency",
     "tampil": "Currency"
    },
    null,
    null
   ],
   "aksiUbah": [
    [
     {
      "aksi": "postValue"
     },
     {
      "aksi": "runActivity",
      "aktivitas": "SetCurrName_Act"
     },
     {
      "aksi": "refresh"
     }
    ],
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "blok",
   "at": 788648,
   "judul": "Event Limits",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 798616,
     "judul": "Event Limits",
     "syarat": [],
     "anak": [
      {
       "t": "medan",
       "at": 828145,
       "label": "RSMD Limit",
       "dari": "sisi",
       "kunci": "CurrencyRSMD",
       "format": "pxDropdown",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = 1",
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ],
       "pilihan": {
        "sumber": "reportdefinition",
        "rd": "BrowseCurrencyTreatyIn_RD",
        "nilai": "Currency",
        "tampil": "Currency"
       },
       "aksiUbah": [
        {
         "aksi": "postValue"
        }
       ]
      },
      {
       "t": "medan",
       "at": 853160,
       "label": "",
       "dari": "sisi",
       "kunci": "RSMDLimit",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = 1",
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ]
      },
      {
       "t": "medan",
       "at": 884088,
       "label": "Earthquake Limit",
       "dari": "sisi",
       "kunci": "CurrencyEarthquake",
       "format": "pxDropdown",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = 1",
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ],
       "pilihan": {
        "sumber": "reportdefinition",
        "rd": "BrowseCurrencyTreatyIn_RD",
        "nilai": "Currency",
        "tampil": "Currency"
       },
       "aksiUbah": [
        {
         "aksi": "postValue"
        }
       ]
      },
      {
       "t": "medan",
       "at": 909906,
       "label": "",
       "dari": "sisi",
       "kunci": "Earthquake",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ]
      },
      {
       "t": "medan",
       "at": 940793,
       "label": "Flood Limit (Jabodetabek)",
       "dari": "sisi",
       "kunci": "CurrencyFloodJab",
       "format": "pxDropdown",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = 1",
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ],
       "pilihan": {
        "sumber": "reportdefinition",
        "rd": "BrowseCurrencyTreatyIn_RD",
        "nilai": "Currency",
        "tampil": "Currency"
       },
       "aksiUbah": [
        {
         "aksi": "postValue"
        }
       ]
      },
      {
       "t": "medan",
       "at": 966619,
       "label": "",
       "dari": "sisi",
       "kunci": "FloodJab",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ]
      },
      {
       "t": "medan",
       "at": 997509,
       "label": "Flood Limit (Nationwide)",
       "dari": "sisi",
       "kunci": "CurrencyFloodNat",
       "format": "pxDropdown",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = 1",
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ],
       "pilihan": {
        "sumber": "reportdefinition",
        "rd": "BrowseCurrencyTreatyIn_RD",
        "nilai": "Currency",
        "tampil": "Currency"
       },
       "aksiUbah": [
        {
         "aksi": "postValue"
        }
       ]
      },
      {
       "t": "medan",
       "at": 1023334,
       "label": "",
       "dari": "sisi",
       "kunci": "FloodNation",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ]
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1054344,
   "judul": "Deduction In A",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1064748,
     "judul": "Deduction In A",
     "syarat": [],
     "anak": [
      {
       "t": "medan",
       "at": 1076343,
       "label": "% OGR",
       "dari": "sisi",
       "kunci": "RIOGR",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ]
      },
      {
       "t": "medan",
       "at": 1082313,
       "label": "% ONR",
       "dari": "sisi",
       "kunci": "RIONR",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ]
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1100591,
   "judul": "Deduction",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1110557,
     "judul": "Deduction Details",
     "syarat": [],
     "anak": [
      {
       "t": "grid",
       "at": 1127846,
       "prop": ".DeductionList",
       "dari": "sisi",
       "larik": "DeductionList",
       "syarat": [],
       "kolom": [
        "Description",
        "Currency",
        "Deduction",
        "or",
        "Deduction %",
        ""
       ],
       "kunci": [
        "Comment",
        "CurrencyID",
        "Deduction",
        "",
        "DeductionPct",
        ""
       ],
       "lebar": [
        206,
        197,
        201,
        31,
        213,
        116
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
        "pxTextInput",
        "pxDropdown",
        "pxNumber",
        "",
        "pxNumber",
        "pxButton"
       ],
       "syaratSel": [
        null,
        null,
        null,
        null,
        null,
        "TreatyIn.ViewState !='1'"
       ],
       "atSel": [
        1161683,
        1167810,
        1182252,
        1191285,
        1193782,
        1202685
       ],
       "baca": [
        [
         "TreatyIn.ViewState = 1",
         "TreatyIn.EDMMaterialType = 2"
        ],
        [
         "TreatyIn.ViewState = 1",
         "TreatyIn.EDMMaterialType = 2"
        ],
        [
         "TreatyIn.ViewState = 1",
         "TreatyIn.EDMMaterialType = 2"
        ],
        null,
        [
         "TreatyIn.ViewState = 1",
         "TreatyIn.EDMMaterialType = 2"
        ],
        "selalu"
       ],
       "tombol": [
        null,
        null,
        null,
        null,
        null,
        {
         "t": "tombol",
         "at": 1202685,
         "label": "Remove",
         "syarat": [
          "TreatyIn.ViewState !='1'"
         ],
         "aksi": [
          {
           "aksi": "deleteRow"
          },
          {
           "aksi": "refresh",
           "aktivitas": "CalculateDeduction",
           "param": {
            "sts": "",
            "index": ".pxListSubscript"
           }
          }
         ],
         "nonaktif": [
          "TreatyIn.EDMMaterialType = 2"
         ]
        }
       ],
       "tombolKepala": [
        null,
        null,
        null,
        null,
        null,
        {
         "t": "tombol",
         "at": 1152149,
         "label": "Add",
         "syarat": [
          "TreatyIn.ViewState !='1'"
         ],
         "aksi": [
          {
           "aksi": "refresh",
           "aktivitas": "AddDeduction"
          }
         ],
         "nonaktif": [
          "TreatyIn.EDMMaterialType = 2"
         ]
        }
       ],
       "pilihan": [
        null,
        {
         "sumber": "reportdefinition",
         "rd": "BrowseCurrency_RD",
         "nilai": "ID",
         "tampil": "Currency"
        },
        null,
        null,
        null,
        null
       ],
       "aksiUbah": [
        null,
        [
         {
          "aksi": "postValue"
         },
         {
          "aksi": "runActivity",
          "aktivitas": "SetCurrName_Act"
         },
         {
          "aksi": "refresh",
          "aktivitas": "CalculateDeduction",
          "param": {
           "sts": "",
           "index": ".pxListSubscript"
          }
         }
        ],
        [
         {
          "aksi": "refresh",
          "aktivitas": "CalculateDeduction",
          "param": {
           "sts": "val",
           "index": ".pxListSubscript"
          }
         }
        ],
        null,
        [
         {
          "aksi": "refresh",
          "aktivitas": "CalculateDeduction",
          "param": {
           "sts": "pct",
           "index": ".pxListSubscript"
          }
         }
        ],
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInDeduction!pyGridModalTemplate"
      }
     ]
    },
    {
     "t": "grid",
     "at": 1259471,
     "prop": ".DeductionTotalList",
     "dari": "sisi",
     "larik": "DeductionTotalList",
     "syarat": [],
     "kolom": [
      "Deductions",
      ""
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      196,
      356
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
      1271748,
      1276886
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    }
   ]
  },
  {
   "t": "blok",
   "at": 1311728,
   "judul": "Reserve",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1322124,
     "judul": "Premium Reserve",
     "syarat": [],
     "anak": [
      {
       "t": "medan",
       "at": 1337879,
       "label": "% Premium Reserve",
       "dari": "sisi",
       "kunci": "PremiumReservePct",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState ='1' || TreatyIn.EDMMaterialType = 2"
       ],
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "PremiumReserveCalculate"
        }
       ]
      },
      {
       "t": "teks",
       "at": 1389912,
       "teks": "Premium Reserve",
       "syarat": []
      },
      {
       "t": "grid",
       "at": 1427153,
       "prop": ".ReserveList",
       "dari": "sisi",
       "larik": "ReserveList",
       "syarat": [],
       "kolom": [
        "",
        "",
        ""
       ],
       "kunci": [
        "Currency",
        "Value",
        ""
       ],
       "lebar": [
        198,
        363,
        101
       ],
       "desimal": [
        null,
        2,
        null
       ],
       "format": [
        "pxDropdown",
        "pxNumber",
        "pxButton"
       ],
       "syaratSel": [
        null,
        null,
        "TreatyIn.ViewState !='1'"
       ],
       "atSel": [
        1448396,
        1459946,
        1465633
       ],
       "baca": [
        [
         "TreatyIn.ViewState = 1",
         "TreatyIn.EDMMaterialType = 2"
        ],
        [
         "TreatyIn.ViewState = 1",
         "TreatyIn.EDMMaterialType = 2"
        ],
        "selalu"
       ],
       "tombol": [
        null,
        null,
        {
         "t": "tombol",
         "at": 1465633,
         "label": "Remove",
         "syarat": [
          "TreatyIn.ViewState !='1'"
         ],
         "aksi": [
          {
           "aksi": "deleteRow"
          }
         ],
         "nonaktif": [
          "TreatyIn.EDMMaterialType = 2"
         ]
        }
       ],
       "tombolKepala": [
        null,
        null,
        {
         "t": "tombol",
         "at": 1438825,
         "label": "Add",
         "syarat": [
          "TreatyIn.ViewState !='1'"
         ],
         "aksi": [
          {
           "aksi": "refresh",
           "aktivitas": "AddValue",
           "param": {
            "type": "\"reserve\""
           }
          }
         ],
         "nonaktif": [
          "TreatyIn.EDMMaterialType = 2"
         ]
        }
       ],
       "pilihan": [
        {
         "sumber": "reportdefinition",
         "rd": "BrowseCurrencyTreatyIn_RD",
         "nilai": "Currency",
         "tampil": "Currency"
        },
        null,
        null
       ],
       "aksiUbah": [
        [
         {
          "aksi": "postValue"
         },
         {
          "aksi": "runActivity",
          "aktivitas": "SetCurrName_Act"
         }
        ],
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1515314,
   "judul": "Experience Premium Refund",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1525729,
     "judul": "Experience Premium Refund",
     "syarat": [],
     "anak": [
      {
       "t": "medan",
       "at": 1537335,
       "label": "% Commision",
       "dari": "sisi",
       "kunci": "ProfitCommision",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = 1",
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ]
      },
      {
       "t": "medan",
       "at": 1543380,
       "label": "% ME",
       "dari": "sisi",
       "kunci": "ProfitME",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ]
      },
      {
       "t": "medan",
       "at": 1549354,
       "label": "YDCF",
       "dari": "sisi",
       "kunci": "ProfitYDCF",
       "format": "pxTextInput",
       "desimal": -999,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
       ]
      },
      {
       "t": "teks",
       "at": 1555372,
       "teks": "Note:",
       "syarat": []
      },
      {
       "t": "teks",
       "at": 1560523,
       "teks": "ME = Management Expense",
       "syarat": []
      },
      {
       "t": "teks",
       "at": 1565692,
       "teks": "YDCF = Years Deficit Carried Forward",
       "syarat": []
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1583182,
   "judul": "PLA",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1593574,
     "judul": "PLA",
     "syarat": [],
     "anak": [
      {
       "t": "teks",
       "at": 1636337,
       "teks": "PLA",
       "syarat": []
      },
      {
       "t": "grid",
       "at": 1673566,
       "prop": ".PLAList",
       "dari": "sisi",
       "larik": "PLAList",
       "syarat": [],
       "kolom": [
        "",
        "",
        ""
       ],
       "kunci": [
        "Currency",
        "Value",
        ""
       ],
       "lebar": [
        197,
        361,
        101
       ],
       "desimal": [
        null,
        2,
        null
       ],
       "format": [
        "pxDropdown",
        "pxNumber",
        "pxButton"
       ],
       "syaratSel": [
        null,
        null,
        "TreatyIn.ViewState !='1'"
       ],
       "atSel": [
        1694797,
        1705957,
        1711644
       ],
       "baca": [
        [
         "TreatyIn.ViewState = 1",
         "TreatyIn.EDMMaterialType = 2"
        ],
        [
         "TreatyIn.ViewState = 1",
         "TreatyIn.EDMMaterialType = 2"
        ],
        "selalu"
       ],
       "tombol": [
        null,
        null,
        {
         "t": "tombol",
         "at": 1711644,
         "label": "Remove",
         "syarat": [
          "TreatyIn.ViewState !='1'"
         ],
         "aksi": [
          {
           "aksi": "deleteRow"
          }
         ],
         "nonaktif": [
          "TreatyIn.EDMMaterialType = 2"
         ]
        }
       ],
       "tombolKepala": [
        null,
        null,
        {
         "t": "tombol",
         "at": 1685234,
         "label": "Add",
         "syarat": [
          "TreatyIn.ViewState !='1'"
         ],
         "aksi": [
          {
           "aksi": "refresh",
           "aktivitas": "AddValue",
           "param": {
            "type": "\"pla\""
           }
          }
         ],
         "nonaktif": [
          "TreatyIn.EDMMaterialType = 2"
         ]
        }
       ],
       "pilihan": [
        {
         "sumber": "reportdefinition",
         "rd": "BrowseCurrencyTreatyIn_RD",
         "nilai": "Currency",
         "tampil": "Currency"
        },
        null,
        null
       ],
       "aksiUbah": [
        [
         {
          "aksi": "postValue"
         },
         {
          "aksi": "runActivity",
          "aktivitas": "SetCurrName_Act"
         }
        ],
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1761325,
   "judul": "Cash Loss Limit",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1771730,
     "judul": "Cash Loss Limit",
     "syarat": [],
     "anak": [
      {
       "t": "teks",
       "at": 1814504,
       "teks": "Cash Loss Limit",
       "syarat": []
      },
      {
       "t": "grid",
       "at": 1851745,
       "prop": ".CashLossList",
       "dari": "sisi",
       "larik": "CashLossList",
       "syarat": [],
       "kolom": [
        "",
        "",
        ""
       ],
       "kunci": [
        "Currency",
        "Value",
        ""
       ],
       "lebar": [
        197,
        361,
        101
       ],
       "desimal": [
        null,
        2,
        null
       ],
       "format": [
        "pxDropdown",
        "pxNumber",
        "pxButton"
       ],
       "syaratSel": [
        null,
        null,
        "TreatyIn.ViewState !='1'"
       ],
       "atSel": [
        1872991,
        1884151,
        1889838
       ],
       "baca": [
        [
         "TreatyIn.ViewState = 1",
         "TreatyIn.EDMMaterialType = 2"
        ],
        [
         "TreatyIn.ViewState = 1",
         "TreatyIn.EDMMaterialType = 2"
        ],
        "selalu"
       ],
       "tombol": [
        null,
        null,
        {
         "t": "tombol",
         "at": 1889838,
         "label": "Remove",
         "syarat": [
          "TreatyIn.ViewState !='1'"
         ],
         "aksi": [
          {
           "aksi": "deleteRow"
          }
         ],
         "nonaktif": [
          "TreatyIn.EDMMaterialType = 2"
         ]
        }
       ],
       "tombolKepala": [
        null,
        null,
        {
         "t": "tombol",
         "at": 1863418,
         "label": "Add",
         "syarat": [
          "TreatyIn.ViewState !='1'"
         ],
         "aksi": [
          {
           "aksi": "refresh",
           "aktivitas": "AddValue",
           "param": {
            "type": "\"cashloss\""
           }
          }
         ],
         "nonaktif": [
          "TreatyIn.EDMMaterialType = 2"
         ]
        }
       ],
       "pilihan": [
        {
         "sumber": "reportdefinition",
         "rd": "BrowseCurrencyTreatyIn_RD",
         "nilai": "Currency",
         "tampil": "Currency"
        },
        null,
        null
       ],
       "aksiUbah": [
        [
         {
          "aksi": "postValue"
         },
         {
          "aksi": "runActivity",
          "aktivitas": "SetCurrName_Act"
         }
        ],
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1939519,
   "judul": "Claim Cooperation",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1949927,
     "judul": "Claim Cooperation",
     "syarat": [],
     "anak": [
      {
       "t": "teks",
       "at": 1992703,
       "teks": "Claim Cooperation",
       "syarat": []
      },
      {
       "t": "grid",
       "at": 2029946,
       "prop": ".ClaimCoopList",
       "dari": "sisi",
       "larik": "ClaimCoopList",
       "syarat": [],
       "kolom": [
        "",
        "",
        ""
       ],
       "kunci": [
        "Currency",
        "Value",
        ""
       ],
       "lebar": [
        197,
        361,
        101
       ],
       "desimal": [
        null,
        2,
        null
       ],
       "format": [
        "pxDropdown",
        "pxNumber",
        "pxButton"
       ],
       "syaratSel": [
        null,
        null,
        "TreatyIn.ViewState !='1'"
       ],
       "atSel": [
        2051195,
        2062355,
        2068042
       ],
       "baca": [
        [
         "TreatyIn.ViewState = 1",
         "TreatyIn.EDMMaterialType = 2"
        ],
        [
         "TreatyIn.ViewState = 1",
         "TreatyIn.EDMMaterialType = 2"
        ],
        "selalu"
       ],
       "tombol": [
        null,
        null,
        {
         "t": "tombol",
         "at": 2068042,
         "label": "Remove",
         "syarat": [
          "TreatyIn.ViewState !='1'"
         ],
         "aksi": [
          {
           "aksi": "deleteRow"
          }
         ],
         "nonaktif": [
          "TreatyIn.EDMMaterialType = 2"
         ]
        }
       ],
       "tombolKepala": [
        null,
        null,
        {
         "t": "tombol",
         "at": 2041620,
         "label": "Add",
         "syarat": [
          "TreatyIn.ViewState !='1'"
         ],
         "aksi": [
          {
           "aksi": "refresh",
           "aktivitas": "AddValue",
           "param": {
            "type": "\"claimcoop\""
           }
          }
         ],
         "nonaktif": [
          "TreatyIn.EDMMaterialType = 2"
         ]
        }
       ],
       "pilihan": [
        {
         "sumber": "reportdefinition",
         "rd": "BrowseCurrencyTreatyIn_RD",
         "nilai": "Currency",
         "tampil": "Currency"
        },
        null,
        null
       ],
       "aksiUbah": [
        [
         {
          "aksi": "postValue"
         },
         {
          "aksi": "runActivity",
          "aktivitas": "SetCurrName_Act"
         }
        ],
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 2117723,
   "judul": "LPC",
   "syarat": [],
   "anak": [
    {
     "t": "teks",
     "at": 2134247,
     "teks": "Loss Participation Clause (LPC)",
     "syarat": []
    },
    {
     "t": "medan",
     "at": 2147571,
     "label": "Lower Band",
     "dari": "sisi",
     "kunci": "LowerBand",
     "format": "pxNumber",
     "desimal": null,
     "syarat": [],
     "baca": [
      "TreatyIn.ViewState ='1'"
     ],
     "aksiUbah": [
      {
       "aksi": "postValue"
      }
     ]
    },
    {
     "t": "medan",
     "at": 2154546,
     "label": "Upper Band",
     "dari": "sisi",
     "kunci": "UpperBand",
     "format": "pxNumber",
     "desimal": null,
     "syarat": [],
     "baca": [
      "TreatyIn.ViewState ='1'"
     ],
     "aksiUbah": [
      {
       "aksi": "postValue"
      }
     ]
    },
    {
     "t": "medan",
     "at": 2161521,
     "label": "Reisured Participant",
     "dari": "sisi",
     "kunci": "ReisuredParticipant",
     "format": "pxNumber",
     "desimal": null,
     "syarat": [],
     "baca": [
      "TreatyIn.ViewState ='1'"
     ],
     "aksiUbah": [
      {
       "aksi": "postValue"
      }
     ]
    },
    {
     "t": "medan",
     "at": 2168472,
     "label": "Period (Month)",
     "dari": "sisi",
     "kunci": "Periode",
     "format": "pxNumber",
     "desimal": null,
     "syarat": [],
     "baca": [
      "TreatyIn.ViewState ='1'"
     ],
     "aksiUbah": [
      {
       "aksi": "postValue"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 2194135,
   "judul": "EPI",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 2204108,
     "judul": "EPI",
     "syarat": [],
     "anak": [
      {
       "t": "teks",
       "at": 2247000,
       "teks": "EPI",
       "syarat": []
      },
      {
       "t": "grid",
       "at": 2284378,
       "prop": ".EPIList",
       "dari": "sisi",
       "larik": "EPIList",
       "syarat": [],
       "kolom": [
        "",
        "",
        ""
       ],
       "kunci": [
        "Currency",
        "Value",
        ""
       ],
       "lebar": [
        197,
        361,
        101
       ],
       "desimal": [
        null,
        2,
        null
       ],
       "format": [
        "pxDropdown",
        "pxNumber",
        "pxButton"
       ],
       "syaratSel": [
        null,
        null,
        "TreatyIn.ViewState !='1'"
       ],
       "atSel": [
        2305690,
        2316863,
        2322563
       ],
       "baca": [
        [
         "TreatyIn.ViewState = 1",
         "TreatyIn.EDMMaterialType = 2"
        ],
        [
         "TreatyIn.ViewState = 1",
         "TreatyIn.EDMMaterialType = 2"
        ],
        "selalu"
       ],
       "tombol": [
        null,
        null,
        {
         "t": "tombol",
         "at": 2322563,
         "label": "Remove",
         "syarat": [
          "TreatyIn.ViewState !='1'"
         ],
         "aksi": [
          {
           "aksi": "deleteRow"
          }
         ],
         "nonaktif": [
          "TreatyIn.EDMMaterialType = 2"
         ]
        }
       ],
       "tombolKepala": [
        null,
        null,
        {
         "t": "tombol",
         "at": 2296101,
         "label": "Add",
         "syarat": [
          "TreatyIn.ViewState !='1'"
         ],
         "aksi": [
          {
           "aksi": "refresh",
           "aktivitas": "AddValue",
           "param": {
            "type": "\"epi\""
           }
          }
         ],
         "nonaktif": [
          "TreatyIn.EDMMaterialType = 2"
         ]
        }
       ],
       "pilihan": [
        {
         "sumber": "reportdefinition",
         "rd": "BrowseCurrencyTreatyIn_RD",
         "nilai": "Currency",
         "tampil": "Currency"
        },
        null,
        null
       ],
       "aksiUbah": [
        [
         {
          "aksi": "postValue"
         },
         {
          "aksi": "runActivity",
          "aktivitas": "SetCurrName_Act"
         }
        ],
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 2372416,
   "judul": "Achievement",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 2382269,
     "judul": "Achievement",
     "syarat": [],
     "anak": [
      {
       "t": "grid",
       "at": 2413255,
       "prop": ".AchievementLists",
       "dari": "sisi",
       "larik": "AchievementLists",
       "syarat": [],
       "kolom": [
        "Quarter",
        "Quarter Year",
        "Currency",
        "Premium",
        "R/I COMM",
        "Brokerage",
        "Net Premium Before Claim",
        "Paid Claim",
        "Cash Call Claim",
        "Outstanding Claim",
        "Incured Claim",
        "Total",
        "Net Loss Ratio"
       ],
       "kunci": [
        "Quarter",
        "QUARTERYEAR",
        "Currency",
        "PREMIUM",
        "RICOMM",
        "BROKERAGE",
        "NETPREMIUM",
        "PaidClaim",
        "CASHCALL",
        "OutstandingClaim",
        "IncuredClaim",
        "Total",
        "LossRatio"
       ],
       "lebar": [
        83,
        75,
        93,
        138,
        140,
        159,
        124,
        131,
        130,
        163,
        161,
        113,
        73
       ],
       "desimal": [
        null,
        null,
        null,
        2,
        2,
        2,
        2,
        2,
        2,
        2,
        null,
        null,
        null
       ],
       "format": [
        "pxNumber",
        "",
        "",
        "pxNumber",
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
        null,
        null,
        null,
        null,
        null
       ],
       "atSel": [
        2471653,
        2476638,
        2481054,
        2485443,
        2490471,
        2495498,
        2500528,
        2505559,
        2510589,
        2515618,
        2520656,
        2525614,
        2530533
       ],
       "baca": [
        "selalu",
        "selalu",
        "selalu",
        "selalu",
        "selalu",
        "selalu",
        "selalu",
        "selalu",
        "selalu",
        "selalu",
        "selalu",
        "selalu",
        "selalu"
       ],
       "tombol": [
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
       "tombolKepala": [
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
       "pilihan": [
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
       "aksiUbah": [
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
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInAchievement!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 2844427,
       "prop": ".CurrencyList",
       "dari": "sisi",
       "larik": "CurrencyList",
       "syarat": [],
       "kolom": [
        "Parameter",
        "Based on Gross",
        "Based on Nett"
       ],
       "kunci": [
        "Parameter",
        "AchievementPctGross",
        "LossRatioGross"
       ],
       "lebar": [
        135,
        183,
        183
       ],
       "desimal": [
        null,
        null,
        null
       ],
       "format": [
        "",
        "pxTextInput",
        "pxTextInput"
       ],
       "syaratSel": [
        null,
        null,
        null
       ],
       "atSel": [
        2861163,
        2865651,
        2871415
       ],
       "baca": [
        "selalu",
        "selalu",
        "selalu"
       ],
       "tombol": [
        null,
        null,
        null
       ],
       "tombolKepala": [
        null,
        null,
        null
       ],
       "pilihan": [
        null,
        null,
        null
       ],
       "aksiUbah": [
        null,
        null,
        null
       ],
       "templatBaris": "Data-!pyGridModalTemplate"
      },
      {
       "t": "medan",
       "at": 2929803,
       "label": "As At Quarter",
       "dari": "sesi",
       "kunci": "SearchData.CARI1",
       "format": "pxDropdown",
       "desimal": null,
       "syarat": [],
       "pilihan": {
        "sumber": "pageList",
        "halaman": "TempQuarter.pxResults"
       },
       "aksiUbah": [
        {
         "aksi": "refresh",
         "transformasi": "Reset_DT"
        }
       ]
      },
      {
       "t": "medan",
       "at": 2939138,
       "label": "Quarter Year",
       "dari": "sesi",
       "kunci": "SearchData.CARI2",
       "format": "pxDropdown",
       "desimal": null,
       "syarat": [],
       "pilihan": {
        "sumber": "pageList",
        "halaman": "TempQuarterYear.pxResults"
       },
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "GetAchievement",
         "param": {
          "search": "\"search\""
         }
        }
       ]
      },
      {
       "t": "teks",
       "at": 2961357,
       "teks": "EPI Must Not be Empty",
       "syarat": []
      },
      {
       "t": "tombol",
       "at": 2976066,
       "label": "Refresh",
       "syarat": [],
       "aksi": [
        {
         "aksi": "refresh",
         "aktivitas": "GetAchievement"
        }
       ]
      },
      {
       "t": "tombol",
       "at": 2982462,
       "label": "Generate Excel",
       "syarat": [
        "FlagExcel.CARI1=='1'"
       ],
       "aksi": [
        {
         "aksi": "showHarness",
         "aktivitas": "GenerateCSVTreaty"
        }
       ]
      },
      {
       "t": "tombol",
       "at": 2997863,
       "label": "Submit",
       "syarat": [
        "FlagExcel.CARI1=='1'"
       ],
       "aksi": [
        {
         "aksi": "refresh",
         "aktivitas": "InsertToLogAchievement"
        }
       ]
      }
     ]
    }
   ]
  }
 ],
 "DetailLimitsOldData": [
  {
   "t": "medan",
   "at": 26708,
   "label": "Treaty Group",
   "dari": "sisi",
   "kunci": "TreatyGroup",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseTreatyGroup_RD",
    "nilai": "TreatyGroupName"
   },
   "aksiUbah": [
    {
     "aksi": "runDataTransform",
     "transformasi": "SetTreatyGroupID",
     "paramDT": {
      "treatygroup": ".TreatyGroup",
      "treatygroupid": ".TreatyGroupID"
     }
    },
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh",
     "aktivitas": "FetchQSfromMaster",
     "transformasi": "SetDetailsID",
     "syarat": ".TreatyType = 'QUOTA SHARE'"
    },
    {
     "aksi": "refresh",
     "aktivitas": "LimitCalculation",
     "param": {
      "kindoftreaty": "surplus",
      "add": "",
      "autocalculate": "true"
     },
     "syarat": ".TreatyType = 'SURPLUS'"
    }
   ]
  },
  {
   "t": "grid",
   "at": 84825,
   "prop": ".COBList",
   "dari": "sisi",
   "larik": "COBList",
   "syarat": [],
   "kolom": [
    "Class of Business"
   ],
   "kunci": [
    "ClassOfBusiness"
   ],
   "lebar": [
    270
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
    92932
   ],
   "baca": [
    "selalu"
   ],
   "tombol": [
    null
   ],
   "tombolKepala": [
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseTreatyBusinessWOType_RD",
     "nilai": "BIZNAME"
    }
   ],
   "aksiUbah": [
    [
     {
      "aksi": "runDataTransform",
      "transformasi": "SetCoBID",
      "paramDT": {
       "id": ".ClassOfBusinessID"
      }
     },
     {
      "aksi": "postValue"
     }
    ]
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridRowDetails"
  },
  {
   "t": "blok",
   "at": 149173,
   "judul": "",
   "syarat": [
    ".TreatyType = 'QUOTA SHARE'"
   ],
   "anak": [
    {
     "t": "medan",
     "at": 164901,
     "label": "QS %",
     "dari": "sisi",
     "kunci": "QSPct",
     "format": "pxNumber",
     "desimal": 2,
     "syarat": [],
     "baca": "selalu",
     "aksiUbah": [
      {
       "aksi": "postValue"
      },
      {
       "aksi": "refresh",
       "aktivitas": "LimitCalculation",
       "param": {
        "kindoftreaty": "qs",
        "add": "",
        "autocalculate": "true"
       }
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 190970,
   "judul": "",
   "syarat": [
    ".TreatyType = 'SURPLUS' || .TreatyType = '2ND SURPLUS' || .TreatyType = '3RD SURPLUS'"
   ],
   "anak": [
    {
     "t": "medan",
     "at": 206758,
     "label": "Lines",
     "dari": "sisi",
     "kunci": "Surplus",
     "format": "pxNumber",
     "desimal": 2,
     "syarat": [],
     "baca": "selalu",
     "aksiUbah": [
      {
       "aksi": "postValue"
      },
      {
       "aksi": "refresh",
       "aktivitas": "LimitCalculation",
       "param": {
        "kindoftreaty": "surplus",
        "add": "",
        "autocalculate": "true"
       }
      }
     ]
    }
   ]
  },
  {
   "t": "teks",
   "at": 261931,
   "teks": "100% Limit",
   "syarat": []
  },
  {
   "t": "teks",
   "at": 275978,
   "teks": "100",
   "syarat": [
    ".TreatyType = 'QUOTA SHARE'"
   ]
  },
  {
   "t": "teks",
   "at": 280256,
   "teks": "%",
   "syarat": [
    ".TreatyType = 'QUOTA SHARE'"
   ]
  },
  {
   "t": "grid",
   "at": 323177,
   "prop": ".IOOLimitList",
   "dari": "sisi",
   "larik": "IOOLimitList",
   "syarat": [],
   "kolom": [
    "",
    ""
   ],
   "kunci": [
    "Currency",
    "Value"
   ],
   "lebar": [
    193,
    355
   ],
   "desimal": [
    null,
    2
   ],
   "format": [
    "pxAutoComplete",
    "pxNumber"
   ],
   "syaratSel": [
    null,
    null
   ],
   "atSel": [
    335260,
    351571
   ],
   "baca": [
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    null
   ],
   "tombolKepala": [
    null,
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseCurrencyTreatyIn_RD",
     "nilai": "Currency"
    },
    null
   ],
   "aksiUbah": [
    [
     {
      "aksi": "runDataTransform",
      "transformasi": "SetCurrencyID",
      "paramDT": {
       "Currency": ".Currency",
       "ID": ".CurrencyID"
      }
     },
     {
      "aksi": "postValue"
     },
     {
      "aksi": "refresh",
      "aktivitas": "LimitCalculation",
      "param": {
       "kindoftreaty": "QS",
       "add": "",
       "autocalculate": ".Layer"
      },
      "syarat": ".Note = 'QUOTA SHARE'"
     }
    ],
    [
     {
      "aksi": "refresh",
      "aktivitas": "LimitCalculation",
      "param": {
       "kindoftreaty": "qs",
       "add": "",
       "autocalculate": ".Layer"
      },
      "syarat": ".Note = 'QUOTA SHARE'"
     }
    ]
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "teks",
   "at": 407394,
   "teks": "Retention",
   "syarat": []
  },
  {
   "t": "medan",
   "at": 421440,
   "label": "",
   "dari": "sisi",
   "kunci": "RetentionPct",
   "format": "pxNumber",
   "desimal": null,
   "syarat": [
    ".TreatyType = 'QUOTA SHARE'"
   ],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh",
     "aktivitas": "LimitCalculation",
     "param": {
      "kindoftreaty": "qs",
      "ioo": ".IOOLimit"
     }
    }
   ]
  },
  {
   "t": "teks",
   "at": 428468,
   "teks": "%",
   "syarat": [
    ".TreatyType = 'QUOTA SHARE'"
   ]
  },
  {
   "t": "grid",
   "at": 471389,
   "prop": ".RetentionList",
   "dari": "sisi",
   "larik": "RetentionList",
   "syarat": [],
   "kolom": [
    "",
    ""
   ],
   "kunci": [
    "Currency",
    "Value"
   ],
   "lebar": [
    193,
    352
   ],
   "desimal": [
    null,
    2
   ],
   "format": [
    "pxAutoComplete",
    "pxNumber"
   ],
   "syaratSel": [
    null,
    null
   ],
   "atSel": [
    483477,
    501659
   ],
   "baca": [
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    null
   ],
   "tombolKepala": [
    null,
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseCurrencyTreatyIn_RD",
     "nilai": "Currency"
    },
    null
   ],
   "aksiUbah": [
    [
     {
      "aksi": "runDataTransform",
      "transformasi": "SetCurrencyID",
      "paramDT": {
       "Currency": ".Currency",
       "ID": ".CurrencyID"
      }
     },
     {
      "aksi": "postValue"
     },
     {
      "aksi": "refresh",
      "aktivitas": "LimitCalculation",
      "param": {
       "kindoftreaty": "surplus",
       "add": "man",
       "autocalculate": ".Layer"
      },
      "syarat": ".Note = 'SURPLUS' || .Note = '2ND SURPLUS' || .Note = '3RD SURPLUS'"
     }
    ],
    [
     {
      "aksi": "refresh",
      "aktivitas": "LimitCalculation",
      "param": {
       "kindoftreaty": "surplus",
       "add": "man",
       "autocalculate": ".Layer"
      },
      "syarat": ".Note = 'SURPLUS' || .Note = '2ND SURPLUS' || .Note = '3RD SURPLUS'"
     }
    ]
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "teks",
   "at": 559361,
   "teks": "Cession to R/I",
   "syarat": []
  },
  {
   "t": "medan",
   "at": 573414,
   "label": "",
   "dari": "sisi",
   "kunci": "CessionPct",
   "format": "pxNumber",
   "desimal": null,
   "syarat": [
    ".TreatyType = 'QUOTA SHARE'"
   ],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh",
     "aktivitas": "LimitCalculation",
     "param": {
      "kindoftreaty": "qs",
      "ioo": ".IOOLimit"
     }
    }
   ]
  },
  {
   "t": "teks",
   "at": 580456,
   "teks": "%",
   "syarat": [
    ".TreatyType = 'QUOTA SHARE'"
   ]
  },
  {
   "t": "grid",
   "at": 623386,
   "prop": ".CessionList",
   "dari": "sisi",
   "larik": "CessionList",
   "syarat": [],
   "kolom": [
    "",
    ""
   ],
   "kunci": [
    "Currency",
    "Value"
   ],
   "lebar": [
    193,
    350
   ],
   "desimal": [
    null,
    2
   ],
   "format": [
    "pxAutoComplete",
    "pxNumber"
   ],
   "syaratSel": [
    null,
    null
   ],
   "atSel": [
    635472,
    649796
   ],
   "baca": [
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    null
   ],
   "tombolKepala": [
    null,
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseCurrencyTreatyIn_RD",
     "nilai": "Currency"
    },
    null
   ],
   "aksiUbah": [
    [
     {
      "aksi": "runDataTransform",
      "transformasi": "SetCurrencyID",
      "paramDT": {
       "Currency": ".Currency",
       "ID": ".CurrencyID"
      }
     },
     {
      "aksi": "postValue"
     },
     {
      "aksi": "refresh"
     }
    ],
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "blok",
   "at": 698454,
   "judul": "Event Limits",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 708990,
     "judul": "Event Limits",
     "syarat": [],
     "anak": [
      {
       "t": "medan",
       "at": 738606,
       "label": "RSMD Limit",
       "dari": "sisi",
       "kunci": "CurrencyRSMD",
       "format": "pxAutoComplete",
       "desimal": null,
       "syarat": [],
       "baca": "selalu",
       "pilihan": {
        "sumber": "reportdefinition",
        "rd": "BrowseCurrencyTreatyIn_RD",
        "nilai": "Currency"
       },
       "aksiUbah": [
        {
         "aksi": "postValue"
        }
       ]
      },
      {
       "t": "medan",
       "at": 764339,
       "label": "",
       "dari": "sisi",
       "kunci": "RSMDLimit",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu"
      },
      {
       "t": "medan",
       "at": 795250,
       "label": "Earthquake Limit",
       "dari": "sisi",
       "kunci": "CurrencyEarthquake",
       "format": "pxAutoComplete",
       "desimal": null,
       "syarat": [],
       "baca": "selalu",
       "pilihan": {
        "sumber": "reportdefinition",
        "rd": "BrowseCurrencyTreatyIn_RD",
        "nilai": "Currency"
       },
       "aksiUbah": [
        {
         "aksi": "postValue"
        }
       ]
      },
      {
       "t": "medan",
       "at": 821643,
       "label": "",
       "dari": "sisi",
       "kunci": "Earthquake",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu"
      },
      {
       "t": "medan",
       "at": 852512,
       "label": "Flood Limit (Jabodetabek)",
       "dari": "sisi",
       "kunci": "CurrencyFloodJab",
       "format": "pxAutoComplete",
       "desimal": null,
       "syarat": [],
       "baca": "selalu",
       "pilihan": {
        "sumber": "reportdefinition",
        "rd": "BrowseCurrencyTreatyIn_RD",
        "nilai": "Currency"
       },
       "aksiUbah": [
        {
         "aksi": "postValue"
        }
       ]
      },
      {
       "t": "medan",
       "at": 878315,
       "label": "",
       "dari": "sisi",
       "kunci": "FloodJab",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu"
      },
      {
       "t": "medan",
       "at": 909187,
       "label": "Flood Limit (Nationwide)",
       "dari": "sisi",
       "kunci": "CurrencyFloodNat",
       "format": "pxAutoComplete",
       "desimal": null,
       "syarat": [],
       "baca": "selalu",
       "pilihan": {
        "sumber": "reportdefinition",
        "rd": "BrowseCurrencyTreatyIn_RD",
        "nilai": "Currency"
       },
       "aksiUbah": [
        {
         "aksi": "postValue"
        }
       ]
      },
      {
       "t": "medan",
       "at": 935011,
       "label": "",
       "dari": "sisi",
       "kunci": "FloodNation",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 969202,
   "judul": "RI Comm",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 979302,
     "judul": "RI Comm",
     "syarat": [],
     "anak": [
      {
       "t": "medan",
       "at": 990919,
       "label": "% OGR",
       "dari": "sisi",
       "kunci": "RIOGR",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu"
      },
      {
       "t": "medan",
       "at": 996816,
       "label": "% ONR",
       "dari": "sisi",
       "kunci": "RIONR",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1018275,
   "judul": "Deduction",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1028377,
     "judul": "Deduction Details",
     "syarat": [],
     "anak": [
      {
       "t": "grid",
       "at": 1045789,
       "prop": ".DeductionList",
       "dari": "sisi",
       "larik": "DeductionList",
       "syarat": [],
       "kolom": [
        "Description",
        "Currency",
        "Deduction",
        "or",
        "Deduction %"
       ],
       "kunci": [
        "Comment",
        "Currency",
        "Deduction",
        "",
        "DeductionPct"
       ],
       "lebar": [
        205,
        197,
        200,
        30,
        213
       ],
       "desimal": [
        null,
        null,
        2,
        null,
        2
       ],
       "format": [
        "pxTextInput",
        "pxAutoComplete",
        "pxNumber",
        "",
        "pxNumber"
       ],
       "syaratSel": [
        null,
        null,
        null,
        null,
        null
       ],
       "atSel": [
        1070545,
        1076636,
        1087733,
        1096645,
        1099155
       ],
       "baca": [
        "selalu",
        "selalu",
        "selalu",
        null,
        "selalu"
       ],
       "tombol": [
        null,
        null,
        null,
        null,
        null
       ],
       "tombolKepala": [
        null,
        null,
        null,
        null,
        null
       ],
       "pilihan": [
        null,
        {
         "sumber": "reportdefinition",
         "rd": "BrowseCurrency_RD",
         "nilai": "Currency"
        },
        null,
        null,
        null
       ],
       "aksiUbah": [
        null,
        [
         {
          "aksi": "refresh",
          "aktivitas": "CalculateDeduction",
          "param": {
           "sts": "",
           "index": ".pxListSubscript"
          }
         }
        ],
        [
         {
          "aksi": "refresh",
          "aktivitas": "CalculateDeduction",
          "param": {
           "sts": "val",
           "index": ".pxListSubscript"
          }
         }
        ],
        null,
        [
         {
          "aksi": "refresh",
          "aktivitas": "CalculateDeduction",
          "param": {
           "sts": "pct",
           "index": ".pxListSubscript"
          }
         }
        ]
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInDeduction!pyGridModalTemplate"
      }
     ]
    },
    {
     "t": "grid",
     "at": 1153026,
     "prop": ".DeductionTotalList",
     "dari": "sisi",
     "larik": "DeductionTotalList",
     "syarat": [],
     "kolom": [
      "Deductions",
      ""
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      196,
      356
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
      1165371,
      1170522
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    }
   ]
  },
  {
   "t": "blok",
   "at": 1208751,
   "judul": "Reserve",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1219284,
     "judul": "Premium Reserve",
     "syarat": [],
     "anak": [
      {
       "t": "medan",
       "at": 1235084,
       "label": "% Premium Reserve",
       "dari": "sisi",
       "kunci": "PremiumReservePct",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu",
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "PremiumReserveCalculate"
        }
       ]
      },
      {
       "t": "teks",
       "at": 1287303,
       "teks": "Premium Reserve",
       "syarat": []
      },
      {
       "t": "grid",
       "at": 1324693,
       "prop": ".ReserveList",
       "dari": "sisi",
       "larik": "ReserveList",
       "syarat": [],
       "kolom": [
        "",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        198,
        362
       ],
       "desimal": [
        null,
        2
       ],
       "format": [
        "pxAutoComplete",
        "pxNumber"
       ],
       "syaratSel": [
        null,
        null
       ],
       "atSel": [
        1336779,
        1344785
       ],
       "baca": [
        "selalu",
        "selalu"
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        {
         "sumber": "reportdefinition",
         "rd": "BrowseCurrencyTreatyIn_RD",
         "nilai": "Currency"
        },
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1395461,
   "judul": "Profit Commision",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1406002,
     "judul": "Profit Commision",
     "syarat": [],
     "anak": [
      {
       "t": "medan",
       "at": 1417628,
       "label": "% Commision",
       "dari": "sisi",
       "kunci": "ProfitCommision",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu"
      },
      {
       "t": "medan",
       "at": 1423600,
       "label": "% ME",
       "dari": "sisi",
       "kunci": "ProfitME",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu"
      },
      {
       "t": "medan",
       "at": 1429501,
       "label": "YDCF",
       "dari": "sisi",
       "kunci": "ProfitYDCF",
       "format": "pxTextInput",
       "desimal": -999,
       "syarat": [],
       "baca": "selalu"
      },
      {
       "t": "teks",
       "at": 1435446,
       "teks": "Note:",
       "syarat": []
      },
      {
       "t": "teks",
       "at": 1440610,
       "teks": "ME = Management Expense",
       "syarat": []
      },
      {
       "t": "teks",
       "at": 1445792,
       "teks": "YDCF = Years Deficit Carried Forward",
       "syarat": []
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1466549,
   "judul": "PLA",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1477078,
     "judul": "PLA",
     "syarat": [],
     "anak": [
      {
       "t": "teks",
       "at": 1519968,
       "teks": "PLA",
       "syarat": []
      },
      {
       "t": "grid",
       "at": 1557346,
       "prop": ".PLAList",
       "dari": "sisi",
       "larik": "PLAList",
       "syarat": [],
       "kolom": [
        "",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        197,
        360
       ],
       "desimal": [
        null,
        2
       ],
       "format": [
        "pxAutoComplete",
        "pxNumber"
       ],
       "syaratSel": [
        null,
        null
       ],
       "atSel": [
        1569428,
        1577465
       ],
       "baca": [
        "selalu",
        "selalu"
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        {
         "sumber": "reportdefinition",
         "rd": "BrowseCurrencyTreatyIn_RD",
         "nilai": "Currency"
        },
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1628141,
   "judul": "Cash Loss Limit",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1638681,
     "judul": "Cash Loss Limit",
     "syarat": [],
     "anak": [
      {
       "t": "teks",
       "at": 1681584,
       "teks": "Cash Loss Limit",
       "syarat": []
      },
      {
       "t": "grid",
       "at": 1718974,
       "prop": ".CashLossList",
       "dari": "sisi",
       "larik": "CashLossList",
       "syarat": [],
       "kolom": [
        "",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        195,
        359
       ],
       "desimal": [
        null,
        2
       ],
       "format": [
        "pxAutoComplete",
        "pxNumber"
       ],
       "syaratSel": [
        null,
        null
       ],
       "atSel": [
        1731061,
        1739067
       ],
       "baca": [
        "selalu",
        "selalu"
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        {
         "sumber": "reportdefinition",
         "rd": "BrowseCurrencyTreatyIn_RD",
         "nilai": "Currency"
        },
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1789499,
   "judul": "Claim Cooperation",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1800042,
     "judul": "Claim Cooperation",
     "syarat": [],
     "anak": [
      {
       "t": "teks",
       "at": 1842947,
       "teks": "Claim Cooperation",
       "syarat": []
      },
      {
       "t": "grid",
       "at": 1880339,
       "prop": ".ClaimCoopList",
       "dari": "sisi",
       "larik": "ClaimCoopList",
       "syarat": [],
       "kolom": [
        "",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        197,
        360
       ],
       "desimal": [
        null,
        2
       ],
       "format": [
        "pxAutoComplete",
        "pxNumber"
       ],
       "syaratSel": [
        null,
        null
       ],
       "atSel": [
        1892427,
        1900433
       ],
       "baca": [
        "selalu",
        "selalu"
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        {
         "sumber": "reportdefinition",
         "rd": "BrowseCurrencyTreatyIn_RD",
         "nilai": "Currency"
        },
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1950865,
   "judul": "EPI",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 1960961,
     "judul": "EPI",
     "syarat": [],
     "anak": [
      {
       "t": "teks",
       "at": 2003854,
       "teks": "EPI",
       "syarat": []
      },
      {
       "t": "grid",
       "at": 2041234,
       "prop": ".EPIList",
       "dari": "sisi",
       "larik": "EPIList",
       "syarat": [],
       "kolom": [
        "",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        197,
        360
       ],
       "desimal": [
        null,
        2
       ],
       "format": [
        "pxAutoComplete",
        "pxNumber"
       ],
       "syaratSel": [
        null,
        null
       ],
       "atSel": [
        2053317,
        2061323
       ],
       "baca": [
        "selalu",
        "selalu"
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        {
         "sumber": "reportdefinition",
         "rd": "BrowseCurrencyTreatyIn_RD",
         "nilai": "Currency"
        },
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 2112050,
   "judul": "Achievement",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 2122146,
     "judul": "Achievement",
     "syarat": [],
     "anak": [
      {
       "t": "grid",
       "at": 2154390,
       "prop": ".AchievementLists",
       "dari": "sisi",
       "larik": "AchievementLists",
       "syarat": [],
       "kolom": [
        "Quarter",
        "Quarter Year",
        "Currency",
        "Premium",
        "Paid Claim",
        "Outstanding Claim",
        "Incured Claim",
        "Total",
        "Loss Ratio"
       ],
       "kunci": [
        "Period",
        "",
        "Currency",
        "Value",
        "PaidClaim",
        "OutstandingClaim",
        "IncuredClaim",
        "Total",
        "LossRatio"
       ],
       "lebar": [
        152,
        100,
        154,
        252,
        255,
        266,
        256,
        266,
        128
       ],
       "desimal": [
        null,
        null,
        null,
        2,
        2,
        2,
        null,
        null,
        null
       ],
       "format": [
        "pxNumber",
        "",
        "",
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
        2195534,
        2200528,
        2203040,
        2207439,
        2212506,
        2217577,
        2222655,
        2227653,
        2232612
       ],
       "baca": [
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
       "tombol": [
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
       "tombolKepala": [
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
       "pilihan": [
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
       "aksiUbah": [
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
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInAchievement!pyGridModalTemplate"
      },
      {
       "t": "medan",
       "at": 2365309,
       "label": "Total Premium",
       "dari": "sisi",
       "kunci": "AchievementLists(1).Currency",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "1=1"
       ]
      },
      {
       "t": "medan",
       "at": 2371100,
       "label": "",
       "dari": "sisi",
       "kunci": "TotalAchievPremium",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "1=1"
       ]
      },
      {
       "t": "medan",
       "at": 2385256,
       "label": "Total Paid Claim",
       "dari": "sisi",
       "kunci": "AchievementLists(1).Currency",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "1=1"
       ]
      },
      {
       "t": "medan",
       "at": 2391007,
       "label": "",
       "dari": "sisi",
       "kunci": "TotalAchievPaid",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "1=1"
       ]
      },
      {
       "t": "medan",
       "at": 2405163,
       "label": "Total Outstanding",
       "dari": "sisi",
       "kunci": "AchievementLists(1).Currency",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "1=1"
       ]
      },
      {
       "t": "medan",
       "at": 2410916,
       "label": "",
       "dari": "sisi",
       "kunci": "TotalAchievOuts",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": [
        "1=1"
       ]
      },
      {
       "t": "medan",
       "at": 2429444,
       "label": "Achievement %",
       "dari": "sisi",
       "kunci": "AchievementPct",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu"
      },
      {
       "t": "medan",
       "at": 2434909,
       "label": "Mean Loss Ratio %",
       "dari": "sisi",
       "kunci": "LossRatio",
       "format": "pxTextInput",
       "desimal": null,
       "syarat": [],
       "baca": "selalu"
      },
      {
       "t": "medan",
       "at": 2440758,
       "label": "Achievement Date",
       "dari": "sisi",
       "kunci": "AchievementDate",
       "format": "pxDateTime",
       "desimal": null,
       "syarat": [],
       "baca": "selalu"
      },
      {
       "t": "teks",
       "at": 2446892,
       "teks": "EPI Must Not be Empty",
       "syarat": []
      },
      {
       "t": "tombol",
       "at": 2451270,
       "label": "Refresh",
       "syarat": [],
       "aksi": [
        {
         "aksi": "refresh",
         "aktivitas": "RefreshAchievement"
        }
       ]
      }
     ]
    }
   ]
  }
 ],
 "DetailShare": [
  {
   "t": "medan",
   "at": 29678,
   "label": "% RNM Share",
   "dari": "sisi",
   "kunci": "RNMShare",
   "format": "pxNumber",
   "desimal": 2,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState =='1' || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh",
     "aktivitas": "TreatyInPropshareDetail"
    }
   ]
  },
  {
   "t": "grid",
   "at": 80836,
   "prop": ".RNMShareList",
   "dari": "sisi",
   "larik": "RNMShareList",
   "syarat": [],
   "kolom": [
    "RNM Share",
    ""
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
    94519,
    99913
   ],
   "baca": [
    null,
    null
   ],
   "tombol": [
    null,
    null
   ],
   "tombolKepala": [
    null,
    null
   ],
   "pilihan": [
    null,
    null
   ],
   "aksiUbah": [
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "blok",
   "at": 137662,
   "judul": "Spreading",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 148231,
     "judul": "",
     "syarat": [
      ".SpreadingTypeID != ''"
     ],
     "anak": [
      {
       "t": "medan",
       "at": 154957,
       "label": "Spreading Type",
       "dari": "sisi",
       "kunci": "SpreadingTypeID",
       "format": "pxDropdown",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = '1'"
       ],
       "pilihan": {
        "sumber": "reportdefinition",
        "rd": "BrowseTreatyArrangement_ParentReinsMasterTrt",
        "nilai": "ReinsTypeID",
        "tampil": "ReinsTypeName"
       },
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "FetchQSfromMaster",
         "param": {
          "ParentReinsTypeID": ".SpreadingTypeID"
         }
        }
       ]
      }
     ]
    },
    {
     "t": "blok",
     "at": 184685,
     "judul": "",
     "syarat": [
      ".SpreadingTypeID != ''"
     ],
     "anak": [
      {
       "t": "grid",
       "at": 210120,
       "prop": ".SpreadingList",
       "dari": "sisi",
       "larik": "SpreadingList",
       "syarat": [],
       "kolom": [
        "Reins Type",
        "Pct"
       ],
       "kunci": [
        "ReinsTypeName",
        "Pct"
       ],
       "lebar": [
        200,
        192
       ],
       "desimal": [
        null,
        null
       ],
       "format": [
        "pxTextInput",
        ""
       ],
       "syaratSel": [
        null,
        null
       ],
       "atSel": [
        223229,
        228894
       ],
       "baca": [
        "selalu",
        "selalu"
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitsSpreading!pyGridModalTemplate"
      },
      {
       "t": "teks",
       "at": 246518,
       "teks": "Total Pct",
       "syarat": []
      }
     ]
    },
    {
     "t": "blok",
     "at": 262103,
     "judul": "",
     "syarat": [
      ".SpreadingTypeID == ''"
     ],
     "anak": [
      {
       "t": "grid",
       "at": 287538,
       "prop": ".SpreadingList",
       "dari": "sisi",
       "larik": "SpreadingList",
       "syarat": [],
       "kolom": [
        "Reins Type",
        "Pct Share",
        ".pyTemplateInputBox"
       ],
       "kunci": [
        "ReinsTypeID",
        "Pct",
        "pyTemplateInputBox"
       ],
       "lebar": [
        269,
        192,
        84
       ],
       "desimal": [
        null,
        2,
        null
       ],
       "format": [
        "pxDropdown",
        "pxNumber",
        "pxLink"
       ],
       "syaratSel": [
        null,
        null,
        "TreatyIn.ViewState != '1'"
       ],
       "atSel": [
        310892,
        325413,
        335540
       ],
       "baca": [
        [
         "TreatyIn.ViewState = '1'"
        ],
        [
         "TreatyIn.ViewState = '1'"
        ],
        null
       ],
       "tombol": [
        null,
        null,
        null
       ],
       "tombolKepala": [
        null,
        null,
        null
       ],
       "pilihan": [
        {
         "sumber": "reportdefinition",
         "rd": "BrowseTreatyArrangement_ParentReinsMasterTrt",
         "nilai": "ReinsTypeID",
         "tampil": "ReinsTypeName"
        },
        null,
        null
       ],
       "aksiUbah": [
        [
         {
          "aksi": "postValue"
         },
         {
          "aksi": "refresh",
          "aktivitas": "SetSpreadName"
         }
        ],
        [
         {
          "aksi": "postValue"
         },
         {
          "aksi": "refresh",
          "aktivitas": "SetSpreadName"
         }
        ],
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitsSpreading!pyGridRowDetails"
      },
      {
       "t": "teks",
       "at": 357434,
       "teks": "Total Pct",
       "syarat": []
      }
     ]
    },
    {
     "t": "medan",
     "at": 388630,
     "label": "Total Spreading Pct :",
     "dari": "sisi",
     "kunci": "SpreadingTotalPct",
     "format": "pxNumber",
     "desimal": 2,
     "syarat": [
      ".SpreadingTypeID != ''"
     ],
     "baca": "selalu"
    },
    {
     "t": "medan",
     "at": 395971,
     "label": "Total Share Pct",
     "dari": "sisi",
     "kunci": "SpreadingTotalPct",
     "format": "pxNumber",
     "desimal": 2,
     "syarat": [
      ".SpreadingTypeID = ''"
     ],
     "baca": "selalu"
    },
    {
     "t": "grid",
     "at": 441471,
     "prop": ".RNMSpreadedList",
     "dari": "sisi",
     "larik": "RNMSpreadedList",
     "syarat": [],
     "kolom": [
      "Value Spreading OR",
      ""
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
      453850,
      459001
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 507430,
     "prop": ".RNMSpreadedListRI",
     "dari": "sisi",
     "larik": "RNMSpreadedListRI",
     "syarat": [],
     "kolom": [
      "Value Spreading R/I",
      ""
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
      519812,
      524963
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    }
   ]
  }
 ],
 "DetailShareOldData": [
  {
   "t": "medan",
   "at": 26614,
   "label": "% RNM Share",
   "dari": "sisi",
   "kunci": "RNMShare",
   "format": "pxNumber",
   "desimal": 2,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh",
     "aktivitas": "TreatyInPropshareDetail"
    }
   ]
  },
  {
   "t": "grid",
   "at": 79061,
   "prop": ".RNMShareList",
   "dari": "sisi",
   "larik": "RNMShareList",
   "syarat": [],
   "kolom": [
    "RNM Share",
    ""
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
    92744,
    98138
   ],
   "baca": [
    null,
    null
   ],
   "tombol": [
    null,
    null
   ],
   "tombolKepala": [
    null,
    null
   ],
   "pilihan": [
    null,
    null
   ],
   "aksiUbah": [
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "blok",
   "at": 135887,
   "judul": "Spreading",
   "syarat": [],
   "anak": [
    {
     "t": "blok",
     "at": 146456,
     "judul": "",
     "syarat": [
      ".SpreadingTypeID != ''"
     ],
     "anak": [
      {
       "t": "medan",
       "at": 153182,
       "label": "Spreading Type",
       "dari": "sisi",
       "kunci": "SpreadingTypeID",
       "format": "pxDropdown",
       "desimal": null,
       "syarat": [],
       "baca": [
        "TreatyIn.ViewState = '1'"
       ],
       "pilihan": {
        "sumber": "reportdefinition",
        "rd": "BrowseTreatyArrangement_ParentReinsMasterTrt",
        "nilai": "ReinsTypeID",
        "tampil": "ReinsTypeName"
       },
       "aksiUbah": [
        {
         "aksi": "refresh",
         "aktivitas": "FetchQSfromMaster",
         "param": {
          "ParentReinsTypeID": ".SpreadingTypeID"
         }
        }
       ]
      }
     ]
    },
    {
     "t": "blok",
     "at": 182910,
     "judul": "",
     "syarat": [
      ".SpreadingTypeID != ''"
     ],
     "anak": [
      {
       "t": "grid",
       "at": 208345,
       "prop": ".SpreadingList",
       "dari": "sisi",
       "larik": "SpreadingList",
       "syarat": [],
       "kolom": [
        "Reins Type",
        "Pct"
       ],
       "kunci": [
        "ReinsTypeName",
        "Pct"
       ],
       "lebar": [
        200,
        192
       ],
       "desimal": [
        null,
        null
       ],
       "format": [
        "pxTextInput",
        ""
       ],
       "syaratSel": [
        null,
        null
       ],
       "atSel": [
        221454,
        227158
       ],
       "baca": [
        "selalu",
        "selalu"
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitsSpreading!pyGridModalTemplate"
      },
      {
       "t": "teks",
       "at": 244782,
       "teks": "Total Pct",
       "syarat": []
      }
     ]
    },
    {
     "t": "blok",
     "at": 260367,
     "judul": "",
     "syarat": [
      ".SpreadingTypeID == ''"
     ],
     "anak": [
      {
       "t": "grid",
       "at": 285802,
       "prop": ".SpreadingList",
       "dari": "sisi",
       "larik": "SpreadingList",
       "syarat": [],
       "kolom": [
        "Reins Type",
        "Pct Share"
       ],
       "kunci": [
        "ReinsTypeID",
        "Pct"
       ],
       "lebar": [
        269,
        192
       ],
       "desimal": [
        null,
        2
       ],
       "format": [
        "pxDropdown",
        "pxNumber"
       ],
       "syaratSel": [
        null,
        null
       ],
       "atSel": [
        309137,
        323613
       ],
       "baca": [
        "selalu",
        "selalu"
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        {
         "sumber": "reportdefinition",
         "rd": "BrowseTreatyArrangement_ParentReinsMasterTrt",
         "nilai": "ReinsTypeID",
         "tampil": "ReinsTypeName"
        },
        null
       ],
       "aksiUbah": [
        [
         {
          "aksi": "postValue"
         },
         {
          "aksi": "refresh",
          "aktivitas": "SetSpreadName"
         }
        ],
        [
         {
          "aksi": "postValue"
         },
         {
          "aksi": "refresh",
          "aktivitas": "SetSpreadName"
         }
        ]
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitsSpreading!pyGridRowDetails"
      },
      {
       "t": "teks",
       "at": 355570,
       "teks": "Total Pct",
       "syarat": []
      }
     ]
    },
    {
     "t": "medan",
     "at": 386766,
     "label": "Total Spreading Pct :",
     "dari": "sisi",
     "kunci": "SpreadingTotalPct",
     "format": "pxNumber",
     "desimal": 2,
     "syarat": [
      ".SpreadingTypeID != ''"
     ],
     "baca": "selalu"
    },
    {
     "t": "medan",
     "at": 394107,
     "label": "Total Share Pct",
     "dari": "sisi",
     "kunci": "SpreadingTotalPct",
     "format": "pxNumber",
     "desimal": 2,
     "syarat": [
      ".SpreadingTypeID = ''"
     ],
     "baca": "selalu"
    },
    {
     "t": "grid",
     "at": 439607,
     "prop": ".RNMSpreadedList",
     "dari": "sisi",
     "larik": "RNMSpreadedList",
     "syarat": [],
     "kolom": [
      "Value Spreading OR",
      ""
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
      451986,
      457137
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 505566,
     "prop": ".RNMSpreadedListRI",
     "dari": "sisi",
     "larik": "RNMSpreadedListRI",
     "syarat": [],
     "kolom": [
      "Value Spreading R/I",
      ""
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
      517948,
      523099
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    }
   ]
  }
 ],
 "Installments": [
  {
   "t": "grid",
   "at": 21076,
   "prop": ".InstallmentList",
   "dari": "sisi",
   "larik": "InstallmentList",
   "syarat": [],
   "kolom": [
    "Installment",
    "Due Date",
    "WPC (In Days)",
    "Payment Date",
    "% Installment",
    "Amount"
   ],
   "kunci": [
    "Installment",
    "DueDate",
    "WPC",
    "PaymentDate",
    "InstallmentPct",
    "Amount"
   ],
   "lebar": [
    75,
    187,
    187,
    174,
    188,
    188
   ],
   "desimal": [
    null,
    null,
    null,
    null,
    2,
    2
   ],
   "format": [
    "pxTextInput",
    "pxDateTime",
    "pxTextInput",
    "pxDateTime",
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
    51175,
    56889,
    66398,
    75179,
    81277,
    91479
   ],
   "baca": [
    "selalu",
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ],
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ],
    "selalu",
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ],
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ]
   ],
   "tombol": [
    null,
    null,
    null,
    null,
    null,
    null
   ],
   "tombolKepala": [
    null,
    null,
    null,
    null,
    null,
    null
   ],
   "pilihan": [
    null,
    null,
    null,
    null,
    null,
    null
   ],
   "aksiUbah": [
    null,
    [
     {
      "aksi": "refresh",
      "aktivitas": "TreatyInUpdatePaymentDate"
     }
    ],
    [
     {
      "aksi": "refresh",
      "aktivitas": "TreatyInUpdatePaymentDate_Act"
     }
    ],
    null,
    [
     {
      "aksi": "refresh",
      "aktivitas": "SetTotalInstallment",
      "param": {
       "status": "editpercentage"
      }
     }
    ],
    [
     {
      "aksi": "refresh",
      "aktivitas": "SetTotalInstallment"
     }
    ]
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInInstallment!pyGridModalTemplate"
  },
  {
   "t": "medan",
   "at": 143775,
   "label": "% Total",
   "dari": "sisi",
   "kunci": "PctTotal",
   "format": "pxNumber",
   "desimal": 4,
   "syarat": [],
   "baca": "selalu"
  },
  {
   "t": "medan",
   "at": 149318,
   "label": "Total",
   "dari": "sisi",
   "kunci": "AmountTotal",
   "format": "pxNumber",
   "desimal": 2,
   "syarat": [],
   "baca": "selalu"
  }
 ],
 "Installments_ReadOnly": [
  {
   "t": "grid",
   "at": 20251,
   "prop": ".InstallmentList",
   "dari": "sisi",
   "larik": "InstallmentList",
   "syarat": [],
   "kolom": [
    "Installment",
    "Due Date",
    "WPC (In Days)",
    "Payment Date",
    "% Installment",
    "Amount"
   ],
   "kunci": [
    "Installment",
    "DueDate",
    "WPC",
    "PaymentDate",
    "InstallmentPct",
    "Amount"
   ],
   "lebar": [
    76,
    184,
    184,
    172,
    184,
    184
   ],
   "desimal": [
    null,
    null,
    null,
    null,
    2,
    2
   ],
   "format": [
    "pxTextInput",
    "pxDateTime",
    "pxTextInput",
    "pxDateTime",
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
    50349,
    56063,
    62752,
    68913,
    75011,
    80811
   ],
   "baca": [
    "selalu",
    "selalu",
    "selalu",
    "selalu",
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    null,
    null,
    null,
    null,
    null
   ],
   "tombolKepala": [
    null,
    null,
    null,
    null,
    null,
    null
   ],
   "pilihan": [
    null,
    null,
    null,
    null,
    null,
    null
   ],
   "aksiUbah": [
    null,
    null,
    null,
    null,
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInInstallment!pyGridModalTemplate"
  },
  {
   "t": "medan",
   "at": 128803,
   "label": "% Total",
   "dari": "sisi",
   "kunci": "PctTotal",
   "format": "pxNumber",
   "desimal": 4,
   "syarat": [],
   "baca": "selalu"
  },
  {
   "t": "medan",
   "at": 134346,
   "label": "Total",
   "dari": "sisi",
   "kunci": "AmountTotal",
   "format": "pxNumber",
   "desimal": 2,
   "syarat": [],
   "baca": "selalu"
  }
 ],
 "Layers": [
  {
   "t": "medan",
   "at": 19542,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerType",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "pilihan": {
    "sumber": "associated"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 27061,
   "label": "",
   "dari": "sisi",
   "kunci": "Layer",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "teks",
   "at": 34562,
   "teks": "Part of",
   "syarat": []
  },
  {
   "t": "medan",
   "at": 39230,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerPartType",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "pilihan": {
    "sumber": "associated"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 46737,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerPart",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "blok",
   "at": 60018,
   "judul": "",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 79117,
     "prop": ".TreatyGroupList",
     "dari": "sisi",
     "larik": "TreatyGroupList",
     "syarat": [],
     "kolom": [
      "Treaty Group",
      "",
      "ROL Profile"
     ],
     "kunci": [
      "TreatyGroup",
      "",
      "IsROLProfile"
     ],
     "lebar": [
      202,
      147,
      104
     ],
     "desimal": [
      null,
      null,
      null
     ],
     "format": [
      "pxAutoComplete",
      "pxButton",
      "pxCheckbox"
     ],
     "syaratSel": [
      null,
      "TreatyIn.ViewState !='1'",
      null
     ],
     "atSel": [
      103329,
      116975,
      127532
     ],
     "baca": [
      [
       "TreatyIn.EDMMaterialType = 2"
      ],
      "selalu",
      [
       "TreatyIn.ViewState = '1' || TreatyIn.EDMMaterialType = 2"
      ]
     ],
     "tombol": [
      null,
      {
       "t": "tombol",
       "at": 116975,
       "label": "Delete",
       "syarat": [
        "TreatyIn.ViewState !='1'"
       ],
       "aksi": [
        {
         "aksi": "deleteRow"
        },
        {
         "aksi": "refresh",
         "aktivitas": "TotalEgnpi",
         "param": {
          "idxLimit": ".pxListSubscript"
         }
        }
       ],
       "nonaktif": [
        "TreatyIn.EDMMaterialType = 2"
       ]
      },
      null
     ],
     "tombolKepala": [
      null,
      {
       "t": "tombol",
       "at": 87797,
       "label": "Add Treaty Group",
       "syarat": [
        "TreatyIn.ViewState !='1'"
       ],
       "aksi": [
        {
         "aksi": "addRow"
        },
        {
         "aksi": "runDataTransform",
         "transformasi": "TreatyTypeSetIndex"
        }
       ],
       "nonaktif": [
        "TreatyIn.EDMMaterialType = 2"
       ]
      },
      null
     ],
     "pilihan": [
      {
       "sumber": "reportdefinition",
       "rd": "BrowseTreatyGroup_RD",
       "nilai": "TreatyGroupName"
      },
      null,
      null
     ],
     "aksiUbah": [
      [
       {
        "aksi": "runDataTransform",
        "transformasi": "SetIndexLayer_DT",
        "paramDT": {
         "idxlimit": ""
        }
       },
       {
        "aksi": "refresh",
        "aktivitas": "TotalEgnpi",
        "param": {
         "idxLimit": ".Layer"
        }
       }
      ],
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitLayer!pyGridRowDetails",
     "rincian": "CoBList"
    }
   ]
  },
  {
   "t": "grid",
   "at": 185070,
   "prop": ".EgnpiTotalList",
   "dari": "sisi",
   "larik": "EgnpiTotalList",
   "syarat": [],
   "kolom": [
    "Egnpi this layer",
    ""
   ],
   "kunci": [
    "Currency",
    "Value"
   ],
   "lebar": [
    196,
    184
   ],
   "desimal": [
    null,
    2
   ],
   "format": [
    "pxTextInput",
    "pxNumber"
   ],
   "syaratSel": [
    null,
    null
   ],
   "atSel": [
    197316,
    204185
   ],
   "baca": [
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    null
   ],
   "tombolKepala": [
    null,
    null
   ],
   "pilihan": [
    null,
    null
   ],
   "aksiUbah": [
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridRowDetails"
  },
  {
   "t": "medan",
   "at": 258061,
   "label": "",
   "dari": "sisi",
   "kunci": "Cover",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "pilihan": {
    "sumber": "associated"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 263935,
   "label": "Currency Relation",
   "dari": "sisi",
   "kunci": "CurrencyRelation",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1",
    "TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "associated"
   },
   "aksiUbah": [
    {
     "aksi": "refresh"
    }
   ]
  },
  {
   "t": "medan",
   "at": 305127,
   "label": "100 % Limit",
   "dari": "sisi",
   "kunci": "Currency",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1",
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 330498,
   "label": "",
   "dari": "sisi",
   "kunci": "Limit",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 364116,
   "label": "Agregate Year Limit",
   "dari": "sisi",
   "kunci": "Currency",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 389460,
   "label": "",
   "dari": "sisi",
   "kunci": "AgregateLimit",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 422493,
   "label": "Deductible",
   "dari": "sisi",
   "kunci": "Currency",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 447802,
   "label": "",
   "dari": "sisi",
   "kunci": "Deductible",
   "format": "pxNumber",
   "desimal": 2,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 494821,
   "label": "100 % Limit",
   "dari": "sisi",
   "kunci": "Currency2",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 520151,
   "label": "",
   "dari": "sisi",
   "kunci": "Limit2",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 553771,
   "label": "Agregate Year Limit",
   "dari": "sisi",
   "kunci": "Currency2",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 579119,
   "label": "",
   "dari": "sisi",
   "kunci": "AgregateLimit2",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 612158,
   "label": "Deductible",
   "dari": "sisi",
   "kunci": "Currency2",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 637471,
   "label": "",
   "dari": "sisi",
   "kunci": "Deductible2",
   "format": "pxNumber",
   "desimal": 2,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 682521,
   "label": "",
   "dari": "sisi",
   "kunci": "NoRIPCalculation",
   "format": "pxCheckbox",
   "desimal": null,
   "syarat": [],
   "caption": "No Reinstatement Premium Calculation",
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "teks",
   "at": 704972,
   "teks": "Reinstatement",
   "syarat": []
  },
  {
   "t": "medan",
   "at": 714520,
   "label": "",
   "dari": "sisi",
   "kunci": "ReinstatementValue",
   "format": "pxNumber",
   "desimal": 0,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "refresh",
     "aktivitas": "SetReinstatementPct"
    }
   ]
  },
  {
   "t": "grid",
   "at": 770156,
   "prop": ".Reinstatement_List",
   "dari": "sisi",
   "larik": "Reinstatement_List",
   "syarat": [],
   "kolom": [
    "Reinstatement No.",
    "% Additional Premium",
    "Note",
    "Reinstatement Premium Amount IDR (MDP x % Add Premium)",
    "Reinstatement Premium Amount USD (MDP x % Add Premium)",
    "Reinstatement %",
    "Reinstatement Amount IDR",
    "Reinstatement Amount USD"
   ],
   "kunci": [
    "ReinstatementValue",
    "ReinstatementPct",
    "ReinstatementNote",
    "AdditionalAmount1",
    "AdditionalAmount2",
    "AdditionalPct",
    "ReinstatementAmount1",
    "ReinstatementAmount2"
   ],
   "lebar": [
    156,
    146,
    319,
    236,
    231,
    132,
    302,
    304
   ],
   "desimal": [
    null,
    2,
    null,
    2,
    2,
    2,
    2,
    2
   ],
   "format": [
    "",
    "pxNumber",
    "pxDropdown",
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
    ".Limit2 != 0"
   ],
   "atSel": [
    806441,
    811346,
    820051,
    828206,
    833794,
    838915,
    847709,
    856176
   ],
   "baca": [
    "selalu",
    [
     "TreatyIn.ViewState = '1'",
     "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
    ],
    [
     "TreatyIn.ViewState = '1'",
     "TreatyIn.EDMMaterialType = 2"
    ],
    [
     "TreatyIn.ViewState = '1'",
     "TreatyIn.EDMMaterialType = 2"
    ],
    [
     "TreatyIn.ViewState = '1'",
     "TreatyIn.EDMMaterialType = 2"
    ],
    [
     "TreatyIn.ViewState = '1'",
     "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
    ],
    [
     "TreatyIn.ViewState = '1'",
     "TreatyIn.EDMMaterialType = 2"
    ],
    [
     "TreatyIn.ViewState = '1'",
     "TreatyIn.EDMMaterialType = 2"
    ]
   ],
   "tombol": [
    null,
    null,
    null,
    null,
    null,
    null,
    null,
    null
   ],
   "tombolKepala": [
    null,
    null,
    null,
    null,
    null,
    null,
    null,
    null
   ],
   "pilihan": [
    null,
    null,
    {
     "sumber": "associated"
    },
    null,
    null,
    null,
    null,
    null
   ],
   "aksiUbah": [
    null,
    [
     {
      "aksi": "refresh",
      "transformasi": "ReCalculateReinstatement",
      "paramDT": {
       "Subscript": ".pxListSubscript"
      }
     }
    ],
    [
     {
      "aksi": "refresh"
     }
    ],
    null,
    null,
    [
     {
      "aksi": "refresh",
      "transformasi": "CalculateReinstatement",
      "paramDT": {
       "Subscript": ".pxListSubscript"
      }
     }
    ],
    [
     {
      "aksi": "refresh",
      "transformasi": "CalculateReinstatementPct",
      "paramDT": {
       "Subscript": ".pxListSubscript"
      }
     }
    ],
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridModalTemplate"
  },
  {
   "t": "medan",
   "at": 890972,
   "label": "Adjustment Rate %",
   "dari": "sisi",
   "kunci": "AdjRate",
   "format": "pxNumber",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "refresh",
     "aktivitas": "DetailCalculation",
     "param": {
      "type": "adj"
     }
    }
   ]
  },
  {
   "t": "grid",
   "at": 918915,
   "prop": ".PremiumEarnedList",
   "dari": "sisi",
   "larik": "PremiumEarnedList",
   "syarat": [],
   "kolom": [
    "Premium Earned",
    "Value"
   ],
   "kunci": [
    "Currency",
    "Value"
   ],
   "lebar": [
    197,
    355
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
    930511,
    935250
   ],
   "baca": [
    "selalu",
    [
     "TreatyIn.EDMMaterialType = 2"
    ]
   ],
   "tombol": [
    null,
    null
   ],
   "tombolKepala": [
    null,
    null
   ],
   "pilihan": [
    null,
    null
   ],
   "aksiUbah": [
    null,
    [
     {
      "aksi": "refresh"
     }
    ]
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "medan",
   "at": 966789,
   "label": "Min Premium %",
   "dari": "sisi",
   "kunci": "MDPMinPct",
   "format": "pxNumber",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "refresh",
     "aktivitas": "DetailCalculation",
     "param": {
      "type": "mdp"
     }
    }
   ]
  },
  {
   "t": "grid",
   "at": 992599,
   "prop": ".MDPMinList",
   "dari": "sisi",
   "larik": "MDPMinList",
   "syarat": [],
   "kolom": [
    "Min Premium Amt",
    "Value",
    ""
   ],
   "kunci": [
    "Currency",
    "Value",
    ""
   ],
   "lebar": [
    198,
    363,
    101
   ],
   "desimal": [
    null,
    2,
    null
   ],
   "format": [
    "pxAutoComplete",
    "pxNumber",
    "pxButton"
   ],
   "syaratSel": [
    null,
    null,
    "TreatyIn.ViewState !=1"
   ],
   "atSel": [
    1012631,
    1022882,
    1030653
   ],
   "baca": [
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ],
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ],
    "selalu"
   ],
   "tombol": [
    null,
    null,
    {
     "t": "tombol",
     "at": 1030653,
     "label": "Remove",
     "syarat": [
      "TreatyIn.ViewState !=1"
     ],
     "aksi": [
      {
       "aksi": "deleteRow"
      }
     ],
     "nonaktif": [
      "TreatyIn.EDMMaterialType = 2"
     ]
    }
   ],
   "tombolKepala": [
    null,
    null,
    {
     "t": "tombol",
     "at": 1003822,
     "label": "Add MDP",
     "syarat": [
      "TreatyIn.ViewState !=1"
     ],
     "aksi": [
      {
       "aksi": "addRow"
      }
     ],
     "nonaktif": [
      "TreatyIn.EDMMaterialType = 2"
     ]
    }
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseCurrency_RD",
     "nilai": "Currency"
    },
    null,
    null
   ],
   "aksiUbah": [
    [
     {
      "aksi": "refresh"
     }
    ],
    [
     {
      "aksi": "refresh"
     }
    ],
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "medan",
   "at": 1061946,
   "label": "MDP %",
   "dari": "sisi",
   "kunci": "MDPPct",
   "format": "pxNumber",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "refresh",
     "aktivitas": "DetailCalculation",
     "param": {
      "type": "mdp"
     }
    }
   ]
  },
  {
   "t": "grid",
   "at": 1089868,
   "prop": ".MDPList",
   "dari": "sisi",
   "larik": "MDPList",
   "syarat": [],
   "kolom": [
    "MDP",
    "Value",
    ""
   ],
   "kunci": [
    "Currency",
    "Value",
    ""
   ],
   "lebar": [
    198,
    363,
    100
   ],
   "desimal": [
    null,
    2,
    null
   ],
   "format": [
    "pxAutoComplete",
    "pxNumber",
    "pxButton"
   ],
   "syaratSel": [
    null,
    null,
    "TreatyIn.ViewState !=1"
   ],
   "atSel": [
    1109893,
    1120144,
    1127915
   ],
   "baca": [
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ],
    [
     "TreatyIn.ViewState = 1",
     "TreatyIn.EDMMaterialType = 2"
    ],
    "selalu"
   ],
   "tombol": [
    null,
    null,
    {
     "t": "tombol",
     "at": 1127915,
     "label": "Remove",
     "syarat": [
      "TreatyIn.ViewState !=1"
     ],
     "aksi": [
      {
       "aksi": "deleteRow"
      }
     ],
     "nonaktif": [
      "TreatyIn.EDMMaterialType = 2"
     ]
    }
   ],
   "tombolKepala": [
    null,
    null,
    {
     "t": "tombol",
     "at": 1101084,
     "label": "Add MDP",
     "syarat": [
      "TreatyIn.ViewState !=1"
     ],
     "aksi": [
      {
       "aksi": "addRow"
      }
     ],
     "nonaktif": [
      "TreatyIn.EDMMaterialType = 2"
     ]
    }
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseCurrency_RD",
     "nilai": "Currency"
    },
    null,
    null
   ],
   "aksiUbah": [
    [
     {
      "aksi": "refresh"
     }
    ],
    [
     {
      "aksi": "refresh"
     }
    ],
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "medan",
   "at": 1159213,
   "label": "",
   "dari": "sisi",
   "kunci": "IsCombineMDP",
   "format": "pxCheckbox",
   "desimal": null,
   "syarat": [],
   "caption": "(*) Combine MDP",
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh"
    }
   ]
  },
  {
   "t": "medan",
   "at": 1166404,
   "label": "ROL %",
   "dari": "sisi",
   "kunci": "ROLPct",
   "format": "pxNumber",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  }
 ],
 "LayersEDM": [
  {
   "t": "medan",
   "at": 17792,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerType",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "pilihan": {
    "sumber": "associated"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 25311,
   "label": "",
   "dari": "sisi",
   "kunci": "Layer",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "teks",
   "at": 32812,
   "teks": "Part of",
   "syarat": []
  },
  {
   "t": "medan",
   "at": 37480,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerPartType",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "pilihan": {
    "sumber": "associated"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 44987,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerPart",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "blok",
   "at": 58268,
   "judul": "",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 77367,
     "prop": ".TreatyGroupList",
     "dari": "sisi",
     "larik": "TreatyGroupList",
     "syarat": [],
     "kolom": [
      "Treaty Group",
      "ROL Profile"
     ],
     "kunci": [
      "TreatyGroup",
      "IsROLProfile"
     ],
     "lebar": [
      200,
      105
     ],
     "desimal": [
      null,
      null
     ],
     "format": [
      "pxAutoComplete",
      "pxCheckbox"
     ],
     "syaratSel": [
      null,
      null
     ],
     "atSel": [
      90680,
      104266
     ],
     "baca": [
      "selalu",
      "selalu"
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      {
       "sumber": "reportdefinition",
       "rd": "BrowseTreatyGroup_RD",
       "nilai": "TreatyGroupName"
      },
      null
     ],
     "aksiUbah": [
      [
       {
        "aksi": "runDataTransform",
        "transformasi": "SetIndexLayer_DT",
        "paramDT": {
         "idxlimit": ""
        }
       },
       {
        "aksi": "refresh",
        "aktivitas": "TotalEgnpi",
        "param": {
         "idxLimit": ".Layer"
        }
       }
      ],
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitLayer!pyGridRowDetails",
     "rincian": "CoBList"
    }
   ]
  },
  {
   "t": "grid",
   "at": 161054,
   "prop": ".EgnpiTotalList",
   "dari": "sisi",
   "larik": "EgnpiTotalList",
   "syarat": [],
   "kolom": [
    "Egnpi this layer",
    ""
   ],
   "kunci": [
    "Currency",
    "Value"
   ],
   "lebar": [
    196,
    184
   ],
   "desimal": [
    null,
    2
   ],
   "format": [
    "pxTextInput",
    "pxNumber"
   ],
   "syaratSel": [
    null,
    null
   ],
   "atSel": [
    173300,
    180169
   ],
   "baca": [
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    null
   ],
   "tombolKepala": [
    null,
    null
   ],
   "pilihan": [
    null,
    null
   ],
   "aksiUbah": [
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridRowDetails"
  },
  {
   "t": "medan",
   "at": 234045,
   "label": "",
   "dari": "sisi",
   "kunci": "Cover",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "pilihan": {
    "sumber": "associated"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 239919,
   "label": "Currency Relation",
   "dari": "sisi",
   "kunci": "CurrencyRelation",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1",
    "TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "associated"
   }
  },
  {
   "t": "medan",
   "at": 277924,
   "label": "100 % Limit",
   "dari": "sisi",
   "kunci": "Currency",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 303209,
   "label": "",
   "dari": "sisi",
   "kunci": "Limit",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 336741,
   "label": "Agregate Year Limit",
   "dari": "sisi",
   "kunci": "Currency",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 361999,
   "label": "",
   "dari": "sisi",
   "kunci": "AgregateLimit",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 394946,
   "label": "Deductible",
   "dari": "sisi",
   "kunci": "Currency",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 420169,
   "label": "",
   "dari": "sisi",
   "kunci": "Deductible",
   "format": "pxNumber",
   "desimal": 2,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 467102,
   "label": "100 % Limit",
   "dari": "sisi",
   "kunci": "Currency2",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 492346,
   "label": "",
   "dari": "sisi",
   "kunci": "Limit2",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 525877,
   "label": "Agregate Year Limit",
   "dari": "sisi",
   "kunci": "Currency2",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 551139,
   "label": "",
   "dari": "sisi",
   "kunci": "AgregateLimit2",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 584092,
   "label": "Deductible",
   "dari": "sisi",
   "kunci": "Currency2",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 609319,
   "label": "",
   "dari": "sisi",
   "kunci": "Deductible2",
   "format": "pxNumber",
   "desimal": 2,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "teks",
   "at": 654283,
   "teks": "Reinstatement",
   "syarat": []
  },
  {
   "t": "medan",
   "at": 663831,
   "label": "",
   "dari": "sisi",
   "kunci": "ReinstatementValue",
   "format": "pxNumber",
   "desimal": 0,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "refresh",
     "aktivitas": "SetReinstatementPct"
    }
   ]
  },
  {
   "t": "grid",
   "at": 718829,
   "prop": ".Reinstatement_List",
   "dari": "sisi",
   "larik": "Reinstatement_List",
   "syarat": [],
   "kolom": [
    "Reinstatement No.",
    "Reinstatement %",
    "Note"
   ],
   "kunci": [
    "ReinstatementValue",
    "ReinstatementPct",
    "ReinstatementNote"
   ],
   "lebar": [
    140,
    125,
    272
   ],
   "desimal": [
    null,
    2,
    null
   ],
   "format": [
    "",
    "pxNumber",
    "pxDropdown"
   ],
   "syaratSel": [
    null,
    null,
    null
   ],
   "atSel": [
    734215,
    739120,
    747381
   ],
   "baca": [
    "selalu",
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    null,
    null
   ],
   "tombolKepala": [
    null,
    null,
    null
   ],
   "pilihan": [
    null,
    null,
    {
     "sumber": "associated"
    }
   ],
   "aksiUbah": [
    null,
    null,
    [
     {
      "aksi": "refresh"
     }
    ]
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridModalTemplate"
  },
  {
   "t": "medan",
   "at": 781216,
   "label": "Adjustment Rate %",
   "dari": "sisi",
   "kunci": "AdjRate",
   "format": "pxNumber",
   "desimal": null,
   "syarat": [],
   "aksiUbah": [
    {
     "aksi": "refresh",
     "aktivitas": "DetailCalculation",
     "param": {
      "type": "adj"
     }
    }
   ]
  },
  {
   "t": "grid",
   "at": 809073,
   "prop": ".PremiumEarnedList",
   "dari": "sisi",
   "larik": "PremiumEarnedList",
   "syarat": [],
   "kolom": [
    "Premium Earned",
    "Value"
   ],
   "kunci": [
    "Currency",
    "Value"
   ],
   "lebar": [
    195,
    354
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
    820669,
    825408
   ],
   "baca": [
    "selalu",
    [
     "TreatyIn.ViewState = '1'"
    ]
   ],
   "tombol": [
    null,
    null
   ],
   "tombolKepala": [
    null,
    null
   ],
   "pilihan": [
    null,
    null
   ],
   "aksiUbah": [
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "medan",
   "at": 854171,
   "label": "MDP %",
   "dari": "sisi",
   "kunci": "MDPPct",
   "format": "pxNumber",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "refresh",
     "aktivitas": "DetailCalculation",
     "param": {
      "type": "mdp"
     }
    }
   ]
  },
  {
   "t": "grid",
   "at": 882003,
   "prop": ".MDPList",
   "dari": "sisi",
   "larik": "MDPList",
   "syarat": [],
   "kolom": [
    "MDP",
    "Value"
   ],
   "kunci": [
    "Currency",
    "Value"
   ],
   "lebar": [
    195,
    357
   ],
   "desimal": [
    null,
    2
   ],
   "format": [
    "pxAutoComplete",
    "pxNumber"
   ],
   "syaratSel": [
    null,
    null
   ],
   "atSel": [
    893579,
    906709
   ],
   "baca": [
    [
     "TreatyIn.ViewState = 1"
    ],
    [
     "TreatyIn.ViewState = 1"
    ]
   ],
   "tombol": [
    null,
    null
   ],
   "tombolKepala": [
    null,
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseCurrency_RD",
     "nilai": "Currency"
    },
    null
   ],
   "aksiUbah": [
    [
     {
      "aksi": "postValue"
     },
     {
      "aksi": "refresh"
     }
    ],
    [
     {
      "aksi": "postValue"
     },
     {
      "aksi": "refresh"
     }
    ]
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "medan",
   "at": 940306,
   "label": "",
   "dari": "sisi",
   "kunci": "IsCombineMDP",
   "format": "pxCheckbox",
   "desimal": null,
   "syarat": [],
   "caption": "(*) Combine MDP",
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh"
    }
   ]
  },
  {
   "t": "medan",
   "at": 947495,
   "label": "ROL %",
   "dari": "sisi",
   "kunci": "ROLPct",
   "format": "pxNumber",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  }
 ],
 "LayersOldData": [
  {
   "t": "medan",
   "at": 17554,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerType",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "associated"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 25033,
   "label": "",
   "dari": "sisi",
   "kunci": "Layer",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "teks",
   "at": 32494,
   "teks": "Part of",
   "syarat": []
  },
  {
   "t": "medan",
   "at": 37165,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerPartType",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "associated"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 44632,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerPart",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "blok",
   "at": 57873,
   "judul": "",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 76993,
     "prop": ".TreatyGroupList",
     "dari": "sisi",
     "larik": "TreatyGroupList",
     "syarat": [],
     "kolom": [
      "Treaty Group",
      "ROL Profile"
     ],
     "kunci": [
      "TreatyGroup",
      "IsROLProfile"
     ],
     "lebar": [
      200,
      105
     ],
     "desimal": [
      null,
      null
     ],
     "format": [
      "pxAutoComplete",
      "pxCheckbox"
     ],
     "syaratSel": [
      null,
      null
     ],
     "atSel": [
      90316,
      103877
     ],
     "baca": [
      "selalu",
      "selalu"
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      {
       "sumber": "reportdefinition",
       "rd": "BrowseTreatyGroup_RD",
       "nilai": "TreatyGroupName"
      },
      null
     ],
     "aksiUbah": [
      [
       {
        "aksi": "runDataTransform",
        "transformasi": "SetIndexLayer_DT",
        "paramDT": {
         "idxlimit": ""
        }
       },
       {
        "aksi": "refresh",
        "aktivitas": "TotalEgnpi",
        "param": {
         "idxLimit": ".Layer"
        }
       }
      ],
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitLayer!pyGridRowDetails",
     "rincian": "CoBListOldData"
    }
   ]
  },
  {
   "t": "grid",
   "at": 160712,
   "prop": ".EgnpiTotalList",
   "dari": "sisi",
   "larik": "EgnpiTotalList",
   "syarat": [],
   "kolom": [
    "Egnpi this layer",
    ""
   ],
   "kunci": [
    "Currency",
    "Value"
   ],
   "lebar": [
    196,
    184
   ],
   "desimal": [
    null,
    2
   ],
   "format": [
    "pxTextInput",
    "pxNumber"
   ],
   "syaratSel": [
    null,
    null
   ],
   "atSel": [
    172968,
    179839
   ],
   "baca": [
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    null
   ],
   "tombolKepala": [
    null,
    null
   ],
   "pilihan": [
    null,
    null
   ],
   "aksiUbah": [
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridRowDetails"
  },
  {
   "t": "medan",
   "at": 233741,
   "label": "",
   "dari": "sisi",
   "kunci": "Cover",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "pilihan": {
    "sumber": "associated"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 239617,
   "label": "Currency Relation",
   "dari": "sisi",
   "kunci": "CurrencyRelation",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1",
    "TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "associated"
   }
  },
  {
   "t": "medan",
   "at": 277633,
   "label": "100 % Limit",
   "dari": "sisi",
   "kunci": "Currency",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 302923,
   "label": "",
   "dari": "sisi",
   "kunci": "Limit",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 336464,
   "label": "Agregate Year Limit",
   "dari": "sisi",
   "kunci": "Currency",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 361727,
   "label": "",
   "dari": "sisi",
   "kunci": "AgregateLimit",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 394683,
   "label": "Deductible",
   "dari": "sisi",
   "kunci": "Currency",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 419911,
   "label": "",
   "dari": "sisi",
   "kunci": "Deductible",
   "format": "pxNumber",
   "desimal": 2,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 466856,
   "label": "100 % Limit",
   "dari": "sisi",
   "kunci": "Currency2",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 492105,
   "label": "",
   "dari": "sisi",
   "kunci": "Limit2",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 525645,
   "label": "Agregate Year Limit",
   "dari": "sisi",
   "kunci": "Currency2",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 550910,
   "label": "",
   "dari": "sisi",
   "kunci": "AgregateLimit2",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 583868,
   "label": "Deductible",
   "dari": "sisi",
   "kunci": "Currency2",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 609098,
   "label": "",
   "dari": "sisi",
   "kunci": "Deductible2",
   "format": "pxNumber",
   "desimal": 2,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "teks",
   "at": 654067,
   "teks": "Reinstatement",
   "syarat": []
  },
  {
   "t": "medan",
   "at": 663617,
   "label": "",
   "dari": "sisi",
   "kunci": "ReinstatementValue",
   "format": "pxNumber",
   "desimal": 0,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "refresh",
     "aktivitas": "SetReinstatementPct"
    }
   ]
  },
  {
   "t": "grid",
   "at": 719115,
   "prop": ".Reinstatement_List",
   "dari": "sisi",
   "larik": "Reinstatement_List",
   "syarat": [],
   "kolom": [
    "Reinstatement No.",
    "Reinstatement %",
    "Note"
   ],
   "kunci": [
    "ReinstatementValue",
    "ReinstatementPct",
    "ReinstatementNote"
   ],
   "lebar": [
    140,
    125,
    272
   ],
   "desimal": [
    null,
    2,
    null
   ],
   "format": [
    "",
    "pxNumber",
    "pxDropdown"
   ],
   "syaratSel": [
    null,
    null,
    null
   ],
   "atSel": [
    734507,
    739413,
    747675
   ],
   "baca": [
    "selalu",
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    null,
    null
   ],
   "tombolKepala": [
    null,
    null,
    null
   ],
   "pilihan": [
    null,
    null,
    {
     "sumber": "associated"
    }
   ],
   "aksiUbah": [
    null,
    null,
    [
     {
      "aksi": "refresh"
     }
    ]
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimits!pyGridModalTemplate"
  },
  {
   "t": "medan",
   "at": 779708,
   "label": "Adjustment Rate %",
   "dari": "sisi",
   "kunci": "AdjRate",
   "format": "pxNumber",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "refresh",
     "aktivitas": "DetailCalculation",
     "param": {
      "type": "adj"
     }
    }
   ]
  },
  {
   "t": "grid",
   "at": 807576,
   "prop": ".PremiumEarnedList",
   "dari": "sisi",
   "larik": "PremiumEarnedList",
   "syarat": [],
   "kolom": [
    "Premium Earned",
    "Value"
   ],
   "kunci": [
    "Currency",
    "Value"
   ],
   "lebar": [
    198,
    353
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
    819178,
    823918
   ],
   "baca": [
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    null
   ],
   "tombolKepala": [
    null,
    null
   ],
   "pilihan": [
    null,
    null
   ],
   "aksiUbah": [
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "medan",
   "at": 852609,
   "label": "MDP %",
   "dari": "sisi",
   "kunci": "MDPPct",
   "format": "pxNumber",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "refresh",
     "aktivitas": "DetailCalculation",
     "param": {
      "type": "mdp"
     }
    }
   ]
  },
  {
   "t": "grid",
   "at": 880452,
   "prop": ".MDPList",
   "dari": "sisi",
   "larik": "MDPList",
   "syarat": [],
   "kolom": [
    "MDP",
    "Value"
   ],
   "kunci": [
    "Currency",
    "Value"
   ],
   "lebar": [
    195,
    359
   ],
   "desimal": [
    null,
    2
   ],
   "format": [
    "pxAutoComplete",
    "pxNumber"
   ],
   "syaratSel": [
    null,
    null
   ],
   "atSel": [
    892032,
    899814
   ],
   "baca": [
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    null
   ],
   "tombolKepala": [
    null,
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseCurrency_RD",
     "nilai": "Currency"
    },
    null
   ],
   "aksiUbah": [
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
  },
  {
   "t": "medan",
   "at": 928462,
   "label": "",
   "dari": "sisi",
   "kunci": "IsCombineMDP",
   "format": "pxCheckbox",
   "desimal": null,
   "syarat": [],
   "caption": "(*) Combine MDP",
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh"
    }
   ]
  },
  {
   "t": "medan",
   "at": 935652,
   "label": "ROL %",
   "dari": "sisi",
   "kunci": "ROLPct",
   "format": "pxNumber",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  }
 ],
 "LimitProportional": [
  {
   "t": "medan",
   "at": 17830,
   "label": "Treaty Type",
   "dari": "sisi",
   "kunci": "TreatyTypeID",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseReinsuranceType_RD",
    "nilai": "ID",
    "tampil": "Note"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "refresh",
     "aktivitas": "SetTreatyTypeName_Act"
    }
   ]
  },
  {
   "t": "grid",
   "at": 49737,
   "prop": ".Detail",
   "dari": "sisi",
   "larik": "Detail",
   "syarat": [],
   "kolom": [
    "Treaty Group",
    ""
   ],
   "kunci": [
    "TreatyGroup",
    ""
   ],
   "lebar": [
    1368,
    83
   ],
   "desimal": [
    null,
    null
   ],
   "format": [
    "pxAutoComplete",
    "pxButton"
   ],
   "syaratSel": [
    null,
    "TreatyIn.ViewState !='1'"
   ],
   "atSel": [
    66806,
    75642
   ],
   "baca": [
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    {
     "t": "tombol",
     "at": 75642,
     "label": "Delete",
     "syarat": [
      "TreatyIn.ViewState !='1'"
     ],
     "aksi": [
      {
       "aksi": "deleteRow"
      }
     ]
    }
   ],
   "tombolKepala": [
    null,
    {
     "t": "tombol",
     "at": 57397,
     "label": "Add",
     "syarat": [
      "TreatyIn.ViewState !='1'"
     ],
     "aksi": [
      {
       "aksi": "refresh",
       "aktivitas": "AddClassofBusiness"
      }
     ],
     "ikon": "pyWorkActionsAddWork.png"
    }
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseBusiness_RD",
     "nilai": "Note"
    },
    null
   ],
   "aksiUbah": [
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitsDetail!pyGridRowDetails",
   "rincian": "DetailLimits"
  }
 ],
 "LimitProportionalOldData": [
  {
   "t": "medan",
   "at": 16974,
   "label": "Treaty Type",
   "dari": "sisi",
   "kunci": "TreatyType",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseReinsuranceType_RD",
    "nilai": "Note"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    },
    {
     "aksi": "runDataTransform",
     "transformasi": "TreatyTypeSetIndex"
    }
   ]
  },
  {
   "t": "grid",
   "at": 51554,
   "prop": ".Detail",
   "dari": "sisi",
   "larik": "Detail",
   "syarat": [],
   "kolom": [
    "Treaty Group"
   ],
   "kunci": [
    "TreatyGroup"
   ],
   "lebar": [
    1368
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
    59573
   ],
   "baca": [
    null
   ],
   "tombol": [
    null
   ],
   "tombolKepala": [
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseBusiness_RD",
     "nilai": "Note"
    }
   ],
   "aksiUbah": [
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitsDetail!pyGridRowDetails",
   "rincian": "DetailLimitsOldData"
  }
 ],
 "MaxRetention": [
  {
   "t": "medan",
   "at": 16635,
   "label": "Treaty Group",
   "dari": "sisi",
   "kunci": "TreatyGroup",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.IsEditData = 1",
    "TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseTreatyGroup_RD",
    "nilai": "TreatyGroupName"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 46342,
   "label": "Amount",
   "dari": "sisi",
   "kunci": "Currency",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1",
    "TreatyIn.IsEditData = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 73049,
   "label": "",
   "dari": "sisi",
   "kunci": "Amount",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1",
    "TreatyIn.IsEditData = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 94723,
   "label": "Text Area",
   "dari": "sisi",
   "kunci": "Note",
   "format": "pxTextArea",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.IsEditData = 1",
    "TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  }
 ],
 "MaxRetentionOldData": [
  {
   "t": "medan",
   "at": 15322,
   "label": "Treaty Group",
   "dari": "sisi",
   "kunci": "TreatyGroup",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseTreatyGroup_RD",
    "nilai": "TreatyGroupName"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 43362,
   "label": "Amount",
   "dari": "sisi",
   "kunci": "Currency",
   "format": "pxAutoComplete",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "reportdefinition",
    "rd": "BrowseCurrencyTreatyIn_RD",
    "nilai": "Currency"
   },
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 68494,
   "label": "",
   "dari": "sisi",
   "kunci": "Amount",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  },
  {
   "t": "medan",
   "at": 89497,
   "label": "Text Area",
   "dari": "sisi",
   "kunci": "Note",
   "format": "pxTextArea",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "aksiUbah": [
    {
     "aksi": "postValue"
    }
   ]
  }
 ],
 "Share": [
  {
   "t": "medan",
   "at": 29178,
   "label": "% RNM Share",
   "dari": "sisi",
   "kunci": "RNMShare",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
   ],
   "aksiUbah": [
    {
     "aksi": "refresh",
     "aktivitas": "TreatyInXOLAddSpreadingDetail",
     "param": {
      "idx": ".pxListSubscript"
     }
    },
    {
     "aksi": "refresh",
     "aktivitas": "TreatyInXOLAddSpreadingDetailActual",
     "param": {
      "idx": ".pxListSubscript"
     },
     "syarat": "TreatyIn.EDMState = 3"
    },
    {
     "aksi": "refresh"
    },
    {
     "aksi": "refresh"
    }
   ]
  },
  {
   "t": "medan",
   "at": 63179,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerType",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "associated"
   }
  },
  {
   "t": "medan",
   "at": 68845,
   "label": "",
   "dari": "sisi",
   "kunci": "Layer",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu"
  },
  {
   "t": "teks",
   "at": 74845,
   "teks": "Part of",
   "syarat": []
  },
  {
   "t": "medan",
   "at": 79347,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerPartType",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "associated"
   }
  },
  {
   "t": "medan",
   "at": 85023,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerPart",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu"
  },
  {
   "t": "blok",
   "at": 100108,
   "judul": "",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 117429,
     "prop": ".TreatyGroupList",
     "dari": "sisi",
     "larik": "TreatyGroupList",
     "syarat": [],
     "kolom": [
      "Treaty Group"
     ],
     "kunci": [
      "TreatyGroup"
     ],
     "lebar": [
      196
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
      125547
     ],
     "baca": [
      null
     ],
     "tombol": [
      null
     ],
     "tombolKepala": [
      null
     ],
     "pilihan": [
      {
       "sumber": "reportdefinition",
       "rd": "BrowseTreatyGroup_RD",
       "nilai": "TreatyGroupName"
      }
     ],
     "aksiUbah": [
      [
       {
        "aksi": "runDataTransform",
        "transformasi": "SetIndexLayer_DT",
        "paramDT": {
         "idxlimit": ""
        }
       },
       {
        "aksi": "refresh",
        "aktivitas": "TotalEgnpi",
        "param": {
         "idxLimit": ".Layer"
        }
       }
      ]
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitLayer!pyGridRowDetails",
     "rincian": "CoBListReadOnly"
    }
   ]
  },
  {
   "t": "medan",
   "at": 162558,
   "label": "Dropdown",
   "dari": "sisi",
   "kunci": "Cover",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "pilihan": {
    "sumber": "associated"
   }
  },
  {
   "t": "blok",
   "at": 171027,
   "judul": "Deduction Details",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 188456,
     "prop": ".DeductionList",
     "dari": "sisi",
     "larik": "DeductionList",
     "syarat": [],
     "kolom": [
      "Description",
      "Currency",
      "Deduction",
      "or",
      "Deduction %",
      "Auto Calculate %",
      ""
     ],
     "kunci": [
      "Comment",
      "Currency",
      "Deduction",
      "",
      "DeductionPct",
      "DeductionPctCalculate",
      ""
     ],
     "lebar": [
      203,
      195,
      199,
      30,
      210,
      117,
      118
     ],
     "desimal": [
      null,
      null,
      2,
      null,
      2,
      null,
      null
     ],
     "format": [
      "pxTextInput",
      "pxAutoComplete",
      "pxNumber",
      "",
      "pxNumber",
      "pxCheckbox",
      "pxButton"
     ],
     "syaratSel": [
      null,
      null,
      null,
      null,
      null,
      null,
      "TreatyIn.ViewState !='1'"
     ],
     "atSel": [
      226605,
      232723,
      243915,
      252953,
      255462,
      264370,
      269980
     ],
     "baca": [
      [
       "TreatyIn.ViewState = 1",
       "TreatyIn.EDMMaterialType = 2"
      ],
      [
       "TreatyIn.ViewState = 1",
       "TreatyIn.EDMMaterialType = 2"
      ],
      [
       "TreatyIn.ViewState = 1",
       "TreatyIn.EDMMaterialType = 2"
      ],
      null,
      [
       "TreatyIn.ViewState = 1",
       "TreatyIn.EDMMaterialType = 2"
      ],
      [
       "TreatyIn.ViewState = 1 || TreatyIn.EDMMaterialType = 2"
      ],
      "selalu"
     ],
     "tombol": [
      null,
      null,
      null,
      null,
      null,
      null,
      {
       "t": "tombol",
       "at": 269980,
       "label": "Remove",
       "syarat": [
        "TreatyIn.ViewState !='1'"
       ],
       "aksi": [
        {
         "aksi": "deleteRow"
        },
        {
         "aksi": "refresh",
         "aktivitas": "CalculateDeduction",
         "param": {
          "sts": "",
          "index": ".pxListSubscript"
         }
        }
       ],
       "nonaktif": [
        "TreatyIn.EDMMaterialType = 2"
       ]
      }
     ],
     "tombolKepala": [
      null,
      null,
      null,
      null,
      null,
      null,
      {
       "t": "tombol",
       "at": 217024,
       "label": "Add",
       "syarat": [
        "TreatyIn.ViewState !='1'"
       ],
       "aksi": [
        {
         "aksi": "refresh",
         "aktivitas": "AddDeduction"
        }
       ],
       "nonaktif": [
        "TreatyIn.EDMMaterialType = 2"
       ]
      }
     ],
     "pilihan": [
      null,
      {
       "sumber": "reportdefinition",
       "rd": "BrowseCurrency_RD",
       "nilai": "Currency"
      },
      null,
      null,
      null,
      null,
      null
     ],
     "aksiUbah": [
      null,
      [
       {
        "aksi": "refresh",
        "aktivitas": "CalculateDeduction",
        "param": {
         "sts": "",
         "index": ".pxListSubscript"
        }
       }
      ],
      [
       {
        "aksi": "refresh",
        "aktivitas": "CalculateDeduction",
        "param": {
         "sts": "val",
         "index": ".pxListSubscript"
        }
       }
      ],
      null,
      [
       {
        "aksi": "refresh",
        "aktivitas": "CalculateDeduction",
        "param": {
         "sts": "pct",
         "index": ".pxListSubscript"
        }
       }
      ],
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInDeduction!pyGridModalTemplate"
    }
   ]
  },
  {
   "t": "blok",
   "at": 320627,
   "judul": "",
   "syarat": [
    ".SpreadingTypeXOL!=''"
   ],
   "anak": [
    {
     "t": "medan",
     "at": 327353,
     "label": "Spreading Type",
     "dari": "sisi",
     "kunci": "SpreadingTypeXOL",
     "format": "pxDropdown",
     "desimal": null,
     "syarat": [],
     "baca": [
      "TreatyIn.ViewState = '1'"
     ],
     "pilihan": {
      "sumber": "reportdefinition",
      "rd": "BrowseTreatyArrangement_ParentReinsMasterTrt",
      "nilai": "ReinsTypeName",
      "tampil": "ReinsTypeName"
     },
     "aksiUbah": [
      {
       "aksi": "refresh",
       "aktivitas": "FetchQSfromMasterXOL",
       "param": {
        "ParentReinsTypeID": ".SpreadingTypeXOL",
        "indexShare": ".pxListSubscript",
        "IsUpdate": "0"
       }
      },
      {
       "aksi": "refresh"
      },
      {
       "aksi": "refresh"
      },
      {
       "aksi": "refresh"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 356200,
   "judul": "",
   "syarat": [
    ".SpreadingTypeXOL!=''"
   ],
   "anak": [
    {
     "t": "blok",
     "at": 365306,
     "judul": "Spreading",
     "syarat": [],
     "anak": [
      {
       "t": "blok",
       "at": 374389,
       "judul": "",
       "syarat": [
        ".SpreadingTypeXOL !=''"
       ],
       "anak": [
        {
         "t": "blok",
         "at": 383424,
         "judul": "",
         "syarat": [
          ".SpreadingTypeXOL!=''"
         ],
         "anak": [
          {
           "t": "grid",
           "at": 408855,
           "prop": ".SpreadingListXOL",
           "dari": "sisi",
           "larik": "SpreadingListXOL",
           "syarat": [],
           "kolom": [
            "Reins Type",
            "Pct"
           ],
           "kunci": [
            "ReinsTypeName",
            "Pct"
           ],
           "lebar": [
            142,
            192
           ],
           "desimal": [
            null,
            null
           ],
           "format": [
            "pxTextInput",
            ""
           ],
           "syaratSel": [
            null,
            null
           ],
           "atSel": [
            421971,
            425284
           ],
           "baca": [
            null,
            null
           ],
           "tombol": [
            null,
            null
           ],
           "tombolKepala": [
            null,
            null
           ],
           "pilihan": [
            null,
            null
           ],
           "aksiUbah": [
            null,
            null
           ],
           "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitsSpreading!pyGridModalTemplate"
          },
          {
           "t": "teks",
           "at": 442631,
           "teks": "Total Pct",
           "syarat": []
          }
         ]
        },
        {
         "t": "medan",
         "at": 473720,
         "label": "Spreading Total Pct",
         "dari": "sisi",
         "kunci": "SpreadingTotalPctXOL",
         "format": "Decimal",
         "desimal": null,
         "syarat": [],
         "baca": "selalu"
        },
        {
         "t": "teks",
         "at": 484197,
         "teks": "%",
         "syarat": []
        }
       ]
      }
     ]
    },
    {
     "t": "blok",
     "at": 516234,
     "judul": "",
     "syarat": [],
     "anak": [
      {
       "t": "grid",
       "at": 551558,
       "prop": ".RnmLimitList",
       "dari": "sisi",
       "larik": "RnmLimitList",
       "syarat": [],
       "kolom": [
        "RNM Limit",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        197,
        355
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
        563606,
        568563
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 616908,
       "prop": ".RNMSpreadedListXOL",
       "dari": "sisi",
       "larik": "RNMSpreadedListXOL",
       "syarat": [],
       "kolom": [
        "Spreading OR Limit",
        ""
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
        629276,
        634427
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 682835,
       "prop": ".RNMSpreadedListRIXOL",
       "dari": "sisi",
       "larik": "RNMSpreadedListRIXOL",
       "syarat": [],
       "kolom": [
        "Spreading R/I Limit",
        ""
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
        695206,
        700357
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 764104,
       "prop": ".GrossPremiumMinList",
       "dari": "sisi",
       "larik": "GrossPremiumMinList",
       "syarat": [],
       "kolom": [
        "Gross Min Premium",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        142,
        247
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
        776305,
        781262
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 829607,
       "prop": ".RNMSpreadedListGrossMinXOL",
       "dari": "sisi",
       "larik": "RNMSpreadedListGrossMinXOL",
       "syarat": [],
       "kolom": [
        "Spreading OR Min Premium",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        152,
        237
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
        841684,
        846641
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 894986,
       "prop": ".RNMSpreadedListGrossRIMinXOL",
       "dari": "sisi",
       "larik": "RNMSpreadedListGrossRIMinXOL",
       "syarat": [],
       "kolom": [
        "Spreading R/I Min Premium",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        152,
        237
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
        907066,
        912023
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 975707,
       "prop": ".GrossPremiumList",
       "dari": "sisi",
       "larik": "GrossPremiumList",
       "syarat": [],
       "kolom": [
        "Gross Premium (MDP)",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        195,
        355
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
        987907,
        992864
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 1041209,
       "prop": ".RNMSpreadedListGrossXOL",
       "dari": "sisi",
       "larik": "RNMSpreadedListGrossXOL",
       "syarat": [],
       "kolom": [
        "Spreading OR MDP",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        195,
        355
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
        1053275,
        1058232
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 1106577,
       "prop": ".RNMSpreadedListGrossRIXOL",
       "dari": "sisi",
       "larik": "RNMSpreadedListGrossRIXOL",
       "syarat": [],
       "kolom": [
        "Spreading R/I MDP",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        195,
        355
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
        1118646,
        1123603
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 1187287,
       "prop": ".DeductionTotalList",
       "dari": "sisi",
       "larik": "DeductionTotalList",
       "syarat": [],
       "kolom": [
        "Deductions",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        196,
        356
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
        1199612,
        1204763
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "@baseclass!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 1253369,
       "prop": ".RNMSpreadedListDeductXOL",
       "dari": "sisi",
       "larik": "RNMSpreadedListDeductXOL",
       "syarat": [],
       "kolom": [
        "Spreading OR Deductions",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        196,
        356
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
        1265734,
        1270885
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 1319491,
       "prop": ".RNMSpreadedListDeductRIXOL",
       "dari": "sisi",
       "larik": "RNMSpreadedListDeductRIXOL",
       "syarat": [],
       "kolom": [
        "Spreading R/I Deductions",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        196,
        356
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
        1331859,
        1337010
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 1400955,
       "prop": ".NetPremiumList",
       "dari": "sisi",
       "larik": "NetPremiumList",
       "syarat": [],
       "kolom": [
        "Net Premium",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        195,
        354
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
        1413007,
        1417964
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 1466309,
       "prop": ".RNMSpreadedListNetXOL",
       "dari": "sisi",
       "larik": "RNMSpreadedListNetXOL",
       "syarat": [],
       "kolom": [
        "Spreading OR Net Premi",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        195,
        354
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
        1478379,
        1483336
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      },
      {
       "t": "grid",
       "at": 1531681,
       "prop": ".RNMSpreadedListNetRIXOL",
       "dari": "sisi",
       "larik": "RNMSpreadedListNetRIXOL",
       "syarat": [],
       "kolom": [
        "Spreading R/I Net Premium",
        ""
       ],
       "kunci": [
        "Currency",
        "Value"
       ],
       "lebar": [
        195,
        354
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
        1543756,
        1548713
       ],
       "baca": [
        null,
        null
       ],
       "tombol": [
        null,
        null
       ],
       "tombolKepala": [
        null,
        null
       ],
       "pilihan": [
        null,
        null
       ],
       "aksiUbah": [
        null,
        null
       ],
       "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 1605909,
   "judul": "",
   "syarat": [
    ".SpreadingTypeXOL= ''"
   ],
   "anak": [
    {
     "t": "blok",
     "at": 1615006,
     "judul": "Spreading",
     "syarat": [],
     "anak": [
      {
       "t": "blok",
       "at": 1624090,
       "judul": "",
       "syarat": [
        ".SpreadingTypeXOL =''"
       ],
       "anak": [
        {
         "t": "blok",
         "at": 1633125,
         "judul": "",
         "syarat": [
          ".SpreadingTypeXOL==''"
         ],
         "anak": [
          {
           "t": "grid",
           "at": 1658559,
           "prop": ".SpreadingListXOL",
           "dari": "sisi",
           "larik": "SpreadingListXOL",
           "syarat": [],
           "kolom": [
            "Reins Type",
            "Pct Share",
            ".pyTemplateInputBox"
           ],
           "kunci": [
            "ReinsTypeID",
            "Pct",
            "pyTemplateInputBox"
           ],
           "lebar": [
            195,
            192,
            80
           ],
           "desimal": [
            null,
            2,
            null
           ],
           "format": [
            "pxDropdown",
            "pxNumber",
            "pxLink"
           ],
           "syaratSel": [
            null,
            null,
            "TreatyIn.ViewState != '1'"
           ],
           "atSel": [
            1681285,
            1697979,
            1707977
           ],
           "baca": [
            [
             "TreatyIn.ViewState = '1'"
            ],
            [
             "TreatyIn.ViewState = '1'"
            ],
            null
           ],
           "tombol": [
            null,
            null,
            null
           ],
           "tombolKepala": [
            null,
            null,
            null
           ],
           "pilihan": [
            {
             "sumber": "reportdefinition",
             "rd": "BrowseTreatyArrangement_ParentReinsMasterTrt",
             "nilai": "ReinsTypeID",
             "tampil": "ReinsTypeName"
            },
            null,
            null
           ],
           "aksiUbah": [
            [
             {
              "aksi": "postValue"
             },
             {
              "aksi": "refresh"
             }
            ],
            [
             {
              "aksi": "postValue"
             },
             {
              "aksi": "refresh",
              "aktivitas": "SetSpreadingXOL"
             }
            ],
            null
           ],
           "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitsSpreading!pyGridModalTemplate"
          },
          {
           "t": "teks",
           "at": 1729251,
           "teks": "Total Pct",
           "syarat": []
          }
         ]
        },
        {
         "t": "medan",
         "at": 1760513,
         "label": "Total Spreading Pct :",
         "dari": "sisi",
         "kunci": "SpreadingTotalPctXOL",
         "format": "pxNumber",
         "desimal": 2,
         "syarat": [
          ".SpreadingTypeXOL != ''"
         ],
         "baca": "selalu"
        },
        {
         "t": "medan",
         "at": 1767848,
         "label": "Total Share Pct",
         "dari": "sisi",
         "kunci": "SpreadingTotalPctXOL",
         "format": "pxNumber",
         "desimal": 2,
         "syarat": [
          ".SpreadingTypeXOL = ''"
         ],
         "baca": "selalu"
        }
       ]
      }
     ]
    }
   ]
  }
 ],
 "ShareOldData": [
  {
   "t": "medan",
   "at": 23993,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerType",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "associated"
   }
  },
  {
   "t": "medan",
   "at": 29603,
   "label": "",
   "dari": "sisi",
   "kunci": "Layer",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu"
  },
  {
   "t": "teks",
   "at": 35175,
   "teks": "Part of",
   "syarat": []
  },
  {
   "t": "medan",
   "at": 39656,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerPartType",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "associated"
   }
  },
  {
   "t": "medan",
   "at": 45276,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerPart",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu"
  },
  {
   "t": "blok",
   "at": 59882,
   "judul": "",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 76628,
     "prop": ".TreatyGroupList",
     "dari": "sisi",
     "larik": "TreatyGroupList",
     "syarat": [],
     "kolom": [
      "Treaty Group"
     ],
     "kunci": [
      "TreatyGroup"
     ],
     "lebar": [
      196
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
      84542
     ],
     "baca": [
      null
     ],
     "tombol": [
      null
     ],
     "tombolKepala": [
      null
     ],
     "pilihan": [
      {
       "sumber": "reportdefinition",
       "rd": "BrowseTreatyGroup_RD",
       "nilai": "TreatyGroupName"
      }
     ],
     "aksiUbah": [
      [
       {
        "aksi": "runDataTransform",
        "transformasi": "SetIndexLayer_DT",
        "paramDT": {
         "idxlimit": ""
        }
       },
       {
        "aksi": "refresh",
        "aktivitas": "TotalEgnpi",
        "param": {
         "idxLimit": ".Layer"
        }
       }
      ]
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitLayer!pyGridRowDetails",
     "rincian": "CoBListReadOnly"
    }
   ]
  },
  {
   "t": "medan",
   "at": 121204,
   "label": "Dropdown",
   "dari": "sisi",
   "kunci": "Cover",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "pilihan": {
    "sumber": "associated"
   }
  },
  {
   "t": "blok",
   "at": 129595,
   "judul": "Deduction Details",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 146358,
     "prop": ".DeductionList",
     "dari": "sisi",
     "larik": "DeductionList",
     "syarat": [],
     "kolom": [
      "Description",
      "Currency",
      "Deduction",
      "or",
      "Deduction %",
      "Auto Calculate %"
     ],
     "kunci": [
      "Comment",
      "Currency",
      "Deduction",
      "",
      "DeductionPct",
      "DeductionPctCalculate"
     ],
     "lebar": [
      203,
      195,
      199,
      30,
      210,
      117
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
      "pxTextInput",
      "pxAutoComplete",
      "pxNumber",
      "",
      "pxNumber",
      "pxCheckbox"
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
      174081,
      179451,
      190182,
      198728,
      201238,
      209654
     ],
     "baca": [
      "selalu",
      "selalu",
      "selalu",
      null,
      "selalu",
      "selalu"
     ],
     "tombol": [
      null,
      null,
      null,
      null,
      null,
      null
     ],
     "tombolKepala": [
      null,
      null,
      null,
      null,
      null,
      null
     ],
     "pilihan": [
      null,
      {
       "sumber": "reportdefinition",
       "rd": "BrowseCurrency_RD",
       "nilai": "Currency"
      },
      null,
      null,
      null,
      null
     ],
     "aksiUbah": [
      null,
      [
       {
        "aksi": "refresh",
        "aktivitas": "CalculateDeduction",
        "param": {
         "sts": "",
         "index": ".pxListSubscript"
        }
       }
      ],
      [
       {
        "aksi": "refresh",
        "aktivitas": "CalculateDeduction",
        "param": {
         "sts": "val",
         "index": ".pxListSubscript"
        }
       }
      ],
      null,
      [
       {
        "aksi": "refresh",
        "aktivitas": "CalculateDeduction",
        "param": {
         "sts": "pct",
         "index": ".pxListSubscript"
        }
       }
      ],
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInDeduction!pyGridModalTemplate"
    }
   ]
  },
  {
   "t": "blok",
   "at": 252022,
   "judul": "Spreading",
   "syarat": [],
   "anak": [
    {
     "t": "medan",
     "at": 267038,
     "label": "Spreading Type",
     "dari": "sisi",
     "kunci": "SpreadingTypeXOL",
     "format": "pxDropdown",
     "desimal": null,
     "syarat": [],
     "baca": "selalu",
     "pilihan": {
      "sumber": "reportdefinition",
      "rd": "BrowseTreatyArrangement_ParentReinsMasterTrt",
      "nilai": "ReinsTypeID",
      "tampil": "ReinsTypeName"
     }
    },
    {
     "t": "blok",
     "at": 283028,
     "judul": "",
     "syarat": [
      ".SpreadingTypeXOL !=''"
     ],
     "anak": [
      {
       "t": "blok",
       "at": 291733,
       "judul": "",
       "syarat": [],
       "anak": [
        {
         "t": "grid",
         "at": 315668,
         "prop": ".SpreadingListXOL",
         "dari": "sisi",
         "larik": "SpreadingListXOL",
         "syarat": [],
         "kolom": [
          "Reins Type",
          "Pct"
         ],
         "kunci": [
          "ReinsTypeName",
          "Pct"
         ],
         "lebar": [
          142,
          192
         ],
         "desimal": [
          null,
          null
         ],
         "format": [
          "pxTextInput",
          ""
         ],
         "syaratSel": [
          null,
          null
         ],
         "atSel": [
          328536,
          331153
         ],
         "baca": [
          null,
          null
         ],
         "tombol": [
          null,
          null
         ],
         "tombolKepala": [
          null,
          null
         ],
         "pilihan": [
          null,
          null
         ],
         "aksiUbah": [
          null,
          null
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitsSpreading!pyGridModalTemplate"
        },
        {
         "t": "teks",
         "at": 347978,
         "teks": "Total Pct",
         "syarat": []
        }
       ]
      },
      {
       "t": "medan",
       "at": 377989,
       "label": "Spreading Total Pct",
       "dari": "sisi",
       "kunci": "SpreadingTotalPctXOL",
       "format": "Decimal",
       "desimal": null,
       "syarat": [],
       "baca": "selalu"
      },
      {
       "t": "teks",
       "at": 388385,
       "teks": "%",
       "syarat": []
      }
     ]
    }
   ]
  },
  {
   "t": "blok",
   "at": 420262,
   "judul": "",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 454340,
     "prop": ".RnmLimitList",
     "dari": "sisi",
     "larik": "RnmLimitList",
     "syarat": [],
     "kolom": [
      "RNM Limit",
      ""
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      197,
      355
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
      465972,
      470686
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 518184,
     "prop": ".RNMSpreadedListXOL",
     "dari": "sisi",
     "larik": "RNMSpreadedListXOL",
     "syarat": [],
     "kolom": [
      "Spreading OR Limit",
      ""
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
      530137,
      535045
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 582607,
     "prop": ".RNMSpreadedListRIXOL",
     "dari": "sisi",
     "larik": "RNMSpreadedListRIXOL",
     "syarat": [],
     "kolom": [
      "Spreading R/I Limit",
      ""
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
      594563,
      599471
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInShare!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 662013,
     "prop": ".GrossPremiumList",
     "dari": "sisi",
     "larik": "GrossPremiumList",
     "syarat": [],
     "kolom": [
      "Gross Premium (MDP)",
      ""
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      195,
      355
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
      673660,
      678374
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 725873,
     "prop": ".RNMSpreadedListGrossXOL",
     "dari": "sisi",
     "larik": "RNMSpreadedListGrossXOL",
     "syarat": [],
     "kolom": [
      "Spreading OR MDP",
      ""
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      195,
      355
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
      737524,
      742238
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 789737,
     "prop": ".RNMSpreadedListGrossRIXOL",
     "dari": "sisi",
     "larik": "RNMSpreadedListGrossRIXOL",
     "syarat": [],
     "kolom": [
      "Spreading R/I MDP",
      ""
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      195,
      355
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
      801391,
      806105
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 868584,
     "prop": ".DeductionTotalList",
     "dari": "sisi",
     "larik": "DeductionTotalList",
     "syarat": [],
     "kolom": [
      "Deductions",
      ""
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      196,
      356
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
      880494,
      885402
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "@baseclass!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 933162,
     "prop": ".RNMSpreadedListDeductXOL",
     "dari": "sisi",
     "larik": "RNMSpreadedListDeductXOL",
     "syarat": [],
     "kolom": [
      "Spreading OR Deductions",
      ""
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      196,
      356
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
      945112,
      950020
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 997780,
     "prop": ".RNMSpreadedListDeductRIXOL",
     "dari": "sisi",
     "larik": "RNMSpreadedListDeductRIXOL",
     "syarat": [],
     "kolom": [
      "Spreading R/I Deductions",
      ""
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      196,
      356
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
      1009733,
      1014641
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 1077381,
     "prop": ".NetPremiumList",
     "dari": "sisi",
     "larik": "NetPremiumList",
     "syarat": [],
     "kolom": [
      "Net Premium",
      ""
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      195,
      354
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
      1089018,
      1093732
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 1141231,
     "prop": ".RNMSpreadedListNetXOL",
     "dari": "sisi",
     "larik": "RNMSpreadedListNetXOL",
     "syarat": [],
     "kolom": [
      "Spreading OR Net Premi",
      ""
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      195,
      354
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
      1152886,
      1157600
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    },
    {
     "t": "grid",
     "at": 1205099,
     "prop": ".RNMSpreadedListNetRIXOL",
     "dari": "sisi",
     "larik": "RNMSpreadedListNetRIXOL",
     "syarat": [],
     "kolom": [
      "Spreading R/I Net Premium",
      ""
     ],
     "kunci": [
      "Currency",
      "Value"
     ],
     "lebar": [
      195,
      354
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
      1216759,
      1221473
     ],
     "baca": [
      null,
      null
     ],
     "tombol": [
      null,
      null
     ],
     "tombolKepala": [
      null,
      null
     ],
     "pilihan": [
      null,
      null
     ],
     "aksiUbah": [
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInTotal!pyGridModalTemplate"
    }
   ]
  }
 ],
 "ShareRetro": [
  {
   "t": "medan",
   "at": 26960,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerType",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "associated"
   }
  },
  {
   "t": "medan",
   "at": 32568,
   "label": "",
   "dari": "sisi",
   "kunci": "Layer",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu"
  },
  {
   "t": "teks",
   "at": 38138,
   "teks": "Part of",
   "syarat": []
  },
  {
   "t": "medan",
   "at": 42617,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerPartType",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": "selalu",
   "pilihan": {
    "sumber": "associated"
   }
  },
  {
   "t": "medan",
   "at": 48235,
   "label": "",
   "dari": "sisi",
   "kunci": "LayerPart",
   "format": "pxTextInput",
   "desimal": null,
   "syarat": [],
   "baca": "selalu"
  },
  {
   "t": "blok",
   "at": 62837,
   "judul": "",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 79574,
     "prop": ".TreatyGroupList",
     "dari": "sisi",
     "larik": "TreatyGroupList",
     "syarat": [],
     "kolom": [
      "Treaty Group"
     ],
     "kunci": [
      "TreatyGroup"
     ],
     "lebar": [
      196
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
      87484
     ],
     "baca": [
      null
     ],
     "tombol": [
      null
     ],
     "tombolKepala": [
      null
     ],
     "pilihan": [
      {
       "sumber": "reportdefinition",
       "rd": "BrowseTreatyGroup_RD",
       "nilai": "TreatyGroupName"
      }
     ],
     "aksiUbah": [
      [
       {
        "aksi": "runDataTransform",
        "transformasi": "SetIndexLayer_DT",
        "paramDT": {
         "idxlimit": ""
        }
       },
       {
        "aksi": "refresh",
        "aktivitas": "TotalEgnpi",
        "param": {
         "idxLimit": ".Layer"
        }
       }
      ]
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitLayer!pyGridRowDetails",
     "rincian": "CoBListReadOnly"
    }
   ]
  },
  {
   "t": "medan",
   "at": 124136,
   "label": "Dropdown",
   "dari": "sisi",
   "kunci": "Cover",
   "format": "pxDropdown",
   "desimal": null,
   "syarat": [],
   "baca": [
    "TreatyIn.ViewState = 1"
   ],
   "pilihan": {
    "sumber": "associated"
   }
  },
  {
   "t": "blok",
   "at": 132525,
   "judul": "Deduction Details",
   "syarat": [],
   "anak": [
    {
     "t": "grid",
     "at": 149279,
     "prop": ".DeductionList",
     "dari": "sisi",
     "larik": "DeductionList",
     "syarat": [],
     "kolom": [
      "Description",
      "Currency",
      "Deduction",
      "or",
      "Deduction %",
      "Auto Calculate %",
      ""
     ],
     "kunci": [
      "Comment",
      "Currency",
      "Deduction",
      "",
      "DeductionPct",
      "DeductionPctCalculate",
      ""
     ],
     "lebar": [
      203,
      196,
      199,
      31,
      211,
      117,
      118
     ],
     "desimal": [
      null,
      null,
      2,
      null,
      2,
      null,
      null
     ],
     "format": [
      "pxTextInput",
      "pxAutoComplete",
      "pxNumber",
      "",
      "pxNumber",
      "pxCheckbox",
      "pxButton"
     ],
     "syaratSel": [
      null,
      null,
      null,
      null,
      null,
      null,
      "TreatyIn.ViewState !='1'"
     ],
     "atSel": [
      185782,
      191242,
      202032,
      210668,
      213177,
      221683,
      227187
     ],
     "baca": [
      "selalu",
      "selalu",
      "selalu",
      null,
      "selalu",
      "selalu",
      "selalu"
     ],
     "tombol": [
      null,
      null,
      null,
      null,
      null,
      null,
      {
       "t": "tombol",
       "at": 227187,
       "label": "Remove",
       "syarat": [
        "TreatyIn.ViewState !='1'"
       ],
       "aksi": [
        {
         "aksi": "deleteRow"
        },
        {
         "aksi": "refresh",
         "aktivitas": "CalculateDeduction",
         "param": {
          "sts": "",
          "index": ".pxListSubscript"
         }
        }
       ],
       "nonaktif": [
        "TreatyIn.EDMMaterialType = 2"
       ]
      }
     ],
     "tombolKepala": [
      null,
      null,
      null,
      null,
      null,
      null,
      {
       "t": "tombol",
       "at": 176597,
       "label": "Add",
       "syarat": [
        "TreatyIn.ViewState !='1'"
       ],
       "aksi": [
        {
         "aksi": "refresh",
         "aktivitas": "AddDeduction"
        }
       ],
       "nonaktif": [
        "TreatyIn.EDMMaterialType = 2"
       ]
      }
     ],
     "pilihan": [
      null,
      {
       "sumber": "reportdefinition",
       "rd": "BrowseCurrency_RD",
       "nilai": "Currency"
      },
      null,
      null,
      null,
      null,
      null
     ],
     "aksiUbah": [
      null,
      [
       {
        "aksi": "refresh",
        "aktivitas": "CalculateDeduction",
        "param": {
         "sts": "",
         "index": ".pxListSubscript"
        }
       }
      ],
      [
       {
        "aksi": "refresh",
        "aktivitas": "CalculateDeduction",
        "param": {
         "sts": "val",
         "index": ".pxListSubscript"
        }
       }
      ],
      null,
      [
       {
        "aksi": "refresh",
        "aktivitas": "CalculateDeduction",
        "param": {
         "sts": "pct",
         "index": ".pxListSubscript"
        }
       }
      ],
      null,
      null
     ],
     "templatBaris": "ASM-FW-GISFW-Data-TreatyInDeduction!pyGridModalTemplate"
    }
   ]
  },
  {
   "t": "blok",
   "at": 1310257,
   "judul": "Spreading",
   "syarat": [],
   "anak": [
    {
     "t": "medan",
     "at": 1325294,
     "label": "Spreading Type",
     "dari": "sisi",
     "kunci": "SpreadingTypeXOLRetro",
     "format": "pxAutoComplete",
     "desimal": null,
     "syarat": [],
     "baca": [
      "TreatyIn.ViewState = 1"
     ],
     "pilihan": {
      "sumber": "reportdefinition",
      "rd": "BrowseReinsuranceType_RD",
      "nilai": "Note"
     },
     "aksiUbah": [
      {
       "aksi": "refresh",
       "aktivitas": "DelSpreadingRetro"
      },
      {
       "aksi": "refresh",
       "aktivitas": "GetSpreadingRetro"
      },
      {
       "aksi": "refresh"
      }
     ]
    },
    {
     "t": "blok",
     "at": 1350522,
     "judul": "",
     "syarat": [
      ".SpreadingTypeXOL !=''"
     ],
     "anak": [
      {
       "t": "blok",
       "at": 1359227,
       "judul": "",
       "syarat": [
        "TreatyIn.FacultativeShare >0"
       ],
       "anak": [
        {
         "t": "grid",
         "at": 1376025,
         "prop": ".ShareFacultativeReinsurers",
         "dari": "sisi",
         "larik": "ShareFacultativeReinsurers",
         "syarat": [],
         "kolom": [
          "Other Treaty Retro",
          "% Share",
          "Amount",
          "Amount2"
         ],
         "kunci": [
          "ReinsName",
          "SharePct",
          "Amount",
          "Amount2"
         ],
         "lebar": [
          514,
          85,
          100,
          100
         ],
         "desimal": [
          null,
          2,
          2,
          2
         ],
         "format": [
          "",
          "pxNumber",
          "pxNumber",
          "pxNumber"
         ],
         "syaratSel": [
          null,
          null,
          null,
          null
         ],
         "atSel": [
          1395296,
          1399551,
          1404789,
          1409597
         ],
         "baca": [
          "selalu",
          "selalu",
          "selalu",
          "selalu"
         ],
         "tombol": [
          null,
          null,
          null,
          null
         ],
         "tombolKepala": [
          null,
          null,
          null,
          null
         ],
         "pilihan": [
          null,
          null,
          null,
          null
         ],
         "aksiUbah": [
          null,
          null,
          null,
          null
         ],
         "templatBaris": "ASM-FW-GISFW-Data-TreatyInShareReins!pyGridModalTemplate"
        }
       ]
      }
     ]
    }
   ]
  }
 ],
 "TotalLimits": [
  {
   "t": "grid",
   "at": 35085,
   "prop": ".Detail",
   "dari": "sisi",
   "larik": "Detail",
   "syarat": [],
   "kolom": [
    "Treaty Group",
    "% RNM Share",
    ""
   ],
   "kunci": [
    "TreatyGroup",
    "RNMShare",
    "ShareNote"
   ],
   "lebar": [
    399,
    115,
    101
   ],
   "desimal": [
    null,
    2,
    null
   ],
   "format": [
    "pxAutoComplete",
    "pxNumber",
    ""
   ],
   "syaratSel": [
    null,
    null,
    null
   ],
   "atSel": [
    49235,
    56592,
    61354
   ],
   "baca": [
    null,
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    null,
    null
   ],
   "tombolKepala": [
    null,
    null,
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseBusinessGroup_RD",
     "nilai": "Note"
    },
    null,
    null
   ],
   "aksiUbah": [
    null,
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitsDetail!pyGridRowDetails",
   "rincian": "DetailShare"
  }
 ],
 "TotalLimitsOldData": [
  {
   "t": "grid",
   "at": 35890,
   "prop": ".Detail",
   "dari": "sisi",
   "larik": "Detail",
   "syarat": [],
   "kolom": [
    "Treaty Group",
    "% RNM Share",
    ""
   ],
   "kunci": [
    "TreatyGroup",
    "RNMShare",
    "ShareNote"
   ],
   "lebar": [
    399,
    115,
    101
   ],
   "desimal": [
    null,
    2,
    null
   ],
   "format": [
    "pxAutoComplete",
    "pxNumber",
    ""
   ],
   "syaratSel": [
    null,
    null,
    null
   ],
   "atSel": [
    50456,
    58056,
    63061
   ],
   "baca": [
    null,
    "selalu",
    "selalu"
   ],
   "tombol": [
    null,
    null,
    null
   ],
   "tombolKepala": [
    null,
    null,
    null
   ],
   "pilihan": [
    {
     "sumber": "reportdefinition",
     "rd": "BrowseBusinessGroup_RD",
     "nilai": "Note"
    },
    null,
    null
   ],
   "aksiUbah": [
    null,
    null,
    null
   ],
   "templatBaris": "ASM-FW-GISFW-Data-TreatyInLimitsDetail!pyGridRowDetails",
   "rincian": "DetailShareOldData"
  }
 ]
}

export const DIBUANG: readonly Terbuang[] = [
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
  "alasan": "penjaga mati: ['1=2']"
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
  "at": 541585,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
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
  "berkas": "TreatyInTabsNonProportionalAdjustPremi",
  "at": 302231,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalAdjustPremi",
  "at": 308502,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalAdjustPremi",
  "at": 346023,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalAdjustPremi",
  "at": 352294,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInTabsNonProportionalAdjustPremi",
  "at": 358566,
  "jenis": "sel",
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
  "alasan": "penjaga mati: ['1=2']"
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
  "alasan": "penjaga mati: ['1=2']"
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
  "alasan": "penjaga mati: ['1=2']"
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
  "berkas": "TreatyInActualLimits",
  "at": 16459,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "TreatyInActualLimits",
  "at": 338382,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInActualLimits",
  "at": 398786,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInActualShare",
  "at": 24741,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "TreatyInActualShare",
  "at": 228677,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInActualShare",
  "at": 875414,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInActualShare",
  "at": 909879,
  "jenis": "blok",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "TreatyInActualShare",
  "at": 967465,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInActualSumary",
  "at": 577986,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInActualSumary",
  "at": 639729,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInActualSumary",
  "at": 1307311,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInActualSumary",
  "at": 1336924,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "TreatyInActualSumary",
  "at": 1901477,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailEGNPI",
  "at": 66677,
  "jenis": "medan",
  "alasan": "kontrol tanpa properti: '.pyTemplateRichTextEditor'"
 },
 {
  "berkas": "Layers",
  "at": 345431,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Layers",
  "at": 350741,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Layers",
  "at": 404406,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Layers",
  "at": 409118,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Layers",
  "at": 461752,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Layers",
  "at": 466466,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Layers",
  "at": 535083,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Layers",
  "at": 540394,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Layers",
  "at": 594068,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Layers",
  "at": 598781,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Layers",
  "at": 651470,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Layers",
  "at": 656185,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Layers",
  "at": 709081,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Layers",
  "at": 722805,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Layers",
  "at": 728244,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "Layers",
  "at": 736484,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "Share",
  "at": 390383,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "Share",
  "at": 398993,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Share",
  "at": 446688,
  "jenis": "medan",
  "alasan": "kontrol tanpa properti: '.pyTemplateInputBox'"
 },
 {
  "berkas": "Share",
  "at": 478874,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "Share",
  "at": 1640084,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "Share",
  "at": 1648694,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "Share",
  "at": 1733308,
  "jenis": "medan",
  "alasan": "kontrol tanpa properti: '.pyTemplateInputBox'"
 },
 {
  "berkas": "Share",
  "at": 1775209,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "Share",
  "at": 1813015,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailEGNPIOldData",
  "at": 65431,
  "jenis": "medan",
  "alasan": "kontrol tanpa properti: '.pyTemplateRichTextEditor'"
 },
 {
  "berkas": "LayersOldData",
  "at": 317772,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersOldData",
  "at": 323084,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersOldData",
  "at": 376589,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersOldData",
  "at": 381303,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersOldData",
  "at": 433777,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersOldData",
  "at": 438493,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersOldData",
  "at": 506953,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersOldData",
  "at": 512265,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersOldData",
  "at": 565774,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersOldData",
  "at": 570488,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersOldData",
  "at": 623012,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersOldData",
  "at": 627728,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersOldData",
  "at": 658177,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersOldData",
  "at": 671817,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersOldData",
  "at": 676759,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "LayersOldData",
  "at": 684946,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "AchievementCombine",
  "at": 126286,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "ShareOldData",
  "at": 298309,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "ShareOldData",
  "at": 306595,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "ShareOldData",
  "at": 351827,
  "jenis": "medan",
  "alasan": "kontrol tanpa properti: '.pyTemplateInputBox'"
 },
 {
  "berkas": "ShareOldData",
  "at": 383120,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "ShareOldData",
  "at": 1269022,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 318056,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 323366,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 376859,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 381571,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 434033,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 438747,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 507192,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 512502,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 566002,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 570715,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 623232,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 627947,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 658392,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 672030,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "LayersEDM",
  "at": 676971,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "LayersEDM",
  "at": 685157,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "ShareRetro",
  "at": 267051,
  "jenis": "blok",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "ShareRetro",
  "at": 1441136,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "ShareRetro",
  "at": 1503154,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "CoBList",
  "at": 79969,
  "jenis": "kolom",
  "alasan": "penjaga mati: [\"TreatyIn.ViewState !='1' && 1=2\"]"
 },
 {
  "berkas": "DetailLimits",
  "at": 61842,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 116567,
  "jenis": "kolom",
  "alasan": "penjaga mati: [\"TreatyIn.ViewState !='1' && 1=2\"]"
 },
 {
  "berkas": "DetailLimits",
  "at": 162411,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 270223,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 446824,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 629183,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 805387,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 865479,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 870301,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 922184,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 927006,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 978893,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 983717,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 1071521,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 1353933,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 1376700,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 1532513,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 1600358,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 1623125,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 1778526,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 1801292,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 1956725,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 1979492,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 2210908,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 2233746,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 2389061,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 2568023,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 2637322,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 2724996,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 2901172,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 2907715,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 2914260,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimits",
  "at": 2955111,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 46714,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 59779,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 141350,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 248657,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 393813,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 545778,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 715777,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 776585,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 781420,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 833848,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 838683,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 890516,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 895353,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 986084,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 1251211,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 1274049,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 1412793,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 1483878,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 1506714,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 1645493,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 1668330,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 1806856,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 1829694,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 1967761,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 1990599,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 2128955,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 2267895,
  "jenis": "blok",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 2344009,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 2348269,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 2352529,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 2356789,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 2361049,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 2376736,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 2380996,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 2396639,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 2400901,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailLimitsOldData",
  "at": 2422888,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailShareOldData",
  "at": 36857,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailShareOldData",
  "at": 189876,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "DetailShareOldData",
  "at": 198485,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailShareOldData",
  "at": 248838,
  "jenis": "medan",
  "alasan": "kontrol tanpa properti: '.pyTemplateInputBox'"
 },
 {
  "berkas": "DetailShareOldData",
  "at": 267333,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "DetailShareOldData",
  "at": 275942,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailShareOldData",
  "at": 333695,
  "jenis": "kolom",
  "alasan": "penjaga mati: ['Never']"
 },
 {
  "berkas": "DetailShareOldData",
  "at": 359626,
  "jenis": "medan",
  "alasan": "kontrol tanpa properti: '.pyTemplateInputBox'"
 },
 {
  "berkas": "DetailShareOldData",
  "at": 401474,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "DetailShare",
  "at": 39130,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailShare",
  "at": 191651,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "DetailShare",
  "at": 200260,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailShare",
  "at": 250574,
  "jenis": "medan",
  "alasan": "kontrol tanpa properti: '.pyTemplateInputBox'"
 },
 {
  "berkas": "DetailShare",
  "at": 269069,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 },
 {
  "berkas": "DetailShare",
  "at": 277678,
  "jenis": "sel",
  "alasan": "penjaga mati: 1=2"
 },
 {
  "berkas": "DetailShare",
  "at": 361490,
  "jenis": "medan",
  "alasan": "kontrol tanpa properti: '.pyTemplateInputBox'"
 },
 {
  "berkas": "DetailShare",
  "at": 403338,
  "jenis": "sel",
  "alasan": "penjaga mati: NEVER"
 }
]
