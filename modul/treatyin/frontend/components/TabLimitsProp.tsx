// Tab **Limits** cabang PROPORSIONAL — tiga tingkat yang dapat dibuka,
// dibaca dari ekspor (`labelsLimitsProp.ts`).
//
//   Kind of Treaty  (grid, `Limits[]`)          expandPane → LimitProportional
//   └ Treaty Group  (grid, `Limits[].Detail[]`) expandPane → DetailLimits
//     └ medan · grid · sebelas tab
//
// ⛔ MENGGANTIKAN `PohonLimits` untuk cabang ini saja. Pohon lama memasang
// medan LAYER non-prop (Cover, MDP, ROL …) di bawah `Kind of Treaty`; Section
// prop tidak punya satu pun, sehingga layar berisi deretan "—". Cabang
// non-prop tetap memakai `PohonLimits`.
//
// ⭐ Mode `ubah`: seluruh medan dapat disunting dan tiap grid ber-`Add`
// punya `Add`/`Delete` yang bekerja — atas salinan pohon di layar ini.
// Mode `lihat`: baca-saja, tombol tidak dirender (`TreatyIn.ViewState !='1'`).

import {
  createContext,
  useContext,
  useEffect,
  useState,
  type ReactNode,
} from "react";

import {
  Field,
  Kosong,
  Panel,
  Pilih,
  StripTab,
} from "../../../../inti/frontend/components/ui/dasar";
import { formatNumber } from "../../../../inti/frontend/lib/format";
import {
  ambilOpsiLimits,
  type OpsiLimits,
  type PilihanWarisan,
  type SimpulLimit,
} from "../api";
import type { GolonganAngka } from "../labels";
import {
  GRID_DETAIL,
  KOLOM_KIND_OF_TREATY,
  KOLOM_TREATY_GROUP,
  LIMITS_PROP,
  TAB_DETAIL_LIMIT,
  syaratQS,
  syaratSurplus,
  type KolomLimit,
  type MedanLimit,
} from "../labelsLimitsProp";
import type { ModeForm } from "../mode";
import { selAngka } from "./angka";

/** Teks satu medan simpul — larik dan kunci yang tidak ada → kosong. */
export function teksDari(s: SimpulLimit, kunci: string): string {
  const v = s[kunci];
  return typeof v === "string" ? v : "";
}

/** Larik di simpul — yang tidak ada → kosong. */
function larikDari(s: SimpulLimit, kunci: string): SimpulLimit[] {
  const v = s[kunci];
  return Array.isArray(v) ? v : [];
}

/**
 * Format satu nilai menurut golongan + desimal ekspor.
 *
 * ⛔ Pemformatnya `selAngka` modul (memanggil `format.ts` inti) — bukan
 * pemformat kedua. `-1` = `pyDecimalPlaces = -999`, tak dibatasi.
 */
export function formatLimit(
  golongan: GolonganAngka,
  desimal: number | null,
  nilai: string,
): string {
  if (desimal === -1) return formatNumber(nilai, -1);
  return selAngka(desimal === null ? golongan : [golongan, desimal], nilai);
}

/** Isi dropdown tab Limits — diberikan sekali di akar, dibaca tiap tingkat. */
const OPSI_KOSONG: OpsiLimits = {
  jenisTreaty: [],
  kelompokTreaty: [],
  mataUang: [],
};
const OpsiLimitsCtx = createContext<OpsiLimits>(OPSI_KOSONG);

/** Pilihan untuk `Pilih`; nama kembar diberi kodenya. */
function opsiDari(daftar: readonly PilihanWarisan[], nilai: "id" | "nama") {
  const lihat = new Set<string>();
  const keluar: { value: string; label: string }[] = [];
  for (const o of daftar) {
    const v = nilai === "id" ? o.id : o.nama;
    if (lihat.has(v)) continue;
    lihat.add(v);
    keluar.push({
      value: v,
      label:
        nilai === "id" && o.kembar
          ? `${o.nama} (${o.id} — ${LIMITS_PROP.namaKembar})`
          : o.nama,
    });
  }
  return keluar;
}

/**
 * Dropdown terikat kode (`.TreatyTypeID`, `.TreatyGroupID`, `.CurrencyID`):
 * mode ubah `Pilih`; mode lihat nama baca-saja, seperti dropdown read-only Pega.
 * Memilih mengisi nama pasangannya — `Set…Name_Act` di Pega.
 */
function PilihKode({
  label,
  daftar,
  simpul,
  kunciID,
  kunciNama,
  bisaUbah,
  onUbah,
}: {
  label: string;
  daftar: readonly PilihanWarisan[];
  simpul: SimpulLimit;
  kunciID: string;
  kunciNama: string;
  bisaUbah: boolean;
  onUbah: (s: SimpulLimit) => void;
}) {
  const kode = teksDari(simpul, kunciID);
  if (!bisaUbah) {
    const nama =
      daftar.find((o) => o.id === kode)?.nama ?? teksDari(simpul, kunciNama);
    return (
      <Field
        label={label}
        value={nama !== "" ? nama : kode}
        readOnly
        onChange={() => undefined}
      />
    );
  }
  return (
    <Pilih
      label={label}
      value={kode}
      kosong={LIMITS_PROP.pilihKosong}
      opsi={opsiDari(daftar, "id")}
      onChange={(v) => {
        onUbah({
          ...simpul,
          [kunciID]: v,
          [kunciNama]: daftar.find((o) => o.id === v)?.nama ?? "",
        });
      }}
    />
  );
}

/** Ganti satu simpul di larik, tanpa mengubah larik aslinya. */
function ganti(
  larik: readonly SimpulLimit[],
  i: number,
  baru: SimpulLimit,
): SimpulLimit[] {
  return larik.map((x, j) => (j === i ? baru : x));
}

/** Satu medan simpul: teks baca-saja (lihat) atau kotak isian (ubah). */
function MedanSimpul({
  m,
  simpul,
  bisaUbah,
  onUbah,
}: {
  m: MedanLimit;
  simpul: SimpulLimit;
  bisaUbah: boolean;
  onUbah: (s: SimpulLimit) => void;
}) {
  const opsi = useContext(OpsiLimitsCtx);
  const ada = m.kunci in simpul;
  const nilai = teksDari(simpul, m.kunci);
  return (
    <div className="trin__limit-medan">
      {/* `.Currency…` pxDropdown — `BrowseCurrencyTreatyIn_RD`, nilai `.Currency`. */}
      {m.mataUang !== undefined &&
        (bisaUbah ? (
          <Pilih
            label={m.label}
            value={teksDari(simpul, m.mataUang)}
            kosong={LIMITS_PROP.pilihKosong}
            opsi={opsiDari(opsi.mataUang, "nama")}
            onChange={(v) => {
              onUbah({ ...simpul, [m.mataUang ?? ""]: v });
            }}
          />
        ) : (
          <Field
            label={m.label}
            value={teksDari(simpul, m.mataUang)}
            readOnly
            onChange={() => undefined}
          />
        ))}
      <Field
        label={m.mataUang !== undefined ? "" : m.label}
        // ⛔ Mode ubah: nilai TERSIMPAN apa adanya di kotak; mode lihat:
        // terjemahan tampil. Terjemahan untuk membaca, bukan untuk mengisi.
        value={bisaUbah ? nilai : formatLimit(m.golongan, m.desimal, nilai)}
        readOnly={!bisaUbah}
        onChange={(v) => {
          onUbah({ ...simpul, [m.kunci]: v });
        }}
      />
      {!ada && !bisaUbah && (
        <span className="trin__redup">{LIMITS_PROP.tidakAda}</span>
      )}
    </div>
  );
}

/**
 * Grid satu larik — kolom dari ekspor; `Add`/`Delete` hidup di mode ubah.
 *
 * ⛔ Kolom yang syarat SEL-nya tidak terpenuhi untuk sebuah baris tetap
 * berdiri — selnya yang kosong. Begitu Pega merender syarat sel.
 */
export function GridLimit({
  kolom,
  baris,
  bisaUbah,
  tambah,
  onUbah,
  judul,
}: {
  kolom: readonly KolomLimit[];
  baris: readonly SimpulLimit[];
  bisaUbah: boolean;
  tambah: boolean;
  onUbah: (baris: SimpulLimit[]) => void;
  judul?: string;
}) {
  const opsi = useContext(OpsiLimitsCtx);
  const total = kolom.reduce((a, k) => a + k.lebar, 0) || 1;
  const pakaiTombol = bisaUbah && tambah;
  return (
    <div className="trin__limit-grid">
      {(judul !== undefined || pakaiTombol) && (
        <div className="trin__pohon-kepala">
          <span>{judul ?? ""}</span>
          {pakaiTombol && (
            <button
              type="button"
              className="btn btn--ghost btn--sm"
              onClick={() => {
                onUbah([...baris, {}]);
              }}
            >
              {LIMITS_PROP.tambah}
            </button>
          )}
        </div>
      )}
      <div className="table-wrap">
        <table className="trin__tabel">
          <colgroup>
            {kolom.map((k, i) => (
              <col
                key={i}
                style={{ width: `${((k.lebar / total) * 100).toFixed(2)}%` }}
              />
            ))}
            {pakaiTombol && <col />}
          </colgroup>
          <thead>
            <tr>
              {kolom.map((k, i) => (
                <th key={i} scope="col">
                  {k.label}
                </th>
              ))}
              {pakaiTombol && <th scope="col" aria-label={LIMITS_PROP.hapus} />}
            </tr>
          </thead>
          <tbody>
            {baris.length === 0 && (
              <tr>
                <td colSpan={kolom.length + (pakaiTombol ? 1 : 0)}>
                  <Kosong pesan={LIMITS_PROP.tanpaBaris} />
                </td>
              </tr>
            )}
            {baris.map((b, r) => (
              <tr key={r}>
                {kolom.map((k, i) => {
                  if (
                    k.kunci === "" ||
                    (k.syarat !== undefined && !k.syarat(b))
                  )
                    return <td key={i} />;
                  const v = teksDari(b, k.kunci);
                  return (
                    <td
                      key={i}
                      className={
                        k.golongan === "teks" ? undefined : "trin__angka"
                      }
                    >
                      {bisaUbah && k.mataUang === "id" ? (
                        <PilihKode
                          label=""
                          daftar={opsi.mataUang}
                          simpul={b}
                          kunciID="CurrencyID"
                          kunciNama={k.kunci}
                          bisaUbah
                          onUbah={(x) => {
                            onUbah(ganti(baris, r, x));
                          }}
                        />
                      ) : bisaUbah && k.mataUang === "nama" ? (
                        <Pilih
                          label=""
                          value={v}
                          kosong={LIMITS_PROP.pilihKosong}
                          opsi={opsiDari(opsi.mataUang, "nama")}
                          onChange={(x) => {
                            onUbah(ganti(baris, r, { ...b, [k.kunci]: x }));
                          }}
                        />
                      ) : bisaUbah ? (
                        <input
                          className="field__input"
                          type="text"
                          value={v}
                          aria-label={k.label || k.kunci}
                          onChange={(e) => {
                            onUbah(
                              ganti(baris, r, {
                                ...b,
                                [k.kunci]: e.target.value,
                              }),
                            );
                          }}
                        />
                      ) : (
                        formatLimit(k.golongan, k.desimal, v)
                      )}
                    </td>
                  );
                })}
                {pakaiTombol && (
                  <td>
                    <button
                      type="button"
                      className="btn btn--ghost btn--sm"
                      onClick={() => {
                        onUbah(baris.filter((_, j) => j !== r));
                      }}
                    >
                      {LIMITS_PROP.hapusBaris}
                    </button>
                  </td>
                )}
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}

/** Tingkat 3 — `Section/DetailLimits.xml`. */
function DetailLimits({
  d,
  bisaUbah,
  onUbah,
}: {
  d: SimpulLimit;
  bisaUbah: boolean;
  onUbah: (d: SimpulLimit) => void;
}) {
  const [tab, setTab] = useState<string>(TAB_DETAIL_LIMIT[0]?.judul ?? "");
  const opsiLimits = useContext(OpsiLimitsCtx);
  const jenis = teksDari(d, "TreatyType");
  const qs = syaratQS(jenis);
  const grid = (larik: string, judul?: string) => {
    const g = GRID_DETAIL[larik];
    if (g === undefined) return null;
    return (
      <GridLimit
        key={larik}
        judul={judul}
        kolom={g.kolom}
        baris={larikDari(d, larik)}
        bisaUbah={bisaUbah}
        tambah={g.tambah}
        onUbah={(b) => {
          onUbah({ ...d, [larik]: b });
        }}
      />
    );
  };
  const medan = (m: MedanLimit) => (
    <MedanSimpul
      key={m.kunci}
      m={m}
      simpul={d}
      bisaUbah={bisaUbah}
      onUbah={onUbah}
    />
  );
  const aktif =
    TAB_DETAIL_LIMIT.find((t) => t.judul === tab) ?? TAB_DETAIL_LIMIT[0];
  return (
    <div className="trin__limit-detail">
      {/* `.TreatyGroupID` pxDropdown — `BrowseTreatyGroup_RD` (TREATYGROUP):
          nilai `ID`, label `.TreatyGroupName`; `SetTreatyGroupName_Act`
          mengisi `.TreatyGroup`. ⚠️ `FetchQSfromMaster` (isi spreading dari
          master treaty) yang juga berjalan di Pega BELUM dibangun. */}
      <div className="trin__limit-medan">
        <PilihKode
          label={LIMITS_PROP.treatyGroup}
          daftar={opsiLimits.kelompokTreaty}
          simpul={d}
          kunciID="TreatyGroupID"
          kunciNama="TreatyGroup"
          bisaUbah={bisaUbah}
          onUbah={onUbah}
        />
      </div>
      {grid("COBList")}

      {qs &&
        medan({
          label: LIMITS_PROP.qs,
          kunci: "QSPct",
          golongan: "persen",
          desimal: 2,
        })}
      {syaratSurplus(jenis) &&
        medan({
          label: LIMITS_PROP.lines,
          kunci: "Surplus",
          golongan: "uang",
          desimal: 2,
        })}

      <span className="trin__limit-label">
        {LIMITS_PROP.limit100}
        {qs && ` ${LIMITS_PROP.seratus} ${LIMITS_PROP.persen}`}
      </span>
      {grid("IOOLimitList")}

      <span className="trin__limit-label">{LIMITS_PROP.retensi}</span>
      {qs &&
        medan({
          label: LIMITS_PROP.persen,
          kunci: "RetentionPct",
          golongan: "persen",
          desimal: null,
        })}
      {grid("RetentionList")}

      <span className="trin__limit-label">{LIMITS_PROP.cession}</span>
      {qs &&
        medan({
          label: LIMITS_PROP.persen,
          kunci: "CessionPct",
          golongan: "persen",
          desimal: null,
        })}
      {grid("CessionList")}

      <StripTab
        tab={TAB_DETAIL_LIMIT.map((t) => t.judul)}
        aktif={aktif?.judul ?? ""}
        onPilih={setTab}
      />
      <Panel judul={aktif?.judul ?? ""}>
        {(aktif?.isi ?? []).map((b, i) =>
          b.t === "medan" ? (
            medan(b.m)
          ) : b.t === "grid" ? (
            grid(b.larik, b.judul)
          ) : (
            <span key={i} className="trin__limit-label">
              {b.teks}
            </span>
          ),
        )}
      </Panel>
    </div>
  );
}

/** Satu baris grid yang dapat dibuka (expandPane). */
function BarisBuka({
  judul,
  nomor,
  bisaUbah,
  onHapus,
  children,
}: {
  judul: string;
  /** Nomor baris — grid Treaty Group menampilkannya, Kind of Treaty tidak. */
  nomor?: number;
  bisaUbah: boolean;
  onHapus: () => void;
  children: ReactNode;
}) {
  // Judul kosong = bilah kosong, seperti Pega (baris baru belum punya nama).
  return (
    <details className="trin__pohon-tingkat">
      <summary>
        {nomor !== undefined && <span className="trin__redup">{nomor} </span>}
        {judul}
        {bisaUbah && (
          <button
            type="button"
            className="btn btn--ghost btn--sm trin__limit-hapus"
            onClick={(e) => {
              e.preventDefault();
              onHapus();
            }}
          >
            {LIMITS_PROP.hapus}
          </button>
        )}
      </summary>
      {children}
    </details>
  );
}

/** Tingkat 2 — `Section/LimitProportional.xml`. */
function LimitProportional({
  l,
  bisaUbah,
  onUbah,
}: {
  l: SimpulLimit;
  bisaUbah: boolean;
  onUbah: (l: SimpulLimit) => void;
}) {
  const detail = larikDari(l, "Detail");
  const { jenisTreaty } = useContext(OpsiLimitsCtx);
  return (
    <div className="trin__limit-detail">
      <div className="trin__limit-medan">
        {/* `.TreatyTypeID` pxDropdown — `BrowseReinsuranceType_RD`
            (REINSURANCETYPE): nilai `.ID`, label `.Note`, `Flag=active`;
            `SetTreatyTypeName_Act` mengisi `.TreatyType` (Kind of Treaty).
            ⚠️ Ekspor menandai sel ini `Read-only` tanpa syarat; dropdown
            dihidupkan di mode Edit atas permintaan pemakai 6 Oktober 2026. */}
        <PilihKode
          label={LIMITS_PROP.treatyType}
          daftar={jenisTreaty}
          simpul={l}
          kunciID="TreatyTypeID"
          kunciNama="TreatyType"
          bisaUbah={bisaUbah}
          onUbah={onUbah}
        />
      </div>
      <div className="trin__pohon-kepala">
        <span>{KOLOM_TREATY_GROUP[0]?.label}</span>
        {bisaUbah && (
          <button
            type="button"
            className="btn btn--ghost btn--sm"
            onClick={() => {
              onUbah({
                ...l,
                Detail: [...detail, { TreatyType: teksDari(l, "TreatyType") }],
              });
            }}
          >
            {LIMITS_PROP.tambah}
          </button>
        )}
      </div>
      {detail.length === 0 && <Kosong pesan={LIMITS_PROP.tanpaBaris} />}
      {detail.map((d, i) => (
        <BarisBuka
          key={i}
          nomor={i + 1}
          judul={teksDari(d, "TreatyGroup")}
          bisaUbah={bisaUbah}
          onHapus={() => {
            onUbah({ ...l, Detail: detail.filter((_, j) => j !== i) });
          }}
        >
          <DetailLimits
            d={d}
            bisaUbah={bisaUbah}
            onUbah={(baru) => {
              onUbah({ ...l, Detail: ganti(detail, i, baru) });
            }}
          />
        </BarisBuka>
      ))}
    </div>
  );
}

/**
 * Tab Limits prop. Pohonnya disalin ke keadaan layar sekali (pemanggil
 * memberi `key` per kontrak); suntingan hidup di sini.
 */
export default function TabLimitsProp({
  pohon,
  mode,
  opsi: opsiAwal,
}: {
  pohon: readonly SimpulLimit[];
  mode: ModeForm;
  /** Isi dropdown tab ini; bila tidak diberikan, diminta sendiri. */
  opsi?: OpsiLimits;
}) {
  const [limits, setLimits] = useState<SimpulLimit[]>(() => [...pohon]);
  const [opsi, setOpsi] = useState<OpsiLimits>(opsiAwal ?? OPSI_KOSONG);
  const bisaUbah = mode === "ubah";
  useEffect(() => {
    if (opsiAwal !== undefined) return;
    let dibuang = false;
    ambilOpsiLimits()
      .then((o) => {
        if (!dibuang) setOpsi(o);
      })
      .catch(() => undefined);
    return () => {
      dibuang = true;
    };
  }, [opsiAwal]);
  return (
    <OpsiLimitsCtx.Provider value={opsi}>
      <Panel judul={LIMITS_PROP.judul}>
        <div
          className="trin__pohon"
          role="group"
          aria-label={LIMITS_PROP.judul}
        >
          <div className="trin__pohon-kepala">
            <span>{KOLOM_KIND_OF_TREATY[0]?.label}</span>
            {bisaUbah && (
              <button
                type="button"
                className="btn btn--ghost btn--sm"
                onClick={() => {
                  setLimits([...limits, { Detail: [] }]);
                }}
              >
                {LIMITS_PROP.tambah}
              </button>
            )}
          </div>
          {limits.length === 0 && <Kosong pesan={LIMITS_PROP.tanpaBaris} />}
          {limits.map((l, i) => (
            <BarisBuka
              key={i}
              judul={teksDari(l, "TreatyType")}
              bisaUbah={bisaUbah}
              onHapus={() => {
                setLimits(limits.filter((_, j) => j !== i));
              }}
            >
              {/* Kind of Treaty TIDAK punya isian sendiri — ia `.TreatyType`,
                diisi `SetTreatyTypeName_Act` dari pilihan Treaty Type. */}
              <LimitProportional
                l={l}
                bisaUbah={bisaUbah}
                onUbah={(baru) => {
                  setLimits(ganti(limits, i, baru));
                }}
              />
            </BarisBuka>
          ))}
        </div>
      </Panel>
    </OpsiLimitsCtx.Provider>
  );
}
