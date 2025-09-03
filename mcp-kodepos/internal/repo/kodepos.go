package repo

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"mcp-kodepos/internal/models"
	"strings"
	"time"
)

type KodeposRepo struct {
	db *sql.DB
}

func NewKodeposRepo(db *sql.DB) *KodeposRepo {
	return &KodeposRepo{db: db}
}

func (r *KodeposRepo) Lookup(ctx context.Context, kodepos string) ([]models.ZipRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := r.db.QueryContext(ctx, `
        SELECT kodepos, kelurahan, kecamatan, kota_kab, provinsi, latitude, longitude
        FROM dbo.kodepos WITH (NOLOCK)
        WHERE kodepos = @p1
        ORDER BY kelurahan, kecamatan
    `, kodepos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanRows(rows)
}

func (r *KodeposRepo) Find(ctx context.Context, p models.ZipFindParams) ([]models.ZipRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	var where []string
	var args []any

	if p.Query != "" {
		like := "%" + strings.TrimSpace(p.Query) + "%"
		where = append(where, "(kelurahan LIKE @p? OR kecamatan LIKE @p? OR kota_kab LIKE @p? OR provinsi LIKE @p? OR kodepos LIKE @p?)")
		for i := 0; i < 5; i++ {
			args = append(args, like)
		}
	}
	if p.Provinsi != "" {
		where = append(where, "provinsi = @p?")
		args = append(args, p.Provinsi)
	}
	if p.KotaKab != "" {
		where = append(where, "kota_kab = @p?")
		args = append(args, p.KotaKab)
	}
	if p.Kecamatan != "" {
		where = append(where, "kecamatan = @p?")
		args = append(args, p.Kecamatan)
	}
	if p.Kelurahan != "" {
		where = append(where, "kelurahan = @p?")
		args = append(args, p.Kelurahan)
	}

	q := `
        SELECT TOP(20) kodepos, kelurahan, kecamatan, kota_kab, provinsi, latitude, longitude
        FROM dbo.kodepos 
    `
	if len(where) > 0 {
		q += " WHERE " + bindify(where)
	}
	q += " ORDER BY kota_kab, kecamatan, kelurahan, kodepos"

	// lim := p.Limit
	// if lim <= 0 || lim > 200 {
	// 	lim = 20
	// }
	// args = append([]any{lim}, args...) // @lim = first

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		log.Printf(q, args...)
		log.Printf(err.Error())
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

func (r *KodeposRepo) Suggest(ctx context.Context, prefix string, limit int) ([]models.ZipRecord, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if limit <= 0 || limit > 50 {
		limit = 10
	}
	like := strings.TrimSpace(prefix) + "%"

	rows, err := r.db.QueryContext(ctx, `
        SELECT TOP(@lim) kodepos, kelurahan, kecamatan, kota_kab, provinsi, latitude, longitude
        FROM dbo.kodepos WITH (NOLOCK)
        WHERE kelurahan LIKE @p1 OR kecamatan LIKE @p1 OR kota_kab LIKE @p1 OR provinsi LIKE @p1 OR kodepos LIKE @p1
        ORDER BY kodepos
    `, limit, like)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanRows(rows)
}

// Helpers

func scanRows(rows *sql.Rows) ([]models.ZipRecord, error) {
	var out []models.ZipRecord
	for rows.Next() {
		var rec models.ZipRecord
		var lat, lng sql.NullFloat64
		if err := rows.Scan(&rec.Kodepos, &rec.Kelurahan, &rec.Kecamatan, &rec.KotaKab, &rec.Provinsi, &lat, &lng); err != nil {
			return nil, err
		}
		if lat.Valid {
			rec.Latitude = &lat.Float64
		}
		if lng.Valid {
			rec.Longitude = &lng.Float64
		}
		out = append(out, rec)
	}
	return out, rows.Err()
}

// Replace each "?" with sequential @p1, @p2 ... for sqlserver
func bindify(conds []string) string {
	idx := 1
	res := conds
	for i, c := range conds {
		for strings.Contains(c, "@p?") {
			c = strings.Replace(c, "@p?", fmt.Sprintf("@p%d", idx), 1)
			idx++
		}
		res[i] = c
	}
	return strings.Join(res, " AND ")
}
